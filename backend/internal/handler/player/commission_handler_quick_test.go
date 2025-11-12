package player

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "context"

    "github.com/gin-gonic/gin"

    "gamelink/internal/model"
    commissionrepo "gamelink/internal/repository/commission"
    "gamelink/internal/repository"
    commissionsvc "gamelink/internal/service/commission"
)

type fakeCommissionRepoPlayer struct{}
func (fakeCommissionRepoPlayer) CreateRule(context.Context, *model.CommissionRule) error { return nil }
func (fakeCommissionRepoPlayer) GetRule(context.Context, uint64) (*model.CommissionRule, error) { return nil, repository.ErrNotFound }
func (fakeCommissionRepoPlayer) GetDefaultRule(context.Context) (*model.CommissionRule, error) { return nil, repository.ErrNotFound }
func (fakeCommissionRepoPlayer) GetRuleForOrder(context.Context, *uint64, *uint64, *string) (*model.CommissionRule, error) { return nil, repository.ErrNotFound }
func (fakeCommissionRepoPlayer) ListRules(context.Context, commissionrepo.CommissionRuleListOptions) ([]model.CommissionRule, int64, error) { return nil, 0, nil }
func (fakeCommissionRepoPlayer) UpdateRule(context.Context, *model.CommissionRule) error { return nil }
func (fakeCommissionRepoPlayer) DeleteRule(context.Context, uint64) error { return nil }
func (fakeCommissionRepoPlayer) CreateRecord(context.Context, *model.CommissionRecord) error { return nil }
func (fakeCommissionRepoPlayer) GetRecord(context.Context, uint64) (*model.CommissionRecord, error) { return nil, repository.ErrNotFound }
func (fakeCommissionRepoPlayer) GetRecordByOrderID(context.Context, uint64) (*model.CommissionRecord, error) { return nil, repository.ErrNotFound }
func (fakeCommissionRepoPlayer) ListRecords(context.Context, commissionrepo.CommissionRecordListOptions) ([]model.CommissionRecord, int64, error) { return []model.CommissionRecord{{ID:1, PlayerID:1, TotalAmountCents:1000, CommissionCents:200, PlayerIncomeCents:800}}, 1, nil }
func (fakeCommissionRepoPlayer) UpdateRecord(context.Context, *model.CommissionRecord) error { return nil }
func (fakeCommissionRepoPlayer) CreateSettlement(context.Context, *model.MonthlySettlement) error { return nil }
func (fakeCommissionRepoPlayer) GetSettlement(context.Context, uint64) (*model.MonthlySettlement, error) { return nil, repository.ErrNotFound }
func (fakeCommissionRepoPlayer) GetSettlementByPlayerMonth(context.Context, uint64, string) (*model.MonthlySettlement, error) { return nil, repository.ErrNotFound }
func (fakeCommissionRepoPlayer) ListSettlements(context.Context, commissionrepo.SettlementListOptions) ([]model.MonthlySettlement, int64, error) { return []model.MonthlySettlement{{ID:1, PlayerID:1, SettlementMonth:"2025-01", TotalIncomeCents:800}}, 1, nil }
func (fakeCommissionRepoPlayer) UpdateSettlement(context.Context, *model.MonthlySettlement) error { return nil }
func (fakeCommissionRepoPlayer) GetMonthlyStats(context.Context, string) (*commissionrepo.MonthlyStats, error) { return &commissionrepo.MonthlyStats{}, nil }
func (fakeCommissionRepoPlayer) GetPlayerMonthlyIncome(context.Context, uint64, string) (int64, error) { return 800, nil }

type dummyOrderRepoPlayer struct{}
func (dummyOrderRepoPlayer) Create(context.Context, *model.Order) error { return nil }
func (dummyOrderRepoPlayer) List(context.Context, repository.OrderListOptions) ([]model.Order, int64, error) { return nil, 0, nil }
func (dummyOrderRepoPlayer) Get(context.Context, uint64) (*model.Order, error) { return nil, repository.ErrNotFound }
func (dummyOrderRepoPlayer) Update(context.Context, *model.Order) error { return nil }
func (dummyOrderRepoPlayer) Delete(context.Context, uint64) error { return nil }

type dummyPlayerRepoPlayer struct{}
func (dummyPlayerRepoPlayer) List(context.Context) ([]model.Player, error) { return nil, nil }
func (dummyPlayerRepoPlayer) ListPaged(context.Context, int, int) ([]model.Player, int64, error) { return nil, 0, nil }
func (dummyPlayerRepoPlayer) Get(context.Context, uint64) (*model.Player, error) { return &model.Player{Nickname:"p"}, nil }
func (dummyPlayerRepoPlayer) Create(context.Context, *model.Player) error { return nil }
func (dummyPlayerRepoPlayer) Update(context.Context, *model.Player) error { return nil }
func (dummyPlayerRepoPlayer) Delete(context.Context, uint64) error { return nil }

func setupPlayerCommissionRouter() *gin.Engine {
    r := gin.New()
    r.Use(func(c *gin.Context){ c.Set("user_id", uint64(1)); c.Next() })
    svc := commissionsvc.NewCommissionService(fakeCommissionRepoPlayer{}, dummyOrderRepoPlayer{}, dummyPlayerRepoPlayer{})
    RegisterCommissionRoutes(r, svc, func(c *gin.Context){ c.Next() })
    return r
}

func TestPlayerCommission_Routes(t *testing.T) {
    r := setupPlayerCommissionRouter()
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/player/commission/summary?month=2025-01", nil)
    r.ServeHTTP(w, req)
    if w.Code != http.StatusOK { t.Fatalf("%d", w.Code) }

    w2 := httptest.NewRecorder()
    req2 := httptest.NewRequest(http.MethodGet, "/player/commission/records?page=1&pageSize=10", nil)
    r.ServeHTTP(w2, req2)
    if w2.Code != http.StatusOK { t.Fatalf("%d", w2.Code) }

    w3 := httptest.NewRecorder()
    req3 := httptest.NewRequest(http.MethodGet, "/player/commission/settlements?page=1&pageSize=10", nil)
    r.ServeHTTP(w3, req3)
    if w3.Code != http.StatusOK { t.Fatalf("%d", w3.Code) }
}
