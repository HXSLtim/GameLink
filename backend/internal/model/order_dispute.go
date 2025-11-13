package model

import "time"

// OrderDisputeRaisedBy 表示争议发起方。
type OrderDisputeRaisedBy string

// OrderDisputeResolution 表示争议最终裁决。
type OrderDisputeResolution string

const (
	OrderDisputeRaisedByUser   OrderDisputeRaisedBy = "user"
	OrderDisputeRaisedByPlayer OrderDisputeRaisedBy = "player"
	OrderDisputeRaisedBySystem OrderDisputeRaisedBy = "system"
)

const (
	OrderDisputeResolutionNone     OrderDisputeResolution = "none"
	OrderDisputeResolutionRefund   OrderDisputeResolution = "refund"
	OrderDisputeResolutionReassign OrderDisputeResolution = "reassign"
	OrderDisputeResolutionReject   OrderDisputeResolution = "reject"
)

// OrderDispute 记录订单争议及调解信息。
type OrderDispute struct {
	Base
	OrderID           uint64                 `json:"orderId" gorm:"column:order_id;not null;index"`
	RaisedBy          OrderDisputeRaisedBy   `json:"raisedBy" gorm:"column:raised_by;size:16;not null"`
	RaisedByUserID    *uint64                `json:"raisedByUserId" gorm:"column:raised_by_user_id;index"`
	Reason            string                 `json:"reason" gorm:"type:text"`
	EvidenceURLs      []string               `json:"evidenceUrls" gorm:"column:evidence_urls;serializer:json"`
	Status            OrderDisputeStatus     `json:"status" gorm:"column:status;size:32;index;default:'pending'"`
	Resolution        OrderDisputeResolution `json:"resolution" gorm:"column:resolution;size:32;default:'none'"`
	ResolutionNote    string                 `json:"resolutionNote" gorm:"column:resolution_note;type:text"`
	RefundAmountCents int64                  `json:"refundAmountCents" gorm:"column:refund_amount_cents;default:0"`
	HandledByID       *uint64                `json:"handledById" gorm:"column:handled_by_id;index"`
	HandledAt         *time.Time             `json:"handledAt" gorm:"column:handled_at"`
	ResponseDeadline  time.Time              `json:"responseDeadline" gorm:"column:response_deadline"`
	RespondedAt       *time.Time             `json:"respondedAt" gorm:"column:responded_at"`
	TraceID           string                 `json:"traceId" gorm:"column:trace_id;size:64"`
}

// TableName 自定义表名。
func (OrderDispute) TableName() string {
	return "order_disputes"
}
