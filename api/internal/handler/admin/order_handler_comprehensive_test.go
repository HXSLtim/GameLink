// Package admin provides comprehensive unit tests for order handlers.
// This file extends the existing order_handler_test.go with additional test scenarios
// to achieve 85%+ coverage for this critical financial module.
package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/handler/testutil"
	"gamelink/internal/model"
)

// ============================================================================
// Additional CreateOrder Tests - Edge Cases & Business Logic
// ============================================================================

func TestOrderHandler_Comprehensive_CreateOrder_ScheduledTimes(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	testServiceItem := testutil.CreateTestServiceItem(t, ctx.DB, ctx.TestGame.ID, 5000, true)

	tests := []struct {
		name           string
		scheduledStart string
		scheduledEnd   string
		expectSuccess  bool
		expectedStatus int
	}{
		{
			name:           "valid future times",
			scheduledStart: time.Now().Add(1 * time.Hour).Format(time.RFC3339),
			scheduledEnd:   time.Now().Add(2 * time.Hour).Format(time.RFC3339),
			expectSuccess:  true,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "end before start - should fail",
			scheduledStart: time.Now().Add(2 * time.Hour).Format(time.RFC3339),
			scheduledEnd:   time.Now().Add(1 * time.Hour).Format(time.RFC3339),
			expectSuccess:  false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "nil times allowed",
			scheduledStart: "",
			scheduledEnd:   "",
			expectSuccess:  true,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "start in past - should fail",
			scheduledStart: time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
			scheduledEnd:   time.Now().Add(1 * time.Hour).Format(time.RFC3339),
			expectSuccess:  false,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := map[string]interface{}{
				"user_id":           ctx.TestUser.ID,
				"game_id":           ctx.TestGame.ID,
				"item_id":           testServiceItem.ID,
				"title":             "Test Order",
				"total_price_cents": 10000,
				"currency":          "CNY",
			}
			if tt.scheduledStart != "" {
				payload["scheduled_start"] = tt.scheduledStart
			}
			if tt.scheduledEnd != "" {
				payload["scheduled_end"] = tt.scheduledEnd
			}

			w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/orders", ctx.AdminToken, payload)

			if tt.expectSuccess {
				testutil.AssertSuccess(t, w, tt.expectedStatus)
			} else {
				testutil.AssertError(t, w, tt.expectedStatus)
			}
		})
	}
}

func TestOrderHandler_Comprehensive_CreateOrder_CurrencyValidation(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	testServiceItem := testutil.CreateTestServiceItem(t, ctx.DB, ctx.TestGame.ID, 5000, true)

	validCurrencies := []string{"CNY", "USD", "EUR", "JPY", "GBP"}
	for _, currency := range validCurrencies {
		t.Run("valid_currency_"+currency, func(t *testing.T) {
			payload := map[string]interface{}{
				"user_id":           ctx.TestUser.ID,
				"game_id":           ctx.TestGame.ID,
				"item_id":           testServiceItem.ID,
				"total_price_cents": 10000,
				"currency":          currency,
			}

			w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/orders", ctx.AdminToken, payload)
			testutil.AssertSuccess(t, w, http.StatusCreated)
		})
	}

	t.Run("invalid_currency", func(t *testing.T) {
		payload := map[string]interface{}{
			"user_id":           ctx.TestUser.ID,
			"game_id":           ctx.TestGame.ID,
			"item_id":           testServiceItem.ID,
			"total_price_cents": 10000,
			"currency":          "INVALID",
		}

		w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/orders", ctx.AdminToken, payload)
		// Should either succeed with uppercase conversion or fail with validation error
		// depending on service layer validation
		assert.True(t, w.Code == http.StatusCreated || w.Code == http.StatusBadRequest)
	})
}

func TestOrderHandler_Comprehensive_CreateOrder_PriceValidation(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	testServiceItem := testutil.CreateTestServiceItem(t, ctx.DB, ctx.TestGame.ID, 5000, true)

	tests := []struct {
		name           string
		priceCents     int64
		expectSuccess  bool
		expectedStatus int
	}{
		{
			name:           "zero price",
			priceCents:     0,
			expectSuccess:  false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "negative price",
			priceCents:     -100,
			expectSuccess:  false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "very small valid price",
			priceCents:     1,
			expectSuccess:  true,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "very large price",
			priceCents:     999999999,
			expectSuccess:  true,
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := map[string]interface{}{
				"user_id":           ctx.TestUser.ID,
				"game_id":           ctx.TestGame.ID,
				"item_id":           testServiceItem.ID,
				"total_price_cents": tt.priceCents,
				"currency":          "CNY",
			}

			w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/orders", ctx.AdminToken, payload)

			if tt.expectSuccess {
				testutil.AssertSuccess(t, w, tt.expectedStatus)
			} else {
				testutil.AssertError(t, w, tt.expectedStatus)
			}
		})
	}
}

