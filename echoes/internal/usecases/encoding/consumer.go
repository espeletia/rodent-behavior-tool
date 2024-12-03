package encoding

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"echoes/internal/domain"
	"echoes/internal/ports"
	"echoes/internal/usecases/encoding/video"
	commonDomain "ghiaccio/domain"

	"go.uber.org/zap"
)

type QueueConsumer struct {
	files        ports.FileManager
	videoEncoder *video.VideoMediaEncoder
	tempDir      string
	// statusProducerUsecases *job.StatusProducerUsecases
}

func NewQueueConsumer(
	files ports.FileManager,
	videoEncoder *video.VideoMediaEncoder,
	tempDir string,
	// statusProducerUsecases *job.StatusProducerUsecases,
) *QueueConsumer {
	return &QueueConsumer{
		files:        files,
		videoEncoder: videoEncoder,
		tempDir:      tempDir,
		// statusProducerUsecases: statusProducerUsecases,
	}
}

func (qw *QueueConsumer) ProcessVideoQueue(ctx context.Context, job commonDomain.VideoEncodingMessage) error {
	videoJob := domain.VideoEncodingJob{
		ID:  job.ID,
		URl: job.Url,
	}
	err := qw.processVideoQueue(ctx, videoJob)
	if err != nil {
		zap.L().Error("Error from job status producer usecase", zap.Error(err))
		return err
	}
	return nil
}

func (qw *QueueConsumer) processVideoQueue(ctx context.Context, job domain.VideoEncodingJob) error {
	zap.L().Info(fmt.Sprintf("Starting to handle video job #%v", job.ID), zap.String("ID", job.ID.String()))

	qw.logJobContents(ctx, job)

	dir, err := qw.createTemDir(job)
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	file, err := qw.files.DownloadFile(ctx, job.URl, dir)
	if err != nil {
		return err
	}
	zap.L().Debug(fmt.Sprintf("Encoding video file %v", file.Path()))
	filesToUpload := []domain.JobResult{}

	files, err := qw.videoEncoder.EncodeVideoWith256(ctx, file.Path(), dir)
	if err != nil {
		return err
	}
	filesToUpload = append(filesToUpload, domain.JobResult{
		ID:                 job.ID,
		LocalFileSrc:       files,
		FileDestinationSrc: job.URl,
	})

	err = qw.uploadFinishedFiles(ctx, filesToUpload)
	if err != nil {
		return err
	}

	return nil
}

func (qw *QueueConsumer) uploadFinishedFiles(ctx context.Context, files []domain.JobResult) error {
	for _, file := range files {
		err := qw.files.UploadFile(ctx, file.LocalFileSrc, fmt.Sprintf("S3:/test/videos/1/outputs/boxes_videoplayback/results_%d.mp4", file.ID.String()), "video/mp4")
		if err != nil {
			return err
		}
	}
	return nil
}

func (qw *QueueConsumer) logJobContents(ctx context.Context, job domain.VideoEncodingJob) {
	s, _ := json.MarshalIndent(job, "", "\t")
	zap.L().Info(string(s))
}

func (qw *QueueConsumer) createTemDir(job domain.VideoEncodingJob) (string, error) {
	dir, err := os.MkdirTemp(qw.tempDir, fmt.Sprintf("%d_*", job.ID))
	if err != nil {
		return "", err
	}
	return dir, nil
}
