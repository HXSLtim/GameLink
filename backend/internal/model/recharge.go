package model

import "time"

// RechargeStatus 充值状态
type RechargeStatus string

const (
	RechargeStatusPending  RechargeStatus = "pending"  // 待支付
	RechargeStatusPaid     RechargeStatus = "paid"     // 已支付
	RechargeStatusFailed   RechargeStatus = "failed"   // 支付失败
	RechargeStatusRefunded RechargeStatus = "refunded" // 已退款
	RechargeStatusCanceled RechargeStatus = "canceled" // 已取消
)

// RechargeOption 充值档位配置
type RechargeOption struct {
	Base
	Name            string `json:"name" gorm:"size:64;not null"`                             // 档位名称
	AmountCents     int64  `json:"amountCents" gorm:"column:amount_cents;not null"`          // 充值金额（分）
	BonusCents      int64  `json:"bonusCents" gorm:"column:bonus_cents;default:0"`           // 赠送金额（分）
	TotalCents      int64  `json:"totalCents" gorm:"column:total_cents;not null"`            // 实际到账（分）= AmountCents + BonusCents
	OriginalCents   *int64 `json:"originalCents" gorm:"column:original_cents"`               // 原价（分），用于显示划线价
	DiscountPercent *int   `json:"discountPercent" gorm:"column:discount_percent"`           // 折扣百分比（显示用，如 80 表示 8折）
	Description     string `json:"description" gorm:"size:255"`                              // 描述
	Tag             string `json:"tag" gorm:"size:32"`                                       // 标签（如"推荐"、"热门"）
	IconURL         string `json:"iconUrl" gorm:"column:icon_url;size:255"`                  // 图标URL
	SortOrder       int    `json:"sortOrder" gorm:"column:sort_order;default:0"`             // 排序（越小越靠前）
	IsActive        bool   `json:"isActive" gorm:"column:is_active;default:true"`            // 是否启用
	IsRecommended   bool   `json:"isRecommended" gorm:"column:is_recommended;default:false"` // 是否推荐

	// 优惠券赠送配置
	CouponTemplateID *uint64 `json:"couponTemplateId" gorm:"column:coupon_template_id;index"` // 赠送优惠券模板ID
	CouponCount      int     `json:"couponCount" gorm:"column:coupon_count;default:0"`        // 赠送优惠券数量

	// 限制配置
	MinVipLevel   *uint64 `json:"minVipLevel" gorm:"column:min_vip_level"`              // 最低VIP等级要求（nil=无限制）
	PerUserLimit  int     `json:"perUserLimit" gorm:"column:per_user_limit;default:0"`  // 每人限购次数（0=无限制）
	TotalLimit    int     `json:"totalLimit" gorm:"column:total_limit;default:0"`       // 总限购次数（0=无限制）
	PurchaseCount int     `json:"purchaseCount" gorm:"column:purchase_count;default:0"` // 已购买次数

	// Relations
	CouponTemplate *CouponTemplate `json:"couponTemplate,omitempty" gorm:"foreignKey:CouponTemplateID"`
}

// TableName 指定表名
func (RechargeOption) TableName() string {
	return "recharge_options"
}

// RechargeRecord 充值记录
type RechargeRecord struct {
	Base
	UserID      uint64         `json:"userId" gorm:"column:user_id;not null;index"`     // 用户ID
	OptionID    *uint64        `json:"optionId" gorm:"column:option_id;index"`          // 档位ID（自定义金额时为nil）
	AmountCents int64          `json:"amountCents" gorm:"column:amount_cents;not null"` // 充值金额（分）
	BonusCents  int64          `json:"bonusCents" gorm:"column:bonus_cents;default:0"`  // 赠送金额（分）
	TotalCents  int64          `json:"totalCents" gorm:"column:total_cents;not null"`   // 实际到账（分）
	Status      RechargeStatus `json:"status" gorm:"size:32;default:'pending';index"`   // 状态

	// 订单号（退款用）
	OrderNo         string `json:"orderNo" gorm:"column:order_no;size:64;uniqueIndex"`            // 内部订单号
	MerchantOrderNo string `json:"merchantOrderNo" gorm:"column:merchant_order_no;size:64;index"` // 商户订单号（提交给支付渠道）
	ProviderTradeNo string `json:"providerTradeNo" gorm:"column:provider_trade_no;size:64;index"` // 第三方交易号（支付渠道返回）

	// 收款分流信息
	MerchantID       string `json:"merchantId" gorm:"column:merchant_id;size:64;index"`         // 商户号
	CollectionEntity string `json:"collectionEntity" gorm:"column:collection_entity;size:64"`   // 收款主体
	PaymentChannel   string `json:"paymentChannel" gorm:"column:payment_channel;size:32;index"` // 支付渠道：wechat/alipay

	// 支付信息
	PaymentMethod string     `json:"paymentMethod" gorm:"column:payment_method;size:32"` // 支付方式：wechat_h5/wechat_mini/alipay_h5 等
	PaidAt        *time.Time `json:"paidAt" gorm:"column:paid_at"`                       // 支付时间
	ExpireAt      *time.Time `json:"expireAt" gorm:"column:expire_at"`                   // 过期时间

	// 退款信息
	RefundedAt        *time.Time `json:"refundedAt" gorm:"column:refunded_at"`                          // 退款时间
	RefundAmountCents int64      `json:"refundAmountCents" gorm:"column:refund_amount_cents;default:0"` // 退款金额（分）
	RefundReason      string     `json:"refundReason" gorm:"column:refund_reason;size:255"`             // 退款原因
	RefundProviderNo  string     `json:"refundProviderNo" gorm:"column:refund_provider_no;size:64"`     // 退款第三方单号

	// 优惠券发放记录
	CouponIssued bool   `json:"couponIssued" gorm:"column:coupon_issued;default:false"` // 优惠券是否已发放
	CouponIDs    string `json:"couponIds" gorm:"column:coupon_ids;type:json"`           // 发放的优惠券ID列表（JSON数组）

	// 客户端信息
	ClientIP   string `json:"clientIp" gorm:"column:client_ip;size:64"`       // 客户端IP
	UserAgent  string `json:"userAgent" gorm:"column:user_agent;size:255"`    // User-Agent
	DeviceInfo string `json:"deviceInfo" gorm:"column:device_info;type:json"` // 设备信息（JSON）

	// 备注
	Remark string `json:"remark" gorm:"size:255"` // 备注

	// Relations
	User   *User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Option *RechargeOption `json:"option,omitempty" gorm:"foreignKey:OptionID"`
}

// TableName 指定表名
func (RechargeRecord) TableName() string {
	return "recharge_records"
}

// IsPaid 是否已支付
func (r *RechargeRecord) IsPaid() bool {
	return r.Status == RechargeStatusPaid
}

// CanRefund 是否可退款
func (r *RechargeRecord) CanRefund() bool {
	return r.Status == RechargeStatusPaid && r.RefundAmountCents < r.AmountCents
}
