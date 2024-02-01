package validator

import (
	"errors"

	"github.com/MetaEMK/FGK_PASMAS_backend/model"
)


var (
    ErrPassengerId = errors.New("Id must not empty")
    ErrPassengerWeight = errors.New("Weight must be >= 0")
    ErrPassengerLastName = errors.New("Name must not be empty")
)

func ValidatePassengerForReserve(pass model.Passenger) error {
    if pass.Weight <= 0 {
        return ErrPassengerWeight
    }

    return nil
}

func ValidatePassengerForBooking(pass model.Passenger) error {
    if pass.ID <= 0 {
        return ErrPassengerId
    }

    if pass.LastName == "" {
        return ErrPassengerLastName
    }

    if pass.FirstName == "" {
        return ErrPassengerLastName
    }

    return nil
}
