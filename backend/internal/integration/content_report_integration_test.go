package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	adminhandler "gamelink/internal/handler/admin"
	"gamelink/internal/model"
	feedrepo "gamelink/internal/repository/content"
	operationlogrepo "gamelink/internal/repository/operationlog"
	userrepo "gamelink/internal/repository/user"
	contentservice "gamelink/internal/service/content"
	"gamelink/pkg/testutil"
)

// TestReportProcessingFlow tests the complete report processing workflow
// Validates: Requirements 5.1, 5.2, 5.3, 5.4, 5.5
func TestReportProcessingFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateContentReportModels(t, db)

	ctx := context.Background()
	userRepo := userrepo.NewUserRepository(db)
	feedRepo := feedrepo.NewFeedRepository(db)
	opLogRepo := operationlogrepo.NewOperationLogRepository(db)

	// Create admin user
	admin := &model.User{
		Name:         "ReportAdmin",
		Email:        "report_admin@example.com",
		Phone:        "90000000001",
		PasswordHash: "hashed",
		Status:       model.UserStatusActive,
		Role:         model.RoleAdmin,
	}
	require.NoError(t, userRepo.Create(ctx, admin))

	// Create content author
	author := &model.User{
		Name:         "ContentAuthor",
		Email:        "content_author@example.com",
		Phone:        "90000000002",
		PasswordHash: "hashed",
		Status:       model.UserStatusActive,
		Role:         model.RoleUser,
	}
	require.NoError(t, userRepo.Create(ctx, author))

	// Create reporter
	reporter := &model.User{
		Name:         "Reporter",
		Email:        "reporter@example.com",
		Phone:        "90000000003",
		PasswordHash: "hashed",
		Status:       model.UserStatusActive,
		Role:         model.RoleUser,
	}
	require.NoError(t, userRepo.Create(ctx, reporter))

	// Create feeds to be reported
	feed1 := &model.Feed{
		AuthorID:         author.ID,
		Content:          "Reported content 1",
		Visibility:       model.FeedVisibilityPublic,
		ModerationStatus: model.FeedModerationApproved,
	}
	require.NoError(t, feedRepo.Create(ctx, feed1))

	feed2 := &model.Feed{
		AuthorID:         author.ID,
		Content:          "Reported content 2",
		Visibility:       model.FeedVisibilityPublic,
		ModerationStatus: model.FeedModerationApproved,
	}
	require.NoError(t, feedRepo.Create(ctx, feed2))

	feed3 := &model.Feed{
		AuthorID:         author.ID,
		Content:          "Reported content 3",
		Visibility:       model.FeedVisibilityPublic,
		ModerationStatus: model.FeedModerationApproved,
	}
	require.NoError(t, feedRepo.Create(ctx, feed3))

	// Create reports
	report1 := &model.FeedReport{
		FeedID:   feed1.ID,
		Reporter: reporter.ID,
		Reason:   "Inappropriate content",
		Status:   "pending",
	}
	require.NoError(t, db.Create(report1).Error)

	report2 := &model.FeedReport{
		FeedID:   feed2.ID,
		Reporter: reporter.ID,
		Reason:   "Spam content",
		Status:   "pending",
	}
	require.NoError(t, db.Create(report2).Error)

	report3 := &model.FeedReport{
		FeedID:   feed3.ID,
		Reporter: reporter.ID,
		Reason:   "False report test",
		Status:   "pending",
	}
	require.NoError(t, db.Create(report3).Error)

	// Setup service and handler
	reportSvc := contentservice.NewFeedReportService(feedRepo, opLogRepo)
	contentHandler := adminhandler.NewContentHandler(nil, nil, reportSvc, nil)

	router := gin.New()
	api := router.Group("/api/v1/admin")
	auth := fakeContentReportAdminAuthMiddleware(admin.ID)
	api.Use(auth)
	api.GET("/content/reports", contentHandler.ListFeedReports)
	api.GET("/content/reports/:id", contentHandler.GetFeedReport)
	api.POST("/content/reports/:id/process", contentHandler.ProcessFeedReport)

	t.Run("List pending reports", func(t *testing.T) {
		resp := doJSON(router, http.MethodGet, "/api/v1/admin/content/reports?status=pending", nil, "")
		assert.Equal(t, http.StatusOK, resp.Code)

		var result apiResp[contentservice.ListFeedReportsResponse]
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		assert.True(t, result.Success)
		assert.Equal(t, int64(3), result.Data.Total)
	})

	t.Run("Get report detail", func(t *testing.T) {
		resp := doJSON(router, http.MethodGet, "/api/v1/admin/content/reports/"+uintToStr(report1.ID), nil, "")
		assert.Equal(t, http.StatusOK, resp.Code)

		var result apiResp[contentservice.FeedReportDTO]
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		assert.True(t, result.Success)
		assert.Equal(t, report1.ID, result.Data.ID)
		assert.Equal(t, "Inappropriate content", result.Data.Reason)
	})

	t.Run("Process report - delete content", func(t *testing.T) {
		payload := map[string]interface{}{
			"action": "delete_content",
			"result": "Content violates community guidelines",
		}
		resp := doJSON(router, http.MethodPost, "/api/v1/admin/content/reports/"+uintToStr(report1.ID)+"/process", payload, "")
		assert.Equal(t, http.StatusOK, resp.Code)

		// Verify report status changed
		var updatedReport model.FeedReport
		require.NoError(t, db.First(&updatedReport, report1.ID).Error)
		assert.Equal(t, "processed", updatedReport.Status)
		assert.NotNil(t, updatedReport.HandledBy)
		assert.Equal(t, admin.ID, *updatedReport.HandledBy)

		// Verify feed was removed
		updatedFeed, err := feedRepo.Get(ctx, feed1.ID)
		require.NoError(t, err)
		assert.Equal(t, model.FeedModerationRemoved, updatedFeed.ModerationStatus)
	})

	t.Run("Process report - warn user", func(t *testing.T) {
		payload := map[string]interface{}{
			"action": "warn_user",
			"result": "User warned for posting spam",
		}
		resp := doJSON(router, http.MethodPost, "/api/v1/admin/content/reports/"+uintToStr(report2.ID)+"/process", payload, "")
		assert.Equal(t, http.StatusOK, resp.Code)

		// Verify report status changed
		var updatedReport model.FeedReport
		require.NoError(t, db.First(&updatedReport, report2.ID).Error)
		assert.Equal(t, "processed", updatedReport.Status)
		assert.Equal(t, "User warned for posting spam", updatedReport.Result)
	})

	t.Run("Process report - dismiss", func(t *testing.T) {
		payload := map[string]interface{}{
			"action": "dismiss",
			"result": "Report is not valid",
		}
		resp := doJSON(router, http.MethodPost, "/api/v1/admin/content/reports/"+uintToStr(report3.ID)+"/process", payload, "")
		assert.Equal(t, http.StatusOK, resp.Code)

		// Verify report status changed to dismissed
		var updatedReport model.FeedReport
		require.NoError(t, db.First(&updatedReport, report3.ID).Error)
		assert.Equal(t, "dismissed", updatedReport.Status)
	})

	t.Run("Process report - invalid action should fail", func(t *testing.T) {
		// Create another report
		report4 := &model.FeedReport{
			FeedID:   feed1.ID,
			Reporter: reporter.ID,
			Reason:   "Another report",
			Status:   "pending",
		}
		require.NoError(t, db.Create(report4).Error)

		payload := map[string]interface{}{
			"action": "invalid_action",
			"result": "Test",
		}
		resp := doJSON(router, http.MethodPost, "/api/v1/admin/content/reports/"+uintToStr(report4.ID)+"/process", payload, "")
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("Get non-existent report should fail", func(t *testing.T) {
		resp := doJSON(router, http.MethodGet, "/api/v1/admin/content/reports/99999", nil, "")
		assert.Equal(t, http.StatusNotFound, resp.Code)
	})
}

// TestReportFilteringFlow tests report filtering functionality
// Validates: Requirements 5.1
func TestReportFilteringFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateContentReportModels(t, db)

	ctx := context.Background()
	userRepo := userrepo.NewUserRepository(db)
	feedRepo := feedrepo.NewFeedRepository(db)
	opLogRepo := operationlogrepo.NewOperationLogRepository(db)

	// Create admin user
	admin := &model.User{
		Name:         "FilterAdmin",
		Email:        "filter_admin@example.com",
		Phone:        "90000000010",
		PasswordHash: "hashed",
		Status:       model.UserStatusActive,
		Role:         model.RoleAdmin,
	}
	require.NoError(t, userRepo.Create(ctx, admin))

	// Create users
	author := &model.User{
		Name:         "FilterAuthor",
		Email:        "filter_author@example.com",
		Phone:        "90000000011",
		PasswordHash: "hashed",
		Status:       model.UserStatusActive,
		Role:         model.RoleUser,
	}
	require.NoError(t, userRepo.Create(ctx, author))

	reporter1 := &model.User{
		Name:         "Reporter1",
		Email:        "reporter1@example.com",
		Phone:        "90000000012",
		PasswordHash: "hashed",
		Status:       model.UserStatusActive,
		Role:         model.RoleUser,
	}
	require.NoError(t, userRepo.Create(ctx, reporter1))

	reporter2 := &model.User{
		Name:         "Reporter2",
		Email:        "reporter2@example.com",
		Phone:        "90000000013",
		PasswordHash: "hashed",
		Status:       model.UserStatusActive,
		Role:         model.RoleUser,
	}
	require.NoError(t, userRepo.Create(ctx, reporter2))

	// Create feeds
	feed1 := &model.Feed{
		AuthorID:         author.ID,
		Content:          "Filter test feed 1",
		Visibility:       model.FeedVisibilityPublic,
		ModerationStatus: model.FeedModerationApproved,
	}
	require.NoError(t, feedRepo.Create(ctx, feed1))

	feed2 := &model.Feed{
		AuthorID:         author.ID,
		Content:          "Filter test feed 2",
		Visibility:       model.FeedVisibilityPublic,
		ModerationStatus: model.FeedModerationApproved,
	}
	require.NoError(t, feedRepo.Create(ctx, feed2))

	// Create reports with different statuses
	require.NoError(t, db.Create(&model.FeedReport{
		FeedID:   feed1.ID,
		Reporter: reporter1.ID,
		Reason:   "Pending report 1",
		Status:   "pending",
	}).Error)

	require.NoError(t, db.Create(&model.FeedReport{
		FeedID:   feed1.ID,
		Reporter: reporter2.ID,
		Reason:   "Pending report 2",
		Status:   "pending",
	}).Error)

	require.NoError(t, db.Create(&model.FeedReport{
		FeedID:   feed2.ID,
		Reporter: reporter1.ID,
		Reason:   "Processed report",
		Status:   "processed",
	}).Error)

	require.NoError(t, db.Create(&model.FeedReport{
		FeedID:   feed2.ID,
		Reporter: reporter2.ID,
		Reason:   "Dismissed report",
		Status:   "dismissed",
	}).Error)

	// Setup service and handler
	reportSvc := contentservice.NewFeedReportService(feedRepo, opLogRepo)
	contentHandler := adminhandler.NewContentHandler(nil, nil, reportSvc, nil)

	router := gin.New()
	api := router.Group("/api/v1/admin")
	auth := fakeContentReportAdminAuthMiddleware(admin.ID)
	api.Use(auth)
	api.GET("/content/reports", contentHandler.ListFeedReports)

	t.Run("Filter by status - pending", func(t *testing.T) {
		resp := doJSON(router, http.MethodGet, "/api/v1/admin/content/reports?status=pending", nil, "")
		assert.Equal(t, http.StatusOK, resp.Code)

		var result apiResp[contentservice.ListFeedReportsResponse]
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		assert.True(t, result.Success)
		assert.Equal(t, int64(2), result.Data.Total)
	})

	t.Run("Filter by status - processed", func(t *testing.T) {
		resp := doJSON(router, http.MethodGet, "/api/v1/admin/content/reports?status=processed", nil, "")
		assert.Equal(t, http.StatusOK, resp.Code)

		var result apiResp[contentservice.ListFeedReportsResponse]
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		assert.True(t, result.Success)
		assert.Equal(t, int64(1), result.Data.Total)
	})

	t.Run("Filter by feedId", func(t *testing.T) {
		resp := doJSON(router, http.MethodGet, "/api/v1/admin/content/reports?feedId="+uintToStr(feed1.ID), nil, "")
		assert.Equal(t, http.StatusOK, resp.Code)

		var result apiResp[contentservice.ListFeedReportsResponse]
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		assert.True(t, result.Success)
		assert.Equal(t, int64(2), result.Data.Total)
	})

	t.Run("List all reports", func(t *testing.T) {
		resp := doJSON(router, http.MethodGet, "/api/v1/admin/content/reports", nil, "")
		assert.Equal(t, http.StatusOK, resp.Code)

		var result apiResp[contentservice.ListFeedReportsResponse]
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		assert.True(t, result.Success)
		assert.Equal(t, int64(4), result.Data.Total)
	})
}

func fakeContentReportAdminAuthMiddleware(adminID uint64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", adminID)
		c.Set("userID", adminID)
		c.Set("admin_id", adminID)
		c.Set("adminID", adminID)
		c.Next()
	}
}

func migrateContentReportModels(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.AutoMigrate(
		&model.User{},
		&model.Feed{},
		&model.FeedImage{},
		&model.FeedReport{},
		&model.OperationLog{},
	); err != nil {
		t.Fatalf("migrate content report models: %v", err)
	}
}
