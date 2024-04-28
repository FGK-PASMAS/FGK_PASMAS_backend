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
        err = cerror.ErrFlightStatusDoesNotFitProcess
    }

    flightCreation.Lock()
    defer flightCreation.Unlock()

    if !flightlogic.CheckIfSlotIsFree(blocker.PlaneId, blocker.DepartureTime, blocker.ArrivalTime) {
        err = cerror.ErrSlotIsNotFree
    }

    dh := databasehandler.NewDatabaseHandler()
    defer func() {
        err = dh.CommitOrRollback(err)
    }()

    newBlocker, err = dh.CreateFlight(blocker)
    
    return
}

// deleteBlocker is handled by the deleteFlight function in the flightService.go file
