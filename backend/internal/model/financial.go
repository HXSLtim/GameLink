package model

import (
	"encoding/json"
	"time"
)

// FinancialAccountType 财务科目类型
type FinancialAccountType string

const (
	FinancialAccountTypeAsset     FinancialAccountType = "asset"     // 资产
	FinancialAccountTypeLiability  FinancialAccountType = "liability"  // 负债
	FinancialAccountTypeEquity     FinancialAccountType = "equity"     // 所有者权益
	FinancialAccountTypeRevenue    FinancialAccountType = "revenue"    // 收入
	FinancialAccountTypeExpense    FinancialAccountType = "expense"    // 费用
)

// FinancialAccountLevel 科目级别
type FinancialAccountLevel int

const (
	FinancialAccountLevel1 FinancialAccountLevel = 1 // 一级科目
	FinancialAccountLevel2 FinancialAccountLevel = 2 // 二级科目
	FinancialAccountLevel3 FinancialAccountLevel = 3 // 三级科目
)

// FinancialAccountDirection 余额方向
type FinancialAccountDirection string

const (
	FinancialAccountDirectionDebit  FinancialAccountDirection = "debit"  // 借方
	FinancialAccountDirectionCredit FinancialAccountDirection = "credit" // 贷方
)

// FinancialAccountStatus 科目状态
type FinancialAccountStatus string

const (
	FinancialAccountStatusActive   FinancialAccountStatus = "active"   // 启用
	FinancialAccountStatusInactive FinancialAccountStatus = "inactive" // 停用
)

// FinancialAccount 财务科目
type FinancialAccount struct {
	Base
	Code             string                   `json:"code" gorm:"column:code;size:20;uniqueIndex;not null"`              // 科目编码
	Name             string                   `json:"name" gorm:"column:name;size:100;not null"`                         // 科目名称
	Type             FinancialAccountType     `json:"type" gorm:"column:type;size:20;not null;index"`                      // 科目类型
	Level            FinancialAccountLevel    `json:"level" gorm:"column:level;not null"`                                  // 科目级别
	ParentCode       *string                  `json:"parentCode,omitempty" gorm:"column:parent_code;size:20;index"`          // 父科目编码
	Direction        FinancialAccountDirection `json:"direction" gorm:"column:direction;size:10;not null"`                    // 余额方向
	OpeningBalance   int64                    `json:"openingBalance" gorm:"column:opening_balance;default:0"`               // 期初余额（分）
	CurrentBalance   int64                    `json:"currentBalance" gorm:"column:current_balance;default:0"`               // 当前余额（分）
	Description      string                   `json:"description,omitempty" gorm:"column:description;size:500"`           // 科目描述
	Status           FinancialAccountStatus   `json:"status" gorm:"column:status;size:20;default:'active'"`                // 状态
	IsSystem         bool                     `json:"isSystem" gorm:"column:is_system;default:false"`                      // 是否系统科目

	// 审计字段
	CreatedBy        uint64                   `json:"createdBy" gorm:"column:created_by;not null;index"`                    // 创建人
	UpdatedBy        *uint64                  `json:"updatedBy,omitempty" gorm:"column:updated_by;index"`                   // 更新人

	// 关联关系
	Parent           *FinancialAccount       `json:"parent,omitempty" gorm:"foreignKey:ParentCode;references:Code"`
	Children         []FinancialAccount      `json:"children,omitempty" gorm:"foreignKey:ParentCode;references:Code"`
	Vouchers         []FinancialVoucher      `json:"-" gorm:"foreignKey:AccountCode;references:Code"`
	Transactions     []FinancialTransaction  `json:"-" gorm:"foreignKey:AccountCode;references:Code"`
}

// FinancialVoucherType 凭证类型
type FinancialVoucherType string

