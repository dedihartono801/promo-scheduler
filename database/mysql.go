package database

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMysql() (*gorm.DB, error) {
	database, err := gorm.Open(mysql.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return database, nil
}
