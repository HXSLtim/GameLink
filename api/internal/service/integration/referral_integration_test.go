package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository/referral"
)

// ============================================================================
// Referral Config Tests
// ============================================================================

func TestReferralRepository_SetAndGetConfig(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := referral.NewReferralRepository(db)
	ctx := context.Background()

	// Set config
	err := repo.SetConfig(ctx, "test_key", "test_value", "Test description")
	require.NoError(t, err)

	// Get config
	config, err := repo.GetConfig(ctx, "test_key")
	require.NoError(t, err)
	assert.Equal(t, "test_value", config.ConfigValue)
	assert.Equal(t, "Test description", config.Description)
}

func TestReferralRepository_GetAllConfigs(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := referral.NewReferralRepository(db)
	ctx := context.Background()

	// Set multiple configs
	require.NoError(t, repo.SetConfig(ctx, "config_a", "value_a", "Desc A"))
	require.NoError(t, repo.SetConfig(ctx, "config_b", "value_b", "Desc B"))

	configs, err := repo.GetAllConfigs(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(configs), 2)
}

func TestReferralRepository_UpdateConfig(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := referral.NewReferralRepository(db)
	ctx := context.Background()

	// Set initial config
	require.NoError(t, repo.SetConfig(ctx, "update_key", "initial", "Initial desc"))

	// Update config
	err := repo.SetConfig(ctx, "update_key", "updated", "Updated desc")
	require.NoError(t, err)

	// Verify update
	config, err := repo.GetConfig(ctx, "update_key")
	require.NoError(t, err)
	assert.Equal(t, "updated", config.ConfigValue)
}

// ============================================================================
// Referral Code Tests
// ============================================================================

func TestReferralRepository_CreateCode(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := referral.NewReferralRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "referral_code")

	code := &model.ReferralCode{
		UserID:   user.ID,
		Type:     model.ReferralTypeUserToUser,
		IsActive: true,
		MaxUse:   100,
	}

	err := repo.CreateCode(ctx, code)
	require.NoError(t, err)
	assert.NotZero(t, code.ID)
	assert.NotEmpty(t, code.Code) // Auto-generated
}

func TestReferralRepository_GetCodeByID(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := referral.NewReferralRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "get_code")

	code := &model.ReferralCode{
		UserID:   user.ID,
		Type:     model.ReferralTypeUserToUser,
		IsActive: true,
	}
	require.NoError(t, repo.CreateCode(ctx, code))

	got, err := repo.GetCodeByID(ctx, code.ID)
	require.NoError(t, err)
	assert.Equal(t, code.Code, got.Code)
	assert.Equal(t, user.ID, got.UserID)
}

func TestReferralRepository_GetCodeByCode(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := referral.NewReferralRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "get_by_code")

	code := &model.ReferralCode{
		Code:     "TESTCODE123",
		UserID:   user.ID,
		Type:     model.ReferralTypeUserToUser,
		IsActive: true,
	}
	require.NoError(t, repo.CreateCode(ctx, code))

	got, err := repo.GetCodeByCode(ctx, "TESTCODE123")
	require.NoError(t, err)
	assert.Equal(t, code.ID, got.ID)
}

func TestReferralRepository_GetUserCode(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := referral.NewReferralRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "user_code")

	code := &model.ReferralCode{
		UserID:   user.ID,
		Type:     model.ReferralTypeUserToUser,
		IsActive: true,
	}
	require.NoError(t, repo.CreateCode(ctx, code))

	got, err := repo.GetUserCode(ctx, user.ID, model.ReferralTypeUserToUser)
	require.NoError(t, err)
	assert.Equal(t, code.ID, got.ID)
}

func TestReferralRepository_ListCodes(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := referral.NewReferralRepository(db)
	ctx := context.Background()

	// Create multiple codes
	for i := 0; i < 3; i++ {
		user := CreateUniqueTestUser(t, db, "list_codes")
		code := &model.ReferralCode{
			UserID:   user.ID,
			Type:     model.ReferralTypeUserToUser,
			IsActive: true,
		}
		require.NoError(t, repo.CreateCode(ctx, code))
	}

	codes, total, err := repo.ListCodes(ctx, referral.CodeListOptions{
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(3))
	assert.GreaterOrEqual(t, len(codes), 3)
}

func TestReferralRepository_UpdateCode(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := referral.NewReferralRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "update_code")

	code := &model.ReferralCode{
		UserID:   user.ID,
		Type:     model.ReferralTypeUserToUser,
		IsActive: true,
		MaxUse:   50,
	}
	require.NoError(t, repo.CreateCode(ctx, code))

	// Update
	code.MaxUse = 100
	code.IsActive = false
	err := repo.UpdateCode(ctx, code)
	require.NoError(t, err)

	got, err := repo.GetCodeByID(ctx, code.ID)
	require.NoError(t, err)
	assert.Equal(t, 100, got.MaxUse)
}

