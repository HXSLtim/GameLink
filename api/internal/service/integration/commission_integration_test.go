// Package integration provides integration tests for the commission service.
package integration

import (
	"context"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	commissionrepo "gamelink/internal/repository/commission"
	"gamelink/internal/repository/implementations"
	playerrepo "gamelink/internal/repository/player"
	commissionservice "gamelink/internal/service/commission"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Three-Tier Commission Calculation Tests (三层抽成计算测试)
// ============================================================================

func TestCommissionService_CalculateCommission_PlayerIndividualRate(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	commissionRepo := commissionrepo.NewCommissionRepository(db)
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := playerrepo.NewPlayerRepository(db)
	svc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)

	// Create test data
	user := CreateUniqueTestUser(t, db, "test_user")
	testPlayer := CreateTestPlayer(t, db, user)
	game := CreateTestGame(t, db, "LOL")

	// Create player-specific commission rule (10% - lowest rate)
	playerRule := &model.CommissionRule{
		Name:        "Player Special Rate",
		Type:        model.CommissionRuleTypeSpecial,
		Rate:        10,
		IsActive:    true,
		PlayerID:    &testPlayer.ID,
		Description: "Special rate for this player",
	}
	err := commissionRepo.CreateRule(ctx, playerRule)
	require.NoError(t, err)

	// Create service item with 20% commission
	serviceItem := CreateTestServiceItem(t, db, game, "陪练", 10000)

	// Create completed order
	order := CreateTestOrderWithDetails(t, db, user, testPlayer, game, model.OrderStatusCompleted, 10000)
	order.ItemID = serviceItem.ID

	// Calculate commission
	calc, err := svc.CalculateOrderCommission(ctx, order)
	require.NoError(t, err)

	// Player individual rate (10%) should be used (lowest)
	assert.Equal(t, int64(10000), calc.TotalAmountCents)
	assert.Equal(t, 10, calc.CommissionRate)
	assert.Equal(t, int64(1000), calc.CommissionCents)   // 10% of 10000
	assert.Equal(t, int64(9000), calc.PlayerIncomeCents) // 10000 - 1000
	assert.Equal(t, "陪玩师专属", calc.AppliedRule)
}

func TestCommissionService_CalculateCommission_ServiceItemDefaultRate(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	commissionRepo := commissionrepo.NewCommissionRepository(db)
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := playerrepo.NewPlayerRepository(db)
	svc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)

	// Create test data
	user := CreateUniqueTestUser(t, db, "test_user")
	testPlayer := CreateTestPlayer(t, db, user)
	game := CreateTestGame(t, db, "王者荣耀")

	// Create default commission rule (25%)
	defaultRule := &model.CommissionRule{
		Name:        "Default 25%",
		Type:        model.CommissionRuleTypeDefault,
		Rate:        25,
		IsActive:    true,
		Description: "Platform default rate",
	}
	err := commissionRepo.CreateRule(ctx, defaultRule)
	require.NoError(t, err)

	// Create service item with 20% commission
	serviceItem := CreateTestServiceItem(t, db, game, "代练", 10000)

	// Create completed order
	testOrder := CreateTestOrderWithDetails(t, db, user, testPlayer, game, model.OrderStatusCompleted, 10000)
	testOrder.ItemID = serviceItem.ID

	// Calculate commission
	calc, err := svc.CalculateOrderCommission(ctx, testOrder)
	require.NoError(t, err)

	// Service item rate (20%) should be used (lower than default 25%)
	assert.Equal(t, int64(10000), calc.TotalAmountCents)
	assert.Equal(t, 20, calc.CommissionRate)
	assert.Equal(t, int64(2000), calc.CommissionCents)   // 20% of 10000
	assert.Equal(t, int64(8000), calc.PlayerIncomeCents) // 10000 - 2000
}

func TestCommissionService_CalculateCommission_DefaultRuleOnly(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	commissionRepo := commissionrepo.NewCommissionRepository(db)
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := playerrepo.NewPlayerRepository(db)
	svc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)

	// Create test data
	user := CreateUniqueTestUser(t, db, "test_user")
	testPlayer := CreateTestPlayer(t, db, user)
	_ = CreateTestGame(t, db, "吃鸡")

	// Create default commission rule (20%)
	defaultRule := &model.CommissionRule{
		Name:        "Platform Default",
		Type:        model.CommissionRuleTypeDefault,
		Rate:        20,
		IsActive:    true,
		Description: "Default platform commission",
	}
	err := commissionRepo.CreateRule(ctx, defaultRule)
	require.NoError(t, err)

	// Create order without service item
	testOrder := CreateTestOrder(t, db, user, testPlayer, model.OrderStatusCompleted)
	testOrder.TotalPriceCents = 10000
	db.Save(testOrder)

	// Calculate commission
	calc, err := svc.CalculateOrderCommission(ctx, testOrder)
	require.NoError(t, err)

	// Default rate (20%) should be used
	assert.Equal(t, int64(10000), calc.TotalAmountCents)
	assert.Equal(t, 20, calc.CommissionRate)
	assert.Equal(t, int64(2000), calc.CommissionCents)   // 20% of 10000
	assert.Equal(t, int64(8000), calc.PlayerIncomeCents) // 10000 - 2000
	assert.Equal(t, "默认规则", calc.AppliedRule)
}

