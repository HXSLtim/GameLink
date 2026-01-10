// Package admin provides unit tests for review handlers.
package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/handler/testutil"
	"gamelink/internal/model"
	"gamelink/internal/repository/admin"
	"gamelink/internal/repository/game"
	"gamelink/internal/repository/gamecategory"
	"gamelink/internal/repository/implementations"
	"gamelink/internal/repository/payment"
	"gamelink/internal/repository/player"
	"gamelink/internal/repository/serviceitem"
	"gamelink/internal/repository/stats"
	"gamelink/internal/repository/user"
	adminservice "gamelink/internal/service/admin"
	"gamelink/pkg/cache"
)

// ============================================================================
// Test Setup
// ============================================================================

// ReviewTestContext provides test context for review handler tests.
type ReviewTestContext struct {
	Router     *gin.Engine
	Handler    *ReviewHandler
	Service    *adminservice.AdminService
	DB         *gorm.DB
	AdminUser  *model.User
	AdminToken string
}

// SetupReviewTest initializes test environment for review handler tests.
func SetupReviewTest(t *testing.T) *ReviewTestContext {
	t.Helper()

	// Setup test database
	db := testutil.SetupTestDB(t)

	// Create repositories
	games := game.NewGameRepository(db)
	users := user.NewUserRepository(db)
	players := player.NewPlayerRepository(db)
	ordersRepo := implementations.NewOrderRepository(db)
	payments := payment.NewPaymentRepository(db)
	roles := admin.NewRoleRepository(db)
	serviceItems := serviceitem.NewServiceItemRepository(db)
	permissions := admin.NewPermissionRepository(db)
	menus := admin.NewMenuRepository(db)
	statsRepo := stats.NewStatsRepository(db)
	gameCategories := gamecategory.NewGameCategoryRepository(db)
	c := cache.NewMemory()

	// Create admin service
	svc := adminservice.NewAdminService(
		games, users, players, ordersRepo, payments,
		roles, serviceItems, permissions, menus, statsRepo, nil, gameCategories, c,
	)

	// Setup router
	router := testutil.SetupGinTest(t)
	handler := NewReviewHandler(svc)

	// Create super admin user
	adminUser := testutil.CreateSuperAdmin(t, db)
	adminToken := testutil.GenerateTestToken(adminUser.ID)

	return &ReviewTestContext{
		Router:     router,
		Handler:    handler,
		Service:    svc,
		DB:         db,
		AdminUser:  adminUser,
		AdminToken: adminToken,
	}
}

// RegisterReviewRoutes registers review routes for testing.
func (ctx *ReviewTestContext) RegisterReviewRoutes() {
	group := ctx.Router.Group("/admin/reviews")
	{
		group.GET("", ctx.Handler.ListReviews)
		group.POST("", ctx.Handler.CreateReview)
		group.GET("/pending", ctx.Handler.ListPendingReviews)
		group.PUT("/batch-approve", ctx.Handler.BatchApproveReviews)
		group.PUT("/batch-reject", ctx.Handler.BatchRejectReviews)
		group.PUT("/approve-all-non-sensitive", ctx.Handler.ApproveAllNonSensitiveReviews)
		group.GET("/:id", ctx.Handler.GetReview)
		group.PUT("/:id", ctx.Handler.UpdateReview)
		group.DELETE("/:id", ctx.Handler.DeleteReview)
		group.GET("/:id/logs", ctx.Handler.ListReviewLogs)
		group.PUT("/:id/approve", ctx.Handler.ApproveReview)
		group.PUT("/:id/reject", ctx.Handler.RejectReview)
		group.POST("/:id/reports", ctx.Handler.CreateReviewReport)
	}

	// Review reports routes
	ctx.Router.GET("/admin/review-reports", ctx.Handler.ListReviewReports)
	ctx.Router.GET("/admin/review-reports/:id", ctx.Handler.GetReviewReport)
	ctx.Router.PUT("/admin/review-reports/:id/handle", ctx.Handler.HandleReviewReport)

	// Review replies routes
	ctx.Router.PUT("/admin/review-replies/:id", ctx.Handler.UpdateReply)
	ctx.Router.DELETE("/admin/review-replies/:id", ctx.Handler.DeleteReply)

	// Operation logs routes
	ctx.Router.GET("/admin/operation-logs", ctx.Handler.SearchOperationLogs)
	ctx.Router.GET("/admin/operation-logs/export", ctx.Handler.ExportOperationLogs)
}

