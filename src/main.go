package main

import (
	"github.com/MetaEMK/FGK_PASMAS_backend/database"
	"github.com/MetaEMK/FGK_PASMAS_backend/router"
)

func main() {
    database.SetupDatabaseConnection()

    database.InitDatabaseStructure()

    server := router.InitRouter() 
    server.Run(":8080")
}
