package databasehandler

import (
	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"gorm.io/gorm"
)

func initPassenger() {
    Db.AutoMigrate(&model.Passenger{})
}


func CreatePassenger(db *gorm.DB, pass model.Passenger) (newPassenger model.Passenger, err error) {
    if db == nil {
        db = Db
    }

    pass.ID = 0

    if pass.Weight <= 0 {
        err = cerror.ErrPassengerWeightIsZero
        db.AddError(err)
        return
    }

    if pass.FlightID == 0 && pass.Flight == nil {
        err = cerror.ErrObjectDependencyMissing
        db.AddError(err)
        return
    }

    err = db.Create(&pass).Error
    if err == cerror.ErrObjectNotFound {
        err = cerror.ErrObjectDependencyMissing
    }
    db.AddError(err)

    newPassenger = pass
    return
}

/*
partialUpdatePassenger updates the newPass with all its attributes.

- 0 or "" values mean that the field should be set nil.

- Nil values are not updated in the database. The newPass is updated in the database and returned.
*/
func PartialUpdatePassenger(db *gorm.DB, id uint, newPass *model.Passenger) {
    var oldPass model.Passenger
    if db == nil {
        db = Db
    }

    if newPass.Weight <= 0 {
        db.AddError(cerror.ErrPassengerWeightIsZero)
    }

    err := db.First(&oldPass, id).Error
    switch err {
    case gorm.ErrRecordNotFound:
        db.AddError(cerror.ErrObjectDependencyMissing)
    default:
        db.AddError(err)
    }

    if db.Error != nil{
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

// passengerDelete deletes a passenger from the database
func DeletePassenger(db *gorm.DB, id uint) (passenger model.Passenger, err error) {
    if db == nil {
        db = Db
    }

    passenger = model.Passenger{}
    err = db.First(&passenger, id).Error

    switch err {
    case gorm.ErrRecordNotFound:
        db.AddError(cerror.ErrObjectNotFound)
    default:
        db.AddError(err)
    }

    err = db.Delete(&passenger).Error
    if err != nil {
        db.AddError(err)
    }

    return
}
