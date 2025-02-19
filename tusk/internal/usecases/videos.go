package usecases

import (
	"context"
	"tusk/internal/domain"
	"tusk/internal/middleware"
	"tusk/internal/ports"

	commonDomain "ghiaccio/domain"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type VideoUsecase struct {
	mediaUsecase       *MediaUsecase
	videoDatabaseStore ports.VideoDatabaseStore
	queueHandler       ports.QueueHandler
}

func NewVideoUsecase(media *MediaUsecase, videoDatabaseStore ports.VideoDatabaseStore, queueHandler ports.QueueHandler) *VideoUsecase {
	return &VideoUsecase{
		mediaUsecase:       media,
		videoDatabaseStore: videoDatabaseStore,
		queueHandler:       queueHandler,
	}
}

func (vu *VideoUsecase) CreateNewVideo(ctx context.Context, data domain.CreateVideoDto) error {
	usr, ok := middleware.GetUser(ctx)
	if !ok {
		return domain.Unauthorized
	}
	_, err := vu.createNewVideo(ctx, data, &usr.ID, nil)

	return err
}

func (vu *VideoUsecase) CreateNewCageVideo(ctx context.Context, data domain.CreateVideoDto, cageID uuid.UUID) (*domain.Video, error) {
	return vu.createNewVideo(ctx, data, nil, &cageID)
}

func (vu *VideoUsecase) createNewVideo(ctx context.Context, data domain.CreateVideoDto, ownerID, cageID *uuid.UUID) (*domain.Video, error) {
	if ownerID == nil && cageID == nil {
		return nil, domain.BadRequest
	}
	videoId := uuid.New()
	videoMedia, err := vu.mediaUsecase.ProcessUploadedFile(ctx, data.VideoUrl, "video", videoId.String())
	if err != nil {
		return nil, err
	}
	video, err := vu.videoDatabaseStore.Create(ctx,
		domain.Video{
			ID:            videoId,
			Video:         *videoMedia,
			OwnerId:       ownerID,
			CageId:        cageID,
			Description:   data.Description,
			Name:          data.Name,
			AnalysedVideo: nil,
		},
	)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	err = vu.queueHandler.AddAnalystJob(ctx, commonDomain.AnalystJobMessage{
		ID:      uuid.New(),
		VideoID: videoId,
		Url:     videoMedia.Url,
		MediaID: videoMedia.ID,
	})
	if err != nil {
		return nil, err
	}

	return video, nil
}

func (vu *VideoUsecase) GetVideosCursored(ctx context.Context, offset domain.OffsetLimit) (*domain.VideosCursored, error) {
	usr, ok := middleware.GetUser(ctx)
	if !ok {
		return nil, domain.Unauthorized
	}

	videos, err := vu.videoDatabaseStore.GetVideosCursored(ctx, usr.ID, offset)
	if err != nil {
		return nil, err
	}

	return videos, nil
}

func (vu *VideoUsecase) GetByID(ctx context.Context, id uuid.UUID) (*domain.Video, error) {
	_, ok := middleware.GetUser(ctx)
	if !ok {
		return nil, domain.Unauthorized
	}
	return vu.videoDatabaseStore.GetByID(ctx, id)
}

// maybe rewrite this into a queue consumer and separate from usecase?
func (vu *VideoUsecase) ProcessAnalystJobResultQueue(ctx context.Context, job commonDomain.AnalystJobResultMessage) error {
	result := domain.AnalystResult{
		ID:      job.ID,
		VideoID: job.VideoID,
		MediaID: job.MediaID,
		Url:     job.Url,
	}
	err := vu.processAnalystJobResultQueue(ctx, result)
	if err != nil {
		zap.L().Error("Error from video analysis usecase", zap.Error(err))
		return err
	}
	return nil
}

func (vu *VideoUsecase) ProcessEncodingJobResultQueue(ctx context.Context, job commonDomain.VideoEncodingResultMessage) error {
	result := domain.EncodingResult{
		ID:      job.ID,
		VideoID: job.VideoID,
		MediaID: job.MediaID,
		Url:     job.Url,
	}
	err := vu.processEncodingResult(ctx, result)
	if err != nil {
		zap.L().Error("Error from video analysis usecase", zap.Error(err))
		return err
	}
	return nil
}

func (vu *VideoUsecase) processEncodingResult(ctx context.Context, job domain.EncodingResult) error {
	mediaFile, err := vu.mediaUsecase.ProcessFile(ctx, job.Url, "video", job.VideoID.String(), domain.MediaVariantAnalysedX264, &job.MediaID)
	if err != nil {
		return err
	}

	err = vu.videoDatabaseStore.AddAnalyzedVideo(ctx, job.VideoID, mediaFile.ID)
	if err != nil {
		return err
	}

	return nil
}

func (vu *VideoUsecase) processAnalystJobResultQueue(ctx context.Context, jobResult domain.AnalystResult) error {
	mediaFile, err := vu.mediaUsecase.ProcessFile(ctx, jobResult.Url, "video", jobResult.VideoID.String(), domain.MediaVariantAnalysedRaw, &jobResult.MediaID)
	if err != nil {
		return err
	}

	err = vu.queueHandler.AddEncoderJob(ctx, commonDomain.VideoEncodingMessage{
		ID:      uuid.New(),
		VideoID: jobResult.VideoID,
		MediaID: jobResult.MediaID,
		Url:     mediaFile.Url,
	})
	if err != nil {
		return err
	}

	zap.L().Info("sending videoEncoding message", zap.String("videoID", jobResult.VideoID.String()))

	return nil
}
