package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/longtk26/chat-app/configs"
	"github.com/longtk26/chat-app/internal/domain/entities"
	"github.com/longtk26/chat-app/internal/domain/repo"
)

type IAuthUseCase interface {
	Login(username, password string) (string, error)
	Register(username, password string) error
}

type AuthUseCase struct {
	userRepo  repo.IUserRepo
	appConfig configs.AppConfig
}

func NewAuthUseCase(userRepo repo.IUserRepo, appConfig configs.AppConfig) IAuthUseCase {
	fmt.Println("Creating auth use case")
	return &AuthUseCase{userRepo: userRepo, appConfig: appConfig}
}

func (a *AuthUseCase) Login(username, password string) (string, error) {
	ctx := context.Background()

	user, err := a.userRepo.GetUserByEmail(ctx, username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// TODO: replace with proper password hashing comparison (e.g. bcrypt)
	if user.Password != password {
		return "", errors.New("invalid credentials")
	}

	// TODO: generate a real JWT token
	return "token_" + user.ID, nil
}

func (a *AuthUseCase) Register(username, password string) error {
	ctx := context.Background()

	foundUser, err := a.userRepo.GetUserByEmail(ctx, username)
	if err == nil && foundUser != nil {
		return errors.New("user already exists")
	}

	newUser := &entities.UserEntity{
		ID:       uuid.New().String(),
		Username: username,
		Email:    username,
		Password: password, // TODO: hash password before storing
	}

	if err := a.userRepo.CreateUser(ctx, newUser); err != nil {
		fmt.Printf("Error creating user: %v\n", err)
		return errors.New("could not register user")
	}

	return nil
}
