package passengerhandler

import dh "github.com/MetaEMK/FGK_PASMAS_backend/database/divisionHandler"

// This type is used for for the response of any select from passenger table
type PassengerStructSelect struct {
	Id        int                       `json:"id"`
	LastName  string                    `json:"lastName"`
	FirstName string                    `json:"firstName"`
	Weight    int                       `json:"weight"`
	Division  dh.DivisionStructSelect   `json:"division"`
}

// This type is used to create new passengers in the database
type PassengerStructInsert struct {
	LastName   string `json:"lastName"`
	FirstName  string `json:"firstName"`
	Weight     int    `json:"weight" binding:"required"`
	DivisionId int    `json:"divisionId" binding:"required"`
}

type PassengerStructUpdate struct {
	Id         int    `json:"id" binding:"required"`
	LastName   string `json:"lastName" binding:"required"`
	FirstName  string `json:"firstName"`
	Weight     int    `json:"weight" binding:"required"`
	DivisionId int    `json:"divisionId" binding:"required"`
}

// This type is used to represent the passenger entity in the database
type DatabasePassenger struct {
	Id         int    `json:"id"`
	LastName   string `json:"lastName"`
	FirstName  string `json:"firstName"`
	Weight     int    `json:"weight"`
	DivisionId int    `json:"divisionId"`
}
