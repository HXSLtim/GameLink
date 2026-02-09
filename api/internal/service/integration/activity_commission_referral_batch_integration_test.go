// Package integration provides integration tests for batch operations
// across Activity, Commission, and Referral services.
package integration

import (
	"context"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository/activity"
	"gamelink/internal/repository/commission"
	"gamelink/internal/repository/referral"
	activityservice "gamelink/internal/service/activity"
	commissionservice "gamelink/internal/service/commission"
	referralservice "gamelink/internal/service/referral"

	"gorm.io/gorm"
)

// ============================================================================
// Activity Batch Operations Tests
// ============================================================================

func TestActivityService_BatchDeleteActivities_PartialFailure(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create test activities
	activity1 := createTestActivity(t, db, "Activity 1", model.ActivityStatusDraft)
	activity2 := createTestActivity(t, db, "Activity 2", model.ActivityStatusDraft)
	activity3 := createTestActivity(t, db, "Activity 3", model.ActivityStatusEnded)

	// Create repository and service
	activityRepo := activity.NewActivityRepository(db)
	service := activityservice.NewActivityService(activityRepo, nil)

	// Test: Delete with mix of valid and invalid IDs
	nonExistentID := uint64(99999)
	ids := []uint64{activity1.ID, activity2.ID, nonExistentID}

	result, err := service.BatchDeleteActivities(ctx, ids)
	if err != nil {
		t.Fatalf("BatchDeleteActivities failed: %v", err)
	}

	// Assertions
	if result.SuccessCount != 2 {
		t.Errorf("Expected 2 successful deletions, got %d", result.SuccessCount)
	}
	if result.FailedCount != 1 {
		t.Errorf("Expected 1 failed deletion, got %d", result.FailedCount)
	}
	if len(result.FailedIDs) != 1 {
		t.Errorf("Expected 1 failed ID, got %d", len(result.FailedIDs))
	}
	if len(result.Errors) != 1 {
		t.Errorf("Expected 1 error message, got %d", len(result.Errors))
	}

	// Verify database state
	var count int64
	db.Model(&model.Activity{}).Where("id IN ?", []uint64{activity1.ID, activity2.ID}).Count(&count)
	if count != 0 {
		t.Errorf("Expected activities to be deleted, but found %d remaining", count)
	}

	// Verify activity3 still exists (ended status, should not be deleted in our test)
	var remaining model.Activity
	if err := db.First(&remaining, activity3.ID).Error; err != nil {
		t.Errorf("Activity 3 should still exist: %v", err)
	}
}

func TestActivityService_BatchDeleteActivities_ActiveStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create test activities - one active (cannot delete)
	activity1 := createTestActivity(t, db, "Active Activity", model.ActivityStatusActive)
	activity2 := createTestActivity(t, db, "Draft Activity", model.ActivityStatusDraft)

	activityRepo := activity.NewActivityRepository(db)
	service := activityservice.NewActivityService(activityRepo, nil)

	ids := []uint64{activity1.ID, activity2.ID}

	result, err := service.BatchDeleteActivities(ctx, ids)
	if err != nil {
		t.Fatalf("BatchDeleteActivities failed: %v", err)
	}

	// Assertions - active activity should fail deletion
	if result.SuccessCount != 1 {
		t.Errorf("Expected 1 successful deletion, got %d", result.SuccessCount)
	}
	if result.FailedCount != 1 {
		t.Errorf("Expected 1 failed deletion, got %d", result.FailedCount)
	}

	// Verify active activity still exists
	var remaining model.Activity
	if err := db.First(&remaining, activity1.ID).Error; err != nil {
		t.Errorf("Active activity should still exist: %v", err)
	}
}