func TestCommissionService_CalculateCommission_ThreeTierCalculation(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	commissionRepo := commissionrepo.NewCommissionRepository(db)
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := playerrepo.NewPlayerRepository(db)
	svc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)

	// Create test data
	user := CreateUniqueTestUser(t, db, "test_user")
	testPlayer := CreateTestPlayer(t, db, user)
	game := CreateTestGame(t, db, "LOL")

	// 1. Create player-specific rule (15% - lowest, should win)
	playerRule := &model.CommissionRule{
		Name:        "VIP Player 15%",
		Type:        model.CommissionRuleTypeSpecial,
		Rate:        15,
		IsActive:    true,
		PlayerID:    &testPlayer.ID,
		Description: "VIP player special rate",
	}
	err := commissionRepo.CreateRule(ctx, playerRule)
	require.NoError(t, err)

	// 2. Create service item with 20% commission
	serviceItem := CreateTestServiceItem(t, db, game, "陪练", 10000)

	// 3. Create default rule (25%)
	defaultRule := &model.CommissionRule{
		Name:        "Default 25%",
		Type:        model.CommissionRuleTypeDefault,
		Rate:        25,
		IsActive:    true,
		Description: "Platform default",
	}
	err = commissionRepo.CreateRule(ctx, defaultRule)
	require.NoError(t, err)

	// Create completed order with ¥100 amount
	testOrder := CreateTestOrderWithDetails(t, db, user, testPlayer, game, model.OrderStatusCompleted, 10000)
	testOrder.ItemID = serviceItem.ID

	// Calculate commission
	calc, err := svc.CalculateOrderCommission(ctx, testOrder)
	require.NoError(t, err)

	// Player rate (15%) should win - it's the lowest
	assert.Equal(t, int64(10000), calc.TotalAmountCents)
	assert.Equal(t, 15, calc.CommissionRate)
	assert.Equal(t, int64(1500), calc.CommissionCents)   // 15% of 10000
	assert.Equal(t, int64(8500), calc.PlayerIncomeCents) // 10000 - 1500
	assert.Equal(t, "陪玩师专属", calc.AppliedRule)

	// Verify all candidate rates were considered
	assert.GreaterOrEqual(t, len(calc.CandidateRates), 2)
}

func TestCommissionService_CalculateCommission_NoPlayerAssigned(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	commissionRepo := commissionrepo.NewCommissionRepository(db)
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := playerrepo.NewPlayerRepository(db)
	svc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)

	// Create test data
	user := CreateUniqueTestUser(t, db, "test_user")
	_ = CreateTestGame(t, db, "LOL")

	// Create default rule
	defaultRule := &model.CommissionRule{
		Name:     "Default 20%",
		Type:     model.CommissionRuleTypeDefault,
		Rate:     20,
		IsActive: true,
	}
	err := commissionRepo.CreateRule(ctx, defaultRule)
	require.NoError(t, err)

	// Create order without player
	testOrder := &model.Order{
		Base:            model.Base{ExtJSON: "{}"},
		OrderNo:         "TEST_NO_PLAYER",
		UserID:          user.ID,
		PlayerID:        nil, // No player assigned
		TotalPriceCents: 10000,
		Status:          model.OrderStatusCompleted,
		Currency:        model.CurrencyCNY,
		OrderConfig:     "{}",
	}
	err = db.Create(testOrder).Error
	require.NoError(t, err)

	// Calculate commission - should fail or use default only
	_, err = svc.CalculateOrderCommission(ctx, testOrder)
	// The service should handle this gracefully
	// Either return error or calculate without player-specific rate
	assert.Error(t, err)
}

// ============================================================================
// Commission Recording Tests (佣金记录测试)
// ============================================================================

func TestCommissionService_RecordCommission_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	commissionRepo := commissionrepo.NewCommissionRepository(db)
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := playerrepo.NewPlayerRepository(db)
	svc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)

	// Create test data
	user := CreateUniqueTestUser(t, db, "test_user")
	testPlayer := CreateTestPlayer(t, db, user)
	game := CreateTestGame(t, db, "LOL")

	// Create default commission rule
	defaultRule := &model.CommissionRule{
		Name:     "Default 20%",
		Type:     model.CommissionRuleTypeDefault,
		Rate:     20,
		IsActive: true,
	}
	err := commissionRepo.CreateRule(ctx, defaultRule)
	require.NoError(t, err)

	// Create completed order
	testOrder := CreateTestOrderWithDetails(t, db, user, testPlayer, game, model.OrderStatusCompleted, 10000)

	// Record commission
	err = svc.RecordCommission(ctx, testOrder.ID)
	require.NoError(t, err)

	// Verify commission record was created
	record, err := commissionRepo.GetRecordByOrderID(ctx, testOrder.ID)
	require.NoError(t, err)
	assert.Equal(t, testOrder.ID, record.OrderID)
	assert.Equal(t, testPlayer.ID, record.PlayerID)
	assert.Equal(t, int64(10000), record.TotalAmountCents)
	assert.Equal(t, 20, record.CommissionRate)
	assert.Equal(t, model.SettlementStatusPending, record.SettlementStatus)
	assert.NotEmpty(t, record.SettlementMonth)
}

func TestCommissionService_RecordCommission_PreventDuplicate(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	commissionRepo := commissionrepo.NewCommissionRepository(db)
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := playerrepo.NewPlayerRepository(db)
	svc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)

	// Create test data
	user := CreateUniqueTestUser(t, db, "test_user")
	testPlayer := CreateTestPlayer(t, db, user)
	game := CreateTestGame(t, db, "LOL")

	// Create default commission rule
	defaultRule := &model.CommissionRule{
		Name:     "Default 20%",
		Type:     model.CommissionRuleTypeDefault,
		Rate:     20,
		IsActive: true,
	}
	err := commissionRepo.CreateRule(ctx, defaultRule)
	require.NoError(t, err)

	// Create completed order
	testOrder := CreateTestOrderWithDetails(t, db, user, testPlayer, game, model.OrderStatusCompleted, 10000)

	// Record commission first time
	err = svc.RecordCommission(ctx, testOrder.ID)
	require.NoError(t, err)

	// Try to record again - should fail
	err = svc.RecordCommission(ctx, testOrder.ID)
	assert.Error(t, err)
	assert.Equal(t, commissionservice.ErrAlreadyRecorded, err)
}

