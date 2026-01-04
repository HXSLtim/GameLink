package routingrule

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gamelink/internal/model"
)

// ============================================================================
// Additional Coverage Tests
// These tests target the remaining uncovered lines to reach 100% coverage
// ============================================================================

// TestRoutingEngine_MatchEquals_StringCaseSensitivity tests exact string matching
// RoutingEngine uses case-sensitive matching (not EqualFold like Service)
func TestRoutingEngine_MatchEquals_StringCaseSensitivity(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	// Test case-sensitive matching
	fieldValue := "LOL"
	condValue := json.RawMessage(`"lol"`)

	matched, err := engine.matchEquals(fieldValue, condValue)

	assert.NoError(t, err)
	// Engine uses exact match (not EqualFold), so "LOL" != "lol"
	assert.False(t, matched)
}

// TestRoutingEngine_MatchEquals_StringMatch tests exact match
func TestRoutingEngine_MatchEquals_StringMatch(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	fieldValue := "LOL"
	condValue := json.RawMessage(`"LOL"`)

	matched, err := engine.matchEquals(fieldValue, condValue)

	assert.NoError(t, err)
	assert.True(t, matched)
}

// TestRoutingEngine_MatchEquals_Int64Match tests int64 matching
func TestRoutingEngine_MatchEquals_Int64Match(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	fieldValue := int64(10000)
	condValue := json.RawMessage(`10000`)

	matched, err := engine.matchEquals(fieldValue, condValue)

	assert.NoError(t, err)
	assert.True(t, matched)
}

// TestRoutingEngine_MatchIn_StringInList tests string in list
func TestRoutingEngine_MatchIn_StringInList(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	fieldValue := "LOL"
	condValue := json.RawMessage(`["LOL", "王者荣耀", "PUBG"]`)

	matched, err := engine.matchIn(fieldValue, condValue)

	assert.NoError(t, err)
	assert.True(t, matched)
}

// TestRoutingEngine_MatchIn_StringNotInList tests string not in list
func TestRoutingEngine_MatchIn_StringNotInList(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	fieldValue := "DOTA"
	condValue := json.RawMessage(`["LOL", "王者荣耀", "PUBG"]`)

	matched, err := engine.matchIn(fieldValue, condValue)

	assert.NoError(t, err)
	assert.False(t, matched)
}

// TestRoutingEngine_MatchIn_Int64InList tests int64 in list
func TestRoutingEngine_MatchIn_Int64InList(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	fieldValue := int64(10000)
	condValue := json.RawMessage(`[5000, 10000, 15000]`)

	matched, err := engine.matchIn(fieldValue, condValue)

	assert.NoError(t, err)
	assert.True(t, matched)
}

// TestRoutingEngine_MatchLessThan_True tests less than matching (true case)
func TestRoutingEngine_MatchLessThan_True(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	fieldValue := int64(5000)
	condValue := json.RawMessage(`10000`)

	matched, err := engine.matchLessThan(fieldValue, condValue)

	assert.NoError(t, err)
	assert.True(t, matched)
}

// TestRoutingEngine_MatchLessThan_False tests less than matching (false case)
func TestRoutingEngine_MatchLessThan_False(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	fieldValue := int64(15000)
	condValue := json.RawMessage(`10000`)

	matched, err := engine.matchLessThan(fieldValue, condValue)

	assert.NoError(t, err)
	assert.False(t, matched)
}

// TestRoutingEngine_MatchLessThan_Equal tests less than with equal values
func TestRoutingEngine_MatchLessThan_Equal(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	fieldValue := int64(10000)
	condValue := json.RawMessage(`10000`)

	matched, err := engine.matchLessThan(fieldValue, condValue)

	assert.NoError(t, err)
	assert.False(t, matched)
}

// TestRoutingEngine_MatchBetween_WithinRange tests between matching (within range)
func TestRoutingEngine_MatchBetween_WithinRange(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	fieldValue := int64(7500)
	condValue := json.RawMessage(`[5000, 10000]`)

	matched, err := engine.matchBetween(fieldValue, condValue)

	assert.NoError(t, err)
	assert.True(t, matched)
}

// TestRoutingEngine_MatchBetween_BelowRange tests below range
func TestRoutingEngine_MatchBetween_BelowRange(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	fieldValue := int64(3000)
	condValue := json.RawMessage(`[5000, 10000]`)

	matched, err := engine.matchBetween(fieldValue, condValue)

	assert.NoError(t, err)
	assert.False(t, matched)
}

