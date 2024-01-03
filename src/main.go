package main

import (
	"os"
	"os/exec"

	"github.com/MetaEMK/FGK_PASMAS_backend/database"
	"github.com/MetaEMK/FGK_PASMAS_backend/router"
)

func main() {
    cmd := exec.Command("clear")
    cmd.Stdout = os.Stdout
    cmd.Run()

    database.SetupDatabaseConnection()

    database.InitDatabaseStructure()

    server := router.InitRouter() 
    server.Run(":8080")
}