func TestReferralRepository_DeleteCode(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := referral.NewReferralRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "delete_code")

	code := &model.ReferralCode{
		UserID:   user.ID,
		Type:     model.ReferralTypeUserToUser,
		IsActive: true,
	}
	require.NoError(t, repo.CreateCode(ctx, code))

	err := repo.DeleteCode(ctx, code.ID)
	require.NoError(t, err)

	_, err = repo.GetCodeByID(ctx, code.ID)
	assert.Error(t, err)
}

func TestReferralRepository_IncrementCodeUseCount(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := referral.NewReferralRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "increment_code")

	code := &model.ReferralCode{
		UserID:   user.ID,
		Type:     model.ReferralTypeUserToUser,
		IsActive: true,
		UseCount: 0,
	}
	require.NoError(t, repo.CreateCode(ctx, code))

	err := repo.IncrementCodeUseCount(ctx, code.ID)
	require.NoError(t, err)

	got, err := repo.GetCodeByID(ctx, code.ID)
	require.NoError(t, err)
	assert.Equal(t, 1, got.UseCount)
}


// ============================================================================
// Referral Record Tests
// ============================================================================

func TestReferralRepository_CreateReferral(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := referral.NewReferralRepository(db)
	ctx := context.Background()

	referrer := CreateUniqueTestUser(t, db, "referrer")
	referee := CreateUniqueTestUser(t, db, "referee")

	code := &model.ReferralCode{
		UserID:   referrer.ID,
		Type:     model.ReferralTypeUserToUser,
		IsActive: true,
	}
	require.NoError(t, repo.CreateCode(ctx, code))

	ref := &model.Referral{
		ReferrerID:       referrer.ID,
		RefereeID:        referee.ID,
		CodeID:           &code.ID,
		Type:             model.ReferralTypeUserToUser,
		Level:            1,
		Status:           model.ReferralStatusPending,
		RefereeCondition: "registered",
	}

	err := repo.CreateReferral(ctx, ref)
	require.NoError(t, err)
	assert.NotZero(t, ref.ID)
}

func TestReferralRepository_GetReferralByID(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := referral.NewReferralRepository(db)
	ctx := context.Background()

	referrer := CreateUniqueTestUser(t, db, "get_ref_referrer")
	referee := CreateUniqueTestUser(t, db, "get_ref_referee")

	ref := &model.Referral{
		ReferrerID: referrer.ID,
		RefereeID:  referee.ID,
		Type:       model.ReferralTypeUserToUser,
		Status:     model.ReferralStatusPending,
	}
	require.NoError(t, repo.CreateReferral(ctx, ref))

	got, err := repo.GetReferralByID(ctx, ref.ID)
	require.NoError(t, err)
	assert.Equal(t, referrer.ID, got.ReferrerID)
	assert.Equal(t, referee.ID, got.RefereeID)
}

func TestReferralRepository_GetReferralByReferee(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := referral.NewReferralRepository(db)
	ctx := context.Background()

	referrer := CreateUniqueTestUser(t, db, "by_referee_referrer")
	referee := CreateUniqueTestUser(t, db, "by_referee_referee")

	ref := &model.Referral{
		ReferrerID: referrer.ID,
		RefereeID:  referee.ID,
		Type:       model.ReferralTypeUserToUser,
		Status:     model.ReferralStatusPending,
	}
	require.NoError(t, repo.CreateReferral(ctx, ref))

	got, err := repo.GetReferralByReferee(ctx, referee.ID)
	require.NoError(t, err)
	assert.Equal(t, ref.ID, got.ID)
}

