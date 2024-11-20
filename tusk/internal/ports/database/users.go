package database

import (
	"context"
	"database/sql"
	"tusk/internal/domain"
	"tusk/internal/ports/database/gen/ratt-api/public/model"
	"tusk/internal/ports/database/gen/ratt-api/public/table"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

func NewUserDatabaseStore(db *sql.DB) *UserDatabaseStore {
	return &UserDatabaseStore{
		DB: db,
	}
}

type UserDatabaseStore struct {
	DB *sql.DB
}

type userWithExtraData struct {
	model.Users
	RoleName  string
	GroupName string
}

func (udbs *UserDatabaseStore) CreateUser(ctx context.Context, user domain.UserData) (*domain.User, error) {
	usrModel := model.Users{
		ID:           uuid.New(),
		Email:        user.Email,
		DisplayName:  user.DisplayName,
		Username:     user.Username,
		PasswordHash: user.Hash,
	}
	stmt := table.Users.INSERT(
		table.Users.ID,
		table.Users.Email,
		table.Users.Username,
		table.Users.PasswordHash,
		table.Users.DisplayName).
		MODEL(usrModel).
		RETURNING(
			table.Users.AllColumns,
		)
	dest := []model.Users{}
	err := stmt.QueryContext(ctx, udbs.DB, &dest)
	if err != nil {
		return nil, err
	}
	return mapUserFromDB(dest[0])
}

func (udbs *UserDatabaseStore) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	dest := []model.Users{}
	stmt := table.Users.SELECT(
		table.Users.AllColumns,
	).WHERE(table.Users.Email.EQ(postgres.String(email)))
	err := stmt.QueryContext(ctx, udbs.DB, &dest)
	if err != nil {
		return nil, err
	}
	if len(dest) == 0 {
		return nil, domain.UserNotFound
	}
	return mapUserFromDB(dest[0])
}

func (udbs *UserDatabaseStore) GetUserById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	dest := []model.Users{}
	stmt := table.Users.SELECT(
		table.Users.AllColumns,
	).WHERE(table.Users.ID.EQ(postgres.UUID(id)))
	err := stmt.QueryContext(ctx, udbs.DB, &dest)
	if err != nil {
		return nil, err
	}
	if len(dest) == 0 {
		return nil, domain.UserNotFound
	}
	return mapUserFromDB(dest[0])
}

func mapUserFromDB(usr model.Users) (*domain.User, error) {
	return &domain.User{
		ID:                usr.ID,
		Username:          usr.Username,
		DisplayName:       usr.DisplayName,
		Email:             usr.Email,
		HashedPassword:    usr.PasswordHash,
		ProfilePictureUrl: usr.ProfilePictureURL,
	}, nil
}
