package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository/activity"
	activityservice "gamelink/internal/service/activity"
)

// setupActivityService creates an Activity service for testing.
func setupActivityService(t *testing.T) *activityservice.Service {
	t.Helper()
	db := SetupTestDB(t)
	repo := activity.NewActivityRepository(db)
	// Create service without coupon dependency for basic tests
	return activityservice.NewActivityService(repo, nil)
}

// ============================================================================
// Activity CRUD Tests
// ============================================================================

func TestActivityService_CreateActivity(t *testing.T) {
	SkipIfNoTestDB(t)
	svc := setupActivityService(t)
	ctx := context.Background()

	now := time.Now()
	act := &model.Activity{
		Name:        "Test Activity",
		Description: "Test activity description",
		Type:        model.ActivityTypeCoupon,
		Status:      model.ActivityStatusDraft,
		StartAt:     now.Add(time.Hour),
		EndAt:       now.Add(24 * time.Hour),
		IsVisible:   true,
	}

	err := svc.CreateActivity(ctx, act)
	require.NoError(t, err)
	assert.NotZero(t, act.ID)

	// Verify
	got, err := svc.GetActivity(ctx, act.ID)
	require.NoError(t, err)
	assert.Equal(t, "Test Activity", got.Name)
	assert.Equal(t, model.ActivityStatusDraft, got.Status)
}

func TestActivityService_CreateActivity_InvalidTime(t *testing.T) {
	SkipIfNoTestDB(t)
	svc := setupActivityService(t)
	ctx := context.Background()

	now := time.Now()
	act := &model.Activity{
		Name:    "Invalid Activity",
		Type:    model.ActivityTypeCoupon,
		StartAt: now.Add(24 * time.Hour),
		EndAt:   now.Add(time.Hour), // End before start
	}

	err := svc.CreateActivity(ctx, act)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "结束时间必须晚于开始时间")
}

func TestActivityService_UpdateActivity(t *testing.T) {
	SkipIfNoTestDB(t)
	svc := setupActivityService(t)
	ctx := context.Background()

	now := time.Now()
	act := &model.Activity{
		Name:    "Update Test",
		Type:    model.ActivityTypeCoupon,
		Status:  model.ActivityStatusDraft,
		StartAt: now.Add(time.Hour),
		EndAt:   now.Add(24 * time.Hour),
	}
	require.NoError(t, svc.CreateActivity(ctx, act))

	// Update
	act.Name = "Updated Name"
	act.Description = "Updated description"
	err := svc.UpdateActivity(ctx, act)
	require.NoError(t, err)

	// Verify
	got, err := svc.GetActivity(ctx, act.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", got.Name)
	assert.Equal(t, "Updated description", got.Description)
}

func TestActivityService_DeleteActivity(t *testing.T) {
	SkipIfNoTestDB(t)
	svc := setupActivityService(t)
	ctx := context.Background()

	now := time.Now()
	act := &model.Activity{
		Name:    "Delete Test",
		Type:    model.ActivityTypeCoupon,
		Status:  model.ActivityStatusDraft,
		StartAt: now.Add(time.Hour),
		EndAt:   now.Add(24 * time.Hour),
	}
	require.NoError(t, svc.CreateActivity(ctx, act))

	err := svc.DeleteActivity(ctx, act.ID)
	require.NoError(t, err)

	// Verify deleted
	_, err = svc.GetActivity(ctx, act.ID)
	assert.Error(t, err)
}

func TestActivityService_DeleteActivity_ActiveNotAllowed(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := activity.NewActivityRepository(db)
	svc := activityservice.NewActivityService(repo, nil)
	ctx := context.Background()

	now := time.Now()
	act := &model.Activity{
		Name:    "Active Activity",
		Type:    model.ActivityTypeCoupon,
		Status:  model.ActivityStatusDraft,
		StartAt: now.Add(time.Hour),
		EndAt:   now.Add(24 * time.Hour),
	}
	require.NoError(t, svc.CreateActivity(ctx, act))

	// Update to active status
	require.NoError(t, db.Model(act).Update("status", model.ActivityStatusActive).Error)

	err := svc.DeleteActivity(ctx, act.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "进行中的活动不能删除")
}

// ============================================================================
// Activity Status Tests
// ============================================================================

func TestActivityService_UpdateActivityStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	svc := setupActivityService(t)
	ctx := context.Background()

	now := time.Now()
	act := &model.Activity{
		Name:    "Status Test",
		Type:    model.ActivityTypeCoupon,
		Status:  model.ActivityStatusDraft,
		StartAt: now.Add(time.Hour),
		EndAt:   now.Add(24 * time.Hour),
	}
	require.NoError(t, svc.CreateActivity(ctx, act))

	// Draft -> Active
	err := svc.UpdateActivityStatus(ctx, act.ID, model.ActivityStatusActive)
	require.NoError(t, err)

	got, err := svc.GetActivity(ctx, act.ID)
	require.NoError(t, err)
	assert.Equal(t, model.ActivityStatusActive, got.Status)
}

