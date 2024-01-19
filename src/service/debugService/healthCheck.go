package debugservice

import (
	"github.com/MetaEMK/FGK_PASMAS_backend/database"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
)

func GetHealthCheck() model.HealthCheckResponse {
    res :=  model.HealthCheckResponse{}

    dbErr := database.CheckDatabaseConnection()
    if(dbErr != nil) {
        res.DatabaseConnection = dbErr.Error()
    } else {
        res.DatabaseConnection = "successfull"
    }

    
    return res
}

