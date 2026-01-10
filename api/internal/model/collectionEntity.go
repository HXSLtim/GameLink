package model

import (
	"encoding/json"
	"time"
)

// EntityStatus defines the status of a collection entity
// @Enum active, inactive
type EntityStatus string

const (
	// EntityStatusActive entity is active and can receive payments
	EntityStatusActive EntityStatus = "active"
	// EntityStatusInactive entity is inactive and cannot receive payments
	EntityStatusInactive EntityStatus = "inactive"
)

// CollectionEntity represents a legal entity that receives user payments
// Requirements: 15.1, 16.1
type CollectionEntity struct {
	Base
	Name                 string       `json:"name" gorm:"column:name;size:200;not null"`                           // 公司名称
	CreditCode           string       `json:"creditCode" gorm:"column:credit_code;size:18;uniqueIndex;not null"`   // 统一社会信用代码
	TaxRegistrationNo    string       `json:"taxRegistrationNo" gorm:"column:tax_registration_no;size:20"`         // 税务登记号
	Status               EntityStatus `json:"status" gorm:"column:status;size:20;default:'active';index"`          // 状态
	IsDefault            bool         `json:"isDefault" gorm:"column:is_default;default:false;index"`              // 是否默认收款主体
	TotalCollectionCents int64        `json:"totalCollectionCents" gorm:"column:total_collection_cents;default:0"` // 累计收款金额（分）
	TransactionCount     int64        `json:"transactionCount" gorm:"column:transaction_count;default:0"`          // 交易笔数
	CreatedBy            uint64       `json:"createdBy" gorm:"column:created_by;not null;index"`                   // 创建人
	UpdatedBy            *uint64      `json:"updatedBy,omitempty" gorm:"column:updated_by;index"`                  // 更新人

	// Relations
	PaymentChannels []PaymentChannelConfig `json:"paymentChannels,omitempty" gorm:"foreignKey:CollectionEntityID;references:ID"`
	RoutingRules    []RoutingRule          `json:"routingRules,omitempty" gorm:"foreignKey:TargetEntityID;references:ID"`
	Creator         User                   `json:"-" gorm:"foreignKey:CreatedBy;references:ID"`
	Updater         *User                  `json:"-" gorm:"foreignKey:UpdatedBy;references:ID"`
}

// TableName specifies the table name for CollectionEntity model
func (CollectionEntity) TableName() string {
	return "collection_entities"
}

// IsActive returns true if the entity is active
func (e *CollectionEntity) IsActive() bool {
	return e.Status == EntityStatusActive
}

// HasRequiredFields checks if the collection entity has all required fields
// Property 17: 收款分流记录完整性
func (e *CollectionEntity) HasRequiredFields() bool {
	return e.Name != "" &&
		e.CreditCode != "" &&
		e.Status != "" &&
		!e.CreatedAt.IsZero()
}

// PaymentChannelConfig represents the payment channel configuration for a collection entity
// Requirements: 15.4
type PaymentChannelConfig struct {
	Base
	CollectionEntityID uint64        `json:"collectionEntityId" gorm:"column:collection_entity_id;not null;index;uniqueIndex:idx_entity_channel"` // 收款主体ID
	Channel            PaymentMethod `json:"channel" gorm:"column:channel;size:32;not null;uniqueIndex:idx_entity_channel"`                       // 支付渠道
	MerchantNo         string        `json:"merchantNo" gorm:"column:merchant_no;size:64;not null"`                                               // 商户号
	MerchantKey        string        `json:"-" gorm:"column:merchant_key;size:256"`                                                               // 商户密钥（加密存储，不返回给前端）
	CallbackURL        string        `json:"callbackUrl" gorm:"column:callback_url;size:500"`                                                     // 回调地址
	Enabled            bool          `json:"enabled" gorm:"column:enabled;default:true;index"`                                                    // 是否启用
	Priority           int           `json:"priority" gorm:"column:priority;default:0"`                                                           // 优先级（同渠道多配置时使用）
	Remark             string        `json:"remark,omitempty" gorm:"column:remark;size:500"`                                                      // 备注

	// Relations
	CollectionEntity CollectionEntity `json:"-" gorm:"foreignKey:CollectionEntityID;references:ID"`
}

// TableName specifies the table name for PaymentChannelConfig model
func (PaymentChannelConfig) TableName() string {
	return "payment_channel_configs"
}

// IsEnabled returns true if the channel is enabled
func (c *PaymentChannelConfig) IsEnabled() bool {
	return c.Enabled
}

