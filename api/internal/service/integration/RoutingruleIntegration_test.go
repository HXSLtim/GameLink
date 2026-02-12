// Package integration provides CRUD integration tests for RoutingRule module.
package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository/collectionentity"
	"gamelink/internal/repository/routingrule"
)

// ============================================================================
// RoutingRule CRUD Integration Tests
// ============================================================================

// TestRoutingRuleRepository_Create tests creating a new routing rule
func TestRoutingRuleRepository_Create(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	_ = collectionentity.NewCollectionEntityRepository(db)
	repo := routingrule.NewRoutingRuleRepository(db)
	ctx := context.Background()

	admin := CreateUniqueTestUser(t, db, "admin_rule_create")
	entity := CreateTestCollectionEntity(t, db, "Rule Entity")

	conditions := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldGameType,
			Operator: model.ConditionOperatorEquals,
			Value:    json.RawMessage(`"王者荣耀"`),
		},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rule := &model.RoutingRule{
		Name:           "Test Routing Rule",
		Priority:       1,
		Conditions:     conditionsJSON,
		TargetEntityID: entity.ID,
		Status:         model.RuleStatusActive,
		Description:    "Test rule description",
		CreatedBy:      admin.ID,
	}

	err := repo.Create(ctx, rule)
	require.NoError(t, err)
	assert.NotZero(t, rule.ID)
	assert.NotZero(t, rule.CreatedAt)
}

// TestRoutingRuleRepository_Get tests retrieving a rule by ID
func TestRoutingRuleRepository_Get(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	_ = collectionentity.NewCollectionEntityRepository(db)
	repo := routingrule.NewRoutingRuleRepository(db)
	ctx := context.Background()

	entity := CreateTestCollectionEntity(t, db, "Get Rule Entity")
	rule := CreateTestRoutingRule(t, db, entity, 1)

	got, err := repo.Get(ctx, rule.ID)
	require.NoError(t, err)
	assert.Equal(t, rule.Name, got.Name)
	assert.Equal(t, rule.Priority, got.Priority)
	assert.Equal(t, model.RuleStatusActive, got.Status)
	assert.NotNil(t, got.TargetEntity)
}

// TestRoutingRuleRepository_Get_NonExistent tests getting non-existent rule
func TestRoutingRuleRepository_Get_NonExistent(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := routingrule.NewRoutingRuleRepository(db)
	ctx := context.Background()

	_, err := repo.Get(ctx, 99999)
	assert.Error(t, err)
}

// TestRoutingRuleRepository_Update tests updating a rule
func TestRoutingRuleRepository_Update(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	_ = collectionentity.NewCollectionEntityRepository(db)
	repo := routingrule.NewRoutingRuleRepository(db)
	ctx := context.Background()

	entity := CreateTestCollectionEntity(t, db, "Update Rule Entity")
	rule := CreateTestRoutingRule(t, db, entity, 1)

	// Update rule
	rule.Name = "Updated Rule Name"
	rule.Priority = 10
	rule.Description = "Updated description"
	updater := uint64(999)
	rule.UpdatedBy = &updater

	err := repo.Update(ctx, rule)
	require.NoError(t, err)

	got, err := repo.Get(ctx, rule.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Rule Name", got.Name)
	assert.Equal(t, 10, got.Priority)
}

// TestRoutingRuleRepository_Update_Conditions tests updating rule conditions
func TestRoutingRuleRepository_Update_Conditions(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	_ = collectionentity.NewCollectionEntityRepository(db)
	repo := routingrule.NewRoutingRuleRepository(db)
	ctx := context.Background()

	entity := CreateTestCollectionEntity(t, db, "Conditions Rule Entity")
	rule := CreateTestRoutingRule(t, db, entity, 1)

	// Update with new conditions
	newConditions := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldOrderAmount,
			Operator: model.ConditionOperatorGreaterThan,
			Value:    json.RawMessage(`10000`),
		},
		{
			Field:    model.ConditionFieldServiceType,
			Operator: model.ConditionOperatorIn,
			Value:    json.RawMessage(`["escort","training"]`),
		},
	}
	err := rule.SetConditions(newConditions)
	require.NoError(t, err)

	err = repo.Update(ctx, rule)
	require.NoError(t, err)

	got, err := repo.Get(ctx, rule.ID)
	require.NoError(t, err)
	parsedConditions, err := got.GetConditions()
	require.NoError(t, err)
	assert.Len(t, parsedConditions, 2)
}

