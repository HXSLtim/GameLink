package model

import "time"

// ReviewReportStatus 表示举报的处理状态
type ReviewReportStatus string

const (
	ReviewReportStatusPending  ReviewReportStatus = "pending"  // 待处理
	ReviewReportStatusApproved ReviewReportStatus = "approved" // 已通过（删除评价）
	ReviewReportStatusRejected ReviewReportStatus = "rejected" // 已驳回
)

// Valid 检查举报状态是否合法
func (rrs ReviewReportStatus) Valid() bool {
	switch rrs {
	case ReviewReportStatusPending, ReviewReportStatusApproved, ReviewReportStatusRejected:
		return true
	default:
		return false
	}
}

// ReviewReport 评价举报模型
// @Description 评价举报记录，用于处理用户对不当评价的投诉
type ReviewReport struct {
	Base
	// 被举报的评价ID
	// @Example 1001
	ReviewID uint64 `json:"reviewId" gorm:"column:review_id;index;not null"`
	// 举报人ID
	// @Example 2001
	ReporterID uint64 `json:"reporterId" gorm:"column:reporter_id;index;not null"`
	// 举报原因
	// @Example 评价内容包含不实信息和恶意诋毁
	Reason string `json:"reason" gorm:"column:reason;type:text;not null"`
	// 举报证据（可选，如截图URL等）
	// @Example https://example.com/evidence.jpg
	Evidence string `json:"evidence,omitempty" gorm:"column:evidence;type:text"`
	// 处理状态
	// @Enum pending, approved, rejected
	// @Example pending
	Status ReviewReportStatus `json:"status" gorm:"column:status;type:varchar(20);default:'pending';index"`
	// 处理人ID（管理员ID）
	// @Example 3001
	HandledBy *uint64 `json:"handledBy,omitempty" gorm:"column:handled_by;index"`
	// 处理时间
	// @Example 2024-01-15T14:30:00Z
	HandledAt *time.Time `json:"handledAt,omitempty" gorm:"column:handled_at"`
	// 处理备注
	// @Example 经核实，评价内容确实存在不实信息，已删除该评价
	HandlingNote string `json:"handlingNote,omitempty" gorm:"column:handling_note;type:text"`
}
