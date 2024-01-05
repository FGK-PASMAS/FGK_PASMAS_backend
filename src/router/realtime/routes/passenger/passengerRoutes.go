package passenger

import (
	"github.com/MetaEMK/FGK_PASMAS_backend/router/realtime"
	"github.com/gin-gonic/gin"
)

var PassengerStream = realtime.NewStream()

func SubscribeToPassenger(r *gin.RouterGroup) {
    r.GET("/", realtime.HeadersMiddleware(), PassengerStream.ServeStream(), realtime.StreamToClient)
}
