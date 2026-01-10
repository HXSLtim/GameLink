package model

import "time"

// CouponType 优惠券类型
type CouponType string

const (
	CouponTypeDeduct   CouponType = "deduct"   // 满减券（含无门槛，MinAmount=0）
	CouponTypeDiscount CouponType = "discount" // 折扣券
)

// CouponScope 优惠券适用范围
type CouponScope string

const (
	CouponScopeAll  CouponScope = "all"  // 通用
	CouponScopeGame CouponScope = "game" // 指定游戏
	CouponScopeItem CouponScope = "item" // 指定服务项目
)

// CouponSource 优惠券来源
type CouponSource string

const (
	CouponSourceNewUser  CouponSource = "new_user" // 新用户注册发放
	CouponSourceLink     CouponSource = "link"     // 链接领取
	CouponSourceVip      CouponSource = "vip"      // VIP月度发放
	CouponSourceRecharge CouponSource = "recharge" // 充值赠送
	CouponSourceActivity CouponSource = "activity" // 活动发放
	CouponSourceManual   CouponSource = "manual"   // 手动发放
)

// CouponState 优惠券状态
type CouponState string

const (
	CouponStateAvailable CouponState = "available" // 可用
	CouponStateLocked    CouponState = "locked"    // 已锁定（下单中）
	CouponStateUsed      CouponState = "used"      // 已使用
	CouponStateExpired   CouponState = "expired"   // 已过期
	CouponStateDeleted   CouponState = "deleted"   // 已删除（用户侧）
)

// CouponTemplate 优惠券配置/模板
type CouponTemplate struct {
	Base
	Name        string       `json:"name" gorm:"size:128;not null"`        // 券名称
	Type        CouponType   `json:"type" gorm:"size:32;not null"`         // 类型
	Source      CouponSource `json:"source" gorm:"size:32;not null;index"` // 来源
	Description string       `json:"description" gorm:"size:255"`          // 描述

	// 优惠配置
	MinAmountCents    int64   `json:"minAmountCents" gorm:"column:min_amount_cents;default:0"`       // 最低消费门槛（分），0=无门槛
	DeductAmountCents int64   `json:"deductAmountCents" gorm:"column:deduct_amount_cents;default:0"` // 满减金额（分）- 满减券用
	DiscountRate      float64 `json:"discountRate" gorm:"column:discount_rate;default:1.0"`          // 折扣率（0.9=9折）- 折扣券用
	MaxDiscountCents  int64   `json:"maxDiscountCents" gorm:"column:max_discount_cents;default:0"`   // 最大折扣金额（分）- 折扣券用，0=无上限

	// 适用范围
	Scope   CouponScope `json:"scope" gorm:"size:32;default:'all'"`                    // 适用范围
	GameIDs string      `json:"gameIds" gorm:"column:game_ids;type:json;default:'[]'"` // 指定游戏ID列表（JSON数组）
	ItemIDs string      `json:"itemIds" gorm:"column:item_ids;type:json;default:'[]'"` // 指定服务项目ID列表（JSON数组）

	// 有效期配置
	ValidityType  string     `json:"validityType" gorm:"column:validity_type;size:32;default:'days'"` // days=固定天数, fixed=固定截止日期
	ValidityDays  int        `json:"validityDays" gorm:"column:validity_days;default:30"`             // 有效天数（领取后）
	FixedExpireAt *time.Time `json:"fixedExpireAt" gorm:"column:fixed_expire_at"`                     // 固定截止日期

	// 领取配置
	TotalCount   int    `json:"totalCount" gorm:"column:total_count;default:0"`          // 总发放数量（0=无限制）
	ClaimedCount int    `json:"claimedCount" gorm:"column:claimed_count;default:0"`      // 已领取数量
	PerUserLimit int    `json:"perUserLimit" gorm:"column:per_user_limit;default:1"`     // 每人限领数量
	ClaimLink    string `json:"claimLink" gorm:"column:claim_link;size:128;uniqueIndex"` // 领取链接码（链接领取用）

	// 使用频率（预留）
	UsageFrequency string `json:"usageFrequency" gorm:"column:usage_frequency;size:32;default:'once'"` // once/daily/weekly/monthly

	// 状态
	IsActive bool `json:"isActive" gorm:"column:is_active;default:true"` // 是否启用
}