func TestActivityService_BatchUpdateActivityStatus_AllSuccess(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create test activities
	activity1 := createTestActivity(t, db, "Activity 1", model.ActivityStatusDraft)
	activity2 := createTestActivity(t, db, "Activity 2", model.ActivityStatusDraft)
	activity3 := createTestActivity(t, db, "Activity 3", model.ActivityStatusDraft)

	activityRepo := activity.NewActivityRepository(db)
	service := activityservice.NewActivityService(activityRepo, nil)

	// Test: Publish activities (draft -> preheat)
	ids := []uint64{activity1.ID, activity2.ID, activity3.ID}

	result, err := service.BatchUpdateActivityStatus(ctx, ids, model.ActivityStatusPreheat)
	if err != nil {
		t.Fatalf("BatchUpdateActivityStatus failed: %v", err)
	}

	// Assertions
	if result.SuccessCount != 3 {
		t.Errorf("Expected 3 successful updates, got %d", result.SuccessCount)
	}
	if result.FailedCount != 0 {
		t.Errorf("Expected 0 failed updates, got %d", result.FailedCount)
	}

	// Verify database state
	var activities []model.Activity
	db.Find(&activities, ids)
	for _, a := range activities {
		if a.Status != model.ActivityStatusPreheat {
			t.Errorf("Activity %d status: expected %s, got %s", a.ID, model.ActivityStatusPreheat, a.Status)
		}
	}
}

func TestActivityService_BatchUpdateActivityStatus_InvalidTransition(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create test activities with different statuses
	activity1 := createTestActivity(t, db, "Draft Activity", model.ActivityStatusDraft)
	activity2 := createTestActivity(t, db, "Ended Activity", model.ActivityStatusEnded)

	activityRepo := activity.NewActivityRepository(db)
	service := activityservice.NewActivityService(activityRepo, nil)

	// Test: Try to transition ended activity (should fail - no valid transitions from ended)
	ids := []uint64{activity1.ID, activity2.ID}

	result, err := service.BatchUpdateActivityStatus(ctx, ids, model.ActivityStatusActive)
	if err != nil {
		t.Fatalf("BatchUpdateActivityStatus failed: %v", err)
	}

	// Assertions - ended activity cannot transition
	if result.SuccessCount != 1 {
		t.Errorf("Expected 1 successful update, got %d", result.SuccessCount)
	}
	if result.FailedCount != 1 {
		t.Errorf("Expected 1 failed update, got %d", result.FailedCount)
	}

	// Verify ended activity still has ended status
	var endedActivity model.Activity
	db.First(&endedActivity, activity2.ID)
	if endedActivity.Status != model.ActivityStatusEnded {
		t.Errorf("Ended activity status should not change: expected %s, got %s",
			model.ActivityStatusEnded, endedActivity.Status)
	}
}

func TestActivityService_BatchUpdateActivityStatus_PartialNonExistent(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create test activities
	activity1 := createTestActivity(t, db, "Activity 1", model.ActivityStatusDraft)

	activityRepo := activity.NewActivityRepository(db)
	service := activityservice.NewActivityService(activityRepo, nil)

	// Test: Include non-existent ID
	nonExistentID := uint64(99999)
	ids := []uint64{activity1.ID, nonExistentID}

	result, err := service.BatchUpdateActivityStatus(ctx, ids, model.ActivityStatusActive)
	if err != nil {
		t.Fatalf("BatchUpdateActivityStatus failed: %v", err)
	}

	// Assertions
	if result.SuccessCount != 1 {
		t.Errorf("Expected 1 successful update, got %d", result.SuccessCount)
	}
	if result.FailedCount != 1 {
		t.Errorf("Expected 1 failed update, got %d", result.FailedCount)
	}
	if len(result.FailedIDs) != 1 || result.FailedIDs[0] != nonExistentID {
		t.Errorf("FailedIDs should contain non-existent ID")
	}
}

// ============================================================================
// Commission Batch Operations Tests
// ============================================================================

