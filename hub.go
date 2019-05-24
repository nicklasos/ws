package main

// Hub manages incoming connections
type Hub struct {
	clients map[*Client]bool

	broadcast  chan []byte
	chat       chan *ChatMessage
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		chat:       make(chan *ChatMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) registerClient(client *Client) {
	h.clients[client] = true
}

func (h *Hub) disconnect(client *Client) {
	delete(h.clients, client)
	close(client.send)
}

func (h *Hub) send(client *Client, msg []byte) {
	select {
	case client.send <- msg:
	default:
		h.disconnect(client)
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				h.disconnect(client)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				h.send(client, message)
			}
		case chat := <-h.chat:
			for client := range h.clients {
				// Todo: optimize this
				if inArray(chat.room, client.rooms) {
					h.send(client, chat.messageRaw)
				}
			}
		}
	}
}

func inArray(value string, array []string) bool {
	for _, val := range array {
		if val == value {
			return true
		}
	}

	return false
}
