package routes

import (
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api/routes/debug"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api/routes/division"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api/routes/passenger"
	"github.com/gin-gonic/gin"
)

func InitApiRoutes(r *gin.RouterGroup) {
    divisionRoutes := r.Group("division")
    division.SetupDivisionRoutes(divisionRoutes)

    api_debug := r.Group("debug")
    debug.SetupDebugRoutes(api_debug)

    passengerRoutes := r.Group("passenger")
    passenger.SetupPassengerRoutes(passengerRoutes)

}
