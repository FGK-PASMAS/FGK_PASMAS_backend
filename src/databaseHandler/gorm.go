package databasehandler

import (
	"github.com/MetaEMK/FGK_PASMAS_backend/router/realtime"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitGorm(dbConn *gorm.DB) *gorm.DB {
    Db = dbConn

    initDivision()

    initPlane(nil)
    initPilot(nil)

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
    transaction.Exec("TRUNCATE TABLE pilots RESTART IDENTITY CASCADE")
    transaction.Commit()

    seed := Db.Begin()
    SeedDivision()
    SeedPlane(nil)
    SeedPilot(nil)
    seed.Commit()

    return transaction.Error
}

func CommitOrRollback(db *gorm.DB, rt *realtime.RealtimeHandler) {
    if db.Error == nil {
        err := db.Commit().Error
        if err != nil {
            db.AddError(err)
            db.Rollback()
        } else {
            rt.PublishEvents()
        }
    } else {
        db.Rollback()
    }
} 
