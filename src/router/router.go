package router

import (
	apiRoutes "github.com/MetaEMK/FGK_PASMAS_backend/router/api/routes"
	realtimeRoutes "github.com/MetaEMK/FGK_PASMAS_backend/router/realtime/routes"
	"github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
)

// InitRouter initializes the router
func InitRouter() *gin.Engine {
    r := gin.New()
    r.Use(gin.Logger())
    r.Use(gin.Recovery())
    r.Use(cors.Default())

    apiRouter := r.Group("/api")
    realtimeRouter := r.Group("/realtime")

    apiRoutes.InitApiRoutes(apiRouter)
    realtimeRoutes.InitRealtimeRoutes(realtimeRouter)

    return r
}

