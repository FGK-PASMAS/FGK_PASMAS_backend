package routes

import (
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api/routes/debug"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api/routes/division"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api/routes/flight"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api/routes/passenger"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api/routes/plane"
	"github.com/gin-gonic/gin"
)

func InitApiRoutes(r *gin.RouterGroup) {
    divisionRoutes := r.Group("divisions")
    division.SetupDivisionRoutes(divisionRoutes)

    passengerRoutes := r.Group("passengers")
    passenger.SetupPassengerRoutes(passengerRoutes)

    flightRoutes := r.Group("flights")
    flight.SetupFlightRoutes(*flightRoutes)

    planeRoutes := r.Group("planes")
    plane.SetupPlaneRoutes(planeRoutes)

    api_debug := r.Group("debug")
    debug.SetupDebugRoutes(api_debug)

}
