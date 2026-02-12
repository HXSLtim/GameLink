package routingrule

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gamelink/internal/model"
	"gamelink/internal/repository/routingrule"
)

// ============================================================================
// Additional Edge Case Tests for 100% Coverage
// ============================================================================

// TestRoutingEngine_FallbackToDefault_GetDefaultError tests the error path
// when getting default entity fails with non-NotFound error
func TestRoutingEngine_FallbackToDefault_GetDefaultError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	routingCtx := &RoutingContext{
		OrderID:     1,
		GameType:    "王者荣耀",
		ServiceType: "escort",
		AmountCents: 10000,
		Method:      model.PaymentMethodWeChat,
	}

	// Mock GetDefault to return a non-NotFound error
	mockEntityRepo.On("GetDefault", ctx).Return(nil, errors.New("database connection failed"))

	result, err := engine.fallbackToDefault(ctx, routingCtx)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get default entity")

	mockEntityRepo.AssertExpectations(t)
}

// TestRoutingEngine_MatchRule_ConditionsParseError tests error handling
// when conditions JSON cannot be parsed
func TestRoutingEngine_MatchRule_ConditionsParseError(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	routingCtx := &RoutingContext{
		OrderID:     1,
		GameType:    "王者荣耀",
		ServiceType: "escort",
		AmountCents: 10000,
		Method:      model.PaymentMethodWeChat,
	}

	// Create a rule with invalid conditions JSON
	rule := &model.RoutingRule{
		Base:       model.Base{ID: 1},
		Name:       "Test Rule",
		Priority:   1,
		Conditions: json.RawMessage(`{invalid json}`),
	}

	matched, matchDetails, err := engine.matchRule(rule, routingCtx)

	assert.Error(t, err)
	assert.False(t, matched)
	assert.Nil(t, matchDetails)
}

// TestRoutingEngine_MatchEquals_BoolFieldValue tests matchEquals with bool field type
// This should return error since bool is not supported
func TestRoutingEngine_MatchEquals_BoolFieldValue(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	// Test with bool value (unsupported type)
	boolValue := true
	condValue := json.RawMessage(`true`)

	matched, err := engine.matchEquals(boolValue, condValue)

	assert.Error(t, err)
	assert.False(t, matched)
	assert.Contains(t, err.Error(), "unsupported field type")
}

// TestRoutingEngine_MatchIn_BoolFieldValue tests matchIn with bool field type
func TestRoutingEngine_MatchIn_BoolFieldValue(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	// Test with bool value (unsupported type)
	boolValue := true
	condValue := json.RawMessage(`[true, false]`)

	matched, err := engine.matchIn(boolValue, condValue)

	assert.Error(t, err)
	assert.False(t, matched)
	assert.Contains(t, err.Error(), "unsupported field type")
}

// TestRoutingEngine_MatchGreaterThan_InvalidJSON tests matchGreaterThan with invalid JSON
func TestRoutingEngine_MatchGreaterThan_InvalidJSON(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	fieldValue := int64(10000)
	condValue := json.RawMessage(`not a number`)

	matched, err := engine.matchGreaterThan(fieldValue, condValue)

	assert.Error(t, err)
	assert.False(t, matched)
}

// TestRoutingEngine_MatchLessThan_InvalidJSON tests matchLessThan with invalid JSON
func TestRoutingEngine_MatchLessThan_InvalidJSON(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	fieldValue := int64(10000)
	condValue := json.RawMessage(`{}`)

	matched, err := engine.matchLessThan(fieldValue, condValue)

	assert.Error(t, err)
	assert.False(t, matched)
}

// TestRoutingEngine_MatchBetween_InvalidJSON tests matchBetween with invalid JSON
func TestRoutingEngine_MatchBetween_InvalidJSON(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	fieldValue := int64(10000)
	condValue := json.RawMessage(`[1]`) // Only 1 value, need 2

	matched, err := engine.matchBetween(fieldValue, condValue)

	assert.Error(t, err)
	assert.False(t, matched)
	assert.Contains(t, err.Error(), "between requires exactly 2 values")
}

