package model

import "encoding/json"

// OperationLog 记录后台关键业务操作日志，用于审计与可视化。
type OperationLog struct {
	Base
	EntityType   string          `json:"entityType" gorm:"column:entity_type;size:32;index"` // order | payment
	EntityID     uint64          `json:"entityId" gorm:"column:entity_id;index"`
	ActorUserID  *uint64         `json:"actorUserId" gorm:"column:actor_user_id;index"`
	Action       string          `json:"action" gorm:"size:64;index"`
	Reason       string          `json:"reason,omitempty" gorm:"type:text"`
	MetadataJSON json.RawMessage `json:"metadata,omitempty" gorm:"column:metadata_json;type:json"`
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

	// 争议处理
	OpActionInitiateDispute  OperationAction = "initiate_dispute"
	OpActionAssignDispute    OperationAction = "assign_dispute"
	OpActionMediateDispute   OperationAction = "mediate_dispute"
	OpActionResolveDispute   OperationAction = "resolve_dispute"
	OpActionRollbackDispute  OperationAction = "rollback_dispute"
	OpActionRejectDispute    OperationAction = "reject_dispute"
)

// OperationEntityType 枚举被审计的实体类型。
type OperationEntityType string

const (
	OpEntityOrder    OperationEntityType = "order"
	OpEntityPayment  OperationEntityType = "payment"
	OpEntityPlayer   OperationEntityType = "player"
	OpEntityGame     OperationEntityType = "game"
	OpEntityReview   OperationEntityType = "review"
	OpEntityUser     OperationEntityType = "user"
	OpEntityDispute  OperationEntityType = "dispute"
)
