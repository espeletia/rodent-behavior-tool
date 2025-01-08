package middleware

import (
	"context"
	"net/http"
	"strings"
	"tusk/internal/domain"
	"tusk/internal/usecases/auth"

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

func WithUser(ctx context.Context, user *domain.User) context.Context {
	return context.WithValue(ctx, userCtxKey, user)
}

func GetUser(ctx context.Context) (*domain.User, bool) {
	user, ok := ctx.Value(userCtxKey).(*domain.User)
	return user, ok
}

func WithCageToken(ctx context.Context, cageToken *string) context.Context {
	return context.WithValue(ctx, cageTokenCtxKey, cageToken)
}

func GetCageToken(ctx context.Context) (*string, bool) {
	cageToken, ok := ctx.Value(cageTokenCtxKey).(*string)
	return cageToken, ok
}

func WithCage(ctx context.Context, cage *domain.Cage) context.Context {
	return context.WithValue(ctx, cageTokenCtxKey, cage)
}

func GetCage(ctx context.Context) (*domain.Cage, bool) {
	cageToken, ok := ctx.Value(cageTokenCtxKey).(*domain.Cage)
	return cageToken, ok
}

func Authentication(auth auth.AuthUsecaseInterface) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if token := parseAuthHeader(r); token != "" {
				user, err := auth.Authenticate(r.Context(), token)
				if err == nil {
					zap.S().Infof("User stored to ctx")
					ctx = WithUser(ctx, user)
				} else {
					zap.L().Error("Failed to authenticate user", zap.Error(err))
				}
			}
			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}

func AuthenticationForCages(auth auth.AuthUsecaseInterface) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if token := parseAuthHeader(r); token != "" {
				cage, err := auth.AuthenticateCage(r.Context(), token)
				if err == nil {
					ctx = WithCage(ctx, cage)
				} else {
					zap.L().Error("Failed to authenticate cage", zap.Error(err))
				}
			}
			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}

func parseAuthHeader(r *http.Request) string {
	authHeader := r.Header.Get(authHeader)
	if strings.HasPrefix(authHeader, bearerPrefix) {
		split := strings.Split(authHeader, " ")
		_, token := split[0], split[1]
		return token
	}
	return ""
}
