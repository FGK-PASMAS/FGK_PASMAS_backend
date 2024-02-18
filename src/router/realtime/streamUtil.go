package realtime

import (
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

type Stream struct {
	// Events are pushed to this channel by the main events-gathering routine
	message chan string

	// New client connections
	newClients chan chan string

	// Closed client connections
	closedClients chan chan string

	// Total client connections
	totalClients map[chan string]bool

}


type clientChan chan string

// addClient adds a new client to the clients map.
func (stream *Stream) serveStream() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Init client channel
        clientChan := make(clientChan)

        //Send new client to event server
        stream.newClients <- clientChan

        defer func() {
            stream.closedClients <- clientChan
        }()

        c.Set("clientChan", clientChan)

        c.Next()
    }
}

// SendEvent sends an event to all clients in this stream.
func (stream *Stream) sendEvent(eventMessage string) {
    stream.message <- eventMessage
}

// newStream creates a new event server and returns it.
func newStream() (event *Stream){
    event = &Stream{
        message:       make(chan string),
        newClients:    make(chan chan string),
        closedClients: make(chan chan string),
        totalClients:  make(map[chan string]bool),
    }

    go event.listen()
    return
}

func streamToClient(c *gin.Context) {
        v, err := c.Get("clientChan")
        if !err {
            fmt.Println("Error getting clientChan")
            fmt.Printf("%v \n", err)
            return
        }

        clientChan, err := v.(clientChan)
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
        case client := <-stream.newClients:
            fmt.Println("New client", client)
            stream.totalClients[client] = true

        case client := <-stream.closedClients:
            delete(stream.totalClients, client)
            close(client)

        case eventMessage := <-stream.message:
            for clientMessageChan := range stream.totalClients {
                clientMessageChan <- eventMessage
            }
        }
    }
}

// headersMiddleware sets headers for server-sent events.
func headersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Next()
	}
}

func (s *Stream) publishEvent(event realtimeEvent) {
    bodyString := event.toJson()
    s.sendEvent(bodyString)
}
