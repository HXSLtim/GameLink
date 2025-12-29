// Package integration provides batch operation integration tests for SettlementCompany and RoutingRule modules.
package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository/collectionentity"
	"gamelink/internal/repository/player"
	"gamelink/internal/repository/routingrule"
	"gamelink/internal/repository/settlementcompany"
	settlementcompanyservice "gamelink/internal/service/settlementcompany"
	routingruleservice "gamelink/internal/service/routingrule"
)

// ============================================================================
// SettlementCompany Batch Operations Tests
// ============================================================================

// TestSettlementCompanyService_BatchUpdateCompanyStatus_Success tests batch updating company status successfully
func TestSettlementCompanyService_BatchUpdateCompanyStatus_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := settlementcompany.NewSettlementCompanyRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	svc := settlementcompanyservice.NewSettlementCompanyService(repo, playerRepo)
	ctx := context.Background()

	// Create test companies with no players
	var companyIDs []uint64
	for i := 0; i < 3; i++ {
		company := CreateTestSettlementCompany(t, db, "Test Company "+string(rune('A'+i)))
		companyIDs = append(companyIDs, company.ID)
	}

	// Batch disable companies
	result, err := svc.BatchUpdateCompanyStatus(ctx, companyIDs, false, 1)
	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Len(t, result.SuccessItems, 3)
	assert.Empty(t, result.FailedItems)

	// Verify all companies are disabled
	for _, id := range companyIDs {
		company, err := repo.Get(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, model.CompanyStatusInactive, company.Status)
	}

	// Batch enable companies back
	result, err = svc.BatchUpdateCompanyStatus(ctx, companyIDs, true, 1)
	require.NoError(t, err)
	assert.Equal(t, 3, result.SuccessCount)

	// Verify all companies are enabled
	for _, id := range companyIDs {
		company, err := repo.Get(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, model.CompanyStatusActive, company.Status)
	}
}

// TestSettlementCompanyService_BatchUpdateCompanyStatus_WithActivePlayers tests disabling companies with active players fails
func TestSettlementCompanyService_BatchUpdateCompanyStatus_WithActivePlayers(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := settlementcompany.NewSettlementCompanyRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	svc := settlementcompanyservice.NewSettlementCompanyService(repo, playerRepo)
	ctx := context.Background()

	// Create companies without players
	var companyIDs []uint64
	for i := 0; i < 2; i++ {
		company := CreateTestSettlementCompany(t, db, "Empty Company "+string(rune('A'+i)))
		companyIDs = append(companyIDs, company.ID)
	}

	// Create company with active player
	busyCompany := CreateTestSettlementCompany(t, db, "Busy Company")
	user := CreateUniqueTestUser(t, db, "player_for_company")
	player := CreateTestPlayer(t, db, user)

	// Assign player to company (this increments PlayerCount)
	now := time.Now()
	assignment := &model.PlayerCompanyAssignment{
		PlayerID:            player.ID,
		SettlementCompanyID: busyCompany.ID,
		EffectiveDate:       now,
		Reason:              "Test assignment",
		AssignedBy:          1,
		IsCurrent:           true,
	}
	require.NoError(t, db.Create(assignment).Error)

	// Update company's player count
	db.Model(&model.SettlementCompany{}).Where("id = ?", busyCompany.ID).Update("player_count", 1)

	companyIDs = append(companyIDs, busyCompany.ID)

	// Batch disable companies - the busy one should fail
	result, err := svc.BatchUpdateCompanyStatus(ctx, companyIDs, false, 1)
	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 2, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Len(t, result.SuccessItems, 2)
	assert.Len(t, result.FailedItems, 1)

	// Verify failed item details
	assert.Equal(t, busyCompany.ID, result.FailedItems[0].ID)
	assert.Contains(t, result.FailedItems[0].Message, "cannot disable company with 1 active players")

	// Verify empty companies are disabled
	for _, id := range companyIDs[:2] {
		company, err := repo.Get(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, model.CompanyStatusInactive, company.Status)
	}

	// Verify busy company is still active
	busyCompanyResult, err := repo.Get(ctx, busyCompany.ID)
	require.NoError(t, err)
	assert.Equal(t, model.CompanyStatusActive, busyCompanyResult.Status)
}

