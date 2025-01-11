package handlers

import (
	"context"
	"errors"
	"fmt"
	commonHandlers "ghiaccio/handlers"
	"net/http"
	"valentine/internal/domain"

	"go.uber.org/zap"
)

type CommonHandler struct {
}

func NewCommonHandler() CommonHandler {
	return CommonHandler{}
}

func (ch *CommonHandler) Handle(handler func(hrw http.ResponseWriter, hreq *http.Request) error) http.Handler {

	operationName := commonHandlers.GetFunctionName(handler)
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		err := handler(rw, req)
		if err != nil {
			if errors.Is(err, domain.TokenNotFound) {

				zap.L().Warn(fmt.Sprintf("%s: unauthorized", operationName), zap.Error(err))
				http.Redirect(rw, req, "/login", http.StatusSeeOther)
				return
			}
			if errors.Is(err, context.Canceled) {
				zap.L().Warn(fmt.Sprintf("%s: ctx canceled", operationName), zap.Error(err))
			} else {
				zap.L().Error(fmt.Sprintf("%s: failed", operationName), zap.Error(err))
			}

			commonHandlers.SendError(rw, err)
		}
	})
}
