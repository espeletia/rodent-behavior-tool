package usecases

import (
	"context"
	"fmt"
	"tusk/internal/ports"
	"tusk/internal/ports/filemanager"
)

type MediaUsecase struct {
	fileManager ports.FileManager
}

func NewMediaUsecase(fileManager *filemanager.S3FileManager) *MediaUsecase {
	return &MediaUsecase{
		fileManager: fileManager,
	}
}

func (mu *MediaUsecase) DefaultFileUpload(ctx context.Context, fileSrc string, contentType string, filename string) error {
	err := mu.fileManager.UploadFile(ctx, fileSrc, fmt.Sprintf("s3://test/uploads/%s", filename), contentType)
	if err != nil {
		return err
	}

	return nil
}