func TestCommissionService_BatchDeleteCommissionRules_PartialFailure(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create test commission rules
	rule1 := CreateTestCommissionRule(t, db, model.CommissionRuleTypeDefault, 20)
	rule2 := CreateTestCommissionRule(t, db, model.CommissionRuleTypeSpecial, 15)
	rule3 := CreateTestCommissionRule(t, db, model.CommissionRuleTypeGift, 25)

	// Create repository and service
	commissionRepo := commission.NewCommissionRepository(db)
	service := commissionservice.NewCommissionService(commissionRepo, &mockOrderReader{}, &mockPlayerRepo{})

	// Test: Delete with mix of valid and invalid IDs
	nonExistentID := uint64(99999)
	ids := []uint64{rule1.ID, rule2.ID, nonExistentID}

	result, err := service.BatchDeleteCommissionRules(ctx, ids)
	if err != nil {
		t.Fatalf("BatchDeleteCommissionRules failed: %v", err)
	}

	// Assertions
	if result.SuccessCount != 2 {
		t.Errorf("Expected 2 successful deletions, got %d", result.SuccessCount)
	}
	if result.FailedCount != 1 {
		t.Errorf("Expected 1 failed deletion, got %d", result.FailedCount)
	}
	if len(result.FailedIDs) != 1 {
		t.Errorf("Expected 1 failed ID, got %d", len(result.FailedIDs))
	}
	if len(result.Errors) != 1 {
		t.Errorf("Expected 1 error message, got %d", len(result.Errors))
	}

	// Verify database state
	var count int64
	db.Model(&model.CommissionRule{}).Where("id IN ?", []uint64{rule1.ID, rule2.ID}).Count(&count)
	if count != 0 {
		t.Errorf("Expected rules to be deleted, but found %d remaining", count)
	}

	// Verify rule3 still exists
	var remaining model.CommissionRule
	if err := db.First(&remaining, rule3.ID).Error; err != nil {
		t.Errorf("Rule 3 should still exist: %v", err)
	}
}

func TestCommissionService_BatchUpdateCommissionRuleStatus_AllSuccess(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create test commission rules (all active by default)
	rule1 := CreateTestCommissionRule(t, db, model.CommissionRuleTypeDefault, 20)
	rule2 := CreateTestCommissionRule(t, db, model.CommissionRuleTypeSpecial, 15)
	rule3 := CreateTestCommissionRule(t, db, model.CommissionRuleTypeGift, 25)

	// Create repository and service
	commissionRepo := commission.NewCommissionRepository(db)
	service := commissionservice.NewCommissionService(commissionRepo, &mockOrderReader{}, &mockPlayerRepo{})

	// Test: Disable rules
	ids := []uint64{rule1.ID, rule2.ID, rule3.ID}

	result, err := service.BatchUpdateCommissionRuleStatus(ctx, ids, false)
	if err != nil {
		t.Fatalf("BatchUpdateCommissionRuleStatus failed: %v", err)
	}

	// Assertions
	if result.SuccessCount != 3 {
		t.Errorf("Expected 3 successful updates, got %d", result.SuccessCount)
	}
	if result.FailedCount != 0 {
		t.Errorf("Expected 0 failed updates, got %d", result.FailedCount)
	}

	// Verify database state
	var rules []model.CommissionRule
	db.Find(&rules, ids)
	for _, r := range rules {
		if r.IsActive {
			t.Errorf("Rule %d should be inactive, but IsActive is true", r.ID)
		}
	}
}

func TestCommissionService_BatchUpdateCommissionRuleStatus_Enable(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create test commission rules and set them inactive
	rule1 := CreateTestCommissionRule(t, db, model.CommissionRuleTypeDefault, 20)
	rule2 := CreateTestCommissionRule(t, db, model.CommissionRuleTypeSpecial, 15)

	db.Model(&model.CommissionRule{}).Where("id = ?", rule1.ID).Update("is_active", false)
	db.Model(&model.CommissionRule{}).Where("id = ?", rule2.ID).Update("is_active", false)

	// Create repository and service
	commissionRepo := commission.NewCommissionRepository(db)
	service := commissionservice.NewCommissionService(commissionRepo, &mockOrderReader{}, &mockPlayerRepo{})

	// Test: Enable rules
	ids := []uint64{rule1.ID, rule2.ID}

	result, err := service.BatchUpdateCommissionRuleStatus(ctx, ids, true)
	if err != nil {
		t.Fatalf("BatchUpdateCommissionRuleStatus failed: %v", err)
	}

	// Assertions
	if result.SuccessCount != 2 {
		t.Errorf("Expected 2 successful updates, got %d", result.SuccessCount)
	}

	// Verify database state
	var rules []model.CommissionRule
	db.Find(&rules, ids)
	for _, r := range rules {
		if !r.IsActive {
			t.Errorf("Rule %d should be active, but IsActive is false", r.ID)
		}
	}
}

