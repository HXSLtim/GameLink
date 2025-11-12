package admin

import (
    "context"
    "net/http"
    "net/http/httptest"
    "testing"
    "bytes"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"

    "gamelink/internal/model"
    "gamelink/internal/repository"
    withdrawrepo "gamelink/internal/repository/withdraw"
)

type fakeWithdrawRepo struct{ items map[uint64]model.Withdraw }
func newFakeWithdrawRepo() *fakeWithdrawRepo { return &fakeWithdrawRepo{items: map[uint64]model.Withdraw{}} }
func (f *fakeWithdrawRepo) Create(ctx context.Context, _ *model.Withdraw) error { return nil }
func (f *fakeWithdrawRepo) Get(ctx context.Context, id uint64) (*model.Withdraw, error) { v, ok := f.items[id]; if !ok { return nil, repository.ErrNotFound }; c:=v; return &c, nil }
func (f *fakeWithdrawRepo) Update(ctx context.Context, w *model.Withdraw) error { f.items[w.ID] = *w; return nil }
func (f *fakeWithdrawRepo) List(ctx context.Context, _ withdrawrepo.WithdrawListOptions) ([]model.Withdraw, int64, error) { out := make([]model.Withdraw,0,len(f.items)); for _, v:= range f.items { out = append(out, v) } ; return out, int64(len(out)), nil }
func (f *fakeWithdrawRepo) GetPlayerBalance(ctx context.Context, _ uint64) (*withdrawrepo.PlayerBalance, error) { return &withdrawrepo.PlayerBalance{}, nil }

func setupWithdrawRouter(repo withdrawrepo.WithdrawRepository) *gin.Engine {
    r := newTestEngine()
    r.Use(func(c *gin.Context){ c.Set("user_id", uint64(1)); c.Next() })
    RegisterWithdrawRoutes(r, repo)
    return r
}

func TestWithdraw_List_Get_Approve_Reject_Complete(t *testing.T) {
    repo := newFakeWithdrawRepo()
    repo.items[1] = model.Withdraw{ID:1, PlayerID: 7, UserID: 1, AmountCents: 1000, Method: model.WithdrawMethodAlipay, AccountInfo: "x", Status: model.WithdrawStatusPending}
    r := setupWithdrawRouter(repo)

    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/admin/withdraws?page=1&pageSize=20", nil)
    r.ServeHTTP(w, req)
    assert.Equal(t, http.StatusOK, w.Code)

    w2 := httptest.NewRecorder()
    req2 := httptest.NewRequest(http.MethodGet, "/admin/withdraws/1", nil)
    r.ServeHTTP(w2, req2)
    assert.Equal(t, http.StatusOK, w2.Code)

    w3 := httptest.NewRecorder()
    req3 := httptest.NewRequest(http.MethodPost, "/admin/withdraws/1/approve", nil)
    r.ServeHTTP(w3, req3)
    assert.Equal(t, http.StatusOK, w3.Code)

    repo.items[2] = model.Withdraw{ID:2, PlayerID: 8, UserID: 1, AmountCents: 2000, Method: model.WithdrawMethodAlipay, AccountInfo: "x", Status: model.WithdrawStatusPending}
    w4 := httptest.NewRecorder()
    req4 := httptest.NewRequest(http.MethodPost, "/admin/withdraws/2/reject", nil)
    r.ServeHTTP(w4, req4)
    assert.Equal(t, http.StatusBadRequest, w4.Code)

    rb := httptest.NewRecorder()
    reqb := httptest.NewRequest(http.MethodPost, "/admin/withdraws/2/reject", bytes.NewReader([]byte("{\"reason\":\"dup\"}")))
    reqb.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(rb, reqb)
    assert.Equal(t, http.StatusOK, rb.Code)

    repo.items[3] = model.Withdraw{ID:3, PlayerID: 9, UserID: 1, AmountCents: 3000, Method: model.WithdrawMethodAlipay, AccountInfo: "x", Status: model.WithdrawStatusApproved}
    w5 := httptest.NewRecorder()
    req5 := httptest.NewRequest(http.MethodPost, "/admin/withdraws/3/complete", nil)
    r.ServeHTTP(w5, req5)
    assert.Equal(t, http.StatusOK, w5.Code)
}

type errWithdrawRepo struct{ fakeWithdrawRepo }
func (e *errWithdrawRepo) Update(ctx context.Context, w *model.Withdraw) error { return assert.AnError }

func TestWithdraw_InvalidID_And_StatusErrors(t *testing.T) {
    repo := newFakeWithdrawRepo()
    repo.items[10] = model.Withdraw{ID:10, PlayerID: 1, UserID: 1, AmountCents: 100, Method: model.WithdrawMethodAlipay, AccountInfo: "x", Status: model.WithdrawStatusApproved}
    repo.items[11] = model.Withdraw{ID:11, PlayerID: 1, UserID: 1, AmountCents: 100, Method: model.WithdrawMethodAlipay, AccountInfo: "x", Status: model.WithdrawStatusPending}
    r := setupWithdrawRouter(repo)

    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/admin/withdraws/abc", nil)
    r.ServeHTTP(w, req)
    assert.Equal(t, http.StatusBadRequest, w.Code)

    w2 := httptest.NewRecorder()
    req2 := httptest.NewRequest(http.MethodPost, "/admin/withdraws/10/approve", nil)
    r.ServeHTTP(w2, req2)
    assert.Equal(t, http.StatusBadRequest, w2.Code)

    w3 := httptest.NewRecorder()
    req3 := httptest.NewRequest(http.MethodPost, "/admin/withdraws/11/complete", nil)
    r.ServeHTTP(w3, req3)
    assert.Equal(t, http.StatusBadRequest, w3.Code)
}

func TestWithdraw_UpdateError(t *testing.T) {
    base := newFakeWithdrawRepo()
    base.items[20] = model.Withdraw{ID:20, PlayerID: 1, UserID: 1, AmountCents: 100, Method: model.WithdrawMethodAlipay, AccountInfo: "x", Status: model.WithdrawStatusPending}
    repo := &errWithdrawRepo{fakeWithdrawRepo: *base}
    r := setupWithdrawRouter(repo)
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPost, "/admin/withdraws/20/approve", nil)
    r.ServeHTTP(w, req)
    assert.Equal(t, http.StatusInternalServerError, w.Code)
}
