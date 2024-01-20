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

func ResetDatabase() error {
    transaction := Db.Begin()
    transaction.Exec("TRUNCATE TABLE divisions RESTART IDENTITY CASCADE")
    transaction.Exec("TRUNCATE TABLE passengers RESTART IDENTITY CASCADE")
    transaction.Commit()

    return transaction.Error
}