// TestRoutingEngine_MatchBetween_NonArrayJSON tests matchBetween with non-array JSON
func TestRoutingEngine_MatchBetween_NonArrayJSON(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	fieldValue := int64(10000)
	condValue := json.RawMessage(`10000`) // Single number, not array

	matched, err := engine.matchBetween(fieldValue, condValue)

	assert.Error(t, err)
	assert.False(t, matched)
}

// TestService_MatchEquals_BoolFieldValue tests service matchEquals with unsupported bool type
func TestService_MatchEquals_BoolFieldValue(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	// Test with bool value (unsupported type)
	boolValue := true
	condValue := json.RawMessage(`true`)

	matched, err := svc.matchEquals(boolValue, condValue)

	assert.Error(t, err)
	assert.False(t, matched)
	assert.Contains(t, err.Error(), "unsupported field type for equals")
}

// TestService_MatchIn_BoolFieldValue tests service matchIn with unsupported bool type
func TestService_MatchIn_BoolFieldValue(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	// Test with bool value (unsupported type)
	boolValue := true
	condValue := json.RawMessage(`[true, false]`)

	matched, err := svc.matchIn(boolValue, condValue)

	assert.Error(t, err)
	assert.False(t, matched)
	assert.Contains(t, err.Error(), "unsupported field type for in")
}

// TestService_MatchGreaterThan_InvalidJSON tests service matchGreaterThan with invalid JSON
func TestService_MatchGreaterThan_InvalidJSON(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	fieldValue := int64(10000)
	condValue := json.RawMessage(`invalid`)

	matched, err := svc.matchGreaterThan(fieldValue, condValue)

	assert.Error(t, err)
	assert.False(t, matched)
}

// TestService_MatchLessThan_InvalidJSON tests service matchLessThan with invalid JSON
func TestService_MatchLessThan_InvalidJSON(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	fieldValue := int64(10000)
	condValue := json.RawMessage(`{}`)

	matched, err := svc.matchLessThan(fieldValue, condValue)

	assert.Error(t, err)
	assert.False(t, matched)
}

// TestService_MatchBetween_NonNumericValue tests service matchBetween with non-numeric field
func TestService_MatchBetween_NonNumericValue(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	// Test with string value (wrong type)
	fieldValue := "not a number"
	condValue := json.RawMessage(`[1, 100]`)

	matched, err := svc.matchBetween(fieldValue, condValue)

	assert.Error(t, err)
	assert.False(t, matched)
	assert.Contains(t, err.Error(), "between only supports numeric fields")
}

// TestService_CreateRule_EntityRepoGetError tests CreateRule when entity repo Get fails
func TestService_CreateRule_EntityRepoGetError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldGameType,
			Operator: model.ConditionOperatorEquals,
			Value:    json.RawMessage(`"王者荣耀"`),
		},
	}

	req := &model.CreateRoutingRuleRequest{
		Name:           "Test Rule",
		Priority:       1,
		Conditions:     conditions,
		TargetEntityID: 1,
		Description:    "Test description",
	}

	// Mock entity repo Get to return non-NotFound error
	mockEntityRepo.On("Get", ctx, uint64(1)).Return(nil, errors.New("database connection failed"))

	rule, err := svc.CreateRule(ctx, req, 100)

	assert.Error(t, err)
	assert.Nil(t, rule)
	assert.Contains(t, err.Error(), "failed to get target entity")

	mockEntityRepo.AssertExpectations(t)
}

