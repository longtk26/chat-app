package presenters

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/longtk26/chat-app/internal/presenters/dto"
	"github.com/longtk26/chat-app/internal/usecases"
)

type AuthPresenter struct {
	authUseCase usecases.IAuthUseCase
}

func NewAuthPresenter(authUseCase usecases.IAuthUseCase) *AuthPresenter {
	fmt.Println("Creating AuthPresenter...")
	return &AuthPresenter{authUseCase: authUseCase}
}

func (p *AuthPresenter) Login(c fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	token, err := p.authUseCase.Login(req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	return c.Status(fiber.StatusOK).JSON(dto.LoginResponse{
		AccessToken: token,
		UserName:    req.Username,
	})
}

func (p *AuthPresenter) Register(c fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := p.authUseCase.Register(req.Username, req.Password); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not register user"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User registered successfully"})
}
