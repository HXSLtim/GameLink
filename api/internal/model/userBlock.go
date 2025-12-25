package model

import "time"

// ========== 用户拉黑系统 ==========

// BlockUserType 拉黑用户类型
type BlockUserType string

const (
	// BlockUserTypeUser 普通用户
	BlockUserTypeUser BlockUserType = "user"
	// BlockUserTypePlayer 陪玩师
	BlockUserTypePlayer BlockUserType = "player"
)

// BlockStatus 拉黑状态
type BlockStatus string

const (
	// BlockStatusActive 生效中
	BlockStatusActive BlockStatus = "active"
	// BlockStatusCanceled 已取消（用户主动）
	BlockStatusCanceled BlockStatus = "canceled"
	// BlockStatusAdminCanceled 已取消（管理员强制）
	BlockStatusAdminCanceled BlockStatus = "admin_canceled"
)

// UserBlock 用户拉黑记录
// @Description 用户/陪玩师之间的拉黑关系，支持双向拉黑
type UserBlock struct {
	Base
	// 拉黑发起人ID（User.ID）
	BlockerID uint64 `json:"blockerId" gorm:"column:blocker_id;index;not null"`
	// 发起人类型：user/player
	BlockerType BlockUserType `json:"blockerType" gorm:"column:blocker_type;size:32;index;not null"`
	// 被拉黑人ID（User.ID）
	BlockedID uint64 `json:"blockedId" gorm:"column:blocked_id;index;not null"`
	// 被拉黑人类型：user/player
	BlockedType BlockUserType `json:"blockedType" gorm:"column:blocked_type;size:32;index;not null"`
	// 拉黑原因（可选）
	Reason string `json:"reason,omitempty" gorm:"size:500"`
	// 状态：active/canceled/admin_canceled
	Status BlockStatus `json:"status" gorm:"size:32;index;default:'active'"`
	// 拉黑时间
	BlockedAt time.Time `json:"blockedAt" gorm:"column:blocked_at;not null"`
	// 取消时间
	CanceledAt *time.Time `json:"canceledAt,omitempty" gorm:"column:canceled_at"`
	// 取消人ID（管理员强制解除时记录）
	CanceledBy *uint64 `json:"canceledBy,omitempty" gorm:"column:canceled_by"`
	// 管理员备注
	AdminRemark string `json:"adminRemark,omitempty" gorm:"column:admin_remark;size:500"`

	// 关联
	Blocker  *User `json:"blocker,omitempty" gorm:"foreignKey:BlockerID"`
	Blocked  *User `json:"blocked,omitempty" gorm:"foreignKey:BlockedID"`
	Canceler *User `json:"canceler,omitempty" gorm:"foreignKey:CanceledBy"`
}

// TableName 指定表名
func (UserBlock) TableName() string {
	return "user_blocks"
}

// IsActive 检查拉黑是否生效中
func (b *UserBlock) IsActive() bool {
	return b.Status == BlockStatusActive
}

// Cancel 取消拉黑（用户主动）
func (b *UserBlock) Cancel() {
	now := time.Now()
	b.Status = BlockStatusCanceled
	b.CanceledAt = &now
}

// AdminCancel 管理员强制取消拉黑
func (b *UserBlock) AdminCancel(adminID uint64, remark string) {
	now := time.Now()
	b.Status = BlockStatusAdminCanceled
	b.CanceledAt = &now
	b.CanceledBy = &adminID
	b.AdminRemark = remark
}
