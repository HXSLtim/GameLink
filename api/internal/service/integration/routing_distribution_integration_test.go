// Package integration provides routing distribution integration tests for intelligent player allocation.
package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository/collectionentity"
	"gamelink/internal/repository/routingrule"
	rulingservice "gamelink/internal/service/routingrule"
)

// ============================================================================
// Routing Distribution Integration Tests
// Tests for intelligent routing system with rule matching and distribution
// ============================================================================

// TestRoutingDistribution_MatchByGameType tests routing by game type
func TestRoutingDistribution_MatchByGameType(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	service := rulingservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create collection entity for LOL games
	lolEntity := CreateTestCollectionEntity(t, db, "LOL Collection Entity")

	// Create rule: LOL game -> lolEntity
	conditions := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldGameType,
			Operator: model.ConditionOperatorEquals,
			Value:    json.RawMessage(`"王者荣耀"`),
		},
	}
	rule := CreateTestRoutingRuleWithConditions(t, db, lolEntity, 1, conditions)

	// Test routing request for LOL game
	req := &model.RoutingTestRequest{
		GameType:    "王者荣耀",
		ServiceType: "escort",
		AmountCents: 10000,
		Region:      "beijing",
	}

	result, err := service.MatchCollectionEntity(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, result.MatchedRuleID)
	assert.Equal(t, rule.ID, *result.MatchedRuleID)
	assert.Equal(t, lolEntity.ID, result.CollectionEntityID)
	assert.False(t, result.IsDefault)
}

// TestRoutingDistribution_MatchByServiceType tests routing by service type
func TestRoutingDistribution_MatchByServiceType(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	service := rulingservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create collection entity for training services
	trainingEntity := CreateTestCollectionEntity(t, db, "Training Collection Entity")

	// Create rule: service_type = training -> trainingEntity
	conditions := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldServiceType,
			Operator: model.ConditionOperatorEquals,
			Value:    json.RawMessage(`"training"`),
		},
	}
	CreateTestRoutingRuleWithConditions(t, db, trainingEntity, 1, conditions)

	// Set default entity
	defaultEntity := CreateTestCollectionEntity(t, db, "Default Entity")
	_ = entityRepo.SetDefault(ctx, defaultEntity.ID)

	// Test routing request for training service
	req := &model.RoutingTestRequest{
		GameType:    "王者荣耀",
		ServiceType: "training",
		AmountCents: 10000,
		Region:      "beijing",
	}

	result, err := service.MatchCollectionEntity(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, result.MatchedRuleID)
	assert.Equal(t, trainingEntity.ID, result.CollectionEntityID)
	assert.False(t, result.IsDefault)
}

// TestRoutingDistribution_MatchByOrderAmount tests routing by order amount
func TestRoutingDistribution_MatchByOrderAmount(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	service := rulingservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create collection entity for high value orders
	highValueEntity := CreateTestCollectionEntity(t, db, "High Value Entity")

	// Create rule: amount > 50000 -> highValueEntity
	conditions := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldOrderAmount,
			Operator: model.ConditionOperatorGreaterThan,
			Value:    json.RawMessage(`50000`),
		},
	}
	CreateTestRoutingRuleWithConditions(t, db, highValueEntity, 1, conditions)

	// Set default entity
	defaultEntity := CreateTestCollectionEntity(t, db, "Default Entity")
	_ = entityRepo.SetDefault(ctx, defaultEntity.ID)

	// Test routing request for high value order (60000 cents = 600 yuan)
	req := &model.RoutingTestRequest{
		GameType:    "王者荣耀",
		ServiceType: "escort",
		AmountCents: 60000,
		Region:      "beijing",
	}

	result, err := service.MatchCollectionEntity(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, result.MatchedRuleID)
	assert.Equal(t, highValueEntity.ID, result.CollectionEntityID)
	assert.False(t, result.IsDefault)
}

// TestRoutingDistribution_MatchByRegion tests routing by region
func TestRoutingDistribution_MatchByRegion(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	service := rulingservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create collection entity for Beijing region
	beijingEntity := CreateTestCollectionEntity(t, db, "Beijing Entity")

	// Create rule: region = beijing -> beijingEntity
	conditions := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldRegion,
			Operator: model.ConditionOperatorEquals,
			Value:    json.RawMessage(`"beijing"`),
		},
	}
	CreateTestRoutingRuleWithConditions(t, db, beijingEntity, 1, conditions)

	// Set default entity
	defaultEntity := CreateTestCollectionEntity(t, db, "Default Entity")
	_ = entityRepo.SetDefault(ctx, defaultEntity.ID)

	// Test routing request for Beijing region
	req := &model.RoutingTestRequest{
		GameType:    "王者荣耀",
		ServiceType: "escort",
		AmountCents: 10000,
		Region:      "beijing",
	}

	result, err := service.MatchCollectionEntity(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, result.MatchedRuleID)
	assert.Equal(t, beijingEntity.ID, result.CollectionEntityID)
	assert.False(t, result.IsDefault)
}

