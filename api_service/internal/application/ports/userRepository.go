package ports

import (
	"api-service/internal/domain"
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound      error = errors.New("user not found")
	ErrUserAlreadyExists       = errors.New("user already exists")
	ErrInternal                = errors.New("internal error")
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)

	Create(ctx context.Context, user *domain.User) error

	Delete(ctx context.Context, id uuid.UUID) error
}
