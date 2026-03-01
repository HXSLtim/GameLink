// Package ws provides WebSocket functionality for real-time communication.
package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// Redis channels for WebSocket message distribution
	ChannelBroadcast = "ws:broadcast" // Broadcast to all connected clients
	ChannelRole      = "ws:role:%s"   // Broadcast to specific role (user, player, admin)
	ChannelUser      = "ws:user:%d"   // Send to specific user
	ChannelGroup     = "ws:group:%d"  // Send to specific conversation/group subscribers
	ChannelPresence  = "ws:presence"  // User presence updates
)

// RedisPubSub manages Redis Pub/Sub for cross-instance WebSocket communication.
// This enables horizontal scaling by allowing multiple WebSocket server instances
// to broadcast messages to each other's connected clients.
type RedisPubSub struct {
	client *redis.Client
	hub    *Hub // Local hub for local client broadcasting
	ctx    context.Context
	cancel context.CancelFunc

	// Subscription management
	mu       sync.RWMutex
	pubsub   *redis.PubSub
	channels map[string]struct{}

	// Metrics
	messagesReceived int64
	messagesSent     int64
	lastActivityAt   time.Time
}

// PubSubMessage represents a message transmitted via Redis Pub/Sub.
type PubSubMessage struct {
	Type      string    `json:"type"`   // broadcast, role, user
	UserID    *uint64   `json:"userID"` // Target user ID (for user messages)
	Role      *string   `json:"role"`   // Target role (for role messages)
	GroupID   *uint64   `json:"groupID,omitempty"`
	Data      []byte    `json:"data"` // Actual message payload
	Timestamp time.Time `json:"timestamp"`
}

// NewRedisPubSub creates a new Redis Pub/Sub manager for WebSocket multi-instance support.
//
// Parameters:
//   - client: Redis client for Pub/Sub operations
//   - hub: Local hub instance for broadcasting to local clients
//
// Returns:
//   - *RedisPubSub: Configured Redis Pub/Sub manager
//
// Example:
//
//	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
//	hub := ws.NewHub()
//	redisPS := ws.NewRedisPubSub(redisClient, hub)
//	go redisPS.Subscribe()
func NewRedisPubSub(client *redis.Client, hub *Hub) *RedisPubSub {
	ctx, cancel := context.WithCancel(context.Background())

	return &RedisPubSub{
		client:         client,
		hub:            hub,
		ctx:            ctx,
		cancel:         cancel,
		channels:       make(map[string]struct{}),
		lastActivityAt: time.Now(),
	}
}

// Subscribe subscribes to Redis channels and starts listening for messages.
// This method starts a goroutine that receives messages from Redis and
// broadcasts them to local clients.
//
// Subscribed channels:
//   - ws:broadcast - Broadcast to all connected clients
//   - ws:role:* - Broadcast to clients with specific roles
//   - ws:user:* - Send to specific users
//   - ws:presence - User presence updates
//
// This method should be called in a goroutine:
//
//	go redisPS.Subscribe()
func (rps *RedisPubSub) Subscribe() {
	// Subscribe to all relevant channels
	patterns := []string{
		ChannelBroadcast,
		"ws:role:*", // Pattern for role channels
		"ws:user:*", // Pattern for user channels
		"ws:group:*",
		ChannelPresence,
	}

	rps.mu.Lock()
	rps.pubsub = rps.client.PSubscribe(rps.ctx, patterns...)
	for _, pattern := range patterns {
		rps.channels[pattern] = struct{}{}
	}
	rps.mu.Unlock()

	log.Printf("Redis Pub/Sub: Subscribed to channels: %v", patterns)

	// Start listening for messages
	go rps.listen()
}

// listen receives messages from Redis and forwards them to local clients.
func (rps *RedisPubSub) listen() {
	ch := rps.pubsub.Channel()

	for {
		select {
		case <-rps.ctx.Done():
			log.Println("Redis Pub/Sub: Stopping listener")
			return

		case msg, ok := <-ch:
			if !ok {
				log.Println("Redis Pub/Sub: Channel closed")
				return
			}

			rps.mu.Lock()
			rps.messagesReceived++
			rps.lastActivityAt = time.Now()
			rps.mu.Unlock()

			// Parse and handle the message
			if err := rps.handleMessage(msg); err != nil {
				log.Printf("Redis Pub/Sub: Error handling message: %v", err)
			}
		}
	}
}

// handleMessage processes a received Redis message and broadcasts to local clients.
func (rps *RedisPubSub) handleMessage(msg *redis.Message) error {
	// Parse the PubSubMessage
	var psMsg PubSubMessage
	if err := json.Unmarshal([]byte(msg.Payload), &psMsg); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	// Route based on message type
	// IMPORTANT: Use local-only broadcast methods to avoid infinite loop
	// (Redis message -> hub.Broadcast -> redisPS.Broadcast -> Redis message...)
	switch psMsg.Type {
	case "broadcast":
		// Broadcast to local clients only (don't re-publish to Redis)
		rps.hub.BroadcastLocal(psMsg.Data)

	case "role":
		// Broadcast to local clients with specific role
		if psMsg.Role != nil {
			rps.hub.BroadcastToRoleLocal(psMsg.Data, *psMsg.Role)
		}

	case "user":
		// Send to local client with specific user ID
		if psMsg.UserID != nil {
			rps.hub.BroadcastToUserLocal(psMsg.Data, *psMsg.UserID)
		}

	case "group":
		// Send to local subscribers of group conversation
		if psMsg.GroupID != nil {
			rps.hub.BroadcastToGroupLocal(psMsg.Data, *psMsg.GroupID)
		}

	case "presence":
		// Handle presence updates (user joined/left)
		log.Printf("Redis Pub/Sub: Presence update: %s", string(psMsg.Data))

	default:
		log.Printf("Redis Pub/Sub: Unknown message type: %s", psMsg.Type)
	}

	return nil
}

