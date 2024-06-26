package databasehandler

import (
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/realtime"
)

func initFlight() {
	Db.AutoMigrate(&model.Flight{})
}

func GetFlights(include *FlightInclude, filter *FlightFilter) (flights []model.Flight, err error) {
	db := Db
	db = interpretFlightConfig(db, include, filter)

	err = db.Order("departure_time ASC").Find(&flights).Error

	for i := range flights {
		flights[i].SetTimesToUTC()
	}

	return
}

func GetFlightById(id uint, include *FlightInclude) (flight model.Flight, err error) {
	db := Db
	db = interpretFlightConfig(db, include, nil)

	db.First(&flight, id)
	flight.SetTimesToUTC()

	return flight, db.Error
}

func (dh *DatabaseHandler) CreateFlight(flight model.Flight) (newFlight model.Flight, err error) {
	flight.ID = 0
	flight.Passengers = nil
	err = dh.Db.Create(&flight).Error

	if err != nil {
		dh.Db.AddError(err)
		return
	}

	err = dh.Db.Preload("Plane").Preload("Passengers").Preload("Pilot").First(&newFlight, flight.ID).Error
	if err != nil {
		return
	}

	flight.SetTimesToUTC()

	dh.rt.AddEvent(realtime.FlightStream, realtime.CREATED, newFlight)
	plane := model.Plane{}

	dh.Db.First(&plane, flight.PlaneId)
	stream := realtime.GetFlightStreamForDivisionId(plane.DivisionId)
	dh.rt.AddEvent(stream, realtime.CREATED, newFlight)

	return
}

// partialUpdateFlight updates the newFlight with all set data from newFlight. 0 or "" values means that the field should be set to nil
func (dh *DatabaseHandler) PartialUpdateFlight(id uint, newFlightData model.Flight) (flight model.Flight, err error) {
	err = dh.Db.Preload("Plane").First(&flight, id).Error
	if err != nil {
		dh.Db.AddError(err)
		return
	}

	if newFlightData.Status == model.FsBooked && flight.Status == model.FsReserved {
		flight.Status = newFlightData.Status
	}

	if newFlightData.FlightNo != nil {
		flight.FlightNo = newFlightData.FlightNo
	}

	if newFlightData.Description != nil {
		if *newFlightData.Description == "" {
			flight.Description = nil
		} else {
			flight.Description = newFlightData.Description
		}
	}

	if newFlightData.FuelAtDeparture != nil {
		if *newFlightData.FuelAtDeparture == 0 {
			flight.FuelAtDeparture = nil
		} else {
			flight.FuelAtDeparture = newFlightData.FuelAtDeparture
		}
	}

	err = dh.Db.Updates(&flight).Error
	dh.Db.Preload("Plane").Preload("Passengers").Preload("Pilot").First(&flight, id)
	if err == nil {
		flight.SetTimesToUTC()
		dh.rt.AddEvent(realtime.FlightStream, realtime.UPDATED, flight)
		if flight.Plane != nil {
			stream := realtime.GetFlightStreamForDivisionId(flight.Plane.DivisionId)
			dh.rt.AddEvent(stream, realtime.UPDATED, flight)
		}
	}

	return
}

// DeleteFlight deletes the flight with the given id and all its passengers. It returns the deleted flight and all its passengers.
func (dh *DatabaseHandler) DeleteFlight(id uint) (flight model.Flight, passengers []model.Passenger, err error) {
	err = dh.Db.Preload("Plane").Preload("Passengers").First(&flight, id).Error
	if err != nil {
		return
	}
	passengers = *flight.Passengers

	err = dh.Db.Delete(&flight, id).Error
	if err != nil {
		return
	}
	dh.rt.AddEvent(realtime.FlightStream, realtime.DELETED, flight)
	if flight.Plane != nil {
		stream := realtime.GetFlightStreamForDivisionId(flight.Plane.DivisionId)
		dh.rt.AddEvent(stream, realtime.DELETED, flight)
	}

	for _, p := range passengers {
		dh.DeletePassenger(p.ID)
	}

	return
}