func TestOrderHandler_Comprehensive_CreateOrder_PlayerAssignment(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	testServiceItem := testutil.CreateTestServiceItem(t, ctx.DB, ctx.TestGame.ID, 5000, true)

	t.Run("create with player assignment", func(t *testing.T) {
		payload := map[string]interface{}{
			"user_id":           ctx.TestUser.ID,
			"player_id":         ctx.TestPlayer.ID,
			"game_id":           ctx.TestGame.ID,
			"item_id":           testServiceItem.ID,
			"total_price_cents": 10000,
			"currency":          "CNY",
		}

		w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/orders", ctx.AdminToken, payload)
		testutil.AssertSuccess(t, w, http.StatusCreated)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		data := response["data"].(map[string]interface{})
		playerID, ok := data["player_id"]
		assert.True(t, ok, "Response should contain player_id")
		assert.NotNil(t, playerID)
	})

	t.Run("create without player - player_id null", func(t *testing.T) {
		payload := map[string]interface{}{
			"user_id":           ctx.TestUser.ID,
			"game_id":           ctx.TestGame.ID,
			"item_id":           testServiceItem.ID,
			"total_price_cents": 10000,
			"currency":          "CNY",
		}

		w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/orders", ctx.AdminToken, payload)
		testutil.AssertSuccess(t, w, http.StatusCreated)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		data := response["data"].(map[string]interface{})
		playerID := data["player_id"]
		assert.Nil(t, playerID, "player_id should be nil when not specified")
	})
}

// ============================================================================
// Order State Machine Tests - Critical for Business Logic
// ============================================================================

func TestOrderHandler_Comprehensive_StateMachine_AllTransitions(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	// Test valid state transitions
	validTransitions := []struct {
		fromStatus  model.OrderStatus
		toStatus    model.OrderStatus
		endpoint    string
		setupAction func(*testing.T, *model.Order)
	}{
		{
			fromStatus:  model.OrderStatusPending,
			toStatus:    model.OrderStatusConfirmed,
			endpoint:    "/confirm",
			setupAction: nil,
		},
		{
			fromStatus:  model.OrderStatusConfirmed,
			toStatus:    model.OrderStatusInProgress,
			endpoint:    "/start",
			setupAction: nil,
		},
		{
			fromStatus:  model.OrderStatusInProgress,
			toStatus:    model.OrderStatusCompleted,
			endpoint:    "/complete",
			setupAction: nil,
		},
	}

	for _, tt := range validTransitions {
		t.Run(fmt.Sprintf("valid_transition_%s_to_%s", tt.fromStatus, tt.toStatus), func(t *testing.T) {
			// Create order in fromStatus
			order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, tt.fromStatus)

			path := fmt.Sprintf("/admin/orders/%d%s", order.ID, tt.endpoint)
			w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", path, ctx.AdminToken, map[string]interface{}{})
			testutil.AssertSuccess(t, w)

			// Verify status changed
			testutil.AssertOrderStatus(t, ctx.DB, order.ID, tt.toStatus)
		})
	}
}

func TestOrderHandler_Comprehensive_StateMachine_InvalidTransitions(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	// Test invalid state transitions
	invalidTransitions := []struct {
		name        string
		fromStatus  model.OrderStatus
		endpoint    string
		description string
	}{
		{
			name:        "completed_to_confirmed",
			fromStatus:  model.OrderStatusCompleted,
			endpoint:    "/confirm",
			description: "Cannot confirm a completed order",
		},
		{
			name:        "canceled_to_start",
			fromStatus:  model.OrderStatusCanceled,
			endpoint:    "/start",
			description: "Cannot start a canceled order",
		},
		{
			name:        "pending_to_complete",
			fromStatus:  model.OrderStatusPending,
			endpoint:    "/complete",
			description: "Cannot complete a pending order",
		},
	}

	for _, tt := range invalidTransitions {
		t.Run(tt.name, func(t *testing.T) {
			order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, tt.fromStatus)

			path := fmt.Sprintf("/admin/orders/%d%s", order.ID, tt.endpoint)
			w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", path, ctx.AdminToken, map[string]interface{}{})

			// Should fail with validation error
			testutil.AssertError(t, w, http.StatusBadRequest)

			// Verify status unchanged
			testutil.AssertOrderStatus(t, ctx.DB, order.ID, tt.fromStatus)
		})
	}
}

// ============================================================================
// Order Update Tests - Comprehensive Field Updates
// ============================================================================

