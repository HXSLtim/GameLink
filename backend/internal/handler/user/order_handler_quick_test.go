package user

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
    "context"

    "github.com/gin-gonic/gin"

    "gamelink/internal/model"
    "gamelink/internal/repository"
    commissionrepo "gamelink/internal/repository/commission"
    ordersvc "gamelink/internal/service/order"
)

type fakeOrderRepo struct{ items map[uint64]*model.Order; next uint64 }
func newFakeOrderRepo() *fakeOrderRepo { return &fakeOrderRepo{items: map[uint64]*model.Order{}, next:1} }
func (f *fakeOrderRepo) Create(ctx context.Context, o *model.Order) error { o.ID = f.next; f.next++; f.items[o.ID] = o; return nil }
func (f *fakeOrderRepo) List(ctx context.Context, _ repository.OrderListOptions) ([]model.Order, int64, error) { out := make([]model.Order,0,len(f.items)); for _, v := range f.items { out = append(out, *v) } ; return out, int64(len(out)), nil }
func (f *fakeOrderRepo) Get(ctx context.Context, id uint64) (*model.Order, error) { v := f.items[id]; if v==nil { return nil, repository.ErrNotFound } ; return v, nil }
func (f *fakeOrderRepo) Update(ctx context.Context, o *model.Order) error { f.items[o.ID] = o; return nil }
func (f *fakeOrderRepo) Delete(ctx context.Context, id uint64) error { delete(f.items, id); return nil }

type fakePlayerRepoOrd struct{}
func (fakePlayerRepoOrd) List(context.Context) ([]model.Player, error) { return nil, nil }
func (fakePlayerRepoOrd) ListPaged(context.Context, int, int) ([]model.Player, int64, error) { return []model.Player{{Base: model.Base{ID:1}, UserID:1, Nickname:"p", HourlyRateCents: 1000}}, 1, nil }
func (fakePlayerRepoOrd) Get(context.Context, uint64) (*model.Player, error) { return &model.Player{Base: model.Base{ID:1}, UserID:1, Nickname:"p", HourlyRateCents: 1000}, nil }
func (fakePlayerRepoOrd) Create(context.Context, *model.Player) error { return nil }
func (fakePlayerRepoOrd) Update(context.Context, *model.Player) error { return nil }
func (fakePlayerRepoOrd) Delete(context.Context, uint64) error { return nil }

type fakeUserRepoOrd struct{}
func (fakeUserRepoOrd) List(context.Context) ([]model.User, error) { return nil, nil }
func (fakeUserRepoOrd) ListPaged(context.Context, int, int) ([]model.User, int64, error) { return nil, 0, nil }
func (fakeUserRepoOrd) ListWithFilters(context.Context, repository.UserListOptions) ([]model.User, int64, error) { return nil, 0, nil }
func (fakeUserRepoOrd) Get(context.Context, uint64) (*model.User, error) { return &model.User{Base: model.Base{ID:1}, AvatarURL:"a"}, nil }
func (fakeUserRepoOrd) GetByPhone(context.Context, string) (*model.User, error) { return nil, repository.ErrNotFound }
func (fakeUserRepoOrd) FindByEmail(context.Context, string) (*model.User, error) { return nil, repository.ErrNotFound }
func (fakeUserRepoOrd) FindByPhone(context.Context, string) (*model.User, error) { return nil, repository.ErrNotFound }
func (fakeUserRepoOrd) Create(context.Context, *model.User) error { return nil }
func (fakeUserRepoOrd) Update(context.Context, *model.User) error { return nil }
func (fakeUserRepoOrd) Delete(context.Context, uint64) error { return nil }

