package pasmasservice

import (
	"errors"
	"sync"

	databasehandler "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
)

var (
    ErrSlotIsNotFree = errors.New("Slot is not free")
    ErrFlightStatusDoesNotFitProcess = errors.New("Flight status does not fit process")
    ErrDepartureTimeIsZero = errors.New("Departure time is zero")
    ErrInvalidArrivalTime = errors.New("Invalid arrival time")
)

var flightCreation sync.Mutex

func GetFlights(include *databasehandler.FlightInclude, filter *databasehandler.FlightFilter) (flights []model.Flight, err error) {
    flights, err = databasehandler.GetFlights(include, filter)
    return
}


func FlightCreation(flight *model.Flight, passengers *[]model.Passenger) (newFlight model.Flight, newPassengers []model.Passenger, err error) {
    var plane model.Plane
    flight.Status = model.FsReserved

    plane, err = databasehandler.GetPlaneById(flight.PlaneId, &databasehandler.PlaneInclude{IncludeDivision: true})
    if err != nil {
        if err == ErrObjectNotFound {
            err = ErrObjectDependencyMissing
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

    println(plane.Division.ID)

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
        dh := databasehandler.NewDatabaseHandler()
        newFlight, newPassengers, err = dh.CreateFlight(*flight, paxs)
        dh.CommitOrRollback(err)
    }

    return
}

func FlightUpdate(flightId uint, newFlightData model.Flight) (flight model.Flight, err error) {
    var passengers []model.Passenger
    var plane model.Plane
    dh := databasehandler.NewDatabaseHandler()
    defer func() {
        dh.CommitOrRollback(err)
        if err == nil {
            flight, err = databasehandler.GetFlightById(flightId, &databasehandler.FlightInclude{IncludePassengers: true, IncludePlane: true})
        }
    }()

    if newFlightData.Passengers != nil {
        passengers = *newFlightData.Passengers
    }

    flight, err = databasehandler.GetFlightById(flightId, &databasehandler.FlightInclude{IncludePassengers: true, IncludePlane: true})
    if err != nil {
        return
    }

    if flight.Status != model.FsReserved {
        err = ErrFlightStatusDoesNotFitProcess
        return
    }

    plane, err = databasehandler.GetPlaneById(flight.PlaneId, &databasehandler.PlaneInclude{IncludeDivision: true})
    if err != nil {
        return 
    }


    passTMP := flight.Passengers
    flight, err = dh.PartialUpdateFlight(flightId, newFlightData)
    flight.Passengers = passTMP

    for index := range passengers {
        passengers[index].FlightID = flight.ID
    }
    partialUpdatePassengers(dh, flight.Passengers, &passengers)

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

func DeleteFlights(id uint) (err error){
    dh := databasehandler.NewDatabaseHandler()
    _, _, err = dh.DeleteFlight(id)

    dh.CommitOrRollback(err)
    return
}