func TestCommissionService_RecordCommission_OrderNotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	commissionRepo := commissionrepo.NewCommissionRepository(db)
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := playerrepo.NewPlayerRepository(db)
	svc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)

	// Try to record commission for non-existent order
	err := svc.RecordCommission(ctx, 99999)
	assert.Error(t, err)
}

// ============================================================================
// Monthly Settlement Tests (月度结算测试)
// ============================================================================

func TestCommissionService_SettleMonth_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	commissionRepo := commissionrepo.NewCommissionRepository(db)
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := playerrepo.NewPlayerRepository(db)
	svc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)

	// Create test data
	user := CreateUniqueTestUser(t, db, "test_user")
	testPlayer := CreateTestPlayer(t, db, user)
	game := CreateTestGame(t, db, "LOL")

	// Create default commission rule
	defaultRule := &model.CommissionRule{
		Name:     "Default 20%",
		Type:     model.CommissionRuleTypeDefault,
		Rate:     20,
		IsActive: true,
	}
	err := commissionRepo.CreateRule(ctx, defaultRule)
	require.NoError(t, err)

	// Create completed orders and commission records
	month := time.Now().Format("2006-01")
	for i := 0; i < 3; i++ {
		order := CreateTestOrderWithDetails(t, db, user, testPlayer, game, model.OrderStatusCompleted, 10000)

		// Create commission record
		record := &model.CommissionRecord{
			OrderID:           order.ID,
			PlayerID:          testPlayer.ID,
			TotalAmountCents:  10000,
			CommissionRate:    20,
			CommissionCents:   2000,
			PlayerIncomeCents: 8000,
			SettlementStatus:  model.SettlementStatusPending,
			SettlementMonth:   month,
		}
		err = commissionRepo.CreateRecord(ctx, record)
		require.NoError(t, err)
	}

	// Settle month
	err = svc.SettleMonth(ctx, month)
	require.NoError(t, err)

	// Verify settlement was created
	settlements, total, err := commissionRepo.ListSettlements(ctx, commissionrepo.SettlementListOptions{
		PlayerID:        &testPlayer.ID,
		SettlementMonth: &month,
		Page:            1,
		PageSize:        10,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, settlements, 1)

	settlement := settlements[0]
	assert.Equal(t, testPlayer.ID, settlement.PlayerID)
	assert.Equal(t, month, settlement.SettlementMonth)
	assert.Equal(t, int64(3), settlement.TotalOrderCount)
	assert.Equal(t, int64(30000), settlement.TotalAmountCents)
	assert.Equal(t, int64(6000), settlement.TotalCommissionCents)
	assert.Equal(t, int64(24000), settlement.TotalIncomeCents)

	// Verify commission records were updated to settled
	records, _, err := commissionRepo.ListRecords(ctx, commissionrepo.CommissionRecordListOptions{
		PlayerID:        &testPlayer.ID,
		SettlementMonth: &month,
		Page:            1,
		PageSize:        10,
	})
	require.NoError(t, err)
	for _, record := range records {
		assert.Equal(t, model.SettlementStatusSettled, record.SettlementStatus)
		assert.NotNil(t, record.SettledAt)
	}
}

func TestCommissionService_SettleMonth_AlreadySettled(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	commissionRepo := commissionrepo.NewCommissionRepository(db)
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := playerrepo.NewPlayerRepository(db)
	svc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)

	// Create test data
	user := CreateUniqueTestUser(t, db, "test_user")
	testPlayer := CreateTestPlayer(t, db, user)
	game := CreateTestGame(t, db, "LOL")

	// Create default commission rule
	defaultRule := &model.CommissionRule{
		Name:     "Default 20%",
		Type:     model.CommissionRuleTypeDefault,
		Rate:     20,
		IsActive: true,
	}
	err := commissionRepo.CreateRule(ctx, defaultRule)
	require.NoError(t, err)

	// Create order and commission record
	month := time.Now().Format("2006-01")
	testOrder := CreateTestOrderWithDetails(t, db, user, testPlayer, game, model.OrderStatusCompleted, 10000)

	record := &model.CommissionRecord{
		OrderID:           testOrder.ID,
		PlayerID:          testPlayer.ID,
		TotalAmountCents:  10000,
		CommissionRate:    20,
		CommissionCents:   2000,
		PlayerIncomeCents: 8000,
		SettlementStatus:  model.SettlementStatusPending,
		SettlementMonth:   month,
	}
	err = commissionRepo.CreateRecord(ctx, record)
	require.NoError(t, err)

	// First settlement
	err = svc.SettleMonth(ctx, month)
	require.NoError(t, err)

	// Try to settle again - should fail
	err = svc.SettleMonth(ctx, month)
	assert.Error(t, err)
	assert.Equal(t, commissionservice.ErrAlreadySettled, err)
}

func TestCommissionService_SettleMonth_NoRecords(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	commissionRepo := commissionrepo.NewCommissionRepository(db)
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := playerrepo.NewPlayerRepository(db)
	svc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)

	// Try to settle month with no records
	month := "2024-01"
	err := svc.SettleMonth(ctx, month)
	assert.Error(t, err)
}

