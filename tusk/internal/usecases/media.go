package usecases

import (
	"context"
	"fmt"
	"tusk/internal/ports"
	"tusk/internal/ports/filemanager"
)

type MediaUsecase struct {
	fileManager       ports.FileManager
	s3Url             string
	uploadsPathPrefix string
	bucket            string
}

func NewMediaUsecase(fileManager *filemanager.S3FileManager, s3Url, uploadsPathPrefix, bucket string) *MediaUsecase {
	return &MediaUsecase{
		fileManager:       fileManager,
		s3Url:             s3Url,
		uploadsPathPrefix: uploadsPathPrefix,
		bucket:            bucket,
	}
}

func (mu *MediaUsecase) DefaultFileUpload(ctx context.Context, fileSrc string, contentType string, filename string) (string, error) {
	s3Url := fmt.Sprintf("s3://test/uploads/%s", filename)
	url := fmt.Sprintf("%s/%s/%s/%s", mu.s3Url, mu.bucket, mu.uploadsPathPrefix, filename)
	err := mu.fileManager.UploadFile(ctx, fileSrc, s3Url, contentType)
	if err != nil {
		return "", err
	}

	return url, nil
}