// TestRoutingRuleRepository_Delete tests deleting a rule
func TestRoutingRuleRepository_Delete(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	_ = collectionentity.NewCollectionEntityRepository(db)
	repo := routingrule.NewRoutingRuleRepository(db)
	ctx := context.Background()

	entity := CreateTestCollectionEntity(t, db, "Delete Rule Entity")
	rule := CreateTestRoutingRule(t, db, entity, 1)

	err := repo.Delete(ctx, rule.ID)
	require.NoError(t, err)

	// Verify rule is deleted
	_, err = repo.Get(ctx, rule.ID)
	assert.Error(t, err)
}

// TestRoutingRuleRepository_List tests listing rules
func TestRoutingRuleRepository_List(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	_ = collectionentity.NewCollectionEntityRepository(db)
	repo := routingrule.NewRoutingRuleRepository(db)
	ctx := context.Background()

	entity := CreateTestCollectionEntity(t, db, "List Rules Entity")

	// Create multiple rules
	for i := 0; i < 3; i++ {
		CreateTestRoutingRule(t, db, entity, i+1)
	}

	rules, total, err := repo.List(ctx, routingrule.ListOptions{
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(3))
	assert.GreaterOrEqual(t, len(rules), 3)
}

// TestRoutingRuleRepository_List_ByStatus tests filtering by status
func TestRoutingRuleRepository_List_ByStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	_ = collectionentity.NewCollectionEntityRepository(db)
	repo := routingrule.NewRoutingRuleRepository(db)
	ctx := context.Background()

	entity := CreateTestCollectionEntity(t, db, "Status Filter Entity")

	// Create active rules
	for i := 0; i < 2; i++ {
		rule := CreateTestRoutingRule(t, db, entity, i+1)
		rule.Status = model.RuleStatusActive
		db.Save(rule)
	}

	// Create inactive rule
	inactiveRule := CreateTestRoutingRule(t, db, entity, 10)
	inactiveRule.Status = model.RuleStatusInactive
	db.Save(inactiveRule)

	// List only active rules
	activeStatus := model.RuleStatusActive
	rules, total, err := repo.List(ctx, routingrule.ListOptions{
		Status:   &activeStatus,
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(2))
	for _, r := range rules {
		assert.Equal(t, model.RuleStatusActive, r.Status)
	}
}

// TestRoutingRuleRepository_List_ByTargetEntity tests filtering by target entity
func TestRoutingRuleRepository_List_ByTargetEntity(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	_ = collectionentity.NewCollectionEntityRepository(db)
	repo := routingrule.NewRoutingRuleRepository(db)
	ctx := context.Background()

	entity1 := CreateTestCollectionEntity(t, db, "Target Entity 1")
	entity2 := CreateTestCollectionEntity(t, db, "Target Entity 2")

	// Create rules for entity1
	for i := 0; i < 2; i++ {
		CreateTestRoutingRule(t, db, entity1, i+1)
	}

	// Create rule for entity2
	CreateTestRoutingRule(t, db, entity2, 10)

	// List only entity1's rules
	rules, total, err := repo.List(ctx, routingrule.ListOptions{
		TargetEntityID: &entity1.ID,
		Page:           1,
		PageSize:       10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(2))
	for _, r := range rules {
		assert.Equal(t, entity1.ID, r.TargetEntityID)
	}
}

// TestRoutingRuleRepository_List_ByKeyword tests searching by keyword
func TestRoutingRuleRepository_List_ByKeyword(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	_ = collectionentity.NewCollectionEntityRepository(db)
	repo := routingrule.NewRoutingRuleRepository(db)
	ctx := context.Background()

	entity := CreateTestCollectionEntity(t, db, "Keyword Search Entity")
	uniqueKeyword := fmt.Sprintf("UniqueKeyword_%d", time.Now().UnixNano())

	rule := CreateTestRoutingRule(t, db, entity, 1)
	rule.Name = uniqueKeyword + " Rule Name"
	rule.Description = "Description with " + uniqueKeyword
	db.Save(rule)

	rules, total, err := repo.List(ctx, routingrule.ListOptions{
		Keyword:  uniqueKeyword,
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(1))
	found := false
	for _, r := range rules {
		if r.ID == rule.ID {
			found = true
			break
		}
	}
	assert.True(t, found, "Rule should be found in search results")
}

// TestRoutingRuleRepository_ToggleStatus tests toggling rule status
func TestRoutingRuleRepository_ToggleStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	_ = collectionentity.NewCollectionEntityRepository(db)
	repo := routingrule.NewRoutingRuleRepository(db)
	ctx := context.Background()

	entity := CreateTestCollectionEntity(t, db, "Toggle Status Entity")
	rule := CreateTestRoutingRule(t, db, entity, 1)
	assert.Equal(t, model.RuleStatusActive, rule.Status)

	// Deactivate
	err := repo.ToggleStatus(ctx, rule.ID, model.RuleStatusInactive)
	require.NoError(t, err)

	got, err := repo.Get(ctx, rule.ID)
	require.NoError(t, err)
	assert.Equal(t, model.RuleStatusInactive, got.Status)

	// Reactivate
	err = repo.ToggleStatus(ctx, rule.ID, model.RuleStatusActive)
	require.NoError(t, err)

	got, err = repo.Get(ctx, rule.ID)
	require.NoError(t, err)
	assert.Equal(t, model.RuleStatusActive, got.Status)
}

// TestRoutingRuleRepository_ListActiveByPriority tests listing active rules ordered by priority
func TestRoutingRuleRepository_ListActiveByPriority(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	_ = collectionentity.NewCollectionEntityRepository(db)
	repo := routingrule.NewRoutingRuleRepository(db)
	ctx := context.Background()

	entity := CreateTestCollectionEntity(t, db, "Priority Entity")

	// Create rules with different priorities
	priorities := []int{5, 1, 3, 2, 4}
	for _, priority := range priorities {
		rule := CreateTestRoutingRule(t, db, entity, priority)
		rule.Priority = priority
		db.Save(rule)
	}

	// Get active rules by priority
	rules, err := repo.ListActiveByPriority(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(rules), 5)

	// Verify sorted by priority ascending
	for i := 1; i < len(rules); i++ {
		assert.LessOrEqual(t, rules[i-1].Priority, rules[i].Priority)
	}
}

// TestRoutingRuleRepository_ListActiveByPriority_WithInactive tests that inactive rules are excluded
func TestRoutingRuleRepository_ListActiveByPriority_WithInactive(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	_ = collectionentity.NewCollectionEntityRepository(db)
	repo := routingrule.NewRoutingRuleRepository(db)
	ctx := context.Background()

	entity := CreateTestCollectionEntity(t, db, "Mixed Status Entity")

	// Create active rules
	for i := 0; i < 2; i++ {
		rule := CreateTestRoutingRule(t, db, entity, i+1)
		rule.Status = model.RuleStatusActive
		db.Save(rule)
	}

	// Create inactive rule
	inactiveRule := CreateTestRoutingRule(t, db, entity, 10)
	inactiveRule.Status = model.RuleStatusInactive
	inactiveRule.Priority = 0 // Higher priority but inactive
	db.Save(inactiveRule)

	// Get active rules
	rules, err := repo.ListActiveByPriority(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(rules), 2)

	// Verify all are active
	for _, r := range rules {
		assert.Equal(t, model.RuleStatusActive, r.Status)
	}
}

// TestRoutingRuleRepository_CreateHistory tests creating history record
func TestRoutingRuleRepository_CreateHistory(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	_ = collectionentity.NewCollectionEntityRepository(db)
	repo := routingrule.NewRoutingRuleRepository(db)
	ctx := context.Background()

	entity := CreateTestCollectionEntity(t, db, "History Entity")
	rule := CreateTestRoutingRule(t, db, entity, 1)

	history := &model.RoutingRuleHistory{
		RoutingRuleID: rule.ID,
		FieldName:     "priority",
		OldValue:      "1",
		NewValue:      "10",
		ChangedBy:     1,
	}

	err := repo.CreateHistory(ctx, history)
	require.NoError(t, err)
	assert.NotZero(t, history.ID)

	// Verify history
	histories, err := repo.GetHistory(ctx, rule.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(histories), 1)
}

// TestRoutingRuleRepository_GetHistory tests getting rule history
func TestRoutingRuleRepository_GetHistory(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	_ = collectionentity.NewCollectionEntityRepository(db)
	repo := routingrule.NewRoutingRuleRepository(db)
	ctx := context.Background()

	entity := CreateTestCollectionEntity(t, db, "Get History Entity")
	rule := CreateTestRoutingRule(t, db, entity, 1)

	// Create multiple history records
	for i := 0; i < 3; i++ {
		history := &model.RoutingRuleHistory{
			RoutingRuleID: rule.ID,
			FieldName:     fmt.Sprintf("field_%d", i),
			OldValue:      fmt.Sprintf("old_%d", i),
			NewValue:      fmt.Sprintf("new_%d", i),
			ChangedBy:     1,
		}
		require.NoError(t, repo.CreateHistory(ctx, history))
	}

	// Get history
	histories, err := repo.GetHistory(ctx, rule.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(histories), 3)
}

// TestRoutingRuleRepository_CreateRoutingLog tests creating routing log
func TestRoutingRuleRepository_CreateRoutingLog(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	_ = collectionentity.NewCollectionEntityRepository(db)
	repo := routingrule.NewRoutingRuleRepository(db)
	ctx := context.Background()

	entity := CreateTestCollectionEntity(t, db, "Routing Log Entity")
	user := CreateUniqueTestUser(t, db, "routing_user")
	playerUser := CreateUniqueTestUser(t, db, "routing_player")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "routing_game")
	order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusPending, 10000)
	payment := CreateTestPayment(t, db, order, model.PaymentStatusPaid)

	rule := CreateTestRoutingRule(t, db, entity, 1)

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

	err := repo.CreateRoutingLog(ctx, log)
	require.NoError(t, err)
	assert.NotZero(t, log.ID)
}

