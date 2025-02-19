package usecases

import (
	"context"
	"crypto/rand"
	"fmt"
	commonDomain "ghiaccio/domain"
	"math/big"
	"strings"
	"tusk/internal/domain"
	"tusk/internal/middleware"
	"tusk/internal/ports"
	"tusk/internal/util"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type CagesUsecase struct {
	videoUsecase         *VideoUsecase
	cageDatabaseStore    ports.CagesDatabaseStore
	activationCodeLength int64
	secretTokenLength    int64
	queueHandler         ports.QueueHandler
}

func NewCagesUsecase(videoUsecase *VideoUsecase, cageDatabaseStore ports.CagesDatabaseStore, activationCodeLength int64, secretTokenLength int64, queueHandler ports.QueueHandler) *CagesUsecase {
	return &CagesUsecase{
		cageDatabaseStore:    cageDatabaseStore,
		activationCodeLength: activationCodeLength,
		secretTokenLength:    secretTokenLength,
		queueHandler:         queueHandler,
		videoUsecase:         videoUsecase,
	}
}

func (cu *CagesUsecase) ProcessInternalCageJob(ctx context.Context, job commonDomain.CageMessageVideoAnalysisJob) error {
	zap.L().Info("Process internal cage job", zap.Any("job", job))
	videoDto := domain.CreateVideoDto{
		Name:        fmt.Sprintf("Cage-%d", job.MessageID),
		Description: util.ToPointer(fmt.Sprintf("From cage %s", job.CageID)),
		VideoUrl:    job.Url,
	}
	video, err := cu.videoUsecase.CreateNewCageVideo(ctx, videoDto, job.CageID)
	if err != nil {
		return err
	}
	err = cu.cageDatabaseStore.InsertVideoIDToCageMessage(ctx, job.MessageID, video.ID)
	return err
}

func (cu *CagesUsecase) GetCageMessages(ctx context.Context, cageId uuid.UUID, offset domain.OffsetLimit) (*domain.CageMessasgesCursored, error) {
	usr, ok := middleware.GetUser(ctx)
	if !ok {
		return nil, domain.Unauthorized
	}
	_, err := cu.cageDatabaseStore.GetCageById(ctx, cageId, usr.ID)
	if err != nil {
		return nil, err
	}

	messages, err := cu.cageDatabaseStore.FetchCageMessages(ctx, cageId, offset)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (cu *CagesUsecase) CreateNewCage(ctx context.Context) (string, string, error) {
	zap.L().Info("Creatign new cage")
	activationCode, err := GenerateActivationCode(int(cu.activationCodeLength))
	if err != nil {
		return "", "", err
	}
	secretToken, err := GenerateActivationCode(int(cu.secretTokenLength))
	if err != nil {
		return "", "", err
	}

	err = cu.cageDatabaseStore.CreateNewCage(ctx, activationCode, secretToken)
	if err != nil {
		return "", "", err
	}

	return activationCode, secretToken, nil
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

func (cu *CagesUsecase) CageSelf(ctx context.Context) (*domain.Cage, error) {
	cage, ok := middleware.GetCage(ctx)
	if !ok {
		return nil, domain.Unauthorized
	}
	return cage, nil
}

func (cu *CagesUsecase) GetCageBySecretToken(ctx context.Context, secretToken string) (*domain.Cage, error) {
	cage, err := cu.cageDatabaseStore.GetCageBySecretToken(ctx, secretToken)
	if err != nil {
		return nil, err
	}
	return cage, nil
}

func (cu *CagesUsecase) ProcessCageMessage(ctx context.Context, message domain.CageMessageData) error {
	cage, ok := middleware.GetCage(ctx)
	if !ok {
		return domain.Unauthorized
	}
	storedMessage, err := cu.cageDatabaseStore.InsertNewCageMessage(ctx, message, cage.ID)
	if err != nil {
		return err
	}
	if message.VideoUrl == nil {
		return nil
	}
	job := commonDomain.CageMessageVideoAnalysisJob{
		ID:        uuid.New(),
		CageID:    cage.ID,
		Url:       *message.VideoUrl,
		MessageID: storedMessage.ID,
		Timestamp: message.Timestamp.Unix(),
	}
	zap.L().Info("Preparing internal cage job", zap.Any("job", job))
	err = cu.queueHandler.AddInternalCageJob(ctx, job)
	return err
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
