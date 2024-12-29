package user_test

import (
	"context"
	"testing"

	userRepo "banking/app/repo/mysql/user"
	mysqlModel "banking/model/mysql"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_CreateUser(t *testing.T) {
	if err := mysqlTestDB.Migrator().DropTable(&mysqlModel.User{}); err != nil {
		t.Fatal(err)
	}
	if err := mysqlTestDB.AutoMigrate(&mysqlModel.User{}); err != nil {
		t.Fatal(err)
	}

	user1 := &mysqlModel.User{
		Name:     "user1",
		Email:    "user1@yopmail",
		Password: "password",
		Balance:  decimal.NewFromFloat(100),
	}

	userCommandRepo := userRepo.NewUserCommandRepo(mysqlTestDB)
	err := userCommandRepo.CreateUser(context.Background(), &mysqlModel.User{
		Name:     user1.Name,
		Email:    user1.Email,
		Password: user1.Password,
		Balance:  user1.Balance,
	})

	assert.Nil(t, err)
}
