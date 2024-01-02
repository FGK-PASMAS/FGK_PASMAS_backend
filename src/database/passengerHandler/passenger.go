package passengerhandler

import dh "github.com/MetaEMK/FGK_PASMAS_backend/database/divisionHandler"

type Passenger struct {
    Id int                  `json:"id"`
    LastName string         `json:"lastName"`
    FirstName string        `json:"firstName"`
    Weight int              `json:"weight"`
    Division dh.Division    `json:"division"`
}
