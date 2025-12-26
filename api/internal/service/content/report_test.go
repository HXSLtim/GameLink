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
)

func TestNewFeedReportService(t *testing.T) {
	feedRepo := &MockFeedRepository{}
	opLogRepo := &MockOperationLogRepository{}

	svc := NewFeedReportService(feedRepo, opLogRepo)

	require.NotNil(t, svc)
	assert.Equal(t, feedRepo, svc.feedRepo)
	assert.Equal(t, opLogRepo, svc.opLogRepo)
}

func TestFeedReportService_ListFeedReports(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success with default pagination", func(t *testing.T) {
		feedRepo := &MockFeedRepository{}
		svc := NewFeedReportService(feedRepo, nil)

		reports := []model.FeedReport{
			{Base: model.Base{ID: 1, CreatedAt: now}, FeedID: 100, Reporter: 1, Reason: "spam", Status: "pending"},
			{Base: model.Base{ID: 2, CreatedAt: now}, FeedID: 101, Reporter: 2, Reason: "abuse", Status: "pending"},
		}
		feed := &model.Feed{Base: model.Base{ID: 100, CreatedAt: now}, Content: "Test", ModerationStatus: model.FeedModerationApproved}

		feedRepo.On("ListReports", ctx, mock.AnythingOfType("repository.FeedReportListOptions")).Return(reports, int64(2), nil)
		feedRepo.On("Get", ctx, uint64(100)).Return(feed, nil)
		feedRepo.On("Get", ctx, uint64(101)).Return(feed, nil)

		req := ListFeedReportsRequest{Page: 0, PageSize: 0}
		result, err := svc.ListFeedReports(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Items, 2)
		assert.Equal(t, int64(2), result.Total)
	})

	t.Run("repo error", func(t *testing.T) {
		feedRepo := &MockFeedRepository{}
		svc := NewFeedReportService(feedRepo, nil)

		feedRepo.On("ListReports", ctx, mock.AnythingOfType("repository.FeedReportListOptions")).Return([]model.FeedReport{}, int64(0), errors.New("db error"))

		req := ListFeedReportsRequest{Page: 1, PageSize: 10}
		_, err := svc.ListFeedReports(ctx, req)

		require.Error(t, err)
	})
}

func TestFeedReportService_GetFeedReport(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		feedRepo := &MockFeedRepository{}
		svc := NewFeedReportService(feedRepo, nil)

		report := &model.FeedReport{
			Base:     model.Base{ID: 1, CreatedAt: now},
			FeedID:   100,
			Reporter: 1,
			Reason:   "spam",
			Status:   "pending",
		}
		feedRepo.On("GetReport", ctx, uint64(1)).Return(report, nil)

		result, err := svc.GetFeedReport(ctx, 1)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, uint64(1), result.ID)
		assert.Equal(t, "spam", result.Reason)
	})

	t.Run("not found", func(t *testing.T) {
		feedRepo := &MockFeedRepository{}
		svc := NewFeedReportService(feedRepo, nil)

		feedRepo.On("GetReport", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

		_, err := svc.GetFeedReport(ctx, 999)

		require.Error(t, err)
	})
}

