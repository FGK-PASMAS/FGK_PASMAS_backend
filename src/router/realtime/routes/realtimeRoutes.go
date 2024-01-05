package routes

import (
	"github.com/MetaEMK/FGK_PASMAS_backend/router/realtime/routes/debug"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/realtime/routes/passenger"
	"github.com/gin-gonic/gin"
)


func InitRealtimeRoutes(r *gin.RouterGroup) {
    rt_debug := r.Group("debug")
    debug.SetupDebugRoutes(rt_debug)
    rt_pass := r.Group("passenger")
    passenger.SubscribeToPassenger(rt_pass)
}
