package mysql

import (
	"context"
	"fmt"
	"time"

	"banking/model"

	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	mysql "go.elastic.co/apm/module/apmgormv2/v2/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func InitMySQL(ctx context.Context) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.username"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetString("mysql.dbName"))

	location, err := time.LoadLocation("UTC")
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			TablePrefix:   viper.GetString("mysql.tablePrefix"),
		},
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().In(location)
		},
	})
	if err != nil {
		return nil, err
	}
	db = db.WithContext(ctx)

	// Set up connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(viper.GetInt("mysql.maxIdleConns"))
	sqlDB.SetMaxOpenConns(viper.GetInt("mysql.maxOpenConns"))
	sqlDB.SetConnMaxLifetime(time.Duration(viper.GetInt("mysql.maxLifetime")) * time.Hour)

	// Auto migrate
	if err := db.AutoMigrate(
		&model.User{},
		&model.Transaction{},
	); err != nil {
		return nil, err
	}

	// Seed User data
	seedUsers(db)

	return db, nil
}

// Function to seed User data
func seedUsers(db *gorm.DB) {
	users := []*model.User{
		{
			Name:    "User1",
			Balance: decimal.NewFromFloat(100.00),
		},
		{
			Name:    "User2",
			Balance: decimal.NewFromFloat(200.00),
		},
		{
			Name:    "User3",
			Balance: decimal.NewFromFloat(300.00),
		},
	}
	for _, user := range users {
		db.Create(&user)
	}
}