// ConditionField defines the field types for routing conditions
// @Enum game_type, service_type, order_amount, region
type ConditionField string

const (
	ConditionFieldGameType    ConditionField = "game_type"
	ConditionFieldServiceType ConditionField = "service_type"
	ConditionFieldOrderAmount ConditionField = "order_amount"
	ConditionFieldRegion      ConditionField = "region"
)

// ConditionOperator defines the operators for routing conditions
// @Enum eq, neq, in, not_in, gt, lt, between
type ConditionOperator string

const (
	ConditionOperatorEquals      ConditionOperator = "eq"
	ConditionOperatorNotEquals   ConditionOperator = "neq"
	ConditionOperatorIn          ConditionOperator = "in"
	ConditionOperatorNotIn       ConditionOperator = "not_in"
	ConditionOperatorGreaterThan ConditionOperator = "gt"
	ConditionOperatorLessThan    ConditionOperator = "lt"
	ConditionOperatorBetween     ConditionOperator = "between"
)

// RuleStatus defines the status of a routing rule
// @Enum active, inactive
type RuleStatus string

const (
	RuleStatusActive   RuleStatus = "active"
	RuleStatusInactive RuleStatus = "inactive"
)

// RoutingCondition represents a single condition for routing rule matching
type RoutingCondition struct {
	Field    ConditionField    `json:"field"`    // 条件字段
	Operator ConditionOperator `json:"operator"` // 操作符
	Value    json.RawMessage   `json:"value"`    // 条件值（可以是字符串、数字或数组）
}

// RoutingRule represents a rule for routing payments to collection entities
// Requirements: 16.1, 16.2
type RoutingRule struct {
	Base
	Name           string          `json:"name" gorm:"column:name;size:100;not null"`                    // 规则名称
	Priority       int             `json:"priority" gorm:"column:priority;not null;index"`               // 优先级（数字越小优先级越高）
	Conditions     json.RawMessage `json:"conditions" gorm:"column:conditions;type:json"`                // 匹配条件（JSON数组）
	TargetEntityID uint64          `json:"targetEntityId" gorm:"column:target_entity_id;not null;index"` // 目标收款主体ID
	Status         RuleStatus      `json:"status" gorm:"column:status;size:20;default:'active';index"`   // 状态
	Description    string          `json:"description,omitempty" gorm:"column:description;size:500"`     // 规则描述
	CreatedBy      uint64          `json:"createdBy" gorm:"column:created_by;not null;index"`            // 创建人
	UpdatedBy      *uint64         `json:"updatedBy,omitempty" gorm:"column:updated_by;index"`           // 更新人

	// Relations
	TargetEntity CollectionEntity `json:"targetEntity,omitempty" gorm:"foreignKey:TargetEntityID;references:ID"`
	Creator      User             `json:"-" gorm:"foreignKey:CreatedBy;references:ID"`
	Updater      *User            `json:"-" gorm:"foreignKey:UpdatedBy;references:ID"`
}

// TableName specifies the table name for RoutingRule model
func (RoutingRule) TableName() string {
	return "routing_rules"
}

// IsActive returns true if the rule is active
func (r *RoutingRule) IsActive() bool {
	return r.Status == RuleStatusActive
}

// GetConditions parses and returns the routing conditions
func (r *RoutingRule) GetConditions() ([]RoutingCondition, error) {
	if r.Conditions == nil {
		return nil, nil
	}
	var conditions []RoutingCondition
	if err := json.Unmarshal(r.Conditions, &conditions); err != nil {
		return nil, err
	}
	return conditions, nil
}

// SetConditions sets the routing conditions from a slice
func (r *RoutingRule) SetConditions(conditions []RoutingCondition) error {
	data, err := json.Marshal(conditions)
	if err != nil {
		return err
	}
	r.Conditions = data
	return nil
}

// RoutingRuleHistory represents a history record of routing rule changes
type RoutingRuleHistory struct {
	Base
	RoutingRuleID uint64 `json:"routingRuleId" gorm:"column:routing_rule_id;not null;index"` // 分流规则ID
	FieldName     string `json:"fieldName" gorm:"column:field_name;size:50;not null"`        // 修改字段名
	OldValue      string `json:"oldValue" gorm:"column:old_value;type:text"`                 // 修改前值
	NewValue      string `json:"newValue" gorm:"column:new_value;type:text"`                 // 修改后值
	ChangedBy     uint64 `json:"changedBy" gorm:"column:changed_by;not null;index"`          // 修改人

	// Relations
	RoutingRule RoutingRule `json:"-" gorm:"foreignKey:RoutingRuleID;references:ID"`
	Changer     User        `json:"-" gorm:"foreignKey:ChangedBy;references:ID"`
}

