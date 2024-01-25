package model

import (
	"time"

	"gorm.io/gorm"
)

type Flight struct {
    gorm.Model
    Type                FlightType
    //Plane               interface{}     //TODO: Add Plane reference

    FuelAtDeparture     float32
    DepartureTime       time.Time
    ArrivalTime         time.Time

    //Pilot               interface{}     //TODO: Add Pilot reference
    Passengers          []Passenger         `gorm:"foreignKey:FlightID"`
}


type FlightType string

const (
    FsBlocked = "BLOCKED"
    FsReserved = "RESERVED"
    FsBooked = "BOOKED"
)