func TestCommissionService_BatchUpdateCommissionRuleStatus_PartialNonExistent(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create test commission rule
	rule1 := CreateTestCommissionRule(t, db, model.CommissionRuleTypeDefault, 20)

	// Create repository and service
	commissionRepo := commission.NewCommissionRepository(db)
	service := commissionservice.NewCommissionService(commissionRepo, &mockOrderReader{}, &mockPlayerRepo{})

	// Test: Include non-existent ID
	nonExistentID := uint64(99999)
	ids := []uint64{rule1.ID, nonExistentID}

	result, err := service.BatchUpdateCommissionRuleStatus(ctx, ids, false)
	if err != nil {
		t.Fatalf("BatchUpdateCommissionRuleStatus failed: %v", err)
	}

	// Assertions
	if result.SuccessCount != 1 {
		t.Errorf("Expected 1 successful update, got %d", result.SuccessCount)
	}
	if result.FailedCount != 1 {
		t.Errorf("Expected 1 failed update, got %d", result.FailedCount)
	}
	if len(result.FailedIDs) != 1 || result.FailedIDs[0] != nonExistentID {
		t.Errorf("FailedIDs should contain non-existent ID")
	}
}

// ============================================================================
// Referral Batch Operations Tests
// ============================================================================

func TestReferralService_BatchDeleteReferrals_PartialFailure(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create test users
	user1 := CreateUniqueTestUser(t, db, "referrer1")
	user2 := CreateUniqueTestUser(t, db, "referrer2")
	referee1 := CreateUniqueTestUser(t, db, "referee1")
	referee2 := CreateUniqueTestUser(t, db, "referee2")
	referee3 := CreateUniqueTestUser(t, db, "referee3")

	// Create test referrals
	referral1 := createTestReferral(t, db, user1.ID, referee1.ID, model.ReferralTypeUserToUser)
	referral2 := createTestReferral(t, db, user2.ID, referee2.ID, model.ReferralTypeUserToUser)
	referral3 := createTestReferral(t, db, user1.ID, referee3.ID, model.ReferralTypeUserToUser)

	// Create repository and service
	referralRepo := referral.NewReferralRepository(db)
	service := referralservice.NewReferralService(referralRepo)

	// Test: Delete with mix of valid and invalid IDs
	nonExistentID := uint64(99999)
	ids := []uint64{referral1.ID, referral2.ID, nonExistentID}

	result, err := service.BatchDeleteReferrals(ctx, ids)
	if err != nil {
		t.Fatalf("BatchDeleteReferrals failed: %v", err)
	}

	// Assertions
	if result.SuccessCount != 2 {
		t.Errorf("Expected 2 successful deletions, got %d", result.SuccessCount)
	}
	if result.FailedCount != 1 {
		t.Errorf("Expected 1 failed deletion, got %d", result.FailedCount)
	}
	if len(result.FailedIDs) != 1 {
		t.Errorf("Expected 1 failed ID, got %d", len(result.FailedIDs))
	}
	if result.TotalCount != 3 {
		t.Errorf("Expected total count 3, got %d", result.TotalCount)
	}

	// Verify database state
	var count int64
	db.Model(&model.Referral{}).Where("id IN ?", []uint64{referral1.ID, referral2.ID}).Count(&count)
	if count != 0 {
		t.Errorf("Expected referrals to be deleted, but found %d remaining", count)
	}

	// Verify referral3 still exists
	var remaining model.Referral
	if err := db.First(&remaining, referral3.ID).Error; err != nil {
		t.Errorf("Referral 3 should still exist: %v", err)
	}
}