// TestRoutingDistribution_MultipleConditions_AND tests multiple conditions with AND logic
func TestRoutingDistribution_MultipleConditions_AND(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	service := rulingservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create collection entity for LOL + training + Beijing
	specialEntity := CreateTestCollectionEntity(t, db, "Special Entity")

	// Create rule: game_type = LOL AND service_type = training AND region = beijing
	conditions := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldGameType,
			Operator: model.ConditionOperatorEquals,
			Value:    json.RawMessage(`"王者荣耀"`),
		},
		{
			Field:    model.ConditionFieldServiceType,
			Operator: model.ConditionOperatorEquals,
			Value:    json.RawMessage(`"training"`),
		},
		{
			Field:    model.ConditionFieldRegion,
			Operator: model.ConditionOperatorEquals,
			Value:    json.RawMessage(`"beijing"`),
		},
	}
	CreateTestRoutingRuleWithConditions(t, db, specialEntity, 1, conditions)

	// Set default entity
	defaultEntity := CreateTestCollectionEntity(t, db, "Default Entity")
	_ = entityRepo.SetDefault(ctx, defaultEntity.ID)

	// Test routing request matching all conditions
	req := &model.RoutingTestRequest{
		GameType:    "王者荣耀",
		ServiceType: "training",
		AmountCents: 10000,
		Region:      "beijing",
	}

	result, err := service.MatchCollectionEntity(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, result.MatchedRuleID)
	assert.Equal(t, specialEntity.ID, result.CollectionEntityID)
	assert.False(t, result.IsDefault)
}

// TestRoutingDistribution_MultipleConditions_PartialMatch tests that all conditions must match
func TestRoutingDistribution_MultipleConditions_PartialMatch(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	service := rulingservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create collection entity
	specialEntity := CreateTestCollectionEntity(t, db, "Special Entity")

	// Create rule with multiple conditions
	conditions := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldGameType,
			Operator: model.ConditionOperatorEquals,
			Value:    json.RawMessage(`"王者荣耀"`),
		},
		{
			Field:    model.ConditionFieldServiceType,
			Operator: model.ConditionOperatorEquals,
			Value:    json.RawMessage(`"training"`),
		},
	}
	CreateTestRoutingRuleWithConditions(t, db, specialEntity, 1, conditions)

	// Set default entity
	defaultEntity := CreateTestCollectionEntity(t, db, "Default Entity")
	_ = entityRepo.SetDefault(ctx, defaultEntity.ID)

	// Test routing request with partial match (only game_type matches)
	req := &model.RoutingTestRequest{
		GameType:    "王者荣耀",
		ServiceType: "escort", // Different from rule
		AmountCents: 10000,
		Region:      "beijing",
	}

	result, err := service.MatchCollectionEntity(ctx, req)
	require.NoError(t, err)
	// Should use default since not all conditions match
	assert.Nil(t, result.MatchedRuleID)
	assert.Equal(t, defaultEntity.ID, result.CollectionEntityID)
	assert.True(t, result.IsDefault)
}

// TestRoutingDistribution_PriorityOrder tests that higher priority rules match first
func TestRoutingDistribution_PriorityOrder(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	service := rulingservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create two entities
	highPriorityEntity := CreateTestCollectionEntity(t, db, "High Priority Entity")
	lowPriorityEntity := CreateTestCollectionEntity(t, db, "Low Priority Entity")

	// Create high priority rule (priority = 1)
	conditionsHigh := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldGameType,
			Operator: model.ConditionOperatorEquals,
			Value:    json.RawMessage(`"王者荣耀"`),
		},
	}
	CreateTestRoutingRuleWithConditions(t, db, highPriorityEntity, 1, conditionsHigh)

	// Create low priority rule (priority = 10) with same condition
	conditionsLow := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldGameType,
			Operator: model.ConditionOperatorEquals,
			Value:    json.RawMessage(`"王者荣耀"`),
		},
	}
	CreateTestRoutingRuleWithConditions(t, db, lowPriorityEntity, 10, conditionsLow)

	// Test routing request
	req := &model.RoutingTestRequest{
		GameType:    "王者荣耀",
		ServiceType: "escort",
		AmountCents: 10000,
		Region:      "beijing",
	}

	result, err := service.MatchCollectionEntity(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, result.MatchedRuleID)
	// Should match high priority rule
	assert.Equal(t, highPriorityEntity.ID, result.CollectionEntityID)
	assert.False(t, result.IsDefault)
}

