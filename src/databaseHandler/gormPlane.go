package databasehandler

import (
	"time"

	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	"github.com/MetaEMK/FGK_PASMAS_backend/config"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/realtime"
	"gorm.io/gorm"
)

func initPlane(db *gorm.DB) {
	Db.AutoMigrate(&model.Plane{})

	if config.EnableSeeder {
        log.Debug("Seeding planes")
		SeedPlane(db)
	}
}

func GetPlanes(planeInclude *PlaneInclude, planeFilter *PlaneFilter) ([]model.Plane, error) {
	db := Db
	planes := []model.Plane{}

	db = interpretPlaneConfig(db, planeInclude, planeFilter)

	db = db.Order("id ASC").Find(&planes)

	for i := range planes {
		planes[i].SetTimesToUTC()
	}

	return planes, db.Error
}

func GetPlaneById(id uint, planeInclude *PlaneInclude) (plane model.Plane, err error) {
	db := Db

	db = interpretPlaneConfig(db, planeInclude, nil)

	err = db.First(&plane, id).Error
	plane.SetTimesToUTC()
	return
}

type PartialUpdatePlaneStruct struct {
	PrefPilotId       *uint          `json:"PrefPilotId"`
	FuelburnPerFlight *float32       `json:"FuelburnPerFlight"`
	FlightDuration    *time.Duration `json:"FlightDuration"`
	SlotEndTime       *time.Time     `json:"SlotEndTime"`
	SlotStartTime     *time.Time     `json:"SlotStartTime"`
}

func (dh *DatabaseHandler) PartialUpdatePlane(id uint, updateData PartialUpdatePlaneStruct) (plane model.Plane, err error) {
	plane, err = GetPlaneById(id, &PlaneInclude{IncludeAllowedPilots: true, IncludePrefPilot: true, IncludeDivision: true})
	if err != nil {
		return
	}

	if updateData.PrefPilotId != nil {
		if *updateData.PrefPilotId != *plane.PrefPilotId {
			status := false

			for _, pilot := range *plane.AllowedPilots {
				println("UpdateData.PrefPilotId: ", *updateData.PrefPilotId)
				println("Piolt.ID: ", pilot.ID)
				if *updateData.PrefPilotId == pilot.ID {
					plane.PrefPilotId = updateData.PrefPilotId
					plane.PrefPilot = &pilot
					println("set: ", *plane.PrefPilotId)
					status = true
					break
				}
			}

			if !status {
				err = cerror.ErrPilotNotInAllowedPilots
                err = cerror.NewInvalidFlightLogicError("Pilot is not allowed to fly this plane")
				return
			}
		}
	}

	if updateData.FuelburnPerFlight != nil {
		if updateData.FuelburnPerFlight != &plane.FuelburnPerFlight {
			plane.FuelburnPerFlight = *updateData.FuelburnPerFlight
		}
	}

	if updateData.FlightDuration != nil {
		if updateData.FlightDuration != &plane.FlightDuration {
			plane.FlightDuration = *updateData.FlightDuration
		}
	}

	if updateData.SlotStartTime != nil {
		if updateData.SlotStartTime != &plane.SlotStartTime {
			plane.SlotStartTime = *updateData.SlotStartTime
		}
	}

	if updateData.SlotEndTime != nil {
		if updateData.SlotEndTime != &plane.SlotEndTime {
			plane.SlotEndTime = *updateData.SlotEndTime
		}
	}

	if updateData.SlotStartTime != nil || updateData.SlotEndTime != nil {
		if plane.SlotStartTime.After(plane.SlotEndTime) || plane.SlotStartTime.Equal(plane.SlotEndTime) {
			err = cerror.ErrSlotTimeInvalid
            err = cerror.NewInvalidFlightLogicError("Slot time invalid")
			return
		}
	}

	err = dh.Db.Updates(&plane).Error
	if err != nil {
		return
	}

	dh.rt.AddEvent(realtime.PlaneStream, realtime.UPDATED, &plane)
	return
}

