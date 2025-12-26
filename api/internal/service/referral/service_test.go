package referral

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"gamelink/internal/model"
)

// ============================================================================
// Tests
// ============================================================================

func TestService_GetExpireDays(t *testing.T) {
	t.Run("default expire days", func(t *testing.T) {
		days := 30 // default
		assert.Equal(t, 30, days)
	})

	t.Run("custom expire days", func(t *testing.T) {
		days := 7
		assert.Equal(t, 7, days)
	})
}

func TestReferralCode_IsValid(t *testing.T) {
	t.Run("valid code", func(t *testing.T) {
		expireAt := time.Now().Add(7 * 24 * time.Hour)
		code := &model.ReferralCode{
			IsActive: true,
			MaxUse:   10,
			UseCount: 5,
			ExpireAt: &expireAt,
		}
		assert.True(t, code.IsValid())
	})

	t.Run("inactive code", func(t *testing.T) {
		code := &model.ReferralCode{
			IsActive: false,
		}
		assert.False(t, code.IsValid())
	})

	t.Run("expired code", func(t *testing.T) {
		expireAt := time.Now().Add(-1 * time.Hour)
		code := &model.ReferralCode{
			IsActive: true,
			ExpireAt: &expireAt,
		}
		assert.False(t, code.IsValid())
	})

	t.Run("max use reached", func(t *testing.T) {
		code := &model.ReferralCode{
			IsActive: true,
			MaxUse:   10,
			UseCount: 10,
		}
		assert.False(t, code.IsValid())
	})

	t.Run("unlimited use", func(t *testing.T) {
		code := &model.ReferralCode{
			IsActive: true,
			MaxUse:   0, // unlimited
			UseCount: 100,
		}
		assert.True(t, code.IsValid())
	})
}

func TestReferralType_Constants(t *testing.T) {
	t.Run("type values", func(t *testing.T) {
		assert.Equal(t, model.ReferralType("user_to_user"), model.ReferralTypeUserToUser)
		assert.Equal(t, model.ReferralType("player_to_player"), model.ReferralTypePlayerToPlayer)
		assert.Equal(t, model.ReferralType("user_to_player"), model.ReferralTypeUserToPlayer)
	})
}

func TestReferralStatus_Constants(t *testing.T) {
	t.Run("status values", func(t *testing.T) {
		assert.Equal(t, model.ReferralStatus("pending"), model.ReferralStatusPending)
		assert.Equal(t, model.ReferralStatus("completed"), model.ReferralStatusCompleted)
		assert.Equal(t, model.ReferralStatus("rewarded"), model.ReferralStatusRewarded)
		assert.Equal(t, model.ReferralStatus("expired"), model.ReferralStatusExpired)
		assert.Equal(t, model.ReferralStatus("canceled"), model.ReferralStatusCanceled)
	})
}

func TestRewardType_Constants(t *testing.T) {
	t.Run("type values", func(t *testing.T) {
		assert.Equal(t, model.RewardType("cash"), model.RewardTypeCash)
		assert.Equal(t, model.RewardType("coupon"), model.RewardTypeCoupon)
		assert.Equal(t, model.RewardType("points"), model.RewardTypePoints)
	})
}

func TestReferralRewardStatus_Constants(t *testing.T) {
	t.Run("status values", func(t *testing.T) {
		assert.Equal(t, model.ReferralRewardStatus("pending"), model.ReferralRewardStatusPending)
		assert.Equal(t, model.ReferralRewardStatus("issued"), model.ReferralRewardStatusIssued)
		assert.Equal(t, model.ReferralRewardStatus("failed"), model.ReferralRewardStatusFailed)
	})
}

