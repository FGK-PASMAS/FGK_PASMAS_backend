package databasehandler

import (
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitGorm(dbConn *gorm.DB) *gorm.DB {
    Db = dbConn

    initDivision()
    initPassenger()

    return Db
}
