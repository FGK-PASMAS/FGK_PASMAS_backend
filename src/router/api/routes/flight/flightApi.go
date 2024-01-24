package flight

import (
	"strconv"

	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api"
	pasmasservice "github.com/MetaEMK/FGK_PASMAS_backend/service/pasmasService"
	"github.com/gin-gonic/gin"
)

func getFlights(c *gin.Context) {
    var response interface{}
    var httpCode int = 500

    flights, err := pasmasservice.GetFlights()


    if err != nil {
        res := api.GetErrorResponse(err)
        c.JSON(res.HttpCode, res.ErrorResponse)
    } else {
        response = flights
        httpCode = 200
    }

    c.JSON(httpCode, response)
}

func createFlight(c *gin.Context) {
    var response interface{}
    var httpCode int = 500

    flight := model.Flight{}
    c.ShouldBind(&flight)

    var newFlight *model.Flight
    var err error

    switch flight.Type {
        case model.FsReserved:
            newFlight, err = pasmasservice.ReserveFlight(&flight)
        case model.FsBlocked:
            err = api.ErrNotImplemented
        case model.FsBooked:
            err = api.ErrNotImplemented
        default: 
            err = api.ErrInvalidFlightType
    }

    if err != nil {
        res := api.GetErrorResponse(err)
        httpCode = res.HttpCode
        response = res.ErrorResponse
    } else {
        response = api.SuccessResponse {
            Success: true,
            Response: newFlight,
        }
        httpCode = 201
    }

    c.JSON(httpCode, response)
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
    id, err := strconv.ParseUint(idStr, 10, 64)

    if err != nil {
        res := api.GetErrorResponse(err)
        c.JSON(res.HttpCode, res.ErrorResponse)
    } else {
        err := pasmasservice.DeleteFlights(uint(id))
        if err != nil {
            res := api.GetErrorResponse(err)
            c.JSON(res.HttpCode, res.ErrorResponse)
        } else {
            c.JSON(204, nil)
        }
    }
}
