package user

import (
	"context"

	"banking/app/service/user"
	"banking/model"

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
	result := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (r *userQueryRepo) GetUsers(ctx context.Context) (users []*model.User, err error) {
	result := r.db.WithContext(ctx).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}