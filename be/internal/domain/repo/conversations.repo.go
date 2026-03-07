package repo

import (
	"context"

	"github.com/longtk26/chat-app/internal/domain/entities"
)

type IConversationsRepo interface {
	CreateConversation(ctx context.Context, conversation entities.ConversationEntity, userIds []string) error
}
