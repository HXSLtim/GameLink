package model

import "time"

// OrderStatus defines lifecycle states for an order.
type OrderStatus string

// OrderStatus values define the lifecycle of an order.
const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusInProgress OrderStatus = "in_progress"
	OrderStatusCompleted  OrderStatus = "completed"
	OrderStatusCanceled   OrderStatus = "canceled"
	OrderStatusRefunded   OrderStatus = "refunded"
)

// Order represents a unified order (护航服务 or 礼物)
type Order struct {
	Base
	OrderNo           string  `json:"orderNo" gorm:"column:order_no;size:64;uniqueIndex:idx_order_no,where:order_no != ''"` // 订单号
	UserID            uint64  `json:"userId" gorm:"column:user_id;not null;index"`                                          // 下单用户
	ItemID            uint64  `json:"itemId" gorm:"column:item_id;not null;index"`                                          // 服务项目ID
	PlayerID          *uint64 `json:"playerId,omitempty" gorm:"column:player_id;index"`                                     // 服务陪玩师（礼物订单也可选）
	RecipientPlayerID *uint64 `json:"recipientPlayerId,omitempty" gorm:"column:recipient_player_id;index"`                  // 接收礼物的陪玩师

	// 价格相关
	Quantity          int      `json:"quantity" gorm:"default:1"`                                     // 数量
	UnitPriceCents    int64    `json:"unitPriceCents" gorm:"column:unit_price_cents;not null"`        // 单价（分）
	TotalPriceCents   int64    `json:"totalPriceCents" gorm:"column:total_price_cents;not null"`      // 总价（分）
	CommissionCents   int64    `json:"commissionCents" gorm:"column:commission_cents;default:0"`      // 平台抽成（分）
	PlayerIncomeCents int64    `json:"playerIncomeCents" gorm:"column:player_income_cents;default:0"` // 陪玩师收入（分）
	Currency          Currency `json:"currency,omitempty" gorm:"type:char(3);default:'CNY'"`          // 货币

	// 订单信息
	Status      OrderStatus `json:"status" gorm:"size:32;index;default:'pending'"` // 订单状态
	Title       string      `json:"title,omitempty" gorm:"size:128"`               // 订单标题
	Description string      `json:"description,omitempty" gorm:"type:text"`        // 订单描述

	// 护航服务字段
	GameID         *uint64    `json:"gameId,omitempty" gorm:"column:game_id;index"`           // 游戏ID
	ScheduledStart *time.Time `json:"scheduledStart,omitempty" gorm:"column:scheduled_start"` // 预约开始时间
	ScheduledEnd   *time.Time `json:"scheduledEnd,omitempty" gorm:"column:scheduled_end"`     // 预约结束时间
	StartedAt      *time.Time `json:"startedAt,omitempty" gorm:"column:started_at"`           // 实际开始时间
	CompletedAt    *time.Time `json:"completedAt,omitempty" gorm:"column:completed_at"`       // 完成时间

	// 礼物订单字段
	GiftMessage string     `json:"giftMessage,omitempty" gorm:"column:gift_message;type:text"` // 礼物留言
	IsAnonymous bool       `json:"isAnonymous" gorm:"column:is_anonymous;default:false"`       // 是否匿名
	DeliveredAt *time.Time `json:"deliveredAt,omitempty" gorm:"column:delivered_at"`           // 礼物送达时间

	// 取消/退款信息
	CancelReason      string     `json:"cancelReason,omitempty" gorm:"column:cancel_reason;type:text"`
	RefundAmountCents int64      `json:"refundAmountCents,omitempty" gorm:"column:refund_amount_cents;default:0"`
	RefundReason      string     `json:"refundReason,omitempty" gorm:"column:refund_reason;type:text"`
	RefundedAt        *time.Time `json:"refundedAt,omitempty" gorm:"column:refunded_at"`

	// 扩展字段
	OrderConfig string `json:"orderConfig,omitempty" gorm:"column:order_config;type:json"` // 订单配置（JSON）
	UserNotes   string `json:"userNotes,omitempty" gorm:"column:user_notes;type:text"`     // 用户备注

	// 争议相关字段
	HasDispute  bool   `json:"hasDispute" gorm:"column:has_dispute;default:false;index"`   // 是否有争议
	DisputeID   *uint64 `json:"disputeId,omitempty" gorm:"column:dispute_id;index"`        // 关联的争议ID

	// Relations
	User            User         `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:UserID;references:ID"`
	Player          *Player      `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:PlayerID;references:ID"`
	RecipientPlayer *Player      `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:RecipientPlayerID;references:ID"`
	Game            *Game        `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:GameID;references:ID"`
	ServiceItem     *ServiceItem `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:ItemID;references:ID"`
	Dispute         *OrderDispute `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:DisputeID;references:ID"`
}

// IsGiftOrder 判断是否为礼物订单
func (o *Order) IsGiftOrder() bool {
	return o.RecipientPlayerID != nil && *o.RecipientPlayerID > 0
}

// GetPlayerID 获取陪玩师ID（向后兼容）
func (o *Order) GetPlayerID() uint64 {
	if o.PlayerID != nil {
		return *o.PlayerID
	}
	return 0
}

// GetGameID 获取游戏ID（向后兼容）
func (o *Order) GetGameID() uint64 {
	if o.GameID != nil {
		return *o.GameID
	}
	return 0
}

// GetPriceCents 获取价格（向后兼容）- 返回总价
func (o *Order) GetPriceCents() int64 {
	return o.TotalPriceCents
}

// SetPlayerID 设置陪玩师ID（向后兼容）
func (o *Order) SetPlayerID(playerID uint64) {
	o.PlayerID = &playerID
}

// SetGameID 设置游戏ID（向后兼容）
func (o *Order) SetGameID(gameID uint64) {
	o.GameID = &gameID
}
