package user

import (
    "context"
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"

    "gamelink/internal/model"
    "gamelink/internal/repository"
    commissionrepo "gamelink/internal/repository/commission"
    serviceitemrepo "gamelink/internal/repository/serviceitem"
    giftsvc "gamelink/internal/service/gift"
    itemsvc "gamelink/internal/service/item"
)

type fakeItemRepoGift struct{ items map[uint64]model.ServiceItem }
func (f *fakeItemRepoGift) Create(ctx context.Context, _ *model.ServiceItem) error { return nil }
func (f *fakeItemRepoGift) Get(ctx context.Context, id uint64) (*model.ServiceItem, error) { v := f.items[id]; c:=v; return &c, nil }
func (f *fakeItemRepoGift) GetByCode(ctx context.Context, code string) (*model.ServiceItem, error) { _ = code; return nil, repository.ErrNotFound }
func (f *fakeItemRepoGift) List(ctx context.Context, opts serviceitemrepo.ServiceItemListOptions) ([]model.ServiceItem, int64, error) { _ = opts; out := make([]model.ServiceItem,0,len(f.items)); for _, v:= range f.items { out = append(out, v) }; return out, int64(len(out)), nil }
func (f *fakeItemRepoGift) Update(ctx context.Context, _ *model.ServiceItem) error { return nil }
func (f *fakeItemRepoGift) Delete(ctx context.Context, _ uint64) error { return nil }
func (f *fakeItemRepoGift) BatchUpdateStatus(ctx context.Context, _ []uint64, _ bool) error { return nil }
func (f *fakeItemRepoGift) BatchUpdatePrice(ctx context.Context, _ []uint64, _ int64) error { return nil }
func (f *fakeItemRepoGift) GetGifts(ctx context.Context, page, pageSize int) ([]model.ServiceItem, int64, error) { _=page; _=pageSize; out := make([]model.ServiceItem,0,len(f.items)); for _, v:= range f.items { if v.IsGift() && v.IsActive { out = append(out, v) } } ; return out, int64(len(out)), nil }
func (f *fakeItemRepoGift) GetGameServices(ctx context.Context, _ uint64, _ *model.ServiceItemSubCategory) ([]model.ServiceItem, error) { return nil, nil }

type dummyOrderRepoGift struct{ last *model.Order }
func (d *dummyOrderRepoGift) Create(ctx context.Context, o *model.Order) error { _=ctx; d.last = o; o.ID = 1; return nil }
func (d *dummyOrderRepoGift) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) { _=ctx; _=opts; return nil, 0, nil }
func (d *dummyOrderRepoGift) Get(ctx context.Context, _ uint64) (*model.Order, error) { _=ctx; return d.last, nil }
func (d *dummyOrderRepoGift) Update(ctx context.Context, _ *model.Order) error { _=ctx; return nil }
func (d *dummyOrderRepoGift) Delete(ctx context.Context, _ uint64) error { _=ctx; return nil }

type dummyPlayerRepoGift struct{}
func (dummyPlayerRepoGift) List(ctx context.Context) ([]model.Player, error) { _=ctx; return nil, nil }
func (dummyPlayerRepoGift) ListPaged(ctx context.Context, _ int, _ int) ([]model.Player, int64, error) { _=ctx; return nil, 0, nil }
func (dummyPlayerRepoGift) Get(ctx context.Context, _ uint64) (*model.Player, error) { _=ctx; return &model.Player{Nickname:"p"}, nil }
func (dummyPlayerRepoGift) GetByUserID(ctx context.Context, userID uint64) (*model.Player, error) {
    _ = ctx
    return &model.Player{Base: model.Base{ID: userID}, UserID: userID, Nickname:"p"}, nil
}
func (dummyPlayerRepoGift) Create(ctx context.Context, _ *model.Player) error { _=ctx; return nil }
func (dummyPlayerRepoGift) Update(ctx context.Context, _ *model.Player) error { _=ctx; return nil }
func (dummyPlayerRepoGift) Delete(ctx context.Context, _ uint64) error { _=ctx; return nil }

