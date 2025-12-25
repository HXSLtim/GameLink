package model

import "time"

// ========== 订单超时配置 ==========

// OrderTimeoutConfig 订单超时配置
// @Description 订单超时相关的系统配置
type OrderTimeoutConfig struct {
	Base
	// 配置键（唯一）
	ConfigKey string `json:"configKey" gorm:"column:config_key;size:64;uniqueIndex;not null"`
	// 配置值
	ConfigValue string `json:"configValue" gorm:"column:config_value;size:255;not null"`
	// 描述
	Description string `json:"description,omitempty" gorm:"size:255"`
}

// TableName 指定表名
func (OrderTimeoutConfig) TableName() string {
	return "order_timeout_configs"
}

// 订单超时配置键常量
const (
	// PaymentTimeoutMinutes 支付超时时间（分钟），默认30
	PaymentTimeoutMinutes = "payment_timeout_minutes"
	// OrderAcceptTimeoutMinutes 接单超时时间（分钟），默认30
	OrderAcceptTimeoutMinutes = "order_accept_timeout_minutes"
	// AutoCancelEnabled 是否启用自动取消，默认true
	AutoCancelEnabled = "auto_cancel_enabled"
	// AutoRefundEnabled 是否启用自动退款，默认true
	AutoRefundEnabled = "auto_refund_enabled"
	// AutoAssignServiceEnabled 接单后是否自动分配客服，默认true
	AutoAssignServiceEnabled = "auto_assign_service_enabled"
)

// ========== 订单超时日志 ==========

// OrderTimeoutType 订单超时类型
type OrderTimeoutType string

const (
	// OrderTimeoutTypePayment 支付超时
	OrderTimeoutTypePayment OrderTimeoutType = "payment_timeout"
	// OrderTimeoutTypeAccept 接单超时
	OrderTimeoutTypeAccept OrderTimeoutType = "accept_timeout"
)

// OrderTimeoutAction 订单超时处理动作
type OrderTimeoutAction string

const (
	// OrderTimeoutActionCanceled 已取消
	OrderTimeoutActionCanceled OrderTimeoutAction = "canceled"
	// OrderTimeoutActionRefunded 已退款
	OrderTimeoutActionRefunded OrderTimeoutAction = "refunded"
	// OrderTimeoutActionNotified 已通知
	OrderTimeoutActionNotified OrderTimeoutAction = "notified"
)

// OrderTimeoutLog 订单超时日志
// @Description 记录订单超时处理的日志
type OrderTimeoutLog struct {
	Base
	// 订单ID
	OrderID uint64 `json:"orderId" gorm:"column:order_id;index;not null"`
	// 超时类型
	TimeoutType OrderTimeoutType `json:"timeoutType" gorm:"column:timeout_type;size:32;index;not null"`
	// 超时时间
	TimeoutAt time.Time `json:"timeoutAt" gorm:"column:timeout_at;not null"`
	// 处理动作
	Action OrderTimeoutAction `json:"action" gorm:"size:32;not null"`
	// 退款金额（分）
	RefundAmountCents int64 `json:"refundAmountCents,omitempty" gorm:"column:refund_amount_cents;default:0"`
	// 退款记录ID
	RefundRecordID *uint64 `json:"refundRecordId,omitempty" gorm:"column:refund_record_id"`
	// 备注
	Remark string `json:"remark,omitempty" gorm:"size:500"`

	// 关联
	Order *Order `json:"order,omitempty" gorm:"foreignKey:OrderID"`
}

// TableName 指定表名
func (OrderTimeoutLog) TableName() string {
	return "order_timeout_logs"
}

// ========== 客服分配记录 ==========

// ServiceAssignmentStatus 客服分配状态
type ServiceAssignmentStatus string

const (
	// ServiceAssignmentStatusAssigned 已分配
	ServiceAssignmentStatusAssigned ServiceAssignmentStatus = "assigned"
	// ServiceAssignmentStatusJoined 已加入房间
	ServiceAssignmentStatusJoined ServiceAssignmentStatus = "joined"
	// ServiceAssignmentStatusLeft 已离开
	ServiceAssignmentStatusLeft ServiceAssignmentStatus = "left"
	// ServiceAssignmentStatusCompleted 已完成
	ServiceAssignmentStatusCompleted ServiceAssignmentStatus = "completed"
)

// OrderServiceAssignment 订单客服分配记录
// @Description 记录订单接单后自动分配的客服
type OrderServiceAssignment struct {
	Base
	// 订单ID
	OrderID uint64 `json:"orderId" gorm:"column:order_id;index;not null"`
	// 客服用户ID
	ServiceUserID uint64 `json:"serviceUserId" gorm:"column:service_user_id;index;not null"`
	// 聊天群组ID
	ChatGroupID *uint64 `json:"chatGroupId,omitempty" gorm:"column:chat_group_id;index"`
	// 分配状态
	Status ServiceAssignmentStatus `json:"status" gorm:"size:32;index;default:'assigned'"`
	// 分配时间
	AssignedAt time.Time `json:"assignedAt" gorm:"column:assigned_at;not null"`
	// 加入房间时间
	JoinedAt *time.Time `json:"joinedAt,omitempty" gorm:"column:joined_at"`
	// 离开时间
	LeftAt *time.Time `json:"leftAt,omitempty" gorm:"column:left_at"`
	// 分配方式：auto（自动）/ manual（手动）
	AssignType string `json:"assignType" gorm:"column:assign_type;size:32;default:'auto'"`
	// 备注
	Remark string `json:"remark,omitempty" gorm:"size:500"`

	// 关联
	Order       *Order     `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	ServiceUser *User      `json:"serviceUser,omitempty" gorm:"foreignKey:ServiceUserID"`
	ChatGroup   *ChatGroup `json:"chatGroup,omitempty" gorm:"foreignKey:ChatGroupID"`
}

// TableName 指定表名
func (OrderServiceAssignment) TableName() string {
	return "order_service_assignments"
}
