package testutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type DivisionModel struct {
    Id int                  `json:"id"`
    Name string             `json:"name"`
    PassengerCapacity int   `json:"passengerCapacity"`
}

func ValidateDivision(t *testing.T, div DivisionModel) {
    assert.NotNil(t, div.Id)
    assert.NotNil(t, div.Name)
    assert.GreaterOrEqualf(t, len(div.Name), 1, "Division name is empty")
    assert.GreaterOrEqual(t, div.PassengerCapacity, 1)
}
