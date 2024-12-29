package domain

import "context"

type IAuthService interface {
	JWTConfirmation(ctx context.Context, email, token string) (err error)
	APIKeyConfirmation(ctx context.Context, userID uint, key string, secret string) (err error)
}
