package model

// OrderItemStatus 订单明细状态
type OrderItemStatus string

const (
	OrderItemStatusPending   OrderItemStatus = "pending"   // 待匹配
	OrderItemStatusMatched   OrderItemStatus = "matched"   // 已匹配陪玩师
	OrderItemStatusCompleted OrderItemStatus = "completed" // 已完成
	OrderItemStatusCanceled  OrderItemStatus = "canceled"  // 已取消
)

// OrderItem 订单明细（一个订单可包含多个服务项目/座位）
type OrderItem struct {
	Base
	OrderID        uint64          `json:"orderId" gorm:"column:order_id;not null;index"`                  // 主订单ID
	ItemID         uint64          `json:"itemId" gorm:"column:item_id;not null;index"`                    // 服务项目ID
	Slot           int             `json:"slot" gorm:"column:slot;not null"`                               // 座位号（1, 2, 3...）
	UnitPriceCents int64           `json:"unitPriceCents" gorm:"column:unit_price_cents;not null"`         // 单价（分）
	Quantity       int             `json:"quantity" gorm:"column:quantity;default:1"`                      // 数量
	TotalCents     int64           `json:"totalCents" gorm:"column:total_cents;not null"`                  // 小计（分）
	CommissionRate float64         `json:"commissionRate" gorm:"column:commission_rate;type:decimal(5,2)"` // 抽成比例
	Status         OrderItemStatus `json:"status" gorm:"column:status;size:32;default:'pending';index"`    // pending/matched/completed/canceled

	// 接单后填充
	PlayerID *uint64 `json:"playerId,omitempty" gorm:"column:player_id;index"` // 接单的陪玩师ID

	// 评价关联
	ReviewID *uint64 `json:"reviewId,omitempty" gorm:"column:review_id;index"` // 关联的评价ID（nil=未评价）

	// Relations
	Order       *Order       `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	ServiceItem *ServiceItem `json:"serviceItem,omitempty" gorm:"foreignKey:ItemID"`
	Player      *Player      `json:"player,omitempty" gorm:"foreignKey:PlayerID"`
}

// TableName 指定表名
func (OrderItem) TableName() string {
	return "order_items"
}

// IsMatched 是否已匹配陪玩师
func (oi *OrderItem) IsMatched() bool {
	return oi.PlayerID != nil && *oi.PlayerID > 0
}
