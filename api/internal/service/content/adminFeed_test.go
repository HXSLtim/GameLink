package content

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/service/sensitiveword"
)

// MockOperationLogRepository is a mock implementation of OperationLogRepository
type MockOperationLogRepository struct {
	mock.Mock
}

func (m *MockOperationLogRepository) Append(ctx context.Context, log *model.OperationLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

func (m *MockOperationLogRepository) ListByEntity(ctx context.Context, entityType string, entityID uint64, opts repository.OperationLogListOptions) ([]model.OperationLog, int64, error) {
	args := m.Called(ctx, entityType, entityID, opts)
	return args.Get(0).([]model.OperationLog), args.Get(1).(int64), args.Error(2)
}

func (m *MockOperationLogRepository) List(ctx context.Context, opts repository.OperationLogSearchOptions) ([]model.OperationLog, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.OperationLog), args.Get(1).(int64), args.Error(2)
}

// MockSensitiveWordService is a mock for sensitive word detection
type MockSensitiveWordService struct {
	mock.Mock
}

func (m *MockSensitiveWordService) DetectSensitiveWords(ctx context.Context, req sensitiveword.DetectSensitiveWordsRequest) (*sensitiveword.DetectSensitiveWordsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sensitiveword.DetectSensitiveWordsResponse), args.Error(1)
}

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
	return args.Get(0).(map[model.FeedModerationStatus]int64), args.Error(1)
}

func (m *MockFeedRepository) GetTrend(ctx context.Context, days int) ([]repository.DateValue, error) {
	args := m.Called(ctx, days)
	return args.Get(0).([]repository.DateValue), args.Error(1)
}

func TestNewAdminFeedService(t *testing.T) {
	feedRepo := &MockFeedRepository{}
	opLogRepo := &MockOperationLogRepository{}

	t.Run("with all dependencies", func(t *testing.T) {
		svc := NewAdminFeedService(feedRepo, nil, opLogRepo)
		require.NotNil(t, svc)
		assert.Equal(t, feedRepo, svc.feedRepo)
		assert.Nil(t, svc.sensitiveWord)
		assert.Equal(t, opLogRepo, svc.opLogRepo)
	})

	t.Run("with nil opLogRepo", func(t *testing.T) {
		svc := NewAdminFeedService(feedRepo, nil, nil)
		require.NotNil(t, svc)
		assert.Nil(t, svc.opLogRepo)
	})
}

func TestAdminFeedService_ListFeeds(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success with default pagination", func(t *testing.T) {
		feedRepo := &MockFeedRepository{}
		svc := NewAdminFeedService(feedRepo, nil, nil)

		feeds := []model.Feed{
			{Base: model.Base{ID: 1, CreatedAt: now, UpdatedAt: now}, AuthorID: 1, Content: "Feed 1", ModerationStatus: model.FeedModerationPending},
			{Base: model.Base{ID: 2, CreatedAt: now, UpdatedAt: now}, AuthorID: 2, Content: "Feed 2", ModerationStatus: model.FeedModerationApproved},
		}
		feedRepo.On("ListPaged", ctx, mock.AnythingOfType("repository.FeedPagedListOptions")).Return(feeds, int64(2), nil)

		req := AdminListFeedsRequest{Page: 0, PageSize: 0}
		result, err := svc.ListFeeds(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Items, 2)
		assert.Equal(t, int64(2), result.Total)
		feedRepo.AssertExpectations(t)
	})

	t.Run("success with filters", func(t *testing.T) {
		feedRepo := &MockFeedRepository{}
		svc := NewAdminFeedService(feedRepo, nil, nil)

		authorID := uint64(1)
		categoryID := uint64(5)
		status := model.FeedModerationPending
		feeds := []model.Feed{
			{Base: model.Base{ID: 1, CreatedAt: now, UpdatedAt: now}, AuthorID: 1, Content: "Filtered Feed"},
		}
		feedRepo.On("ListPaged", ctx, mock.AnythingOfType("repository.FeedPagedListOptions")).Return(feeds, int64(1), nil)

		req := AdminListFeedsRequest{
			Page:             1,
			PageSize:         10,
			AuthorID:         &authorID,
			CategoryID:       &categoryID,
			Keyword:          "test",
			ModerationStatus: &status,
		}
		result, err := svc.ListFeeds(ctx, req)

		require.NoError(t, err)
		assert.Len(t, result.Items, 1)
	})

	t.Run("repo error", func(t *testing.T) {
		feedRepo := &MockFeedRepository{}
		svc := NewAdminFeedService(feedRepo, nil, nil)

		feedRepo.On("ListPaged", ctx, mock.AnythingOfType("repository.FeedPagedListOptions")).Return([]model.Feed{}, int64(0), errors.New("db error"))

		req := AdminListFeedsRequest{Page: 1, PageSize: 10}
		_, err := svc.ListFeeds(ctx, req)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})

	t.Run("page size capped at 100", func(t *testing.T) {
		feedRepo := &MockFeedRepository{}
		svc := NewAdminFeedService(feedRepo, nil, nil)

		feedRepo.On("ListPaged", ctx, mock.MatchedBy(func(opts repository.FeedPagedListOptions) bool {
			return opts.PageSize == 20 // Should be capped to default 20 when > 100
		})).Return([]model.Feed{}, int64(0), nil)

		req := AdminListFeedsRequest{Page: 1, PageSize: 200}
		_, err := svc.ListFeeds(ctx, req)

		require.NoError(t, err)
	})
}

