package routingrule

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/collectionentity"
	"gamelink/internal/repository/routingrule"
	"gamelink/pkg/apierr"
)

// RoutingEngine 分流规则匹配引擎
// Requirements: 17.1, 17.2, 17.4, 17.5
// Property 15: 收款分流规则优先级 - 按优先级顺序匹配规则
// Property 16: 收款分流默认主体回退 - 无规则匹配时使用默认主体
// Property 17: 收款分流记录完整性 - 支付记录包含收款主体和商户号
type RoutingEngine struct {
	ruleRepo   routingrule.RoutingRuleRepository
	entityRepo collectionentity.CollectionEntityRepository
}

// NewRoutingEngine 创建分流规则匹配引擎
func NewRoutingEngine(
	ruleRepo routingrule.RoutingRuleRepository,
	entityRepo collectionentity.CollectionEntityRepository,
) *RoutingEngine {
	return &RoutingEngine{
		ruleRepo:   ruleRepo,
		entityRepo: entityRepo,
	}
}

// RoutingContext 分流上下文，包含订单相关信息用于规则匹配
type RoutingContext struct {
	OrderID     uint64              `json:"orderId"`
	GameType    string              `json:"gameType"`
	ServiceType string              `json:"serviceType"`
	AmountCents int64               `json:"amountCents"`
	Region      string              `json:"region"`
	Method      model.PaymentMethod `json:"method"`
}

// RoutingResult 分流结果
type RoutingResult struct {
	CollectionEntityID uint64                   `json:"collectionEntityId"`
	EntityName         string                   `json:"entityName"`
	MerchantNo         string                   `json:"merchantNo"`
	MatchedRuleID      *uint64                  `json:"matchedRuleId,omitempty"`
	MatchedRuleName    string                   `json:"matchedRuleName,omitempty"`
	IsDefault          bool                     `json:"isDefault"`
	IsFallback         bool                     `json:"isFallback"`
	MatchDetails       []model.RoutingCondition `json:"matchDetails,omitempty"`
}

// RoutePayment 执行支付分流
// Requirements: 17.1, 17.2
// Property 15: 收款分流规则优先级 - 按优先级顺序匹配规则，使用第一个匹配的规则
// Property 16: 收款分流默认主体回退 - 无规则匹配时使用默认主体
func (e *RoutingEngine) RoutePayment(ctx context.Context, routingCtx *RoutingContext) (*RoutingResult, error) {
	// 获取所有活跃规则（按优先级排序）
	// Property 15: 按优先级顺序获取规则
	rules, err := e.ruleRepo.ListActiveByPriority(ctx)
	if err != nil {
		return nil, apierr.InternalError("failed to list routing rules").WithDetails(err.Error())
	}

	// 按优先级顺序匹配规则
	for _, rule := range rules {
		matched, matchDetails, err := e.matchRule(&rule, routingCtx)
		if err != nil {
			// 匹配出错跳过该规则，继续下一条
			continue
		}
		if matched {
			// 获取目标收款主体
			entity, err := e.entityRepo.Get(ctx, rule.TargetEntityID)
			if err != nil || entity.Status != model.EntityStatusActive {
				// 主体不可用，继续匹配下一条规则
				// Requirements: 17.5 - 主体不可用时切换备用主体
				continue
			}

			// 获取对应支付方式的商户号
			merchantNo, err := e.getMerchantNo(ctx, entity, routingCtx.Method)
			if err != nil {
				// 没有对应支付方式的配置，继续匹配下一条规则
				continue
			}

			return &RoutingResult{
				CollectionEntityID: entity.ID,
				EntityName:         entity.Name,
				MerchantNo:         merchantNo,
				MatchedRuleID:      &rule.ID,
				MatchedRuleName:    rule.Name,
				IsDefault:          false,
				IsFallback:         false,
				MatchDetails:       matchDetails,
			}, nil
		}
	}

	// 没有规则匹配，使用默认主体
	// Property 16: 收款分流默认主体回退
	// Requirements: 17.4 - 规则匹配失败使用默认主体
	return e.fallbackToDefault(ctx, routingCtx)
}

