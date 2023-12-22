package debug

import (
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

type Event struct {
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
var debugStream = NewStream()

func debug(r *gin.RouterGroup) {
    r.GET("stream", HeadersMiddleware(), debugStream.addClient(), func(c *gin.Context) {
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

    })

    r.POST("stream", func(c *gin.Context) {
        debugStream.Message <- "test"
    })
}

// addClient adds a new client to the clients map.
func (stream *Event) addClient() gin.HandlerFunc {
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


// NewStream creates a new event server and returns it.
func NewStream() (event *Event){
    event = &Event{
        Message:       make(chan string),
        NewClients:    make(chan chan string),
        ClosedClients: make(chan chan string),
        TotalClients:  make(map[chan string]bool),
    }

    go event.listen()
    fmt.Println("New event server created")

    return
}


func (stream *Event) listen() {
    for {
        select {
        case client := <-stream.NewClients:
            stream.TotalClients[client] = true
            fmt.Println("Client added. %d registered clients", len(stream.TotalClients))

        case client := <-stream.ClosedClients:
            delete(stream.TotalClients, client)
            close(client)
            fmt.Println("Removed client. %d registered clients", len(stream.TotalClients))

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
