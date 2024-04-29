package router

import (
	apiRoutes "github.com/MetaEMK/FGK_PASMAS_backend/router/api/routes"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/realtime"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// InitRouter initializes the router
func InitRouter() *gin.Engine {
    r := gin.New()
    r.Use(gin.Logger())
    r.Use(gin.Recovery())

    corsConfig := cors.DefaultConfig()
    corsConfig.AllowAllOrigins = true
    corsConfig.AllowMethods = append(corsConfig.AllowMethods, "OPTIONS")
    corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
    r.Use(cors.New(corsConfig))

    apiRouter := r.Group("/api")
    realtimeRouter := r.Group("/realtime")

    apiRoutes.InitApiRoutes(apiRouter)
    realtime.SetupRealtimeRoutes(realtimeRouter)

    return r
}

