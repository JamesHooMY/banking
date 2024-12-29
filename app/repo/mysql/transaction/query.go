package transaction

import (
	"context"

	"banking/domain"
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

func (r *transactionQueryRepo) GetTransactions(ctx context.Context, userID uint) (transactions []*mysqlModel.Transaction, err error) {
	span, ctx := apm.StartSpan(ctx, "transactionQueryRepo.GetTransactions", "repo")
	defer span.End()

	result := r.db.WithContext(ctx).Where("from_user_id = ?", userID).Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}

	return transactions, nil
}
