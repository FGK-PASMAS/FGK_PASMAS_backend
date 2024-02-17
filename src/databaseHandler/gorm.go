package databasehandler

import (
	"github.com/MetaEMK/FGK_PASMAS_backend/router/realtime"
	"gorm.io/gorm"
)

var Db *gorm.DB

type DatabaseHandler struct {
    Db *gorm.DB
    rt *realtime.RealtimeHandler
}

func NewDatabaseHandler() (dh *DatabaseHandler) {
    dh = &DatabaseHandler{
        Db: Db.Begin(),
        rt: realtime.NewRealtimeHandler(),
    }

    return
}

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

func (dh *DatabaseHandler) CommitOrRollback(err error) {
    if dh.Db.Error == nil && err == nil {
        err := dh.Db.Commit().Error
        if err != nil {
            dh.Db.AddError(err)
            dh.Db.Rollback()
        } else {
            dh.rt.PublishEvents()
        }
    } else {
        dh.Db.Rollback()
    }
} 
