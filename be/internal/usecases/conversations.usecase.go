package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/longtk26/chat-app/internal/domain/entities"
	"github.com/longtk26/chat-app/internal/domain/repo"
	"github.com/longtk26/chat-app/internal/presenters/dto"
	baseDto "github.com/longtk26/chat-app/pkg/dto"
)

var _ IConversationsUsecase = &ConversationsUsecase{}

type IConversationsUsecase interface {
	CreateConversation(c context.Context, payload dto.CreateConversationRequestDto) error
	ListConversations(c context.Context, query dto.ListConversationsQueryDto) (dto.ListConversationsResponseDto, error)
	GetConversationByID(c context.Context, conversationID string) (dto.ConversationDto, error)
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

func (conv *ConversationsUsecase) ListConversations(c context.Context, query dto.ListConversationsQueryDto) (dto.ListConversationsResponseDto, error) {
	page := query.Page
	if page <= 0 {
		page = 1
	}

	pageSize := query.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	result, err := conv.conveRepo.ListConversationsByUserWithPagination(
		c,
		query.UserID,
		pageSize,
		offset,
	)
	if err != nil {
		return dto.ListConversationsResponseDto{}, err
	}

	conversations := make([]dto.ConversationDto, len(result.Conversations))
	for i, con := range result.Conversations {
		users := make([]dto.UserDto, len(con.Users))
		for j, user := range con.Users {
			users[j] = dto.UserDto{
				ID:       user.UserID.String(),
				Username: user.Username,
			}
		}

		conversations[i] = dto.ConversationDto{
			ID:    con.ID.String(),
			Title: con.Title,
			Type:  con.Type,
			Users: users,
		}
	}

	response := dto.ListConversationsResponseDto{
		Conversations: conversations,
		PaginationResponseDto: baseDto.PaginationResponseDto{
			TotalCount: result.TotalItems,
			Page:       page,
		},
	}

	return response, nil
}

func (conv *ConversationsUsecase) GetConversationByID(c context.Context, conversationID string) (dto.ConversationDto, error) {
	conve, err := conv.conveRepo.GetConversationByID(c, conversationID)
	if err != nil {
		return dto.ConversationDto{}, err
	}

	response := dto.ConversationDto{
		ID:    conve.ID.String(),
		Title: conve.Title,
		Type:  conve.Type,
	}

	return response, nil
}
