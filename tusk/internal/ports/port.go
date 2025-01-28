package ports

import (
	"context"
	"tusk/internal/domain"

	commonDomain "ghiaccio/domain"
	"github.com/google/uuid"
)

type UserDatabaseStore interface {
	CreateUser(ctx context.Context, user domain.UserData) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*domain.User, error)
	// GetUsersByUsernamePattern(ctx context.Context, usernamePattern string) ([]domain.User, error)
}

type MediaDatabaseStore interface {
	Create(ctx context.Context, file domain.MediaFile) (*domain.MediaFile, error)
}

type VideoDatabaseStore interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Video, error)
	Create(ctx context.Context, video domain.Video) error
	AddAnalyzedVideo(ctx context.Context, videoID, mediaId uuid.UUID) error
	GetVideosCursored(ctx context.Context, userId uuid.UUID, offsetLimit domain.OffsetLimit) (*domain.VideosCursored, error)
}

type TokenGeneratorAuthInterface interface {
	CreateUserJWT(ctx context.Context, usr domain.User) (string, error)
	ValidateUserJWT(ctx context.Context, token string) (*uuid.UUID, error)
}

type CagesDatabaseStore interface {
	CreateNewCage(ctx context.Context, activation, secretToken string) error
	ActivateCage(ctx context.Context, userId uuid.UUID, activationCode string) error
	GetCagesByUserId(ctx context.Context, userId uuid.UUID) ([]domain.Cage, error)
	GetCageBySecretToken(ctx context.Context, secretToken string) (*domain.Cage, error)
	InsertNewCageMessage(ctx context.Context, cageMessage domain.CageMessageData, cageId uuid.UUID) error
	GetCageById(ctx context.Context, cageId, userId uuid.UUID) (*domain.Cage, error)
	FetchCageMessages(ctx context.Context, cageId uuid.UUID, offsetLimit domain.OffsetLimit) (*domain.CageMessasgesCursored, error)
}

type QueueHandler interface {
	AddAnalystJob(ctx context.Context, job commonDomain.AnalystJobMessage) error
	AddEncoderJob(ctx context.Context, job commonDomain.VideoEncodingMessage) error
}
