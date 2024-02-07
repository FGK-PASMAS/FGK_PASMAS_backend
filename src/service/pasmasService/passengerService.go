package pasmasservice

import (
	"errors"

	dh "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/realtime"
	"gorm.io/gorm"
)

var (
    ErrPassengerWeightIsZero = errors.New("Passenger weight must > 0")
)

func GetPassengers() ([]model.Passenger, error) {
    passengers := []model.Passenger{}
    result := dh.Db.Preload("Flight").Find(&passengers)

    return passengers, result.Error
}

func passengerCreate(db *gorm.DB, pass *model.Passenger) {
    if db == nil {
        db = dh.Db
    }

    pass.ID = 0

    if pass.Weight <= 0 {
        db.AddError(ErrPassengerWeightIsZero)
    }

    if pass.FlightID == 0 && pass.Flight == nil {
        db.AddError(ErrObjectDependencyMissing)
    }

    err := db.Create(pass).Error

    if err == ErrObjectNotFound {
        db.AddError(ErrObjectDependencyMissing)
    }
}

// partialUpdatePassenger updates the newPass with all its attributes.
// 0 or "" values mean that the field should be set nil.
// Nil values are not updated in the database. The newPass is updated in the database and returned.
func partialUpdatePassenger(db *gorm.DB, id uint, newPass *model.Passenger) {
    var oldPass model.Passenger
    if db == nil {
        db = dh.Db
    }

    if newPass.Weight <= 0 {
        db.AddError(ErrPassengerWeightIsZero)
    }

    err := db.First(&oldPass, id).Error
    if err != nil {
        db.AddError(err)
        return
    }

    if newPass.Weight > 0 {
        oldPass.Weight = newPass.Weight
    }

    if newPass.LastName != "" {
        oldPass.LastName = newPass.LastName
    }

    if newPass.FirstName != "" {
        oldPass.FirstName = newPass.FirstName
    }

    db.Updates(&oldPass)
    newPass = &oldPass
}

func DeletePassenger(id uint) error {
    pass := model.Passenger{}
    result := dh.Db.First(&pass, id)
    if result.Error != nil {
        return result.Error
    }

    result = dh.Db.Delete(&pass)
    realtime.PassengerStream.PublishEvent(realtime.DELETED, pass)
    return result.Error
}
