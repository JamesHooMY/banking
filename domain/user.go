package domain

import (
	"context"

	mysqlModel "banking/model/mysql"

	"github.com/gin-gonic/gin"
)

//go:generate mockgen -destination ./mock/user.go -source=./user.go -package=mock

type IUserHandler interface {
	CreateUser() gin.HandlerFunc
	GetUser() gin.HandlerFunc
	GetUsers() gin.HandlerFunc
}

type IUserService interface {
	CreateUser(ctx context.Context, user *mysqlModel.User) (err error)
	GetUser(ctx context.Context, userID uint) (user *mysqlModel.User, err error)
	GetUsers(ctx context.Context) (users []*mysqlModel.User, err error)
}

type IUserQueryRepo interface {
	GetUser(ctx context.Context, userID uint) (user *mysqlModel.User, err error)
	GetUsers(ctx context.Context) (users []*mysqlModel.User, err error)
}

type IUserCommandRepo interface {
	CreateUser(ctx context.Context, user *mysqlModel.User) (err error)
}
