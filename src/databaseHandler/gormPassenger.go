package databasehandler

import (
	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/realtime"
	"gorm.io/gorm"
)

func initPassenger() {
    Db.AutoMigrate(&model.Passenger{})
}

// GetPassengers returns all passengers from the database
func GetPassengers() (passengers []model.Passenger, err error) {
    err = Db.Preload("Flight").Order("id ASC").Find(&passengers).Error

    for i := range passengers {
        passengers[i].SetTimesToUTC()
    }

    return
}

func (dh *DatabaseHandler) CreatePassenger(pass model.Passenger) (newPassenger model.Passenger, err error) {
    pass.ID = 0

    pass.SetTimesToUTC()

    if pass.Weight <= 0 {
        err = cerror.NewInvalidFlightLogicError("Passenger weight is zero")
        dh.Db.AddError(err)
        return
    }

    if pass.FlightID == 0 && pass.Flight == nil {
        err = cerror.NewDependencyNotFoundError("CreatePassenger: related flight was not found")
        dh.Db.AddError(err)
        return
    }

    err = dh.Db.Create(&pass).Error
    if err == gorm.ErrRecordNotFound {
        err = cerror.NewUnknownError("CreatePassenger: Could not create passenger: " + err.Error())
    }
    dh.Db.AddError(err)

    if err == nil {
        newPassenger = model.Passenger{}
        err = dh.Db.Preload("Flight").First(&newPassenger, pass.ID).Error
        if err != nil {
            return
        }

        newPassenger.SetTimesToUTC()

        dh.rt.AddEvent(realtime.PassengerStream, realtime.CREATED, newPassenger)
    }
    return
}

// passengerDelete deletes a passenger from the database
func (dh *DatabaseHandler) DeletePassenger(id uint) (passenger model.Passenger, err error) {
    passenger = model.Passenger{}
    err = dh.Db.First(&passenger, id).Error
    if err != nil {
        return
    }

    err = dh.Db.Delete(&passenger).Error
    if err != nil {
        dh.Db.AddError(err)
    } else {
        dh.rt.AddEvent(realtime.PassengerStream, realtime.DELETED, passenger)
    }

    return
}
