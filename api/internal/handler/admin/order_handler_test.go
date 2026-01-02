// Package admin provides unit tests for order handlers.
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
	adminservice "gamelink/internal/service/admin"
	"gamelink/internal/repository/admin"
	"gamelink/internal/repository/game"
	"gamelink/internal/repository/gamecategory"
	"gamelink/internal/repository/implementations"
	"gamelink/internal/repository/payment"
	"gamelink/internal/repository/player"
	"gamelink/internal/repository/serviceitem"
	"gamelink/internal/repository/stats"
	"gamelink/internal/repository/user"
	"gamelink/pkg/cache"
)

// ============================================================================
// Test Setup
// ============================================================================

// OrderTestContext provides test context for order handler tests.
type OrderTestContext struct {
	Router      *gin.Engine
	Handler     *OrderHandler
	Service     *adminservice.AdminService
	DB          *gorm.DB
	AdminUser   *model.User
	AdminToken  string
	TestUser    *model.User
	TestPlayer  *model.Player
	TestGame    *model.Game
}

// SetupOrderTest initializes test environment for order handler tests.
func SetupOrderTest(t *testing.T) *OrderTestContext {
	t.Helper()

	// Setup test database
	db := testutil.SetupTestDB(t)

	// Create repositories
	games := game.NewGameRepository(db)
	users := user.NewUserRepository(db)
	players := player.NewPlayerRepository(db)
	orders := implementations.NewOrderRepository(db)
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
		games, users, players, orders, payments,
		roles, serviceItems, permissions, menus, statsRepo, nil, gameCategories, c,
	)

	// Setup router
	router := testutil.SetupGinTest(t)
	handler := NewOrderHandler(svc)

	// Create super admin user
	adminUser := testutil.CreateSuperAdmin(t, db)
	adminToken := testutil.GenerateTestToken(adminUser.ID)

	// Create test data
	testUser := testutil.CreateAdminUser(t, db, model.RoleUser)
	testGame := testutil.CreateTestGame(t, db)

	testPlayer := &model.Player{
		Base: model.Base{
			ExtJSON: "{}",
		},
		UserID:             testUser.ID,
		Nickname:           "Test Player",
		Bio:                "Test bio",
		HourlyRateCents:    5000,
		MainGameID:         testGame.ID,
		VerificationStatus: model.VerificationVerified,
	}
	require.NoError(t, db.Create(testPlayer).Error)

	return &OrderTestContext{
		Router:     router,
		Handler:    handler,
		Service:    svc,
		DB:         db,
		AdminUser:  adminUser,
		AdminToken: adminToken,
		TestUser:   testUser,
		TestPlayer: testPlayer,
		TestGame:   testGame,
	}
}

// RegisterOrderRoutes registers order routes for testing.
func (ctx *OrderTestContext) RegisterOrderRoutes() {
	group := ctx.Router.Group("/admin/orders")
	{
		group.POST("", ctx.Handler.CreateOrder)
		group.GET("", ctx.Handler.ListOrders)
		group.GET("/:id", ctx.Handler.GetOrder)
		group.PUT("/:id", ctx.Handler.UpdateOrder)
		group.DELETE("/:id", ctx.Handler.DeleteOrder)
		group.POST("/:id/assign", ctx.Handler.AssignOrder)
		group.POST("/:id/confirm", ctx.Handler.ConfirmOrder)
		group.POST("/:id/start", ctx.Handler.StartOrder)
		group.POST("/:id/complete", ctx.Handler.CompleteOrder)
		group.POST("/:id/cancel", ctx.Handler.CancelOrder)
		group.POST("/:id/review", ctx.Handler.ReviewOrder)
		group.POST("/:id/refund", ctx.Handler.RefundOrder)
		group.GET("/:id/timeline", ctx.Handler.GetOrderTimeline)
		group.GET("/:id/payments", ctx.Handler.ListOrderPayments)
		group.GET("/:id/refunds", ctx.Handler.ListOrderRefunds)
		group.GET("/:id/reviews", ctx.Handler.ListOrderReviews)
		group.GET("/:id/logs", ctx.Handler.ListOrderLogs)
	}
}

// ============================================================================
// CreateOrder Tests
// ============================================================================

