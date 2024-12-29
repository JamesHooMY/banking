package user

import (
	"context"
	"errors"

	"banking/domain"
	mysqlModel "banking/model/mysql"

	"go.elastic.co/apm/v2"
	"gorm.io/gorm"
)

type userCommandRepo struct {
	db *gorm.DB
}

func NewUserCommandRepo(db *gorm.DB) domain.IUserCommandRepo {
	return &userCommandRepo{
		db: db,
	}
}

func (r *userCommandRepo) CreateUser(ctx context.Context, user *mysqlModel.User) (err error) {
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
