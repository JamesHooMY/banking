package transaction_test

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

	user1 := &mysqlModel.User{
		Model:   gorm.Model{ID: 1},
		Name:    "user1",
		Email:   "user1@yopmail",
		Balance: decimal.NewFromFloat(100),
	}

	user2 := &mysqlModel.User{
		Model:   gorm.Model{ID: 2},
		Name:    "user2",
		Email:   "user2@yopmail",
		Balance: decimal.NewFromFloat(200),
	}

	if err := mysqlTestDB.Create(user1).Error; err != nil {
		t.Fatal(err)
	}
	if err := mysqlTestDB.Create(user2).Error; err != nil {
		t.Fatal(err)
	}

	transactionCommandRepo := transactionRepo.NewTransactionCommandRepo(mysqlTestDB)
	transaction, err := transactionCommandRepo.Transfer(context.Background(), user1.Model.ID, user2.Model.ID, decimal.NewFromFloat(50))

	assert.Nil(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, transaction.FromUserID, user1.Model.ID)
	assert.Equal(t, transaction.ToUserID, user2.Model.ID)
	assert.True(t, transaction.Amount.Equal(decimal.NewFromFloat(50)))
	assert.True(t, transaction.FromUserBalance.Equal(user1.Balance.Sub(decimal.NewFromFloat(50))))
	assert.True(t, transaction.ToUserBalance.Equal(user2.Balance.Add(decimal.NewFromFloat(50))))
	assert.Equal(t, transaction.TransactionType, mysqlModel.Transfer)
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

	user1 := &mysqlModel.User{
		Model:   gorm.Model{ID: 1},
		Name:    "user1",
		Balance: decimal.NewFromFloat(100),
	}

	result := mysqlTestDB.Create(user1)
	if result.Error != nil {
		t.Fatal(result.Error)
	}

	transactionCommandRepo := transactionRepo.NewTransactionCommandRepo(mysqlTestDB)
	transaction, err := transactionCommandRepo.Deposit(context.Background(), 1, decimal.NewFromFloat(50))

	assert.Nil(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, user1.Model.ID, transaction.FromUserID)
	assert.Equal(t, user1.Model.ID, transaction.ToUserID)
	assert.True(t, user1.Balance.Add(decimal.NewFromFloat(50)).Equal(transaction.FromUserBalance))
	assert.True(t, user1.Balance.Add(decimal.NewFromFloat(50)).Equal(transaction.ToUserBalance))
	assert.True(t, decimal.NewFromFloat(50).Equal(transaction.Amount))
	assert.Equal(t, mysqlModel.Deposit, transaction.TransactionType)
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

	user1 := &mysqlModel.User{
		Model:   gorm.Model{ID: 1},
		Name:    "user1",
		Balance: decimal.NewFromFloat(100),
	}

	if result := mysqlTestDB.Create(user1); result.Error != nil {
		t.Fatal(result.Error)
	}

	transactionCommandRepo := transactionRepo.NewTransactionCommandRepo(mysqlTestDB)
	transaction, err := transactionCommandRepo.Withdraw(context.Background(), 1, decimal.NewFromFloat(50))

	assert.Nil(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, user1.Model.ID, transaction.FromUserID)
	assert.Equal(t, user1.Model.ID, transaction.ToUserID)
	assert.True(t, user1.Balance.Sub(decimal.NewFromFloat(50)).Equal(transaction.FromUserBalance))
	assert.True(t, user1.Balance.Sub(decimal.NewFromFloat(50)).Equal(transaction.ToUserBalance))
	assert.True(t, decimal.NewFromFloat(50).Equal(transaction.Amount))
	assert.Equal(t, mysqlModel.Withdraw, transaction.TransactionType)
}