func TestActivityService_UpdateActivityStatus_InvalidTransition(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := activity.NewActivityRepository(db)
	svc := activityservice.NewActivityService(repo, nil)
	ctx := context.Background()

	now := time.Now()
	act := &model.Activity{
		Name:    "Invalid Transition",
		Type:    model.ActivityTypeCoupon,
		Status:  model.ActivityStatusDraft,
		StartAt: now.Add(time.Hour),
		EndAt:   now.Add(24 * time.Hour),
	}
	require.NoError(t, svc.CreateActivity(ctx, act))

	// Update to ended status directly
	require.NoError(t, db.Model(act).Update("status", model.ActivityStatusEnded).Error)

	// Ended -> Active (invalid)
	err := svc.UpdateActivityStatus(ctx, act.ID, model.ActivityStatusActive)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "不允许")
}

// ============================================================================
// Activity Reward Tests
// ============================================================================

func TestActivityService_CreateReward(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := activity.NewActivityRepository(db)
	svc := activityservice.NewActivityService(repo, nil)
	ctx := context.Background()

	// Create coupon template first
	template := CreateTestCouponTemplate(t, db, "Activity Reward Template", 500)

	now := time.Now()
	act := &model.Activity{
		Name:    "Reward Test",
		Type:    model.ActivityTypeCoupon,
		Status:  model.ActivityStatusDraft,
		StartAt: now.Add(time.Hour),
		EndAt:   now.Add(24 * time.Hour),
	}
	require.NoError(t, svc.CreateActivity(ctx, act))

	reward := &model.ActivityReward{
		ActivityID:       act.ID,
		CouponTemplateID: template.ID,
		CouponCount:      1,
		Probability:      100,
		TotalStock:       100,
	}

	err := svc.CreateReward(ctx, reward)
	require.NoError(t, err)
	assert.NotZero(t, reward.ID)
	assert.Equal(t, 100, reward.RemainingStock)
}

func TestActivityService_GetRewardsByActivityID(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := activity.NewActivityRepository(db)
	svc := activityservice.NewActivityService(repo, nil)
	ctx := context.Background()

	template := CreateTestCouponTemplate(t, db, "Rewards List Template", 500)

	now := time.Now()
	act := &model.Activity{
		Name:    "Rewards List Test",
		Type:    model.ActivityTypeCoupon,
		Status:  model.ActivityStatusDraft,
		StartAt: now.Add(time.Hour),
		EndAt:   now.Add(24 * time.Hour),
	}
	require.NoError(t, svc.CreateActivity(ctx, act))

	// Create multiple rewards
	for i := 0; i < 3; i++ {
		reward := &model.ActivityReward{
			ActivityID:       act.ID,
			CouponTemplateID: template.ID,
			CouponCount:      i + 1,
			Probability:      100,
			TotalStock:       50,
			SortOrder:        i,
		}
		require.NoError(t, svc.CreateReward(ctx, reward))
	}

	rewards, err := svc.GetRewardsByActivityID(ctx, act.ID)
	require.NoError(t, err)
	assert.Len(t, rewards, 3)
}

// ============================================================================
// Activity Participation Tests
// ============================================================================

func TestActivityService_ParticipateActivity(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := activity.NewActivityRepository(db)
	svc := activityservice.NewActivityService(repo, nil)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "activity_participate")
	template := CreateTestCouponTemplate(t, db, "Participate Template", 500)

	now := time.Now()
	act := &model.Activity{
		Name:      "Participate Test",
		Type:      model.ActivityTypeCoupon,
		Status:    model.ActivityStatusActive,
		StartAt:   now.Add(-time.Hour),
		EndAt:     now.Add(24 * time.Hour),
		IsVisible: true,
	}
	require.NoError(t, svc.CreateActivity(ctx, act))

	reward := &model.ActivityReward{
		ActivityID:       act.ID,
		CouponTemplateID: template.ID,
		CouponCount:      1,
		Probability:      100,
		TotalStock:       100,
	}
	require.NoError(t, svc.CreateReward(ctx, reward))

	// Participate
	participation, err := svc.ParticipateActivity(ctx, user.ID, act.ID, reward.ID, "127.0.0.1")
	require.NoError(t, err)
	assert.NotZero(t, participation.ID)
	assert.Equal(t, user.ID, participation.UserID)
	assert.Equal(t, act.ID, participation.ActivityID)
}

func TestActivityService_ParticipateActivity_NotActive(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := activity.NewActivityRepository(db)
	svc := activityservice.NewActivityService(repo, nil)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "activity_not_active")
	template := CreateTestCouponTemplate(t, db, "Not Active Template", 500)

	now := time.Now()
	act := &model.Activity{
		Name:    "Not Active Test",
		Type:    model.ActivityTypeCoupon,
		Status:  model.ActivityStatusDraft, // Not active
		StartAt: now.Add(time.Hour),
		EndAt:   now.Add(24 * time.Hour),
	}
	require.NoError(t, svc.CreateActivity(ctx, act))

	reward := &model.ActivityReward{
		ActivityID:       act.ID,
		CouponTemplateID: template.ID,
		CouponCount:      1,
		Probability:      100,
	}
	require.NoError(t, svc.CreateReward(ctx, reward))

	_, err := svc.ParticipateActivity(ctx, user.ID, act.ID, reward.ID, "127.0.0.1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "未开始或已结束")
}

