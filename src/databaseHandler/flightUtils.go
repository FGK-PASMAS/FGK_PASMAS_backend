package databasehandler

import (
	"strconv"
	"time"

	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


type FlightInclude struct {
    IncludePassengers bool
    IncludePlane bool
    IncludePilot bool
}

type FlightFilter struct {
    ByDivisionId uint
    ByPlaneId uint
    ByDepartureTime time.Time
}

func ParseFlightInclude(c *gin.Context) (*FlightInclude, error) {
    incPassStr := c.Query("includePassengers")
    incPlaneStr := c.Query("includePlane")
    incPilotStr := c.Query("includePilot")

    include := FlightInclude{}

    if incPassStr != "" {
        var err error
        include.IncludePassengers, err = strconv.ParseBool(incPassStr)

        if err != nil {
            return nil, cerror.NewNotValidParametersError("Include not supported")
        }
    }

    if incPlaneStr != "" {
        var err error
        include.IncludePlane, err = strconv.ParseBool(incPlaneStr)

        if err != nil {
            return nil, cerror.NewNotValidParametersError("Include not supported")
        }
    }

    if incPilotStr != "" {
        var err error
        include.IncludePilot, err = strconv.ParseBool(incPilotStr)

        if err != nil {
            return nil, cerror.NewNotValidParametersError("Include not supported")
        }
    }

    return &include, nil
}

func ParseFlightFilter(c *gin.Context) (*FlightFilter, error) {
    divIdStr := c.Query("byDivisionId")
    planeIdStr := c.Query("byPlaneId")
    departTimeStr := c.Query("byDepartureTime")

    filter := FlightFilter{}

    if divIdStr != "" {
        var err error
        id, err := strconv.ParseUint(divIdStr, 10, 64)
        filter.ByDivisionId = uint(id)

        if err != nil {
            return nil, cerror.NewNotValidParametersError("Filter not supported")
        }
    }

    if planeIdStr != "" {
        var err error
        d, err := strconv.ParseUint(planeIdStr, 10, 64)
        filter.ByPlaneId = uint(d)

        if err != nil {
            return nil, cerror.NewNotValidParametersError("Filter not supported")
        }
    }

    if departTimeStr != "" {
        var err error
        filter.ByDepartureTime, err = time.Parse(time.RFC3339, departTimeStr)
        filter.ByDepartureTime = filter.ByDepartureTime.UTC()

        if err != nil {
            return nil, cerror.NewNotValidParametersError("Filter not supported")
        }
    }

    return &filter, nil
}

func interpretFlightConfig(db *gorm.DB, flightInclude *FlightInclude, flightFilter *FlightFilter) (*gorm.DB) {
    if flightInclude != nil {
        if flightInclude.IncludePassengers {
            db = db.Preload("Passengers")
        }

        if flightInclude.IncludePlane {
            db = db.Preload("Plane")
        }

        if flightInclude.IncludePilot {
            db = db.Preload("Pilot")
        }
    }

    if flightFilter != nil {
        if flightFilter.ByDivisionId != 0 {
            db = db.Joins("Plane").Where("division_id = ?", flightFilter.ByDivisionId)
        }

        if flightFilter.ByPlaneId != 0 {
            db = db.Where("plane_id = ?", flightFilter.ByPlaneId)
        }

        if !flightFilter.ByDepartureTime.IsZero() {
            db = db.Where("departure_time >= ?", flightFilter.ByDepartureTime)
        }
    }

    return db
}
