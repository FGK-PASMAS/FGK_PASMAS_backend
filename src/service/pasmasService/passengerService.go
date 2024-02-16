package pasmasservice

import (
	"errors"

	dh "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
)

var (
    ErrPassengerWeightIsZero = errors.New("Passenger weight must > 0")
)

func GetPassengers() ([]model.Passenger, error) {
    passengers := []model.Passenger{}
    result := dh.Db.Preload("Flight").Find(&passengers)

    return passengers, result.Error
}

func DeletePassenger(id uint) (passenger model.Passenger, err error){
    passenger, err = dh.DeletePassenger(nil, nil, id)

    return
}
