package routers

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/longtk26/chat-app/internal/presenters"
)

var _ RouteHandler = &ConversationsRoute{}

type ConversationsRoute struct {
	presenter *presenters.ConversationsPresenter
}

func NewConversationsRoute(presenter *presenters.ConversationsPresenter) RouteHandler {
	fmt.Println("Creating ConversationsRoute...")
	return &ConversationsRoute{
		presenter: presenter,
	}
}

func (r *ConversationsRoute) Register(app *fiber.App) {
	conversations := app.Group("/api/v1/conversations")
	conversations.Get("/", r.presenter.ListConversations)
	conversations.Post("/", r.presenter.CreateConversation)
	conversations.Get("/:id", r.presenter.GetConversation)
}
