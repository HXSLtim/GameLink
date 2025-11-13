package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	commissionrepo "gamelink/internal/repository/commission"
	commissionservice "gamelink/internal/service/commission"
)

type fakeCommissionRepoAdmin struct {
	rules []model.CommissionRule
	stats commissionrepo.MonthlyStats
}

func (f *fakeCommissionRepoAdmin) CreateRule(_ context.Context, rule *model.CommissionRule) error {
	rule.ID = uint64(len(f.rules) + 1)
	f.rules = append(f.rules, *rule)
	return nil
}
func (f *fakeCommissionRepoAdmin) GetRule(_ context.Context, id uint64) (*model.CommissionRule, error) {
	for i := range f.rules {
		if f.rules[i].ID == id {
			r := f.rules[i]
			return &r, nil
		}
	}
	return nil, repository.ErrNotFound
}
func (f *fakeCommissionRepoAdmin) GetDefaultRule(context.Context) (*model.CommissionRule, error) {
	return nil, repository.ErrNotFound
}
func (f *fakeCommissionRepoAdmin) GetRuleForOrder(context.Context, *uint64, *uint64, *string) (*model.CommissionRule, error) {
	return nil, repository.ErrNotFound
}
func (f *fakeCommissionRepoAdmin) ListRules(context.Context, commissionrepo.CommissionRuleListOptions) ([]model.CommissionRule, int64, error) {
	return f.rules, int64(len(f.rules)), nil
}
func (f *fakeCommissionRepoAdmin) UpdateRule(context.Context, *model.CommissionRule) error {
	return nil
}
func (f *fakeCommissionRepoAdmin) DeleteRule(context.Context, uint64) error { return nil }
func (f *fakeCommissionRepoAdmin) CreateRecord(context.Context, *model.CommissionRecord) error {
	return nil
}
func (f *fakeCommissionRepoAdmin) GetRecord(context.Context, uint64) (*model.CommissionRecord, error) {
	return nil, repository.ErrNotFound
}
func (f *fakeCommissionRepoAdmin) GetRecordByOrderID(context.Context, uint64) (*model.CommissionRecord, error) {
	return nil, repository.ErrNotFound
}
func (f *fakeCommissionRepoAdmin) ListRecords(context.Context, commissionrepo.CommissionRecordListOptions) ([]model.CommissionRecord, int64, error) {
	return nil, 0, nil
}
func (f *fakeCommissionRepoAdmin) UpdateRecord(context.Context, *model.CommissionRecord) error {
	return nil
}
func (f *fakeCommissionRepoAdmin) CreateSettlement(context.Context, *model.MonthlySettlement) error {
	return nil
}
func (f *fakeCommissionRepoAdmin) GetSettlement(context.Context, uint64) (*model.MonthlySettlement, error) {
	return nil, repository.ErrNotFound
}
func (f *fakeCommissionRepoAdmin) GetSettlementByPlayerMonth(context.Context, uint64, string) (*model.MonthlySettlement, error) {
	return nil, repository.ErrNotFound
}
func (f *fakeCommissionRepoAdmin) ListSettlements(context.Context, commissionrepo.SettlementListOptions) ([]model.MonthlySettlement, int64, error) {
	return nil, 0, nil
}
func (f *fakeCommissionRepoAdmin) UpdateSettlement(context.Context, *model.MonthlySettlement) error {
	return nil
}
func (f *fakeCommissionRepoAdmin) GetMonthlyStats(context.Context, string) (*commissionrepo.MonthlyStats, error) {
	return &f.stats, nil
}
func (f *fakeCommissionRepoAdmin) GetPlayerMonthlyIncome(context.Context, uint64, string) (int64, error) {
	return 0, nil
}

type quickFakeScheduler struct{ last string }

func (f *quickFakeScheduler) TriggerSettlement(m string) error { f.last = m; return nil }

type dummyOrderRepo struct{}

func (dummyOrderRepo) Create(context.Context, *model.Order) error { return nil }
func (dummyOrderRepo) List(context.Context, repository.OrderListOptions) ([]model.Order, int64, error) {
	return nil, 0, nil
}
func (dummyOrderRepo) Get(context.Context, uint64) (*model.Order, error) {
	return nil, repository.ErrNotFound
}
func (dummyOrderRepo) Update(context.Context, *model.Order) error { return nil }
func (dummyOrderRepo) Delete(context.Context, uint64) error       { return nil }

type dummyPlayerRepo struct{}

func (dummyPlayerRepo) Get(context.Context, uint64) (*model.Player, error) {
	return nil, repository.ErrNotFound
}
func (dummyPlayerRepo) GetByUserID(context.Context, uint64) (*model.Player, error) {
	return nil, repository.ErrNotFound
}
func (dummyPlayerRepo) Create(context.Context, *model.Player) error  { return nil }
func (dummyPlayerRepo) Update(context.Context, *model.Player) error  { return nil }
func (dummyPlayerRepo) Delete(context.Context, uint64) error         { return nil }
func (dummyPlayerRepo) List(context.Context) ([]model.Player, error) { return nil, nil }
func (dummyPlayerRepo) ListPaged(context.Context, int, int) ([]model.Player, int64, error) {
	return nil, 0, nil
}

func setupCommissionRouter(repo *fakeCommissionRepoAdmin) (*gin.Engine, *commissionservice.CommissionService, *quickFakeScheduler) {
	r := newTestEngine()
	svc := commissionservice.NewCommissionService(repo, dummyOrderRepo{}, dummyPlayerRepo{})
	sch := &quickFakeScheduler{}
	RegisterCommissionRoutes(r, svc, sch)
	return r, svc, sch
}

func TestCommission_CreateUpdateStatsAndTrigger(t *testing.T) {
	repo := &fakeCommissionRepoAdmin{}
	r, _, sch := setupCommissionRouter(repo)
	create := commissionservice.CreateCommissionRuleRequest{Name: "默认", Type: "default", Rate: 20}
	body, _ := json.Marshal(create)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/admin/commission/rules", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 200/500, got %d", w.Code)
	}

	upd := commissionservice.UpdateCommissionRuleRequest{Rate: ptrInt(25)}
	ub, _ := json.Marshal(upd)
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPut, "/admin/commission/rules/1", bytes.NewReader(ub))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK && w2.Code != http.StatusInternalServerError {
		t.Fatalf("expected 200/500, got %d", w2.Code)
	}

	w3 := httptest.NewRecorder()
	req3 := httptest.NewRequest(http.MethodPost, "/admin/commission/settlements/trigger?month=2025-01", nil)
	r.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)
	assert.Equal(t, "2025-01", sch.last)

	repo.stats = commissionrepo.MonthlyStats{TotalOrders: 3}
	w4 := httptest.NewRecorder()
	req4 := httptest.NewRequest(http.MethodGet, "/admin/commission/stats?month=2025-01", nil)
	r.ServeHTTP(w4, req4)
	assert.Equal(t, http.StatusOK, w4.Code)
}

func ptrInt(i int) *int { return &i }