// TestRoutingDistribution_PriorityStop tests that routing stops after first match
func TestRoutingDistribution_PriorityStop(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	service := rulingservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create two entities
	firstEntity := CreateTestCollectionEntity(t, db, "First Match Entity")
	secondEntity := CreateTestCollectionEntity(t, db, "Second Match Entity")

	// Create high priority rule for LOL
	conditionsFirst := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldGameType,
			Operator: model.ConditionOperatorEquals,
			Value:    json.RawMessage(`"王者荣耀"`),
		},
	}
	CreateTestRoutingRuleWithConditions(t, db, firstEntity, 1, conditionsFirst)

	// Create low priority rule for LOL (should never match)
	conditionsSecond := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldGameType,
			Operator: model.ConditionOperatorEquals,
			Value:    json.RawMessage(`"王者荣耀"`),
		},
	}
	CreateTestRoutingRuleWithConditions(t, db, secondEntity, 10, conditionsSecond)

	// Test routing request
	req := &model.RoutingTestRequest{
		GameType:    "王者荣耀",
		ServiceType: "escort",
		AmountCents: 10000,
		Region:      "beijing",
	}

	result, err := service.MatchCollectionEntity(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, result.MatchedRuleID)
	// Should match first rule and stop
	assert.Equal(t, firstEntity.ID, result.CollectionEntityID)
	assert.NotEqual(t, secondEntity.ID, result.CollectionEntityID)
}

// TestRoutingDistribution_NoMatch_Default tests default entity when no rules match
func TestRoutingDistribution_NoMatch_Default(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	service := rulingservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create entity with rule for LOL only
	lolEntity := CreateTestCollectionEntity(t, db, "LOL Entity")
	conditions := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldGameType,
			Operator: model.ConditionOperatorEquals,
			Value:    json.RawMessage(`"王者荣耀"`),
		},
	}
	CreateTestRoutingRuleWithConditions(t, db, lolEntity, 1, conditions)

	// Set default entity
	defaultEntity := CreateTestCollectionEntity(t, db, "Default Entity")
	err := entityRepo.SetDefault(ctx, defaultEntity.ID)
	require.NoError(t, err)

	// Test routing request for different game (no rule match)
	req := &model.RoutingTestRequest{
		GameType:    "英雄联盟",
		ServiceType: "escort",
		AmountCents: 10000,
		Region:      "beijing",
	}

	result, err := service.MatchCollectionEntity(ctx, req)
	require.NoError(t, err)
	// Should use default entity
	assert.Nil(t, result.MatchedRuleID)
	assert.Equal(t, defaultEntity.ID, result.CollectionEntityID)
	assert.True(t, result.IsDefault)
}

// TestRoutingDistribution_NoMatch_NoDefault tests error when no rules match and no default
func TestRoutingDistribution_NoMatch_NoDefault(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	service := rulingservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create entity with rule for LOL only
	lolEntity := CreateTestCollectionEntity(t, db, "LOL Entity")
	conditions := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldGameType,
			Operator: model.ConditionOperatorEquals,
			Value:    json.RawMessage(`"王者荣耀"`),
		},
	}
	CreateTestRoutingRuleWithConditions(t, db, lolEntity, 1, conditions)

	// No default entity set

	// Test routing request for different game (no rule match, no default)
	req := &model.RoutingTestRequest{
		GameType:    "英雄联盟",
		ServiceType: "escort",
		AmountCents: 10000,
		Region:      "beijing",
	}

	_, err := service.MatchCollectionEntity(ctx, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no default")
}

// TestRoutingDistribution_OperatorIn tests the In operator
func TestRoutingDistribution_OperatorIn(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	service := rulingservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create entity
	mobileGameEntity := CreateTestCollectionEntity(t, db, "Mobile Game Entity")

	// Create rule: game_type in ["王者荣耀", "英雄联盟", "和平精英"]
	conditions := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldGameType,
			Operator: model.ConditionOperatorIn,
			Value:    json.RawMessage(`["王者荣耀","英雄联盟","和平精英"]`),
		},
	}
	CreateTestRoutingRuleWithConditions(t, db, mobileGameEntity, 1, conditions)

	// Set default entity
	defaultEntity := CreateTestCollectionEntity(t, db, "Default Entity")
	_ = entityRepo.SetDefault(ctx, defaultEntity.ID)

	// Test with game in list
	req := &model.RoutingTestRequest{
		GameType:    "和平精英",
		ServiceType: "escort",
		AmountCents: 10000,
		Region:      "beijing",
	}

	result, err := service.MatchCollectionEntity(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, result.MatchedRuleID)
	assert.Equal(t, mobileGameEntity.ID, result.CollectionEntityID)
	assert.False(t, result.IsDefault)
}

// TestRoutingDistribution_OperatorNotIn tests the NotIn operator
func TestRoutingDistribution_OperatorNotIn(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	service := rulingservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create entity for non-LOL games
	otherEntity := CreateTestCollectionEntity(t, db, "Other Game Entity")

	// Create rule: game_type not in ["王者荣耀"]
	conditions := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldGameType,
			Operator: model.ConditionOperatorNotIn,
			Value:    json.RawMessage(`["王者荣耀"]`),
		},
	}
	CreateTestRoutingRuleWithConditions(t, db, otherEntity, 1, conditions)

	// Set default entity
	defaultEntity := CreateTestCollectionEntity(t, db, "Default Entity")
	_ = entityRepo.SetDefault(ctx, defaultEntity.ID)

	// Test with game not in exclusion list
	req := &model.RoutingTestRequest{
		GameType:    "绝地求生",
		ServiceType: "escort",
		AmountCents: 10000,
		Region:      "beijing",
	}

	result, err := service.MatchCollectionEntity(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, result.MatchedRuleID)
	assert.Equal(t, otherEntity.ID, result.CollectionEntityID)
	assert.False(t, result.IsDefault)
}

