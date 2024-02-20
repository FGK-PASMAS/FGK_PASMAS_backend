package databasehandler

import (
	"time"

	"github.com/MetaEMK/FGK_PASMAS_backend/logging"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"gorm.io/gorm"
)

func initPlane(db *gorm.DB) {
    Db.AutoMigrate(&model.Plane{})
    SeedPlane(db)
}

func GetPlanes(planeInclude *PlaneInclude, planeFilter *PlaneFilter) ([]model.Plane, error) {
    db := Db
    planes := []model.Plane{}

    db =interpretPlaneConfig(db, planeInclude, planeFilter)

    db = db.Find(&planes)
    return planes, db.Error
}

func GetPlaneById(id uint, planeInclude *PlaneInclude) (plane model.Plane, err error) {
    db := Db

    db = interpretPlaneConfig(db, planeInclude, nil)

    err = db.First(&plane, id).Error
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
        logging.DbLogger.Error("Error while seeding planes: " + motorErr.Error.Error() + " " + motsegErr.Error.Error() + " " + segelErr.Error.Error())
        return
    }

    db.FirstOrCreate(&model.Plane{}, model.Plane{
        Registration: "D-ELXX",
        AircraftType: "C172",
        FlightDuration: time.Duration(24 * time.Minute),
        FuelStartAmount: 100,
        FuelMaxCapacity: 140,
        FuelburnPerFlight: 20,
        FuelConversionFactor: 0.72,
        MaxSeatPayload: -1,
        MTOW: 1050,
        EmptyWeight: 650,
        DivisionId: motorflug.ID,
    })

    db.FirstOrCreate(&model.Plane{}, model.Plane{
        Registration: "D-ESXX",
        AircraftType: "C172",
        FlightDuration: time.Duration(24 * time.Minute),
        FuelStartAmount: 80,
        FuelMaxCapacity: 120,
        FuelburnPerFlight: 15,
        FuelConversionFactor: 0.72,
        MaxSeatPayload: -1,
        MTOW: 1083,
        EmptyWeight: 756,
        DivisionId: motorflug.ID,
    })

    db.FirstOrCreate(&model.Plane{}, model.Plane{
        Registration: "D-ESYY",
        AircraftType: "C172",
        FlightDuration: time.Duration(24 * time.Minute),
        FuelStartAmount: 160,
        FuelMaxCapacity: 180,
        FuelburnPerFlight: 20,
        FuelConversionFactor: 0.72,
        MaxSeatPayload: -1,
        MTOW: 1100,
        EmptyWeight: 734,
        DivisionId: motorflug.ID,
    })

    db.FirstOrCreate(&model.Plane{}, model.Plane{
        Registration: "D-EFXX",
        AircraftType: "PA28",
        FlightDuration: time.Duration(24 * time.Minute),
        FuelStartAmount: 100,
        FuelMaxCapacity: 140,
        FuelburnPerFlight: 20,
        FuelConversionFactor: 0.72,
        MaxSeatPayload: -1,
        MTOW: 1050,
        EmptyWeight: 663,
        DivisionId: motorflug.ID,
    })

    db.FirstOrCreate(&model.Plane{}, model.Plane{
        Registration: "D-KOXX",
        AircraftType: "HK36",
        FlightDuration: time.Duration(24 * time.Minute),
        FuelStartAmount: 80,
        FuelMaxCapacity: 80,
        FuelburnPerFlight: 10,
        FuelConversionFactor: 0.72,
        MaxSeatPayload: -1,
        MTOW: 800,
        EmptyWeight: 600,
        DivisionId: motorsegler.ID,
    })

    db.FirstOrCreate(&model.Plane{}, model.Plane{
        Registration: "D-KIXX",
        AircraftType: "SF25C",
        FlightDuration: time.Duration(24 * time.Minute),
        FuelStartAmount: 40,
        FuelMaxCapacity: 40,
        FuelburnPerFlight: 5,
        FuelConversionFactor: 0.72,
        MaxSeatPayload: -1,
        MTOW: 450,
        EmptyWeight: 300,
        DivisionId: motorsegler.ID,
    })

    db.FirstOrCreate(&model.Plane{}, model.Plane{
        Registration: "D-0761",
        AircraftType: "ASK21",
        FlightDuration: time.Duration(10 * time.Minute),
        FuelStartAmount: 0,
        FuelMaxCapacity: -1,
        FuelburnPerFlight: -1,
        FuelConversionFactor: -1,
        MaxSeatPayload: 110,
        MTOW: 500,
        EmptyWeight: 300,
        DivisionId: segelflug.ID,
    })

    db.FirstOrCreate(&model.Plane{}, model.Plane{
        Registration: "D-7208",
        AircraftType: "Duo Discus",
        FlightDuration: time.Duration(10 * time.Minute),
        FuelStartAmount: 0,
        FuelMaxCapacity: -1,
        FuelburnPerFlight: -1,
        FuelConversionFactor: -1,
        MaxSeatPayload: 110,
        MTOW: 520,
        EmptyWeight: 300,
        DivisionId: segelflug.ID,
    })
}
