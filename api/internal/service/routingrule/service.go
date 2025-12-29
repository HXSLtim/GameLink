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
	"strconv"
	"strings"
)

var (
	// ErrNotFound 分流规则不存在
	ErrNotFound = apierr.NotFound("routing rule not found")
	// ErrValidation 表示输入校验失败
	ErrValidation = apierr.BadRequest("validation failed")
	// ErrTargetEntityNotFound 目标收款主体不存在
	ErrTargetEntityNotFound = apierr.BadRequest("target collection entity not found")
	// ErrTargetEntityInactive 目标收款主体已禁用
	ErrTargetEntityInactive = apierr.BadRequest("target collection entity is inactive")
	// ErrNoDefaultEntity 没有默认收款主体
	ErrNoDefaultEntity = apierr.NotFound("no default collection entity found")
	// ErrInvalidCondition 无效的条件配置
	ErrInvalidCondition = apierr.BadRequest("invalid routing condition")
	// ErrDuplicatePriority 优先级重复
	ErrDuplicatePriority = apierr.Conflict("priority already exists")
)

// RoutingRuleService 分流规则服务
// Requirements: 16.1, 16.2, 16.3, 16.4, 16.5
type RoutingRuleService struct {
	ruleRepo   routingrule.RoutingRuleRepository
	entityRepo collectionentity.CollectionEntityRepository
}

// NewRoutingRuleService 创建分流规则服务
func NewRoutingRuleService(
	ruleRepo routingrule.RoutingRuleRepository,
	entityRepo collectionentity.CollectionEntityRepository,
) *RoutingRuleService {
	return &RoutingRuleService{
		ruleRepo:   ruleRepo,
		entityRepo: entityRepo,
	}
}

// CreateRule 创建分流规则
// Requirements: 16.1
func (s *RoutingRuleService) CreateRule(ctx context.Context, req *model.CreateRoutingRuleRequest, createdBy uint64) (*model.RoutingRule, error) {
	// 验证目标收款主体存在且活跃
	entity, err := s.entityRepo.Get(ctx, req.TargetEntityID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrTargetEntityNotFound.WithDetails(fmt.Sprintf("entity ID %d not found", req.TargetEntityID))
		}
		return nil, apierr.InternalError("failed to get target entity").WithDetails(err.Error())
	}
	if entity.Status != model.EntityStatusActive {
		return nil, ErrTargetEntityInactive.WithDetails("cannot create rule for inactive entity")
	}

	// 验证条件配置
	if err := s.validateConditions(req.Conditions); err != nil {
		return nil, err
	}

	// 序列化条件
	conditionsJSON, err := json.Marshal(req.Conditions)
	if err != nil {
		return nil, apierr.InternalError("failed to serialize conditions").WithDetails(err.Error())
	}

	// 创建规则
	rule := &model.RoutingRule{
		Name:           req.Name,
		Priority:       req.Priority,
		Conditions:     conditionsJSON,
		TargetEntityID: req.TargetEntityID,
		Status:         model.RuleStatusActive,
		Description:    req.Description,
		CreatedBy:      createdBy,
	}

	if err := s.ruleRepo.Create(ctx, rule); err != nil {
		return nil, apierr.InternalError("failed to create routing rule").WithDetails(err.Error())
	}

	// 重新加载以获取关联数据
	return s.ruleRepo.Get(ctx, rule.ID)
}

// validateConditions 验证条件配置
func (s *RoutingRuleService) validateConditions(conditions []model.RoutingCondition) error {
	if len(conditions) == 0 {
		return ErrInvalidCondition.WithDetails("at least one condition is required")
	}

	for i, cond := range conditions {
		// 验证字段
		switch cond.Field {
		case model.ConditionFieldGameType, model.ConditionFieldServiceType, model.ConditionFieldRegion:
			// 字符串类型字段
		case model.ConditionFieldOrderAmount:
			// 数值类型字段
		default:
			return ErrInvalidCondition.WithDetails(fmt.Sprintf("invalid field at condition %d: %s", i, cond.Field))
		}

		// 验证操作符
		switch cond.Operator {
		case model.ConditionOperatorEquals, model.ConditionOperatorNotEquals,
			model.ConditionOperatorIn, model.ConditionOperatorNotIn,
			model.ConditionOperatorGreaterThan, model.ConditionOperatorLessThan,
			model.ConditionOperatorBetween:
			// 有效操作符
		default:
			return ErrInvalidCondition.WithDetails(fmt.Sprintf("invalid operator at condition %d: %s", i, cond.Operator))
		}

		// 验证值不为空
		if len(cond.Value) == 0 {
			return ErrInvalidCondition.WithDetails(fmt.Sprintf("value is required at condition %d", i))
		}
	}

	return nil
}

