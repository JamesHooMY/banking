package mysql

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name    string          `gorm:"type:varchar(20);unique;index;not null" json:"name"`
	Balance decimal.Decimal `gorm:"type:decimal(10,2);unsigned;not null;default:'0'" json:"balance"`
}
