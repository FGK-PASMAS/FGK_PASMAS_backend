package realtime

import "github.com/MetaEMK/FGK_PASMAS_backend/model"

var flightStreamByDivision = map[uint]*Stream{}

func InitAllFlightByDivisionEndpoints(divisions []model.Division) {
    for k := range flightStreamByDivision {
        DeleteFlightStreamEndpoint(k)
    }
    
    for _, d := range divisions {
        AddFlightStreamWithDivisionId(d.ID)
    }
}

// GetStreamForDivisionId return the correct stream for this flight.
// returns nil if no Stream is found
func GetFlightStreamForDivisionId(id uint) (stream *Stream) {
    stream = flightStreamByDivision[id]
    return
}

func AddFlightStreamWithDivisionId(id uint) (ok bool) {
    if GetFlightStreamForDivisionId(id) == nil {
        flightStreamByDivision[id] = newStream()
        ok = true
    }

    return
}

func DeleteFlightStreamEndpoint(id uint) (ok bool) {
    stream := GetFlightStreamForDivisionId(id)
    if stream != nil {
        stream.publishEvent(realtimeEvent{
            Stream: stream,
            Action: OTHER,
            Data: "Stream will be deleted",
        })

        flightStreamByDivision[id] = nil
    }

    return
}