func TestReferralRepository_ListReferrals(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := referral.NewReferralRepository(db)
	ctx := context.Background()

	referrer := CreateUniqueTestUser(t, db, "list_ref_referrer")

	// Create multiple referrals
	for i := 0; i < 3; i++ {
		referee := CreateUniqueTestUser(t, db, "list_ref_referee")
		ref := &model.Referral{
			ReferrerID: referrer.ID,
			RefereeID:  referee.ID,
			Type:       model.ReferralTypeUserToUser,
			Status:     model.ReferralStatusPending,
		}
		require.NoError(t, repo.CreateReferral(ctx, ref))
	}

	refs, total, err := repo.ListReferrals(ctx, referral.ReferralListOptions{
		Page:       1,
		PageSize:   10,
		ReferrerID: &referrer.ID,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, refs, 3)
}

func TestReferralRepository_UpdateReferralStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := referral.NewReferralRepository(db)
	ctx := context.Background()

	referrer := CreateUniqueTestUser(t, db, "status_referrer")
	referee := CreateUniqueTestUser(t, db, "status_referee")

	ref := &model.Referral{
		ReferrerID: referrer.ID,
		RefereeID:  referee.ID,
		Type:       model.ReferralTypeUserToUser,
		Status:     model.ReferralStatusPending,
	}
	require.NoError(t, repo.CreateReferral(ctx, ref))

	// Update to completed
	err := repo.UpdateReferralStatus(ctx, ref.ID, model.ReferralStatusCompleted)
	require.NoError(t, err)

	got, err := repo.GetReferralByID(ctx, ref.ID)
	require.NoError(t, err)
	assert.Equal(t, model.ReferralStatusCompleted, got.Status)
	assert.NotNil(t, got.CompletedAt)
}

func TestReferralRepository_GetUserReferrals(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := referral.NewReferralRepository(db)
	ctx := context.Background()

	referrer := CreateUniqueTestUser(t, db, "user_refs_referrer")

	for i := 0; i < 2; i++ {
		referee := CreateUniqueTestUser(t, db, "user_refs_referee")
		ref := &model.Referral{
			ReferrerID: referrer.ID,
			RefereeID:  referee.ID,
			Type:       model.ReferralTypeUserToUser,
			Status:     model.ReferralStatusCompleted,
		}
		require.NoError(t, repo.CreateReferral(ctx, ref))
	}

	refs, err := repo.GetUserReferrals(ctx, referrer.ID, 10)
	require.NoError(t, err)
	assert.Len(t, refs, 2)
}

func TestReferralRepository_CountUserReferrals(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := referral.NewReferralRepository(db)
	ctx := context.Background()

	referrer := CreateUniqueTestUser(t, db, "count_refs_referrer")

	for i := 0; i < 3; i++ {
		referee := CreateUniqueTestUser(t, db, "count_refs_referee")
		ref := &model.Referral{
			ReferrerID: referrer.ID,
			RefereeID:  referee.ID,
			Type:       model.ReferralTypeUserToUser,
			Status:     model.ReferralStatusCompleted,
		}
		require.NoError(t, repo.CreateReferral(ctx, ref))
	}

	count, err := repo.CountUserReferrals(ctx, referrer.ID, nil)
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)

	// Count by status
	status := model.ReferralStatusCompleted
	count, err = repo.CountUserReferrals(ctx, referrer.ID, &status)
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
}

// ============================================================================
// Referral Reward Tests
// ============================================================================

func TestReferralRepository_CreateReward(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := referral.NewReferralRepository(db)
	ctx := context.Background()

	referrer := CreateUniqueTestUser(t, db, "reward_referrer")
	referee := CreateUniqueTestUser(t, db, "reward_referee")

	ref := &model.Referral{
		ReferrerID: referrer.ID,
		RefereeID:  referee.ID,
		Type:       model.ReferralTypeUserToUser,
		Status:     model.ReferralStatusCompleted,
	}
	require.NoError(t, repo.CreateReferral(ctx, ref))

	reward := &model.ReferralReward{
		ReferralID:  ref.ID,
		UserID:      referrer.ID,
		Type:        model.RewardTypeCash,
		AmountCents: 1000,
		Status:      model.ReferralRewardStatusPending,
	}

	err := repo.CreateReward(ctx, reward)
	require.NoError(t, err)
	assert.NotZero(t, reward.ID)
}

func TestReferralRepository_GetRewardByID(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := referral.NewReferralRepository(db)
	ctx := context.Background()

	referrer := CreateUniqueTestUser(t, db, "get_reward_referrer")
	referee := CreateUniqueTestUser(t, db, "get_reward_referee")

	ref := &model.Referral{
		ReferrerID: referrer.ID,
		RefereeID:  referee.ID,
		Type:       model.ReferralTypeUserToUser,
		Status:     model.ReferralStatusCompleted,
	}
	require.NoError(t, repo.CreateReferral(ctx, ref))

	reward := &model.ReferralReward{
		ReferralID:  ref.ID,
		UserID:      referrer.ID,
		Type:        model.RewardTypeCash,
		AmountCents: 500,
		Status:      model.ReferralRewardStatusPending,
	}
	require.NoError(t, repo.CreateReward(ctx, reward))

	got, err := repo.GetRewardByID(ctx, reward.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(500), got.AmountCents)
}

