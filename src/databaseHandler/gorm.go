package databasehandler

import (
	//"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitGorm(dbConn *gorm.DB) *gorm.DB {
    Db = dbConn


    initDivision()
    initPlane()
    initPilot()
    initFlight()
    initPassenger()

    return Db
}

func ResetDatabase() error {
    transaction := Db.Begin()
    transaction.Exec("TRUNCATE TABLE divisions RESTART IDENTITY CASCADE")
    transaction.Exec("TRUNCATE TABLE passengers RESTART IDENTITY CASCADE")
    transaction.Exec("TRUNCATE TABLE flights RESTART IDENTITY CASCADE")
    transaction.Exec("TRUNCATE TABLE planes RESTART IDENTITY CASCADE")
    transaction.Commit()

    SeedDivision()

    return transaction.Error
}