// TestRoutingRuleRepository_GetRoutingLogByPayment tests getting routing log by payment
func TestRoutingRuleRepository_GetRoutingLogByPayment(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	_ = collectionentity.NewCollectionEntityRepository(db)
	repo := routingrule.NewRoutingRuleRepository(db)
	ctx := context.Background()

	entity := CreateTestCollectionEntity(t, db, "Log By Payment Entity")
	user := CreateUniqueTestUser(t, db, "log_by_payment_user")
	playerUser := CreateUniqueTestUser(t, db, "log_by_payment_player")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "log_by_payment_game")
	order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusPending, 10000)
	payment := CreateTestPayment(t, db, order, model.PaymentStatusPaid)

	log := &model.RoutingLog{
		PaymentID:          payment.ID,
		OrderID:            order.ID,
		CollectionEntityID: entity.ID,
		MerchantNo:         "MERCHANT_001",
		IsDefault:          true,
	}
	require.NoError(t, repo.CreateRoutingLog(ctx, log))

	// Get log by payment
	retrieved, err := repo.GetRoutingLogByPayment(ctx, payment.ID)
	require.NoError(t, err)
	assert.Equal(t, payment.ID, retrieved.PaymentID)
	assert.Equal(t, order.ID, retrieved.OrderID)
	assert.NotNil(t, retrieved.CollectionEntity)
}

