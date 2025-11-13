package player

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "context"

    "github.com/gin-gonic/gin"

    "gamelink/internal/model"
    "gamelink/internal/repository"
    commissionrepo "gamelink/internal/repository/commission"
    serviceitemrepo "gamelink/internal/repository/serviceitem"
    "gamelink/internal/service/gift"
)

type fakeItemRepoPG struct{ items map[uint64]model.ServiceItem }
func (f *fakeItemRepoPG) Create(ctx context.Context, it *model.ServiceItem) error { return nil }
func (f *fakeItemRepoPG) Get(ctx context.Context, id uint64) (*model.ServiceItem, error) { v := f.items[id]; c:=v; return &c, nil }
func (f *fakeItemRepoPG) GetByCode(ctx context.Context, code string) (*model.ServiceItem, error) { return nil, repository.ErrNotFound }
func (f *fakeItemRepoPG) List(ctx context.Context, opts serviceitemrepo.ServiceItemListOptions) ([]model.ServiceItem, int64, error) { out := make([]model.ServiceItem,0,len(f.items)); for _, v:= range f.items { out = append(out, v) } ; return out, int64(len(out)), nil }
func (f *fakeItemRepoPG) Update(ctx context.Context, it *model.ServiceItem) error { return nil }
func (f *fakeItemRepoPG) Delete(ctx context.Context, id uint64) error { return nil }
func (f *fakeItemRepoPG) BatchUpdateStatus(ctx context.Context, ids []uint64, isActive bool) error { return nil }
func (f *fakeItemRepoPG) BatchUpdatePrice(ctx context.Context, ids []uint64, basePriceCents int64) error { return nil }
func (f *fakeItemRepoPG) GetGifts(ctx context.Context, page, pageSize int) ([]model.ServiceItem, int64, error) { out := make([]model.ServiceItem,0,len(f.items)); for _, v:= range f.items { if v.IsGift() && v.IsActive { out = append(out, v) } } ; return out, int64(len(out)), nil }
func (f *fakeItemRepoPG) GetGameServices(ctx context.Context, gameID uint64, subCategory *model.ServiceItemSubCategory) ([]model.ServiceItem, error) { return nil, nil }

type fakeOrderRepoPG struct{ orders []model.Order }
func (f *fakeOrderRepoPG) Create(ctx context.Context, o *model.Order) error { return nil }
func (f *fakeOrderRepoPG) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) { return f.orders, int64(len(f.orders)), nil }
func (f *fakeOrderRepoPG) Get(ctx context.Context, id uint64) (*model.Order, error) { return nil, repository.ErrNotFound }
func (f *fakeOrderRepoPG) Update(ctx context.Context, o *model.Order) error { return nil }
func (f *fakeOrderRepoPG) Delete(ctx context.Context, id uint64) error { return nil }

type fakePlayerRepoPG struct{}
func (fakePlayerRepoPG) List(ctx context.Context) ([]model.Player, error) { return nil, nil }
func (fakePlayerRepoPG) ListPaged(ctx context.Context, page, pageSize int) ([]model.Player, int64, error) { return nil, 0, nil }
func (fakePlayerRepoPG) Get(ctx context.Context, id uint64) (*model.Player, error) { return &model.Player{Base: model.Base{ID: id}, UserID:id, Nickname:"p"}, nil }
func (fakePlayerRepoPG) GetByUserID(ctx context.Context, userID uint64) (*model.Player, error) {
    return &model.Player{Base: model.Base{ID: userID}, UserID: userID, Nickname:"p"}, nil
}
func (fakePlayerRepoPG) Create(ctx context.Context, p *model.Player) error { return nil }
func (fakePlayerRepoPG) Update(ctx context.Context, p *model.Player) error { return nil }
func (fakePlayerRepoPG) Delete(ctx context.Context, id uint64) error { return nil }

