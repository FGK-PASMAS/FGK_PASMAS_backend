package model

import (
	"gorm.io/gorm"
)

type Passenger struct {
    gorm.Model
    LastName            string
    FirstName           string
    Weight              uint            `gorm:"not null"`

    FlightID            *uint           `gorm:"index"`
    Flight              *Flight          `gorm:"foreignKey:FlightID;OnUpdate:CASCADE;OnDelete:RESTRICT"`
}

