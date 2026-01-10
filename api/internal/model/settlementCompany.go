package model

import (
	"regexp"
	"time"
)

// CompanyStatus defines the status of a settlement company
// @Enum active, inactive
type CompanyStatus string

const (
	// CompanyStatusActive company is active and can receive new assignments
	CompanyStatusActive CompanyStatus = "active"
	// CompanyStatusInactive company is inactive and cannot receive new assignments
	CompanyStatusInactive CompanyStatus = "inactive"
)

// SettlementCompany represents a legal entity responsible for paying salaries to players
// Requirements: 11.1, 12.1
type SettlementCompany struct {
	Base
	Name              string        `json:"name" gorm:"column:name;size:200;not null"`                         // 公司名称
	CreditCode        string        `json:"creditCode" gorm:"column:credit_code;size:18;uniqueIndex;not null"` // 统一社会信用代码
	TaxRegistrationNo string        `json:"taxRegistrationNo" gorm:"column:tax_registration_no;size:20"`       // 税务登记号
	BankName          string        `json:"bankName" gorm:"column:bank_name;size:100"`                         // 开户银行
	BankAccount       string        `json:"bankAccount" gorm:"column:bank_account;size:30"`                    // 银行账号
	BankBranch        string        `json:"bankBranch" gorm:"column:bank_branch;size:200"`                     // 开户支行
	ContactName       string        `json:"contactName" gorm:"column:contact_name;size:50"`                    // 联系人
	ContactPhone      string        `json:"contactPhone" gorm:"column:contact_phone;size:20"`                  // 联系电话
	Address           string        `json:"address" gorm:"column:address;size:500"`                            // 公司地址
	Status            CompanyStatus `json:"status" gorm:"column:status;size:20;default:'active';index"`        // 状态
	PlayerCount       int           `json:"playerCount" gorm:"column:player_count;default:0"`                  // 关联陪玩师数量
	TotalPayoutCents  int64         `json:"totalPayoutCents" gorm:"column:total_payout_cents;default:0"`       // 累计发放金额（分）
	CreatedBy         uint64        `json:"createdBy" gorm:"column:created_by;not null;index"`                 // 创建人
	UpdatedBy         *uint64       `json:"updatedBy,omitempty" gorm:"column:updated_by;index"`                // 更新人

	// Relations
	Assignments []PlayerCompanyAssignment `json:"assignments,omitempty" gorm:"foreignKey:SettlementCompanyID;references:ID"`
	Creator     User                      `json:"-" gorm:"foreignKey:CreatedBy;references:ID"`
	Updater     *User                     `json:"-" gorm:"foreignKey:UpdatedBy;references:ID"`
}

// TableName specifies the table name for SettlementCompany model
func (SettlementCompany) TableName() string {
	return "settlement_companies"
}

// IsActive returns true if the company is active
func (s *SettlementCompany) IsActive() bool {
	return s.Status == CompanyStatusActive
}

// CreditCodeRegex is the regex pattern for validating Chinese Unified Social Credit Code
// The code consists of 18 characters:
// - Position 1: Registration management department code (1 digit)
// - Position 2: Organization category code (1 digit)
// - Position 3-8: Registration management authority administrative division code (6 digits)
// - Position 9-17: Subject identifier code (9 characters, alphanumeric)
// - Position 18: Check code (1 character, alphanumeric)
var CreditCodeRegex = regexp.MustCompile(`^[0-9A-HJ-NPQRTUWXY]{2}\d{6}[0-9A-HJ-NPQRTUWXY]{10}$`)

// ValidateCreditCode validates the format of a Chinese Unified Social Credit Code
// Requirements: 11.2, 15.2
// Property 11: 统一社会信用代码格式验证
func ValidateCreditCode(code string) bool {
	if len(code) != 18 {
		return false
	}
	return CreditCodeRegex.MatchString(code)
}

