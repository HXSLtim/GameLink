package model

import "time"

// FollowStatus 关注状态
type FollowStatus string

const (
	// FollowStatusActive 关注中
	FollowStatusActive FollowStatus = "active"
	// FollowStatusBlocked 已拉黑
	FollowStatusBlocked FollowStatus = "blocked"
)

// Follow 关注关系
type Follow struct {
	ID         uint64       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     uint64       `gorm:"not null;index:idx_user_player" json:"userId"`
	PlayerID   uint64       `gorm:"not null;index:idx_user_player" json:"playerId"`
	Status     FollowStatus `gorm:"type:varchar(32);not null;default:'active'" json:"status"`
	NotifyNewService bool   `gorm:"default:true" json:"notifyNewService"` // 新服务通知
	NotifyOnline     bool   `gorm:"default:true" json:"notifyOnline"`     // 上线通知
	CreatedAt  time.Time    `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time    `gorm:"autoUpdateTime" json:"updatedAt"`
}

// TableName 指定表名
func (Follow) TableName() string {
	return "follows"
}

// Friendship 好友关系
type Friendship struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID1    uint64    `gorm:"not null;index:idx_users" json:"userId1"` // 较小的用户ID
	UserID2    uint64    `gorm:"not null;index:idx_users" json:"userId2"` // 较大的用户ID
	Status     string    `gorm:"type:varchar(32);not null;default:'pending'" json:"status"` // pending/accepted/rejected
	InitiatorID uint64   `gorm:"not null" json:"initiatorId"` // 发起人ID
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	AcceptedAt *time.Time `json:"acceptedAt"`
}

// TableName 指定表名
func (Friendship) TableName() string {
	return "friendships"
}

// Message 私信
type Message struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	SenderID   uint64    `gorm:"not null;index" json:"senderId"`
	ReceiverID uint64    `gorm:"not null;index" json:"receiverId"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	IsRead     bool      `gorm:"default:false" json:"isRead"`
	ReadAt     *time.Time `json:"readAt"`
	CreatedAt  time.Time `gorm:"autoCreateTime;index" json:"createdAt"`
}

// TableName 指定表名
func (Message) TableName() string {
	return "messages"
}

// Notification 通知
type Notification struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     uint64    `gorm:"not null;index" json:"userId"`
	Type       string    `gorm:"type:varchar(64);not null;index" json:"type"` // order/payment/follow/system
	Title      string    `gorm:"type:varchar(255);not null" json:"title"`
	Content    string    `gorm:"type:text" json:"content"`
	RelatedID  *uint64   `json:"relatedId"` // 关联的实体ID
	IsRead     bool      `gorm:"default:false;index" json:"isRead"`
	ReadAt     *time.Time `json:"readAt"`
	CreatedAt  time.Time `gorm:"autoCreateTime;index" json:"createdAt"`
}

// TableName 指定表名
func (Notification) TableName() string {
	return "notifications"
}

// PlayerMoment 陪玩师动态
type PlayerMoment struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	PlayerID  uint64    `gorm:"not null;index" json:"playerId"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Images    string    `gorm:"type:text" json:"images"` // JSON数组
	LikeCount int64     `gorm:"default:0" json:"likeCount"`
	ViewCount int64     `gorm:"default:0" json:"viewCount"`
	CreatedAt time.Time `gorm:"autoCreateTime;index" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

// TableName 指定表名
func (PlayerMoment) TableName() string {
	return "player_moments"
}

// MomentLike 动态点赞
type MomentLike struct {
	ID       uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	MomentID uint64    `gorm:"not null;index:idx_moment_user" json:"momentId"`
	UserID   uint64    `gorm:"not null;index:idx_moment_user" json:"userId"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

// TableName 指定表名
func (MomentLike) TableName() string {
	return "moment_likes"
}

// MomentComment 动态评论
type MomentComment struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	MomentID  uint64    `gorm:"not null;index" json:"momentId"`
	UserID    uint64    `gorm:"not null;index" json:"userId"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	ParentID  *uint64   `gorm:"index" json:"parentId"` // 父评论ID（用于回复）
	CreatedAt time.Time `gorm:"autoCreateTime;index" json:"createdAt"`
}

// TableName 指定表名
func (MomentComment) TableName() string {
	return "moment_comments"
}