// TestRoutingDistribution_OperatorNotEquals tests the NotEquals operator
func TestRoutingDistribution_OperatorNotEquals(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	service := rulingservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create entity for non-escort services
	trainingEntity := CreateTestCollectionEntity(t, db, "Training Entity")

	// Create rule: service_type != "escort"
	conditions := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldServiceType,
			Operator: model.ConditionOperatorNotEquals,
			Value:    json.RawMessage(`"escort"`),
		},
	}
	CreateTestRoutingRuleWithConditions(t, db, trainingEntity, 1, conditions)

	// Set default entity
	defaultEntity := CreateTestCollectionEntity(t, db, "Default Entity")
	_ = entityRepo.SetDefault(ctx, defaultEntity.ID)

	// Test with training service
	req := &model.RoutingTestRequest{
		GameType:    "王者荣耀",
		ServiceType: "training",
		AmountCents: 10000,
		Region:      "beijing",
	}

	result, err := service.MatchCollectionEntity(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, result.MatchedRuleID)
	assert.Equal(t, trainingEntity.ID, result.CollectionEntityID)
	assert.False(t, result.IsDefault)
}

// TestRoutingDistribution_OperatorLessThan tests the LessThan operator
func TestRoutingDistribution_OperatorLessThan(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	service := rulingservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create entity for small orders
	smallOrderEntity := CreateTestCollectionEntity(t, db, "Small Order Entity")

	// Create rule: amount < 30000
	conditions := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldOrderAmount,
			Operator: model.ConditionOperatorLessThan,
			Value:    json.RawMessage(`30000`),
		},
	}
	CreateTestRoutingRuleWithConditions(t, db, smallOrderEntity, 1, conditions)

	// Set default entity
	defaultEntity := CreateTestCollectionEntity(t, db, "Default Entity")
	_ = entityRepo.SetDefault(ctx, defaultEntity.ID)

	// Test with small amount
	req := &model.RoutingTestRequest{
		GameType:    "王者荣耀",
		ServiceType: "escort",
		AmountCents: 20000,
		Region:      "beijing",
	}

	result, err := service.MatchCollectionEntity(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, result.MatchedRuleID)
	assert.Equal(t, smallOrderEntity.ID, result.CollectionEntityID)
	assert.False(t, result.IsDefault)
}

// TestRoutingDistribution_OperatorBetween tests the Between operator
func TestRoutingDistribution_OperatorBetween(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	service := rulingservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create entity for medium orders
	mediumOrderEntity := CreateTestCollectionEntity(t, db, "Medium Order Entity")

	// Create rule: amount between [20000, 50000]
	conditions := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldOrderAmount,
			Operator: model.ConditionOperatorBetween,
			Value:    json.RawMessage(`[20000,50000]`),
		},
	}
	CreateTestRoutingRuleWithConditions(t, db, mediumOrderEntity, 1, conditions)

	// Set default entity
	defaultEntity := CreateTestCollectionEntity(t, db, "Default Entity")
	_ = entityRepo.SetDefault(ctx, defaultEntity.ID)

	// Test with amount in range
	req := &model.RoutingTestRequest{
		GameType:    "王者荣耀",
		ServiceType: "escort",
		AmountCents: 35000,
		Region:      "beijing",
	}

	result, err := service.MatchCollectionEntity(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, result.MatchedRuleID)
	assert.Equal(t, mediumOrderEntity.ID, result.CollectionEntityID)
	assert.False(t, result.IsDefault)
}

// TestRoutingDistribution_InactiveRule tests that inactive rules are skipped
func TestRoutingDistribution_InactiveRule(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	service := rulingservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create entities
	inactiveEntity := CreateTestCollectionEntity(t, db, "Inactive Rule Entity")
	defaultEntity := CreateTestCollectionEntity(t, db, "Default Entity")

	// Create inactive rule
	conditions := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldGameType,
			Operator: model.ConditionOperatorEquals,
			Value:    json.RawMessage(`"王者荣耀"`),
		},
	}
	rule := CreateTestRoutingRuleWithConditions(t, db, inactiveEntity, 1, conditions)
	rule.Status = model.RuleStatusInactive
	db.Save(rule)

	// Set default entity
	_ = entityRepo.SetDefault(ctx, defaultEntity.ID)

	// Test routing request
	req := &model.RoutingTestRequest{
		GameType:    "王者荣耀",
		ServiceType: "escort",
		AmountCents: 10000,
		Region:      "beijing",
	}

	result, err := service.MatchCollectionEntity(ctx, req)
	require.NoError(t, err)
	// Should skip inactive rule and use default
	assert.Nil(t, result.MatchedRuleID)
	assert.Equal(t, defaultEntity.ID, result.CollectionEntityID)
	assert.True(t, result.IsDefault)
}

