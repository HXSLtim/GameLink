package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// DisputeStatus defines the lifecycle states for a dispute.
// @Enum pending, assigned, mediating, resolved, rejected, canceled
type DisputeStatus string

// DisputeStatus values define the lifecycle of a dispute.
const (
	DisputeStatusPending   DisputeStatus = "pending"   // 待处理
	DisputeStatusAssigned  DisputeStatus = "assigned"  // 已指派
	DisputeStatusMediating DisputeStatus = "mediating" // 调解中
	DisputeStatusResolved  DisputeStatus = "resolved"  // 已解决
	DisputeStatusRejected  DisputeStatus = "rejected"  // 已驳回
	DisputeStatusCanceled  DisputeStatus = "canceled"  // 已取消
)

// DisputeResolution defines the resolution decision for a dispute.
type DisputeResolution string

// DisputeResolution values define possible resolution outcomes.
const (
	ResolutionRefund   DisputeResolution = "refund"   // 全额退款
	ResolutionPartial  DisputeResolution = "partial"  // 部分退款
	ResolutionReassign DisputeResolution = "reassign" // 重新指派
	ResolutionReject   DisputeResolution = "reject"   // 驳回
	ResolutionPending  DisputeResolution = "pending"  // 待决定
)

// DisputeInitiatorType 争议发起人类型
type DisputeInitiatorType string

const (
	DisputeInitiatorUser   DisputeInitiatorType = "user"   // 用户发起
	DisputeInitiatorPlayer DisputeInitiatorType = "player" // 陪玩师发起
)

// DisputeType 争议类型
type DisputeType string

const (
	DisputeTypeServiceQuality     DisputeType = "service_quality"      // 服务质量问题
	DisputeTypeBadAttitude        DisputeType = "bad_attitude"         // 态度问题
	DisputeTypeIncompleteService  DisputeType = "incomplete_service"   // 未完成服务
	DisputeTypeUserNotCooperative DisputeType = "user_not_cooperative" // 用户不配合/不听指挥
	DisputeTypeUserHarassment     DisputeType = "user_harassment"      // 用户骚扰
	DisputeTypeOther              DisputeType = "other"                // 其他
)

// AssignmentSource defines the source of an assignment.
type AssignmentSource string

// AssignmentSource values define where an assignment came from.
const (
	AssignmentSourceSystem AssignmentSource = "system" // 系统推荐
	AssignmentSourceManual AssignmentSource = "manual" // 人工指定
	AssignmentSourceTeam   AssignmentSource = "team"   // 车队分配
)

// OrderDispute represents a customer service dispute for an order.
type OrderDispute struct {
	Base
	OrderID       uint64               `json:"orderId" gorm:"column:order_id;not null;index"`               // 订单ID
	InitiatorID   uint64               `json:"initiatorId" gorm:"column:initiator_id;not null;index"`       // 发起人ID
	InitiatorType DisputeInitiatorType `json:"initiatorType" gorm:"column:initiator_type;size:32;not null"` // 发起人类型：user/player
	Type          DisputeType          `json:"type" gorm:"column:type;size:32;not null"`                    // 争议类型
	Status        DisputeStatus        `json:"status" gorm:"column:status;size:32;index;default:'pending'"` // 争议状态
	Reason        string               `json:"reason" gorm:"column:reason;type:text;not null"`              // 争议原因（用户填写）
	EvidenceURLs  EvidenceURLArray     `json:"evidenceUrls" gorm:"column:evidence_urls;type:json"`          // 证据截图URL列表（最多5张）
	EvidenceText  string               `json:"evidenceText" gorm:"column:evidence_text;type:text"`          // 文字证据说明

	// 聊天记录快照
	ChatSnapshotID *uint64 `json:"chatSnapshotId,omitempty" gorm:"column:chat_snapshot_id;index"` // 聊天记录快照ID

	// 双客服机制
	OriginalServiceID *uint64 `json:"originalServiceId,omitempty" gorm:"column:original_service_id;index"` // 原客服ID（订单原有的客服）
	AssignedServiceID *uint64 `json:"assignedServiceId,omitempty" gorm:"column:assigned_service_id;index"` // 分配的独立客服ID（保证公正）

	// SLA 信息
	SLADeadline   *time.Time `json:"slaDeadline" gorm:"column:sla_deadline;index"`         // SLA 截止时间（默认30分钟）
	SLABreached   bool       `json:"slaBreached" gorm:"column:sla_breached;default:false"` // 是否超过SLA
	SLABreachedAt *time.Time `json:"slaBreachedAt" gorm:"column:sla_breached_at"`          // 超过SLA的时间

	// 处理信息
	Resolution    DisputeResolution `json:"resolution" gorm:"column:resolution;size:32;default:'pending'"` // 处理决定：refund/reject
	ResolvedBy    *uint64           `json:"resolvedBy,omitempty" gorm:"column:resolved_by"`                // 处理人ID
	ResolvedAt    *time.Time        `json:"resolvedAt,omitempty" gorm:"column:resolved_at"`                // 处理时间
	ResolveRemark string            `json:"resolveRemark" gorm:"column:resolve_remark;type:text"`          // 处理备注

	// 回退信息（保留向后兼容）
	RolledBackAt       *time.Time `json:"rolledBackAt,omitempty" gorm:"column:rolled_back_at"`               // 回退时间
	RolledBackByUserID *uint64    `json:"rolledBackByUserId,omitempty" gorm:"column:rolled_back_by_user_id"` // 回退人ID
	RollbackReason     string     `json:"rollbackReason,omitempty" gorm:"column:rollback_reason;type:text"`  // 回退原因

	// 追踪信息
	TraceID string `json:"traceId" gorm:"column:trace_id;size:64;index"` // 追踪ID

	// Relations
	Order            Order         `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:OrderID;references:ID"`
	Initiator        User          `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:InitiatorID;references:ID"`
	ChatSnapshot     *ChatSnapshot `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:ChatSnapshotID;references:ID"`
	OriginalService  *User         `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:OriginalServiceID;references:ID"`
	AssignedService  *User         `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:AssignedServiceID;references:ID"`
	ResolvedByUser   *User         `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:ResolvedBy;references:ID"`
	RolledBackByUser *User         `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:RolledBackByUserID;references:ID"`
}