func TestReferralRepository_ListRewards(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := referral.NewReferralRepository(db)
	ctx := context.Background()

	referrer := CreateUniqueTestUser(t, db, "list_rewards_referrer")

	// Create multiple rewards
	for i := 0; i < 3; i++ {
		referee := CreateUniqueTestUser(t, db, "list_rewards_referee")
		ref := &model.Referral{
			ReferrerID: referrer.ID,
			RefereeID:  referee.ID,
			Type:       model.ReferralTypeUserToUser,
			Status:     model.ReferralStatusCompleted,
		}
		require.NoError(t, repo.CreateReferral(ctx, ref))

		reward := &model.ReferralReward{
			ReferralID:  ref.ID,
			UserID:      referrer.ID,
			Type:        model.RewardTypeCash,
			AmountCents: 1000,
			Status:      model.ReferralRewardStatusPending,
		}
		require.NoError(t, repo.CreateReward(ctx, reward))
	}

	rewards, total, err := repo.ListRewards(ctx, referral.RewardListOptions{
		Page:     1,
		PageSize: 10,
		UserID:   &referrer.ID,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, rewards, 3)
}

func TestReferralRepository_UpdateRewardStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := referral.NewReferralRepository(db)
	ctx := context.Background()

	referrer := CreateUniqueTestUser(t, db, "update_reward_referrer")
	referee := CreateUniqueTestUser(t, db, "update_reward_referee")

	ref := &model.Referral{
		ReferrerID: referrer.ID,
		RefereeID:  referee.ID,
		Type:       model.ReferralTypeUserToUser,
		Status:     model.ReferralStatusCompleted,
	}
	require.NoError(t, repo.CreateReferral(ctx, ref))

	reward := &model.ReferralReward{
		ReferralID:  ref.ID,
		UserID:      referrer.ID,
		Type:        model.RewardTypeCash,
		AmountCents: 1000,
		Status:      model.ReferralRewardStatusPending,
	}
	require.NoError(t, repo.CreateReward(ctx, reward))

	// Update to issued
	err := repo.UpdateRewardStatus(ctx, reward.ID, model.ReferralRewardStatusIssued, "")
	require.NoError(t, err)

	got, err := repo.GetRewardByID(ctx, reward.ID)
	require.NoError(t, err)
	assert.Equal(t, model.ReferralRewardStatusIssued, got.Status)
	assert.NotNil(t, got.IssuedAt)
}

func TestReferralRepository_UpdateRewardStatus_Failed(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := referral.NewReferralRepository(db)
	ctx := context.Background()

	referrer := CreateUniqueTestUser(t, db, "failed_reward_referrer")
	referee := CreateUniqueTestUser(t, db, "failed_reward_referee")

	ref := &model.Referral{
		ReferrerID: referrer.ID,
		RefereeID:  referee.ID,
		Type:       model.ReferralTypeUserToUser,
		Status:     model.ReferralStatusCompleted,
	}
	require.NoError(t, repo.CreateReferral(ctx, ref))

	reward := &model.ReferralReward{
		ReferralID:  ref.ID,
		UserID:      referrer.ID,
		Type:        model.RewardTypeCash,
		AmountCents: 1000,
		Status:      model.ReferralRewardStatusPending,
	}
	require.NoError(t, repo.CreateReward(ctx, reward))

	// Update to failed with reason
	err := repo.UpdateRewardStatus(ctx, reward.ID, model.ReferralRewardStatusFailed, "Insufficient balance")
	require.NoError(t, err)

	got, err := repo.GetRewardByID(ctx, reward.ID)
	require.NoError(t, err)
	assert.Equal(t, model.ReferralRewardStatusFailed, got.Status)
	assert.Equal(t, "Insufficient balance", got.FailureReason)
}

