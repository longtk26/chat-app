package presenters

import "github.com/gofiber/fiber/v3"

type MessagesPresenter struct{}

func NewMessagesPresenter() *MessagesPresenter {
	return &MessagesPresenter{}
}

func (p *MessagesPresenter) ListMessages(c fiber.Ctx) error {
	return nil
}

func (p *MessagesPresenter) SendMessage(c fiber.Ctx) error {
	return nil
}

func (p *MessagesPresenter) UpdateMessage(c fiber.Ctx) error {
	return nil
}

func (p *MessagesPresenter) DeleteMessage(c fiber.Ctx) error {
	return nil
}
