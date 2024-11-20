package models

import "github.com/google/uuid"

// swagger:model
type UserData struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Password    string `json:"password"`
}

// swagger:model
type User struct {
	ID          uuid.UUID `json:"id"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
}

// swagger:model
type LoginCreds struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// swagger:model
type LoginResp struct {
	Token string `json:"token"`
}
