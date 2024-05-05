package passengerService

import (
	"errors"

	databasehandler "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
)

var (
    ErrPassengerWeightIsZero = errors.New("Passenger weight must > 0")
)

func GetPassengers() ([]model.Passenger, error) {
    passengers, err := databasehandler.GetPassengers()

    return passengers, err
}

func DeletePassenger(user model.UserJwtBody, id uint) (passenger model.Passenger, err error){
    if err = user.ValidateRole(model.Vendor); err != nil {
        return
    }

    dh := databasehandler.NewDatabaseHandler(user)
    passenger, err = dh.DeletePassenger(id)

    dh.CommitOrRollback(err)
    return
}
