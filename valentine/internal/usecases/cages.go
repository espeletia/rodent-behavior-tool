package usecases

import (
	"context"
	"ghiaccio/models"
	"net/http"
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
	cages, err := GenericFetch[models.Cages](ctx, cu.apiUrl+"/v1/cages", "GET", nil)
	if err != nil {
		return nil, err
	}

	return cages.Data, nil
}

func (cu *CageUsecase) GetCageMessages(ctx context.Context, cageId string) (*models.CageMessagesCursored, error) {
	result, err := GenericFetch[models.CageMessagesCursored](ctx, cu.apiUrl+"/v1/cages/"+cageId+"/messages", "GET", nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}
