package user_test

import (
	"context"
	"testing"

	userRepo "banking/app/repo/mysql/user"
	userModel "banking/model/mysql"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func Test_GetUsers(t *testing.T) {
	if err := mysqlTestDB.Migrator().DropTable(&userModel.User{}); err != nil {
		t.Fatal(err)
	}
	if err := mysqlTestDB.AutoMigrate(&userModel.User{}); err != nil {
		t.Fatal(err)
	}

	user1 := &userModel.User{
		Model:    gorm.Model{ID: 1},
		Name:     "user1",
		Email:    "user1@yopmail",
		Password: "password",
		Balance:  decimal.NewFromFloat(100),
	}

	user2 := &userModel.User{
		Model:    gorm.Model{ID: 2},
		Name:     "user2",
		Email:    "user2@yopmail",
		Password: "password",
		Balance:  decimal.NewFromFloat(200),
	}

	if err := mysqlTestDB.Create(user1).Error; err != nil {
		t.Fatal("Error creating user1:", err)
	}
	if err := mysqlTestDB.Create(user2).Error; err != nil {
		t.Fatal("Error creating user2:", err)
	}

	userQueryRepo := userRepo.NewUserQueryRepo(mysqlTestDB)
	users, err := userQueryRepo.GetUsers(context.Background(), 0)

	assert.Nil(t, err)
	assert.NotNil(t, users)
	assert.Equal(t, 2, len(users))
	assert.Equal(t, user1.Name, users[0].Name)
	assert.Equal(t, user2.Name, users[1].Name)
}
