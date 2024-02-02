package databasehandler

import (
	"fmt"

	"github.com/MetaEMK/FGK_PASMAS_backend/logging"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
)

var log = logging.DbLogger

func initPilot() {
    Db.AutoMigrate(&model.Pilot{})

    SeedPilot()
}

func SeedPilot() {
    pilots := []model.Pilot{
        {FirstName: "Tobias", LastName: "Kornwestheim", Weight: 68, AllowedPilots: &[]model.Plane{{Registration: "D-0761"}, {Registration: "D-7208"}}},
        {FirstName: "Dennis", LastName: "Kornwestheim", Weight: 72, AllowedPilots: &[]model.Plane{{Registration: "D-0761"}, {Registration: "D-7208"}, {Registration: "D-KOXX"}, {Registration: "D-ELXX"}}},
        {FirstName: "Yannic", LastName: "Kornwestheim", Weight: 90, AllowedPilots: &[]model.Plane{{Registration: "D-0761"}, {Registration: "D-7208"}, {Registration: "D-KOXX"}, {Registration: "D-ELXX"}}},
        {FirstName: "Michel", LastName: "Kornwestheim", Weight: 90, AllowedPilots: &[]model.Plane{{Registration: "D-0761"}, {Registration: "D-7208"}, {Registration: "D-KOXX"}, {Registration: "D-ELXX"}}},
        {FirstName: "Horst", LastName: "Kornwestheim", Weight: 85, AllowedPilots: &[]model.Plane{{Registration: "D-KIXX"}, {Registration: "D-KOXX"}}},
        {FirstName: "Katja", LastName: "Motorflug", Weight: 100, AllowedPilots: &[]model.Plane{{Registration: "D-EFXX"}}},
        {FirstName: "Markus", LastName: "Motorflug", Weight: 100, AllowedPilots: &[]model.Plane{{Registration: "D-EFXX"}}},
        {FirstName: "Max", LastName: "Stuttgart", Weight: 100, AllowedPilots: &[]model.Plane{{Registration: "D-ESXX"}}},
        {FirstName: "Mario", LastName: "Stuttgart", Weight: 100, AllowedPilots: &[]model.Plane{{Registration: "D-ESXX"}, {Registration: "D-ESYY"}}},
        {FirstName: "Bernd", LastName: "Stuttgart", Weight: 100, AllowedPilots: &[]model.Plane{{Registration: "D-ESXX"}, {Registration: "D-ESYY"}}},
    }

    for _, pilot := range pilots {
        err := CreateOrUpdatePilot(&pilot)

        if err != nil {
            log.Warn(fmt.Sprintf("PilotSeeder: Error while seeding pilot %s %s: %s", pilot.FirstName, pilot.LastName, err))
        }
    }
}

func CreateOrUpdatePilot(pilot *model.Pilot) error {
    aircrafts := []model.Plane{}
    for _, registration := range *pilot.AllowedPilots{
        aircraft := model.Plane{}
        err := Db.Where("registration = ?", registration.Registration).First(&aircraft).Error
        if err != nil {
            log.Info(fmt.Sprintf("PilotSeeder: No aircraft with registration %s found - creating pilot without aircraft binding", registration.Registration))
        }
        if aircraft.ID != 0 {
            aircrafts = append(aircrafts, aircraft)
        }
    }

    pilot.AllowedPilots = &aircrafts
    
    err := Db.Where("last_name = ?", pilot.LastName).Where("first_name = ?", pilot.FirstName).FirstOrCreate(&pilot).Error
    return err
}
