package pasmasservice

import (
	"errors"
	"strconv"
	"sync"

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
    ByDivisionId uint
    ByPlaneId uint
}

var (
    ErrSlotIsNotFree = errors.New("Slot is not free")
    ErrFlightStatusDoesNotFitProcess = errors.New("Flight status does not fit process")
    ErrDepartureTimeIsZero = errors.New("Departure time is zero")
    ErrInvalidArrivalTime = errors.New("Invalid arrival time")
)

var flightCreation sync.Mutex

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
    divIdStr := c.Query("byDivisionId")
    planeIdStr := c.Query("byPlaneId")

    filter := FlightFilter{}

    if divIdStr != "" {
        var err error
        id, err := strconv.ParseUint(divIdStr, 10, 64)
        filter.ByDivisionId = uint(id)

        if err != nil {
            return nil, ErrIncludeNotSupported
        }
    }

    if planeIdStr != "" {
        var err error
        d, err := strconv.ParseUint(planeIdStr, 10, 64)
        filter.ByPlaneId = uint(d)

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
        if filter.ByDivisionId != 0 {
            res = res.Joins("Plane").Where("division_id = ?", filter.ByDivisionId)
        }

        if filter.ByPlaneId != 0 {
            res = res.Where("plane_id = ?", filter.ByPlaneId)
        }
    }

    res.Find(flights)
    return flights, res.Error
}

func FlightCreation(flight *model.Flight, passengers *[]model.Passenger) error {
    var err error
    var plane model.Plane
    flight.Status = model.FsReserved

    err = dh.Db.Preload("Division").First(&plane, flight.PlaneId).Error
    if err != nil {
        if err == ErrObjectNotFound {
            return ErrObjectDependencyMissing
        }

        return err
    }

    if flight.DepartureTime.IsZero() {
        return ErrDepartureTimeIsZero
    }

    if flight.ArrivalTime.IsZero() {
        flight.ArrivalTime = flight.DepartureTime.Add(plane.FlightDuration)
    } else {
        if flight.ArrivalTime.Before(flight.DepartureTime) {
            return ErrInvalidArrivalTime
        }
    }

    fuelAmount, err := calculateFuelAtDeparture(flight, plane)
    if err != nil {
        return err
    }

    if passengers == nil {
        passengers = &[]model.Passenger{}
    }
    passWeight, err := checkPassengerAndCalcWeight(*passengers, plane.MaxSeatPayload, 0, plane.Division.PassengerCapacity, false)
    if err != nil {
        return err
    }

    pilot, err := calculatePilot(passWeight, fuelAmount, plane)
    if err != nil {
        return err
    }
    flight.Pilot = &pilot

    flightCreation.Lock()
    defer flightCreation.Unlock()
    if(!checkIfSlotIsFree(flight.PlaneId, flight.DepartureTime, flight.ArrivalTime)) {
        return ErrSlotIsNotFree
    }

    if err == nil {
        db := dh.Db.Begin()
        db.Create(flight)

        for index := range *passengers {
            (*passengers)[index].FlightID = flight.ID
            passengerCreate(db, &(*passengers)[index])
        }

        if db.Error != nil {
            db.Rollback()
            return db.Error
        } else {
            err = db.Commit().Error
            if err != nil {
                return err
            }
            flight.Passengers = passengers

            go realtime.FlightStream.PublishEvent(realtime.CREATED, flight)
            go realtime.PassengerStream.PublishEvent(realtime.CREATED, *passengers)
        }
    }

    return err
}

func FlightUpdate(flightId uint, newFlightData *model.Flight) (error) {
    var flight model.Flight
    var passengers []model.Passenger
    if newFlightData.Passengers != nil {
        passengers = *newFlightData.Passengers
    }

    err := dh.Db.Preload("Plane").Preload("Passengers").First(&flight, flightId).Error
    if err != nil {
        return err
    }

    if flight.Status != model.FsReserved {
        return ErrFlightStatusDoesNotFitProcess
    }

    err = dh.Db.Preload("Division").First(&flight.Plane, flight.PlaneId).Error
    if err != nil {
        return err
    }

    db := dh.Db.Begin()
    partialUpdateFlight(db, flightId, newFlightData)
    for index := range passengers {
        passengers[index].FlightID = flight.ID
    }
    partialUpdatePassengers(db, flight.Passengers, &passengers)


    var minPass uint
    var fullPassCheck bool = false
    if newFlightData.Status == model.FsBooked {
        minPass = 1
        fullPassCheck = true
    }

    passWeight, err := checkPassengerAndCalcWeight(*flight.Passengers, flight.Plane.MaxSeatPayload, minPass, flight.Plane.Division.PassengerCapacity, fullPassCheck)
    if err != nil {
        db.Rollback()
        return err
    }

    fuelAmount, err := calculateFuelAtDeparture(&flight, *flight.Plane)
    if err != nil {
        db.Rollback()
        return err
    }

    pilot, err := calculatePilot(passWeight, fuelAmount, *flight.Plane)
    if err != nil {
        db.Rollback()
        return err
    }
    flight.PilotId = &pilot.ID

    if newFlightData.Description != nil {
        flight.Description = newFlightData.Description
    }

    err = checkFlightValidation(flight)
    if err != nil {
        db.Rollback()
        return err
    }

    err = db.Commit().Error
    if err != nil {
        db.Rollback()
        return err
    }

    dh.Db.Preload("Passengers").First(&newFlightData, flightId)
    go realtime.FlightStream.PublishEvent(realtime.UPDATED, flight)
    go sendRealtimeEventsForPassengers(passengers, realtime.PING)

    return err
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

