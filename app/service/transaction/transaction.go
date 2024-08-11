package transaction

import (
	"context"

	"banking/domain"
	mysqlModel "banking/model/mysql"

	"github.com/shopspring/decimal"
	"go.elastic.co/apm/v2"
)

type transactionService struct {
	transactionCmdRepo domain.ITransactionCommandRepo
}

func NewTransactionService(transactionCmdRepo domain.ITransactionCommandRepo) domain.ITransactionService {
	return &transactionService{
		transactionCmdRepo: transactionCmdRepo,
	}
}

func (s *transactionService) Transfer(ctx context.Context, fromUserID, toUserID uint, amount decimal.Decimal) (user *mysqlModel.User, err error) {
	span, ctx := apm.StartSpan(ctx, "userService.Transfer", "service")
	defer span.End()

	return s.transactionCmdRepo.Transfer(ctx, fromUserID, toUserID, amount)
}

func (s *transactionService) Deposit(ctx context.Context, userID uint, amount decimal.Decimal) (user *mysqlModel.User, err error) {
	span, ctx := apm.StartSpan(ctx, "userService.Deposit", "service")
	defer span.End()

	return s.transactionCmdRepo.Deposit(ctx, userID, amount)
}

func (s *transactionService) Withdraw(ctx context.Context, userID uint, amount decimal.Decimal) (user *mysqlModel.User, err error) {
	span, ctx := apm.StartSpan(ctx, "userService.Withdraw", "service")
	defer span.End()

	return s.transactionCmdRepo.Withdraw(ctx, userID, amount)
}
