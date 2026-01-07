package model

import "time"

// OrderGroupStatus 主订单状态
type OrderGroupStatus string

const (
	OrderGroupStatusPending    OrderGroupStatus = "pending"     // 待支付
	OrderGroupStatusPaid       OrderGroupStatus = "paid"        // 已支付，待接单
	OrderGroupStatusInProgress OrderGroupStatus = "in_progress" // 进行中
	OrderGroupStatusCompleted  OrderGroupStatus = "completed"   // 已完成
	OrderGroupStatusCanceled   OrderGroupStatus = "canceled"    // 已取消
	OrderGroupStatusRefunded   OrderGroupStatus = "refunded"    // 已退款
	OrderGroupStatusPartial    OrderGroupStatus = "partial"     // 部分完成（有子订单转单）
)

// OrderGroup 主订单（用户视角的订单）
// 用户下单 N 小时，系统拆分成 N 个子订单，但用户只看到这一个主订单
type OrderGroup struct {
	Base
	GroupNo         string           `json:"groupNo" gorm:"column:group_no;size:64;uniqueIndex"` // 主订单号（用户看到的）
	UserID          uint64           `json:"userId" gorm:"column:user_id;not null;index"`        // 下单用户
	GameID          uint64           `json:"gameId" gorm:"column:game_id;not null;index"`        // 游戏ID
	ItemID          uint64           `json:"itemId" gorm:"column:item_id;not null;index"`        // 服务项目ID
	OriginalPlayer  uint64           `json:"originalPlayer" gorm:"column:original_player_id"`    // 原始陪玩师ID
	TotalPriceCents int64            `json:"totalPriceCents" gorm:"column:total_price_cents"`    // 总价（分）
	TotalHours      int              `json:"totalHours" gorm:"column:total_hours"`               // 总时长（小时）
	CompletedHours  int              `json:"completedHours" gorm:"column:completed_hours"`       // 已完成时长
	Status          OrderGroupStatus `json:"status" gorm:"size:32;index;default:'pending'"`      // 主订单状态
	Title           string           `json:"title" gorm:"size:128"`                              // 订单标题
	Description     string           `json:"description" gorm:"type:text"`                       // 订单描述
	ScheduledStart  *time.Time       `json:"scheduledStart" gorm:"column:scheduled_start"`       // 预约开始时间
	ScheduledEnd    *time.Time       `json:"scheduledEnd" gorm:"column:scheduled_end"`           // 预约结束时间
	Currency        Currency         `json:"currency" gorm:"type:char(3);default:'CNY'"`         // 货币
	UserNotes       string           `json:"userNotes" gorm:"column:user_notes;type:text"`       // 用户备注

	// 取消/退款
	CancelReason      string     `json:"cancelReason,omitempty" gorm:"column:cancel_reason;type:text"`
	RefundAmountCents int64      `json:"refundAmountCents" gorm:"column:refund_amount_cents;default:0"`
	RefundedAt        *time.Time `json:"refundedAt,omitempty" gorm:"column:refunded_at"`

	// Relations
	User        User         `json:"-" gorm:"foreignKey:UserID"`
	Game        *Game        `json:"-" gorm:"foreignKey:GameID"`
	ServiceItem *ServiceItem `json:"-" gorm:"foreignKey:ItemID"`
	SubOrders   []Order      `json:"subOrders,omitempty" gorm:"foreignKey:GroupID"` // 子订单列表
}

// TableName 指定表名
func (OrderGroup) TableName() string {
	return "order_groups"
}

// GenerateGroupOrderNo 生成主订单号
func GenerateGroupOrderNo() string {
	return GenerateOrderNo("G") // G 开头表示 Group
}

// UpdateStatusFromSubOrders 根据子订单状态更新主订单状态
func (g *OrderGroup) UpdateStatusFromSubOrders(subOrders []Order) {
	if len(subOrders) == 0 {
		return
	}

	completedCount := 0
	inProgressCount := 0
	canceledCount := 0
	pendingCount := 0

	for _, sub := range subOrders {
		switch sub.Status {
		case OrderStatusCompleted:
			completedCount++
		case OrderStatusInProgress:
			inProgressCount++
		case OrderStatusCanceled, OrderStatusRefunded:
			canceledCount++
		case OrderStatusPending, OrderStatusConfirmed:
			pendingCount++
		}
	}

	g.CompletedHours = completedCount

	// 状态判断逻辑
	if completedCount == len(subOrders) {
		g.Status = OrderGroupStatusCompleted
	} else if canceledCount == len(subOrders) {
		g.Status = OrderGroupStatusCanceled
	} else if inProgressCount > 0 {
		g.Status = OrderGroupStatusInProgress
	} else if completedCount > 0 && (pendingCount > 0 || inProgressCount > 0) {
		g.Status = OrderGroupStatusPartial
	} else if pendingCount == len(subOrders) {
		g.Status = OrderGroupStatusPending
	}
}
