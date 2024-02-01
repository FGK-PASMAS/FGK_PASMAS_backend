package flight

import "github.com/gin-gonic/gin"

func SetupFlightRoutes(gr gin.RouterGroup) {
    gr.GET("", getFlights)
    gr.POST("/planning", flightPlanning)
    gr.POST("/reservation/:id", flightReservation)
    gr.POST("/booking/:id", flightBooking)
    gr.DELETE("/:id", deleteFlight)
}
