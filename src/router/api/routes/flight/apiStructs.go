package flight

import (
	"time"

	"github.com/MetaEMK/FGK_PASMAS_backend/model"
)

type FlightCreation struct {
    Description     *string
    DepartureTime   time.Time
    ArrivalTime     *time.Time

    PlaneId         uint
    Passengers      *[]model.Passenger
}

type FlightReservation struct {
    Description     *string
    Passengers      []model.Passenger
}

type FlightBooking struct {
    Description     *string
    Passengers      []model.Passenger
}
