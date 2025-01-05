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

func (cdbs *CagesDatabaseStore) CreateNewCage(ctx context.Context, activation string) error {
	insertModel := model.RodentCages{
		ID:             uuid.New(),
		ActivationCode: activation,
		Name:           "Inactive " + activation,
	}

	insertStmt := table.RodentCages.INSERT(
		table.RodentCages.ID,
		table.RodentCages.ActivationCode,
		table.RodentCages.Name).
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

func mapCageToDomain(cage model.RodentCages) domain.Cage {
	registered := time.Now() // TODO: make this better
	if cage.RegisteredAt != nil {
		registered = *cage.RegisteredAt
	}
	return domain.Cage{
		ID:          cage.ID,
		Name:        cage.Name,
		Description: cage.Description,
		Register:    registered,
	}
}