// TestRoutingRuleRepository_ListRoutingLogs tests listing routing logs
func TestRoutingRuleRepository_ListRoutingLogs(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	_ = collectionentity.NewCollectionEntityRepository(db)
	repo := routingrule.NewRoutingRuleRepository(db)
	ctx := context.Background()

	entity := CreateTestCollectionEntity(t, db, "List Logs Entity")
	user := CreateUniqueTestUser(t, db, "list_logs_user")
	playerUser := CreateUniqueTestUser(t, db, "list_logs_player")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "list_logs_game")

	// Create multiple logs
	for i := 0; i < 3; i++ {
		order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusPending, 10000)
		payment := CreateTestPayment(t, db, order, model.PaymentStatusPaid)

		log := &model.RoutingLog{
			PaymentID:          payment.ID,
			OrderID:            order.ID,
			CollectionEntityID: entity.ID,
			MerchantNo:         fmt.Sprintf("MERCHANT_%d", i),
			IsDefault:          false,
		}
		require.NoError(t, repo.CreateRoutingLog(ctx, log))
	}

	// List logs
	logs, total, err := repo.ListRoutingLogs(ctx, routingrule.RoutingLogListOptions{
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(3))
	assert.GreaterOrEqual(t, len(logs), 3)
}

// TestRoutingRuleRepository_ListRoutingLogs_ByCollectionEntity tests filtering logs by entity
func TestRoutingRuleRepository_ListRoutingLogs_ByCollectionEntity(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	_ = collectionentity.NewCollectionEntityRepository(db)
	repo := routingrule.NewRoutingRuleRepository(db)
	ctx := context.Background()

	entity1 := CreateTestCollectionEntity(t, db, "Filter Entity 1")
	entity2 := CreateTestCollectionEntity(t, db, "Filter Entity 2")
	user := CreateUniqueTestUser(t, db, "filter_logs_user")
	playerUser := CreateUniqueTestUser(t, db, "filter_logs_player")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "filter_logs_game")

	// Create logs for entity1
	for i := 0; i < 2; i++ {
		order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusPending, 10000)
		payment := CreateTestPayment(t, db, order, model.PaymentStatusPaid)

		log := &model.RoutingLog{
			PaymentID:          payment.ID,
			OrderID:            order.ID,
			CollectionEntityID: entity1.ID,
			MerchantNo:         fmt.Sprintf("ENTITY1_%d", i),
		}
		require.NoError(t, repo.CreateRoutingLog(ctx, log))
	}

	// Create log for entity2
	order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusPending, 10000)
	payment := CreateTestPayment(t, db, order, model.PaymentStatusPaid)

	log := &model.RoutingLog{
		PaymentID:          payment.ID,
		OrderID:            order.ID,
		CollectionEntityID: entity2.ID,
		MerchantNo:         "ENTITY2_1",
	}
	require.NoError(t, repo.CreateRoutingLog(ctx, log))

	// List only entity1's logs
	logs, total, err := repo.ListRoutingLogs(ctx, routingrule.RoutingLogListOptions{
		CollectionEntityID: &entity1.ID,
		Page:               1,
		PageSize:           10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(2))
	for _, l := range logs {
		assert.Equal(t, entity1.ID, l.CollectionEntityID)
	}
}

