package mysql

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string          `gorm:"type:varchar(20);not null" json:"name"`
	Email    string          `gorm:"type:varchar(100);unique;index;not null" json:"email"`
	Password string          `gorm:"type:varchar(255);not null" json:"password"`
	Balance  decimal.Decimal `gorm:"type:decimal(10,2);unsigned;not null;default:'0'" json:"balance"`
	IsAdmin  bool            `gorm:"type:tinyint(1);default:false" json:"isAdmin"`
}
