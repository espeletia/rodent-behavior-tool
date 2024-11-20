package ports

import (
	"context"
	"tusk/internal/domain"

	"github.com/google/uuid"
)

type UserDatabaseStore interface {
	CreateUser(ctx context.Context, user domain.UserData) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*domain.User, error)
	// GetUsersByUsernamePattern(ctx context.Context, usernamePattern string) ([]domain.User, error)
}

type TokenGeneratorAuthInterface interface {
	CreateUserJWT(ctx context.Context, usr domain.User) (string, error)
	ValidateUserJWT(ctx context.Context, token string) (*uuid.UUID, error)
}
