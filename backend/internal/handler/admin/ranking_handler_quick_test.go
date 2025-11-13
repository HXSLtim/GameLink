package admin

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "context"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"

    "gamelink/internal/model"
    rankingrepo "gamelink/internal/repository/ranking"
)

type fakeRankingCfgRepo struct{ store map[uint64]model.RankingCommissionConfig; next uint64 }
func newFakeRankingCfgRepo() *fakeRankingCfgRepo { return &fakeRankingCfgRepo{store: map[uint64]model.RankingCommissionConfig{}, next:1} }
func (f *fakeRankingCfgRepo) CreateConfig(_ context.Context, c *model.RankingCommissionConfig) error { c.ID = f.next; f.next++; f.store[c.ID] = *c; return nil }
func (f *fakeRankingCfgRepo) GetConfig(_ context.Context, id uint64) (*model.RankingCommissionConfig, error) { v, ok := f.store[id]; if !ok { return nil, rankingrepo.ErrNotFound }; c:=v; return &c, nil }
func (f *fakeRankingCfgRepo) GetActiveConfigForMonth(_ context.Context, _ model.RankingType, _ string) (*model.RankingCommissionConfig, error) { return nil, rankingrepo.ErrNotFound }
func (f *fakeRankingCfgRepo) ListConfigs(_ context.Context, _ rankingrepo.RankingCommissionConfigListOptions) ([]model.RankingCommissionConfig, int64, error) { out := make([]model.RankingCommissionConfig,0,len(f.store)); for _, v:= range f.store { out = append(out, v) } ; return out, int64(len(out)), nil }
func (f *fakeRankingCfgRepo) UpdateConfig(_ context.Context, c *model.RankingCommissionConfig) error { f.store[c.ID] = *c; return nil }
func (f *fakeRankingCfgRepo) DeleteConfig(_ context.Context, id uint64) error { delete(f.store, id); return nil }

func setupRankingCfgRouter(repo rankingrepo.RankingCommissionRepository) *gin.Engine {
    r := newTestEngine()
    RegisterRankingCommissionRoutes(r, repo)
    return r
}

func TestRankingCommission_CRUD(t *testing.T) {
    r := setupRankingCfgRouter(newFakeRankingCfgRepo())
    rules := []model.RankingCommissionRule{{RankStart:1, RankEnd:3, CommissionRate:10}}
    payload := map[string]any{"name":"cfg","rankingType":"income","month":"2025-01","rules":rules}
    b, _ := json.Marshal(payload)
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPost, "/admin/ranking-commission/configs", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    assert.Equal(t, http.StatusOK, w.Code)

    w2 := httptest.NewRecorder()
    req2 := httptest.NewRequest(http.MethodGet, "/admin/ranking-commission/configs?page=1&pageSize=20", nil)
    r.ServeHTTP(w2, req2)
    assert.Equal(t, http.StatusOK, w2.Code)

    w3 := httptest.NewRecorder()
    req3 := httptest.NewRequest(http.MethodGet, "/admin/ranking-commission/configs/1", nil)
    r.ServeHTTP(w3, req3)
    assert.Equal(t, http.StatusOK, w3.Code)

    nm := "cfg2"
    up := map[string]any{"name":nm}
    ub, _ := json.Marshal(up)
    w4 := httptest.NewRecorder()
    req4 := httptest.NewRequest(http.MethodPut, "/admin/ranking-commission/configs/1", bytes.NewReader(ub))
    req4.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w4, req4)
    assert.Equal(t, http.StatusOK, w4.Code)

    w5 := httptest.NewRecorder()
    req5 := httptest.NewRequest(http.MethodDelete, "/admin/ranking-commission/configs/1", nil)
    r.ServeHTTP(w5, req5)
    assert.Equal(t, http.StatusOK, w5.Code)
}

func TestRankingCommission_InvalidRules(t *testing.T) {
    r := setupRankingCfgRouter(newFakeRankingCfgRepo())
    rules := []model.RankingCommissionRule{{RankStart:0, RankEnd:3, CommissionRate:200}}
    payload := map[string]any{"name":"bad","rankingType":"income","month":"2025-01","rules":rules}
    b, _ := json.Marshal(payload)
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPost, "/admin/ranking-commission/configs", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRankingCommission_Update_InvalidID(t *testing.T) {
    r := setupRankingCfgRouter(newFakeRankingCfgRepo())
    b := []byte(`{"name":"x"}`)
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPut, "/admin/ranking-commission/configs/abc", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRankingCommission_Get_Delete_InvalidID(t *testing.T) {
    r := setupRankingCfgRouter(newFakeRankingCfgRepo())
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/admin/ranking-commission/configs/abc", nil)
    r.ServeHTTP(w, req)
    assert.Equal(t, http.StatusBadRequest, w.Code)

    w2 := httptest.NewRecorder()
    req2 := httptest.NewRequest(http.MethodDelete, "/admin/ranking-commission/configs/abc", nil)
    r.ServeHTTP(w2, req2)
    assert.Equal(t, http.StatusBadRequest, w2.Code)
}

func TestRankingCommission_Update_InvalidRules(t *testing.T) {
    repo := newFakeRankingCfgRepo()
    r := setupRankingCfgRouter(repo)
    rules := []model.RankingCommissionRule{{RankStart:1, RankEnd:3, CommissionRate:10}}
    payload := map[string]any{"name":"cfg","rankingType":"income","month":"2025-01","rules":rules}
    b, _ := json.Marshal(payload)
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPost, "/admin/ranking-commission/configs", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    assert.Equal(t, http.StatusOK, w.Code)

    bad := []model.RankingCommissionRule{{RankStart:0, RankEnd:3, CommissionRate:200}}
    ub, _ := json.Marshal(map[string]any{"rules": bad})
    w2 := httptest.NewRecorder()
    req2 := httptest.NewRequest(http.MethodPut, "/admin/ranking-commission/configs/1", bytes.NewReader(ub))
    req2.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w2, req2)
    assert.Equal(t, http.StatusBadRequest, w2.Code)
}

func TestRankingCommission_List_Filters(t *testing.T) {
    repo := newFakeRankingCfgRepo()
    repo.CreateConfig(context.Background(), &model.RankingCommissionConfig{Name:"a", RankingType:model.RankingTypeIncome, Period:"monthly", Month:"2025-01", RulesJSON:"[]", IsActive:true})
    repo.CreateConfig(context.Background(), &model.RankingCommissionConfig{Name:"b", RankingType:model.RankingTypeOrderCount, Period:"monthly", Month:"2025-02", RulesJSON:"[]", IsActive:true})
    r := setupRankingCfgRouter(repo)
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/admin/ranking-commission/configs?month=2025-01&rankingType=income&page=1&pageSize=10", nil)
    r.ServeHTTP(w, req)
    assert.Equal(t, http.StatusOK, w.Code)
}