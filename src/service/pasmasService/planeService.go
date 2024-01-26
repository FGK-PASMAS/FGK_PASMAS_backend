package pasmasservice

import (
	"strconv"

	dh "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/gin-gonic/gin"
)

type PlaneInclude struct {
    IncludeFlights bool
    IncludePilots bool
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
        include.IncludePilots, err = strconv.ParseBool(incPilotStr)

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
    divIdStr := c.Query("divisionId")

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
            res = dh.Db.Preload("Flights")
        }

        if planeInclude.IncludePilots {
            res = dh.Db.Preload("Pilots")
        }

        if planeInclude.IncludePrefPilot {
            res = dh.Db.Preload("PrefPilot")
        }
    }

    if planeFilter != nil {
        if planeFilter.DivisionId != 0 {
            res = dh.Db.Where("division_id = ?", planeFilter.DivisionId)
        }
    }

    res = res.Find(planes)
    return planes, res.Error
}
