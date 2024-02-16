package databasehandler

import (
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"gorm.io/gorm"
)

func initFlight() {
    Db.AutoMigrate(&model.Flight{})
}

func CreateFlight(db *gorm.DB, flight *model.Flight, passengers *[]model.Passenger) {
    if db == nil {
        db = Db
    }

    flight.ID = 0
    flight.Passengers = nil
    db.Create(flight)

    for index := range *passengers {
        (*passengers)[index].FlightID = flight.ID
        CreatePassenger(db, &(*passengers)[index])
    }
}

func UpdateFlight() {

}

func DeleteFlight() {

}