// TestSettlementCompanyService_BatchUpdateCompanyStatus_PartialNotFound tests batch updating with some companies not found
func TestSettlementCompanyService_BatchUpdateCompanyStatus_PartialNotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := settlementcompany.NewSettlementCompanyRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	svc := settlementcompanyservice.NewSettlementCompanyService(repo, playerRepo)
	ctx := context.Background()

	// Create one valid company
	validCompany := CreateTestSettlementCompany(t, db, "Valid Company")
	companyIDs := []uint64{validCompany.ID, 99999, 88888}

	// Batch disable companies
	result, err := svc.BatchUpdateCompanyStatus(ctx, companyIDs, false, 1)
	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 2, result.FailedCount)
	assert.Len(t, result.SuccessItems, 1)
	assert.Len(t, result.FailedItems, 2)

	// Verify failed items are for non-existent companies
	failedIDs := make([]uint64, 0)
	for _, item := range result.FailedItems {
		failedIDs = append(failedIDs, item.ID)
		assert.Contains(t, item.Message, "not found")
	}
	assert.Contains(t, failedIDs, uint64(99999))
	assert.Contains(t, failedIDs, uint64(88888))

	// Verify valid company is disabled
	validCompanyResult, err := repo.Get(ctx, validCompany.ID)
	require.NoError(t, err)
	assert.Equal(t, model.CompanyStatusInactive, validCompanyResult.Status)
}

// TestSettlementCompanyService_BatchDeleteCompanies_Success tests batch deleting companies successfully
func TestSettlementCompanyService_BatchDeleteCompanies_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := settlementcompany.NewSettlementCompanyRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	svc := settlementcompanyservice.NewSettlementCompanyService(repo, playerRepo)
	ctx := context.Background()

	// Create test companies with no players
	var companyIDs []uint64
	for i := 0; i < 3; i++ {
		company := CreateTestSettlementCompany(t, db, "Deletable Company "+string(rune('A'+i)))
		companyIDs = append(companyIDs, company.ID)
	}

	// Batch delete companies
	result, err := svc.BatchDeleteCompanies(ctx, companyIDs)
	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Len(t, result.SuccessItems, 3)
	assert.Empty(t, result.FailedItems)

	// Verify all companies are deleted
	for _, id := range companyIDs {
		_, err := repo.Get(ctx, id)
		assert.Error(t, err)
	}
}

// TestSettlementCompanyService_BatchDeleteCompanies_WithActivePlayers tests deleting companies with active players fails
func TestSettlementCompanyService_BatchDeleteCompanies_WithActivePlayers(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := settlementcompany.NewSettlementCompanyRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	svc := settlementcompanyservice.NewSettlementCompanyService(repo, playerRepo)
	ctx := context.Background()

	// Create companies without players
	var companyIDs []uint64
	for i := 0; i < 2; i++ {
		company := CreateTestSettlementCompany(t, db, "Empty Del Company "+string(rune('A'+i)))
		companyIDs = append(companyIDs, company.ID)
	}

	// Create company with active player
	busyCompany := CreateTestSettlementCompany(t, db, "Busy Del Company")
	user := CreateUniqueTestUser(t, db, "player_for_del_company")
	player := CreateTestPlayer(t, db, user)

	// Assign player to company
	now := time.Now()
	assignment := &model.PlayerCompanyAssignment{
		PlayerID:            player.ID,
		SettlementCompanyID: busyCompany.ID,
		EffectiveDate:       now,
		Reason:              "Test assignment",
		AssignedBy:          1,
		IsCurrent:           true,
	}
	require.NoError(t, db.Create(assignment).Error)

	// Update company's player count
	db.Model(&model.SettlementCompany{}).Where("id = ?", busyCompany.ID).Update("player_count", 1)

	companyIDs = append(companyIDs, busyCompany.ID)

	// Batch delete companies - the busy one should fail
	result, err := svc.BatchDeleteCompanies(ctx, companyIDs)
	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 2, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Len(t, result.SuccessItems, 2)
	assert.Len(t, result.FailedItems, 1)

	// Verify failed item details
	assert.Equal(t, busyCompany.ID, result.FailedItems[0].ID)
	assert.Contains(t, result.FailedItems[0].Message, "cannot delete company with 1 active players")

	// Verify empty companies are deleted
	for _, id := range companyIDs[:2] {
		_, err := repo.Get(ctx, id)
		assert.Error(t, err)
	}

	// Verify busy company still exists
	busyCompanyResult, err := repo.Get(ctx, busyCompany.ID)
	require.NoError(t, err)
	assert.Equal(t, busyCompany.ID, busyCompanyResult.ID)
}