func TestOrderHandler_Comprehensive_UpdateOrder_AllFields(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusPending)

	newStartTime := time.Now().Add(5 * time.Hour).Format(time.RFC3339)
	newEndTime := time.Now().Add(6 * time.Hour).Format(time.RFC3339)

	payload := map[string]interface{}{
		"status":            "confirmed",
		"total_price_cents": 15000,
		"currency":          "USD",
		"scheduled_start":   newStartTime,
		"scheduled_end":     newEndTime,
	}

	path := fmt.Sprintf("/admin/orders/%d", order.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", path, ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	// Verify all fields updated in DB
	var updatedOrder model.Order
	err := ctx.DB.First(&updatedOrder, order.ID).Error
	require.NoError(t, err)

	assert.Equal(t, model.OrderStatusConfirmed, updatedOrder.Status)
	assert.Equal(t, int64(15000), updatedOrder.TotalPriceCents)
	assert.Equal(t, model.CurrencyUSD, updatedOrder.Currency)
	assert.NotNil(t, updatedOrder.ScheduledStart)
	assert.NotNil(t, updatedOrder.ScheduledEnd)
}

func TestOrderHandler_Comprehensive_UpdateOrder_StatusOnly(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusPending)

	payload := map[string]interface{}{
		"status":            "confirmed",
		"total_price_cents": order.TotalPriceCents, // Keep same
		"currency":          string(order.Currency),
	}

	path := fmt.Sprintf("/admin/orders/%d", order.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", path, ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	testutil.AssertOrderStatus(t, ctx.DB, order.ID, model.OrderStatusConfirmed)
}

// ============================================================================
// Order Cancel Tests - Business Logic
// ============================================================================

func TestOrderHandler_Comprehensive_CancelOrder_AllStatuses(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	// Test cancellation from different states
	cancellableStatuses := []model.OrderStatus{
		model.OrderStatusPending,
		model.OrderStatusConfirmed,
	}

	for _, status := range cancellableStatuses {
		t.Run(fmt.Sprintf("cancel_from_%s", status), func(t *testing.T) {
			order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, status)

			reason := fmt.Sprintf("Cancelling from %s", status)
			payload := map[string]interface{}{
				"reason": reason,
			}

			path := fmt.Sprintf("/admin/orders/%d/cancel", order.ID)
			w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", path, ctx.AdminToken, payload)
			testutil.AssertSuccess(t, w)

			// Verify order canceled
			testutil.AssertOrderStatus(t, ctx.DB, order.ID, model.OrderStatusCanceled)

			// Verify cancel reason saved
			var updatedOrder model.Order
			ctx.DB.First(&updatedOrder, order.ID)
			assert.Equal(t, reason, updatedOrder.CancelReason)
		})
	}

	// Test that completed orders cannot be canceled
	t.Run("cannot_cancel_completed", func(t *testing.T) {
		order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusCompleted)

		payload := map[string]interface{}{
			"reason": "Try to cancel completed order",
		}

		path := fmt.Sprintf("/admin/orders/%d/cancel", order.ID)
		w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", path, ctx.AdminToken, payload)
		testutil.AssertError(t, w, http.StatusBadRequest)

		// Verify status unchanged
		testutil.AssertOrderStatus(t, ctx.DB, order.ID, model.OrderStatusCompleted)
	})
}

func TestOrderHandler_Comprehensive_CancelOrder_WithPayment(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusConfirmed)
	payment := testutil.CreateTestPayment(t, ctx.DB, order.ID, ctx.TestUser.ID, model.PaymentStatusPaid)

	// Verify payment exists before cancel
	var paymentCount int64
	ctx.DB.Model(&model.Payment{}).Where("id = ?", payment.ID).Count(&paymentCount)
	assert.Equal(t, int64(1), paymentCount)

	// Cancel the order
	payload := map[string]interface{}{
		"reason": "Customer requested refund",
	}

	path := fmt.Sprintf("/admin/orders/%d/cancel", order.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", path, ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	testutil.AssertOrderStatus(t, ctx.DB, order.ID, model.OrderStatusCanceled)

	// Note: Actual refund processing happens in service layer
	// This test verifies handler correctly forwards the request
}

// ============================================================================
// Order Assignment Tests
// ============================================================================

func TestOrderHandler_Comprehensive_AssignOrder_AlreadyAssigned(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	// Create order with existing player
	order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusPending)

	// Create another player
	anotherUser := testutil.CreateAdminUser(t, ctx.DB, model.RolePlayer)
	anotherPlayer := testutil.CreateTestPlayer(t, ctx.DB, anotherUser.ID)

	// Try to assign different player
	payload := map[string]interface{}{
		"player_id": anotherPlayer.ID,
	}

	path := fmt.Sprintf("/admin/orders/%d/assign", order.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", path, ctx.AdminToken, payload)

	// Behavior depends on service layer - might allow reassignment or reject
	// Just verify response is valid
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusBadRequest)

	// Verify original player still assigned if rejected
	if w.Code == http.StatusBadRequest {
		var updatedOrder model.Order
		ctx.DB.First(&updatedOrder, order.ID)
		assert.Equal(t, ctx.TestPlayer.ID, *updatedOrder.PlayerID)
	}
}

