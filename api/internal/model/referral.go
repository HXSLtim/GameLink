package model

import "time"

// ReferralType 推荐类型
type ReferralType string

const (
	ReferralTypeUserToUser     ReferralType = "user_to_user"     // 用户邀请用户
	ReferralTypePlayerToPlayer ReferralType = "player_to_player" // 陪玩师邀请陪玩师
	ReferralTypeUserToPlayer   ReferralType = "user_to_player"   // 用户转陪玩师
)

// ReferralStatus 推荐状态
type ReferralStatus string

const (
	ReferralStatusPending   ReferralStatus = "pending"   // 待完成（被邀请人未注册/未满足条件）
	ReferralStatusCompleted ReferralStatus = "completed" // 已完成
	ReferralStatusRewarded  ReferralStatus = "rewarded"  // 已发放奖励
	ReferralStatusExpired   ReferralStatus = "expired"   // 已过期
	ReferralStatusCanceled  ReferralStatus = "canceled"  // 已取消
)

// RewardType 奖励类型
type RewardType string

const (
	RewardTypeCash   RewardType = "cash"   // 现金（到钱包）
	RewardTypeCoupon RewardType = "coupon" // 优惠券
	RewardTypePoints RewardType = "points" // 积分（预留）
)

// ReferralConfig 推荐配置（全局配置）
type ReferralConfig struct {
	Base
	ConfigKey   string `json:"configKey" gorm:"column:config_key;size:64;uniqueIndex"` // 配置键
	ConfigValue string `json:"configValue" gorm:"column:config_value;type:text"`       // 配置值
	Description string `json:"description" gorm:"size:255"`                            // 描述
}

// TableName 指定表名
func (ReferralConfig) TableName() string {
	return "referral_configs"
}

// 推荐配置键常量
const (
	ReferralConfigEnabled            = "enabled"              // 是否启用推荐系统
	ReferralConfigExpireDays         = "expire_days"          // 邀请链接过期天数
	ReferralConfigMaxLevel           = "max_level"            // 最大推荐层级（1=直推，2=二级分销）
	ReferralConfigUserRewardType     = "user_reward_type"     // 用户邀请奖励类型
	ReferralConfigUserRewardAmount   = "user_reward_amount"   // 用户邀请奖励金额（分）
	ReferralConfigPlayerRewardType   = "player_reward_type"   // 陪玩师邀请奖励类型
	ReferralConfigPlayerRewardAmount = "player_reward_amount" // 陪玩师邀请奖励金额（分）
)

// ReferralCode 邀请码
type ReferralCode struct {
	Base
	Code     string       `json:"code" gorm:"size:32;uniqueIndex"`               // 邀请码
	UserID   uint64       `json:"userId" gorm:"column:user_id;not null;index"`   // 所属用户ID
	Type     ReferralType `json:"type" gorm:"size:32;index"`                     // 推荐类型
	IsActive bool         `json:"isActive" gorm:"column:is_active;default:true"` // 是否启用
	ExpireAt *time.Time   `json:"expireAt,omitempty" gorm:"column:expire_at"`    // 过期时间（nil=永久）
	UseCount int          `json:"useCount" gorm:"column:use_count;default:0"`    // 使用次数
	MaxUse   int          `json:"maxUse" gorm:"column:max_use;default:0"`        // 最大使用次数（0=无限制）

	// Relations
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (ReferralCode) TableName() string {
	return "referral_codes"
}

// IsValid 检查邀请码是否有效
func (rc *ReferralCode) IsValid() bool {
	if !rc.IsActive {
		return false
	}
	if rc.ExpireAt != nil && time.Now().After(*rc.ExpireAt) {
		return false
	}
	if rc.MaxUse > 0 && rc.UseCount >= rc.MaxUse {
		return false
	}
	return true
}

// Referral 推荐记录
type Referral struct {
	Base
	ReferrerID  uint64         `json:"referrerId" gorm:"column:referrer_id;not null;index"` // 推荐人ID
	RefereeID   uint64         `json:"refereeId" gorm:"column:referee_id;not null;index"`   // 被推荐人ID
	CodeID      *uint64        `json:"codeId,omitempty" gorm:"column:code_id;index"`        // 使用的邀请码ID
	Type        ReferralType   `json:"type" gorm:"size:32;index"`                           // 推荐类型
	Level       int            `json:"level" gorm:"default:1"`                              // 推荐层级（1=直推，2=二级）
	Status      ReferralStatus `json:"status" gorm:"size:32;default:'pending';index"`       // 状态
	CompletedAt *time.Time     `json:"completedAt,omitempty" gorm:"column:completed_at"`    // 完成时间

	// 奖励信息
	RewardType        RewardType `json:"rewardType" gorm:"column:reward_type;size:32"`                  // 奖励类型
	RewardAmountCents int64      `json:"rewardAmountCents" gorm:"column:reward_amount_cents;default:0"` // 奖励金额（分）
	RewardedAt        *time.Time `json:"rewardedAt,omitempty" gorm:"column:rewarded_at"`                // 奖励发放时间
	RewardNote        string     `json:"rewardNote,omitempty" gorm:"column:reward_note;size:255"`       // 奖励备注

	// 被推荐人条件（用于判断是否完成）
	RefereeCondition string `json:"refereeCondition,omitempty" gorm:"column:referee_condition;size:64"` // 完成条件：registered/first_order/first_recharge

	// Relations
	Referrer *User         `json:"referrer,omitempty" gorm:"foreignKey:ReferrerID"`
	Referee  *User         `json:"referee,omitempty" gorm:"foreignKey:RefereeID"`
	Code     *ReferralCode `json:"code,omitempty" gorm:"foreignKey:CodeID"`
}

// TableName 指定表名
func (Referral) TableName() string {
	return "referrals"
}

// ReferralReward 推荐奖励记录
type ReferralReward struct {
	Base
	ReferralID    uint64     `json:"referralId" gorm:"column:referral_id;not null;index"`           // 推荐记录ID
	UserID        uint64     `json:"userId" gorm:"column:user_id;not null;index"`                   // 获得奖励的用户ID
	Type          RewardType `json:"type" gorm:"size:32"`                                           // 奖励类型
	AmountCents   int64      `json:"amountCents" gorm:"column:amount_cents;default:0"`              // 奖励金额（分）
	CouponID      *uint64    `json:"couponId,omitempty" gorm:"column:coupon_id;index"`              // 发放的优惠券ID
	Status        string     `json:"status" gorm:"size:32;default:'pending'"`                       // pending/issued/failed
	IssuedAt      *time.Time `json:"issuedAt,omitempty" gorm:"column:issued_at"`                    // 发放时间
	FailureReason string     `json:"failureReason,omitempty" gorm:"column:failure_reason;size:255"` // 失败原因

	// Relations
	Referral *Referral `json:"referral,omitempty" gorm:"foreignKey:ReferralID"`
	User     *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Coupon   *Coupon   `json:"coupon,omitempty" gorm:"foreignKey:CouponID"`
}

// TableName 指定表名
func (ReferralReward) TableName() string {
	return "referral_rewards"
}
