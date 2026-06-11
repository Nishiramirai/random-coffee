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
	TelegramID      int64
	Username        string
	FullName        string
	About           string
	PreferredFormat MeetingFormat
	State           string
	IsActive        bool
	RegisteredAt    time.Time
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
