package pasmasservice

import (
	"errors"
	"strconv"
	"time"

	dh "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/realtime"
	"github.com/MetaEMK/FGK_PASMAS_backend/validator"
	"github.com/gin-gonic/gin"
)

type FlightInclude struct {
    IncludePassengers bool
    IncludePlane bool
    IncludePilot bool
}

type FlightFilter struct {
    DivisionId uint
    PlaneId uint
}

var (
    ErrSlotIsNotFree = errors.New("Slot is not free")
)

func ParseFlightInclude(c *gin.Context) (*FlightInclude, error) {
    incPassStr := c.Query("includePassengers")
    incPlaneStr := c.Query("includePlane")
    incPilotStr := c.Query("includePilot")

    include := FlightInclude{}

    if incPassStr != "" {
        var err error
        include.IncludePassengers, err = strconv.ParseBool(incPassStr)

        if err != nil {
            return nil, ErrIncludeNotSupported
        }
    }

    if incPlaneStr != "" {
        var err error
        include.IncludePlane, err = strconv.ParseBool(incPlaneStr)

        if err != nil {
            return nil, ErrIncludeNotSupported
        }
    }

    if incPilotStr != "" {
        var err error
        include.IncludePilot, err = strconv.ParseBool(incPilotStr)

        if err != nil {
            return nil, ErrIncludeNotSupported
        }
    }

    return &include, nil
}

func ParseFlightFilter(c *gin.Context) (*FlightFilter, error) {
    divIdStr := c.Query("divisionId")
    planeIdStr := c.Query("planeId")

    filter := FlightFilter{}

    if divIdStr != "" {
        var err error
        id, err := strconv.ParseUint(divIdStr, 10, 64)
        filter.DivisionId = uint(id)

        if err != nil {
            return nil, ErrIncludeNotSupported
        }
    }

    if planeIdStr != "" {
        var err error
        d, err := strconv.ParseUint(planeIdStr, 10, 64)
        filter.PlaneId = uint(d)

        if err != nil {
            return nil, ErrIncludeNotSupported
        }
    }

    return &filter, nil
}

func GetFlights(include *FlightInclude, filter *FlightFilter) (*[]model.Flight, error) {
    res := dh.Db
    flights := &[]model.Flight{}

    if include != nil {
        if include.IncludePassengers {
            res = res.Preload("Passengers")
        }

        if include.IncludePlane {
            res = res.Joins("Plane")
        }

        if include.IncludePilot {
            res = res.Joins("Pilot")
        }
    }

    if filter != nil {
        if filter.DivisionId != 0 {
            res = res.Joins("Plane").Where("division_id = ?", filter.DivisionId)
        }

        if filter.PlaneId != 0 {
            res = res.Where("plane_id = ?", filter.PlaneId)
        }
    }

    res.Find(flights)
    return flights, res.Error
}

func ReserveFlight(flight *model.Flight) (*model.Flight, error) {
    err := validator.ValidateFlightReservation(flight)
    if err != nil {
        return &model.Flight{}, err
    }

    if flight.ArrivalTime.IsZero() {
        flight.ArrivalTime = flight.DepartureTime.Add(flightTimeDuration)
    }

    if !CheckIfSlotIsFree(flight) {
        return &model.Flight{}, ErrSlotIsNotFree
    }
    //TODO: Check parameters for this flight
    //fuel check

    result := dh.Db.Create(flight)

    realtime.FlightStream.PublishEvent(realtime.CREATED, flight)
    return flight, result.Error
}

func BookFlight(id uint, passengers *[]model.Passenger) (*model.Flight, error) {
    for _, pass := range *passengers {
        err := validator.ValidatePassenger(pass)
        if err != nil {
            return &model.Flight{}, err
        }
    }

    flight := model.Flight{}
    res := dh.Db.First(&flight, id)

    if res.Error != nil {
        return &model.Flight{}, res.Error
    }

    flight.Passengers = passengers

    result := dh.Db.Model(&flight).Updates(&flight)
    if result.Error != nil {
        return &model.Flight{}, result.Error
    } else {
        realtime.FlightStream.PublishEvent(realtime.UPDATED, &flight)
        return &flight, nil
    }
}

func DeleteFlights(id uint) error {
    flight := model.Flight{}

    result := dh.Db.Delete(&flight, id)

    if result.RowsAffected != 1 {
        return ErrObjectNotFound
    }

    return result.Error
}

func CheckIfSlotIsFree(flight *model.Flight) bool {
    //TODO: Plane

    flights := []model.Flight{}
    arr_time := flight.ArrivalTime.Truncate(time.Minute).Local()
    dep_time := flight.DepartureTime.Truncate(time.Minute).Local()
    result := dh.Db.Where("arrival_time >= ?", dep_time).Where("departure_time <= ?", arr_time).Find(&flights)

    if result.Error != nil {
        return false
    }

    if len(flights) == 0 {
        return true
    }

    return false
}

