package model

import "time"

// FeedVisibility defines who can see the feed.
type FeedVisibility string

const (
	// FeedVisibilityPublic means visible to everyone.
	FeedVisibilityPublic FeedVisibility = "public"
	// FeedVisibilityFollowers restricts to followers only.
	FeedVisibilityFollowers FeedVisibility = "followers"
	// FeedVisibilityPrivate hides the feed from others.
	FeedVisibilityPrivate FeedVisibility = "private"
)

// FeedModerationStatus captures moderation pipeline state.
type FeedModerationStatus string

const (
	// FeedModerationPending indicates awaiting review.
	FeedModerationPending FeedModerationStatus = "pending"
	// FeedModerationApproved indicates feed is visible.
	FeedModerationApproved FeedModerationStatus = "approved"
	// FeedModerationRejected indicates content rejected.
	FeedModerationRejected FeedModerationStatus = "rejected"
	// FeedModerationRemoved indicates removed after publish.
	FeedModerationRemoved FeedModerationStatus = "removed"
)

// FeedMetricFields stores counters for feed interactions.
type FeedMetricFields struct {
	LikeCount   uint64 `json:"likeCount" gorm:"column:metrics_like_count;default:0"`
	ReplyCount  uint64 `json:"replyCount" gorm:"column:metrics_reply_count;default:0"`
	ReportCount uint64 `json:"reportCount" gorm:"column:metrics_report_count;default:0"`
	ViewCount   uint64 `json:"viewCount" gorm:"column:metrics_view_count;default:0"`
	ShareCount  uint64 `json:"shareCount" gorm:"column:metrics_share_count;default:0"`
}

// Feed represents a community feed item.
type Feed struct {
	Base
	AuthorID          uint64               `json:"authorId" gorm:"column:author_id;index"`
	Content           string               `json:"content" gorm:"column:content;type:text"`
	CategoryID        *uint64              `json:"categoryId,omitempty" gorm:"column:category_id;index"`
	Visibility        FeedVisibility       `json:"visibility" gorm:"column:visibility;type:varchar(32);default:'public'"`
	ModerationStatus  FeedModerationStatus `json:"moderationStatus" gorm:"column:moderation_status;type:varchar(32);index;default:'pending'"`
	ModerationNote    string               `json:"moderationNote,omitempty" gorm:"column:moderation_note;type:text"`
	ModeratedBy       *uint64              `json:"moderatedBy,omitempty" gorm:"column:moderated_by"`
	AutoModeratedAt   *time.Time           `json:"autoModeratedAt,omitempty" gorm:"column:auto_moderated_at"`
	ManualModeratedAt *time.Time           `json:"manualModeratedAt,omitempty" gorm:"column:manual_moderated_at"`
	Metrics           FeedMetricFields     `json:"metrics" gorm:"embedded"`
	Images            []FeedImage          `json:"images" gorm:"foreignKey:FeedID;constraint:OnDelete:CASCADE"`
	Category          *ContentCategory     `json:"category,omitempty" gorm:"foreignKey:CategoryID;references:ID"`
}

// TableName implements gorm tabler.
func (Feed) TableName() string { return "feeds" }

// FeedImage stores metadata for feed images.
type FeedImage struct {
	Base
	FeedID    uint64 `json:"feedId" gorm:"column:feed_id;index"`
	URL       string `json:"url" gorm:"column:url;type:text"`
	Order     int    `json:"order" gorm:"column:display_order;default:0"`
	Width     int    `json:"width" gorm:"column:width;default:0"`
	Height    int    `json:"height" gorm:"column:height;default:0"`
	SizeBytes int64  `json:"sizeBytes" gorm:"column:size_bytes;default:0"`
}

// TableName implements gorm tabler.
func (FeedImage) TableName() string { return "feed_images" }

// FeedReport records user reports for moderation.
type FeedReport struct {
	Base
	FeedID    uint64     `json:"feedId" gorm:"column:feed_id;index"`
	Reporter  uint64     `json:"reporterId" gorm:"column:reporter_id;index"`
	Reason    string     `json:"reason" gorm:"column:reason;type:text"`
	Status    string     `json:"status" gorm:"column:status;type:varchar(32);default:'pending'"`
	Result    string     `json:"result" gorm:"column:result;type:text"`
	HandledBy *uint64    `json:"handledBy,omitempty" gorm:"column:handled_by"`
	HandledAt *time.Time `json:"handledAt,omitempty" gorm:"column:handled_at"`
}

