package ports

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"echoes/internal/domain"
	"go.uber.org/zap"
)

type Queue interface {
	HandleVideoJob(ctx context.Context, handler func(ctx context.Context, job domain.VideoEncodingJob) error, errChan chan error) error
}

type File interface {
	Path() string
	Delete()
}

type FileManager interface {
	DownloadFile(ctx context.Context, fileSrc string, localDir string) (File, error)
	UploadFile(ctx context.Context, fileSrc string, dest string, contentType string) error
}

type TempFile struct {
	path string
}

func NewTempFile(dir string, data io.ReadCloser) (*TempFile, error) {
	file, err := ioutil.TempFile(dir, "input.*.data")
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
