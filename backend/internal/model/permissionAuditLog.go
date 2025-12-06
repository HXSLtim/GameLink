package model

import "time"

// AuditAction defines the type of audit action performed.
type AuditAction string

// AuditAction values define the types of audit actions.
const (
	AuditActionCreate        AuditAction = "create"
	AuditActionUpdate        AuditAction = "update"
	AuditActionDelete        AuditAction = "delete"
	AuditActionAssign        AuditAction = "assign"
	AuditActionRevoke        AuditAction = "revoke"
	AuditActionBatchAssign   AuditAction = "batch_assign"
	AuditActionBatchRevoke   AuditAction = "batch_revoke"
	AuditActionInheritChange AuditAction = "inherit_change"
)

// AuditTargetType defines the type of target being audited.
type AuditTargetType string

// AuditTargetType values define the types of audit targets.
const (
	AuditTargetTypeRole       AuditTargetType = "role"
	AuditTargetTypeUser       AuditTargetType = "user"
	AuditTargetTypePermission AuditTargetType = "permission"
)

// AuditLogRetentionDays defines how many days audit logs are kept online.
const AuditLogRetentionDays = 90

// AuditLogArchiveDays defines how many days audit logs are kept in archive.
const AuditLogArchiveDays = 365

// PermissionAuditLog records all permission-related changes for audit purposes.
// It captures the operator, target, action, and before/after data snapshots.
type PermissionAuditLog struct {
	ID           uint64          `json:"id" gorm:"primaryKey"`
	CreatedAt    time.Time       `json:"createdAt" gorm:"column:created_at;index"`
	OperatorID   uint64          `json:"operatorId" gorm:"index;not null;comment:操作者ID"`
	OperatorName string          `json:"operatorName" gorm:"size:128;comment:操作者名称"`
	TargetType   AuditTargetType `json:"targetType" gorm:"size:32;index;comment:目标类型(role/user/permission)"`
	TargetID     uint64          `json:"targetId" gorm:"index;comment:目标ID"`
	TargetName   string          `json:"targetName" gorm:"size:128;comment:目标名称"`
	Action       AuditAction     `json:"action" gorm:"size:32;index;comment:操作类型"`
	BeforeData   string          `json:"beforeData" gorm:"type:text;comment:操作前数据快照(JSON)"`
	AfterData    string          `json:"afterData" gorm:"type:text;comment:操作后数据快照(JSON)"`
	IPAddress    string          `json:"ipAddress" gorm:"size:64;comment:操作IP"`
	UserAgent    string          `json:"userAgent" gorm:"size:512;comment:用户代理"`
	RequestID    string          `json:"requestId" gorm:"size:64;index;comment:请求追踪ID"`
}

// TableName specifies the table name for PermissionAuditLog.
func (PermissionAuditLog) TableName() string {
	return "permission_audit_logs"
}
