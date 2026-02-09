package model

import "time"

// ChatGroupType represents the classification of a chat group.
type ChatGroupType string

// Supported chat group types.
const (
	ChatGroupTypePublic  ChatGroupType = "public"  // 公开群组
	ChatGroupTypePrivate ChatGroupType = "private" // 私聊群组
	ChatGroupTypeOrder   ChatGroupType = "order"   // 订单房间
	ChatGroupTypeTeam    ChatGroupType = "team"    // 组队房间
	ChatGroupTypeLFG     ChatGroupType = "lfg"     // 快速匹配房间
	ChatGroupTypeCustom  ChatGroupType = "custom"  // 自定义房间
)

// ChatGroupStatus 房间状态 (用于游戏房间)
// @Enum waiting, ready, in_game, finished, canceled
type ChatGroupStatus string

const (
	ChatGroupStatusWaiting  ChatGroupStatus = "waiting"  // 等待中
	ChatGroupStatusReady    ChatGroupStatus = "ready"    // 准备就绪
	ChatGroupStatusInGame   ChatGroupStatus = "in_game"  // 游戏中
	ChatGroupStatusFinished ChatGroupStatus = "finished" // 已结束
	ChatGroupStatusCanceled ChatGroupStatus = "canceled" // 已取消
)

// ChatMessageType enumerates supported chat message payload categories.
type ChatMessageType string

// Supported chat message types.
const (
	ChatMessageTypeText   ChatMessageType = "text"
	ChatMessageTypeImage  ChatMessageType = "image"
	ChatMessageTypeFile   ChatMessageType = "file"
	ChatMessageTypeSystem ChatMessageType = "system"
	ChatMessageTypeVoice  ChatMessageType = "voice" // 语音消息（预留）
	ChatMessageTypeEmoji  ChatMessageType = "emoji" // 表情包（预留）
)

// ChatMessageAuditStatus represents moderation state of a message.
type ChatMessageAuditStatus string

// Supported chat message audit statuses.
const (
	ChatMessageAuditPending  ChatMessageAuditStatus = "pending"
	ChatMessageAuditApproved ChatMessageAuditStatus = "approved"
	ChatMessageAuditRejected ChatMessageAuditStatus = "rejected"
	ChatMessageAuditDeleted  ChatMessageAuditStatus = "deleted"
)

