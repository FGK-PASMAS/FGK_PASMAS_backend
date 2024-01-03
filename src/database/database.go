package database

import (
	"database/sql"

	"github.com/MetaEMK/FGK_PASMAS_backend/logging"
	_ "github.com/lib/pq"
)

var Db *sql.DB;
var log = logging.DbLogger

func SetupDatabaseConnection() {
    log.Info("Trying to connect to the database")

    connectionString, err := getConnectionString()
    if err != nil {
        log.Error("Failed to generate the connectionString")
        panic(err)
    }

    db, err := sql.Open("postgres", connectionString)

    if err != nil {
        log.Error("Failed to open a new connection")
        panic(err)
    } else {
        log.Debug("Successfully opened a new connection")
        Db = db
    }

    CheckDatabaseConnection()
}

func CheckDatabaseConnection() error {
    err := Db.Ping()

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
        Db.Close()
        return err
    } 

    _, err = Db.Exec(statements)

    if(err != nil) {
        log.Error("Failed to create the database structure")
        Db.Close()
        panic(err)
    }

    log.Debug("Successfully created the database structure")

    return nil
}

func CloseDatabase() {
    log.Warn("Closing database connection")
    err := Db.Close()
    if err != nil {
        println(err)
    }
}
