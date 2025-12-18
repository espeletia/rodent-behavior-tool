package usecases

import (
	"bytes"
	"context"
	"encoding/json"
	"ghiaccio/models"
	"net/http"

	"go.uber.org/zap"
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

func (vu *VideoUsecase) CreateVideo(ctx context.Context, video models.CreateVideoDto) error {
	jsonData, err := json.Marshal(video)
	if err != nil {
		return err
	}

	jsonReader := bytes.NewReader(jsonData)
	err = GenericFetch204(ctx, vu.apiUrl+"/v1/video", "PUT", jsonReader)
	zap.L().Info("ok")
	if err != nil {
		zap.L().Info("error", zap.Error(err))
		return err
	}
	return nil

}