// TestRoutingDistribution_InactiveTargetEntity tests that inactive target entities are skipped
func TestRoutingDistribution_InactiveTargetEntity(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	service := rulingservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create inactive entity
	inactiveEntity := CreateTestCollectionEntity(t, db, "Inactive Entity")
	inactiveEntity.Status = model.EntityStatusInactive
	db.Save(inactiveEntity)

	// Create rule pointing to inactive entity
	conditions := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldGameType,
			Operator: model.ConditionOperatorEquals,
			Value:    json.RawMessage(`"王者荣耀"`),
		},
	}
	CreateTestRoutingRuleWithConditions(t, db, inactiveEntity, 1, conditions)

	// Set default entity
	defaultEntity := CreateTestCollectionEntity(t, db, "Default Entity")
	_ = entityRepo.SetDefault(ctx, defaultEntity.ID)

	// Test routing request
	req := &model.RoutingTestRequest{
		GameType:    "王者荣耀",
		ServiceType: "escort",
		AmountCents: 10000,
		Region:      "beijing",
	}

	result, err := service.MatchCollectionEntity(ctx, req)
	require.NoError(t, err)
	// Should skip rule with inactive entity and use default
	assert.Nil(t, result.MatchedRuleID)
	assert.Equal(t, defaultEntity.ID, result.CollectionEntityID)
	assert.True(t, result.IsDefault)
}

// TestRoutingDistribution_CreateRoutingLog tests creating routing log
func TestRoutingDistribution_CreateRoutingLog(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	service := rulingservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create test data
	user := CreateUniqueTestUser(t, db, "routing_log_user")
	playerUser := CreateUniqueTestUser(t, db, "routing_log_player")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "LOL")
	order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusPending, 10000)
	payment := CreateTestPayment(t, db, order, model.PaymentStatusPaid)

	entity := CreateTestCollectionEntity(t, db, "Log Test Entity")
	rule := CreateTestRoutingRule(t, db, entity, 1)

	// Create routing log
	log := &model.RoutingLog{
		PaymentID:          payment.ID,
		OrderID:            order.ID,
		MatchedRuleID:      &rule.ID,
		CollectionEntityID: entity.ID,
		MerchantNo:         "TEST_MERCHANT_001",
		IsDefault:          false,
		IsFallback:         false,
		MatchDetails:       `{"game":"王者荣耀","amount":10000}`,
	}

	err := service.CreateRoutingLog(ctx, log)
	require.NoError(t, err)
	assert.NotZero(t, log.ID)

	// Verify log was created
	retrieved, err := service.GetRoutingLogByPayment(ctx, payment.ID)
	require.NoError(t, err)
	assert.Equal(t, payment.ID, retrieved.PaymentID)
	assert.Equal(t, order.ID, retrieved.OrderID)
	assert.Equal(t, rule.ID, *retrieved.MatchedRuleID)
}

// TestRoutingDistribution_RoutingLogForDefault tests routing log for default entity
func TestRoutingDistribution_RoutingLogForDefault(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	service := rulingservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create test data
	user := CreateUniqueTestUser(t, db, "default_log_user")
	playerUser := CreateUniqueTestUser(t, db, "default_log_player")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "LOL")
	order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusPending, 10000)
	payment := CreateTestPayment(t, db, order, model.PaymentStatusPaid)

	// Create default entity
	defaultEntity := CreateTestCollectionEntity(t, db, "Default Entity")
	_ = entityRepo.SetDefault(ctx, defaultEntity.ID)

	// Create routing log for default routing
	log := &model.RoutingLog{
		PaymentID:          payment.ID,
		OrderID:            order.ID,
		MatchedRuleID:      nil, // No rule matched
		CollectionEntityID: defaultEntity.ID,
		MerchantNo:         "DEFAULT_MERCHANT",
		IsDefault:          true,
		IsFallback:         false,
		MatchDetails:       "",
	}

	err := service.CreateRoutingLog(ctx, log)
	require.NoError(t, err)

	// Verify log
	retrieved, err := service.GetRoutingLogByPayment(ctx, payment.ID)
	require.NoError(t, err)
	assert.Nil(t, retrieved.MatchedRuleID)
	assert.True(t, retrieved.IsDefault)
}

