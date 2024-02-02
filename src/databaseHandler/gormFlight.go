package databasehandler

import "github.com/MetaEMK/FGK_PASMAS_backend/model"

func initFlight() {
    Db.AutoMigrate(&model.Flight{})
}
