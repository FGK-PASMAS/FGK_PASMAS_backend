package debugService

import (
	databasehandler "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/service/flightService"
)

func TruncateData() error {
    flightService.LockAll()
    err := databasehandler.ResetDatabase()
    flightService.UnlockAll()
    return err
}