func SeedPlane(db *gorm.DB) {
	if db == nil {
		db = Db
	}
	motorflug := model.Division{}
	motorsegler := model.Division{}
	segelflug := model.Division{}

	motorErr := Db.First(&motorflug, "name = ?", "Motorflug")
	motsegErr := Db.First(&motorsegler, "name = ?", "Motorsegler")
	segelErr := Db.First(&segelflug, "name = ?", "Segelflug")

	if motorErr.Error != nil || motsegErr.Error != nil || segelErr.Error != nil {
		log.Error("Error while seeding planes: " + motorErr.Error.Error() + " " + motsegErr.Error.Error() + " " + segelErr.Error.Error())
		return
	}

	startTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 9, 0, 0, 0, time.UTC)
	endTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 17, 0, 0, 0, time.UTC)

	db.FirstOrCreate(&model.Plane{}, model.Plane{
		Registration:         "D-ELXX",
		AircraftType:         "C172",
		FlightDuration:       time.Duration(24 * time.Minute),
		FuelStartAmount:      100,
		FuelMaxCapacity:      140,
		FuelburnPerFlight:    20,
		FuelConversionFactor: 0.72,
		MaxSeatPayload:       -1,
		MTOW:                 1050,
		EmptyWeight:          650,
		DivisionId:           motorflug.ID,
		SlotStartTime:        startTime,
		SlotEndTime:          endTime,
		PassNoBase:           500,
	})

	db.FirstOrCreate(&model.Plane{}, model.Plane{
		Registration:         "D-ESXX",
		AircraftType:         "C172",
		FlightDuration:       time.Duration(24 * time.Minute),
		FuelStartAmount:      80,
		FuelMaxCapacity:      120,
		FuelburnPerFlight:    15,
		FuelConversionFactor: 0.72,
		MaxSeatPayload:       -1,
		MTOW:                 1083,
		EmptyWeight:          756,
		DivisionId:           motorflug.ID,
		SlotStartTime:        startTime,
		SlotEndTime:          endTime,
		PassNoBase:           600,
	})

	db.FirstOrCreate(&model.Plane{}, model.Plane{
		Registration:         "D-ESYY",
		AircraftType:         "C172",
		FlightDuration:       time.Duration(24 * time.Minute),
		FuelStartAmount:      160,
		FuelMaxCapacity:      180,
		FuelburnPerFlight:    20,
		FuelConversionFactor: 0.72,
		MaxSeatPayload:       -1,
		MTOW:                 1100,
		EmptyWeight:          734,
		DivisionId:           motorflug.ID,
		SlotStartTime:        startTime,
		SlotEndTime:          endTime,
	})

	db.FirstOrCreate(&model.Plane{}, model.Plane{
		Registration:         "D-EFXX",
		AircraftType:         "PA28",
		FlightDuration:       time.Duration(24 * time.Minute),
		FuelStartAmount:      100,
		FuelMaxCapacity:      140,
		FuelburnPerFlight:    20,
		FuelConversionFactor: 0.72,
		MaxSeatPayload:       -1,
		MTOW:                 1050,
		EmptyWeight:          663,
		DivisionId:           motorflug.ID,
		SlotStartTime:        startTime,
		SlotEndTime:          endTime,
		PassNoBase:           700,
	})

	db.FirstOrCreate(&model.Plane{}, model.Plane{
		Registration:         "D-KOXX",
		AircraftType:         "HK36",
		FlightDuration:       time.Duration(24 * time.Minute),
		FuelStartAmount:      80,
		FuelMaxCapacity:      80,
		FuelburnPerFlight:    10,
		FuelConversionFactor: 0.72,
		MaxSeatPayload:       -1,
		MTOW:                 800,
		EmptyWeight:          600,
		DivisionId:           motorsegler.ID,
		SlotStartTime:        startTime,
		SlotEndTime:          endTime,
	})

	db.FirstOrCreate(&model.Plane{}, model.Plane{
		Registration:         "D-KIXX",
		AircraftType:         "SF25C",
		FlightDuration:       time.Duration(24 * time.Minute),
		FuelStartAmount:      40,
		FuelMaxCapacity:      40,
		FuelburnPerFlight:    5,
		FuelConversionFactor: 0.72,
		MaxSeatPayload:       -1,
		MTOW:                 450,
		EmptyWeight:          300,
		DivisionId:           motorsegler.ID,
		SlotStartTime:        startTime,
		SlotEndTime:          endTime,
		PassNoBase:           400,
	})

	db.FirstOrCreate(&model.Plane{}, model.Plane{
		Registration:         "D-0761",
		AircraftType:         "ASK21",
		FlightDuration:       time.Duration(10 * time.Minute),
		FuelStartAmount:      0,
		FuelMaxCapacity:      -1,
		FuelburnPerFlight:    -1,
		FuelConversionFactor: -1,
		MaxSeatPayload:       110,
		MTOW:                 500,
		EmptyWeight:          300,
		DivisionId:           segelflug.ID,
		SlotStartTime:        startTime,
		SlotEndTime:          endTime.Add(2 * time.Hour),
		PassNoBase:           100,
	})

	db.FirstOrCreate(&model.Plane{}, model.Plane{
		Registration:         "D-7208",
		AircraftType:         "Duo Discus",
		FlightDuration:       time.Duration(10 * time.Minute),
		FuelStartAmount:      0,
		FuelMaxCapacity:      -1,
		FuelburnPerFlight:    -1,
		FuelConversionFactor: -1,
		MaxSeatPayload:       110,
		MTOW:                 520,
		EmptyWeight:          300,
		DivisionId:           segelflug.ID,
		SlotStartTime:        startTime,
		SlotEndTime:          endTime.Add(2 * time.Hour),
		PassNoBase:           100,
	})
}
