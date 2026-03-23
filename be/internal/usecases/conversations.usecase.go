package usecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/longtk26/chat-app/internal/domain/entities"
	"github.com/longtk26/chat-app/internal/domain/repo"
	"github.com/longtk26/chat-app/internal/presenters/dto"
	baseDto "github.com/longtk26/chat-app/pkg/dto"
)

var _ IConversationsUsecase = &ConversationsUsecase{}

type IConversationsUsecase interface {
	CreateConversation(c context.Context, payload dto.CreateConversationRequestDto) (dto.CreateConversationResponseDto, error)
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

func toConversationDto(conversation entities.ConversationEntity, fallbackUserIDs []string) dto.ConversationDto {
	users := make([]dto.UserDto, 0, len(conversation.Users))
	for _, user := range conversation.Users {
		users = append(users, dto.UserDto{
			ID:       user.UserID.String(),
			Username: user.Username,
		})
	}

	if len(users) == 0 {
		users = make([]dto.UserDto, 0, len(fallbackUserIDs))
		for _, userID := range fallbackUserIDs {
			users = append(users, dto.UserDto{ID: userID})
		}
	}

	return dto.ConversationDto{
		ID:    conversation.ID.String(),
		Title: conversation.Title,
		Type:  conversation.Type,
		Users: users,
	}
}

func (conv *ConversationsUsecase) CreateConversation(
	c context.Context,
	payload dto.CreateConversationRequestDto,
) (dto.CreateConversationResponseDto, error) {
	if len(payload.UserIDs) != 2 {
		return dto.CreateConversationResponseDto{}, fmt.Errorf("private conversation requires exactly 2 users")
	}

	userA := payload.UserIDs[0]
	userB := payload.UserIDs[1]

	existingConve, err := conv.conveRepo.GetPrivateConversationBetweenUsers(c, userA, userB)

	// Log debug info about existing conversation check
	if err != nil {
		fmt.Printf("Error checking for existing private conversation between users %s and %s: %v\n", userA, userB, err)
	} else if existingConve.ID != uuid.Nil {
		fmt.Printf("Existing private conversation found between users %s and %s: Conversation ID %s\n", userA, userB, existingConve.ID.String())
	}

	if err == nil && existingConve.ID != uuid.Nil {
		return dto.CreateConversationResponseDto{
			Conversation: toConversationDto(existingConve, payload.UserIDs),
			Message:      "Private conversation already exists between these users",
		}, nil
	}

	// 2. Create new conversation
	uid, err := uuid.NewRandom()
	if err != nil {
		return dto.CreateConversationResponseDto{}, err
	}

	conEntity := entities.NewConversationEntity(
		uid,
		payload.Title,
		payload.Type,
	)

	err = conv.conveRepo.CreateConversation(
		c,
		conEntity,
		payload.UserIDs,
	)
	if err != nil {
		existingConve, getErr := conv.conveRepo.GetPrivateConversationBetweenUsers(c, userA, userB)
		if getErr == nil && existingConve.ID != uuid.Nil {
			return dto.CreateConversationResponseDto{
				Conversation: toConversationDto(existingConve, payload.UserIDs),
				Message:      "Private conversation already exists between these users",
			}, nil
		}
		return dto.CreateConversationResponseDto{}, err
	}

	return dto.CreateConversationResponseDto{
		Conversation: toConversationDto(conEntity, payload.UserIDs),
		Message:      "Conversation created successfully",
	}, nil
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

	users := make([]dto.UserDto, len(conve.Users))
	for i, user := range conve.Users {
		users[i] = dto.UserDto{
			ID:       user.UserID.String(),
			Username: user.Username,
		}
	}
	fmt.Printf("last message ID %s", conve.LastMessageID)
	response := dto.ConversationDto{
		ID:            conve.ID.String(),
		Title:         conve.Title,
		Type:          conve.Type,
		Users:         users,
		LastMessageID: conve.LastMessageID,
	}

	return response, nil
}
