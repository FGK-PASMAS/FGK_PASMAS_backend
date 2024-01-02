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
