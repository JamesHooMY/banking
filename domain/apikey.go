package domain

import (
	"context"

	mysqlModel "banking/model/mysql"
)

//go:generate mockgen -destination ./mock/apikey.go -source=./apikey.go -package=mock

type IAPIKeyService interface {
	CreateAPIKey(ctx context.Context, userID uint) (key string, secret string, err error)
	DeleteAPIKey(ctx context.Context, userID uint, key string) (err error)
	// GetAPIKey(ctx context.Context, userID uint, key string) (secret string, err error)
	GetAPIKeys(ctx context.Context, userID uint, key string) (apiKeys []*mysqlModel.APIKey, err error)
}

type IRedisAPIKeyQueryRepo interface {
	GetRedisAPIKey(ctx context.Context, userID uint, key string) (secret string, err error)
}

type IRedisAPIKeyCommandRepo interface {
	SetRedisAPIKey(ctx context.Context, userID uint, key string, secret string) (err error)
	DeleteRedisAPIKey(ctx context.Context, userID uint, key string) (err error)
}

type IAPIKeyQueryRepo interface {
	// GetAPIKey(ctx context.Context, userID uint, key string) (apiKey *mysqlModel.APIKey, err error)
	GetAPIKeys(ctx context.Context, userID uint, key string) (apiKeys []*mysqlModel.APIKey, err error)
}

type IAPIKeyCommandRepo interface {
	CreateAPIKey(ctx context.Context, userID uint, key string, secret string) (err error)
	DeleteAPIKey(ctx context.Context, userID uint, key string) (err error)
}
