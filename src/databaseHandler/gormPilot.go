package databasehandler

import (
	"fmt"

	"github.com/MetaEMK/FGK_PASMAS_backend/config"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"gorm.io/gorm"
)

func initPilot(db *gorm.DB) {
	Db.AutoMigrate(&model.Pilot{})

	if config.EnableSeeder {
		log.Debug("Seeding pilots")
		SeedPilot(db)
	}
}

func SeedPilot(db *gorm.DB) {
	pilots := []model.Pilot{
		{FirstName: "Alex", LastName: "Beispielpilot", Weight: 75, AllowedPilots: &[]model.Plane{{Registration: "D-0761"}, {Registration: "D-7208"}}},
		{FirstName: "Felix", LastName: "Musterflug", Weight: 68, AllowedPilots: &[]model.Plane{{Registration: "D-0761"}, {Registration: "D-7208"}, {Registration: "D-KOXX"}, {Registration: "D-ELXX"}}},
		{FirstName: "Sebastian", LastName: "Testpilot", Weight: 85, AllowedPilots: &[]model.Plane{{Registration: "D-0761"}, {Registration: "D-7208"}, {Registration: "D-KOXX"}, {Registration: "D-ELXX"}}},
		{FirstName: "Lukas", LastName: "Beispielflug", Weight: 88, AllowedPilots: &[]model.Plane{{Registration: "D-0761"}, {Registration: "D-7208"}, {Registration: "D-KOXX"}, {Registration: "D-ELXX"}}},
		{FirstName: "Thomas", LastName: "Testflug", Weight: 82, AllowedPilots: &[]model.Plane{{Registration: "D-KIXX"}, {Registration: "D-KOXX"}}},
		{FirstName: "Sophia", LastName: "Beispiel", Weight: 95, AllowedPilots: &[]model.Plane{{Registration: "D-EFXX"}}},
		{FirstName: "Anna", LastName: "Mustermann", Weight: 97, AllowedPilots: &[]model.Plane{{Registration: "D-EFXX"}}},
		{FirstName: "Paul", LastName: "Testflug", Weight: 93, AllowedPilots: &[]model.Plane{{Registration: "D-ESXX"}}},
		{FirstName: "Julian", LastName: "Musterpilot", Weight: 89, AllowedPilots: &[]model.Plane{{Registration: "D-ESXX"}, {Registration: "D-ESYY"}}},
		{FirstName: "Maximilian", LastName: "Beispielflug", Weight: 92, AllowedPilots: &[]model.Plane{{Registration: "D-ESXX"}, {Registration: "D-ESYY"}}},
	}

	for _, pilot := range pilots {
		err := createOrUpdatePilot(db, &pilot)

		if err != nil {
			log.Warn(fmt.Sprintf("PilotSeeder: Error while seeding pilot %s %s: %s", pilot.FirstName, pilot.LastName, err))
		}
	}
}

func createOrUpdatePilot(db *gorm.DB, pilot *model.Pilot) error {
	pilot.SetTimesToUTC()
	if db == nil {
		db = Db
	}

	aircrafts := []model.Plane{}
	for _, registration := range *pilot.AllowedPilots {
		aircraft := model.Plane{}
		err := Db.Preload("PrefPilot").Where("registration = ?", registration.Registration).First(&aircraft).Error
		if err != nil {
			log.Info(fmt.Sprintf("PilotSeeder: No aircraft with registration %s found - creating pilot without aircraft binding", registration.Registration))
		}
		if aircraft.ID != 0 {
			aircrafts = append(aircrafts, aircraft)
		}
	}

	pilot.AllowedPilots = &aircrafts

	err := db.Where("last_name = ?", pilot.LastName).Where("first_name = ?", pilot.FirstName).FirstOrCreate(&pilot).Error

	for _, a := range aircrafts {
		if a.PrefPilot == nil || a.PrefPilot.Weight < pilot.Weight {
			db.Model(&model.Plane{}).Where("id = ?", a.ID).UpdateColumn("PrefPilotId", pilot.ID)
		}
	}

	return err
}