// TableName implements gorm tabler.
func (FeedReport) TableName() string { return "feed_reports" }

// NotificationPriority defines severity of notification event.
type NotificationPriority string

const (
	// NotificationPriorityLow indicates informational message.
	NotificationPriorityLow NotificationPriority = "low"
	// NotificationPriorityNormal is default priority.
	NotificationPriorityNormal NotificationPriority = "normal"
	// NotificationPriorityHigh is urgent message.
	NotificationPriorityHigh NotificationPriority = "high"
)

// NotificationEvent captures notification payload for a user.
type NotificationEvent struct {
	Base
	UserID        uint64               `json:"userId" gorm:"column:user_id;index"`
	Title         string               `json:"title" gorm:"column:title;type:varchar(255)"`
	Message       string               `json:"message" gorm:"column:message;type:text"`
	Channel       string               `json:"channel" gorm:"column:channel;type:varchar(32);default:'web'"`
	Priority      NotificationPriority `json:"priority" gorm:"column:priority;type:varchar(16);default:'normal'"`
	ReferenceType string               `json:"referenceType" gorm:"column:reference_type;type:varchar(64)"`
	ReferenceID   *uint64              `json:"referenceId,omitempty" gorm:"column:reference_id"`
	ReadAt        *time.Time           `json:"readAt,omitempty" gorm:"column:read_at;index"`
	Metadata      string               `json:"metadata,omitempty" gorm:"column:metadata;type:text"`
}

// TableName implements gorm tabler.
func (NotificationEvent) TableName() string { return "notification_events" }

// ReviewReply represents replies from players to reviews.
type ReviewReply struct {
	Base
	ReviewID       uint64            `json:"reviewId" gorm:"column:review_id;not null;index"`
	PlayerID       uint64            `json:"playerId" gorm:"column:player_id;not null;index"` // 陪玩师ID
	AuthorID       uint64            `json:"authorId" gorm:"column:author_id;index"`          // 向后兼容
	Content        string            `json:"content" gorm:"column:content;type:text;not null"`
	ReplyCount     int               `json:"replyCount" gorm:"column:reply_count;default:1"`              // 回复次数（最多3次）
	Status         ReviewReplyStatus `json:"status" gorm:"column:status;size:20;default:'pending';index"` // 状态
	ModerationNote string            `json:"moderationNote,omitempty" gorm:"column:moderation_note;type:text"`
	ModeratedAt    *time.Time        `json:"moderatedAt,omitempty" gorm:"column:moderated_at"`

	// Relations
	Review *Review `json:"review,omitempty" gorm:"foreignKey:ReviewID"`
	Player *Player `json:"player,omitempty" gorm:"foreignKey:PlayerID"`
}

// TableName implements gorm tabler.
func (ReviewReply) TableName() string { return "review_replies" }

// ReviewAppeal 评价申诉（陪玩师对差评的申诉）
type ReviewAppeal struct {
	Base
	ReviewID     uint64       `json:"reviewId" gorm:"column:review_id;not null;index"`              // 评价ID
	PlayerID     uint64       `json:"playerId" gorm:"column:player_id;not null;index"`              // 陪玩师ID
	Reason       string       `json:"reason" gorm:"type:text;not null"`                             // 申诉原因
	EvidenceURLs string       `json:"evidenceUrls,omitempty" gorm:"column:evidence_urls;type:json"` // 证据截图（JSON数组）
	Status       AppealStatus `json:"status" gorm:"size:20;default:'pending';index"`                // 状态
	HandledBy    *uint64      `json:"handledBy,omitempty" gorm:"column:handled_by"`                 // 处理人ID
	HandledAt    *time.Time   `json:"handledAt,omitempty" gorm:"column:handled_at"`                 // 处理时间
	HandleRemark string       `json:"handleRemark,omitempty" gorm:"column:handle_remark;type:text"` // 处理备注

	// Relations
	Review  *Review `json:"review,omitempty" gorm:"foreignKey:ReviewID"`
	Player  *Player `json:"player,omitempty" gorm:"foreignKey:PlayerID"`
	Handler *User   `json:"handler,omitempty" gorm:"foreignKey:HandledBy"`
}

// TableName implements gorm tabler.
func (ReviewAppeal) TableName() string { return "review_appeals" }
