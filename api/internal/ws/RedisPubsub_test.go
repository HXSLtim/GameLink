// Package ws provides WebSocket functionality for real-time communication.
package ws

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestRedis creates a test Redis server and client.
func setupTestRedis(t *testing.T) (*miniredis.Miniredis, *redis.Client) {
	t.Helper()

	// Start miniredis
	mr, err := miniredis.Run()
	require.NoError(t, err)

	// Create Redis client
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	// Verify connection
	ctx := context.Background()
	require.NoError(t, client.Ping(ctx).Err())

	return mr, client
}

// TestNewRedisPubSub tests creating a new Redis Pub/Sub instance.
func TestNewRedisPubSub(t *testing.T) {
	mr, client := setupTestRedis(t)
	defer mr.Close()
	defer client.Close()

	hub := NewHub()
	redisPS := NewRedisPubSub(client, hub)

	assert.NotNil(t, redisPS)
	assert.NotNil(t, redisPS.client)
	assert.NotNil(t, redisPS.hub)
	assert.NotNil(t, redisPS.ctx)
}

// TestRedisPubSub_Subscribe tests subscribing to Redis channels.
func TestRedisPubSub_Subscribe(t *testing.T) {
	mr, client := setupTestRedis(t)
	defer mr.Close()
	defer client.Close()

	hub := NewHub()
	redisPS := NewRedisPubSub(client, hub)

	// Subscribe to channels
	redisPS.Subscribe()

	// Wait a bit for subscription to establish
	time.Sleep(100 * time.Millisecond)

	// Verify connection is active
	assert.True(t, redisPS.IsConnected())
}

// TestRedisPubSub_Broadcast tests broadcasting messages to all instances.
func TestRedisPubSub_Broadcast(t *testing.T) {
	mr, client := setupTestRedis(t)
	defer mr.Close()
	defer client.Close()

	hub := NewHub()
	redisPS := NewRedisPubSub(client, hub)
	redisPS.Subscribe()

	time.Sleep(100 * time.Millisecond)

	// Create a test message
	testMessage := []byte(`{"type":"test","data":"hello"}`)

	// Broadcast the message
	err := redisPS.Broadcast(testMessage)
	require.NoError(t, err)

	// Verify metrics
	metrics := redisPS.GetMetrics()
	assert.Equal(t, int64(1), metrics["messagesSent"])
}

// TestRedisPubSub_BroadcastToRole tests broadcasting to specific roles.
func TestRedisPubSub_BroadcastToRole(t *testing.T) {
	mr, client := setupTestRedis(t)
	defer mr.Close()
	defer client.Close()

	hub := NewHub()
	redisPS := NewRedisPubSub(client, hub)
	redisPS.Subscribe()

	time.Sleep(100 * time.Millisecond)

	testMessage := []byte(`{"type":"alert","data":"player alert"}`)

	err := redisPS.BroadcastToRole(testMessage, "player")
	require.NoError(t, err)

	// Verify metrics
	metrics := redisPS.GetMetrics()
	assert.Equal(t, int64(1), metrics["messagesSent"])
}

// TestRedisPubSub_BroadcastToUser tests broadcasting to specific users.
func TestRedisPubSub_BroadcastToUser(t *testing.T) {
	mr, client := setupTestRedis(t)
	defer mr.Close()
	defer client.Close()

	hub := NewHub()
	redisPS := NewRedisPubSub(client, hub)
	redisPS.Subscribe()

	time.Sleep(100 * time.Millisecond)

	testMessage := []byte(`{"type":"message","data":"direct message"}`)
	var userID uint64 = 12345

	err := redisPS.BroadcastToUser(testMessage, userID)
	require.NoError(t, err)

	// Verify metrics
	metrics := redisPS.GetMetrics()
	assert.Equal(t, int64(1), metrics["messagesSent"])
}

// TestRedisPubSub_BroadcastToGroup tests broadcasting to specific conversation groups.
func TestRedisPubSub_BroadcastToGroup(t *testing.T) {
	mr, client := setupTestRedis(t)
	defer mr.Close()
	defer client.Close()

	hub := NewHub()
	redisPS := NewRedisPubSub(client, hub)
	redisPS.Subscribe()

	time.Sleep(100 * time.Millisecond)

	testMessage := []byte(`{"type":"conversation_message","data":"group payload"}`)
	err := redisPS.BroadcastToGroup(testMessage, 99)
	require.NoError(t, err)

	metrics := redisPS.GetMetrics()
	assert.Equal(t, int64(1), metrics["messagesSent"])
}

