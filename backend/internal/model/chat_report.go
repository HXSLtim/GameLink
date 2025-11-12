package model

import "time"

// ChatReport captures user reports against chat messages.
type ChatReport struct {
	Base
	MessageID  uint64  `json:"messageId" gorm:"column:message_id;not null;index"`
	ReporterID uint64  `json:"reporterId" gorm:"column:reporter_id;not null;index"`
	Reason     string  `json:"reason" gorm:"column:reason;type:text;not null"`
	Evidence   string  `json:"evidence" gorm:"column:evidence;type:text"`
	Status     string  `json:"status" gorm:"column:status;type:varchar(16);default:'pending';index"`
	HandledBy  *uint64 `json:"handledBy" gorm:"column:handled_by"`
	HandledAt  *time.Time `json:"handledAt" gorm:"column:handled_at"`
	Notes      string  `json:"notes" gorm:"column:notes;type:text"`
}

// TableName overrides default table name for chat reports.
func (ChatReport) TableName() string { return "chat_reports" }
