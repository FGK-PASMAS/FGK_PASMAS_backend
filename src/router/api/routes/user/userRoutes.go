package user

import (
	"github.com/MetaEMK/FGK_PASMAS_backend/router/middleware"
	"github.com/gin-gonic/gin"
)

func InitUserRoutes(r *gin.RouterGroup) {
    r.POST("", validateUser)
    r.POST("/createUser", middleware.ValidateJwt, createNewUser)
}
