package domain

import (
	"context"

	mysqlModel "banking/model/mysql"

	"github.com/gin-gonic/gin"
)

//go:generate mockgen -destination ./mock/user.go -source=./user.go -package=mock

type IUserHandler interface {
	CreateUser() gin.HandlerFunc
	// GetUser() gin.HandlerFunc
	GetUsers() gin.HandlerFunc
	Login() gin.HandlerFunc
	CreateAPIKey() gin.HandlerFunc
	DeleteAPIKey() gin.HandlerFunc
	GetAPIKeys() gin.HandlerFunc
}

type IUserService interface {
	CreateUser(ctx context.Context, user *mysqlModel.User) (err error)
	GetUsers(ctx context.Context, userID uint) (users []*mysqlModel.User, err error)
	Login(ctx context.Context, email, password string) (token string, err error)
}

type IUserQueryRepo interface {
	GetUsers(ctx context.Context, userID uint) (users []*mysqlModel.User, err error)
	GetUserByEmail(ctx context.Context, email string) (user *mysqlModel.User, err error)
}

type IUserCommandRepo interface {
	CreateUser(ctx context.Context, user *mysqlModel.User) (err error)
}
