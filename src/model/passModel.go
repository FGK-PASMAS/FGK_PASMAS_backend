package model

type PassengerStructSelect struct {
	Id        int64                       `json:"id"`
	LastName  string                    `json:"lastName"`
	FirstName string                    `json:"firstName"`
	Weight    int                       `json:"weight"`
}

// This type is used to create new passengers in the database
type PassengerStructInsert struct {
	LastName   string `json:"lastName"`
	FirstName  string `json:"firstName"`
	Weight     int    `json:"weight" binding:"required"`
}

type PassengerStructUpdate struct {
	LastName   string `json:"lastName" binding:"required"`
	FirstName  string `json:"firstName"`
	Weight     int    `json:"weight" binding:"required"`
}

// This type is used to represent the passenger entity in the database
type DatabasePassenger struct {
	Id         int64    `json:"id"`
	LastName   string `json:"lastName"`
	FirstName  string `json:"firstName"`
	Weight     int    `json:"weight"`
}