const (
	FinancialVoucherTypeManual    FinancialVoucherType = "manual"    // 手工凭证
	FinancialVoucherTypeAuto      FinancialVoucherType = "auto"      // 自动凭证
	FinancialVoucherTypeAdjust    FinancialVoucherType = "adjust"    // 调整凭证
	FinancialVoucherTypeClosing   FinancialVoucherType = "closing"   // 结账凭证
	FinancialVoucherTypeReversal  FinancialVoucherType = "reversal"  // 红字冲销凭证
)

// FinancialVoucherStatus 凭证状态
type FinancialVoucherStatus string

const (
	FinancialVoucherStatusDraft     FinancialVoucherStatus = "draft"     // 草稿
	FinancialVoucherStatusPending   FinancialVoucherStatus = "pending"   // 待审核
	FinancialVoucherStatusApproved  FinancialVoucherStatus = "approved"  // 已审核
	FinancialVoucherStatusRejected  FinancialVoucherStatus = "rejected"  // 已驳回
	FinancialVoucherStatusPosted    FinancialVoucherStatus = "posted"    // 已过账
	FinancialVoucherStatusCancelled FinancialVoucherStatus = "cancelled" // 已取消
)

// FinancialVoucher 财务凭证
type FinancialVoucher struct {
	Base
	VoucherNo        string                  `json:"voucherNo" gorm:"column:voucher_no;size:32;uniqueIndex;not null"` // 凭证号
	VoucherDate      time.Time               `json:"voucherDate" gorm:"column:voucher_date;not null;index"`             // 凭证日期
	Type             FinancialVoucherType    `json:"type" gorm:"column:type;size:20;not null;index"`                    // 凭证类型
	Abstract         string                  `json:"abstract" gorm:"column:abstract;size:200;not null"`                   // 摘要
	AttachmentCount  int                     `json:"attachmentCount" gorm:"column:attachment_count;default:0"`             // 附件数量
	Status           FinancialVoucherStatus  `json:"status" gorm:"column:status;size:20;default:'draft'"`                // 状态
	TotalAmount      int64                   `json:"totalAmount" gorm:"column:total_amount;default:0"`                     // 总金额（分）

	// 业务关联
	BusinessType     string                  `json:"businessType,omitempty" gorm:"column:business_type;size:50"`           // 业务类型
	BusinessNo       string                  `json:"businessNo,omitempty" gorm:"column:business_no;size:64"`               // 业务单号

	// 审计字段
	CreatedBy        uint64                  `json:"createdBy" gorm:"column:created_by;not null;index"`                    // 制单人
	ReviewedBy       *uint64                 `json:"reviewedBy,omitempty" gorm:"column:reviewed_by;index"`                 // 审核人
	PostedBy         *uint64                 `json:"postedBy,omitempty" gorm:"column:posted_by;index"`                    // 过账人
	ReviewedAt       *time.Time              `json:"reviewedAt,omitempty" gorm:"column:reviewed_at"`                        // 审核时间
	PostedAt         *time.Time              `json:"postedAt,omitempty" gorm:"column:posted_at"`                            // 过账时间

	// 关联关系
	Entries          []FinancialVoucherEntry `json:"entries,omitempty" gorm:"foreignKey:VoucherID;references:ID"`
	Creator          User                    `json:"-" gorm:"foreignKey:CreatedBy;references:ID"`
	Reviewer         *User                   `json:"-" gorm:"foreignKey:ReviewedBy;references:ID"`
	Poster           *User                   `json:"-" gorm:"foreignKey:PostedBy;references:ID"`
}

