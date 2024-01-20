package model

import (
	"gorm.io/gorm"
)

type Passenger struct {
    gorm.Model
    LastName            string
    FirstName           string
    Weight              uint            `gorm:"not null"`
}

