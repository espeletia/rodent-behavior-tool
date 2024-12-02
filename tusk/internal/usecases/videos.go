package usecases

import (
	"context"
	"tusk/internal/domain"
	"tusk/internal/middleware"
	"tusk/internal/ports"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type VideoUsecase struct {
	mediaUsecase       *MediaUsecase
	videoDatabaseStore ports.VideoDatabaseStore
}

func NewVideoUsecase(media *MediaUsecase, videoDatabaseStore ports.VideoDatabaseStore) *VideoUsecase {
	return &VideoUsecase{
		mediaUsecase:       media,
		videoDatabaseStore: videoDatabaseStore,
	}
}

func (vu *VideoUsecase) CreateNewVideo(ctx context.Context, data domain.CreateVideoDto) error {
	usr, ok := middleware.GetUser(ctx)
	if !ok {
		return domain.Unauthorized
	}
	videoId := uuid.New()
	videoMedia, err := vu.mediaUsecase.ProcessUploadedFile(ctx, data.VideoUrl, "video", videoId.String())
	if err != nil {
		return err
	}
	err = vu.videoDatabaseStore.Create(ctx,
		domain.Video{
			ID:            videoId,
			Video:         *videoMedia,
			OwnerId:       usr.ID,
			Description:   data.Description,
			Name:          data.Name,
			AnalysedVideo: nil,
		},
	)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	return nil

}

func (vu *VideoUsecase) GetByID(ctx context.Context, id uuid.UUID) (*domain.Video, error) {
	_, ok := middleware.GetUser(ctx)
	if !ok {
		return nil, domain.Unauthorized
	}
	return vu.videoDatabaseStore.GetByID(ctx, id)
}
