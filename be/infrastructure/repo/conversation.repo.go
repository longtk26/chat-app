package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/longtk26/chat-app/infrastructure/db/sqlc/out"
	"github.com/longtk26/chat-app/internal/domain/entities"
	"github.com/longtk26/chat-app/internal/domain/repo"
)

var _ repo.IConversationsRepo = &ConversationsRepo{}

type ConversationsRepo struct {
	queries *out.Queries
}

func NewConversationsRepo(queries *out.Queries) repo.IConversationsRepo {
	return &ConversationsRepo{
		queries: queries,
	}
}

func (r *ConversationsRepo) CreateConversation(ctx context.Context, conversation entities.ConversationEntity, userIds []string) error {
	var err error
	createdCon, err := r.queries.CreateConversation(ctx, out.CreateConversationParams{
		Title: pgtype.Text{String: conversation.Title, Valid: true},
		Type:  conversation.Type,
	})

	if err != nil {
		return fmt.Errorf("failed to create conversation: %w", err)
	}

	userIDs := make([]pgtype.UUID, len(userIds))
	for i, userId := range userIds {
		uid, err := uuid.Parse(userId)
		if err != nil {
			return fmt.Errorf("invalid user ID: %w", err)
		}
		userIDs[i] = pgtype.UUID{Bytes: uid, Valid: true}
	}

	_, err = r.queries.AddConversationParticipants(ctx, out.AddConversationParticipantsParams{
		ConversationID: createdCon.ID,
		Column2:        userIDs,
	})
	if err != nil {
		return fmt.Errorf("failed to add participants: %w", err)
	}

	return nil
}
