package repo

import (
	"context"

	"github.com/longtk26/chat-app/internal/domain/entities"
)

type IUserRepo interface {
	CreateUser(ctx context.Context, user *entities.UserEntity) error
	GetUserByEmail(ctx context.Context, email string) (*entities.UserEntity, error)
	ListUsers(ctx context.Context) ([]*entities.UserEntity, error)
}
