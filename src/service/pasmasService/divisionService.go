package pasmasservice

import (
	divisionhandler "github.com/MetaEMK/FGK_PASMAS_backend/database/divisionHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
)

func GetDivisions() ([]model.DivisionStructSelect, error) {
    return divisionhandler.GetDivisions()
}
