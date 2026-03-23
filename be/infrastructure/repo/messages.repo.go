package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/longtk26/chat-app/infrastructure/db/sqlc/out"
	"github.com/longtk26/chat-app/internal/domain/entities"
	"github.com/longtk26/chat-app/internal/domain/repo"
)

var _ repo.IMessagesRepo = &MessagesRepo{}

type MessagesRepo struct {
	queries *out.Queries
}

func NewMessagesRepo(queries *out.Queries) repo.IMessagesRepo {
	return &MessagesRepo{queries: queries}
}

func (r *MessagesRepo) ListMessagesByConversation(ctx context.Context, conversationID, lastMessageID string, cursorTime *time.Time, limit int) (repo.ListMessagesCursorResult, error) {
	conversationUUID, err := uuid.Parse(conversationID)
	if err != nil {
		return repo.ListMessagesCursorResult{}, fmt.Errorf("invalid conversation ID: %w", err)
	}

	if limit <= 0 {
		limit = 20
	}

	lastMessageUUID, err := parseOptionalUUID(lastMessageID)
	if err != nil {
		return repo.ListMessagesCursorResult{}, fmt.Errorf("invalid last_message_id: %w", err)
	}

	cursor := pgtype.Timestamptz{}
	if cursorTime != nil {
		cursor = pgtype.Timestamptz{Time: cursorTime.UTC(), Valid: true}
	}

	rows, err := r.queries.ListMessagesByConversationCursor(ctx, out.ListMessagesByConversationCursorParams{
		ConversationID: pgtype.UUID{Bytes: conversationUUID, Valid: true},
		Column2:        lastMessageUUID,
		Column3:        cursor,
		Limit:          int32(limit + 1),
	})
	if err != nil {
		return repo.ListMessagesCursorResult{}, fmt.Errorf("failed to list messages: %w", err)
	}

	hasMore := len(rows) > limit
	if hasMore {
		rows = rows[:limit]
	}

	messages := make([]entities.MessageEntity, 0, len(rows))
	for _, row := range rows {
		message, err := toMessageEntity(row)
		if err != nil {
			return repo.ListMessagesCursorResult{}, err
		}
		messages = append(messages, message)
	}

	result := repo.ListMessagesCursorResult{
		Messages: messages,
		HasMore:  hasMore,
	}

	if len(messages) > 0 {
		last := messages[len(messages)-1].SentAt
		result.NextCursorTime = &last
	}

	return result, nil
}

func (r *MessagesRepo) CreateMessage(ctx context.Context, payload entities.MessageEntity) (entities.MessageEntity, error) {
	senderUUID, err := uuid.Parse(payload.SenderID)
	if err != nil {
		return entities.MessageEntity{}, fmt.Errorf("invalid sender ID: %w", err)
	}

	// recipient_id is still required by the current schema, so we keep it server-managed.
	recipientUUID := senderUUID

	conversationUUID, err := uuid.Parse(payload.ConversationID)
	if err != nil {
		return entities.MessageEntity{}, fmt.Errorf("invalid conversation ID: %w", err)
	}

	conversationPGUUID := pgtype.UUID{Bytes: conversationUUID, Valid: true}

	message, err := r.queries.CreateMessage(ctx, out.CreateMessageParams{
		SenderID:       pgtype.UUID{Bytes: senderUUID, Valid: true},
		RecipientID:    pgtype.UUID{Bytes: recipientUUID, Valid: true},
		ConversationID: conversationPGUUID,
		Content:        payload.Content,
	})

	if err != nil {
		return entities.MessageEntity{}, fmt.Errorf("failed to create message: %w", err)
	}

	if err := r.queries.UpdateConversationLastMessage(ctx, out.UpdateConversationLastMessageParams{
		ID:            conversationPGUUID,
		LastMessageID: message.ID,
	}); err != nil {
		return entities.MessageEntity{}, fmt.Errorf("failed to update conversation last message: %w", err)
	}

	return toMessageEntity(message)
}

func (r *MessagesRepo) UpdateMessageContent(ctx context.Context, messageID, content string) (entities.MessageEntity, error) {
	messageUUID, err := uuid.Parse(messageID)
	if err != nil {
		return entities.MessageEntity{}, fmt.Errorf("invalid message ID: %w", err)
	}

	message, err := r.queries.UpdateMessageContent(ctx, out.UpdateMessageContentParams{
		ID:      pgtype.UUID{Bytes: messageUUID, Valid: true},
		Content: content,
	})
	if err != nil {
		return entities.MessageEntity{}, fmt.Errorf("failed to update message: %w", err)
	}

	return toMessageEntity(message)
}

func (r *MessagesRepo) DeleteMessage(ctx context.Context, messageID string) error {
	messageUUID, err := uuid.Parse(messageID)
	if err != nil {
		return fmt.Errorf("invalid message ID: %w", err)
	}

	if err := r.queries.SoftDeleteMessage(ctx, pgtype.UUID{Bytes: messageUUID, Valid: true}); err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	return nil
}

func parseOptionalUUID(raw string) (pgtype.UUID, error) {
	if raw == "" {
		return pgtype.UUID{}, nil
	}

	parsedUUID, err := uuid.Parse(raw)
	if err != nil {
		return pgtype.UUID{}, err
	}

	return pgtype.UUID{Bytes: parsedUUID, Valid: true}, nil
}

func toMessageEntity(message out.Message) (entities.MessageEntity, error) {
	messageID, err := uuidFromPGUUID(message.ID)
	if err != nil {
		return entities.MessageEntity{}, fmt.Errorf("failed to parse message ID: %w", err)
	}

	senderID, err := uuidFromPGUUID(message.SenderID)
	if err != nil {
		return entities.MessageEntity{}, fmt.Errorf("failed to parse sender ID: %w", err)
	}

	conversationID, err := uuidFromPGUUID(message.ConversationID)
	if err != nil {
		return entities.MessageEntity{}, fmt.Errorf("failed to parse conversation ID: %w", err)
	}

	return entities.MessageEntity{
		ID:             messageID,
		SenderID:       senderID,
		ConversationID: conversationID,
		Content:        message.Content,
		SentAt:         message.SentAt.Time.UTC(),
		UpdatedAt:      message.UpdatedAt.Time.UTC(),
	}, nil
}

func uuidFromPGUUID(value pgtype.UUID) (string, error) {
	if !value.Valid {
		return "", fmt.Errorf("uuid is not valid")
	}

	parsed, err := uuid.FromBytes(value.Bytes[:])
	if err != nil {
		return "", err
	}

	return parsed.String(), nil
}
