package transaction

import (
	"banking/model/mysql"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	FromUserID      uint                  `json:"fromUserId"`
	FromUserBalance decimal.Decimal       `json:"fromUserBalance"`
	ToUserID        uint                  `json:"toUserId"`
	ToUserBalance   decimal.Decimal       `json:"toUserBalance,omitempty"`
	Amount          decimal.Decimal       `json:"amount"`
	TransactionType mysql.TransactionType `json:"transactionType"`
	Details         string                `json:"details"`
}

type TransferResp struct {
	Data *Transaction `json:"data"`
}

type DepositResp struct {
	Data *Transaction `json:"data"`
}

type WithdrawResp struct {
	Data *Transaction `json:"data"`
}

type GetTransactionsResp struct {
	Data []*Transaction `json:"data"`
}