// TestRoutingDistribution_ListRoutingLogs tests listing routing logs
func TestRoutingDistribution_ListRoutingLogs(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	service := rulingservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create test data
	user := CreateUniqueTestUser(t, db, "list_logs_user")
	playerUser := CreateUniqueTestUser(t, db, "list_logs_player")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "LOG")
	entity := CreateTestCollectionEntity(t, db, "List Logs Entity")

	// Create multiple routing logs
	for i := 0; i < 3; i++ {
		order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusPending, 10000)
		payment := CreateTestPayment(t, db, order, model.PaymentStatusPaid)

		log := &model.RoutingLog{
			PaymentID:          payment.ID,
			OrderID:            order.ID,
			CollectionEntityID: entity.ID,
			MerchantNo:         fmt.Sprintf("MERCHANT_%d", i),
			IsDefault:          i%2 == 0,
		}
		require.NoError(t, service.CreateRoutingLog(ctx, log))
	}

	// List logs
	logs, total, err := service.ListRoutingLogs(ctx, routingrule.RoutingLogListOptions{
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(3))
	assert.GreaterOrEqual(t, len(logs), 3)
}

// TestRoutingDistribution_DynamicRuleUpdate tests that rule changes take effect immediately
func TestRoutingDistribution_DynamicRuleUpdate(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	service := rulingservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create entities
	firstEntity := CreateTestCollectionEntity(t, db, "First Entity")
	secondEntity := CreateTestCollectionEntity(t, db, "Second Entity")
	defaultEntity := CreateTestCollectionEntity(t, db, "Default Entity")
	_ = entityRepo.SetDefault(ctx, defaultEntity.ID)

	// Create rule for LOL
	conditions := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldGameType,
			Operator: model.ConditionOperatorEquals,
			Value:    json.RawMessage(`"王者荣耀"`),
		},
	}
	rule := CreateTestRoutingRuleWithConditions(t, db, firstEntity, 1, conditions)

	// Test initial routing
	req := &model.RoutingTestRequest{
		GameType:    "王者荣耀",
		ServiceType: "escort",
		AmountCents: 10000,
		Region:      "beijing",
	}

	result, err := service.MatchCollectionEntity(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, firstEntity.ID, result.CollectionEntityID)

	// Update rule to target different entity
	rule.TargetEntityID = secondEntity.ID
	db.Save(rule)

	// Test routing after update
	result, err = service.MatchCollectionEntity(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, secondEntity.ID, result.CollectionEntityID)
}

// TestRoutingDistribution_RuleHistoryTracksChanges tests that rule changes are tracked in history
func TestRoutingDistribution_RuleHistoryTracksChanges(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	service := rulingservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create entities
	firstEntity := CreateTestCollectionEntity(t, db, "First Entity")
	_ = CreateTestCollectionEntity(t, db, "Second Entity")
	adminUser := CreateUniqueTestUser(t, db, "admin_history")

	// Create rule
	conditions := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldGameType,
			Operator: model.ConditionOperatorEquals,
			Value:    json.RawMessage(`"王者荣耀"`),
		},
	}
	rule := CreateTestRoutingRuleWithConditions(t, db, firstEntity, 1, conditions)

	// Update rule
	newName := "Updated Rule Name"
	newPriority := 10
	_, err := service.UpdateRule(ctx, rule.ID, &model.UpdateRoutingRuleRequest{
		Name:     &newName,
		Priority: &newPriority,
	}, adminUser.ID)
	require.NoError(t, err)

	// Check history
	histories, err := service.GetRuleHistory(ctx, rule.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(histories), 2) // name and priority changes
}

// TestRoutingDistribution_CaseInsensitiveMatch tests case-insensitive matching
func TestRoutingDistribution_CaseInsensitiveMatch(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	service := rulingservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create entity
	entity := CreateTestCollectionEntity(t, db, "Case Test Entity")

	// Create rule with lowercase game type
	conditions := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldGameType,
			Operator: model.ConditionOperatorEquals,
			Value:    json.RawMessage(`"王者荣耀"`),
		},
	}
	CreateTestRoutingRuleWithConditions(t, db, entity, 1, conditions)

	// Set default entity
	defaultEntity := CreateTestCollectionEntity(t, db, "Default Entity")
	_ = entityRepo.SetDefault(ctx, defaultEntity.ID)

	// Test with different case (should still match due to case-insensitive comparison)
	req := &model.RoutingTestRequest{
		GameType:    "王者荣耀",
		ServiceType: "escort",
		AmountCents: 10000,
		Region:      "beijing",
	}

	result, err := service.MatchCollectionEntity(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, result.MatchedRuleID)
	assert.Equal(t, entity.ID, result.CollectionEntityID)
	assert.False(t, result.IsDefault)
}

