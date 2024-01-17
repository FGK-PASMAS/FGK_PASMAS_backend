package testutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type PassengerModel struct {
	Id        int                       `json:"id"`
	LastName  string                    `json:"lastName"`
	FirstName string                    `json:"firstName"`
	Weight    int                       `json:"weight"`
	Division  DivisionModel             `json:"division"`
}


func ValidatePassengerModel(t *testing.T, pass PassengerModel) {
    assert.NotNil(t, pass.Id)
    assert.NotNil(t, pass.LastName)
    assert.GreaterOrEqual(t, len(pass.LastName), 1, "LastName should not be empty")
    assert.NotNil(t, pass.FirstName)
    assert.GreaterOrEqual(t, pass.Weight, 1, "Weight should not be empty")
    ValidateDivision(t, pass.Division)
}
