package model

import (
	"time"
)

// WithdrawStatus defines withdrawal status
// @Enum pending, approved, rejected, completed, failed
type WithdrawStatus string

const (
	// WithdrawStatusPending withdrawal is pending processing
	WithdrawStatusPending WithdrawStatus = "pending"
	// WithdrawStatusApproved withdrawal has been approved
	WithdrawStatusApproved WithdrawStatus = "approved"
	// WithdrawStatusRejected withdrawal has been rejected
	WithdrawStatusRejected WithdrawStatus = "rejected"
	// WithdrawStatusCompleted withdrawal has been completed
	WithdrawStatusCompleted WithdrawStatus = "completed"
	// WithdrawStatusFailed withdrawal processing failed
	WithdrawStatusFailed WithdrawStatus = "failed"
)

// WithdrawMethod defines withdrawal methods
// @Enum alipay, wechat, bank
type WithdrawMethod string

const (
	// WithdrawMethodAlipay Alipay payment method
	WithdrawMethodAlipay WithdrawMethod = "alipay"
	// WithdrawMethodWeChat WeChat payment method
	WithdrawMethodWeChat WithdrawMethod = "wechat"
	// WithdrawMethodBank bank transfer method
	WithdrawMethodBank WithdrawMethod = "bank"
)

// Withdraw represents a withdrawal request
// Requirements: 13.2, 13.5
type Withdraw struct {
	ID           uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	PlayerID     uint64         `gorm:"not null;index" json:"playerId"`
	UserID       uint64         `gorm:"not null;index" json:"userId"` // redundant field for convenient queries
	AmountCents  int64          `gorm:"not null" json:"amountCents"`  // withdrawal amount in cents
	Method       WithdrawMethod `gorm:"type:varchar(32);not null" json:"method"`
	AccountInfo  string         `gorm:"type:varchar(255);not null" json:"accountInfo"` // account info (encrypted storage)
	Status       WithdrawStatus `gorm:"type:varchar(32);not null;default:'pending'" json:"status"`
	RejectReason string         `gorm:"type:text" json:"rejectReason"` // reason for rejection
	AdminRemark  string         `gorm:"type:text" json:"adminRemark"`  // admin remarks
	ProcessedBy  *uint64        `gorm:"index" json:"processedBy"`      // processor ID
	ProcessedAt  *time.Time     `json:"processedAt"`                   // processing time
	CompletedAt  *time.Time     `json:"completedAt"`                   // completion time
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`

	// 提现分流相关字段 - Requirements: 13.2, 13.5
	SettlementCompanyID   *uint64    `gorm:"index" json:"settlementCompanyId,omitempty"`               // 结算公司ID
	SettlementCompanyName string     `gorm:"type:varchar(200)" json:"settlementCompanyName,omitempty"` // 结算公司名称（冗余字段，便于查询）
	PaymentBankAccount    string     `gorm:"type:varchar(30)" json:"paymentBankAccount,omitempty"`     // 付款银行账户
	BankTransactionNo     string     `gorm:"type:varchar(64)" json:"bankTransactionNo,omitempty"`      // 银行流水号
	TaxDeductedCents      int64      `gorm:"default:0" json:"taxDeductedCents"`                        // 代扣个税金额（分）
	ActualAmountCents     int64      `gorm:"default:0" json:"actualAmountCents"`                       // 实际到账金额（分）
	PaidAt                *time.Time `json:"paidAt,omitempty"`                                         // 付款时间

	// Relations
	SettlementCompany *SettlementCompany `json:"-" gorm:"foreignKey:SettlementCompanyID;references:ID"`
}

// TableName specifies the table name for Withdraw model
func (Withdraw) TableName() string {
	return "withdraws"
}

// CalculateActualAmount 计算实际到账金额
// 实际到账金额 = 提现金额 - 代扣个税
func (w *Withdraw) CalculateActualAmount() int64 {
	return w.AmountCents - w.TaxDeductedCents
}

// SetRoutingInfo 设置提现分流信息
// Requirements: 13.2
func (w *Withdraw) SetRoutingInfo(company *SettlementCompany) {
	if company != nil {
		w.SettlementCompanyID = &company.ID
		w.SettlementCompanyName = company.Name
		w.PaymentBankAccount = company.BankAccount
	}
}

// SalaryPaymentRecord 工资发放记录
// Requirements: 13.4
type SalaryPaymentRecord struct {
	Base
	WithdrawID            uint64     `json:"withdrawId" gorm:"column:withdraw_id;not null;index"`                    // 关联提现ID
	PlayerID              uint64     `json:"playerId" gorm:"column:player_id;not null;index"`                        // 陪玩师ID
	SettlementCompanyID   uint64     `json:"settlementCompanyId" gorm:"column:settlement_company_id;not null;index"` // 结算公司ID
	SettlementCompanyName string     `json:"settlementCompanyName" gorm:"column:settlement_company_name;size:200"`   // 结算公司名称
	AmountCents           int64      `json:"amountCents" gorm:"column:amount_cents;not null"`                        // 发放金额（分）
	TaxDeductedCents      int64      `json:"taxDeductedCents" gorm:"column:tax_deducted_cents;default:0"`            // 代扣个税（分）
	ActualAmountCents     int64      `json:"actualAmountCents" gorm:"column:actual_amount_cents;not null"`           // 实际到账（分）
	BankAccount           string     `json:"bankAccount" gorm:"column:bank_account;size:30"`                         // 付款银行账户
	BankTransactionNo     string     `json:"bankTransactionNo" gorm:"column:bank_transaction_no;size:64"`            // 银行流水号
	Status                string     `json:"status" gorm:"column:status;size:20;default:'pending'"`                  // 状态: pending, completed, failed
	PaidAt                *time.Time `json:"paidAt,omitempty" gorm:"column:paid_at"`                                 // 付款时间
	Remark                string     `json:"remark,omitempty" gorm:"column:remark;type:text"`                        // 备注

	// Relations
	Withdraw          *Withdraw          `json:"-" gorm:"foreignKey:WithdrawID;references:ID"`
	SettlementCompany *SettlementCompany `json:"-" gorm:"foreignKey:SettlementCompanyID;references:ID"`
}

// TableName specifies the table name for SalaryPaymentRecord model
func (SalaryPaymentRecord) TableName() string {
	return "salary_payment_records"
}

// WithdrawRoutingStats 提现分流统计
// Requirements: 14.1
type WithdrawRoutingStats struct {
	SettlementCompanyID    uint64  `json:"settlementCompanyId"`
	SettlementCompanyName  string  `json:"settlementCompanyName"`
	TotalWithdrawals       int64   `json:"totalWithdrawals"`       // 提现笔数
	TotalAmountCents       int64   `json:"totalAmountCents"`       // 提现总额（分）
	TotalTaxDeductedCents  int64   `json:"totalTaxDeductedCents"`  // 代扣个税总额（分）
	TotalActualAmountCents int64   `json:"totalActualAmountCents"` // 实际发放总额（分）
	AverageAmountCents     int64   `json:"averageAmountCents"`     // 平均提现金额（分）
	Percentage             float64 `json:"percentage"`             // 占比
}

// WithdrawRoutingReport 提现分流报表
// Requirements: 14.3, 14.4
type WithdrawRoutingReport struct {
	ID                     uint64                 `json:"id"`
	ReportType             string                 `json:"reportType"` // monthly, quarterly, yearly
	Year                   int                    `json:"year"`
	Month                  int                    `json:"month,omitempty"`
	Quarter                int                    `json:"quarter,omitempty"`
	TotalWithdrawals       int64                  `json:"totalWithdrawals"`
	TotalAmountCents       int64                  `json:"totalAmountCents"`
	TotalTaxDeductedCents  int64                  `json:"totalTaxDeductedCents"`
	TotalActualAmountCents int64                  `json:"totalActualAmountCents"`
	ByCompany              []WithdrawRoutingStats `json:"byCompany"`
	GeneratedAt            time.Time              `json:"generatedAt"`
}

// ListWithdrawsByCompanyRequest 按公司查询提现列表请求
type ListWithdrawsByCompanyRequest struct {
	SettlementCompanyID *uint64         `form:"settlementCompanyId"`
	Status              *WithdrawStatus `form:"status"`
	DateFrom            *time.Time      `form:"dateFrom"`
	DateTo              *time.Time      `form:"dateTo"`
	Page                int             `form:"page" binding:"omitempty,min=1"`
	PageSize            int             `form:"pageSize" binding:"omitempty,min=1,max=100"`
}

// ListWithdrawsByCompanyResponse 按公司查询提现列表响应
type ListWithdrawsByCompanyResponse struct {
	Total     int64      `json:"total"`
	Page      int        `json:"page"`
	PageSize  int        `json:"pageSize"`
	Withdraws []Withdraw `json:"withdraws"`
}

// WithdrawRoutingStatsRequest 提现分流统计请求
type WithdrawRoutingStatsRequest struct {
	DateFrom *time.Time `form:"dateFrom"`
	DateTo   *time.Time `form:"dateTo"`
}

// WithdrawRoutingStatsResponse 提现分流统计响应
type WithdrawRoutingStatsResponse struct {
	TotalWithdrawals       int64                  `json:"totalWithdrawals"`
	TotalAmountCents       int64                  `json:"totalAmountCents"`
	TotalTaxDeductedCents  int64                  `json:"totalTaxDeductedCents"`
	TotalActualAmountCents int64                  `json:"totalActualAmountCents"`
	ByCompany              []WithdrawRoutingStats `json:"byCompany"`
}

// WithdrawRoutingReportRequest 提现分流报表请求
type WithdrawRoutingReportRequest struct {
	ReportType string `form:"reportType" binding:"required,oneof=monthly quarterly yearly"`
	Year       int    `form:"year" binding:"required,min=2020,max=2100"`
	Month      int    `form:"month" binding:"omitempty,min=1,max=12"`
	Quarter    int    `form:"quarter" binding:"omitempty,min=1,max=4"`
}
