package validator

import (
	"errors"

	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
)


func ValidateFlightReservation(flight *model.Flight) error {
    if flight.Status != model.FsReserved {
        return cerror.NewInvalidFlightLogicError("FlightStatus not valid")
    }

    if flight.PlaneId == 0 {
        return cerror.NewInvalidFlightLogicError("PlaneID not valid")
    }

    if flight.DepartureTime.IsZero() {
        return cerror.NewInvalidFlightLogicError("DepartureTime not valid")
    }

    var err error
    if len(*flight.Passengers) == 0 {
        return cerror.NewInvalidFlightLogicError("Passengers not valid")
    } else {
        for _, p := range *flight.Passengers {
            err = errors.Join(err, ValidatePassengerForReserve(p))
        }
    }

    return err
}

func ValidateFlightPlaning(flight *model.Flight) error {
    var err error

    if flight.PlaneId == 0 {
        err = errors.Join(err, cerror.NewInvalidFlightLogicError("PlaneId is not valid"))
    }

    if flight.DepartureTime.IsZero() {
        err = errors.Join(err, cerror.NewInvalidFlightLogicError("DepartureTime is not valid"))
    }

    return err
}