// fallbackToDefault 回退到默认收款主体
// Requirements: 17.4
// Property 16: 收款分流默认主体回退
func (e *RoutingEngine) fallbackToDefault(ctx context.Context, routingCtx *RoutingContext) (*RoutingResult, error) {
	defaultEntity, err := e.entityRepo.GetDefault(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNoDefaultEntity.WithDetails("no matching rule and no default entity configured")
		}
		return nil, apierr.InternalError("failed to get default entity").WithDetails(err.Error())
	}

	// 验证默认主体是活跃状态
	if defaultEntity.Status != model.EntityStatusActive {
		return nil, ErrTargetEntityInactive.WithDetails("default entity is inactive")
	}

	// 获取对应支付方式的商户号
	merchantNo, err := e.getMerchantNo(ctx, defaultEntity, routingCtx.Method)
	if err != nil {
		// 尝试获取任意可用的商户号
		merchantNo = e.getAnyMerchantNo(defaultEntity)
		if merchantNo == "" {
			return nil, apierr.InternalError("no payment channel configured for default entity")
		}
	}

	return &RoutingResult{
		CollectionEntityID: defaultEntity.ID,
		EntityName:         defaultEntity.Name,
		MerchantNo:         merchantNo,
		MatchedRuleID:      nil,
		MatchedRuleName:    "",
		IsDefault:          true,
		IsFallback:         true,
		MatchDetails:       nil,
	}, nil
}

// getMerchantNo 获取指定支付方式的商户号
// Requirements: 17.2
func (e *RoutingEngine) getMerchantNo(ctx context.Context, entity *model.CollectionEntity, method model.PaymentMethod) (string, error) {
	channel, err := e.entityRepo.GetChannelByEntityAndMethod(ctx, entity.ID, method)
	if err != nil {
		return "", err
	}
	if !channel.Enabled {
		return "", fmt.Errorf("payment channel %s is disabled", method)
	}
	return channel.MerchantNo, nil
}

// getAnyMerchantNo 获取任意可用的商户号
func (e *RoutingEngine) getAnyMerchantNo(entity *model.CollectionEntity) string {
	for _, channel := range entity.PaymentChannels {
		if channel.Enabled && channel.MerchantNo != "" {
			return channel.MerchantNo
		}
	}
	return ""
}

// matchRule 匹配单条规则
// Requirements: 17.1 - 支持多条件组合匹配
func (e *RoutingEngine) matchRule(rule *model.RoutingRule, routingCtx *RoutingContext) (bool, []model.RoutingCondition, error) {
	conditions, err := rule.GetConditions()
	if err != nil {
		return false, nil, err
	}

	if len(conditions) == 0 {
		// 没有条件的规则不匹配任何请求
		return false, nil, nil
	}

	var matchedConditions []model.RoutingCondition

	// 所有条件都必须匹配（AND逻辑）
	for _, cond := range conditions {
		matched, err := e.matchCondition(&cond, routingCtx)
		if err != nil {
			return false, nil, err
		}
		if !matched {
			return false, nil, nil
		}
		matchedConditions = append(matchedConditions, cond)
	}

	return true, matchedConditions, nil
}

// matchCondition 匹配单个条件
// Requirements: 17.1
func (e *RoutingEngine) matchCondition(cond *model.RoutingCondition, routingCtx *RoutingContext) (bool, error) {
	var fieldValue interface{}

	// 获取字段值
	switch cond.Field {
	case model.ConditionFieldGameType:
		fieldValue = routingCtx.GameType
	case model.ConditionFieldServiceType:
		fieldValue = routingCtx.ServiceType
	case model.ConditionFieldOrderAmount:
		fieldValue = routingCtx.AmountCents
	case model.ConditionFieldRegion:
		fieldValue = routingCtx.Region
	default:
		return false, fmt.Errorf("unknown field: %s", cond.Field)
	}

	// 根据操作符进行匹配
	return e.evaluateCondition(fieldValue, cond.Operator, cond.Value)
}

