package database

import (
	"context"
	"database/sql"
	"errors"
	"tusk/internal/domain"
	"tusk/internal/ports/database/gen/ratt-api/public/model"
	"tusk/internal/ports/database/gen/ratt-api/public/table"
	"tusk/internal/util"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type VideoDatabaseStore struct {
	DB *sql.DB
}

func NewVideoDatabaseStore(db *sql.DB) *VideoDatabaseStore {
	return &VideoDatabaseStore{
		DB: db,
	}
}

type video struct {
	model.VideoAnalysis
	Media []model.Media
}

func (vdbs *VideoDatabaseStore) Create(ctx context.Context, video domain.Video) (*domain.Video, error) {
	zap.L().Info("Inserting video")
	insertModel := MapVideoToDB(video)

	insertStmt := table.VideoAnalysis.INSERT(
		table.VideoAnalysis.AllColumns.Except(
			table.VideoAnalysis.CreatedAt,
			table.VideoAnalysis.UpdatedAt,
			table.VideoAnalysis.AnalysedVideo,
		),
	).MODEL(insertModel).RETURNING(table.VideoAnalysis.AllColumns)

	dest := []model.VideoAnalysis{}

	err := insertStmt.QueryContext(ctx, vdbs.DB, &dest)
	if err != nil {
		return nil, err
	}
	if len(dest) != 1 {
		return nil, errors.New("failed to insert new video")
	}

	result := MapVideoToDomain(dest[0])

	return &result, nil
}

func (vdbs *VideoDatabaseStore) AddAnalyzedVideo(ctx context.Context, videoID, mediaId uuid.UUID) error {
	updateModel := model.VideoAnalysis{
		AnalysedVideo: &mediaId,
	}

	updateStmt := table.VideoAnalysis.UPDATE(table.VideoAnalysis.AnalysedVideo).
		MODEL(updateModel).
		WHERE(table.VideoAnalysis.ID.EQ(postgres.UUID(videoID)))

	r, err := updateStmt.ExecContext(ctx, vdbs.DB)
	if err != nil {
		return err
	}
	rows, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return errors.New("failed to add analysed video")
	}

	return nil
}

func (vdbs *VideoDatabaseStore) GetByID(ctx context.Context, id uuid.UUID) (*domain.Video, error) {
	selectStmt := table.VideoAnalysis.SELECT(
		table.VideoAnalysis.AllColumns,
		table.Media.AllColumns,
	).FROM(
		table.VideoAnalysis.LEFT_JOIN(
			table.Media, postgres.OR(
				table.Media.ID.EQ(table.VideoAnalysis.MediaID),
				table.Media.ID.EQ(table.VideoAnalysis.AnalysedVideo),
			),
		)).WHERE(
		table.VideoAnalysis.ID.EQ(postgres.UUID(id)),
	).GROUP_BY(
		table.VideoAnalysis.ID,
		table.Media.ID,
	)

	dest := []video{}

	err := selectStmt.QueryContext(ctx, vdbs.DB, &dest)
	if err != nil {
		return nil, err
	}

	if len(dest) == 0 {
		return nil, domain.VideoNotFound
	}

	return validateAndMapVideoAnalysis(dest[0]), nil
}

func (vdbs *VideoDatabaseStore) GetVideosCursored(ctx context.Context, userId uuid.UUID, offsetLimit domain.OffsetLimit) (*domain.VideosCursored, error) {
	selectStmt := table.VideoAnalysis.SELECT(
		table.VideoAnalysis.AllColumns,
		table.Media.AllColumns,
	).FROM(
		table.VideoAnalysis.LEFT_JOIN(
			table.Media, postgres.OR(
				table.Media.ID.EQ(table.VideoAnalysis.MediaID),
				table.Media.ID.EQ(table.VideoAnalysis.AnalysedVideo),
			),
		)).WHERE(
		table.VideoAnalysis.OwnerID.EQ(postgres.UUID(userId)),
	).GROUP_BY(
		table.VideoAnalysis.ID,
		table.Media.ID,
	).LIMIT(int64(offsetLimit.Limit)).OFFSET(offsetLimit.Offset)

	dest := []video{}

	err := selectStmt.QueryContext(ctx, vdbs.DB, &dest)
	if err != nil {
		return nil, err
	}
	data := []domain.Video{}
	for _, video := range dest {
		data = append(data, *validateAndMapVideoAnalysis(video))
	}
	return &domain.VideosCursored{
		Data:   data,
		Cursor: util.BuildCursorWithOffsetCursor(data, offsetLimit.Offset, offsetLimit.Limit),
	}, nil

}

// Todo Imporove and add more fields that need no validation nor complicated queries
func MapVideoToDomain(video model.VideoAnalysis) domain.Video {
	return domain.Video{
		ID:          video.ID,
		Name:        video.Name,
		Description: video.Description,
	}

}

func MapVideoToDB(video domain.Video) model.VideoAnalysis {
	return model.VideoAnalysis{
		ID:          video.ID,
		MediaID:     video.Video.ID,
		Name:        video.Name,
		Description: video.Description,
		OwnerID:     video.OwnerId,
		CageID:      video.CageId,
	}
}

func validateAndMapVideoAnalysis(v video) *domain.Video {
	var videoMedia *domain.MediaFile
	var analysisMedia *domain.MediaFile
	for _, file := range v.Media {
		if file.ID == v.MediaID {
			videoMedia = mapDBMediaFileToDomain(file)
		}
		if v.AnalysedVideo != nil {
			if file.ID == *v.AnalysedVideo {
				analysisMedia = mapDBMediaFileToDomain(file)
			}
		}
	}
	return &domain.Video{
		ID:            v.ID,
		Video:         *videoMedia,
		OwnerId:       v.OwnerID,
		Description:   v.Description,
		Name:          v.Name,
		AnalysedVideo: analysisMedia,
	}
}

func mapDBMediaFileToDomain(m model.Media) *domain.MediaFile {
	var duration *int64 = nil
	if m.Duration != nil {
		duration = util.ToPointer(int64(*m.Duration))
	}

	master := false
	if m.MasterID == nil {
		master = true
	}
	return &domain.MediaFile{
		ID:         m.ID,
		MimeType:   m.Mimetype,
		Variant:    m.Variant,
		EntityType: m.EntityType,
		MasterID:   m.MasterID,
		Url:        m.URL,
		Created:    m.CreatedAt,
		Duration:   duration,
		Size:       int64(m.Size),
		Width:      int64(m.Width),
		Height:     int64(m.Height),
		Master:     master,
	}
}
