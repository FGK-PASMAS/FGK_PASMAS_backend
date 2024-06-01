package flightService

import (
	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	databasehandler "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	flightlogic "github.com/MetaEMK/FGK_PASMAS_backend/service/flightService/flightLogic"
)

// CreateBlocker creates a new blocker for a plane
//
// The user must be an admin
//
// this uses the flightCreation mutex
func CreateBlocker(user model.UserJwtBody, blocker model.Flight) (newBlocker model.Flight, err error) {
    if err = user.ValidateRole(model.Admin); err != nil {
        return
    }

    if blocker.Status != model.FsBlocked {
        err = cerror.NewInvalidFlightLogicError("Flight status does not fit current process")
    }

    if blocker.ArrivalTime.IsZero() {
        err = cerror.NewInvalidRequestBodyError("ArrivalTime is zero")
        return
    }

    LockFlightCreation()
    defer UnlockFlightCreation()

    if flightlogic.CheckIfSlotIsFree(blocker.PlaneId, blocker.DepartureTime, blocker.ArrivalTime) == false {
        err = cerror.NewInvalidFlightLogicError("Slot is not free")
        return
    }

    dh := databasehandler.NewDatabaseHandler(user)
    defer func() {
        err = dh.CommitOrRollback(err)
    }()

    newBlocker, err = dh.CreateFlight(blocker)

    return
}
