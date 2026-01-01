package ws

import (
	"encoding/json"
	"time"
)

// WSMessage represents a WebSocket message structure.
type WSMessage struct {
	Type      string      `json:"type"`
	Timestamp string      `json:"timestamp"`
	Data      interface{} `json:"data,omitempty"`
}

// MessageType defines the type of WebSocket message.
type MessageType string

const (
	// System status message types
	MessageTypeSystemStatus MessageType = "system_status"
	MessageTypeOnlineUsers  MessageType = "online_users"
	MessageTypeOrderQueue   MessageType = "order_queue"
	MessageTypeAlert        MessageType = "alert"

	// Client message types
	MessageTypePing      MessageType = "ping"
	MessageTypePong      MessageType = "pong"
	MessageTypeSubscribe MessageType = "subscribe"
)

// SystemStatus represents system health metrics.
type SystemStatus struct {
	CPUUsage       float64       `json:"cpuUsage"`
	MemoryUsage    float64       `json:"memoryUsage"`
	MemoryTotal    uint64        `json:"memoryTotal"`
	MemoryUsed     uint64        `json:"memoryUsed"`
	Goroutines     int           `json:"goroutines"`
	DBConnections  DBConnections `json:"dbConnections"`
	Uptime         int64         `json:"uptime"` // seconds
	RequestsPerSec float64       `json:"requestsPerSec"`
	Status         string        `json:"status"` // healthy, degraded, critical
}

// DBConnections represents database connection pool status.
type DBConnections struct {
	Active int `json:"active"`
	Idle   int `json:"idle"`
	Max    int `json:"max"`
}

// OnlineUsers represents online user statistics.
type OnlineUsers struct {
	Total     int            `json:"total"`
	Peak      int            `json:"peak"`
	ByRole    map[string]int `json:"byRole"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

// OrderQueue represents order processing queue status.
type OrderQueue struct {
	Pending         int     `json:"pending"`
	Processing      int     `json:"processing"`
	Completed       int     `json:"completed"`
	ProcessingSpeed float64 `json:"processingSpeed"` // orders per minute
	AverageWaitTime float64 `json:"averageWaitTime"` // seconds
	HasBacklog      bool    `json:"hasBacklog"`
}

// Alert represents a system or business alert.
type Alert struct {
	ID        string    `json:"id"`
	Level     string    `json:"level"` // high, medium, low
	Type      string    `json:"type"`  // system, business, security
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"createdAt"`
	IsRead    bool      `json:"isRead"`
}

// AlertLevel defines alert severity levels.
type AlertLevel string

const (
	AlertLevelHigh   AlertLevel = "high"
	AlertLevelMedium AlertLevel = "medium"
	AlertLevelLow    AlertLevel = "low"
)

// AlertType defines alert categories.
type AlertType string

const (
	AlertTypeSystem   AlertType = "system"
	AlertTypeBusiness AlertType = "business"
	AlertTypeSecurity AlertType = "security"
)

// NewWSMessage creates a new WebSocket message.
func NewWSMessage(msgType MessageType, data interface{}) *WSMessage {
	return &WSMessage{
		Type:      string(msgType),
		Timestamp: time.Now().Format(time.RFC3339),
		Data:      data,
	}
}

// ToJSON converts the message to JSON bytes.
func (m *WSMessage) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// NewSystemStatusMessage creates a system status message.
func NewSystemStatusMessage(status *SystemStatus) ([]byte, error) {
	msg := NewWSMessage(MessageTypeSystemStatus, status)
	return msg.ToJSON()
}

// NewOnlineUsersMessage creates an online users message.
func NewOnlineUsersMessage(users *OnlineUsers) ([]byte, error) {
	msg := NewWSMessage(MessageTypeOnlineUsers, users)
	return msg.ToJSON()
}

// NewOrderQueueMessage creates an order queue message.
func NewOrderQueueMessage(queue *OrderQueue) ([]byte, error) {
	msg := NewWSMessage(MessageTypeOrderQueue, queue)
	return msg.ToJSON()
}

// NewAlertMessage creates an alert message.
func NewAlertMessage(alert *Alert) ([]byte, error) {
	msg := NewWSMessage(MessageTypeAlert, alert)
	return msg.ToJSON()
}