// EvidenceURLArray is a custom type for storing evidence URLs as JSON array.
type EvidenceURLArray []string

// Scan implements the sql.Scanner interface.
func (e *EvidenceURLArray) Scan(value any) error {
	if value == nil {
		*e = EvidenceURLArray{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, e)
	case string:
		return json.Unmarshal([]byte(v), e)
	default:
		return errors.New("unsupported type for EvidenceURLArray")
	}
}

// Value implements the driver.Valuer interface.
func (e EvidenceURLArray) Value() (driver.Value, error) {
	return json.Marshal(e)
}

// TableName specifies the table name for OrderDispute.
func (OrderDispute) TableName() string {
	return "order_disputes"
}

// IsOverSLA checks if the dispute has exceeded the SLA deadline.
func (d *OrderDispute) IsOverSLA() bool {
	if d.SLADeadline == nil {
		return false
	}
	return time.Now().After(*d.SLADeadline)
}

// GetSLARemaining returns the remaining time until SLA deadline in seconds.
func (d *OrderDispute) GetSLARemaining() int64 {
	if d.SLADeadline == nil {
		return 0
	}
	remaining := d.SLADeadline.Unix() - time.Now().Unix()
	if remaining < 0 {
		return 0
	}
	return remaining
}

// CanInitiateDispute checks if a dispute can be initiated for the given order.
// Disputes can only be initiated during service or within 7 days of order completion.
func CanInitiateDispute(order *Order) bool {
	if order.CompletedAt == nil {
		// Can also initiate during service (in_progress status)
		return order.Status == OrderStatusInProgress
	}
	// 7 days after completion (售后期)
	return time.Since(*order.CompletedAt) <= 7*24*time.Hour
}

// ========== 争议类型模板 ==========

// DisputeTemplate 争议类型模板
// @Description 预设的争议类型，用户选择后再填写具体原因
type DisputeTemplate struct {
	Base
	Code          string               `json:"code" gorm:"size:64;uniqueIndex;not null"`                 // 模板编码（唯一）
	Name          string               `json:"name" gorm:"size:128;not null"`                            // 模板名称
	InitiatorType DisputeInitiatorType `json:"initiatorType" gorm:"column:initiator_type;size:32;index"` // 适用发起人：user/player
	Description   string               `json:"description" gorm:"size:500"`                              // 描述说明
	SortOrder     int                  `json:"sortOrder" gorm:"column:sort_order;default:0"`             // 排序
	IsActive      bool                 `json:"isActive" gorm:"column:is_active;default:true"`            // 是否启用
}

// TableName specifies the table name for DisputeTemplate.
func (DisputeTemplate) TableName() string {
	return "dispute_templates"
}

// ========== 聊天记录快照 ==========

// ChatSnapshot 聊天记录快照
// @Description 争议发起时自动保存的聊天记录快照
type ChatSnapshot struct {
	Base
	DisputeID   uint64    `json:"disputeId" gorm:"column:dispute_id;index;not null"` // 争议ID
	OrderID     uint64    `json:"orderId" gorm:"column:order_id;index;not null"`     // 订单ID
	ChatGroupID uint64    `json:"chatGroupId" gorm:"column:chat_group_id;not null"`  // 聊天群组ID
	Messages    string    `json:"messages" gorm:"type:text;not null"`                // 聊天记录（JSON数组）
	SnapshotAt  time.Time `json:"snapshotAt" gorm:"column:snapshot_at;not null"`     // 快照时间
}

// TableName specifies the table name for ChatSnapshot.
func (ChatSnapshot) TableName() string {
	return "chat_snapshots"
}
