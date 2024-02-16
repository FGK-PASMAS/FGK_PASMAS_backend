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


func CreatePassenger(db *gorm.DB, rt *realtime.RealtimeHandler, pass model.Passenger) (newPassenger model.Passenger, err error) {
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

    if err == nil {
        newPassenger = pass
        rt.AddEvent(realtime.PassengerStream, realtime.CREATED, newPassenger)
    }
    return
}

/*
partialUpdatePassenger updates the newPass with all its attributes.

- 0 or "" values mean that the field should be set nil.

- Nil values are not updated in the database. The newPass is updated in the database and returned.
*/
func PartialUpdatePassenger(db *gorm.DB, rt *realtime.RealtimeHandler, id uint, newPass *model.Passenger) {
    var oldPass model.Passenger
    if db == nil {
        db = Db
        rt = realtime.NewRealtimeHandler()
        defer func() {
            if db.Error == nil {
                rt.PublishEvents()
            }
        }()
    }

    if rt == nil {
        return
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

    err = db.Updates(&oldPass).Error

    if err != nil {
        db.AddError(err)
    } else {
        *newPass = oldPass
        rt.AddEvent(realtime.PassengerStream, realtime.UPDATED, oldPass)
    }
}

// passengerDelete deletes a passenger from the database
func DeletePassenger(db *gorm.DB, rt *realtime.RealtimeHandler, id uint) (passenger model.Passenger, err error) {
    if db == nil {
        db = Db
        rt = realtime.NewRealtimeHandler()
        defer func() {
            if db.Error == nil {
                rt.PublishEvents()
            }
        }()
    }

    if rt == nil {
        err = cerror.ErrNoRealtimeHandlerFound
        return
    }

    passenger = model.Passenger{}
    err = db.First(&passenger, id).Error

    if err == cerror.ErrObjectNotFound {
        err = cerror.ErrObjectDependencyMissing
    }

    if err != nil {
        return
    }

    err = db.Delete(&passenger).Error
    if err != nil {
        db.AddError(err)
    } else {
        rt.AddEvent(realtime.PassengerStream, realtime.DELETED, passenger)
    }
    return
}
