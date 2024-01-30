package realtime

import (
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

type Stream struct {
	// Events are pushed to this channel by the main events-gathering routine
	Message chan string

	// New client connections
	NewClients chan chan string

	// Closed client connections
	ClosedClients chan chan string

	// Total client connections
	TotalClients map[chan string]bool

}


type ClientChan chan string

// addClient adds a new client to the clients map.
func (stream *Stream) ServeStream() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Init client channel
        clientChan := make(ClientChan)

        //Send new client to event server
        stream.NewClients <- clientChan

        defer func() {
            stream.ClosedClients <- clientChan
        }()

        c.Set("clientChan", clientChan)

        c.Next()
    }
}

// SendEvent sends an event to all clients in this stream.
func (stream *Stream) sendEvent(eventMessage string) {
    stream.Message <- eventMessage
}

// NewStream creates a new event server and returns it.
func NewStream() (event *Stream){
    event = &Stream{
        Message:       make(chan string),
        NewClients:    make(chan chan string),
        ClosedClients: make(chan chan string),
        TotalClients:  make(map[chan string]bool),
    }

    go event.listen()
    return
}

func StreamToClient(c *gin.Context) {
        v, err := c.Get("clientChan")
        if !err {
            fmt.Println("Error getting clientChan")
            return
        }

        clientChan, err := v.(ClientChan)
        if !err {
            fmt.Println("Error asserting clientChan")
            return
        }

        c.Stream(func(w io.Writer) bool {
            //Stream message to client
            if msg, ok := <-clientChan; ok {
                c.SSEvent("message", msg)
                return true
            }
            return false
        })
}


func (stream *Stream) listen() {
    for {
        select {
        case client := <-stream.NewClients:
            fmt.Println("New client", client)
            stream.TotalClients[client] = true

        case client := <-stream.ClosedClients:
            delete(stream.TotalClients, client)
            close(client)

        case eventMessage := <-stream.Message:
            for clientMessageChan := range stream.TotalClients {
                clientMessageChan <- eventMessage
            }
        }
    }
}

// HeadersMiddleware sets headers for server-sent events.
func HeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Next()
	}
}

func (s *Stream) PublishEvent(actionType ActionType, data interface{}) {
    body := RealtimeBodyModel {
        Action: actionType,
        Data: data,
    }

    bodyString := body.ToJson()

    s.sendEvent(bodyString)
}
