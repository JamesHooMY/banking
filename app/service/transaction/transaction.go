package transaction

import (
	"context"

	"banking/domain"
	mysqlModel "banking/model/mysql"

	"github.com/shopspring/decimal"
	"go.elastic.co/apm/v2"
)

type transactionService struct {
	transactionCmdRepo   domain.ITransactionCommandRepo
	transactionQueryRepo domain.ITransactionQueryRepo
}

func NewTransactionService(TransactionCmdRepo domain.ITransactionCommandRepo, TransactionQueryRepo domain.ITransactionQueryRepo) domain.ITransactionService {
	return &transactionService{
		transactionCmdRepo:   TransactionCmdRepo,
		transactionQueryRepo: TransactionQueryRepo,
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

func (s *transactionService) GetTransactions(ctx context.Context) (transactions []*mysqlModel.Transaction, err error) {
	span, ctx := apm.StartSpan(ctx, "userService.GetTransactions", "service")
	defer span.End()

	return s.transactionQueryRepo.GetTransactions(ctx)
}
