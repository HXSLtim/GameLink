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
	MessageTypePing        MessageType = "ping"
	MessageTypePong        MessageType = "pong"
	MessageTypeSubscribe   MessageType = "subscribe"
	MessageTypeUnsubscribe MessageType = "unsubscribe"

	// User realtime message types
	MessageTypeChatMessage         MessageType = "chat_message"
	MessageTypeConversationMessage MessageType = "conversation_message"
	MessageTypeConversationClosed  MessageType = "conversation_closed"
	MessageTypeNotification        MessageType = "notification"
	MessageTypeOrderStatus         MessageType = "order_status"
	MessageTypeOrderNew            MessageType = "order_new"

	// Presence message types (Discord/Kook style)
	MessageTypePresenceUpdate    MessageType = "presence_update"    // 状态变更
	MessageTypePresenceSubscribe MessageType = "presence_subscribe" // 订阅状态
	MessageTypePresenceBatch     MessageType = "presence_batch"     // 批量状态更新

	// Room message types
	MessageTypeRoomCreated      MessageType = "room_created"       // 房间创建
	MessageTypeRoomUpdated      MessageType = "room_updated"       // 房间更新
	MessageTypeRoomClosed       MessageType = "room_closed"        // 房间关闭
	MessageTypeRoomMemberJoined MessageType = "room_member_joined" // 成员加入
	MessageTypeRoomMemberLeft   MessageType = "room_member_left"   // 成员离开
	MessageTypeRoomMemberReady  MessageType = "room_member_ready"  // 成员准备
	MessageTypeRoomStarted      MessageType = "room_started"       // 游戏开始
	MessageTypeRoomFinished     MessageType = "room_finished"      // 游戏结束

	// LFG (Looking For Group) message types
	MessageTypeLFGNew      MessageType = "lfg_new"      // 新匹配请求
	MessageTypeLFGMatched  MessageType = "lfg_matched"  // 匹配成功
	MessageTypeLFGExpired  MessageType = "lfg_expired"  // 请求过期
	MessageTypeLFGCanceled MessageType = "lfg_canceled" // 请求取消

	// Voice message types (TRTC)
	MessageTypeVoiceStarted      MessageType = "voice_started"       // 语音开启
	MessageTypeVoiceStopped      MessageType = "voice_stopped"       // 语音关闭
	MessageTypeVoiceMemberJoined MessageType = "voice_member_joined" // 成员加入语音
	MessageTypeVoiceMemberLeft   MessageType = "voice_member_left"   // 成员离开语音
	MessageTypeVoiceMemberMuted  MessageType = "voice_member_muted"  // 成员静音状态变更
)

// ChatMessageEvent represents a real-time chat message event.
type ChatMessageEvent struct {
	GroupID uint64      `json:"groupId"`
	Message interface{} `json:"message"`
}

// OrderStatusEvent represents order status changes.
type OrderStatusEvent struct {
	OrderID   uint64 `json:"orderId"`
	Status    string `json:"status"`
	Message   string `json:"message,omitempty"`
	UpdatedAt string `json:"updatedAt"`
}

// OrderNewEvent represents a new order available event.
type OrderNewEvent struct {
	OrderID        uint64     `json:"orderId"`
	Title          string     `json:"title"`
	PriceCents     int64      `json:"priceCents"`
	ScheduledStart *time.Time `json:"scheduledStart,omitempty"`
	GameID         *uint64    `json:"gameId,omitempty"`
}

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

// ============================================================================
// Presence Message Types (Discord/Kook style)
// ============================================================================

// PresenceUpdate represents a player presence status change.
type PresenceUpdate struct {
	PlayerID        uint64 `json:"playerId"`
	Status          string `json:"status"`
	CurrentGameID   uint64 `json:"currentGameId,omitempty"`
	CurrentGameName string `json:"currentGameName,omitempty"`
	CustomStatus    string `json:"customStatus,omitempty"`
	CurrentRoomID   uint64 `json:"currentRoomId,omitempty"`
	UpdatedAt       string `json:"updatedAt"`
}

// PresenceBatch represents multiple presence updates.
type PresenceBatch struct {
	Presences []PresenceUpdate `json:"presences"`
}

// NewPresenceUpdateMessage creates a presence update message.
func NewPresenceUpdateMessage(update *PresenceUpdate) ([]byte, error) {
	msg := NewWSMessage(MessageTypePresenceUpdate, update)
	return msg.ToJSON()
}

// NewPresenceBatchMessage creates a batch presence update message.
func NewPresenceBatchMessage(batch *PresenceBatch) ([]byte, error) {
	msg := NewWSMessage(MessageTypePresenceBatch, batch)
	return msg.ToJSON()
}

// ============================================================================
// Room Message Types
// ============================================================================

// RoomEvent represents a room-related event.
type RoomEvent struct {
	RoomID         uint64 `json:"roomId"`
	RoomName       string `json:"roomName"`
	RoomType       string `json:"roomType"`
	GameID         uint64 `json:"gameId"`
	GameName       string `json:"gameName,omitempty"`
	HostUserID     uint64 `json:"hostUserId"`
	Status         string `json:"status"`
	CurrentMembers int    `json:"currentMembers"`
	MaxMembers     int    `json:"maxMembers"`
}

// RoomMemberEvent represents a room member event.
type RoomMemberEvent struct {
	RoomID   uint64 `json:"roomId"`
	UserID   uint64 `json:"userId"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar,omitempty"`
	Role     string `json:"role"`
	IsReady  bool   `json:"isReady"`
}

