package repo

import (
	"context"
	"time"

	"github.com/longtk26/chat-app/internal/domain/entities"
)

type ListMessagesCursorResult struct {
	Messages       []entities.MessageEntity
	NextCursorTime *time.Time
	HasMore        bool
}

type IMessagesRepo interface {
	ListMessagesByConversation(ctx context.Context, conversationID, lastMessageID string, cursorTime *time.Time, limit int) (ListMessagesCursorResult, error)
	CreateMessage(ctx context.Context, payload entities.MessageEntity) (entities.MessageEntity, error)
	UpdateMessageContent(ctx context.Context, messageID, content string) (entities.MessageEntity, error)
	DeleteMessage(ctx context.Context, messageID string) error
}
