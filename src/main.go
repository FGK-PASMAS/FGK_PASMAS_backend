package main

import (
	"os/exec"

	"github.com/MetaEMK/FGK_PASMAS_backend/database"
	"github.com/MetaEMK/FGK_PASMAS_backend/router"
)

func main() {
    exec.Command("clear")
    database.SetupDatabaseConnection()

    database.InitDatabaseStructure()

    server := router.InitRouter() 
    server.Run(":8080")
}
