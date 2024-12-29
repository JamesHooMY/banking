package apikey

import (
	"context"
	"fmt"

	"banking/domain"

	"github.com/go-redis/redis/v8"
)

type apikeyRedisQueryRepo struct {
	redisClient *redis.Client
}

func NewRedisAPIKeyQueryRepo(redisClient *redis.Client) domain.IRedisAPIKeyQueryRepo {
	return &apikeyRedisQueryRepo{redisClient: redisClient}
}

func (r *apikeyRedisQueryRepo) GetRedisAPIKey(ctx context.Context, userID uint, key string) (secret string, err error) {
	cacheKey := fmt.Sprintf("apiAuthKey:%d:%s", userID, key)

	secret, err = r.redisClient.Get(r.redisClient.Context(), cacheKey).Result()
	if err != nil {
		return "", err
	}

	return secret, nil
}
