package domain

import "context"

//go:generate mockgen -destination ./mock/auth.go -source=./auth.go -package=mock

type IAuthService interface {
	JWTConfirmation(ctx context.Context, email, token string) (err error)
	APIKeyConfirmation(ctx context.Context, userID uint, key string, secret string) (err error)
}
