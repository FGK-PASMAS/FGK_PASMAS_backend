package api

import (
	"errors"
	"net/http"

	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	"github.com/MetaEMK/FGK_PASMAS_backend/validator"
)

type ApiError struct {
    HttpCode int
    ErrorResponse ErrorResponse
}

var (
    unknownError = ApiError { HttpCode: http.StatusInternalServerError, ErrorResponse: ErrorResponse { Success: false, Type: "UNKNOWN_ERROR"} }
    ObjectAlreadyExists = ApiError { HttpCode: http.StatusConflict, ErrorResponse: ErrorResponse { Success: false, Type: "OBJECT_ALREADY_EXISTS"} }
    InvalidRequestBody = ApiError { HttpCode: http.StatusBadRequest, ErrorResponse: ErrorResponse { Success: false, Type: "INVALID_REQUEST_BODY"} }
    invalidFlightLogic = ApiError { HttpCode: http.StatusBadRequest, ErrorResponse: ErrorResponse { Success: false, Type: "INVALID_FLIGHT_LOGIC"} }
    objectNotFound = ApiError { HttpCode: http.StatusNotFound, ErrorResponse: ErrorResponse { Success: false, Type: "OBJECT_NOT_FOUND"} }
    dependencyNotFound = ApiError { HttpCode: http.StatusBadRequest, ErrorResponse: ErrorResponse { Success: false, Type: "DEPENDENCY_NOT_FOUND"} }
    notImplemented = ApiError { HttpCode: http.StatusNotImplemented, ErrorResponse: ErrorResponse { Success: false, Type: "NOT_IMPLEMENTED"} }
    notValidParameters = ApiError { HttpCode: http.StatusBadRequest, ErrorResponse: ErrorResponse { Success: false, Type: "NOT_VALID_PARAMETERS"} }
    unauthorized = ApiError { HttpCode: http.StatusUnauthorized, ErrorResponse: ErrorResponse { Success: false, Type: "UNAUTHORIZED"} }
    forbidden = ApiError { HttpCode: http.StatusForbidden, ErrorResponse: ErrorResponse { Success: false, Type: "FORBIDDEN"} }
)

var (
    ErrInvalidFlightStatus = errors.New("Flight status is not valid")
    ErrInvalidPlane = errors.New("Plane is not valid")
    ErrNotImplemented = errors.New("Functionality not implemented")
)

func GetErrorResponse(err error) ApiError {
    var obj ApiError

    switch err {
        case cerror.ErrObjectNotFound:
            obj = objectNotFound

        case // InvalidRequestBody
            ErrInvalidFlightStatus:
                obj = InvalidRequestBody

        // dependencyNotFound
        case 
            validator.ErrInvalidPilot,
            validator.ErrInvalidPlane,
            cerror.ErrObjectDependencyMissing:
                obj = dependencyNotFound

        // invalidFlightLogic
        case 
            validator.ErrPassengerWeight,
            validator.ErrPassengerLastName,
            validator.ErrInvalidDepartureTime,
            cerror.ErrFlightStatusDoesNotFitProcess,
            cerror.ErrNoPilotAvailable,
            cerror.ErrNoStartFuelFound,
            cerror.ErrMaxSeatPayload,
            cerror.ErrTooManyPassenger,
            cerror.ErrTooLessPassenger,
            cerror.ErrTooMuchFuel,
            cerror.ErrTooLessFuel,
            cerror.ErrOverloaded,
            cerror.ErrSlotIsNotFree,
            cerror.ErrTimeSlotForPlaneNotValid,
            cerror.ErrDepartureTimeIsZero:
                obj = invalidFlightLogic

        case cerror.ErrIncludeNotSupported:
            obj = notValidParameters

        case cerror.ErrForbidden:
            obj = forbidden

        case cerror.ErrInvalidCredentials:
            obj = unauthorized

        case cerror.ErrUserAlreadyExists:
            obj = ObjectAlreadyExists

        default:
            obj = unknownError
    }

    obj.ErrorResponse.Message = err.Error()
    return obj
}
