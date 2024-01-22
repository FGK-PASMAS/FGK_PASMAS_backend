package flight

import (
	"strconv"

	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api"
	"github.com/gin-gonic/gin"
)

func getFlights(c *gin.Context) {
    flights := model.Flight{}

    response := api.SuccessResponse {
        Success: true,
        Response: flights,
    }

    c.JSON(200, response)
}

func createFlight(c *gin.Context) {
    flight := model.Flight{}
    c.ShouldBind(&flight)

    response := api.SuccessResponse {
        Success: true,
        Response: flight,
    }

    c.JSON(201, response)
}

func updateFlight(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 64)

    flight := model.Flight{}
    c.ShouldBind(&flight)

    if err != nil {
        res := api.GetErrorResponse(err)
        c.JSON(res.HttpCode, res.ErrorResponse)
    } else {
        flight.ID = uint(id)
        response := api.SuccessResponse {
            Success: true,
            Response: flight,
        }
        c.JSON(200, response)
    }
}

func deleteFlight(c *gin.Context) {
    idStr := c.Param("id")
    _, err := strconv.ParseUint(idStr, 10, 64)

    if err != nil {
        res := api.GetErrorResponse(err)
        c.JSON(res.HttpCode, res.ErrorResponse)
    } else {
        c.JSON(204, nil)
    }
}
