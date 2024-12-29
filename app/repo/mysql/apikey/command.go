package apikey

import (
	"context"

	"banking/domain"
	mysqlModel "banking/model/mysql"

	"go.elastic.co/apm/v2"
	"gorm.io/gorm"
)

type apikeyCommandRepo struct {
	db *gorm.DB
}

func NewAPIKeyCommandRepo(db *gorm.DB) domain.IAPIKeyCommandRepo {
	return &apikeyCommandRepo{db: db}
}

func (r *apikeyCommandRepo) CreateAPIKey(ctx context.Context, userID uint, key string, secret string) error {
	span, ctx := apm.StartSpan(ctx, "apikeyQueryRepo.GetAPIKey", "repo")
	defer span.End()

	apiKey := &mysqlModel.APIKey{
		UserID: userID,
		APIKey:    key,
		Secret: secret,
	}
	result := r.db.WithContext(ctx).Create(apiKey)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *apikeyCommandRepo) DeleteAPIKey(ctx context.Context, userID uint, key string) error {
	span, ctx := apm.StartSpan(ctx, "apikeyQueryRepo.GetAPIKey", "repo")
	defer span.End()

	result := r.db.WithContext(ctx).Where("userId = ? AND key = ?", userID, key).Delete(&mysqlModel.APIKey{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
