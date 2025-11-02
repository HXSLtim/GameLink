package player

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/service/order"
)

// ---- Fake OrderRepository for player_order tests ----

type mockOrderRepoForPlayerOrder struct {
	orders map[uint64]*model.Order
}

func newMockOrderRepoForPlayerOrder() *mockOrderRepoForPlayerOrder {
	return &mockOrderRepoForPlayerOrder{
		orders: map[uint64]*model.Order{
			1: {Base: model.Base{ID: 1}, UserID: 100, GameID: 10, Status: model.OrderStatusConfirmed, PriceCents: 5000},
			2: {Base: model.Base{ID: 2}, UserID: 101, GameID: 10, Status: model.OrderStatusPending, PriceCents: 8000},
			3: {Base: model.Base{ID: 3}, UserID: 102, PlayerID: 1, GameID: 20, Status: model.OrderStatusInProgress, PriceCents: 3000},
		},
	}
}

func (m *mockOrderRepoForPlayerOrder) Create(ctx context.Context, o *model.Order) error {
	o.ID = uint64(len(m.orders) + 1)
	m.orders[o.ID] = o
	return nil
}

func (m *mockOrderRepoForPlayerOrder) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) {
	var result []model.Order
	for _, o := range m.orders {
		// Filter by player if specified
		if opts.PlayerID != nil && *opts.PlayerID != o.PlayerID {
			continue
		}
		// Filter by game if specified
		if opts.GameID != nil && *opts.GameID != o.GameID {
			continue
		}
		result = append(result, *o)
	}
	return result, int64(len(result)), nil
}

func (m *mockOrderRepoForPlayerOrder) Get(ctx context.Context, id uint64) (*model.Order, error) {
	if o, ok := m.orders[id]; ok {
		return o, nil
	}
	return nil, repository.ErrNotFound
}

func (m *mockOrderRepoForPlayerOrder) Update(ctx context.Context, o *model.Order) error {
	m.orders[o.ID] = o
	return nil
}

func (m *mockOrderRepoForPlayerOrder) Delete(ctx context.Context, id uint64) error {
	delete(m.orders, id)
	return nil
}

// ---- Fake PlayerRepository for order tests ----

type fakePlayerRepositoryForOrder struct{}

func (m *fakePlayerRepositoryForOrder) List(ctx context.Context) ([]model.Player, error) {
	return []model.Player{}, nil
}

func (m *fakePlayerRepositoryForOrder) ListPaged(ctx context.Context, page, pageSize int) ([]model.Player, int64, error) {
	return []model.Player{
		{Base: model.Base{ID: 1}, UserID: 200, Nickname: "TestPlayer"},
	}, 1, nil
}

func (m *fakePlayerRepositoryForOrder) Get(ctx context.Context, id uint64) (*model.Player, error) {
	return &model.Player{Base: model.Base{ID: id}, UserID: 200, Nickname: "TestPlayer"}, nil
}

func (m *fakePlayerRepositoryForOrder) GetByUserID(ctx context.Context, userID uint64) (*model.Player, error) {
	return &model.Player{Base: model.Base{ID: 1}, UserID: userID, Nickname: "TestPlayer"}, nil
}

func (m *fakePlayerRepositoryForOrder) Create(ctx context.Context, player *model.Player) error {
	return nil
}

func (m *fakePlayerRepositoryForOrder) Update(ctx context.Context, player *model.Player) error {
	return nil
}

func (m *fakePlayerRepositoryForOrder) Delete(ctx context.Context, id uint64) error {
	return nil
}

func (m *fakePlayerRepositoryForOrder) ListByGameID(ctx context.Context, gameID uint64) ([]model.Player, error) {
	return []model.Player{}, nil
}

// ---- Tests for player_order.go ----

func TestGetAvailableOrdersHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	orderRepo := newMockOrderRepoForPlayerOrder()
	orderSvc := order.NewOrderService(orderRepo, &fakePlayerRepositoryForOrder{}, &fakeUserRepository{}, &fakeGameRepository{}, &fakePaymentRepository{}, &fakeReviewRepository{})

	router := gin.New()
	router.GET("/player/orders/available", func(c *gin.Context) {
		c.Set("user_id", uint64(200))
		getAvailableOrdersHandler(c, orderSvc)
	})

	req := httptest.NewRequest(http.MethodGet, "/player/orders/available?page=1&pageSize=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp model.APIResponse[map[string]any]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !resp.Success {
		t.Fatalf("Expected success=true, got %v", resp.Success)
	}
}

