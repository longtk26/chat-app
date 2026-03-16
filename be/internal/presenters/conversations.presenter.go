package presenters

import (
	"github.com/gofiber/fiber/v3"
	"github.com/longtk26/chat-app/internal/presenters/dto"
	"github.com/longtk26/chat-app/internal/usecases"
)

type ConversationsPresenter struct {
	usecase usecases.IConversationsUsecase
}

func NewConversationsPresenter(usecase usecases.IConversationsUsecase) *ConversationsPresenter {
	return &ConversationsPresenter{
		usecase: usecase,
	}
}

func (p *ConversationsPresenter) ListConversations(c fiber.Ctx) {
	var query dto.ListConversationsQueryDto

	if err := c.Bind().Query(&query); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid query parameters"})
		return
	}

	resp, err := p.usecase.ListConversations(c.Context(), query)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to list conversations"})
		return
	}

	c.JSON(resp)
}

func (p *ConversationsPresenter) CreateConversation(c fiber.Ctx) {
	var req dto.CreateConversationRequestDto
	if err := c.Bind().Body(&req); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		return
	}

	if err := p.usecase.CreateConversation(c.Context(), req); err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create conversation"})
		return
	}

	c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Conversation created successfully"})
}

func (p *ConversationsPresenter) GetConversation(c fiber.Ctx) {}

func (p *ConversationsPresenter) DeleteConversation(c fiber.Ctx) {}