type dummyCommissionRepoGift struct{}
func (dummyCommissionRepoGift) CreateRule(ctx context.Context, _ *model.CommissionRule) error { _=ctx; return nil }
func (dummyCommissionRepoGift) GetRule(ctx context.Context, _ uint64) (*model.CommissionRule, error) { _=ctx; return nil, repository.ErrNotFound }
func (dummyCommissionRepoGift) GetDefaultRule(ctx context.Context) (*model.CommissionRule, error) { _=ctx; return nil, repository.ErrNotFound }
func (dummyCommissionRepoGift) GetRuleForOrder(ctx context.Context, _ *uint64, _ *uint64, _ *string) (*model.CommissionRule, error) { _=ctx; return nil, repository.ErrNotFound }
func (dummyCommissionRepoGift) ListRules(ctx context.Context, _ commissionrepo.CommissionRuleListOptions) ([]model.CommissionRule, int64, error) { _=ctx; return nil, 0, nil }
func (dummyCommissionRepoGift) UpdateRule(ctx context.Context, _ *model.CommissionRule) error { _=ctx; return nil }
func (dummyCommissionRepoGift) DeleteRule(ctx context.Context, _ uint64) error { _=ctx; return nil }
func (dummyCommissionRepoGift) CreateRecord(ctx context.Context, _ *model.CommissionRecord) error { _=ctx; return nil }
func (dummyCommissionRepoGift) GetRecord(ctx context.Context, _ uint64) (*model.CommissionRecord, error) { _=ctx; return nil, repository.ErrNotFound }
func (dummyCommissionRepoGift) GetRecordByOrderID(ctx context.Context, _ uint64) (*model.CommissionRecord, error) { _=ctx; return nil, repository.ErrNotFound }
func (dummyCommissionRepoGift) ListRecords(ctx context.Context, _ commissionrepo.CommissionRecordListOptions) ([]model.CommissionRecord, int64, error) { _=ctx; return nil, 0, nil }
func (dummyCommissionRepoGift) UpdateRecord(ctx context.Context, _ *model.CommissionRecord) error { _=ctx; return nil }
func (dummyCommissionRepoGift) CreateSettlement(ctx context.Context, _ *model.MonthlySettlement) error { _=ctx; return nil }
func (dummyCommissionRepoGift) GetSettlement(ctx context.Context, _ uint64) (*model.MonthlySettlement, error) { _=ctx; return nil, repository.ErrNotFound }
func (dummyCommissionRepoGift) GetSettlementByPlayerMonth(ctx context.Context, _ uint64, _ string) (*model.MonthlySettlement, error) { _=ctx; return nil, repository.ErrNotFound }
func (dummyCommissionRepoGift) ListSettlements(ctx context.Context, _ commissionrepo.SettlementListOptions) ([]model.MonthlySettlement, int64, error) { _=ctx; return nil, 0, nil }
func (dummyCommissionRepoGift) UpdateSettlement(ctx context.Context, _ *model.MonthlySettlement) error { _=ctx; return nil }
func (dummyCommissionRepoGift) GetMonthlyStats(ctx context.Context, _ string) (*commissionrepo.MonthlyStats, error) { _=ctx; return &commissionrepo.MonthlyStats{}, nil }
func (dummyCommissionRepoGift) GetPlayerMonthlyIncome(ctx context.Context, _ uint64, _ string) (int64, error) { _=ctx; return 0, nil }

func setupUserGiftRouter() *gin.Engine {
    r := gin.New()
    r.Use(func(c *gin.Context){ c.Set("user_id", uint64(1)); c.Next() })
    itemRepo := &fakeItemRepoGift{items: map[uint64]model.ServiceItem{1: {ID:1, ItemCode:"G1", Name:"Gift", Category:"escort", SubCategory:model.SubCategoryGift, BasePriceCents:1000, CommissionRate:0.2, IsActive:true}}}
    itemSvc := itemsvc.NewServiceItemService(itemRepo, dummyGameRepo{}, dummyPlayerRepoGift{})
    giftSvc := giftsvc.NewGiftService(itemRepo, &dummyOrderRepoGift{}, dummyPlayerRepoGift{}, &dummyCommissionRepoGift{})
    RegisterGiftRoutes(r, giftSvc, itemSvc, func(c *gin.Context){ c.Next() })
    return r
}

type dummyGameRepo struct{}
func (dummyGameRepo) List(ctx context.Context) ([]model.Game, error) { _=ctx; return nil, nil }
func (dummyGameRepo) ListPaged(ctx context.Context, _ int, _ int) ([]model.Game, int64, error) { _=ctx; return nil, 0, nil }
func (dummyGameRepo) Get(ctx context.Context, _ uint64) (*model.Game, error) { _=ctx; return &model.Game{Name:"g"}, nil }
func (dummyGameRepo) Create(ctx context.Context, _ *model.Game) error { _=ctx; return nil }
func (dummyGameRepo) Update(ctx context.Context, _ *model.Game) error { _=ctx; return nil }
func (dummyGameRepo) Delete(ctx context.Context, _ uint64) error { _=ctx; return nil }

func TestUserGift_ListAndSend(t *testing.T) {
    r := setupUserGiftRouter()
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/user/gifts?page=1&pageSize=10", nil)
    r.ServeHTTP(w, req)
    if w.Code != http.StatusOK { t.Fatalf("%d", w.Code) }

    payload := giftsvc.SendGiftRequest{PlayerID:1, GiftItemID:1, Quantity:1, Message:"m"}
    b, _ := json.Marshal(payload)
    w2 := httptest.NewRecorder()
    req2 := httptest.NewRequest(http.MethodPost, "/user/gifts/send", bytes.NewReader(b))
    req2.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w2, req2)
    if w2.Code != http.StatusOK { t.Fatalf("%d", w2.Code) }
}