func TestCommissionService_SettleMonth_MultiplePlayers(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	commissionRepo := commissionrepo.NewCommissionRepository(db)
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := playerrepo.NewPlayerRepository(db)
	svc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)

	// Create test data for multiple players
	user := CreateUniqueTestUser(t, db, "test_user")
	player1 := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "player1"))
	player2 := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "player2"))
	game := CreateTestGame(t, db, "LOL")

	// Create default commission rule
	defaultRule := &model.CommissionRule{
		Name:     "Default 20%",
		Type:     model.CommissionRuleTypeDefault,
		Rate:     20,
		IsActive: true,
	}
	err := commissionRepo.CreateRule(ctx, defaultRule)
	require.NoError(t, err)

	// Create orders and records for both players
	month := time.Now().Format("2006-01")

	// Player1: 2 orders
	for i := 0; i < 2; i++ {
		order := CreateTestOrderWithDetails(t, db, user, player1, game, model.OrderStatusCompleted, 10000)
		record := &model.CommissionRecord{
			OrderID:           order.ID,
			PlayerID:          player1.ID,
			TotalAmountCents:  10000,
			CommissionRate:    20,
			CommissionCents:   2000,
			PlayerIncomeCents: 8000,
			SettlementStatus:  model.SettlementStatusPending,
			SettlementMonth:   month,
		}
		err = commissionRepo.CreateRecord(ctx, record)
		require.NoError(t, err)
	}

	// Player2: 3 orders
	for i := 0; i < 3; i++ {
		order := CreateTestOrderWithDetails(t, db, user, player2, game, model.OrderStatusCompleted, 15000)
		record := &model.CommissionRecord{
			OrderID:           order.ID,
			PlayerID:          player2.ID,
			TotalAmountCents:  15000,
			CommissionRate:    20,
			CommissionCents:   3000,
			PlayerIncomeCents: 12000,
			SettlementStatus:  model.SettlementStatusPending,
			SettlementMonth:   month,
		}
		err = commissionRepo.CreateRecord(ctx, record)
		require.NoError(t, err)
	}

	// Settle month
	err = svc.SettleMonth(ctx, month)
	require.NoError(t, err)

	// Verify both players have settlements
	settlements, total, err := commissionRepo.ListSettlements(ctx, commissionrepo.SettlementListOptions{
		SettlementMonth: &month,
		Page:            1,
		PageSize:        10,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, settlements, 2)

	// Verify player1 settlement
	player1Settlement, err := commissionRepo.GetSettlementByPlayerMonth(ctx, player1.ID, month)
	require.NoError(t, err)
	assert.Equal(t, int64(2), player1Settlement.TotalOrderCount)
	assert.Equal(t, int64(20000), player1Settlement.TotalAmountCents)

	// Verify player2 settlement
	player2Settlement, err := commissionRepo.GetSettlementByPlayerMonth(ctx, player2.ID, month)
	require.NoError(t, err)
	assert.Equal(t, int64(3), player2Settlement.TotalOrderCount)
	assert.Equal(t, int64(45000), player2Settlement.TotalAmountCents)
}

// ============================================================================
// Commission Rule Management Tests (费率规则管理测试)
// ============================================================================

func TestCommissionService_CreateCommissionRule_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	commissionRepo := commissionrepo.NewCommissionRepository(db)
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := playerrepo.NewPlayerRepository(db)
	svc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)

	// Create test data
	user := CreateUniqueTestUser(t, db, "test_user")
	testPlayer := CreateTestPlayer(t, db, user)

	// Create commission rule
	req := commissionservice.CreateCommissionRuleRequest{
		Name:        "Test Rule",
		Description: "Test description",
		Type:        model.CommissionRuleTypeSpecial,
		Rate:        15,
		PlayerID:    &testPlayer.ID,
	}

	rule, err := svc.CreateCommissionRule(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, rule)
	assert.Equal(t, "Test Rule", rule.Name)
	assert.Equal(t, 15, rule.Rate)
	assert.Equal(t, model.CommissionRuleTypeSpecial, rule.Type)
	assert.True(t, rule.IsActive)
}

func TestCommissionService_CreateCommissionRule_InvalidRate(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	commissionRepo := commissionrepo.NewCommissionRepository(db)
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := playerrepo.NewPlayerRepository(db)
	svc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)

	// Try to create rule with invalid rate (> 100%)
	req := commissionservice.CreateCommissionRuleRequest{
		Name: "Invalid Rule",
		Type: model.CommissionRuleTypeDefault,
		Rate: 150, // Invalid: > 100
	}

	_, err := svc.CreateCommissionRule(ctx, req)
	assert.Error(t, err)

	// Try negative rate
	req.Rate = -10
	_, err = svc.CreateCommissionRule(ctx, req)
	assert.Error(t, err)
}

func TestCommissionService_UpdateCommissionRule_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	commissionRepo := commissionrepo.NewCommissionRepository(db)
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := playerrepo.NewPlayerRepository(db)
	svc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)

	// Create rule
	rule := CreateTestCommissionRule(t, db, model.CommissionRuleTypeDefault, 20)

	// Update rule
	newName := "Updated Rule"
	newRate := 25
	isActive := false

	req := commissionservice.UpdateCommissionRuleRequest{
		Name:     &newName,
		Rate:     &newRate,
		IsActive: &isActive,
	}

	err := svc.UpdateCommissionRule(ctx, rule.ID, req)
	require.NoError(t, err)

	// Verify update
	updatedRule, err := commissionRepo.GetRule(ctx, rule.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Rule", updatedRule.Name)
	assert.Equal(t, 25, updatedRule.Rate)
	assert.False(t, updatedRule.IsActive)
}

