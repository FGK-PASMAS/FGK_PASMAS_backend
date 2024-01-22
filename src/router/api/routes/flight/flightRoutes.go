package flight

import "github.com/gin-gonic/gin"

func SetupFlightRoutes(gr gin.RouterGroup) {
    gr.GET("", getFlights)
    gr.POST("", createFlight)
    gr.PUT("/:id", updateFlight)
    gr.DELETE("/:id", deleteFlight)
}
