package mysql

import (
	"github.com/lukedever/api"
	driver "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDB open mysql connection
func NewDB(dsn, mode string) (*gorm.DB, error) {
	db, err := gorm.Open(driver.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if mode == "debug" {
		db.Logger = logger.Default.LogMode(logger.Info)
		_ = db.AutoMigrate(&api.User{})
	}

	return db, err
}
