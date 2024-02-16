package databasehandler

import (
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"gorm.io/gorm"
)

func initFlight() {
    Db.AutoMigrate(&model.Flight{})
}

func CreateFlight(db *gorm.DB, flight model.Flight, passengers []model.Passenger) (newFlight model.Flight, newPassengers []model.Passenger, err error) {
    if db == nil {
        db = Db.Begin()
        defer func() {
            CommitOrRollback(db)
            err = db.Error
        }()
    }

    flight.ID = 0
    flight.Passengers = nil
    db.Create(&flight)

    for index := range passengers {
        passengers[index].FlightID = flight.ID
        pass, err := CreatePassenger(db, passengers[index])

        db.AddError(err)
        newPassengers = append(newPassengers, pass)
    }

    newFlight = flight

    return
}

// partialUpdateFlight updates the newFlight with all set data from newFlight. 0 or "" values means that the field should be set to nil
func PartialUpdateFlight(db *gorm.DB, id uint, newFlightData model.Flight) (flight model.Flight, err error) {
    if db == nil {
        db = Db
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

    return
}

// DeleteFlight deletes the flight with the given id and all its passengers. It returns the deleted flight and all its passengers.
func DeleteFlight(db *gorm.DB, id uint) (flight model.Flight, passengers []model.Passenger, err error) {
    if db == nil {
        db = Db.Begin()
        defer func() {
            err = db.Error
        }()
        defer CommitOrRollback(db)
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

    for _, p := range passengers{
        _, delErr := DeletePassenger(db, p.ID)
        db.AddError(delErr)
    }

    return
}
