package pasmasservice

import (
	"errors"
	"fmt"
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
    result := dh.Db.Preload("Passengers").Find(&flights)

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
    //fuel check

    result := dh.Db.Create(flight)

    realtime.FlightStream.PublishEvent(realtime.CREATED, flight)
    return flight, result.Error
}

func BookFlight(id uint, passengers *[]model.Passenger) (*model.Flight, error) {
    for _, pass := range *passengers {
        err := validator.ValidatePassenger(pass)
        if err != nil {
            println(fmt.Sprintf("Passenger %v is invalid", pass))
            return &model.Flight{}, err
        }
    }

    flight := model.Flight{}
    res := dh.Db.First(&flight, id)

    if res.Error != nil {
        return &model.Flight{}, res.Error
    }

    // db := dh.Db.Begin()
    // for _, pass := range *passengers {
    //     p, err := CreatePassenger(pass)
    //     if err != nil {
    //         db.Rollback()
    //         return &model.Flight{}, err
    //     } else {
    //         flight.Passengers = append(flight.Passengers, p)
    //     }
    // }

    flight.Passengers = *passengers

    result := dh.Db.Model(&flight).Updates(&flight)
    if result.Error != nil {
        return &model.Flight{}, result.Error
    } else {
        realtime.FlightStream.PublishEvent(realtime.UPDATED, &flight)
        return &flight, nil
    }
}

func DeleteFlights(id uint) error {
    flight := model.Flight{}

    result := dh.Db.Delete(&flight, id)

    if result.RowsAffected != 1 {
        return ErrObjectNotFound
    }

    return result.Error
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

