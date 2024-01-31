package pasmasservice

import (
	"errors"
	"time"

	dh "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"gorm.io/gorm"
)

var (
    ErrNoPilotAvailable = errors.New("No valid pilot available")
    ErrNoStartFuelFound = errors.New("No start fuel found")
    ErrMaxSeatPayload = errors.New("maxSeatPayload was exceeded")
    ErrOverloaded = errors.New("MTOW is exceeded")
)

func checkIfSlotIsFree(flight *model.Flight) bool {
    flights := []model.Flight{}
    arr_time := flight.ArrivalTime.Truncate(time.Minute).Local()
    dep_time := flight.DepartureTime.Truncate(time.Minute).Local()
    result := dh.Db.Where("plane_id = ?", flight.PlaneId).Where("arrival_time >= ?", dep_time).Where("departure_time <= ?", arr_time).Find(&flights)

    if result.Error != nil {
        return false
    }

    if len(flights) == 0 {
        return true
    }

    return false
}

func calculatePilot(passWeight uint, fuelAmount float32, plane model.Plane) (model.Pilot, error) {
    var baseETOW uint = 0
    pilot := model.Pilot{}

    err := dh.Db.Preload("AllowedPilots").Preload("PrefPilot").First(&plane).Error
    if err != nil {
        return model.Pilot{}, err
    }

    if plane.PrefPilot == nil {
        if len(*plane.AllowedPilots) > 0 {
            pilot = (*plane.AllowedPilots)[0]
        } else {
            return model.Pilot{}, ErrNoPilotAvailable
        }

    } else {
        pilot = *plane.PrefPilot
    }

    println(pilot.ID)

    baseETOW += passWeight
    baseETOW += plane.EmptyWeight
    baseETOW += uint(fuelAmount * plane.FuelConversionFactor)

    if plane.MTOW < baseETOW + pilot.Weight {
        newPilot := model.Pilot{}

        if len(*plane.AllowedPilots) == 0 {
            return model.Pilot{}, ErrNoPilotAvailable
        }

        for _, p := range *plane.AllowedPilots {
            if plane.MTOW >= baseETOW + p.Weight {
                newPilot = p
                break
            }
        }

        if newPilot.ID == 0 {
            return model.Pilot{}, ErrOverloaded
        }

        pilot = newPilot
    }

    return pilot, err
}

func checkFlightValidation(flight model.Flight) error {
    var err error
    plane := model.Plane{}
    pilot := model.Pilot{}

    planeErr := dh.Db.First(&plane, flight.PlaneId).Error
    pilotErr := dh.Db.First(&pilot, flight.PilotId).Error

    err = errors.Join(err, planeErr, pilotErr)
    if err != nil {
        return err
    }

    // Validate if flight is overweight
    var etow float32 = 0
    etow += float32(plane.EmptyWeight)
    etow += flight.FuelAtDeparture * plane.FuelConversionFactor
    etow += float32(pilot.Weight)

    for _, p := range *flight.Passengers {
        if plane.MaxSeatPayload > 0 {
            if p.Weight > uint(plane.MaxSeatPayload) {
                err = errors.Join(err, ErrMaxSeatPayload)
            }
        }
        etow += float32(p.Weight)
    }

    return err
}


func calculatePassWeight(passengers []model.Passenger, maxSeatPayload int) (uint, error) {
    weight := uint(0)
    for _, p := range passengers {
        if maxSeatPayload > 0 && p.Weight > uint(maxSeatPayload){
            return 0, ErrMaxSeatPayload
        }
        weight += p.Weight
    }

    return weight, nil
}

func calculateFuelAtDeparture(flight model.Flight, plane model.Plane) (float32, error) {
    if flight.FuelAtDeparture != 0 {
        if flight.FuelAtDeparture > float32(plane.FuelMaxCapacity) {
            return 0, ErrTooMuchFuel
        }
        return flight.FuelAtDeparture, nil
    }

    // Get one flight before this
    beforeFlight := model.Flight{}
    err := dh.Db.Not("status = ?", model.FsBlocked).Where("plane_id = ?", flight.PlaneId).Where("departure_time < ?", flight.DepartureTime).Order("departure_time DESC").First(&beforeFlight).Error
    if err == gorm.ErrRecordNotFound {
        return float32(plane.FuelStartAmount), nil
    }

    value, err := calculateFuelAtDeparture(beforeFlight, plane)
    value -= plane.FuelburnPerFlight

    if value <= 0 {
        return 0, ErrTooLessFuel
    }

    return value, nil
}