// evaluateCondition 评估条件
func (e *RoutingEngine) evaluateCondition(fieldValue interface{}, operator model.ConditionOperator, condValue json.RawMessage) (bool, error) {
	switch operator {
	case model.ConditionOperatorEquals:
		return e.matchEquals(fieldValue, condValue)
	case model.ConditionOperatorNotEquals:
		matched, err := e.matchEquals(fieldValue, condValue)
		return !matched, err
	case model.ConditionOperatorIn:
		return e.matchIn(fieldValue, condValue)
	case model.ConditionOperatorNotIn:
		matched, err := e.matchIn(fieldValue, condValue)
		return !matched, err
	case model.ConditionOperatorGreaterThan:
		return e.matchGreaterThan(fieldValue, condValue)
	case model.ConditionOperatorLessThan:
		return e.matchLessThan(fieldValue, condValue)
	case model.ConditionOperatorBetween:
		return e.matchBetween(fieldValue, condValue)
	default:
		return false, fmt.Errorf("unknown operator: %s", operator)
	}
}

// matchEquals 等于匹配
func (e *RoutingEngine) matchEquals(fieldValue interface{}, condValue json.RawMessage) (bool, error) {
	switch v := fieldValue.(type) {
	case string:
		var expected string
		if err := json.Unmarshal(condValue, &expected); err != nil {
			return false, err
		}
		return v == expected, nil
	case int64:
		var expected int64
		if err := json.Unmarshal(condValue, &expected); err != nil {
			return false, err
		}
		return v == expected, nil
	default:
		return false, fmt.Errorf("unsupported field type for equals: %T", fieldValue)
	}
}

// matchIn 包含匹配
func (e *RoutingEngine) matchIn(fieldValue interface{}, condValue json.RawMessage) (bool, error) {
	switch v := fieldValue.(type) {
	case string:
		var expected []string
		if err := json.Unmarshal(condValue, &expected); err != nil {
			return false, err
		}
		for _, exp := range expected {
			if v == exp {
				return true, nil
			}
		}
		return false, nil
	case int64:
		var expected []int64
		if err := json.Unmarshal(condValue, &expected); err != nil {
			return false, err
		}
		for _, exp := range expected {
			if v == exp {
				return true, nil
			}
		}
		return false, nil
	default:
		return false, fmt.Errorf("unsupported field type for in: %T", fieldValue)
	}
}

// matchGreaterThan 大于匹配
func (e *RoutingEngine) matchGreaterThan(fieldValue interface{}, condValue json.RawMessage) (bool, error) {
	v, ok := fieldValue.(int64)
	if !ok {
		return false, fmt.Errorf("greater than only supports numeric fields")
	}
	var expected int64
	if err := json.Unmarshal(condValue, &expected); err != nil {
		return false, err
	}
	return v > expected, nil
}

// matchLessThan 小于匹配
func (e *RoutingEngine) matchLessThan(fieldValue interface{}, condValue json.RawMessage) (bool, error) {
	v, ok := fieldValue.(int64)
	if !ok {
		return false, fmt.Errorf("less than only supports numeric fields")
	}
	var expected int64
	if err := json.Unmarshal(condValue, &expected); err != nil {
		return false, err
	}
	return v < expected, nil
}

// matchBetween 区间匹配
func (e *RoutingEngine) matchBetween(fieldValue interface{}, condValue json.RawMessage) (bool, error) {
	v, ok := fieldValue.(int64)
	if !ok {
		return false, fmt.Errorf("between only supports numeric fields")
	}
	var expected []int64
	if err := json.Unmarshal(condValue, &expected); err != nil {
		return false, err
	}
	if len(expected) != 2 {
		return false, fmt.Errorf("between requires exactly 2 values")
	}
	return v >= expected[0] && v <= expected[1], nil
}