// GetRule 获取分流规则
func (s *RoutingRuleService) GetRule(ctx context.Context, id uint64) (*model.RoutingRule, error) {
	rule, err := s.ruleRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, apierr.InternalError("failed to get routing rule").WithDetails(err.Error())
	}
	return rule, nil
}

// UpdateRule 更新分流规则
// Requirements: 16.1
func (s *RoutingRuleService) UpdateRule(ctx context.Context, id uint64, req *model.UpdateRoutingRuleRequest, updatedBy uint64) (*model.RoutingRule, error) {
	rule, err := s.ruleRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, apierr.InternalError("failed to get routing rule").WithDetails(err.Error())
	}

	// 记录修改历史
	changes := s.detectChanges(rule, req)

	// 应用更新
	if req.Name != nil {
		rule.Name = *req.Name
	}
	if req.Priority != nil {
		rule.Priority = *req.Priority
	}
	if req.Conditions != nil {
		if err := s.validateConditions(*req.Conditions); err != nil {
			return nil, err
		}
		conditionsJSON, err := json.Marshal(*req.Conditions)
		if err != nil {
			return nil, apierr.InternalError("failed to serialize conditions").WithDetails(err.Error())
		}
		rule.Conditions = conditionsJSON
	}
	if req.TargetEntityID != nil {
		// 验证新的目标收款主体
		entity, err := s.entityRepo.Get(ctx, *req.TargetEntityID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, ErrTargetEntityNotFound
			}
			return nil, apierr.InternalError("failed to get target entity").WithDetails(err.Error())
		}
		if entity.Status != model.EntityStatusActive {
			return nil, ErrTargetEntityInactive
		}
		rule.TargetEntityID = *req.TargetEntityID
	}
	if req.Description != nil {
		rule.Description = *req.Description
	}
	rule.UpdatedBy = &updatedBy

	if err := s.ruleRepo.Update(ctx, rule); err != nil {
		return nil, apierr.InternalError("failed to update routing rule").WithDetails(err.Error())
	}

	// 保存修改历史
	for _, change := range changes {
		change.RoutingRuleID = id
		change.ChangedBy = updatedBy
		if err := s.ruleRepo.CreateHistory(ctx, &change); err != nil {
			// 记录历史失败不影响主流程
			continue
		}
	}

	return s.ruleRepo.Get(ctx, id)
}

// detectChanges 检测字段变更
func (s *RoutingRuleService) detectChanges(rule *model.RoutingRule, req *model.UpdateRoutingRuleRequest) []model.RoutingRuleHistory {
	var changes []model.RoutingRuleHistory

	if req.Name != nil && *req.Name != rule.Name {
		changes = append(changes, model.RoutingRuleHistory{
			FieldName: "name",
			OldValue:  rule.Name,
			NewValue:  *req.Name,
		})
	}
	if req.Priority != nil && *req.Priority != rule.Priority {
		changes = append(changes, model.RoutingRuleHistory{
			FieldName: "priority",
			OldValue:  fmt.Sprintf("%d", rule.Priority),
			NewValue:  fmt.Sprintf("%d", *req.Priority),
		})
	}
	if req.TargetEntityID != nil && *req.TargetEntityID != rule.TargetEntityID {
		changes = append(changes, model.RoutingRuleHistory{
			FieldName: "target_entity_id",
			OldValue:  fmt.Sprintf("%d", rule.TargetEntityID),
			NewValue:  fmt.Sprintf("%d", *req.TargetEntityID),
		})
	}
	if req.Description != nil && *req.Description != rule.Description {
		changes = append(changes, model.RoutingRuleHistory{
			FieldName: "description",
			OldValue:  rule.Description,
			NewValue:  *req.Description,
		})
	}

	return changes
}