func TestOrderHandler_Unit_CreateOrder_Success(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	payload := map[string]interface{}{
		"user_id":           ctx.TestUser.ID,
		"player_id":         ctx.TestPlayer.ID,
		"game_id":           ctx.TestGame.ID,
		"title":             "Test Order",
		"description":       "Test description",
		"total_price_cents": 10000,
		"currency":          "USD",
		"scheduled_start":   time.Now().Add(1 * time.Hour).Format(time.RFC3339),
		"scheduled_end":     time.Now().Add(2 * time.Hour).Format(time.RFC3339),
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/orders", ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w, http.StatusCreated)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.NotEmpty(t, data["id"])
	assert.Equal(t, "Test Order", data["title"])
}

func TestOrderHandler_Unit_CreateOrder_ValidationError(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	// Missing required fields
	payload := map[string]interface{}{
		"title": "Test Order",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/orders", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

func TestOrderHandler_Unit_CreateOrder_InvalidTime(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	payload := map[string]interface{}{
		"user_id":           ctx.TestUser.ID,
		"game_id":           ctx.TestGame.ID,
		"title":             "Test Order",
		"total_price_cents": 10000,
		"currency":          "USD",
		"scheduled_start":   "invalid-time",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/orders", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

// ============================================================================
// ListOrders Tests
// ============================================================================

func TestOrderHandler_Unit_ListOrders_Success(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	// Create test order
	_ = testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusPending)

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/orders", ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	items, pagination := testutil.GetResponseList(t, w)
	assert.GreaterOrEqual(t, len(items), 1)
	assert.Equal(t, float64(1), pagination["page"])
}

func TestOrderHandler_Unit_ListOrders_WithFilters(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	// Create test orders with different statuses
	testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusPending)
	testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusCompleted)

	// Filter by status
	params := map[string]string{"status": "pending"}
	w := testutil.MakeRequest(t, ctx.Router, "GET", "/admin/orders?status=pending", nil, testutil.WithAuth(ctx.AdminToken), testutil.WithQuery(params))
	testutil.AssertSuccess(t, w)

	items, _ := testutil.GetResponseList(t, w)
	// Should only return pending orders
	for _, item := range items {
		order := item.(map[string]interface{})
		assert.Equal(t, "pending", order["status"])
	}
}

func TestOrderHandler_Unit_ListOrders_WithPagination(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	// Create multiple orders
	for i := 0; i < 25; i++ {
		testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusPending)
	}

	w := testutil.MakeRequest(t, ctx.Router, "GET", "/admin/orders?page=1&page_size=10", nil, testutil.WithAuth(ctx.AdminToken))
	testutil.AssertSuccess(t, w)

	items, pagination := testutil.GetResponseList(t, w)
	assert.LessOrEqual(t, len(items), 10)
	assert.Equal(t, float64(1), pagination["page"])
	assert.Equal(t, float64(10), pagination["page_size"])
}

// ============================================================================
// GetOrder Tests
// ============================================================================

func TestOrderHandler_Unit_GetOrder_Success(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusPending)

	path := fmt.Sprintf("/admin/orders/%d", order.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", path, ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(order.ID), data["id"])
	assert.Equal(t, "Test Order", data["title"])
}

func TestOrderHandler_Unit_GetOrder_NotFound(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/orders/999999", ctx.AdminToken, nil)
	testutil.AssertError(t, w, http.StatusNotFound)
}

func TestOrderHandler_Unit_GetOrder_InvalidID(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/orders/invalid", ctx.AdminToken, nil)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

// ============================================================================
// UpdateOrder Tests
// ============================================================================

func TestOrderHandler_Unit_UpdateOrder_Success(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusPending)

	payload := map[string]interface{}{
		"status":            "confirmed",
		"total_price_cents": 15000,
		"currency":          "USD",
	}

	path := fmt.Sprintf("/admin/orders/%d", order.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", path, ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	// Verify status changed in DB
	testutil.AssertOrderStatus(t, ctx.DB, order.ID, model.OrderStatusConfirmed)
}

func TestOrderHandler_Unit_UpdateOrder_InvalidTransition(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusCompleted)

	// Try to transition from completed to pending (invalid)
	payload := map[string]interface{}{
		"status":            "pending",
		"total_price_cents": 10000,
		"currency":          "USD",
	}

	path := fmt.Sprintf("/admin/orders/%d", order.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", path, ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

// ============================================================================
// DeleteOrder Tests
// ============================================================================

func TestOrderHandler_Unit_DeleteOrder_Success(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusPending)

	path := fmt.Sprintf("/admin/orders/%d", order.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "DELETE", path, ctx.AdminToken, nil)
	testutil.AssertDeleted(t, w)

	// Verify order deleted
	var count int64
	ctx.DB.Model(&model.Order{}).Where("id = ?", order.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestOrderHandler_Unit_DeleteOrder_NotFound(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "DELETE", "/admin/orders/999999", ctx.AdminToken, nil)
	testutil.AssertError(t, w, http.StatusNotFound)
}

// ============================================================================
// AssignOrder Tests
// ============================================================================

func TestOrderHandler_Unit_AssignOrder_Success(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, 0, ctx.TestGame.ID, model.OrderStatusPending)

	payload := map[string]interface{}{
		"player_id": ctx.TestPlayer.ID,
	}

	path := fmt.Sprintf("/admin/orders/%d/assign", order.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", path, ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	// Verify player assigned
	var updatedOrder model.Order
	ctx.DB.First(&updatedOrder, order.ID)
	assert.NotNil(t, updatedOrder.PlayerID)
	assert.Equal(t, ctx.TestPlayer.ID, *updatedOrder.PlayerID)
}

func TestOrderHandler_Unit_AssignOrder_PlayerNotFound(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, 0, ctx.TestGame.ID, model.OrderStatusPending)

	payload := map[string]interface{}{
		"player_id": uint64(999999),
	}

	path := fmt.Sprintf("/admin/orders/%d/assign", order.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", path, ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusNotFound)
}

// ============================================================================
// ConfirmOrder Tests
// ============================================================================

func TestOrderHandler_Unit_ConfirmOrder_Success(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusPending)

	payload := map[string]interface{}{
		"note": "Order confirmed",
	}

	path := fmt.Sprintf("/admin/orders/%d/confirm", order.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", path, ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	testutil.AssertOrderStatus(t, ctx.DB, order.ID, model.OrderStatusConfirmed)
}

func TestOrderHandler_Unit_ConfirmOrder_InvalidTransition(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusCompleted)

	path := fmt.Sprintf("/admin/orders/%d/confirm", order.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", path, ctx.AdminToken, map[string]interface{}{})
	testutil.AssertError(t, w, http.StatusBadRequest)
}

// ============================================================================
// StartOrder Tests
// ============================================================================

func TestOrderHandler_Unit_StartOrder_Success(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusConfirmed)

	path := fmt.Sprintf("/admin/orders/%d/start", order.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", path, ctx.AdminToken, map[string]interface{}{})
	testutil.AssertSuccess(t, w)

	testutil.AssertOrderStatus(t, ctx.DB, order.ID, model.OrderStatusInProgress)
}

// ============================================================================
// CompleteOrder Tests
// ============================================================================

func TestOrderHandler_Unit_CompleteOrder_Success(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusInProgress)

	path := fmt.Sprintf("/admin/orders/%d/complete", order.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", path, ctx.AdminToken, map[string]interface{}{})
	testutil.AssertSuccess(t, w)

	testutil.AssertOrderStatus(t, ctx.DB, order.ID, model.OrderStatusCompleted)
}

// ============================================================================
// CancelOrder Tests
// ============================================================================

func TestOrderHandler_Unit_CancelOrder_Success(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusConfirmed)

	payload := map[string]interface{}{
		"reason": "Customer requested cancellation",
	}

	path := fmt.Sprintf("/admin/orders/%d/cancel", order.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", path, ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	testutil.AssertOrderStatus(t, ctx.DB, order.ID, model.OrderStatusCanceled)
}

// ============================================================================
// ReviewOrder Tests
// ============================================================================

func TestOrderHandler_Unit_ReviewOrder_Approve(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusPending)

	payload := map[string]interface{}{
		"approved": true,
		"reason":   "Order looks good",
	}

	path := fmt.Sprintf("/admin/orders/%d/review", order.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", path, ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	testutil.AssertOrderStatus(t, ctx.DB, order.ID, model.OrderStatusConfirmed)
}

func TestOrderHandler_Unit_ReviewOrder_Reject(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusPending)

	payload := map[string]interface{}{
		"approved": false,
		"reason":   "Invalid information",
	}

	path := fmt.Sprintf("/admin/orders/%d/review", order.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", path, ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	testutil.AssertOrderStatus(t, ctx.DB, order.ID, model.OrderStatusCanceled)
}

// ============================================================================
// RefundOrder Tests
// ============================================================================

func TestOrderHandler_Unit_RefundOrder_Success(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusCompleted)
	_ = testutil.CreateTestPayment(t, ctx.DB, order.ID, ctx.TestUser.ID, model.PaymentStatusPaid)

	payload := map[string]interface{}{
		"reason":       "Customer service refund",
		"amount_cents": int64(5000),
		"note":         "Partial refund",
	}

	path := fmt.Sprintf("/admin/orders/%d/refund", order.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", path, ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)
}

// ============================================================================
// GetOrderTimeline Tests
// ============================================================================

func TestOrderHandler_Unit_GetOrderTimeline_Success(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusCompleted)

	path := fmt.Sprintf("/admin/orders/%d/timeline", order.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", path, ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"]
	assert.NotNil(t, data)
}

// ============================================================================
// ListOrderPayments Tests
// ============================================================================

func TestOrderHandler_Unit_ListOrderPayments_Success(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusPending)
	_ = testutil.CreateTestPayment(t, ctx.DB, order.ID, ctx.TestUser.ID, model.PaymentStatusPending)

	path := fmt.Sprintf("/admin/orders/%d/payments", order.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", path, ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].([]interface{})
	assert.GreaterOrEqual(t, len(data), 1)
}

// ============================================================================
// ListOrderLogs Tests
// ============================================================================

func TestOrderHandler_Unit_ListOrderLogs_Success(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusPending)

	path := fmt.Sprintf("/admin/orders/%d/logs", order.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", path, ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	items, pagination := testutil.GetResponseList(t, w)
	assert.NotNil(t, items)
	assert.NotNil(t, pagination)
}
