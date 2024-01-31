package validator_test

import (
	"testing"

	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/MetaEMK/FGK_PASMAS_backend/validator"
	"github.com/stretchr/testify/assert"
)

func TestValidatePassenger(t *testing.T) {
    pass := model.Passenger{
        LastName: "Test",
        FirstName: "Test",
        Weight: 42,
    }

    pass.LastName = ""
    err := validator.ValidatePassengerForReserve(pass)
    assert.Equal(t, validator.ErrPassengerLastName, err)
    pass.LastName = "Test"

    pass.Weight = 0
    err = validator.ValidatePassengerForReserve(pass)
    assert.Equal(t, validator.ErrPassengerWeight, err)
    pass.Weight = 42

    pass.FirstName = ""
    err = validator.ValidatePassengerForReserve(pass)
    assert.Nil(t, err)
    pass.FirstName = "Test"

    err = validator.ValidatePassengerForReserve(pass)
    assert.Nil(t, err)
}
