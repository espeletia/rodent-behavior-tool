package domain

import (
	"net/http"

	commonDomain "ghiaccio/domain"
)

var (
	ErrInvalidEmail          = Error{Message: "Invalid email", Code: http.StatusNotAcceptable}
	InvalidCredentials       = Error{Message: "Invalid credentials", Code: http.StatusUnauthorized}
	Unauthorized             = Error{Message: "Unauthorized", Code: http.StatusUnauthorized}
	UserNotFound             = Error{Message: "User not found", Code: http.StatusNotFound}
	VideoNotFound            = Error{Message: "Video not found", Code: http.StatusNotFound}
	InvalidUrlError          = Error{Message: "Invalid url", Code: http.StatusBadRequest}
	URLIsNotUploadFoundError = Error{Message: "Url is not upload url", Code: http.StatusBadRequest}
	UrlNotFoundError         = Error{Message: "Url not found", Code: http.StatusBadRequest}
	InvalidFileType          = Error{Message: "Invalid file type, allowed types: image/jpeg, image/png, video/mp4, video/mpeg", Code: http.StatusBadRequest}
	CageNotFound             = Error{Message: "Cage not found", Code: http.StatusNotFound}
)
