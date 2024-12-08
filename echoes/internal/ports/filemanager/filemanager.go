package filemanager

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"echoes/internal/ports"
	"echoes/internal/util"

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

func (hfm *S3FileManager) DownloadFile(ctx context.Context, src util.S3URI, dir string) (ports.File, error) {
	zap.L().Info("fetching file", zap.String("bucket", *src.Bucket), zap.String("key", *src.Key))

	obj, err := hfm.s3client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(*src.Bucket),
		Key:    aws.String(*src.Key),
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

func (hfm *S3FileManager) UploadFile(ctx context.Context, fileSrc string, dest util.S3URI, contentType string) error {
	data, err := os.Open(filepath.Clean(fileSrc))
	if err != nil {
		return err
	}
	zap.L().Info("uploading file", zap.String("host", *dest.Bucket), zap.String("key", *dest.Key))
	uploader := manager.NewUploader(hfm.s3client)
	result, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(*dest.Bucket),
		Key:         aws.String(*dest.Key),
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
