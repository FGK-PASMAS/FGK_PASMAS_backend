package planeService

import (
	databasehandler "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/MetaEMK/FGK_PASMAS_backend/service/flightService"
)


func UpdatePlane(user model.UserJwtBody, id uint, updateData databasehandler.PartialUpdatePlaneStruct) (newPlane model.Plane, err error) {
    if err = user.ValidateRole(model.Admin); err != nil {
        return
    }

    _, err = databasehandler.GetPlaneById(id, &databasehandler.PlaneInclude{IncludePrefPilot: true, IncludeAllowedPilots: true})
    if err != nil {
        return
    }

    flightService.LockAll()
    dh := databasehandler.NewDatabaseHandler()
    defer func() {
        dh.CommitOrRollback(err)
        flightService.UnlockAll()
    }()

    newPlane, err = dh.PartialUpdatePlane(id, updateData)
    println(*newPlane.PrefPilotId)

    return
}
