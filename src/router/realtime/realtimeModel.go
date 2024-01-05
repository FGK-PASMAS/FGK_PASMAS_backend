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
)

var actionTypeStrings = []string{"CREATED", "UPDATED", "DELETED"}

func (a ActionType) MarshalJSON() ([]byte, error) {
    if a < CREATED || a > DELETED {
        return nil, fmt.Errorf("invalid action type")
    }
    return json.Marshal(actionTypeStrings[a])
}

func (rtBody *RealtimeBodyModel) ToJson() (string) {
    jsonBytes, err := json.Marshal(rtBody)

    if err != nil {
        return ""
    }

    return string(jsonBytes)
}
