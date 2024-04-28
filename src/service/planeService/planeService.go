package planeService

import (
	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	dh "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
)

func GetPlanes(includes *dh.PlaneInclude, filters *dh.PlaneFilter) ([]model.Plane, error) {
    planes, err := dh.GetPlanes(includes, filters)

    return planes, err
}

func UpdatePrefPilot(planeId uint, pilotId uint) (*model.Plane, error) {
    plane := &model.Plane{}
    pilot := &model.Pilot{}

    err := dh.Db.First(plane, planeId).Error
    if err != nil {
        return &model.Plane{}, cerror.ErrObjectNotFound
    }

    err = dh.Db.First(pilot, pilotId).Error
    if err != nil {
        return &model.Plane{}, cerror.ErrObjectDependencyMissing
    }

    //plane.PrefPilotId = pilotId
    plane.PrefPilot = pilot

    err = dh.Db.Model(plane).Updates(plane).Preload("PrefPilot").Error
    return plane, err
}
