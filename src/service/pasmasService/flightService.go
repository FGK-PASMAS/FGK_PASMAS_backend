package pasmasservice

import (
	"errors"
	"strconv"
	"time"

	dh "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/realtime"
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
    ErrFlightStatusDoesNotFitProcess = errors.New("Flight status does not fit process")
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

func FlightPlanning(planeId uint, departureTime time.Time, arrivalTime *time.Time, description *string) (*model.Flight, error) {
    var plane model.Plane
    var err error
    var newFlight model.Flight

    newFlight.Status = model.FsPlanned
    newFlight.Description = description
    newFlight.DepartureTime = departureTime

    err = dh.Db.First(&plane, planeId).Error
    if err == ErrObjectNotFound {
        return &model.Flight{}, ErrObjectDependencyMissing
    }
    newFlight.PlaneId = plane.ID

    if arrivalTime == nil || arrivalTime.IsZero() {
        arrTime := departureTime.Add(plane.FlightDuration)
        newFlight.ArrivalTime = arrTime
    } else {
        newFlight.ArrivalTime = *arrivalTime
    }

    if !checkIfSlotIsFree(newFlight.PlaneId, newFlight.DepartureTime, newFlight.ArrivalTime) {
        return &model.Flight{}, ErrSlotIsNotFree
    }

    fuelAmount, err := calculateFuelAtDeparture(&newFlight, plane)
    if err != nil {
        return &model.Flight{},err
    }

    if err == nil {
        err = dh.Db.Create(&newFlight).Error
        if err == nil {
            newFlight.FuelAtDeparture = &fuelAmount
            go realtime.FlightStream.PublishEvent(realtime.CREATED, &newFlight)
            return &newFlight, nil
        }
    }
    
    return &model.Flight{}, err
}

func FlightReservation(flightId uint, passengers *[]model.Passenger, description *string) (*model.Flight, error) {
    var flight model.Flight
    flight.ID = flightId

    err := dh.Db.Preload("Plane").First(&flight, flightId).Error
    if err != nil {
        return &model.Flight{}, err
    }

    if flight.Status != model.FsPlanned {
        return &model.Flight{}, ErrFlightStatusDoesNotFitProcess
    }

    flight.Status = model.FsReserved

    passWeight, err := calculatePassWeight(*passengers, flight.Plane.MaxSeatPayload)
    if err != nil {
        return &model.Flight{}, err
    }
    flight.Passengers = passengers

    fuelAmount, err := calculateFuelAtDeparture(&flight, *flight.Plane)
    if err != nil {
        return &model.Flight{}, err
    }

    pilot, err := calculatePilot(passWeight, fuelAmount, *flight.Plane)
    if err != nil {
        return &model.Flight{}, err
    }
    flight.PilotId = &pilot.ID

    if description != nil {
        flight.Description = description
    }

    err = dh.Db.Updates(&flight).Error

    if err == nil {
        go realtime.FlightStream.PublishEvent(realtime.CREATED, flight)
        go realtime.PassengerStream.PublishEvent(realtime.CREATED, flight.Passengers)
    }

    return &flight, err
}

func BookFlight(id uint, passengers *[]model.Passenger, description *string) (*model.Flight, error) {
    var flight model.Flight
    var err error

    err = dh.Db.Preload("Plane").Preload("Passengers").First(&flight).Error
    if err != nil {
        return &model.Flight{}, ErrObjectNotFound
    }

    if flight.Status != model.FsReserved {
        return &model.Flight{}, ErrFlightStatusDoesNotFitProcess
    }
    flight.Status = model.FsBooked

    if description != nil {
        if *description == "" {
            flight.Description = nil
        } else {
            flight.Description = description
        }
    }

    for _, p := range *passengers{
        if !partialUpdatePassenger(flight.Passengers, p) {
            return &model.Flight{}, ErrObjectDependencyMissing
        }
    }

    err = checkFlightValidation(flight)
    if err != nil {
        return &model.Flight{}, err
    }

    db := dh.Db.Begin().Debug()

    for index := range *flight.Passengers {
        println((*flight.Passengers)[index].ID)
        err = dh.Db.Debug().Updates(&(*flight.Passengers)[index]).Error
    }

    db = db.Updates(&flight)

    if db.Error != nil {
        db.Rollback()
        return &model.Flight{}, db.Error
    } else {
        err = db.Commit().Error
    }

    if err != nil {
        return &model.Flight{}, err
    } else {
        go realtime.FlightStream.PublishEvent(realtime.UPDATED, flight)
        go realtime.PassengerStream.PublishEvent(realtime.UPDATED, flight.Passengers)

        return &flight, nil
    }
}

func DeleteFlights(id uint) error {
    var err error
    flight := model.Flight{}

    dh.Db.Preload("Passengers").First(&flight, id)
    result := dh.Db.Delete(&flight, id)

    if result.RowsAffected != 1 {
        return ErrObjectNotFound
    }

    err = errors.Join(err, result.Error)

    result = dh.Db.Delete(&flight.Passengers)

    if result.Error != nil {
        go realtime.PassengerStream.PublishEvent(realtime.DELETED, flight.Passengers)
        go realtime.FlightStream.PublishEvent(realtime.DELETED, flight)
    }

    return result.Error
}

