package db

import (
	"context"
	"errors"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Nishiramirai/random-coffee/internal/model"
)

// Storage — слой доступа к данным поверх пула соединений pgx.
type Storage struct {
	pool *pgxpool.Pool
}

// New создаёт пул соединений и проверяет доступность базы данных.
func New(ctx context.Context, dsn string) (*Storage, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	return &Storage{pool: pool}, nil
}

func (s *Storage) Close() { s.pool.Close() }

// RunMigrations применяет SQL-миграции из каталога migrations.
func RunMigrations(dsn string) error {
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

// GetUser возвращает участника по Telegram ID или nil, если он не найден.
func (s *Storage) GetUser(ctx context.Context, id int64) (*model.User, error) {
	const q = `SELECT telegram_id, username, full_name, about,
	                  preferred_format, state, is_active, registered_at
	           FROM users WHERE telegram_id = $1`
	var u model.User
	err := s.pool.QueryRow(ctx, q, id).Scan(&u.TelegramID, &u.Username,
		&u.FullName, &u.About, &u.PreferredFormat, &u.State,
		&u.IsActive, &u.RegisteredAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// UpsertUser создаёт или обновляет анкету участника.
func (s *Storage) UpsertUser(ctx context.Context, u *model.User) error {
	const q = `INSERT INTO users (telegram_id, username, full_name, about,
	               preferred_format, state, is_active, registered_at)
	           VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	           ON CONFLICT (telegram_id) DO UPDATE SET
	               username = $2, full_name = $3, about = $4,
	               preferred_format = $5, state = $6, is_active = $7`
	_, err := s.pool.Exec(ctx, q, u.TelegramID, u.Username, u.FullName,
		u.About, u.PreferredFormat, u.State, u.IsActive, u.RegisteredAt)
	return err
}

func (s *Storage) SetState(ctx context.Context, id int64, state string) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE users SET state = $2 WHERE telegram_id = $1`, id, state)
	return err
}

func (s *Storage) SetActive(ctx context.Context, id int64, active bool) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE users SET is_active = $2 WHERE telegram_id = $1`, id, active)
	return err
}

// GetActiveUsers возвращает всех участников, не находящихся на паузе.
func (s *Storage) GetActiveUsers(ctx context.Context) ([]model.User, error) {
	const q = `SELECT telegram_id, username, full_name, about,
	                  preferred_format, state, is_active, registered_at
	           FROM users WHERE is_active = TRUE`
	rows, err := s.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.TelegramID, &u.Username, &u.FullName,
			&u.About, &u.PreferredFormat, &u.State, &u.IsActive,
			&u.RegisteredAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

// CreateRoundWithMatches сохраняет раунд и его пары в одной транзакции.
func (s *Storage) CreateRoundWithMatches(
	ctx context.Context, pairs [][2]model.User) (int, error) {

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var roundID int
	err = tx.QueryRow(ctx,
		`INSERT INTO rounds (started_at, participants_count)
		 VALUES ($1, $2) RETURNING id`,
		time.Now(), len(pairs)*2).Scan(&roundID)
	if err != nil {
		return 0, err
	}

	for _, p := range pairs {
		_, err = tx.Exec(ctx,
			`INSERT INTO matches (round_id, user1_id, user2_id, created_at)
			 VALUES ($1, $2, $3, $4)`,
			roundID, p[0].TelegramID, p[1].TelegramID, time.Now())
		if err != nil {
			return 0, err
		}
	}
	return roundID, tx.Commit(ctx)
}

// GetMatchHistory возвращает множество пар из последних n раундов.
func (s *Storage) GetMatchHistory(
	ctx context.Context, n int) (map[[2]int64]bool, error) {

	const q = `SELECT user1_id, user2_id FROM matches
	           WHERE round_id > (SELECT COALESCE(MAX(id), 0) - $1 FROM rounds)`
	rows, err := s.pool.Query(ctx, q, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hist := make(map[[2]int64]bool)
	for rows.Next() {
		var a, b int64
		if err := rows.Scan(&a, &b); err != nil {
			return nil, err
		}
		if a > b {
			a, b = b, a
		}
		hist[[2]int64{a, b}] = true
	}
	return hist, rows.Err()
}

// Stats — агрегированная статистика сообщества.
type Stats struct {
	Total  int
	Active int
	Rounds int
}

func (s *Storage) Stats(ctx context.Context) (Stats, error) {
	var st Stats
	err := s.pool.QueryRow(ctx, `SELECT
	        (SELECT COUNT(*) FROM users),
	        (SELECT COUNT(*) FROM users WHERE is_active),
	        (SELECT COUNT(*) FROM rounds)`).Scan(
		&st.Total, &st.Active, &st.Rounds)
	return st, err
}

// SaveFeedback сохраняет результат обратной связи в последней встрече участника.
func (s *Storage) SaveFeedback(ctx context.Context,
	userID int64, code string) error {
	const q = `UPDATE matches SET
	    feedback_u1 = CASE WHEN user1_id = $1 THEN $2 ELSE feedback_u1 END,
	    feedback_u2 = CASE WHEN user2_id = $1 THEN $2 ELSE feedback_u2 END
	    WHERE id = (SELECT id FROM matches
	                WHERE user1_id = $1 OR user2_id = $1
	                ORDER BY created_at DESC LIMIT 1)`
	_, err := s.pool.Exec(ctx, q, userID, code)
	return err
}
