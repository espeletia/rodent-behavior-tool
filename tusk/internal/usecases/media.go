package usecases

import (
	"context"
	"fmt"
	"strings"

	"net/url"
	"tusk/internal/domain"
	"tusk/internal/ports"
	"tusk/internal/util"

	"tusk/internal/ports/filemanager"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type MediaUsecase struct {
	mediaDatabaseStore ports.MediaDatabaseStore
	fileManager        ports.FileManager
	s3Url              string
	uploadsPathPrefix  string
	bucket             string
}

func NewMediaUsecase(mediaDatabaseStore ports.MediaDatabaseStore, fileManager *filemanager.S3FileManager, s3Url, uploadsPathPrefix, bucket string) *MediaUsecase {
	return &MediaUsecase{
		mediaDatabaseStore: mediaDatabaseStore,
		fileManager:        fileManager,
		s3Url:              s3Url,
		uploadsPathPrefix:  uploadsPathPrefix,
		bucket:             bucket,
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

func (mu *MediaUsecase) copyUploadedFile(ctx context.Context, sourceURL string, targetPath string) (string, *domain.FileMetadata, error) {
	metadata, err := mu.fileManager.GetFileMetadata(ctx, sourceURL)
	if err != nil {
		return "", nil, err
	}

	s3Source, err := util.ParseString(sourceURL)
	if err != nil {
		return "", nil, err
	}

	zap.S().Infof(fmt.Sprintf("%s/%s/%s%s", mu.s3Url, mu.bucket, targetPath, metadata.FileExtension))
	s3Dest, err := util.ParseString(fmt.Sprintf("%s/%s/%s%s", mu.s3Url, mu.bucket, targetPath, metadata.FileExtension))
	if err != nil {
		return "", nil, err
	}

	if !mu.validateUrlIsUpload(*s3Source) {
		return "", nil, domain.URLIsNotUploadFoundError
	}

	err = mu.fileManager.CopyS3URI(ctx, *s3Source, *s3Dest)
	if err != nil {
		return "", nil, err
	}

	// err = mu.fileManager.Delete(ctx, *s3Source.Bucket, *s3Source.Key)
	// if err != nil {
	// 	return "", nil, err
	// }
	//
	return s3Dest.URI().String(), metadata, nil
}

func (mu *MediaUsecase) copyFile(ctx context.Context, sourceURL string, targetPath string) (string, *domain.FileMetadata, error) {
	metadata, err := mu.fileManager.GetFileMetadata(ctx, sourceURL)
	if err != nil {
		return "", nil, err
	}

	s3Source, err := util.ParseString(sourceURL)
	if err != nil {
		return "", nil, err
	}

	zap.S().Infof(fmt.Sprintf("%s/%s/%s%s", mu.s3Url, mu.bucket, targetPath, metadata.FileExtension))
	s3Dest, err := util.ParseString(fmt.Sprintf("%s/%s/%s%s", mu.s3Url, mu.bucket, targetPath, metadata.FileExtension))
	if err != nil {
		return "", nil, err
	}

	err = mu.fileManager.CopyS3URI(ctx, *s3Source, *s3Dest)
	if err != nil {
		return "", nil, err
	}

	// err = mu.fileManager.Delete(ctx, *s3Source.Bucket, *s3Source.Key)
	// if err != nil {
	// 	return "", nil, err
	// }

	return s3Dest.URI().String(), metadata, nil
}

func (mu *MediaUsecase) ProcessFile(ctx context.Context, sourceURL string, entity string, entityID string, variant string, masterMediaID *uuid.UUID) (*domain.MediaFile, error) {
	parsedUrl, err := url.ParseRequestURI(sourceURL)
	if err != nil {
		zap.L().Error("MediaUsecases.ProcessUploadedFile: Error parsing url", zap.Error(err))
		return nil, domain.InvalidUrlError
	}
	mediaId := uuid.New()
	fileUrl, metadata, err := mu.copyFile(ctx, parsedUrl.String(), fmt.Sprintf("%s/%s/%s/%s", entity, entityID, variant, mediaId.String()))
	if err != nil {
		return nil, err
	}

	mediaFile, err := mu.mediaDatabaseStore.Create(ctx, domain.MediaFile{
		ID:         mediaId,
		Url:        fileUrl,
		Variant:    variant,
		MimeType:   metadata.ContentType,
		EntityType: entity,
		MasterID:   masterMediaID,
		Size:       metadata.Size,
	})
	if err != nil {
		return nil, err
	}

	return mediaFile, nil
}

func (mu *MediaUsecase) ProcessUploadedFile(ctx context.Context, sourceURL string, entity string, entityID string) (*domain.MediaFile, error) {
	parsedUrl, err := url.ParseRequestURI(sourceURL)
	if err != nil {
		zap.L().Error("MediaUsecases.ProcessUploadedFile: Error parsing url", zap.Error(err))
		return nil, domain.InvalidUrlError
	}
	mediaId := uuid.New()
	fileUrl, metadata, err := mu.copyUploadedFile(ctx, parsedUrl.String(), fmt.Sprintf("%s/%s/%s", entity, entityID, mediaId.String()))
	if err != nil {
		return nil, err
	}

	mediaFile, err := mu.mediaDatabaseStore.Create(ctx, domain.MediaFile{
		ID:         mediaId,
		Url:        fileUrl,
		Variant:    domain.MediaVariantOriginal,
		MimeType:   metadata.ContentType,
		EntityType: entity,
		Size:       metadata.Size,
	})
	if err != nil {
		return nil, err
	}

	return mediaFile, nil
}

func (mu *MediaUsecase) validateUrlIsUpload(url util.S3URI) bool {
	return strings.HasPrefix(*url.Key, mu.uploadsPathPrefix)
}
