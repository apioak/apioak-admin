package packages

import (
	"gorm.io/gorm"
)

var (
	dbConnection *gorm.DB
)

func SetDb(db *gorm.DB) {
	dbConnection = db
}

func GetDb() *gorm.DB {
	return dbConnection
}
