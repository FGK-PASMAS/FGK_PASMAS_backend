package database

import (
	"os"
)

var sep = string(os.PathSeparator)
// ConnectionsString Paramenters

var hostname = os.Getenv("DATABASE_HOSTNAME")
var user = os.Getenv("DATABASE_USER")
var password = os.Getenv("DATABASE_PASSWORD")
var database = os.Getenv("DATABASE_NAME")

// Parameters for DatabaseStructure Tasks
var sqlFolder = "database" + sep + "sqlScripts"
var createStructureFile = "createDatabaseStructure.sql"

func getConnectionString() (string, error) {
    //TODO: Logging

    if password == "" {
        //TODO: Adding Error for no password
        //BUG: How to Handle Passwords?

        //return "", errors.New("No password set for database")
        password = "password"
    }

    if hostname == "" {
        hostname = "localhost"
    } 

    if user == "" {
        user = "fgk_pasmas"
    }

    if database == "" {
        database = "pasmas"
    }

    connectionString := "postgresql://" + user + ":" + password + "@" + hostname + ":5432" + "/" + database + "?sslmode=disable"

    return connectionString, nil
}

func getInitDatabaseStructure() (string, error) {
    pwd, err := os.Getwd()
    filepath := pwd + sep + sqlFolder + sep + createStructureFile
    statements, err := os.ReadFile(filepath)

    return string(statements), err
}
