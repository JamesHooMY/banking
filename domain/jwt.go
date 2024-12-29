package domain

import "context"

type IRedisJWTCommandRepo interface {
	SetRedisJWT(ctx context.Context, email, token string) (err error)
	DeleteRedisJWT(ctx context.Context, email string) (err error)
}

type IRedisJWTQueryRepo interface {
	GetRedisJWT(ctx context.Context, email string) (token string, err error)
}
