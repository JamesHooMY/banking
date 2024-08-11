package mysql

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type TransactionType string

const (
	Deposit  TransactionType = "deposit"
	Withdraw TransactionType = "withdraw"
	Transfer TransactionType = "transfer"
)

type Transaction struct {
	gorm.Model
	FromUser        User            `gorm:"foreignKey:FromUserID" json:"-"`
	FromUserID      uint            `gorm:"type:int;unsigned;index;not null" json:"fromUserId"`
	FromUserBalance decimal.Decimal `gorm:"type:decimal(10,2);unsigned;not null" json:"fromUserBalance"`
	ToUser          User            `gorm:"foreignKey:ToUserID" json:"-"`
	ToUserID        uint            `gorm:"type:int;unsigned;index;not null" json:"toUserId"`
	ToUserBalance   decimal.Decimal `gorm:"type:decimal(10,2);unsigned;not null" json:"toUserBalance"`
	Amount          decimal.Decimal `gorm:"type:decimal(10,2);unsigned;not null" json:"amount"`
	TransactionType TransactionType `gorm:"type:enum('deposit','withdraw','transfer');not null" json:"transactionType"`
	Details         string          `gorm:"type:text" json:"details"`
}
