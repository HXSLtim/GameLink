package player

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/service/earnings"
)

// ---- Fake Repositories for earnings tests ----
// Note: We reuse the fake repositories from user_order_test.go for player and order repositories

// ---- Tests for player_earnings.go ----

func TestGetEarningsSummaryHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	earningsSvc := earnings.NewEarningsService(&fakePlayerRepository{}, newFakeOrderRepository())

	router := gin.New()
	router.GET("/player/earnings/summary", func(c *gin.Context) {
		c.Set("user_id", uint64(200))
		getEarningsSummaryHandler(c, earningsSvc)
	})

	req := httptest.NewRequest(http.MethodGet, "/player/earnings/summary", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp model.APIResponse[earnings.EarningsSummaryResponse]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !resp.Success {
		t.Fatalf("Expected success=true, got %v", resp.Success)
	}
}

func TestGetEarningsTrendHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	earningsSvc := earnings.NewEarningsService(&fakePlayerRepository{}, newFakeOrderRepository())

	router := gin.New()
	router.GET("/player/earnings/trend", func(c *gin.Context) {
		c.Set("user_id", uint64(200))
		getEarningsTrendHandler(c, earningsSvc)
	})

	req := httptest.NewRequest(http.MethodGet, "/player/earnings/trend?days=7", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp model.APIResponse[earnings.EarningsTrendResponse]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !resp.Success {
		t.Fatalf("Expected success=true, got %v", resp.Success)
	}
}

func TestGetEarningsTrendHandler_InvalidDays(t *testing.T) {
	gin.SetMode(gin.TestMode)

	earningsSvc := earnings.NewEarningsService(&fakePlayerRepository{}, newFakeOrderRepository())

	router := gin.New()
	router.GET("/player/earnings/trend", func(c *gin.Context) {
		c.Set("user_id", uint64(200))
		getEarningsTrendHandler(c, earningsSvc)
	})

	// Test invalid days (too small)
	req := httptest.NewRequest(http.MethodGet, "/player/earnings/trend?days=5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400 for days=5, got %d", w.Code)
	}

	// Test invalid days (too large)
	req = httptest.NewRequest(http.MethodGet, "/player/earnings/trend?days=100", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400 for days=100, got %d", w.Code)
	}

	// Test invalid days (not a number)
	req = httptest.NewRequest(http.MethodGet, "/player/earnings/trend?days=abc", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400 for days=abc, got %d", w.Code)
	}
}

func TestRequestWithdrawHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	orderRepo := newFakeOrderRepository()
	// Create some completed orders to give the player earnings
	for i := 0; i < 3; i++ {
		order := &model.Order{
			UserID:     100 + uint64(i),
			PlayerID:   1, // Player with UserID 200
			Status:     model.OrderStatusCompleted,
			PriceCents: 5000, // Total: 15000 cents
			GameID:     1,
		}
		orderRepo.Create(context.Background(), order)
	}

	earningsSvc := earnings.NewEarningsService(&fakePlayerRepository{}, orderRepo)

	router := gin.New()
	router.POST("/player/earnings/withdraw", func(c *gin.Context) {
		c.Set("user_id", uint64(200))
		requestWithdrawHandler(c, earningsSvc)
	})

	reqBody := earnings.WithdrawRequest{
		AmountCents: 10000,
		Method:      "alipay",
		AccountInfo: "test@example.com",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/player/earnings/withdraw", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestRequestWithdrawHandler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	earningsSvc := earnings.NewEarningsService(&fakePlayerRepository{}, newFakeOrderRepository())

	router := gin.New()
	router.POST("/player/earnings/withdraw", func(c *gin.Context) {
		c.Set("user_id", uint64(200))
		requestWithdrawHandler(c, earningsSvc)
	})

	req := httptest.NewRequest(http.MethodPost, "/player/earnings/withdraw", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d", w.Code)
	}
}

func TestGetWithdrawHistoryHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	earningsSvc := earnings.NewEarningsService(&fakePlayerRepository{}, newFakeOrderRepository())

	router := gin.New()
	router.GET("/player/earnings/withdraw-history", func(c *gin.Context) {
		c.Set("user_id", uint64(200))
		getWithdrawHistoryHandler(c, earningsSvc)
	})

	req := httptest.NewRequest(http.MethodGet, "/player/earnings/withdraw-history?page=1&pageSize=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp model.APIResponse[earnings.WithdrawHistoryResponse]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !resp.Success {
		t.Fatalf("Expected success=true, got %v", resp.Success)
	}
}
