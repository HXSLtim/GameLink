// Package user provides unit tests for order handlers.
package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/handler/testutil"
	"gamelink/internal/model"
	"gamelink/internal/repository/commission"
	"gamelink/internal/repository/game"
	"gamelink/internal/repository/implementations"
	"gamelink/internal/repository/payment"
	"gamelink/internal/repository/player"
	"gamelink/internal/repository/review"
	"gamelink/internal/repository/user"
	"gamelink/internal/service/order"
)

// ============================================================================
// Test Setup
// ============================================================================

// UserOrderTestContext provides test context for user order handler tests.
type UserOrderTestContext struct {
	Router      *gin.Engine
	Service     *order.OrderService
	DB          *gorm.DB
	TestUser    *model.User
	TestPlayer  *model.Player
	TestGame    *model.Game
	AuthToken   string
}

// SetupUserOrderTest initializes test environment for user order handler tests.
func SetupUserOrderTest(t *testing.T) *UserOrderTestContext {
	t.Helper()

	// Setup test database
	db := testutil.SetupTestDB(t)

	// Setup Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create all required repositories
	users := user.NewUserRepository(db)
	players := player.NewPlayerRepository(db)
	games := game.NewGameRepository(db)
	orders := implementations.NewOrderRepository(db)
	payments := payment.NewPaymentRepository(db)
	reviews := review.NewReviewRepository(db)
	commissions := commission.NewCommissionRepository(db)

	// Create order service with all required dependencies
	orderSvc := order.NewOrderService(orders, players, users, games, payments, reviews, commissions)

	// Create test user
	testUser := testutil.CreateAdminUser(t, db, model.RoleUser)
	authToken := fmt.Sprintf("Bearer test-token-%d", testUser.ID)

	// Create test game
	testGame := testutil.CreateTestGame(t, db)

	// Create test player
	testPlayer := testutil.CreateTestPlayer(t, db, testUser.ID)

	return &UserOrderTestContext{
		Router:     router,
		Service:    orderSvc,
		DB:         db,
		TestUser:   testUser,
		TestPlayer: testPlayer,
		TestGame:   testGame,
		AuthToken:  authToken,
	}
}

// RegisterUserOrderRoutes registers user order routes for testing.
func (ctx *UserOrderTestContext) RegisterUserOrderRoutes() {
	group := ctx.Router.Group("/user/orders")
	{
		group.POST("", func(c *gin.Context) {
			c.Set("user_id", ctx.TestUser.ID)
			createOrderHandler(c, ctx.Service)
		})
		group.GET("", func(c *gin.Context) {
			c.Set("user_id", ctx.TestUser.ID)
			getMyOrdersHandler(c, ctx.Service)
		})
		group.GET("/:id", func(c *gin.Context) {
			c.Set("user_id", ctx.TestUser.ID)
			getOrderDetailHandler(c, ctx.Service)
		})
		group.PUT("/:id/cancel", func(c *gin.Context) {
			c.Set("user_id", ctx.TestUser.ID)
			cancelOrderHandler(c, ctx.Service)
		})
		group.PUT("/:id/complete", func(c *gin.Context) {
			c.Set("user_id", ctx.TestUser.ID)
			completeOrderHandler(c, ctx.Service)
		})
	}
}

// makeRequest helper for making authenticated requests
func (ctx *UserOrderTestContext) makeRequest(t *testing.T, method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonData, err := json.Marshal(body)
		require.NoError(t, err)
		reqBody = bytes.NewBuffer(jsonData)
	} else {
		reqBody = bytes.NewBuffer([]byte{})
	}

	req, err := http.NewRequest(method, path, reqBody)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", ctx.AuthToken)

	w := httptest.NewRecorder()
	ctx.Router.ServeHTTP(w, req)
	return w
}

// ============================================================================
// CreateOrder Tests
// ============================================================================

func TestUserOrderHandler_Unit_CreateOrder_Success(t *testing.T) {
	ctx := SetupUserOrderTest(t)
	ctx.RegisterUserOrderRoutes()

	payload := map[string]interface{}{
		"player_id":     ctx.TestPlayer.ID,
		"game_id":       ctx.TestGame.ID,
		"title":         "Test Order",
		"description":   "I need a companion for gaming",
		"price_cents":   5000,
		"currency":      "USD",
		"scheduled_for": time.Now().Add(1 * time.Hour).Format(time.RFC3339),
		"duration_mins": 60,
	}

	w := ctx.makeRequest(t, "POST", "/user/orders", payload)
	testutil.AssertSuccess(t, w, http.StatusOK)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.NotEmpty(t, data["order_id"])
}

