package cerror

import (
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

type ApiError struct {
	HttpCode int `json:"-"`
	Type     string
	Message  string
}

type ErrorResponse struct {
	Success bool
	Type    string
}

func (e ApiError) Error() string {
	return fmt.Sprintf("Error %d: %s - %s", e.HttpCode, e.Type, e.Message)
}

func New(httpCode int, typeCode string, description string) ApiError {
	return ApiError{
		HttpCode: httpCode,
		Type:     typeCode,
		Message:  description,
	}
}

func InterpretError(err error) ApiError {
	if apiErr, ok := err.(ApiError); ok {
		return apiErr
	} else {
        if err == gorm.ErrRecordNotFound {
            return NewObjectNotFoundError(err.Error())
        }
		return NewUnknownError(err.Error())
	}
}

func NewUnknownError(description string) ApiError {
	return New(http.StatusInternalServerError, "UNKNOWN_ERROR", description)
}

func NewObjectAlreadyExistsError(description string) ApiError {
	return New(http.StatusConflict, "OBJECT_ALREADY_EXISTS", description)
}

func NewInvalidRequestBodyError(description string) ApiError {
	return New(http.StatusBadRequest, "INVALID_REQUEST_BODY", description)
}

func NewInvalidFlightLogicError(description string) ApiError {
	return New(http.StatusBadRequest, "INVALID_FLIGHT_LOGIC", description)
}

func NewObjectNotFoundError(description string) ApiError {
	return New(http.StatusNotFound, "OBJECT_NOT_FOUND", description)
}

func NewDependencyNotFoundError(description string) ApiError {
	return New(http.StatusBadRequest, "DEPENDENCY_NOT_FOUND", description)
}

func NewNotImplementedError(description string) ApiError {
	return New(http.StatusNotImplemented, "NOT_IMPLEMENTED", description)
}

func NewNotValidParametersError(description string) ApiError {
	return New(http.StatusBadRequest, "NOT_VALID_PARAMETERS", description)
}

func NewAuthenticationError(description string) ApiError {
	return New(http.StatusUnauthorized, "UNAUTHORIZED", description)
}

func NewForbiddenError(description string) ApiError {
	return New(http.StatusForbidden, "FORBIDDEN", description)
}
