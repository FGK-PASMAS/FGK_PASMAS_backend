package model

import (
	"time"

	"gorm.io/gorm"
)

type Plane struct {
	gorm.Model

	// Aricraft Registration f.E. D-ELHN
    Registration string `gorm:"unique;not null"`

	// Aircraft Type f.E. C172
	AircraftType string `gorm:"not null"`

    // Default FlightDuration in Default FlightDuration in NanoSeconds
    FlightDuration time.Duration `gorm:"not null"`

	// Maximum amount of fuel this plane can take in liters; -1 if not applicable
	FuelMaxCapacity int `gorm:"not null"`

    // Fuel amount to start with
    FuelStartAmount uint 

	// Fuel consumption per flight in liters; -1 if not applicable
	FuelburnPerFlight float32 `gorm:"not null"`

	// Conversion factor of fuel in kg/l; -1 if not applicable
	FuelConversionFactor float32 `gorm:"not null"`

	// Maximum takeoff weight in kg
	MTOW uint `gorm:"not null"`

	// Empty weight of the aircraft in kg
	EmptyWeight uint `gorm:"not null"`

    // Maximum weight of one seat in kg; -1 if not applicable
    MaxSeatPayload int `gorm:"not null"`

	// Aircrafts division f.E. "Motorflug"
	DivisionId uint     `gorm:"index"`
	Division   *Division `gorm:"foreignKey:DivisionId;OnUpdate:CASCADE;OnDelete:RESTRICT"`

    // Contains all pilots who are allowed to fly this aircraft
    AllowedPilots *[]Pilot `gorm:"many2many:AllowedPilots;"`

    // Contains the pilot who should fly all new flights by default
    PrefPilotId *uint 
    // Contains the pilot who should fly all new flights by default
    PrefPilot *Pilot `gorm:"foreignKey:PrefPilotId"`

    // Contains all flights flown by this aircraft
    Flights *[]Flight `gorm:"foreignKey:PlaneId"`

    // Start time slot?
    SlotStartTime time.Time `gorm:"not null"`

    // End time slot?
    SlotEndTime time.Time `gorm:"not null"`

    // Base value for passNo
    PassNoBase uint `gorm:"not null" json:"-"`
}

func (p * Plane) SetTimesToUTC() {
    p.CreatedAt = p.CreatedAt.UTC()
    p.UpdatedAt = p.UpdatedAt.UTC()
    p.SlotStartTime = p.SlotStartTime.UTC()
    p.SlotEndTime = p.SlotEndTime.UTC()

    if p.Division != nil {
        p.Division.SetTimesToUTC()
    }

    if p.PrefPilot != nil {
        p.PrefPilot.SetTimesToUTC()
    }

    if p.AllowedPilots != nil {
        for _, a := range *p.AllowedPilots {
            a.SetTimesToUTC()
        }
    }

    if p.Flights != nil {
        for _, f := range *p.Flights {
            f.SetTimesToUTC()
        }
    }
}
