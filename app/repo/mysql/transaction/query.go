package transaction

import (
	"context"

	domain "banking/domain"
	mysqlModel "banking/model/mysql"

	"go.elastic.co/apm/v2"
	"gorm.io/gorm"
)

type transactionQueryRepo struct {
	db *gorm.DB
}

func NewTransactionQueryRepo(db *gorm.DB) domain.ITransactionQueryRepo {
	return &transactionQueryRepo{
		db: db,
	}
}

func (r *transactionQueryRepo) GetTransactions(ctx context.Context) (transactions []*mysqlModel.Transaction, err error) {
	span, ctx := apm.StartSpan(ctx, "transactionQueryRepo.GetTransactions", "repo")
	defer span.End()

	result := r.db.WithContext(ctx).Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}

	return transactions, nil
}