func TestAdminFeedService_GetFeed(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		feedRepo := &MockFeedRepository{}
		svc := NewAdminFeedService(feedRepo, nil, nil)

		feed := &model.Feed{
			Base:             model.Base{ID: 1, CreatedAt: now, UpdatedAt: now},
			AuthorID:         100,
			Content:          "Test content",
			ModerationStatus: model.FeedModerationApproved,
			Images: []model.FeedImage{
				{Base: model.Base{ID: 1}, URL: "https://example.com/1.jpg", Order: 0, Width: 800, Height: 600},
			},
		}
		feedRepo.On("Get", ctx, uint64(1)).Return(feed, nil)

		result, err := svc.GetFeed(ctx, 1)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, uint64(1), result.ID)
		assert.Equal(t, "Test content", result.Content)
		assert.Len(t, result.Images, 1)
		feedRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		feedRepo := &MockFeedRepository{}
		svc := NewAdminFeedService(feedRepo, nil, nil)

		feedRepo.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

		_, err := svc.GetFeed(ctx, 999)

		require.Error(t, err)
		assert.Equal(t, repository.ErrNotFound, err)
	})
}

func TestAdminFeedService_ApproveFeed(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		feedRepo := &MockFeedRepository{}
		opLogRepo := &MockOperationLogRepository{}
		svc := NewAdminFeedService(feedRepo, nil, opLogRepo)

		moderatorID := uint64(1)
		feedRepo.On("UpdateModeration", ctx, uint64(1), model.FeedModerationApproved, "approved", &moderatorID).Return(nil)
		opLogRepo.On("Append", ctx, mock.AnythingOfType("*model.OperationLog")).Return(nil)

		err := svc.ApproveFeed(ctx, 1, moderatorID, "approved")

		require.NoError(t, err)
		feedRepo.AssertExpectations(t)
		opLogRepo.AssertExpectations(t)
	})

	t.Run("success without note", func(t *testing.T) {
		feedRepo := &MockFeedRepository{}
		opLogRepo := &MockOperationLogRepository{}
		svc := NewAdminFeedService(feedRepo, nil, opLogRepo)

		moderatorID := uint64(1)
		feedRepo.On("UpdateModeration", ctx, uint64(1), model.FeedModerationApproved, "", &moderatorID).Return(nil)
		opLogRepo.On("Append", ctx, mock.AnythingOfType("*model.OperationLog")).Return(nil)

		err := svc.ApproveFeed(ctx, 1, moderatorID, "")

		require.NoError(t, err)
	})

	t.Run("repo error", func(t *testing.T) {
		feedRepo := &MockFeedRepository{}
		svc := NewAdminFeedService(feedRepo, nil, nil)

		moderatorID := uint64(1)
		feedRepo.On("UpdateModeration", ctx, uint64(1), model.FeedModerationApproved, "note", &moderatorID).Return(errors.New("db error"))

		err := svc.ApproveFeed(ctx, 1, moderatorID, "note")

		require.Error(t, err)
	})

	t.Run("without opLogRepo", func(t *testing.T) {
		feedRepo := &MockFeedRepository{}
		svc := NewAdminFeedService(feedRepo, nil, nil)

		moderatorID := uint64(1)
		feedRepo.On("UpdateModeration", ctx, uint64(1), model.FeedModerationApproved, "note", &moderatorID).Return(nil)

		err := svc.ApproveFeed(ctx, 1, moderatorID, "note")

		require.NoError(t, err)
	})
}