// ============================================================================
// Test Data Helpers
// ============================================================================

// CreateTestReviewWithOrder creates a test review with associated order.
func (ctx *ReviewTestContext) CreateTestReviewWithOrder(t *testing.T, score model.Rating, status model.ReviewStatus) *model.Review {
	t.Helper()

	// Create test game
	testGame := testutil.CreateTestGame(t, ctx.DB)

	// Create test service item
	testServiceItem := testutil.CreateTestServiceItem(t, ctx.DB, testGame.ID, 5000, true)

	// Create test player
	testPlayer := testutil.CreateTestPlayer(t, ctx.DB, ctx.AdminUser.ID)

	// Create test user
	testUser := testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)

	// Create test order
	testOrder := testutil.CreateTestOrderWithItem(t, ctx.DB, testUser.ID, testPlayer.ID, testServiceItem.ID, model.OrderStatusCompleted)

	// Create test review
	review := &model.Review{
		Base:        model.Base{ExtJSON: "{}"},
		OrderID:     testOrder.ID,
		UserID:      testUser.ID,
		PlayerID:    testPlayer.ID,
		Score:       score,
		Content:     "Test review content",
		Status:      status,
		IsPublic:    true,
		IsAnonymous: false,
		Images:      model.StringArray{},
		ExpireAt:    time.Now().Add(7 * 24 * time.Hour),
	}

	require.NoError(t, ctx.DB.Create(review).Error, "Failed to create test review")
	return review
}

// ============================================================================
// ListReviews Tests
// ============================================================================