// DeleteRule 删除分流规则
// Requirements: 16.1
func (s *RoutingRuleService) DeleteRule(ctx context.Context, id uint64) error {
	_, err := s.ruleRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return apierr.InternalError("failed to get routing rule").WithDetails(err.Error())
	}

	if err := s.ruleRepo.Delete(ctx, id); err != nil {
		return apierr.InternalError("failed to delete routing rule").WithDetails(err.Error())
	}

	return nil
}

// ListRules 查询分流规则列表
// Requirements: 16.5
func (s *RoutingRuleService) ListRules(ctx context.Context, req *model.ListRoutingRulesRequest) (*model.ListRoutingRulesResponse, error) {
	opts := routingrule.ListOptions{
		Keyword:        req.Keyword,
		TargetEntityID: req.TargetEntityID,
		Page:           req.Page,
		PageSize:       req.PageSize,
	}

	if req.Status != "" {
		opts.Status = &req.Status
	}

	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.PageSize < 1 || opts.PageSize > 100 {
		opts.PageSize = 20
	}

	rules, total, err := s.ruleRepo.List(ctx, opts)
	if err != nil {
		return nil, apierr.InternalError("failed to list routing rules").WithDetails(err.Error())
	}

	return &model.ListRoutingRulesResponse{
		Total:    total,
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Rules:    rules,
	}, nil
}

// ToggleRuleStatus 切换分流规则状态
// Requirements: 16.4
func (s *RoutingRuleService) ToggleRuleStatus(ctx context.Context, id uint64, enabled bool, updatedBy uint64) error {
	rule, err := s.ruleRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return apierr.InternalError("failed to get routing rule").WithDetails(err.Error())
	}

	var newStatus model.RuleStatus
	if enabled {
		newStatus = model.RuleStatusActive
	} else {
		newStatus = model.RuleStatusInactive
	}

	// 记录状态变更历史
	if rule.Status != newStatus {
		history := &model.RoutingRuleHistory{
			RoutingRuleID: id,
			FieldName:     "status",
			OldValue:      string(rule.Status),
			NewValue:      string(newStatus),
			ChangedBy:     updatedBy,
		}
		if err := s.ruleRepo.CreateHistory(ctx, history); err != nil {
			// 记录历史失败不影响主流程
		}
	}

	if err := s.ruleRepo.ToggleStatus(ctx, id, newStatus); err != nil {
		return apierr.InternalError("failed to toggle routing rule status").WithDetails(err.Error())
	}

	return nil
}

// GetRuleHistory 获取分流规则修改历史
func (s *RoutingRuleService) GetRuleHistory(ctx context.Context, ruleID uint64) ([]model.RoutingRuleHistory, error) {
	// 先验证规则存在
	_, err := s.ruleRepo.Get(ctx, ruleID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, apierr.InternalError("failed to get routing rule").WithDetails(err.Error())
	}

	histories, err := s.ruleRepo.GetHistory(ctx, ruleID)
	if err != nil {
		return nil, apierr.InternalError("failed to get routing rule history").WithDetails(err.Error())
	}

	return histories, nil
}

// SetDefaultEntity 设置默认收款主体
// Requirements: 16.3
func (s *RoutingRuleService) SetDefaultEntity(ctx context.Context, entityID uint64, updatedBy uint64) error {
	entity, err := s.entityRepo.Get(ctx, entityID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrTargetEntityNotFound
		}
		return apierr.InternalError("failed to get collection entity").WithDetails(err.Error())
	}

	// 验证主体是活跃状态
	if entity.Status != model.EntityStatusActive {
		return ErrTargetEntityInactive.WithDetails("cannot set inactive entity as default")
	}

	if err := s.entityRepo.SetDefault(ctx, entityID); err != nil {
		return apierr.InternalError("failed to set default collection entity").WithDetails(err.Error())
	}

	return nil
}

