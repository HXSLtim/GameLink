// Package admin provides unit tests for payment handlers.
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

// PaymentTestContext provides test context for payment handler tests.
type PaymentTestContext struct {
	Router     *gin.Engine
	Handler    *PaymentHandler
	Service    *adminservice.AdminService
	DB         *gorm.DB
	AdminUser  *model.User
	AdminToken string
	TestUser   *model.User
	TestGame   *model.Game
	TestOrder  *model.Order
}

// SetupPaymentTest initializes test environment for payment handler tests.
func SetupPaymentTest(t *testing.T) *PaymentTestContext {
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
	svc := adminservice.NewAdminService(adminservice.AdminDeps{
		Games: games, Users: users, Players: players, Orders: orders, Payments: payments,
		Refunds: payment.NewRefundRecordRepository(db),
		Roles:   roles, ServiceItems: serviceItems, Permissions: permissions, Menus: menus,
		Stats: statsRepo, GameCategories: gameCategories, Cache: c,
	})

	// Setup router
	router := testutil.SetupGinTest(t)
	handler := NewPaymentHandler(svc)

	// Create super admin user
	adminUser := testutil.CreateSuperAdmin(t, db)
	adminToken := testutil.GenerateTestToken(adminUser.ID)

	// Create test data
	testUser := testutil.CreateAdminUser(t, db, model.RoleUser)
	testGame := testutil.CreateTestGame(t, db)
	testPlayer := testutil.CreateTestPlayer(t, db, testUser.ID)
	testOrder := testutil.CreateTestOrder(t, db, testUser.ID, testPlayer.ID, testGame.ID, model.OrderStatusPending)

	return &PaymentTestContext{
		Router:     router,
		Handler:    handler,
		Service:    svc,
		DB:         db,
		AdminUser:  adminUser,
		AdminToken: adminToken,
		TestUser:   testUser,
		TestGame:   testGame,
		TestOrder:  testOrder,
	}
}

// RegisterPaymentRoutes registers payment routes for testing.
func (ctx *PaymentTestContext) RegisterPaymentRoutes() {
	group := ctx.Router.Group("/admin/payments")
	{
		group.GET("", ctx.Handler.ListPayments)
		group.POST("", ctx.Handler.CreatePayment)
		group.GET("/:id", ctx.Handler.GetPayment)
		group.PUT("/:id", ctx.Handler.UpdatePayment)
		group.DELETE("/:id", ctx.Handler.DeletePayment)
		group.POST("/:id/capture", ctx.Handler.CapturePayment)
		group.POST("/:id/refund", ctx.Handler.RefundPayment)
		group.GET("/:id/refunds", ctx.Handler.GetRefundHistory)
		group.GET("/:id/logs", ctx.Handler.ListPaymentLogs)
		group.POST("/batch/capture", ctx.Handler.BatchCapturePayments)
		group.POST("/batch/refund", ctx.Handler.BatchRefundPayments)
		group.POST("/batch/cancel", ctx.Handler.BatchCancelPayments)
		group.PUT("/batch/status", ctx.Handler.BatchUpdatePaymentsStatus)
	}
}

// ============================================================================
// CreatePayment Tests
// ============================================================================

