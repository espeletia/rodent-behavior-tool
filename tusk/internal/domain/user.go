package domain

import "github.com/google/uuid"

type User struct {
	ID                uuid.UUID
	DisplayName       string
	Email             string
	Username          string
	HashedPassword    string
	ProfilePictureUrl *string
}

type UserData struct {
	DisplayName string
	Email       string
	Username    string
	Hash        string
}

type LoginCreds struct {
	Email    string
	Password string
}
