package ws

import (
	"encoding/json"
	"log"
	"sync"

	"ticketBooking/internal/messaging/dto"
)

// Client represents a single websocket connection of a user.
type Client struct {
	UserID uint
	Send   chan []byte
	Hub    *Hub
}

// type client struct{
// 	UserID uint
// 	Send   chan []byte
// 	clients map[uint]map[*Client]bool

// 	mu      sync.RWMutex

// 	register   chan *Client
// 	unregister chan *Client
// 	broadcast  chan *targetedMessage
// }

// Hub keeps track of all active connections.
// One user can have multiple tabs/devices → multiple clients.
type Hub struct {
	// userID → set of clients
	clients map[uint]map[*Client]bool

	mu sync.RWMutex

	register   chan *Client
	unregister chan *Client
	broadcast  chan *targetedMessage
}

type targetedMessage struct {
	userID  uint
	payload []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uint]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *targetedMessage, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.UserID] == nil {
				h.clients[client.UserID] = make(map[*Client]bool)
			}
			h.clients[client.UserID][client] = true
			h.mu.Unlock()
			log.Printf("[ws] user %d connected (total clients for user: %d)", client.UserID, len(h.clients[client.UserID]))

		case client := <-h.unregister:
			h.mu.Lock()
			if conns, ok := h.clients[client.UserID]; ok {
				if _, exists := conns[client]; exists {
					delete(conns, client)
					close(client.Send)
					if len(conns) == 0 {
						delete(h.clients, client.UserID)
					}
				}
			}
			h.mu.Unlock()
			log.Printf("[ws] user %d disconnected", client.UserID)

		case msg := <-h.broadcast:
			h.mu.RLock()
			conns := h.clients[msg.userID]
			h.mu.RUnlock()
			for client := range conns {
				select {
				case client.Send <- msg.payload:
				default:
					// slow client – drop & cleanup
					h.mu.Lock()
					delete(h.clients[msg.userID], client)
					close(client.Send)
					h.mu.Unlock()
				}
			}
		}
	}
}

// SendToUser pushes a message to all active connections of a user.
func (h *Hub) SendToUser(userID uint, outgoing dto.WSOutgoing) {
	data, err := json.Marshal(outgoing)
	if err != nil {
		return
	}
	h.broadcast <- &targetedMessage{userID: userID, payload: data}
}

// IsOnline checks if a user currently has at least one active connection.
func (h *Hub) IsOnline(userID uint) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients[userID]) > 0
}

func (h *Hub) Register(c *Client)   { h.register <- c }
func (h *Hub) Unregister(c *Client) { h.unregister <- c }
