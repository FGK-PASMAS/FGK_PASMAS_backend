package pasmasservice_test

import (
	"testing"
	"time"

	dh "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	pasmasservice "github.com/MetaEMK/FGK_PASMAS_backend/service/pasmasService"
	"github.com/stretchr/testify/assert"
)

func TestReserveFlight(t *testing.T) {
    t.Skip()
}

func TestGetFlights(t *testing.T) {
    initDB(t)

    // Preparing Databases for test cases
    flight := model.Flight {
        Type: model.FsBooked,
        FuelAtDeparture: 69,
        DepartureTime: time.Now(),
        ArrivalTime: time.Now().Add(20 * time.Minute),
    }
    dh.Db.Create(&flight)

    flight = model.Flight {
        Type: model.FsReserved,
        FuelAtDeparture: 42,
        DepartureTime: time.Now().Add(40 * time.Minute),
        ArrivalTime: time.Now().Add(60 * time.Minute),
    }
    dh.Db.Create(&flight)

    flight = model.Flight {
        Type: model.FsBlocked,
        FuelAtDeparture: 42,
        DepartureTime: time.Now().Add(time.Hour),
        ArrivalTime: time.Now().Add(20 * time.Minute + time.Hour),
    }
    dh.Db.Create(&flight)
    dh.Db.Delete(&flight)

    flights, err := pasmasservice.GetFlights()

    if err != nil {
        assert.FailNowf(t, "Error while getting flights: %v", err.Error())
    }

    assert.Equal(t, 2, len(*flights), "There should be 2 flights in database")
}

func TestCheckIfSlotIsFree(t *testing.T) {
    initDB(t)

    // Preparing Databases for test cases
    timeNow := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)

    flight := model.Flight{
        Type: model.FsBooked,
        FuelAtDeparture: 69,
        DepartureTime: timeNow,                         // 12:00
        ArrivalTime: timeNow.Add(19 * time.Minute),     // 12:20
    }
    dh.Db.Create(&flight)

    flight = model.Flight{
        Type: model.FsBooked,
        FuelAtDeparture: 42,
        DepartureTime: timeNow.Add(40 * time.Minute),   // 12:40
        ArrivalTime: timeNow.Add(59 * time.Minute),     // 13:00
    }
    dh.Db.Create(&flight)

    // Flight Validation
    // Too early
    checkFlight := model.Flight{
        Type: model.FsReserved,
        FuelAtDeparture: 42,
        DepartureTime: timeNow.Add(10 * time.Minute),   // 12:10
        ArrivalTime: timeNow.Add(29 * time.Minute),     // 12:30
    }
    res := pasmasservice.CheckIfSlotIsFree(&checkFlight)
    assert.Equalf(t, false, res, "This flight should be to early %v", checkFlight)

    // Too late
    checkFlight.DepartureTime = timeNow.Add(30 * time.Minute)   //12:30
    checkFlight.ArrivalTime = timeNow.Add(49 * time.Minute)     //12:50
    res = pasmasservice.CheckIfSlotIsFree(&checkFlight)
    assert.Equalf(t, false, res, "This flight should be to late %v", checkFlight)

    // should not fail
    checkFlight.DepartureTime = timeNow.Add(20 * time.Minute)
    checkFlight.ArrivalTime = timeNow.Add(39 * time.Minute)
    res = pasmasservice.CheckIfSlotIsFree(&checkFlight)
    assert.Equalf(t, true, res, "This flight is perfectly scheduled: %v", checkFlight)

    checkFlight.DepartureTime = timeNow.Add(time.Hour)
    checkFlight.ArrivalTime = timeNow.Add(19 * time.Minute + time.Hour)
    res = pasmasservice.CheckIfSlotIsFree(&checkFlight)
    assert.Equalf(t, true, res, "This flight is perfectly scheduled %v", checkFlight)

    checkFlight.DepartureTime = timeNow.Add(2 *time.Hour)
    checkFlight.ArrivalTime = timeNow.Add(19 * time.Minute + 2 * time.Hour)
    res = pasmasservice.CheckIfSlotIsFree(&checkFlight)
    assert.Equalf(t, true, res, "This flight is scheduled correctly %v", checkFlight)

}
