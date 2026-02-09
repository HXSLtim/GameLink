package ws

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewHub verifies hub initialization with default values.
func TestNewHub(t *testing.T) {
	hub := NewHub()

	require.NotNil(t, hub)
	assert.NotNil(t, hub.clients)
	assert.NotNil(t, hub.broadcast)
	assert.NotNil(t, hub.register)
	assert.NotNil(t, hub.unregister)
	assert.NotNil(t, hub.metrics)

	// Verify initial metrics
	assert.Equal(t, int64(0), hub.metrics.TotalConnections)
	assert.Equal(t, 0, hub.metrics.ActiveConnections)
	assert.False(t, hub.metrics.LastActivityAt.IsZero())
}

// TestHubRun verifies the hub's main event loop processes channels correctly.
func TestHubRun(t *testing.T) {
	hub := NewHub()

	// Start hub in background
	go hub.Run()

	// Create a test client
	client := &Client{
		UserID: 123,
		Role:   "admin",
		hub:    hub,
		send:   make(chan []byte, 256),
	}

	// Test registration
	hub.register <- client
	time.Sleep(10 * time.Millisecond)

	hub.mu.RLock()
	assert.Len(t, hub.clients, 1)
	assert.True(t, hub.clients[client])
	hub.mu.RUnlock()

	// Test broadcast
	hub.broadcast <- []byte(`{"type":"test","data":"message"}`)
	time.Sleep(10 * time.Millisecond)

	// Verify message sent
	select {
	case msg := <-client.send:
		assert.Equal(t, []byte(`{"type":"test","data":"message"}`), msg)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("expected broadcast message not received")
	}

	// Test unregistration
	hub.unregister <- client
	time.Sleep(10 * time.Millisecond)

	hub.mu.RLock()
	assert.Len(t, hub.clients, 0)
	hub.mu.RUnlock()
}

// TestHubClientRegistration verifies client registration and unregistration.
func TestHubClientRegistration(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Create test clients
	client1 := &Client{
		UserID: 1,
		Role:   "admin",
		hub:    hub,
		send:   make(chan []byte, 256),
	}
	client2 := &Client{
		UserID: 2,
		Role:   "user",
		hub:    hub,
		send:   make(chan []byte, 256),
	}

	// Register first client
	hub.register <- client1
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, 1, hub.GetOnlineCount())

	// Register second client
	hub.register <- client2
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, 2, hub.GetOnlineCount())

	// Unregister first client
	hub.unregister <- client1
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, 1, hub.GetOnlineCount())

	// Verify send channel closed after unregister
	select {
	case _, ok := <-client1.send:
		assert.False(t, ok, "send channel should be closed")
	case <-time.After(100 * time.Millisecond):
		// Channel closed, no message to read
	}
}

// TestHubDuplicateRegistration verifies handling of duplicate client registration.
func TestHubDuplicateRegistration(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client := &Client{
		UserID: 1,
		Role:   "admin",
		hub:    hub,
		send:   make(chan []byte, 256),
	}

	// Register client twice
	hub.register <- client
	time.Sleep(10 * time.Millisecond)
	hub.register <- client
	time.Sleep(10 * time.Millisecond)

	// Should still only have one entry
	hub.mu.RLock()
	count := len(hub.clients)
	hub.mu.RUnlock()

	assert.Equal(t, 1, count, "duplicate registration should not increase count")
}

// TestHubBroadcast verifies message broadcasting to all connected clients.
func TestHubBroadcast(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Create multiple clients
	var clients []*Client
	for i := 0; i < 5; i++ {
		client := &Client{
			UserID: uint64(i + 1),
			Role:   "user",
			hub:    hub,
			send:   make(chan []byte, 256),
		}
		clients = append(clients, client)
		hub.register <- client
	}

	time.Sleep(50 * time.Millisecond)

	// Broadcast message
	message := []byte(`{"type":"broadcast","data":"test"}`)
	hub.Broadcast(message)

	time.Sleep(50 * time.Millisecond)

	// Verify all clients received the message
	for _, client := range clients {
		select {
		case msg := <-client.send:
			assert.Equal(t, message, msg)
		case <-time.After(100 * time.Millisecond):
			t.Errorf("client %d did not receive broadcast message", client.UserID)
		}
	}
}

// TestHubBroadcastToRole verifies role-based message broadcasting.
func TestHubBroadcastToRole(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Create clients with different roles
	admin1 := &Client{
		UserID: 1,
		Role:   "admin",
		hub:    hub,
		send:   make(chan []byte, 256),
	}
	admin2 := &Client{
		UserID: 2,
		Role:   "admin",
		hub:    hub,
		send:   make(chan []byte, 256),
	}
	user1 := &Client{
		UserID: 3,
		Role:   "user",
		hub:    hub,
		send:   make(chan []byte, 256),
	}

	hub.register <- admin1
	hub.register <- admin2
	hub.register <- user1

	time.Sleep(50 * time.Millisecond)

	// Broadcast to admins only
	message := []byte(`{"type":"admin_alert","data":"system update"}`)
	hub.BroadcastToRole(message, "admin")

	time.Sleep(50 * time.Millisecond)

	// Verify admins received the message
	select {
	case msg := <-admin1.send:
		assert.Equal(t, message, msg)
	case <-time.After(100 * time.Millisecond):
		t.Error("admin1 did not receive role-based broadcast")
	}

	select {
	case msg := <-admin2.send:
		assert.Equal(t, message, msg)
	case <-time.After(100 * time.Millisecond):
		t.Error("admin2 did not receive role-based broadcast")
	}

	// Verify user did not receive the message
	select {
	case msg := <-user1.send:
		t.Errorf("user should not receive admin broadcast, got: %s", string(msg))
	case <-time.After(50 * time.Millisecond):
		// Expected timeout - user should not receive message
	}
}

