package passenger_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/MetaEMK/FGK_PASMAS_backend/database/debug"
	passengerhandler "github.com/MetaEMK/FGK_PASMAS_backend/database/passengerHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	testutils "github.com/MetaEMK/FGK_PASMAS_backend/testUtils"
	"github.com/stretchr/testify/assert"
)

var endpoint = "/api/passengers"

func TestGetPassengers(t *testing.T) {
    env := testutils.InitRouter(true)
    req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
    pass1 := testutils.CreateDummyPassengerCreate()
    pass2 := testutils.CreateDummyPassengerCreate()

    res := env.SendTestingRequestSuccess (
        t,
        req,
        func() {
            passengerhandler.CreatePassenger(pass1)
            passengerhandler.CreatePassenger(pass2)
        },
        http.StatusOK,
        true,
    )

    var passenger []testutils.PassengerModel
    jsonBytes, _ := json.Marshal(res.Response)

    err := json.Unmarshal(jsonBytes, &passenger)
    assert.Nil(t, err)

    for _, pass := range passenger {
        testutils.ValidatePassengerModel(t, pass)
    }
}

func TestCreatePassenger(t *testing.T) {
    env := testutils.InitRouter(true)
    pass := testutils.CreateDummyPassengerCreate()
    passJson, _ := json.Marshal(pass)

    req, _ := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(passJson))
    req.Header.Set("Content-Type", "application/json")

    res := env.SendTestingRequestSuccess (
        t,
        req,
        func() {},
        http.StatusCreated,
        true,
    )

    // Depenendy on divsion should fail
    var passenger testutils.PassengerModel
    jsonBytes, _ := json.Marshal(res.Response)
    err := json.Unmarshal(jsonBytes, &passenger)
    assert.Nil(t, err)
    testutils.ValidatePassengerModel(t, passenger)

    // Body validation should fail
    passWrongBody := []byte(`{"lastName": "test", "firstName": "test", "divisionId": 1}`)
    reqWrongBody, _ := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(passWrongBody))
    reqWrongBody.Header.Set("Content-Type", "application/json")

    env.SendTestingRequestError (
        t,
        reqWrongBody,
        func() {},
        http.StatusBadRequest,
        "INVALID_REQUEST_BODY",
    )
}

func TestUpdatePassenger(t *testing.T) {
    env := testutils.InitRouter(true)
    passUpdateCorrect(t, &env)
    passUpdateError(
        t,
        &env,
        testutils.CreateDummyPassengerCreate(),
        []byte(`{"lastName": "test", "firstName": "test", "weight": 100}`),
        endpoint + "/2",
        http.StatusNotFound,
        "OBJECT_NOT_FOUND",
    )
    passUpdateError(
        t,
        &env,
        testutils.CreateDummyPassengerCreate(),
        []byte(`{"firstName": "test", "weight": 100}`),
        endpoint + "/1",
        http.StatusBadRequest,
        "INVALID_REQUEST_BODY",
    )
    passUpdateError(
        t,
        &env,
        testutils.CreateDummyPassengerCreate(),
        []byte(`{"lastName": "test", "firstName": "test"}`),
        endpoint + "/1",
        http.StatusBadRequest,
        "INVALID_REQUEST_BODY",
    )
}

func passUpdateCorrect(t *testing.T, env *testutils.TestEnv) {
    pass := testutils.CreateDummyPassengerCreate()
    passUpdate := testutils.DummyUpdatePassenger()
    passUpdateJson, _ := json.Marshal(passUpdate)

    req, _ := http.NewRequest(http.MethodPut, endpoint + "/1", bytes.NewBuffer(passUpdateJson))
    req.Header.Set("Content-Type", "application/json")

    res := env.SendTestingRequestSuccess (
        t,
        req,
        func() {
            passengerhandler.CreatePassenger(pass)
        },
        http.StatusOK,
        true,
    )

    var passenger testutils.PassengerModel
    jsonBytes, _ := json.Marshal(res.Response)

    err := json.Unmarshal(jsonBytes, &passenger)
    assert.Nil(t, err)

    testutils.ValidatePassengerModel(t, passenger)
}

func passUpdateError(
    t *testing.T,
    env *testutils.TestEnv,
    passCreate model.PassengerStructInsert,
    passUpdateJson []byte,
    url string,
    expectedHttpStatusCode int,
    expectedErrorType string,
    ) {

    req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(passUpdateJson))
    req.Header.Set("Content-Type", "application/json")

    env.SendTestingRequestError (
        t,
        req,
        func() {
            debug.TruncateDatabase()
            passengerhandler.CreatePassenger(passCreate)
        },
        expectedHttpStatusCode,
        expectedErrorType,
    )
}

func TestDeletePassenger(t *testing.T) {
    env := testutils.InitRouter(true)
    pass := testutils.CreateDummyPassengerCreate()

    req, _ := http.NewRequest(http.MethodDelete, endpoint + "/1", nil)

    env.SendTestingRequestError (
        t,
        req,
        func() {},
        http.StatusNotFound,
        "OBJECT_NOT_FOUND",
    )

    env.SendTestingRequestSuccess(
        t, 
        req,
        func() {
            passengerhandler.CreatePassenger(pass)
        },
        http.StatusNoContent,
        false,
    )

}
