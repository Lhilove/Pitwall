package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	ws "github.com/lhilove/pitwall/internal/websocket"
)

// upgrader upgrades a normal HTTP connection to a WebSocket connection
// CheckOrigin returns true to allow all connections (you may want to restrict this in production)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// ServeWs upgrades the HTTP connection and registers the client with the hub
func ServeWs(hub *ws.Hub, c *gin.Context) {
	// upgrade HTTP to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	// create a new client for this connection
	client := &ws.Client{
		Hub:  hub,
		Conn: conn,
		Send: make(chan []byte, 256), // buffered so slow clients don't block others
	}

	// register the client with the hub
	hub.Register <- client

	// start listening for messages to send to this client
	go client.WritePump()
}