// TestRoutingRuleRepository_ListRoutingLogs_ByIsDefault tests filtering by is_default
func TestRoutingRuleRepository_ListRoutingLogs_ByIsDefault(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	_ = collectionentity.NewCollectionEntityRepository(db)
	repo := routingrule.NewRoutingRuleRepository(db)
	ctx := context.Background()

	entity := CreateTestCollectionEntity(t, db, "Default Log Entity")
	user := CreateUniqueTestUser(t, db, "default_log_user")
	playerUser := CreateUniqueTestUser(t, db, "default_log_player")
	player := CreateTestPlayer(t, db, playerUser)
	game := CreateTestGame(t, db, "default_log_game")

	// Create default logs
	for i := 0; i < 2; i++ {
		order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusPending, 10000)
		payment := CreateTestPayment(t, db, order, model.PaymentStatusPaid)

		log := &model.RoutingLog{
			PaymentID:          payment.ID,
			OrderID:            order.ID,
			CollectionEntityID: entity.ID,
			MerchantNo:         fmt.Sprintf("DEFAULT_%d", i),
			IsDefault:          true,
		}
		require.NoError(t, repo.CreateRoutingLog(ctx, log))
	}

	// Create non-default log
	order := CreateTestOrderWithDetails(t, db, user, player, game, model.OrderStatusPending, 10000)
	payment := CreateTestPayment(t, db, order, model.PaymentStatusPaid)

	log := &model.RoutingLog{
		PaymentID:          payment.ID,
		OrderID:            order.ID,
		CollectionEntityID: entity.ID,
		MerchantNo:         "RULE_MATCHED",
		IsDefault:          false,
	}
	require.NoError(t, repo.CreateRoutingLog(ctx, log))

	// List only default logs
	isDefault := true
	logs, total, err := repo.ListRoutingLogs(ctx, routingrule.RoutingLogListOptions{
		IsDefault: &isDefault,
		Page:      1,
		PageSize:  10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(2))
	for _, l := range logs {
		assert.True(t, l.IsDefault)
	}
}

