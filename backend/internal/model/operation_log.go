package model

import "encoding/json"

// OperationLog 记录后台关键业务操作日志，用于审计与可视化。
type OperationLog struct {
	Base
	EntityType   string          `json:"entity_type" gorm:"size:32;index"` // order | payment
	EntityID     uint64          `json:"entity_id" gorm:"index"`
	ActorUserID  *uint64         `json:"actor_user_id" gorm:"index"`
	Action       string          `json:"action" gorm:"size:64;index"`
	Reason       string          `json:"reason,omitempty" gorm:"type:text"`
	MetadataJSON json.RawMessage `json:"metadata,omitempty" gorm:"type:json"`
}

// OperationAction 枚举所有标准化的审计动作。
type OperationAction string

const (
	// 通用
	OpActionCreate       OperationAction = "create"
	OpActionUpdateStatus OperationAction = "update_status"
	OpActionDelete       OperationAction = "delete"

	// 订单
	OpActionAssignPlayer OperationAction = "assign_player"
	OpActionCancel       OperationAction = "cancel"
	OpActionConfirm      OperationAction = "confirm"
	OpActionStart        OperationAction = "start"
	OpActionComplete     OperationAction = "complete"

	// 支付
	OpActionCapture    OperationAction = "capture"
	OpActionRefund     OperationAction = "refund"
	OpActionUpdate     OperationAction = "update"
	OpActionUpdateRole OperationAction = "update_role"
)

// OperationEntityType 枚举被审计的实体类型。
type OperationEntityType string

const (
	OpEntityOrder   OperationEntityType = "order"
	OpEntityPayment OperationEntityType = "payment"
	OpEntityPlayer  OperationEntityType = "player"
	OpEntityGame    OperationEntityType = "game"
	OpEntityReview  OperationEntityType = "review"
	OpEntityUser    OperationEntityType = "user"
)
