package databasehandler

import (
	"github.com/MetaEMK/FGK_PASMAS_backend/config"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/realtime"
)

func initDivision() {
	Db.AutoMigrate(&model.Division{})

    err := SeedDivision()
    if err != nil {
        log.Warn(err.Error())
    }
}

func SeedDivision() error {
	if config.EnableSeeder {
		log.Debug("Seeding divisions")

		divs := []model.Division{
			{Name: "Segelflug", PassengerCapacity: 1},
			{Name: "Motorsegler", PassengerCapacity: 1},
			{Name: "Motorflug", PassengerCapacity: 3},
		}

		for _, div := range divs {
			d := model.Division{}
			res := Db.Where("name = ?", div.Name).Find(&d)

			if res.RowsAffected == 0 {
				Db.Create(&div)
			} else {
				Db.Model(&d).Updates(div)
			}
		}
	}

	divisions := []model.Division{}
	err := Db.Find(&divisions).Error
	if err != nil {
		return err
	}

	realtime.InitAllFlightByDivisionEndpoints(divisions)

	return nil
}
