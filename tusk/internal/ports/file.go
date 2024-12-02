package ports

import (
	"context"
	"fmt"
	"io"
	"os"
	"tusk/internal/domain"
	"tusk/internal/util"

	"go.uber.org/zap"
)

type File interface {
	Path() string
	Delete()
}

type FileManager interface {
	DownloadFile(ctx context.Context, fileSrc string, localDir string) (File, error)
	UploadFile(ctx context.Context, fileSrc string, dest string, contentType string) error
	Delete(ctx context.Context, bucket string, key string) error
	CopyS3URI(ctx context.Context, sourceURL, destURL util.S3URI) error
	GetFileMetadata(ctx context.Context, url string) (*domain.FileMetadata, error)
}

type TempFile struct {
	path string
}

func NewTempFile(dir string, data io.ReadCloser) (*TempFile, error) {
	file, err := os.CreateTemp(dir, "input.*.data")
	if err != nil {
		return nil, err
	}
	defer file.Close() // #nosec G307
	defer data.Close()
	_, err = io.Copy(file, data)
	if err != nil {
		return nil, err
	}
	return &TempFile{
		file.Name(),
	}, nil
}

func (hf *TempFile) Path() string {
	return hf.path
}

func (hf *TempFile) Delete() {
	err := os.Remove(hf.path)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Failed to remove file %v", hf.path))
	}
}
