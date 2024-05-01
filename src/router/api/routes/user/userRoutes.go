package user

import (
	"github.com/gin-gonic/gin"
)

func InitUserRoutes(r *gin.RouterGroup) {
    r.GET("", getAllUsers)
    r.POST("/createUser", createNewUser)
    r.DELETE("/:id", deleteUser)
}