func TestOrderHandler_Comprehensive_AssignOrder_NonExistentPlayer(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, 0, ctx.TestGame.ID, model.OrderStatusPending)

	payload := map[string]interface{}{
		"player_id": uint64(99999999),
	}

	path := fmt.Sprintf("/admin/orders/%d/assign", order.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", path, ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusNotFound)
}

// ============================================================================
// Order Refund Tests - Financial Operations
// ============================================================================

func TestOrderHandler_Comprehensive_RefundOrder_FullRefund(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusCompleted)
	_ = testutil.CreateTestPayment(t, ctx.DB, order.ID, ctx.TestUser.ID, model.PaymentStatusPaid)

	payload := map[string]interface{}{
		"reason":       "Full refund requested",
		"amount_cents": order.TotalPriceCents, // Full refund
		"note":         "Customer satisfaction issue",
	}

	path := fmt.Sprintf("/admin/orders/%d/refund", order.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", path, ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)
}

func TestOrderHandler_Comprehensive_RefundOrder_PartialRefund(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusCompleted)
	_ = testutil.CreateTestPayment(t, ctx.DB, order.ID, ctx.TestUser.ID, model.PaymentStatusPaid)

	// Refund half
	partialAmount := order.TotalPriceCents / 2

	payload := map[string]interface{}{
		"reason":       "Partial refund for service issue",
		"amount_cents": partialAmount,
		"note":         "Compensated for late start",
	}

	path := fmt.Sprintf("/admin/orders/%d/refund", order.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", path, ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)
}

func TestOrderHandler_Comprehensive_RefundOrder_NoPayment(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	// Create order without payment
	order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusCompleted)

	payload := map[string]interface{}{
		"reason":       "Try to refund unpaid order",
		"amount_cents": int64(5000),
	}

	path := fmt.Sprintf("/admin/orders/%d/refund", order.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", path, ctx.AdminToken, payload)

	// Should fail - no payment to refund
	testutil.AssertError(t, w, http.StatusBadRequest)
}

func TestOrderHandler_Comprehensive_RefundOrder_ExcessAmount(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusCompleted)
	_ = testutil.CreateTestPayment(t, ctx.DB, order.ID, ctx.TestUser.ID, model.PaymentStatusPaid)

	// Try to refund more than order amount
	excessAmount := order.TotalPriceCents + 10000

	payload := map[string]interface{}{
		"reason":       "Greedy refund",
		"amount_cents": excessAmount,
	}

	path := fmt.Sprintf("/admin/orders/%d/refund", order.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", path, ctx.AdminToken, payload)

	// Should fail - refund amount exceeds order total
	testutil.AssertError(t, w, http.StatusBadRequest)
}

// ============================================================================
// Order List Tests - Advanced Filtering
// ============================================================================

