package model

import "time"

// LFGRequestStatus 快速匹配请求状态
// @Enum pending, matched, expired, canceled
type LFGRequestStatus string

const (
	LFGPending  LFGRequestStatus = "pending"  // 等待匹配
	LFGMatched  LFGRequestStatus = "matched"  // 已匹配
	LFGExpired  LFGRequestStatus = "expired"  // 已过期
	LFGCanceled LFGRequestStatus = "canceled" // 已取消
)

// LFGRequestType 匹配请求类型
// @Enum find_player, find_team
type LFGRequestType string

const (
	LFGFindPlayer LFGRequestType = "find_player" // 找陪玩
	LFGFindTeam   LFGRequestType = "find_team"   // 找队伍
)

// LFGRequest 快速匹配请求 (Looking For Group)
type LFGRequest struct {
	Base
	UserID          uint64           `json:"userId" gorm:"column:user_id;not null;index"`
	GameID          uint64           `json:"gameId" gorm:"column:game_id;not null;index"`
	RequestType     LFGRequestType   `json:"requestType" gorm:"column:request_type;size:32;not null;index"`
	Title           string           `json:"title" gorm:"column:title;size:64"`
	Description     string           `json:"description,omitempty" gorm:"column:description;size:256"`
	RequiredPlayers int              `json:"requiredPlayers" gorm:"column:required_players;default:1"`
	MinRank         string           `json:"minRank,omitempty" gorm:"column:min_rank;size:32"`
	MaxPriceCents   int64            `json:"maxPriceCents,omitempty" gorm:"column:max_price_cents"`
	Status          LFGRequestStatus `json:"status" gorm:"column:status;size:32;default:'pending';index"`
	ExpiresAt       time.Time        `json:"expiresAt" gorm:"column:expires_at;index"`
	MatchedRoomID   *uint64          `json:"matchedRoomId,omitempty" gorm:"column:matched_room_id;index"` // 关联 ChatGroup ID
	MatchedAt       *time.Time       `json:"matchedAt,omitempty" gorm:"column:matched_at"`

	// Relations
	User        *User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Game        *Game      `json:"game,omitempty" gorm:"foreignKey:GameID"`
	MatchedRoom *ChatGroup `json:"matchedRoom,omitempty" gorm:"foreignKey:MatchedRoomID"` // 匹配成功后创建的聊天房间
}

// TableName specifies the table name for LFGRequest
func (LFGRequest) TableName() string {
	return "lfg_requests"
}

// IsActive returns true if the request is still active
func (r *LFGRequest) IsActive() bool {
	return r.Status == LFGPending && time.Now().Before(r.ExpiresAt)
}

// IsExpired returns true if the request has expired
func (r *LFGRequest) IsExpired() bool {
	return time.Now().After(r.ExpiresAt)
}

// CanMatch returns true if this request can be matched
func (r *LFGRequest) CanMatch() bool {
	return r.Status == LFGPending && !r.IsExpired()
}

// GetTimeRemaining returns the time remaining before expiration
func (r *LFGRequest) GetTimeRemaining() time.Duration {
	if r.IsExpired() {
		return 0
	}
	return time.Until(r.ExpiresAt)
}