// TestRoutingEngine_MatchBetween_AboveRange tests above range
func TestRoutingEngine_MatchBetween_AboveRange(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	fieldValue := int64(15000)
	condValue := json.RawMessage(`[5000, 10000]`)

	matched, err := engine.matchBetween(fieldValue, condValue)

	assert.NoError(t, err)
	assert.False(t, matched)
}

// TestRoutingEngine_MatchBetween_BoundaryLow tests lower boundary (inclusive)
func TestRoutingEngine_MatchBetween_BoundaryLow(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	fieldValue := int64(5000)
	condValue := json.RawMessage(`[5000, 10000]`)

	matched, err := engine.matchBetween(fieldValue, condValue)

	assert.NoError(t, err)
	assert.True(t, matched)
}

// TestRoutingEngine_MatchBetween_BoundaryHigh tests upper boundary (inclusive)
func TestRoutingEngine_MatchBetween_BoundaryHigh(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	fieldValue := int64(10000)
	condValue := json.RawMessage(`[5000, 10000]`)

	matched, err := engine.matchBetween(fieldValue, condValue)

	assert.NoError(t, err)
	assert.True(t, matched)
}

// TestService_CreateRule_CreateError tests CreateRule when rule creation fails
func TestService_CreateRule_CreateError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	entity := &model.CollectionEntity{
		Base:   model.Base{ID: 1},
		Name:   "Test Entity",
		Status: model.EntityStatusActive,
	}

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorEquals, Value: mustJSON("LOL")},
	}

	req := &model.CreateRoutingRuleRequest{
		Name:           "Test Rule",
		Priority:       1,
		Conditions:     conditions,
		TargetEntityID: 1,
		Description:    "Test description",
	}

	mockEntityRepo.On("Get", ctx, uint64(1)).Return(entity, nil)
	mockRuleRepo.On("Create", ctx, mock.AnythingOfType("*model.RoutingRule")).Return(errors.New("database write failed"))

	result, err := svc.CreateRule(ctx, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create routing rule")

	mockEntityRepo.AssertExpectations(t)
	mockRuleRepo.AssertExpectations(t)
}

// TestService_CreateRule_ReloadError tests CreateRule when reloading rule fails
func TestService_CreateRule_ReloadError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	entity := &model.CollectionEntity{
		Base:   model.Base{ID: 1},
		Name:   "Test Entity",
		Status: model.EntityStatusActive,
	}

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorEquals, Value: mustJSON("LOL")},
	}

	req := &model.CreateRoutingRuleRequest{
		Name:           "Test Rule",
		Priority:       1,
		Conditions:     conditions,
		TargetEntityID: 1,
		Description:    "Test description",
	}

	mockEntityRepo.On("Get", ctx, uint64(1)).Return(entity, nil)
	mockRuleRepo.On("Create", ctx, mock.AnythingOfType("*model.RoutingRule")).Return(nil)
	mockRuleRepo.On("Get", ctx, mock.AnythingOfType("uint64")).Return(nil, errors.New("failed to reload"))

	result, err := svc.CreateRule(ctx, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockEntityRepo.AssertExpectations(t)
	mockRuleRepo.AssertExpectations(t)
}

// TestService_UpdateRule_UpdateError tests UpdateRule when update fails
func TestService_UpdateRule_UpdateError(t *testing.T) {
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

	newName := "New Name"
	req := &model.UpdateRoutingRuleRequest{
		Name: &newName,
	}

	mockRuleRepo.On("Get", ctx, uint64(1)).Return(existingRule, nil)
	mockRuleRepo.On("Update", ctx, mock.Anything).Return(errors.New("database update failed"))

	rule, err := svc.UpdateRule(ctx, 1, req, 100)

	assert.Error(t, err)
	assert.Nil(t, rule)
	assert.Contains(t, err.Error(), "failed to update routing rule")

	mockRuleRepo.AssertExpectations(t)
}

// TestService_UpdateRule_ReloadError tests UpdateRule when reloading fails
func TestService_UpdateRule_ReloadError(t *testing.T) {
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

	newName := "New Name"
	req := &model.UpdateRoutingRuleRequest{
		Name: &newName,
	}

	mockRuleRepo.On("Get", ctx, uint64(1)).Return(existingRule, nil).Once()
	mockRuleRepo.On("CreateHistory", ctx, mock.Anything).Return(nil)
	mockRuleRepo.On("Update", ctx, mock.Anything).Return(nil)
	mockRuleRepo.On("Get", ctx, uint64(1)).Return(nil, errors.New("reload failed"))

	rule, err := svc.UpdateRule(ctx, 1, req, 100)

	assert.Error(t, err)
	assert.Nil(t, rule)

	mockRuleRepo.AssertExpectations(t)
}