func TestPaymentHandler_Unit_CreatePayment_Success(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	payload := map[string]interface{}{
		"order_id":     ctx.TestOrder.ID,
		"user_id":      ctx.TestUser.ID,
		"method":       "alipay",
		"amount_cents": 10000,
		"currency":     "USD",
		"provider_raw": json.RawMessage(`{"txn_id":"test123"}`),
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/payments", ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w, http.StatusCreated)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.NotEmpty(t, data["id"])
	assert.Equal(t, float64(10000), data["amount_cents"])
}

func TestPaymentHandler_Unit_CreatePayment_ValidationError(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	// Missing required fields
	payload := map[string]interface{}{
		"amount_cents": 10000,
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/payments", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

func TestPaymentHandler_Unit_CreatePayment_OrderNotFound(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	payload := map[string]interface{}{
		"order_id":     999999,
		"user_id":      ctx.TestUser.ID,
		"method":       "alipay",
		"amount_cents": 10000,
		"currency":     "USD",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/payments", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusNotFound)
}

// ============================================================================
// ListPayments Tests
// ============================================================================

func TestPaymentHandler_Unit_ListPayments_Success(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	// Create test payments
	for i := 0; i < 5; i++ {
		testutil.CreateTestPayment(t, ctx.DB, ctx.TestOrder.ID, ctx.TestUser.ID, model.PaymentStatusPending)
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/payments", ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	items, pagination := testutil.GetResponseList(t, w)
	assert.GreaterOrEqual(t, len(items), 5)
	assert.Equal(t, float64(1), pagination["page"])
}

func TestPaymentHandler_Unit_ListPayments_WithPagination(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	// Create test payments
	for i := 0; i < 25; i++ {
		testutil.CreateTestPayment(t, ctx.DB, ctx.TestOrder.ID, ctx.TestUser.ID, model.PaymentStatusPending)
	}

	w := testutil.MakeRequest(t, ctx.Router, "GET", "/admin/payments?page=1&page_size=10", nil, testutil.WithAuth(ctx.AdminToken))
	testutil.AssertSuccess(t, w)

	items, pagination := testutil.GetResponseList(t, w)
	assert.LessOrEqual(t, len(items), 10)
	assert.Equal(t, float64(1), pagination["page"])
	assert.Equal(t, float64(10), pagination["page_size"])
}

func TestPaymentHandler_Unit_ListPayments_WithStatusFilter(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	// Create payments with different statuses
	testutil.CreateTestPayment(t, ctx.DB, ctx.TestOrder.ID, ctx.TestUser.ID, model.PaymentStatusPending)
	testutil.CreateTestPayment(t, ctx.DB, ctx.TestOrder.ID, ctx.TestUser.ID, model.PaymentStatusPaid)

	w := testutil.MakeRequest(t, ctx.Router, "GET", "/admin/payments?status=pending", nil, testutil.WithAuth(ctx.AdminToken))
	testutil.AssertSuccess(t, w)

	items, _ := testutil.GetResponseList(t, w)
	for _, item := range items {
		payment := item.(map[string]interface{})
		assert.Equal(t, "pending", payment["status"])
	}
}

func TestPaymentHandler_Unit_ListPayments_WithMethodFilter(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	// Create payments with different methods
	testutil.CreateTestPayment(t, ctx.DB, ctx.TestOrder.ID, ctx.TestUser.ID, model.PaymentStatusPending)

	w := testutil.MakeRequest(t, ctx.Router, "GET", "/admin/payments?method=alipay", nil, testutil.WithAuth(ctx.AdminToken))
	testutil.AssertSuccess(t, w)

	items, _ := testutil.GetResponseList(t, w)
	assert.GreaterOrEqual(t, len(items), 1)
}

func TestPaymentHandler_Unit_ListPayments_WithUserFilter(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	// Create payment for specific user
	testutil.CreateTestPayment(t, ctx.DB, ctx.TestOrder.ID, ctx.TestUser.ID, model.PaymentStatusPending)

	w := testutil.MakeRequest(t, ctx.Router, "GET", fmt.Sprintf("/admin/payments?user_id=%d", ctx.TestUser.ID), nil, testutil.WithAuth(ctx.AdminToken))
	testutil.AssertSuccess(t, w)

	items, _ := testutil.GetResponseList(t, w)
	assert.GreaterOrEqual(t, len(items), 1)
}

func TestPaymentHandler_Unit_ListPayments_WithOrderFilter(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	// Create payment for specific order
	testutil.CreateTestPayment(t, ctx.DB, ctx.TestOrder.ID, ctx.TestUser.ID, model.PaymentStatusPending)

	w := testutil.MakeRequest(t, ctx.Router, "GET", fmt.Sprintf("/admin/payments?order_id=%d", ctx.TestOrder.ID), nil, testutil.WithAuth(ctx.AdminToken))
	testutil.AssertSuccess(t, w)

	items, _ := testutil.GetResponseList(t, w)
	assert.GreaterOrEqual(t, len(items), 1)
}

// ============================================================================
// GetPayment Tests
// ============================================================================

func TestPaymentHandler_Unit_GetPayment_Success(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	testPayment := testutil.CreateTestPayment(t, ctx.DB, ctx.TestOrder.ID, ctx.TestUser.ID, model.PaymentStatusPending)

	path := fmt.Sprintf("/admin/payments/%d", testPayment.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", path, ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(testPayment.ID), data["id"])
}

func TestPaymentHandler_Unit_GetPayment_NotFound(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/payments/999999", ctx.AdminToken, nil)
	testutil.AssertError(t, w, http.StatusNotFound)
}

func TestPaymentHandler_Unit_GetPayment_InvalidID(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/payments/invalid", ctx.AdminToken, nil)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

// ============================================================================
// UpdatePayment Tests
// ============================================================================

func TestPaymentHandler_Unit_UpdatePayment_Success(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	testPayment := testutil.CreateTestPayment(t, ctx.DB, ctx.TestOrder.ID, ctx.TestUser.ID, model.PaymentStatusPending)

	payload := map[string]interface{}{
		"status":            "paid",
		"provider_trade_no": "ALI123456",
		"paid_at":           time.Now().Format(time.RFC3339),
	}

	path := fmt.Sprintf("/admin/payments/%d", testPayment.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", path, ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	// Verify status update
	testutil.AssertPaymentStatus(t, ctx.DB, testPayment.ID, model.PaymentStatusPaid)
}

func TestPaymentHandler_Unit_UpdatePayment_InvalidTime(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	testPayment := testutil.CreateTestPayment(t, ctx.DB, ctx.TestOrder.ID, ctx.TestUser.ID, model.PaymentStatusPending)

	payload := map[string]interface{}{
		"status":  "paid",
		"paid_at": "invalid-time",
	}

	path := fmt.Sprintf("/admin/payments/%d", testPayment.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", path, ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

func TestPaymentHandler_Unit_UpdatePayment_NotFound(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	payload := map[string]interface{}{
		"status": "paid",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", "/admin/payments/999999", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusNotFound)
}

// ============================================================================
// DeletePayment Tests
// ============================================================================

func TestPaymentHandler_Unit_DeletePayment_Success(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	testPayment := testutil.CreateTestPayment(t, ctx.DB, ctx.TestOrder.ID, ctx.TestUser.ID, model.PaymentStatusPending)

	path := fmt.Sprintf("/admin/payments/%d", testPayment.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "DELETE", path, ctx.AdminToken, nil)
	testutil.AssertDeleted(t, w)

	// Verify deletion
	var count int64
	ctx.DB.Model(&model.Payment{}).Where("id = ?", testPayment.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestPaymentHandler_Unit_DeletePayment_NotFound(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "DELETE", "/admin/payments/999999", ctx.AdminToken, nil)
	testutil.AssertError(t, w, http.StatusNotFound)
}

// ============================================================================
// CapturePayment Tests
// ============================================================================

func TestPaymentHandler_Unit_CapturePayment_Success(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	testPayment := testutil.CreateTestPayment(t, ctx.DB, ctx.TestOrder.ID, ctx.TestUser.ID, model.PaymentStatusPending)

	payload := map[string]interface{}{
		"provider_trade_no": "ALI123456",
		"provider_raw":      json.RawMessage(`{"status":"success"}`),
		"paid_at":           time.Now().Format(time.RFC3339),
	}

	path := fmt.Sprintf("/admin/payments/%d/capture", testPayment.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", path, ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	// Verify payment status changed
	testutil.AssertPaymentStatus(t, ctx.DB, testPayment.ID, model.PaymentStatusPaid)
}

func TestPaymentHandler_Unit_CapturePayment_InvalidTime(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	testPayment := testutil.CreateTestPayment(t, ctx.DB, ctx.TestOrder.ID, ctx.TestUser.ID, model.PaymentStatusPending)

	payload := map[string]interface{}{
		"paid_at": "invalid-time",
	}

	path := fmt.Sprintf("/admin/payments/%d/capture", testPayment.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", path, ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

func TestPaymentHandler_Unit_CapturePayment_NotFound(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	payload := map[string]interface{}{
		"provider_trade_no": "ALI123456",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/payments/999999/capture", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusNotFound)
}

// ============================================================================
// RefundPayment Tests
// ============================================================================

func TestPaymentHandler_Unit_RefundPayment_Success(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	testPayment := testutil.CreateTestPayment(t, ctx.DB, ctx.TestOrder.ID, ctx.TestUser.ID, model.PaymentStatusPaid)

	payload := map[string]interface{}{
		"amount_cents": int64(5000),
		"reason":       "Customer requested refund",
		"note":         "Partial refund approved",
	}

	path := fmt.Sprintf("/admin/payments/%d/refund", testPayment.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", path, ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	// Verify refund recorded
	var updatedPayment model.Payment
	ctx.DB.First(&updatedPayment, testPayment.ID)
	assert.Equal(t, int64(5000), updatedPayment.RefundedAmountCents)
}

func TestPaymentHandler_Unit_RefundPayment_FullRefund(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	testPayment := testutil.CreateTestPayment(t, ctx.DB, ctx.TestOrder.ID, ctx.TestUser.ID, model.PaymentStatusPaid)

	payload := map[string]interface{}{
		"amount_cents": testPayment.AmountCents,
		"reason":       "Full refund",
	}

	path := fmt.Sprintf("/admin/payments/%d/refund", testPayment.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", path, ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	// Verify payment marked as refunded
	var updatedPayment model.Payment
	ctx.DB.First(&updatedPayment, testPayment.ID)
	assert.Equal(t, model.PaymentStatusRefunded, updatedPayment.Status)
}

func TestPaymentHandler_Unit_RefundPayment_CumulativePartial(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	// 退款语义在已支付订单上更清晰
	require.NoError(t, ctx.DB.Model(&model.Order{}).
		Where("id = ?", ctx.TestOrder.ID).
		Update("status", model.OrderStatusCompleted).Error)

	testPayment := testutil.CreateTestPayment(t, ctx.DB, ctx.TestOrder.ID, ctx.TestUser.ID, model.PaymentStatusPaid)

	firstPayload := map[string]interface{}{
		"amount_cents": int64(3000),
		"reason":       "First partial refund",
	}
	firstPath := fmt.Sprintf("/admin/payments/%d/refund", testPayment.ID)
	firstResp := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", firstPath, ctx.AdminToken, firstPayload)
	testutil.AssertSuccess(t, firstResp)

	secondPayload := map[string]interface{}{
		"amount_cents": int64(2000),
		"reason":       "Second partial refund",
	}
	secondPath := fmt.Sprintf("/admin/payments/%d/refund", testPayment.ID)
	secondResp := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", secondPath, ctx.AdminToken, secondPayload)
	testutil.AssertSuccess(t, secondResp)

	var updatedPayment model.Payment
	require.NoError(t, ctx.DB.First(&updatedPayment, testPayment.ID).Error)
	assert.Equal(t, int64(5000), updatedPayment.RefundedAmountCents)
	assert.Equal(t, model.PaymentStatusPaid, updatedPayment.Status)

	var updatedOrder model.Order
	require.NoError(t, ctx.DB.First(&updatedOrder, ctx.TestOrder.ID).Error)
	assert.Equal(t, int64(5000), updatedOrder.RefundAmountCents)
	assert.Equal(t, model.OrderStatusCompleted, updatedOrder.Status)
}

func TestPaymentHandler_Unit_RefundPayment_InvalidAmount(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	testPayment := testutil.CreateTestPayment(t, ctx.DB, ctx.TestOrder.ID, ctx.TestUser.ID, model.PaymentStatusPaid)

	payload := map[string]interface{}{
		"amount_cents": int64(999999), // More than payment amount
		"reason":       "Excessive refund",
	}

	path := fmt.Sprintf("/admin/payments/%d/refund", testPayment.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", path, ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

func TestPaymentHandler_Unit_RefundPayment_NotFound(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	payload := map[string]interface{}{
		"amount_cents": int64(5000),
		"reason":       "Test refund",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/payments/999999/refund", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusNotFound)
}

// ============================================================================
// GetRefundHistory Tests
// ============================================================================

func TestPaymentHandler_Unit_GetRefundHistory_Success(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	testPayment := testutil.CreateTestPayment(t, ctx.DB, ctx.TestOrder.ID, ctx.TestUser.ID, model.PaymentStatusPaid)

	// Create one refund first, then verify history returns it
	refundPayload := map[string]interface{}{
		"amount_cents": int64(1000),
		"reason":       "Test refund history",
	}
	refundPath := fmt.Sprintf("/admin/payments/%d/refund", testPayment.ID)
	refundResp := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", refundPath, ctx.AdminToken, refundPayload)
	testutil.AssertSuccess(t, refundResp)

	path := fmt.Sprintf("/admin/payments/%d/refunds", testPayment.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", path, ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].([]interface{})
	assert.GreaterOrEqual(t, len(data), 1)
}

func TestPaymentHandler_Unit_GetRefundHistory_NotFound(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/payments/999999/refunds", ctx.AdminToken, nil)
	testutil.AssertError(t, w, http.StatusNotFound)
}

// ============================================================================
// ListPaymentLogs Tests
// ============================================================================

func TestPaymentHandler_Unit_ListPaymentLogs_Success(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	testPayment := testutil.CreateTestPayment(t, ctx.DB, ctx.TestOrder.ID, ctx.TestUser.ID, model.PaymentStatusPaid)

	path := fmt.Sprintf("/admin/payments/%d/logs", testPayment.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", path, ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	items, pagination := testutil.GetResponseList(t, w)
	assert.NotNil(t, items)
	assert.NotNil(t, pagination)
}

func TestPaymentHandler_Unit_ListPaymentLogs_WithFilters(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	testPayment := testutil.CreateTestPayment(t, ctx.DB, ctx.TestOrder.ID, ctx.TestUser.ID, model.PaymentStatusPaid)

	w := testutil.MakeRequest(t, ctx.Router, "GET", fmt.Sprintf("/admin/payments/%d/logs?action=create", testPayment.ID), nil, testutil.WithAuth(ctx.AdminToken))
	testutil.AssertSuccess(t, w)

	items, _ := testutil.GetResponseList(t, w)
	assert.NotNil(t, items)
}

func TestPaymentHandler_Unit_ListPaymentLogs_NotFound(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/payments/999999/logs", ctx.AdminToken, nil)
	testutil.AssertError(t, w, http.StatusNotFound)
}

func TestPaymentHandler_Unit_BatchCapturePayments_Success(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	p1 := testutil.CreateTestPayment(t, ctx.DB, ctx.TestOrder.ID, ctx.TestUser.ID, model.PaymentStatusPending)
	p2 := testutil.CreateTestPayment(t, ctx.DB, ctx.TestOrder.ID, ctx.TestUser.ID, model.PaymentStatusPending)

	payload := map[string]interface{}{
		"paymentIds": []uint64{p1.ID, p2.ID},
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/payments/batch/capture", ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["successCount"])
	assert.Equal(t, float64(0), data["failedCount"])

	testutil.AssertPaymentStatus(t, ctx.DB, p1.ID, model.PaymentStatusPaid)
	testutil.AssertPaymentStatus(t, ctx.DB, p2.ID, model.PaymentStatusPaid)
}

func TestPaymentHandler_Unit_BatchRefundPayments_Success(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	p1 := testutil.CreateTestPayment(t, ctx.DB, ctx.TestOrder.ID, ctx.TestUser.ID, model.PaymentStatusPaid)
	p2 := testutil.CreateTestPayment(t, ctx.DB, ctx.TestOrder.ID, ctx.TestUser.ID, model.PaymentStatusPaid)

	payload := map[string]interface{}{
		"paymentIds": []uint64{p1.ID, p2.ID},
		"reason":     "batch refund",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/payments/batch/refund", ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["successCount"])
	assert.Equal(t, float64(0), data["failedCount"])

	testutil.AssertPaymentStatus(t, ctx.DB, p1.ID, model.PaymentStatusRefunded)
	testutil.AssertPaymentStatus(t, ctx.DB, p2.ID, model.PaymentStatusRefunded)
}

func TestPaymentHandler_Unit_BatchCancelPayments_Success(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	p1 := testutil.CreateTestPayment(t, ctx.DB, ctx.TestOrder.ID, ctx.TestUser.ID, model.PaymentStatusPending)
	p2 := testutil.CreateTestPayment(t, ctx.DB, ctx.TestOrder.ID, ctx.TestUser.ID, model.PaymentStatusPending)

	payload := map[string]interface{}{
		"paymentIds": []uint64{p1.ID, p2.ID},
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/payments/batch/cancel", ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["successCount"])
	assert.Equal(t, float64(0), data["failedCount"])

	testutil.AssertPaymentStatus(t, ctx.DB, p1.ID, model.PaymentStatusFailed)
	testutil.AssertPaymentStatus(t, ctx.DB, p2.ID, model.PaymentStatusFailed)
}

func TestPaymentHandler_Unit_BatchUpdatePaymentsStatus_Success(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	p1 := testutil.CreateTestPayment(t, ctx.DB, ctx.TestOrder.ID, ctx.TestUser.ID, model.PaymentStatusPending)
	p2 := testutil.CreateTestPayment(t, ctx.DB, ctx.TestOrder.ID, ctx.TestUser.ID, model.PaymentStatusPending)

	payload := map[string]interface{}{
		"paymentIds": []uint64{p1.ID, p2.ID},
		"status":     "paid",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", "/admin/payments/batch/status", ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["successCount"])
	assert.Equal(t, float64(0), data["failedCount"])

	testutil.AssertPaymentStatus(t, ctx.DB, p1.ID, model.PaymentStatusPaid)
	testutil.AssertPaymentStatus(t, ctx.DB, p2.ID, model.PaymentStatusPaid)
}

func TestPaymentHandler_Unit_BatchUpdatePaymentsStatus_InvalidStatus(t *testing.T) {
	ctx := SetupPaymentTest(t)
	ctx.RegisterPaymentRoutes()

	p1 := testutil.CreateTestPayment(t, ctx.DB, ctx.TestOrder.ID, ctx.TestUser.ID, model.PaymentStatusPending)

	payload := map[string]interface{}{
		"paymentIds": []uint64{p1.ID},
		"status":     "invalid",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", "/admin/payments/batch/status", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}
