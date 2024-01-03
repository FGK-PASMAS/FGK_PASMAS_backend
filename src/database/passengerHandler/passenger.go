package passengerhandler

import (
	dh "github.com/MetaEMK/FGK_PASMAS_backend/database/divisionHandler"
)

// This type is used to create new passengers in the database
type InsertPassenger struct {
    LastName string         `json:"lastName"`
    FirstName string        `json:"firstName"`
    Weight int              `json:"weight" binding:"required"`
    DivisionId int          `json:"divisionId" binding:"required"`
}

// This type is used for for the response of any select from passenger table
type SelectPassenger struct {
    Id int                  `json:"id"`
    LastName string         `json:"lastName"`
    FirstName string        `json:"firstName"`
    Weight int              `json:"weight"`
    Division dh.Division    `json:"division"`
}

// This type is used to represent the passenger entity in the database
type DatabasePassenger struct {
    Id int                  `json:"id"`
    LastName string         `json:"lastName"`
    FirstName string        `json:"firstName"`
    Weight int              `json:"weight"`
    DivisionId int          `json:"divisionId"`
}