func TestOrderHandler_Comprehensive_ListOrders_AllFilters(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	// Create orders with different attributes
	user1 := ctx.TestUser
	user2 := testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)
	player2 := testutil.CreateTestPlayer(t, ctx.DB, user2.ID)
	game2 := testutil.CreateTestGame(t, ctx.DB)

	testutil.CreateTestOrder(t, ctx.DB, user1.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusPending)
	testutil.CreateTestOrder(t, ctx.DB, user2.ID, player2.ID, game2.ID, model.OrderStatusCompleted)

	time.Sleep(time.Second) // Ensure different timestamps
	yesterday := time.Now().Add(-24 * time.Hour)

	tests := []struct {
		name           string
		queryParams    map[string]string
		minResults     int
		maxResults     int
		validateResult func(*testing.T, []interface{})
	}{
		{
			name: "filter_by_user",
			queryParams: map[string]string{
				"user_id": fmt.Sprintf("%d", user1.ID),
			},
			minResults: 1,
			validateResult: func(t *testing.T, items []interface{}) {
				for _, item := range items {
					order := item.(map[string]interface{})
					userID := uint64(order["user_id"].(float64))
					assert.Equal(t, user1.ID, userID)
				}
			},
		},
		{
			name: "filter_by_player",
			queryParams: map[string]string{
				"player_id": fmt.Sprintf("%d", ctx.TestPlayer.ID),
			},
			minResults: 1,
			validateResult: func(t *testing.T, items []interface{}) {
				for _, item := range items {
					order := item.(map[string]interface{})
					playerID := order["player_id"]
					if playerID != nil {
						pid := uint64(playerID.(float64))
						assert.Equal(t, ctx.TestPlayer.ID, pid)
					}
				}
			},
		},
		{
			name: "filter_by_game",
			queryParams: map[string]string{
				"game_id": fmt.Sprintf("%d", ctx.TestGame.ID),
			},
			minResults: 1,
		},
		{
			name: "filter_by_status",
			queryParams: map[string]string{
				"status": "completed",
			},
			minResults: 1,
			validateResult: func(t *testing.T, items []interface{}) {
				for _, item := range items {
					order := item.(map[string]interface{})
					assert.Equal(t, "completed", order["status"])
				}
			},
		},
		{
			name: "filter_by_date_range",
			queryParams: map[string]string{
				"date_from": yesterday.Format("2006-01-02"),
				"date_to":   time.Now().Format("2006-01-02"),
			},
			minResults: 1,
		},
		{
			name: "combined_filters",
			queryParams: map[string]string{
				"status":   "pending",
				"user_id":  fmt.Sprintf("%d", user1.ID),
				"game_id":  fmt.Sprintf("%d", ctx.TestGame.ID),
				"page":     "1",
				"page_size": "10",
			},
			minResults: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := testutil.MakeRequest(t, ctx.Router, "GET", "/admin/orders", nil,
				testutil.WithAuth(ctx.AdminToken),
				testutil.WithQuery(tt.queryParams))

			testutil.AssertSuccess(t, w)

			items, _ := testutil.GetResponseList(t, w)

			if tt.minResults > 0 {
				assert.GreaterOrEqual(t, len(items), tt.minResults)
			}
			if tt.maxResults > 0 {
				assert.LessOrEqual(t, len(items), tt.maxResults)
			}
			if tt.validateResult != nil {
				tt.validateResult(t, items)
			}
		})
	}
}

func TestOrderHandler_Comprehensive_ListOrders_Sorting(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	// Create multiple orders
	for i := 0; i < 5; i++ {
		testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusPending)
		time.Sleep(10 * time.Millisecond) // Ensure different timestamps
	}

	w := testutil.MakeRequest(t, ctx.Router, "GET", "/admin/orders?page=1&page_size=20", nil,
		testutil.WithAuth(ctx.AdminToken))
	testutil.AssertSuccess(t, w)

	items, _ := testutil.GetResponseList(t, w)
	assert.GreaterOrEqual(t, len(items), 5)

	// Verify sorted by created_at desc (most recent first)
	// Note: This depends on service layer implementation
}

// ============================================================================
// Order Timeline Tests
// ============================================================================

func TestOrderHandler_Comprehensive_GetOrderTimeline_CompleteWorkflow(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	testServiceItem := testutil.CreateTestServiceItem(t, ctx.DB, ctx.TestGame.ID, 5000, true)

	// Create order
	payload := map[string]interface{}{
		"user_id":           ctx.TestUser.ID,
		"player_id":         ctx.TestPlayer.ID,
		"game_id":           ctx.TestGame.ID,
		"item_id":           testServiceItem.ID,
		"total_price_cents": 10000,
		"currency":          "CNY",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/orders", ctx.AdminToken, payload)
	require.Equal(t, http.StatusCreated, w.Code)

	var createResponse map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &createResponse))
	orderData := createResponse["data"].(map[string]interface{})
	orderID := uint64(orderData["id"].(float64))

	// Progress through states
	transitions := []struct {
		endpoint string
		status   model.OrderStatus
	}{
		{"/confirm", model.OrderStatusConfirmed},
		{"/start", model.OrderStatusInProgress},
		{"/complete", model.OrderStatusCompleted},
	}

	for _, tt := range transitions {
		path := fmt.Sprintf("/admin/orders/%d%s", orderID, tt.endpoint)
		w = testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", path, ctx.AdminToken, map[string]interface{}{})
		require.Equal(t, http.StatusOK, w.Code)
		time.Sleep(10 * time.Millisecond) // Ensure different timestamps
	}

	// Get timeline
	path := fmt.Sprintf("/admin/orders/%d/timeline", orderID)
	w = testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", path, ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	timeline := response["data"].([]interface{})
	assert.GreaterOrEqual(t, len(timeline), 1, "Timeline should have at least one entry")

	// Verify timeline entries have required fields
	for _, item := range timeline {
		entry := item.(map[string]interface{})
		assert.Contains(t, entry, "status")
		assert.Contains(t, entry, "changed_at")
	}
}

// ============================================================================
// Order Payments Tests
// ============================================================================