func TestGetAvailableOrdersHandler_WithFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)

	orderRepo := newMockOrderRepoForPlayerOrder()
	orderSvc := order.NewOrderService(orderRepo, &fakePlayerRepositoryForOrder{}, &fakeUserRepository{}, &fakeGameRepository{}, &fakePaymentRepository{}, &fakeReviewRepository{})

	router := gin.New()
	router.GET("/player/orders/available", func(c *gin.Context) {
		c.Set("user_id", uint64(200))
		getAvailableOrdersHandler(c, orderSvc)
	})

	req := httptest.NewRequest(http.MethodGet, "/player/orders/available?gameId=10&page=1&pageSize=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
}

func TestAcceptOrderHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	orderRepo := newMockOrderRepoForPlayerOrder()
	// Order 1 has status Pending and no PlayerID, so it can be accepted
	orderSvc := order.NewOrderService(orderRepo, &fakePlayerRepositoryForOrder{}, &fakeUserRepository{}, &fakeGameRepository{}, &fakePaymentRepository{}, &fakeReviewRepository{})

	router := gin.New()
	router.POST("/player/orders/:id/accept", func(c *gin.Context) {
		c.Set("user_id", uint64(200))
		acceptOrderHandler(c, orderSvc)
	})

	// Use order 1 which has status OrderStatusPending
	req := httptest.NewRequest(http.MethodPost, "/player/orders/1/accept", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAcceptOrderHandler_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	orderRepo := newMockOrderRepoForPlayerOrder()
	orderSvc := order.NewOrderService(orderRepo, &fakePlayerRepositoryForOrder{}, &fakeUserRepository{}, &fakeGameRepository{}, &fakePaymentRepository{}, &fakeReviewRepository{})

	router := gin.New()
	router.POST("/player/orders/:id/accept", func(c *gin.Context) {
		c.Set("user_id", uint64(200))
		acceptOrderHandler(c, orderSvc)
	})

	req := httptest.NewRequest(http.MethodPost, "/player/orders/invalid/accept", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d", w.Code)
	}
}

func TestGetMyAcceptedOrdersHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	orderRepo := newMockOrderRepoForPlayerOrder()
	orderSvc := order.NewOrderService(orderRepo, &fakePlayerRepositoryForOrder{}, &fakeUserRepository{}, &fakeGameRepository{}, &fakePaymentRepository{}, &fakeReviewRepository{})

	router := gin.New()
	router.GET("/player/orders/my", func(c *gin.Context) {
		c.Set("user_id", uint64(200))
		getMyAcceptedOrdersHandler(c, orderSvc)
	})

	req := httptest.NewRequest(http.MethodGet, "/player/orders/my?page=1&pageSize=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp model.APIResponse[order.MyOrderListResponse]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !resp.Success {
		t.Fatalf("Expected success=true, got %v", resp.Success)
	}
}

func TestGetMyAcceptedOrdersHandler_WithStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	orderRepo := newMockOrderRepoForPlayerOrder()
	orderSvc := order.NewOrderService(orderRepo, &fakePlayerRepositoryForOrder{}, &fakeUserRepository{}, &fakeGameRepository{}, &fakePaymentRepository{}, &fakeReviewRepository{})

	router := gin.New()
	router.GET("/player/orders/my", func(c *gin.Context) {
		c.Set("user_id", uint64(200))
		getMyAcceptedOrdersHandler(c, orderSvc)
	})

	req := httptest.NewRequest(http.MethodGet, "/player/orders/my?status=confirmed&page=1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
}

func TestCompleteOrderByPlayerHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	orderRepo := newMockOrderRepoForPlayerOrder()
	// Order 3 has PlayerID=200 and status=Confirmed, owned by player with UserID=200
	orderSvc := order.NewOrderService(orderRepo, &fakePlayerRepositoryForOrder{}, &fakeUserRepository{}, &fakeGameRepository{}, &fakePaymentRepository{}, &fakeReviewRepository{})

	router := gin.New()
	router.PUT("/player/orders/:id/complete", func(c *gin.Context) {
		c.Set("user_id", uint64(200))
		completeOrderByPlayerHandler(c, orderSvc)
	})

	// Use order 3 which has PlayerID=200 and status=OrderStatusConfirmed
	req := httptest.NewRequest(http.MethodPut, "/player/orders/3/complete", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCompleteOrderByPlayerHandler_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	orderRepo := newMockOrderRepoForPlayerOrder()
	orderSvc := order.NewOrderService(orderRepo, &fakePlayerRepositoryForOrder{}, &fakeUserRepository{}, &fakeGameRepository{}, &fakePaymentRepository{}, &fakeReviewRepository{})

	router := gin.New()
	router.PUT("/player/orders/:id/complete", func(c *gin.Context) {
		c.Set("user_id", uint64(200))
		completeOrderByPlayerHandler(c, orderSvc)
	})

	req := httptest.NewRequest(http.MethodPut, "/player/orders/invalid/complete", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d", w.Code)
	}
}
