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
	if err := mysqlTestDB.Migrator().DropTable(&mysqlModel.Transaction{}); err != nil {
		t.Fatal(err)
	}

	if err := mysqlTestDB.AutoMigrate(&mysqlModel.Transaction{}); err != nil {
		t.Fatal(err)
	}

	// transaction := &mysqlModel.Transaction{
	// 	FromUserID:      userID,
	// 	ToUserID:        userID,
	// 	Amount:          amount,
	// 	FromUserBalance: user.Balance,
	// 	ToUserBalance:   user.Balance,
	// 	TransactionType: mysqlModel.Deposit,
	// }

	mysqlTestDB.Create(&mysqlModel.Transaction{
		FromUserID:      1,
		ToUserID:        1,
		Amount:          decimal.NewFromInt(100),
		FromUserBalance: decimal.NewFromInt(100),
		ToUserBalance:   decimal.NewFromInt(100),
		TransactionType: mysqlModel.Deposit,
	})

	transactionQueryRepo := transactionRepo.NewTransactionQueryRepo(mysqlTestDB)
	transactions, err := transactionQueryRepo.GetTransactions(context.Background())

	assert.Nil(t, err)
	assert.NotNil(t, transactions)
	assert.Equal(t, 1, len(transactions))
	assert.Equal(t, decimal.NewFromInt(100), transactions[0].Amount)
	assert.Equal(t, decimal.NewFromInt(100), transactions[0].FromUserBalance)
	assert.Equal(t, decimal.NewFromInt(100), transactions[0].ToUserBalance)
	assert.Equal(t, mysqlModel.Deposit, transactions[0].TransactionType)
}
