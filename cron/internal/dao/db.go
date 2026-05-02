package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMySQLDB(dsn string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
