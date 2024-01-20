package api

import (
	"net/http"

	pasmasservice "github.com/MetaEMK/FGK_PASMAS_backend/service/pasmasService"
	"github.com/MetaEMK/FGK_PASMAS_backend/validator"
)

type ApiError struct {
    HttpCode int
    ErrorResponse ErrorResponse
}

var (
    unknownError = ApiError { HttpCode: http.StatusInternalServerError, ErrorResponse: ErrorResponse { Success: false, Type: "UNKNOWN_ERROR"} }
    InvalidRequestBody = ApiError { HttpCode: http.StatusBadRequest, ErrorResponse: ErrorResponse { Success: false, Type: "INVALID_REQUEST_BODY"} }
    InvalidObjectParamenter = ApiError { HttpCode: http.StatusBadRequest, ErrorResponse: ErrorResponse { Success: false, Type: "INVALID_OBJECT_PARAMETER"} }
    invalidObjectDependency = ApiError { HttpCode: http.StatusBadRequest, ErrorResponse: ErrorResponse { Success: false, Type: "INVALID_OBJECT_DEPENDENCY"} }
    objectNotFound = ApiError { HttpCode: http.StatusNotFound, ErrorResponse: ErrorResponse { Success: false, Type: "OBJECT_NOT_FOUND"} }
)

func GetErrorResponse(err error) ApiError {
    switch err {
        case pasmasservice.ErrObjectDependencyDivisionMissing:
           return invalidObjectDependency

        case pasmasservice.ErrObjectNotFound:
            return objectNotFound

        case validator.ErrPassengerWeight:
            return InvalidRequestBody

        case validator.ErrPassengerLastName:
            return InvalidRequestBody

        default:
            return unknownError
    }
}
