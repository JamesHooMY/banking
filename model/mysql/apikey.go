package mysql

import "gorm.io/gorm"

type APIKey struct {
	gorm.Model
	UserID uint   `gorm:"not null" json:"userId"`
	APIKey string `gorm:"type:varchar(255);unique;not null" json:"key"`
	Secret string `gorm:"type:varchar(255);not null" json:"secret"`
	User   User   `gorm:"foreignKey:UserID;" json:"-"` // Foreign key to User
}
