package mysql

import (
	driver "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewDB open mysql connection
func NewDB(dsn, mode string) (*gorm.DB, error) {
	db, err := gorm.Open(driver.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if mode == "debug" {
		db.Debug()
	}

	return db, err
}
