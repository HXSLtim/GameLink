package content

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// MockFeedRepository is a mock implementation of FeedRepository
type MockFeedRepository struct {
	mock.Mock
}

func (m *MockFeedRepository) Create(ctx context.Context, feed *model.Feed) error {
	args := m.Called(ctx, feed)
	return args.Error(0)
}

func (m *MockFeedRepository) Get(ctx context.Context, id uint64) (*model.Feed, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Feed), args.Error(1)
}

func (m *MockFeedRepository) List(ctx context.Context, opts repository.FeedListOptions) ([]model.Feed, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Feed), args.Error(1)
}

func (m *MockFeedRepository) ListPaged(ctx context.Context, opts repository.FeedPagedListOptions) ([]model.Feed, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.Feed), args.Get(1).(int64), args.Error(2)
}

func (m *MockFeedRepository) Update(ctx context.Context, feed *model.Feed) error {
	args := m.Called(ctx, feed)
	return args.Error(0)
}

func (m *MockFeedRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockFeedRepository) UpdateModeration(ctx context.Context, feedID uint64, status model.FeedModerationStatus, note string, moderatorID *uint64) error {
	args := m.Called(ctx, feedID, status, note, moderatorID)
	return args.Error(0)
}

func (m *MockFeedRepository) BatchUpdateModeration(ctx context.Context, feedIDs []uint64, status model.FeedModerationStatus, note string, moderatorID *uint64) error {
	args := m.Called(ctx, feedIDs, status, note, moderatorID)
	return args.Error(0)
}

func (m *MockFeedRepository) CreateReport(ctx context.Context, report *model.FeedReport) error {
	args := m.Called(ctx, report)
	return args.Error(0)
}

func (m *MockFeedRepository) GetReport(ctx context.Context, id uint64) (*model.FeedReport, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.FeedReport), args.Error(1)
}

func (m *MockFeedRepository) ListReports(ctx context.Context, opts repository.FeedReportListOptions) ([]model.FeedReport, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.FeedReport), args.Get(1).(int64), args.Error(2)
}

func (m *MockFeedRepository) UpdateReport(ctx context.Context, report *model.FeedReport) error {
	args := m.Called(ctx, report)
	return args.Error(0)
}

func (m *MockFeedRepository) CountByStatus(ctx context.Context) (map[model.FeedModerationStatus]int64, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[model.FeedModerationStatus]int64), args.Error(1)
}

func (m *MockFeedRepository) GetTrend(ctx context.Context, days int) ([]repository.DateValue, error) {
	args := m.Called(ctx, days)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repository.DateValue), args.Error(1)
}

// MockChatMessageRepository is a mock implementation of ChatMessageRepository
type MockChatMessageRepository struct {
	mock.Mock
}

func (m *MockChatMessageRepository) Create(ctx context.Context, msg *model.ChatMessage) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}

func (m *MockChatMessageRepository) CreateBatch(ctx context.Context, messages []*model.ChatMessage) error {
	args := m.Called(ctx, messages)
	return args.Error(0)
}

func (m *MockChatMessageRepository) Get(ctx context.Context, id uint64) (*model.ChatMessage, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ChatMessage), args.Error(1)
}

func (m *MockChatMessageRepository) ListByGroup(ctx context.Context, opts repository.ChatMessageListOptions) ([]model.ChatMessage, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.ChatMessage), args.Get(1).(int64), args.Error(2)
}

func (m *MockChatMessageRepository) MarkDeleted(ctx context.Context, id uint64, deletedBy uint64) error {
	args := m.Called(ctx, id, deletedBy)
	return args.Error(0)
}

func (m *MockChatMessageRepository) ListForModeration(ctx context.Context, opts repository.ChatMessageModerationListOptions) ([]model.ChatMessage, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.ChatMessage), args.Get(1).(int64), args.Error(2)
}

func (m *MockChatMessageRepository) UpdateAuditStatus(ctx context.Context, id uint64, status model.ChatMessageAuditStatus, moderatorID *uint64, reason string) error {
	args := m.Called(ctx, id, status, moderatorID, reason)
	return args.Error(0)
}

