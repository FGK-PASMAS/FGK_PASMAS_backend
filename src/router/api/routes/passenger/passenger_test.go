package passenger_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	passengerhandler "github.com/MetaEMK/FGK_PASMAS_backend/database/passengerHandler"
	testutils "github.com/MetaEMK/FGK_PASMAS_backend/testUtils"
	"github.com/stretchr/testify/assert"
)

func TestGetPassengers(t *testing.T) {
    req, _ := http.NewRequest(http.MethodGet, "/api/passenger/", nil)
    pass1 := passengerhandler.PassengerStructInsert{LastName: "test", FirstName: "test", Weight: 100, DivisionId: 1}
    pass2 := passengerhandler.PassengerStructInsert{LastName: "test", FirstName: "", Weight: 100, DivisionId: 1}

    w := testutils.SendTestingRequest(t, req, func() {
        passengerhandler.CreatePassenger(pass1)
        passengerhandler.CreatePassenger(pass2)
    })

    assert.Equal(t, w.Code, http.StatusOK)

    res := testutils.ParseAndValidateResponse(t, w)
    assert.Equal(t, res.Success, true)

    var passenger []testutils.PassengerModel
    jsonBytes, _ := json.Marshal(res.Response)

    err := json.Unmarshal(jsonBytes, &passenger)
    assert.Nil(t, err)

    for _, pass := range passenger {
        testutils.ValidatePassengerModel(t, pass)
    }
}

func CreatePassenger(t *testing.T) {
    pass := passengerhandler.PassengerStructInsert{LastName: "test", FirstName: "test", Weight: 100, DivisionId: 1}
    passJson, _ := json.Marshal(pass)

    req, _ := http.NewRequest(http.MethodPost, "/api/passenger/", bytes.NewBuffer(passJson))
    req.Header.Set("Content-Type", "application/json")

    w := testutils.SendTestingRequest(t, req)

    assert.Equal(t, w.Code, http.StatusCreated)

    res := testutils.ParseAndValidateResponse(t, w)
    assert.Equal(t, res.Success, true)

    var passenger testutils.PassengerModel
    jsonBytes, _ := json.Marshal(res.Response)

    err := json.Unmarshal(jsonBytes, &passenger)
    assert.Nil(t, err)

    testutils.ValidatePassengerModel(t, passenger)
}

func UpdatePassenger(t *testing.T) {
    pass := passengerhandler.PassengerStructUpdate{LastName: "test", FirstName: "test", Weight: 100, DivisionId: 1}
    passJson, _ := json.Marshal(pass)

    req, _ := http.NewRequest(http.MethodPut, "/api/passenger/", bytes.NewBuffer(passJson))
    req.Header.Set("Content-Type", "application/json")

    w := testutils.SendTestingRequest(t, req)

    assert.Equal(t, w.Code, http.StatusOK)

    res := testutils.ParseAndValidateResponse(t, w)
    assert.Equal(t, res.Success, true)

    var passenger testutils.PassengerModel
    jsonBytes, _ := json.Marshal(res.Response)

    err := json.Unmarshal(jsonBytes, &passenger)
    assert.Nil(t, err)

    testutils.ValidatePassengerModel(t, passenger)
}

func DeletePassenger(t *testing.T) {
    t.Skip("Not implemented")
}
