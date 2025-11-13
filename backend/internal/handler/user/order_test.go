package user

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
        commissionrepo "gamelink/internal/repository/commission"
        assignmentservice "gamelink/internal/service/assignment"
        "gamelink/internal/service/order"
)

// ---- Fake repositories for handler tests ----

type fakeOrderRepository struct{ orders map[uint64]*model.Order }

func newFakeOrderRepository() *fakeOrderRepository {
	return &fakeOrderRepository{orders: make(map[uint64]*model.Order)}
}

func (m *fakeOrderRepository) Create(ctx context.Context, o *model.Order) error {
	o.ID = uint64(len(m.orders) + 1)
	m.orders[o.ID] = o
	return nil
}

func (m *fakeOrderRepository) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) {
	// Minimal implementation sufficient for tests here
	var res []model.Order
	for _, o := range m.orders {
		res = append(res, *o)
	}
	return res, int64(len(res)), nil
}

func (m *fakeOrderRepository) Get(ctx context.Context, id uint64) (*model.Order, error) {
	if o, ok := m.orders[id]; ok {
		return o, nil
	}
	return nil, repository.ErrNotFound
}

func (m *fakeOrderRepository) Update(ctx context.Context, o *model.Order) error {
	m.orders[o.ID] = o
	return nil
}

func (m *fakeOrderRepository) Delete(ctx context.Context, id uint64) error {
	delete(m.orders, id)
	return nil
}

type fakePlayerRepository struct{}

func (m *fakePlayerRepository) List(ctx context.Context) ([]model.Player, error) {
	return []model.Player{
		{Base: model.Base{ID: 1}, UserID: 200, Nickname: "TestPlayer"},
	}, nil
}
func (m *fakePlayerRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.Player, int64, error) {
	players := []model.Player{
		{Base: model.Base{ID: 1}, UserID: 200, Nickname: "TestPlayer"},
	}
	return players, int64(len(players)), nil
}
func (m *fakePlayerRepository) Get(ctx context.Context, id uint64) (*model.Player, error) {
	return &model.Player{Base: model.Base{ID: id}, UserID: 1, Nickname: "TestPlayer"}, nil
}
func (m *fakePlayerRepository) GetByUserID(ctx context.Context, userID uint64) (*model.Player, error) {
	return &model.Player{Base: model.Base{ID: 1}, UserID: userID, Nickname: "TestPlayer"}, nil
}
func (m *fakePlayerRepository) Create(ctx context.Context, player *model.Player) error { return nil }
func (m *fakePlayerRepository) Update(ctx context.Context, player *model.Player) error { return nil }
func (m *fakePlayerRepository) Delete(ctx context.Context, id uint64) error            { return nil }
func (m *fakePlayerRepository) ListByGameID(ctx context.Context, gameID uint64) ([]model.Player, error) {
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
	return &model.User{Base: model.Base{ID: id}, AvatarURL: "http://avatar.test"}, nil
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

type fakeCommissionRepository struct{}

// CommissionRule methods
func (m *fakeCommissionRepository) CreateRule(ctx context.Context, rule *model.CommissionRule) error {
	return nil
}
func (m *fakeCommissionRepository) GetRule(ctx context.Context, id uint64) (*model.CommissionRule, error) {
	return &model.CommissionRule{ID: 1, Rate: 20}, nil
}
func (m *fakeCommissionRepository) GetDefaultRule(ctx context.Context) (*model.CommissionRule, error) {
	return &model.CommissionRule{ID: 1, Rate: 20, Type: "default", IsActive: true}, nil
}
func (m *fakeCommissionRepository) GetRuleForOrder(ctx context.Context, gameID *uint64, playerID *uint64, serviceType *string) (*model.CommissionRule, error) {
	return &model.CommissionRule{Rate: 20}, nil
}
func (m *fakeCommissionRepository) ListRules(ctx context.Context, opts commissionrepo.CommissionRuleListOptions) ([]model.CommissionRule, int64, error) {
	return []model.CommissionRule{}, 0, nil
}
func (m *fakeCommissionRepository) UpdateRule(ctx context.Context, rule *model.CommissionRule) error {
	return nil
}
func (m *fakeCommissionRepository) DeleteRule(ctx context.Context, id uint64) error {
	return nil
}

// CommissionRecord methods
func (m *fakeCommissionRepository) CreateRecord(ctx context.Context, record *model.CommissionRecord) error {
	return nil
}
func (m *fakeCommissionRepository) GetRecord(ctx context.Context, id uint64) (*model.CommissionRecord, error) {
	return &model.CommissionRecord{}, nil
}
func (m *fakeCommissionRepository) GetRecordByOrderID(ctx context.Context, orderID uint64) (*model.CommissionRecord, error) {
	return &model.CommissionRecord{}, nil
}
func (m *fakeCommissionRepository) ListRecords(ctx context.Context, opts commissionrepo.CommissionRecordListOptions) ([]model.CommissionRecord, int64, error) {
	return []model.CommissionRecord{}, 0, nil
}
func (m *fakeCommissionRepository) UpdateRecord(ctx context.Context, record *model.CommissionRecord) error {
	return nil
}

// MonthlySettlement methods
func (m *fakeCommissionRepository) CreateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error {
	return nil
}
func (m *fakeCommissionRepository) GetSettlement(ctx context.Context, id uint64) (*model.MonthlySettlement, error) {
	return &model.MonthlySettlement{}, nil
}
func (m *fakeCommissionRepository) GetSettlementByPlayerMonth(ctx context.Context, playerID uint64, month string) (*model.MonthlySettlement, error) {
	return &model.MonthlySettlement{}, nil
}
func (m *fakeCommissionRepository) ListSettlements(ctx context.Context, opts commissionrepo.SettlementListOptions) ([]model.MonthlySettlement, int64, error) {
	return []model.MonthlySettlement{}, 0, nil
}
func (m *fakeCommissionRepository) UpdateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error {
	return nil
}

// Stats methods
func (m *fakeCommissionRepository) GetMonthlyStats(ctx context.Context, month string) (*commissionrepo.MonthlyStats, error) {
	return &commissionrepo.MonthlyStats{}, nil
}
func (m *fakeCommissionRepository) GetPlayerMonthlyIncome(ctx context.Context, playerID uint64, month string) (int64, error) {
	return 0, nil
}

// ---- Helpers ----

func setupOrderTestService() (*order.OrderService, *fakeOrderRepository, *assignmentservice.Service) {
        orders := newFakeOrderRepository()
        svc := order.NewOrderService(orders, &fakePlayerRepository{}, &fakeUserRepository{}, &fakeGameRepository{}, &fakePaymentRepository{}, &fakeReviewRepository{}, &fakeCommissionRepository{})
        assignSvc := newAssignmentServiceStub(orders, &fakePlayerRepository{})
        return svc, orders, assignSvc
}

func fakeAuthMiddleware(userID uint64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	}
}

