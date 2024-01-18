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

type ErrorType string


func (e *InternalError) ToJson() (string) {
    jsonBytes, err := json.Marshal(e)

    if err != nil {
        return ""
    }

    return string(jsonBytes)
}


func (e InternalError) Error() string {
    str := fmt.Sprintf("Type: %s, Message: %s", e.Type, e.Message)
    return str
}

const (
    //This error is returned when an unknown error occurs
    ErrorUnknownError ErrorType = "UnknownError"


    //Database errors

    //This error is returned when the database is not setup correctly
    ErrorSetupDatabaseError ErrorType = "SetupDatabaseError"

    //This error is returned when the connection to the database is lost
    ErrorDatabaseConnectionError ErrorType = "DatabaseConnectionError"

    //This error is returned when the query to the database fails
    ErrorDatabaseQueryError ErrorType = "DatabaseQueryError"

    //This error is returned when the data from the database is not in the expected format (could not be parsed to the expected struct)
    ErrorParseError ErrorType = "ParseError"


    //Object errors

    //This error is returned when there is no object for an action/operation is found
    ErrorObjectNotFoundError ErrorType = "ObjectNotFoundError"

    //This error is returned when the object already exists
    ErrorObjectAlreadyExistsError ErrorType = "ObjectAlreadyExistsError"

    //This error is returned when the related elements of the object are not found
    ErrorRelatedObjectNotFoundError ErrorType = "RelatedObjectNotFoundError"
)
