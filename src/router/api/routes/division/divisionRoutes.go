package division

import "github.com/gin-gonic/gin"

func SetupDivisionRoutes(r *gin.RouterGroup) {
    r.GET("/", getDivisions)
}
