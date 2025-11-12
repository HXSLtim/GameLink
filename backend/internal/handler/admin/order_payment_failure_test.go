package admin

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"
    "context"

    "github.com/gin-gonic/gin"

    "gamelink/internal/cache"
    "gamelink/internal/model"
    "gamelink/internal/repository"
    adminservice "gamelink/internal/service/admin"
)

type failOrderRepo struct{ items map[uint64]*model.Order }
func newFailOrderRepo() *failOrderRepo { return &failOrderRepo{items: map[uint64]*model.Order{}} }
func (r *failOrderRepo) Create(context.Context, *model.Order) error { return nil }
func (r *failOrderRepo) List(context.Context, repository.OrderListOptions) ([]model.Order, int64, error) { return nil, 0, nil }
func (r *failOrderRepo) Get(context.Context, uint64) (*model.Order, error) { return &model.Order{Base: model.Base{ID:1}, Status:model.OrderStatusPending}, nil }
func (r *failOrderRepo) Update(context.Context, *model.Order) error { return nil }
func (r *failOrderRepo) Delete(context.Context, uint64) error { return nil }

type failPayRepo struct{}
func (failPayRepo) Create(context.Context, *model.Payment) error { return nil }
func (failPayRepo) List(context.Context, repository.PaymentListOptions) ([]model.Payment, int64, error) { return nil, 0, nil }
func (failPayRepo) Get(context.Context, uint64) (*model.Payment, error) { return &model.Payment{Base: model.Base{ID:1}, Status:model.PaymentStatusPaid}, nil }
func (failPayRepo) Update(context.Context, *model.Payment) error { return nil }
func (failPayRepo) Delete(context.Context, uint64) error { return nil }

func setupOrderPaymentHandlers() (*OrderHandler, *PaymentHandler) {
    svc := adminservice.NewAdminService(nil, nil, nil, newFailOrderRepo(), failPayRepo{}, nil, cache.NewMemory())
    return NewOrderHandler(svc), NewPaymentHandler(svc)
}

func TestAdminOrder_Confirm_BadJSON(t *testing.T) {
    oh, _ := setupOrderPaymentHandlers()
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Params = gin.Params{{Key:"id", Value:"1"}}
    c.Request = httptest.NewRequest(http.MethodPost, "/admin/orders/1/confirm", bytes.NewReader([]byte("{")))
    c.Request.Header.Set("Content-Type", "application/json")
    oh.ConfirmOrder(c)
    if w.Code != http.StatusBadRequest { t.Fatalf("%d", w.Code) }
}

func TestAdminOrder_Start_InvalidID(t *testing.T) {
    oh, _ := setupOrderPaymentHandlers()
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Params = gin.Params{{Key:"id", Value:"abc"}}
    c.Request = httptest.NewRequest(http.MethodPost, "/admin/orders/abc/start", nil)
    oh.StartOrder(c)
    if w.Code != http.StatusBadRequest { t.Fatalf("%d", w.Code) }
}

func TestAdminOrder_Refund_BadJSON(t *testing.T) {
    oh, _ := setupOrderPaymentHandlers()
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Params = gin.Params{{Key:"id", Value:"1"}}
    c.Request = httptest.NewRequest(http.MethodPost, "/admin/orders/1/refund", bytes.NewReader([]byte("{")))
    c.Request.Header.Set("Content-Type", "application/json")
    oh.RefundOrder(c)
    if w.Code != http.StatusBadRequest { t.Fatalf("%d", w.Code) }
}

func TestAdminPayment_Capture_InvalidPaidAt(t *testing.T) {
    _, ph := setupOrderPaymentHandlers()
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Params = gin.Params{{Key:"id", Value:"1"}}
    c.Request = httptest.NewRequest(http.MethodPost, "/admin/payments/1/capture", bytes.NewReader([]byte(`{"paid_at":"bad"}`)))
    c.Request.Header.Set("Content-Type", "application/json")
    ph.CapturePayment(c)
    if w.Code != http.StatusBadRequest { t.Fatalf("%d", w.Code) }
}

