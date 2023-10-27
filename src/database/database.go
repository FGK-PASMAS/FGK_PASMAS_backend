package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var Db *sql.DB;

func SetupDatabase() (*sql.DB, error) {
    connectionString, err := getConnectionString()
    if err != nil {
        return nil, err
    }

    db, err := sql.Open("postgres", connectionString)

    if err != nil {
        println(err)
        return nil, err
    } else {
        println("Connected to database")
        Db = db
    }

    return db, nil
}
