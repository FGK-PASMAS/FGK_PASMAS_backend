package pilots

import "github.com/gin-gonic/gin"

func SetupPilotRoutes(gr *gin.RouterGroup) {
    gr.GET("", getPilots)
}
