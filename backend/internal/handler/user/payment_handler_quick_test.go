package user

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "context"

    "github.com/gin-gonic/gin"

    "gamelink/internal/model"
    "gamelink/internal/repository"
    paysvc "gamelink/internal/service/payment"
)

type fakePaymentRepo struct{ items map[uint64]*model.Payment; next uint64 }
func newFakePaymentRepo() *fakePaymentRepo { return &fakePaymentRepo{items: map[uint64]*model.Payment{}, next:1} }
func (f *fakePaymentRepo) Create(ctx context.Context, p *model.Payment) error { p.ID=f.next; f.next++; f.items[p.ID]=p; return nil }
func (f *fakePaymentRepo) List(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) { out:=make([]model.Payment,0,len(f.items)); for _, v:= range f.items { out = append(out, *v) } ; return out, int64(len(out)), nil }
func (f *fakePaymentRepo) Get(ctx context.Context, id uint64) (*model.Payment, error) { v := f.items[id]; if v==nil { return nil, repository.ErrNotFound } ; return v, nil }
func (f *fakePaymentRepo) Update(ctx context.Context, p *model.Payment) error { f.items[p.ID]=p; return nil }
func (f *fakePaymentRepo) Delete(ctx context.Context, id uint64) error { delete(f.items, id); return nil }

type fakeOrderRepoPay struct{ items map[uint64]*model.Order }
func (f *fakeOrderRepoPay) Create(ctx context.Context, o *model.Order) error { _=ctx; _=o; return nil }
func (f *fakeOrderRepoPay) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) { _=ctx; _=opts; return nil, 0, nil }
func (f *fakeOrderRepoPay) Get(ctx context.Context, id uint64) (*model.Order, error) { _=ctx; o := f.items[id]; if o==nil { return nil, repository.ErrNotFound } ; return o, nil }
func (f *fakeOrderRepoPay) Update(ctx context.Context, o *model.Order) error { _=ctx; f.items[o.ID]=o; return nil }
func (f *fakeOrderRepoPay) Delete(ctx context.Context, id uint64) error { _=ctx; _=id; return nil }

func setupUserPaymentRouter(payRepo *fakePaymentRepo, ordRepo *fakeOrderRepoPay) *gin.Engine {
    r := gin.New()
    r.Use(func(c *gin.Context){ c.Set("user_id", uint64(1)); c.Next() })
    svc := paysvc.NewPaymentService(payRepo, ordRepo)
    RegisterPaymentRoutes(r, svc, func(c *gin.Context){ c.Next() })
    return r
}

func TestUserPayment_CreateStatusCancel(t *testing.T) {
    payRepo := newFakePaymentRepo()
    ordRepo := &fakeOrderRepoPay{items: map[uint64]*model.Order{1: {Base: model.Base{ID:1}, UserID:1, TotalPriceCents:1000, Currency:model.CurrencyCNY, Status:model.OrderStatusPending}}}
    r := setupUserPaymentRouter(payRepo, ordRepo)

    payload := paysvc.CreatePaymentRequest{OrderID:1, Method:model.PaymentMethodWeChat}
    b, _ := json.Marshal(payload)
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPost, "/user/payments", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    if w.Code != http.StatusOK { t.Fatalf("%d", w.Code) }

    w2 := httptest.NewRecorder()
    req2 := httptest.NewRequest(http.MethodGet, "/user/payments/abc", nil)
    r.ServeHTTP(w2, req2)
    if w2.Code != http.StatusBadRequest { t.Fatalf("%d", w2.Code) }

    pid := uint64(1)
    payRepo.items[pid] = &model.Payment{Base: model.Base{ID:pid}, UserID:1, OrderID:1, Method:model.PaymentMethodWeChat, AmountCents:1000, Currency:model.CurrencyCNY, Status:model.PaymentStatusPending}
    w3 := httptest.NewRecorder()
    req3 := httptest.NewRequest(http.MethodPost, "/user/payments/1/cancel", nil)
    r.ServeHTTP(w3, req3)
    if w3.Code != http.StatusOK { t.Fatalf("%d", w3.Code) }
}

func TestUserPayment_StatusNotFound(t *testing.T) {
    payRepo := newFakePaymentRepo()
    ordRepo := &fakeOrderRepoPay{items: map[uint64]*model.Order{}}
    r := setupUserPaymentRouter(payRepo, ordRepo)
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/user/payments/999", nil)
    r.ServeHTTP(w, req)
    if w.Code != http.StatusNotFound { t.Fatalf("%d", w.Code) }
}

func TestUserPayment_CancelFailWhenPaid(t *testing.T) {
    payRepo := newFakePaymentRepo()
    ordRepo := &fakeOrderRepoPay{items: map[uint64]*model.Order{1: {Base: model.Base{ID:1}, UserID:1, TotalPriceCents:1000, Currency:model.CurrencyCNY, Status:model.OrderStatusConfirmed}}}
    r := setupUserPaymentRouter(payRepo, ordRepo)
    payRepo.items[2] = &model.Payment{Base: model.Base{ID:2}, UserID:1, OrderID:1, Method:model.PaymentMethodWeChat, AmountCents:1000, Currency:model.CurrencyCNY, Status:model.PaymentStatusPaid}
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPost, "/user/payments/2/cancel", nil)
    r.ServeHTTP(w, req)
    if w.Code != http.StatusInternalServerError { t.Fatalf("%d", w.Code) }
}
