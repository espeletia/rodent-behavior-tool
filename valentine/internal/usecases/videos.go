package usecases

import (
	"context"
	"ghiaccio/models"
	"net/http"
)

type VideoUsecase struct {
	client http.Client
	apiUrl string
}

func NewVideoUsecase(client http.Client, url string) *VideoUsecase {
	return &VideoUsecase{
		client: client,
		apiUrl: url,
	}
}

func (vu *VideoUsecase) GetVideos(ctx context.Context) (*models.CursoredVideoAnalysis, error) {
	result, err := GenericFetch[models.CursoredVideoAnalysis](ctx, vu.apiUrl+"/v1/videos", "GET", nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}
