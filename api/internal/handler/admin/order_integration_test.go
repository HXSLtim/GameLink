// Package admin provides integration tests for order admin handlers.
package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/service/integration"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestOrderHandler_ListOrders tests the list orders endpoint with various filters.
func TestOrderHandler_ListOrders(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test data
	testUser := integration.CreateUniqueTestUser(t, helper.db, "list_user")
	playerUser := integration.CreateUniqueTestUser(t, helper.db, "list_player")
	testPlayer := integration.CreateTestPlayer(t, helper.db, playerUser)
	testGame := integration.CreateTestGame(t, helper.db, "list_game")

	// Create orders with different statuses
	_ = integration.CreateTestOrderWithDetails(t, helper.db, testUser, testPlayer, testGame, model.OrderStatusPending, 5000)
	_ = integration.CreateTestOrderWithDetails(t, helper.db, testUser, testPlayer, testGame, model.OrderStatusConfirmed, 10000)
	_ = integration.CreateTestOrderWithDetails(t, helper.db, testUser, testPlayer, testGame, model.OrderStatusInProgress, 15000)
	_ = integration.CreateTestOrderWithDetails(t, helper.db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 20000)

	t.Run("ListAllOrders", func(t *testing.T) {
		w := helper.MakeRequest("GET", "/admin/orders", nil, true)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))
		data := response["data"].(map[string]interface{})
		items := data["items"].([]interface{})
		assert.GreaterOrEqual(t, len(items), 4)
	})

	t.Run("ListOrdersWithStatusFilter", func(t *testing.T) {
		w := helper.MakeRequest("GET", "/admin/orders?status=pending", nil, true)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		data := response["data"].(map[string]interface{})
		items := data["items"].([]interface{})
		assert.Greater(t, len(items), 0)
	})

	t.Run("ListOrdersWithUserFilter", func(t *testing.T) {
		url := fmt.Sprintf("/admin/orders?user_id=%d", testUser.ID)
		w := helper.MakeRequest("GET", url, nil, true)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		data := response["data"].(map[string]interface{})
		items := data["items"].([]interface{})
		assert.Greater(t, len(items), 0)
	})

	t.Run("ListOrdersWithPlayerFilter", func(t *testing.T) {
		url := fmt.Sprintf("/admin/orders?player_id=%d", testPlayer.ID)
		w := helper.MakeRequest("GET", url, nil, true)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		data := response["data"].(map[string]interface{})
		items := data["items"].([]interface{})
		assert.Greater(t, len(items), 0)
	})

	t.Run("ListOrdersWithGameFilter", func(t *testing.T) {
		url := fmt.Sprintf("/admin/orders?game_id=%d", testGame.ID)
		w := helper.MakeRequest("GET", url, nil, true)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		data := response["data"].(map[string]interface{})
		items := data["items"].([]interface{})
		assert.Greater(t, len(items), 0)
	})

	t.Run("ListOrdersWithDateRange", func(t *testing.T) {
		today := time.Now().Format("2006-01-02")
		tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
		url := fmt.Sprintf("/admin/orders?date_from=%s&date_to=%s", today, tomorrow)
		w := helper.MakeRequest("GET", url, nil, true)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("ListOrdersWithPagination", func(t *testing.T) {
		w := helper.MakeRequest("GET", "/admin/orders?page=1&page_size=2", nil, true)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		data := response["data"].(map[string]interface{})
		items := data["items"].([]interface{})
		assert.LessOrEqual(t, len(items), 2)
	})
}

// TestOrderHandler_CreateOrder tests creating an order.
func TestOrderHandler_CreateOrder(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test data
	testUser := integration.CreateUniqueTestUser(t, helper.db, "create_user")
	testGame := integration.CreateTestGame(t, helper.db, "create_game")

	t.Run("CreateOrderSuccess", func(t *testing.T) {
		payload := map[string]interface{}{
			"user_id":           testUser.ID,
			"game_id":           testGame.ID,
			"title":             "Test Order",
			"description":       "Test order description",
			"total_price_cents": 10000,
			"currency":          "CNY",
		}

		w := helper.MakeRequest("POST", "/admin/orders", payload, true)
		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))
		data := response["data"].(map[string]interface{})
		assert.Equal(t, "Test Order", data["title"])
		assert.Equal(t, int64(10000), data["total_price_cents"])
	})

	t.Run("CreateOrderWithPlayer", func(t *testing.T) {
		playerUser := integration.CreateUniqueTestUser(t, helper.db, "create_player_user")
		testPlayer := integration.CreateTestPlayer(t, helper.db, playerUser)

		payload := map[string]interface{}{
			"user_id":           testUser.ID,
			"player_id":         testPlayer.ID,
			"game_id":           testGame.ID,
			"title":             "Order with Player",
			"total_price_cents": 15000,
			"currency":          "CNY",
		}

		w := helper.MakeRequest("POST", "/admin/orders", payload, true)
		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))
		data := response["data"].(map[string]interface{})
		assert.Equal(t, "Order with Player", data["title"])
	})

	t.Run("CreateOrderWithSchedule", func(t *testing.T) {
		startTime := time.Now().Add(time.Hour).Format(time.RFC3339)
		endTime := time.Now().Add(2 * time.Hour).Format(time.RFC3339)

		payload := map[string]interface{}{
			"user_id":           testUser.ID,
			"game_id":           testGame.ID,
			"title":             "Scheduled Order",
			"total_price_cents": 12000,
			"currency":          "CNY",
			"scheduled_start":   startTime,
			"scheduled_end":     endTime,
		}

		w := helper.MakeRequest("POST", "/admin/orders", payload, true)
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("CreateOrderValidationFailure", func(t *testing.T) {
		payload := map[string]interface{}{
			"user_id": testUser.ID,
			// missing required game_id
			"title":             "Invalid Order",
			"total_price_cents": 10000,
			"currency":          "CNY",
		}

		w := helper.MakeRequest("POST", "/admin/orders", payload, true)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("CreateOrderInvalidTimeFormat", func(t *testing.T) {
		payload := map[string]interface{}{
			"user_id":           testUser.ID,
			"game_id":           testGame.ID,
			"title":             "Invalid Time Order",
			"total_price_cents": 10000,
			"currency":          "CNY",
			"scheduled_start":   "invalid-time",
		}

		w := helper.MakeRequest("POST", "/admin/orders", payload, true)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestOrderHandler_GetOrder tests getting a single order.
func TestOrderHandler_GetOrder(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test data
	testUser := integration.CreateUniqueTestUser(t, helper.db, "get_user")
	playerUser := integration.CreateUniqueTestUser(t, helper.db, "get_player")
	testPlayer := integration.CreateTestPlayer(t, helper.db, playerUser)
	testGame := integration.CreateTestGame(t, helper.db, "get_game")
	testOrder := integration.CreateTestOrderWithDetails(t, helper.db, testUser, testPlayer, testGame, model.OrderStatusConfirmed, 10000)

	t.Run("GetOrderSuccess", func(t *testing.T) {
		url := fmt.Sprintf("/admin/orders/%d", testOrder.ID)
		w := helper.MakeRequest("GET", url, nil, true)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))
		data := response["data"].(map[string]interface{})
		assert.Equal(t, float64(testOrder.ID), data["id"])
	})

	t.Run("GetOrderNotFound", func(t *testing.T) {
		w := helper.MakeRequest("GET", "/admin/orders/999999", nil, true)
		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.False(t, response["success"].(bool))
	})

	t.Run("GetOrderInvalidID", func(t *testing.T) {
		w := helper.MakeRequest("GET", "/admin/orders/invalid", nil, true)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestOrderHandler_UpdateOrder tests updating an order.
func TestOrderHandler_UpdateOrder(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test data
	testUser := integration.CreateUniqueTestUser(t, helper.db, "update_user")
	playerUser := integration.CreateUniqueTestUser(t, helper.db, "update_player")
	testPlayer := integration.CreateTestPlayer(t, helper.db, playerUser)
	testGame := integration.CreateTestGame(t, helper.db, "update_game")
	testOrder := integration.CreateTestOrderWithDetails(t, helper.db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)

	t.Run("UpdateOrderStatus", func(t *testing.T) {
		payload := map[string]interface{}{
			"status":            "confirmed",
			"total_price_cents": 10000,
			"currency":          "CNY",
		}

		url := fmt.Sprintf("/admin/orders/%d", testOrder.ID)
		w := helper.MakeRequest("PUT", url, payload, true)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))
		data := response["data"].(map[string]interface{})
		assert.Equal(t, "confirmed", data["status"])
	})

	t.Run("UpdateOrderPrice", func(t *testing.T) {
		payload := map[string]interface{}{
			"status":            "pending",
			"total_price_cents": 15000,
			"currency":          "CNY",
		}

		url := fmt.Sprintf("/admin/orders/%d", testOrder.ID)
		w := helper.MakeRequest("PUT", url, payload, true)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		data := response["data"].(map[string]interface{})
		assert.Equal(t, int64(15000), data["total_price_cents"])
	})

	t.Run("UpdateOrderWithSchedule", func(t *testing.T) {
		startTime := time.Now().Add(2 * time.Hour).Format(time.RFC3339)
		endTime := time.Now().Add(4 * time.Hour).Format(time.RFC3339)

		payload := map[string]interface{}{
			"status":            "pending",
			"total_price_cents": 10000,
			"currency":          "CNY",
			"scheduled_start":   startTime,
			"scheduled_end":     endTime,
		}

		url := fmt.Sprintf("/admin/orders/%d", testOrder.ID)
		w := helper.MakeRequest("PUT", url, payload, true)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("UpdateOrderNotFound", func(t *testing.T) {
		payload := map[string]interface{}{
			"status":            "confirmed",
			"total_price_cents": 10000,
			"currency":          "CNY",
		}

		w := helper.MakeRequest("PUT", "/admin/orders/999999", payload, true)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("UpdateOrderInvalidStatus", func(t *testing.T) {
		// Create a completed order
		completedOrder := integration.CreateTestOrderWithDetails(t, helper.db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)

		// Try to update to pending (invalid transition)
		payload := map[string]interface{}{
			"status":            "pending",
			"total_price_cents": 10000,
			"currency":          "CNY",
		}

		url := fmt.Sprintf("/admin/orders/%d", completedOrder.ID)
		w := helper.MakeRequest("PUT", url, payload, true)
		// Should return validation error
		assert.NotEqual(t, http.StatusOK, w.Code)
	})
}

// TestOrderHandler_DeleteOrder tests deleting an order.
func TestOrderHandler_DeleteOrder(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test data
	testUser := integration.CreateUniqueTestUser(t, helper.db, "delete_user")
	playerUser := integration.CreateUniqueTestUser(t, helper.db, "delete_player")
	testPlayer := integration.CreateTestPlayer(t, helper.db, playerUser)
	testGame := integration.CreateTestGame(t, helper.db, "delete_game")
	testOrder := integration.CreateTestOrderWithDetails(t, helper.db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)

	t.Run("DeleteOrderSuccess", func(t *testing.T) {
		url := fmt.Sprintf("/admin/orders/%d", testOrder.ID)
		w := helper.MakeRequest("DELETE", url, nil, true)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))
	})

	t.Run("DeleteOrderNotFound", func(t *testing.T) {
		w := helper.MakeRequest("DELETE", "/admin/orders/999999", nil, true)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// TestOrderHandler_ReviewOrder tests reviewing an order.
func TestOrderHandler_ReviewOrder(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test data
	testUser := integration.CreateUniqueTestUser(t, helper.db, "review_user")
	playerUser := integration.CreateUniqueTestUser(t, helper.db, "review_player")
	testPlayer := integration.CreateTestPlayer(t, helper.db, playerUser)
	testGame := integration.CreateTestGame(t, helper.db, "review_game")
	testOrder := integration.CreateTestOrderWithDetails(t, helper.db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)

	t.Run("ReviewOrderApprove", func(t *testing.T) {
		payload := map[string]interface{}{
			"approved": true,
			"reason":   "Order verified",
		}

		url := fmt.Sprintf("/admin/orders/%d/review", testOrder.ID)
		w := helper.MakeRequest("POST", url, payload, true)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))
		data := response["data"].(map[string]interface{})
		assert.Equal(t, "confirmed", data["status"])
	})

	t.Run("ReviewOrderReject", func(t *testing.T) {
		rejectOrder := integration.CreateTestOrderWithDetails(t, helper.db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)

		payload := map[string]interface{}{
			"approved": false,
			"reason":   "Invalid order information",
		}

		url := fmt.Sprintf("/admin/orders/%d/review", rejectOrder.ID)
		w := helper.MakeRequest("POST", url, payload, true)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		data := response["data"].(map[string]interface{})
		assert.Equal(t, "canceled", data["status"])
		assert.Equal(t, "Invalid order information", data["cancel_reason"])
	})
}

// TestOrderHandler_CancelOrder tests canceling an order.
func TestOrderHandler_CancelOrder(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test data
	testUser := integration.CreateUniqueTestUser(t, helper.db, "cancel_user")
	playerUser := integration.CreateUniqueTestUser(t, helper.db, "cancel_player")
	testPlayer := integration.CreateTestPlayer(t, helper.db, playerUser)
	testGame := integration.CreateTestGame(t, helper.db, "cancel_game")
	testOrder := integration.CreateTestOrderWithDetails(t, helper.db, testUser, testPlayer, testGame, model.OrderStatusConfirmed, 10000)

	t.Run("CancelOrderSuccess", func(t *testing.T) {
		payload := map[string]interface{}{
			"reason": "Customer request",
		}

		url := fmt.Sprintf("/admin/orders/%d/cancel", testOrder.ID)
		w := helper.MakeRequest("POST", url, payload, true)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))
		data := response["data"].(map[string]interface{})
		assert.Equal(t, "canceled", data["status"])
		assert.Equal(t, "Customer request", data["cancel_reason"])
	})
}

// TestOrderHandler_AssignOrder tests assigning a player to an order.
func TestOrderHandler_AssignOrder(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test data
	testUser := integration.CreateUniqueTestUser(t, helper.db, "assign_user")
	playerUser1 := integration.CreateUniqueTestUser(t, helper.db, "assign_player1")
	playerUser2 := integration.CreateUniqueTestUser(t, helper.db, "assign_player2")
	testPlayer1 := integration.CreateTestPlayer(t, helper.db, playerUser1)
	testPlayer2 := integration.CreateTestPlayer(t, helper.db, playerUser2)
	testGame := integration.CreateTestGame(t, helper.db, "assign_game")
	testOrder := integration.CreateTestOrderWithDetails(t, helper.db, testUser, testPlayer1, testGame, model.OrderStatusPending, 10000)

	t.Run("AssignOrderSuccess", func(t *testing.T) {
		payload := map[string]interface{}{
			"player_id": testPlayer2.ID,
		}

		url := fmt.Sprintf("/admin/orders/%d/assign", testOrder.ID)
		w := helper.MakeRequest("POST", url, payload, true)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))
		data := response["data"].(map[string]interface{})
		assert.Equal(t, float64(testPlayer2.ID), data["player_id"])
	})

	t.Run("AssignOrderNotFound", func(t *testing.T) {
		payload := map[string]interface{}{
			"player_id": testPlayer1.ID,
		}

		w := helper.MakeRequest("POST", "/admin/orders/999999/assign", payload, true)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// TestOrderHandler_ConfirmOrder tests confirming an order.
func TestOrderHandler_ConfirmOrder(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test data
	testUser := integration.CreateUniqueTestUser(t, helper.db, "confirm_user")
	playerUser := integration.CreateUniqueTestUser(t, helper.db, "confirm_player")
	testPlayer := integration.CreateTestPlayer(t, helper.db, playerUser)
	testGame := integration.CreateTestGame(t, helper.db, "confirm_game")
	testOrder := integration.CreateTestOrderWithDetails(t, helper.db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)

	t.Run("ConfirmOrderSuccess", func(t *testing.T) {
		payload := map[string]interface{}{
			"note": "Order confirmed by admin",
		}

		url := fmt.Sprintf("/admin/orders/%d/confirm", testOrder.ID)
		w := helper.MakeRequest("POST", url, payload, true)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))
		data := response["data"].(map[string]interface{})
		assert.Equal(t, "confirmed", data["status"])
	})

	t.Run("ConfirmOrderWithoutNote", func(t *testing.T) {
		confirmOrder := integration.CreateTestOrderWithDetails(t, helper.db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)

		url := fmt.Sprintf("/admin/orders/%d/confirm", confirmOrder.ID)
		w := helper.MakeRequest("POST", url, nil, true)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

// TestOrderHandler_StartOrder tests starting an order.
func TestOrderHandler_StartOrder(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test data
	testUser := integration.CreateUniqueTestUser(t, helper.db, "start_user")
	playerUser := integration.CreateUniqueTestUser(t, helper.db, "start_player")
	testPlayer := integration.CreateTestPlayer(t, helper.db, playerUser)
	testGame := integration.CreateTestGame(t, helper.db, "start_game")
	testOrder := integration.CreateTestOrderWithDetails(t, helper.db, testUser, testPlayer, testGame, model.OrderStatusConfirmed, 10000)

	t.Run("StartOrderSuccess", func(t *testing.T) {
		payload := map[string]interface{}{
			"note": "Service started",
		}

		url := fmt.Sprintf("/admin/orders/%d/start", testOrder.ID)
		w := helper.MakeRequest("POST", url, payload, true)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))
		data := response["data"].(map[string]interface{})
		assert.Equal(t, "in_progress", data["status"])
		assert.NotNil(t, data["started_at"])
	})
}

// TestOrderHandler_CompleteOrder tests completing an order.
func TestOrderHandler_CompleteOrder(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test data
	testUser := integration.CreateUniqueTestUser(t, helper.db, "complete_user")
	playerUser := integration.CreateUniqueTestUser(t, helper.db, "complete_player")
	testPlayer := integration.CreateTestPlayer(t, helper.db, playerUser)
	testGame := integration.CreateTestGame(t, helper.db, "complete_game")
	testOrder := integration.CreateTestOrderWithDetails(t, helper.db, testUser, testPlayer, testGame, model.OrderStatusInProgress, 10000)

	t.Run("CompleteOrderSuccess", func(t *testing.T) {
		payload := map[string]interface{}{
			"note": "Service completed successfully",
		}

		url := fmt.Sprintf("/admin/orders/%d/complete", testOrder.ID)
		w := helper.MakeRequest("POST", url, payload, true)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))
		data := response["data"].(map[string]interface{})
		assert.Equal(t, "completed", data["status"])
		assert.NotNil(t, data["completed_at"])
	})
}

// TestOrderHandler_RefundOrder tests refunding an order.
func TestOrderHandler_RefundOrder(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test data
	testUser := integration.CreateUniqueTestUser(t, helper.db, "refund_user")
	playerUser := integration.CreateUniqueTestUser(t, helper.db, "refund_player")
	testPlayer := integration.CreateTestPlayer(t, helper.db, playerUser)
	testGame := integration.CreateTestGame(t, helper.db, "refund_game")
	testOrder := integration.CreateTestOrderWithDetails(t, helper.db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)

	t.Run("RefundOrderSuccess", func(t *testing.T) {
		payload := map[string]interface{}{
			"reason":       "Customer satisfaction issue",
			"amount_cents": int64(5000),
			"note":         "Partial refund approved",
		}

		url := fmt.Sprintf("/admin/orders/%d/refund", testOrder.ID)
		w := helper.MakeRequest("POST", url, payload, true)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))
		data := response["data"].(map[string]interface{})
		assert.Equal(t, "refunded", data["status"])
	})

	t.Run("RefundOrderFullRefund", func(t *testing.T) {
		fullRefundOrder := integration.CreateTestOrderWithDetails(t, helper.db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 8000)

		payload := map[string]interface{}{
			"reason": "Service not provided",
			"note":   "Full refund",
		}

		url := fmt.Sprintf("/admin/orders/%d/refund", fullRefundOrder.ID)
		w := helper.MakeRequest("POST", url, payload, true)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))
	})
}

