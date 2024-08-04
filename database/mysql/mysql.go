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

var (
	masterDB *gorm.DB
	slaveDB  *gorm.DB
)

func InitMySQL(ctx context.Context) error {
	// Initialize master DB
	masterDSN := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.master.username"),
		viper.GetString("mysql.master.password"),
		viper.GetString("mysql.master.host"),
		viper.GetString("mysql.master.dbName"))

	if err := initDB(ctx, masterDSN, true); err != nil {
		return err
	}

	// Initialize slave DB
	slaveDSN := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.slave.username"),
		viper.GetString("mysql.slave.password"),
		viper.GetString("mysql.slave.host"),
		viper.GetString("mysql.slave.dbName"))

	if err := initDB(ctx, slaveDSN, false); err != nil {
		return err
	}

	// Auto migrate on master
	if err := masterDB.AutoMigrate(
		&model.User{},
		&model.Transaction{},
	); err != nil {
		return err
	}

	// Seed User data
	seedUsers(masterDB)

	return nil
}

func initDB(ctx context.Context, dsn string, isMaster bool) error {
	location, err := time.LoadLocation("UTC")
	if err != nil {
		return err
	}

	var db *gorm.DB
	err = retry(ctx, func() error {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
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
			return err
		}
		sqlDB, errDB := db.DB()
		if errDB != nil {
			return errDB
		}

		// Use the context for the ping operation
		ctxWT, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		return sqlDB.PingContext(ctxWT)
	}, 5, 5*time.Second)
	if err != nil {
		return err
	}

	if isMaster {
		masterDB = db
	} else {
		slaveDB = db
	}

	return nil
}

func retry(ctx context.Context, action func() error, attempts int, sleep time.Duration) error {
	for i := 0; i < attempts; i++ {
		err := action()
		if err == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(sleep):
		}
	}
	return fmt.Errorf("failed after %d attempts", attempts)
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

// GetMasterDB returns the master database connection
func GetMasterDB() *gorm.DB {
	return masterDB
}

// GetSlaveDB returns the slave database connection
func GetSlaveDB() *gorm.DB {
	return slaveDB
}
