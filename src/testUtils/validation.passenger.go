package testutils

import (
	"testing"

	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/stretchr/testify/assert"
)

type PassengerModel struct {
	Id        int                       `json:"id"`
	LastName  string                    `json:"lastName"`
	FirstName string                    `json:"firstName"`
	Weight    int                       `json:"weight"`
}


func CreateDummyPassengerCreate() model.PassengerStructInsert {
    return model.PassengerStructInsert{
        LastName: "test",
        FirstName: "test",
        Weight: 100,
    }
}

func DummyUpdatePassenger() model.PassengerStructUpdate {
    return model.PassengerStructUpdate{
        LastName: "test",
        FirstName: "test",
        Weight: 100,
    }
}


func ValidatePassengerModel(t *testing.T, pass PassengerModel) {
    assert.NotNil(t, pass.Id)
    assert.NotNil(t, pass.LastName)
    assert.GreaterOrEqual(t, len(pass.LastName), 1, "LastName should not be empty")
    assert.NotNil(t, pass.FirstName)
    assert.GreaterOrEqual(t, pass.Weight, 1, "Weight should not be empty")
}
