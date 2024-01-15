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
var seedDatabaseFile = "seed"

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
func getInitDatabaseStructure() (string, error) {
    statements, err := readStatementsFromFile(createStructureFile)
    return string(statements), err
}


// getSeedDatabaseQueries return the SQL-Statements as a string to seed the given table
func getSeedDatabaseQueries(table string) (string, error) {
    filepath := seedDatabaseFile + sep + table + ".sql"
    statements, err := readStatementsFromFile(filepath)
    return string(statements), err
}

// readStatementsFromFile reads the SQL-Statements from a file and returns them as a string
func readStatementsFromFile(filepath string) (string, error) {
    pwd, err := os.Getwd()
    path := pwd + sep + sqlFolder + sep + filepath

    statements, err := os.ReadFile(path)

    return string(statements), err
}
