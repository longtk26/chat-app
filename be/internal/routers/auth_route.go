package routers

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/longtk26/chat-app/internal/presenters"
)

var _ RouteHandler = &AuthRoute{}

type AuthRoute struct {
	presenter *presenters.AuthPresenter
}

func NewAuthRoute(presenter *presenters.AuthPresenter) RouteHandler {
	fmt.Println("Creating AuthRoute...")
	return &AuthRoute{presenter: presenter}
}

func (r *AuthRoute) Register(app *fiber.App) {
	fmt.Println("Registering auth routes...")
	auth := app.Group("/api/v1/auth")
	auth.Post("/login", r.presenter.Login)
	auth.Post("/register", r.presenter.Register)
}