// TestService_DeleteRule_GetError tests DeleteRule when Get fails with non-NotFound error
func TestService_DeleteRule_GetError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	mockRuleRepo.On("Get", ctx, uint64(1)).Return(nil, errors.New("database connection failed"))

	err := svc.DeleteRule(ctx, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get routing rule")

	mockRuleRepo.AssertExpectations(t)
}

// TestService_DeleteRule_DeleteError tests DeleteRule when Delete operation fails
func TestService_DeleteRule_DeleteError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	rule := &model.RoutingRule{Base: model.Base{ID: 1}}
	mockRuleRepo.On("Get", ctx, uint64(1)).Return(rule, nil)
	mockRuleRepo.On("Delete", ctx, uint64(1)).Return(errors.New("database delete failed"))

	err := svc.DeleteRule(ctx, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete routing rule")

	mockRuleRepo.AssertExpectations(t)
}

// TestService_ValidateConditions_InvalidField tests validation with invalid field
func TestService_ValidateConditions_InvalidField(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{
			Field:    model.ConditionField("invalid_field"),
			Operator: model.ConditionOperatorEquals,
			Value:    json.RawMessage(`"test"`),
		},
	}

	err := svc.validateConditions(conditions)

	assert.Error(t, err)
	// Error message is "invalid routing condition" with details
	assert.Contains(t, err.Error(), "invalid routing condition")
}

// TestService_ValidateConditions_InvalidOperator tests validation with invalid operator
func TestService_ValidateConditions_InvalidOperator(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldGameType,
			Operator: model.ConditionOperator("invalid_op"),
			Value:    json.RawMessage(`"test"`),
		},
	}

	err := svc.validateConditions(conditions)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid routing condition")
}

// TestService_ValidateConditions_EmptyValue tests validation with empty value
func TestService_ValidateConditions_EmptyValue(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldGameType,
			Operator: model.ConditionOperatorEquals,
			Value:    json.RawMessage(``),
		},
	}

	err := svc.validateConditions(conditions)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid routing condition")
}

// TestService_MatchCondition_UnknownField tests matchCondition with unknown field
func TestService_MatchCondition_UnknownField(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	cond := &model.RoutingCondition{
		Field:    model.ConditionField("unknown_field"),
		Operator: model.ConditionOperatorEquals,
		Value:    json.RawMessage(`"test"`),
	}

	req := &model.RoutingTestRequest{
		GameType:    "LOL",
		ServiceType: "escort",
		AmountCents: 10000,
		Region:      "beijing",
	}

	matched, err := svc.matchCondition(cond, req)

	assert.Error(t, err)
	assert.False(t, matched)
	assert.Contains(t, err.Error(), "unknown field")
}

// TestService_MatchCondition_UnknownOperator tests matchCondition with unknown operator
func TestService_MatchCondition_UnknownOperator(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	cond := &model.RoutingCondition{
		Field:    model.ConditionFieldGameType,
		Operator: model.ConditionOperator("unknown_op"),
		Value:    json.RawMessage(`"LOL"`),
	}

	req := &model.RoutingTestRequest{
		GameType:    "LOL",
		ServiceType: "escort",
		AmountCents: 10000,
		Region:      "beijing",
	}

	matched, err := svc.matchCondition(cond, req)

	assert.Error(t, err)
	assert.False(t, matched)
	assert.Contains(t, err.Error(), "unknown operator")
}

// TestService_MatchRule_ConditionsParseError tests service matchRule with parse error
func TestService_MatchRule_ConditionsParseError(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	rule := &model.RoutingRule{
		Base:       model.Base{ID: 1},
		Name:       "Test Rule",
		Priority:   1,
		Conditions: json.RawMessage(`{invalid json}`),
	}

	req := &model.RoutingTestRequest{
		GameType:    "LOL",
		ServiceType: "escort",
		AmountCents: 10000,
		Region:      "beijing",
	}

	matched, matchDetails, err := svc.matchRule(rule, req)

	assert.Error(t, err)
	assert.False(t, matched)
	assert.Nil(t, matchDetails)
}

