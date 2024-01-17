package passenger

import (
	"github.com/MetaEMK/FGK_PASMAS_backend/router/realtime"
)

func PublishPassengerEvent(actionType realtime.ActionType, data interface{}) {
    body := realtime.RealtimeBodyModel {
        Action: actionType,
        Data: data,
    }

    bodyString := body.ToJson()

    PassengerStream.SendEvent(bodyString)
}
