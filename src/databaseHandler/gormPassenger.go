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
    err = Db.Preload("Flight").Find(&passengers).Error

    for i := range passengers {
        passengers[i].SetTimesToUTC()
    }

    return
}

func (dh *DatabaseHandler) CreatePassenger(pass model.Passenger) (newPassenger model.Passenger, err error) {
    pass.ID = 0

    pass.SetTimesToUTC()

    if pass.Weight <= 0 {
        err = cerror.ErrPassengerWeightIsZero
        dh.Db.AddError(err)
        return
    }

    if pass.FlightID == 0 && pass.Flight == nil {
        err = cerror.ErrObjectDependencyMissing
        dh.Db.AddError(err)
        return
    }

    err = dh.Db.Create(&pass).Error
    if err == cerror.ErrObjectNotFound {
        err = cerror.ErrObjectDependencyMissing
    }
    dh.Db.AddError(err)

    if err == nil {
        newPassenger := model.Passenger{}
        err = dh.Db.Preload("Flight").First(&newPassenger, pass.ID).Error
        newPassenger.SetTimesToUTC()

        dh.rt.AddEvent(realtime.PassengerStream, realtime.CREATED, newPassenger)
    }
    return
}

/*
partialUpdatePassenger updates the newPass with all its attributes.

- 0 or "" values mean that the field should be set nil.

- Nil values are not updated in the database. The newPass is updated in the database and returned.
*/
func (dh DatabaseHandler) PartialUpdatePassenger(id uint, newPass *model.Passenger) {
    var oldPass model.Passenger

    if newPass.Weight <= 0 {
        dh.Db.AddError(cerror.ErrPassengerWeightIsZero)
    }

    err := dh.Db.First(&oldPass, id).Error
    switch err {
    case gorm.ErrRecordNotFound:
        dh.Db.AddError(cerror.ErrObjectDependencyMissing)
    default:
        dh.Db.AddError(err)
    }

    if dh.Db.Error != nil{
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

    if newPass.PassNo != 0 {
        oldPass.PassNo = newPass.PassNo
    }

    err = dh.Db.Updates(&oldPass).Error

    if err != nil {
        dh.Db.AddError(err)
    } else {
        *newPass = oldPass
        newPass.SetTimesToUTC()

        tmpPass := model.Passenger{}
        err = dh.Db.Preload("Flight").First(&tmpPass, oldPass.ID).Error
        tmpPass.SetTimesToUTC()
        dh.rt.AddEvent(realtime.PassengerStream, realtime.UPDATED, tmpPass)
    }
}

// passengerDelete deletes a passenger from the database
func (dh *DatabaseHandler) DeletePassenger(id uint) (passenger model.Passenger, err error) {
    passenger = model.Passenger{}
    err = dh.Db.First(&passenger, id).Error

    if err == cerror.ErrObjectNotFound {
        err = cerror.ErrObjectDependencyMissing
    }

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
