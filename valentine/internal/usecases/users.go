package usecases

import (
	"bytes"
	"context"
	"encoding/json"
	"ghiaccio/models"
	"net/http"
	"valentine/internal/domain"
	"valentine/internal/middleware"
)

type UserUsecase struct {
	client http.Client
	apiUrl string
}

func NewUserUsecase(client http.Client, apiUrl string) *UserUsecase {
	return &UserUsecase{client: client, apiUrl: apiUrl}
}

func (uu *UserUsecase) Login(ctx context.Context, email, password string) (*string, error) {
	request := models.LoginCreds{
		Email:    email,
		Password: password,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	tokenResp, err := http.Post(uu.apiUrl+"/v1/login", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer tokenResp.Body.Close()

	resultToken := models.LoginResp{}
	err = json.NewDecoder(tokenResp.Body).Decode(&resultToken)
	if err != nil {
		return nil, err
	}

	return &resultToken.Token, nil
}

func (uu *UserUsecase) Me(ctx context.Context) (*models.User, error) {
	token, ok := middleware.GetUserToken(ctx)
	if !ok {
		return nil, domain.TokenNotFound
	}
	req, err := http.NewRequest("GET", "http://tusk:8080/v1/me", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+*token)

	resp, err := uu.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	usr := models.User{}

	err = json.NewDecoder(resp.Body).Decode(&usr)
	if err != nil {
		return nil, err
	}

	return &usr, nil
}