package domain

import (
	"context"

	mysqlModel "banking/model/mysql"

	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
)

//go:generate mockgen -destination ./mock/transaction.go -source=./transaction.go -package=mock

type ITransactionHandler interface {
	Transfer() gin.HandlerFunc
	Deposit() gin.HandlerFunc
	Withdraw() gin.HandlerFunc
	GetTransactions() gin.HandlerFunc
}

type ITransactionService interface {
	Transfer(ctx context.Context, fromUserID, toUserID uint, amount decimal.Decimal) (user *mysqlModel.User, err error)
	Deposit(ctx context.Context, userID uint, amount decimal.Decimal) (user *mysqlModel.User, err error)
	Withdraw(ctx context.Context, userID uint, amount decimal.Decimal) (user *mysqlModel.User, err error)
	GetTransactions(ctx context.Context) (transactions []*mysqlModel.Transaction, err error)
}

type ITransactionQueryRepo interface {
	GetTransactions(ctx context.Context) (transactions []*mysqlModel.Transaction, err error)
}

type ITransactionCommandRepo interface {
	Transfer(ctx context.Context, fromUserID, toUserID uint, amount decimal.Decimal) (user *mysqlModel.User, err error)
	Deposit(ctx context.Context, userID uint, amount decimal.Decimal) (user *mysqlModel.User, err error)
	Withdraw(ctx context.Context, userID uint, amount decimal.Decimal) (user *mysqlModel.User, err error)
}