type fakeCommissionRepoPG struct{}
func (fakeCommissionRepoPG) CreateRule(ctx context.Context, rule *model.CommissionRule) error { return nil }
func (fakeCommissionRepoPG) GetRule(ctx context.Context, id uint64) (*model.CommissionRule, error) { return nil, repository.ErrNotFound }
func (fakeCommissionRepoPG) GetDefaultRule(ctx context.Context) (*model.CommissionRule, error) { return nil, repository.ErrNotFound }
func (fakeCommissionRepoPG) GetRuleForOrder(ctx context.Context, gameID *uint64, playerID *uint64, serviceType *string) (*model.CommissionRule, error) { return nil, repository.ErrNotFound }
func (fakeCommissionRepoPG) ListRules(ctx context.Context, opts commissionrepo.CommissionRuleListOptions) ([]model.CommissionRule, int64, error) { return nil, 0, nil }
func (fakeCommissionRepoPG) UpdateRule(ctx context.Context, rule *model.CommissionRule) error { return nil }
func (fakeCommissionRepoPG) DeleteRule(ctx context.Context, id uint64) error { return nil }
func (fakeCommissionRepoPG) CreateRecord(ctx context.Context, record *model.CommissionRecord) error { return nil }
func (fakeCommissionRepoPG) GetRecord(ctx context.Context, id uint64) (*model.CommissionRecord, error) { return nil, repository.ErrNotFound }
func (fakeCommissionRepoPG) GetRecordByOrderID(ctx context.Context, orderID uint64) (*model.CommissionRecord, error) { return nil, repository.ErrNotFound }
func (fakeCommissionRepoPG) ListRecords(ctx context.Context, opts commissionrepo.CommissionRecordListOptions) ([]model.CommissionRecord, int64, error) { return nil, 0, nil }
func (fakeCommissionRepoPG) UpdateRecord(ctx context.Context, record *model.CommissionRecord) error { return nil }
func (fakeCommissionRepoPG) CreateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error { return nil }
func (fakeCommissionRepoPG) GetSettlement(ctx context.Context, id uint64) (*model.MonthlySettlement, error) { return nil, repository.ErrNotFound }
func (fakeCommissionRepoPG) GetSettlementByPlayerMonth(ctx context.Context, playerID uint64, month string) (*model.MonthlySettlement, error) { return nil, repository.ErrNotFound }
func (fakeCommissionRepoPG) ListSettlements(ctx context.Context, opts commissionrepo.SettlementListOptions) ([]model.MonthlySettlement, int64, error) { return nil, 0, nil }
func (fakeCommissionRepoPG) UpdateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error { return nil }
func (fakeCommissionRepoPG) GetMonthlyStats(ctx context.Context, month string) (*commissionrepo.MonthlyStats, error) { return &commissionrepo.MonthlyStats{}, nil }
func (fakeCommissionRepoPG) GetPlayerMonthlyIncome(ctx context.Context, playerID uint64, month string) (int64, error) { return 0, nil }

func setupPlayerGiftRouter() *gin.Engine {
    r := gin.New()
    r.Use(func(c *gin.Context){ c.Set("user_id", uint64(1)); c.Next() })
    itemRepo := &fakeItemRepoPG{items: map[uint64]model.ServiceItem{1: {ID:1, ItemCode:"G1", Name:"Gift", Category:"escort", SubCategory:model.SubCategoryGift, BasePriceCents:1000, CommissionRate:0.20, IsActive:true}}}
    pid := uint64(1)
    orders := []model.Order{{Base: model.Base{ID:1}, ItemID:1, PlayerID:&pid, RecipientPlayerID:&pid, Quantity:1, TotalPriceCents:1000, PlayerIncomeCents:800, OrderNo:"GIFT1"}}
    orderRepo := &fakeOrderRepoPG{orders: orders}
    svc := gift.NewGiftService(itemRepo, orderRepo, fakePlayerRepoPG{}, &fakeCommissionRepoPG{})
    RegisterGiftRoutes(r, svc, func(c *gin.Context){ c.Next() })
    return r
}

func TestPlayerGift_ReceivedAndStats(t *testing.T) {
    r := setupPlayerGiftRouter()
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/player/gifts/received?page=1&pageSize=10", nil)
    r.ServeHTTP(w, req)
    if w.Code != http.StatusOK { t.Fatalf("%d", w.Code) }

    w2 := httptest.NewRecorder()
    req2 := httptest.NewRequest(http.MethodGet, "/player/gifts/stats", nil)
    r.ServeHTTP(w2, req2)
    if w2.Code != http.StatusOK { t.Fatalf("%d", w2.Code) }
}
