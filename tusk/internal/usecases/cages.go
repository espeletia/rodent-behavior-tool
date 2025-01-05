package usecases

import (
	"context"
	"crypto/rand"
	"math/big"
	"strings"
	"tusk/internal/domain"
	"tusk/internal/middleware"
	"tusk/internal/ports"

	"go.uber.org/zap"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type CagesUsecase struct {
	cageDatabaseStore    ports.CagesDatabaseStore
	activationCodeLength int64
}

func NewCagesUsecase(cageDatabaseStore ports.CagesDatabaseStore, activationCodeLength int64) *CagesUsecase {
	return &CagesUsecase{
		cageDatabaseStore:    cageDatabaseStore,
		activationCodeLength: activationCodeLength,
	}
}

func (cu *CagesUsecase) CreateNewCage(ctx context.Context) (string, error) {
	zap.L().Info("Creatign new cage")
	activationCode, err := GenerateActivationCode(int(cu.activationCodeLength))
	if err != nil {
		return "", err
	}
	err = cu.cageDatabaseStore.CreateNewCage(ctx, activationCode)
	if err != nil {
		return "", err
	}

	return activationCode, nil
}

func (cu *CagesUsecase) RegisterCage(ctx context.Context, activationCode string) error {
	zap.L().Info("Registering cage")
	usr, ok := middleware.GetUser(ctx)
	if !ok {
		return domain.Unauthorized
	}
	err := cu.cageDatabaseStore.ActivateCage(ctx, usr.ID, activationCode)
	if err != nil {
		return err
	}
	return nil
}

func (cu *CagesUsecase) GetCagesForUser(ctx context.Context) ([]domain.Cage, error) {
	usr, ok := middleware.GetUser(ctx)
	if !ok {
		return nil, domain.Unauthorized
	}
	cages, err := cu.cageDatabaseStore.GetCagesByUserId(ctx, usr.ID)
	if err != nil {
		return nil, err
	}
	return cages, nil
}

func GenerateActivationCode(length int) (string, error) {
	var code strings.Builder
	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code.WriteByte(charset[index.Int64()])
	}
	return code.String(), nil
}