// FinancialVoucherEntry 凭证分录
type FinancialVoucherEntry struct {
	Base
	VoucherID        uint64                   `json:"voucherId" gorm:"column:voucher_id;not null;index"`                       // 凭证ID
	LineNo           int                      `json:"lineNo" gorm:"column:line_no;not null"`                                    // 分录行号
	AccountCode      string                   `json:"accountCode" gorm:"column:account_code;size:20;not null;index"`           // 科目编码
	Abstract         string                   `json:"abstract" gorm:"column:abstract;size:200"`                                  // 摘要
	DebitAmount      int64                    `json:"debitAmount" gorm:"column:debit_amount;default:0"`                          // 借方金额（分）
	CreditAmount     int64                    `json:"creditAmount" gorm:"column:credit_amount;default:0"`                         // 贷方金额（分）
	Currency         Currency                 `json:"currency,omitempty" gorm:"type:char(3);default:'CNY'"`                    // 货币

	// 业务关联
	BusinessType     string                   `json:"businessType,omitempty" gorm:"column:business_type;size:50"`              // 业务类型
	BusinessNo       string                   `json:"businessNo,omitempty" gorm:"column:business_no;size:64"`                // 业务单号
	RelatedEntity    string                   `json:"relatedEntity,omitempty" gorm:"column:related_entity;size:50"`           // 关联实体
	RelatedEntityID  *uint64                  `json:"relatedEntityId,omitempty" gorm:"column:related_entity_id"`              // 关联实体ID

	// 关联关系
	Voucher          FinancialVoucher         `json:"-" gorm:"foreignKey:VoucherID;references:ID"`
	Account          FinancialAccount        `json:"-" gorm:"foreignKey:AccountCode;references:Code"`
}

// FinancialTransactionType 交易类型
type FinancialTransactionType string

const (
	FinancialTransactionTypeRevenue    FinancialTransactionType = "revenue"    // 收入
	FinancialTransactionTypeExpense    FinancialTransactionType = "expense"    // 支出
	FinancialTransactionTypeTransfer   FinancialTransactionType = "transfer"   // 转账
	FinancialTransactionTypeAdjustment FinancialTransactionType = "adjustment" // 调整
	FinancialTransactionTypeOpening    FinancialTransactionType = "opening"    // 开账
	FinancialTransactionTypeClosing    FinancialTransactionType = "closing"    // 结账
)

// FinancialTransaction 财务流水
type FinancialTransaction struct {
	Base
	TransactionNo    string                     `json:"transactionNo" gorm:"column:transaction_no;size:32;uniqueIndex;not null"` // 交易流水号
	TransactionDate  time.Time                  `json:"transactionDate" gorm:"column:transaction_date;not null;index"`             // 交易日期
	AccountCode      string                     `json:"accountCode" gorm:"column:account_code;size:20;not null;index"`           // 科目编码
	Type             FinancialTransactionType   `json:"type" gorm:"column:type;size:20;not null;index"`                        // 交易类型
	Direction        FinancialAccountDirection   `json:"direction" gorm:"column:direction;size:10;not null"`                      // 借贷方向
	Amount           int64                      `json:"amount" gorm:"column:amount;not null"`                                      // 金额（分）
	BalanceBefore    int64                      `json:"balanceBefore" gorm:"column:balance_before;not null"`                       // 交易前余额（分）
	BalanceAfter     int64                      `json:"balanceAfter" gorm:"column:balance_after;not null"`                         // 交易后余额（分）
	Currency         Currency                   `json:"currency,omitempty" gorm:"type:char(3);default:'CNY'"`                    // 货币
	Abstract         string                     `json:"abstract" gorm:"column:abstract;size:200"`                                  // 摘要

	// 业务关联
	BusinessType     string                     `json:"businessType,omitempty" gorm:"column:business_type;size:50"`              // 业务类型
	BusinessNo       string                     `json:"businessNo,omitempty" gorm:"column:business_no;size:64"`                // 业务单号
	VoucherID        *uint64                    `json:"voucherId,omitempty" gorm:"column:voucher_id;index"`                      // 关联凭证ID
	VoucherEntryID   *uint64                    `json:"voucherEntryId,omitempty" gorm:"column:voucher_entry_id"`                 // 关联凭证分录ID

	// 关联关系
	Account          FinancialAccount           `json:"-" gorm:"foreignKey:AccountCode;references:Code"`
	Voucher          *FinancialVoucher         `json:"-" gorm:"foreignKey:VoucherID;references:ID"`
	VoucherEntry     *FinancialVoucherEntry    `json:"-" gorm:"foreignKey:VoucherEntryID;references:ID"`
}

// ReconciliationType 对账类型
type ReconciliationType string

