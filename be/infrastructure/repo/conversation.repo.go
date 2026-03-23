package repositories

import (
	"context"
	"encoding/json"
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
		conID, err := uuid.FromBytes(row.Conversation.ID.Bytes[:])
		if err != nil {
			return repo.ListConversationsResult{}, fmt.Errorf("failed to parse conversation ID: %w", err)
		}

		var participants []struct {
			UserID   uuid.UUID `json:"id"`
			Username string    `json:"username"`
		}

		if err := json.Unmarshal(row.Participants, &participants); err != nil {
			return repo.ListConversationsResult{}, fmt.Errorf("failed to unmarshal participants: %w", err)
		}

		conversationUsers := make([]entities.UserConversationEntity, 0, len(participants))
		for _, userRow := range participants {
			conversationUsers = append(conversationUsers, entities.UserConversationEntity{
				UserID:   userRow.UserID,
				Username: userRow.Username,
			})
		}

		lastMessageID := ""
		if row.Conversation.LastMessageID.Valid {
			lmID, err := uuid.FromBytes(row.Conversation.LastMessageID.Bytes[:])
			if err != nil {
				return repo.ListConversationsResult{}, fmt.Errorf("failed to parse last message ID: %w", err)
			}
			lastMessageID = lmID.String()
		}

		conversations = append(conversations, entities.ConversationEntity{
			ID:            conID,
			Title:         row.Conversation.Title.String,
			Type:          row.Conversation.Type,
			Users:         conversationUsers,
			LastMessageID: lastMessageID,
		})
	}

	totalItems := 0
	if len(rows) > 0 {
		totalItems = int(rows[0].TotalCount)
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

func (r *ConversationsRepo) GetConversationByID(ctx context.Context, conversationID string) (entities.ConversationEntity, error) {
	// convert to pgtype.UUID
	convID, err := uuid.Parse(conversationID)

	if err != nil {
		return entities.ConversationEntity{}, fmt.Errorf("invalid conversation ID: %w", err)
	}
	dbID := pgtype.UUID{
		Bytes: [16]byte(convID),
		Valid: true,
	}

	row, err := r.queries.GetConversationByID(ctx, dbID)
	if err != nil {
		return entities.ConversationEntity{}, fmt.Errorf("failed to get conversation: %w", err)
	}

	lastMessageID := ""
	if row.LastMessageID.Valid {
		lmID, err := uuid.FromBytes(row.LastMessageID.Bytes[:])
		if err != nil {
			return entities.ConversationEntity{}, fmt.Errorf("failed to parse last message ID: %w", err)
		}
		lastMessageID = lmID.String()
	}

	users := make([]entities.UserConversationEntity, 0, len(row.Participants))
	var participants []struct {
		UserID   uuid.UUID `json:"id"`
		Username string    `json:"username"`
	}

	if err := json.Unmarshal(row.Participants, &participants); err != nil {
		return entities.ConversationEntity{}, fmt.Errorf("failed to unmarshal participants: %w", err)
	}

	users = make([]entities.UserConversationEntity, 0, len(participants))
	for _, participant := range participants {
		users = append(users, entities.UserConversationEntity{
			UserID:   participant.UserID,
			Username: participant.Username,
		})
	}

	return entities.ConversationEntity{
		ID:            convID,
		Title:         row.Title.String,
		Type:          row.Type,
		LastMessageID: lastMessageID,
		Users:         users,
	}, nil
}

func (r *ConversationsRepo) GetPrivateConversationBetweenUsers(ctx context.Context, userID1, userID2 string) (entities.ConversationEntity, error) {
	uid1, err := uuid.Parse(userID1)
	if err != nil {
		return entities.ConversationEntity{}, fmt.Errorf("invalid user ID 1: %w", err)
	}

	uid2, err := uuid.Parse(userID2)
	if err != nil {
		return entities.ConversationEntity{}, fmt.Errorf("invalid user ID 2: %w", err)
	}

	row, err := r.queries.GetPrivateConversationBetweenUsers(ctx, out.GetPrivateConversationBetweenUsersParams{
		UserID:   pgtype.UUID{Bytes: uid1, Valid: true},
		UserID_2: pgtype.UUID{Bytes: uid2, Valid: true},
	})
	if err != nil {
		return entities.ConversationEntity{}, fmt.Errorf("failed to get private conversation: %w", err)
	}

	lastMessageID := ""
	if row.LastMessageID.Valid {
		lmID, err := uuid.FromBytes(row.LastMessageID.Bytes[:])
		if err != nil {
			return entities.ConversationEntity{}, fmt.Errorf("failed to parse last message ID: %w", err)
		}
		lastMessageID = lmID.String()
	}

	users := make([]entities.UserConversationEntity, 0, len(row.Participants))
	var participants []struct {
		UserID   uuid.UUID `json:"id"`
		Username string    `json:"username"`
	}

	if err := json.Unmarshal(row.Participants, &participants); err != nil {
		return entities.ConversationEntity{}, fmt.Errorf("failed to unmarshal participants: %w", err)
	}

	users = make([]entities.UserConversationEntity, 0, len(participants))
	for _, participant := range participants {
		users = append(users, entities.UserConversationEntity{
			UserID:   participant.UserID,
			Username: participant.Username,
		})
	}

	return entities.ConversationEntity{
		ID:            uuid.Must(uuid.FromBytes(row.ID.Bytes[:])),
		Title:         row.Title.String,
		Type:          row.Type,
		LastMessageID: lastMessageID,
		Users:         users,
	}, nil
}