// TableName 指定表名
func (CouponTemplate) TableName() string {
	return "coupon_templates"
}

// Coupon 用户优惠券
type Coupon struct {
	Base
	TemplateID uint64      `json:"templateId" gorm:"column:template_id;not null;index"` // 模板ID
	UserID     uint64      `json:"userId" gorm:"column:user_id;not null;index"`         // 用户ID
	State      CouponState `json:"state" gorm:"size:32;default:'available';index"`      // 状态

	// 冗余字段（方便查询，从模板复制）
	Name              string       `json:"name" gorm:"size:128"`                                  // 券名称
	Type              CouponType   `json:"type" gorm:"size:32"`                                   // 类型
	Source            CouponSource `json:"source" gorm:"size:32"`                                 // 来源
	MinAmountCents    int64        `json:"minAmountCents" gorm:"column:min_amount_cents"`         // 最低消费门槛
	DeductAmountCents int64        `json:"deductAmountCents" gorm:"column:deduct_amount_cents"`   // 满减金额
	DiscountRate      float64      `json:"discountRate" gorm:"column:discount_rate"`              // 折扣率
	MaxDiscountCents  int64        `json:"maxDiscountCents" gorm:"column:max_discount_cents"`     // 最大折扣金额
	Scope             CouponScope  `json:"scope" gorm:"size:32"`                                  // 适用范围
	GameIDs           string       `json:"gameIds" gorm:"column:game_ids;type:json;default:'[]'"` // 指定游戏ID
	ItemIDs           string       `json:"itemIds" gorm:"column:item_ids;type:json;default:'[]'"` // 指定服务项目ID

	// 时间
	ClaimedAt *time.Time `json:"claimedAt" gorm:"column:claimed_at"`     // 领取时间
	ExpireAt  time.Time  `json:"expireAt" gorm:"column:expire_at;index"` // 过期时间
	UsedAt    *time.Time `json:"usedAt" gorm:"column:used_at"`           // 使用时间

	// 锁定信息
	LockedByOrderID *uint64    `json:"lockedByOrderId" gorm:"column:locked_by_order_id;index"` // 锁定的订单ID
	LockedAt        *time.Time `json:"lockedAt" gorm:"column:locked_at"`                       // 锁定时间

	// 使用信息
	UsedOrderID   *uint64 `json:"usedOrderId" gorm:"column:used_order_id;index"`        // 使用的订单ID
	DiscountCents int64   `json:"discountCents" gorm:"column:discount_cents;default:0"` // 实际折扣金额（分）

	// Relations
	Template *CouponTemplate `json:"template,omitempty" gorm:"foreignKey:TemplateID"`
	User     *User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (Coupon) TableName() string {
	return "coupons"
}

// IsValid 检查优惠券是否可用
func (c *Coupon) IsValid() bool {
	return c.State == CouponStateAvailable && time.Now().Before(c.ExpireAt)
}

// CalculateDiscount 计算折扣金额
func (c *Coupon) CalculateDiscount(orderAmountCents int64) int64 {
	// 检查门槛
	if orderAmountCents < c.MinAmountCents {
		return 0
	}

	if c.Type == CouponTypeDeduct {
		// 满减券
		if c.DeductAmountCents > orderAmountCents {
			return orderAmountCents
		}
		return c.DeductAmountCents
	}

	// 折扣券
	discount := int64(float64(orderAmountCents) * (1 - c.DiscountRate))
	if c.MaxDiscountCents > 0 && discount > c.MaxDiscountCents {
		return c.MaxDiscountCents
	}
	return discount
}
