package pasmasservice

import (
	dh "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/MetaEMK/FGK_PASMAS_backend/validator"
)
func GetPassengers() ([]model.Passenger, error) {
    passengers := []model.Passenger{}
    result := dh.Db.Find(&passengers)

    return passengers, result.Error
}

func CreatePassenger(pass model.Passenger) (model.Passenger, error) {
    err := validator.ValidatePassenger(pass)
    if err != nil {
        return model.Passenger{}, err
    }

    result := dh.Db.Create(&pass)
    return pass, result.Error
}

func UpdatePassenger(id uint, pass model.Passenger) (model.Passenger, error) {
    err := validator.ValidatePassenger(pass)
    if err != nil {
        return model.Passenger{}, err
    }

    oldPass := model.Passenger{}
    result := dh.Db.First(&oldPass, id)
    if result.Error != nil {
        return model.Passenger{}, result.Error
    }

    result = dh.Db.Model(&oldPass).Updates(pass)

    return oldPass, nil
}

func DeletePassenger(id int64) error {
    pass := model.Passenger{}
    result := dh.Db.First(&pass, id)
    if result.Error != nil {
        return result.Error
    }

    result = dh.Db.Delete(&pass)

    return result.Error
}
