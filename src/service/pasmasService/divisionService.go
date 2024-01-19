package pasmasservice

import (
	"errors"

	databasehandler "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
)

func GetDivisions() ([]model.Division, error) {
    if databasehandler.Db == nil {
        return []model.Division{}, errors.New("Database is null")
    }

    divisions := []model.Division{}
    result := databasehandler.Db.Find(&divisions)

    if result.Error != nil {
        return []model.Division{}, result.Error
    }

    return divisions, nil
}
