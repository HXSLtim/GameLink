package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestActivity_TableName(t *testing.T) {
	activity := Activity{}
	assert.Equal(t, "activities", activity.TableName())
}

func TestActivity_IsInPreheat(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		activity Activity
		expected bool
	}{
		{
			name: "no preheat time set",
			activity: Activity{
				PreheatAt: nil,
				StartAt:   now.Add(time.Hour),
			},
			expected: false,
		},
		{
			name: "in preheat period",
			activity: Activity{
				PreheatAt: func() *time.Time { t := now.Add(-time.Hour); return &t }(),
				StartAt:   now.Add(time.Hour),
			},
			expected: true,
		},
		{
			name: "before preheat",
			activity: Activity{
				PreheatAt: func() *time.Time { t := now.Add(time.Hour); return &t }(),
				StartAt:   now.Add(2 * time.Hour),
			},
			expected: false,
		},
		{
			name: "after start (not in preheat)",
			activity: Activity{
				PreheatAt: func() *time.Time { t := now.Add(-2 * time.Hour); return &t }(),
				StartAt:   now.Add(-time.Hour),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.activity.IsInPreheat())
		})
	}
}

func TestActivity_IsActive(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		activity Activity
		expected bool
	}{
		{
			name: "active and in time range",
			activity: Activity{
				Status:  ActivityStatusActive,
				StartAt: now.Add(-time.Hour),
				EndAt:   now.Add(time.Hour),
			},
			expected: true,
		},
		{
			name: "active but before start",
			activity: Activity{
				Status:  ActivityStatusActive,
				StartAt: now.Add(time.Hour),
				EndAt:   now.Add(2 * time.Hour),
			},
			expected: false,
		},
		{
			name: "active but after end",
			activity: Activity{
				Status:  ActivityStatusActive,
				StartAt: now.Add(-2 * time.Hour),
				EndAt:   now.Add(-time.Hour),
			},
			expected: false,
		},
		{
			name: "not active status",
			activity: Activity{
				Status:  ActivityStatusPaused,
				StartAt: now.Add(-time.Hour),
				EndAt:   now.Add(time.Hour),
			},
			expected: false,
		},
		{
			name: "draft status",
			activity: Activity{
				Status:  ActivityStatusDraft,
				StartAt: now.Add(-time.Hour),
				EndAt:   now.Add(time.Hour),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.activity.IsActive())
		})
	}
}

func TestActivityReward_TableName(t *testing.T) {
	reward := ActivityReward{}
	assert.Equal(t, "activity_rewards", reward.TableName())
}

func TestActivityParticipation_TableName(t *testing.T) {
	participation := ActivityParticipation{}
	assert.Equal(t, "activity_participations", participation.TableName())
}

func TestActivityDailyStats_TableName(t *testing.T) {
	stats := ActivityDailyStats{}
	assert.Equal(t, "activity_daily_stats", stats.TableName())
}

func TestActivityStatusConstants(t *testing.T) {
	assert.Equal(t, ActivityStatus("draft"), ActivityStatusDraft)
	assert.Equal(t, ActivityStatus("preheat"), ActivityStatusPreheat)
	assert.Equal(t, ActivityStatus("active"), ActivityStatusActive)
	assert.Equal(t, ActivityStatus("paused"), ActivityStatusPaused)
	assert.Equal(t, ActivityStatus("ended"), ActivityStatusEnded)
	assert.Equal(t, ActivityStatus("canceled"), ActivityStatusCanceled)
}

func TestActivityTypeConstants(t *testing.T) {
	assert.Equal(t, ActivityType("coupon"), ActivityTypeCoupon)
	assert.Equal(t, ActivityType("discount"), ActivityTypeDiscount)
	assert.Equal(t, ActivityType("gift"), ActivityTypeGift)
}
