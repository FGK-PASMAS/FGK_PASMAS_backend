package divisionService

import (
	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	databasehandler "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
)

func GetDivisions() ([]model.Division, error) {
    if databasehandler.Db == nil {
        return []model.Division{}, cerror.NewUnknownError("Database connection lost")
    }
    divisions := []model.Division{}
    result := databasehandler.Db.Order("id ASC").Find(&divisions)

    if result.Error != nil {
        return []model.Division{}, result.Error
    }

    return divisions, nil
}
