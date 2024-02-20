package realtime

import (
	"time"

	"github.com/gin-gonic/gin"
)

var FlightStream = newStream()
var PassengerStream = newStream()
var pingStream = newStream()

func SetupRealtimeRoutes(r *gin.RouterGroup) {
    subscribeToStream(r, "/passengers", PassengerStream)
    subscribeToStream(r, "/pings", pingStream)
    subscribeToFlightByDivisionEndpoint(r)

    go sendPings()
}

func subscribeToStream(r *gin.RouterGroup, url string, stream *Stream) {
    r.GET(url, headersMiddleware(), stream.serveStream(), streamToClient)
}

func sendPings() {
    for {
        var res struct {
            Description string
            Value       string
        }
        res.Description = "The The current time is:"
        res.Value = time.Now().In(time.UTC).String()

        event := realtimeEvent {
            Stream: pingStream,
            Action: OTHER,
            Data: res,
        }

        event.publishEvent()
        time.Sleep(2 * time.Second)
    }
}
