package routes

import (
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api/routes/auth"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api/routes/debug"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api/routes/division"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api/routes/flight"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api/routes/passenger"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api/routes/pilots"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api/routes/plane"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api/routes/user"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/middleware"
	"github.com/gin-gonic/gin"
)

func InitApiRoutes(r *gin.RouterGroup) {
    authRoutes := r.Group("auth")
    auth.InitAuthRoutes(authRoutes)

    userRoutes := r.Group("users")
    userRoutes.Use(middleware.ValidateJwt)
    user.InitUserRoutes(userRoutes)

    divisionRoutes := r.Group("divisions")
    divisionRoutes.Use(middleware.ValidateJwt)
    division.SetupDivisionRoutes(divisionRoutes)

    passengerRoutes := r.Group("passengers")
    passengerRoutes.Use(middleware.ValidateJwt)
    passenger.SetupPassengerRoutes(passengerRoutes)

    flightRoutes := r.Group("flights")
    flightRoutes.Use(middleware.ValidateJwt)
    flight.SetupFlightRoutes(*flightRoutes)

    planeRoutes := r.Group("planes")
    planeRoutes.Use(middleware.ValidateJwt)
    plane.SetupPlaneRoutes(planeRoutes)

    pilotRoutes := r.Group("pilots")
    pilotRoutes.Use(middleware.ValidateJwt)
    pilots.SetupPilotRoutes(pilotRoutes)

    api_debug := r.Group("debug")
    debug.SetupDebugRoutes(api_debug)
}