// Broadcast sends a message to all connected clients across all instances.
// This publishes the message to Redis, which then forwards it to all instances.
func (rps *RedisPubSub) Broadcast(message []byte) error {
	psMsg := PubSubMessage{
		Type:      "broadcast",
		Data:      message,
		Timestamp: time.Now(),
	}

	data, err := json.Marshal(psMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %w", err)
	}

	if err := rps.client.Publish(rps.ctx, ChannelBroadcast, data).Err(); err != nil {
		return fmt.Errorf("failed to publish broadcast: %w", err)
	}

	rps.mu.Lock()
	rps.messagesSent++
	rps.lastActivityAt = time.Now()
	rps.mu.Unlock()

	return nil
}

// BroadcastToRole sends a message to all clients with a specific role across all instances.
func (rps *RedisPubSub) BroadcastToRole(message []byte, role string) error {
	psMsg := PubSubMessage{
		Type:      "role",
		Role:      &role,
		Data:      message,
		Timestamp: time.Now(),
	}

	data, err := json.Marshal(psMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal role message: %w", err)
	}

	channel := fmt.Sprintf(ChannelRole, role)
	if err := rps.client.Publish(rps.ctx, channel, data).Err(); err != nil {
		return fmt.Errorf("failed to publish role message: %w", err)
	}

	rps.mu.Lock()
	rps.messagesSent++
	rps.lastActivityAt = time.Now()
	rps.mu.Unlock()

	return nil
}

// BroadcastToUser sends a message to a specific user across all instances.
// If the user is connected to multiple instances, they will receive the message
// from the instance they are currently connected to.
func (rps *RedisPubSub) BroadcastToUser(message []byte, userID uint64) error {
	psMsg := PubSubMessage{
		Type:      "user",
		UserID:    &userID,
		Data:      message,
		Timestamp: time.Now(),
	}

	data, err := json.Marshal(psMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal user message: %w", err)
	}

	channel := fmt.Sprintf(ChannelUser, userID)
	if err := rps.client.Publish(rps.ctx, channel, data).Err(); err != nil {
		return fmt.Errorf("failed to publish user message: %w", err)
	}

	rps.mu.Lock()
	rps.messagesSent++
	rps.lastActivityAt = time.Now()
	rps.mu.Unlock()

	return nil
}

// BroadcastToGroup sends a message to subscribers of a group across all instances.
func (rps *RedisPubSub) BroadcastToGroup(message []byte, groupID uint64) error {
	psMsg := PubSubMessage{
		Type:      "group",
		GroupID:   &groupID,
		Data:      message,
		Timestamp: time.Now(),
	}

	data, err := json.Marshal(psMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal group message: %w", err)
	}

	channel := fmt.Sprintf(ChannelGroup, groupID)
	if err := rps.client.Publish(rps.ctx, channel, data).Err(); err != nil {
		return fmt.Errorf("failed to publish group message: %w", err)
	}

	rps.mu.Lock()
	rps.messagesSent++
	rps.lastActivityAt = time.Now()
	rps.mu.Unlock()

	return nil
}

// PublishPresence publishes a user presence update (joined/left) to all instances.
func (rps *RedisPubSub) PublishPresence(userID uint64, role string, action string) error {
	presence := map[string]interface{}{
		"userID":    userID,
		"role":      role,
		"action":    action, // "joined" or "left"
		"timestamp": time.Now().Format(time.RFC3339),
	}

	data, err := json.Marshal(presence)
	if err != nil {
		return fmt.Errorf("failed to marshal presence message: %w", err)
	}

	psMsg := PubSubMessage{
		Type:      "presence",
		Data:      data,
		Timestamp: time.Now(),
	}

	psData, err := json.Marshal(psMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal presence wrapper: %w", err)
	}

	if err := rps.client.Publish(rps.ctx, ChannelPresence, psData).Err(); err != nil {
		return fmt.Errorf("failed to publish presence: %w", err)
	}

	return nil
}

// Close gracefully shuts down the Redis Pub/Sub connection.
func (rps *RedisPubSub) Close() error {
	rps.cancel()

	rps.mu.Lock()
	defer rps.mu.Unlock()

	if rps.pubsub != nil {
		if err := rps.pubsub.Close(); err != nil {
			return fmt.Errorf("failed to close pubsub: %w", err)
		}
	}

	return nil
}

// GetMetrics returns current Redis Pub/Sub metrics.
func (rps *RedisPubSub) GetMetrics() map[string]interface{} {
	rps.mu.RLock()
	defer rps.mu.RUnlock()

	return map[string]interface{}{
		"messagesReceived": rps.messagesReceived,
		"messagesSent":     rps.messagesSent,
		"lastActivityAt":   rps.lastActivityAt,
		"subscribedChannels": func() []string {
			channels := make([]string, 0, len(rps.channels))
			for ch := range rps.channels {
				channels = append(channels, ch)
			}
			return channels
		}(),
	}
}

// IsConnected checks if Redis Pub/Sub is active.
func (rps *RedisPubSub) IsConnected() bool {
	rps.mu.RLock()
	defer rps.mu.RUnlock()

	if rps.pubsub == nil {
		return false
	}

	// Ping Redis to verify connection
	ctx, cancel := context.WithTimeout(rps.ctx, 1*time.Second)
	defer cancel()

	return rps.client.Ping(ctx).Err() == nil
}
