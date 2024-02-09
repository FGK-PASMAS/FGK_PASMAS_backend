package pasmasservice

import (
	dh "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"gorm.io/gorm"
)


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

func passengerDelete(db *gorm.DB, id uint) {
    if db == nil {
        db = dh.Db
    }

    pass := model.Passenger{}
    err := db.First(&pass, id).Error

    if err != nil {
        db.AddError(err)
        return
    }

    err = db.Delete(&pass).Error

    if err != nil {
        db.AddError(err)
    }
}