func (m *MockChatMessageRepository) DeleteByGroupIDs(ctx context.Context, groupIDs []uint64) error {
	args := m.Called(ctx, groupIDs)
	return args.Error(0)
}

func TestNewContentStatsService(t *testing.T) {
	feedRepo := &MockFeedRepository{}
	messageRepo := &MockChatMessageRepository{}

	svc := NewContentStatsService(feedRepo, messageRepo)
	require.NotNil(t, svc)
	assert.Equal(t, feedRepo, svc.feedRepo)
	assert.Equal(t, messageRepo, svc.messageRepo)
}

func TestContentStatsService_GetStats(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		trendDays     int
		setupMocks    func(*MockFeedRepository, *MockChatMessageRepository)
		wantErr       bool
		validateStats func(*testing.T, *ContentStatsDTO)
	}{
		{
			name:      "successful stats retrieval",
			trendDays: 7,
			setupMocks: func(feedRepo *MockFeedRepository, msgRepo *MockChatMessageRepository) {
				feedRepo.On("CountByStatus", ctx).Return(map[model.FeedModerationStatus]int64{
					model.FeedModerationPending:  10,
					model.FeedModerationApproved: 100,
					model.FeedModerationRejected: 5,
				}, nil)
				feedRepo.On("GetTrend", ctx, 7).Return([]repository.DateValue{
					{Date: "2025-12-01", Value: 15},
					{Date: "2025-12-02", Value: 20},
				}, nil)
				msgRepo.On("ListForModeration", ctx, mock.Anything).Return([]model.ChatMessage{}, int64(500), nil)
			},
			validateStats: func(t *testing.T, stats *ContentStatsDTO) {
				assert.Equal(t, int64(115), stats.TotalFeeds)
				assert.Equal(t, int64(10), stats.PendingFeeds)
				assert.Equal(t, int64(100), stats.ApprovedFeeds)
				assert.Equal(t, int64(5), stats.RejectedFeeds)
				assert.Equal(t, int64(500), stats.TotalMessages)
				assert.Len(t, stats.FeedTrend, 2)
			},
		},
		{
			name:      "default trend days when less than 1",
			trendDays: 0,
			setupMocks: func(feedRepo *MockFeedRepository, msgRepo *MockChatMessageRepository) {
				feedRepo.On("CountByStatus", ctx).Return(map[model.FeedModerationStatus]int64{
					model.FeedModerationApproved: 50,
				}, nil)
				feedRepo.On("GetTrend", ctx, 30).Return([]repository.DateValue{}, nil) // Should use 30 as default
				msgRepo.On("ListForModeration", ctx, mock.Anything).Return([]model.ChatMessage{}, int64(0), nil)
			},
			validateStats: func(t *testing.T, stats *ContentStatsDTO) {
				assert.Equal(t, int64(50), stats.TotalFeeds)
			},
		},
		{
			name:      "count by status error",
			trendDays: 7,
			setupMocks: func(feedRepo *MockFeedRepository, msgRepo *MockChatMessageRepository) {
				feedRepo.On("CountByStatus", ctx).Return(nil, errors.New("database error"))
			},
			wantErr: true,
		},
		{
			name:      "trend error should not fail - returns empty trend",
			trendDays: 7,
			setupMocks: func(feedRepo *MockFeedRepository, msgRepo *MockChatMessageRepository) {
				feedRepo.On("CountByStatus", ctx).Return(map[model.FeedModerationStatus]int64{
					model.FeedModerationApproved: 20,
				}, nil)
				feedRepo.On("GetTrend", ctx, 7).Return(nil, errors.New("trend error"))
				msgRepo.On("ListForModeration", ctx, mock.Anything).Return([]model.ChatMessage{}, int64(100), nil)
			},
			validateStats: func(t *testing.T, stats *ContentStatsDTO) {
				assert.Equal(t, int64(20), stats.TotalFeeds)
				assert.Empty(t, stats.FeedTrend)
			},
		},
		{
			name:      "empty stats",
			trendDays: 7,
			setupMocks: func(feedRepo *MockFeedRepository, msgRepo *MockChatMessageRepository) {
				feedRepo.On("CountByStatus", ctx).Return(map[model.FeedModerationStatus]int64{}, nil)
				feedRepo.On("GetTrend", ctx, 7).Return([]repository.DateValue{}, nil)
				msgRepo.On("ListForModeration", ctx, mock.Anything).Return([]model.ChatMessage{}, int64(0), nil)
			},
			validateStats: func(t *testing.T, stats *ContentStatsDTO) {
				assert.Equal(t, int64(0), stats.TotalFeeds)
				assert.Equal(t, int64(0), stats.PendingFeeds)
				assert.Equal(t, int64(0), stats.ApprovedFeeds)
				assert.Equal(t, int64(0), stats.RejectedFeeds)
				assert.Equal(t, int64(0), stats.TotalMessages)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			feedRepo := &MockFeedRepository{}
			msgRepo := &MockChatMessageRepository{}
			tt.setupMocks(feedRepo, msgRepo)

			svc := NewContentStatsService(feedRepo, msgRepo)
			stats, err := svc.GetStats(ctx, tt.trendDays)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, stats)
			if tt.validateStats != nil {
				tt.validateStats(t, stats)
			}

			feedRepo.AssertExpectations(t)
			msgRepo.AssertExpectations(t)
		})
	}
}

