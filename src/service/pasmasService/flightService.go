package pasmasservice

import (
	"errors"
	"time"

	dh "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/realtime"
	"github.com/MetaEMK/FGK_PASMAS_backend/validator"
)

var (
    ErrSlotIsNotFree = errors.New("Slot is not free")
)

func GetFlights() (*[]model.Flight, error) {
    flights := []model.Flight{}
    result := dh.Db.Find(&flights)

    return &flights, result.Error
}

func ReserveFlight(flight *model.Flight) (*model.Flight, error) {
    err := validator.ValidateFlightReservation(flight)
    if err != nil {
        return &model.Flight{}, err
    }

    if flight.ArrivalTime.IsZero() {
        flight.ArrivalTime = flight.DepartureTime.Add(flightTimeDuration)
    }

    if !CheckIfSlotIsFree(flight) {
        return &model.Flight{}, ErrSlotIsNotFree
    }
    //TODO: Check parameters for this flight

    result := dh.Db.Create(flight)

    realtime.FlightStream.PublishEvent(realtime.CREATED, flight)
    return flight, result.Error
}

func CheckIfSlotIsFree(flight *model.Flight) bool {
    //TODO: Plane

    flights := []model.Flight{}
    arr_time := flight.ArrivalTime.Truncate(time.Minute).Local()
    dep_time := flight.DepartureTime.Truncate(time.Minute).Local()
    result := dh.Db.Where("arrival_time >= ?", dep_time).Where("departure_time <= ?", arr_time).Find(&flights)

    if result.Error != nil {
        return false
    }

    if len(flights) == 0 {
        return true
    }

    return false
}