// TableName specifies the table name for RoutingRuleHistory model
func (RoutingRuleHistory) TableName() string {
	return "routing_rule_histories"
}

// CollectionEntityHistory represents a history record of collection entity changes
type CollectionEntityHistory struct {
	Base
	CollectionEntityID uint64 `json:"collectionEntityId" gorm:"column:collection_entity_id;not null;index"` // 收款主体ID
	FieldName          string `json:"fieldName" gorm:"column:field_name;size:50;not null"`                  // 修改字段名
	OldValue           string `json:"oldValue" gorm:"column:old_value;type:text"`                           // 修改前值
	NewValue           string `json:"newValue" gorm:"column:new_value;type:text"`                           // 修改后值
	ChangedBy          uint64 `json:"changedBy" gorm:"column:changed_by;not null;index"`                    // 修改人

	// Relations
	CollectionEntity CollectionEntity `json:"-" gorm:"foreignKey:CollectionEntityID;references:ID"`
	Changer          User             `json:"-" gorm:"foreignKey:ChangedBy;references:ID"`
}

// TableName specifies the table name for CollectionEntityHistory model
func (CollectionEntityHistory) TableName() string {
	return "collection_entity_histories"
}

// RoutingLog represents a log entry for payment routing decisions
// Requirements: 17.3
type RoutingLog struct {
	Base
	PaymentID          uint64  `json:"paymentId" gorm:"column:payment_id;not null;index"`                    // 支付记录ID
	OrderID            uint64  `json:"orderId" gorm:"column:order_id;not null;index"`                        // 订单ID
	MatchedRuleID      *uint64 `json:"matchedRuleId,omitempty" gorm:"column:matched_rule_id;index"`          // 匹配的规则ID（null表示使用默认主体）
	CollectionEntityID uint64  `json:"collectionEntityId" gorm:"column:collection_entity_id;not null;index"` // 最终使用的收款主体ID
	MerchantNo         string  `json:"merchantNo" gorm:"column:merchant_no;size:64;not null"`                // 使用的商户号
	IsDefault          bool    `json:"isDefault" gorm:"column:is_default;default:false"`                     // 是否使用默认主体
	IsFallback         bool    `json:"isFallback" gorm:"column:is_fallback;default:false"`                   // 是否使用备用主体
	MatchDetails       string  `json:"matchDetails,omitempty" gorm:"column:match_details;type:text"`         // 匹配详情（JSON）
	ErrorMessage       string  `json:"errorMessage,omitempty" gorm:"column:error_message;type:text"`         // 错误信息

	// Relations
	Payment          Payment          `json:"-" gorm:"foreignKey:PaymentID;references:ID"`
	MatchedRule      *RoutingRule     `json:"-" gorm:"foreignKey:MatchedRuleID;references:ID"`
	CollectionEntity CollectionEntity `json:"-" gorm:"foreignKey:CollectionEntityID;references:ID"`
}

// TableName specifies the table name for RoutingLog model
func (RoutingLog) TableName() string {
	return "routing_logs"
}

// Request/Response DTOs

// CreateCollectionEntityRequest represents a request to create a collection entity
type CreateCollectionEntityRequest struct {
	Name              string `json:"name" binding:"required,max=200"`
	CreditCode        string `json:"creditCode" binding:"required,len=18"`
	TaxRegistrationNo string `json:"taxRegistrationNo" binding:"max=20"`
	IsDefault         bool   `json:"isDefault"`
}

// UpdateCollectionEntityRequest represents a request to update a collection entity
type UpdateCollectionEntityRequest struct {
	Name              *string `json:"name" binding:"omitempty,max=200"`
	TaxRegistrationNo *string `json:"taxRegistrationNo" binding:"omitempty,max=20"`
}

// ConfigurePaymentChannelRequest represents a request to configure a payment channel
type ConfigurePaymentChannelRequest struct {
	Channel     PaymentMethod `json:"channel" binding:"required,oneof=wechat alipay"`
	MerchantNo  string        `json:"merchantNo" binding:"required,max=64"`
	MerchantKey string        `json:"merchantKey" binding:"required,max=256"`
	CallbackURL string        `json:"callbackUrl" binding:"max=500"`
	Enabled     bool          `json:"enabled"`
	Priority    int           `json:"priority"`
	Remark      string        `json:"remark" binding:"max=500"`
}

