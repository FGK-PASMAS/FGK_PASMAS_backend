package flightlogic

import (
	"time"

	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	databasehandler "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/MetaEMK/FGK_PASMAS_backend/validator"
	"gorm.io/gorm"
)

func FlightLogicProcess(flight model.Flight, plane model.Plane, division model.Division, checkSlot bool) (newFlight model.Flight, err error) {
    var passengers []model.Passenger
    var fullValidation bool = false
    var minPass = 0

    if flight.DepartureTime.IsZero() {
        err = cerror.ErrDepartureTimeIsZero
        return
    }

    if flight.ArrivalTime.IsZero() {
        flight.ArrivalTime = flight.DepartureTime.Add(plane.FlightDuration).UTC()
    }

    if flight.ArrivalTime.Before(flight.DepartureTime) {
        err = cerror.ErrInvalidArrivalTime
        return
    }

    if flight.Passengers != nil {
        passengers = *flight.Passengers
    } else {
        passengers = make([]model.Passenger, 0)
    }

    if flight.Status == model.FsBooked {
        fullValidation = true
        minPass = 1
    }


    passWeight := CalcPassWeight(passengers)
    err = CheckPassenger(passengers, plane.MaxSeatPayload, uint(minPass), division.PassengerCapacity, fullValidation)
    if err != nil {
        return
    }

    var fuelAmount float32 = 0
    if plane.FuelburnPerFlight > 0 {
        fuelAmount, err = calculateFuelAtDeparture(flight, plane)
        if err != nil {
            return
        }
    }

    pilot, err := calculatePilot(passWeight, fuelAmount, plane)
    if err != nil {
        return
    }
    flight.PilotId = &pilot.ID
    flight.Pilot = &pilot

    if checkSlot {
        if CheckIfSlotIsFree(plane.ID, flight.DepartureTime, flight.ArrivalTime) == false {
            err = cerror.ErrSlotIsNotFree
            return
        }
    }

    newFlight = flight

    return
}

// checkIfSlotIsFree checks if the slot is free for the given planeid, departureTime and arrivalTime
//
// returns true if the slot is free, false if not
func CheckIfSlotIsFree(planeId uint, departureTime time.Time, arrivalTime time.Time) bool {
    var count int64
    result := databasehandler.Db.Model(&model.Flight{}).Where("plane_id = ?", planeId).Where("departure_time < ? AND arrival_time > ?", arrivalTime, departureTime).Count(&count)

    if result.Error != nil {
        println(result.Error.Error())
        return false
    }

    return count == 0
}

func checkIfTimeSlotIsValid(plane model.Plane,departureTime time.Time) error {
    if departureTime.After(plane.SlotStartTime) {
        if departureTime.Before(plane.SlotEndTime) {
            return nil
        }
    }

    return cerror.ErrTimeSlotForPlaneNotValid
}

func calculatePilot(passWeight uint, fuelAmount float32, plane model.Plane) (model.Pilot, error) {
    var baseETOW uint = 0
    pilot := model.Pilot{}

    //err := dh.Db.Preload("AllowedPilots").Preload("PrefPilot").First(&plane).Error
    plane, err := databasehandler.GetPlaneById(plane.ID, &databasehandler.PlaneInclude{IncludeAllowedPilots: true, IncludePrefPilot: true})
    if err != nil {
        return model.Pilot{}, err
    }

    if plane.PrefPilot == nil {
        if len(*plane.AllowedPilots) > 0 {
            pilot = (*plane.AllowedPilots)[0]
        } else {
            return model.Pilot{}, cerror.ErrNoPilotAvailable
        }

    } else {
        pilot = *plane.PrefPilot
    }


    baseETOW += passWeight
    baseETOW += plane.EmptyWeight
    baseETOW += uint(fuelAmount * plane.FuelConversionFactor)

    if plane.MTOW < baseETOW + pilot.Weight {
        newPilot := model.Pilot{}

        if len(*plane.AllowedPilots) == 0 {
            return model.Pilot{}, cerror.ErrNoPilotAvailable
        }

        for _, p := range *plane.AllowedPilots {
            if plane.MTOW >= baseETOW + p.Weight {
                newPilot = p
                break
            }
        }

        if newPilot.ID == 0 {
            return model.Pilot{}, cerror.ErrOverloaded
        }

        pilot = newPilot
    }

    return pilot, err
}

func CalcPassWeight(passengers []model.Passenger) uint {
    var passWeight uint = 0

    for _, p := range passengers {
        passWeight += p.Weight
    }

    return passWeight
}

func CheckPassenger(passengers []model.Passenger, maxSeatPayload int, min uint, max uint, fullPassCheck bool) error {
    if len(passengers) > int(max) {
        return cerror.ErrTooManyPassenger
    }

    if len(passengers) < int(min) {
        return cerror.ErrTooLessPassenger
    }

    for _, p := range passengers {
        if maxSeatPayload > 0 && p.Weight > uint(maxSeatPayload){
            return cerror.ErrMaxSeatPayload
        }

        if fullPassCheck {
            err := validator.ValidatePassengerForBooking(p)
            if err != nil {
                return err
            }
        }

    }

    return nil
}

func calculateFuelAtDeparture(flight model.Flight, plane model.Plane) (float32, error) {
    if flight.FuelAtDeparture != nil && *flight.FuelAtDeparture != 0 {
        if *flight.FuelAtDeparture > float32(plane.FuelMaxCapacity) {
            return 0, cerror.ErrTooMuchFuel
        }
        return *flight.FuelAtDeparture, nil
    }

    // Get one flight before this
    beforeFlight := model.Flight{}
    err := databasehandler.Db.Not("status = ?", model.FsBlocked) .Where("plane_id = ?", flight.PlaneId) .Where("departure_time < ?", flight.DepartureTime) .Order("departure_time DESC").First(&beforeFlight).Error

    if err == gorm.ErrRecordNotFound {
        fuel := float32(plane.FuelStartAmount)
        flight.FuelAtDeparture = &fuel
        return float32(plane.FuelStartAmount), nil
    }

    value, err := calculateFuelAtDeparture(beforeFlight, plane)
    value -= plane.FuelburnPerFlight

    if value <= 0 {
        return 0, cerror.ErrTooLessFuel
    }

    return value, nil
}

