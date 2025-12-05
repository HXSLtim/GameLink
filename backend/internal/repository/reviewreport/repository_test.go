package reviewreport

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/pkg/testutil"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.User{},
		&model.Player{},
		&model.Order{},
		&model.Review{},
		&model.ReviewReport{},
	)
	return db
}

func TestReviewReportRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer testutil.CleanDB(t, db)

	repo := NewReviewReportRepository(db)
	ctx := context.Background()

	report := &model.ReviewReport{
		ReviewID:   1001,
		ReporterID: 2001,
		Reason:     "评价内容包含不实信息",
		Evidence:   "https://example.com/evidence.jpg",
		Status:     model.ReviewReportStatusPending,
	}

	err := repo.Create(ctx, report)
	require.NoError(t, err)
	assert.NotZero(t, report.ID)
	assert.NotZero(t, report.CreatedAt)
}

func TestReviewReportRepository_Get(t *testing.T) {
	db := setupTestDB(t)
	defer testutil.CleanDB(t, db)

	repo := NewReviewReportRepository(db)
	ctx := context.Background()

	// Create a report
	report := &model.ReviewReport{
		ReviewID:   1001,
		ReporterID: 2001,
		Reason:     "评价内容包含不实信息",
		Status:     model.ReviewReportStatusPending,
	}
	err := repo.Create(ctx, report)
	require.NoError(t, err)

	// Get the report
	retrieved, err := repo.Get(ctx, report.ID)
	require.NoError(t, err)
	assert.Equal(t, report.ID, retrieved.ID)
	assert.Equal(t, report.ReviewID, retrieved.ReviewID)
	assert.Equal(t, report.ReporterID, retrieved.ReporterID)
	assert.Equal(t, report.Reason, retrieved.Reason)
	assert.Equal(t, report.Status, retrieved.Status)

	// Test not found
	_, err = repo.Get(ctx, 99999)
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestReviewReportRepository_List(t *testing.T) {
	db := setupTestDB(t)
	defer testutil.CleanDB(t, db)

	repo := NewReviewReportRepository(db)
	ctx := context.Background()

	// Create multiple reports
	reports := []*model.ReviewReport{
		{
			ReviewID:   1001,
			ReporterID: 2001,
			Reason:     "不实信息",
			Status:     model.ReviewReportStatusPending,
		},
		{
			ReviewID:   1002,
			ReporterID: 2002,
			Reason:     "恶意诋毁",
			Status:     model.ReviewReportStatusApproved,
		},
		{
			ReviewID:   1003,
			ReporterID: 2001,
			Reason:     "广告内容",
			Status:     model.ReviewReportStatusPending,
		},
	}

	for _, report := range reports {
		err := repo.Create(ctx, report)
		require.NoError(t, err)
	}

	t.Run("list all reports", func(t *testing.T) {
		results, total, err := repo.List(ctx, repository.ReviewReportListOptions{
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, results, 3)
	})

	t.Run("filter by status", func(t *testing.T) {
		status := model.ReviewReportStatusPending
		results, total, err := repo.List(ctx, repository.ReviewReportListOptions{
			Page:     1,
			PageSize: 10,
			Status:   &status,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, results, 2)
		for _, r := range results {
			assert.Equal(t, model.ReviewReportStatusPending, r.Status)
		}
	})

	t.Run("filter by reporter", func(t *testing.T) {
		reporterID := uint64(2001)
		results, total, err := repo.List(ctx, repository.ReviewReportListOptions{
			Page:       1,
			PageSize:   10,
			ReporterID: &reporterID,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, results, 2)
		for _, r := range results {
			assert.Equal(t, uint64(2001), r.ReporterID)
		}
	})

	t.Run("filter by review", func(t *testing.T) {
		reviewID := uint64(1001)
		results, total, err := repo.List(ctx, repository.ReviewReportListOptions{
			Page:     1,
			PageSize: 10,
			ReviewID: &reviewID,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, results, 1)
		assert.Equal(t, uint64(1001), results[0].ReviewID)
	})

	t.Run("pagination", func(t *testing.T) {
		results, total, err := repo.List(ctx, repository.ReviewReportListOptions{
			Page:     1,
			PageSize: 2,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, results, 2)

		results, total, err = repo.List(ctx, repository.ReviewReportListOptions{
			Page:     2,
			PageSize: 2,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, results, 1)
	})
}

func TestReviewReportRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer testutil.CleanDB(t, db)

	repo := NewReviewReportRepository(db)
	ctx := context.Background()

	// Create a report
	report := &model.ReviewReport{
		ReviewID:   1001,
		ReporterID: 2001,
		Reason:     "评价内容包含不实信息",
		Status:     model.ReviewReportStatusPending,
	}
	err := repo.Create(ctx, report)
	require.NoError(t, err)

	// Update the report
	handlerID := uint64(3001)
	handledAt := time.Now()
	report.Status = model.ReviewReportStatusApproved
	report.HandledBy = &handlerID
	report.HandledAt = &handledAt
	report.HandlingNote = "经核实，评价内容确实存在不实信息"

	err = repo.Update(ctx, report)
	require.NoError(t, err)

	// Verify the update
	retrieved, err := repo.Get(ctx, report.ID)
	require.NoError(t, err)
	assert.Equal(t, model.ReviewReportStatusApproved, retrieved.Status)
	assert.NotNil(t, retrieved.HandledBy)
	assert.Equal(t, uint64(3001), *retrieved.HandledBy)
	assert.NotNil(t, retrieved.HandledAt)
	assert.Equal(t, "经核实，评价内容确实存在不实信息", retrieved.HandlingNote)

	// Test update non-existent report
	nonExistent := &model.ReviewReport{
		Base: model.Base{ID: 99999},
	}
	err = repo.Update(ctx, nonExistent)
	assert.ErrorIs(t, err, repository.ErrNotFound)
}