func TestFeedReportService_ProcessFeedReport(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("delete content action", func(t *testing.T) {
		feedRepo := &MockFeedRepository{}
		opLogRepo := &MockOperationLogRepository{}
		svc := NewFeedReportService(feedRepo, opLogRepo)

		report := &model.FeedReport{
			Base:     model.Base{ID: 1, CreatedAt: now},
			FeedID:   100,
			Reporter: 1,
			Reason:   "spam",
			Status:   "pending",
		}
		handlerID := uint64(10)

		feedRepo.On("GetReport", ctx, uint64(1)).Return(report, nil)
		feedRepo.On("UpdateModeration", ctx, uint64(100), model.FeedModerationRemoved, "举报处理：删除内容", &handlerID).Return(nil)
		feedRepo.On("UpdateReport", ctx, mock.AnythingOfType("*model.FeedReport")).Return(nil)
		opLogRepo.On("Append", ctx, mock.AnythingOfType("*model.OperationLog")).Return(nil)

		req := ProcessReportRequest{Action: ReportActionDeleteContent, Result: "违规内容已删除"}
		err := svc.ProcessFeedReport(ctx, 1, req, handlerID)

		require.NoError(t, err)
		feedRepo.AssertExpectations(t)
	})

	t.Run("warn user action", func(t *testing.T) {
		feedRepo := &MockFeedRepository{}
		opLogRepo := &MockOperationLogRepository{}
		svc := NewFeedReportService(feedRepo, opLogRepo)

		report := &model.FeedReport{
			Base:     model.Base{ID: 1, CreatedAt: now},
			FeedID:   100,
			Reporter: 1,
			Status:   "pending",
		}
		handlerID := uint64(10)

		feedRepo.On("GetReport", ctx, uint64(1)).Return(report, nil)
		feedRepo.On("UpdateReport", ctx, mock.AnythingOfType("*model.FeedReport")).Return(nil)
		opLogRepo.On("Append", ctx, mock.AnythingOfType("*model.OperationLog")).Return(nil)

		req := ProcessReportRequest{Action: ReportActionWarnUser, Result: "已警告用户"}
		err := svc.ProcessFeedReport(ctx, 1, req, handlerID)

		require.NoError(t, err)
	})

	t.Run("dismiss action", func(t *testing.T) {
		feedRepo := &MockFeedRepository{}
		opLogRepo := &MockOperationLogRepository{}
		svc := NewFeedReportService(feedRepo, opLogRepo)

		report := &model.FeedReport{
			Base:     model.Base{ID: 1, CreatedAt: now},
			FeedID:   100,
			Reporter: 1,
			Status:   "pending",
		}
		handlerID := uint64(10)

		feedRepo.On("GetReport", ctx, uint64(1)).Return(report, nil)
		feedRepo.On("UpdateReport", ctx, mock.AnythingOfType("*model.FeedReport")).Return(nil)
		opLogRepo.On("Append", ctx, mock.AnythingOfType("*model.OperationLog")).Return(nil)

		req := ProcessReportRequest{Action: ReportActionDismiss, Result: "举报不成立"}
		err := svc.ProcessFeedReport(ctx, 1, req, handlerID)

		require.NoError(t, err)
	})

	t.Run("invalid action", func(t *testing.T) {
		feedRepo := &MockFeedRepository{}
		svc := NewFeedReportService(feedRepo, nil)

		report := &model.FeedReport{
			Base:   model.Base{ID: 1, CreatedAt: now},
			FeedID: 100,
			Status: "pending",
		}

		feedRepo.On("GetReport", ctx, uint64(1)).Return(report, nil)

		req := ProcessReportRequest{Action: "invalid_action"}
		err := svc.ProcessFeedReport(ctx, 1, req, 10)

		require.Error(t, err)
		assert.Equal(t, ErrAdminValidation, err)
	})

	t.Run("report not found", func(t *testing.T) {
		feedRepo := &MockFeedRepository{}
		svc := NewFeedReportService(feedRepo, nil)

		feedRepo.On("GetReport", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

		req := ProcessReportRequest{Action: ReportActionDismiss}
		err := svc.ProcessFeedReport(ctx, 999, req, 10)

		require.Error(t, err)
	})
}

func TestFeedReportService_toReportDTO(t *testing.T) {
	now := time.Now()
	handledAt := now.Add(1 * time.Hour)
	handlerID := uint64(10)

	t.Run("basic report", func(t *testing.T) {
		svc := &FeedReportService{}
		report := &model.FeedReport{
			Base:     model.Base{ID: 1, CreatedAt: now},
			FeedID:   100,
			Reporter: 1,
			Reason:   "spam",
			Status:   "pending",
		}

		dto := svc.toReportDTO(report)

		assert.Equal(t, uint64(1), dto.ID)
		assert.Equal(t, uint64(100), dto.FeedID)
		assert.Equal(t, "spam", dto.Reason)
		assert.Equal(t, "pending", dto.Status)
		assert.Empty(t, dto.HandledAt)
	})

	t.Run("handled report", func(t *testing.T) {
		svc := &FeedReportService{}
		report := &model.FeedReport{
			Base:      model.Base{ID: 1, CreatedAt: now},
			FeedID:    100,
			Reporter:  1,
			Reason:    "spam",
			Status:    "processed",
			Result:    "已删除",
			HandledBy: &handlerID,
			HandledAt: &handledAt,
		}

		dto := svc.toReportDTO(report)

		assert.Equal(t, "processed", dto.Status)
		assert.Equal(t, "已删除", dto.Result)
		assert.Equal(t, &handlerID, dto.HandledBy)
		assert.NotEmpty(t, dto.HandledAt)
	})
}
