package transaction

import (
	"context"
	"time"

	domain "banking/domain"
	"banking/global"
	mysqlModel "banking/model/mysql"

	"github.com/shopspring/decimal"
	"go.elastic.co/apm/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type transactionCommandRepo struct {
	db *gorm.DB
}

func NewTransactionCommandRepo(db *gorm.DB) domain.ITransactionCommandRepo {
	return &transactionCommandRepo{
		db: db,
	}
}

// func (r *transactionCommandRepo) Transfer(ctx context.Context, fromUserID, toUserID uint, amount decimal.Decimal) (transaction *mysqlModel.Transaction, err error) {
// 	span, ctx := apm.StartSpan(ctx, "userCommandRepo.Transfer", "repo")
// 	defer span.End()

// 	tx := r.db.WithContext(ctx).Begin()
// 	if err := tx.Error; err != nil {
// 		return nil, err
// 	}

// 	fromUser := &mysqlModel.User{}
// 	// result := tx.Model(&mysqlModel.User{}).Where("id = ?", fromUserID).First(fromUser)
// 	result := tx.Where("id = ?", fromUserID).Take(fromUser)
// 	if err := result.Error; err != nil {
// 		tx.Rollback()
// 		return nil, err
// 	}

// 	toUser := &mysqlModel.User{}
// 	// result = tx.Model(&mysqlModel.User{}).Where("id = ?", toUserID).First(toUser)
// 	result = tx.Where("id = ?", toUserID).Take(toUser)
// 	if err := result.Error; err != nil {
// 		tx.Rollback()
// 		return nil, err
// 	}

// 	if fromUser.Balance.LessThan(amount) {
// 		tx.Rollback()
// 		return nil, ErrInsufficientBalance
// 	}

// 	fromUser.Balance = fromUser.Balance.Sub(amount)
// 	toUser.Balance = toUser.Balance.Add(amount)

// 	result = tx.Save(fromUser)
// 	if err := result.Error; err != nil {
// 		tx.Rollback()
// 		return nil, err
// 	}

// 	result = tx.Save(toUser)
// 	if err := result.Error; err != nil {
// 		tx.Rollback()
// 		return nil, err
// 	}

// 	transaction = &mysqlModel.Transaction{
// 		FromUserID:      fromUserID,
// 		ToUserID:        toUserID,
// 		Amount:          amount,
// 		FromUserBalance: fromUser.Balance,
// 		ToUserBalance:   toUser.Balance,
// 		TransactionType: mysqlModel.Transfer,
// 	}

// 	result = tx.Create(transaction)
// 	if err := result.Error; err != nil {
// 		tx.Rollback()
// 		return nil, err
// 	}

// 	if err := tx.Commit().Error; err != nil {
// 		tx.Rollback()
// 		return nil, err
// 	}

// 	return transaction, nil
// }

// clause lock
func (r *transactionCommandRepo) Transfer(ctx context.Context, fromUserID, toUserID uint, amount decimal.Decimal) (transaction *mysqlModel.Transaction, err error) {
	span, ctx := apm.StartSpan(ctx, "userCommandRepo.Transfer", "repo")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx := r.db.WithContext(ctx).Begin()
	if err = tx.Error; err != nil {
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()

	fromUser := &mysqlModel.User{}
	result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", fromUserID).Take(fromUser)
	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected == 0 {
		return nil, ErrUserNotFound
	} else if fromUser.Balance.LessThan(amount) {
		return nil, ErrInsufficientBalance
	}

	toUser := &mysqlModel.User{}
	result = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", toUserID).Take(toUser)
	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected == 0 {
		return nil, ErrUserNotFound
	}

	calculatedFromUser := fromUser
	calculatedToUser := toUser
	calculatedFromUser.Balance = fromUser.Balance.Sub(amount)
	calculatedToUser.Balance = toUser.Balance.Add(amount)

	// Update the fromUser balance
	result = tx.Save(fromUser)
	if err := result.Error; err != nil {
		return nil, err
	}

	// Update the toUser balance
	result = tx.Save(toUser)
	if err := result.Error; err != nil {
		return nil, err
	}

	transaction = &mysqlModel.Transaction{
		FromUserID:      fromUserID,
		ToUserID:        toUserID,
		Amount:          amount,
		FromUserBalance: calculatedFromUser.Balance,
		ToUserBalance:   calculatedToUser.Balance,
		TransactionType: mysqlModel.Transfer,
	}

	result = tx.Create(transaction)
	if err := result.Error; err != nil {
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return transaction, nil
}

func (r *transactionCommandRepo) Deposit(ctx context.Context, userID uint, amount decimal.Decimal) (transaction *mysqlModel.Transaction, err error) {
	span, ctx := apm.StartSpan(ctx, "userCommandRepo.Deposit", "repo")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx := r.db.WithContext(ctx).Begin()
	if err = tx.Error; err != nil {
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			global.Logger.Errorf("panic: %v", r)
			tx.Rollback()
		} else if err != nil {
			tx.Rollback()
		}
	}()

	// Lock the user row for update to prevent concurrent updates
	user := &mysqlModel.User{}
	result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", userID).Take(user)
	if result.Error != nil {
		return nil, result.Error
	}

	// Update the user balance
	calculatedBalance := user.Balance.Add(amount)
	result = tx.Model(user).Update("balance", calculatedBalance)
	if err = result.Error; err != nil {
		return nil, err
	}

	transaction = &mysqlModel.Transaction{
		FromUserID:      userID,
		ToUserID:        userID,
		Amount:          amount,
		FromUserBalance: calculatedBalance,
		ToUserBalance:   calculatedBalance,
		TransactionType: mysqlModel.Deposit,
	}

	if err := tx.Create(transaction).Error; err != nil {
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return transaction, nil
}

func (r *transactionCommandRepo) Withdraw(ctx context.Context, userID uint, amount decimal.Decimal) (transaction *mysqlModel.Transaction, err error) {
	span, ctx := apm.StartSpan(ctx, "userCommandRepo.Withdraw", "repo")
	defer span.End()

	tx := r.db.WithContext(ctx).Begin()
	if err = tx.Error; err != nil {
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			global.Logger.Errorf("panic: %v", r)
			tx.Rollback()
		} else if err != nil {
			tx.Rollback()
		}
	}()

	user := &mysqlModel.User{}
	result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", userID).Take(user)
	if err := result.Error; err != nil {
		return nil, err
	} else if result.RowsAffected == 0 {
		return nil, ErrInsufficientBalance
	} else if user.Balance.LessThan(amount) {
		return nil, ErrInsufficientBalance
	}

	// Update the user balance
	calculatedBalance := user.Balance.Sub(amount)
	result = tx.Model(user).Update("balance", calculatedBalance)
	if err := result.Error; err != nil {
		return nil, err
	}

	transaction = &mysqlModel.Transaction{
		FromUserID:      userID,
		ToUserID:        userID,
		Amount:          amount,
		FromUserBalance: calculatedBalance,
		ToUserBalance:   calculatedBalance,
		TransactionType: mysqlModel.Withdraw,
	}

	if err := tx.Create(transaction).Error; err != nil {
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return transaction, nil
}