func TestActivityService_ParticipateActivity_PerUserLimit(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := activity.NewActivityRepository(db)
	svc := activityservice.NewActivityService(repo, nil)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "activity_limit")
	template := CreateTestCouponTemplate(t, db, "Limit Template", 500)

	now := time.Now()
	act := &model.Activity{
		Name:         "Limit Test",
		Type:         model.ActivityTypeCoupon,
		Status:       model.ActivityStatusActive,
		StartAt:      now.Add(-time.Hour),
		EndAt:        now.Add(24 * time.Hour),
		PerUserLimit: 1, // Only 1 per user
	}
	require.NoError(t, svc.CreateActivity(ctx, act))

	reward := &model.ActivityReward{
		ActivityID:       act.ID,
		CouponTemplateID: template.ID,
		CouponCount:      1,
		Probability:      100,
		TotalStock:       100,
	}
	require.NoError(t, svc.CreateReward(ctx, reward))

	// First participation - should succeed
	_, err := svc.ParticipateActivity(ctx, user.ID, act.ID, reward.ID, "127.0.0.1")
	require.NoError(t, err)

	// Second participation - should fail
	_, err = svc.ParticipateActivity(ctx, user.ID, act.ID, reward.ID, "127.0.0.1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "参与上限")
}

func TestActivityService_ParticipateActivity_StockExhausted(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := activity.NewActivityRepository(db)
	svc := activityservice.NewActivityService(repo, nil)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "activity_stock")
	template := CreateTestCouponTemplate(t, db, "Stock Template", 500)

	now := time.Now()
	act := &model.Activity{
		Name:    "Stock Test",
		Type:    model.ActivityTypeCoupon,
		Status:  model.ActivityStatusActive,
		StartAt: now.Add(-time.Hour),
		EndAt:   now.Add(24 * time.Hour),
	}
	require.NoError(t, svc.CreateActivity(ctx, act))

	reward := &model.ActivityReward{
		ActivityID:       act.ID,
		CouponTemplateID: template.ID,
		CouponCount:      1,
		Probability:      100,
		TotalStock:       1,
	}
	require.NoError(t, svc.CreateReward(ctx, reward))

	// Manually set remaining stock to 0 to simulate exhausted
	require.NoError(t, db.Model(reward).Update("remaining_stock", 0).Error)

	_, err := svc.ParticipateActivity(ctx, user.ID, act.ID, reward.ID, "127.0.0.1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "库存不足")
}

func TestActivityService_GetUserParticipations(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := activity.NewActivityRepository(db)
	svc := activityservice.NewActivityService(repo, nil)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "activity_history")
	template := CreateTestCouponTemplate(t, db, "History Template", 500)

	now := time.Now()
	act := &model.Activity{
		Name:    "History Test",
		Type:    model.ActivityTypeCoupon,
		Status:  model.ActivityStatusActive,
		StartAt: now.Add(-time.Hour),
		EndAt:   now.Add(24 * time.Hour),
	}
	require.NoError(t, svc.CreateActivity(ctx, act))

	reward := &model.ActivityReward{
		ActivityID:       act.ID,
		CouponTemplateID: template.ID,
		CouponCount:      1,
		Probability:      100,
		TotalStock:       100,
	}
	require.NoError(t, svc.CreateReward(ctx, reward))

	// Participate multiple times (different activities would be needed for real scenario)
	_, err := svc.ParticipateActivity(ctx, user.ID, act.ID, reward.ID, "127.0.0.1")
	require.NoError(t, err)

	participations, err := svc.GetUserParticipations(ctx, user.ID, 10)
	require.NoError(t, err)
	assert.Len(t, participations, 1)
}

// ============================================================================
// Activity Stats Tests
// ============================================================================

func TestActivityService_GetActivityStats(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := activity.NewActivityRepository(db)
	svc := activityservice.NewActivityService(repo, nil)
	ctx := context.Background()

	template := CreateTestCouponTemplate(t, db, "Stats Template", 500)

	now := time.Now()
	act := &model.Activity{
		Name:    "Stats Test",
		Type:    model.ActivityTypeCoupon,
		Status:  model.ActivityStatusActive,
		StartAt: now.Add(-time.Hour),
		EndAt:   now.Add(24 * time.Hour),
	}
	require.NoError(t, svc.CreateActivity(ctx, act))

	reward := &model.ActivityReward{
		ActivityID:       act.ID,
		CouponTemplateID: template.ID,
		CouponCount:      1,
		Probability:      100,
		TotalStock:       100,
	}
	require.NoError(t, svc.CreateReward(ctx, reward))

	stats, err := svc.GetActivityStats(ctx, act.ID)
	require.NoError(t, err)
	assert.NotNil(t, stats)
}