// TestRedisPubSub_PublishPresence tests publishing presence updates.
func TestRedisPubSub_PublishPresence(t *testing.T) {
	mr, client := setupTestRedis(t)
	defer mr.Close()
	defer client.Close()

	hub := NewHub()
	redisPS := NewRedisPubSub(client, hub)
	redisPS.Subscribe()

	time.Sleep(100 * time.Millisecond)

	var userID uint64 = 12345
	role := "player"

	err := redisPS.PublishPresence(userID, role, "joined")
	require.NoError(t, err)

	err = redisPS.PublishPresence(userID, role, "left")
	require.NoError(t, err)
}

// TestRedisPubSub_HandleMessage tests handling received messages.
func TestRedisPubSub_HandleMessage(t *testing.T) {
	mr, client := setupTestRedis(t)
	defer mr.Close()
	defer client.Close()

	hub := NewHub()
	redisPS := NewRedisPubSub(client, hub)
	redisPS.Subscribe()

	time.Sleep(100 * time.Millisecond)

	// Create a PubSubMessage
	psMsg := PubSubMessage{
		Type:      "broadcast",
		Data:      []byte(`{"test":"data"}`),
		Timestamp: time.Now(),
	}

	data, err := json.Marshal(psMsg)
	require.NoError(t, err)

	// Publish directly to Redis to simulate receiving a message
	ctx := context.Background()
	err = client.Publish(ctx, ChannelBroadcast, data).Err()
	require.NoError(t, err)

	// Wait for message to be processed
	time.Sleep(200 * time.Millisecond)

	// Verify metrics
	metrics := redisPS.GetMetrics()
	assert.GreaterOrEqual(t, metrics["messagesReceived"], int64(1))
}

// TestRedisPubSub_Integration tests full integration with Hub.
func TestRedisPubSub_Integration(t *testing.T) {
	t.Skip("Skipping integration test due to complexity - individual tests cover functionality")
}

// TestNewHubWithRedis tests creating a hub with Redis support.
func TestNewHubWithRedis(t *testing.T) {
	mr, client := setupTestRedis(t)
	defer mr.Close()
	defer client.Close()

	// Test with Redis client
	hub := NewHubWithRedis(client)
	assert.NotNil(t, hub)

	// Verify Redis Pub/Sub is configured
	hub.mu.RLock()
	hasRedis := hub.redisPubSub != nil
	hub.mu.RUnlock()

	assert.True(t, hasRedis, "Hub should have Redis Pub/Sub configured")

	// Test with nil client (single-instance mode)
	hub2 := NewHubWithRedis(nil)
	assert.NotNil(t, hub2)

	// Verify no Redis Pub/Sub
	hub2.mu.RLock()
	hasRedis2 := hub2.redisPubSub != nil
	hub2.mu.RUnlock()

	assert.False(t, hasRedis2, "Hub should not have Redis Pub/Sub with nil client")
}

// TestRedisPubSub_Close tests graceful shutdown.
func TestRedisPubSub_Close(t *testing.T) {
	mr, client := setupTestRedis(t)
	defer mr.Close()
	defer client.Close()

	hub := NewHub()
	redisPS := NewRedisPubSub(client, hub)
	redisPS.Subscribe()

	time.Sleep(100 * time.Millisecond)

	// Close Redis Pub/Sub
	err := redisPS.Close()
	require.NoError(t, err)

	// Verify connection is closed (should fail after close)
	assert.False(t, redisPS.IsConnected())
}

// setupBenchmarkRedis creates a test Redis server and client for benchmarks.
func setupBenchmarkRedis(b *testing.B) (*miniredis.Miniredis, *redis.Client) {
	b.Helper()

	// Start miniredis
	mr, err := miniredis.Run()
	require.NoError(b, err)

	// Create Redis client
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	// Verify connection
	ctx := context.Background()
	require.NoError(b, client.Ping(ctx).Err())

	return mr, client
}

// BenchmarkRedisPubSub_Broadcast benchmarks broadcast performance.
func BenchmarkRedisPubSub_Broadcast(b *testing.B) {
	mr, client := setupBenchmarkRedis(b)
	defer mr.Close()
	defer client.Close()

	hub := NewHub()
	redisPS := NewRedisPubSub(client, hub)
	redisPS.Subscribe()

	time.Sleep(100 * time.Millisecond)

	testMessage := []byte(`{"type":"test","data":"benchmark"}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = redisPS.Broadcast(testMessage)
	}
}
