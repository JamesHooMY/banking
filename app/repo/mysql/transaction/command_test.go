package user_test

import (
	"context"
	"testing"

	transactionRepo "banking/app/repo/mysql/transaction"
	mysqlModel "banking/model/mysql"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func Test_Transfer(t *testing.T) {
	if err := mysqlTestDB.Migrator().DropTable(
		&mysqlModel.User{},
		&mysqlModel.Transaction{},
	); err != nil {
		t.Fatal(err)
	}
	if err := mysqlTestDB.AutoMigrate(
		&mysqlModel.User{},
		&mysqlModel.Transaction{},
	); err != nil {
		t.Fatal(err)
	}
	mysqlTestDB.CreateInBatches([]*mysqlModel.User{
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

	transactionCommandRepo := transactionRepo.NewTransactionCommandRepo(mysqlTestDB)
	user, err := transactionCommandRepo.Transfer(context.Background(), 1, 2, decimal.NewFromFloat(50))

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "user1", user.Name)
	assert.True(t, user.Balance.Equal(decimal.NewFromFloat(50)))

	user2 := &mysqlModel.User{}
	mysqlTestDB.First(user2, "name = ?", "user2")

	assert.NotNil(t, user2)
	assert.Equal(t, "user2", user2.Name)
	assert.True(t, user2.Balance.Equal(decimal.NewFromFloat(250)))
}

func Test_Deposit(t *testing.T) {
	if err := mysqlTestDB.Migrator().DropTable(
		&mysqlModel.User{},
		&mysqlModel.Transaction{},
	); err != nil {
		t.Fatal(err)
	}
	if err := mysqlTestDB.AutoMigrate(
		&mysqlModel.User{},
		&mysqlModel.Transaction{},
	); err != nil {
		t.Fatal(err)
	}
	mysqlTestDB.Create(&mysqlModel.User{
		Model:   gorm.Model{ID: 1},
		Name:    "user1",
		Balance: decimal.NewFromFloat(100),
	})

	transactionCommandRepo := transactionRepo.NewTransactionCommandRepo(mysqlTestDB)
	user, err := transactionCommandRepo.Deposit(context.Background(), 1, decimal.NewFromFloat(50))

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "user1", user.Name)
	assert.True(t, user.Balance.Equal(decimal.NewFromFloat(150)))
}

func Test_Withdraw(t *testing.T) {
	if err := mysqlTestDB.Migrator().DropTable(
		&mysqlModel.User{},
		&mysqlModel.Transaction{},
	); err != nil {
		t.Fatal(err)
	}
	if err := mysqlTestDB.AutoMigrate(
		&mysqlModel.User{},
		&mysqlModel.Transaction{},
	); err != nil {
		t.Fatal(err)
	}
	mysqlTestDB.Create(&mysqlModel.User{
		Model:   gorm.Model{ID: 1},
		Name:    "user1",
		Balance: decimal.NewFromFloat(100),
	})

	transactionCommandRepo := transactionRepo.NewTransactionCommandRepo(mysqlTestDB)
	user, err := transactionCommandRepo.Withdraw(context.Background(), 1, decimal.NewFromFloat(50))

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "user1", user.Name)
	assert.True(t, user.Balance.Equal(decimal.NewFromFloat(50)))
}
