package user

import (
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

type invalidOrderRepo struct{ items map[uint64]*model.Order }
func newInvalidOrderRepo() *invalidOrderRepo { return &invalidOrderRepo{items: map[uint64]*model.Order{}} }
func (r *invalidOrderRepo) Create(ctx context.Context, o *model.Order) error { r.items[o.ID]=o; return nil }
func (r *invalidOrderRepo) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) { out:=[]model.Order{}; for _, v:= range r.items { out = append(out, *v) } ; return out, int64(len(out)), nil }
func (r *invalidOrderRepo) Get(ctx context.Context, id uint64) (*model.Order, error) { v:=r.items[id]; if v==nil { return nil, repository.ErrNotFound } ; return v, nil }
func (r *invalidOrderRepo) Update(ctx context.Context, o *model.Order) error { r.items[o.ID]=o; return nil }
func (r *invalidOrderRepo) Delete(ctx context.Context, id uint64) error { delete(r.items, id); return nil }

type invalidPlayers struct{}
func (invalidPlayers) List(context.Context) ([]model.Player, error) { return nil, nil }
func (invalidPlayers) ListPaged(context.Context, int, int) ([]model.Player, int64, error) { return []model.Player{{Base: model.Base{ID:1}, UserID:1, Nickname:"p", HourlyRateCents: 1000}}, 1, nil }
func (invalidPlayers) Get(context.Context, uint64) (*model.Player, error) { return &model.Player{Base: model.Base{ID:1}, UserID:1, Nickname:"p", HourlyRateCents: 1000}, nil }
func (invalidPlayers) Create(context.Context, *model.Player) error { return nil }
func (invalidPlayers) Update(context.Context, *model.Player) error { return nil }
func (invalidPlayers) Delete(context.Context, uint64) error { return nil }
func (invalidPlayers) GetByUserID(context.Context, uint64) (*model.Player, error) { return nil, repository.ErrNotFound }

type invalidUsers struct{}
func (invalidUsers) List(context.Context) ([]model.User, error) { return nil, nil }
func (invalidUsers) ListPaged(context.Context, int, int) ([]model.User, int64, error) { return nil, 0, nil }
func (invalidUsers) ListWithFilters(context.Context, repository.UserListOptions) ([]model.User, int64, error) { return nil, 0, nil }
func (invalidUsers) Get(context.Context, uint64) (*model.User, error) { return &model.User{Base: model.Base{ID:1}, AvatarURL:"a"}, nil }
func (invalidUsers) GetByPhone(context.Context, string) (*model.User, error) { return nil, repository.ErrNotFound }
func (invalidUsers) FindByEmail(context.Context, string) (*model.User, error) { return nil, repository.ErrNotFound }
func (invalidUsers) FindByPhone(context.Context, string) (*model.User, error) { return nil, repository.ErrNotFound }
func (invalidUsers) Create(context.Context, *model.User) error { return nil }
func (invalidUsers) Update(context.Context, *model.User) error { return nil }
func (invalidUsers) Delete(context.Context, uint64) error { return nil }

type invalidGames struct{}
func (invalidGames) List(context.Context) ([]model.Game, error) { return nil, nil }
func (invalidGames) ListPaged(context.Context, int, int) ([]model.Game, int64, error) { return nil, 0, nil }
func (invalidGames) Get(context.Context, uint64) (*model.Game, error) { return &model.Game{Base: model.Base{ID:1}, Name:"g"}, nil }
func (invalidGames) Create(context.Context, *model.Game) error { return nil }
func (invalidGames) Update(context.Context, *model.Game) error { return nil }
func (invalidGames) Delete(context.Context, uint64) error { return nil }

type invalidPayments struct{}
func (invalidPayments) Create(context.Context, *model.Payment) error { return nil }
func (invalidPayments) List(context.Context, repository.PaymentListOptions) ([]model.Payment, int64, error) { return nil, 0, nil }
func (invalidPayments) Get(context.Context, uint64) (*model.Payment, error) { return nil, repository.ErrNotFound }
func (invalidPayments) Update(context.Context, *model.Payment) error { return nil }
func (invalidPayments) Delete(context.Context, uint64) error { return nil }