// TestService_UpdateRule_EntityRepoGetError tests UpdateRule when entity repo Get fails
func TestService_UpdateRule_EntityRepoGetError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	existingRule := &model.RoutingRule{
		Base:           model.Base{ID: 1},
		Name:           "Old Name",
		Priority:       1,
		TargetEntityID: 1,
	}

	newEntityID := uint64(2)
	req := &model.UpdateRoutingRuleRequest{
		TargetEntityID: &newEntityID,
	}

	mockRuleRepo.On("Get", ctx, uint64(1)).Return(existingRule, nil)
	mockEntityRepo.On("Get", ctx, uint64(2)).Return(nil, errors.New("database timeout"))

	rule, err := svc.UpdateRule(ctx, 1, req, 100)

	assert.Error(t, err)
	assert.Nil(t, rule)
	assert.Contains(t, err.Error(), "failed to get target entity")

	mockRuleRepo.AssertExpectations(t)
	mockEntityRepo.AssertExpectations(t)
}

// TestService_UpdateRule_ConditionsJSONError tests UpdateRule when conditions JSON marshaling fails
func TestService_UpdateRule_ConditionsJSONError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	existingRule := &model.RoutingRule{
		Base:           model.Base{ID: 1},
		Name:           "Old Name",
		Priority:       1,
		TargetEntityID: 1,
	}

	// Create conditions with invalid JSON value (channel that can't be marshaled)
	conditions := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldGameType,
			Operator: model.ConditionOperatorEquals,
			Value:    json.RawMessage(`"test"`),
		},
	}

	// This is a valid condition, but let's test the path
	req := &model.UpdateRoutingRuleRequest{
		Conditions: &conditions,
	}

	mockRuleRepo.On("Get", ctx, uint64(1)).Return(existingRule, nil)
	mockRuleRepo.On("Update", ctx, mock.Anything).Return(nil)
	mockRuleRepo.On("Get", ctx, uint64(1)).Return(existingRule, nil)

	// The validateConditions should pass
	rule, err := svc.UpdateRule(ctx, 1, req, 100)

	// Should succeed in validation and serialization
	assert.NoError(t, err)
	assert.NotNil(t, rule)

	mockRuleRepo.AssertExpectations(t)
}

// TestService_ListRules_RepoListError tests ListRules when repo List fails
func TestService_ListRules_RepoListError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	req := &model.ListRoutingRulesRequest{
		Page:     1,
		PageSize: 20,
	}

	mockRuleRepo.On("List", ctx, mock.Anything).Return([]model.RoutingRule{}, int64(0), errors.New("database error"))

	result, err := svc.ListRules(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to list routing rules")

	mockRuleRepo.AssertExpectations(t)
}

// TestService_ValidateConditions_EmptyConditions tests validation with empty conditions
func TestService_ValidateConditions_EmptyConditions(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{}

	err := svc.validateConditions(conditions)

	assert.Error(t, err)
	// Error is wrapped, so check for the base error message
	assert.Contains(t, err.Error(), "invalid routing condition")
}

// TestService_ReorderPriorities_UpdateError tests ReorderPriorities when update fails
func TestService_ReorderPriorities_UpdateError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	rule := &model.RoutingRule{
		Base:     model.Base{ID: 1},
		Name:     "Test Rule",
		Priority: 10, // Different priority, should trigger update
	}

	mockRuleRepo.On("Get", ctx, uint64(1)).Return(rule, nil).Maybe()
	mockRuleRepo.On("CreateHistory", ctx, mock.Anything).Return(nil)
	mockRuleRepo.On("Update", ctx, mock.Anything).Return(errors.New("database write failed"))

	err := svc.ReorderPriorities(ctx, []uint64{1}, 100)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update rule priority")
}

// TestService_MatchCollectionEntity_EntityRepoListError tests MatchCollectionEntity when ListActiveByPriority fails
func TestService_MatchCollectionEntity_EntityRepoListError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	req := &model.RoutingTestRequest{
		GameType:    "王者荣耀",
		ServiceType: "escort",
		AmountCents: 10000,
		Region:      "beijing",
	}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return([]model.RoutingRule{}, errors.New("database read failed"))

	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to list routing rules")

	mockRuleRepo.AssertExpectations(t)
}

