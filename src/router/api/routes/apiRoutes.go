package routes

import (
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api/routes/debug"
	"github.com/gin-gonic/gin"
)

func InitApiRoutes(r *gin.RouterGroup) {
    api_debug := r.Group("debug")
    debug.SetupDebugRoutes(api_debug)
}
