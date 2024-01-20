package main

import (
	databasehandler "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/router"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var mode = "DEBUG"

func main() {

    dsn := databasehandler.GetConnectionString()
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{

    })
    if err != nil {
        panic("failed to connect database")
    }

    databasehandler.InitGorm(db)

    server := router.InitRouter() 
    server.Run(":8080")
}
