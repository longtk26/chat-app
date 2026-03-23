package entities

import "time"

type MessageEntity struct {
	ID             string
	SenderID       string
	ConversationID string
	Content        string
	SentAt         time.Time
	UpdatedAt      time.Time
}
