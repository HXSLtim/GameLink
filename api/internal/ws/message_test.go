package ws

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWSMessageMarshal verifies JSON marshaling of WebSocket messages.
func TestWSMessageMarshal(t *testing.T) {
	msg := WSMessage{
		Type:      "test",
		Timestamp: time.Now().Format(time.RFC3339),
		Data:      map[string]string{"key": "value"},
	}

	bytes, err := json.Marshal(msg)
	require.NoError(t, err)

	var decoded WSMessage
	err = json.Unmarshal(bytes, &decoded)
	require.NoError(t, err)

	assert.Equal(t, msg.Type, decoded.Type)
	assert.Equal(t, msg.Timestamp, decoded.Timestamp)
	assert.NotNil(t, decoded.Data)
}

// TestWSMessageMarshalWithoutData verifies marshaling messages without data.
func TestWSMessageMarshalWithoutData(t *testing.T) {
	msg := WSMessage{
		Type:      "ping",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	bytes, err := json.Marshal(msg)
	require.NoError(t, err)

	var decoded WSMessage
	err = json.Unmarshal(bytes, &decoded)
	require.NoError(t, err)

	assert.Equal(t, msg.Type, decoded.Type)
	assert.Equal(t, msg.Timestamp, decoded.Timestamp)
}

// TestNewWSMessage verifies creating a new WebSocket message.
func TestNewWSMessage(t *testing.T) {
	msgType := MessageTypeSystemStatus
	data := SystemStatus{
		CPUUsage:   75.5,
		Status:     "healthy",
		Uptime:     3600,
		Goroutines: 100,
	}

	msg := NewWSMessage(msgType, data)

	require.NotNil(t, msg)
	assert.Equal(t, string(msgType), msg.Type)
	assert.NotEmpty(t, msg.Timestamp)
	assert.NotNil(t, msg.Data)
}

// TestNewWSMessageWithNilData verifies creating message with nil data.
func TestNewWSMessageWithNilData(t *testing.T) {
	msg := NewWSMessage(MessageTypePing, nil)

	require.NotNil(t, msg)
	assert.Equal(t, string(MessageTypePing), msg.Type)
	assert.NotEmpty(t, msg.Timestamp)
}

// TestWSMessageToJSON verifies converting message to JSON bytes.
func TestWSMessageToJSON(t *testing.T) {
	msg := WSMessage{
		Type:      "test",
		Timestamp: time.Now().Format(time.RFC3339),
		Data:      map[string]string{"message": "hello"},
	}

	bytes, err := msg.ToJSON()
	require.NoError(t, err)

	var decoded map[string]interface{}
	err = json.Unmarshal(bytes, &decoded)
	require.NoError(t, err)

	assert.Equal(t, "test", decoded["type"])
	assert.NotNil(t, decoded["timestamp"])
	assert.NotNil(t, decoded["data"])
}

// TestNewSystemStatusMessage verifies creating system status message.
func TestNewSystemStatusMessage(t *testing.T) {
	status := &SystemStatus{
		CPUUsage:       45.2,
		MemoryUsage:    62.8,
		MemoryTotal:    8589934592, // 8GB
		MemoryUsed:     5390110000, // ~5GB
		Goroutines:     150,
		DBConnections:  DBConnections{Active: 5, Idle: 10, Max: 20},
		Uptime:         7200,
		RequestsPerSec: 125.5,
		Status:         "healthy",
	}

	bytes, err := NewSystemStatusMessage(status)
	require.NoError(t, err)

	var msg WSMessage
	err = json.Unmarshal(bytes, &msg)
	require.NoError(t, err)

	assert.Equal(t, string(MessageTypeSystemStatus), msg.Type)

	var decodedStatus SystemStatus
	statusBytes, _ := json.Marshal(msg.Data)
	err = json.Unmarshal(statusBytes, &decodedStatus)
	require.NoError(t, err)

	assert.Equal(t, 45.2, decodedStatus.CPUUsage)
	assert.Equal(t, "healthy", decodedStatus.Status)
	assert.Equal(t, 150, decodedStatus.Goroutines)
}

// TestNewOnlineUsersMessage verifies creating online users message.
func TestNewOnlineUsersMessage(t *testing.T) {
	users := &OnlineUsers{
		Total: 125,
		Peak:  200,
		ByRole: map[string]int{
			"admin":  5,
			"user":   100,
			"player": 20,
		},
		UpdatedAt: time.Now(),
	}

	bytes, err := NewOnlineUsersMessage(users)
	require.NoError(t, err)

	var msg WSMessage
	err = json.Unmarshal(bytes, &msg)
	require.NoError(t, err)

	assert.Equal(t, string(MessageTypeOnlineUsers), msg.Type)

	var decodedUsers OnlineUsers
	usersBytes, _ := json.Marshal(msg.Data)
	err = json.Unmarshal(usersBytes, &decodedUsers)
	require.NoError(t, err)

	assert.Equal(t, 125, decodedUsers.Total)
	assert.Equal(t, 200, decodedUsers.Peak)
	assert.Equal(t, 5, decodedUsers.ByRole["admin"])
	assert.Equal(t, 100, decodedUsers.ByRole["user"])
	assert.Equal(t, 20, decodedUsers.ByRole["player"])
}

// TestNewOrderQueueMessage verifies creating order queue message.
func TestNewOrderQueueMessage(t *testing.T) {
	queue := &OrderQueue{
		Pending:         15,
		Processing:      8,
		Completed:       150,
		ProcessingSpeed: 45.5,
		AverageWaitTime: 120.5,
		HasBacklog:      true,
	}

	bytes, err := NewOrderQueueMessage(queue)
	require.NoError(t, err)

	var msg WSMessage
	err = json.Unmarshal(bytes, &msg)
	require.NoError(t, err)

	assert.Equal(t, string(MessageTypeOrderQueue), msg.Type)

	var decodedQueue OrderQueue
	queueBytes, _ := json.Marshal(msg.Data)
	err = json.Unmarshal(queueBytes, &decodedQueue)
	require.NoError(t, err)

	assert.Equal(t, 15, decodedQueue.Pending)
	assert.Equal(t, 8, decodedQueue.Processing)
	assert.Equal(t, 150, decodedQueue.Completed)
	assert.Equal(t, 45.5, decodedQueue.ProcessingSpeed)
	assert.True(t, decodedQueue.HasBacklog)
}

// TestNewAlertMessage verifies creating alert message.
func TestNewAlertMessage(t *testing.T) {
	alert := &Alert{
		ID:        "alert-123",
		Level:     string(AlertLevelHigh),
		Type:      string(AlertTypeSystem),
		Title:     "High CPU Usage",
		Message:   "CPU usage exceeded 90%",
		Source:    "monitoring",
		CreatedAt: time.Now(),
		IsRead:    false,
	}

	bytes, err := NewAlertMessage(alert)
	require.NoError(t, err)

	var msg WSMessage
	err = json.Unmarshal(bytes, &msg)
	require.NoError(t, err)

	assert.Equal(t, string(MessageTypeAlert), msg.Type)

	var decodedAlert Alert
	alertBytes, _ := json.Marshal(msg.Data)
	err = json.Unmarshal(alertBytes, &decodedAlert)
	require.NoError(t, err)

	assert.Equal(t, "alert-123", decodedAlert.ID)
	assert.Equal(t, string(AlertLevelHigh), decodedAlert.Level)
	assert.Equal(t, string(AlertTypeSystem), decodedAlert.Type)
	assert.Equal(t, "High CPU Usage", decodedAlert.Title)
	assert.False(t, decodedAlert.IsRead)
}

// TestAlertLevels verifies all alert level constants.
func TestAlertLevels(t *testing.T) {
	assert.Equal(t, AlertLevel("high"), AlertLevelHigh)
	assert.Equal(t, AlertLevel("medium"), AlertLevelMedium)
	assert.Equal(t, AlertLevel("low"), AlertLevelLow)
}

// TestAlertTypes verifies all alert type constants.
func TestAlertTypes(t *testing.T) {
	assert.Equal(t, AlertType("system"), AlertTypeSystem)
	assert.Equal(t, AlertType("business"), AlertTypeBusiness)
	assert.Equal(t, AlertType("security"), AlertTypeSecurity)
}

// TestMessageTypes verifies all message type constants.
func TestMessageTypes(t *testing.T) {
	assert.Equal(t, MessageType("system_status"), MessageTypeSystemStatus)
	assert.Equal(t, MessageType("online_users"), MessageTypeOnlineUsers)
	assert.Equal(t, MessageType("order_queue"), MessageTypeOrderQueue)
	assert.Equal(t, MessageType("alert"), MessageTypeAlert)
	assert.Equal(t, MessageType("ping"), MessageTypePing)
	assert.Equal(t, MessageType("pong"), MessageTypePong)
	assert.Equal(t, MessageType("subscribe"), MessageTypeSubscribe)
}

// TestSystemStatusFields verifies system status structure.
func TestSystemStatusFields(t *testing.T) {
	status := SystemStatus{
		CPUUsage:       75.5,
		MemoryUsage:    62.8,
		MemoryTotal:    8589934592,
		MemoryUsed:     5390110000,
		Goroutines:     150,
		DBConnections:  DBConnections{Active: 5, Idle: 10, Max: 20},
		Uptime:         7200,
		RequestsPerSec: 125.5,
		Status:         "healthy",
	}

	// Verify JSON tags
	bytes, err := json.Marshal(status)
	require.NoError(t, err)

	var decoded map[string]interface{}
	err = json.Unmarshal(bytes, &decoded)
	require.NoError(t, err)

	assert.Contains(t, decoded, "cpuUsage")
	assert.Contains(t, decoded, "memoryUsage")
	assert.Contains(t, decoded, "memoryTotal")
	assert.Contains(t, decoded, "goroutines")
	assert.Contains(t, decoded, "uptime")
	assert.Contains(t, decoded, "status")
}

// TestOnlineUsersFields verifies online users structure.
func TestOnlineUsersFields(t *testing.T) {
	users := OnlineUsers{
		Total:     125,
		Peak:      200,
		ByRole:    map[string]int{"admin": 5, "user": 120},
		UpdatedAt: time.Now(),
	}

	bytes, err := json.Marshal(users)
	require.NoError(t, err)

	var decoded map[string]interface{}
	err = json.Unmarshal(bytes, &decoded)
	require.NoError(t, err)

	assert.Contains(t, decoded, "total")
	assert.Contains(t, decoded, "peak")
	assert.Contains(t, decoded, "byRole")
	assert.Contains(t, decoded, "updatedAt")
}

// TestOrderQueueFields verifies order queue structure.
func TestOrderQueueFields(t *testing.T) {
	queue := OrderQueue{
		Pending:         15,
		Processing:      8,
		Completed:       150,
		ProcessingSpeed: 45.5,
		AverageWaitTime: 120.5,
		HasBacklog:      true,
	}

	bytes, err := json.Marshal(queue)
	require.NoError(t, err)

	var decoded map[string]interface{}
	err = json.Unmarshal(bytes, &decoded)
	require.NoError(t, err)

	assert.Contains(t, decoded, "pending")
	assert.Contains(t, decoded, "processing")
	assert.Contains(t, decoded, "completed")
	assert.Contains(t, decoded, "processingSpeed")
	assert.Contains(t, decoded, "averageWaitTime")
	assert.Contains(t, decoded, "hasBacklog")
}

// TestAlertFields verifies alert structure.
func TestAlertFields(t *testing.T) {
	alert := Alert{
		ID:        "alert-123",
		Level:     "high",
		Type:      "system",
		Title:     "Test Alert",
		Message:   "Test message",
		Source:    "test",
		CreatedAt: time.Now(),
		IsRead:    false,
	}

	bytes, err := json.Marshal(alert)
	require.NoError(t, err)

	var decoded map[string]interface{}
	err = json.Unmarshal(bytes, &decoded)
	require.NoError(t, err)

	assert.Contains(t, decoded, "id")
	assert.Contains(t, decoded, "level")
	assert.Contains(t, decoded, "type")
	assert.Contains(t, decoded, "title")
	assert.Contains(t, decoded, "message")
	assert.Contains(t, decoded, "source")
	assert.Contains(t, decoded, "createdAt")
	assert.Contains(t, decoded, "isRead")
}

// TestWSMessageWithComplexData verifies handling complex nested data structures.
func TestWSMessageWithComplexData(t *testing.T) {
	type NestedData struct {
		Name string   `json:"name"`
		Tags []string `json:"tags"`
		Meta struct {
			Enabled bool `json:"enabled"`
		} `json:"meta"`
	}

	data := NestedData{
		Name: "test",
		Tags: []string{"tag1", "tag2", "tag3"},
		Meta: struct {
			Enabled bool `json:"enabled"`
		}{
			Enabled: true,
		},
	}

	msg := NewWSMessage("complex", data)

	bytes, err := msg.ToJSON()
	require.NoError(t, err)

	var decoded WSMessage
	err = json.Unmarshal(bytes, &decoded)
	require.NoError(t, err)

	assert.Equal(t, "complex", decoded.Type)
	assert.NotNil(t, decoded.Data)
}

// TestWSMessageTimestampFormat verifies timestamp format.
func TestWSMessageTimestampFormat(t *testing.T) {
	msg := NewWSMessage(MessageTypePing, nil)

	// Verify timestamp is in RFC3339 format
	_, err := time.Parse(time.RFC3339, msg.Timestamp)
	assert.NoError(t, err, "timestamp should be in RFC3339 format")
}

// TestWSMessageOmitEmptyData verifies data field omits empty when nil.
func TestWSMessageOmitEmptyData(t *testing.T) {
	msg := WSMessage{
		Type:      "test",
		Timestamp: time.Now().Format(time.RFC3339),
		// Data is nil
	}

	bytes, err := json.Marshal(msg)
	require.NoError(t, err)

	var decoded map[string]interface{}
	err = json.Unmarshal(bytes, &decoded)
	require.NoError(t, err)

	// When data is nil, it should be omitted or null
	// The omitempty tag means it won't be in the JSON
	_, hasData := decoded["data"]
	assert.False(t, hasData, "data field should be omitted when nil")
}

// TestDBConnectionsStructure verifies DB connections structure.
func TestDBConnectionsStructure(t *testing.T) {
	dbConns := DBConnections{
		Active: 10,
		Idle:   5,
		Max:    20,
	}

	bytes, err := json.Marshal(dbConns)
	require.NoError(t, err)

	var decoded DBConnections
	err = json.Unmarshal(bytes, &decoded)
	require.NoError(t, err)

	assert.Equal(t, 10, decoded.Active)
	assert.Equal(t, 5, decoded.Idle)
	assert.Equal(t, 20, decoded.Max)
}

// TestAlertReadStatus verifies alert read status tracking.
func TestAlertReadStatus(t *testing.T) {
	alert := &Alert{
		ID:        "alert-123",
		Level:     string(AlertLevelMedium),
		Type:      string(AlertTypeBusiness),
		Title:     "New Order",
		Message:   "A new order has been placed",
		Source:    "order_service",
		CreatedAt: time.Now(),
		IsRead:    false,
	}

	// Initial state
	assert.False(t, alert.IsRead)

	// Mark as read
	alert.IsRead = true
	assert.True(t, alert.IsRead)

	// Verify it serializes correctly
	bytes, err := json.Marshal(alert)
	require.NoError(t, err)

	var decoded Alert
	err = json.Unmarshal(bytes, &decoded)
	require.NoError(t, err)
	assert.True(t, decoded.IsRead)
}

// TestMessageFactoryFunctions verifies all message factory functions.
func TestMessageFactoryFunctions(t *testing.T) {
	t.Run("SystemStatusMessage", func(t *testing.T) {
		status := &SystemStatus{Status: "healthy"}
		bytes, err := NewSystemStatusMessage(status)
		require.NoError(t, err)
		assert.NotEmpty(t, bytes)
	})

	t.Run("OnlineUsersMessage", func(t *testing.T) {
		users := &OnlineUsers{Total: 100}
		bytes, err := NewOnlineUsersMessage(users)
		require.NoError(t, err)
		assert.NotEmpty(t, bytes)
	})

	t.Run("OrderQueueMessage", func(t *testing.T) {
		queue := &OrderQueue{Pending: 5}
		bytes, err := NewOrderQueueMessage(queue)
		require.NoError(t, err)
		assert.NotEmpty(t, bytes)
	})

	t.Run("AlertMessage", func(t *testing.T) {
		alert := &Alert{ID: "123", Title: "Test"}
		bytes, err := NewAlertMessage(alert)
		require.NoError(t, err)
		assert.NotEmpty(t, bytes)
	})
}

// TestConcurrentMessageCreation verifies thread-safe message creation.
func TestConcurrentMessageCreation(t *testing.T) {
	status := &SystemStatus{Status: "healthy"}
	users := &OnlineUsers{Total: 100}
	queue := &OrderQueue{Pending: 5}
	alert := &Alert{ID: "123", Title: "Test"}

	done := make(chan bool, 4)

	// Create messages concurrently
	go func() {
		_, _ = NewSystemStatusMessage(status)
		done <- true
	}()

	go func() {
		_, _ = NewOnlineUsersMessage(users)
		done <- true
	}()

	go func() {
		_, _ = NewOrderQueueMessage(queue)
		done <- true
	}()

	go func() {
		_, _ = NewAlertMessage(alert)
		done <- true
	}()

	for i := 0; i < 4; i++ {
		<-done
	}

	// All should complete without race condition (detected by -race flag)
}