type invalidReviews struct{}
func (invalidReviews) List(context.Context, repository.ReviewListOptions) ([]model.Review, int64, error) { return nil, 0, nil }
func (invalidReviews) Get(context.Context, uint64) (*model.Review, error) { return nil, repository.ErrNotFound }
func (invalidReviews) Create(context.Context, *model.Review) error { return nil }
func (invalidReviews) Update(context.Context, *model.Review) error { return nil }
func (invalidReviews) Delete(context.Context, uint64) error { return nil }

type invalidCommissions struct{}
func (invalidCommissions) CreateRule(context.Context, *model.CommissionRule) error { return nil }
func (invalidCommissions) GetRule(context.Context, uint64) (*model.CommissionRule, error) { return nil, repository.ErrNotFound }
func (invalidCommissions) GetDefaultRule(context.Context) (*model.CommissionRule, error) { return &model.CommissionRule{Rate:20}, nil }
func (invalidCommissions) GetRuleForOrder(context.Context, *uint64, *uint64, *string) (*model.CommissionRule, error) { return nil, repository.ErrNotFound }
func (invalidCommissions) ListRules(context.Context, commissionrepo.CommissionRuleListOptions) ([]model.CommissionRule, int64, error) { return nil, 0, nil }
func (invalidCommissions) UpdateRule(context.Context, *model.CommissionRule) error { return nil }
func (invalidCommissions) DeleteRule(context.Context, uint64) error { return nil }
func (invalidCommissions) CreateRecord(context.Context, *model.CommissionRecord) error { return nil }
func (invalidCommissions) GetRecord(context.Context, uint64) (*model.CommissionRecord, error) { return nil, repository.ErrNotFound }
func (invalidCommissions) GetRecordByOrderID(context.Context, uint64) (*model.CommissionRecord, error) { return nil, repository.ErrNotFound }
func (invalidCommissions) ListRecords(context.Context, commissionrepo.CommissionRecordListOptions) ([]model.CommissionRecord, int64, error) { return nil, 0, nil }
func (invalidCommissions) UpdateRecord(context.Context, *model.CommissionRecord) error { return nil }
func (invalidCommissions) CreateSettlement(context.Context, *model.MonthlySettlement) error { return nil }
func (invalidCommissions) GetSettlement(context.Context, uint64) (*model.MonthlySettlement, error) { return nil, repository.ErrNotFound }
func (invalidCommissions) GetSettlementByPlayerMonth(context.Context, uint64, string) (*model.MonthlySettlement, error) { return nil, repository.ErrNotFound }
func (invalidCommissions) ListSettlements(context.Context, commissionrepo.SettlementListOptions) ([]model.MonthlySettlement, int64, error) { return nil, 0, nil }
func (invalidCommissions) UpdateSettlement(context.Context, *model.MonthlySettlement) error { return nil }
func (invalidCommissions) GetMonthlyStats(context.Context, string) (*commissionrepo.MonthlyStats, error) { return &commissionrepo.MonthlyStats{}, nil }
func (invalidCommissions) GetPlayerMonthlyIncome(context.Context, uint64, string) (int64, error) { return 0, nil }

func setupOrderInvalidRouter(repo *invalidOrderRepo) *gin.Engine {
    r := gin.New()
    r.Use(func(c *gin.Context){ c.Set("user_id", uint64(1)); c.Next() })
    svc := ordersvc.NewOrderService(repo, invalidPlayers{}, invalidUsers{}, invalidGames{}, invalidPayments{}, invalidReviews{}, invalidCommissions{})
    RegisterOrderRoutes(r, svc, func(c *gin.Context){ c.Next() })
    return r
}

func TestUserOrder_Detail_Unauthorized(t *testing.T) {
    repo := newInvalidOrderRepo()
    repo.items[10] = &model.Order{Base: model.Base{ID:10}, UserID:2, Title:"x", Status:model.OrderStatusPending}
    r := setupOrderInvalidRouter(repo)
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/user/orders/10", nil)
    r.ServeHTTP(w, req)
    if w.Code != http.StatusForbidden { t.Fatalf("%d", w.Code) }
}

func TestUserOrder_Complete_InvalidTransition(t *testing.T) {
    repo := newInvalidOrderRepo()
    repo.items[11] = &model.Order{Base: model.Base{ID:11}, UserID:1, Title:"x", Status:model.OrderStatusPending}
    r := setupOrderInvalidRouter(repo)
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPut, "/user/orders/11/complete", nil)
    r.ServeHTTP(w, req)
    if w.Code != http.StatusBadRequest { t.Fatalf("%d", w.Code) }
}