// ValidateCreditCodeWithChecksum validates the credit code format and checksum
// The checksum is calculated using a weighted sum algorithm
func ValidateCreditCodeWithChecksum(code string) bool {
	if !ValidateCreditCode(code) {
		return false
	}

	// Character mapping for checksum calculation
	charMap := map[rune]int{
		'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9,
		'A': 10, 'B': 11, 'C': 12, 'D': 13, 'E': 14, 'F': 15, 'G': 16, 'H': 17,
		'J': 18, 'K': 19, 'L': 20, 'M': 21, 'N': 22, 'P': 23, 'Q': 24, 'R': 25,
		'T': 26, 'U': 27, 'W': 28, 'X': 29, 'Y': 30,
	}

	// Weights for each position
	weights := []int{1, 3, 9, 27, 19, 26, 16, 17, 20, 29, 25, 13, 8, 24, 10, 30, 28}

	// Calculate weighted sum
	sum := 0
	runes := []rune(code)
	for i := 0; i < 17; i++ {
		val, ok := charMap[runes[i]]
		if !ok {
			return false
		}
		sum += val * weights[i]
	}

	// Calculate check digit
	remainder := sum % 31
	checkDigit := (31 - remainder) % 31

	// Get expected check character
	reverseMap := "0123456789ABCDEFGHJKLMNPQRTUWXY"
	expectedCheck := rune(reverseMap[checkDigit])

	return runes[17] == expectedCheck
}

// PlayerCompanyAssignment represents the assignment relationship between a player and a settlement company
// Requirements: 12.1
// Property 12: 陪玩师结算公司分配唯一性 - Each player can only have one current assignment
type PlayerCompanyAssignment struct {
	Base
	PlayerID            uint64     `json:"playerId" gorm:"column:player_id;not null;index;uniqueIndex:idx_player_current_unique,where:is_current = true"` // 陪玩师ID (PostgreSQL requires boolean comparison)
	SettlementCompanyID uint64     `json:"settlementCompanyId" gorm:"column:settlement_company_id;not null;index"`                                        // 结算公司ID
	EffectiveDate       time.Time  `json:"effectiveDate" gorm:"column:effective_date;not null;index"`                                                     // 生效日期
	EndDate             *time.Time `json:"endDate,omitempty" gorm:"column:end_date;index"`                                                                // 结束日期
	Reason              string     `json:"reason" gorm:"column:reason;size:500"`                                                                          // 分配原因
	AssignedBy          uint64     `json:"assignedBy" gorm:"column:assigned_by;not null;index"`                                                           // 分配操作人
	IsCurrent           bool       `json:"isCurrent" gorm:"column:is_current;default:false;index"`                                                        // 是否当前生效

	// Relations
	Player            Player            `json:"-" gorm:"foreignKey:PlayerID;references:ID"`
	SettlementCompany SettlementCompany `json:"-" gorm:"foreignKey:SettlementCompanyID;references:ID"`
	Assigner          User              `json:"-" gorm:"foreignKey:AssignedBy;references:ID"`
}

// TableName specifies the table name for PlayerCompanyAssignment model
func (PlayerCompanyAssignment) TableName() string {
	return "player_company_assignments"
}

// IsEffective checks if the assignment is currently effective
func (a *PlayerCompanyAssignment) IsEffective() bool {
	now := time.Now()
	if a.EffectiveDate.After(now) {
		return false
	}
	if a.EndDate != nil && a.EndDate.Before(now) {
		return false
	}
	return true
}

// SettlementCompanyHistory represents a history record of settlement company changes
type SettlementCompanyHistory struct {
	Base
	SettlementCompanyID uint64 `json:"settlementCompanyId" gorm:"column:settlement_company_id;not null;index"` // 结算公司ID
	FieldName           string `json:"fieldName" gorm:"column:field_name;size:50;not null"`                    // 修改字段名
	OldValue            string `json:"oldValue" gorm:"column:old_value;type:text"`                             // 修改前值
	NewValue            string `json:"newValue" gorm:"column:new_value;type:text"`                             // 修改后值
	ChangedBy           uint64 `json:"changedBy" gorm:"column:changed_by;not null;index"`                      // 修改人

	// Relations
	SettlementCompany SettlementCompany `json:"-" gorm:"foreignKey:SettlementCompanyID;references:ID"`
	Changer           User              `json:"-" gorm:"foreignKey:ChangedBy;references:ID"`
}

