package databasehandler

import (
	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"gorm.io/gorm"
)

func initPassenger() {
    Db.AutoMigrate(&model.Passenger{})
}


func CreatePassenger(db *gorm.DB, pass *model.Passenger) {
    if db == nil {
        db = Db
    }

    pass.ID = 0

    if pass.Weight <= 0 {
        db.AddError(cerror.ErrPassengerWeightIsZero)
    }

    if pass.FlightID == 0 && pass.Flight == nil {
        db.AddError(cerror.ErrObjectDependencyMissing)
    }

    err := db.Create(pass).Error
    switch err {
    case gorm.ErrRecordNotFound:
        db.AddError(cerror.ErrObjectDependencyMissing)
    default:
        db.AddError(err)
    }
}

// passengerDelete deletes a passenger from the database
func DeletePassenger(db *gorm.DB, id uint) {
    if db == nil {
        db = Db
    }

    pass := model.Passenger{}
    err := db.First(&pass, id).Error

    switch err {
    case gorm.ErrRecordNotFound:
        db.AddError(cerror.ErrObjectNotFound)
    default:
        db.AddError(err)
    }


    err = db.Delete(&pass).Error
    if err != nil {
        db.AddError(err)
    }
}
