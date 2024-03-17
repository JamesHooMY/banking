package user

import (
	"context"
	"errors"

	"banking/app/service/user"
	"banking/model"

	"github.com/shopspring/decimal"
	"go.elastic.co/apm/v2"
	"gorm.io/gorm"
)

type userCommandRepo struct {
	db *gorm.DB
}

func NewUserCommandRepo(db *gorm.DB) user.IUserCommandRepo {
	return &userCommandRepo{
		db: db,
	}
}

func (r *userCommandRepo) CreateUser(ctx context.Context, user *model.User) (err error) {
	span, ctx := apm.StartSpan(ctx, "userCommandRepo.CreateUser", "repo")
	defer span.End()

	result := r.db.WithContext(ctx).Create(user)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserExisted
		}
		return err
	}

	return nil
}

func (r *userCommandRepo) Transfer(ctx context.Context, fromUserID, toUserID uint, amount decimal.Decimal) (user *model.User, err error) {
	span, ctx := apm.StartSpan(ctx, "userCommandRepo.Transfer", "repo")
	defer span.End()

	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return nil, err
	}

	fromUser := &model.User{}
	// result := tx.Model(&model.User{}).Where("id = ?", fromUserID).First(fromUser)
	result := tx.Where("id = ?", fromUserID).Take(fromUser)
	if err := result.Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	toUser := &model.User{}
	// result = tx.Model(&model.User{}).Where("id = ?", toUserID).First(toUser)
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

	transaction := &model.Transaction{
		FromUserID:      fromUserID,
		ToUserID:        toUserID,
		Amount:          amount,
		FromUserBalance: fromUser.Balance,
		ToUserBalance:   toUser.Balance,
		TransactionType: model.Transfer,
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

func (r *userCommandRepo) Deposit(ctx context.Context, userID uint, amount decimal.Decimal) (user *model.User, err error) {
	span, ctx := apm.StartSpan(ctx, "userCommandRepo.Deposit", "repo")
	defer span.End()

	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return nil, err
	}

	user = &model.User{}
	result := tx.Model(&model.User{}).Where("id = ?", userID).First(user)
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

	transaction := &model.Transaction{
		FromUserID:      userID,
		ToUserID:        userID,
		Amount:          amount,
		FromUserBalance: user.Balance,
		ToUserBalance:   user.Balance,
		TransactionType: model.Deposit,
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

func (r *userCommandRepo) Withdraw(ctx context.Context, userID uint, amount decimal.Decimal) (user *model.User, err error) {
	span, ctx := apm.StartSpan(ctx, "userCommandRepo.Withdraw", "repo")
	defer span.End()

	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return nil, err
	}

	user = &model.User{}
	result := tx.Model(&model.User{}).Where("id = ?", userID).First(user)
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

	transaction := &model.Transaction{
		FromUserID:      userID,
		ToUserID:        userID,
		Amount:          amount,
		FromUserBalance: user.Balance,
		ToUserBalance:   user.Balance,
		TransactionType: model.Withdraw,
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
