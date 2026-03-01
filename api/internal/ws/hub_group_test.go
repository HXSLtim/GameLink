package ws

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHub_BroadcastToGroup(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	inGroup := &Client{
		UserID: 1,
		Role:   "user",
		hub:    hub,
		send:   make(chan []byte, 10),
	}
	outGroup := &Client{
		UserID: 2,
		Role:   "user",
		hub:    hub,
		send:   make(chan []byte, 10),
	}

	hub.register <- inGroup
	hub.register <- outGroup
	time.Sleep(20 * time.Millisecond)

	hub.SubscribeClientToGroup(inGroup, 99)
	hub.BroadcastToGroup([]byte(`{"type":"conversation_message"}`), 99)
	time.Sleep(20 * time.Millisecond)

	select {
	case payload := <-inGroup.send:
		assert.NotEmpty(t, payload)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("expected message for subscribed client")
	}

	select {
	case payload := <-outGroup.send:
		t.Fatalf("unexpected payload for unsubscribed client: %s", string(payload))
	case <-time.After(60 * time.Millisecond):
	}
}

func TestHub_GroupSubscriptionRemovedOnUnregister(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client := &Client{
		UserID: 3,
		Role:   "user",
		hub:    hub,
		send:   make(chan []byte, 10),
	}

	hub.register <- client
	time.Sleep(10 * time.Millisecond)
	hub.SubscribeClientToGroup(client, 100)
	assert.Equal(t, 1, hub.GetOnlineCountByGroup(100))

	hub.unregister <- client
	time.Sleep(20 * time.Millisecond)
	assert.Equal(t, 0, hub.GetOnlineCountByGroup(100))
}
