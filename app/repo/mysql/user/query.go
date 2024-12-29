package user

import (
	"context"
	"time"

	domain "banking/domain"
	mysqlModel "banking/model/mysql"

	"go.elastic.co/apm/v2"
	"gorm.io/gorm"
)

type userQueryRepo struct {
	db *gorm.DB
}

func NewUserQueryRepo(db *gorm.DB) domain.IUserQueryRepo {
	return &userQueryRepo{
		db: db,
	}
}

func (r *userQueryRepo) GetUsers(ctx context.Context, userID uint) (users []*mysqlModel.User, err error) {
	span, ctx := apm.StartSpan(ctx, "userQueryRepo.GetUsers", "repo")
	defer span.End()

	if userID != 0 {
		user := &mysqlModel.User{}
		result := r.db.WithContext(ctx).Where("id = ?", userID).Take(&user)
		if result.Error != nil {
			return nil, result.Error
		}

		users = append(users, user)
		return users, nil
	}

	result := r.db.WithContext(ctx).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	// Simulate slow query for demo
	time.Sleep(2 * time.Millisecond)

	return users, nil
}

func (r *userQueryRepo) GetUserByEmail(ctx context.Context, email string) (user *mysqlModel.User, err error) {
	span, ctx := apm.StartSpan(ctx, "userQueryRepo.GetUserByEmail", "repo")
	defer span.End()

	result := r.db.WithContext(ctx).Where("email = ?", email).Take(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}
