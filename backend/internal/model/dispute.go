package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// DisputeStatus defines the lifecycle states for a dispute.
type DisputeStatus string

// DisputeStatus values define the lifecycle of a dispute.
const (
	DisputeStatusPending    DisputeStatus = "pending"     // 待处理
	DisputeStatusAssigned   DisputeStatus = "assigned"    // 已指派
	DisputeStatusMediating  DisputeStatus = "mediating"   // 调解中
	DisputeStatusResolved   DisputeStatus = "resolved"    // 已解决
	DisputeStatusRejected   DisputeStatus = "rejected"    // 已驳回
	DisputeStatusCanceled   DisputeStatus = "canceled"    // 已取消
)

// DisputeResolution defines the resolution decision for a dispute.
type DisputeResolution string

// DisputeResolution values define possible resolution outcomes.
const (
	ResolutionRefund    DisputeResolution = "refund"     // 全额退款
	ResolutionPartial   DisputeResolution = "partial"    // 部分退款
	ResolutionReassign  DisputeResolution = "reassign"   // 重新指派
	ResolutionReject    DisputeResolution = "reject"     // 驳回
	ResolutionPending   DisputeResolution = "pending"    // 待决定
)

// AssignmentSource defines the source of an assignment.
type AssignmentSource string

// AssignmentSource values define where an assignment came from.
const (
	AssignmentSourceSystem AssignmentSource = "system"   // 系统推荐
	AssignmentSourceManual AssignmentSource = "manual"   // 人工指定
	AssignmentSourceTeam   AssignmentSource = "team"     // 车队分配
)

// OrderDispute represents a customer service dispute for an order.
type OrderDispute struct {
	Base
	OrderID              uint64            `json:"orderId" gorm:"column:order_id;not null;index"`                    // 订单ID
	UserID               uint64            `json:"userId" gorm:"column:user_id;not null;index"`                      // 发起用户ID
	Status               DisputeStatus     `json:"status" gorm:"column:status;size:32;index;default:'pending'"`      // 争议状态
	Reason               string            `json:"reason" gorm:"column:reason;type:text;not null"`                   // 争议原因
	Description          string            `json:"description" gorm:"column:description;type:text"`                  // 详细描述
	EvidenceURLs         EvidenceURLArray  `json:"evidenceUrls" gorm:"column:evidence_urls;type:json"`               // 证据截图URL列表
	
	// 指派信息
	AssignedToUserID     *uint64           `json:"assignedToUserId" gorm:"column:assigned_to_user_id;index"`         // 指派给的客服ID
	AssignmentSource     AssignmentSource  `json:"assignmentSource" gorm:"column:assignment_source;size:32"`         // 指派来源
	AssignedAt           *time.Time        `json:"assignedAt" gorm:"column:assigned_at"`                             // 指派时间
	
	// SLA 信息
	SLADeadline          *time.Time        `json:"slaDeadline" gorm:"column:sla_deadline;index"`                     // SLA 截止时间（默认30分钟）
	SLABreached          bool              `json:"slaBreached" gorm:"column:sla_breached;default:false"`             // 是否超过SLA
	SLABreachedAt        *time.Time        `json:"slaBreachedAt" gorm:"column:sla_breached_at"`                      // 超过SLA的时间
	
	// 处理信息
	Resolution           DisputeResolution `json:"resolution" gorm:"column:resolution;size:32;default:'pending'"`   // 处理决定
	ResolutionAmount     int64             `json:"resolutionAmount" gorm:"column:resolution_amount;default:0"`       // 退款金额（分）
	ResolutionNotes      string            `json:"resolutionNotes" gorm:"column:resolution_notes;type:text"`         // 处理备注
	ResolvedAt           *time.Time        `json:"resolvedAt" gorm:"column:resolved_at"`                             // 解决时间
	ResolvedByUserID     *uint64           `json:"resolvedByUserId" gorm:"column:resolved_by_user_id"`               // 处理人ID
	
	// 回退信息
	RolledBackAt         *time.Time        `json:"rolledBackAt" gorm:"column:rolled_back_at"`                        // 回退时间
	RolledBackByUserID   *uint64           `json:"rolledBackByUserId" gorm:"column:rolled_back_by_user_id"`          // 回退人ID
	RollbackReason       string            `json:"rollbackReason" gorm:"column:rollback_reason;type:text"`           // 回退原因
	
	// 追踪信息
	TraceID              string            `json:"traceId" gorm:"column:trace_id;size:64;index"`                     // 追踪ID
	
	// Relations
	Order                Order             `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:OrderID;references:ID"`
	User                 User              `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:UserID;references:ID"`
	AssignedToUser       *User             `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:AssignedToUserID;references:ID"`
	ResolvedByUser       *User             `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:ResolvedByUserID;references:ID"`
	RolledBackByUser     *User             `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:RolledBackByUserID;references:ID"`
}

// EvidenceURLArray is a custom type for storing evidence URLs as JSON array.
type EvidenceURLArray []string

// Scan implements the sql.Scanner interface.
func (e *EvidenceURLArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion failed")
	}
	return json.Unmarshal(bytes, &e)
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
// Disputes can only be initiated within 24 hours of order completion.
func CanInitiateDispute(order *Order) bool {
	if order.CompletedAt == nil {
		// Can also initiate during service (in_progress status)
		return order.Status == OrderStatusInProgress
	}
	// 24 hours after completion
	return time.Since(*order.CompletedAt) <= 24*time.Hour
}