// TestRoutingRule_GetConditions tests parsing and getting conditions
func TestRoutingRule_GetConditions(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	entity := CreateTestCollectionEntity(t, db, "Conditions Parse Entity")
	rule := CreateTestRoutingRule(t, db, entity, 1)

	// Get conditions
	conditions, err := rule.GetConditions()
	require.NoError(t, err)
	assert.NotEmpty(t, conditions)
	assert.Equal(t, model.ConditionFieldGameType, conditions[0].Field)
}

// TestRoutingRule_SetConditions tests setting conditions
func TestRoutingRule_SetConditions(t *testing.T) {
	SkipIfNoTestDB(t)

	rule := &model.RoutingRule{}

	conditions := []model.RoutingCondition{
		{
			Field:    model.ConditionFieldOrderAmount,
			Operator: model.ConditionOperatorGreaterThan,
			Value:    json.RawMessage(`5000`),
		},
		{
			Field:    model.ConditionFieldRegion,
			Operator: model.ConditionOperatorEquals,
			Value:    json.RawMessage(`"beijing"`),
		},
	}

	err := rule.SetConditions(conditions)
	require.NoError(t, err)
	assert.NotNil(t, rule.Conditions)

	// Verify can parse back
	parsed, err := rule.GetConditions()
	require.NoError(t, err)
	assert.Len(t, parsed, 2)
}

// TestRoutingRule_AllConditionFields tests all condition field types
func TestRoutingRule_AllConditionFields(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	entity := CreateTestCollectionEntity(t, db, "All Fields Entity")
	admin := CreateUniqueTestUser(t, db, "admin_all_fields")

	// Test all condition fields
	fields := []model.ConditionField{
		model.ConditionFieldGameType,
		model.ConditionFieldServiceType,
		model.ConditionFieldOrderAmount,
		model.ConditionFieldRegion,
	}

	for _, field := range fields {
		var value json.RawMessage
		switch field {
		case model.ConditionFieldGameType, model.ConditionFieldServiceType, model.ConditionFieldRegion:
			value = json.RawMessage(`"test_value"`)
		case model.ConditionFieldOrderAmount:
			value = json.RawMessage(`10000`)
		}

		conditions := []model.RoutingCondition{
			{Field: field, Operator: model.ConditionOperatorEquals, Value: value},
		}
		conditionsJSON, _ := json.Marshal(conditions)

		rule := &model.RoutingRule{
			Name:           fmt.Sprintf("Test %s", field),
			Priority:       1,
			Conditions:     conditionsJSON,
			TargetEntityID: entity.ID,
			Status:         model.RuleStatusActive,
			CreatedBy:      admin.ID,
		}
		require.NoError(t, db.Create(rule).Error)
		assert.NotZero(t, rule.ID)
	}
}

// TestRoutingRule_AllConditionOperators tests all condition operators
func TestRoutingRule_AllConditionOperators(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	entity := CreateTestCollectionEntity(t, db, "All Operators Entity")
	admin := CreateUniqueTestUser(t, db, "admin_all_operators")

	// Test all operators
	operators := []model.ConditionOperator{
		model.ConditionOperatorEquals,
		model.ConditionOperatorNotEquals,
		model.ConditionOperatorIn,
		model.ConditionOperatorNotIn,
		model.ConditionOperatorGreaterThan,
		model.ConditionOperatorLessThan,
		model.ConditionOperatorBetween,
	}

	for _, operator := range operators {
		var value json.RawMessage
		switch operator {
		case model.ConditionOperatorIn, model.ConditionOperatorNotIn, model.ConditionOperatorBetween:
			value = json.RawMessage(`["value1","value2"]`)
		default:
			value = json.RawMessage(`"test_value"`)
		}

		conditions := []model.RoutingCondition{
			{Field: model.ConditionFieldGameType, Operator: operator, Value: value},
		}
		conditionsJSON, _ := json.Marshal(conditions)

		rule := &model.RoutingRule{
			Name:           fmt.Sprintf("Test %s", operator),
			Priority:       1,
			Conditions:     conditionsJSON,
			TargetEntityID: entity.ID,
			Status:         model.RuleStatusActive,
			CreatedBy:      admin.ID,
		}
		require.NoError(t, db.Create(rule).Error)
		assert.NotZero(t, rule.ID)
	}
}