// TestRoutingDistribution_ComplexRoutingScenario tests a complex real-world routing scenario
func TestRoutingDistribution_ComplexRoutingScenario(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	service := rulingservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create entities for different scenarios
	vipEntity := CreateTestCollectionEntity(t, db, "VIP High Value Entity")
	normalEntity := CreateTestCollectionEntity(t, db, "Normal Entity")
	trainingEntity := CreateTestCollectionEntity(t, db, "Training Entity")
	defaultEntity := CreateTestCollectionEntity(t, db, "Default Entity")

	// Rule 1 (Priority 1): High value orders (> 50000 cents)
	conditions1 := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldOrderAmount,
			Operator: model.ConditionOperatorGreaterThan,
			Value:    json.RawMessage(`50000`),
		},
	}
	CreateTestRoutingRuleWithConditions(t, db, vipEntity, 1, conditions1)

	// Rule 2 (Priority 2): Training service
	conditions2 := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldServiceType,
			Operator: model.ConditionOperatorEquals,
			Value:    json.RawMessage(`"training"`),
		},
	}
	CreateTestRoutingRuleWithConditions(t, db, trainingEntity, 2, conditions2)

	// Rule 3 (Priority 3): LOL games
	conditions3 := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldGameType,
			Operator: model.ConditionOperatorEquals,
			Value:    json.RawMessage(`"王者荣耀"`),
		},
	}
	CreateTestRoutingRuleWithConditions(t, db, normalEntity, 3, conditions3)

	// Set default
	_ = entityRepo.SetDefault(ctx, defaultEntity.ID)

	testCases := []struct {
		name             string
		request          *model.RoutingTestRequest
		expectedEntityID uint64
		isDefault        bool
	}{
		{
			name: "High value LOL order should use Rule 1 (high value priority)",
			request: &model.RoutingTestRequest{
				GameType:    "王者荣耀",
				ServiceType: "escort",
				AmountCents: 60000,
				Region:      "beijing",
			},
			expectedEntityID: vipEntity.ID,
			isDefault:        false,
		},
		{
			name: "Training order should use Rule 2 (training)",
			request: &model.RoutingTestRequest{
				GameType:    "王者荣耀",
				ServiceType: "training",
				AmountCents: 10000,
				Region:      "beijing",
			},
			expectedEntityID: trainingEntity.ID,
			isDefault:        false,
		},
		{
			name: "Normal LOL order should use Rule 3 (LOL)",
			request: &model.RoutingTestRequest{
				GameType:    "王者荣耀",
				ServiceType: "escort",
				AmountCents: 10000,
				Region:      "beijing",
			},
			expectedEntityID: normalEntity.ID,
			isDefault:        false,
		},
		{
			name: "Other game should use default",
			request: &model.RoutingTestRequest{
				GameType:    "绝地求生",
				ServiceType: "escort",
				AmountCents: 10000,
				Region:      "beijing",
			},
			expectedEntityID: defaultEntity.ID,
			isDefault:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.MatchCollectionEntity(ctx, tc.request)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedEntityID, result.CollectionEntityID)
			assert.Equal(t, tc.isDefault, result.IsDefault)
		})
	}
}

// TestRoutingDistribution_MatchDetailsReturned tests that match details are returned
func TestRoutingDistribution_MatchDetailsReturned(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	service := rulingservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create entity
	entity := CreateTestCollectionEntity(t, db, "Match Details Entity")

	// Create rule with multiple conditions
	conditions := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldGameType,
			Operator: model.ConditionOperatorEquals,
			Value:    json.RawMessage(`"王者荣耀"`),
		},
		{
			Field:    model.ConditionFieldServiceType,
			Operator: model.ConditionOperatorEquals,
			Value:    json.RawMessage(`"training"`),
		},
	}
	rule := CreateTestRoutingRuleWithConditions(t, db, entity, 1, conditions)

	// Test routing
	req := &model.RoutingTestRequest{
		GameType:    "王者荣耀",
		ServiceType: "training",
		AmountCents: 10000,
		Region:      "beijing",
	}

	result, err := service.MatchCollectionEntity(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, result.MatchDetails)
	assert.Len(t, result.MatchDetails, 2) // Should return matched conditions
	assert.Equal(t, rule.ID, *result.MatchedRuleID)
}

// TestRoutingDistribution_ReorderPriorities tests reordering rule priorities
func TestRoutingDistribution_ReorderPriorities(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	service := rulingservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create entities
	entity1 := CreateTestCollectionEntity(t, db, "Entity 1")
	entity2 := CreateTestCollectionEntity(t, db, "Entity 2")
	entity3 := CreateTestCollectionEntity(t, db, "Entity 3")

	// Create rules
	conditions := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldGameType,
			Operator: model.ConditionOperatorEquals,
			Value:    json.RawMessage(`"王者荣耀"`),
		},
	}
	rule1 := CreateTestRoutingRuleWithConditions(t, db, entity1, 1, conditions)
	rule2 := CreateTestRoutingRuleWithConditions(t, db, entity2, 2, conditions)
	rule3 := CreateTestRoutingRuleWithConditions(t, db, entity3, 3, conditions)

	adminUser := CreateUniqueTestUser(t, db, "admin_reorder")

	// Reorder: rule3 -> priority 1, rule1 -> priority 2, rule2 -> priority 3
	err := service.ReorderPriorities(ctx, []uint64{rule3.ID, rule1.ID, rule2.ID}, adminUser.ID)
	require.NoError(t, err)

	// Verify new order
	rules, err := service.ListActiveRulesByPriority(ctx)
	require.NoError(t, err)

	// Find our rules in the result
	var ruleIDs []uint64
	for _, r := range rules {
		if r.ID == rule1.ID || r.ID == rule2.ID || r.ID == rule3.ID {
			ruleIDs = append(ruleIDs, r.ID)
		}
	}

	assert.Equal(t, []uint64{rule3.ID, rule1.ID, rule2.ID}, ruleIDs)
}

