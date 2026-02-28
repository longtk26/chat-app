package routers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/longtk26/chat-app/internal/presenters"
)

var _ RouteHandler = (*AuthRoute)(nil)

type AuthRoute struct {
	presenter *presenters.AuthPresenter
}

func NewAuthRoute(presenter *presenters.AuthPresenter) RouteHandler {
	return &AuthRoute{presenter: presenter}
}

func (r *AuthRoute) Register(app *fiber.App) {
	auth := app.Group("/api/v1/auth")
	auth.Post("/login", r.presenter.Login)
	auth.Post("/register", r.presenter.Register)
}
