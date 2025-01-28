package filemanager

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"tusk/internal/domain"
	"tusk/internal/ports"
	"tusk/internal/util"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gabriel-vasile/mimetype"
	"go.uber.org/zap"
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
	zap.L().Info("url", zap.String("url", src))
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

func (hfm *S3FileManager) GetFileMetadata(ctx context.Context, url string) (*domain.FileMetadata, error) {
	zap.L().Info("Getting file metadata", zap.String("URL", url))

	s3Url, err := util.ParseString(url)
	if err != nil {
		return nil, err
	}

	zap.L().Info("url params", zap.Stringp("bucket", s3Url.Bucket), zap.Stringp("key", s3Url.Key))

	obj, err := hfm.s3client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: s3Url.Bucket,
		Key:    s3Url.Key,
	})
	if err != nil {
		zap.L().Error("Error getting file metadata", zap.Stringp("Bucket", s3Url.Bucket), zap.Stringp("Key", s3Url.Key), zap.Error(err))
		if strings.Contains(err.Error(), "404") {
			return nil, domain.UrlNotFoundError
		}
		return nil, err
	}

	mime := mimetype.Lookup(*obj.ContentType)
	if mime == nil {
		return nil, fmt.Errorf("Unknown content type")
	}

	size := int64(0)
	if obj.ContentLength != nil {
		size = *obj.ContentLength
	}

	return &domain.FileMetadata{
		ContentType:   mime.String(),
		FileExtension: mime.Extension(),
		Size:          size,
	}, nil
}

func (hfm *S3FileManager) CopyS3URI(ctx context.Context, sourceURL, destURL util.S3URI) error {
	// _, err := hfm.s3client.CopyObject(ctx, &s3.CopyObjectInput{
	// 	Bucket:     destURL.Bucket,
	// 	CopySource: aws.String(fmt.Sprintf("%s/%s", *sourceURL.Bucket, *sourceURL.Key)),
	// 	Key:        destURL.Key,
	// 	ACL:        types.ObjectCannedACLPublicRead,
	// })

	obj, err := hfm.s3client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: sourceURL.Bucket,
		Key:    sourceURL.Key,
	})
	if err != nil {
		return err
	}

	uploader := manager.NewUploader(hfm.s3client)
	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      destURL.Bucket,
		Key:         destURL.Key,
		Body:        obj.Body,
		ContentType: obj.ContentType,
		ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func (hfm *S3FileManager) Delete(ctx context.Context, bucket string, key string) error {
	zap.L().Info("Deleting file in S3", zap.String("bucket", bucket), zap.String("key", key))
	_, err := hfm.s3client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return err
	}

	return nil
}
