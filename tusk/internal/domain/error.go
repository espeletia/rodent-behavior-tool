package domain

import "net/http"

type Error struct {
	Message string
	Code    int
}

func (e Error) Error() string { return e.Message }

var (
	ErrInvalidEmail    = Error{Message: "Invalid email", Code: http.StatusNotAcceptable}
	InvalidCredentials = Error{Message: "Invalid credentials", Code: http.StatusUnauthorized}
	Unauthorized       = Error{Message: "Unauthorized", Code: http.StatusUnauthorized}
	UserNotFound       = Error{Message: "User not found", Code: http.StatusNotFound}
)
