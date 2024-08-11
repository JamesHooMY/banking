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

	userCommandRepo := userRepo.NewUserCommandRepo(mysqlTestDB)
	err := userCommandRepo.CreateUser(context.Background(), &mysqlModel.User{
		Name:    "user1",
		Balance: decimal.NewFromFloat(100),
	})

	assert.Nil(t, err)

	user := &mysqlModel.User{}
	mysqlTestDB.First(user, "name = ?", "user1")

	assert.NotNil(t, user)
	assert.Equal(t, "user1", user.Name)
	assert.True(t, user.Balance.Equal(decimal.NewFromFloat(100)))
}
