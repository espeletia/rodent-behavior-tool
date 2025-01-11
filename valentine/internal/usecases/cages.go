package usecases

import (
	"context"
	"encoding/json"
	"ghiaccio/models"
	"net/http"
	"valentine/internal/domain"
	"valentine/internal/middleware"
)

type CageUsecase struct {
	client http.Client
	apiUrl string
}

func NewCageUsecase(client http.Client, url string) *CageUsecase {
	return &CageUsecase{
		client: client,
		apiUrl: url,
	}
}

func (cu *CageUsecase) GetCages(ctx context.Context) ([]models.Cage, error) {
	token, ok := middleware.GetUserToken(ctx)
	if !ok {
		return nil, domain.TokenNotFound
	}
	cageReq, err := http.NewRequest("GET", cu.apiUrl+"/v1/cages", nil)
	if err != nil {
		return nil, err
	}
	cageReq.Header.Set("Authorization", "Bearer "+*token)
	client := http.Client{}
	cageResp, err := client.Do(cageReq)
	if err != nil {
		return nil, err
	}
	defer cageResp.Body.Close()
	cages := models.Cages{
		Data: []models.Cage{},
	}
	err = json.NewDecoder(cageResp.Body).Decode(&cages)
	if err != nil {
		return nil, err
	}

	return cages.Data, nil
}
