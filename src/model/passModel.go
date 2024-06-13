package model

import (
	"gorm.io/gorm"
)

type Passenger struct {
    gorm.Model
    LastName            string
    FirstName           string
    Weight              uint            `gorm:"not null"`
    PassNo              uint            `gorm:"uniqueIndex"`

    FlightID            uint            `gorm:"index"`
    Flight              *Flight         `gorm:"foreignKey:FlightID"`

    // This virtual field is used in the api to determine what action to take on this passenger
    //Action              Action          `gorm:"-"`
}

type Action string

const (
    ActionCreate Action = "CREATE"
    ActionUpdate Action = "UPDATE"
    ActionDelete Action = "DELETE"
)

func (p * Passenger) SetTimesToUTC() {
    p.CreatedAt = p.CreatedAt.UTC()
    p.UpdatedAt = p.UpdatedAt.UTC()

    if p.Flight != nil {
        p.Flight.SetTimesToUTC()
    }
}
