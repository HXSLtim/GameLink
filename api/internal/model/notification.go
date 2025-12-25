package model

import "time"

// NotificationType 通知类型
type NotificationType string

const (
	NotificationTypeOrderStatus   NotificationType = "order_status"   // 订单状态变更
	NotificationTypeVipExpire     NotificationType = "vip_expire"     // VIP到期提醒
	NotificationTypeCouponExpire  NotificationType = "coupon_expire"  // 优惠券过期提醒
	NotificationTypeActivityStart NotificationType = "activity_start" // 活动开始提醒
	NotificationTypeActivityEnd   NotificationType = "activity_end"   // 活动结束提醒
	NotificationTypeSystem        NotificationType = "system"         // 系统公告
	NotificationTypePromotion     NotificationType = "promotion"      // 营销推广（预留）
	NotificationTypeChat          NotificationType = "chat"           // 聊天消息（预留）
	NotificationTypeReview        NotificationType = "review"         // 评价提醒
	NotificationTypeReviewReply   NotificationType = "review_reply"   // 评价回复
)

// NotificationChannel 通知渠道
type NotificationChannel string

const (
	NotificationChannelInApp  NotificationChannel = "in_app" // 站内消息
	NotificationChannelPush   NotificationChannel = "push"   // App推送
	NotificationChannelSMS    NotificationChannel = "sms"    // 短信（预留）
	NotificationChannelWechat NotificationChannel = "wechat" // 微信模板消息（预留）
	NotificationChannelEmail  NotificationChannel = "email"  // 邮件（预留）
)

// NotificationStatus 通知状态
type NotificationStatus string

const (
	NotificationStatusPending  NotificationStatus = "pending"  // 待发送
	NotificationStatusSent     NotificationStatus = "sent"     // 已发送
	NotificationStatusRead     NotificationStatus = "read"     // 已读
	NotificationStatusFailed   NotificationStatus = "failed"   // 发送失败
	NotificationStatusCanceled NotificationStatus = "canceled" // 已取消
)

// NotificationTemplate 通知模板
type NotificationTemplate struct {
	Base
	Code      string           `json:"code" gorm:"size:64;uniqueIndex"`                // 模板编码（唯一）
	Name      string           `json:"name" gorm:"size:128;not null"`                  // 模板名称
	Type      NotificationType `json:"type" gorm:"size:32;not null;index"`             // 通知类型
	Title     string           `json:"title" gorm:"size:256"`                          // 标题模板（支持变量）
	Content   string           `json:"content" gorm:"type:text;not null"`              // 内容模板（支持变量）
	Channels  string           `json:"channels" gorm:"size:256"`                       // 支持的渠道（JSON数组）
	Variables string           `json:"variables,omitempty" gorm:"type:text"`           // 变量说明（JSON）
	IsActive  bool             `json:"isActive" gorm:"column:is_active;default:true"`  // 是否启用
	IsSystem  bool             `json:"isSystem" gorm:"column:is_system;default:false"` // 是否系统模板（不可删除）

	// 推送配置
	PushTitle   string `json:"pushTitle,omitempty" gorm:"column:push_title;size:128"`     // 推送标题
	PushContent string `json:"pushContent,omitempty" gorm:"column:push_content;size:256"` // 推送内容

	// 短信配置（预留）
	SMSTemplateID string `json:"smsTemplateId,omitempty" gorm:"column:sms_template_id;size:64"` // 短信模板ID

	// 微信配置（预留）
	WechatTemplateID string `json:"wechatTemplateId,omitempty" gorm:"column:wechat_template_id;size:64"` // 微信模板ID
}

// TableName 指定表名
func (NotificationTemplate) TableName() string {
	return "notification_templates"
}

