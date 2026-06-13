package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config — конфигурация приложения, заполняемая из переменных окружения.
type Config struct {
	BotToken        string `env:"BOT_TOKEN" env-required:"true"`
	AdminTelegramID int64  `env:"ADMIN_TELEGRAM_ID" env-required:"true"`
	// MatchSchedule задаёт расписание автоматического раунда в формате "DOW HH:MM",
	// например "MON 10:00". Если не задан, автоматический запуск отключён.
	MatchSchedule string `env:"MATCH_SCHEDULE"`
	// VideoCallBase — адрес публичного сервера видеовстреч (Jitsi), например
	// "https://meet.ffmuc.net/". Если не задан, используется значение по умолчанию.
	VideoCallBase string `env:"VIDEO_CALL_BASE"`
	DB            DBConfig
}

type DBConfig struct {
	Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
	Port     int    `env:"POSTGRES_PORT" env-default:"5432"`
	Name     string `env:"POSTGRES_DB" env-required:"true"`
	User     string `env:"POSTGRES_USER" env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
}

// DSN формирует строку подключения к PostgreSQL.
func (d DBConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		d.User, d.Password, d.Host, d.Port, d.Name)
}

// Load читает конфигурацию из окружения и валидирует обязательные поля.
func Load() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	return &cfg, nil
}
