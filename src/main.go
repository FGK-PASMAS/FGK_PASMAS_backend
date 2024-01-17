package main

import (
	"os"

	"github.com/MetaEMK/FGK_PASMAS_backend/database"
	"github.com/MetaEMK/FGK_PASMAS_backend/router"
)

var mode = "DEBUG"

func main() {
    os.Setenv("MODE", mode)
    os.Setenv("GIN_MODE", mode)

    err := database.SetupDatabaseConnection()
    if err != nil {
        panic(err)
    }

    database.InitDatabaseStructure()

    database.SeedDatabase()

    go database.AutoReconnectForDatabaseConnection()

    server := router.InitRouter() 
    server.Run(":8080")
}
