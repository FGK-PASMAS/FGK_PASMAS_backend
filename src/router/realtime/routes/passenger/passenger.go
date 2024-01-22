package passenger

import (
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/realtime"
)

func PublishPassengerEvent(actionType realtime.ActionType, data *model.Passenger) {
    body := realtime.RealtimeBodyModel {
        Action: actionType,
        Data: data,
    }

    bodyString := body.ToJson()

    PassengerStream.SendEvent(bodyString)
}
