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

// TestContentModerationFlow tests the complete content moderation workflow
// Validates: Requirements 1.3, 1.4, 2.4, 2.5
func TestContentModerationFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateContentModerationModels(t, db)

	ctx := context.Background()
	userRepo := userrepo.NewUserRepository(db)
	feedRepo := feedrepo.NewFeedRepository(db)
	opLogRepo := operationlogrepo.NewOperationLogRepository(db)

	// Create admin user
	admin := &model.User{
		Name:         "ContentAdmin",
		Email:        "content_admin@example.com",
		Phone:        "80000000001",
		PasswordHash: "hashed",
		Status:       model.UserStatusActive,
		Role:         model.RoleAdmin,
	}
	require.NoError(t, userRepo.Create(ctx, admin))

	// Create regular user
	user := &model.User{
		Name:         "ContentUser",
		Email:        "content_user@example.com",
		Phone:        "80000000002",
		PasswordHash: "hashed",
		Status:       model.UserStatusActive,
		Role:         model.RoleUser,
	}
	require.NoError(t, userRepo.Create(ctx, user))

	// Create pending feeds
	feed1 := &model.Feed{
		AuthorID:         user.ID,
		Content:          "Test feed content 1",
		Visibility:       model.FeedVisibilityPublic,
		ModerationStatus: model.FeedModerationPending,
	}
	require.NoError(t, feedRepo.Create(ctx, feed1))

	feed2 := &model.Feed{
		AuthorID:         user.ID,
		Content:          "Test feed content 2",
		Visibility:       model.FeedVisibilityPublic,
		ModerationStatus: model.FeedModerationPending,
	}
	require.NoError(t, feedRepo.Create(ctx, feed2))

	// Setup service and handler
	adminFeedSvc := contentservice.NewAdminFeedService(feedRepo, nil, opLogRepo)
	contentHandler := adminhandler.NewContentHandler(adminFeedSvc, nil, nil, nil)

	router := gin.New()
	api := router.Group("/api/v1/admin")
	auth := fakeContentAdminAuthMiddleware(admin.ID)
	api.Use(auth)
	api.GET("/content/feeds", contentHandler.ListFeeds)
	api.GET("/content/feeds/:id", contentHandler.GetFeed)
	api.PUT("/content/feeds/:id/approve", contentHandler.ApproveFeed)
	api.PUT("/content/feeds/:id/reject", contentHandler.RejectFeed)
	api.DELETE("/content/feeds/:id", contentHandler.DeleteFeed)
	api.POST("/content/feeds/batch-approve", contentHandler.BatchApproveFeed)
	api.POST("/content/feeds/batch-reject", contentHandler.BatchRejectFeed)

	t.Run("List pending feeds", func(t *testing.T) {
		resp := doJSON(router, http.MethodGet, "/api/v1/admin/content/feeds?moderationStatus=pending", nil, "")
		assert.Equal(t, http.StatusOK, resp.Code)

		var result apiResp[contentservice.AdminListFeedsResponse]
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		assert.True(t, result.Success)
		assert.Equal(t, int64(2), result.Data.Total)
	})

	t.Run("Approve feed", func(t *testing.T) {
		payload := map[string]interface{}{
			"note": "Content approved",
		}
		resp := doJSON(router, http.MethodPut, "/api/v1/admin/content/feeds/"+uintToStr(feed1.ID)+"/approve", payload, "")
		assert.Equal(t, http.StatusOK, resp.Code)

		// Verify feed status changed
		updatedFeed, err := feedRepo.Get(ctx, feed1.ID)
		require.NoError(t, err)
		assert.Equal(t, model.FeedModerationApproved, updatedFeed.ModerationStatus)
	})

	t.Run("Reject feed with reason", func(t *testing.T) {
		payload := map[string]interface{}{
			"note": "Content violates community guidelines",
		}
		resp := doJSON(router, http.MethodPut, "/api/v1/admin/content/feeds/"+uintToStr(feed2.ID)+"/reject", payload, "")
		assert.Equal(t, http.StatusOK, resp.Code)

		// Verify feed status changed
		updatedFeed, err := feedRepo.Get(ctx, feed2.ID)
		require.NoError(t, err)
		assert.Equal(t, model.FeedModerationRejected, updatedFeed.ModerationStatus)
	})

	t.Run("Reject feed without reason should fail", func(t *testing.T) {
		// Create another pending feed
		feed3 := &model.Feed{
			AuthorID:         user.ID,
			Content:          "Test feed content 3",
			Visibility:       model.FeedVisibilityPublic,
			ModerationStatus: model.FeedModerationPending,
		}
		require.NoError(t, feedRepo.Create(ctx, feed3))

		payload := map[string]interface{}{
			"note": "", // Empty reason
		}
		resp := doJSON(router, http.MethodPut, "/api/v1/admin/content/feeds/"+uintToStr(feed3.ID)+"/reject", payload, "")
		// Should fail with 400 Bad Request due to validation
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}

// TestBatchModerationFlow tests batch moderation operations
// Validates: Requirements 2.5
func TestBatchModerationFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateContentModerationModels(t, db)

	ctx := context.Background()
	userRepo := userrepo.NewUserRepository(db)
	feedRepo := feedrepo.NewFeedRepository(db)
	opLogRepo := operationlogrepo.NewOperationLogRepository(db)

	// Create admin user
	admin := &model.User{
		Name:         "BatchAdmin",
		Email:        "batch_admin@example.com",
		Phone:        "80000000003",
		PasswordHash: "hashed",
		Status:       model.UserStatusActive,
		Role:         model.RoleAdmin,
	}
	require.NoError(t, userRepo.Create(ctx, admin))

	// Create regular user
	user := &model.User{
		Name:         "BatchUser",
		Email:        "batch_user@example.com",
		Phone:        "80000000004",
		PasswordHash: "hashed",
		Status:       model.UserStatusActive,
		Role:         model.RoleUser,
	}
	require.NoError(t, userRepo.Create(ctx, user))

	// Create multiple pending feeds
	var feedIDs []uint64
	for i := 0; i < 5; i++ {
		feed := &model.Feed{
			AuthorID:         user.ID,
			Content:          "Batch test feed " + uintToStr(uint64(i)),
			Visibility:       model.FeedVisibilityPublic,
			ModerationStatus: model.FeedModerationPending,
		}
		require.NoError(t, feedRepo.Create(ctx, feed))
		feedIDs = append(feedIDs, feed.ID)
	}

	// Setup service and handler
	adminFeedSvc := contentservice.NewAdminFeedService(feedRepo, nil, opLogRepo)
	contentHandler := adminhandler.NewContentHandler(adminFeedSvc, nil, nil, nil)

	router := gin.New()
	api := router.Group("/api/v1/admin")
	auth := fakeContentAdminAuthMiddleware(admin.ID)
	api.Use(auth)
	api.POST("/content/feeds/batch-approve", contentHandler.BatchApproveFeed)
	api.POST("/content/feeds/batch-reject", contentHandler.BatchRejectFeed)

	t.Run("Batch approve feeds", func(t *testing.T) {
		payload := map[string]interface{}{
			"feedIds": feedIDs[:3], // Approve first 3 feeds
			"note":    "Batch approved",
		}
		resp := doJSON(router, http.MethodPost, "/api/v1/admin/content/feeds/batch-approve", payload, "")
		assert.Equal(t, http.StatusOK, resp.Code)

		// Verify all 3 feeds are approved
		for _, id := range feedIDs[:3] {
			feed, err := feedRepo.Get(ctx, id)
			require.NoError(t, err)
			assert.Equal(t, model.FeedModerationApproved, feed.ModerationStatus)
		}
	})

	t.Run("Batch reject feeds", func(t *testing.T) {
		payload := map[string]interface{}{
			"feedIds": feedIDs[3:], // Reject remaining 2 feeds
			"note":    "Batch rejected for policy violation",
		}
		resp := doJSON(router, http.MethodPost, "/api/v1/admin/content/feeds/batch-reject", payload, "")
		assert.Equal(t, http.StatusOK, resp.Code)

		// Verify all 2 feeds are rejected
		for _, id := range feedIDs[3:] {
			feed, err := feedRepo.Get(ctx, id)
			require.NoError(t, err)
			assert.Equal(t, model.FeedModerationRejected, feed.ModerationStatus)
		}
	})

	t.Run("Batch reject without reason should fail", func(t *testing.T) {
		// Create new pending feeds
		var newFeedIDs []uint64
		for i := 0; i < 2; i++ {
			feed := &model.Feed{
				AuthorID:         user.ID,
				Content:          "New batch test feed " + uintToStr(uint64(i)),
				Visibility:       model.FeedVisibilityPublic,
				ModerationStatus: model.FeedModerationPending,
			}
			require.NoError(t, feedRepo.Create(ctx, feed))
			newFeedIDs = append(newFeedIDs, feed.ID)
		}

		payload := map[string]interface{}{
			"feedIds": newFeedIDs,
			"note":    "", // Empty reason
		}
		resp := doJSON(router, http.MethodPost, "/api/v1/admin/content/feeds/batch-reject", payload, "")
		// Should fail with 400 Bad Request due to validation
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Batch approve empty list should fail", func(t *testing.T) {
		payload := map[string]interface{}{
			"feedIds": []uint64{},
			"note":    "Empty batch",
		}
		resp := doJSON(router, http.MethodPost, "/api/v1/admin/content/feeds/batch-approve", payload, "")
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}

func fakeContentAdminAuthMiddleware(adminID uint64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", adminID)
		c.Set("userID", adminID)
		c.Set("admin_id", adminID)
		c.Set("adminID", adminID)
		c.Next()
	}
}

func migrateContentModerationModels(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.AutoMigrate(
		&model.User{},
		&model.Feed{},
		&model.FeedImage{},
		&model.FeedReport{},
		&model.OperationLog{},
		&model.ContentCategory{},
	); err != nil {
		t.Fatalf("migrate content moderation models: %v", err)
	}
}
