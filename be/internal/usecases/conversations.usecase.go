package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/longtk26/chat-app/internal/domain/entities"
	"github.com/longtk26/chat-app/internal/domain/repo"
	"github.com/longtk26/chat-app/internal/presenters/dto"
)

var _ IConversationsUsecase = &ConversationsUsecase{}

type IConversationsUsecase interface {
	CreateConversation(c context.Context, payload dto.CreateConversationRequestDto) error
}

type ConversationsUsecase struct {
	conveRepo repo.IConversationsRepo
}

func NewConversationsUsecase(conveRepo repo.IConversationsRepo) IConversationsUsecase {
	return &ConversationsUsecase{
		conveRepo: conveRepo,
	}
}

func (conv *ConversationsUsecase) CreateConversation(c context.Context, payload dto.CreateConversationRequestDto) error {
	uid, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	conEntity := entities.NewConversationEntity(
		uid,
		payload.Title,
		payload.Type,
	)

	conv.conveRepo.CreateConversation(
		c,
		conEntity,
		payload.UserIDs,
	)

	return nil
}
