package dberr

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

var (
    ErrUnknown = errors.New("UnknownError")

    // Error occurs when connection to database is lost
    ErrNoConnection = errors.New("Unable to connect to database")
    
    //Error occurs when query fails to execute
    ErrQuery = errors.New("Unable to execute query")

    //Error occurs when rows are expected but none are returned
    ErrNoRows = pgx.ErrNoRows

    //Error occurs when more rows are returned than expected
    ErrTooManyRows = pgx.ErrTooManyRows
)
