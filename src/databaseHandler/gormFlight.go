package databasehandler

import (
	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/realtime"
	"gorm.io/gorm"
)

func initFlight() {
    Db.AutoMigrate(&model.Flight{})
}

func CreateFlight(db *gorm.DB, rt *realtime.RealtimeHandler, flight model.Flight, passengers []model.Passenger) (newFlight model.Flight, newPassengers []model.Passenger, err error) {
    if db == nil {
        db = Db.Begin()
        rt = realtime.NewRealtimeHandler()
        defer func() {
            CommitOrRollback(db, rt)
            err = db.Error
        }()
    }

    if rt == nil {
        err = cerror.ErrNoRealtimeHandlerFound
        return
    }

    flight.ID = 0
    flight.Passengers = nil
    err = db.Create(&flight).Error

    if err != nil {
        db.AddError(err)
        return
    } 

    rt.AddEvent(realtime.FlightStream, realtime.CREATED, flight)

    for index := range passengers {
        passengers[index].FlightID = flight.ID
        pass, err := CreatePassenger(db, rt, passengers[index])

        db.AddError(err)
        newPassengers = append(newPassengers, pass)
    }

    newFlight = flight

    return
}

// partialUpdateFlight updates the newFlight with all set data from newFlight. 0 or "" values means that the field should be set to nil
func PartialUpdateFlight(db *gorm.DB, rt *realtime.RealtimeHandler, id uint, newFlightData model.Flight) (flight model.Flight, err error) {
    if db == nil {
        db = Db
        rt = realtime.NewRealtimeHandler()
        defer func() {
            if err == nil {
                rt.PublishEvents()
            }
        }()
    }

    if rt == nil {
        err = cerror.ErrNoRealtimeHandlerFound
        return
    }
    
    err = Db.First(&flight, id).Error
    if err != nil {
        db.AddError(err)
        return
    }

    if newFlightData.Status == model.FsBooked && flight.Status == model.FsReserved {
        flight.Status = newFlightData.Status
    }

    if newFlightData.Description != nil {
        if *newFlightData.Description == "" {
            flight.Description = nil
        } else {
            flight.Description = newFlightData.Description
        }
    } 

    if newFlightData.FuelAtDeparture != nil {
        if *newFlightData.FuelAtDeparture == 0 {
            flight.FuelAtDeparture = nil
        } else {
            flight.FuelAtDeparture = newFlightData.FuelAtDeparture
        }
    }

    err = db.Updates(&flight).Error
    if err == nil {
        rt.AddEvent(realtime.FlightStream, realtime.UPDATED, flight)
    }

    return
}

// DeleteFlight deletes the flight with the given id and all its passengers. It returns the deleted flight and all its passengers.
func DeleteFlight(db *gorm.DB, rt *realtime.RealtimeHandler, id uint) (flight model.Flight, passengers []model.Passenger, err error) {
    if db == nil {
        db = Db.Begin()
        rt = realtime.NewRealtimeHandler()
        defer func() {
            err = db.Error
        }()
        defer CommitOrRollback(db, rt)
    }

    if rt == nil {
        err = cerror.ErrNoRealtimeHandlerFound
        return
    }

    err = db.Preload("Passengers").First(&flight, id).Error
    if err != nil {
        return
    }
    passengers = *flight.Passengers

    err = db.Delete(&flight, id).Error
    if err != nil {
        return
    } 
    rt.AddEvent(realtime.FlightStream, realtime.DELETED, flight)

    for _, p := range passengers{
        _, delErr := DeletePassenger(db, rt, p.ID)
        db.AddError(delErr)
    }

    return
}