// TestOrderHandler_ListOrderLogs tests listing order operation logs.
func TestOrderHandler_ListOrderLogs(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test data
	testUser := integration.CreateUniqueTestUser(t, helper.db, "logs_user")
	playerUser := integration.CreateUniqueTestUser(t, helper.db, "logs_player")
	testPlayer := integration.CreateTestPlayer(t, helper.db, playerUser)
	testGame := integration.CreateTestGame(t, helper.db, "logs_game")
	testOrder := integration.CreateTestOrderWithDetails(t, helper.db, testUser, testPlayer, testGame, model.OrderStatusConfirmed, 10000)

	t.Run("ListOrderLogsSuccess", func(t *testing.T) {
		url := fmt.Sprintf("/admin/orders/%d/logs", testOrder.ID)
		w := helper.MakeRequest("GET", url, nil, true)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))
		data := response["data"].(map[string]interface{})
		items := data["items"].([]interface{})
		assert.NotNil(t, items)
	})

	t.Run("ListOrderLogsWithPagination", func(t *testing.T) {
		url := fmt.Sprintf("/admin/orders/%d/logs?page=1&page_size=10", testOrder.ID)
		w := helper.MakeRequest("GET", url, nil, true)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("ListOrderLogsByAction", func(t *testing.T) {
		url := fmt.Sprintf("/admin/orders/%d/logs?action=update_status", testOrder.ID)
		w := helper.MakeRequest("GET", url, nil, true)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

// TestOrderHandler_GetOrderTimeline tests getting order timeline.
func TestOrderHandler_GetOrderTimeline(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test data
	testUser := integration.CreateUniqueTestUser(t, helper.db, "timeline_user")
	playerUser := integration.CreateUniqueTestUser(t, helper.db, "timeline_player")
	testPlayer := integration.CreateTestPlayer(t, helper.db, playerUser)
	testGame := integration.CreateTestGame(t, helper.db, "timeline_game")
	testOrder := integration.CreateTestOrderWithDetails(t, helper.db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)

	t.Run("GetOrderTimelineSuccess", func(t *testing.T) {
		url := fmt.Sprintf("/admin/orders/%d/timeline", testOrder.ID)
		w := helper.MakeRequest("GET", url, nil, true)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))
		data := response["data"].([]interface{})
		assert.NotNil(t, data)
	})

	t.Run("GetOrderTimelineNotFound", func(t *testing.T) {
		w := helper.MakeRequest("GET", "/admin/orders/999999/timeline", nil, true)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// TestOrderHandler_ListOrderPayments tests listing order payments.
func TestOrderHandler_ListOrderPayments(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test data
	testUser := integration.CreateUniqueTestUser(t, helper.db, "payments_user")
	playerUser := integration.CreateUniqueTestUser(t, helper.db, "payments_player")
	testPlayer := integration.CreateTestPlayer(t, helper.db, playerUser)
	testGame := integration.CreateTestGame(t, helper.db, "payments_game")
	testOrder := integration.CreateTestOrderWithDetails(t, helper.db, testUser, testPlayer, testGame, model.OrderStatusConfirmed, 10000)
	_ = integration.CreateTestPayment(t, helper.db, testOrder, model.PaymentStatusPaid)

	t.Run("ListOrderPaymentsSuccess", func(t *testing.T) {
		url := fmt.Sprintf("/admin/orders/%d/payments", testOrder.ID)
		w := helper.MakeRequest("GET", url, nil, true)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))
		data := response["data"].([]interface{})
		assert.Greater(t, len(data), 0)
	})

	t.Run("ListOrderPaymentsNotFound", func(t *testing.T) {
		w := helper.MakeRequest("GET", "/admin/orders/999999/payments", nil, true)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// TestOrderHandler_ListOrderRefunds tests listing order refunds.
func TestOrderHandler_ListOrderRefunds(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test data
	testUser := integration.CreateUniqueTestUser(t, helper.db, "refunds_user")
	playerUser := integration.CreateUniqueTestUser(t, helper.db, "refunds_player")
	testPlayer := integration.CreateTestPlayer(t, helper.db, playerUser)
	testGame := integration.CreateTestGame(t, helper.db, "refunds_game")
	testOrder := integration.CreateTestOrderWithDetails(t, helper.db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)

	t.Run("ListOrderRefundsSuccess", func(t *testing.T) {
		url := fmt.Sprintf("/admin/orders/%d/refunds", testOrder.ID)
		w := helper.MakeRequest("GET", url, nil, true)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))
		data := response["data"].([]interface{})
		assert.NotNil(t, data)
	})
}

// TestOrderHandler_ListOrderReviews tests listing order reviews.
func TestOrderHandler_ListOrderReviews(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test data
	testUser := integration.CreateUniqueTestUser(t, helper.db, "reviews_user")
	playerUser := integration.CreateUniqueTestUser(t, helper.db, "reviews_player")
	testPlayer := integration.CreateTestPlayer(t, helper.db, playerUser)
	testGame := integration.CreateTestGame(t, helper.db, "reviews_game")
	testOrder := integration.CreateTestOrderWithDetails(t, helper.db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)
	_ = integration.CreateTestReview(t, helper.db, testOrder, model.Rating5)

	t.Run("ListOrderReviewsSuccess", func(t *testing.T) {
		url := fmt.Sprintf("/admin/orders/%d/reviews", testOrder.ID)
		w := helper.MakeRequest("GET", url, nil, true)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))
		data := response["data"].([]interface{})
		assert.Greater(t, len(data), 0)
	})

	t.Run("ListOrderReviewsEmpty", func(t *testing.T) {
		emptyOrder := integration.CreateTestOrderWithDetails(t, helper.db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)

		url := fmt.Sprintf("/admin/orders/%d/reviews", emptyOrder.ID)
		w := helper.MakeRequest("GET", url, nil, true)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))
		data := response["data"].([]interface{})
		assert.Equal(t, 0, len(data))
	})
}