// ChatGroup defines a chat room entity.
type ChatGroup struct {
	Base
	GroupName            string        `json:"groupName" gorm:"size:128;not null"`
	GroupType            ChatGroupType `json:"groupType" gorm:"type:varchar(32);not null;index"`
	RelatedOrderID       *uint64       `json:"relatedOrderId,omitempty" gorm:"column:related_order_id;index"`
	CreatedBy            uint64        `json:"createdBy" gorm:"column:created_by;not null;index"`
	MaxMembers           int           `json:"maxMembers" gorm:"column:max_members;default:100"`
	IsActive             bool          `json:"isActive" gorm:"column:is_active;default:true;index"`
	AutoDestroy          bool          `json:"autoDestroy" gorm:"column:auto_destroy;default:false"`
	DeactivatedAt        *time.Time    `json:"deactivatedAt" gorm:"column:deactivated_at;index"`
	AvatarURL            string        `json:"avatarUrl" gorm:"column:avatar_url;size:255"`
	Description          string        `json:"description" gorm:"type:text"`
	Settings             string        `json:"settings" gorm:"type:json;default:'{}'"`
	MessageRetentionDays int           `json:"messageRetentionDays" gorm:"column:message_retention_days;default:30"` // 消息保留天数（默认30）

	// 游戏房间字段 (Discord/Kook 风格)
	GameID         *uint64         `json:"gameId,omitempty" gorm:"column:game_id;index"`                         // 关联游戏
	RoomStatus     ChatGroupStatus `json:"roomStatus" gorm:"column:room_status;size:32;default:'waiting';index"` // 房间状态
	IsPrivate      bool            `json:"isPrivate" gorm:"column:is_private;default:false"`                     // 是否私密房间
	Password       string          `json:"-" gorm:"column:password;size:64"`                                     // 房间密码（不返回给前端）
	CurrentMembers int             `json:"currentMembers" gorm:"column:current_members;default:0"`               // 当前成员数
	RelatedTeamID  *uint64         `json:"relatedTeamId,omitempty" gorm:"column:related_team_id;index"`          // 关联战队
	RelatedLFGID   *uint64         `json:"relatedLfgId,omitempty" gorm:"column:related_lfg_id;index"`            // 关联LFG请求
	StartedAt      *time.Time      `json:"startedAt,omitempty" gorm:"column:started_at"`                         // 游戏开始时间
	FinishedAt     *time.Time      `json:"finishedAt,omitempty" gorm:"column:finished_at"`                       // 游戏结束时间

	// 语音服务字段 (TRTC)
	VoiceEnabled    bool       `json:"voiceEnabled" gorm:"column:voice_enabled;default:false"`     // 是否启用语音
	VoiceRoomID     string     `json:"voiceRoomId" gorm:"column:voice_room_id;size:128;index"`     // 语音房间ID（第三方服务）
	VoiceProvider   string     `json:"voiceProvider" gorm:"column:voice_provider;size:32"`         // 语音服务商（trtc）
	VoiceSDKAppID   uint64     `json:"voiceSdkAppId" gorm:"column:voice_sdk_app_id"`               // TRTC SDK AppID
	VoiceStartedAt  *time.Time `json:"voiceStartedAt" gorm:"column:voice_started_at"`              // 语音开始时间
	VoiceEndedAt    *time.Time `json:"voiceEndedAt" gorm:"column:voice_ended_at"`                  // 语音结束时间
	VoiceDuration   int        `json:"voiceDuration" gorm:"column:voice_duration;default:0"`       // 语音时长（秒）
	VoiceRecordURL  string     `json:"voiceRecordUrl" gorm:"column:voice_record_url;size:255"`     // 语音录制文件URL
	VoiceMaxMembers int        `json:"voiceMaxMembers" gorm:"column:voice_max_members;default:10"` // 语音最大人数

	// Relations
	Members []ChatGroupMember `json:"members" gorm:"foreignKey:GroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Game    *Game             `json:"game,omitempty" gorm:"foreignKey:GameID"`
}

// ChatMemberRole 聊天成员角色
type ChatMemberRole string

const (
	ChatMemberRoleOwner  ChatMemberRole = "owner"  // 群主
	ChatMemberRoleAdmin  ChatMemberRole = "admin"  // 管理员
	ChatMemberRoleMember ChatMemberRole = "member" // 普通成员
)

// ChatGroupMember binds users to chat groups.
type ChatGroupMember struct {
	Base
	GroupID           uint64         `json:"groupId" gorm:"column:group_id;not null;index;uniqueIndex:idx_group_user"`
	UserID            uint64         `json:"userId" gorm:"column:user_id;not null;index;uniqueIndex:idx_group_user"`
	Role              ChatMemberRole `json:"role" gorm:"size:32;default:'member'"`
	Nickname          string         `json:"nickname" gorm:"size:64"`
	JoinedAt          time.Time      `json:"joinedAt" gorm:"column:joined_at;index"`
	LastReadAt        *time.Time     `json:"lastReadAt" gorm:"column:last_read_at"`
	LastReadMessageID *uint64        `json:"lastReadMessageId" gorm:"column:last_read_message_id"`
	IsMuted           bool           `json:"isMuted" gorm:"column:is_muted;default:false"`
	MutedUntil        *time.Time     `json:"mutedUntil,omitempty" gorm:"column:muted_until;index"`
	MutedBy           *uint64        `json:"mutedBy,omitempty" gorm:"column:muted_by"`
	MuteReason        string         `json:"muteReason,omitempty" gorm:"column:mute_reason;type:text"`
	IsActive          bool           `json:"isActive" gorm:"column:is_active;default:true"`
	IsReady           bool           `json:"isReady" gorm:"column:is_ready;default:false"` // 游戏房间准备状态

	Group ChatGroup `json:"-" gorm:"foreignKey:GroupID;references:ID"`
}

// ChatMessage represents persisted chat messages.
//
// Covering Index Notes:
//   - idx_chat_messages_group_sent_covering: PostgreSQL covering index for message list queries
//   - Created via migration: api/migrations/0001_add_covering_indexes.sql
//   - Covers: SELECT id, content, sender_id FROM chat_messages WHERE group_id = ? ORDER BY created_at DESC LIMIT 50
//   - Index columns: (group_id, created_at DESC)
//   - INCLUDE columns: (id, content, sender_id, message_type, audit_status)
//   - Benefit: Index-only scan for chat history, no heap fetch needed
type ChatMessage struct {
	Base
	GroupID      uint64                 `json:"groupId" gorm:"column:group_id;not null;index"`
	SenderID     uint64                 `json:"senderId" gorm:"column:sender_id;not null;index"`
	Content      string                 `json:"content" gorm:"type:text;not null"`
	MessageType  ChatMessageType        `json:"messageType" gorm:"column:message_type;type:varchar(16);default:'text'"`
	ReplyToID    *uint64                `json:"replyToId" gorm:"column:reply_to_id"`
	ImageURL     string                 `json:"imageUrl" gorm:"column:image_url;size:255"`
	Metadata     string                 `json:"metadata" gorm:"type:json;default:'{}'"`
	IsDeleted    bool                   `json:"isDeleted" gorm:"column:is_deleted;default:false"`
	AuditStatus  ChatMessageAuditStatus `json:"auditStatus" gorm:"column:audit_status;type:varchar(16);default:'pending';index"`
	ModeratedBy  *uint64                `json:"moderatedBy" gorm:"column:moderated_by"`
	ModeratedAt  *time.Time             `json:"moderatedAt" gorm:"column:moderated_at"`
	RejectReason string                 `json:"rejectReason" gorm:"column:reject_reason;type:text"`

	Group ChatGroup `json:"-" gorm:"foreignKey:GroupID;references:ID"`
}

// TableName overrides default table name for chat messages.
func (ChatMessage) TableName() string { return "chat_messages" }

// TableName overrides default table name for chat groups.
func (ChatGroup) TableName() string { return "chat_groups" }
