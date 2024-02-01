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

func flightPlanning(c *gin.Context) {
    var response interface{}
    var httpCode int

    var input FlightPlanning
    err := c.ShouldBind(&input)

    flight, err := pasmasservice.FlightPlanning(input.PlaneId, input.DepartureTime, input.ArrivalTime, input.Description)

    if err != nil {
        res := api.GetErrorResponse(err)
        httpCode = res.HttpCode
        response = res.ErrorResponse
    } else {
        response = api.SuccessResponse {
            Success: true,
            Response: flight,
        }
        httpCode = http.StatusCreated
    }

    c.JSON(httpCode, response)
}

func flightReservation(c *gin.Context) {
    var response interface{}
    var httpCode int
    var err error
    var flight *model.Flight

    var input FlightReservation
    err = c.ShouldBind(&input)

    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err == nil {
        flight, err = pasmasservice.FlightReservation(uint(id), &input.Passengers, input.Description)
    }

    if err != nil {
        res := api.GetErrorResponse(err)
        httpCode = res.HttpCode
        response = res.ErrorResponse
    } else {
        response = api.SuccessResponse {
            Success: true,
            Response: flight,
        }
        httpCode = http.StatusCreated
    }

    c.JSON(httpCode, response)
}

func flightBooking(c *gin.Context) {
    var flight *model.Flight
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 64)

    var input FlightBooking
    c.ShouldBind(&input)

    if err == nil {
        flight, err = pasmasservice.BookFlight(uint(id), &input.Passengers, input.Description)
    }

    if err != nil {
        res := api.GetErrorResponse(err)
        c.JSON(res.HttpCode, res.ErrorResponse)
    } else {
        response := api.SuccessResponse {
            Success: true,
            Response: flight,
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