// TestSettlementCompanyService_BatchDeleteCompanies_PartialNotFound tests batch deleting with some companies not found
func TestSettlementCompanyService_BatchDeleteCompanies_PartialNotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := settlementcompany.NewSettlementCompanyRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	svc := settlementcompanyservice.NewSettlementCompanyService(repo, playerRepo)
	ctx := context.Background()

	// Create one valid company
	validCompany := CreateTestSettlementCompany(t, db, "Valid Del Company")
	companyIDs := []uint64{validCompany.ID, 99999, 88888}

	// Batch delete companies
	result, err := svc.BatchDeleteCompanies(ctx, companyIDs)
	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 2, result.FailedCount)
	assert.Len(t, result.SuccessItems, 1)
	assert.Len(t, result.FailedItems, 2)

	// Verify failed items
	failedIDs := make([]uint64, 0)
	for _, item := range result.FailedItems {
		failedIDs = append(failedIDs, item.ID)
		assert.Contains(t, item.Message, "not found")
	}
	assert.Contains(t, failedIDs, uint64(99999))
	assert.Contains(t, failedIDs, uint64(88888))

	// Verify valid company is deleted
	_, err = repo.Get(ctx, validCompany.ID)
	assert.Error(t, err)
}

// TestSettlementCompanyService_BatchUpdateCompanyStatus_EmptyIDs tests batch update status with empty IDs
func TestSettlementCompanyService_BatchUpdateCompanyStatus_EmptyIDs(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := settlementcompany.NewSettlementCompanyRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	svc := settlementcompanyservice.NewSettlementCompanyService(repo, playerRepo)
	ctx := context.Background()

	// Empty ID list should return success with zero counts
	result, err := svc.BatchUpdateCompanyStatus(ctx, []uint64{}, false, 1)
	require.NoError(t, err)
	assert.Equal(t, 0, result.TotalCount)
	assert.Equal(t, 0, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
}

// TestSettlementCompanyService_BatchDeleteCompanies_EmptyIDs tests batch delete with empty IDs
func TestSettlementCompanyService_BatchDeleteCompanies_EmptyIDs(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := settlementcompany.NewSettlementCompanyRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	svc := settlementcompanyservice.NewSettlementCompanyService(repo, playerRepo)
	ctx := context.Background()

	// Empty ID list should return success with zero counts
	result, err := svc.BatchDeleteCompanies(ctx, []uint64{})
	require.NoError(t, err)
	assert.Equal(t, 0, result.TotalCount)
	assert.Equal(t, 0, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
}

// ============================================================================
// RoutingRule Batch Operations Tests
// ============================================================================

// TestRoutingRuleService_BatchUpdateRuleStatus_Success tests batch updating rule status successfully
func TestRoutingRuleService_BatchUpdateRuleStatus_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	svc := routingruleservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create test entity and rules
	entity := CreateTestCollectionEntity(t, db, "Test Entity")
	var ruleIDs []uint64
	for i := 0; i < 3; i++ {
		rule := CreateTestRoutingRule(t, db, entity, i+1)
		ruleIDs = append(ruleIDs, rule.ID)
	}

	// Batch disable rules
	result, err := svc.BatchUpdateRuleStatus(ctx, ruleIDs, false, 1)
	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Len(t, result.SuccessItems, 3)
	assert.Empty(t, result.FailedItems)

	// Verify all rules are disabled
	for _, id := range ruleIDs {
		rule, err := ruleRepo.Get(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, model.RuleStatusInactive, rule.Status)
	}

	// Batch enable rules back
	result, err = svc.BatchUpdateRuleStatus(ctx, ruleIDs, true, 1)
	require.NoError(t, err)
	assert.Equal(t, 3, result.SuccessCount)

	// Verify all rules are enabled
	for _, id := range ruleIDs {
		rule, err := ruleRepo.Get(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, model.RuleStatusActive, rule.Status)
	}
}

// TestRoutingRuleService_BatchUpdateRuleStatus_PartialNotFound tests batch updating with some rules not found
func TestRoutingRuleService_BatchUpdateRuleStatus_PartialNotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	svc := routingruleservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create one valid rule
	entity := CreateTestCollectionEntity(t, db, "Test Entity Partial")
	validRule := CreateTestRoutingRule(t, db, entity, 1)
	ruleIDs := []uint64{validRule.ID, 99999, 88888}

	// Batch disable rules
	result, err := svc.BatchUpdateRuleStatus(ctx, ruleIDs, false, 1)
	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 2, result.FailedCount)
	assert.Len(t, result.SuccessItems, 1)
	assert.Len(t, result.FailedItems, 2)

	// Verify failed items
	failedIDs := make([]uint64, 0)
	for _, item := range result.FailedItems {
		failedIDs = append(failedIDs, item.ID)
		assert.Contains(t, item.Message, "not found")
	}
	assert.Contains(t, failedIDs, uint64(99999))
	assert.Contains(t, failedIDs, uint64(88888))

	// Verify valid rule is disabled
	validRuleResult, err := ruleRepo.Get(ctx, validRule.ID)
	require.NoError(t, err)
	assert.Equal(t, model.RuleStatusInactive, validRuleResult.Status)
}

