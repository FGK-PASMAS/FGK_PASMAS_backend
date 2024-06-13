package model

import (
	"gorm.io/gorm"
)

type Division struct {
	gorm.Model
	Name              string
	PassengerCapacity uint
	Planes            []Plane `gorm:"foreignKey:DivisionId"`
}

func (d *Division) SetTimesToUTC() {
	d.CreatedAt = d.CreatedAt.UTC()
	d.UpdatedAt = d.UpdatedAt.UTC()

	for _, p := range d.Planes {
		p.SetTimesToUTC()
	}
}
