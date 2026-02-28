package internal

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/longtk26/chat-app/internal/routers"
	"github.com/matzefriedrich/parsley/pkg/bootstrap"
)

var _ bootstrap.Application = (*App)(nil)

type App struct {
	app      *fiber.App
	handlers []routers.RouteHandler
}

func NewApp(app *fiber.App, handlers []routers.RouteHandler) bootstrap.Application {
	for _, h := range handlers {
		h.Register(app)
	}
	return &App{app: app, handlers: handlers}
}

func (a *App) Run(_ context.Context) error {
	return a.app.Listen(":3000")
}