func setupOrderTestRouter(svc *order.OrderService, assignSvc *assignmentservice.Service, userID uint64) *gin.Engine {
        gin.SetMode(gin.TestMode)
        r := gin.New()
        RegisterOrderRoutes(r, svc, assignSvc, fakeAuthMiddleware(userID))
        return r
}

// ---- Tests ----

func TestUserOrder_GetOrderDetail_Success(t *testing.T) {
        svc, orders, assignSvc := setupOrderTestService()
	gameID := uint64(1)
	orders.orders[1] = &model.Order{
		Base:            model.Base{ID: 1},
		UserID:          100,
		GameID:          &gameID,
		ItemID:          1,
		OrderNo:         "USER-TEST-001",
		Title:           "Test Order",
		Status:          model.OrderStatusPending,
		TotalPriceCents: 1000,
		UnitPriceCents:  1000,
	}

        r := setupOrderTestRouter(svc, assignSvc, 100)
	req, _ := http.NewRequest(http.MethodGet, "/user/orders/1", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp model.APIResponse[order.OrderDetailResponse]
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if !resp.Success || resp.Code != http.StatusOK || resp.Message != "OK" {
		t.Fatalf("unexpected response envelope: success=%v code=%d message=%s", resp.Success, resp.Code, resp.Message)
	}
	if resp.Data.Order.ID != 1 {
		t.Fatalf("expected order id 1, got %d", resp.Data.Order.ID)
	}
}

func TestUserOrder_GetOrderDetail_NotFound(t *testing.T) {
        svc, _, assignSvc := setupOrderTestService()
        r := setupOrderTestRouter(svc, assignSvc, 123)
	req, _ := http.NewRequest(http.MethodGet, "/user/orders/9999", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
	var resp model.APIResponse[any]
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.Message != repository.ErrNotFound.Error() {
		t.Fatalf("expected message '%s', got '%s'", repository.ErrNotFound.Error(), resp.Message)
	}
}

func TestUserOrder_GetOrderDetail_Forbidden(t *testing.T) {
        svc, orders, assignSvc := setupOrderTestService()
	gameID := uint64(1)
	orders.orders[1] = &model.Order{Base: model.Base{ID: 1}, UserID: 100, GameID: &gameID, ItemID: 1, OrderNo: "USER-TEST-002", Title: "Test Order"}

	// user_id 200 is neither the order's user nor player
        r := setupOrderTestRouter(svc, assignSvc, 200)
	req, _ := http.NewRequest(http.MethodGet, "/user/orders/1", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", rec.Code)
	}
	var resp model.APIResponse[any]
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.Message != order.ErrUnauthorized.Error() {
		t.Fatalf("expected message '%s', got '%s'", order.ErrUnauthorized.Error(), resp.Message)
	}
}

func TestUserOrder_GetOrderDetail_InvalidID(t *testing.T) {
        svc, _, assignSvc := setupOrderTestService()
        r := setupOrderTestRouter(svc, assignSvc, 123)
	req, _ := http.NewRequest(http.MethodGet, "/user/orders/abc", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
	var resp model.APIResponse[any]
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	// Check that error message is reasonable for invalid ID
	if resp.Message == "" {
		t.Fatalf("expected non-empty error message")
	}
}

// ---- Tests for createOrderHandler ----

func TestUserOrder_CreateOrder_Success(t *testing.T) {
        svc, _, assignSvc := setupOrderTestService()
        r := setupOrderTestRouter(svc, assignSvc, 100)

	reqBody := `{"playerId":1,"gameId":1,"title":"Test Order","scheduledStart":"2025-01-15T10:00:00Z","durationHours":2.0}`
	req, _ := http.NewRequest(http.MethodPost, "/user/orders", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp model.APIResponse[order.CreateOrderResponse]
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if !resp.Success || resp.Code != http.StatusOK {
		t.Fatalf("unexpected response: success=%v code=%d", resp.Success, resp.Code)
	}
}

func TestUserOrder_CreateOrder_InvalidJSON(t *testing.T) {
        svc, _, assignSvc := setupOrderTestService()
        r := setupOrderTestRouter(svc, assignSvc, 100)

	req, _ := http.NewRequest(http.MethodPost, "/user/orders", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

// ---- Tests for getMyOrdersHandler ----

func TestUserOrder_GetMyOrders_Success(t *testing.T) {
        svc, orders, assignSvc := setupOrderTestService()
	gameID := uint64(1)
	orders.orders[1] = &model.Order{
		Base:            model.Base{ID: 1},
		UserID:          100,
		GameID:          &gameID,
		ItemID:          1,
		OrderNo:         "TEST-001",
		Status:          model.OrderStatusPending,
		TotalPriceCents: 2000,
	}
	orders.orders[2] = &model.Order{
		Base:            model.Base{ID: 2},
		UserID:          100,
		GameID:          &gameID,
		ItemID:          1,
		OrderNo:         "TEST-002",
		Status:          model.OrderStatusConfirmed,
		TotalPriceCents: 3000,
	}

        r := setupOrderTestRouter(svc, assignSvc, 100)
	req, _ := http.NewRequest(http.MethodGet, "/user/orders?page=1&pageSize=10", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp model.APIResponse[order.MyOrderListResponse]
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if !resp.Success {
		t.Fatalf("expected success=true, got %v", resp.Success)
	}
}

func TestUserOrder_GetMyOrders_WithStatusFilter(t *testing.T) {
        svc, orders, assignSvc := setupOrderTestService()
	gameID := uint64(1)
	orders.orders[1] = &model.Order{
		Base:   model.Base{ID: 1},
		UserID: 100,
		GameID: &gameID,
		ItemID: 1,
		Status: model.OrderStatusPending,
	}

        r := setupOrderTestRouter(svc, assignSvc, 100)
	req, _ := http.NewRequest(http.MethodGet, "/user/orders?status=pending&page=1", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
}

func TestUserOrder_GetMyOrders_InvalidQuery(t *testing.T) {
        svc, _, assignSvc := setupOrderTestService()
        r := setupOrderTestRouter(svc, assignSvc, 100)

	// Invalid page parameter
	req, _ := http.NewRequest(http.MethodGet, "/user/orders?page=invalid", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

// ---- Tests for cancelOrderHandler ----

func TestUserOrder_CancelOrder_Success(t *testing.T) {
        svc, orders, assignSvc := setupOrderTestService()
	gameID := uint64(1)
	orders.orders[1] = &model.Order{
		Base:   model.Base{ID: 1},
		UserID: 100,
		GameID: &gameID,
		ItemID: 1,
		Status: model.OrderStatusPending,
	}

        r := setupOrderTestRouter(svc, assignSvc, 100)
	reqBody := `{"reason":"不想要了"}`
	req, _ := http.NewRequest(http.MethodPut, "/user/orders/1/cancel", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp model.APIResponse[any]
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if !resp.Success {
		t.Fatalf("expected success=true, got %v", resp.Success)
	}
}

func TestUserOrder_CancelOrder_InvalidID(t *testing.T) {
        svc, _, assignSvc := setupOrderTestService()
        r := setupOrderTestRouter(svc, assignSvc, 100)

	reqBody := `{"reason":"test"}`
	req, _ := http.NewRequest(http.MethodPut, "/user/orders/invalid/cancel", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestUserOrder_CancelOrder_InvalidJSON(t *testing.T) {
        svc, orders, assignSvc := setupOrderTestService()
	gameID := uint64(1)
	orders.orders[1] = &model.Order{
		Base:   model.Base{ID: 1},
		UserID: 100,
		GameID: &gameID,
		Status: model.OrderStatusPending,
	}

        r := setupOrderTestRouter(svc, assignSvc, 100)
	req, _ := http.NewRequest(http.MethodPut, "/user/orders/1/cancel", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestUserOrder_CancelOrder_Unauthorized(t *testing.T) {
        svc, orders, assignSvc := setupOrderTestService()
	gameID := uint64(1)
	orders.orders[1] = &model.Order{
		Base:   model.Base{ID: 1},
		UserID: 100,
		GameID: &gameID,
		Status: model.OrderStatusPending,
	}

	// Different user trying to cancel
        r := setupOrderTestRouter(svc, assignSvc, 999)
	reqBody := `{"reason":"test"}`
	req, _ := http.NewRequest(http.MethodPut, "/user/orders/1/cancel", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", rec.Code)
	}
}

func TestUserOrder_CancelOrder_InvalidTransition(t *testing.T) {
        svc, orders, assignSvc := setupOrderTestService()
	gameID := uint64(1)
	orders.orders[1] = &model.Order{
		Base:   model.Base{ID: 1},
		UserID: 100,
		GameID: &gameID,
		Status: model.OrderStatusCompleted, // Already completed
	}

        r := setupOrderTestRouter(svc, assignSvc, 100)
	reqBody := `{"reason":"test"}`
	req, _ := http.NewRequest(http.MethodPut, "/user/orders/1/cancel", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

// ---- Tests for completeOrderHandler ----

func TestUserOrder_CompleteOrder_Success(t *testing.T) {
        svc, orders, assignSvc := setupOrderTestService()
	gameID := uint64(1)
	orders.orders[1] = &model.Order{
		Base:   model.Base{ID: 1},
		UserID: 100,
		GameID: &gameID,
		ItemID: 1,
		Status: model.OrderStatusInProgress,
	}

        r := setupOrderTestRouter(svc, assignSvc, 100)
	req, _ := http.NewRequest(http.MethodPut, "/user/orders/1/complete", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp model.APIResponse[any]
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if !resp.Success {
		t.Fatalf("expected success=true, got %v", resp.Success)
	}
}

func TestUserOrder_CompleteOrder_InvalidID(t *testing.T) {
        svc, _, assignSvc := setupOrderTestService()
        r := setupOrderTestRouter(svc, assignSvc, 100)

	req, _ := http.NewRequest(http.MethodPut, "/user/orders/invalid/complete", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestUserOrder_CompleteOrder_Unauthorized(t *testing.T) {
        svc, orders, assignSvc := setupOrderTestService()
	gameID := uint64(1)
	orders.orders[1] = &model.Order{
		Base:   model.Base{ID: 1},
		UserID: 100,
		GameID: &gameID,
		Status: model.OrderStatusInProgress,
	}

	// Different user trying to complete
        r := setupOrderTestRouter(svc, assignSvc, 999)
	req, _ := http.NewRequest(http.MethodPut, "/user/orders/1/complete", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", rec.Code)
	}
}

func TestUserOrder_CompleteOrder_InvalidTransition(t *testing.T) {
        svc, orders, assignSvc := setupOrderTestService()
	gameID := uint64(1)
	orders.orders[1] = &model.Order{
		Base:   model.Base{ID: 1},
		UserID: 100,
		GameID: &gameID,
		Status: model.OrderStatusPending, // Wrong status
	}

        r := setupOrderTestRouter(svc, assignSvc, 100)
	req, _ := http.NewRequest(http.MethodPut, "/user/orders/1/complete", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

// ---- Tests for getUserIDFromContext ----

func TestGetUserIDFromContext_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("user_id", uint64(123))

	userID := getUserIDFromContext(c)
	if userID != 123 {
		t.Fatalf("expected userID 123, got %d", userID)
	}
}

func TestGetUserIDFromContext_NotExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	userID := getUserIDFromContext(c)
	if userID != 0 {
		t.Fatalf("expected userID 0, got %d", userID)
	}
}

func TestGetUserIDFromContext_WrongType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("user_id", "not a uint64")

	userID := getUserIDFromContext(c)
	if userID != 0 {
		t.Fatalf("expected userID 0, got %d", userID)
	}
}
