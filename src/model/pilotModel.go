package model

import "gorm.io/gorm"

type Pilot struct {
    gorm.Model
    FirstName string `gorm:"uniqueIndex:idx_name"`
    LastName  string `gorm:"uniqueIndex:idx_name"`
    Weight    uint
    AllowedPilots *[]Plane `gorm:"many2many:AllowedPilots" json:"AllowedPlanes"`
}

func (p * Pilot) SetTimesToUTC() {
    p.CreatedAt = p.CreatedAt.UTC()
    p.UpdatedAt = p.UpdatedAt.UTC()

    if p.AllowedPilots != nil {
        for _, a := range *p.AllowedPilots {
            a.SetTimesToUTC()
        }
    }
}