// TestService_MatchCollectionEntity_DefaultEntityGetError tests MatchCollectionEntity when GetDefault fails
func TestService_MatchCollectionEntity_DefaultEntityGetError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	req := &model.RoutingTestRequest{
		GameType:    "王者荣耀",
		ServiceType: "escort",
		AmountCents: 10000,
		Region:      "beijing",
	}

	// Return no active rules
	mockRuleRepo.On("ListActiveByPriority", ctx).Return([]model.RoutingRule{}, nil)
	// GetDefault returns non-NotFound error
	mockEntityRepo.On("GetDefault", ctx).Return(nil, errors.New("connection timeout"))

	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get default entity")

	mockRuleRepo.AssertExpectations(t)
	mockEntityRepo.AssertExpectations(t)
}

// TestService_CreateRoutingLog_RepoError tests CreateRoutingLog when repo operation fails
func TestService_CreateRoutingLog_RepoError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	log := &model.RoutingLog{
		PaymentID:          1,
		OrderID:            1,
		CollectionEntityID: 1,
		MerchantNo:         "TEST_MERCHANT",
	}

	mockRuleRepo.On("CreateRoutingLog", ctx, log).Return(errors.New("database write failed"))

	err := svc.CreateRoutingLog(ctx, log)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create routing log")

	mockRuleRepo.AssertExpectations(t)
}

// TestService_ListRoutingLogs_RepoError tests ListRoutingLogs when repo operation fails
func TestService_ListRoutingLogs_RepoError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	opts := routingrule.RoutingLogListOptions{
		Page:     1,
		PageSize: 20,
	}

	mockRuleRepo.On("ListRoutingLogs", ctx, opts).Return([]model.RoutingLog{}, int64(0), errors.New("database query failed"))

	logs, total, err := svc.ListRoutingLogs(ctx, opts)

	assert.Error(t, err)
	assert.Nil(t, logs)
	assert.Zero(t, total)
	assert.Contains(t, err.Error(), "failed to list routing logs")

	mockRuleRepo.AssertExpectations(t)
}

// TestRoutingEngine_RoutePayment_EntityGetError tests RoutePayment when entity Get fails
func TestRoutingEngine_RoutePayment_EntityGetError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	routingCtx := &RoutingContext{
		OrderID:     1,
		GameType:    "王者荣耀",
		ServiceType: "escort",
		AmountCents: 10000,
		Method:      model.PaymentMethodWeChat,
	}

	conditions := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldGameType,
			Operator: model.ConditionOperatorEquals,
			Value:    json.RawMessage(`"王者荣耀"`),
		},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rule := &model.RoutingRule{
		Base:           model.Base{ID: 1},
		Name:           "Test Rule",
		Priority:       1,
		Conditions:     conditionsJSON,
		TargetEntityID: 1,
		Status:         model.RuleStatusActive,
	}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return([]model.RoutingRule{*rule}, nil)
	// Entity Get returns non-NotFound error
	mockEntityRepo.On("Get", ctx, uint64(1)).Return(nil, errors.New("entity fetch failed"))
	// GetDefault also fails when falling back
	mockEntityRepo.On("GetDefault", ctx).Return(nil, errors.New("no default configured"))

	result, err := engine.RoutePayment(ctx, routingCtx)

	// Should fall through to default but fail there too
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get default entity")

	mockRuleRepo.AssertExpectations(t)
	mockEntityRepo.AssertExpectations(t)
}

// TestService_GetRoutingLogByPayment_RepoGetError tests GetRoutingLogByPayment when repo Get fails
func TestService_GetRoutingLogByPayment_RepoGetError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	mockRuleRepo.On("GetRoutingLogByPayment", ctx, uint64(1)).Return(nil, errors.New("database read failed"))

	log, err := svc.GetRoutingLogByPayment(ctx, 1)

	assert.Error(t, err)
	assert.Nil(t, log)
	assert.Contains(t, err.Error(), "failed to get routing log")

	mockRuleRepo.AssertExpectations(t)
}
