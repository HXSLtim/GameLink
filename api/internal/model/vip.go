package model

// VipLevel VIP等级配置（完全可配置）
type VipLevel struct {
	Base
	Slug        string `json:"slug" gorm:"size:32;uniqueIndex"`                           // 等级标识（可配置，如 vip1/vip2/svip1）
	Title       string `json:"title" gorm:"size:64;not null"`                             // 等级名称（显示用）
	ExpRequired int64  `json:"expRequired" gorm:"column:exp_required;not null;default:0"` // 升级所需累计消费/经验（分）

	// 永久折扣
	OrderDiscount float64 `json:"orderDiscount" gorm:"column:order_discount;default:1.0"` // 下单永久折扣 (0.98 = 98折, 1.0 = 无折扣)

	// 月度券配置（绑定优惠券模板）
	MonthlyCouponTemplateID *uint64 `json:"monthlyCouponTemplateId" gorm:"column:monthly_coupon_template_id;index"` // 月度券模板ID
	MonthlyCouponCount      int     `json:"monthlyCouponCount" gorm:"column:monthly_coupon_count;default:0"`        // 每月发放数量

	// 展示配置
	IconURL   string `json:"iconUrl" gorm:"column:icon_url;size:255"`      // 等级图标
	Color     string `json:"color" gorm:"size:32"`                         // 等级颜色（前端展示）
	Benefits  string `json:"benefits" gorm:"type:json;default:'{}'"`       // 其他权益描述(JSON)
	SortOrder int    `json:"sortOrder" gorm:"column:sort_order;default:0"` // 排序（越小越靠前）

	// 状态
	IsDefault bool `json:"isDefault" gorm:"column:is_default;default:false"` // 是否默认等级（新用户解锁后）
	IsActive  bool `json:"isActive" gorm:"column:is_active;default:true"`    // 是否启用
}

// TableName 指定表名
func (VipLevel) TableName() string {
	return "vip_levels"
}

// HasOrderDiscount 是否有永久折扣
func (v *VipLevel) HasOrderDiscount() bool {
	return v.OrderDiscount > 0 && v.OrderDiscount < 1.0
}

// HasMonthlyCoupon 是否有月度券
func (v *VipLevel) HasMonthlyCoupon() bool {
	return v.MonthlyCouponTemplateID != nil && v.MonthlyCouponCount > 0
}

// VipConfig VIP系统配置（全局配置）
type VipConfig struct {
	Base
	ConfigKey   string `json:"configKey" gorm:"column:config_key;size:64;uniqueIndex"` // 配置键
	ConfigValue string `json:"configValue" gorm:"column:config_value;type:text"`       // 配置值
	Description string `json:"description" gorm:"size:255"`                            // 描述
}

// TableName 指定表名
func (VipConfig) TableName() string {
	return "vip_configs"
}

// VIP配置键常量
const (
	VipConfigUnlockByConsume  = "unlock_by_consume"  // 累计消费解锁门槛（分）
	VipConfigUnlockByRecharge = "unlock_by_recharge" // 累计充值解锁门槛（分）
	VipConfigExpireDays       = "expire_days"        // VIP过期天数（0=永久）
)
