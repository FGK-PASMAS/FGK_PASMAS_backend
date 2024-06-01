package databasehandler

import (
	"runtime"

	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	"github.com/MetaEMK/FGK_PASMAS_backend/config"
	"github.com/MetaEMK/FGK_PASMAS_backend/logging"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/realtime"
	"gorm.io/gorm"
)

var Db *gorm.DB

var log = logging.NewLogger("DATABASE", config.GetGlobalLogLevel())

type DatabaseHandler struct {
	Db       *gorm.DB
	rt       *realtime.RealtimeHandler
	isClosed bool
}

func InitGorm(dbConn *gorm.DB) *gorm.DB {
    log.Info("Initalizing dtabasehandler")
	Db = dbConn

	initUser()

	initDivision()

	initPlane(nil)
	initPilot(nil)

	initFlight()
	initPassenger()

    log.Debug("Database init finished successfully")
	return Db
}

func NewDatabaseHandler(user model.UserJwtBody) (dh *DatabaseHandler) {
	dh = &DatabaseHandler{
		Db: Db.Begin(),
		rt: realtime.NewRealtimeHandler(user),
	}

	runtime.SetFinalizer(dh, finalize)
	return
}

func (dh *DatabaseHandler) CommitOrRollback(err error) error {
	if dh.isClosed {
		log.Error("DatabaseHandler already closed")
		return nil
	}

	dh.isClosed = true
	if dh.Db.Error == nil && err == nil {
		err := dh.Db.Commit().Error
		if err != nil {
			dh.Db.AddError(err)
			dh.Db.Rollback()
			log.Warn("Commit failed, rolling back")
		} else {
			log.Debug("Commit successful")
			dh.rt.PublishEvents()
		}
	} else {
		log.Warn("Rolling back")
		dh.Db.Rollback()
		dh.Db.AddError(err)
	}

	return dh.Db.Error
}

func finalize(dh *DatabaseHandler) {
	if dh.isClosed == false {
		log.Error(cerror.ErrDatabaseHandlerDestroy.Error())
		dh.CommitOrRollback(cerror.ErrDatabaseHandlerDestroy)
	}
}
