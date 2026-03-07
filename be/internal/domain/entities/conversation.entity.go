package entities

import "github.com/google/uuid"

type ConversationEntity struct {
	ID    uuid.UUID
	Title string
	Type  string // "private" or "group"

	LastMessageID string
}

func NewConversationEntity(id uuid.UUID, title, ctype string) ConversationEntity {
	return ConversationEntity{
		ID:    id,
		Title: title,
		Type:  ctype,
	}
}