// CreateRoutingRuleRequest represents a request to create a routing rule
type CreateRoutingRuleRequest struct {
	Name           string             `json:"name" binding:"required,max=100"`
	Priority       int                `json:"priority" binding:"required,min=1"`
	Conditions     []RoutingCondition `json:"conditions" binding:"required,min=1"`
	TargetEntityID uint64             `json:"targetEntityId" binding:"required"`
	Description    string             `json:"description" binding:"max=500"`
}

// UpdateRoutingRuleRequest represents a request to update a routing rule
type UpdateRoutingRuleRequest struct {
	Name           *string             `json:"name" binding:"omitempty,max=100"`
	Priority       *int                `json:"priority" binding:"omitempty,min=1"`
	Conditions     *[]RoutingCondition `json:"conditions" binding:"omitempty,min=1"`
	TargetEntityID *uint64             `json:"targetEntityId"`
	Description    *string             `json:"description" binding:"omitempty,max=500"`
}

// ListCollectionEntitiesRequest represents a request to list collection entities
type ListCollectionEntitiesRequest struct {
	Page      int          `form:"page" binding:"omitempty,min=1"`
	PageSize  int          `form:"pageSize" binding:"omitempty,min=1,max=100"`
	Status    EntityStatus `form:"status" binding:"omitempty,oneof=active inactive"`
	Keyword   string       `form:"keyword" binding:"omitempty,max=100"`
	SortBy    string       `form:"sortBy" binding:"omitempty,oneof=name created_at total_collection_cents"`
	SortOrder string       `form:"sortOrder" binding:"omitempty,oneof=asc desc"`
}

// ListCollectionEntitiesResponse represents a response for listing collection entities
type ListCollectionEntitiesResponse struct {
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"pageSize"`
	Entities []CollectionEntity `json:"entities"`
}

// ListRoutingRulesRequest represents a request to list routing rules
type ListRoutingRulesRequest struct {
	Page           int        `form:"page" binding:"omitempty,min=1"`
	PageSize       int        `form:"pageSize" binding:"omitempty,min=1,max=100"`
	Status         RuleStatus `form:"status" binding:"omitempty,oneof=active inactive"`
	TargetEntityID *uint64    `form:"targetEntityId"`
	Keyword        string     `form:"keyword" binding:"omitempty,max=100"`
}

// ListRoutingRulesResponse represents a response for listing routing rules
type ListRoutingRulesResponse struct {
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"pageSize"`
	Rules    []RoutingRule `json:"rules"`
}

// RoutingTestRequest represents a request to test routing rules
type RoutingTestRequest struct {
	GameType    string `json:"gameType"`
	ServiceType string `json:"serviceType"`
	AmountCents int64  `json:"amountCents"`
	Region      string `json:"region"`
}

// RoutingTestResponse represents a response for routing test
type RoutingTestResponse struct {
	MatchedRuleID      *uint64            `json:"matchedRuleId,omitempty"`
	MatchedRuleName    string             `json:"matchedRuleName,omitempty"`
	CollectionEntityID uint64             `json:"collectionEntityId"`
	EntityName         string             `json:"entityName"`
	MerchantNo         string             `json:"merchantNo"`
	IsDefault          bool               `json:"isDefault"`
	MatchDetails       []RoutingCondition `json:"matchDetails,omitempty"`
}

// PaymentRoutingStats represents statistics for payment routing by entity
// Requirements: 18.1
type PaymentRoutingStats struct {
	CollectionEntityID  uint64  `json:"collectionEntityId"`
	EntityName          string  `json:"entityName"`
	TotalAmountCents    int64   `json:"totalAmountCents"`
	RefundedAmountCents int64   `json:"refundedAmountCents"`
	NetAmountCents      int64   `json:"netAmountCents"`
	TransactionCount    int64   `json:"transactionCount"`
	AverageAmountCents  int64   `json:"averageAmountCents"`
	Percentage          float64 `json:"percentage"` // 占比
}

// PaymentRoutingReport represents a monthly payment routing report
type PaymentRoutingReport struct {
	Period           string                `json:"period"` // YYYY-MM format
	TotalAmountCents int64                 `json:"totalAmountCents"`
	TotalRefundCents int64                 `json:"totalRefundCents"`
	NetAmountCents   int64                 `json:"netAmountCents"`
	TransactionCount int64                 `json:"transactionCount"`
	ByEntity         []PaymentRoutingStats `json:"byEntity"`
	GeneratedAt      time.Time             `json:"generatedAt"`
}