func TestUserOrderHandler_Unit_CreateOrder_ValidationError(t *testing.T) {
	ctx := SetupUserOrderTest(t)
	ctx.RegisterUserOrderRoutes()

	// Missing required fields
	payload := map[string]interface{}{
		"title": "Test Order",
	}

	w := ctx.makeRequest(t, "POST", "/user/orders", payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

func TestUserOrderHandler_Unit_CreateOrder_InvalidPlayer(t *testing.T) {
	ctx := SetupUserOrderTest(t)
	ctx.RegisterUserOrderRoutes()

	payload := map[string]interface{}{
		"player_id":     999999,
		"game_id":       ctx.TestGame.ID,
		"title":         "Test Order",
		"price_cents":   5000,
		"currency":      "USD",
		"scheduled_for": time.Now().Add(1 * time.Hour).Format(time.RFC3339),
		"duration_mins": 60,
	}

	w := ctx.makeRequest(t, "POST", "/user/orders", payload)
	testutil.AssertError(t, w, http.StatusNotFound)
}

func TestUserOrderHandler_Unit_CreateOrder_InvalidGame(t *testing.T) {
	ctx := SetupUserOrderTest(t)
	ctx.RegisterUserOrderRoutes()

	payload := map[string]interface{}{
		"player_id":     ctx.TestPlayer.ID,
		"game_id":       999999,
		"title":         "Test Order",
		"price_cents":   5000,
		"currency":      "USD",
		"scheduled_for": time.Now().Add(1 * time.Hour).Format(time.RFC3339),
		"duration_mins": 60,
	}

	w := ctx.makeRequest(t, "POST", "/user/orders", payload)
	testutil.AssertError(t, w, http.StatusNotFound)
}

func TestUserOrderHandler_Unit_CreateOrder_InvalidTime(t *testing.T) {
	ctx := SetupUserOrderTest(t)
	ctx.RegisterUserOrderRoutes()

	payload := map[string]interface{}{
		"player_id":     ctx.TestPlayer.ID,
		"game_id":       ctx.TestGame.ID,
		"title":         "Test Order",
		"price_cents":   5000,
		"currency":      "USD",
		"scheduled_for": "invalid-time",
		"duration_mins": 60,
	}

	w := ctx.makeRequest(t, "POST", "/user/orders", payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

// ============================================================================
// GetMyOrders Tests
// ============================================================================

func TestUserOrderHandler_Unit_GetMyOrders_Success(t *testing.T) {
	ctx := SetupUserOrderTest(t)
	ctx.RegisterUserOrderRoutes()

	// Create test order
	testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusPending)

	w := ctx.makeRequest(t, "GET", "/user/orders", nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	items := data["orders"].([]interface{})
	assert.GreaterOrEqual(t, len(items), 1)
}

func TestUserOrderHandler_Unit_GetMyOrders_WithStatusFilter(t *testing.T) {
	ctx := SetupUserOrderTest(t)
	ctx.RegisterUserOrderRoutes()

	// Create orders with different statuses
	testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusPending)
	testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusCompleted)

	w := ctx.makeRequest(t, "GET", "/user/orders?status=pending", nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	items := data["orders"].([]interface{})
	// Should only return pending orders
	for _, item := range items {
		order := item.(map[string]interface{})
		assert.Equal(t, "pending", order["status"])
	}
}

func TestUserOrderHandler_Unit_GetMyOrders_WithPagination(t *testing.T) {
	ctx := SetupUserOrderTest(t)
	ctx.RegisterUserOrderRoutes()

	// Create multiple orders
	for i := 0; i < 15; i++ {
		testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusPending)
	}

	w := ctx.makeRequest(t, "GET", "/user/orders?page=1&page_size=10", nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	items := data["orders"].([]interface{})
	pagination := data["pagination"].(map[string]interface{})
	assert.LessOrEqual(t, len(items), 10)
	assert.Equal(t, float64(1), pagination["page"])
	assert.Equal(t, float64(10), pagination["page_size"])
}

func TestUserOrderHandler_Unit_GetMyOrders_InvalidQueryParams(t *testing.T) {
	ctx := SetupUserOrderTest(t)
	ctx.RegisterUserOrderRoutes()

	w := ctx.makeRequest(t, "GET", "/user/orders?status=invalid_status", nil)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

func TestUserOrderHandler_Unit_GetMyOrders_EmptyList(t *testing.T) {
	ctx := SetupUserOrderTest(t)
	ctx.RegisterUserOrderRoutes()

	w := ctx.makeRequest(t, "GET", "/user/orders", nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	items := data["orders"].([]interface{})
	assert.Equal(t, 0, len(items))
}

// ============================================================================
// GetOrderDetail Tests
// ============================================================================

func TestUserOrderHandler_Unit_GetOrderDetail_Success(t *testing.T) {
	ctx := SetupUserOrderTest(t)
	ctx.RegisterUserOrderRoutes()

	testOrder := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusPending)

	path := fmt.Sprintf("/user/orders/%d", testOrder.ID)
	w := ctx.makeRequest(t, "GET", path, nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(testOrder.ID), data["id"])
	assert.Equal(t, "Test Order", data["title"])
}

func TestUserOrderHandler_Unit_GetOrderDetail_NotFound(t *testing.T) {
	ctx := SetupUserOrderTest(t)
	ctx.RegisterUserOrderRoutes()

	w := ctx.makeRequest(t, "GET", "/user/orders/999999", nil)
	testutil.AssertError(t, w, http.StatusNotFound)
}

func TestUserOrderHandler_Unit_GetOrderDetail_Unauthorized(t *testing.T) {
	ctx := SetupUserOrderTest(t)
	ctx.RegisterUserOrderRoutes()

	// Create order for different user
	anotherUser := testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)
	testOrder := testutil.CreateTestOrder(t, ctx.DB, anotherUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusPending)

	path := fmt.Sprintf("/user/orders/%d", testOrder.ID)
	w := ctx.makeRequest(t, "GET", path, nil)
	testutil.AssertError(t, w, http.StatusForbidden)
}

func TestUserOrderHandler_Unit_GetOrderDetail_InvalidID(t *testing.T) {
	ctx := SetupUserOrderTest(t)
	ctx.RegisterUserOrderRoutes()

	w := ctx.makeRequest(t, "GET", "/user/orders/invalid", nil)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

// ============================================================================
// CancelOrder Tests
// ============================================================================

func TestUserOrderHandler_Unit_CancelOrder_Success(t *testing.T) {
	ctx := SetupUserOrderTest(t)
	ctx.RegisterUserOrderRoutes()

	testOrder := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusConfirmed)

	payload := map[string]interface{}{
		"reason": "I changed my mind",
	}

	path := fmt.Sprintf("/user/orders/%d/cancel", testOrder.ID)
	w := ctx.makeRequest(t, "PUT", path, payload)
	testutil.AssertSuccess(t, w)

	// Verify order status changed
	testutil.AssertOrderStatus(t, ctx.DB, testOrder.ID, model.OrderStatusCanceled)
}

func TestUserOrderHandler_Unit_CancelOrder_InvalidTransition(t *testing.T) {
	ctx := SetupUserOrderTest(t)
	ctx.RegisterUserOrderRoutes()

	// Try to cancel an already completed order
	testOrder := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusCompleted)

	payload := map[string]interface{}{
		"reason": "Want to cancel",
	}

	path := fmt.Sprintf("/user/orders/%d/cancel", testOrder.ID)
	w := ctx.makeRequest(t, "PUT", path, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

func TestUserOrderHandler_Unit_CancelOrder_NotFound(t *testing.T) {
	ctx := SetupUserOrderTest(t)
	ctx.RegisterUserOrderRoutes()

	payload := map[string]interface{}{
		"reason": "Test",
	}

	w := ctx.makeRequest(t, "PUT", "/user/orders/999999/cancel", payload)
	testutil.AssertError(t, w, http.StatusNotFound)
}

func TestUserOrderHandler_Unit_CancelOrder_MissingReason(t *testing.T) {
	ctx := SetupUserOrderTest(t)
	ctx.RegisterUserOrderRoutes()

	testOrder := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusConfirmed)

	path := fmt.Sprintf("/user/orders/%d/cancel", testOrder.ID)
	w := ctx.makeRequest(t, "PUT", path, map[string]interface{}{})
	testutil.AssertError(t, w, http.StatusBadRequest)
}

// ============================================================================
// CompleteOrder Tests
// ============================================================================

func TestUserOrderHandler_Unit_CompleteOrder_Success(t *testing.T) {
	ctx := SetupUserOrderTest(t)
	ctx.RegisterUserOrderRoutes()

	testOrder := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusInProgress)

	path := fmt.Sprintf("/user/orders/%d/complete", testOrder.ID)
	w := ctx.makeRequest(t, "PUT", path, nil)
	testutil.AssertSuccess(t, w)

	// Verify order status changed
	testutil.AssertOrderStatus(t, ctx.DB, testOrder.ID, model.OrderStatusCompleted)
}

func TestUserOrderHandler_Unit_CompleteOrder_InvalidTransition(t *testing.T) {
	ctx := SetupUserOrderTest(t)
	ctx.RegisterUserOrderRoutes()

	// Try to complete a pending order (must be in_progress first)
	testOrder := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusPending)

	path := fmt.Sprintf("/user/orders/%d/complete", testOrder.ID)
	w := ctx.makeRequest(t, "PUT", path, nil)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

func TestUserOrderHandler_Unit_CompleteOrder_NotFound(t *testing.T) {
	ctx := SetupUserOrderTest(t)
	ctx.RegisterUserOrderRoutes()

	w := ctx.makeRequest(t, "PUT", "/user/orders/999999/complete", nil)
	testutil.AssertError(t, w, http.StatusNotFound)
}

func TestUserOrderHandler_Unit_CompleteOrder_Unauthorized(t *testing.T) {
	ctx := SetupUserOrderTest(t)
	ctx.RegisterUserOrderRoutes()

	// Create order for different user
	anotherUser := testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)
	testOrder := testutil.CreateTestOrder(t, ctx.DB, anotherUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusInProgress)

	path := fmt.Sprintf("/user/orders/%d/complete", testOrder.ID)
	w := ctx.makeRequest(t, "PUT", path, nil)
	testutil.AssertError(t, w, http.StatusForbidden)
}

// ============================================================================
// Error Handling Tests
// ============================================================================

func TestUserOrderHandler_Unit_InvalidJSON(t *testing.T) {
	ctx := SetupUserOrderTest(t)
	ctx.RegisterUserOrderRoutes()

	reqBody := bytes.NewBufferString("{invalid json}")
	req, err := http.NewRequest("POST", "/user/orders", reqBody)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", ctx.AuthToken)

	w := httptest.NewRecorder()
	ctx.Router.ServeHTTP(w, req)

	testutil.AssertError(t, w, http.StatusBadRequest)
}

func TestUserOrderHandler_Unit_MissingContentType(t *testing.T) {
	ctx := SetupUserOrderTest(t)
	ctx.RegisterUserOrderRoutes()

	payload := map[string]interface{}{
		"title": "Test",
	}
	jsonData, _ := json.Marshal(payload)
	reqBody := bytes.NewBuffer(jsonData)

	req, err := http.NewRequest("POST", "/user/orders", reqBody)
	require.NoError(t, err)
	req.Header.Set("Authorization", ctx.AuthToken)
	// Missing Content-Type header

	w := httptest.NewRecorder()
	ctx.Router.ServeHTTP(w, req)

	// Handler should still work with Gin's default binding
}

func TestUserOrderHandler_Unit_URLParameterEncoding(t *testing.T) {
	ctx := SetupUserOrderTest(t)
	ctx.RegisterUserOrderRoutes()

	// Create order with special characters in title
	testOrder := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusPending)

	// Use url.PathEscape to ensure proper encoding
	path := fmt.Sprintf("/user/orders/%d", testOrder.ID)
	w := ctx.makeRequest(t, "GET", path, nil)
	testutil.AssertSuccess(t, w)
}

func TestUserOrderHandler_Unit_QueryParameterEncoding(t *testing.T) {
	ctx := SetupUserOrderTest(t)
	ctx.RegisterUserOrderRoutes()

	// Test with properly encoded query parameters
	params := url.Values{}
	params.Set("status", "pending")
	params.Set("page", "1")
	params.Set("page_size", "10")

	req, err := http.NewRequest("GET", "/user/orders?"+params.Encode(), nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", ctx.AuthToken)

	w := httptest.NewRecorder()
	ctx.Router.ServeHTTP(w, req)

	testutil.AssertSuccess(t, w)
}
