package pasmasservice

import databasehandler "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"

func TruncateData() error {
    flightCreation.Lock()
    err := databasehandler.ResetDatabase()
    flightCreation.Unlock()
    return err
}
