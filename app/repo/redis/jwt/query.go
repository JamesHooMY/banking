package jwt

import (
	"context"
	"fmt"

	"banking/domain"

	"github.com/go-redis/redis/v8"
)

type jwtRedisQueryRepo struct {
	redisClient *redis.Client
}

func NewRedisJWTQueryRepo(redisClient *redis.Client) domain.IRedisJWTQueryRepo {
	return &jwtRedisQueryRepo{redisClient: redisClient}
}

func (r *jwtRedisQueryRepo) GetRedisJWT(ctx context.Context, email string) (token string, err error) {
	cacheKey := fmt.Sprintf("jwt:%s", email)

	token, err = r.redisClient.Get(r.redisClient.Context(), cacheKey).Result()
	if err != nil {
		return "", err
	}

	return token, nil
}
