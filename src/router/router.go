package router

import (
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/realtime"
	"github.com/gin-gonic/gin"
)

// InitRouter initializes the router
func InitRouter() *gin.Engine {
    r := gin.New()
    r.Use(gin.Logger())
    r.Use(gin.Recovery())

    apiRouter := r.Group("/api")
    realtimeRouter := r.Group("/realtime")

    api.InitApiRoutes(apiRouter)
    realtime.InitRealtimeRoutes(realtimeRouter)

    return r
}

