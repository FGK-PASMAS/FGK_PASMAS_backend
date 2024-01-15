package internalerror

import (
	"encoding/json"
	"fmt"
)

type InternalError struct {
    //The type of the error
    Type ErrorType          `json:"type"`

    //A description of the error
    Message string          `json:"message"`

    //The error that caused this error
    Body error              `json:"error"`
}

type ErrorType int

const (
    //This error is returned when an unknown error occurs
    UnknownError ErrorType = iota

    //This error is returned when the database is not setup correctly
    SetupDatabaseError

    //This error is returned when the connection to the database is lost
    DatabaseConnectionError

    //This error is returned when the query to the database fails
    DatabaseQueryError

    //This error is returned when the data from the database is not in the expected format (could not be parsed to the expected struct)
    ParseError
)

func (e *InternalError) ToJson() (string) {
    jsonBytes, err := json.Marshal(e)

    if err != nil {
        return ""
    }

    return string(jsonBytes)
}


func (e InternalError) Error() string {
    str := fmt.Sprintf("Type: %d, Message: %s", e.Type, e.Message)
    return str
}

