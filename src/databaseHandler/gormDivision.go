package databasehandler

import "github.com/MetaEMK/FGK_PASMAS_backend/model"

func initDivision() {
    Db.AutoMigrate(&model.Division{})

    SeedDivision()
}

func SeedDivision() error {
    divs := []model.Division{
        { Name: "Segelflug", PassengerCapacity: 1},
        { Name: "Motorsegler", PassengerCapacity: 1},
        { Name: "Motorflug", PassengerCapacity: 3},
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

    return nil
}
