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

func CloseDatabase() {
    log.Warn("Closing database connection")
    err := PgConn.Close(context.Background())
    if err != nil {
        println(err)
    }
}
