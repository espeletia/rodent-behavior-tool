package middleware

import (
	"context"
	"errors"
	"net/http"
	"valentine/internal/domain"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const (
	authHeader      = "Authorization"
	userCtxKey      = "auth"
	userTokenCtxKey = "token"
	cageTokenCtxKey = "cageToken"
	bearerPrefix    = "Bearer "
)

func WithUserToken(ctx context.Context, user *string) context.Context {
	return context.WithValue(ctx, userTokenCtxKey, user)
}

func GetUserToken(ctx context.Context) (*string, bool) {
	user, ok := ctx.Value(userTokenCtxKey).(*string)
	return user, ok
}

func Authentication() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			token, err := GetAuthToken(r)
			if err != nil {
				if !errors.Is(err, domain.TokenNotFound) {
					zap.L().Error("Failed to authenticate user", zap.Error(err))
				}
			} else {
				ctx = WithUserToken(ctx, token)
			}
			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}

func GetAuthToken(r *http.Request) (*string, error) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		if err == http.ErrNoCookie {
			return nil, domain.TokenNotFound
		}
		return nil, err
	}
	return &cookie.Value, nil
}

//palma-protect-jacob-status-clean-272
