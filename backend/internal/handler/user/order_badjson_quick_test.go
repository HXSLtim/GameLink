package user

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"
    "context"

    "github.com/gin-gonic/gin"

    "gamelink/internal/model"
    "gamelink/internal/repository"
    commissionrepo "gamelink/internal/repository/commission"
    ordersvc "gamelink/internal/service/order"
)

type obOrderRepo struct{ items map[uint64]*model.Order }
func newObOrderRepo() *obOrderRepo { return &obOrderRepo{items: map[uint64]*model.Order{}} }
func (r *obOrderRepo) Create(ctx context.Context, o *model.Order) error { if o.ID==0 { o.ID=uint64(len(r.items)+1) } ; r.items[o.ID]=o; return nil }
func (r *obOrderRepo) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) { out:=[]model.Order{}; for _, v:= range r.items { out = append(out, *v) } ; return out, int64(len(out)), nil }
func (r *obOrderRepo) Get(ctx context.Context, id uint64) (*model.Order, error) { v:=r.items[id]; if v==nil { return nil, repository.ErrNotFound } ; return v, nil }
func (r *obOrderRepo) Update(ctx context.Context, o *model.Order) error { r.items[o.ID]=o; return nil }
func (r *obOrderRepo) Delete(ctx context.Context, id uint64) error { delete(r.items, id); return nil }

type obPlayers struct{}
func (obPlayers) List(context.Context) ([]model.Player, error) { return nil, nil }
func (obPlayers) ListPaged(context.Context, int, int) ([]model.Player, int64, error) { return []model.Player{{Base: model.Base{ID:1}, UserID:1, Nickname:"p", HourlyRateCents: 1000}}, 1, nil }
func (obPlayers) Get(context.Context, uint64) (*model.Player, error) { return &model.Player{Base: model.Base{ID:1}, UserID:1, Nickname:"p", HourlyRateCents: 1000}, nil }
func (obPlayers) Create(context.Context, *model.Player) error { return nil }
func (obPlayers) Update(context.Context, *model.Player) error { return nil }
func (obPlayers) Delete(context.Context, uint64) error { return nil }

type obUsers struct{}
func (obUsers) List(context.Context) ([]model.User, error) { return nil, nil }
func (obUsers) ListPaged(context.Context, int, int) ([]model.User, int64, error) { return nil, 0, nil }
func (obUsers) ListWithFilters(context.Context, repository.UserListOptions) ([]model.User, int64, error) { return nil, 0, nil }
func (obUsers) Get(context.Context, uint64) (*model.User, error) { return &model.User{Base: model.Base{ID:1}, AvatarURL:"a"}, nil }
func (obUsers) GetByPhone(context.Context, string) (*model.User, error) { return nil, repository.ErrNotFound }
func (obUsers) FindByEmail(context.Context, string) (*model.User, error) { return nil, repository.ErrNotFound }
func (obUsers) FindByPhone(context.Context, string) (*model.User, error) { return nil, repository.ErrNotFound }
func (obUsers) Create(context.Context, *model.User) error { return nil }
func (obUsers) Update(context.Context, *model.User) error { return nil }
func (obUsers) Delete(context.Context, uint64) error { return nil }

type obGames struct{}
func (obGames) List(context.Context) ([]model.Game, error) { return nil, nil }
func (obGames) ListPaged(context.Context, int, int) ([]model.Game, int64, error) { return nil, 0, nil }
func (obGames) Get(context.Context, uint64) (*model.Game, error) { return &model.Game{Base: model.Base{ID:1}, Name:"g"}, nil }
func (obGames) Create(context.Context, *model.Game) error { return nil }
func (obGames) Update(context.Context, *model.Game) error { return nil }
func (obGames) Delete(context.Context, uint64) error { return nil }

type obPayments struct{}
func (obPayments) Create(context.Context, *model.Payment) error { return nil }
func (obPayments) List(context.Context, repository.PaymentListOptions) ([]model.Payment, int64, error) { return nil, 0, nil }
func (obPayments) Get(context.Context, uint64) (*model.Payment, error) { return nil, repository.ErrNotFound }
func (obPayments) Update(context.Context, *model.Payment) error { return nil }
func (obPayments) Delete(context.Context, uint64) error { return nil }

