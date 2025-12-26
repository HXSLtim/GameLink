package model

import "time"

// ActivityStatus 活动状态
type ActivityStatus string

const (
	ActivityStatusDraft    ActivityStatus = "draft"    // 草稿
	ActivityStatusPreheat  ActivityStatus = "preheat"  // 预热中（可见但不可参与）
	ActivityStatusActive   ActivityStatus = "active"   // 进行中
	ActivityStatusPaused   ActivityStatus = "paused"   // 已暂停
	ActivityStatusEnded    ActivityStatus = "ended"    // 已结束
	ActivityStatusCanceled ActivityStatus = "canceled" // 已取消
)

// ActivityType 活动类型
type ActivityType string

const (
	ActivityTypeCoupon   ActivityType = "coupon"   // 优惠券发放活动
	ActivityTypeDiscount ActivityType = "discount" // 限时折扣活动（预留）
	ActivityTypeGift     ActivityType = "gift"     // 赠品活动（预留）
)

// Activity 活动
type Activity struct {
	Base
	Name        string         `json:"name" gorm:"size:128;not null"`                         // 活动名称
	Description string         `json:"description" gorm:"type:text"`                          // 活动描述
	Type        ActivityType   `json:"type" gorm:"size:32;not null;index"`                    // 活动类型
	Status      ActivityStatus `json:"status" gorm:"size:32;default:'draft';index"`           // 活动状态
	CoverURL    string         `json:"coverUrl,omitempty" gorm:"column:cover_url;size:512"`   // 活动封面图
	BannerURL   string         `json:"bannerUrl,omitempty" gorm:"column:banner_url;size:512"` // 活动Banner图

	// 时间控制
	PreheatAt *time.Time `json:"preheatAt,omitempty" gorm:"column:preheat_at"` // 预热开始时间
	StartAt   time.Time  `json:"startAt" gorm:"column:start_at;not null"`      // 活动开始时间
	EndAt     time.Time  `json:"endAt" gorm:"column:end_at;not null"`          // 活动结束时间

	// 参与限制
	TotalLimit   int `json:"totalLimit" gorm:"column:total_limit;default:0"`      // 总参与次数限制（0=无限制）
	DailyLimit   int `json:"dailyLimit" gorm:"column:daily_limit;default:0"`      // 每日参与次数限制（0=无限制）
	PerUserLimit int `json:"perUserLimit" gorm:"column:per_user_limit;default:1"` // 每人参与次数限制

	// 统计
	TotalParticipants int `json:"totalParticipants" gorm:"column:total_participants;default:0"` // 总参与人数
	TodayParticipants int `json:"todayParticipants" gorm:"column:today_participants;default:0"` // 今日参与人数
	TotalClaimed      int `json:"totalClaimed" gorm:"column:total_claimed;default:0"`           // 总领取次数

	// 配置
	AllowVipStack bool   `json:"allowVipStack" gorm:"column:allow_vip_stack;default:false"` // 是否允许与VIP折扣叠加
	Rules         string `json:"rules,omitempty" gorm:"type:text"`                          // 活动规则说明
	SortOrder     int    `json:"sortOrder" gorm:"column:sort_order;default:0"`              // 排序（越小越靠前）
	IsVisible     bool   `json:"isVisible" gorm:"column:is_visible;default:true"`           // 是否在前端展示

	// 关联
	Rewards []ActivityReward `json:"rewards,omitempty" gorm:"foreignKey:ActivityID"`
}

// TableName 指定表名
func (Activity) TableName() string {
	return "activities"
}

// IsInPreheat 是否在预热期
func (a *Activity) IsInPreheat() bool {
	now := time.Now()
	if a.PreheatAt == nil {
		return false
	}
	return now.After(*a.PreheatAt) && now.Before(a.StartAt)
}

// IsActive 是否在活动期间
func (a *Activity) IsActive() bool {
	now := time.Now()
	return a.Status == ActivityStatusActive && now.After(a.StartAt) && now.Before(a.EndAt)
}

// ActivityReward 活动奖励配置
type ActivityReward struct {
	Base
	ActivityID       uint64 `json:"activityId" gorm:"column:activity_id;not null;index"`        // 活动ID
	CouponTemplateID uint64 `json:"couponTemplateId" gorm:"column:coupon_template_id;not null"` // 优惠券模板ID
	CouponCount      int    `json:"couponCount" gorm:"column:coupon_count;default:1"`           // 每次发放数量
	Probability      int    `json:"probability" gorm:"default:100"`                             // 发放概率（1-100，预留抽奖用）
	TotalStock       int    `json:"totalStock" gorm:"column:total_stock;default:0"`             // 总库存（0=无限制）
	RemainingStock   int    `json:"remainingStock" gorm:"column:remaining_stock;default:0"`     // 剩余库存
	SortOrder        int    `json:"sortOrder" gorm:"column:sort_order;default:0"`               // 排序

	// Relations
	Activity       *Activity       `json:"activity,omitempty" gorm:"foreignKey:ActivityID"`
	CouponTemplate *CouponTemplate `json:"couponTemplate,omitempty" gorm:"foreignKey:CouponTemplateID"`
}

// TableName 指定表名
func (ActivityReward) TableName() string {
	return "activity_rewards"
}

// ActivityParticipation 活动参与记录
type ActivityParticipation struct {
	Base
	ActivityID uint64    `json:"activityId" gorm:"column:activity_id;not null;index"`   // 活动ID
	UserID     uint64    `json:"userId" gorm:"column:user_id;not null;index"`           // 用户ID
	RewardID   uint64    `json:"rewardId" gorm:"column:reward_id;not null"`             // 奖励配置ID
	CouponIDs  string    `json:"couponIds,omitempty" gorm:"column:coupon_ids;size:512"` // 发放的优惠券ID列表（JSON数组）
	ClaimedAt  time.Time `json:"claimedAt" gorm:"column:claimed_at;not null"`           // 领取时间
	ClientIP   string    `json:"clientIp,omitempty" gorm:"column:client_ip;size:64"`    // 客户端IP

	// Relations
	Activity *Activity       `json:"activity,omitempty" gorm:"foreignKey:ActivityID"`
	User     *User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Reward   *ActivityReward `json:"reward,omitempty" gorm:"foreignKey:RewardID"`
}

// TableName 指定表名
func (ActivityParticipation) TableName() string {
	return "activity_participations"
}

// ActivityDailyStats 活动每日统计（用于每日限制重置）
type ActivityDailyStats struct {
	Base
	ActivityID   uint64    `json:"activityId" gorm:"column:activity_id;not null;uniqueIndex:idx_activity_date"` // 活动ID
	StatsDate    time.Time `json:"statsDate" gorm:"column:stats_date;not null;uniqueIndex:idx_activity_date"`   // 统计日期
	Participants int       `json:"participants" gorm:"default:0"`                                               // 当日参与人数
	ClaimCount   int       `json:"claimCount" gorm:"column:claim_count;default:0"`                              // 当日领取次数
}

// TableName 指定表名
func (ActivityDailyStats) TableName() string {
	return "activity_daily_stats"
}