// TestRoutingRuleService_BatchDeleteRules_Success tests batch deleting rules successfully
func TestRoutingRuleService_BatchDeleteRules_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	svc := routingruleservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create test entity and rules
	entity := CreateTestCollectionEntity(t, db, "Test Entity Del")
	var ruleIDs []uint64
	for i := 0; i < 3; i++ {
		rule := CreateTestRoutingRule(t, db, entity, i+1)
		ruleIDs = append(ruleIDs, rule.ID)
	}

	// Batch delete rules
	result, err := svc.BatchDeleteRules(ctx, ruleIDs)
	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Len(t, result.SuccessItems, 3)
	assert.Empty(t, result.FailedItems)

	// Verify all rules are deleted
	for _, id := range ruleIDs {
		_, err := ruleRepo.Get(ctx, id)
		assert.Error(t, err)
	}
}

// TestRoutingRuleService_BatchDeleteRules_PartialNotFound tests batch deleting with some rules not found
func TestRoutingRuleService_BatchDeleteRules_PartialNotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	svc := routingruleservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create one valid rule
	entity := CreateTestCollectionEntity(t, db, "Test Entity Partial Del")
	validRule := CreateTestRoutingRule(t, db, entity, 1)
	ruleIDs := []uint64{validRule.ID, 99999, 88888}

	// Batch delete rules
	result, err := svc.BatchDeleteRules(ctx, ruleIDs)
	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 2, result.FailedCount)
	assert.Len(t, result.SuccessItems, 1)
	assert.Len(t, result.FailedItems, 2)

	// Verify failed items
	failedIDs := make([]uint64, 0)
	for _, item := range result.FailedItems {
		failedIDs = append(failedIDs, item.ID)
		assert.Contains(t, item.Message, "not found")
	}
	assert.Contains(t, failedIDs, uint64(99999))
	assert.Contains(t, failedIDs, uint64(88888))

	// Verify valid rule is deleted
	_, err = ruleRepo.Get(ctx, validRule.ID)
	assert.Error(t, err)
}

