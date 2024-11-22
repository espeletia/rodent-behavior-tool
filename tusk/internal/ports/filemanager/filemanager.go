package filemanager

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"tusk/internal/ports"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3FileManager struct {
	s3client *s3.Client
}

func NewS3FileManager(s3client *s3.Client) *S3FileManager {
	return &S3FileManager{
		s3client: s3client,
	}
}

func (hfm *S3FileManager) DownloadFile(ctx context.Context, src string, dir string) (ports.File, error) {
	u, err := url.Parse(src)
	if err != nil {
		return nil, err
	}
	zap.L().Info("fetching file", zap.String("host", u.Host), zap.String("key", u.Path))

	obj, err := hfm.s3client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(u.Host),
		Key:    aws.String(u.Path),
	})
	if err != nil {
		return nil, err
	}

	tempFile, err := ports.NewTempFile(dir, obj.Body)
	if err != nil {
		return nil, err
	}
	return tempFile, nil
}

func (hfm *S3FileManager) UploadFile(ctx context.Context, fileSrc string, dest string, contentType string) error {
	u, err := url.Parse(dest)
	if err != nil {
		return err
	}

	data, err := os.Open(filepath.Clean(fileSrc))
	if err != nil {
		return err
	}
	key := strings.TrimPrefix(u.Path, "/")
	zap.L().Info("uploading file", zap.String("host", u.Host), zap.String("key", key))
	uploader := manager.NewUploader(hfm.s3client)
	result, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(u.Host),
		Key:         aws.String(key),
		Body:        data,
		ContentType: &contentType,
		ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return err
	}

	zap.L().Info(fmt.Sprintf("File uploaded %s", result.Location), zap.String("url", result.Location))

	return nil
}
