//go:build ignore

/**
 * @file payment batch integration test
 * @description Integration tests for payment batch operations
 * @note This test file is currently disabled due to missing TestHelper implementation
 * TODO: Rewrite this test to use SetupAdminTest like other integration tests
 */

package admin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/handler/middleware"
	"gamelink/internal/model"
	"gamelink/internal/repository/gamecategory"
	"gamelink/internal/service/admin"
	"gamelink/internal/service/integration"
	"gamelink/pkg/auth"
)

// setupPaymentBatchTest sets up test environment for payment batch operations
func setupPaymentBatchTest(t *testing.T) (*httptest.Server, *integration.TestHelper, func()) {
	t.Helper()

	helper := integration.NewTestHelper(t)
	helper.Setup()

	// Create test data
	user := integration.CreateUniqueTestUser(t, helper.DB, "batch_user")
	player := integration.CreateUniqueTestUser(t, helper.DB, "batch_player")
	playerModel := integration.CreateTestPlayer(t, helper.DB, player.ID)
	game := integration.CreateTestGame(t, helper.DB, "batch_game")
	serviceItem := integration.CreateTestServiceItem(t, helper.DB, game.ID)

	// Create orders and payments
	pendingPaymentIDs := make([]uint64, 0)
	paidPaymentIDs := make([]uint64, 0)

	for i := 0; i < 3; i++ {
		order := integration.CreateTestOrderForPlayer(t, helper.DB, user.ID, playerModel.ID, game.ID, serviceItem.ID)
		payment := integration.CreateTestPayment(t, helper.DB, order.ID, user.ID, model.PaymentStatusPending)
		pendingPaymentIDs = append(pendingPaymentIDs, payment.ID)
	}

	for i := 0; i < 2; i++ {
		order := integration.CreateTestOrderForPlayer(t, helper.DB, user.ID, playerModel.ID, game.ID, serviceItem.ID)
		payment := integration.CreateTestPayment(t, helper.DB, order.ID, user.ID, model.PaymentStatusPaid)
		paidPaymentIDs = append(paidPaymentIDs, payment.ID)
	}

	// Setup router
	adminSvc := admin.NewAdminService(helper.DB, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, gamecategory.NewGameCategoryRepository(helper.DB), nil)
	jwtService := auth.NewJWTService("test-secret-key-12345678901234567890", time.Hour)
	pm := middleware.NewPermissionMiddleware(adminSvc, jwtService)

	router := setupTestRouter(adminSvc, pm)

	server := httptest.NewServer(router)

	cleanup := func() {
		server.Close()
		helper.TearDown()
	}

	return server, helper, cleanup
}

func TestPaymentBatch_BatchCapture(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	integration.SkipIfNoTestDB(t)

	server, _, cleanup := setupPaymentBatchTest(t)
	defer cleanup()

	// Get pending payment IDs from test setup
	paymentIDs := []uint64{1, 2, 3} // These would be the actual IDs from setup

	t.Run("capture pending payments successfully", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"paymentIds": paymentIDs,
		}
		body, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", server.URL+"/admin/payments/batch/capture", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+getTestAdminToken(t))

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.True(t, result["success"].(bool))
		assert.True(t, result["successCount"].(float64) > 0)
	})

	t.Run("capture with non-existent payment IDs", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"paymentIds": []uint64{99999, 99998},
		}
		body, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", server.URL+"/admin/payments/batch/capture", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+getTestAdminToken(t))

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.True(t, result["success"].(bool))
		assert.Equal(t, float64(0), result["successCount"].(float64))
		assert.Equal(t, float64(2), result["failedCount"].(float64))
	})

	t.Run("capture with invalid status payments", func(t *testing.T) {
		paidIDs := []uint64{4, 5} // Paid payment IDs
		reqBody := map[string]interface{}{
			"paymentIds": paidIDs,
		}
		body, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", server.URL+"/admin/payments/batch/capture", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+getTestAdminToken(t))

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.True(t, result["success"].(bool))
		assert.Equal(t, float64(0), result["successCount"].(float64))
		assert.Equal(t, float64(2), result["failedCount"].(float64))
	})

	t.Run("capture with empty payment IDs", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"paymentIds": []uint64{},
		}
		body, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", server.URL+"/admin/payments/batch/capture", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+getTestAdminToken(t))

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestPaymentBatch_BatchRefund(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	integration.SkipIfNoTestDB(t)

	server, _, cleanup := setupPaymentBatchTest(t)
	defer cleanup()

	paidIDs := []uint64{4, 5} // Paid payment IDs

	t.Run("refund paid payments successfully", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"paymentIds": paidIDs,
			"reason":     "用户主动退款",
		}
		body, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", server.URL+"/admin/payments/batch/refund", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+getTestAdminToken(t))

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.True(t, result["success"].(bool))
		assert.True(t, result["successCount"].(float64) > 0)
	})

	t.Run("refund with missing reason", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"paymentIds": paidIDs,
		}
		body, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", server.URL+"/admin/payments/batch/refund", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+getTestAdminToken(t))

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("refund pending payments (should fail)", func(t *testing.T) {
		pendingIDs := []uint64{1, 2, 3} // Pending payment IDs
		reqBody := map[string]interface{}{
			"paymentIds": pendingIDs,
			"reason":     "test refund",
		}
		body, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", server.URL+"/admin/payments/batch/refund", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+getTestAdminToken(t))

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.True(t, result["success"].(bool))
		assert.Equal(t, float64(0), result["successCount"].(float64))
		assert.Equal(t, float64(3), result["failedCount"].(float64))
	})
}

