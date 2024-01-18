package passenger

import (
	"net/http"
	"strconv"

	passengerHandler "github.com/MetaEMK/FGK_PASMAS_backend/database/passengerHandler"
	internalerror "github.com/MetaEMK/FGK_PASMAS_backend/internalError"
	"github.com/MetaEMK/FGK_PASMAS_backend/logging"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api"
	"github.com/gin-gonic/gin"
)

var log = logging.ApiLogger
type intError = internalerror.InternalError

func getPassengers(c *gin.Context) {
    passengers, err := passengerHandler.GetPassengers()
    var statusCode int
    var response interface{}

    if(err != nil) {
        statusCode = http.StatusInternalServerError
        response = api.ErrorResponse {
            Success: false,
            ErrorCode: 500,
            ErrorBody: err,
        }
    } else {
        statusCode = http.StatusOK
        response = api.SuccessResponse {
            Success: true,
            Response: passengers,
        }
    }

    c.JSON(statusCode, response)
}

func createPassenger(c *gin.Context) {
    var response any
    var statusCode int

    var body passengerHandler.PassengerStructInsert
    parseErr := c.ShouldBind(&body)

    if parseErr != nil {
        statusCode = http.StatusBadRequest
        errMessage := "Failed to parse request body"
        log.Debug(errMessage + ": " + parseErr.Error())
        error := intError{Type: internalerror.ErrorParseError, Message: errMessage, Body: parseErr}

        response = api.ErrorResponse {
            Success: false,
            ErrorBody: error,
        }
    } else {
        newPass, err := passengerHandler.CreatePassenger(body)
        if err != nil {
            statusCode = http.StatusInternalServerError
            response = api.ErrorResponse {
                Success: false,
                ErrorCode: 500,
                ErrorBody: err,
            }
        } else  {
            statusCode = http.StatusCreated
            response = api.SuccessResponse {
                Success: true,
                Response: newPass,
            }
        }
    }
    c.JSON(statusCode, response)
}

func updatePassenger(c *gin.Context) {
    var response any
    var statusCode int

    var body passengerHandler.PassengerStructUpdate
    parseErr := c.ShouldBind(&body)

    if parseErr != nil {
        statusCode = http.StatusBadRequest
        errMessage := "Failed to parse request body"
        log.Debug(errMessage + ": " + parseErr.Error())
        error := intError{Type: internalerror.ErrorParseError, Message: errMessage, Body: parseErr}

        response = api.ErrorResponse {
            Success: false,
            ErrorBody: error,
        }
    } else {
        newPass, err := passengerHandler.UpdatePassenger(body)
        if err != nil {
            statusCode = http.StatusInternalServerError
            response = api.ErrorResponse {
                Success: false,
                ErrorCode: 500,
                ErrorBody: err,
            }
        } else  {
            statusCode = http.StatusOK
            response = api.SuccessResponse {
                Success: true,
                Response: newPass,
            }
        }
    }
    c.JSON(statusCode, response)
}

func deletePassenger(c *gin.Context) {
    var response any
    var statusCode int

    idStr := c.Param("id")

    id, err := strconv.Atoi(idStr)

    if err != nil {
        statusCode = http.StatusBadRequest
        errMessage := "Failed to parse id"
        log.Debug(errMessage + ": " + err.Error())
        error := intError{Type: internalerror.ErrorParseError, Message: errMessage, Body: err}
        response = api.ErrorResponse {
            Success: false,
            ErrorCode: 400,
            ErrorBody: error,
        }
    } else {
        err := passengerHandler.DeletePassenger(id)
        if err != nil {
            statusCode = http.StatusInternalServerError
            response = api.ErrorResponse {
                Success: false,
                ErrorCode: 500,
                ErrorBody: err,
            }
        } else  {
            statusCode = http.StatusNoContent
        }
    }
    c.JSON(statusCode, response)
} 
