package model

import (
	"gorm.io/gorm"
)

type Division struct {
    gorm.Model
    Name                string
    PassengerCapacity   uint
}
