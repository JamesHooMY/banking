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

func Test_CreateUser(t *testing.T) {
	if err := mysqlDB.Migrator().DropTable(&model.User{}); err != nil {
		t.Fatal(err)
	}
	if err := mysqlDB.AutoMigrate(&model.User{}); err != nil {
		t.Fatal(err)
	}

	userCommandRepo := user.NewUserCommandRepo(mysqlDB)
	err := userCommandRepo.CreateUser(context.Background(), &model.User{
		Name:    "user1",
		Balance: decimal.NewFromFloat(100),
	})

	assert.Nil(t, err)

	user := &model.User{}
	mysqlDB.First(user, "name = ?", "user1")

	assert.NotNil(t, user)
	assert.Equal(t, "user1", user.Name)
	assert.True(t, user.Balance.Equal(decimal.NewFromFloat(100)))
}

func Test_Transfer(t *testing.T) {
	if err := mysqlDB.Migrator().DropTable(
		&model.User{},
		&model.Transaction{},
	); err != nil {
		t.Fatal(err)
	}
	if err := mysqlDB.AutoMigrate(
		&model.User{},
		&model.Transaction{},
	); err != nil {
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

	userCommandRepo := user.NewUserCommandRepo(mysqlDB)
	user, err := userCommandRepo.Transfer(context.Background(), 1, 2, decimal.NewFromFloat(50))

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "user1", user.Name)
	assert.True(t, user.Balance.Equal(decimal.NewFromFloat(50)))

	user2 := &model.User{}
	mysqlDB.First(user2, "name = ?", "user2")

	assert.NotNil(t, user2)
	assert.Equal(t, "user2", user2.Name)
	assert.True(t, user2.Balance.Equal(decimal.NewFromFloat(250)))
}

func Test_Deposit(t *testing.T) {
	if err := mysqlDB.Migrator().DropTable(
		&model.User{},
		&model.Transaction{},
	); err != nil {
		t.Fatal(err)
	}
	if err := mysqlDB.AutoMigrate(
		&model.User{},
		&model.Transaction{},
	); err != nil {
		t.Fatal(err)
	}
	mysqlDB.Create(&model.User{
		Model:   gorm.Model{ID: 1},
		Name:    "user1",
		Balance: decimal.NewFromFloat(100),
	})

	userCommandRepo := user.NewUserCommandRepo(mysqlDB)
	user, err := userCommandRepo.Deposit(context.Background(), 1, decimal.NewFromFloat(50))

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "user1", user.Name)
	assert.True(t, user.Balance.Equal(decimal.NewFromFloat(150)))
}

func Test_Withdraw(t *testing.T) {
	if err := mysqlDB.Migrator().DropTable(
		&model.User{},
		&model.Transaction{},
	); err != nil {
		t.Fatal(err)
	}
	if err := mysqlDB.AutoMigrate(
		&model.User{},
		&model.Transaction{},
	); err != nil {
		t.Fatal(err)
	}
	mysqlDB.Create(&model.User{
		Model:   gorm.Model{ID: 1},
		Name:    "user1",
		Balance: decimal.NewFromFloat(100),
	})

	userCommandRepo := user.NewUserCommandRepo(mysqlDB)
	user, err := userCommandRepo.Withdraw(context.Background(), 1, decimal.NewFromFloat(50))

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "user1", user.Name)
	assert.True(t, user.Balance.Equal(decimal.NewFromFloat(50)))
}
