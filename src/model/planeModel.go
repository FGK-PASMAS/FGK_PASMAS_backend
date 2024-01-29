package model

import "gorm.io/gorm"

type Plane struct {
	gorm.Model

	// Aricraft Registration f.E. D-ELHN
    Registration string `gorm:"unique;not null"`

	// Aircraft Type f.E. C172
	AircraftType string `gorm:"not null"`

	// Maximum amount of fuel this plane can take in liters; -1 if not applicable
	FuelMaxCapacity int `gorm:"not null"`

	// Fuel consumption per flight in liters; -1 if not applicable
	FuelburnPerFlight float32 `gorm:"not null"`

	// Conversion factor of fuel in kg/l; -1 if not applicable
	FuelConversionFactor float32 `gorm:"not null"`

	// Maximum takeoff weight in kg
	MTOW int `gorm:"not null"`

	// Empty weight of the aircraft in kg
	EmptyWeight int `gorm:"not null"`

    // Maximum weight of one seat in kg; -1 if not applicable
    MaxSeatPayload int `gorm:"not null"`

	// Aircrafts division f.E. "Motorflug"
	DivisionId uint     `json:"-" gorm:"index"`
	Division   *Division `gorm:"foreignKey:DivisionId;OnUpdate:CASCADE;OnDelete:RESTRICT"`

    AllowedPilots *[]Pilot `gorm:"many2many:AllowedPilots;"`
    //PrefPilot *Pilot `gorm:"foreignKey:PrefPilot;OnUpdate:CASCADE;OnDelete:RESTRICT"`

    Flights *[]Flight `gorm:"foreignKey:PlaneId"`
}