func TestCommissionService_DeleteCommissionRule_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	commissionRepo := commissionrepo.NewCommissionRepository(db)

	// Create rule
	rule := CreateTestCommissionRule(t, db, model.CommissionRuleTypeDefault, 20)

	// Delete rule
	err := commissionRepo.DeleteRule(ctx, rule.ID)
	require.NoError(t, err)

	// Verify deletion
	_, err = commissionRepo.GetRule(ctx, rule.ID)
	assert.Error(t, err)
	assert.Equal(t, repository.ErrNotFound, err)
}

// ============================================================================
// Statistics and Queries Tests (统计查询测试)
// ============================================================================

func TestCommissionService_GetPlayerCommissionSummary(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	commissionRepo := commissionrepo.NewCommissionRepository(db)
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := playerrepo.NewPlayerRepository(db)
	svc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)

	// Create test data
	user := CreateUniqueTestUser(t, db, "test_user")
	testPlayer := CreateTestPlayer(t, db, user)
	game := CreateTestGame(t, db, "LOL")

	// Create default commission rule
	defaultRule := &model.CommissionRule{
		Name:     "Default 20%",
		Type:     model.CommissionRuleTypeDefault,
		Rate:     20,
		IsActive: true,
	}
	err := commissionRepo.CreateRule(ctx, defaultRule)
	require.NoError(t, err)

	// Create commission records
	month := time.Now().Format("2006-01")
	for i := 0; i < 5; i++ {
		order := CreateTestOrderWithDetails(t, db, user, testPlayer, game, model.OrderStatusCompleted, 10000)
		record := &model.CommissionRecord{
			OrderID:           order.ID,
			PlayerID:          testPlayer.ID,
			TotalAmountCents:  10000,
			CommissionRate:    20,
			CommissionCents:   2000,
			PlayerIncomeCents: 8000,
			SettlementStatus:  model.SettlementStatusPending,
			SettlementMonth:   month,
		}
		err = commissionRepo.CreateRecord(ctx, record)
		require.NoError(t, err)
	}

	// Get summary
	summary, err := svc.GetPlayerCommissionSummary(ctx, testPlayer.ID, month)
	require.NoError(t, err)
	assert.Equal(t, int64(40000), summary.MonthlyIncome)   // 5 * 8000
	assert.Equal(t, int64(10000), summary.TotalCommission) // 5 * 2000
	assert.Equal(t, int64(40000), summary.TotalIncome)     // 5 * 8000
	assert.Equal(t, int64(5), summary.TotalOrders)
}

func TestCommissionService_GetCommissionRecords_Pagination(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	commissionRepo := commissionrepo.NewCommissionRepository(db)
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := playerrepo.NewPlayerRepository(db)
	svc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)

	// Create test data
	user := CreateUniqueTestUser(t, db, "test_user")
	testPlayer := CreateTestPlayer(t, db, user)
	game := CreateTestGame(t, db, "LOL")

	// Create commission records
	month := time.Now().Format("2006-01")
	for i := 0; i < 15; i++ {
		order := CreateTestOrderWithDetails(t, db, user, testPlayer, game, model.OrderStatusCompleted, 10000)
		record := &model.CommissionRecord{
			OrderID:           order.ID,
			PlayerID:          testPlayer.ID,
			TotalAmountCents:  10000,
			CommissionRate:    20,
			CommissionCents:   2000,
			PlayerIncomeCents: 8000,
			SettlementStatus:  model.SettlementStatusPending,
			SettlementMonth:   month,
		}
		err := commissionRepo.CreateRecord(ctx, record)
		require.NoError(t, err)
	}

	// Get first page
	resp, err := svc.GetCommissionRecords(ctx, testPlayer.ID, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(15), resp.Total)
	assert.Len(t, resp.Records, 10)

	// Get second page
	resp, err = svc.GetCommissionRecords(ctx, testPlayer.ID, 2, 10)
	require.NoError(t, err)
	assert.Len(t, resp.Records, 5)
}

func TestCommissionService_GetMonthlySettlements_PlayerFilter(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	commissionRepo := commissionrepo.NewCommissionRepository(db)
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := playerrepo.NewPlayerRepository(db)
	svc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)

	// Create test data
	user := CreateUniqueTestUser(t, db, "test_user")
	testPlayer := CreateTestPlayer(t, db, user)
	game := CreateTestGame(t, db, "LOL")

	// Create default commission rule
	defaultRule := &model.CommissionRule{
		Name:     "Default 20%",
		Type:     model.CommissionRuleTypeDefault,
		Rate:     20,
		IsActive: true,
	}
	err := commissionRepo.CreateRule(ctx, defaultRule)
	require.NoError(t, err)

	// Create and settle month
	month := time.Now().Format("2006-01")
	order := CreateTestOrderWithDetails(t, db, user, testPlayer, game, model.OrderStatusCompleted, 10000)
	record := &model.CommissionRecord{
		OrderID:           order.ID,
		PlayerID:          testPlayer.ID,
		TotalAmountCents:  10000,
		CommissionRate:    20,
		CommissionCents:   2000,
		PlayerIncomeCents: 8000,
		SettlementStatus:  model.SettlementStatusPending,
		SettlementMonth:   month,
	}
	err = commissionRepo.CreateRecord(ctx, record)
	require.NoError(t, err)

	err = svc.SettleMonth(ctx, month)
	require.NoError(t, err)

	// Get settlements for player
	resp, err := svc.GetMonthlySettlements(ctx, testPlayer.ID, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), resp.Total)
	assert.Len(t, resp.Settlements, 1)
	assert.Equal(t, month, resp.Settlements[0].SettlementMonth)
}

