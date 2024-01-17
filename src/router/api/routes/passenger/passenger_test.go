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

    assert.Equal(t, http.StatusOK, w.Code)

    res := testutils.ParseAndValidateResponse(t, w)
    assert.Equal(t, true, res.Success)

    var passenger []testutils.PassengerModel
    jsonBytes, _ := json.Marshal(res.Response)

    err := json.Unmarshal(jsonBytes, &passenger)
    assert.Nil(t, err)

    for _, pass := range passenger {
        testutils.ValidatePassengerModel(t, pass)
    }
}

func TestCreatePassenger(t *testing.T) {
    pass := passengerhandler.PassengerStructInsert{LastName: "test", FirstName: "test", Weight: 100, DivisionId: 1}
    passJson, _ := json.Marshal(pass)

    req, _ := http.NewRequest(http.MethodPost, "/api/passenger/", bytes.NewBuffer(passJson))
    req.Header.Set("Content-Type", "application/json")

    w := testutils.SendTestingRequest(t, req)

    assert.Equal(t, http.StatusCreated, w.Code)

    res := testutils.ParseAndValidateResponse(t, w)
    assert.Equal(t, true, res.Success)

    var passenger testutils.PassengerModel
    jsonBytes, _ := json.Marshal(res.Response)

    err := json.Unmarshal(jsonBytes, &passenger)
    assert.Nil(t, err)

    testutils.ValidatePassengerModel(t, passenger)
}

func TestUpdatePassenger(t *testing.T) {
    passUpdateCorrect(t)
    passUpdateWrongId(t)
}

func passUpdateCorrect(t *testing.T) {
    pass := passengerhandler.PassengerStructInsert{LastName: "test", FirstName: "test", Weight: 100, DivisionId: 1}
    passUpdate := passengerhandler.PassengerStructUpdate{Id: 1, LastName: "test", FirstName: "test", Weight: 100, DivisionId: 1}
    passUpdateJson, _ := json.Marshal(passUpdate)

    req, _ := http.NewRequest(http.MethodPut, "/api/passenger/", bytes.NewBuffer(passUpdateJson))
    req.Header.Set("Content-Type", "application/json")

    w := testutils.SendTestingRequest(t, req, func() {
        passengerhandler.CreatePassenger(pass)
    })

    assert.Equal(t, http.StatusOK, w.Code)

    res := testutils.ParseAndValidateResponse(t, w)
    assert.Equal(t, true, res.Success)

    var passenger testutils.PassengerModel
    jsonBytes, _ := json.Marshal(res.Response)

    err := json.Unmarshal(jsonBytes, &passenger)
    assert.Nil(t, err)

    testutils.ValidatePassengerModel(t, passenger)
}

func passUpdateWrongId(t *testing.T) {
    //TODO: improve implementation
    passUpdate := passengerhandler.PassengerStructUpdate{Id: 1, LastName: "test", FirstName: "test", Weight: 100, DivisionId: 1}
    passUpdateJson, _ := json.Marshal(passUpdate)

    req, _ := http.NewRequest(http.MethodPut, "/api/passenger/", bytes.NewBuffer(passUpdateJson))
    req.Header.Set("Content-Type", "application/json")

    w := testutils.SendTestingRequest(t, req)

    assert.Equal(t, http.StatusInternalServerError, w.Code)

    res := testutils.ParseAndValidateResponse(t, w)
    assert.Equal(t, res.Success, false)
}

func passUpdateMissingBody(t *testing.T) {
    t.Skip("Not implemented")
    //TODO: implement
}

func TestDeletePassenger(t *testing.T) {
    t.Skip("Not implemented")
    //TODO: implement
}
