package user

import (
	"context"
	"time"

	"banking/app/service/user"
	"banking/model"

	"go.elastic.co/apm/v2"
	"gorm.io/gorm"
)

type userQueryRepo struct {
	db *gorm.DB
}

func NewUserQueryRepo(db *gorm.DB) user.IUserQueryRepo {
	return &userQueryRepo{
		db: db,
	}
}

func (r *userQueryRepo) GetUser(ctx context.Context, userID uint) (user *model.User, err error) {
	span, ctx := apm.StartSpan(ctx, "userQueryRepo.GetUser", "repo")
	defer span.End()

	// result := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).First(&user)
	result := r.db.WithContext(ctx).Where("id = ?", userID).Take(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	time.Sleep(2 * time.Millisecond)

	return user, nil
}

func (r *userQueryRepo) GetUsers(ctx context.Context) (users []*model.User, err error) {
	span, ctx := apm.StartSpan(ctx, "userQueryRepo.GetUsers", "repo")
	defer span.End()

	result := r.db.WithContext(ctx).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	// Simulate slow query for demo
	time.Sleep(2 * time.Millisecond)

	return users, nil
}
