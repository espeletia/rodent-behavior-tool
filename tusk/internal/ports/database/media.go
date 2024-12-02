package database

import (
	"context"
	"database/sql"
	"fmt"
	"tusk/internal/domain"
	"tusk/internal/ports/database/gen/ratt-api/public/model"
	"tusk/internal/ports/database/gen/ratt-api/public/table"
	"tusk/internal/util"

	"github.com/google/uuid"
	// "github.com/go-jet/jet/v2/postgres"
	// "github.com/google/uuid"
)

type MediaDatabaseStore struct {
	DB *sql.DB
}

func NewMediaDatabaseStore(db *sql.DB) *MediaDatabaseStore {
	return &MediaDatabaseStore{
		DB: db,
	}
}

func (fdbs *MediaDatabaseStore) Create(ctx context.Context, file domain.MediaFile) (*domain.MediaFile, error) {
	insertModel, err := mapMediaFileToDB(file)
	if err != nil {
		return nil, err
	}
	insertModel.ID = uuid.New()
	insertStm := table.Media.INSERT(table.Media.AllColumns.Except(table.Media.CreatedAt, table.Media.UpdatedAt)).
		MODEL(insertModel).
		RETURNING(table.Media.AllColumns)

	files := []model.Media{}

	err = insertStm.QueryContext(ctx, fdbs.DB, &files)
	if err != nil {
		return nil, err
	}

	if len(files) != 1 {
		err := fmt.Errorf("failed to insert media file")
		return nil, err
	}

	return mapDBMediaFile(files[0])
}

func mapDBMediaFile(file model.Media) (*domain.MediaFile, error) {
	master := false
	if file.MasterID == nil {
		master = true
	}
	var duration *int64 = nil
	if file.Duration != nil {
		duration = util.ToPointer(int64(*file.Duration))
	}

	return &domain.MediaFile{
		ID:         file.ID,
		MimeType:   file.Mimetype,
		Variant:    file.Variant,
		EntityType: file.EntityType,
		MasterID:   file.MasterID,
		Url:        file.URL,
		Created:    file.CreatedAt,
		Duration:   duration,
		Size:       int64(file.Size),
		Width:      int64(file.Width),
		Height:     int64(file.Height),
		Master:     master,
	}, nil
}

func mapMediaFileToDB(file domain.MediaFile) (*model.Media, error) {
	var duration *int32 = nil
	if file.Duration != nil {
		duration = util.ToPointer(int32(*file.Duration))
	}
	return &model.Media{
		ID:         file.ID,
		Mimetype:   file.MimeType,
		Variant:    file.Variant,
		EntityType: file.EntityType,
		MasterID:   file.MasterID,
		URL:        file.Url,
		Duration:   duration,
		Size:       int32(file.Size),
		Width:      int32(file.Width),
		Height:     int32(file.Height),
	}, nil
}
