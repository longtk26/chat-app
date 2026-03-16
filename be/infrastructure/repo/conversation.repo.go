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

func (r *ConversationsRepo) ListConversationsByUserWithPagination(ctx context.Context, userID string, limit, offset int) (repo.ListConversationsResult, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return repo.ListConversationsResult{}, fmt.Errorf("invalid user ID: %w", err)
	}

	if limit <= 0 {
		limit = 20
	}

	rows, err := r.queries.ListConversationsByUserWithPagination(ctx, out.ListConversationsByUserWithPaginationParams{
		UserID: pgtype.UUID{Bytes: uid, Valid: true},
		Limit:  int32(limit),
		Offset: int32(offset),
	})

	if err != nil {
		return repo.ListConversationsResult{}, fmt.Errorf("failed to list conversations: %w", err)
	}

	conversations := make([]entities.ConversationEntity, 0, len(rows))
	for _, row := range rows {
		conID, err := uuid.FromBytes(row.ID.Bytes[:])
		if err != nil {
			return repo.ListConversationsResult{}, fmt.Errorf("failed to parse conversation ID: %w", err)
		}

		conversationUsersRows, err := r.queries.ListConversationUsers(ctx, row.ID)
		if err != nil {
			return repo.ListConversationsResult{}, fmt.Errorf("failed to list conversation users: %w", err)
		}

		conversationUsers := make([]entities.UserConversationEntity, 0, len(conversationUsersRows))
		for _, userRow := range conversationUsersRows {
			if !userRow.ID.Valid {
				continue
			}

			userID, err := uuid.FromBytes(userRow.ID.Bytes[:])
			if err != nil {
				return repo.ListConversationsResult{}, fmt.Errorf("failed to parse user ID: %w", err)
			}

			conversationUsers = append(conversationUsers, entities.UserConversationEntity{
				UserID:   userID,
				Username: userRow.Username,
			})
		}

		lastMessageID := ""
		if row.LastMessageID.Valid {
			lmID, err := uuid.FromBytes(row.LastMessageID.Bytes[:])
			if err != nil {
				return repo.ListConversationsResult{}, fmt.Errorf("failed to parse last message ID: %w", err)
			}
			lastMessageID = lmID.String()
		}

		title := ""
		if row.Title.Valid {
			title = row.Title.String
		}

		conversations = append(conversations, entities.ConversationEntity{
			ID:            conID,
			Title:         title,
			Type:          row.Type,
			Users:         conversationUsers,
			LastMessageID: lastMessageID,
		})
	}

	totalItems := 0
	if len(rows) > 0 {
		totalItems = int(rows[0].TotalConversations)
	}

	totalPages := 0
	if totalItems > 0 {
		totalPages = (totalItems + limit - 1) / limit
	}

	return repo.ListConversationsResult{
		Conversations: conversations,
		TotalItems:    totalItems,
		TotalPages:    totalPages,
	}, nil
}