// NewRoomCreatedMessage creates a room created message.
func NewRoomCreatedMessage(event *RoomEvent) ([]byte, error) {
	msg := NewWSMessage(MessageTypeRoomCreated, event)
	return msg.ToJSON()
}

// NewRoomUpdatedMessage creates a room updated message.
func NewRoomUpdatedMessage(event *RoomEvent) ([]byte, error) {
	msg := NewWSMessage(MessageTypeRoomUpdated, event)
	return msg.ToJSON()
}

// NewRoomClosedMessage creates a room closed message.
func NewRoomClosedMessage(roomID uint64) ([]byte, error) {
	msg := NewWSMessage(MessageTypeRoomClosed, map[string]uint64{"roomId": roomID})
	return msg.ToJSON()
}

// NewRoomMemberJoinedMessage creates a member joined message.
func NewRoomMemberJoinedMessage(event *RoomMemberEvent) ([]byte, error) {
	msg := NewWSMessage(MessageTypeRoomMemberJoined, event)
	return msg.ToJSON()
}

// NewRoomMemberLeftMessage creates a member left message.
func NewRoomMemberLeftMessage(event *RoomMemberEvent) ([]byte, error) {
	msg := NewWSMessage(MessageTypeRoomMemberLeft, event)
	return msg.ToJSON()
}

// NewRoomMemberReadyMessage creates a member ready message.
func NewRoomMemberReadyMessage(event *RoomMemberEvent) ([]byte, error) {
	msg := NewWSMessage(MessageTypeRoomMemberReady, event)
	return msg.ToJSON()
}

// NewRoomStartedMessage creates a room started message.
func NewRoomStartedMessage(roomID uint64) ([]byte, error) {
	msg := NewWSMessage(MessageTypeRoomStarted, map[string]uint64{"roomId": roomID})
	return msg.ToJSON()
}

// ============================================================================
// LFG Message Types
// ============================================================================

// LFGEvent represents an LFG request event.
type LFGEvent struct {
	RequestID       uint64 `json:"requestId"`
	UserID          uint64 `json:"userId"`
	UserNickname    string `json:"userNickname,omitempty"`
	GameID          uint64 `json:"gameId"`
	GameName        string `json:"gameName,omitempty"`
	RequestType     string `json:"requestType"`
	Title           string `json:"title"`
	Description     string `json:"description,omitempty"`
	RequiredPlayers int    `json:"requiredPlayers"`
	Status          string `json:"status"`
	MatchedRoomID   uint64 `json:"matchedRoomId,omitempty"`
}

// NewLFGNewMessage creates a new LFG request message.
func NewLFGNewMessage(event *LFGEvent) ([]byte, error) {
	msg := NewWSMessage(MessageTypeLFGNew, event)
	return msg.ToJSON()
}

// NewLFGMatchedMessage creates an LFG matched message.
func NewLFGMatchedMessage(event *LFGEvent) ([]byte, error) {
	msg := NewWSMessage(MessageTypeLFGMatched, event)
	return msg.ToJSON()
}

// NewLFGExpiredMessage creates an LFG expired message.
func NewLFGExpiredMessage(requestID uint64) ([]byte, error) {
	msg := NewWSMessage(MessageTypeLFGExpired, map[string]uint64{"requestId": requestID})
	return msg.ToJSON()
}

// NewLFGCanceledMessage creates an LFG canceled message.
func NewLFGCanceledMessage(requestID uint64) ([]byte, error) {
	msg := NewWSMessage(MessageTypeLFGCanceled, map[string]uint64{"requestId": requestID})
	return msg.ToJSON()
}

// ============================================================================
// Voice Message Types (TRTC)
// ============================================================================

// VoiceEvent represents a voice channel event.
type VoiceEvent struct {
	RoomID      uint64 `json:"roomId"`
	VoiceRoomID string `json:"voiceRoomId"`
	SDKAppID    uint64 `json:"sdkAppId,omitempty"`
}

// VoiceMemberEvent represents a voice member event.
type VoiceMemberEvent struct {
	RoomID   uint64 `json:"roomId"`
	UserID   uint64 `json:"userId"`
	Nickname string `json:"nickname"`
	InVoice  bool   `json:"inVoice"`
	IsMuted  bool   `json:"isMuted"`
}

// NewVoiceStartedMessage creates a voice started message.
func NewVoiceStartedMessage(event *VoiceEvent) ([]byte, error) {
	msg := NewWSMessage(MessageTypeVoiceStarted, event)
	return msg.ToJSON()
}

// NewVoiceStoppedMessage creates a voice stopped message.
func NewVoiceStoppedMessage(roomID uint64) ([]byte, error) {
	msg := NewWSMessage(MessageTypeVoiceStopped, map[string]uint64{"roomId": roomID})
	return msg.ToJSON()
}

// NewVoiceMemberJoinedMessage creates a voice member joined message.
func NewVoiceMemberJoinedMessage(event *VoiceMemberEvent) ([]byte, error) {
	msg := NewWSMessage(MessageTypeVoiceMemberJoined, event)
	return msg.ToJSON()
}

// NewVoiceMemberLeftMessage creates a voice member left message.
func NewVoiceMemberLeftMessage(event *VoiceMemberEvent) ([]byte, error) {
	msg := NewWSMessage(MessageTypeVoiceMemberLeft, event)
	return msg.ToJSON()
}

// NewVoiceMemberMutedMessage creates a voice member muted message.
func NewVoiceMemberMutedMessage(event *VoiceMemberEvent) ([]byte, error) {
	msg := NewWSMessage(MessageTypeVoiceMemberMuted, event)
	return msg.ToJSON()
}