// UserNotification 用户通知（扩展版）
// 注意：与 social.go 中的 Notification 共存，此模型用于更复杂的通知场景
type UserNotification struct {
	Base
	UserID     uint64              `json:"userId" gorm:"column:user_id;not null;index"`          // 用户ID
	TemplateID *uint64             `json:"templateId,omitempty" gorm:"column:template_id;index"` // 模板ID
	Type       NotificationType    `json:"type" gorm:"size:32;not null;index"`                   // 通知类型
	Channel    NotificationChannel `json:"channel" gorm:"size:32;not null;index"`                // 通知渠道
	Title      string              `json:"title" gorm:"size:256"`                                // 标题
	Content    string              `json:"content" gorm:"type:text;not null"`                    // 内容
	Status     NotificationStatus  `json:"status" gorm:"size:32;default:'pending';index"`        // 状态
	ReadAt     *time.Time          `json:"readAt,omitempty" gorm:"column:read_at"`               // 已读时间
	SentAt     *time.Time          `json:"sentAt,omitempty" gorm:"column:sent_at"`               // 发送时间

	// 关联数据
	RelatedType string  `json:"relatedType,omitempty" gorm:"column:related_type;size:32;index"` // 关联类型：order/coupon/activity/vip
	RelatedID   *uint64 `json:"relatedId,omitempty" gorm:"column:related_id;index"`             // 关联ID

	// 推送相关
	PushID        string `json:"pushId,omitempty" gorm:"column:push_id;size:128"`               // 推送ID（第三方返回）
	FailureReason string `json:"failureReason,omitempty" gorm:"column:failure_reason;size:512"` // 失败原因

	// Relations
	User     *User                 `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Template *NotificationTemplate `json:"template,omitempty" gorm:"foreignKey:TemplateID"`
}

// TableName 指定表名
func (UserNotification) TableName() string {
	return "user_notifications"
}

// IsRead 是否已读
func (n *UserNotification) IsRead() bool {
	return n.Status == NotificationStatusRead || n.ReadAt != nil
}

// UserNotificationSetting 用户通知设置
type UserNotificationSetting struct {
	Base
	UserID uint64 `json:"userId" gorm:"column:user_id;uniqueIndex"` // 用户ID（唯一）

	// 通知类型开关
	OrderStatusEnabled  bool `json:"orderStatusEnabled" gorm:"column:order_status_enabled;default:true"`   // 订单状态通知
	VipExpireEnabled    bool `json:"vipExpireEnabled" gorm:"column:vip_expire_enabled;default:true"`       // VIP到期提醒
	CouponExpireEnabled bool `json:"couponExpireEnabled" gorm:"column:coupon_expire_enabled;default:true"` // 优惠券过期提醒
	ActivityEnabled     bool `json:"activityEnabled" gorm:"column:activity_enabled;default:true"`          // 活动提醒
	SystemEnabled       bool `json:"systemEnabled" gorm:"column:system_enabled;default:true"`              // 系统公告
	PromotionEnabled    bool `json:"promotionEnabled" gorm:"column:promotion_enabled;default:true"`        // 营销推广（预留）
	ChatEnabled         bool `json:"chatEnabled" gorm:"column:chat_enabled;default:true"`                  // 聊天消息（预留）

	// 渠道开关
	InAppEnabled  bool `json:"inAppEnabled" gorm:"column:in_app_enabled;default:true"`   // 站内消息
	PushEnabled   bool `json:"pushEnabled" gorm:"column:push_enabled;default:true"`      // App推送
	SMSEnabled    bool `json:"smsEnabled" gorm:"column:sms_enabled;default:false"`       // 短信（预留）
	WechatEnabled bool `json:"wechatEnabled" gorm:"column:wechat_enabled;default:false"` // 微信（预留）
	EmailEnabled  bool `json:"emailEnabled" gorm:"column:email_enabled;default:false"`   // 邮件（预留）

	// 免打扰设置
	DoNotDisturbEnabled bool   `json:"doNotDisturbEnabled" gorm:"column:do_not_disturb_enabled;default:false"` // 是否启用免打扰
	DoNotDisturbStart   string `json:"doNotDisturbStart,omitempty" gorm:"column:do_not_disturb_start;size:8"`  // 免打扰开始时间（HH:MM）
	DoNotDisturbEnd     string `json:"doNotDisturbEnd,omitempty" gorm:"column:do_not_disturb_end;size:8"`      // 免打扰结束时间（HH:MM）

	// Relations
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (UserNotificationSetting) TableName() string {
	return "user_notification_settings"
}

// IsInDoNotDisturbPeriod 是否在免打扰时段
func (s *UserNotificationSetting) IsInDoNotDisturbPeriod() bool {
	if !s.DoNotDisturbEnabled || s.DoNotDisturbStart == "" || s.DoNotDisturbEnd == "" {
		return false
	}
	now := time.Now().Format("15:04")
	// 简单比较，不处理跨天情况
	if s.DoNotDisturbStart <= s.DoNotDisturbEnd {
		return now >= s.DoNotDisturbStart && now <= s.DoNotDisturbEnd
	}
	// 跨天情况（如 22:00 - 08:00）
	return now >= s.DoNotDisturbStart || now <= s.DoNotDisturbEnd
}

// NotificationConfig 通知系统配置
type NotificationConfig struct {
	Base
	ConfigKey   string `json:"configKey" gorm:"column:config_key;size:64;uniqueIndex"` // 配置键
	ConfigValue string `json:"configValue" gorm:"column:config_value;type:text"`       // 配置值
	Description string `json:"description" gorm:"size:255"`                            // 描述
}

// TableName 指定表名
func (NotificationConfig) TableName() string {
	return "notification_configs"
}

// 通知配置键常量
const (
	NotificationConfigVipExpireDays    = "vip_expire_days"    // VIP到期提前提醒天数（JSON数组，如 [7,3,1]）
	NotificationConfigCouponExpireDays = "coupon_expire_days" // 优惠券过期提前提醒天数（JSON数组）
	NotificationConfigPushProvider     = "push_provider"      // 推送服务商：jpush/getui/umeng
	NotificationConfigSMSProvider      = "sms_provider"       // 短信服务商（预留）
)

// NotificationScheduleStatus 定时通知任务状态
type NotificationScheduleStatus string

const (
	NotificationScheduleStatusPending    NotificationScheduleStatus = "pending"    // 待执行
	NotificationScheduleStatusProcessing NotificationScheduleStatus = "processing" // 执行中
	NotificationScheduleStatusCompleted  NotificationScheduleStatus = "completed"  // 已完成
	NotificationScheduleStatusFailed     NotificationScheduleStatus = "failed"     // 失败
)

// NotificationSchedule 定时通知任务
type NotificationSchedule struct {
	Base
	Name        string                     `json:"name" gorm:"size:128;not null"`                          // 任务名称
	Type        NotificationType           `json:"type" gorm:"size:32;not null;index"`                     // 通知类型
	TemplateID  uint64                     `json:"templateId" gorm:"column:template_id;not null"`          // 模板ID
	ScheduleAt  time.Time                  `json:"scheduleAt" gorm:"column:schedule_at;not null;index"`    // 计划发送时间
	Status      NotificationScheduleStatus `json:"status" gorm:"size:32;default:'pending';index"`          // 状态
	TargetType  string                     `json:"targetType" gorm:"column:target_type;size:32"`           // 目标类型：all/vip/specific
	TargetIDs   string                     `json:"targetIds,omitempty" gorm:"column:target_ids;type:text"` // 目标用户ID列表（JSON数组）
	TotalCount  int                        `json:"totalCount" gorm:"column:total_count;default:0"`         // 总发送数
	SentCount   int                        `json:"sentCount" gorm:"column:sent_count;default:0"`           // 已发送数
	FailedCount int                        `json:"failedCount" gorm:"column:failed_count;default:0"`       // 失败数
	StartedAt   *time.Time                 `json:"startedAt,omitempty" gorm:"column:started_at"`           // 开始时间
	CompletedAt *time.Time                 `json:"completedAt,omitempty" gorm:"column:completed_at"`       // 完成时间

	// Relations
	Template *NotificationTemplate `json:"template,omitempty" gorm:"foreignKey:TemplateID"`
}

// TableName 指定表名
func (NotificationSchedule) TableName() string {
	return "notification_schedules"
}
