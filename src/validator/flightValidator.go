package validator

import (
	"errors"

	"github.com/MetaEMK/FGK_PASMAS_backend/model"
)


func ValidateFlightReservation(flight *model.Flight) error {
    if flight.Status != model.FsReserved {
        return ErrInvalidFlightType
    }

    if flight.PlaneId == 0 {
        return ErrInvalidPlane
    }
    
    if flight.DepartureTime.IsZero() {
        return ErrInvalidDepartureTime
    }

    var err error
    if len(*flight.Passengers) == 0 {
        return ErrInvalidPassenger
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
        err = errors.Join(err, ErrInvalidPlane)
    }

    if flight.DepartureTime.IsZero() {
        err = errors.Join(err, ErrInvalidDepartureTime)
    }

    return err
}

var (
    ErrInvalidFlightType = errors.New("Type is not reserved")
    ErrInvalidPlane = errors.New("PlaneId is not valid")
    ErrInvalidDepartureTime = errors.New("Departure time is not valid")
    ErrInvalidPilot = errors.New("PilotId is not valid")
    ErrInvalidPassenger = errors.New("Invalid or no passenger")
)
