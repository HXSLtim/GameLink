package user

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

// ---- Helpers ----

func setupOrderTestService() (*order.OrderService, *fakeOrderRepository) {
	orders := newFakeOrderRepository()
	svc := order.NewOrderService(orders, &fakePlayerRepository{}, &fakeUserRepository{}, &fakeGameRepository{}, &fakePaymentRepository{}, &fakeReviewRepository{})
	return svc, orders
}

func fakeAuthMiddleware(userID uint64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	}
}

func setupOrderTestRouter(svc *order.OrderService, userID uint64) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	RegisterOrderRoutes(r, svc, fakeAuthMiddleware(userID))
	return r
}

// ---- Tests ----

func TestUserOrder_GetOrderDetail_Success(t *testing.T) {
	svc, orders := setupOrderTestService()
	orders.orders[1] = &model.Order{
		Base:       model.Base{ID: 1},
		UserID:     100,
		PlayerID:   0,
		GameID:     1,
		Title:      "Test Order",
		Status:     model.OrderStatusPending,
		PriceCents: 1000,
	}

	r := setupOrderTestRouter(svc, 100)
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
	svc, _ := setupOrderTestService()
	r := setupOrderTestRouter(svc, 123)
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
	svc, orders := setupOrderTestService()
	orders.orders[1] = &model.Order{Base: model.Base{ID: 1}, UserID: 100, PlayerID: 0, GameID: 1, Title: "Test Order"}

	// user_id 200 is neither the order's user nor player
	r := setupOrderTestRouter(svc, 200)
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
	svc, _ := setupOrderTestService()
	r := setupOrderTestRouter(svc, 123)
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
	if resp.Message != ErrInvalidID {
		t.Fatalf("expected message '%s', got '%s'", ErrInvalidID, resp.Message)
	}
}
