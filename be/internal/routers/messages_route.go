package routers

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/longtk26/chat-app/internal/presenters"
)

var _ RouteHandler = &MessagesRoute{}

type MessagesRoute struct {
	presenter *presenters.MessagesPresenter
}

func NewMessagesRoute(presenter *presenters.MessagesPresenter) RouteHandler {
	fmt.Println("Creating MessagesRoute...")
	return &MessagesRoute{
		presenter: presenter,
	}
}

func (r *MessagesRoute) Register(app *fiber.App) {
	messages := app.Group("/api/v1/messages")
	messages.Get("/", r.presenter.ListMessages)
	messages.Post("/", r.presenter.SendMessage)
	messages.Put("/:id", r.presenter.UpdateMessage)
	messages.Delete("/:id", r.presenter.DeleteMessage)
}
