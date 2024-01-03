package passenger

import (
	"net/http"

	passengerHandler "github.com/MetaEMK/FGK_PASMAS_backend/database/passengerHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api"
	"github.com/gin-gonic/gin"
)

func getPassengers(c *gin.Context) {
    passengers, err := passengerHandler.GetPassengers()
    var statusCode int
    var response interface{}

    if(err != nil) {
        statusCode = http.StatusInternalServerError
        response = api.ErrorResponse {
            Success: false,
            ErrorCode: 500,
            ErrorBody: err.Error(),
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

    var body passengerHandler.InsertPassenger
    parseErr := c.ShouldBind(&body)

    if parseErr != nil {
        statusCode = http.StatusBadRequest
        response = api.ErrorResponse {
            Success: false,
            ErrorBody: parseErr.Error(),
        }
    } else {
        err := passengerHandler.CreatePassenger(body)
        if err != nil {
            statusCode = 500
            response = err.Error()
        }
        statusCode = http.StatusOK
        response = api.SuccessResponse {
            Success: true,
        }
    }

    c.JSON(statusCode, response)
}
