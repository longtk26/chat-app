package usecases

import (
	"context"

	"github.com/longtk26/chat-app/internal/domain/entities"
	"github.com/longtk26/chat-app/internal/domain/repo"
)

type IUsersUseCase interface {
	ListUsers(ctx context.Context) ([]*entities.UserEntity, error)
}

var _ IUsersUseCase = &UsersUseCase{}

type UsersUseCase struct {
	userRepo repo.IUserRepo
}

func NewUsersUseCase(userRepo repo.IUserRepo) IUsersUseCase {
	return &UsersUseCase{
		userRepo: userRepo,
	}
}

func (u *UsersUseCase) ListUsers(c context.Context) ([]*entities.UserEntity, error) {
	return u.userRepo.ListUsers(c)
}
