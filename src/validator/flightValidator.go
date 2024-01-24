package validator

import (
	"errors"

	"github.com/MetaEMK/FGK_PASMAS_backend/model"
)


func ValidateFlightReservation(flight *model.Flight) error {
    if flight.Type != model.FsReserved {
        return ErrInvalidFlightType
    }

    // if flight.Plane == nil {
    //     return ErrInvalidPlane
    // }
    
    if flight.DepartureTime.IsZero() {
        return ErrInvalidDepartureTime
    }

    // if flight.Pilot == nil {
    //     return ErrInvalidPilot
    // }

    return nil
}

var (
    ErrInvalidFlightType = errors.New("Type is not reserved")
    ErrInvalidPlane = errors.New("Plane is not valid")
    ErrInvalidDepartureTime = errors.New("Departure time is not valid")
    ErrInvalidPilot = errors.New("Pilot is not valid")
)
