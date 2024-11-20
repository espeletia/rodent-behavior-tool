package tokens

import (
	"context"
	"log"
	"time"
	"tusk/internal/domain"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	expirationClaim = "expiration"
	idClaim         = "id"
	// roleClaim	= "role"
)

type TokenGenerator struct {
	JWTSecret     string
	JWTExpiration time.Duration
}

func NewTokenGenerator(jwtSecret string, JWTExpiration time.Duration) *TokenGenerator {
	return &TokenGenerator{
		JWTSecret:     jwtSecret,
		JWTExpiration: JWTExpiration,
	}
}

func (au *TokenGenerator) CreateUserJWT(ctx context.Context, usr domain.User) (string, error) {
	log.Println("Creating JWT token")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims[idClaim] = usr.ID

	claims[expirationClaim] = time.Now().UTC().Add(au.JWTExpiration).Unix()

	resultToken, err := token.SignedString([]byte(au.JWTSecret))
	if err != nil {
		log.Println("Error creating JWT token")
		return "", err
	}
	log.Println("Successfully created JWT token")
	return resultToken, nil
}

func (au *TokenGenerator) ValidateUserJWT(ctx context.Context, token string) (*uuid.UUID, error) {
	log.Println("Authenticating JWT token")
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, domain.Unauthorized
		}
		return []byte(au.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if parsedToken == nil {
		return nil, domain.Unauthorized
	}

	log.Println("Successfully parsed JWT token")
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("Error parsing claims")
		return nil, domain.Unauthorized
	}

	//check expiration
	expiration, ok := claims[expirationClaim].(float64)
	if !ok {
		log.Println("Error parsing expiration")
		return nil, domain.Unauthorized
	}

	if int64(expiration) < time.Now().UTC().Unix() {
		log.Println("token is expired")
		return nil, domain.Unauthorized
	}

	userIDString, ok := claims[idClaim].(string)
	if !ok {
		log.Println("Error parsing user ID")
		return nil, domain.Unauthorized
	}
	zap.L().Info(userIDString)
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return nil, err
	}

	return &userID, nil
}
