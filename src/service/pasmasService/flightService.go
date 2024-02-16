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

func FlightCreation(flight *model.Flight, passengers *[]model.Passenger) (newFlight model.Flight, newPassengers []model.Passenger, err error) {
    var plane model.Plane
    flight.Status = model.FsReserved

    err = dh.Db.Preload("Division").First(&plane, flight.PlaneId).Error
    if err != nil {
        if err == ErrObjectNotFound {
            err = ErrObjectDependencyMissing
            return 
        }

        return 
    }

    if flight.DepartureTime.IsZero() {
        err = ErrDepartureTimeIsZero
        return 
    }

    if flight.ArrivalTime.IsZero() {
        flight.ArrivalTime = flight.DepartureTime.Add(plane.FlightDuration)
    } else {
        if flight.ArrivalTime.Before(flight.DepartureTime) {
            err = ErrInvalidArrivalTime
            return 
        }
    }

    fuelAmount, err := calculateFuelAtDeparture(flight, plane)
    if err != nil {
        return 
    }

    var paxs []model.Passenger
    if passengers != nil {
        paxs = *passengers
    }


    passWeight, err := checkPassengerAndCalcWeight(paxs, plane.MaxSeatPayload, 0, plane.Division.PassengerCapacity, false)
    if err != nil {
        return 
    }

    pilot, err := calculatePilot(passWeight, fuelAmount, plane)
    if err != nil {
        return 
    }
    flight.Pilot = &pilot

    flightCreation.Lock()
    defer flightCreation.Unlock()
    if(!checkIfSlotIsFree(flight.PlaneId, flight.DepartureTime, flight.ArrivalTime)) {
        err = ErrSlotIsNotFree
        return
    }

    if err == nil {
        newFlight, newPassengers, err = dh.CreateFlight(nil, *flight, paxs)
        if err != nil {
            return 
        }

        go realtime.FlightStream.PublishEvent(realtime.CREATED, flight)
        go sendRealtimeEventsForPassengers(*passengers, realtime.CREATED)
    }

    return
}

func FlightUpdate(flightId uint, newFlightData model.Flight) (flight model.Flight, passengers []model.Passenger, err error) {
    var plane model.Plane
    if newFlightData.Passengers != nil {
        passengers = *newFlightData.Passengers
    }

    err = dh.Db.Preload("Passengers").Preload("Plane").First(&flight, flightId).Error
    if err != nil {
        return
    }

    if flight.Status != model.FsReserved {
        err = ErrFlightStatusDoesNotFitProcess
        return
    }

    err = dh.Db.Preload("Division").First(&plane, flight.PlaneId).Error
    if err != nil {
        return 
    }

    db := dh.Db.Begin()
    defer func() {
        if err != nil {
            db.Rollback()
        } else {
            err = db.Commit().Error
            if err != nil {
                dh.Db.Preload("Passengers").First(&newFlightData, flightId)
                go realtime.FlightStream.PublishEvent(realtime.UPDATED, flight)
                go sendRealtimeEventsForPassengers(passengers, realtime.PING)
            }
        }
    }()

    passTMP := flight.Passengers
    flight, err = dh.PartialUpdateFlight(db, flightId, newFlightData)
    flight.Passengers = passTMP

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

    passWeight, err := checkPassengerAndCalcWeight(
        *flight.Passengers,
        plane.MaxSeatPayload,
        minPass,
        plane.Division.PassengerCapacity,
        fullPassCheck,
    )
    if err != nil {
        return 
    }

    fuelAmount, err := calculateFuelAtDeparture(&flight, plane)
    if err != nil {
        return
    }

    pilot, err := calculatePilot(passWeight, fuelAmount, plane)
    if err != nil {
        return
    }
    flight.PilotId = &pilot.ID

    if newFlightData.Description != nil {
        flight.Description = newFlightData.Description
    }

    err = checkFlightValidation(flight)
    return 
}

func DeleteFlights(id uint) error {
    flight := model.Flight{}

    dh.Db.Preload("Passengers").First(&flight, id)
    result := dh.Db.Delete(&flight, id)

    if result.RowsAffected != 1 {
        return ErrObjectNotFound
    }

    if len(*flight.Passengers) > 0 {
        result = dh.Db.Delete(&flight.Passengers)
    }

    if result.Error != nil {
        go realtime.PassengerStream.PublishEvent(realtime.DELETED, flight.Passengers)
        go realtime.FlightStream.PublishEvent(realtime.DELETED, flight)
    }

    return result.Error
}