// TestHubBroadcastToUser verifies direct user messaging.
func TestHubBroadcastToUser(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Create multiple clients
	client1 := &Client{
		UserID: 1,
		Role:   "user",
		hub:    hub,
		send:   make(chan []byte, 256),
	}
	client2 := &Client{
		UserID: 2,
		Role:   "user",
		hub:    hub,
		send:   make(chan []byte, 256),
	}

	hub.register <- client1
	hub.register <- client2

	time.Sleep(50 * time.Millisecond)

	// Send message to specific user
	message := []byte(`{"type":"direct_message","data":"hello"}`)
	hub.BroadcastToUser(message, 1)

	time.Sleep(50 * time.Millisecond)

	// Verify client1 received the message
	select {
	case msg := <-client1.send:
		assert.Equal(t, message, msg)
	case <-time.After(100 * time.Millisecond):
		t.Error("client1 did not receive direct message")
	}

	// Verify client2 did not receive the message
	select {
	case msg := <-client2.send:
		t.Errorf("client2 should not receive message to client1, got: %s", string(msg))
	case <-time.After(50 * time.Millisecond):
		// Expected timeout - client2 should not receive message
	}
}

// TestHubGetMetrics verifies metrics tracking and retrieval.
func TestHubGetMetrics(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Initial metrics
	metrics := hub.GetMetrics()
	assert.Equal(t, int64(0), metrics.TotalConnections)
	assert.Equal(t, 0, metrics.ActiveConnections)
	assert.Equal(t, int64(0), metrics.MessagesSent)
	assert.Equal(t, int64(0), metrics.MessagesReceived)

	// Register clients
	for i := 0; i < 3; i++ {
		client := &Client{
			UserID: uint64(i + 1),
			Role:   "user",
			hub:    hub,
			send:   make(chan []byte, 256),
		}
		hub.register <- client
	}

	time.Sleep(50 * time.Millisecond)

	metrics = hub.GetMetrics()
	assert.Equal(t, int64(3), metrics.TotalConnections)
	assert.Equal(t, 3, metrics.ActiveConnections)

	// Send broadcast to increment message count
	hub.Broadcast([]byte(`{"type":"test"}`))
	time.Sleep(50 * time.Millisecond)

	metrics = hub.GetMetrics()
	assert.Greater(t, metrics.MessagesSent, int64(0))
	assert.False(t, metrics.LastActivityAt.IsZero())
}

// TestHubGetOnlineCount verifies online client count tracking.
func TestHubGetOnlineCount(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	assert.Equal(t, 0, hub.GetOnlineCount())

	// Add clients
	for i := 0; i < 5; i++ {
		client := &Client{
			UserID: uint64(i + 1),
			Role:   "user",
			hub:    hub,
			send:   make(chan []byte, 256),
		}
		hub.register <- client
	}

	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, 5, hub.GetOnlineCount())

	// Remove clients
	client := &Client{
		UserID: 999,
		hub:    hub,
		send:   make(chan []byte, 256),
	}
	hub.unregister <- client
	time.Sleep(10 * time.Millisecond)

	// Non-existent client unregister should not affect count
	assert.Equal(t, 5, hub.GetOnlineCount())
}

// TestHubGetOnlineCountByRole verifies role-based online count tracking.
func TestHubGetOnlineCountByRole(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Create clients with different roles
	clients := []*Client{
		{UserID: 1, Role: "admin", hub: hub, send: make(chan []byte, 256)},
		{UserID: 2, Role: "admin", hub: hub, send: make(chan []byte, 256)},
		{UserID: 3, Role: "user", hub: hub, send: make(chan []byte, 256)},
		{UserID: 4, Role: "user", hub: hub, send: make(chan []byte, 256)},
		{UserID: 5, Role: "user", hub: hub, send: make(chan []byte, 256)},
		{UserID: 6, Role: "player", hub: hub, send: make(chan []byte, 256)},
	}

	for _, client := range clients {
		hub.register <- client
	}

	time.Sleep(50 * time.Millisecond)

	counts := hub.GetOnlineCountByRole()
	assert.Equal(t, 2, counts["admin"])
	assert.Equal(t, 3, counts["user"])
	assert.Equal(t, 1, counts["player"])
}

