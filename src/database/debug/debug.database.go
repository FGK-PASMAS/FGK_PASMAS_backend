package debug

import (
	"context"

	"github.com/MetaEMK/FGK_PASMAS_backend/database"
	dberr "github.com/MetaEMK/FGK_PASMAS_backend/database/dbErr"
	"github.com/MetaEMK/FGK_PASMAS_backend/logging"
)

// TODO: REMOVE THIS THING: THIS IS FOR DEBUG PURPOSES ONLY
// IMPORTANT: DO NOT USE THIS IN PRODUCTION

var log = logging.DbDebugLogger

var mode = "DEBUG"

// TruncateDatabase truncates the database and seeds it with default values
func TruncateDatabase() error {
    if mode!= "DEBUG" {
        log.Debug(mode)
        return dberr.ErrUnknown
    }

    log.Warn("TRUNCATING DATABASE")

    connErr := database.CheckDatabaseConnection()
    if connErr != nil {
        return dberr.ErrNoConnection
    }

    query := `
        truncate table passenger restart identity cascade;
        truncate table division restart identity cascade; `

    _, err := database.PgConn.Exec(context.Background(), query)

    if err != nil {
        return dberr.ErrQuery
    }

    log.Warn("TRUNCATING FINISHED - seeding")
    database.SeedDatabase()

    return nil
}
