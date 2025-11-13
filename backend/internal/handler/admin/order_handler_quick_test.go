package admin

import (
    "bytes"
    "encoding/json"
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

type fakeOrderRepo2 struct{ m map[uint64]*model.Order }
func newFakeOrderRepo2() *fakeOrderRepo2 { return &fakeOrderRepo2{m: map[uint64]*model.Order{}} }
func (r *fakeOrderRepo2) Create(ctx context.Context, o *model.Order) error { r.m[o.ID] = o; if o.ID==0 { o.ID=1; r.m[o.ID]=o }; return nil }
func (r *fakeOrderRepo2) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) { out:=make([]model.Order,0,len(r.m)); for _, v:= range r.m { out = append(out, *v) } ; return out, int64(len(out)), nil }
func (r *fakeOrderRepo2) Get(ctx context.Context, id uint64) (*model.Order, error) { v:=r.m[id]; if v==nil { return nil, repository.ErrNotFound }; return v, nil }
func (r *fakeOrderRepo2) Update(ctx context.Context, o *model.Order) error { r.m[o.ID] = o; return nil }
func (r *fakeOrderRepo2) Delete(ctx context.Context, id uint64) error { delete(r.m, id); return nil }

type fakePaymentRepo2 struct{ m map[uint64]*model.Payment }
func newFakePaymentRepo2() *fakePaymentRepo2 { return &fakePaymentRepo2{m: map[uint64]*model.Payment{}} }
func (r *fakePaymentRepo2) Create(ctx context.Context, p *model.Payment) error { if p.ID==0 { p.ID = uint64(len(r.m)+1) } ; r.m[p.ID]=p; return nil }
func (r *fakePaymentRepo2) List(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) { out:=[]model.Payment{}; for _, v:= range r.m { if opts.OrderID!=nil && v.OrderID==*opts.OrderID { out = append(out, *v) } } ; return out, int64(len(out)), nil }
func (r *fakePaymentRepo2) Get(ctx context.Context, id uint64) (*model.Payment, error) { v:=r.m[id]; if v==nil { return nil, repository.ErrNotFound } ; return v, nil }
func (r *fakePaymentRepo2) Update(ctx context.Context, p *model.Payment) error { r.m[p.ID]=p; return nil }
func (r *fakePaymentRepo2) Delete(ctx context.Context, id uint64) error { delete(r.m, id); return nil }

type fakeUserRepo3 struct{}
func (fakeUserRepo3) List(context.Context) ([]model.User, error) { return nil, nil }
func (fakeUserRepo3) ListPaged(context.Context, int, int) ([]model.User, int64, error) { return nil, 0, nil }
func (fakeUserRepo3) ListWithFilters(context.Context, repository.UserListOptions) ([]model.User, int64, error) { return nil, 0, nil }
func (fakeUserRepo3) Get(context.Context, uint64) (*model.User, error) { return &model.User{Base: model.Base{ID:1}, Name:"u", Role:model.RoleUser, Status:model.UserStatusActive}, nil }
func (fakeUserRepo3) GetByPhone(context.Context, string) (*model.User, error) { return nil, repository.ErrNotFound }
func (fakeUserRepo3) FindByEmail(context.Context, string) (*model.User, error) { return nil, repository.ErrNotFound }
func (fakeUserRepo3) FindByPhone(context.Context, string) (*model.User, error) { return nil, repository.ErrNotFound }
func (fakeUserRepo3) Create(context.Context, *model.User) error { return nil }
func (fakeUserRepo3) Update(context.Context, *model.User) error { return nil }
func (fakeUserRepo3) Delete(context.Context, uint64) error { return nil }

type fakePlayerRepo3 struct{}
func (fakePlayerRepo3) List(context.Context) ([]model.Player, error) { return nil, nil }
func (fakePlayerRepo3) ListPaged(context.Context, int, int) ([]model.Player, int64, error) { return nil, 0, nil }
func (fakePlayerRepo3) Get(context.Context, uint64) (*model.Player, error) { return &model.Player{Base: model.Base{ID:1}, UserID:1, Nickname:"p"}, nil }
func (fakePlayerRepo3) Create(context.Context, *model.Player) error { return nil }
func (fakePlayerRepo3) Update(context.Context, *model.Player) error { return nil }
func (fakePlayerRepo3) Delete(context.Context, uint64) error { return nil }

