package planeService

import (
	databasehandler "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
)

func GetPlanes(includes *databasehandler.PlaneInclude, filters *databasehandler.PlaneFilter) ([]model.Plane, error) {
    planes, err := databasehandler.GetPlanes(includes, filters)

    return planes, err
}

func GetPlaneById(id uint, includes *databasehandler.PlaneInclude) (model.Plane, error) {
    plane, err := databasehandler.GetPlaneById(id, includes)

    return plane, err
}

