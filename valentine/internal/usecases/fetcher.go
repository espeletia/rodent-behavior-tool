package usecases

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"ghiaccio/models"
	"io"
	"mime/multipart"
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

func GenericFetch204(ctx context.Context, endpoint, method string, body io.Reader) error {
	token, ok := middleware.GetUserToken(ctx)
	if !ok {
		return domain.TokenNotFound
	}
	zap.L().Info("token", zap.Stringp("token", token))
	request, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "Bearer "+*token)
	client := http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	zap.L().Info("status", zap.Int("statuscode", response.StatusCode))

	defer response.Body.Close()

	return nil
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

func (rf *RestFetcher) GenericFileUpload(ctx context.Context, reader io.Reader, filename string) (string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Create form file field with provided filename
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}

	// Copy the contents of the reader to the multipart form part
	_, err = io.Copy(part, reader)
	if err != nil {
		return "", err
	}

	// Close the multipart writer to finalize the body
	err = writer.Close()
	if err != nil {
		return "", err
	}

	// Create the HTTP request
	req, err := http.NewRequest("PUT", rf.apiUrl+"/v1/upload", &body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload failed: %s", responseBody)
	}

	result := models.UploadResponse{}
	err = json.NewDecoder(resp.Body).Decode(&result)

	zap.L().Info("upload successful", zap.String("url", result.UploadUrl))

	return result.UploadUrl, nil
}