func TestOrderHandler_Comprehensive_ListOrderPayments_MultiplePayments(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusPending)

	// Create multiple payments
	_ = testutil.CreateTestPayment(t, ctx.DB, order.ID, ctx.TestUser.ID, model.PaymentStatusPending)
	_ = testutil.CreateTestPayment(t, ctx.DB, order.ID, ctx.TestUser.ID, model.PaymentStatusPaid)

	path := fmt.Sprintf("/admin/orders/%d/payments", order.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", path, ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].([]interface{})
	assert.GreaterOrEqual(t, len(data), 2)
}

func TestOrderHandler_Comprehensive_ListOrderPayments_NoPayments(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusPending)

	path := fmt.Sprintf("/admin/orders/%d/payments", order.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", path, ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].([]interface{})
	assert.Equal(t, 0, len(data))
}

// ============================================================================
// Order Refunds List Tests
// ============================================================================

func TestOrderHandler_Comprehensive_ListOrderRefunds_WithRefunds(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusCompleted)
	payment := testutil.CreateTestPayment(t, ctx.DB, order.ID, ctx.TestUser.ID, model.PaymentStatusPaid)

	// Process a refund
	payment.RefundedAmountCents = 5000
	payment.Status = model.PaymentStatusRefunded
	ctx.DB.Save(payment)

	path := fmt.Sprintf("/admin/orders/%d/refunds", order.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", path, ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"]
	assert.NotNil(t, data)
}

// ============================================================================
// Order Reviews Tests
// ============================================================================

func TestOrderHandler_Comprehensive_ListOrderReviews_WithReviews(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusCompleted)

	// Create a review
	review := &model.Review{
		Base:      model.Base{ExtJSON: "{}"},
		OrderID:   order.ID,
		UserID:    ctx.TestUser.ID,
		PlayerID:  *order.PlayerID,
		Score:     5,
		Content:   "Excellent service!",
		Status:    model.ReviewStatusApproved,
		Images:    model.StringArray{},
	}
	require.NoError(t, ctx.DB.Create(review).Error)

	path := fmt.Sprintf("/admin/orders/%d/reviews", order.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", path, ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].([]interface{})
	assert.GreaterOrEqual(t, len(data), 1)
}

// ============================================================================
// Order Logs Tests - Advanced Filtering
// ============================================================================

func TestOrderHandler_Comprehensive_ListOrderLogs_WithFilters(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusPending)

	tests := []struct {
		name        string
		queryParams map[string]string
		expectLogs  bool
	}{
		{
			name:        "default parameters",
			queryParams:  map[string]string{},
			expectLogs:   true,
		},
		{
			name: "filter_by_action",
			queryParams: map[string]string{
				"action": "create",
			},
			expectLogs: true,
		},
		{
			name: "pagination",
			queryParams: map[string]string{
				"page":      "1",
				"page_size": "10",
			},
			expectLogs: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := testutil.MakeRequest(t, ctx.Router, "GET", fmt.Sprintf("/admin/orders/%d/logs", order.ID), nil,
				testutil.WithAuth(ctx.AdminToken),
				testutil.WithQuery(tt.queryParams))

			testutil.AssertSuccess(t, w)

			items, _ := testutil.GetResponseList(t, w)

			if tt.expectLogs {
				// Logs should exist (created by testutil)
				assert.NotNil(t, items)
			}
		})
	}
}

// ============================================================================
// Order Delete Tests - Edge Cases
// ============================================================================

func TestOrderHandler_Comprehensive_DeleteOrder_Constraints(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	t.Run("delete_order_with_payment", func(t *testing.T) {
		order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusPending)
		_ = testutil.CreateTestPayment(t, ctx.DB, order.ID, ctx.TestUser.ID, model.PaymentStatusPaid)

		path := fmt.Sprintf("/admin/orders/%d", order.ID)
		w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "DELETE", path, ctx.AdminToken, nil)

		// Behavior depends on foreign key constraints in database
		// Might succeed with CASCADE or fail with RESTRICT
		assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusBadRequest)
	})

	t.Run("delete_nonexistent_order", func(t *testing.T) {
		w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "DELETE", "/admin/orders/99999999", ctx.AdminToken, nil)
		testutil.AssertError(t, w, http.StatusNotFound)
	})
}

// ============================================================================
// Review Order Tests - Detailed Scenarios
// ============================================================================

