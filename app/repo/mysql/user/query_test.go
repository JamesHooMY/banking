package user_test

import (
	"context"
	"testing"

	"banking/app/repo/mysql/user"
	"banking/model"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func Test_GetUser(t *testing.T) {
	if err := mysqlDB.Migrator().DropTable(&model.User{}); err != nil {
		t.Fatal(err)
	}
	if err := mysqlDB.AutoMigrate(&model.User{}); err != nil {
		t.Fatal(err)
	}
	mysqlDB.Create(&model.User{
		Model:   gorm.Model{ID: 1},
		Name:    "user1",
		Balance: decimal.NewFromFloat(100),
	})

	userQueryRepo := user.NewUserQueryRepo(mysqlDB)
	user, err := userQueryRepo.GetUser(context.Background(), 1)

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "user1", user.Name)
}

func Test_GetUsers(t *testing.T) {
	if err := mysqlDB.Migrator().DropTable(&model.User{}); err != nil {
		t.Fatal(err)
	}
	if err := mysqlDB.AutoMigrate(&model.User{}); err != nil {
		t.Fatal(err)
	}
	mysqlDB.CreateInBatches([]*model.User{
		{
			Model:   gorm.Model{ID: 1},
			Name:    "user1",
			Balance: decimal.NewFromFloat(100),
		},
		{
			Model:   gorm.Model{ID: 2},
			Name:    "user2",
			Balance: decimal.NewFromFloat(200),
		},
	}, 2)

	userQueryRepo := user.NewUserQueryRepo(mysqlDB)
	users, err := userQueryRepo.GetUsers(context.Background())

	assert.Nil(t, err)
	assert.NotNil(t, users)
	assert.Equal(t, 2, len(users))
	assert.Equal(t, "user1", users[0].Name)
	assert.Equal(t, "user2", users[1].Name)
}