func setupOrderHandlerWithOrder(o *model.Order, payments []*model.Payment) *OrderHandler {
    orders := newFakeOrderRepo2()
    orders.m[1] = o
    pays := newFakePaymentRepo2()
    for _, p := range payments { _ = pays.Create(context.Background(), p) }
    svc := adminservice.NewAdminService(nil, fakeUserRepo3{}, fakePlayerRepo3{}, orders, pays, nil, cache.NewMemory())
    return NewOrderHandler(svc)
}

func TestAdminOrder_ConfirmStartCompleteRefundFlow(t *testing.T) {
    gin.SetMode(gin.TestMode)
    order := &model.Order{Base: model.Base{ID:1}, UserID:1, TotalPriceCents:1000, Currency:model.CurrencyCNY, Status:model.OrderStatusPending}
    pay := &model.Payment{Base: model.Base{ID:1}, OrderID:1, UserID:1, Method:model.PaymentMethodWeChat, AmountCents:1000, Currency:model.CurrencyCNY, Status:model.PaymentStatusPaid}
    h := setupOrderHandlerWithOrder(order, []*model.Payment{pay})

    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Params = gin.Params{{Key:"id", Value:"1"}}
    c.Request = httptest.NewRequest(http.MethodPost, "/admin/orders/1/confirm", bytes.NewReader([]byte(`{"note":"n"}`)))
    c.Request.Header.Set("Content-Type", "application/json")
    h.ConfirmOrder(c)
    if w.Code != http.StatusOK { t.Fatalf("%d", w.Code) }

    w2 := httptest.NewRecorder()
    c2, _ := gin.CreateTestContext(w2)
    c2.Params = gin.Params{{Key:"id", Value:"1"}}
    c2.Request = httptest.NewRequest(http.MethodPost, "/admin/orders/1/start", bytes.NewReader([]byte(`{"note":"n"}`)))
    c2.Request.Header.Set("Content-Type", "application/json")
    h.StartOrder(c2)
    if w2.Code != http.StatusOK { t.Fatalf("%d", w2.Code) }

    w3 := httptest.NewRecorder()
    c3, _ := gin.CreateTestContext(w3)
    c3.Params = gin.Params{{Key:"id", Value:"1"}}
    c3.Request = httptest.NewRequest(http.MethodPost, "/admin/orders/1/complete", bytes.NewReader([]byte(`{"note":"n"}`)))
    c3.Request.Header.Set("Content-Type", "application/json")
    h.CompleteOrder(c3)
    if w3.Code != http.StatusOK { t.Fatalf("%d", w3.Code) }

    b, _ := json.Marshal(orderRefundPayload{Reason:"dup", AmountCents:nil, Note:"n"})
    w4 := httptest.NewRecorder()
    c4, _ := gin.CreateTestContext(w4)
    c4.Params = gin.Params{{Key:"id", Value:"1"}}
    c4.Request = httptest.NewRequest(http.MethodPost, "/admin/orders/1/refund", bytes.NewReader(b))
    c4.Request.Header.Set("Content-Type", "application/json")
    h.RefundOrder(c4)
    if w4.Code != http.StatusOK { t.Fatalf("%d", w4.Code) }
}

func TestAdminOrder_InvalidID(t *testing.T) {
    h := setupOrderHandlerWithOrder(&model.Order{Base: model.Base{ID:1}, UserID:1, TotalPriceCents:1000, Currency:model.CurrencyCNY, Status:model.OrderStatusPending}, nil)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Params = gin.Params{{Key:"id", Value:"abc"}}
    h.ConfirmOrder(c)
    if w.Code != http.StatusBadRequest { t.Fatalf("%d", w.Code) }
}