func TestReviewHandler_ListReviews_Success(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	// Create test reviews
	for i := 0; i < 5; i++ {
		ctx.CreateTestReviewWithOrder(t, model.Rating(5), model.ReviewStatusApproved)
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/reviews", ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	items, pagination := testutil.GetResponseList(t, w)
	assert.GreaterOrEqual(t, len(items), 5)
	assert.Equal(t, float64(1), pagination["page"])
}

func TestReviewHandler_ListReviews_WithPagination(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	// Create test reviews
	for i := 0; i < 25; i++ {
		ctx.CreateTestReviewWithOrder(t, model.Rating(5), model.ReviewStatusApproved)
	}

	w := testutil.MakeRequest(t, ctx.Router, "GET", "/admin/reviews?page=1&page_size=10", nil, testutil.WithAuth(ctx.AdminToken))
	testutil.AssertSuccess(t, w)

	items, pagination := testutil.GetResponseList(t, w)
	assert.LessOrEqual(t, len(items), 10)
	assert.Equal(t, float64(1), pagination["page"])
	assert.Equal(t, float64(10), pagination["page_size"])
}

func TestReviewHandler_ListReviews_WithOrderIDFilter(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	review := ctx.CreateTestReviewWithOrder(t, model.Rating(5), model.ReviewStatusApproved)

	w := testutil.MakeRequest(t, ctx.Router, "GET", fmt.Sprintf("/admin/reviews?order_id=%d", review.OrderID), nil, testutil.WithAuth(ctx.AdminToken))
	testutil.AssertSuccess(t, w)

	items, _ := testutil.GetResponseList(t, w)
	assert.GreaterOrEqual(t, len(items), 1)
}

func TestReviewHandler_ListReviews_WithUserIDFilter(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	review := ctx.CreateTestReviewWithOrder(t, model.Rating(5), model.ReviewStatusApproved)

	w := testutil.MakeRequest(t, ctx.Router, "GET", fmt.Sprintf("/admin/reviews?user_id=%d", review.UserID), nil, testutil.WithAuth(ctx.AdminToken))
	testutil.AssertSuccess(t, w)

	items, _ := testutil.GetResponseList(t, w)
	assert.GreaterOrEqual(t, len(items), 1)
}

func TestReviewHandler_ListReviews_WithPlayerIDFilter(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	review := ctx.CreateTestReviewWithOrder(t, model.Rating(5), model.ReviewStatusApproved)

	w := testutil.MakeRequest(t, ctx.Router, "GET", fmt.Sprintf("/admin/reviews?player_id=%d", review.PlayerID), nil, testutil.WithAuth(ctx.AdminToken))
	testutil.AssertSuccess(t, w)

	items, _ := testutil.GetResponseList(t, w)
	assert.GreaterOrEqual(t, len(items), 1)
}

func TestReviewHandler_ListReviews_WithDateFilter(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	ctx.CreateTestReviewWithOrder(t, model.Rating(5), model.ReviewStatusApproved)

	now := time.Now()
	dateFrom := now.Format("2006-01-02")
	dateTo := now.Add(24 * time.Hour).Format("2006-01-02")

	w := testutil.MakeRequest(t, ctx.Router, "GET", fmt.Sprintf("/admin/reviews?date_from=%s&date_to=%s", dateFrom, dateTo), nil, testutil.WithAuth(ctx.AdminToken))
	testutil.AssertSuccess(t, w)

	items, _ := testutil.GetResponseList(t, w)
	assert.GreaterOrEqual(t, len(items), 1)
}

// ============================================================================
// GetReview Tests
// ============================================================================

func TestReviewHandler_GetReview_Success(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	review := ctx.CreateTestReviewWithOrder(t, model.Rating(5), model.ReviewStatusApproved)

	path := fmt.Sprintf("/admin/reviews/%d", review.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", path, ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(review.ID), data["id"])
	assert.Equal(t, "approved", data["status"])
}

func TestReviewHandler_GetReview_NotFound(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/reviews/999999", ctx.AdminToken, nil)
	testutil.AssertError(t, w, http.StatusNotFound)
}

func TestReviewHandler_GetReview_InvalidID(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/reviews/invalid", ctx.AdminToken, nil)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

// ============================================================================
// CreateReview Tests
// ============================================================================

func TestReviewHandler_CreateReview_Success(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	// Create test game and service item
	testGame := testutil.CreateTestGame(t, ctx.DB)
	testServiceItem := testutil.CreateTestServiceItem(t, ctx.DB, testGame.ID, 5000, true)

	// Create test player and user
	testPlayer := testutil.CreateTestPlayer(t, ctx.DB, ctx.AdminUser.ID)
	testUser := testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)
	testOrder := testutil.CreateTestOrderWithItem(t, ctx.DB, testUser.ID, testPlayer.ID, testServiceItem.ID, model.OrderStatusCompleted)

	payload := map[string]interface{}{
		"order_id":  testOrder.ID,
		"user_id":   testUser.ID,
		"player_id": testPlayer.ID,
		"score":     5,
		"content":   "Great service!",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/reviews", ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w, http.StatusCreated)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.NotEmpty(t, data["id"])
	assert.Equal(t, "Great service!", data["comment"])
}

func TestReviewHandler_CreateReview_ValidationError(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	// Missing required fields
	payload := map[string]interface{}{
		"score": 5,
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/reviews", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

func TestReviewHandler_CreateReview_InvalidScore(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	payload := map[string]interface{}{
		"order_id":  1,
		"user_id":   1,
		"player_id": 1,
		"score":     10, // Invalid score (should be 1-5)
		"content":   "Test",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/reviews", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

// ============================================================================
// UpdateReview Tests
// ============================================================================

func TestReviewHandler_UpdateReview_Success(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	review := ctx.CreateTestReviewWithOrder(t, model.Rating(3), model.ReviewStatusApproved)

	payload := map[string]interface{}{
		"score":   5,
		"content": "Updated review content",
	}

	path := fmt.Sprintf("/admin/reviews/%d", review.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", path, ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	// Verify update in DB
	var updatedReview model.Review
	ctx.DB.First(&updatedReview, review.ID)
	assert.Equal(t, "Updated review content", updatedReview.Content)
}

func TestReviewHandler_UpdateReview_NotFound(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	payload := map[string]interface{}{
		"score":   5,
		"content": "Updated content",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", "/admin/reviews/999999", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusNotFound)
}

// ============================================================================
// DeleteReview Tests
// ============================================================================

func TestReviewHandler_DeleteReview_Success(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	review := ctx.CreateTestReviewWithOrder(t, model.Rating(5), model.ReviewStatusApproved)

	path := fmt.Sprintf("/admin/reviews/%d", review.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "DELETE", path, ctx.AdminToken, nil)
	testutil.AssertDeleted(t, w)

	// Verify soft deletion (record should still exist but with deleted status)
	var deletedReview model.Review
	err := ctx.DB.Unscoped().First(&deletedReview, review.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.ReviewStatusDeleted, deletedReview.Status)
}

func TestReviewHandler_DeleteReview_NotFound(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "DELETE", "/admin/reviews/999999", ctx.AdminToken, nil)
	testutil.AssertError(t, w, http.StatusNotFound)
}

// ============================================================================
// ListPendingReviews Tests
// ============================================================================

func TestReviewHandler_ListPendingReviews_Success(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	// Create pending reviews
	for i := 0; i < 5; i++ {
		ctx.CreateTestReviewWithOrder(t, model.Rating(5), model.ReviewStatusPending)
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/reviews/pending", ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	items := data["items"].([]interface{})
	assert.GreaterOrEqual(t, len(items), 5)
}

func TestReviewHandler_ListPendingReviews_WithSensitiveWordsFilter(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	ctx.CreateTestReviewWithOrder(t, model.Rating(5), model.ReviewStatusPending)

	w := testutil.MakeRequest(t, ctx.Router, "GET", "/admin/reviews/pending?hasSensitiveWords=false", nil, testutil.WithAuth(ctx.AdminToken))
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	items := data["items"].([]interface{})
	// Should return reviews without sensitive words
	assert.NotNil(t, items)
}

// ============================================================================
// ApproveReview Tests
// ============================================================================

func TestReviewHandler_ApproveReview_Success(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	review := ctx.CreateTestReviewWithOrder(t, model.Rating(5), model.ReviewStatusPending)

	payload := map[string]interface{}{
		"reason": "Review approved",
	}

	path := fmt.Sprintf("/admin/reviews/%d/approve", review.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", path, ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	// Verify status update
	var updatedReview model.Review
	ctx.DB.First(&updatedReview, review.ID)
	assert.Equal(t, model.ReviewStatusApproved, updatedReview.Status)
}

func TestReviewHandler_ApproveReview_WithoutReason(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	review := ctx.CreateTestReviewWithOrder(t, model.Rating(5), model.ReviewStatusPending)

	path := fmt.Sprintf("/admin/reviews/%d/approve", review.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", path, ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	// Verify status update
	var updatedReview model.Review
	ctx.DB.First(&updatedReview, review.ID)
	assert.Equal(t, model.ReviewStatusApproved, updatedReview.Status)
}

// ============================================================================
// RejectReview Tests
// ============================================================================

func TestReviewHandler_RejectReview_Success(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	review := ctx.CreateTestReviewWithOrder(t, model.Rating(5), model.ReviewStatusPending)

	payload := map[string]interface{}{
		"reason": "Inappropriate content",
	}

	path := fmt.Sprintf("/admin/reviews/%d/reject", review.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", path, ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	// Verify status update
	var updatedReview model.Review
	ctx.DB.First(&updatedReview, review.ID)
	assert.Equal(t, model.ReviewStatusRejected, updatedReview.Status)
	assert.Equal(t, "Inappropriate content", updatedReview.RejectionReason)
}

func TestReviewHandler_RejectReview_MissingReason(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	review := ctx.CreateTestReviewWithOrder(t, model.Rating(5), model.ReviewStatusPending)

	payload := map[string]interface{}{}

	path := fmt.Sprintf("/admin/reviews/%d/reject", review.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", path, ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

// ============================================================================
// BatchApproveReviews Tests
// ============================================================================

func TestReviewHandler_BatchApproveReviews_Success(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	var reviewIDs []uint64
	for i := 0; i < 3; i++ {
		review := ctx.CreateTestReviewWithOrder(t, model.Rating(5), model.ReviewStatusPending)
		reviewIDs = append(reviewIDs, review.ID)
	}

	payload := map[string]interface{}{
		"reviewIds": reviewIDs,
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", "/admin/reviews/batch-approve", ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	// Verify all reviews are approved
	for _, id := range reviewIDs {
		var review model.Review
		ctx.DB.First(&review, id)
		assert.Equal(t, model.ReviewStatusApproved, review.Status)
	}
}

func TestReviewHandler_BatchApproveReviews_EmptyList(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	payload := map[string]interface{}{
		"reviewIds": []uint64{},
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", "/admin/reviews/batch-approve", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

// ============================================================================
// BatchRejectReviews Tests
// ============================================================================

func TestReviewHandler_BatchRejectReviews_Success(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	var reviewIDs []uint64
	for i := 0; i < 3; i++ {
		review := ctx.CreateTestReviewWithOrder(t, model.Rating(5), model.ReviewStatusPending)
		reviewIDs = append(reviewIDs, review.ID)
	}

	payload := map[string]interface{}{
		"reviewIds": reviewIDs,
		"reason":    "Batch rejection",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", "/admin/reviews/batch-reject", ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	// Verify all reviews are rejected
	for _, id := range reviewIDs {
		var review model.Review
		ctx.DB.First(&review, id)
		assert.Equal(t, model.ReviewStatusRejected, review.Status)
		assert.Equal(t, "Batch rejection", review.RejectionReason)
	}
}

func TestReviewHandler_BatchRejectReviews_MissingReason(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	review := ctx.CreateTestReviewWithOrder(t, model.Rating(5), model.ReviewStatusPending)

	payload := map[string]interface{}{
		"reviewIds": []uint64{review.ID},
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", "/admin/reviews/batch-reject", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

// ============================================================================
// ListReviewLogs Tests
// ============================================================================

func TestReviewHandler_ListReviewLogs_Success(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	review := ctx.CreateTestReviewWithOrder(t, model.Rating(5), model.ReviewStatusApproved)

	path := fmt.Sprintf("/admin/reviews/%d/logs", review.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", path, ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	items := data["items"].([]interface{})
	assert.NotNil(t, items)
}

func TestReviewHandler_ListReviewLogs_WithPagination(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	review := ctx.CreateTestReviewWithOrder(t, model.Rating(5), model.ReviewStatusApproved)

	path := fmt.Sprintf("/admin/reviews/%d/logs?page=1&page_size=10", review.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", path, ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	items, pagination := testutil.GetResponseList(t, w)
	assert.NotNil(t, items)
	assert.Equal(t, float64(1), pagination["page"])
	assert.Equal(t, float64(10), pagination["page_size"])
}

func TestReviewHandler_ListReviewLogs_WithActionFilter(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	review := ctx.CreateTestReviewWithOrder(t, model.Rating(5), model.ReviewStatusApproved)

	path := fmt.Sprintf("/admin/reviews/%d/logs?action=approve", review.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", path, ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)
}

// ============================================================================
// ListPlayerReviews Tests
// ============================================================================

func TestReviewHandler_ListPlayerReviews_Success(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	review := ctx.CreateTestReviewWithOrder(t, model.Rating(5), model.ReviewStatusApproved)

	path := fmt.Sprintf("/admin/players/%d/reviews", review.PlayerID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", path, ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	items, _ := testutil.GetResponseList(t, w)
	assert.GreaterOrEqual(t, len(items), 1)
}

func TestReviewHandler_ListPlayerReviews_NotFound(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/players/999999/reviews", ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w) // Empty list, not error

	items, _ := testutil.GetResponseList(t, w)
	assert.Equal(t, 0, len(items))
}

// ============================================================================
// CreateReviewReport Tests
// ============================================================================

func TestReviewHandler_CreateReviewReport_Success(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	review := ctx.CreateTestReviewWithOrder(t, model.Rating(5), model.ReviewStatusApproved)

	payload := map[string]interface{}{
		"reason":   "Inappropriate content",
		"evidence": "Evidence URL",
	}

	path := fmt.Sprintf("/admin/reviews/%d/reports", review.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", path, ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w, http.StatusCreated)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.NotEmpty(t, data["reportId"])
}

func TestReviewHandler_CreateReviewReport_MissingReason(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	review := ctx.CreateTestReviewWithOrder(t, model.Rating(5), model.ReviewStatusApproved)

	payload := map[string]interface{}{
		"evidence": "Evidence URL",
	}

	path := fmt.Sprintf("/admin/reviews/%d/reports", review.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", path, ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

// ============================================================================
// ListReviewReports Tests
// ============================================================================

func TestReviewHandler_ListReviewReports_Success(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/review-reports", ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	items := data["items"].([]interface{})
	assert.NotNil(t, items)
}

func TestReviewHandler_ListReviewReports_WithStatusFilter(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	w := testutil.MakeRequest(t, ctx.Router, "GET", "/admin/review-reports?status=pending", nil, testutil.WithAuth(ctx.AdminToken))
	testutil.AssertSuccess(t, w)
}

func TestReviewHandler_ListReviewReports_WithReviewIDFilter(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	review := ctx.CreateTestReviewWithOrder(t, model.Rating(5), model.ReviewStatusApproved)

	w := testutil.MakeRequest(t, ctx.Router, "GET", fmt.Sprintf("/admin/review-reports?review_id=%d", review.ID), nil, testutil.WithAuth(ctx.AdminToken))
	testutil.AssertSuccess(t, w)
}

// ============================================================================
// GetReviewReport Tests
// ============================================================================

func TestReviewHandler_GetReviewReport_Success(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	// Create a review report first
	review := ctx.CreateTestReviewWithOrder(t, model.Rating(5), model.ReviewStatusApproved)

	// Create report directly in DB
	report := &model.ReviewReport{
		Base:        model.Base{ExtJSON: "{}"},
		ReviewID:    review.ID,
		ReporterID:  ctx.AdminUser.ID,
		Reason:      "Test report",
		Status:      model.ReviewReportStatusPending,
	}

	require.NoError(t, ctx.DB.Create(report).Error)

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", fmt.Sprintf("/admin/review-reports/%d", report.ID), ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(report.ID), data["id"])
}

func TestReviewHandler_GetReviewReport_NotFound(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/review-reports/999999", ctx.AdminToken, nil)
	testutil.AssertError(t, w, http.StatusNotFound)
}

// ============================================================================
// HandleReviewReport Tests
// ============================================================================

func TestReviewHandler_HandleReviewReport_DeleteAction(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	review := ctx.CreateTestReviewWithOrder(t, model.Rating(5), model.ReviewStatusApproved)

	report := &model.ReviewReport{
		Base:        model.Base{ExtJSON: "{}"},
		ReviewID:    review.ID,
		ReporterID:  ctx.AdminUser.ID,
		Reason:      "Test report",
		Status:      model.ReviewReportStatusPending,
	}

	require.NoError(t, ctx.DB.Create(report).Error)

	payload := map[string]interface{}{
		"action": "delete",
		"note":   "Deleted after review",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", fmt.Sprintf("/admin/review-reports/%d/handle", report.ID), ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	// Verify report status updated
	var updatedReport model.ReviewReport
	ctx.DB.First(&updatedReport, report.ID)
	assert.Equal(t, model.ReviewReportStatusApproved, updatedReport.Status)
}

func TestReviewHandler_HandleReviewReport_WarnAction(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	review := ctx.CreateTestReviewWithOrder(t, model.Rating(5), model.ReviewStatusApproved)

	report := &model.ReviewReport{
		Base:        model.Base{ExtJSON: "{}"},
		ReviewID:    review.ID,
		ReporterID:  ctx.AdminUser.ID,
		Reason:      "Test report",
		Status:      model.ReviewReportStatusPending,
	}

	require.NoError(t, ctx.DB.Create(report).Error)

	payload := map[string]interface{}{
		"action": "warn",
		"note":   "Warning issued",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", fmt.Sprintf("/admin/review-reports/%d/handle", report.ID), ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)
}

func TestReviewHandler_HandleReviewReport_RejectAction(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	review := ctx.CreateTestReviewWithOrder(t, model.Rating(5), model.ReviewStatusApproved)

	report := &model.ReviewReport{
		Base:        model.Base{ExtJSON: "{}"},
		ReviewID:    review.ID,
		ReporterID:  ctx.AdminUser.ID,
		Reason:      "Test report",
		Status:      model.ReviewReportStatusPending,
	}

	require.NoError(t, ctx.DB.Create(report).Error)

	payload := map[string]interface{}{
		"action": "reject",
		"note":   "Report rejected",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", fmt.Sprintf("/admin/review-reports/%d/handle", report.ID), ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	// Verify report status
	var updatedReport model.ReviewReport
	ctx.DB.First(&updatedReport, report.ID)
	assert.Equal(t, model.ReviewReportStatusRejected, updatedReport.Status)
}

func TestReviewHandler_HandleReviewReport_InvalidAction(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	review := ctx.CreateTestReviewWithOrder(t, model.Rating(5), model.ReviewStatusApproved)

	report := &model.ReviewReport{
		Base:        model.Base{ExtJSON: "{}"},
		ReviewID:    review.ID,
		ReporterID:  ctx.AdminUser.ID,
		Reason:      "Test report",
		Status:      model.ReviewReportStatusPending,
	}

	require.NoError(t, ctx.DB.Create(report).Error)

	payload := map[string]interface{}{
		"action": "invalid_action",
		"note":   "Test",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", fmt.Sprintf("/admin/review-reports/%d/handle", report.ID), ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

// ============================================================================
// ApproveAllNonSensitiveReviews Tests
// ============================================================================

func TestReviewHandler_ApproveAllNonSensitiveReviews_Success(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	// Create pending reviews
	for i := 0; i < 3; i++ {
		ctx.CreateTestReviewWithOrder(t, model.Rating(5), model.ReviewStatusPending)
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", "/admin/reviews/approve-all-non-sensitive", ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.GreaterOrEqual(t, int(data["count"].(float64)), 0)
}

func TestReviewHandler_ApproveAllNonSensitiveReviews_NoReviews(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", "/admin/reviews/approve-all-non-sensitive", ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(0), data["count"])
}

// ============================================================================
// UpdateReply Tests
// ============================================================================

func TestReviewHandler_UpdateReply_Success(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	// Create a review reply
	review := ctx.CreateTestReviewWithOrder(t, model.Rating(5), model.ReviewStatusApproved)

	reply := &model.ReviewReply{
		Base:        model.Base{ExtJSON: "{}"},
		ReviewID:    review.ID,
		PlayerID:    review.PlayerID,
		Content:     "Original reply",
		Status:      model.ReviewReplyStatusApproved,
	}

	require.NoError(t, ctx.DB.Create(reply).Error)

	payload := map[string]interface{}{
		"content": "Updated reply content",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", fmt.Sprintf("/admin/review-replies/%d", reply.ID), ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	// Verify update
	var updatedReply model.ReviewReply
	ctx.DB.First(&updatedReply, reply.ID)
	assert.Equal(t, "Updated reply content", updatedReply.Content)
}

func TestReviewHandler_UpdateReply_NotFound(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	payload := map[string]interface{}{
		"content": "Updated reply",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", "/admin/review-replies/999999", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusNotFound)
}

// ============================================================================
// DeleteReply Tests
// ============================================================================

func TestReviewHandler_DeleteReply_Success(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	review := ctx.CreateTestReviewWithOrder(t, model.Rating(5), model.ReviewStatusApproved)

	reply := &model.ReviewReply{
		Base:        model.Base{ExtJSON: "{}"},
		ReviewID:    review.ID,
		PlayerID:    review.PlayerID,
		Content:     "Reply to delete",
		Status:      model.ReviewReplyStatusApproved,
	}

	require.NoError(t, ctx.DB.Create(reply).Error)

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "DELETE", fmt.Sprintf("/admin/review-replies/%d", reply.ID), ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	// Verify deletion
	var count int64
	ctx.DB.Model(&model.ReviewReply{}).Where("id = ?", reply.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestReviewHandler_DeleteReply_NotFound(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "DELETE", "/admin/review-replies/999999", ctx.AdminToken, nil)
	testutil.AssertError(t, w, http.StatusNotFound)
}

// ============================================================================
// SearchOperationLogs Tests
// ============================================================================

func TestReviewHandler_SearchOperationLogs_Success(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/operation-logs", ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	items := data["items"].([]interface{})
	assert.NotNil(t, items)
}

func TestReviewHandler_SearchOperationLogs_WithEntityTypeFilter(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	w := testutil.MakeRequest(t, ctx.Router, "GET", "/admin/operation-logs?entity_type=review", nil, testutil.WithAuth(ctx.AdminToken))
	testutil.AssertSuccess(t, w)
}

func TestReviewHandler_SearchOperationLogs_WithActionFilter(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	w := testutil.MakeRequest(t, ctx.Router, "GET", "/admin/operation-logs?action=create", nil, testutil.WithAuth(ctx.AdminToken))
	testutil.AssertSuccess(t, w)
}

func TestReviewHandler_SearchOperationLogs_WithDateFilter(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	now := time.Now()
	dateFrom := now.Format("2006-01-02")
	dateTo := now.Add(24 * time.Hour).Format("2006-01-02")

	w := testutil.MakeRequest(t, ctx.Router, "GET", fmt.Sprintf("/admin/operation-logs?date_from=%s&date_to=%s", dateFrom, dateTo), nil, testutil.WithAuth(ctx.AdminToken))
	testutil.AssertSuccess(t, w)
}

// ============================================================================
// ExportOperationLogs Tests
// ============================================================================

func TestReviewHandler_ExportOperationLogs_Success(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	w := testutil.MakeRequest(t, ctx.Router, "GET", "/admin/operation-logs/export?export=csv", nil, testutil.WithAuth(ctx.AdminToken))
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/csv", w.Header().Get("Content-Type"))
}

func TestReviewHandler_ExportOperationLogs_WithFilters(t *testing.T) {
	ctx := SetupReviewTest(t)
	ctx.RegisterReviewRoutes()

	w := testutil.MakeRequest(t, ctx.Router, "GET", "/admin/operation-logs/export?entity_type=review", nil, testutil.WithAuth(ctx.AdminToken))
	assert.Equal(t, http.StatusOK, w.Code)
}
