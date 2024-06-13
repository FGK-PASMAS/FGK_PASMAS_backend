package flight

import (
	"net/http"
	"strconv"

	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	databasehandler "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api"
	"github.com/MetaEMK/FGK_PASMAS_backend/service/flightService"
	"github.com/gin-gonic/gin"
)

func getFlights(c *gin.Context) {
	var response interface{}
	var httpCode int = 500
	var err error
	var flights []model.Flight

	includes, incErr := databasehandler.ParseFlightInclude(c)
	filters, filtErr := databasehandler.ParseFlightFilter(c)

	if incErr == nil && filtErr == nil {
		flights, err = flightService.GetFlights(includes, filters)
	} else {
		if incErr != nil {
			err = incErr
		} else {
			err = filtErr
		}
	}

	if err != nil {
		e := cerror.InterpretError(err)
		httpCode = e.HttpCode
		response = e
	} else {
		response = api.SuccessResponse{
			Success:  true,
			Response: flights,
		}
		httpCode = http.StatusOK
	}

	c.JSON(httpCode, response)
}

func flightCreation(c *gin.Context) {
	var response interface{}
	var httpCode int

	var flight model.Flight
	err := c.ShouldBind(&flight)

	user := c.Keys["user"].(model.UserJwtBody)

	passengers := flight.Passengers
	flight.Passengers = nil

	var newFlight model.Flight

	if flight.Status == model.FsBlocked {
		newFlight, err = flightService.CreateBlocker(user, flight)
	} else {
		newFlight, _, err = flightService.FlightCreation(user, flight, passengers)
	}

	if err != nil {
		e := cerror.InterpretError(err)
		httpCode = e.HttpCode
		response = e
	} else {
		response = api.SuccessResponse{
			Success:  true,
			Response: newFlight,
		}
		httpCode = http.StatusCreated
	}

	c.JSON(httpCode, response)
}

func flightUpdate(c *gin.Context) {
	var response interface{}
	var httpCode int
	var err error
	var newFlight model.Flight

	var flight model.Flight
	err = c.ShouldBind(&flight)

	user := c.Keys["user"].(model.UserJwtBody)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err == nil {
		newFlight, err = flightService.FlightBooking(user, uint(id), flight)
	}

	if err != nil {
		e := cerror.InterpretError(err)
		httpCode = e.HttpCode
		response = e
	} else {
		response = api.SuccessResponse{
			Success:  true,
			Response: newFlight,
		}
		httpCode = http.StatusOK
	}

	c.JSON(httpCode, response)
}

func deleteFlight(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)

	user := c.Keys["user"].(model.UserJwtBody)

	if err != nil {
		res := cerror.InterpretError(err)
		c.JSON(res.HttpCode, res)
	} else {
		err := flightService.DeleteFlights(user, uint(id))
		if err != nil {
			res := cerror.InterpretError(err)
			c.JSON(res.HttpCode, res)
		} else {
			c.JSON(204, nil)
		}
	}
}
