package model

import "time"

// PlayerPresenceStatus 陪玩师丰富在线状态
// @Enum online, accepting, in_game, matching, resting, offline, invisible
type PlayerPresenceStatus string

const (
	PresenceOnline    PlayerPresenceStatus = "online"    // 在线空闲
	PresenceAccepting PlayerPresenceStatus = "accepting" // 接单中
	PresenceInGame    PlayerPresenceStatus = "in_game"   // 游戏中
	PresenceMatching  PlayerPresenceStatus = "matching"  // 匹配中
	PresenceResting   PlayerPresenceStatus = "resting"   // 休息中
	PresenceOffline   PlayerPresenceStatus = "offline"   // 离线
	PresenceInvisible PlayerPresenceStatus = "invisible" // 隐身
)

// PlayerPresence 陪玩师实时状态（丰富状态系统）
// 用于实现 Discord/Kook 风格的在线状态显示
type PlayerPresence struct {
	Base
	PlayerID        uint64               `json:"playerId" gorm:"column:player_id;uniqueIndex;not null"`
	Status          PlayerPresenceStatus `json:"status" gorm:"column:status;size:32;index;default:'offline'"`
	CurrentGameID   *uint64              `json:"currentGameId,omitempty" gorm:"column:current_game_id;index"`
	CurrentGameName string               `json:"currentGameName,omitempty" gorm:"column:current_game_name;size:64"`
	CustomStatus    string               `json:"customStatus,omitempty" gorm:"column:custom_status;size:128"` // 自定义状态文字
	CurrentOrderID  *uint64              `json:"currentOrderId,omitempty" gorm:"column:current_order_id;index"`
	CurrentRoomID   *uint64              `json:"currentRoomId,omitempty" gorm:"column:current_room_id;index"`
	LastHeartbeatAt time.Time            `json:"lastHeartbeatAt" gorm:"column:last_heartbeat_at;index"`
	DeviceType      string               `json:"deviceType,omitempty" gorm:"column:device_type;size:32"` // web/ios/android/mini

	// Relations
	Player      *Player `json:"player,omitempty" gorm:"foreignKey:PlayerID"`
	CurrentGame *Game   `json:"currentGame,omitempty" gorm:"foreignKey:CurrentGameID"`
}

// TableName specifies the table name for PlayerPresence
func (PlayerPresence) TableName() string {
	return "player_presences"
}

// IsOnline returns true if the player is in any online state
func (p *PlayerPresence) IsOnline() bool {
	return p.Status != PresenceOffline && p.Status != PresenceInvisible
}

// IsAvailable returns true if the player can accept orders
func (p *PlayerPresence) IsAvailable() bool {
	return p.Status == PresenceOnline || p.Status == PresenceAccepting
}

// IsVisible returns true if the player should be shown to others
func (p *PlayerPresence) IsVisible() bool {
	return p.Status != PresenceInvisible
}

// GetStatusDisplay returns a human-readable status string
func (p *PlayerPresence) GetStatusDisplay() string {
	switch p.Status {
	case PresenceOnline:
		return "在线"
	case PresenceAccepting:
		return "接单中"
	case PresenceInGame:
		if p.CurrentGameName != "" {
			return "正在玩 " + p.CurrentGameName
		}
		return "游戏中"
	case PresenceMatching:
		return "匹配中"
	case PresenceResting:
		return "休息中"
	case PresenceOffline:
		return "离线"
	case PresenceInvisible:
		return "隐身"
	default:
		return "未知"
	}
}

// GetStatusColor returns a color code for the status (for UI)
func (p *PlayerPresence) GetStatusColor() string {
	switch p.Status {
	case PresenceOnline:
		return "#22c55e" // green
	case PresenceAccepting:
		return "#3b82f6" // blue
	case PresenceInGame:
		return "#a855f7" // purple
	case PresenceMatching:
		return "#f59e0b" // amber
	case PresenceResting:
		return "#6b7280" // gray
	case PresenceOffline:
		return "#9ca3af" // light gray
	case PresenceInvisible:
		return "#9ca3af" // light gray
	default:
		return "#9ca3af"
	}
}