type obReviews struct{}
func (obReviews) List(context.Context, repository.ReviewListOptions) ([]model.Review, int64, error) { return nil, 0, nil }
func (obReviews) Get(context.Context, uint64) (*model.Review, error) { return nil, repository.ErrNotFound }
func (obReviews) Create(context.Context, *model.Review) error { return nil }
func (obReviews) Update(context.Context, *model.Review) error { return nil }
func (obReviews) Delete(context.Context, uint64) error { return nil }

type obCommissions struct{}
func (obCommissions) CreateRule(context.Context, *model.CommissionRule) error { return nil }
func (obCommissions) GetRule(context.Context, uint64) (*model.CommissionRule, error) { return nil, repository.ErrNotFound }
func (obCommissions) GetDefaultRule(context.Context) (*model.CommissionRule, error) { return &model.CommissionRule{Rate:20}, nil }
func (obCommissions) GetRuleForOrder(context.Context, *uint64, *uint64, *string) (*model.CommissionRule, error) { return nil, repository.ErrNotFound }
func (obCommissions) ListRules(context.Context, commissionrepo.CommissionRuleListOptions) ([]model.CommissionRule, int64, error) { return nil, 0, nil }
func (obCommissions) UpdateRule(context.Context, *model.CommissionRule) error { return nil }
func (obCommissions) DeleteRule(context.Context, uint64) error { return nil }
func (obCommissions) CreateRecord(context.Context, *model.CommissionRecord) error { return nil }
func (obCommissions) GetRecord(context.Context, uint64) (*model.CommissionRecord, error) { return nil, repository.ErrNotFound }
func (obCommissions) GetRecordByOrderID(context.Context, uint64) (*model.CommissionRecord, error) { return nil, repository.ErrNotFound }
func (obCommissions) ListRecords(context.Context, commissionrepo.CommissionRecordListOptions) ([]model.CommissionRecord, int64, error) { return nil, 0, nil }
func (obCommissions) UpdateRecord(context.Context, *model.CommissionRecord) error { return nil }
func (obCommissions) CreateSettlement(context.Context, *model.MonthlySettlement) error { return nil }
func (obCommissions) GetSettlement(context.Context, uint64) (*model.MonthlySettlement, error) { return nil, repository.ErrNotFound }
func (obCommissions) GetSettlementByPlayerMonth(context.Context, uint64, string) (*model.MonthlySettlement, error) { return nil, repository.ErrNotFound }
func (obCommissions) ListSettlements(context.Context, commissionrepo.SettlementListOptions) ([]model.MonthlySettlement, int64, error) { return nil, 0, nil }
func (obCommissions) UpdateSettlement(context.Context, *model.MonthlySettlement) error { return nil }
func (obCommissions) GetMonthlyStats(context.Context, string) (*commissionrepo.MonthlyStats, error) { return &commissionrepo.MonthlyStats{}, nil }
func (obCommissions) GetPlayerMonthlyIncome(context.Context, uint64, string) (int64, error) { return 0, nil }

func setupOrderBadRouter(repo *obOrderRepo) *gin.Engine {
    r := gin.New()
    r.Use(func(c *gin.Context){ c.Set("user_id", uint64(1)); c.Next() })
    svc := ordersvc.NewOrderService(repo, obPlayers{}, obUsers{}, obGames{}, obPayments{}, obReviews{}, obCommissions{})
    RegisterOrderRoutes(r, svc, func(c *gin.Context){ c.Next() })
    return r
}

func TestUserOrder_Create_BadJSON(t *testing.T) {
    r := setupOrderBadRouter(newObOrderRepo())
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPost, "/user/orders", bytes.NewReader([]byte("{")))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    if w.Code != http.StatusBadRequest { t.Fatalf("%d", w.Code) }
}

func TestUserOrder_Cancel_InvalidID(t *testing.T) {
    r := setupOrderBadRouter(newObOrderRepo())
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPut, "/user/orders/abc/cancel", bytes.NewReader([]byte(`{"reason":"r"}`)))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    if w.Code != http.StatusBadRequest { t.Fatalf("%d", w.Code) }
}

func TestUserOrder_Cancel_Unauthorized(t *testing.T) {
    repo := newObOrderRepo()
    repo.items[7] = &model.Order{Base: model.Base{ID:7}, UserID:2, Title:"x", Status:model.OrderStatusPending}
    r := setupOrderBadRouter(repo)
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPut, "/user/orders/7/cancel", bytes.NewReader([]byte(`{"reason":"r"}`)))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    if w.Code != http.StatusForbidden { t.Fatalf("%d", w.Code) }
}

