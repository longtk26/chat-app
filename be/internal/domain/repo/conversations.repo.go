package repo

import (
	"context"

	"github.com/longtk26/chat-app/internal/domain/entities"
)

type ListConversationsResult struct {
	Conversations []entities.ConversationEntity
	TotalItems    int
	TotalPages    int
}

type IConversationsRepo interface {
	CreateConversation(ctx context.Context, conversation entities.ConversationEntity, userIds []string) error
	ListConversationsByUserWithPagination(ctx context.Context, userId string, limit, offset int) (
		ListConversationsResult,
		error,
	)
	GetConversationByID(ctx context.Context, conversationID string) (entities.ConversationEntity, error)
	GetPrivateConversationBetweenUsers(ctx context.Context, userID1, userID2 string) (entities.ConversationEntity, error)
}
