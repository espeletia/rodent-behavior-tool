package domain

import "net/http"

type Error struct {
	Message string
	Code    int
}

func (e Error) Error() string { return e.Message }

var (
	ErrInvalidEmail          = Error{Message: "Invalid email", Code: http.StatusNotAcceptable}
	InvalidCredentials       = Error{Message: "Invalid credentials", Code: http.StatusUnauthorized}
	Unauthorized             = Error{Message: "Unauthorized", Code: http.StatusUnauthorized}
	UserNotFound             = Error{Message: "User not found", Code: http.StatusNotFound}
	VideoNotFound            = Error{Message: "Video not found", Code: http.StatusNotFound}
	InvalidUrlError          = Error{Message: "Invalid url", Code: http.StatusBadRequest}
	URLIsNotUploadFoundError = Error{Message: "Url is not upload url", Code: http.StatusBadRequest}
	UrlNotFoundError         = Error{Message: "Url not found", Code: http.StatusBadRequest}
)