func TestContentStatsService_ExportStats(t *testing.T) {
	ctx := context.Background()

	t.Run("successful export", func(t *testing.T) {
		feedRepo := &MockFeedRepository{}
		msgRepo := &MockChatMessageRepository{}

		feedRepo.On("CountByStatus", ctx).Return(map[model.FeedModerationStatus]int64{
			model.FeedModerationPending:  5,
			model.FeedModerationApproved: 50,
			model.FeedModerationRejected: 2,
		}, nil)
		feedRepo.On("GetTrend", ctx, 7).Return([]repository.DateValue{
			{Date: "2025-12-01", Value: 10},
			{Date: "2025-12-02", Value: 15},
		}, nil)
		msgRepo.On("ListForModeration", ctx, mock.Anything).Return([]model.ChatMessage{}, int64(200), nil)

		svc := NewContentStatsService(feedRepo, msgRepo)
		buf, filename, err := svc.ExportStats(ctx, 7)

		require.NoError(t, err)
		require.NotNil(t, buf)
		assert.True(t, buf.Len() > 0, "buffer should not be empty")
		assert.Contains(t, filename, "content_stats_")
		assert.Contains(t, filename, ".xlsx")

		feedRepo.AssertExpectations(t)
		msgRepo.AssertExpectations(t)
	})

	t.Run("export with stats error", func(t *testing.T) {
		feedRepo := &MockFeedRepository{}
		msgRepo := &MockChatMessageRepository{}

		feedRepo.On("CountByStatus", ctx).Return(nil, errors.New("database error"))

		svc := NewContentStatsService(feedRepo, msgRepo)
		buf, filename, err := svc.ExportStats(ctx, 7)

		require.Error(t, err)
		assert.Nil(t, buf)
		assert.Empty(t, filename)
	})
}

func TestContentStatsDTO_Fields(t *testing.T) {
	dto := ContentStatsDTO{
		TotalFeeds:    100,
		PendingFeeds:  10,
		ApprovedFeeds: 80,
		RejectedFeeds: 10,
		TotalMessages: 500,
		FeedsByStatus: map[model.FeedModerationStatus]int64{
			model.FeedModerationPending:  10,
			model.FeedModerationApproved: 80,
			model.FeedModerationRejected: 10,
		},
		FeedTrend: []repository.DateValue{
			{Date: "2025-12-01", Value: 20},
		},
	}

	assert.Equal(t, int64(100), dto.TotalFeeds)
	assert.Equal(t, int64(10), dto.PendingFeeds)
	assert.Equal(t, int64(80), dto.ApprovedFeeds)
	assert.Equal(t, int64(10), dto.RejectedFeeds)
	assert.Equal(t, int64(500), dto.TotalMessages)
	assert.Len(t, dto.FeedsByStatus, 3)
	assert.Len(t, dto.FeedTrend, 1)
}
