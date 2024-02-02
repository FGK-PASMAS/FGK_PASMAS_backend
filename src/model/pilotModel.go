package model

import "gorm.io/gorm"

type Pilot struct {
    gorm.Model
    FirstName string `gorm:"uniqueIndex:idx_name"`
    LastName  string `gorm:"uniqueIndex:idx_name"`
    Weight    uint
    AllowedPilots *[]Plane `gorm:"many2many:AllowedPilots" json:"AllowedPlanes"`
}

