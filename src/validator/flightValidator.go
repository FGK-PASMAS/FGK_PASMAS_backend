package validator

import (
	"errors"

	"github.com/MetaEMK/FGK_PASMAS_backend/model"
)


func ValidateFlightValues(flight model.Flight) {
}

var (
    ErrInvalidFlightType = errors.New("Type is invalid")
)