// TableName specifies the table name for SettlementCompanyHistory model
func (SettlementCompanyHistory) TableName() string {
	return "settlement_company_histories"
}

// CreateSettlementCompanyRequest represents a request to create a settlement company
type CreateSettlementCompanyRequest struct {
	Name              string `json:"name" binding:"required,max=200"`
	CreditCode        string `json:"creditCode" binding:"required,len=18"`
	TaxRegistrationNo string `json:"taxRegistrationNo" binding:"max=20"`
	BankName          string `json:"bankName" binding:"max=100"`
	BankAccount       string `json:"bankAccount" binding:"max=30"`
	BankBranch        string `json:"bankBranch" binding:"max=200"`
	ContactName       string `json:"contactName" binding:"max=50"`
	ContactPhone      string `json:"contactPhone" binding:"max=20"`
	Address           string `json:"address" binding:"max=500"`
}

// UpdateSettlementCompanyRequest represents a request to update a settlement company
type UpdateSettlementCompanyRequest struct {
	Name              *string `json:"name" binding:"omitempty,max=200"`
	TaxRegistrationNo *string `json:"taxRegistrationNo" binding:"omitempty,max=20"`
	BankName          *string `json:"bankName" binding:"omitempty,max=100"`
	BankAccount       *string `json:"bankAccount" binding:"omitempty,max=30"`
	BankBranch        *string `json:"bankBranch" binding:"omitempty,max=200"`
	ContactName       *string `json:"contactName" binding:"omitempty,max=50"`
	ContactPhone      *string `json:"contactPhone" binding:"omitempty,max=20"`
	Address           *string `json:"address" binding:"omitempty,max=500"`
}

// AssignPlayerToCompanyRequest represents a request to assign a player to a settlement company
type AssignPlayerToCompanyRequest struct {
	PlayerID            uint64    `json:"playerId" binding:"required"`
	SettlementCompanyID uint64    `json:"settlementCompanyId" binding:"required"`
	EffectiveDate       time.Time `json:"effectiveDate" binding:"required"`
	Reason              string    `json:"reason" binding:"required,max=500"`
}

// BatchAssignPlayersRequest represents a request to batch assign players to a settlement company
type BatchAssignPlayersRequest struct {
	PlayerIDs           []uint64  `json:"playerIds" binding:"required,min=1"`
	SettlementCompanyID uint64    `json:"settlementCompanyId" binding:"required"`
	EffectiveDate       time.Time `json:"effectiveDate" binding:"required"`
	Reason              string    `json:"reason" binding:"required,max=500"`
}

// ListSettlementCompaniesRequest represents a request to list settlement companies
type ListSettlementCompaniesRequest struct {
	Page      int           `form:"page" binding:"omitempty,min=1"`
	PageSize  int           `form:"pageSize" binding:"omitempty,min=1,max=100"`
	Status    CompanyStatus `form:"status" binding:"omitempty,oneof=active inactive"`
	Keyword   string        `form:"keyword" binding:"omitempty,max=100"`
	SortBy    string        `form:"sortBy" binding:"omitempty,oneof=name created_at player_count"`
	SortOrder string        `form:"sortOrder" binding:"omitempty,oneof=asc desc"`
}

// ListSettlementCompaniesResponse represents a response for listing settlement companies
type ListSettlementCompaniesResponse struct {
	Total     int64               `json:"total"`
	Page      int                 `json:"page"`
	PageSize  int                 `json:"pageSize"`
	Companies []SettlementCompany `json:"companies"`
}

// PlayerAssignmentHistoryResponse represents a response for player assignment history
type PlayerAssignmentHistoryResponse struct {
	Total       int64                     `json:"total"`
	Assignments []PlayerCompanyAssignment `json:"assignments"`
}
