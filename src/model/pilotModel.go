package model

import "gorm.io/gorm"

type Pilot struct {
    gorm.Model
    FirstName string 
    LastName string
    Weight int
    AllowedPlanes *[]Plane `gorm:"many2many:AllowedPlanes;"`
}