// TestRoutingEngine_MatchRule_NoConditions tests matchRule with no conditions
func TestRoutingEngine_MatchRule_NoConditions(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	rule := &model.RoutingRule{
		Base:       model.Base{ID: 1},
		Name:       "Test Rule",
		Priority:   1,
		Conditions: json.RawMessage(`[]`), // Empty conditions
	}

	routingCtx := &RoutingContext{
		GameType:    "LOL",
		ServiceType: "escort",
		AmountCents: 10000,
		Method:      model.PaymentMethodWeChat,
	}

	matched, matchDetails, err := engine.matchRule(rule, routingCtx)

	assert.NoError(t, err)
	assert.False(t, matched) // Empty conditions should not match
	assert.Nil(t, matchDetails)
}

// TestRoutingEngine_FallbackToDefault_InactiveDefaultEntity tests fallback when default entity is inactive
func TestRoutingEngine_FallbackToDefault_InactiveDefaultEntity(t *testing.T) {
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

	entity := &model.CollectionEntity{
		Base:   model.Base{ID: 1},
		Status: model.EntityStatusInactive, // Inactive
		PaymentChannels: []model.PaymentChannelConfig{
			{
				Channel:    model.PaymentMethodWeChat,
				MerchantNo: "TEST_MERCHANT",
				Enabled:    true,
			},
		},
	}

	mockEntityRepo.On("GetDefault", ctx).Return(entity, nil)

	result, err := engine.fallbackToDefault(ctx, routingCtx)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "target collection entity is inactive")

	mockEntityRepo.AssertExpectations(t)
}

// TestRoutingEngine_FallbackToDefault_NoChannels tests fallback when no channels configured
func TestRoutingEngine_FallbackToDefault_NoChannels(t *testing.T) {
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

	entity := &model.CollectionEntity{
		Base:            model.Base{ID: 1},
		Status:          model.EntityStatusActive,
		PaymentChannels: []model.PaymentChannelConfig{}, // No channels
	}

	mockEntityRepo.On("GetDefault", ctx).Return(entity, nil)
	// Add mock for GetChannelByEntityAndMethod which will be called first
	mockEntityRepo.On("GetChannelByEntityAndMethod", ctx, uint64(1), model.PaymentMethodWeChat).Return(nil, errors.New("no channel"))

	result, err := engine.fallbackToDefault(ctx, routingCtx)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no payment channel configured")

	mockEntityRepo.AssertExpectations(t)
}

// TestRoutingEngine_GetMerchantNo_NoChannel tests getMerchantNo when channel not found
func TestRoutingEngine_GetMerchantNo_NoChannel(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	entity := &model.CollectionEntity{
		Base:   model.Base{ID: 1},
		Status: model.EntityStatusActive,
		PaymentChannels: []model.PaymentChannelConfig{
			{
				Channel:    model.PaymentMethodAlipay,
				MerchantNo: "ALIPAY_MERCHANT",
				Enabled:    true,
			},
		},
	}

	// Request WeChat but only Alipay configured
	mockEntityRepo.On("GetChannelByEntityAndMethod", ctx, uint64(1), model.PaymentMethodWeChat).Return(nil, errors.New("channel not found"))

	merchantNo, err := engine.getMerchantNo(ctx, entity, model.PaymentMethodWeChat)

	assert.Error(t, err)
	assert.Empty(t, merchantNo)

	mockEntityRepo.AssertExpectations(t)
}

// TestRoutingEngine_GetMerchantNo_ChannelDisabled tests getMerchantNo when channel is disabled
func TestRoutingEngine_GetMerchantNo_ChannelDisabled(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	entity := &model.CollectionEntity{
		Base:   model.Base{ID: 1},
		Status: model.EntityStatusActive,
	}

	channel := &model.PaymentChannelConfig{
		Base:        model.Base{ID: 1},
		Channel:     model.PaymentMethodWeChat,
		MerchantNo:  "TEST_MERCHANT",
		Enabled:     false, // Disabled
	}

	mockEntityRepo.On("GetChannelByEntityAndMethod", ctx, uint64(1), model.PaymentMethodWeChat).Return(channel, nil)

	merchantNo, err := engine.getMerchantNo(ctx, entity, model.PaymentMethodWeChat)

	assert.Error(t, err)
	assert.Empty(t, merchantNo)
	assert.Contains(t, err.Error(), "is disabled")

	mockEntityRepo.AssertExpectations(t)
}
