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

func Test_GetTransactions(t *testing.T) {
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

	user := &mysqlModel.User{
		Model:    gorm.Model{ID: 1},
		Name:     "user1",
		Email:    "user1@yopmail",
		Password: "password",
		Balance:  decimal.NewFromFloat(100),
	}

	if err := mysqlTestDB.Create(user).Error; err != nil {
		t.Fatal(err)
	}

	expectedTransaction := &mysqlModel.Transaction{
		FromUserID:      user.Model.ID,
		ToUserID:        user.Model.ID,
		Amount:          decimal.NewFromFloat(100),
		FromUserBalance: user.Balance.Add(decimal.NewFromFloat(100)),
		ToUserBalance:   user.Balance.Add(decimal.NewFromFloat(100)),
		TransactionType: mysqlModel.Deposit,
	}

	if err := mysqlTestDB.Create(expectedTransaction).Error; err != nil {
		t.Fatal(err)
	}

	transactionQueryRepo := transactionRepo.NewTransactionQueryRepo(mysqlTestDB)
	transactions, err := transactionQueryRepo.GetTransactions(context.Background(), user.Model.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, transactions)
	assert.Equal(t, 1, len(transactions))
	assert.Equal(t, expectedTransaction.FromUserID, transactions[0].FromUserID)
	assert.Equal(t, expectedTransaction.ToUserID, transactions[0].ToUserID)
	assert.True(t, expectedTransaction.Amount.Equal(transactions[0].Amount))
	assert.True(t, expectedTransaction.FromUserBalance.Equal(transactions[0].FromUserBalance))
	assert.True(t, expectedTransaction.ToUserBalance.Equal(transactions[0].ToUserBalance))
	assert.Equal(t, expectedTransaction.TransactionType, transactions[0].TransactionType)
}
