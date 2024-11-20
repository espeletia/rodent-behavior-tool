package usecases

import (
	"context"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"regexp"
	"tusk/internal/domain"
	"tusk/internal/middleware"
	"tusk/internal/ports"
)

func NewUserUsecase(usi ports.UserDatabaseStore,
) *UserUsecase {
	return &UserUsecase{
		store:      usi,
		emailRegex: *regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`),
	}
}

type UserUsecase struct {
	store      ports.UserDatabaseStore
	salt       string
	emailRegex regexp.Regexp
}

func (uu *UserUsecase) CreateUser(ctx context.Context, user domain.UserData, password string) (*domain.User, error) {
	if !uu.emailRegex.MatchString(user.Email) {
		return nil, domain.ErrInvalidEmail
	}
	hashedPassword, err := uu.HashPassword(password)
	if err != nil {
		return nil, err
	}
	user.Hash = hashedPassword
	stored, err := uu.store.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return stored, nil
}

func (uu *UserUsecase) Me(ctx context.Context) (*domain.User, error) {
	usr, ok := middleware.GetUser(ctx)
	if !ok {
		return nil, domain.Unauthorized
	}
	return usr, nil
}

func (uu *UserUsecase) GetUserById(ctx context.Context, Id uuid.UUID) (*domain.User, error) {
	return uu.store.GetUserById(ctx, Id)
}

func (uu *UserUsecase) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return uu.store.GetUserByEmail(ctx, email)
}

// HashPassword hashes a password using bcrypt
func (uu *UserUsecase) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ComparePassword compares a plain password with a hashed password
func (uu *UserUsecase) ComparePassword(hashedPassword, password string) error {
	// bcrypt.CompareHashAndPassword returns nil if the password matches the hash
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
