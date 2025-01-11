package domain

import (
	commonDomain "ghiaccio/domain"
	"net/http"
)

var (
	TokenNotFound = commonDomain.Error{Message: "missing token cookie", Code: http.StatusUnauthorized}
)