// CreateRoutingLog 创建分流日志
// Requirements: 17.3
// Property 17: 收款分流记录完整性
func (e *RoutingEngine) CreateRoutingLog(ctx context.Context, paymentID, orderID uint64, result *RoutingResult) error {
	matchDetails := ""
	if len(result.MatchDetails) > 0 {
		detailsJSON, _ := json.Marshal(result.MatchDetails)
		matchDetails = string(detailsJSON)
	}

	log := &model.RoutingLog{
		PaymentID:          paymentID,
		OrderID:            orderID,
		MatchedRuleID:      result.MatchedRuleID,
		CollectionEntityID: result.CollectionEntityID,
		MerchantNo:         result.MerchantNo,
		IsDefault:          result.IsDefault,
		IsFallback:         result.IsFallback,
		MatchDetails:       matchDetails,
	}

	return e.ruleRepo.CreateRoutingLog(ctx, log)
}

// GetRoutingLogByPayment 根据支付ID获取分流日志
func (e *RoutingEngine) GetRoutingLogByPayment(ctx context.Context, paymentID uint64) (*model.RoutingLog, error) {
	return e.ruleRepo.GetRoutingLogByPayment(ctx, paymentID)
}

// RoutePaymentWithFallback 执行支付分流，带完整容错处理
// Requirements: 17.4, 17.5
// 当规则匹配失败时使用默认主体，当主体不可用时切换备用主体
func (e *RoutingEngine) RoutePaymentWithFallback(ctx context.Context, routingCtx *RoutingContext) (*RoutingResult, error) {
	// 首先尝试正常分流
	result, err := e.RoutePayment(ctx, routingCtx)
	if err == nil {
		return result, nil
	}

	// 如果正常分流失败，尝试获取任意可用的活跃主体
	// Requirements: 17.5 - 主体不可用时切换备用主体
	return e.findAnyAvailableEntity(ctx, routingCtx)
}

// findAnyAvailableEntity 查找任意可用的活跃收款主体
// Requirements: 17.5
func (e *RoutingEngine) findAnyAvailableEntity(ctx context.Context, routingCtx *RoutingContext) (*RoutingResult, error) {
	// 获取所有活跃的收款主体
	entities, err := e.entityRepo.ListActive(ctx)
	if err != nil {
		return nil, apierr.InternalError("failed to list active entities").WithDetails(err.Error())
	}

	// 遍历查找有可用支付渠道的主体
	for _, entity := range entities {
		merchantNo, err := e.getMerchantNo(ctx, &entity, routingCtx.Method)
		if err != nil {
			// 尝试获取任意可用的商户号
			merchantNo = e.getAnyMerchantNo(&entity)
		}
		if merchantNo != "" {
			return &RoutingResult{
				CollectionEntityID: entity.ID,
				EntityName:         entity.Name,
				MerchantNo:         merchantNo,
				MatchedRuleID:      nil,
				MatchedRuleName:    "",
				IsDefault:          entity.IsDefault,
				IsFallback:         true,
				MatchDetails:       nil,
			}, nil
		}
	}

	return nil, apierr.InternalError("no available collection entity found")
}

// ValidateRoutingConfiguration 验证分流配置是否完整
// 用于系统启动时或配置变更时检查
func (e *RoutingEngine) ValidateRoutingConfiguration(ctx context.Context) error {
	// 检查是否有默认收款主体
	defaultEntity, err := e.entityRepo.GetDefault(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return apierr.BadRequest("no default collection entity configured")
		}
		return apierr.InternalError("failed to check default entity").WithDetails(err.Error())
	}

	// 检查默认主体是否活跃
	if defaultEntity.Status != model.EntityStatusActive {
		return apierr.BadRequest("default collection entity is inactive")
	}

	// 检查默认主体是否有可用的支付渠道
	channels, err := e.entityRepo.ListChannelsByEntity(ctx, defaultEntity.ID)
	if err != nil {
		return apierr.InternalError("failed to check payment channels").WithDetails(err.Error())
	}

	hasEnabledChannel := false
	for _, channel := range channels {
		if channel.Enabled && channel.MerchantNo != "" {
			hasEnabledChannel = true
			break
		}
	}

	if !hasEnabledChannel {
		return apierr.BadRequest("default collection entity has no enabled payment channel")
	}

	return nil
}
