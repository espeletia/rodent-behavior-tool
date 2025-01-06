package auth

import (
	"context"
	"tusk/internal/domain"
)

type AuthUsecaseInterface interface {
	Login(ctx context.Context, creds domain.LoginCreds) (string, error)
	Authenticate(ctx context.Context, token string) (*domain.User, error)
	AuthenticateCage(ctx context.Context, secretToken string) (*domain.Cage, error)
}
