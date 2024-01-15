package database

import (
	"context"

	"github.com/MetaEMK/FGK_PASMAS_backend/logging"
	"github.com/jackc/pgx/v5"
)

var log = logging.DbLogger
var PgConn *pgx.Conn

func SetupDatabaseConnection() {
    log.Info("Trying to connect to the database")

    connectionString, err := getConnectionString()
    if err != nil {
        log.Error("Failed to generate the connectionString")
        panic(err)
    }

    pgx, err := pgx.Connect(context.Background(), connectionString)

    if err != nil {
        log.Error("Failed to open a new connection")
        panic(err)
    } else {
        log.Debug("Successfully opened a new connection")
        PgConn = pgx
    }

    CheckDatabaseConnection()
}

func CheckDatabaseConnection() error {
    err := PgConn.Ping(context.Background())

    if err != nil {
        log.Warn("Failed to ping the database")
    } 

    return err
}

func InitDatabaseStructure() (error){
    log.Info("Trying to create the database structure")
    statements, err := getInitDatabaseStructure()
    if(err != nil) {
        log.Error("Failed to retrieve the sql statements for creating the database structure")
        PgConn.Close(context.Background())
        return err
    } 

    _, err = PgConn.Exec(context.Background(), statements)

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

        statements, err := getSeedDatabaseQueries("division")
        if err != nil {
            log.Error("Failed to retrieve the sql statements for seeding the division table")

        } else {
            _, err = PgConn.Exec(context.Background(), statements)
            if err != nil {
                log.Error("Failed to seed the division table")
            } else {
                log.Debug("Successfully seeded the division table")
            }
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