// TestHubClientCount verifies client count tracking.
func TestHubClientCount(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	assert.Equal(t, 0, hub.ClientCount())

	// Add clients
	for i := 0; i < 10; i++ {
		client := &Client{
			UserID: uint64(i + 1),
			Role:   "user",
			hub:    hub,
			send:   make(chan []byte, 256),
		}
		hub.register <- client
	}

	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, 10, hub.ClientCount())
}

// TestHubConcurrentRegistration verifies thread-safe concurrent client registration.
func TestHubConcurrentRegistration(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	numClients := 100
	var wg sync.WaitGroup

	// Register clients concurrently
	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			client := &Client{
				UserID: uint64(id + 1),
				Role:   "user",
				hub:    hub,
				send:   make(chan []byte, 256),
			}
			hub.register <- client
		}(i)
	}

	wg.Wait()
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, numClients, hub.ClientCount())
}

// TestHubConcurrentBroadcast verifies thread-safe concurrent broadcasting.
func TestHubConcurrentBroadcast(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Register clients
	numClients := 50
	for i := 0; i < numClients; i++ {
		client := &Client{
			UserID: uint64(i + 1),
			Role:   "user",
			hub:    hub,
			send:   make(chan []byte, 512), // Larger buffer for concurrent sends
		}
		hub.register <- client
	}

	time.Sleep(50 * time.Millisecond)

	// Broadcast concurrently
	var wg sync.WaitGroup
	numBroadcasts := 20
	for i := 0; i < numBroadcasts; i++ {
		wg.Add(1)
		go func(msgNum int) {
			defer wg.Done()
			message := []byte(`{"type":"concurrent","data":"` + string(rune(msgNum)) + `"}`)
			hub.Broadcast(message)
		}(i)
	}

	wg.Wait()
	time.Sleep(100 * time.Millisecond)

	// Verify messages sent
	metrics := hub.GetMetrics()
	assert.GreaterOrEqual(t, metrics.MessagesSent, int64(numBroadcasts))
}

// TestHubClientSendBufferFull verifies handling of full client send buffer.
func TestHubClientSendBufferFull(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Create client with small send buffer
	client := &Client{
		UserID: 1,
		Role:   "user",
		hub:    hub,
		send:   make(chan []byte, 1), // Very small buffer
	}

	hub.register <- client
	time.Sleep(50 * time.Millisecond)

	// Fill the buffer
	client.send <- []byte("message1")

	// Broadcast more messages (should trigger buffer full handling)
	for i := 0; i < 10; i++ {
		hub.Broadcast([]byte(`{"type":"test"}`))
	}

	time.Sleep(100 * time.Millisecond)

	// Client should be unregistered due to full buffer
	assert.Eventually(t, func() bool {
		return hub.ClientCount() == 0
	}, 500*time.Millisecond, 50*time.Millisecond, "client should be unregistered after buffer overflow")
}

// TestHubMetricsUpdate verifies metrics are updated correctly over time.
func TestHubMetricsUpdate(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	initialMetrics := hub.GetMetrics()
	initialTime := initialMetrics.LastActivityAt

	// Wait a bit and register a client
	time.Sleep(10 * time.Millisecond)
	client := &Client{
		UserID: 1,
		Role:   "user",
		hub:    hub,
		send:   make(chan []byte, 256),
	}
	hub.register <- client

	time.Sleep(10 * time.Millisecond)
	updatedMetrics := hub.GetMetrics()

	assert.Equal(t, int64(1), updatedMetrics.TotalConnections)
	assert.Equal(t, 1, updatedMetrics.ActiveConnections)
	assert.True(t, updatedMetrics.LastActivityAt.After(initialTime))
}

// TestHubRaceConditions verifies no race conditions with concurrent operations.
func TestHubRaceConditions(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	var wg sync.WaitGroup

	// Concurrent registrations
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			client := &Client{
				UserID: uint64(id + 1),
				Role:   "user",
				hub:    hub,
				send:   make(chan []byte, 256),
			}
			hub.register <- client
		}(i)
	}

	// Concurrent broadcasts
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			hub.Broadcast([]byte(`{"type":"race_test"}`))
		}(i)
	}

	// Concurrent role broadcasts
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			hub.BroadcastToRole([]byte(`{"type":"role_test"}`), "user")
		}(i)
	}

	// Concurrent metric queries
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_ = hub.GetMetrics()
			_ = hub.GetOnlineCount()
			_ = hub.GetOnlineCountByRole()
		}(i)
	}

	wg.Wait()
	time.Sleep(100 * time.Millisecond)

	// Should complete without race condition (detected by -race flag)
	assert.Greater(t, hub.ClientCount(), 0)
}

// TestHubUnregisterNonExistentClient verifies unregistering non-existent client is safe.
func TestHubUnregisterNonExistentClient(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client := &Client{
		UserID: 999,
		hub:    hub,
		send:   make(chan []byte, 256),
	}

	// Unregister without registering
	hub.unregister <- client
	time.Sleep(10 * time.Millisecond)

	// Should not panic or cause issues
	assert.Equal(t, 0, hub.ClientCount())
}
