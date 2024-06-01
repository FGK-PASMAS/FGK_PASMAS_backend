package databasehandler

import (
	"strconv"

	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PlaneInclude struct {
    IncludeDivision bool
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
    incDivStr := c.Query("includeDivision")

    include := PlaneInclude{}

    if incDivStr != "" {
        var err error
        include.IncludeDivision, err = strconv.ParseBool(incDivStr)

        if err != nil {
            return nil, cerror.NewNotValidParametersError("division include not valid")
        }
    }

    if incFlightStr != "" {
        var err error
        include.IncludeFlights, err = strconv.ParseBool(incFlightStr)

        if err != nil {
            return nil, cerror.NewNotValidParametersError("flight include not valid")
        }
    }

    if incPilotStr != "" {
        var err error
        include.IncludeAllowedPilots, err = strconv.ParseBool(incPilotStr)

        if err != nil {
            return nil, cerror.NewNotValidParametersError("pilot include not valid")
        }
    }

    if incPrefPilotStr != "" {
        var err error
        include.IncludePrefPilot, err = strconv.ParseBool(incPrefPilotStr)

        if err != nil {
            return nil, cerror.NewNotValidParametersError("PrefPilot include not valid")
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
            return nil, cerror.NewNotValidParametersError("plane filter not valid")
        }
    }

    return &filter, nil
}


func interpretPlaneConfig(db *gorm.DB, planeInclude *PlaneInclude, planeFilter *PlaneFilter) (*gorm.DB) {
    if planeInclude != nil {
        if planeInclude.IncludeFlights {
            db = db.Preload("Flights")
        }

        if planeInclude.IncludeAllowedPilots {
            db = db.Model(&model.Plane{}).Preload("AllowedPilots")
        }

        if planeInclude.IncludePrefPilot {
            db = db.Preload("PrefPilot")
        }

        if planeInclude.IncludeDivision {
            db = db.Preload("Division")
        }
    }

    if planeFilter != nil {
        if planeFilter.DivisionId != 0 {
            db = db.Where("division_id = ?", planeFilter.DivisionId)
        }
    }

    return db
}

