package pasmasservice_test

import (
	"reflect"
	"testing"

	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	pasmasservice "github.com/MetaEMK/FGK_PASMAS_backend/service/pasmasService"
	"github.com/stretchr/testify/assert"
)

func TestGetPassengers(t *testing.T) {
    initDB(t)

    pass, err := pasmasservice.GetPassengers()
    if err != nil {
        assert.Nil(t, err)
        t.FailNow()
    }

    assert.IsType(t, []model.Passenger{}, pass)
}

func TestCreatePassenger(t *testing.T) {
    initDB(t)

    pass := model.Passenger{
        LastName: "TestCreatePassenger",
        FirstName: "pasmasServiceTest",
        Weight: 42,
    }

    newPass, err := pasmasservice.CreatePassenger(pass)
    if err != nil {
        assert.Nil(t, err)
        t.FailNow()
    }

    assert.IsTypef(t, model.Passenger{}, newPass, "Expected: %s\n Actual: %s\n", reflect.TypeOf(model.Passenger{}), reflect.TypeOf(newPass))
    assert.GreaterOrEqual(t, 1, newPass.ID)
}
