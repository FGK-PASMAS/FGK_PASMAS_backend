package plane

import "github.com/gin-gonic/gin"

func SetupPlaneRoutes(r *gin.RouterGroup) {
    r.GET("", getPlanes)
}