func TestReferralService_BatchUpdateReferralsStatus_AllSuccess(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create test users
	user1 := CreateUniqueTestUser(t, db, "referrer1")
	referee1 := CreateUniqueTestUser(t, db, "referee1")
	referee2 := CreateUniqueTestUser(t, db, "referee2")
	referee3 := CreateUniqueTestUser(t, db, "referee3")

	// Create test referrals (all pending by default)
	referral1 := createTestReferral(t, db, user1.ID, referee1.ID, model.ReferralTypeUserToUser)
	referral2 := createTestReferral(t, db, user1.ID, referee2.ID, model.ReferralTypeUserToUser)
	referral3 := createTestReferral(t, db, user1.ID, referee3.ID, model.ReferralTypeUserToUser)

	// Create repository and service
	referralRepo := referral.NewReferralRepository(db)
	service := referralservice.NewReferralService(referralRepo)

	// Test: Complete referrals
	ids := []uint64{referral1.ID, referral2.ID, referral3.ID}

	result, err := service.BatchUpdateReferralsStatus(ctx, ids, model.ReferralStatusCompleted)
	if err != nil {
		t.Fatalf("BatchUpdateReferralsStatus failed: %v", err)
	}

	// Assertions
	if result.SuccessCount != 3 {
		t.Errorf("Expected 3 successful updates, got %d", result.SuccessCount)
	}
	if result.FailedCount != 0 {
		t.Errorf("Expected 0 failed updates, got %d", result.FailedCount)
	}
	if result.TotalCount != 3 {
		t.Errorf("Expected total count 3, got %d", result.TotalCount)
	}

	// Verify database state
	var referrals []model.Referral
	db.Find(&referrals, ids)
	for _, r := range referrals {
		if r.Status != model.ReferralStatusCompleted {
			t.Errorf("Referral %d status: expected %s, got %s", r.ID,
				model.ReferralStatusCompleted, r.Status)
		}
	}
}

func TestReferralService_BatchUpdateReferralsStatus_PartialNonExistent(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create test users and referral
	user1 := CreateUniqueTestUser(t, db, "referrer1")
	referee1 := CreateUniqueTestUser(t, db, "referee1")
	referral1 := createTestReferral(t, db, user1.ID, referee1.ID, model.ReferralTypeUserToUser)

	// Create repository and service
	referralRepo := referral.NewReferralRepository(db)
	service := referralservice.NewReferralService(referralRepo)

	// Test: Include non-existent ID
	nonExistentID := uint64(99999)
	ids := []uint64{referral1.ID, nonExistentID}

	result, err := service.BatchUpdateReferralsStatus(ctx, ids, model.ReferralStatusCompleted)
	if err != nil {
		t.Fatalf("BatchUpdateReferralsStatus failed: %v", err)
	}

	// Assertions
	if result.SuccessCount != 1 {
		t.Errorf("Expected 1 successful update, got %d", result.SuccessCount)
	}
	if result.FailedCount != 1 {
		t.Errorf("Expected 1 failed update, got %d", result.FailedCount)
	}
	if len(result.FailedIDs) != 1 || result.FailedIDs[0] != nonExistentID {
		t.Errorf("FailedIDs should contain non-existent ID")
	}
	if result.TotalCount != 2 {
		t.Errorf("Expected total count 2, got %d", result.TotalCount)
	}

	// Verify valid referral was updated
	var referral model.Referral
	db.First(&referral, referral1.ID)
	if referral.Status != model.ReferralStatusCompleted {
		t.Errorf("Referral status: expected %s, got %s",
			model.ReferralStatusCompleted, referral.Status)
	}
}

func TestReferralService_BatchUpdateReferralsStatus_Rewarded(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create test users
	user1 := CreateUniqueTestUser(t, db, "referrer1")
	referee1 := CreateUniqueTestUser(t, db, "referee1")
	referee2 := CreateUniqueTestUser(t, db, "referee2")

	// Create test referrals
	referral1 := createTestReferral(t, db, user1.ID, referee1.ID, model.ReferralTypeUserToUser)
	referral2 := createTestReferral(t, db, user1.ID, referee2.ID, model.ReferralTypeUserToUser)

	// Create repository and service
	referralRepo := referral.NewReferralRepository(db)
	service := referralservice.NewReferralService(referralRepo)

	// Test: Mark referrals as rewarded
	ids := []uint64{referral1.ID, referral2.ID}

	result, err := service.BatchUpdateReferralsStatus(ctx, ids, model.ReferralStatusRewarded)
	if err != nil {
		t.Fatalf("BatchUpdateReferralsStatus failed: %v", err)
	}

	// Assertions
	if result.SuccessCount != 2 {
		t.Errorf("Expected 2 successful updates, got %d", result.SuccessCount)
	}

	// Verify database state
	var referrals []model.Referral
	db.Find(&referrals, ids)
	for _, r := range referrals {
		if r.Status != model.ReferralStatusRewarded {
			t.Errorf("Referral %d status: expected %s, got %s", r.ID,
				model.ReferralStatusRewarded, r.Status)
		}
	}
}

// ============================================================================
// Helper Functions
// ============================================================================