func TestCommissionService_GetPlatformStats(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	commissionRepo := commissionrepo.NewCommissionRepository(db)
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := playerrepo.NewPlayerRepository(db)
	svc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)

	// Create test data
	user := CreateUniqueTestUser(t, db, "test_user")
	testPlayer := CreateTestPlayer(t, db, user)
	game := CreateTestGame(t, db, "LOL")

	// Create default commission rule
	defaultRule := &model.CommissionRule{
		Name:     "Default 20%",
		Type:     model.CommissionRuleTypeDefault,
		Rate:     20,
		IsActive: true,
	}
	err := commissionRepo.CreateRule(ctx, defaultRule)
	require.NoError(t, err)

	// Create settled records
	month := time.Now().Format("2006-01")
	for i := 0; i < 3; i++ {
		order := CreateTestOrderWithDetails(t, db, user, testPlayer, game, model.OrderStatusCompleted, 10000)
		record := &model.CommissionRecord{
			OrderID:           order.ID,
			PlayerID:          testPlayer.ID,
			TotalAmountCents:  10000,
			CommissionRate:    20,
			CommissionCents:   2000,
			PlayerIncomeCents: 8000,
			SettlementStatus:  model.SettlementStatusSettled,
			SettlementMonth:   month,
		}
		err = commissionRepo.CreateRecord(ctx, record)
		require.NoError(t, err)
	}

	// Get platform stats
	stats, err := svc.GetPlatformStats(ctx, month)
	require.NoError(t, err)
	assert.Equal(t, month, stats.Month)
	assert.Equal(t, int64(3), stats.TotalOrders)
	assert.Equal(t, int64(30000), stats.TotalIncome)
	assert.Equal(t, int64(6000), stats.TotalCommission)
	assert.Equal(t, int64(24000), stats.TotalPlayerIncome)
}

// ============================================================================
// Batch Operations Tests (批量操作测试)
// ============================================================================

func TestCommissionService_BatchDeleteCommissionRules_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	commissionRepo := commissionrepo.NewCommissionRepository(db)
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := playerrepo.NewPlayerRepository(db)
	svc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)

	// Create rules
	var ruleIDs []uint64
	for i := 0; i < 3; i++ {
		rule := CreateTestCommissionRule(t, db, model.CommissionRuleTypeDefault, 20+i*5)
		ruleIDs = append(ruleIDs, rule.ID)
	}

	// Batch delete
	result, err := svc.BatchDeleteCommissionRules(ctx, ruleIDs)
	require.NoError(t, err)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Empty(t, result.FailedIDs)

	// Verify deletion
	for _, id := range ruleIDs {
		_, err := commissionRepo.GetRule(ctx, id)
		assert.Error(t, err)
	}
}

func TestCommissionService_BatchDeleteCommissionRules_MixedSuccess(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	commissionRepo := commissionrepo.NewCommissionRepository(db)
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := playerrepo.NewPlayerRepository(db)
	svc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)

	// Create some rules
	rule1 := CreateTestCommissionRule(t, db, model.CommissionRuleTypeDefault, 20)
	rule2 := CreateTestCommissionRule(t, db, model.CommissionRuleTypeSpecial, 15)

	// Include non-existent IDs
	ruleIDs := []uint64{rule1.ID, rule2.ID, 99999, 88888}

	// Batch delete
	result, err := svc.BatchDeleteCommissionRules(ctx, ruleIDs)
	require.NoError(t, err)
	assert.Equal(t, 2, result.SuccessCount)
	assert.Equal(t, 2, result.FailedCount)
	assert.Len(t, result.FailedIDs, 2)
	assert.Contains(t, result.FailedIDs, uint64(99999))
	assert.Contains(t, result.FailedIDs, uint64(88888))
}

func TestCommissionService_BatchDeleteCommissionRules_EmptyList(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	commissionRepo := commissionrepo.NewCommissionRepository(db)
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := playerrepo.NewPlayerRepository(db)
	svc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)

	// Try to delete empty list
	_, err := svc.BatchDeleteCommissionRules(ctx, []uint64{})
	assert.Error(t, err)
}

func TestCommissionService_BatchUpdateCommissionRuleStatus_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	commissionRepo := commissionrepo.NewCommissionRepository(db)
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := playerrepo.NewPlayerRepository(db)
	svc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)

	// Create active rules
	var ruleIDs []uint64
	for i := 0; i < 3; i++ {
		rule := CreateTestCommissionRule(t, db, model.CommissionRuleTypeDefault, 20)
		ruleIDs = append(ruleIDs, rule.ID)
	}

	// Batch deactivate
	result, err := svc.BatchUpdateCommissionRuleStatus(ctx, ruleIDs, false)
	require.NoError(t, err)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)

	// Verify all are deactivated
	for _, id := range ruleIDs {
		rule, err := commissionRepo.GetRule(ctx, id)
		require.NoError(t, err)
		assert.False(t, rule.IsActive)
	}
}

func TestCommissionService_BatchUpdateCommissionRuleStatus_Activate(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	commissionRepo := commissionrepo.NewCommissionRepository(db)
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := playerrepo.NewPlayerRepository(db)
	svc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)

	// Create inactive rules
	var ruleIDs []uint64
	for i := 0; i < 2; i++ {
		rule := CreateTestCommissionRule(t, db, model.CommissionRuleTypeDefault, 20)
		rule.IsActive = false
		db.Save(rule)
		ruleIDs = append(ruleIDs, rule.ID)
	}

	// Batch activate
	result, err := svc.BatchUpdateCommissionRuleStatus(ctx, ruleIDs, true)
	require.NoError(t, err)
	assert.Equal(t, 2, result.SuccessCount)

	// Verify all are activated
	for _, id := range ruleIDs {
		rule, err := commissionRepo.GetRule(ctx, id)
		require.NoError(t, err)
		assert.True(t, rule.IsActive)
	}
}

