package pasmasservice

import (
	"sync"

	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	databasehandler "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	flightlogic "github.com/MetaEMK/FGK_PASMAS_backend/service/pasmasService/flightLogic"
	"github.com/MetaEMK/FGK_PASMAS_backend/service/pasmasService/noGen"
)

var (
)

var flightCreation sync.Mutex
var flightUpdate sync.Mutex

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

    flight.Passengers = passengers
    flightLogicData, err := flightlogic.FlightLogicProcess(flight, plane, *plane.Division, true)
    flight.ArrivalTime = flightLogicData.ArrivalTime
    flight.Pilot = flightLogicData.Pilot
    flight.PilotId = flightLogicData.ID

    if err == nil {
        dh := databasehandler.NewDatabaseHandler()
        newFlight, err = dh.CreateFlight(flight)
        defer func() {
            err = dh.CommitOrRollback(err)
            if err == nil {
                newFlight, err = databasehandler.GetFlightById(newFlight.ID, &databasehandler.FlightInclude{IncludePassengers: true, IncludePlane: true, IncludePilot: true})
                newFlight.FuelAtDeparture = flightLogicData.FuelAtDeparture
                newFlight.Passengers = passengers
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
    dh := databasehandler.NewDatabaseHandler()
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

    flight, err = dh.PartialUpdateFlight(flightId, newFlightData)

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

    newFlightData, err = flightlogic.FlightLogicProcess(flight, plane, *plane.Division, false)
    if err != nil {
        return
    }

    // TODO: Why is this here?
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

// DEPRECATED 
//
// updates oldPass with data from newPass
func partialUpdatePassengers(dh *databasehandler.DatabaseHandler, oldPass *[]model.Passenger, newPass *[]model.Passenger) {
    if oldPass == nil || newPass == nil {
        return
    }

    if dh == nil {
        dh = databasehandler.NewDatabaseHandler()
        defer dh.CommitOrRollback(nil)
    }

    for i := range *newPass {
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