const (
	ReconciliationTypePayment     ReconciliationType = "payment"     // 支付对账
	ReconciliationTypeInternal    ReconciliationType = "internal"    // 内部对账
	ReconciliationTypeBank        ReconciliationType = "bank"        // 银行对账
	ReconciliationTypeManual      ReconciliationType = "manual"      // 手工对账
)

// ReconciliationStatus 对账状态
type ReconciliationStatus string

const (
	ReconciliationStatusPending   ReconciliationStatus = "pending"   // 待对账
	ReconciliationStatusProgress  ReconciliationStatus = "progress"  // 对账中
	ReconciliationStatusSuccess   ReconciliationStatus = "success"   // 对账成功
	ReconciliationStatusFailed    ReconciliationStatus = "failed"    // 对账失败
	ReconciliationStatusException ReconciliationStatus = "exception" // 异常
)

// Reconciliation 对账单
type Reconciliation struct {
	Base
	ReconciliationNo string               `json:"reconciliationNo" gorm:"column:reconciliation_no;size:32;uniqueIndex;not null"` // 对账单号
	ReconciliationDate time.Time           `json:"reconciliationDate" gorm:"column:reconciliation_date;not null;index"`             // 对账日期
	Type             ReconciliationType   `json:"type" gorm:"column:type;size:20;not null;index"`                        // 对账类型
	Status           ReconciliationStatus `json:"status" gorm:"column:status;size:20;default:'pending'"`                      // 对账状态
	PeriodStart      time.Time            `json:"periodStart" gorm:"column:period_start;not null"`                              // 对账期间开始
	PeriodEnd        time.Time            `json:"periodEnd" gorm:"column:period_end;not null"`                                // 对账期间结束
	TotalRecords     int                  `json:"totalRecords" gorm:"column:total_records;default:0"`                         // 总记录数
	MatchedRecords   int                  `json:"matchedRecords" gorm:"column:matched_records;default:0"`                       // 匹配记录数
	DifferenceAmount int64                `json:"differenceAmount" gorm:"column:difference_amount;default:0"`                   // 差异金额（分）
	Abstract         string               `json:"abstract" gorm:"column:abstract;size:500"`                                   // 摘要

	// 处理信息
	ProcessedAt      *time.Time           `json:"processedAt,omitempty" gorm:"column:processed_at"`                           // 处理时间
	ProcessedBy      *uint64              `json:"processedBy,omitempty" gorm:"column:processed_by;index"`                     // 处理人

	// 关联关系
	Details          []ReconciliationDetail `json:"details,omitempty" gorm:"foreignKey:ReconciliationID;references:ID"`
	Processor        *User                  `json:"-" gorm:"foreignKey:ProcessedBy;references:ID"`
}

// ReconciliationDetail 对账明细
type ReconciliationDetail struct {
	Base
	ReconciliationID uint64 `json:"reconciliationId" gorm:"column:reconciliation_id;not null;index"` // 对账单ID
	LineNo           int    `json:"lineNo" gorm:"column:line_no;not null"`                            // 行号

	// 外部数据
	ExternalType     string `json:"externalType" gorm:"column:external_type;size:50"`                 // 外部类型
	ExternalNo       string `json:"externalNo" gorm:"column:external_no;size:64"`                     // 外部单号
	ExternalAmount   int64  `json:"externalAmount" gorm:"column:external_amount;not null"`             // 外部金额（分）
	ExternalDate     time.Time `json:"externalDate" gorm:"column:external_date;not null"`             // 外部日期

	// 内部数据
	InternalType     string `json:"internalType" gorm:"column:internal_type;size:50"`                 // 内部类型
	InternalNo       string `json:"internalNo" gorm:"column:internal_no;size:64"`                     // 内部单号
	InternalAmount   int64  `json:"internalAmount" gorm:"column:internal_amount;not null"`             // 内部金额（分）
	InternalDate     time.Time `json:"internalDate" gorm:"column:internal_date;not null"`             // 内部日期

	// 对账结果
	Status           string `json:"status" gorm:"column:status;size:20;default:'pending'"`           // 对账状态
	DifferenceAmount int64  `json:"differenceAmount" gorm:"column:difference_amount;default:0"`       // 差异金额（分）
	Remark           string `json:"remark,omitempty" gorm:"column:remark;size:500"`                   // 备注

	// 关联关系
	Reconciliation   Reconciliation `json:"-" gorm:"foreignKey:ReconciliationID;references:ID"`
}