// TestRoutingDistribution_BatchToggleStatus tests batch toggling rule status
func TestRoutingDistribution_BatchToggleStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	service := rulingservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create rules
	entity := CreateTestCollectionEntity(t, db, "Batch Entity")
	conditions := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldGameType,
			Operator: model.ConditionOperatorEquals,
			Value:    json.RawMessage(`"王者荣耀"`),
		},
	}
	rule1 := CreateTestRoutingRuleWithConditions(t, db, entity, 1, conditions)
	rule2 := CreateTestRoutingRuleWithConditions(t, db, entity, 2, conditions)
	rule3 := CreateTestRoutingRuleWithConditions(t, db, entity, 3, conditions)

	adminUser := CreateUniqueTestUser(t, db, "admin_batch")

	// Batch deactivate
	response, err := service.BatchUpdateRuleStatus(ctx, []uint64{rule1.ID, rule2.ID}, false, adminUser.ID)
	require.NoError(t, err)
	assert.Equal(t, 2, response.SuccessCount)
	assert.Equal(t, 0, response.FailedCount)

	// Verify rules are deactivated
	rules, err := service.ListActiveRulesByPriority(ctx)
	require.NoError(t, err)

	// rule3 should still be active
	var activeRuleIDs []uint64
	for _, r := range rules {
		if r.ID == rule1.ID || r.ID == rule2.ID || r.ID == rule3.ID {
			activeRuleIDs = append(activeRuleIDs, r.ID)
		}
	}

	assert.Contains(t, activeRuleIDs, rule3.ID)
	assert.NotContains(t, activeRuleIDs, rule1.ID)
	assert.NotContains(t, activeRuleIDs, rule2.ID)
}

// ============================================================================
// Helper Functions
// ============================================================================

// RoutingResult represents the result of a routing simulation
type RoutingResult struct {
	EntityID   uint64
	MerchantNo string
	LogID      uint64
}

// CreateTestRoutingRuleWithConditions creates a test routing rule with custom conditions
func CreateTestRoutingRuleWithConditions(t *testing.T, db *gorm.DB, entity *model.CollectionEntity, priority int, conditions []model.RoutingCondition) *model.RoutingRule {
	t.Helper()
	adminUser := CreateUniqueTestUser(t, db, fmt.Sprintf("admin_rule_%d", time.Now().UnixNano()))

	conditionsJSON, err := json.Marshal(conditions)
	require.NoError(t, err)

	rule := &model.RoutingRule{
		Base: model.Base{
			ExtJSON: "{}",
		},
		Name:           fmt.Sprintf("Test Rule %d", priority),
		Priority:       priority,
		Conditions:     conditionsJSON,
		TargetEntityID: entity.ID,
		Status:         model.RuleStatusActive,
		Description:    "Test routing rule",
		CreatedBy:      adminUser.ID,
	}
	if err := db.Create(rule).Error; err != nil {
		t.Fatalf("Failed to create test routing rule: %v", err)
	}
	return rule
}

// SimulateRouting simulates the routing decision process
func SimulateRouting(t *testing.T, db *gorm.DB, order *model.Order, payment *model.Payment) *RoutingResult {
	t.Helper()

	// This is a helper function that simulates what the payment service would do
	// when routing a payment to a collection entity
	ctx := context.Background()

	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	service := rulingservice.NewRoutingRuleService(ruleRepo, entityRepo)

	// Build routing request from order
	req := &model.RoutingTestRequest{
		GameType:    "王者荣耀",
		ServiceType: "escort",
		AmountCents: payment.AmountCents,
		Region:      "beijing",
	}

	result, err := service.MatchCollectionEntity(ctx, req)
	if err != nil {
		t.Fatalf("Routing simulation failed: %v", err)
	}

	// Create routing log
	log := &model.RoutingLog{
		PaymentID:          payment.ID,
		OrderID:            order.ID,
		MatchedRuleID:      result.MatchedRuleID,
		CollectionEntityID: result.CollectionEntityID,
		MerchantNo:         result.MerchantNo,
		IsDefault:          result.IsDefault,
		IsFallback:         false,
	}

	if result.MatchDetails != nil {
		detailsJSON, _ := json.Marshal(result.MatchDetails)
		log.MatchDetails = string(detailsJSON)
	}

	if err := service.CreateRoutingLog(ctx, log); err != nil {
		t.Fatalf("Failed to create routing log: %v", err)
	}

	return &RoutingResult{
		EntityID:   result.CollectionEntityID,
		MerchantNo: result.MerchantNo,
		LogID:      log.ID,
	}
}