// ============================================================================
// Ranking Commission Helper Functions Tests (排名抽成辅助函数测试)
// ============================================================================

func TestCommissionService_ParseRankingCommissionRules_Valid(t *testing.T) {
	// Valid JSON rules
	rulesJSON := `[
		{"rankStart": 1, "rankEnd": 3, "commissionRate": 10},
		{"rankStart": 4, "rankEnd": 10, "commissionRate": 12},
		{"rankStart": 11, "rankEnd": 50, "commissionRate": 15}
	]`

	rules, err := commissionservice.ParseRankingCommissionRules(rulesJSON)
	require.NoError(t, err)
	assert.Len(t, rules, 3)
	assert.Equal(t, 1, rules[0].RankStart)
	assert.Equal(t, 3, rules[0].RankEnd)
	assert.Equal(t, 10, rules[0].CommissionRate)
}

func TestCommissionService_ParseRankingCommissionRules_Invalid(t *testing.T) {
	// Invalid JSON
	rulesJSON := `{invalid json}`

	_, err := commissionservice.ParseRankingCommissionRules(rulesJSON)
	assert.Error(t, err)
}

func TestCommissionService_FindCommissionRateForRank_Match(t *testing.T) {
	rules := []model.RankingCommissionRule{
		{RankStart: 1, RankEnd: 3, CommissionRate: 10},
		{RankStart: 4, RankEnd: 10, CommissionRate: 12},
		{RankStart: 11, RankEnd: 50, CommissionRate: 15},
	}

	// Test rank 5 (should match 4-10 bracket)
	rate := commissionservice.FindCommissionRateForRank(rules, 5)
	assert.Equal(t, 12, rate)

	// Test rank 1 (should match 1-3 bracket)
	rate = commissionservice.FindCommissionRateForRank(rules, 1)
	assert.Equal(t, 10, rate)

	// Test rank 20 (should match 11-50 bracket)
	rate = commissionservice.FindCommissionRateForRank(rules, 20)
	assert.Equal(t, 15, rate)
}

func TestCommissionService_FindCommissionRateForRank_NoMatch(t *testing.T) {
	rules := []model.RankingCommissionRule{
		{RankStart: 1, RankEnd: 10, CommissionRate: 15},
	}

	// Test rank 50 (no match)
	rate := commissionservice.FindCommissionRateForRank(rules, 50)
	assert.Equal(t, 0, rate)
}

func TestCommissionService_ValidateRankingRules_Valid(t *testing.T) {
	rules := []model.RankingCommissionRule{
		{RankStart: 1, RankEnd: 10, CommissionRate: 15},
		{RankStart: 11, RankEnd: 50, CommissionRate: 20},
	}

	err := commissionservice.ValidateRankingRules(rules)
	assert.NoError(t, err)
}

func TestCommissionService_ValidateRankingRules_InvalidRange(t *testing.T) {
	// RankStart < 1
	rules := []model.RankingCommissionRule{
		{RankStart: 0, RankEnd: 10, CommissionRate: 15},
	}

	err := commissionservice.ValidateRankingRules(rules)
	assert.Error(t, err)
	assert.Equal(t, commissionservice.ErrValidation, err)
}

func TestCommissionService_ValidateRankingRules_InvalidRate(t *testing.T) {
	// Rate > 100
	rules := []model.RankingCommissionRule{
		{RankStart: 1, RankEnd: 10, CommissionRate: 150},
	}

	err := commissionservice.ValidateRankingRules(rules)
	assert.Error(t, err)
}

func TestCommissionService_ValidateRankingRules_Overlapping(t *testing.T) {
	// Overlapping ranges
	rules := []model.RankingCommissionRule{
		{RankStart: 1, RankEnd: 10, CommissionRate: 15},
		{RankStart: 5, RankEnd: 20, CommissionRate: 20}, // Overlaps with first
	}

	err := commissionservice.ValidateRankingRules(rules)
	assert.Error(t, err)
}

func TestCommissionService_ValidateRankingRules_InvertedRange(t *testing.T) {
	// RankEnd < RankStart
	rules := []model.RankingCommissionRule{
		{RankStart: 10, RankEnd: 5, CommissionRate: 15},
	}

	err := commissionservice.ValidateRankingRules(rules)
	assert.Error(t, err)
}

// ============================================================================
// Edge Cases and Boundary Tests (边界情况测试)
// ============================================================================

func TestCommissionService_CalculateCommission_ZeroRate(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	commissionRepo := commissionrepo.NewCommissionRepository(db)
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := playerrepo.NewPlayerRepository(db)
	svc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)

	// Create test data
	user := CreateUniqueTestUser(t, db, "test_user")
	testPlayer := CreateTestPlayer(t, db, user)
	game := CreateTestGame(t, db, "LOL")

	// Create player-specific rule with 0% commission
	playerRule := &model.CommissionRule{
		Name:        "Zero Commission",
		Type:        model.CommissionRuleTypeSpecial,
		Rate:        0,
		IsActive:    true,
		PlayerID:    &testPlayer.ID,
		Description: "Special zero commission",
	}
	err := commissionRepo.CreateRule(ctx, playerRule)
	require.NoError(t, err)

	// Create order
	testOrder := CreateTestOrderWithDetails(t, db, user, testPlayer, game, model.OrderStatusCompleted, 10000)

	// Calculate commission
	calc, err := svc.CalculateOrderCommission(ctx, testOrder)
	require.NoError(t, err)

	// Zero rate should be used (player gets 100%)
	assert.Equal(t, 0, calc.CommissionRate)
	assert.Equal(t, int64(0), calc.CommissionCents)
	assert.Equal(t, int64(10000), calc.PlayerIncomeCents)
}

