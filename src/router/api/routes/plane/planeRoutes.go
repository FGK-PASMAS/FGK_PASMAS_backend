package plane

import "github.com/gin-gonic/gin"

func SetupPlaneRoutes(r *gin.RouterGroup) {
    r.GET("", getPlanes)
    r.GET("/:id", getPlaneById)
    r.PATCH("/:id", updatePlane)
}

