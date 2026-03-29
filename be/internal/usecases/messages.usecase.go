package usecases

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/longtk26/chat-app/internal/domain/entities"
	"github.com/longtk26/chat-app/internal/domain/repo"
	"github.com/longtk26/chat-app/internal/presenters/dto"
	baseDto "github.com/longtk26/chat-app/pkg/dto"
)

var _ IMessagesUseCase = &MessagesUseCase{}

type IMessagesUseCase interface {
	ListMessages(ctx context.Context, query dto.ListMessagesQueryDto) (dto.ListMessagesResponseDto, error)
	SendMessage(ctx context.Context, payload dto.SendMessageRequestDto) (dto.SendMessageResponseDto, error)
	UpdateMessage(ctx context.Context, messageID string, payload dto.UpdateMessageRequestDto) (dto.UpdateMessageResponseDto, error)
	DeleteMessage(ctx context.Context, messageID string) error
}

type MessagesUseCase struct {
	messageRepo repo.IMessagesRepo
}

func NewMessagesUseCase(messageRepo repo.IMessagesRepo) IMessagesUseCase {
	return &MessagesUseCase{messageRepo: messageRepo}
}

func (u *MessagesUseCase) ListMessages(ctx context.Context, query dto.ListMessagesQueryDto) (dto.ListMessagesResponseDto, error) {
	if query.ConversationID == "" {
		return dto.ListMessagesResponseDto{}, fmt.Errorf("conversation_id is required")
	}

	limit := query.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	var cursorTime *time.Time
	if query.CursorTime != "" {
		parsed, err := parseCursorTime(query.CursorTime)
		if err != nil {
			return dto.ListMessagesResponseDto{}, fmt.Errorf("invalid cursor_time: %w", err)
		}
		cursorTime = &parsed
	}

	result, err := u.messageRepo.ListMessagesByConversation(
		ctx,
		query.ConversationID,
		query.LastMessageID,
		cursorTime,
		limit,
	)
	if err != nil {
		return dto.ListMessagesResponseDto{}, err
	}

	messages := make([]dto.MessageDto, 0, len(result.Messages))
	for _, message := range result.Messages {
		messages = append(messages, toMessageDto(message))
	}

	response := dto.ListMessagesResponseDto{
		Messages: messages,
		CursorPaginationResponseDto: baseDto.CursorPaginationResponseDto{
			Limit:   limit,
			HasMore: result.HasMore,
		},
	}

	if result.NextCursorTime != nil {
		response.NextCursorTime = result.NextCursorTime.UTC().Format(time.RFC3339Nano)
	}

	return response, nil
}

func (u *MessagesUseCase) SendMessage(ctx context.Context, payload dto.SendMessageRequestDto) (dto.SendMessageResponseDto, error) {
	payload.Content = strings.TrimSpace(payload.Content)
	if payload.SenderID == "" || payload.ConversationID == "" || payload.Content == "" {
		return dto.SendMessageResponseDto{}, fmt.Errorf("sender_id, conversation_id and content are required")
	}

	createdMessage, err := u.messageRepo.CreateMessage(ctx, entities.MessageEntity{
		SenderID:       payload.SenderID,
		ConversationID: payload.ConversationID,
		Content:        payload.Content,
	})

	if err != nil {
		return dto.SendMessageResponseDto{}, err
	}

	entityMessage := entities.MessageEntity{
		ID:             createdMessage.ID,
		SenderID:       createdMessage.SenderID,
		ConversationID: createdMessage.ConversationID,
		Content:        createdMessage.Content,
		SentAt:         createdMessage.SentAt,
		UpdatedAt:      createdMessage.UpdatedAt,
		SenderName:     payload.SenderName,
	}

	return dto.SendMessageResponseDto{Message: toMessageDto(entityMessage)}, nil
}

func (u *MessagesUseCase) UpdateMessage(ctx context.Context, messageID string, payload dto.UpdateMessageRequestDto) (dto.UpdateMessageResponseDto, error) {
	payload.Content = strings.TrimSpace(payload.Content)
	if messageID == "" {
		return dto.UpdateMessageResponseDto{}, fmt.Errorf("message id is required")
	}
	if payload.Content == "" {
		return dto.UpdateMessageResponseDto{}, fmt.Errorf("content is required")
	}

	updatedMessage, err := u.messageRepo.UpdateMessageContent(ctx, messageID, payload.Content)
	if err != nil {
		return dto.UpdateMessageResponseDto{}, err
	}

	return dto.UpdateMessageResponseDto{Message: toMessageDto(updatedMessage)}, nil
}

func (u *MessagesUseCase) DeleteMessage(ctx context.Context, messageID string) error {
	if messageID == "" {
		return fmt.Errorf("message id is required")
	}

	return u.messageRepo.DeleteMessage(ctx, messageID)
}

func parseCursorTime(raw string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339Nano, raw); err == nil {
		return t.UTC(), nil
	}

	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return time.Time{}, err
	}

	return t.UTC(), nil
}

func toMessageDto(message entities.MessageEntity) dto.MessageDto {
	return dto.MessageDto{
		ID:             message.ID,
		SenderID:       message.SenderID,
		ConversationID: message.ConversationID,
		Content:        message.Content,
		SenderName:     message.SenderName,
		SentAt:         message.SentAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:      message.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
}
