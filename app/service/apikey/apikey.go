package apikey

import (
	"context"

	"banking/domain"
	mysqlModel "banking/model/mysql"
	"banking/utils"

	"github.com/go-redis/redis/v8"
)

type apikeyService struct {
	apikeyRedisCmdRepo   domain.IRedisAPIKeyCommandRepo
	apikeyRedisQueryRepo domain.IRedisAPIKeyQueryRepo
	apikeyCmdRepo        domain.IAPIKeyCommandRepo
	apikeyQueryRepo      domain.IAPIKeyQueryRepo
}

func NewAPIKeyService(APIKeyRedisCmdRepo domain.IRedisAPIKeyCommandRepo, APIKeyRedisQueryRepo domain.IRedisAPIKeyQueryRepo, APIKeyCmdRepo domain.IAPIKeyCommandRepo, APIKeyQueryRepo domain.IAPIKeyQueryRepo) domain.IAPIKeyService {
	return &apikeyService{
		apikeyRedisCmdRepo:   APIKeyRedisCmdRepo,
		apikeyRedisQueryRepo: APIKeyRedisQueryRepo,
		apikeyCmdRepo:        APIKeyCmdRepo,
		apikeyQueryRepo:      APIKeyQueryRepo,
	}
}

func (s *apikeyService) CreateAPIKey(ctx context.Context, userID uint) (key string, secret string, err error) {
	// Generate key and secret
	key = utils.GenerateRandomAPIKey()
	secret = utils.GenerateRandomSecretKey()

	hashedSecret, err := utils.GenerateHashedSecretKey(secret)
	if err != nil {
		return "", "", err
	}

	// create api key
	err = s.apikeyCmdRepo.CreateAPIKey(ctx, userID, key, hashedSecret)
	if err != nil {
		return "", "", err
	}

	// set api key in redis
	err = s.apikeyRedisCmdRepo.SetRedisAPIKey(ctx, userID, key, hashedSecret)
	if err != nil {
		return "", "", err
	}

	return key, secret, nil
}

// func (s *apikeyService) GetAPIKey(ctx context.Context, userID uint, key string) (secret string, err error) {
// 	// get api key from redis
// 	secret, err = s.apikeyRedisQueryRepo.GetRedisAPIKey(ctx, userID, key)
// 	if err != redis.Nil && err != nil {
// 		return "", err
// 	}

// 	// if not found in redis, get from database
// 	if secret == "" {
// 		apiKeys, err := s.apikeyQueryRepo.GetAPIKeys(ctx, userID, key)
// 		if err != nil {
// 			return "", err
// 		}

// 		err = s.apikeyRedisCmdRepo.SetRedisAPIKey(ctx, userID, key, apiKeys[0].Secret)
// 		if err != nil {
// 			return "", err
// 		}

// 		secret = apiKeys[0].Secret
// 	}

// 	return secret, nil
// }

func (s *apikeyService) GetAPIKeys(ctx context.Context, userID uint, key string) ([]*mysqlModel.APIKey, error) {
	// get api key from redis
	if key != "" && userID != 0 {
		secret, err := s.apikeyRedisQueryRepo.GetRedisAPIKey(ctx, userID, key)
		if err != redis.Nil && err != nil {
			return nil, err
		}

		if secret != "" {
			apiKey := &mysqlModel.APIKey{
				UserID: userID,
				APIKey: key,
				Secret: secret,
			}

			return []*mysqlModel.APIKey{apiKey}, nil
		}
	}

	// get api keys from database
	apiKeys, err := s.apikeyQueryRepo.GetAPIKeys(ctx, userID, key)
	if err != nil {
		return nil, err
	}

	// set api keys in redis
	if key != "" && userID != 0 {
		for _, apiKey := range apiKeys {
			err = s.apikeyRedisCmdRepo.SetRedisAPIKey(ctx, apiKey.UserID, apiKey.APIKey, apiKey.Secret)
			if err != nil {
				return nil, err
			}
		}
	}

	return apiKeys, nil
}

func (s *apikeyService) DeleteAPIKey(ctx context.Context, userID uint, key string) (err error) {
	// delete api key from redis
	err = s.apikeyRedisCmdRepo.DeleteRedisAPIKey(ctx, userID, key)
	if err != nil {
		return err
	}

	// delete api key from database
	err = s.apikeyCmdRepo.DeleteAPIKey(ctx, userID, key)
	if err != nil {
		return err
	}

	return nil
}
