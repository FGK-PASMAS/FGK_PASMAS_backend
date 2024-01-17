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
        user = "pasmas"
    }

    if database == "" {
        database = "pasmas"
    }

    connectionString := "postgresql://" + user + ":" + password + "@" + hostname + ":5432" + "/" + database + "?sslmode=disable"

    return connectionString, nil
}

// getInitDatabaseStructure returns the SQL-Statements as a string to create the database structure
func getInitDatabaseStructure() string {
    return structureStatements
}


// getSeedDatabaseQueries return the SQL-Statements as a string to seed the given table
func getSeedDatabaseQueries(table string) string  {
    switch table {
    case "division":
        return seedDivisionStatements
    default:
        return ""
    }
}

