package realtime

import (
	"encoding/json"
	"fmt"

	"github.com/MetaEMK/FGK_PASMAS_backend/model"
)

type RealtimeHandler struct {
	events     []realtimeEvent
	actionUser model.UserJwtBody
}

type realtimeEvent struct {
	Stream     *Stream           `json:"-"`
	Action     ActionType        `json:"action"`
	ActionUser model.UserJwtBody `json:"actionUser"`
	Data       interface{}       `json:"data"`
}

type ActionType int

const (
	CREATED ActionType = iota
	UPDATED
	DELETED
	OTHER
)

var actionTypeStrings = []string{"CREATED", "UPDATED", "DELETED", "OTHER"}

func (a ActionType) MarshalJSON() ([]byte, error) {
	if a < CREATED || a > OTHER {
		return nil, fmt.Errorf("invalid action type")
	}
	return json.Marshal(actionTypeStrings[a])
}

func (rtBody *realtimeEvent) toJson() string {
	jsonBytes, err := json.Marshal(rtBody)

	if err != nil {
		fmt.Println("Error marshalling realtime body model")
		return ""
	}

	return string(jsonBytes)
}

func NewRealtimeHandler(user model.UserJwtBody) *RealtimeHandler {
	return &RealtimeHandler{
		events: make([]realtimeEvent, 0),
        actionUser: user,
	}
}

func (rt *RealtimeHandler) AddEvent(stream *Stream, action ActionType, data ...interface{}) {
	if stream == nil {
		return
	}

	for _, d := range data {
		rt.events = append(rt.events, realtimeEvent{
			Stream:     stream,
			Action:     action,
			ActionUser: rt.actionUser,
			Data:       d,
		})
	}
}

func (rt *RealtimeHandler) PublishEvents() {
	go func() {
		for _, event := range rt.events {
			event.publishEvent()
		}
	}()
}

func (ev *realtimeEvent) publishEvent() {
	bodyJson := ev.toJson()
	ev.Stream.sendEvent(bodyJson)
}
