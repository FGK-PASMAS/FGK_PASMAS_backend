package pasmasservice_test

import (
	"testing"

	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	pasmasservice "github.com/MetaEMK/FGK_PASMAS_backend/service/pasmasService"
	"github.com/stretchr/testify/assert"
)

func comparePassengers(t *testing.T, expectedPass *model.Passenger, actualPass *model.Passenger) {
    assert.Equalf(t, expectedPass.ID, actualPass.ID, "ID %d is not equal to %d", expectedPass.ID, actualPass.ID)
    assert.Equalf(t, expectedPass.LastName, actualPass.LastName, "LastName %s is not equal to %s", expectedPass.LastName, actualPass.LastName)
    assert.Equalf(t, expectedPass.FirstName, actualPass.FirstName, "FirstName %s is not equal to %s", expectedPass.FirstName, actualPass.FirstName)
    assert.Equalf(t, expectedPass.Weight, actualPass.Weight, "Weight %d is not equal to %d", expectedPass.Weight, actualPass.Weight)
}

func TestGetPassengers(t *testing.T) {
    initDB(t)

    pass, err := pasmasservice.GetPassengers()
    if err != nil {
        assert.Nil(t, err)
        t.FailNow()
    }

    assert.IsType(t, []model.Passenger{}, pass)

    for _, p := range pass {
        assert.Nilf(t, p.DeletedAt, "Passenger %d is marked as deleted", p.ID)
    }
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

    pass.ID = newPass.ID
    comparePassengers(t, &pass, &newPass)
}

func TestUpdatePassenger(t *testing.T) {
    initDB(t)

    pass := model.Passenger {
        LastName: "TestUpdatePassenger",
        FirstName: "pasmasServiceTest",
        Weight: 42,
    }

    newPass, err := pasmasservice.CreatePassenger(pass)
    if err != nil {
        assert.Nil(t, err, "Could not create Passenger for updating")
        t.FailNow()
    }

    newPass.FirstName = "TestUpdatePassenger - Updated"
    newPass.LastName = "pasmasServiceTest - updated"
    newPass.Weight = 420

    updatedPass, err := pasmasservice.UpdatePassenger(newPass.ID ,newPass)
    if err != nil {
        assert.Nil(t, err)
        t.FailNow()
    }

    comparePassengers(t, &newPass, &updatedPass)
}

func TestDeletePassenger(t *testing.T) {
    initDB(t)

    pass := model.Passenger {
        LastName: "TestDeletePassenger",
        FirstName: "pasmasServiceTest",
        Weight: 42,
    }

    newPass, err := pasmasservice.CreatePassenger(pass)
    if err != nil {
        assert.Nil(t, err, "Could not create Passenger for deleting")
        t.FailNow()
    }

    err = pasmasservice.DeletePassenger(newPass.ID)
    assert.Nil(t, err, "Could not delete Passenger")
}