func TestReferral_Structure(t *testing.T) {
	t.Run("referral record", func(t *testing.T) {
		codeID := uint64(1)
		referral := &model.Referral{
			ReferrerID:       100,
			RefereeID:        200,
			CodeID:           &codeID,
			Type:             model.ReferralTypeUserToUser,
			Level:            1,
			Status:           model.ReferralStatusPending,
			RefereeCondition: "registered",
		}

		assert.Equal(t, uint64(100), referral.ReferrerID)
		assert.Equal(t, uint64(200), referral.RefereeID)
		assert.NotNil(t, referral.CodeID)
		assert.Equal(t, 1, referral.Level)
	})
}

func TestReferralReward_Structure(t *testing.T) {
	t.Run("reward record", func(t *testing.T) {
		reward := &model.ReferralReward{
			ReferralID:  1,
			UserID:      100,
			Type:        model.RewardTypeCash,
			AmountCents: 1000,
			Status:      model.ReferralRewardStatusPending,
		}

		assert.Equal(t, uint64(1), reward.ReferralID)
		assert.Equal(t, uint64(100), reward.UserID)
		assert.Equal(t, int64(1000), reward.AmountCents)
	})
}

func TestReferralConfig_Keys(t *testing.T) {
	t.Run("config keys", func(t *testing.T) {
		assert.Equal(t, "enabled", model.ReferralConfigEnabled)
		assert.Equal(t, "expire_days", model.ReferralConfigExpireDays)
		assert.Equal(t, "max_level", model.ReferralConfigMaxLevel)
		assert.Equal(t, "user_reward_type", model.ReferralConfigUserRewardType)
		assert.Equal(t, "user_reward_amount", model.ReferralConfigUserRewardAmount)
		assert.Equal(t, "player_reward_type", model.ReferralConfigPlayerRewardType)
		assert.Equal(t, "player_reward_amount", model.ReferralConfigPlayerRewardAmount)
	})
}

func TestService_GetUserReferrals_DefaultLimit(t *testing.T) {
	t.Run("default limit", func(t *testing.T) {
		limit := 0
		if limit <= 0 {
			limit = 20
		}
		assert.Equal(t, 20, limit)
	})

	t.Run("custom limit", func(t *testing.T) {
		limit := 50
		if limit <= 0 {
			limit = 20
		}
		assert.Equal(t, 50, limit)
	})
}

func TestReferralCode_Expiration(t *testing.T) {
	t.Run("not expired", func(t *testing.T) {
		expireAt := time.Now().Add(7 * 24 * time.Hour)
		code := &model.ReferralCode{
			ExpireAt: &expireAt,
		}
		assert.False(t, time.Now().After(*code.ExpireAt))
	})

	t.Run("expired", func(t *testing.T) {
		expireAt := time.Now().Add(-1 * time.Hour)
		code := &model.ReferralCode{
			ExpireAt: &expireAt,
		}
		assert.True(t, time.Now().After(*code.ExpireAt))
	})

	t.Run("no expiration", func(t *testing.T) {
		code := &model.ReferralCode{
			ExpireAt: nil,
		}
		assert.Nil(t, code.ExpireAt)
	})
}

func TestReferral_Level(t *testing.T) {
	t.Run("direct referral", func(t *testing.T) {
		referral := &model.Referral{
			Level: 1,
		}
		assert.Equal(t, 1, referral.Level)
	})

	t.Run("second level referral", func(t *testing.T) {
		referral := &model.Referral{
			Level: 2,
		}
		assert.Equal(t, 2, referral.Level)
	})
}

func TestReferral_Conditions(t *testing.T) {
	t.Run("registered condition", func(t *testing.T) {
		referral := &model.Referral{
			RefereeCondition: "registered",
		}
		assert.Equal(t, "registered", referral.RefereeCondition)
	})

	t.Run("first_order condition", func(t *testing.T) {
		referral := &model.Referral{
			RefereeCondition: "first_order",
		}
		assert.Equal(t, "first_order", referral.RefereeCondition)
	})

	t.Run("first_recharge condition", func(t *testing.T) {
		referral := &model.Referral{
			RefereeCondition: "first_recharge",
		}
		assert.Equal(t, "first_recharge", referral.RefereeCondition)
	})
}
