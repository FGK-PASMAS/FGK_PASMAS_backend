package main

import (
	"github.com/MetaEMK/FGK_PASMAS_backend/database"
	"github.com/MetaEMK/FGK_PASMAS_backend/router"
)

func main() {
    _, err := database.SetupDatabase()

    if err != nil {
        panic(err)
    }

    _, dberr := database.Db.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, username VARCHAR(50) UNIQUE NOT NULL, password VARCHAR(50) NOT NULL, email VARCHAR(50) UNIQUE NOT NULL, role VARCHAR(50) NOT NULL, created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP)")
    
    if dberr != nil {
        panic(dberr)
    }

    server := router.InitRouter() 
    server.Run(":8080")
}
