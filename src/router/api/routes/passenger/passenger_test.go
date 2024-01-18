package passenger_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	passengerhandler "github.com/MetaEMK/FGK_PASMAS_backend/database/passengerHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api"
	testutils "github.com/MetaEMK/FGK_PASMAS_backend/testUtils"
	"github.com/stretchr/testify/assert"
)

func TestGetPassengers(t *testing.T) {
    req, _ := http.NewRequest(http.MethodGet, "/api/passenger/", nil)
    pass1 := model.PassengerStructInsert{LastName: "test", FirstName: "test", Weight: 100, DivisionId: 1}
    pass2 := model.PassengerStructInsert{LastName: "test", FirstName: "", Weight: 100, DivisionId: 1}

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
    pass := model.PassengerStructInsert{LastName: "test", FirstName: "test", Weight: 100, DivisionId: 1}
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
    passUpdateError(
        t,
        model.PassengerStructInsert{LastName: "test", FirstName: "test", Weight: 100, DivisionId: 1},
        []byte(`{"lastName": "test", "firstName": "test", "weight": 100, "divisionId": 1}`),
        "/api/passenger/2",
        http.StatusNotFound,
        "OBJECT_NOT_FOUND",
    )
    passUpdateError(
        t,
        model.PassengerStructInsert{LastName: "test", FirstName: "test", Weight: 100, DivisionId: 1},
        []byte(`{"firstName": "test", "weight": 100, "divisionId": 1}`),
        "/api/passenger/1",
        http.StatusBadRequest,
        "INVALID_REQUEST_BODY",
    )
    passUpdateError(
        t,
        model.PassengerStructInsert{LastName: "test", FirstName: "test", Weight: 100, DivisionId: 1},
        []byte(`{"lastName": "test", "firstName": "test", "divisionId": 1}`),
        "/api/passenger/1",
        http.StatusBadRequest,
        "INVALID_REQUEST_BODY",
    )

    passUpdateError(
        t,
        model.PassengerStructInsert{LastName: "test", FirstName: "test", Weight: 100, DivisionId: 1},
        []byte(`{"lastName": "test", "firstName": "test", "weight": 100, "divisionId": 5}`),
        "/api/passenger/1",
        http.StatusBadRequest,
        "INVALID_OBJECT_DEPENDENCY",
    )
}

func passUpdateCorrect(t *testing.T) {
    pass := model.PassengerStructInsert{LastName: "test", FirstName: "test", Weight: 100, DivisionId: 1}
    passUpdate := model.PassengerStructUpdate{LastName: "test", FirstName: "test", Weight: 100, DivisionId: 1}
    passUpdateJson, _ := json.Marshal(passUpdate)

    req, _ := http.NewRequest(http.MethodPut, "/api/passenger/1", bytes.NewBuffer(passUpdateJson))
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

func passUpdateError(t *testing.T,
    passCreate model.PassengerStructInsert,
    passUpdateJson []byte,
    url string,
    expectedHttpStatusCode int,
    expectedErrorType string,
    ) {

    req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(passUpdateJson))
    req.Header.Set("Content-Type", "application/json")

    w := testutils.SendTestingRequest(t, req, func() {
        passengerhandler.CreatePassenger(passCreate)
    })

    assert.Equal(t, expectedHttpStatusCode, w.Code)

    res := api.ErrorResponse{}
    err := json.Unmarshal(w.Body.Bytes(), &res)
    assert.Nil(t, err)

    assert.Equal(t, res.Success, false)
    assert.Equal(t, expectedErrorType, res.Type, string(w.Body.Bytes()))
}

func passUpdateMissingBody(t *testing.T) {
    t.Skip("Not implemented")
    //TODO: implement
}

func TestDeletePassenger(t *testing.T) {
    t.Skip("Not implemented")
    //TODO: implement
}