// FinancialReportType 财务报表类型
type FinancialReportType string

const (
	FinancialReportTypeBalanceSheet    FinancialReportType = "balance_sheet"    // 资产负债表
	FinancialReportTypeIncomeStatement  FinancialReportType = "income_statement"  // 利润表
	FinancialReportTypeCashFlow        FinancialReportType = "cash_flow"        // 现金流量表
	FinancialReportTypeTrialBalance     FinancialReportType = "trial_balance"     // 试算平衡表
	FinancialReportTypeCustom          FinancialReportType = "custom"           // 自定义报表
)

// FinancialReport 财务报表
type FinancialReport struct {
	Base
	ReportNo         string              `json:"reportNo" gorm:"column:report_no;size:32;uniqueIndex;not null"` // 报表编号
	ReportName       string              `json:"reportName" gorm:"column:report_name;size:100;not null"`         // 报表名称
	Type             FinancialReportType `json:"type" gorm:"column:type;size:50;not null;index"`                // 报表类型
	PeriodStart      time.Time           `json:"periodStart" gorm:"column:period_start;not null"`                // 报告期间开始
	PeriodEnd        time.Time           `json:"periodEnd" gorm:"column:period_end;not null"`                  // 报告期间结束
	Currency         Currency            `json:"currency,omitempty" gorm:"type:char(3);default:'CNY'"`         // 货币单位
	ReportData       json.RawMessage    `json:"reportData" gorm:"column:report_data;type:json"`                 // 报表数据(JSON)
	Status           string              `json:"status" gorm:"column:status;size:20;default:'draft'"`           // 报表状态

	// 生成信息
	GeneratedAt      *time.Time          `json:"generatedAt,omitempty" gorm:"column:generated_at"`               // 生成时间
	GeneratedBy      *uint64             `json:"generatedBy,omitempty" gorm:"column:generated_by;index"`          // 生成人

	// 关联关系
	Generator        *User               `json:"-" gorm:"foreignKey:GeneratedBy;references:ID"`
}

// FinancialAccountSetting 科目设置
type FinancialAccountSetting struct {
	Base
	AccountCode      string    `json:"accountCode" gorm:"column:account_code;size:20;not null;index"` // 科目编码
	AllowManual      bool      `json:"allowManual" gorm:"column:allow_manual;default:true"`            // 允许手工录入
	RequireApproval  bool      `json:"requireApproval" gorm:"column:require_approval;default:false"`      // 需要审批
	MaxAmount        int64     `json:"maxAmount" gorm:"column:max_amount;default:0"`                      // 最大金额（分）
	MinAmount        int64     `json:"minAmount" gorm:"column:min_amount;default:0"`                      // 最小金额（分）
	Description      string    `json:"description,omitempty" gorm:"column:description;size:500"`         // 描述

	// 关联关系
	Account          FinancialAccount `json:"-" gorm:"foreignKey:AccountCode;references:Code"`
}

// TableName 返回表名
func (FinancialAccount) TableName() string {
	return "financial_accounts"
}

func (FinancialVoucher) TableName() string {
	return "financial_vouchers"
}

func (FinancialVoucherEntry) TableName() string {
	return "financial_voucher_entries"
}

func (FinancialTransaction) TableName() string {
	return "financial_transactions"
}

func (Reconciliation) TableName() string {
	return "reconciliations"
}

func (ReconciliationDetail) TableName() string {
	return "reconciliation_details"
}

func (FinancialReport) TableName() string {
	return "financial_reports"
}

func (FinancialAccountSetting) TableName() string {
	return "financial_account_settings"
}