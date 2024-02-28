package pasmasservice

import (
	"sync"

	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	databasehandler "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	flightlogic "github.com/MetaEMK/FGK_PASMAS_backend/service/pasmasService/flightLogic"
)

var (
)

var flightCreation sync.Mutex

func GetFlights(include *databasehandler.FlightInclude, filter *databasehandler.FlightFilter) (flights []model.Flight, err error) {
    flights, err = databasehandler.GetFlights(include, filter)
    return
}


func FlightCreation(user model.UserJwtBody, flight model.Flight, passengers *[]model.Passenger) (newFlight model.Flight, newPassengers []model.Passenger, err error) {

    if err = user.ValidateRole(model.Vendor); err != nil {
        return
    }

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

    var paxs []model.Passenger
    if passengers != nil {
        paxs = *passengers
    }

    flightLogicData, err := flightlogic.FlightLogicProcess(flight, plane, *plane.Division, false)

    if err == nil {
        dh := databasehandler.NewDatabaseHandler()
        newFlight, newPassengers, err = dh.CreateFlight(flight, paxs)
        newFlight.FuelAtDeparture = flightLogicData.FuelAtDeparture
        err = dh.CommitOrRollback(err)
    }

    return
}

func FlightUpdate(user model.UserJwtBody, flightId uint, newFlightData model.Flight) (flight model.Flight, err error) {
    if err = user.ValidateRole(model.Vendor); err != nil {
        return
    }

    var passengers []model.Passenger
    var plane model.Plane
    dh := databasehandler.NewDatabaseHandler()
    defer func() {
        err = dh.CommitOrRollback(err)
        if err == nil {
            flight, err = databasehandler.GetFlightById(flightId, &databasehandler.FlightInclude{IncludePassengers: true, IncludePlane: true})
            flight.FuelAtDeparture = newFlightData.FuelAtDeparture
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
        err = cerror.ErrFlightStatusDoesNotFitProcess
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

    var fullValidation bool = false
    if newFlightData.Status == model.FsBooked {
        fullValidation = true
    }

    if newFlightData.Description != nil {
        flight.Description = newFlightData.Description
    }

    newFlightData, err = flightlogic.FlightLogicProcess(flight, plane, *plane.Division, fullValidation)
    if err != nil {
        return
    }

    flight.PilotId = newFlightData.PilotId
    flight.Pilot = newFlightData.Pilot

    return 
}

func DeleteFlights(user model.UserJwtBody, id uint) (err error){
    if err = user.ValidateRole(model.Vendor); err != nil {
        return
    }

    dh := databasehandler.NewDatabaseHandler()
    _, _, err = dh.DeleteFlight(id)

    err = dh.CommitOrRollback(err)
    return
}

func partialUpdatePassengers(dh *databasehandler.DatabaseHandler, oldPass *[]model.Passenger, newPass *[]model.Passenger) {
    if oldPass == nil || newPass == nil {
        return
    }

    if dh == nil {
        dh = databasehandler.NewDatabaseHandler()
        defer dh.CommitOrRollback(nil)
    }

    for i := range *newPass {
        println((*newPass)[i].Action)
        switch (*newPass)[i].Action {
        case model.ActionCreate:
            pass, err := dh.CreatePassenger((*newPass)[i])
            if err == nil {
                tmp := append(*oldPass, pass)
                *oldPass = tmp
            } else {
                dh.Db.AddError(err)
            }
        case model.ActionUpdate:
            status := false
            for j := range *oldPass {
                if (*newPass)[i].ID == (*oldPass)[j].ID {
                    dh.PartialUpdatePassenger((*oldPass)[j].ID, &(*newPass)[i])
                    (*oldPass)[j] = (*newPass)[i]
                    status = true
                }
            }

            if !status {
                dh.Db.AddError(cerror.ErrObjectDependencyMissing)
            }
        case model.ActionDelete:
            dh.DeletePassenger((*newPass)[i].ID)
        
        default:
            dh.Db.AddError(cerror.ErrPassengerActionNotValid)
        }
    }
}
