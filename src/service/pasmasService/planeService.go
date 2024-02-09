package pasmasservice

import (
	"strconv"

	dh "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/gin-gonic/gin"
)

type PlaneInclude struct {
    IncludeFlights bool
    IncludeAllowedPilots bool
    IncludePrefPilot bool
}

type PlaneFilter struct {
    DivisionId uint
}

func ParsePlaneInclude(c *gin.Context) (*PlaneInclude, error) {
    incFlightStr := c.Query("includeFlights")
    incPilotStr := c.Query("includePilots")
    incPrefPilotStr := c.Query("includePrefPilot")

    include := PlaneInclude{}

    if incFlightStr != "" {
        var err error
        include.IncludeFlights, err = strconv.ParseBool(incFlightStr)

        if err != nil {
            return nil, ErrIncludeNotSupported
        }
    }

    if incPilotStr != "" {
        var err error
        include.IncludeAllowedPilots, err = strconv.ParseBool(incPilotStr)

        if err != nil {
            return nil, ErrIncludeNotSupported
        }
    }

    if incPrefPilotStr != "" {
        var err error
        include.IncludePrefPilot, err = strconv.ParseBool(incPrefPilotStr)

        if err != nil {
            return nil, ErrIncludeNotSupported
        }
    }

    return &include, nil
}

func ParsePlaneFilter(c *gin.Context) (*PlaneFilter, error) {
    divIdStr := c.Query("byDivisionId")

    filter := PlaneFilter{}

    if divIdStr != "" {
        var err error
        id, err := strconv.ParseUint(divIdStr, 10, 32)
        filter.DivisionId = uint(id)

        if err != nil {
            return nil, ErrIncludeNotSupported
        }
    }

    return &filter, nil
}

func GetPlanes(planeInclude *PlaneInclude, planeFilter *PlaneFilter) (*[]model.Plane, error) {
    res := dh.Db
    planes := &[]model.Plane{}

    if planeInclude != nil {
        if planeInclude.IncludeFlights {
            res = res.Preload("Flights")
        }

        if planeInclude.IncludeAllowedPilots {
            res = res.Model(&model.Plane{}).Preload("AllowedPilots")
        }

        if planeInclude.IncludePrefPilot {
            res = res.Preload("PrefPilot")
        }
    }

    if planeFilter != nil {
        if planeFilter.DivisionId != 0 {
            res = res.Where("division_id = ?", planeFilter.DivisionId)
        }
    }

    res = res.Find(planes)
    return planes, res.Error
}

func UpdatePrefPilot(planeId uint, pilotId uint) (*model.Plane, error) {
    plane := &model.Plane{}
    pilot := &model.Pilot{}

    err := dh.Db.First(plane, planeId).Error
    if err != nil {
        return &model.Plane{}, ErrObjectNotFound
    }

    err = dh.Db.First(pilot, pilotId).Error
    if err != nil {
        return &model.Plane{}, ErrObjectDependencyMissing
    }

    //plane.PrefPilotId = pilotId
    plane.PrefPilot = pilot

    err = dh.Db.Model(plane).Updates(plane).Preload("PrefPilot").Error
    return plane, err
}