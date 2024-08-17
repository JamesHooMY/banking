package transaction

import "banking/model/mysql"

type TransferResp struct {
	Data *mysql.User `json:"data"`
}

type DepositResp struct {
	Data *mysql.User `json:"data"`
}

type WithdrawResp struct {
	Data *mysql.User `json:"data"`
}

type GetTransactionsResp struct {
	Data []*mysql.Transaction `json:"data"`
}