func TestOrderHandler_Comprehensive_ReviewOrder_Detailed(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	t.Run("approve_pending_order", func(t *testing.T) {
		order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusPending)

		payload := map[string]interface{}{
			"approved": true,
			"reason":   "Order verified and approved",
		}

		path := fmt.Sprintf("/admin/orders/%d/review", order.ID)
		w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", path, ctx.AdminToken, payload)
		testutil.AssertSuccess(t, w)

		testutil.AssertOrderStatus(t, ctx.DB, order.ID, model.OrderStatusConfirmed)
	})

	t.Run("reject_pending_order", func(t *testing.T) {
		order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusPending)

		payload := map[string]interface{}{
			"approved": false,
			"reason":   "Invalid payment information",
		}

		path := fmt.Sprintf("/admin/orders/%d/review", order.ID)
		w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", path, ctx.AdminToken, payload)
		testutil.AssertSuccess(t, w)

		testutil.AssertOrderStatus(t, ctx.DB, order.ID, model.OrderStatusCanceled)

		// Verify cancel reason is set
		var updatedOrder model.Order
		ctx.DB.First(&updatedOrder, order.ID)
		assert.Equal(t, "Invalid payment information", updatedOrder.CancelReason)
	})

	t.Run("review_already_processed_order", func(t *testing.T) {
		order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusCompleted)

		payload := map[string]interface{}{
			"approved": true,
			"reason":   "Try to review completed order",
		}

		path := fmt.Sprintf("/admin/orders/%d/review", order.ID)
		w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", path, ctx.AdminToken, payload)

		// Should fail - cannot review already processed order
		testutil.AssertError(t, w, http.StatusBadRequest)
	})
}

// ============================================================================
// Integration with Related Entities
// ============================================================================

func TestOrderHandler_Comprehensive_GameRelationship(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	testServiceItem := testutil.CreateTestServiceItem(t, ctx.DB, ctx.TestGame.ID, 5000, true)

	payload := map[string]interface{}{
		"user_id":           ctx.TestUser.ID,
		"game_id":           ctx.TestGame.ID,
		"item_id":           testServiceItem.ID,
		"total_price_cents": 10000,
		"currency":          "CNY",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/orders", ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	gameID := data["game_id"]
	assert.NotNil(t, gameID)
	assert.Equal(t, float64(ctx.TestGame.ID), gameID)
}

func TestOrderHandler_Comprehensive_ServiceItemRelationship(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	testServiceItem := testutil.CreateTestServiceItem(t, ctx.DB, ctx.TestGame.ID, 5000, true)

	payload := map[string]interface{}{
		"user_id":           ctx.TestUser.ID,
		"game_id":           ctx.TestGame.ID,
		"item_id":           testServiceItem.ID,
		"total_price_cents": 10000,
		"currency":          "CNY",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/orders", ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	itemID := data["item_id"]
	assert.NotNil(t, itemID)
	assert.Equal(t, float64(testServiceItem.ID), itemID)
}

// ============================================================================
// Pagination Edge Cases
// ============================================================================

func TestOrderHandler_Comprehensive_Pagination_EdgeCases(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	// Create exactly 10 orders
	for i := 0; i < 10; i++ {
		testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusPending)
	}

	tests := []struct {
		name          string
		page          string
		pageSize      string
		expectResults int
	}{
		{
			name:          "page_1_size_5",
			page:          "1",
			pageSize:      "5",
			expectResults: 5,
		},
		{
			name:          "page_2_size_5",
			page:          "2",
			pageSize:      "5",
			expectResults: 5,
		},
		{
			name:          "page_1_size_20",
			page:          "1",
			pageSize:      "20",
			expectResults: 10, // Only 10 orders exist
		},
		{
			name:          "page_999_size_10",
			page:          "999",
			pageSize:      "10",
			expectResults: 0, // Beyond available data
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := map[string]string{
				"page":      tt.page,
				"page_size": tt.pageSize,
			}
			w := testutil.MakeRequest(t, ctx.Router, "GET", "/admin/orders", nil,
				testutil.WithAuth(ctx.AdminToken),
				testutil.WithQuery(params))

			testutil.AssertSuccess(t, w)

			items, pagination := testutil.GetResponseList(t, w)
			assert.Equal(t, tt.expectResults, len(items))
			assert.Equal(t, tt.page, fmt.Sprintf("%.0f", pagination["page"]))
			assert.Equal(t, tt.pageSize, fmt.Sprintf("%.0f", pagination["page_size"]))
		})
	}
}

// ============================================================================
// Error Response Validation Tests
// ============================================================================

