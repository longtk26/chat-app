package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/longtk26/chat-app/infrastructure/db/sqlc/out"
	"github.com/longtk26/chat-app/internal/domain/entities"
	"github.com/longtk26/chat-app/internal/domain/repo"
)

var _ repo.IUserRepo = &UserRepository{}

type UserRepository struct {
	queries *out.Queries
}

func NewUserRepository(queries *out.Queries) repo.IUserRepo {
	fmt.Println("Creating UserRepository...")
	return &UserRepository{queries: queries}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *entities.UserEntity) error {
	uid, err := uuid.Parse(user.ID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	pgUUID := pgtype.UUID{
		Bytes: uid,
		Valid: true,
	}

	_, err = r.queries.CreateUser(ctx, out.CreateUserParams{
		ID:       pgUUID,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entities.UserEntity, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	uid, err := uuid.FromBytes(user.ID.Bytes[:])
	if err != nil {
		return nil, fmt.Errorf("failed to parse user ID: %w", err)
	}

	return &entities.UserEntity{
		ID:       uid.String(),
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}

func (r *UserRepository) ListUsers(ctx context.Context) ([]*entities.UserEntity, error) {
	users, err := r.queries.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	var userEntities []*entities.UserEntity
	for _, user := range users {
		uid, err := uuid.FromBytes(user.ID.Bytes[:])
		if err != nil {
			return nil, fmt.Errorf("failed to parse user ID: %w", err)
		}

		userEntities = append(userEntities, &entities.UserEntity{
			ID:       uid.String(),
			Username: user.Username,
			Email:    user.Email,
			Password: user.Password,
		})
	}

	return userEntities, nil
}
