package model

import "time"

// MeetingFormat — предпочтительный формат встречи участника.
type MeetingFormat string

const (
	FormatOnline  MeetingFormat = "online"
	FormatOffline MeetingFormat = "offline"
	FormatAny     MeetingFormat = "any"
)

// User — анкета участника сообщества.
type User struct {
	TelegramID int64
	Username   string
	// Name — имя участника (фамилия не требуется для нетворкинга).
	Name  string
	About string
	// City — город участника; запрашивается при выборе офлайн-формата
	// и используется алгоритмом для подбора пар на личную встречу.
	City            string
	PreferredFormat MeetingFormat
	State           string
	IsActive        bool
	RegisteredAt    time.Time
}

// Pair — сформированная пара участников раунда вместе с разрешённым
// форматом встречи. Online означает онлайн-встречу (с ссылкой RoomURL);
// Fallback указывает, что пара собрана несмотря на несовпадение
// предпочтений (подходящего партнёра в нужном формате/городе не нашлось).
type Pair struct {
	A, B     User
	Online   bool
	RoomURL  string
	Fallback bool
}

// Round — раунд матчинга.
type Round struct {
	ID                int
	StartedAt         time.Time
	ParticipantsCount int
}

// Match — пара участников, сформированная в рамках раунда.
type Match struct {
	ID         int
	RoundID    int
	User1ID    int64
	User2ID    int64
	FeedbackU1 string
	FeedbackU2 string
	CreatedAt  time.Time
}
