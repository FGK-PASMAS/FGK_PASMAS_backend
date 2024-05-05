package model

import (
	"time"

	"gorm.io/gorm"
)

type Flight struct {
    gorm.Model

    Status               FlightType
    FlightNo            *string     `gorm:"uniqueIndex"`

    Description         *string
    FuelAtDeparture     *float32
    DepartureTime       time.Time
    ArrivalTime         time.Time

    PlaneId             uint                
    Plane               *Plane                  `gorm:"foreignKey:PlaneId"`
    PilotId             *uint
    Pilot               *Pilot                  `gorm:"foreignKey:PilotId"`
    Passengers          *[]Passenger            `gorm:"foreignKey:FlightID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}


type FlightType string

const (
    FsPlanned = "PLANNED"
    FsReserved = "RESERVED"
    FsBooked = "BOOKED"
    FsBlocked = "BLOCKED"
)

func (f * Flight) SetTimesToUTC() {
    f.CreatedAt = f.CreatedAt.UTC()
    f.UpdatedAt = f.UpdatedAt.UTC()
    f.DepartureTime = f.DepartureTime.UTC()
    f.ArrivalTime = f.ArrivalTime.UTC()

    if f.Plane != nil {
        f.Plane.SetTimesToUTC()
    }

    if f.Pilot != nil {
        f.Pilot.SetTimesToUTC()
    }

    if f.Passengers != nil {
        for _, p := range *f.Passengers {
            p.SetTimesToUTC()
        }
    }
}
