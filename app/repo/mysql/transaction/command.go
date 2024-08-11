package transaction

import (
	"context"

	domain "banking/domain"
	mysqlModel "banking/model/mysql"

	"github.com/shopspring/decimal"
	"go.elastic.co/apm/v2"
	"gorm.io/gorm"
)

type transactionCommandRepo struct {
	db *gorm.DB
}

func NewTransactionCommandRepo(db *gorm.DB) domain.ITransactionCommandRepo {
	return &transactionCommandRepo{
		db: db,
	}
}

func (r *transactionCommandRepo) Transfer(ctx context.Context, fromUserID, toUserID uint, amount decimal.Decimal) (user *mysqlModel.User, err error) {
	span, ctx := apm.StartSpan(ctx, "userCommandRepo.Transfer", "repo")
	defer span.End()

	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return nil, err
	}

	fromUser := &mysqlModel.User{}
	// result := tx.Model(&mysqlModel.User{}).Where("id = ?", fromUserID).First(fromUser)
	result := tx.Where("id = ?", fromUserID).Take(fromUser)
	if err := result.Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	toUser := &mysqlModel.User{}
	// result = tx.Model(&mysqlModel.User{}).Where("id = ?", toUserID).First(toUser)
	result = tx.Where("id = ?", toUserID).Take(toUser)
	if err := result.Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if fromUser.Balance.LessThan(amount) {
		tx.Rollback()
		return nil, ErrInsufficientBalance
	}

	fromUser.Balance = fromUser.Balance.Sub(amount)
	toUser.Balance = toUser.Balance.Add(amount)

	result = tx.Save(fromUser)
	if err := result.Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	result = tx.Save(toUser)
	if err := result.Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	transaction := &mysqlModel.Transaction{
		FromUserID:      fromUserID,
		ToUserID:        toUserID,
		Amount:          amount,
		FromUserBalance: fromUser.Balance,
		ToUserBalance:   toUser.Balance,
		TransactionType: mysqlModel.Transfer,
	}

	result = tx.Create(transaction)
	if err := result.Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return fromUser, nil
}

func (r *transactionCommandRepo) Deposit(ctx context.Context, userID uint, amount decimal.Decimal) (user *mysqlModel.User, err error) {
	span, ctx := apm.StartSpan(ctx, "userCommandRepo.Deposit", "repo")
	defer span.End()

	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return nil, err
	}

	user = &mysqlModel.User{}
	result := tx.Model(&mysqlModel.User{}).Where("id = ?", userID).First(user)
	if err := result.Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	user.Balance = user.Balance.Add(amount)

	result = tx.Save(user)
	if err := result.Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	transaction := &mysqlModel.Transaction{
		FromUserID:      userID,
		ToUserID:        userID,
		Amount:          amount,
		FromUserBalance: user.Balance,
		ToUserBalance:   user.Balance,
		TransactionType: mysqlModel.Deposit,
	}

	result = tx.Create(transaction)
	if err := result.Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return user, nil
}

func (r *transactionCommandRepo) Withdraw(ctx context.Context, userID uint, amount decimal.Decimal) (user *mysqlModel.User, err error) {
	span, ctx := apm.StartSpan(ctx, "userCommandRepo.Withdraw", "repo")
	defer span.End()

	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return nil, err
	}

	user = &mysqlModel.User{}
	result := tx.Model(&mysqlModel.User{}).Where("id = ?", userID).First(user)
	if err := result.Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if user.Balance.LessThan(amount) {
		tx.Rollback()
		return nil, ErrInsufficientBalance
	}

	user.Balance = user.Balance.Sub(amount)

	result = tx.Save(user)
	if err := result.Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	transaction := &mysqlModel.Transaction{
		FromUserID:      userID,
		ToUserID:        userID,
		Amount:          amount,
		FromUserBalance: user.Balance,
		ToUserBalance:   user.Balance,
		TransactionType: mysqlModel.Withdraw,
	}

	result = tx.Create(transaction)
	if err := result.Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return user, nil
}
