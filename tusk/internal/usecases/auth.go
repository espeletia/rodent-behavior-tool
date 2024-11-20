package usecases

import (
	"context"
	"tusk/internal/domain"
	"tusk/internal/ports"

	"go.uber.org/zap"
)

type AuthUsecase struct {
	Users  *UserUsecase
	Tokens ports.TokenGeneratorAuthInterface
}

func NewAuthUsecase(uu *UserUsecase, tokenGenerator ports.TokenGeneratorAuthInterface) *AuthUsecase {
	return &AuthUsecase{
		Users:  uu,
		Tokens: tokenGenerator,
	}
}

// TODO: handle status better
func (au *AuthUsecase) Login(ctx context.Context, creds domain.LoginCreds) (string, error) {
	usr, err := au.Users.GetUserByEmail(ctx, creds.Email)
	if err != nil {
		zap.L().Error("Error whilst fetching from database", zap.Error(err))
		return "failure", err
	}

	err = au.Users.ComparePassword(usr.HashedPassword, creds.Password)
	if err != nil {
		zap.L().Info("Password Mismatch")
		return "You suck", err
	}
	zap.L().Info("Credentials match")
	token, err := au.Tokens.CreateUserJWT(ctx, *usr)
	if err != nil {
		return "skill issue", err
	}
	return token, nil
}

func (au *AuthUsecase) Authenticate(ctx context.Context, token string) (*domain.User, error) {
	userId, err := au.Tokens.ValidateUserJWT(ctx, token)
	if err != nil {
		return nil, err
	}

	user, err := au.Users.GetUserById(ctx, *userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}
