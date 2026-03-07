package routers

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/longtk26/chat-app/internal/presenters"
)

var _ RouteHandler = &UsersRoute{}

type UsersRoute struct {
	presenter *presenters.UsersPresenter
}

func NewUsersRoute(presenter *presenters.UsersPresenter) RouteHandler {
	fmt.Println("Creating UsersRoute...")
	return &UsersRoute{presenter: presenter}
}

func (r *UsersRoute) Register(app *fiber.App) {
	fmt.Println("Registering users routes...")
	users := app.Group("/api/v1/users")
	users.Get("/", r.presenter.ListUsers)
}