func TestCommissionService_CalculateCommission_MaxRate(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	commissionRepo := commissionrepo.NewCommissionRepository(db)
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := playerrepo.NewPlayerRepository(db)
	svc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)

	// Create test data
	user := CreateUniqueTestUser(t, db, "test_user")
	testPlayer := CreateTestPlayer(t, db, user)
	game := CreateTestGame(t, db, "LOL")

	// Create player-specific rule with 100% commission
	playerRule := &model.CommissionRule{
		Name:        "Max Commission",
		Type:        model.CommissionRuleTypeSpecial,
		Rate:        100,
		IsActive:    true,
		PlayerID:    &testPlayer.ID,
		Description: "Platform takes all",
	}
	err := commissionRepo.CreateRule(ctx, playerRule)
	require.NoError(t, err)

	// Create order
	testOrder := CreateTestOrderWithDetails(t, db, user, testPlayer, game, model.OrderStatusCompleted, 10000)

	// Calculate commission
	calc, err := svc.CalculateOrderCommission(ctx, testOrder)
	require.NoError(t, err)

	// 100% rate - platform takes all
	assert.Equal(t, 100, calc.CommissionRate)
	assert.Equal(t, int64(10000), calc.CommissionCents)
	assert.Equal(t, int64(0), calc.PlayerIncomeCents)
}

func TestCommissionService_CalculateCommission_VerySmallAmount(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	commissionRepo := commissionrepo.NewCommissionRepository(db)
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := playerrepo.NewPlayerRepository(db)
	svc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)

	// Create test data
	user := CreateUniqueTestUser(t, db, "test_user")
	testPlayer := CreateTestPlayer(t, db, user)
	game := CreateTestGame(t, db, "LOL")

	// Create default rule
	defaultRule := &model.CommissionRule{
		Name:     "Default 20%",
		Type:     model.CommissionRuleTypeDefault,
		Rate:     20,
		IsActive: true,
	}
	err := commissionRepo.CreateRule(ctx, defaultRule)
	require.NoError(t, err)

	// Create order with very small amount (1 cent)
	testOrder := CreateTestOrderWithDetails(t, db, user, testPlayer, game, model.OrderStatusCompleted, 1)

	// Calculate commission
	calc, err := svc.CalculateOrderCommission(ctx, testOrder)
	require.NoError(t, err)

	assert.Equal(t, int64(1), calc.TotalAmountCents)
	// 20% of 1 cent = 0.2 cents, integer division should give 0
	assert.Equal(t, int64(0), calc.CommissionCents)
	assert.Equal(t, int64(1), calc.PlayerIncomeCents)
}

func TestCommissionService_CalculateCommission_LargeAmount(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	commissionRepo := commissionrepo.NewCommissionRepository(db)
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := playerrepo.NewPlayerRepository(db)
	svc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)

	// Create test data
	user := CreateUniqueTestUser(t, db, "test_user")
	testPlayer := CreateTestPlayer(t, db, user)
	game := CreateTestGame(t, db, "LOL")

	// Create default rule
	defaultRule := &model.CommissionRule{
		Name:     "Default 20%",
		Type:     model.CommissionRuleTypeDefault,
		Rate:     20,
		IsActive: true,
	}
	err := commissionRepo.CreateRule(ctx, defaultRule)
	require.NoError(t, err)

	// Create order with large amount (1 million cents = ¥10,000)
	testOrder := CreateTestOrderWithDetails(t, db, user, testPlayer, game, model.OrderStatusCompleted, 1000000)

	// Calculate commission
	calc, err := svc.CalculateOrderCommission(ctx, testOrder)
	require.NoError(t, err)

	assert.Equal(t, int64(1000000), calc.TotalAmountCents)
	assert.Equal(t, int64(200000), calc.CommissionCents)
	assert.Equal(t, int64(800000), calc.PlayerIncomeCents)
}

func TestCommissionService_InactiveRuleNotApplied(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup service
	commissionRepo := commissionrepo.NewCommissionRepository(db)
	orderRepo := implementations.NewOrderRepository(db)
	playerRepo := playerrepo.NewPlayerRepository(db)
	svc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)

	// Create test data
	user := CreateUniqueTestUser(t, db, "test_user")
	testPlayer := CreateTestPlayer(t, db, user)
	game := CreateTestGame(t, db, "LOL")

	// Create inactive player-specific rule
	playerRule := &model.CommissionRule{
		Name:        "Inactive Special Rate",
		Type:        model.CommissionRuleTypeSpecial,
		Rate:        10,
		IsActive:    false, // Inactive
		PlayerID:    &testPlayer.ID,
		Description: "This should not be applied",
	}
	err := commissionRepo.CreateRule(ctx, playerRule)
	require.NoError(t, err)

	// Create active default rule
	defaultRule := &model.CommissionRule{
		Name:     "Default 20%",
		Type:     model.CommissionRuleTypeDefault,
		Rate:     20,
		IsActive: true,
	}
	err = commissionRepo.CreateRule(ctx, defaultRule)
	require.NoError(t, err)

	// Create order
	testOrder := CreateTestOrderWithDetails(t, db, user, testPlayer, game, model.OrderStatusCompleted, 10000)

	// Calculate commission
	calc, err := svc.CalculateOrderCommission(ctx, testOrder)
	require.NoError(t, err)

	// Default rate (20%) should be used, not inactive player rule (10%)
	assert.Equal(t, 20, calc.CommissionRate)
	assert.Equal(t, "默认规则", calc.AppliedRule)
}
