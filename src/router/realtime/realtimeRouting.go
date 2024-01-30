package realtime

import (
	"time"

	"github.com/gin-gonic/gin"
)


var FlightStream = NewStream()
var PassengerStream = NewStream()
var pingStream = NewStream()

func SetupRealtimeRoutes(r *gin.RouterGroup) {
    subscribeToStream(r, "/passengers", PassengerStream)
    subscribeToStream(r, "/flights", FlightStream)
    subscribeToStream(r, "/pings", pingStream)

    go sendPings()
}

func subscribeToStream(r *gin.RouterGroup, url string, stream *Stream) {
    r.GET(url, HeadersMiddleware(), stream.ServeStream(), StreamToClient)
}

func sendPings() {
    for {
        var res struct {
            Description string
            Value       string
        }

        res.Description = "The The current time is:"
        res.Value = time.Now().In(time.UTC).String()
        pingStream.PublishEvent(PING, res)

        time.Sleep(2 * time.Second)
    }
}
