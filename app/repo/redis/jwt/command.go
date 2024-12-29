package jwt

import (
	"context"
	"fmt"
	"time"

	"banking/domain"

	"github.com/go-redis/redis/v8"
)

type jwtCommandRepo struct {
	redisClient *redis.Client
}

func NewRedisJWTCommandRepo(redisClient *redis.Client) domain.IRedisJWTCommandRepo {
	return &jwtCommandRepo{redisClient: redisClient}
}

func (r *jwtCommandRepo) SetRedisJWT(ctx context.Context, email, token string) (err error) {
	cacheKey := fmt.Sprintf("jwt:%s", email)

	if err := r.redisClient.Set(r.redisClient.Context(), cacheKey, token, 1*time.Hour).Err(); err != nil {
		return err
	}

	return nil
}

func (r *jwtCommandRepo) DeleteRedisJWT(ctx context.Context, email string) (err error) {
	cacheKey := fmt.Sprintf("jwt:%s", email)

	if err := r.redisClient.Del(r.redisClient.Context(), cacheKey).Err(); err != nil {
		return err
	}

	return nil
}
