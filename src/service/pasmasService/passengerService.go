package pasmasservice

import (
	"errors"

	dh "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/realtime"
	"github.com/MetaEMK/FGK_PASMAS_backend/validator"
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

func PassengerCreate(db *gorm.DB, pass *model.Passenger) *gorm.DB {
    if db == nil {
        db = dh.Db
    }

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

    return db
}

func UpdatePassenger(id uint, pass model.Passenger) (model.Passenger, error) {
    err := validator.ValidatePassengerForReserve(pass)
    if err != nil {
        return model.Passenger{}, err
    }

    oldPass := model.Passenger{}
    result := dh.Db.First(&oldPass, id)
    if result.Error != nil {
        return model.Passenger{}, result.Error
    }

    result = dh.Db.Model(&oldPass).Updates(pass)
    if result.Error != nil {
        return model.Passenger{}, result.Error
    }

    realtime.PassengerStream.PublishEvent(realtime.UPDATED, oldPass)
    return oldPass, nil
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
