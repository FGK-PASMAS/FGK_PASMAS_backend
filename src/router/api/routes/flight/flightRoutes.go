package flight

import (
	"github.com/gin-gonic/gin"
)

func SetupFlightRoutes(gr gin.RouterGroup) {
    gr.GET("", getFlights)
    gr.POST("", flightCreation)
    gr.PATCH("/:id", flightUpdate)
    gr.DELETE("/:id", deleteFlight)
}
