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
	commissionrepo "gamelink/internal/repository/commission"
	"gamelink/internal/service/order"
)

// ---- Fake OrderRepository for player_order tests ----

type mockOrderRepoForPlayerOrder struct {
	orders map[uint64]*model.Order
}

func newMockOrderRepoForPlayerOrder() *mockOrderRepoForPlayerOrder {
	gameID1 := uint64(10)
	gameID2 := uint64(20)
	playerID := uint64(1)
	return &mockOrderRepoForPlayerOrder{
		orders: map[uint64]*model.Order{
			1: {Base: model.Base{ID: 1}, UserID: 100, GameID: &gameID1, ItemID: 1, OrderNo: "TEST-001", Status: model.OrderStatusConfirmed, TotalPriceCents: 5000, UnitPriceCents: 5000},
			2: {Base: model.Base{ID: 2}, UserID: 101, GameID: &gameID1, ItemID: 1, OrderNo: "TEST-002", Status: model.OrderStatusPending, TotalPriceCents: 8000, UnitPriceCents: 8000},
			3: {Base: model.Base{ID: 3}, UserID: 102, PlayerID: &playerID, GameID: &gameID2, ItemID: 1, OrderNo: "TEST-003", Status: model.OrderStatusInProgress, TotalPriceCents: 3000, UnitPriceCents: 3000},
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
		if opts.PlayerID != nil {
			if o.PlayerID == nil || *opts.PlayerID != *o.PlayerID {
				continue
			}
		}
		// Filter by game if specified
		if opts.GameID != nil {
			if o.GameID == nil || *opts.GameID != *o.GameID {
				continue
			}
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

type fakeUserRepository struct{}

func (m *fakeUserRepository) List(ctx context.Context) ([]model.User, error) {
	return []model.User{}, nil
}
func (m *fakeUserRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	return []model.User{}, 0, nil
}
func (m *fakeUserRepository) ListWithFilters(ctx context.Context, opts repository.UserListOptions) ([]model.User, int64, error) {
	return []model.User{}, 0, nil
}
func (m *fakeUserRepository) Get(ctx context.Context, id uint64) (*model.User, error) {
	return &model.User{Base: model.Base{ID: id}, Name: "TestUser"}, nil
}
func (m *fakeUserRepository) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	return nil, nil
}
func (m *fakeUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return nil, nil
}
func (m *fakeUserRepository) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	return nil, nil
}
func (m *fakeUserRepository) Create(ctx context.Context, user *model.User) error { return nil }
func (m *fakeUserRepository) Update(ctx context.Context, user *model.User) error { return nil }
func (m *fakeUserRepository) Delete(ctx context.Context, id uint64) error        { return nil }

type fakeGameRepository struct{}

func (m *fakeGameRepository) List(ctx context.Context) ([]model.Game, error) {
	return []model.Game{}, nil
}
func (m *fakeGameRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.Game, int64, error) {
	return []model.Game{}, 0, nil
}
func (m *fakeGameRepository) Get(ctx context.Context, id uint64) (*model.Game, error) {
	return &model.Game{Base: model.Base{ID: id}, Name: "TestGame"}, nil
}
func (m *fakeGameRepository) Create(ctx context.Context, game *model.Game) error { return nil }
func (m *fakeGameRepository) Update(ctx context.Context, game *model.Game) error { return nil }
func (m *fakeGameRepository) Delete(ctx context.Context, id uint64) error        { return nil }

type fakePaymentRepository struct{}

func (m *fakePaymentRepository) Create(ctx context.Context, payment *model.Payment) error { return nil }
func (m *fakePaymentRepository) List(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
	return []model.Payment{}, 0, nil
}
func (m *fakePaymentRepository) Get(ctx context.Context, id uint64) (*model.Payment, error) {
	return &model.Payment{}, nil
}
func (m *fakePaymentRepository) Update(ctx context.Context, payment *model.Payment) error { return nil }
func (m *fakePaymentRepository) Delete(ctx context.Context, id uint64) error              { return nil }

type fakeReviewRepository struct{}

func (m *fakeReviewRepository) List(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error) {
	return []model.Review{}, 0, nil
}
func (m *fakeReviewRepository) Get(ctx context.Context, id uint64) (*model.Review, error) {
	return &model.Review{}, nil
}
func (m *fakeReviewRepository) GetByOrderID(ctx context.Context, orderID uint64) (*model.Review, error) {
	return nil, repository.ErrNotFound
}
func (m *fakeReviewRepository) Create(ctx context.Context, review *model.Review) error { return nil }
func (m *fakeReviewRepository) Update(ctx context.Context, review *model.Review) error { return nil }
func (m *fakeReviewRepository) Delete(ctx context.Context, id uint64) error            { return nil }

type fakeCommissionRepositoryForOrder struct{}

// CommissionRule methods
func (m *fakeCommissionRepositoryForOrder) CreateRule(ctx context.Context, rule *model.CommissionRule) error {
	return nil
}
func (m *fakeCommissionRepositoryForOrder) GetRule(ctx context.Context, id uint64) (*model.CommissionRule, error) {
	return &model.CommissionRule{ID: 1, Rate: 20}, nil
}
func (m *fakeCommissionRepositoryForOrder) GetDefaultRule(ctx context.Context) (*model.CommissionRule, error) {
	return &model.CommissionRule{ID: 1, Rate: 20, Type: "default", IsActive: true}, nil
}
func (m *fakeCommissionRepositoryForOrder) GetRuleForOrder(ctx context.Context, gameID *uint64, playerID *uint64, serviceType *string) (*model.CommissionRule, error) {
	return &model.CommissionRule{Rate: 20}, nil
}
func (m *fakeCommissionRepositoryForOrder) ListRules(ctx context.Context, opts commissionrepo.CommissionRuleListOptions) ([]model.CommissionRule, int64, error) {
	return []model.CommissionRule{}, 0, nil
}
func (m *fakeCommissionRepositoryForOrder) UpdateRule(ctx context.Context, rule *model.CommissionRule) error {
	return nil
}
func (m *fakeCommissionRepositoryForOrder) DeleteRule(ctx context.Context, id uint64) error {
	return nil
}

// CommissionRecord methods
func (m *fakeCommissionRepositoryForOrder) CreateRecord(ctx context.Context, record *model.CommissionRecord) error {
	return nil
}
func (m *fakeCommissionRepositoryForOrder) GetRecord(ctx context.Context, id uint64) (*model.CommissionRecord, error) {
	return &model.CommissionRecord{}, nil
}
func (m *fakeCommissionRepositoryForOrder) GetRecordByOrderID(ctx context.Context, orderID uint64) (*model.CommissionRecord, error) {
	return &model.CommissionRecord{}, nil
}
func (m *fakeCommissionRepositoryForOrder) ListRecords(ctx context.Context, opts commissionrepo.CommissionRecordListOptions) ([]model.CommissionRecord, int64, error) {
	return []model.CommissionRecord{}, 0, nil
}
func (m *fakeCommissionRepositoryForOrder) UpdateRecord(ctx context.Context, record *model.CommissionRecord) error {
	return nil
}

// MonthlySettlement methods
func (m *fakeCommissionRepositoryForOrder) CreateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error {
	return nil
}
func (m *fakeCommissionRepositoryForOrder) GetSettlement(ctx context.Context, id uint64) (*model.MonthlySettlement, error) {
	return &model.MonthlySettlement{}, nil
}
func (m *fakeCommissionRepositoryForOrder) GetSettlementByPlayerMonth(ctx context.Context, playerID uint64, month string) (*model.MonthlySettlement, error) {
	return &model.MonthlySettlement{}, nil
}
func (m *fakeCommissionRepositoryForOrder) ListSettlements(ctx context.Context, opts commissionrepo.SettlementListOptions) ([]model.MonthlySettlement, int64, error) {
	return []model.MonthlySettlement{}, 0, nil
}
func (m *fakeCommissionRepositoryForOrder) UpdateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error {
	return nil
}

// Stats methods
func (m *fakeCommissionRepositoryForOrder) GetMonthlyStats(ctx context.Context, month string) (*commissionrepo.MonthlyStats, error) {
	return &commissionrepo.MonthlyStats{}, nil
}
func (m *fakeCommissionRepositoryForOrder) GetPlayerMonthlyIncome(ctx context.Context, playerID uint64, month string) (int64, error) {
	return 0, nil
}

// ---- Tests for player_order.go ----

func TestGetAvailableOrdersHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	orderRepo := newMockOrderRepoForPlayerOrder()
	orderSvc := order.NewOrderService(orderRepo, &fakePlayerRepositoryForOrder{}, &fakeUserRepository{}, &fakeGameRepository{}, &fakePaymentRepository{}, &fakeReviewRepository{}, &fakeCommissionRepositoryForOrder{})

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
	orderSvc := order.NewOrderService(orderRepo, &fakePlayerRepositoryForOrder{}, &fakeUserRepository{}, &fakeGameRepository{}, &fakePaymentRepository{}, &fakeReviewRepository{}, &fakeCommissionRepositoryForOrder{})

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
	orderSvc := order.NewOrderService(orderRepo, &fakePlayerRepositoryForOrder{}, &fakeUserRepository{}, &fakeGameRepository{}, &fakePaymentRepository{}, &fakeReviewRepository{}, &fakeCommissionRepositoryForOrder{})

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
	orderSvc := order.NewOrderService(orderRepo, &fakePlayerRepositoryForOrder{}, &fakeUserRepository{}, &fakeGameRepository{}, &fakePaymentRepository{}, &fakeReviewRepository{}, &fakeCommissionRepositoryForOrder{})

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
	orderSvc := order.NewOrderService(orderRepo, &fakePlayerRepositoryForOrder{}, &fakeUserRepository{}, &fakeGameRepository{}, &fakePaymentRepository{}, &fakeReviewRepository{}, &fakeCommissionRepositoryForOrder{})

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
	orderSvc := order.NewOrderService(orderRepo, &fakePlayerRepositoryForOrder{}, &fakeUserRepository{}, &fakeGameRepository{}, &fakePaymentRepository{}, &fakeReviewRepository{}, &fakeCommissionRepositoryForOrder{})

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
	orderSvc := order.NewOrderService(orderRepo, &fakePlayerRepositoryForOrder{}, &fakeUserRepository{}, &fakeGameRepository{}, &fakePaymentRepository{}, &fakeReviewRepository{}, &fakeCommissionRepositoryForOrder{})

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
	orderSvc := order.NewOrderService(orderRepo, &fakePlayerRepositoryForOrder{}, &fakeUserRepository{}, &fakeGameRepository{}, &fakePaymentRepository{}, &fakeReviewRepository{}, &fakeCommissionRepositoryForOrder{})

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
