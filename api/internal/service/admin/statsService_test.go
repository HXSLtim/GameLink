package admin

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/repository"
	repomock "gamelink/internal/repository/mocks"
)

func TestStatsService_UserBehaviorStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repomock.NewMockStatsRepository(ctrl)
	mockRepo.EXPECT().UserBehaviorStats(gomock.Any()).Return(repository.UserBehaviorMetrics{
		DAU:                     128,
		AvgOnlineDurationSecond: 5400,
		AvgConsumptionCents:     12345,
	}, nil)

	svc := NewStatsService(mockRepo)
	result, err := svc.UserBehaviorStats(context.Background())
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, int64(128), result.DAU)
	assert.Equal(t, "1h30m", result.AvgOnlineTime)
	assert.Equal(t, 123.45, result.AvgConsumption)
}

func TestStatsService_UserActivityTrend_DefaultDays(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := []repository.DateValue{
		{Date: "2026-03-01", Value: 12},
		{Date: "2026-03-02", Value: 18},
	}

	mockRepo := repomock.NewMockStatsRepository(ctrl)
	mockRepo.EXPECT().UserActivityTrend(gomock.Any(), 7).Return(expected, nil)

	svc := NewStatsService(mockRepo)
	result, err := svc.UserActivityTrend(context.Background(), 0)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestStatsService_UserDistribution(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repomock.NewMockStatsRepository(ctrl)
	mockRepo.EXPECT().UserDistribution(gomock.Any()).Return(repository.UserDistributionMetrics{
		ByRegion: []repository.DistributionValue{
			{Name: "上海", Value: 23},
		},
		ByAge: []repository.DistributionValue{
			{Name: "0-30天", Value: 11},
			{Name: "31-90天", Value: 7},
		},
	}, nil)

	svc := NewStatsService(mockRepo)
	result, err := svc.UserDistribution(context.Background())
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Len(t, result.ByRegion, 1)
	assert.Equal(t, "上海", result.ByRegion[0].Name)
	assert.Equal(t, int64(23), result.ByRegion[0].Value)
	assert.Len(t, result.ByAge, 2)
	assert.Equal(t, "31-90天", result.ByAge[1].Name)
	assert.Equal(t, int64(7), result.ByAge[1].Value)
}

func TestFormatOnlineDuration(t *testing.T) {
	tests := []struct {
		name     string
		seconds  float64
		expected string
	}{
		{name: "zero", seconds: 0, expected: "0m"},
		{name: "under one minute rounds up", seconds: 45, expected: "1m"},
		{name: "minutes", seconds: 600, expected: "10m"},
		{name: "hour", seconds: 3600, expected: "1h"},
		{name: "hour and minute", seconds: 3660, expected: "1h1m"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, formatOnlineDuration(tc.seconds))
		})
	}
}
