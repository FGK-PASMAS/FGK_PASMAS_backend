package databasehandler

import "github.com/MetaEMK/FGK_PASMAS_backend/model"

func initPassenger() {
    Db.AutoMigrate(&model.Passenger{})
}
