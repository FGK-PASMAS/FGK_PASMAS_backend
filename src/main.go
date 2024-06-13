package main

import (
	"os"

	"github.com/MetaEMK/FGK_PASMAS_backend/config"
	databasehandler "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/logging"
	"github.com/MetaEMK/FGK_PASMAS_backend/router"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var mode = "DEBUG"

func main() {
    log := logging.NewLogger("MAIN", config.GetGlobalLogLevel())
    log.Info("FGK_PASMAS_backend starting")

    config.LoadAuthConfig()
    config.InitDbConfig()

    dsn := config.GetConnectionString()
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{

    })
    if err != nil {
        log.Error("Database connection failed")
        os.Exit(1)
    }

    databasehandler.InitGorm(db)

    server := router.InitRouter()

    tlsConfig := config.LoadTlsConfig()
    if tlsConfig != nil {
        server.RunTLS(":8080", tlsConfig.CertPath, tlsConfig.KeyPath)
    } else {
        server.Run(":8080")
    }
}
