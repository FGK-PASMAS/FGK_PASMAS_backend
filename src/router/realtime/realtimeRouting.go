package realtime

import "github.com/gin-gonic/gin"


var FlightStream = NewStream()
var PassengerStream = NewStream()

func SetupRealtimeRoutes(r *gin.RouterGroup) {
    subscribeToStream(r, "/passengers", PassengerStream)
    subscribeToStream(r, "/flights", FlightStream)
}

func subscribeToStream(r *gin.RouterGroup, url string, stream *Stream) {
    r.GET(url, HeadersMiddleware(), stream.ServeStream(), StreamToClient)
}