type fakeGameRepoOrd struct{}
func (fakeGameRepoOrd) List(context.Context) ([]model.Game, error) { return nil, nil }
func (fakeGameRepoOrd) ListPaged(context.Context, int, int) ([]model.Game, int64, error) { return nil, 0, nil }
func (fakeGameRepoOrd) Get(context.Context, uint64) (*model.Game, error) { return &model.Game{Base: model.Base{ID:1}, Name:"g"}, nil }
func (fakeGameRepoOrd) Create(context.Context, *model.Game) error { return nil }
func (fakeGameRepoOrd) Update(context.Context, *model.Game) error { return nil }
func (fakeGameRepoOrd) Delete(context.Context, uint64) error { return nil }

type fakePaymentRepoOrd struct{ items map[uint64]model.Payment }
func (f *fakePaymentRepoOrd) Create(context.Context, *model.Payment) error { return nil }
func (f *fakePaymentRepoOrd) List(context.Context, repository.PaymentListOptions) ([]model.Payment, int64, error) { out:=make([]model.Payment,0,len(f.items)); for _, v:= range f.items { out = append(out, v) } ; return out, int64(len(out)), nil }
func (f *fakePaymentRepoOrd) Get(context.Context, uint64) (*model.Payment, error) { return nil, repository.ErrNotFound }
func (f *fakePaymentRepoOrd) Update(context.Context, *model.Payment) error { return nil }
func (f *fakePaymentRepoOrd) Delete(context.Context, uint64) error { return nil }

type fakeReviewRepoOrd struct{}
func (fakeReviewRepoOrd) List(context.Context, repository.ReviewListOptions) ([]model.Review, int64, error) { return nil, 0, nil }
func (fakeReviewRepoOrd) Get(context.Context, uint64) (*model.Review, error) { return nil, repository.ErrNotFound }
func (fakeReviewRepoOrd) Create(context.Context, *model.Review) error { return nil }
func (fakeReviewRepoOrd) Update(context.Context, *model.Review) error { return nil }
func (fakeReviewRepoOrd) Delete(context.Context, uint64) error { return nil }

type fakeCommissionRepoOrd struct{}
func (fakeCommissionRepoOrd) CreateRule(context.Context, *model.CommissionRule) error { return nil }
func (fakeCommissionRepoOrd) GetRule(context.Context, uint64) (*model.CommissionRule, error) { return nil, repository.ErrNotFound }
func (fakeCommissionRepoOrd) GetDefaultRule(context.Context) (*model.CommissionRule, error) { return &model.CommissionRule{Rate:20}, nil }
func (fakeCommissionRepoOrd) GetRuleForOrder(context.Context, *uint64, *uint64, *string) (*model.CommissionRule, error) { return nil, repository.ErrNotFound }
func (fakeCommissionRepoOrd) ListRules(context.Context, commissionrepo.CommissionRuleListOptions) ([]model.CommissionRule, int64, error) { return nil, 0, nil }
func (fakeCommissionRepoOrd) UpdateRule(context.Context, *model.CommissionRule) error { return nil }
func (fakeCommissionRepoOrd) DeleteRule(context.Context, uint64) error { return nil }
func (fakeCommissionRepoOrd) CreateRecord(context.Context, *model.CommissionRecord) error { return nil }
func (fakeCommissionRepoOrd) GetRecord(context.Context, uint64) (*model.CommissionRecord, error) { return nil, repository.ErrNotFound }
func (fakeCommissionRepoOrd) GetRecordByOrderID(context.Context, uint64) (*model.CommissionRecord, error) { return nil, repository.ErrNotFound }
func (fakeCommissionRepoOrd) ListRecords(context.Context, commissionrepo.CommissionRecordListOptions) ([]model.CommissionRecord, int64, error) { return nil, 0, nil }
func (fakeCommissionRepoOrd) UpdateRecord(context.Context, *model.CommissionRecord) error { return nil }
func (fakeCommissionRepoOrd) CreateSettlement(context.Context, *model.MonthlySettlement) error { return nil }
func (fakeCommissionRepoOrd) GetSettlement(context.Context, uint64) (*model.MonthlySettlement, error) { return nil, repository.ErrNotFound }
func (fakeCommissionRepoOrd) GetSettlementByPlayerMonth(context.Context, uint64, string) (*model.MonthlySettlement, error) { return nil, repository.ErrNotFound }
func (fakeCommissionRepoOrd) ListSettlements(context.Context, commissionrepo.SettlementListOptions) ([]model.MonthlySettlement, int64, error) { return nil, 0, nil }
func (fakeCommissionRepoOrd) UpdateSettlement(context.Context, *model.MonthlySettlement) error { return nil }
func (fakeCommissionRepoOrd) GetMonthlyStats(context.Context, string) (*commissionrepo.MonthlyStats, error) { return &commissionrepo.MonthlyStats{}, nil }
func (fakeCommissionRepoOrd) GetPlayerMonthlyIncome(context.Context, uint64, string) (int64, error) { return 0, nil }

