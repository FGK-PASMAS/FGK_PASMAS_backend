package passenger

import "github.com/gin-gonic/gin"

func SetupPassengerRoutes(r *gin.RouterGroup) {
    r.GET("", getPassengers)
    // r.POST("", createPassenger)
    // r.PUT("/:id", updatePassenger)
    // r.DELETE("/:id", deletePassenger)
}
