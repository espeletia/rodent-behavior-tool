package encoding

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"echoes/internal/domain"
	"echoes/internal/ports"
	"echoes/internal/usecases/encoding/video"
	"echoes/internal/util"
	commonDomain "ghiaccio/domain"

	"go.uber.org/zap"
)

type QueueConsumer struct {
	files        ports.FileManager
	queueHandler ports.Queue
	videoEncoder *video.VideoMediaEncoder
	tempDir      string

	s3Url  string
	bucket string
}

func NewQueueConsumer(
	files ports.FileManager,
	videoEncoder *video.VideoMediaEncoder,
	tempDir string,
	s3Url string,
	bucket string,
	queueHandler ports.Queue,
) *QueueConsumer {
	return &QueueConsumer{
		files:        files,
		videoEncoder: videoEncoder,
		tempDir:      tempDir,
		s3Url:        s3Url,
		bucket:       bucket,
		queueHandler: queueHandler,
	}
}

func (qw *QueueConsumer) ProcessVideoQueue(ctx context.Context, job commonDomain.VideoEncodingMessage) error {
	videoJob := domain.VideoEncodingJob{
		ID:      job.ID,
		VideoID: job.VideoID,
		MediaID: job.MediaID,
		URl:     job.Url,
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

	parsedUrl, err := util.ParseString(job.URl)
	if err != nil {
		return err
	}

	file, err := qw.files.DownloadFile(ctx, *parsedUrl, dir)
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

	zap.L().Debug(fmt.Sprintf("uploading video file", file.Path()))

	urls, err := qw.uploadFinishedFiles(ctx, filesToUpload)
	if err != nil {
		return err
	}

	zap.L().Debug("Sending result message")
	err = qw.queueHandler.AddEncodingJobResult(ctx, commonDomain.VideoEncodingResultMessage{
		ID:      job.ID,
		VideoID: job.VideoID,
		MediaID: job.MediaID,
		Url:     urls[0],
	})
	if err != nil {
		return err
	}

	return nil
}

func (qw *QueueConsumer) uploadFinishedFiles(ctx context.Context, files []domain.JobResult) ([]string, error) {
	urls := []string{}
	for _, file := range files {
		parsedUrl, err := url.ParseRequestURI(file.FileDestinationSrc)
		if err != nil {
			return nil, err
		}
		url := fmt.Sprintf("%s/%s/%s%s", qw.s3Url, qw.bucket, "encoding_results", parsedUrl.Path)
		uploadUrl, err := util.ParseString(url)
		if err != nil {
			return nil, err
		}
		err = qw.files.UploadFile(ctx, file.LocalFileSrc, *uploadUrl, "video/mp4")
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}

func (qw *QueueConsumer) logJobContents(ctx context.Context, job domain.VideoEncodingJob) {
	s, _ := json.MarshalIndent(job, "", "\t")
	zap.L().Info(string(s))
}

func (qw *QueueConsumer) createTemDir(job domain.VideoEncodingJob) (string, error) {
	dir, err := os.MkdirTemp(qw.tempDir, fmt.Sprintf("%s_*", job.ID.String()))
	if err != nil {
		return "", err
	}
	return dir, nil
}
