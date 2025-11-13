package player

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
    withdrawrepo "gamelink/internal/repository/withdraw"
    "gamelink/internal/service/earnings"
)

type fakePlayerRepoEarn struct{}
func (fakePlayerRepoEarn) List(ctx context.Context) ([]model.Player, error) { return nil, nil }
func (fakePlayerRepoEarn) ListPaged(ctx context.Context, page int, pageSize int) ([]model.Player, int64, error) { return []model.Player{{Base: model.Base{ID:1}, UserID:1, Nickname:"p"}}, 1, nil }
func (fakePlayerRepoEarn) Get(ctx context.Context, id uint64) (*model.Player, error) { return &model.Player{Base: model.Base{ID:1}, UserID:1, Nickname:"p"}, nil }
func (fakePlayerRepoEarn) Create(ctx context.Context, p *model.Player) error { return nil }
func (fakePlayerRepoEarn) Update(ctx context.Context, p *model.Player) error { return nil }
func (fakePlayerRepoEarn) Delete(ctx context.Context, id uint64) error { return nil }

type fakeOrderRepoEarn struct{}
func (fakeOrderRepoEarn) Create(ctx context.Context, o *model.Order) error { return nil }
func (fakeOrderRepoEarn) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) {
    _ = ctx
    _ = opts
    pid := uint64(1)
    return []model.Order{{Base: model.Base{ID:1}, PlayerID:&pid, TotalPriceCents:1000, Status:model.OrderStatusCompleted}}, 1, nil
}
func (fakeOrderRepoEarn) Get(ctx context.Context, id uint64) (*model.Order, error) { return nil, repository.ErrNotFound }
func (fakeOrderRepoEarn) Update(ctx context.Context, o *model.Order) error { return nil }
func (fakeOrderRepoEarn) Delete(ctx context.Context, id uint64) error { return nil }

type fakeWithdrawRepoEarn struct{ balance *withdrawrepo.PlayerBalance; list []model.Withdraw }
func (f *fakeWithdrawRepoEarn) Create(ctx context.Context, w *model.Withdraw) error { w.ID = 1; f.list = append(f.list, *w); return nil }
func (f *fakeWithdrawRepoEarn) Get(ctx context.Context, id uint64) (*model.Withdraw, error) { return nil, repository.ErrNotFound }
func (f *fakeWithdrawRepoEarn) Update(ctx context.Context, w *model.Withdraw) error { return nil }
func (f *fakeWithdrawRepoEarn) List(ctx context.Context, opts withdrawrepo.WithdrawListOptions) ([]model.Withdraw, int64, error) {
    _ = ctx
    _ = opts
    if f.list == nil { return []model.Withdraw{}, 0, nil }
    return f.list, int64(len(f.list)), nil
}
func (f *fakeWithdrawRepoEarn) GetPlayerBalance(ctx context.Context, id uint64) (*withdrawrepo.PlayerBalance, error) {
    if f.balance != nil { return f.balance, nil }
    return &withdrawrepo.PlayerBalance{TotalEarnings: 5000, WithdrawTotal: 0, PendingWithdraw: 0, AvailableBalance: 4000, PendingBalance: 1000}, nil
}

func setupEarningsRouter(wr *fakeWithdrawRepoEarn) *gin.Engine {
    r := gin.New()
    r.Use(func(c *gin.Context){ c.Set("user_id", uint64(1)); c.Next() })
    svc := earnings.NewEarningsService(fakePlayerRepoEarn{}, fakeOrderRepoEarn{}, wr)
    RegisterEarningsRoutes(r, svc, func(c *gin.Context){ c.Next() })
    return r
}

func TestEarnings_SummaryTrend(t *testing.T) {
    r := setupEarningsRouter(&fakeWithdrawRepoEarn{})
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/player/earnings/summary", nil)
    r.ServeHTTP(w, req)
    if w.Code != http.StatusOK { t.Fatalf("%d", w.Code) }

    w2 := httptest.NewRecorder()
    req2 := httptest.NewRequest(http.MethodGet, "/player/earnings/trend?days=30", nil)
    r.ServeHTTP(w2, req2)
    if w2.Code != http.StatusOK { t.Fatalf("%d", w2.Code) }

    w3 := httptest.NewRecorder()
    req3 := httptest.NewRequest(http.MethodGet, "/player/earnings/trend?days=1", nil)
    r.ServeHTTP(w3, req3)
    if w3.Code != http.StatusBadRequest { t.Fatalf("%d", w3.Code) }
}

func TestEarnings_WithdrawAndHistory(t *testing.T) {
    wr := &fakeWithdrawRepoEarn{}
    wr.balance = &withdrawrepo.PlayerBalance{AvailableBalance: 50000}
    r := setupEarningsRouter(wr)
    payload := earnings.WithdrawRequest{AmountCents: 10000, Method: "alipay", AccountInfo: "acc"}
    b, _ := json.Marshal(payload)
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPost, "/player/earnings/withdraw", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    if w.Code != http.StatusOK { t.Fatalf("%d", w.Code) }

    wr.balance = &withdrawrepo.PlayerBalance{AvailableBalance: 500}
    b2, _ := json.Marshal(payload)
    w2 := httptest.NewRecorder()
    req2 := httptest.NewRequest(http.MethodPost, "/player/earnings/withdraw", bytes.NewReader(b2))
    req2.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w2, req2)
    if w2.Code != http.StatusBadRequest { t.Fatalf("%d", w2.Code) }

    tm := time.Now()
    wr.list = []model.Withdraw{{ID:2, PlayerID:1, UserID:1, AmountCents:1000, Method:model.WithdrawMethodAlipay, Status:model.WithdrawStatusPending, CreatedAt: tm}}
    w3 := httptest.NewRecorder()
    req3 := httptest.NewRequest(http.MethodGet, "/player/earnings/withdraw-history?page=1&pageSize=10", nil)
    r.ServeHTTP(w3, req3)
    if w3.Code != http.StatusOK { t.Fatalf("%d", w3.Code) }
}