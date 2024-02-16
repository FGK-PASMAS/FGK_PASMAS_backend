package pasmasservice

import (
	"errors"

	dh "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/realtime"
)

var (
    ErrPassengerWeightIsZero = errors.New("Passenger weight must > 0")
)

func GetPassengers() ([]model.Passenger, error) {
    passengers := []model.Passenger{}
    result := dh.Db.Preload("Flight").Find(&passengers)

    return passengers, result.Error
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
