package scheduler

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "gamelink/internal/model"
    commissionrepo "gamelink/internal/repository/commission"
    commissionsvc "gamelink/internal/service/commission"
)

type fakeCommissionRepo struct {
    lastMonth                 string
    listSettlementsCalled     int
    listRecordsCalled         int
    createSettlementCalled    int
    updateRecordCalled        int
}

var _ commissionrepo.CommissionRepository = (*fakeCommissionRepo)(nil)

func (f *fakeCommissionRepo) CreateRule(context.Context, *model.CommissionRule) error { return nil }
func (f *fakeCommissionRepo) GetRule(context.Context, uint64) (*model.CommissionRule, error) { return nil, commissionsvc.ErrNotFound }
func (f *fakeCommissionRepo) GetDefaultRule(context.Context) (*model.CommissionRule, error) { return nil, commissionsvc.ErrNotFound }
func (f *fakeCommissionRepo) GetRuleForOrder(context.Context, *uint64, *uint64, *string) (*model.CommissionRule, error) {
    return nil, commissionsvc.ErrNotFound
}
func (f *fakeCommissionRepo) ListRules(context.Context, commissionrepo.CommissionRuleListOptions) ([]model.CommissionRule, int64, error) {
    return nil, 0, nil
}
func (f *fakeCommissionRepo) UpdateRule(context.Context, *model.CommissionRule) error { return nil }
func (f *fakeCommissionRepo) DeleteRule(context.Context, uint64) error { return nil }

func (f *fakeCommissionRepo) CreateRecord(context.Context, *model.CommissionRecord) error { return nil }
func (f *fakeCommissionRepo) GetRecord(context.Context, uint64) (*model.CommissionRecord, error) { return nil, commissionsvc.ErrNotFound }
func (f *fakeCommissionRepo) GetRecordByOrderID(context.Context, uint64) (*model.CommissionRecord, error) {
    return nil, commissionsvc.ErrNotFound
}
func (f *fakeCommissionRepo) ListRecords(ctx context.Context, opts commissionrepo.CommissionRecordListOptions) ([]model.CommissionRecord, int64, error) {
    f.listRecordsCalled++
    rec := model.CommissionRecord{ID: 1, PlayerID: 100, TotalAmountCents: 1000, CommissionCents: 200, PlayerIncomeCents: 800}
    return []model.CommissionRecord{rec}, 1, nil
}
func (f *fakeCommissionRepo) UpdateRecord(context.Context, *model.CommissionRecord) error {
    f.updateRecordCalled++
    return nil
}

func (f *fakeCommissionRepo) CreateSettlement(context.Context, *model.MonthlySettlement) error {
    f.createSettlementCalled++
    return nil
}
func (f *fakeCommissionRepo) GetSettlement(context.Context, uint64) (*model.MonthlySettlement, error) {
    return nil, commissionsvc.ErrNotFound
}
func (f *fakeCommissionRepo) GetSettlementByPlayerMonth(context.Context, uint64, string) (*model.MonthlySettlement, error) {
    return nil, commissionsvc.ErrNotFound
}
func (f *fakeCommissionRepo) ListSettlements(ctx context.Context, opts commissionrepo.SettlementListOptions) ([]model.MonthlySettlement, int64, error) {
    f.listSettlementsCalled++
    if opts.SettlementMonth != nil {
        f.lastMonth = *opts.SettlementMonth
    }
    return []model.MonthlySettlement{}, 0, nil
}
func (f *fakeCommissionRepo) UpdateSettlement(context.Context, *model.MonthlySettlement) error { return nil }
func (f *fakeCommissionRepo) GetMonthlyStats(context.Context, string) (*commissionrepo.MonthlyStats, error) { return &commissionrepo.MonthlyStats{}, nil }
func (f *fakeCommissionRepo) GetPlayerMonthlyIncome(context.Context, uint64, string) (int64, error) { return 0, nil }

func TestTriggerSettlement_CallsSettleMonth(t *testing.T) {
    repo := &fakeCommissionRepo{}
    svc := commissionsvc.NewCommissionService(repo, nil, nil)
    s := NewSettlementScheduler(svc)

    month := "2025-01"
    err := s.TriggerSettlement(month)
    require.NoError(t, err)

    assert.Equal(t, month, repo.lastMonth)
    assert.Equal(t, 1, repo.listSettlementsCalled)
    assert.Equal(t, 1, repo.createSettlementCalled)
    assert.Equal(t, 1, repo.updateRecordCalled)
}

func TestStartRegistersCronJob(t *testing.T) {
    repo := &fakeCommissionRepo{}
    svc := commissionsvc.NewCommissionService(repo, nil, nil)
    s := NewSettlementScheduler(svc)

    s.Start()
    defer s.Stop()

    entries := s.cron.Entries()
    require.NotEmpty(t, entries)
    assert.True(t, entries[0].Next.After(time.Now()))
}

func TestMonthlySettlement_ComputesLastMonth(t *testing.T) {
    repo := &fakeCommissionRepo{}
    svc := commissionsvc.NewCommissionService(repo, nil, nil)
    s := NewSettlementScheduler(svc)

    s.monthlySettlement()

    expected := time.Now().AddDate(0, -1, 0).Format("2006-01")
    assert.Equal(t, expected, repo.lastMonth)
}

func TestGetNextRunTime_AfterStart(t *testing.T) {
    repo := &fakeCommissionRepo{}
    svc := commissionsvc.NewCommissionService(repo, nil, nil)
    s := NewSettlementScheduler(svc)

    s.Start()
    defer s.Stop()

    next := s.GetNextRunTime()
    assert.False(t, next.IsZero())
}