// GetDefaultEntity 获取默认收款主体
// Requirements: 16.3, 17.4
// Property 16: 收款分流默认主体回退
func (s *RoutingRuleService) GetDefaultEntity(ctx context.Context) (*model.CollectionEntity, error) {
	entity, err := s.entityRepo.GetDefault(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNoDefaultEntity
		}
		return nil, apierr.InternalError("failed to get default collection entity").WithDetails(err.Error())
	}
	return entity, nil
}

// ListActiveRulesByPriority 按优先级顺序获取所有活跃规则
// Requirements: 16.2
// Property 15: 收款分流规则优先级
func (s *RoutingRuleService) ListActiveRulesByPriority(ctx context.Context) ([]model.RoutingRule, error) {
	rules, err := s.ruleRepo.ListActiveByPriority(ctx)
	if err != nil {
		return nil, apierr.InternalError("failed to list active routing rules").WithDetails(err.Error())
	}
	return rules, nil
}

// MatchCollectionEntity 匹配收款主体
// Requirements: 17.1, 17.2
// Property 15: 收款分流规则优先级 - 按优先级顺序匹配规则
// Property 16: 收款分流默认主体回退 - 无规则匹配时使用默认主体
func (s *RoutingRuleService) MatchCollectionEntity(ctx context.Context, req *model.RoutingTestRequest) (*model.RoutingTestResponse, error) {
	// 获取所有活跃规则（按优先级排序）
	rules, err := s.ruleRepo.ListActiveByPriority(ctx)
	if err != nil {
		return nil, apierr.InternalError("failed to list routing rules").WithDetails(err.Error())
	}

	// 按优先级顺序匹配规则
	for _, rule := range rules {
		matched, matchDetails, err := s.matchRule(&rule, req)
		if err != nil {
			continue // 匹配出错跳过该规则
		}
		if matched {
			// 获取目标收款主体
			entity, err := s.entityRepo.Get(ctx, rule.TargetEntityID)
			if err != nil || entity.Status != model.EntityStatusActive {
				continue // 主体不可用，继续匹配下一条规则
			}

			// 获取商户号
			merchantNo := ""
			if len(entity.PaymentChannels) > 0 {
				merchantNo = entity.PaymentChannels[0].MerchantNo
			}

			return &model.RoutingTestResponse{
				MatchedRuleID:      &rule.ID,
				MatchedRuleName:    rule.Name,
				CollectionEntityID: entity.ID,
				EntityName:         entity.Name,
				MerchantNo:         merchantNo,
				IsDefault:          false,
				MatchDetails:       matchDetails,
			}, nil
		}
	}

	// 没有规则匹配，使用默认主体
	// Property 16: 收款分流默认主体回退
	defaultEntity, err := s.entityRepo.GetDefault(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNoDefaultEntity.WithDetails("no matching rule and no default entity configured")
		}
		return nil, apierr.InternalError("failed to get default entity").WithDetails(err.Error())
	}

	merchantNo := ""
	if len(defaultEntity.PaymentChannels) > 0 {
		merchantNo = defaultEntity.PaymentChannels[0].MerchantNo
	}

	return &model.RoutingTestResponse{
		MatchedRuleID:      nil,
		MatchedRuleName:    "",
		CollectionEntityID: defaultEntity.ID,
		EntityName:         defaultEntity.Name,
		MerchantNo:         merchantNo,
		IsDefault:          true,
		MatchDetails:       nil,
	}, nil
}

