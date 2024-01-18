package database

import (
	"context"
	"time"

	internalerror "github.com/MetaEMK/FGK_PASMAS_backend/internalError"
	"github.com/MetaEMK/FGK_PASMAS_backend/logging"
	"github.com/jackc/pgx/v5"
)

var log = logging.DbLogger
var PgConn *pgx.Conn
var isDatabaseConnected bool = false

func SetupDatabaseConnection() error {
    log.Info("Trying to connect to the database")

    connectionString, err := getConnectionString()
    if err != nil {
        log.Error("Failed to generate the connectionString")
        return internalerror.InternalError{Type: internalerror.ErrorDatabaseConnectionError, Message: "Failed to generate the connectionString", Body: err}
    }

    pgx, err := pgx.Connect(context.Background(), connectionString)

    if err != nil {
        log.Error("Failed to open a new connection")
        return internalerror.InternalError{Type: internalerror.ErrorDatabaseConnectionError, Message: "Failed to open a new connection", Body: err}
    } else {
        log.Debug("Successfully opened a new connection")
        PgConn = pgx
    }

    return nil
}

func CheckDatabaseConnection() error {
    err := PgConn.Ping(context.Background())

    if err != nil {
        error := internalerror.InternalError{Type: internalerror.ErrorDatabaseConnectionError, Message: "Failed to ping the database", Body: err}
        return error
    } 

    return nil
}

func AutoReconnectForDatabaseConnection() {
    for {
        err := CheckDatabaseConnection()
        if err != nil {
            isDatabaseConnected = false
            log.Error("Failed to ping the database - trying to reconnect")

            //CloseDatabase()

            SetupDatabaseConnection()

            err = CheckDatabaseConnection()
            if err == nil {
                isDatabaseConnected = true
                log.Debug("Successfully reconnected to the database")
            }
        }
        time.Sleep(2 * time.Second)
    }
}


func InitDatabaseStructure() (error){
    log.Info("Trying to create the database structure")
    statements := getInitDatabaseStructure()

    _, err := PgConn.Exec(context.Background(), statements)

    if(err != nil) {
        log.Error("Failed to create the database structure")
        PgConn.Close(context.Background())
        panic(err)
    }

    log.Debug("Successfully created the database structure")

    return nil
}

// SeedDatabase seeds the database with data
// Division: If at least one division is in the database, the division table is considered as seeded
func SeedDatabase() {
    log.Info("Start seeding the database")

    // Division
    log.Info("Seeding the division table")
    rows, err := PgConn.Query(context.Background(), "SELECT id FROM division")
    defer rows.Close()

    var divisions []int
    for rows.Next() {
        var division int
        err = rows.Scan(&division)
        if err == nil {
            divisions = append(divisions, division)
        }
    }

    if err != nil {
        log.Warn("Failed to get the divisions from the database")
    } else if len(divisions) == 0 {
        log.Debug("No divisions found in the database - seeding the division table")

        statements := getSeedDatabaseQueries("division")
        _, err = PgConn.Exec(context.Background(), statements)
        if err != nil {
            log.Error("Failed to seed the division table")
        } else {
            log.Debug("Successfully seeded the division table")
        }
    } else {
        log.Debug("Division table is already seeded")
    }
}

func CloseDatabase() {
    log.Warn("Closing database connection")
    err := PgConn.Close(context.Background())
    if err != nil {
        println(err)
    }
}