func setupUserOrderRouter(repo *fakeOrderRepo) *gin.Engine {
    r := gin.New()
    r.Use(func(c *gin.Context){ c.Set("user_id", uint64(1)); c.Next() })
    svc := ordersvc.NewOrderService(repo, fakePlayerRepoOrd{}, fakeUserRepoOrd{}, fakeGameRepoOrd{}, &fakePaymentRepoOrd{}, fakeReviewRepoOrd{}, &fakeCommissionRepoOrd{})
    RegisterOrderRoutes(r, svc, func(c *gin.Context){ c.Next() })
    return r
}

func TestUserOrder_CRUDFlow(t *testing.T) {
    r := setupUserOrderRouter(newFakeOrderRepo())
    start := time.Now().Add(1 * time.Hour)
    payload := ordersvc.CreateOrderRequest{PlayerID:1, GameID:1, Title:"t", Description:"d", ScheduledStart:&start, DurationHours:1}
    b, _ := json.Marshal(payload)
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPost, "/user/orders", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    if w.Code != http.StatusOK { t.Fatalf("%d", w.Code) }

    w2 := httptest.NewRecorder()
    req2 := httptest.NewRequest(http.MethodGet, "/user/orders?page=1&pageSize=10", nil)
    r.ServeHTTP(w2, req2)
    if w2.Code != http.StatusOK { t.Fatalf("%d", w2.Code) }

    w3 := httptest.NewRecorder()
    req3 := httptest.NewRequest(http.MethodGet, "/user/orders/abc", nil)
    r.ServeHTTP(w3, req3)
    if w3.Code != http.StatusBadRequest { t.Fatalf("%d", w3.Code) }
}

func TestUserOrder_DetailCancelComplete(t *testing.T) {
    repo := newFakeOrderRepo()
    repo.Create(context.Background(), &model.Order{Base: model.Base{ID:1}, UserID:1, Status:model.OrderStatusPending, Title:"x"})
    r := setupUserOrderRouter(repo)

    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/user/orders/1", nil)
    r.ServeHTTP(w, req)
    if w.Code != http.StatusOK { t.Fatalf("%d", w.Code) }

    cb, _ := json.Marshal(ordersvc.CancelOrderRequest{Reason:"r"})
    w2 := httptest.NewRecorder()
    req2 := httptest.NewRequest(http.MethodPut, "/user/orders/1/cancel", bytes.NewReader(cb))
    req2.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w2, req2)
    if w2.Code != http.StatusOK { t.Fatalf("%d", w2.Code) }

    repo.items[2] = &model.Order{Base: model.Base{ID:2}, UserID:1, Status:model.OrderStatusPending}
    w3 := httptest.NewRecorder()
    req3 := httptest.NewRequest(http.MethodPut, "/user/orders/2/complete", nil)
    r.ServeHTTP(w3, req3)
    if w3.Code != http.StatusBadRequest { t.Fatalf("%d", w3.Code) }

    repo.items[2].Status = model.OrderStatusInProgress
    w4 := httptest.NewRecorder()
    req4 := httptest.NewRequest(http.MethodPut, "/user/orders/2/complete", nil)
    r.ServeHTTP(w4, req4)
    if w4.Code != http.StatusOK { t.Fatalf("%d", w4.Code) }
}
