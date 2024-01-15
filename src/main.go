package main

import (
	"github.com/MetaEMK/FGK_PASMAS_backend/database"
	"github.com/MetaEMK/FGK_PASMAS_backend/router"
)

func main() {
    err := database.SetupDatabaseConnection()
    if err != nil {
        panic(err)
    }

    database.InitDatabaseStructure()

    database.SeedDatabase()

    go database.AutoReconnectForDatabaseConnection()

    server := router.InitRouter() 
    server.Run(":8081")
}
