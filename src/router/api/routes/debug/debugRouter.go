package debug

import (
	"github.com/MetaEMK/FGK_PASMAS_backend/router/middleware"
	"github.com/gin-gonic/gin"
)

// SetupDebugRoutes sets up the debug routes
func SetupDebugRoutes(g *gin.RouterGroup) {
    g.GET("/ping", middleware.ValidateJwt, ping)
    g.GET("/healthcheck", middleware.ValidateJwt, healthCheck)
    // g.POST("/reset", resetDatabase)
}

