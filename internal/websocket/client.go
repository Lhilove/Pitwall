package websocket

import (
	"github.com/gorilla/websocket"
)

// Client represents a single WebSocket connection to the server
type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
}

// WritePump listens for incoming messages from the WebSocket connection
func (c *Client) WritePump() {
	defer c.Conn.Close()
	for msg := range c.Send {
		c.Conn.WriteMessage(websocket.TextMessage, msg)
	}
}
