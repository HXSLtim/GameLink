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
	"gamelink/internal/repository"
	withdrawrepo "gamelink/internal/repository/withdraw"
	"gamelink/internal/service/earnings"
)

// ---- Fake Repositories for earnings tests ----

type fakePlayerRepositoryForEarnings struct{}

func (m *fakePlayerRepositoryForEarnings) List(ctx context.Context) ([]model.Player, error) {
	return []model.Player{{Base: model.Base{ID: 1}, UserID: 200, Nickname: "TestPlayer"}}, nil
}
func (m *fakePlayerRepositoryForEarnings) ListPaged(ctx context.Context, page, pageSize int) ([]model.Player, int64, error) {
	return []model.Player{{Base: model.Base{ID: 1}, UserID: 200, Nickname: "TestPlayer"}}, 1, nil
}
func (m *fakePlayerRepositoryForEarnings) Get(ctx context.Context, id uint64) (*model.Player, error) {
	return &model.Player{Base: model.Base{ID: id}, UserID: 200, Nickname: "TestPlayer"}, nil
}
func (m *fakePlayerRepositoryForEarnings) GetByUserID(ctx context.Context, userID uint64) (*model.Player, error) {
	return &model.Player{Base: model.Base{ID: 1}, UserID: userID, Nickname: "TestPlayer"}, nil
}
func (m *fakePlayerRepositoryForEarnings) Create(ctx context.Context, player *model.Player) error {
	return nil
}
func (m *fakePlayerRepositoryForEarnings) Update(ctx context.Context, player *model.Player) error {
	return nil
}
func (m *fakePlayerRepositoryForEarnings) Delete(ctx context.Context, id uint64) error { return nil }
func (m *fakePlayerRepositoryForEarnings) ListByGameID(ctx context.Context, gameID uint64) ([]model.Player, error) {
	return []model.Player{}, nil
}

type fakeOrderRepositoryForEarnings struct {
	orders map[uint64]*model.Order
}

func newFakeOrderRepositoryForEarnings() *fakeOrderRepositoryForEarnings {
	return &fakeOrderRepositoryForEarnings{orders: make(map[uint64]*model.Order)}
}

func (m *fakeOrderRepositoryForEarnings) Create(ctx context.Context, o *model.Order) error {
	o.ID = uint64(len(m.orders) + 1)
	m.orders[o.ID] = o
	return nil
}

func (m *fakeOrderRepositoryForEarnings) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) {
	var res []model.Order
	for _, o := range m.orders {
		res = append(res, *o)
	}
	return res, int64(len(res)), nil
}

func (m *fakeOrderRepositoryForEarnings) Get(ctx context.Context, id uint64) (*model.Order, error) {
	if o, ok := m.orders[id]; ok {
		return o, nil
	}
	return nil, repository.ErrNotFound
}

func (m *fakeOrderRepositoryForEarnings) Update(ctx context.Context, o *model.Order) error {
	m.orders[o.ID] = o
	return nil
}

func (m *fakeOrderRepositoryForEarnings) Delete(ctx context.Context, id uint64) error {
	delete(m.orders, id)
	return nil
}

type fakeWithdrawRepositoryForEarnings struct{}

func (m *fakeWithdrawRepositoryForEarnings) Create(ctx context.Context, withdraw *model.Withdraw) error {
	return nil
}
func (m *fakeWithdrawRepositoryForEarnings) Get(ctx context.Context, id uint64) (*model.Withdraw, error) {
	return &model.Withdraw{}, nil
}
func (m *fakeWithdrawRepositoryForEarnings) List(ctx context.Context, opts withdrawrepo.WithdrawListOptions) ([]model.Withdraw, int64, error) {
	return []model.Withdraw{}, 0, nil
}
func (m *fakeWithdrawRepositoryForEarnings) Update(ctx context.Context, withdraw *model.Withdraw) error {
	return nil
}
func (m *fakeWithdrawRepositoryForEarnings) Delete(ctx context.Context, id uint64) error {
	return nil
}
func (m *fakeWithdrawRepositoryForEarnings) GetPlayerBalance(ctx context.Context, playerID uint64) (*withdrawrepo.PlayerBalance, error) {
	return &withdrawrepo.PlayerBalance{
		TotalEarnings:    15000,
		WithdrawTotal:    0,
		PendingWithdraw:  0,
		AvailableBalance: 15000,
		PendingBalance:   0,
	}, nil
}

// ---- Tests for player_earnings.go ----

func TestGetEarningsSummaryHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	earningsSvc := earnings.NewEarningsService(&fakePlayerRepositoryForEarnings{}, newFakeOrderRepositoryForEarnings(), &fakeWithdrawRepositoryForEarnings{})

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

	earningsSvc := earnings.NewEarningsService(&fakePlayerRepositoryForEarnings{}, newFakeOrderRepositoryForEarnings(), &fakeWithdrawRepositoryForEarnings{})

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

	earningsSvc := earnings.NewEarningsService(&fakePlayerRepositoryForEarnings{}, newFakeOrderRepositoryForEarnings(), &fakeWithdrawRepositoryForEarnings{})

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

	orderRepo := newFakeOrderRepositoryForEarnings()
	// Create some completed orders to give the player earnings
	for i := 0; i < 3; i++ {
		playerID := uint64(1) // Player with UserID 200
		gameID := uint64(1)
		order := &model.Order{
			UserID:          100 + uint64(i),
			PlayerID:        &playerID,
			GameID:          &gameID,
			ItemID:          1,
			OrderNo:         "EARNINGS-TEST-" + string(rune('A'+i)),
			Status:          model.OrderStatusCompleted,
			TotalPriceCents: 5000, // Total: 15000 cents
			UnitPriceCents:  5000,
		}
		orderRepo.Create(context.Background(), order)
	}

	earningsSvc := earnings.NewEarningsService(&fakePlayerRepositoryForEarnings{}, orderRepo, &fakeWithdrawRepositoryForEarnings{})

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

	earningsSvc := earnings.NewEarningsService(&fakePlayerRepositoryForEarnings{}, newFakeOrderRepositoryForEarnings(), &fakeWithdrawRepositoryForEarnings{})

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

	earningsSvc := earnings.NewEarningsService(&fakePlayerRepositoryForEarnings{}, newFakeOrderRepositoryForEarnings(), &fakeWithdrawRepositoryForEarnings{})

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
