package admin

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
	"gamelink/internal/service"
)

// buildTestRouter constructs a Gin engine with admin routes and a provided service.
func buildTestRouter(svc *service.AdminService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api/v1")
	RegisterRoutes(api, svc)
	return r
}

func TestAdminRoutes_UnauthorizedWhenTokenConfigured(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("ADMIN_TOKEN", "secret")

	svc := service.NewAdminService(&fakeGameRepo{}, &fakeUserRepo{}, &fakePlayerRepo{}, &fakeOrderRepo{}, &fakePaymentRepo{}, nil)
	r := buildTestRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/games", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAdminRoutes_ListGames_Envelope(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("ADMIN_TOKEN", "")

	svc := service.NewAdminService(&fakeGameRepo{items: []model.Game{{Name: "A"}}}, &fakeUserRepo{}, &fakePlayerRepo{}, &fakeOrderRepo{}, &fakePaymentRepo{}, nil)
	r := buildTestRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/games?page=1&page_size=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var body struct {
		Success bool         `json:"success"`
		Code    int          `json:"code"`
		Data    []model.Game `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if !body.Success || body.Code != 200 || len(body.Data) == 0 {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}
}

func TestAdminRoutes_UpdateOrder_AcceptsCancelledSpelling(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("ADMIN_TOKEN", "")

	o := &model.Order{Base: model.Base{ID: 1}, Status: model.OrderStatusPending, Currency: model.CurrencyCNY}
	oRepo := &fakeOrderRepo{obj: o}
	svc := service.NewAdminService(&fakeGameRepo{}, &fakeUserRepo{}, &fakePlayerRepo{}, oRepo, &fakePaymentRepo{}, nil)
	r := buildTestRouter(svc)

	payload := map[string]any{
		"status":      "cancelled", //nolint:misspell // intentionally testing legacy spelling
		"price_cents": 100,
		"currency":    "USD",
	}
	buf, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/orders/1", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body=%s", w.Code, w.Body.String())
	}
	if oRepo.obj.Status != model.OrderStatusCanceled {
		t.Fatalf("expected status normalized to 'canceled', got %q", oRepo.obj.Status)
	}
}

// ---- minimal fakes for integration tests ----

type fakeGameRepo struct{ items []model.Game }

func (f *fakeGameRepo) List(ctx context.Context) ([]model.Game, error) {
	return append([]model.Game(nil), f.items...), nil
}
func (f *fakeGameRepo) ListPaged(ctx context.Context, page, size int) ([]model.Game, int64, error) {
	return append([]model.Game(nil), f.items...), int64(len(f.items)), nil
}
func (f *fakeGameRepo) Get(ctx context.Context, id uint64) (*model.Game, error) {
	return nil, repository.ErrNotFound
}
func (f *fakeGameRepo) Create(ctx context.Context, g *model.Game) error { return nil }
func (f *fakeGameRepo) Update(ctx context.Context, g *model.Game) error { return nil }
func (f *fakeGameRepo) Delete(ctx context.Context, id uint64) error     { return nil }

type fakeUserRepo struct{}

func (f *fakeUserRepo) List(ctx context.Context) ([]model.User, error) { return nil, nil }
func (f *fakeUserRepo) ListPaged(ctx context.Context, page, size int) ([]model.User, int64, error) {
	return nil, 0, nil
}
func (f *fakeUserRepo) Get(ctx context.Context, id uint64) (*model.User, error) {
	return nil, repository.ErrNotFound
}
func (f *fakeUserRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return nil, repository.ErrNotFound
}
func (f *fakeUserRepo) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	return nil, repository.ErrNotFound
}
func (f *fakeUserRepo) Create(ctx context.Context, u *model.User) error { return nil }
func (f *fakeUserRepo) Update(ctx context.Context, u *model.User) error { return nil }
func (f *fakeUserRepo) Delete(ctx context.Context, id uint64) error     { return nil }

type fakePlayerRepo struct{}

func (f *fakePlayerRepo) List(ctx context.Context) ([]model.Player, error) { return nil, nil }
func (f *fakePlayerRepo) ListPaged(ctx context.Context, page, size int) ([]model.Player, int64, error) {
	return nil, 0, nil
}
func (f *fakePlayerRepo) Get(ctx context.Context, id uint64) (*model.Player, error) {
	return nil, repository.ErrNotFound
}
func (f *fakePlayerRepo) Create(ctx context.Context, p *model.Player) error { return nil }
func (f *fakePlayerRepo) Update(ctx context.Context, p *model.Player) error { return nil }
func (f *fakePlayerRepo) Delete(ctx context.Context, id uint64) error       { return nil }

type fakeOrderRepo struct{ obj *model.Order }

func (f *fakeOrderRepo) List(ctx context.Context, _ repository.OrderListOptions) ([]model.Order, int64, error) {
	return nil, 0, nil
}
func (f *fakeOrderRepo) Get(ctx context.Context, id uint64) (*model.Order, error) {
	if f.obj == nil {
		return nil, repository.ErrNotFound
	}
	return f.obj, nil
}
func (f *fakeOrderRepo) Update(ctx context.Context, o *model.Order) error { f.obj = o; return nil }
func (f *fakeOrderRepo) Delete(ctx context.Context, id uint64) error      { return nil }

type fakePaymentRepo struct{}

func (f *fakePaymentRepo) List(ctx context.Context, _ repository.PaymentListOptions) ([]model.Payment, int64, error) {
	return nil, 0, nil
}
func (f *fakePaymentRepo) Get(ctx context.Context, id uint64) (*model.Payment, error) {
	return nil, repository.ErrNotFound
}
func (f *fakePaymentRepo) Update(ctx context.Context, p *model.Payment) error { return nil }
func (f *fakePaymentRepo) Delete(ctx context.Context, id uint64) error        { return nil }

func TestPaymentHandler_InvalidTime_Returns400(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("ADMIN_TOKEN", "")

	pay := &model.Payment{Base: model.Base{ID: 1}, Status: model.PaymentStatusPending}
	pRepo := &fakePaymentRepoWithObj{obj: pay}
	svc := service.NewAdminService(&fakeGameRepo{}, &fakeUserRepo{}, &fakePlayerRepo{}, &fakeOrderRepo{}, pRepo, nil)
	r := buildTestRouter(svc)

	payload := map[string]any{
		"status":  "paid",
		"paid_at": "not-a-time",
	}
	buf, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/payments/1", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid paid_at, got %d", w.Code)
	}
}

type fakePaymentRepoWithObj struct{ obj *model.Payment }

func (f *fakePaymentRepoWithObj) List(ctx context.Context, _ repository.PaymentListOptions) ([]model.Payment, int64, error) {
	return nil, 0, nil
}
func (f *fakePaymentRepoWithObj) Get(ctx context.Context, id uint64) (*model.Payment, error) {
	if f.obj == nil {
		return nil, repository.ErrNotFound
	}
	return f.obj, nil
}
func (f *fakePaymentRepoWithObj) Update(ctx context.Context, p *model.Payment) error {
	f.obj = p
	return nil
}
func (f *fakePaymentRepoWithObj) Delete(ctx context.Context, id uint64) error { return nil }
