package api

import (
	"net/http"

	pasmasservice "github.com/MetaEMK/FGK_PASMAS_backend/service/pasmasService"
)

type ApiError struct {
    HttpCode int
    ErrorResponse ErrorResponse
}

var (
    unknownError = ApiError { HttpCode: http.StatusInternalServerError, ErrorResponse: ErrorResponse { Success: false, Type: "UNKNOWN_ERROR"} }
    databaseError = ApiError { HttpCode: http.StatusInternalServerError, ErrorResponse: ErrorResponse { Success: false, Type: "DATABASE_ERROR"} }
    InvalidRequestBody = ApiError { HttpCode: http.StatusBadRequest, ErrorResponse: ErrorResponse { Success: false, Type: "INVALID_REQUEST_BODY"} }
    invalidObjectDependency = ApiError { HttpCode: http.StatusBadRequest, ErrorResponse: ErrorResponse { Success: false, Type: "INVALID_OBJECT_DEPENDENCY"} }
    objectNotFound = ApiError { HttpCode: http.StatusNotFound, ErrorResponse: ErrorResponse { Success: false, Type: "OBJECT_NOT_FOUND"} }
)

func GetErrorResponse(err error) ApiError {
    switch err {
        case pasmasservice.ErrObjectDependencyDivisionMissing:
           return invalidObjectDependency

        case pasmasservice.ErrObjectNotFound:
            return objectNotFound

        case pasmasservice.ErrDbQuery:
            return databaseError
        case pasmasservice.ErrNoDbConnection:
            return databaseError

        default:
            return unknownError
    }
}
