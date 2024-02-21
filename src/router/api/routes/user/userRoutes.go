package user

import "github.com/gin-gonic/gin"

func InitUserRoutes(r *gin.RouterGroup) {
    r.POST("", ValidateUser)
}
