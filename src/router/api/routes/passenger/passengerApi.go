package passenger

import (
	"net/http"
	"strconv"

	"github.com/MetaEMK/FGK_PASMAS_backend/logging"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api"
	pasmasservice "github.com/MetaEMK/FGK_PASMAS_backend/service/pasmasService"
	"github.com/gin-gonic/gin"
)

var log = logging.ApiLogger

func getPassengers(c *gin.Context) {
    passengers, err := pasmasservice.GetPassengers()

    if err != nil {
        apiErr := api.GetErrorResponse(err)
        apiErr.ErrorResponse.Message = err.Error()
        c.JSON(apiErr.HttpCode, apiErr.ErrorResponse)
    } else {
        c.JSON(http.StatusOK, api.SuccessResponse { Success: true, Response: passengers })
    }

}

func createPassenger(c *gin.Context) {
    var body model.Passenger
    parseErr := c.ShouldBind(&body)

    if parseErr != nil {
        apiErr := api.InvalidRequestBody
        apiErr.ErrorResponse.Message = parseErr.Error()
        c.JSON(apiErr.HttpCode, apiErr.ErrorResponse)
        return
    }

    newPass, err := pasmasservice.CreatePassenger(body)
    if err != nil {
        apiErr := api.GetErrorResponse(err)
        apiErr.ErrorResponse.Message = err.Error()
        c.JSON(apiErr.HttpCode, apiErr.ErrorResponse)
    } else {
        c.JSON(http.StatusCreated, api.SuccessResponse { Success: true, Response: newPass })
    }

}

func updatePassenger(c *gin.Context) {
    var body model.Passenger
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        apiErr := api.InvalidRequestBody
        apiErr.ErrorResponse.Message = "Failed to read the id as int64 from url paramenters"
        c.JSON(apiErr.HttpCode, apiErr.ErrorResponse)
        return
    }

    parseErr := c.ShouldBind(&body)
    if parseErr != nil {
        apiErr := api.InvalidRequestBody
        apiErr.ErrorResponse.Message = parseErr.Error()
        c.JSON(apiErr.HttpCode, apiErr.ErrorResponse)
        return
    }

    newPass, err := pasmasservice.UpdatePassenger(uint(id), body)
    if err != nil {
        apiErr := api.GetErrorResponse(err)
        apiErr.ErrorResponse.Message = err.Error()
        c.JSON(apiErr.HttpCode, apiErr.ErrorResponse)
    } else {
        c.JSON(http.StatusOK, api.SuccessResponse { Success: true, Response: newPass })
    }
}

func deletePassenger(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseInt(idStr, 10, 64)

    if err != nil {
        apiErr := api.InvalidRequestBody
        apiErr.ErrorResponse.Message = "Failed to read the id as int64 from url paramenters"
        c.JSON(apiErr.HttpCode, apiErr.ErrorResponse)
        return
    }

    err = pasmasservice.DeletePassenger(uint(id))
    if err != nil {
        apiErr := api.GetErrorResponse(err)
        apiErr.ErrorResponse.Message = err.Error()
        c.JSON(apiErr.HttpCode, apiErr.ErrorResponse)
    } else {
        c.JSON(http.StatusNoContent, nil)
    }

} 