func TestReferralRepository_GetUserRewards(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := referral.NewReferralRepository(db)
	ctx := context.Background()

	referrer := CreateUniqueTestUser(t, db, "user_rewards_referrer")

	for i := 0; i < 2; i++ {
		referee := CreateUniqueTestUser(t, db, "user_rewards_referee")
		ref := &model.Referral{
			ReferrerID: referrer.ID,
			RefereeID:  referee.ID,
			Type:       model.ReferralTypeUserToUser,
			Status:     model.ReferralStatusCompleted,
		}
		require.NoError(t, repo.CreateReferral(ctx, ref))

		reward := &model.ReferralReward{
			ReferralID:  ref.ID,
			UserID:      referrer.ID,
			Type:        model.RewardTypeCash,
			AmountCents: 500,
			Status:      model.ReferralRewardStatusIssued,
		}
		require.NoError(t, repo.CreateReward(ctx, reward))
	}

	rewards, err := repo.GetUserRewards(ctx, referrer.ID, 10)
	require.NoError(t, err)
	assert.Len(t, rewards, 2)
}

func TestReferralRepository_SumUserRewards(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := referral.NewReferralRepository(db)
	ctx := context.Background()

	referrer := CreateUniqueTestUser(t, db, "sum_rewards_referrer")

	for i := 0; i < 3; i++ {
		referee := CreateUniqueTestUser(t, db, "sum_rewards_referee")
		ref := &model.Referral{
			ReferrerID: referrer.ID,
			RefereeID:  referee.ID,
			Type:       model.ReferralTypeUserToUser,
			Status:     model.ReferralStatusCompleted,
		}
		require.NoError(t, repo.CreateReferral(ctx, ref))

		reward := &model.ReferralReward{
			ReferralID:  ref.ID,
			UserID:      referrer.ID,
			Type:        model.RewardTypeCash,
			AmountCents: 1000,
			Status:      model.ReferralRewardStatusIssued,
		}
		require.NoError(t, repo.CreateReward(ctx, reward))
		// Update status to issued with issued_at
		now := time.Now()
		db.Model(reward).Updates(map[string]any{
			"status":    model.ReferralRewardStatusIssued,
			"issued_at": &now,
		})
	}

	sum, err := repo.SumUserRewards(ctx, referrer.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(3000), sum)
}

// ============================================================================
// Referral Stats Tests
// ============================================================================

func TestReferralRepository_GetReferralStats(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := referral.NewReferralRepository(db)
	ctx := context.Background()

	// Create some data
	referrer := CreateUniqueTestUser(t, db, "stats_referrer")
	referee := CreateUniqueTestUser(t, db, "stats_referee")

	code := &model.ReferralCode{
		UserID:   referrer.ID,
		Type:     model.ReferralTypeUserToUser,
		IsActive: true,
	}
	require.NoError(t, repo.CreateCode(ctx, code))

	ref := &model.Referral{
		ReferrerID: referrer.ID,
		RefereeID:  referee.ID,
		Type:       model.ReferralTypeUserToUser,
		Status:     model.ReferralStatusCompleted,
	}
	require.NoError(t, repo.CreateReferral(ctx, ref))

	stats, err := repo.GetReferralStats(ctx)
	require.NoError(t, err)
	assert.NotNil(t, stats["totalCount"])
	assert.NotNil(t, stats["completedCount"])
	assert.NotNil(t, stats["activeCodeCount"])
}

func TestReferralRepository_GetUserReferralStats(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := referral.NewReferralRepository(db)
	ctx := context.Background()

	referrer := CreateUniqueTestUser(t, db, "user_stats_referrer")

	for i := 0; i < 2; i++ {
		referee := CreateUniqueTestUser(t, db, "user_stats_referee")
		ref := &model.Referral{
			ReferrerID: referrer.ID,
			RefereeID:  referee.ID,
			Type:       model.ReferralTypeUserToUser,
			Status:     model.ReferralStatusCompleted,
		}
		require.NoError(t, repo.CreateReferral(ctx, ref))
	}

	stats, err := repo.GetUserReferralStats(ctx, referrer.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(2), stats["totalCount"])
	assert.Equal(t, int64(2), stats["completedCount"])
}

// ============================================================================
// Code Validation Tests
// ============================================================================

func TestReferralCode_IsValid(t *testing.T) {
	// Test active code
	code := &model.ReferralCode{
		IsActive: true,
		MaxUse:   10,
		UseCount: 5,
	}
	assert.True(t, code.IsValid())

	// Test inactive code
	code.IsActive = false
	assert.False(t, code.IsValid())

	// Test expired code
	code.IsActive = true
	expired := time.Now().Add(-time.Hour)
	code.ExpireAt = &expired
	assert.False(t, code.IsValid())

	// Test max use reached
	code.ExpireAt = nil
	code.UseCount = 10
	assert.False(t, code.IsValid())
}
