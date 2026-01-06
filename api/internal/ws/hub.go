// Package ws provides WebSocket functionality for real-time communication.
package ws

import (
	"log"
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

	// Redis Pub/Sub for cross-instance broadcasting (optional)
	// If nil, operates in single-instance mode (backward compatible)
	redisPubSub *RedisPubSub
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
// By default, operates in single-instance mode without Redis.
// Use SetRedisPubSub to enable multi-instance support.
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

// SetRedisPubSub sets the Redis Pub/Sub manager for cross-instance broadcasting.
// This enables horizontal scaling by allowing multiple WebSocket server instances
// to broadcast messages to each other's connected clients.
//
// Call this method before starting the hub with Run():
//
//	hub := ws.NewHub()
//	redisPS := ws.NewRedisPubSub(redisClient, hub)
//	hub.SetRedisPubSub(redisPS)
//	go redisPS.Subscribe()
//	go hub.Run()
func (h *Hub) SetRedisPubSub(redisPS *RedisPubSub) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.redisPubSub = redisPS
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

			// Publish presence update if Redis Pub/Sub is enabled
			if h.redisPubSub != nil {
				go func(c *Client) {
					_ = h.redisPubSub.PublishPresence(c.UserID, c.Role, "joined")
				}(client)

			}

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				h.metrics.ActiveConnections = len(h.clients)
				h.metrics.LastActivityAt = time.Now()
			}
			h.mu.Unlock()

			// Publish presence update if Redis Pub/Sub is enabled
			if h.redisPubSub != nil {
				go func(c *Client) {
					_ = h.redisPubSub.PublishPresence(c.UserID, c.Role, "left")
				}(client)
			}

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
// If Redis Pub/Sub is enabled, also broadcasts to all instances.
func (h *Hub) Broadcast(message []byte) {
	// Broadcast to local clients
	h.broadcast <- message

	// If Redis Pub/Sub is enabled, broadcast to all instances
	h.mu.RLock()
	redisPS := h.redisPubSub
	h.mu.RUnlock()

	if redisPS != nil {
		go func() {
			if err := redisPS.Broadcast(message); err != nil {
				log.Printf("Failed to broadcast via Redis: %v", err)
			}
		}()
	}
}

// BroadcastLocal sends a message to local clients only, without publishing to Redis.
// This is used by Redis Pub/Sub handler to avoid infinite loops.
func (h *Hub) BroadcastLocal(message []byte) {
	h.broadcast <- message
}

// BroadcastToRole sends a message to clients with a specific role.
// If Redis Pub/Sub is enabled, also broadcasts to all instances.
func (h *Hub) BroadcastToRole(message []byte, role string) {
	// Broadcast to local clients
	h.broadcastToRoleLocal(message, role)

	// If Redis Pub/Sub is enabled, broadcast to all instances
	h.mu.RLock()
	redisPS := h.redisPubSub
	h.mu.RUnlock()

	if redisPS != nil {
		go func() {
			if err := redisPS.BroadcastToRole(message, role); err != nil {
				log.Printf("Failed to broadcast to role via Redis: %v", err)
			}
		}()
	}
}

// BroadcastToRoleLocal sends a message to local clients with a specific role only.
// This is used by Redis Pub/Sub handler to avoid infinite loops.
func (h *Hub) BroadcastToRoleLocal(message []byte, role string) {
	h.broadcastToRoleLocal(message, role)
}

// broadcastToRoleLocal is the internal implementation for role-based local broadcast.
func (h *Hub) broadcastToRoleLocal(message []byte, role string) {
	h.mu.RLock()
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
	h.mu.RUnlock()
}

// BroadcastToUser sends a message to a specific user.
// If Redis Pub/Sub is enabled, also broadcasts to all instances.
func (h *Hub) BroadcastToUser(message []byte, userID uint64) {
	// Broadcast to local clients
	h.broadcastToUserLocal(message, userID)

	// If Redis Pub/Sub is enabled, broadcast to all instances
	h.mu.RLock()
	redisPS := h.redisPubSub
	h.mu.RUnlock()

	if redisPS != nil {
		go func() {
			if err := redisPS.BroadcastToUser(message, userID); err != nil {
				log.Printf("Failed to broadcast to user via Redis: %v", err)
			}
		}()
	}
}

// BroadcastToUserLocal sends a message to a specific local user only.
// This is used by Redis Pub/Sub handler to avoid infinite loops.
func (h *Hub) BroadcastToUserLocal(message []byte, userID uint64) {
	h.broadcastToUserLocal(message, userID)
}

// broadcastToUserLocal is the internal implementation for user-specific local broadcast.
func (h *Hub) broadcastToUserLocal(message []byte, userID uint64) {
	h.mu.RLock()
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
	h.mu.RUnlock()
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
