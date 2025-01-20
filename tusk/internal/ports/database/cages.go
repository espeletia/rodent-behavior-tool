package database

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"tusk/internal/domain"
	"tusk/internal/ports/database/gen/ratt-api/public/model"
	"tusk/internal/ports/database/gen/ratt-api/public/table"
	"tusk/internal/util"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

type CagesDatabaseStore struct {
	db *sql.DB
}

func NewCagesDatabaseStore(db *sql.DB) *CagesDatabaseStore {
	return &CagesDatabaseStore{
		db: db,
	}
}

func (cdbs *CagesDatabaseStore) CreateNewCage(ctx context.Context, activation, secretToken string) error {
	insertModel := model.RodentCages{
		ID:             uuid.New(),
		ActivationCode: activation,
		Name:           "Inactive " + activation,
		SecretToken:    secretToken,
	}

	insertStmt := table.RodentCages.INSERT(
		table.RodentCages.ID,
		table.RodentCages.ActivationCode,
		table.RodentCages.Name,
		table.RodentCages.SecretToken).
		MODEL(insertModel)

	r, err := insertStmt.ExecContext(ctx, cdbs.db)
	if err != nil {
		return err
	}
	rows, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return errors.New("failed to insert new cage")
	}
	return nil
}

func (cdbs *CagesDatabaseStore) ActivateCage(ctx context.Context, userId uuid.UUID, activationCode string) error {
	updateModel := model.RodentCages{
		UserID:       util.ToPointer(userId),
		RegisteredAt: util.ToPointer(time.Now()),
	}

	updateStmt := table.RodentCages.UPDATE(
		table.RodentCages.UserID,
		table.RodentCages.RegisteredAt).
		MODEL(updateModel).
		WHERE(table.RodentCages.ActivationCode.EQ(postgres.String(activationCode)))
	r, err := updateStmt.ExecContext(ctx, cdbs.db)
	if err != nil {
		return err
	}
	rows, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return domain.CageNotFound
	}

	return nil
}

func (cdbs *CagesDatabaseStore) GetCageBySecretToken(ctx context.Context, secretToken string) (*domain.Cage, error) {
	dest := []model.RodentCages{}

	selectStmt := table.RodentCages.SELECT(table.RodentCages.AllColumns).WHERE(table.RodentCages.SecretToken.EQ(postgres.String(secretToken)))

	err := selectStmt.QueryContext(ctx, cdbs.db, &dest)
	if err != nil {
		return nil, err
	}

	if len(dest) != 1 {
		return nil, domain.CageNotFound
	}

	mappedCage := mapCageToDomain(dest[0])
	return &mappedCage, nil
}

func (cdbs *CagesDatabaseStore) GetCagesByUserId(ctx context.Context, userId uuid.UUID) ([]domain.Cage, error) {
	dest := []model.RodentCages{}

	selectStmt := table.RodentCages.SELECT(table.RodentCages.AllColumns).WHERE(table.RodentCages.UserID.EQ(postgres.UUID(userId)))

	err := selectStmt.QueryContext(ctx, cdbs.db, &dest)
	if err != nil {
		return nil, err
	}

	result := []domain.Cage{}
	for _, cage := range dest {
		result = append(result, mapCageToDomain(cage))
	}

	return result, nil
}

func (cdbs *CagesDatabaseStore) GetCageById(ctx context.Context, cageId, userId uuid.UUID) (*domain.Cage, error) {
	dest := []model.RodentCages{}

	selectStmt := table.RodentCages.SELECT(table.RodentCages.AllColumns).WHERE(table.RodentCages.UserID.EQ(postgres.UUID(userId)).AND(table.RodentCages.ID.EQ(postgres.UUID(cageId))))
	err := selectStmt.QueryContext(ctx, cdbs.db, &dest)
	if err != nil {
		return nil, err
	}

	if len(dest) != 1 {
		return nil, domain.CageNotFound
	}

	mappedCage := mapCageToDomain(dest[0])
	return &mappedCage, nil

}

func (cdbs *CagesDatabaseStore) InsertNewCageMessage(ctx context.Context, cageMessage domain.CageMessageData, cageID uuid.UUID) error {
	insertModel := model.CageMessages{
		CageID:   cageID,
		Revision: int32(cageMessage.Revision),
		Water:    int32(cageMessage.Water),
		Food:     int32(cageMessage.Food),
		Light:    int32(cageMessage.Light),
		Temp:     int32(cageMessage.Temp),
		Humidity: int32(cageMessage.Humidity),
		TimeSent: cageMessage.Timestamp,
		VideoURL: cageMessage.VideoUrl,
	}

	insertStmt := table.CageMessages.INSERT(table.CageMessages.AllColumns.Except(
		table.CageMessages.ID,
		table.CageMessages.VideoID,
		table.CageMessages.UpdatedAt,
		table.CageMessages.CreatedAt,
	)).MODEL(insertModel)

	r, err := insertStmt.ExecContext(ctx, cdbs.db)
	if err != nil {
		return err
	}
	rows, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return errors.New("failed to insert cage message")
	}

	return nil
}

func (cdbs *CagesDatabaseStore) FetchCageMessages(ctx context.Context, cageId uuid.UUID, offsetLimit domain.OffsetLimit) (*domain.CageMessasgesCursored, error) {
	stmt := table.CageMessages.SELECT(table.CageMessages.AllColumns).
		WHERE(table.CageMessages.CageID.EQ(postgres.UUID(cageId))).
		ORDER_BY(table.CageMessages.TimeSent.DESC()).
		OFFSET(offsetLimit.Offset).
		LIMIT(int64(offsetLimit.Limit))

	dest := []model.CageMessages{}

	err := stmt.QueryContext(ctx, cdbs.db, &dest)
	if err != nil {
		return nil, err
	}
	data := []domain.CageMessage{}
	for _, message := range dest {
		data = append(data, mapCageMessageToDomain(message))
	}
	return &domain.CageMessasgesCursored{
		Data:   data,
		Cursor: util.BuildCursorWithOffsetCursor(data, offsetLimit.Offset, offsetLimit.Limit),
	}, nil
}

func mapCageMessageToDomain(message model.CageMessages) domain.CageMessage {
	return domain.CageMessage{
		CageID:    message.CageID,
		Revision:  int64(message.Revision),
		Water:     int64(message.Water),
		Food:      int64(message.Food),
		Light:     int64(message.Light),
		Temp:      int64(message.Temp),
		Humidity:  int64(message.Humidity),
		VideoUrl:  message.VideoURL,
		VideoID:   message.VideoID,
		Timestamp: message.TimeSent,
	}
}

func mapCageToDomain(cage model.RodentCages) domain.Cage {
	registered := time.Now() // TODO: make this better
	if cage.RegisteredAt != nil {
		registered = *cage.RegisteredAt
	}
	return domain.Cage{
		ID:          cage.ID,
		UserID:      cage.UserID,
		Name:        cage.Name,
		Description: cage.Description,
		Register:    registered,
	}
}
