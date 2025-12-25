package model

// OrderItem 订单明细（一个订单可包含多个服务项目/座位）
type OrderItem struct {
	Base
	OrderID        uint64  `json:"orderId" gorm:"column:order_id;not null;index"`                  // 主订单ID
	ItemID         uint64  `json:"itemId" gorm:"column:item_id;not null;index"`                    // 服务项目ID
	Slot           int     `json:"slot" gorm:"column:slot;not null"`                               // 座位号（1, 2, 3...）
	UnitPriceCents int64   `json:"unitPriceCents" gorm:"column:unit_price_cents;not null"`         // 单价（分）
	Quantity       int     `json:"quantity" gorm:"column:quantity;default:1"`                      // 数量
	TotalCents     int64   `json:"totalCents" gorm:"column:total_cents;not null"`                  // 小计（分）
	CommissionRate float64 `json:"commissionRate" gorm:"column:commission_rate;type:decimal(5,2)"` // 抽成比例
	Status         string  `json:"status" gorm:"column:status;size:32;default:'pending'"`          // pending/matched/completed

	// 接单后填充
	PlayerID *uint64 `json:"playerId,omitempty" gorm:"column:player_id;index"` // 接单的陪玩师ID

	// Relations
	Order       *Order       `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	ServiceItem *ServiceItem `json:"serviceItem,omitempty" gorm:"foreignKey:ItemID"`
	Player      *Player      `json:"player,omitempty" gorm:"foreignKey:PlayerID"`
}

// TableName 指定表名
func (OrderItem) TableName() string {
	return "order_items"
}

// OrderItemStatus 状态常量
const (
	OrderItemStatusPending   = "pending"   // 待匹配
	OrderItemStatusMatched   = "matched"   // 已匹配陪玩师
	OrderItemStatusCompleted = "completed" // 已完成
	OrderItemStatusCanceled  = "canceled"  // 已取消
)

// IsMatched 是否已匹配陪玩师
func (oi *OrderItem) IsMatched() bool {
	return oi.PlayerID != nil && *oi.PlayerID > 0
}
