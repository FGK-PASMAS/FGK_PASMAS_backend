package flight

import (
	"net/http"
	"strconv"

	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api"
	pasmasservice "github.com/MetaEMK/FGK_PASMAS_backend/service/pasmasService"
	"github.com/gin-gonic/gin"
)

func getFlights(c *gin.Context) {
    var response interface{}
    var httpCode int = 500
    var err error
    var flights *[]model.Flight

    includes, incErr := pasmasservice.ParseFlightInclude(c)
    filters, filtErr := pasmasservice.ParseFlightFilter(c)

    if incErr == nil && filtErr == nil {
        flights, err = pasmasservice.GetFlights(includes, filters)
    } else {
        if incErr != nil {
            err = incErr
        } else {
            err = filtErr
        }
    }

    if err != nil {
        res := api.GetErrorResponse(err)
        response = res.ErrorResponse
        httpCode = res.HttpCode
    } else {
        response = api.SuccessResponse {
            Success: true,
            Response: flights,
        }
        httpCode = http.StatusOK
    }

    c.JSON(httpCode, response)
}

func createFlight(c *gin.Context) {
    var response interface{}
    var httpCode int

    flight := model.Flight{}
    c.ShouldBind(&flight)

    var newFlight *model.Flight
    var err error

    switch flight.Status {
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
        httpCode = http.StatusCreated
    }

    c.JSON(httpCode, response)
}

func bookFlight(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 64)

    flight := model.Flight{}
    c.ShouldBind(&flight)

    var newFlight *model.Flight

    switch flight.Status {
        case model.FsReserved:
            err = api.ErrNotImplemented
        case model.FsBlocked:
            err = api.ErrNotImplemented
        case model.FsBooked:
            newFlight, err = pasmasservice.BookFlight(uint(id), flight.Passengers)
        default: 
            err = api.ErrInvalidFlightType
    }

    if err != nil {
        res := api.GetErrorResponse(err)
        c.JSON(res.HttpCode, res.ErrorResponse)
    } else {
        flight.ID = uint(id)
        response := api.SuccessResponse {
            Success: true,
            Response: newFlight,
        }

        c.JSON(http.StatusOK, response)
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
