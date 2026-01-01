// Package ws provides WebSocket functionality for real-time communication.
package ws

import (
	"sync"
	"time"
)

// Hub maintains the set of active clients and broadcasts messages to clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Mutex for thread-safe operations
	mu sync.RWMutex

	// Metrics
	metrics *HubMetrics
}

// HubMetrics contains WebSocket hub statistics.
type HubMetrics struct {
	TotalConnections  int64     `json:"totalConnections"`
	ActiveConnections int       `json:"activeConnections"`
	MessagesSent      int64     `json:"messagesSent"`
	MessagesReceived  int64     `json:"messagesReceived"`
	LastActivityAt    time.Time `json:"lastActivityAt"`
}

// NewHub creates a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		metrics: &HubMetrics{
			LastActivityAt: time.Now(),
		},
	}
}

// Run starts the hub's main loop.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.metrics.TotalConnections++
			h.metrics.ActiveConnections = len(h.clients)
			h.metrics.LastActivityAt = time.Now()
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				h.metrics.ActiveConnections = len(h.clients)
				h.metrics.LastActivityAt = time.Now()
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			h.metrics.MessagesSent++
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					// Client's send buffer is full, close it
					go func(c *Client) {
						h.unregister <- c
					}(client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Broadcast sends a message to all connected clients.
func (h *Hub) Broadcast(message []byte) {
	h.broadcast <- message
}

// BroadcastToRole sends a message to clients with a specific role.
func (h *Hub) BroadcastToRole(message []byte, role string) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.Role == role {
			select {
			case client.send <- message:
			default:
				go func(c *Client) {
					h.unregister <- c
				}(client)
			}
		}
	}
}

// BroadcastToUser sends a message to a specific user.
func (h *Hub) BroadcastToUser(message []byte, userID uint64) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.UserID == userID {
			select {
			case client.send <- message:
			default:
				go func(c *Client) {
					h.unregister <- c
				}(client)
			}
		}
	}
}

// GetMetrics returns current hub metrics.
func (h *Hub) GetMetrics() *HubMetrics {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return &HubMetrics{
		TotalConnections:  h.metrics.TotalConnections,
		ActiveConnections: h.metrics.ActiveConnections,
		MessagesSent:      h.metrics.MessagesSent,
		MessagesReceived:  h.metrics.MessagesReceived,
		LastActivityAt:    h.metrics.LastActivityAt,
	}
}

// GetOnlineCount returns the number of currently connected clients.
func (h *Hub) GetOnlineCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// GetOnlineCountByRole returns online counts grouped by role.
func (h *Hub) GetOnlineCountByRole() map[string]int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	counts := make(map[string]int)
	for client := range h.clients {
		counts[client.Role]++
	}
	return counts
}

// ClientCount returns the total number of connected clients.
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
