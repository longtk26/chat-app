package internal

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/longtk26/chat-app/internal/routers"
	"github.com/matzefriedrich/parsley/pkg/bootstrap"
)

type application struct {
	app *fiber.App
}

var _ bootstrap.Application = &application{}

func NewApp(app *fiber.App, routeHandlers []routers.RouteHandler) bootstrap.Application {
	fmt.Println("Registering route handlers...", routeHandlers)
	for _, routeHandler := range routeHandlers {
		routeHandler.Register(app)
	}

	return &application{
		app: app,
	}
}

func (a *application) Run(_ context.Context) error {
	return a.app.Listen(":5502")
}
