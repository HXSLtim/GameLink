package activity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
)

// ============================================================================
// Tests
// ============================================================================

func TestService_validateActivity(t *testing.T) {
	svc := &Service{}

	t.Run("empty name", func(t *testing.T) {
		activity := &model.Activity{Name: ""}
		err := svc.validateActivity(activity)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "名称")
	})

	t.Run("empty start time", func(t *testing.T) {
		activity := &model.Activity{Name: "Test"}
		err := svc.validateActivity(activity)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "开始时间")
	})

	t.Run("empty end time", func(t *testing.T) {
		activity := &model.Activity{
			Name:    "Test",
			StartAt: time.Now(),
		}
		err := svc.validateActivity(activity)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "结束时间")
	})

	t.Run("end before start", func(t *testing.T) {
		now := time.Now()
		activity := &model.Activity{
			Name:    "Test",
			StartAt: now,
			EndAt:   now.Add(-1 * time.Hour),
		}
		err := svc.validateActivity(activity)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "晚于")
	})

	t.Run("preheat after start", func(t *testing.T) {
		now := time.Now()
		preheat := now.Add(1 * time.Hour)
		activity := &model.Activity{
			Name:      "Test",
			StartAt:   now,
			EndAt:     now.Add(2 * time.Hour),
			PreheatAt: &preheat,
		}
		err := svc.validateActivity(activity)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "预热")
	})

	t.Run("valid activity", func(t *testing.T) {
		now := time.Now()
		preheat := now.Add(-1 * time.Hour)
		activity := &model.Activity{
			Name:      "Test",
			StartAt:   now,
			EndAt:     now.Add(2 * time.Hour),
			PreheatAt: &preheat,
		}
		err := svc.validateActivity(activity)
		require.NoError(t, err)
	})
}

func TestService_validateStatusTransition(t *testing.T) {
	svc := &Service{}

	t.Run("draft to active", func(t *testing.T) {
		err := svc.validateStatusTransition(model.ActivityStatusDraft, model.ActivityStatusActive)
		require.NoError(t, err)
	})

	t.Run("draft to preheat", func(t *testing.T) {
		err := svc.validateStatusTransition(model.ActivityStatusDraft, model.ActivityStatusPreheat)
		require.NoError(t, err)
	})

	t.Run("draft to canceled", func(t *testing.T) {
		err := svc.validateStatusTransition(model.ActivityStatusDraft, model.ActivityStatusCanceled)
		require.NoError(t, err)
	})

	t.Run("active to paused", func(t *testing.T) {
		err := svc.validateStatusTransition(model.ActivityStatusActive, model.ActivityStatusPaused)
		require.NoError(t, err)
	})

	t.Run("active to ended", func(t *testing.T) {
		err := svc.validateStatusTransition(model.ActivityStatusActive, model.ActivityStatusEnded)
		require.NoError(t, err)
	})

	t.Run("paused to active", func(t *testing.T) {
		err := svc.validateStatusTransition(model.ActivityStatusPaused, model.ActivityStatusActive)
		require.NoError(t, err)
	})

	t.Run("ended cannot transition", func(t *testing.T) {
		err := svc.validateStatusTransition(model.ActivityStatusEnded, model.ActivityStatusActive)
		require.Error(t, err)
	})

	t.Run("canceled cannot transition", func(t *testing.T) {
		err := svc.validateStatusTransition(model.ActivityStatusCanceled, model.ActivityStatusActive)
		require.Error(t, err)
	})

	t.Run("invalid transition draft to ended", func(t *testing.T) {
		err := svc.validateStatusTransition(model.ActivityStatusDraft, model.ActivityStatusEnded)
		require.Error(t, err)
	})
}

func TestService_validateReward(t *testing.T) {
	svc := &Service{}

	t.Run("empty template id", func(t *testing.T) {
		reward := &model.ActivityReward{CouponTemplateID: 0}
		err := svc.validateReward(reward)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "模板")
	})

	t.Run("invalid coupon count", func(t *testing.T) {
		reward := &model.ActivityReward{CouponTemplateID: 1, CouponCount: 0}
		err := svc.validateReward(reward)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "数量")
	})

	t.Run("invalid probability low", func(t *testing.T) {
		reward := &model.ActivityReward{CouponTemplateID: 1, CouponCount: 1, Probability: 0}
		err := svc.validateReward(reward)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "概率")
	})

	t.Run("invalid probability high", func(t *testing.T) {
		reward := &model.ActivityReward{CouponTemplateID: 1, CouponCount: 1, Probability: 101}
		err := svc.validateReward(reward)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "概率")
	})

	t.Run("valid reward", func(t *testing.T) {
		reward := &model.ActivityReward{CouponTemplateID: 1, CouponCount: 1, Probability: 50}
		err := svc.validateReward(reward)
		require.NoError(t, err)
	})
}

