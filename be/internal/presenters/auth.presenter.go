package presenters

import (
	"github.com/gofiber/fiber/v3"
	"github.com/longtk26/chat-app/be/internal/presenters/dto"
	"github.com/longtk26/chat-app/be/internal/usecases"
)

type AuthPresenter struct {
	authUseCase usecases.IAuthUseCase
}

func NewAuthPresenter(authUseCase usecases.IAuthUseCase) *AuthPresenter {
	return &AuthPresenter{authUseCase: authUseCase}
}

func (p *AuthPresenter) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	token, err := p.authUseCase.Login(req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	return c.JSON(fiber.Map{"token": token})
}

func (p *AuthPresenter) Register(c *fiber.Ctx) error {
	type RegisterRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := p.authUseCase.Register(req.Username, req.Password); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not register user"})
	}

	return c.JSON(fiber.Map{"message": "User registered successfully"})
}