// TestOrderHandler_OrderWorkflow tests a complete order workflow.
func TestOrderHandler_OrderWorkflow(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test data
	testUser := integration.CreateUniqueTestUser(t, helper.db, "workflow_user")
	playerUser := integration.CreateUniqueTestUser(t, helper.db, "workflow_player")
	testPlayer := integration.CreateTestPlayer(t, helper.db, playerUser)
	testGame := integration.CreateTestGame(t, helper.db, "workflow_game")

	t.Run("CompleteOrderWorkflow", func(t *testing.T) {
		// 1. Create order
		createPayload := map[string]interface{}{
			"user_id":           testUser.ID,
			"player_id":         testPlayer.ID,
			"game_id":           testGame.ID,
			"title":             "Workflow Test Order",
			"total_price_cents": 10000,
			"currency":          "CNY",
		}

		w := helper.MakeRequest("POST", "/admin/orders", createPayload, true)
		require.Equal(t, http.StatusCreated, w.Code)

		var createResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &createResponse)
		require.NoError(t, err)
		createData := createResponse["data"].(map[string]interface{})
		orderID := uint64(createData["id"].(float64))

		// 2. Confirm order
		confirmPayload := map[string]interface{}{"note": "Order confirmed"}
		url := fmt.Sprintf("/admin/orders/%d/confirm", orderID)
		w = helper.MakeRequest("POST", url, confirmPayload, true)
		assert.Equal(t, http.StatusOK, w.Code)

		// 3. Start order
		startPayload := map[string]interface{}{"note": "Service started"}
		url = fmt.Sprintf("/admin/orders/%d/start", orderID)
		w = helper.MakeRequest("POST", url, startPayload, true)
		assert.Equal(t, http.StatusOK, w.Code)

		// 4. Complete order
		completePayload := map[string]interface{}{"note": "Service completed"}
		url = fmt.Sprintf("/admin/orders/%d/complete", orderID)
		w = helper.MakeRequest("POST", url, completePayload, true)
		assert.Equal(t, http.StatusOK, w.Code)

		// 5. Get order to verify final status
		url = fmt.Sprintf("/admin/orders/%d", orderID)
		w = helper.MakeRequest("GET", url, nil, true)
		assert.Equal(t, http.StatusOK, w.Code)

		var getResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &getResponse)
		require.NoError(t, err)
		getData := getResponse["data"].(map[string]interface{})
		assert.Equal(t, "completed", getData["status"])
	})
}