func TestAdminFeedService_RejectFeed(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		feedRepo := &MockFeedRepository{}
		opLogRepo := &MockOperationLogRepository{}
		svc := NewAdminFeedService(feedRepo, nil, opLogRepo)

		moderatorID := uint64(1)
		feedRepo.On("UpdateModeration", ctx, uint64(1), model.FeedModerationRejected, "违规内容", &moderatorID).Return(nil)
		opLogRepo.On("Append", ctx, mock.AnythingOfType("*model.OperationLog")).Return(nil)

		err := svc.RejectFeed(ctx, 1, moderatorID, "违规内容")

		require.NoError(t, err)
		feedRepo.AssertExpectations(t)
	})

	t.Run("empty note returns validation error", func(t *testing.T) {
		feedRepo := &MockFeedRepository{}
		svc := NewAdminFeedService(feedRepo, nil, nil)

		err := svc.RejectFeed(ctx, 1, 1, "")

		require.Error(t, err)
		assert.Equal(t, ErrAdminValidation, err)
	})

	t.Run("repo error", func(t *testing.T) {
		feedRepo := &MockFeedRepository{}
		svc := NewAdminFeedService(feedRepo, nil, nil)

		moderatorID := uint64(1)
		feedRepo.On("UpdateModeration", ctx, uint64(1), model.FeedModerationRejected, "reason", &moderatorID).Return(errors.New("db error"))

		err := svc.RejectFeed(ctx, 1, moderatorID, "reason")

		require.Error(t, err)
	})
}

func TestAdminFeedService_DeleteFeed(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		feedRepo := &MockFeedRepository{}
		opLogRepo := &MockOperationLogRepository{}
		svc := NewAdminFeedService(feedRepo, nil, opLogRepo)

		moderatorID := uint64(1)
		feedRepo.On("UpdateModeration", ctx, uint64(1), model.FeedModerationRemoved, "spam", &moderatorID).Return(nil)
		opLogRepo.On("Append", ctx, mock.AnythingOfType("*model.OperationLog")).Return(nil)

		err := svc.DeleteFeed(ctx, 1, moderatorID, "spam")

		require.NoError(t, err)
		feedRepo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		feedRepo := &MockFeedRepository{}
		svc := NewAdminFeedService(feedRepo, nil, nil)

		moderatorID := uint64(1)
		feedRepo.On("UpdateModeration", ctx, uint64(1), model.FeedModerationRemoved, "reason", &moderatorID).Return(errors.New("db error"))

		err := svc.DeleteFeed(ctx, 1, moderatorID, "reason")

		require.Error(t, err)
	})
}

func TestAdminFeedService_BatchApproveFeed(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		feedRepo := &MockFeedRepository{}
		opLogRepo := &MockOperationLogRepository{}
		svc := NewAdminFeedService(feedRepo, nil, opLogRepo)

		feedIDs := []uint64{1, 2, 3}
		moderatorID := uint64(1)
		feedRepo.On("BatchUpdateModeration", ctx, feedIDs, model.FeedModerationApproved, "batch approved", &moderatorID).Return(nil)
		opLogRepo.On("Append", ctx, mock.AnythingOfType("*model.OperationLog")).Return(nil).Times(3)

		err := svc.BatchApproveFeed(ctx, feedIDs, moderatorID, "batch approved")

		require.NoError(t, err)
		feedRepo.AssertExpectations(t)
		opLogRepo.AssertExpectations(t)
	})

	t.Run("empty feedIDs returns validation error", func(t *testing.T) {
		feedRepo := &MockFeedRepository{}
		svc := NewAdminFeedService(feedRepo, nil, nil)

		err := svc.BatchApproveFeed(ctx, []uint64{}, 1, "note")

		require.Error(t, err)
		assert.Equal(t, ErrAdminValidation, err)
	})

	t.Run("repo error", func(t *testing.T) {
		feedRepo := &MockFeedRepository{}
		svc := NewAdminFeedService(feedRepo, nil, nil)

		feedIDs := []uint64{1, 2}
		moderatorID := uint64(1)
		feedRepo.On("BatchUpdateModeration", ctx, feedIDs, model.FeedModerationApproved, "note", &moderatorID).Return(errors.New("db error"))

		err := svc.BatchApproveFeed(ctx, feedIDs, moderatorID, "note")

		require.Error(t, err)
	})
}

