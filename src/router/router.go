package router

import (
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api/debug"
	"github.com/gin-gonic/gin"
)

// InitRouter initializes the router
func InitRouter() *gin.Engine {
    r := gin.New()
    r.Use(gin.Logger())
    r.Use(gin.Recovery())

    api_debug := r.Group("debug")
    debug.SetupDebugRoutes(api_debug)

    return r
}

