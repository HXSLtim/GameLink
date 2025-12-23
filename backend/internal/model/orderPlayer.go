package model

import "time"

// OrderPlayer 订单陪玩师关联表（记录陪玩师接单详情）
type OrderPlayer struct {
	Base
	OrderID     uint64    `json:"orderId" gorm:"column:order_id;not null;index"`                                      // 订单ID
	OrderItemID uint64    `json:"orderItemId" gorm:"column:order_item_id;not null;index;uniqueIndex:idx_item_player"` // 订单明细ID
	PlayerID    uint64    `json:"playerId" gorm:"column:player_id;not null;index;uniqueIndex:idx_item_player"`        // 陪玩师ID
	JoinedAt    time.Time `json:"joinedAt" gorm:"column:joined_at;not null"`                                          // 接单时间
	TeamID      *uint64   `json:"teamId,omitempty" gorm:"column:team_id;index"`                                       // 团队ID（团队接单时）

	// 收入分配
	IncomeCents     int64 `json:"incomeCents" gorm:"column:income_cents;default:0"`         // 该陪玩师收入（分）
	CommissionCents int64 `json:"commissionCents" gorm:"column:commission_cents;default:0"` // 该陪玩师抽成（分）

	// 状态
	Status string `json:"status" gorm:"column:status;size:32;default:'joined'"` // joined/left/completed

	// Relations
	Order     *Order     `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	OrderItem *OrderItem `json:"orderItem,omitempty" gorm:"foreignKey:OrderItemID"`
	Player    *Player    `json:"player,omitempty" gorm:"foreignKey:PlayerID"`
}

// TableName 指定表名
func (OrderPlayer) TableName() string {
	return "order_players"
}

// OrderPlayerStatus 状态常量
const (
	OrderPlayerStatusJoined    = "joined"    // 已加入
	OrderPlayerStatusLeft      = "left"      // 已离开（中途退出）
	OrderPlayerStatusCompleted = "completed" // 已完成
)
