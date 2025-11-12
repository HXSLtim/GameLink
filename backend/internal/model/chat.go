package model

import "time"

// ChatGroupType represents the classification of a chat group.
type ChatGroupType string

// Supported chat group types.
const (
	ChatGroupTypePublic ChatGroupType = "public"
	ChatGroupTypeOrder  ChatGroupType = "order"
)

// ChatMessageType enumerates supported chat message payload categories.
type ChatMessageType string

// Supported chat message types.
const (
	ChatMessageTypeText   ChatMessageType = "text"
	ChatMessageTypeImage  ChatMessageType = "image"
	ChatMessageTypeFile   ChatMessageType = "file"
	ChatMessageTypeSystem ChatMessageType = "system"
)

// ChatMessageAuditStatus represents moderation state of a message.
type ChatMessageAuditStatus string

// Supported chat message audit statuses.
const (
	ChatMessageAuditPending  ChatMessageAuditStatus = "pending"
	ChatMessageAuditApproved ChatMessageAuditStatus = "approved"
	ChatMessageAuditRejected ChatMessageAuditStatus = "rejected"
)

// ChatGroup defines a chat room entity.
type ChatGroup struct {
	Base
	GroupName      string        `json:"groupName" gorm:"size:128;not null"`
	GroupType      ChatGroupType `json:"groupType" gorm:"type:varchar(32);not null;index"`
	RelatedOrderID *uint64       `json:"relatedOrderId,omitempty" gorm:"column:related_order_id;index"`
	CreatedBy      uint64        `json:"createdBy" gorm:"column:created_by;not null;index"`
	MaxMembers     int           `json:"maxMembers" gorm:"column:max_members;default:100"`
	IsActive       bool          `json:"isActive" gorm:"column:is_active;default:true;index"`
	AutoDestroy    bool          `json:"autoDestroy" gorm:"column:auto_destroy;default:false"`
	DeactivatedAt  *time.Time    `json:"deactivatedAt" gorm:"column:deactivated_at;index"`
	AvatarURL      string        `json:"avatarUrl" gorm:"column:avatar_url;size:255"`
	Description    string        `json:"description" gorm:"type:text"`
	Settings       string        `json:"settings" gorm:"type:json"`

	Members []ChatGroupMember `json:"members" gorm:"foreignKey:GroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// ChatGroupMember binds users to chat groups.
type ChatGroupMember struct {
	Base
	GroupID           uint64     `json:"groupId" gorm:"column:group_id;not null;index"`
	UserID            uint64     `json:"userId" gorm:"column:user_id;not null;index"`
	Role              string     `json:"role" gorm:"size:32;default:'member'"`
	Nickname          string     `json:"nickname" gorm:"size:64"`
	JoinedAt          time.Time  `json:"joinedAt" gorm:"column:joined_at;index"`
	LastReadAt        *time.Time `json:"lastReadAt" gorm:"column:last_read_at"`
	LastReadMessageID *uint64    `json:"lastReadMessageId" gorm:"column:last_read_message_id"`
	IsMuted           bool       `json:"isMuted" gorm:"column:is_muted;default:false"`
	IsActive          bool       `json:"isActive" gorm:"column:is_active;default:true"`

	Group ChatGroup `json:"-" gorm:"foreignKey:GroupID;references:ID"`
}

// ChatMessage represents persisted chat messages.
type ChatMessage struct {
	Base
	GroupID     uint64          `json:"groupId" gorm:"column:group_id;not null;index"`
	SenderID    uint64          `json:"senderId" gorm:"column:sender_id;not null;index"`
	Content     string          `json:"content" gorm:"type:text;not null"`
	MessageType ChatMessageType `json:"messageType" gorm:"column:message_type;type:varchar(16);default:'text'"`
	ReplyToID   *uint64         `json:"replyToId" gorm:"column:reply_to_id"`
	ImageURL    string          `json:"imageUrl" gorm:"column:image_url;size:255"`
	Metadata    string          `json:"metadata" gorm:"type:json"`
	IsDeleted   bool            `json:"isDeleted" gorm:"column:is_deleted;default:false"`
	AuditStatus ChatMessageAuditStatus `json:"auditStatus" gorm:"column:audit_status;type:varchar(16);default:'pending';index"`
	ModeratedBy *uint64                 `json:"moderatedBy" gorm:"column:moderated_by"`
	ModeratedAt *time.Time              `json:"moderatedAt" gorm:"column:moderated_at"`
	RejectReason string                 `json:"rejectReason" gorm:"column:reject_reason;type:text"`

	Group ChatGroup `json:"-" gorm:"foreignKey:GroupID;references:ID"`
}

// TableName overrides default table name for chat messages.
func (ChatMessage) TableName() string { return "chat_messages" }

// TableName overrides default table name for chat groups.
func (ChatGroup) TableName() string { return "chat_groups" }
