package realtime

import (
	"encoding/json"
	"fmt"
)

type RealtimeBodyModel struct {
    Action ActionType `json:"action"`
    Data interface{} `json:"data"`
}

type ActionType int

const (
    CREATED ActionType = iota
    UPDATED
    DELETED
    PING
)

var actionTypeStrings = []string{"CREATED", "UPDATED", "DELETED", "PING"}

func (a ActionType) MarshalJSON() ([]byte, error) {
    if a < CREATED || a > PING {
        return nil, fmt.Errorf("invalid action type")
    }
    return json.Marshal(actionTypeStrings[a])
}

func (rtBody *RealtimeBodyModel) ToJson() (string) {
    jsonBytes, err := json.Marshal(rtBody)

    if err != nil {
        fmt.Println("Error marshalling realtime body model")
        return ""
    }

    return string(jsonBytes)
}