// matchRule 匹配单条规则
func (s *RoutingRuleService) matchRule(rule *model.RoutingRule, req *model.RoutingTestRequest) (bool, []model.RoutingCondition, error) {
	conditions, err := rule.GetConditions()
	if err != nil {
		return false, nil, err
	}

	var matchedConditions []model.RoutingCondition

	// 所有条件都必须匹配
	for _, cond := range conditions {
		matched, err := s.matchCondition(&cond, req)
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
func (s *RoutingRuleService) matchCondition(cond *model.RoutingCondition, req *model.RoutingTestRequest) (bool, error) {
	var fieldValue interface{}

	// 获取字段值
	switch cond.Field {
	case model.ConditionFieldGameType:
		fieldValue = req.GameType
	case model.ConditionFieldServiceType:
		fieldValue = req.ServiceType
	case model.ConditionFieldOrderAmount:
		fieldValue = req.AmountCents
	case model.ConditionFieldRegion:
		fieldValue = req.Region
	default:
		return false, fmt.Errorf("unknown field: %s", cond.Field)
	}

	// 根据操作符进行匹配
	switch cond.Operator {
	case model.ConditionOperatorEquals:
		return s.matchEquals(fieldValue, cond.Value)
	case model.ConditionOperatorNotEquals:
		matched, err := s.matchEquals(fieldValue, cond.Value)
		return !matched, err
	case model.ConditionOperatorIn:
		return s.matchIn(fieldValue, cond.Value)
	case model.ConditionOperatorNotIn:
		matched, err := s.matchIn(fieldValue, cond.Value)
		return !matched, err
	case model.ConditionOperatorGreaterThan:
		return s.matchGreaterThan(fieldValue, cond.Value)
	case model.ConditionOperatorLessThan:
		return s.matchLessThan(fieldValue, cond.Value)
	case model.ConditionOperatorBetween:
		return s.matchBetween(fieldValue, cond.Value)
	default:
		return false, fmt.Errorf("unknown operator: %s", cond.Operator)
	}
}

// matchEquals 等于匹配
func (s *RoutingRuleService) matchEquals(fieldValue interface{}, condValue json.RawMessage) (bool, error) {
	switch v := fieldValue.(type) {
	case string:
		var expected string
		if err := json.Unmarshal(condValue, &expected); err != nil {
			return false, err
		}
		return strings.EqualFold(v, expected), nil
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
func (s *RoutingRuleService) matchIn(fieldValue interface{}, condValue json.RawMessage) (bool, error) {
	switch v := fieldValue.(type) {
	case string:
		var expected []string
		if err := json.Unmarshal(condValue, &expected); err != nil {
			return false, err
		}
		for _, e := range expected {
			if strings.EqualFold(v, e) {
				return true, nil
			}
		}
		return false, nil
	case int64:
		var expected []int64
		if err := json.Unmarshal(condValue, &expected); err != nil {
			return false, err
		}
		for _, e := range expected {
			if v == e {
				return true, nil
			}
		}
		return false, nil
	default:
		return false, fmt.Errorf("unsupported field type for in: %T", fieldValue)
	}
}

// matchGreaterThan 大于匹配
func (s *RoutingRuleService) matchGreaterThan(fieldValue interface{}, condValue json.RawMessage) (bool, error) {
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
func (s *RoutingRuleService) matchLessThan(fieldValue interface{}, condValue json.RawMessage) (bool, error) {
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
func (s *RoutingRuleService) matchBetween(fieldValue interface{}, condValue json.RawMessage) (bool, error) {
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

// TestRouting 测试分流规则
// Requirements: 16.5
func (s *RoutingRuleService) TestRouting(ctx context.Context, req *model.RoutingTestRequest) (*model.RoutingTestResponse, error) {
	return s.MatchCollectionEntity(ctx, req)
}

// CreateRoutingLog 创建分流日志
// Requirements: 17.3
func (s *RoutingRuleService) CreateRoutingLog(ctx context.Context, log *model.RoutingLog) error {
	if err := s.ruleRepo.CreateRoutingLog(ctx, log); err != nil {
		return apierr.InternalError("failed to create routing log").WithDetails(err.Error())
	}
	return nil
}

// GetRoutingLogByPayment 根据支付ID获取分流日志
func (s *RoutingRuleService) GetRoutingLogByPayment(ctx context.Context, paymentID uint64) (*model.RoutingLog, error) {
	log, err := s.ruleRepo.GetRoutingLogByPayment(ctx, paymentID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apierr.NotFound("routing log not found")
		}
		return nil, apierr.InternalError("failed to get routing log").WithDetails(err.Error())
	}
	return log, nil
}

// ListRoutingLogs 查询分流日志列表
func (s *RoutingRuleService) ListRoutingLogs(ctx context.Context, opts routingrule.RoutingLogListOptions) ([]model.RoutingLog, int64, error) {
	logs, total, err := s.ruleRepo.ListRoutingLogs(ctx, opts)
	if err != nil {
		return nil, 0, apierr.InternalError("failed to list routing logs").WithDetails(err.Error())
	}
	return logs, total, nil
}

// ReorderPriorities 重新排序规则优先级
// Requirements: 16.2
func (s *RoutingRuleService) ReorderPriorities(ctx context.Context, ruleIDs []uint64, updatedBy uint64) error {
	for i, ruleID := range ruleIDs {
		newPriority := i + 1
		rule, err := s.ruleRepo.Get(ctx, ruleID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				continue
			}
			return apierr.InternalError("failed to get routing rule").WithDetails(err.Error())
		}

		if rule.Priority != newPriority {
			// 记录历史
			history := &model.RoutingRuleHistory{
				RoutingRuleID: ruleID,
				FieldName:     "priority",
				OldValue:      strconv.Itoa(rule.Priority),
				NewValue:      strconv.Itoa(newPriority),
				ChangedBy:     updatedBy,
			}
			if err := s.ruleRepo.CreateHistory(ctx, history); err != nil {
				// 记录历史失败不影响主流程
			}

			rule.Priority = newPriority
			rule.UpdatedBy = &updatedBy
			if err := s.ruleRepo.Update(ctx, rule); err != nil {
				return apierr.InternalError("failed to update rule priority").WithDetails(err.Error())
			}
		}
	}
	return nil
}

// BatchOperationResponse 批量操作响应
type BatchOperationResponse struct {
	SuccessCount int              `json:"success_count"`
	FailedCount  int              `json:"failed_count"`
	TotalCount   int              `json:"total_count"`
	FailedItems  []BatchErrorItem `json:"failed_items,omitempty"`
	SuccessItems []uint64         `json:"success_items,omitempty"`
}

// BatchErrorItem 单个操作错误详情
type BatchErrorItem struct {
	ID      uint64 `json:"id"`
	Message string `json:"message"`
}

// BatchUpdateRuleStatus 批量更新分流规则状态
func (s *RoutingRuleService) BatchUpdateRuleStatus(ctx context.Context, ruleIDs []uint64, isActive bool, updatedBy uint64) (*BatchOperationResponse, error) {
	response := &BatchOperationResponse{
		TotalCount:   len(ruleIDs),
		SuccessItems: make([]uint64, 0),
		FailedItems:  make([]BatchErrorItem, 0),
	}

	for _, ruleID := range ruleIDs {
		err := s.ToggleRuleStatus(ctx, ruleID, isActive, updatedBy)
		if err != nil {
			response.FailedItems = append(response.FailedItems, BatchErrorItem{
				ID:      ruleID,
				Message: err.Error(),
			})
			response.FailedCount++
			continue
		}

		response.SuccessItems = append(response.SuccessItems, ruleID)
		response.SuccessCount++
	}

	return response, nil
}

// BatchDeleteRules 批量删除分流规则
func (s *RoutingRuleService) BatchDeleteRules(ctx context.Context, ruleIDs []uint64) (*BatchOperationResponse, error) {
	response := &BatchOperationResponse{
		TotalCount:   len(ruleIDs),
		SuccessItems: make([]uint64, 0),
		FailedItems:  make([]BatchErrorItem, 0),
	}

	for _, ruleID := range ruleIDs {
		err := s.DeleteRule(ctx, ruleID)
		if err != nil {
			response.FailedItems = append(response.FailedItems, BatchErrorItem{
				ID:      ruleID,
				Message: err.Error(),
			})
			response.FailedCount++
			continue
		}

		response.SuccessItems = append(response.SuccessItems, ruleID)
		response.SuccessCount++
	}

	return response, nil
}
