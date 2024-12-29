package apikey

import (
	"context"

	"banking/domain"
	mysqlModel "banking/model/mysql"

	"go.elastic.co/apm/v2"
	"gorm.io/gorm"
)

type apikeyQueryRepo struct {
	db *gorm.DB
}

func NewAPIKeyQueryRepo(db *gorm.DB) domain.IAPIKeyQueryRepo {
	return &apikeyQueryRepo{db: db}
}

// func (r *apikeyQueryRepo) GetAPIKey(ctx context.Context, userID uint, key string) (*mysqlModel.APIKey, error) {
// 	span, ctx := apm.StartSpan(ctx, "apikeyQueryRepo.GetAPIKey", "repo")
// 	defer span.End()

// 	apiKey := &mysqlModel.APIKey{}
// 	result := r.db.WithContext(ctx).Where("userId = ? AND key = ?", userID, key).First(apiKey)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}

// 	return apiKey, nil
// }

func (r *apikeyQueryRepo) GetAPIKeys(ctx context.Context, userID uint, key string) ([]*mysqlModel.APIKey, error) {
	span, ctx := apm.StartSpan(ctx, "apikeyQueryRepo.GetAPIKeys", "repo")
	defer span.End()

	if key != "" {
		apiKey := &mysqlModel.APIKey{}
		result := r.db.WithContext(ctx).Where("api_key = ?", key).First(apiKey)
		if result.Error != nil {
			return nil, result.Error
		}

		return []*mysqlModel.APIKey{apiKey}, nil
	}

	if userID != 0 {
		var apiKeys []*mysqlModel.APIKey
		result := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&apiKeys)
		if result.Error != nil {
			return nil, result.Error
		}

		return apiKeys, nil
	}

	var apiKeys []*mysqlModel.APIKey
	result := r.db.WithContext(ctx).Find(&apiKeys)
	if result.Error != nil {
		return nil, result.Error
	}

	return apiKeys, nil
}