// createTestActivity creates a test activity with given name and status
func createTestActivity(t *testing.T, db *gorm.DB, name string, status model.ActivityStatus) *model.Activity {
	t.Helper()
	now := time.Now()
	startAt := now.Add(time.Hour * 24)
	endAt := startAt.Add(time.Hour * 168) // 7 days

	activity := &model.Activity{
		Name:         name,
		Description:  "Test activity description",
		Status:       status,
		StartAt:      startAt,
		EndAt:        endAt,
		PreheatAt:    &now,
		IsVisible:    true,
		TotalLimit:   1000,
		PerUserLimit: 1,
		DailyLimit:   100,
	}

	if err := db.Create(activity).Error; err != nil {
		t.Fatalf("Failed to create test activity: %v", err)
	}
	return activity
}

// createTestReferral creates a test referral record
func createTestReferral(t *testing.T, db *gorm.DB, referrerID, refereeID uint64, refType model.ReferralType) *model.Referral {
	t.Helper()
	referral := &model.Referral{
		ReferrerID:       referrerID,
		RefereeID:        refereeID,
		Type:             refType,
		Level:            1,
		Status:           model.ReferralStatusPending,
		RefereeCondition: "registered",
	}

	if err := db.Create(referral).Error; err != nil {
		t.Fatalf("Failed to create test referral: %v", err)
	}
	return referral
}

// ============================================================================
// Mock Implementations for Commission Service Tests
// ============================================================================

// mockOrderReader is a minimal mock implementation of OrderReader
type mockOrderReader struct{}

func (m *mockOrderReader) Get(ctx context.Context, id uint64) (*model.Order, error) {
	gameID := uint64(1)
	order := &model.Order{
		TotalPriceCents: 10000,
		GameID:          &gameID,
		ItemID:          1,
	}
	order.ID = id
	return order, nil
}

// mockPlayerRepo is a minimal mock implementation of PlayerRepository
type mockPlayerRepo struct{}

func (m *mockPlayerRepo) Get(ctx context.Context, id uint64) (*model.Player, error) {
	player := &model.Player{}
	player.ID = id
	return player, nil
}

func (m *mockPlayerRepo) GetByUserID(ctx context.Context, userID uint64) (*model.Player, error) {
	player := &model.Player{
		UserID: userID,
	}
	player.ID = 1
	return player, nil
}

func (m *mockPlayerRepo) List(ctx context.Context) ([]model.Player, error) {
	return []model.Player{}, nil
}

func (m *mockPlayerRepo) ListPaged(ctx context.Context, page, pageSize int) ([]model.Player, int64, error) {
	return []model.Player{}, 0, nil
}

func (m *mockPlayerRepo) ListPagedWithFilter(ctx context.Context, page, pageSize int, keyword string, status *model.VerificationStatus) ([]model.Player, int64, error) {
	return []model.Player{}, 0, nil
}

func (m *mockPlayerRepo) ListFeatured(ctx context.Context, limit int, status *model.VerificationStatus) ([]model.Player, int64, error) {
	return []model.Player{}, 0, nil
}

func (m *mockPlayerRepo) Create(ctx context.Context, player *model.Player) error {
	return nil
}

func (m *mockPlayerRepo) Update(ctx context.Context, player *model.Player) error {
	return nil
}

func (m *mockPlayerRepo) Delete(ctx context.Context, id uint64) error {
	return nil
}

func (m *mockPlayerRepo) BatchUpdateRank(ctx context.Context, ids []uint64, rank string) (int64, error) {
	return int64(len(ids)), nil
}

func (m *mockPlayerRepo) BatchUpdateHourlyRate(ctx context.Context, ids []uint64, rateCents int64) (int64, error) {
	return int64(len(ids)), nil
}

func (m *mockPlayerRepo) BatchUpdateStatus(ctx context.Context, ids []uint64, status model.VerificationStatus) (int64, error) {
	return int64(len(ids)), nil
}

func (m *mockPlayerRepo) BatchDelete(ctx context.Context, ids []uint64) (int64, error) {
	return int64(len(ids)), nil
}

func (m *mockPlayerRepo) GetByIDs(ctx context.Context, ids []uint64) ([]model.Player, error) {
	players := make([]model.Player, len(ids))
	for i, id := range ids {
		players[i] = model.Player{}
		players[i].ID = id
	}
	return players, nil
}
