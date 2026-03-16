package entities

import "github.com/google/uuid"

type UserConversationEntity struct {
	UserID   uuid.UUID
	Username string
}

type ConversationEntity struct {
	ID    uuid.UUID
	Title string
	Type  string // "private" or "group"

	Users         []UserConversationEntity
	LastMessageID string
}

func NewConversationEntity(id uuid.UUID, title, ctype string) ConversationEntity {
	return ConversationEntity{
		ID:    id,
		Title: title,
		Type:  ctype,
	}
}
