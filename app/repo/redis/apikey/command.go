package apikey

import (
	"context"
	"fmt"
	"time"

	"banking/domain"

	"github.com/go-redis/redis/v8"
)

type apikeyCommandRepo struct {
	redisClient *redis.Client
}

func NewRedisAPIKeyCommandRepo(redisClient *redis.Client) domain.IRedisAPIKeyCommandRepo {
	return &apikeyCommandRepo{redisClient: redisClient}
}

func (r *apikeyCommandRepo) SetRedisAPIKey(ctx context.Context, userID uint, key string, secret string) (err error) {
	cacheKey := fmt.Sprintf("apiAuthKey:%d:%s", userID, key)

	if err := r.redisClient.Set(r.redisClient.Context(), cacheKey, secret, 1*time.Hour).Err(); err != nil {
		return err
	}

	return nil
}

func (r *apikeyCommandRepo) DeleteRedisAPIKey(ctx context.Context, userID uint, key string) (err error) {
	cacheKey := fmt.Sprintf("apiAuthKey:%d:%s", userID, key)

	if err := r.redisClient.Del(r.redisClient.Context(), cacheKey).Err(); err != nil {
		return err
	}

	return nil
}
