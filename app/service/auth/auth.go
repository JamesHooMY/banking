package auth

import (
	"context"
	"errors"

	"banking/domain"
	"banking/utils"

	"github.com/go-redis/redis/v8"
)

type authService struct {
	apikeyRedisCmdRepo   domain.IRedisAPIKeyCommandRepo
	apikeyRedisQueryRepo domain.IRedisAPIKeyQueryRepo
	apikeyQueryRepo      domain.IAPIKeyQueryRepo
	jwtRedisCmdRepo      domain.IRedisJWTCommandRepo
	jwtRedisQueryRepo    domain.IRedisJWTQueryRepo
}

func NewAuthService(APIKeyRedisCmdRepo domain.IRedisAPIKeyCommandRepo, APIKeyRedisQueryRepo domain.IRedisAPIKeyQueryRepo, APIKeyQueryRepo domain.IAPIKeyQueryRepo) domain.IAuthService {
	return &authService{
		apikeyRedisCmdRepo:   APIKeyRedisCmdRepo,
		apikeyRedisQueryRepo: APIKeyRedisQueryRepo,
		apikeyQueryRepo:      APIKeyQueryRepo,
	}
}

func (s *authService) JWTConfirmation(ctx context.Context, email, token string) (err error) {
	// Validate the JWT token
	// This is a placeholder for future implementation
	tokenRedis, err := s.jwtRedisQueryRepo.GetRedisJWT(ctx, email)
	if err != nil {
		return err
	}

	if tokenRedis != token {
		return errors.New("Invalid JWT token")
	}

	return nil
}

func (s *authService) APIKeyConfirmation(ctx context.Context, userID uint, key string, secret string) (err error) {
	// get api key from redis
	hashedSecret, err := s.apikeyRedisQueryRepo.GetRedisAPIKey(ctx, userID, key)
	if err != redis.Nil && err != nil {
		return err
	}

	// If the API key is not found in Redis, fall back to the database
	if hashedSecret == "" {
		// Query the database for the API key and secret key
		apiKeys, err := s.apikeyQueryRepo.GetAPIKeys(ctx, userID, key)
		if err != nil {
			return err
		}

		hashedSecret = apiKeys[0].Secret

		// Store the secret in Redis for future requests
		err = s.apikeyRedisCmdRepo.SetRedisAPIKey(ctx, userID, key, hashedSecret)
		if err != nil {
			// Log the error, but continue the request
			// We don't want to reject the request if Redis set fails
			return err
		}
	}

	// validate the secret key
	if !utils.VerifySecretKey(hashedSecret, secret) {
		return errors.New("Unverified API Key or Secret Key")
	}

	return nil
}
