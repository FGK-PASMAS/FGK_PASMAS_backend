package databasehandler

import (
	"fmt"
	"os"
)

var hostname = os.Getenv("DATABASE_HOSTNAME")
var user = os.Getenv("DATABASE_USER")
var password = os.Getenv("DATABASE_PASSWORD")
var database = os.Getenv("DATABASE_NAME")

func GetConnectionString() string {
    if password == "" {
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

    return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", hostname, user, password, database)
}
