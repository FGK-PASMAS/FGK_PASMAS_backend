package pasmasservice

import (
	"strconv"

	dh "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/gin-gonic/gin"
)

type PilotInclude struct {
    Plane bool
}

type PilotFilter struct {
    PlaneId uint
}

func ParsePilotInclude(c *gin.Context) (*PilotInclude, error) {
    println("Include")
    planeStr := c.Query("includePlane")

    pilotInclude := PilotInclude{}

    if planeStr != "" {
        var err error
        pilotInclude.Plane, err = strconv.ParseBool(planeStr)
        if err != nil {
            return &PilotInclude{}, ErrIncludeNotSupported
        }
    }

    return &pilotInclude, nil
}

func ParsePilotFilter(c *gin.Context) (*PilotFilter, error) {
    println("Filter")
    planeIdStr := c.Query("planeId")

    pilotfilter := PilotFilter{}

    if planeIdStr != "" {
        value, err := strconv.ParseUint(planeIdStr, 10, 64)
        if err != nil {
            return &PilotFilter{}, ErrIncludeNotSupported
        } else {
            pilotfilter.PlaneId = uint(value)
        }
    }

    return &pilotfilter, nil
} 

func GetPilots(include *PilotInclude, filter *PilotFilter) (*[]model.Pilot, error) {
    db := dh.Db
    var pilots *[]model.Pilot

    if include != nil {
        if include.Plane {
            db = db.Model(&model.Pilot{}).Preload("AllowedPilots")
        }
    }

    if filter != nil {
        if filter.PlaneId != 0 {
            db = db.Preload("AllowedPlanes", "id = ?", filter.PlaneId)
        }
    }

    res := db.Find(&pilots)
    return pilots, res.Error
}
