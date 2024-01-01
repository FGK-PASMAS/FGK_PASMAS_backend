package debug

import "github.com/gin-gonic/gin"

func SetupDebugRoutes(r *gin.RouterGroup) {
    debug(r)
    ping(r)

    go sendPings()
}
