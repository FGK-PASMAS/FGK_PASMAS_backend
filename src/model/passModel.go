package model

type PassengerStructSelect struct {
	Id        int64                       `json:"id"`
	LastName  string                    `json:"lastName"`
	FirstName string                    `json:"firstName"`
	Weight    int                       `json:"weight"`
	Division  DivisionStructSelect      `json:"division"`
}

// This type is used to create new passengers in the database
type PassengerStructInsert struct {
	LastName   string `json:"lastName"`
	FirstName  string `json:"firstName"`
	Weight     int    `json:"weight" binding:"required"`
	DivisionId int    `json:"divisionId" binding:"required"`
}

type PassengerStructUpdate struct {
	Id         int64    `json:"id" binding:"required"`
	LastName   string `json:"lastName" binding:"required"`
	FirstName  string `json:"firstName"`
	Weight     int    `json:"weight" binding:"required"`
	DivisionId int    `json:"divisionId" binding:"required"`
}

// This type is used to represent the passenger entity in the database
type DatabasePassenger struct {
	Id         int64    `json:"id"`
	LastName   string `json:"lastName"`
	FirstName  string `json:"firstName"`
	Weight     int    `json:"weight"`
	DivisionId int    `json:"divisionId"`
}
