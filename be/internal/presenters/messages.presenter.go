package presenters

import (
	"github.com/gofiber/fiber/v3"
	"github.com/longtk26/chat-app/internal/hub"
	"github.com/longtk26/chat-app/internal/presenters/dto"
	"github.com/longtk26/chat-app/internal/usecases"
)

type MessagesPresenter struct {
	usecase usecases.IMessagesUseCase
	hub     *hub.SocketHub
}

func NewMessagesPresenter(usecase usecases.IMessagesUseCase, h *hub.SocketHub) *MessagesPresenter {
	return &MessagesPresenter{usecase: usecase, hub: h}
}

func (p *MessagesPresenter) ListMessages(c fiber.Ctx) {
	var query dto.ListMessagesQueryDto
	if err := c.Bind().Query(&query); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid query parameters"})
		return
	}

	resp, err := p.usecase.ListMessages(c.Context(), query)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		return
	}

	c.Status(fiber.StatusOK).JSON(resp)
}

func (p *MessagesPresenter) SendMessage(c fiber.Ctx) {
	var payload dto.SendMessageRequestDto
	if err := c.Bind().Body(&payload); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		return
	}

	resp, err := p.usecase.SendMessage(c.Context(), payload)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		return
	}

	// Broadcast the new message to all other participants in the conversation room.
	senderUUID := p.hub.GetSocketUUIDByUserID(resp.Message.ConversationID, resp.Message.SenderID)
	p.hub.BroadcastToRoom(resp.Message.ConversationID, senderUUID, "new_message", resp.Message)

	c.Status(fiber.StatusCreated).JSON(resp)
}

func (p *MessagesPresenter) UpdateMessage(c fiber.Ctx) {
	messageID := c.Params("id")

	var payload dto.UpdateMessageRequestDto
	if err := c.Bind().Body(&payload); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		return
	}

	resp, err := p.usecase.UpdateMessage(c.Context(), messageID, payload)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		return
	}

	c.Status(fiber.StatusOK).JSON(resp)
}

func (p *MessagesPresenter) DeleteMessage(c fiber.Ctx) {
	messageID := c.Params("id")

	if err := p.usecase.DeleteMessage(c.Context(), messageID); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		return
	}

	c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Message deleted"})
}