func TestPaymentBatch_BatchCancel(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	integration.SkipIfNoTestDB(t)

	server, _, cleanup := setupPaymentBatchTest(t)
	defer cleanup()

	pendingIDs := []uint64{1, 2, 3} // Pending payment IDs

	t.Run("cancel pending payments successfully", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"paymentIds": pendingIDs,
		}
		body, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", server.URL+"/admin/payments/batch/cancel", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+getTestAdminToken(t))

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.True(t, result["success"].(bool))
		assert.True(t, result["successCount"].(float64) > 0)
	})

	t.Run("cancel paid payments (should fail)", func(t *testing.T) {
		paidIDs := []uint64{4, 5} // Paid payment IDs
		reqBody := map[string]interface{}{
			"paymentIds": paidIDs,
		}
		body, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", server.URL+"/admin/payments/batch/cancel", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+getTestAdminToken(t))

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.True(t, result["success"].(bool))
		assert.Equal(t, float64(0), result["successCount"].(float64))
		assert.Equal(t, float64(2), result["failedCount"].(float64))
	})
}

func TestPaymentBatch_BatchUpdateStatus(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	integration.SkipIfNoTestDB(t)

	server, _, cleanup := setupPaymentBatchTest(t)
	defer cleanup()

	t.Run("update status to paid", func(t *testing.T) {
		pendingIDs := []uint64{1, 2, 3} // Pending payment IDs
		reqBody := map[string]interface{}{
			"paymentIds": pendingIDs,
			"status":     "paid",
		}
		body, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("PUT", server.URL+"/admin/payments/batch/status", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+getTestAdminToken(t))

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.True(t, result["success"].(bool))
		assert.True(t, result["successCount"].(float64) > 0)
	})

	t.Run("update status with invalid transition", func(t *testing.T) {
		paidIDs := []uint64{4, 5} // Paid payment IDs - can't go to pending
		reqBody := map[string]interface{}{
			"paymentIds": paidIDs,
			"status":     "pending",
		}
		body, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("PUT", server.URL+"/admin/payments/batch/status", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+getTestAdminToken(t))

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.True(t, result["success"].(bool))
		assert.Equal(t, float64(0), result["successCount"].(float64))
		assert.Equal(t, float64(2), result["failedCount"].(float64))
	})

	t.Run("update status with invalid status value", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"paymentIds": []uint64{1, 2},
			"status":     "invalid_status",
		}
		body, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("PUT", server.URL+"/admin/payments/batch/status", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+getTestAdminToken(t))

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

// Helper functions

func setupTestRouter(adminSvc *admin.AdminService, pm *middleware.PermissionMiddleware) *gin.Engine {
	router := gin.New()
	RegisterRoutes(router, adminSvc, nil, pm)
	return router
}

func getTestAdminToken(t *testing.T) string {
	t.Helper()
	// Generate a test JWT token for admin user
	// In a real test environment, you'd create an admin user and get actual token
	// For now, return a mock token format
	return "mock_admin_token"
}
