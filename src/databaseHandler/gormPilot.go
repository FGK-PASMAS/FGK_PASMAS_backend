package databasehandler

import "github.com/MetaEMK/FGK_PASMAS_backend/model"

func initPilot() {
    Db.AutoMigrate(&model.Pilot{})
}

func SeedPilot() {

}