func TestAdminFeedService_BatchRejectFeed(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		feedRepo := &MockFeedRepository{}
		opLogRepo := &MockOperationLogRepository{}
		svc := NewAdminFeedService(feedRepo, nil, opLogRepo)

		feedIDs := []uint64{1, 2}
		moderatorID := uint64(1)
		feedRepo.On("BatchUpdateModeration", ctx, feedIDs, model.FeedModerationRejected, "违规内容", &moderatorID).Return(nil)
		opLogRepo.On("Append", ctx, mock.AnythingOfType("*model.OperationLog")).Return(nil).Times(2)

		err := svc.BatchRejectFeed(ctx, feedIDs, moderatorID, "违规内容")

		require.NoError(t, err)
		feedRepo.AssertExpectations(t)
	})

	t.Run("empty feedIDs returns validation error", func(t *testing.T) {
		feedRepo := &MockFeedRepository{}
		svc := NewAdminFeedService(feedRepo, nil, nil)

		err := svc.BatchRejectFeed(ctx, []uint64{}, 1, "reason")

		require.Error(t, err)
		assert.Equal(t, ErrAdminValidation, err)
	})

	t.Run("empty note returns validation error", func(t *testing.T) {
		feedRepo := &MockFeedRepository{}
		svc := NewAdminFeedService(feedRepo, nil, nil)

		err := svc.BatchRejectFeed(ctx, []uint64{1, 2}, 1, "")

		require.Error(t, err)
		assert.Equal(t, ErrAdminValidation, err)
	})

	t.Run("repo error", func(t *testing.T) {
		feedRepo := &MockFeedRepository{}
		svc := NewAdminFeedService(feedRepo, nil, nil)

		feedIDs := []uint64{1, 2}
		moderatorID := uint64(1)
		feedRepo.On("BatchUpdateModeration", ctx, feedIDs, model.FeedModerationRejected, "reason", &moderatorID).Return(errors.New("db error"))

		err := svc.BatchRejectFeed(ctx, feedIDs, moderatorID, "reason")

		require.Error(t, err)
	})
}

func TestAdminFeedService_toDTO(t *testing.T) {
	now := time.Now()

	t.Run("basic feed", func(t *testing.T) {
		svc := &AdminFeedService{}
		feed := &model.Feed{
			Base:             model.Base{ID: 1, CreatedAt: now, UpdatedAt: now},
			AuthorID:         100,
			Content:          "Test content",
			ModerationStatus: model.FeedModerationApproved,
			ModerationNote:   "auto approved",
			Visibility:       model.FeedVisibilityPublic,
		}

		dto := svc.toDTO(feed)

		assert.Equal(t, uint64(1), dto.ID)
		assert.Equal(t, uint64(100), dto.AuthorID)
		assert.Equal(t, "Test content", dto.Content)
		assert.Equal(t, model.FeedModerationApproved, dto.ModerationStatus)
		assert.Equal(t, "auto approved", dto.ModerationNote)
		assert.Empty(t, dto.Images)
	})

	t.Run("feed with category", func(t *testing.T) {
		svc := &AdminFeedService{}
		categoryID := uint64(5)
		feed := &model.Feed{
			Base:       model.Base{ID: 1, CreatedAt: now, UpdatedAt: now},
			CategoryID: &categoryID,
			Category:   &model.ContentCategory{Name: "游戏分享"},
		}

		dto := svc.toDTO(feed)

		assert.Equal(t, &categoryID, dto.CategoryID)
		assert.Equal(t, "游戏分享", dto.CategoryName)
	})

	t.Run("feed with images", func(t *testing.T) {
		svc := &AdminFeedService{}
		feed := &model.Feed{
			Base: model.Base{ID: 1, CreatedAt: now, UpdatedAt: now},
			Images: []model.FeedImage{
				{Base: model.Base{ID: 1}, URL: "https://example.com/1.jpg", Order: 0, Width: 800, Height: 600},
				{Base: model.Base{ID: 2}, URL: "https://example.com/2.jpg", Order: 1, Width: 1024, Height: 768},
			},
		}

		dto := svc.toDTO(feed)

		assert.Len(t, dto.Images, 2)
		assert.Equal(t, uint64(1), dto.Images[0].ID)
		assert.Equal(t, "https://example.com/1.jpg", dto.Images[0].URL)
		assert.Equal(t, 800, dto.Images[0].Width)
	})
}
