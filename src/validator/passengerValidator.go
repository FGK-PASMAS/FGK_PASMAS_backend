package validator

import (
	"errors"

	"github.com/MetaEMK/FGK_PASMAS_backend/model"
)


var (
    ErrPassengerWeight = errors.New("Weight must be >= 0")
    ErrPassengerLastName = errors.New("Name must not be empty")
)

func ValidatePassenger(pass model.Passenger) error {
    if pass.Weight <= 0 {
        return ErrPassengerWeight
    }

    if len(pass.LastName) == 0 {
        return ErrPassengerLastName
    }

    return nil
}
