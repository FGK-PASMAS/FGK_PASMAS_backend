package user

import (
	"github.com/gin-gonic/gin"
)

func InitUserRoutes(r *gin.RouterGroup) {
    r.GET("", getAllUsers)
    r.POST("", createNewUser)
    r.DELETE("/:id", deleteUser)
}
