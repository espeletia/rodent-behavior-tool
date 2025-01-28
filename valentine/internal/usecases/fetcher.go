package usecases

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"valentine/internal/domain"
	"valentine/internal/middleware"

	"go.uber.org/zap"
)

type RestFetcher struct {
	client http.Client
	apiUrl string
}

func NewRestFetcher(client http.Client, apiUrl string) *RestFetcher {
	return &RestFetcher{
		client: client,
		apiUrl: apiUrl,
	}
}

func GenericFetch[T any](ctx context.Context, endpoint, method string, body io.Reader) (*T, error) {
	token, ok := middleware.GetUserToken(ctx)
	if !ok {
		return nil, domain.TokenNotFound
	}
	zap.L().Info("token", zap.Stringp("token", token))
	request, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+*token)
	client := http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	zap.L().Info("status", zap.Int("statuscode", response.StatusCode))

	defer response.Body.Close()
	var result T
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func GenericFetchNoAuth[T any](ctx context.Context, endpoint, method string, body io.Reader) (*T, error) {
	request, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}
	client := http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	zap.L().Info("status", zap.Int("statuscode", response.StatusCode))

	defer response.Body.Close()
	var result T
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
