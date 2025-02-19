package domain

import (
	"net/http"

	commonDomain "ghiaccio/domain"
)

var (
	ErrInvalidEmail          = commonDomain.Error{Message: "Invalid email", Code: http.StatusNotAcceptable}
	BadRequest               = commonDomain.Error{Message: "bad request", Code: http.StatusBadRequest}
	InvalidCredentials       = commonDomain.Error{Message: "Invalid credentials", Code: http.StatusUnauthorized}
	Unauthorized             = commonDomain.Error{Message: "Unauthorized", Code: http.StatusUnauthorized}
	UserNotFound             = commonDomain.Error{Message: "User not found", Code: http.StatusNotFound}
	VideoNotFound            = commonDomain.Error{Message: "Video not found", Code: http.StatusNotFound}
	InvalidUrlError          = commonDomain.Error{Message: "Invalid url", Code: http.StatusBadRequest}
	URLIsNotUploadFoundError = commonDomain.Error{Message: "Url is not upload url", Code: http.StatusBadRequest}
	UrlNotFoundError         = commonDomain.Error{Message: "Url not found", Code: http.StatusBadRequest}
	CageNotFound             = commonDomain.Error{Message: "Cage not found", Code: http.StatusNotFound}
	InvalidFileType          = commonDomain.Error{Message: "Invalid file type, allowed types: image/jpeg, image/png, video/mp4, video/mpeg", Code: http.StatusBadRequest}
)
