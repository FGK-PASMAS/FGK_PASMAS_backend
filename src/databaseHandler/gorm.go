package databasehandler

import (
	"runtime"

	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	"github.com/MetaEMK/FGK_PASMAS_backend/logging"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/realtime"
	"gorm.io/gorm"
)

var Db *gorm.DB

var log = logging.DbLogger

type DatabaseHandler struct {
    Db              *gorm.DB
    rt              *realtime.RealtimeHandler
    isClosed      bool
}

func InitGorm(dbConn *gorm.DB) *gorm.DB {
    Db = dbConn

    initUser()

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

func NewDatabaseHandler() (dh *DatabaseHandler) {
    dh = &DatabaseHandler{
        Db: Db.Begin(),
        rt: realtime.NewRealtimeHandler(),
    }

    runtime.SetFinalizer(dh, finalize)
    return
}

func (dh *DatabaseHandler) CommitOrRollback(err error) error {
    if dh.isClosed {
        return nil
    }

    dh.isClosed = true
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
    return dh.Db.Error
} 

func finalize(dh *DatabaseHandler) {
    if dh.isClosed == false {
        log.Error(cerror.ErrDatabaseHandlerDestroy.Error())
        dh.CommitOrRollback(cerror.ErrDatabaseHandlerDestroy)
    }
}
