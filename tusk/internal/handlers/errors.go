package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"tusk/internal/domain"

	"go.uber.org/zap"
)

type ErrResp struct {
	Error string `json:"error"`
}

func SendError(rw http.ResponseWriter, err error) {
	responseCode := http.StatusInternalServerError
	responseMsg := err.Error()

	var httpError domain.Error
	if errors.As(err, &httpError) {
		responseCode = httpError.Code
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(responseCode)
	err = json.NewEncoder(rw).Encode(ErrResp{responseMsg})
	if err != nil {
		zap.L().Error("Internal error:", zap.Error(err))
	}
}

func SendTooManyRequests(rw http.ResponseWriter, err error) {
	responseCode := http.StatusTooManyRequests
	rw.WriteHeader(responseCode)
}

func SendForbidden(rw http.ResponseWriter) {
	responseCode := http.StatusForbidden
	rw.WriteHeader(responseCode)
}
