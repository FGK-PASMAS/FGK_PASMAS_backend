package api

import (
	"errors"
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
    objectNotFound = ApiError { HttpCode: http.StatusNotFound, ErrorResponse: ErrorResponse { Success: false, Type: "OBJECT_NOT_FOUND"} }
    notImplemented = ApiError { HttpCode: http.StatusNotImplemented, ErrorResponse: ErrorResponse { Success: false, Type: "NOT_IMPLEMENTED"} }
)

var (
    ErrInvalidFlightType = errors.New("Flight type is not valid")
    ErrInvalidPlane = errors.New("Plane is not valid")
    ErrNotImplemented = errors.New("Functionality not implemented")
)

func GetErrorResponse(err error) ApiError {
    var obj ApiError

    switch err {
        case pasmasservice.ErrObjectNotFound:
            obj = objectNotFound

        case // InvalidRequestBody
            validator.ErrPassengerWeight,
            validator.ErrPassengerLastName,

            validator.ErrInvalidDepartureTime,
            ErrInvalidFlightType:
                obj = InvalidRequestBody

        case pasmasservice.ErrSlotIsNotFree:
            obj = InvalidRequestBody


        default:
            obj = unknownError
    }

    obj.ErrorResponse.Message = err.Error()
    return obj
}