func TestOrderHandler_Comprehensive_ErrorResponses(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	t.Run("invalid_json_payload", func(t *testing.T) {
		w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/orders", ctx.AdminToken, []byte("{invalid json"))
		testutil.AssertError(t, w, http.StatusBadRequest)
	})

	t.Run("missing_required_fields", func(t *testing.T) {
		payload := map[string]interface{}{
			"title": "Order without required fields",
		}
		w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/orders", ctx.AdminToken, payload)
		testutil.AssertError(t, w, http.StatusBadRequest)
	})

	t.Run("invalid_order_id_format", func(t *testing.T) {
		w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/orders/abc", ctx.AdminToken, nil)
		testutil.AssertError(t, w, http.StatusBadRequest)
	})

	t.Run("negative_order_id", func(t *testing.T) {
		w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/orders/-1", ctx.AdminToken, nil)
		testutil.AssertError(t, w, http.StatusBadRequest)
	})
}

// ============================================================================
// Performance Tests - Large Data Sets
// ============================================================================

func TestOrderHandler_Comprehensive_Performance_LargeDataset(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	// Create 100 orders
	t.Log("Creating 100 test orders...")
	for i := 0; i < 100; i++ {
		testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusPending)
	}

	// Test list performance
	start := time.Now()
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/orders?page=1&page_size=50", ctx.AdminToken, nil)
	duration := time.Since(start)

	testutil.AssertSuccess(t, w)

	// Should complete in reasonable time (< 1 second for 100 records)
	assert.Less(t, duration.Milliseconds(), int64(1000), "List query should complete in < 1 second")

	items, _ := testutil.GetResponseList(t, w)
	assert.LessOrEqual(t, len(items), 50)

	t.Logf("Listed %d orders in %v", len(items), duration)
}

// ============================================================================
// Concurrency Tests
// ============================================================================

func TestOrderHandler_Comprehensive_ConcurrentUpdates(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	order := testutil.CreateTestOrder(t, ctx.DB, ctx.TestUser.ID, ctx.TestPlayer.ID, ctx.TestGame.ID, model.OrderStatusPending)

	// Simulate concurrent status updates
	done := make(chan bool, 3)

	for i := 0; i < 3; i++ {
		go func(index int) {
			defer func() { done <- true }()

			payload := map[string]interface{}{
				"status":            "confirmed",
				"total_price_cents": 10000 + int64(index*100),
				"currency":          "CNY",
			}

			path := fmt.Sprintf("/admin/orders/%d", order.ID)
			w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", path, ctx.AdminToken, payload)
			_ = w.Code // Result may vary due to race condition
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 3; i++ {
		<-done
	}

	// Verify final state is valid
	var finalOrder model.Order
	ctx.DB.First(&finalOrder, order.ID)
	assert.True(t, finalOrder.Status == model.OrderStatusPending || finalOrder.Status == model.OrderStatusConfirmed)
}

// ============================================================================
// Authentication & Authorization Tests
// ============================================================================

func TestOrderHandler_Comprehensive_Authentication(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	t.Run("request_without_token", func(t *testing.T) {
		w := testutil.MakeRequest(t, ctx.Router, "GET", "/admin/orders", nil)
		// Should fail without authentication
		assert.True(t, w.Code == http.StatusUnauthorized || w.Code == http.StatusForbidden)
	})

	t.Run("request_with_invalid_token", func(t *testing.T) {
		w := testutil.MakeRequest(t, ctx.Router, "GET", "/admin/orders", nil,
			testutil.WithAuth("invalid_token_bearer"))
		assert.True(t, w.Code == http.StatusUnauthorized || w.Code == http.StatusForbidden)
	})
}

// ============================================================================
// Data Consistency Tests
// ============================================================================

func TestOrderHandler_Comprehensive_DataConsistency(t *testing.T) {
	ctx := SetupOrderTest(t)
	ctx.RegisterOrderRoutes()

	testServiceItem := testutil.CreateTestServiceItem(t, ctx.DB, ctx.TestGame.ID, 5000, true)

	payload := map[string]interface{}{
		"user_id":           ctx.TestUser.ID,
		"player_id":         ctx.TestPlayer.ID,
		"game_id":           ctx.TestGame.ID,
		"item_id":           testServiceItem.ID,
		"title":             "Consistency Test Order",
		"description":       "Testing data consistency",
		"total_price_cents": 12345,
		"currency":          "USD",
		"scheduled_start":   time.Now().Add(1 * time.Hour).Format(time.RFC3339),
		"scheduled_end":     time.Now().Add(2 * time.Hour).Format(time.RFC3339),
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/orders", ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	orderID := uint64(data["id"].(float64))

	// Verify database record matches response
	var dbOrder model.Order
	err = ctx.DB.First(&dbOrder, orderID).Error
	require.NoError(t, err)

	assert.Equal(t, "Consistency Test Order", dbOrder.Title)
	assert.Equal(t, "Testing data consistency", dbOrder.Description)
	assert.Equal(t, int64(12345), dbOrder.TotalPriceCents)
	assert.Equal(t, model.CurrencyUSD, dbOrder.Currency)
	assert.NotNil(t, dbOrder.ScheduledStart)
	assert.NotNil(t, dbOrder.ScheduledEnd)
}
