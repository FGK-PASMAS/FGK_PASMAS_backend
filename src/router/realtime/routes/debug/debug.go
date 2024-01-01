package debug

import (
	"fmt"
	"time"

	"github.com/MetaEMK/FGK_PASMAS_backend/router/realtime"
	"github.com/gin-gonic/gin"
)


var DebugStream = realtime.NewStream()
var pingStream = realtime.NewStream()

func debug(r *gin.RouterGroup) {
    r.GET("stream", realtime.HeadersMiddleware(), DebugStream.ServeStream(), func(c *gin.Context) {
        realtime.StreamToClient(c)
    })

    r.POST("stream", func(c *gin.Context) {
        fmt.Println("Sending event to stream")
        DebugStream.SendEvent("Hello, world!")
        fmt.Println("Event sent to", DebugStream.TotalClients, "clients")
        for _, client := range DebugStream.TotalClients {
            fmt.Println("Client", client, "is connected:")
        }
    })
}

func ping(r *gin.RouterGroup) {
    r.GET("ping", realtime.HeadersMiddleware(), pingStream.ServeStream(), func(c *gin.Context) {
        realtime.StreamToClient(c)
    })
}

func sendPings() {
    for {
        var timestamp = time.Now().Format("15:04:05")
        pingStream.SendEvent("Ping at " + timestamp)
        time.Sleep(2 * time.Second)
    }
}