// TestRoutingRuleService_BatchUpdateRuleStatus_EmptyIDs tests batch update status with empty IDs
func TestRoutingRuleService_BatchUpdateRuleStatus_EmptyIDs(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	svc := routingruleservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Empty ID list should return success with zero counts
	result, err := svc.BatchUpdateRuleStatus(ctx, []uint64{}, false, 1)
	require.NoError(t, err)
	assert.Equal(t, 0, result.TotalCount)
	assert.Equal(t, 0, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
}

// TestRoutingRuleService_BatchDeleteRules_EmptyIDs tests batch delete with empty IDs
func TestRoutingRuleService_BatchDeleteRules_EmptyIDs(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	svc := routingruleservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Empty ID list should return success with zero counts
	result, err := svc.BatchDeleteRules(ctx, []uint64{})
	require.NoError(t, err)
	assert.Equal(t, 0, result.TotalCount)
	assert.Equal(t, 0, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
}

// ============================================================================
// Combined Edge Cases Tests
// ============================================================================

// TestSettlementCompanyAndRoutingRule_BatchOperations_LargeScale tests batch operations with many items
func TestSettlementCompanyAndRoutingRule_BatchOperations_LargeScale(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Test SettlementCompany with 50 items
	settlementRepo := settlementcompany.NewSettlementCompanyRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	settlementSvc := settlementcompanyservice.NewSettlementCompanyService(settlementRepo, playerRepo)

	var companyIDs []uint64
	for i := 0; i < 50; i++ {
		company := CreateTestSettlementCompany(t, db, fmt.Sprintf("Large Scale Company %d", i))
		companyIDs = append(companyIDs, company.ID)
	}

	result, err := settlementSvc.BatchUpdateCompanyStatus(ctx, companyIDs, false, 1)
	require.NoError(t, err)
	assert.Equal(t, 50, result.TotalCount)
	assert.Equal(t, 50, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)

	// Test RoutingRule with 50 items
	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	ruleSvc := routingruleservice.NewRoutingRuleService(ruleRepo, entityRepo)

	entity := CreateTestCollectionEntity(t, db, "Large Scale Entity")
	var ruleIDs []uint64
	for i := 0; i < 50; i++ {
		rule := CreateTestRoutingRule(t, db, entity, i+1)
		ruleIDs = append(ruleIDs, rule.ID)
	}

	result2, err := ruleSvc.BatchDeleteRules(ctx, ruleIDs)
	require.NoError(t, err)
	assert.Equal(t, 50, result2.TotalCount)
	assert.Equal(t, 50, result2.SuccessCount)
	assert.Equal(t, 0, result2.FailedCount)
}

// TestSettlementCompanyService_BatchUpdateCompanyStatus_HistoryLogging tests that status changes create history records
func TestSettlementCompanyService_BatchUpdateCompanyStatus_HistoryLogging(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := settlementcompany.NewSettlementCompanyRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	svc := settlementcompanyservice.NewSettlementCompanyService(repo, playerRepo)
	ctx := context.Background()

	// Create company
	company := CreateTestSettlementCompany(t, db, "History Test Company")

	// Batch update status
	_, err := svc.BatchUpdateCompanyStatus(ctx, []uint64{company.ID}, false, 1)
	require.NoError(t, err)

	// Verify history record was created
	histories, err := repo.GetHistory(ctx, company.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(histories), 1)

	// Find the status change history
	statusChangeFound := false
	for _, history := range histories {
		if history.FieldName == "status" {
			statusChangeFound = true
			assert.Equal(t, "active", history.OldValue)
			assert.Equal(t, "inactive", history.NewValue)
			assert.Equal(t, uint64(1), history.ChangedBy)
		}
	}
	assert.True(t, statusChangeFound, "Status change history should be recorded")
}

// TestRoutingRuleService_BatchUpdateRuleStatus_HistoryLogging tests that rule status changes create history records
func TestRoutingRuleService_BatchUpdateRuleStatus_HistoryLogging(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	svc := routingruleservice.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create rule
	entity := CreateTestCollectionEntity(t, db, "History Test Entity")
	rule := CreateTestRoutingRule(t, db, entity, 1)

	// Batch update status
	_, err := svc.BatchUpdateRuleStatus(ctx, []uint64{rule.ID}, false, 1)
	require.NoError(t, err)

	// Verify history record was created
	histories, err := ruleRepo.GetHistory(ctx, rule.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(histories), 1)

	// Find the status change history
	statusChangeFound := false
	for _, history := range histories {
		if history.FieldName == "status" {
			statusChangeFound = true
			assert.Equal(t, "active", history.OldValue)
			assert.Equal(t, "inactive", history.NewValue)
			assert.Equal(t, uint64(1), history.ChangedBy)
		}
	}
	assert.True(t, statusChangeFound, "Status change history should be recorded")
}
