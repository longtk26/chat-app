package presenters

import (
	"github.com/gofiber/fiber/v3"
	"github.com/longtk26/chat-app/internal/usecases"
)

type UsersPresenter struct {
	usecase usecases.IUsersUseCase
}

func NewUsersPresenter(usecase usecases.IUsersUseCase) *UsersPresenter {
	return &UsersPresenter{usecase: usecase}
}

func (p *UsersPresenter) ListUsers(c fiber.Ctx) {
	users, err := p.usecase.ListUsers(c.Context())
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to list users"})
		return
	}

	c.Status(fiber.StatusOK).JSON(users)
}
