package router

import (
	apiRoutes "github.com/MetaEMK/FGK_PASMAS_backend/router/api/routes"
	realtimeRoutes "github.com/MetaEMK/FGK_PASMAS_backend/router/realtime/routes"
	"github.com/gin-gonic/gin"
)

// InitRouter initializes the router
func InitRouter() *gin.Engine {
    r := gin.New()
    r.Use(gin.Logger())
    r.Use(gin.Recovery())

    apiRouter := r.Group("/api")
    realtimeRouter := r.Group("/realtime")

    apiRoutes.InitApiRoutes(apiRouter)
    realtimeRoutes.InitRealtimeRoutes(realtimeRouter)

    return r
}

