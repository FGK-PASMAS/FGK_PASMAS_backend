package validator

import (
	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
)

func ValidatePassengerForReserve(pass model.Passenger) error {
    if pass.Weight <= 0 {
        return cerror.NewInvalidFlightLogicError("PassengerWeight must not be empty")
    }

    return nil
}

func ValidatePassengerForBooking(pass model.Passenger) error {
    if pass.ID <= 0 {
        return cerror.NewInvalidFlightLogicError("PassengerID must be not empty")
    }

    if pass.LastName == "" {
        return cerror.NewInvalidFlightLogicError("PassengerName must not be empty")
    }

    if pass.FirstName == "" {
        return cerror.NewInvalidFlightLogicError("PassengerName must not be empty")
    }

    if pass.Weight <= 0 {
        return cerror.NewInvalidFlightLogicError("PassengerWeight must not be empty")
    }

    return nil
}