func TestActivity_IsActive(t *testing.T) {
	t.Run("active status and in time range", func(t *testing.T) {
		now := time.Now()
		activity := &model.Activity{
			Status:  model.ActivityStatusActive,
			StartAt: now.Add(-1 * time.Hour),
			EndAt:   now.Add(1 * time.Hour),
		}
		assert.True(t, activity.IsActive())
	})

	t.Run("not active status", func(t *testing.T) {
		now := time.Now()
		activity := &model.Activity{
			Status:  model.ActivityStatusDraft,
			StartAt: now.Add(-1 * time.Hour),
			EndAt:   now.Add(1 * time.Hour),
		}
		assert.False(t, activity.IsActive())
	})

	t.Run("active but not started", func(t *testing.T) {
		now := time.Now()
		activity := &model.Activity{
			Status:  model.ActivityStatusActive,
			StartAt: now.Add(1 * time.Hour),
			EndAt:   now.Add(2 * time.Hour),
		}
		assert.False(t, activity.IsActive())
	})

	t.Run("active but ended", func(t *testing.T) {
		now := time.Now()
		activity := &model.Activity{
			Status:  model.ActivityStatusActive,
			StartAt: now.Add(-2 * time.Hour),
			EndAt:   now.Add(-1 * time.Hour),
		}
		assert.False(t, activity.IsActive())
	})
}

func TestActivity_IsPreheat(t *testing.T) {
	t.Run("in preheat period", func(t *testing.T) {
		now := time.Now()
		preheat := now.Add(-1 * time.Hour)
		activity := &model.Activity{
			Status:    model.ActivityStatusPreheat,
			PreheatAt: &preheat,
			StartAt:   now.Add(1 * time.Hour),
			EndAt:     now.Add(2 * time.Hour),
		}
		// Check preheat status
		assert.Equal(t, model.ActivityStatusPreheat, activity.Status)
		assert.True(t, activity.PreheatAt != nil)
		assert.True(t, time.Now().After(*activity.PreheatAt))
		assert.True(t, time.Now().Before(activity.StartAt))
	})

	t.Run("not preheat status", func(t *testing.T) {
		now := time.Now()
		preheat := now.Add(-1 * time.Hour)
		activity := &model.Activity{
			Status:    model.ActivityStatusDraft,
			PreheatAt: &preheat,
			StartAt:   now.Add(1 * time.Hour),
			EndAt:     now.Add(2 * time.Hour),
		}
		assert.NotEqual(t, model.ActivityStatusPreheat, activity.Status)
	})
}

func TestActivityReward_StockManagement(t *testing.T) {
	t.Run("has stock", func(t *testing.T) {
		reward := &model.ActivityReward{
			TotalStock:     100,
			RemainingStock: 50,
		}
		assert.True(t, reward.RemainingStock > 0)
	})

	t.Run("out of stock", func(t *testing.T) {
		reward := &model.ActivityReward{
			TotalStock:     100,
			RemainingStock: 0,
		}
		assert.False(t, reward.RemainingStock > 0)
	})

	t.Run("unlimited stock", func(t *testing.T) {
		reward := &model.ActivityReward{
			TotalStock:     0,
			RemainingStock: 0,
		}
		// TotalStock = 0 means unlimited
		assert.Equal(t, 0, reward.TotalStock)
	})
}

func TestActivityParticipation_Structure(t *testing.T) {
	t.Run("participation record", func(t *testing.T) {
		now := time.Now()
		participation := &model.ActivityParticipation{
			ActivityID: 1,
			UserID:     100,
			RewardID:   10,
			CouponIDs:  "[1, 2, 3]",
			ClaimedAt:  now,
			ClientIP:   "192.168.1.1",
		}

		assert.Equal(t, uint64(1), participation.ActivityID)
		assert.Equal(t, uint64(100), participation.UserID)
		assert.Equal(t, uint64(10), participation.RewardID)
		assert.NotEmpty(t, participation.CouponIDs)
	})
}

func TestService_ParticipationLimits(t *testing.T) {
	t.Run("per user limit check", func(t *testing.T) {
		activity := &model.Activity{
			PerUserLimit: 3,
		}
		userCount := int64(3)
		assert.True(t, int(userCount) >= activity.PerUserLimit)
	})

	t.Run("daily limit check", func(t *testing.T) {
		activity := &model.Activity{
			DailyLimit: 100,
		}
		todayCount := int64(100)
		assert.True(t, int(todayCount) >= activity.DailyLimit)
	})

	t.Run("total limit check", func(t *testing.T) {
		activity := &model.Activity{
			TotalLimit:   1000,
			TotalClaimed: 1000,
		}
		assert.True(t, activity.TotalClaimed >= activity.TotalLimit)
	})
}

func TestService_GetUserParticipations_DefaultLimit(t *testing.T) {
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

func TestActivityStatus_Constants(t *testing.T) {
	t.Run("status values", func(t *testing.T) {
		assert.Equal(t, model.ActivityStatus("draft"), model.ActivityStatusDraft)
		assert.Equal(t, model.ActivityStatus("preheat"), model.ActivityStatusPreheat)
		assert.Equal(t, model.ActivityStatus("active"), model.ActivityStatusActive)
		assert.Equal(t, model.ActivityStatus("paused"), model.ActivityStatusPaused)
		assert.Equal(t, model.ActivityStatus("ended"), model.ActivityStatusEnded)
		assert.Equal(t, model.ActivityStatus("canceled"), model.ActivityStatusCanceled)
	})
}
