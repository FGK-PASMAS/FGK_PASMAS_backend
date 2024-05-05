package flightService

import (
	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	databasehandler "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	flightlogic "github.com/MetaEMK/FGK_PASMAS_backend/service/flightService/flightLogic"
	"github.com/MetaEMK/FGK_PASMAS_backend/service/flightService/noGen"
)

func GetFlights(include *databasehandler.FlightInclude, filter *databasehandler.FlightFilter) (flights []model.Flight, err error) {
    flights, err = databasehandler.GetFlights(include, filter)
    return
}


func FlightCreation(user model.UserJwtBody, flight model.Flight, passengers *[]model.Passenger) (newFlight model.Flight, newPassengers []model.Passenger, err error) {
    if err = user.ValidateRole(model.Vendor); err != nil {
        return
    }

    flight.DepartureTime = flight.DepartureTime.UTC()
    flight.ArrivalTime = flight.ArrivalTime.UTC()

    var plane model.Plane
    flight.Status = model.FsReserved

    plane, err = databasehandler.GetPlaneById(flight.PlaneId, &databasehandler.PlaneInclude{IncludeDivision: true})
    if err != nil {
        if err == ErrObjectNotFound {
            err = ErrObjectDependencyMissing
        }
        return 
    }

    flightCreation.Lock()
    defer flightCreation.Unlock()

    flightLogicData, err := flightlogic.FlightLogicProcess(flight, plane, *plane.Division, true)
    flight.ArrivalTime = flightLogicData.ArrivalTime.UTC()
    flight.Pilot = flightLogicData.Pilot
    flight.PilotId = flightLogicData.PilotId

    if err == nil {
        dh := databasehandler.NewDatabaseHandler(user)
        newFlight, err = dh.CreateFlight(flight)
        defer func() {
            err = dh.CommitOrRollback(err)
            if err == nil {
                newFlight, err = databasehandler.GetFlightById(newFlight.ID, &databasehandler.FlightInclude{IncludePassengers: true, IncludePlane: true, IncludePilot: true})
                newFlight.FuelAtDeparture = flightLogicData.FuelAtDeparture
            }
        }()
    }

    return
}

func FlightBooking(user model.UserJwtBody, flightId uint, newFlightData model.Flight) (flight model.Flight, err error) {
    if err = user.ValidateRole(model.Vendor); err != nil {
        return
    }

    var passengers []model.Passenger
    var plane model.Plane
    dh := databasehandler.NewDatabaseHandler(user)
    flightUpdate.Lock()
    defer func() {
        err = dh.CommitOrRollback(err)
        if err == nil {
            flight, err = databasehandler.GetFlightById(flightId, &databasehandler.FlightInclude{IncludePassengers: true, IncludePlane: true, IncludePilot: true})
            flight.FuelAtDeparture = newFlightData.FuelAtDeparture
        }

        flightUpdate.Unlock()
    }()

    if flight.Status != model.FsReserved && newFlightData.Status != model.FsBooked {
        err = cerror.ErrFlightStatusDoesNotFitProcess
        return
    }

    if newFlightData.Passengers != nil {
        passengers = *newFlightData.Passengers
    }

    flight, err = databasehandler.GetFlightById(flightId, &databasehandler.FlightInclude{IncludePlane: true})
    if err != nil {
        return
    }

    plane, err = databasehandler.GetPlaneById(flight.PlaneId, &databasehandler.PlaneInclude{IncludeDivision: true})
    if err != nil {
        return 
    }

    flightNo, err := noGen.GenerateFlightNo(plane)
    if err != nil {
        return
    }
    newFlightData.FlightNo = &flightNo
    
    if newFlightData.Description != nil {
        flight.Description = newFlightData.Description
    }


    for index := range passengers {
        passengers[index].FlightID = flight.ID
    }

    err = noGen.GeneratePassNo(plane, &passengers)
    if err != nil {
        return
    }

    flight.Passengers = new([]model.Passenger)
    for _, p := range passengers {
        var newPass model.Passenger
        newPass, err = dh.CreatePassenger(p)
        if err != nil {
            return
        }
        *flight.Passengers = append(*flight.Passengers, newPass)
    }

    flight, err = dh.PartialUpdateFlight(flightId, newFlightData)
    if err != nil {
        return
    }

    newFlightData, err = flightlogic.FlightLogicProcess(flight, plane, *plane.Division, false)
    if err != nil {
        return
    }

    // TODO: Why is this here?
    flight.PilotId = newFlightData.PilotId
    flight.Pilot = newFlightData.Pilot

    return 
}

// deletes a flight or blocker
//
// if the flight is a blocker, the user must be an admin
//
// if the flight is a normal flight, the user must be a vendor
func DeleteFlights(user model.UserJwtBody, id uint) (err error){
    flight, err := databasehandler.GetFlightById(id, nil)

    // if the flight is a blocker, the user must be an admin
    // if the flight is a normal flight, the user must be a vendor
    if flight.Status == model.FsBlocked {
        if err = user.ValidateRole(model.Admin); err != nil {
            return
        } else {
            if err = user.ValidateRole(model.Vendor); err != nil {
                return
            }
        }
    }

    dh := databasehandler.NewDatabaseHandler(user)
    defer func() {
        err = dh.CommitOrRollback(err)
    }()

    _, _, err = dh.DeleteFlight(id)
    return
}

