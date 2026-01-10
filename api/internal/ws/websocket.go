// Package ws provides WebSocket functionality for real-time communication.
package ws

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewHubWithRedis creates a new WebSocket Hub with Redis Pub/Sub support for horizontal scaling.
//
// This function sets up both the local Hub and the Redis Pub/Sub manager for cross-instance
// message broadcasting. It automatically starts the Redis subscription listener.
//
// Parameters:
//   - redisClient: Redis client for Pub/Sub operations (can be nil for single-instance mode)
//
// Returns:
//   - *Hub: Configured hub ready for use with NewHandler
//
// Usage:
//
//	// For single-instance mode (development, no Redis)
//	hub := ws.NewHubWithRedis(nil)
//	go hub.Run()
//
//	// For multi-instance mode (production, with Redis)
//	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
//	hub := ws.NewHubWithRedis(redisClient)
//	go hub.Run()
//
//	// Then create handler
//	handler := ws.NewHandler(hub)
func NewHubWithRedis(redisClient *redis.Client) *Hub {
	hub := NewHub()

	// If Redis client is provided, enable multi-instance support
	if redisClient != nil {
		// Verify Redis connection
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := redisClient.Ping(ctx).Err(); err != nil {
			log.Printf("WebSocket: Redis connection failed, falling back to single-instance mode: %v", err)
			// Return hub without Redis support (backward compatible)
			return hub
		}

		// Create and configure Redis Pub/Sub
		redisPS := NewRedisPubSub(redisClient, hub)
		hub.SetRedisPubSub(redisPS)

		// Start Redis subscription in background
		go redisPS.Subscribe()

		log.Println("WebSocket: Redis Pub/Sub enabled for multi-instance support")
	} else {
		log.Println("WebSocket: Running in single-instance mode (no Redis)")
	}

	return hub
}
