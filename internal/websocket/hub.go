package websocket

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	Clients map[*Client]bool

	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
}

// NewHub initializes a new Hub instance
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}
}

// Run starts the hub to listen for register, unregister, and broadcast events
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			delete(h.Clients, client)
			close(client.Send)
		case msg := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- msg:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}
