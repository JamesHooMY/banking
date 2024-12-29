package transaction_test

import (
	"context"
	"testing"

	transactionRepo "banking/app/repo/mysql/transaction"
	mysqlModel "banking/model/mysql"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
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
		Name:    "user1",
		Balance: decimal.NewFromFloat(100.00).Round(2),
	}
	mysqlTestDB.Create(user)

	expectedTransaction := &mysqlModel.Transaction{
		FromUserID:      user.ID,
		ToUserID:        user.ID,
		Amount:          decimal.NewFromFloat(100.00).Round(2),
		FromUserBalance: user.Balance,
		ToUserBalance:   user.Balance.Add(decimal.NewFromFloat(100.00)).Round(2),
		TransactionType: mysqlModel.Deposit,
	}
	mysqlTestDB.Create(expectedTransaction)

	transactionQueryRepo := transactionRepo.NewTransactionQueryRepo(mysqlTestDB)
	transactions, err := transactionQueryRepo.GetTransactions(context.Background(), user.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, transactions)
	assert.Equal(t, 1, len(transactions))
	assert.Equal(t, expectedTransaction.FromUserID, transactions[0].FromUserID)
	assert.Equal(t, expectedTransaction.ToUserID, transactions[0].ToUserID)
	assert.Equal(t, expectedTransaction.Amount, transactions[0].Amount)
	assert.Equal(t, expectedTransaction.FromUserBalance, transactions[0].FromUserBalance)
	assert.Equal(t, expectedTransaction.ToUserBalance, transactions[0].ToUserBalance)
	assert.Equal(t, expectedTransaction.TransactionType, transactions[0].TransactionType)
}
