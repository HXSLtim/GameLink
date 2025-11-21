package commission

import (
	"context"
	"testing"
	"time"

	"gamelink/internal/model"
	commissionrepo "gamelink/internal/repository/commission"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCommissionService_GetPlayerCommissionSummary(t *testing.T) {
	ctx := context.Background()
	commissionRepo := new(MockCommissionRepo)
	orderRepo := new(MockOrderRepo)
	playerRepo := new(MockPlayerRepo)

	svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

	playerID := uint64(42)
	month := "2025-01"

	commissionRepo.On("GetPlayerMonthlyIncome", ctx, playerID, month).Return(int64(120000), nil)
	commissionRepo.
		On("ListRecords", ctx, mock.MatchedBy(func(opts commissionrepo.CommissionRecordListOptions) bool {
			return opts.PlayerID != nil && *opts.PlayerID == playerID && opts.Page == 1 && opts.PageSize == 1
		})).
		Return([]model.CommissionRecord{
			{CommissionCents: 1000, PlayerIncomeCents: 4000},
		}, int64(3), nil).Once()
	commissionRepo.
		On("ListRecords", ctx, mock.MatchedBy(func(opts commissionrepo.CommissionRecordListOptions) bool {
			return opts.PlayerID != nil && *opts.PlayerID == playerID && opts.Page == 1 && opts.PageSize == 10000
		})).
		Return([]model.CommissionRecord{
			{CommissionCents: 1000, PlayerIncomeCents: 4000},
			{CommissionCents: 1500, PlayerIncomeCents: 6000},
		}, int64(3), nil).Once()

	resp, err := svc.GetPlayerCommissionSummary(ctx, playerID, month)
	assert.NoError(t, err)
	assert.Equal(t, int64(120000), resp.MonthlyIncome)
	assert.Equal(t, int64(2500), resp.TotalCommission)
	assert.Equal(t, int64(10000), resp.TotalIncome)
	assert.Equal(t, int64(3), resp.TotalOrders)
}

func TestCommissionService_GetPlayerCommissionSummary_NoRecords(t *testing.T) {
	ctx := context.Background()
	commissionRepo := new(MockCommissionRepo)
	orderRepo := new(MockOrderRepo)
	playerRepo := new(MockPlayerRepo)

	svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

	playerID := uint64(7)
	month := "2025-02"

	commissionRepo.On("GetPlayerMonthlyIncome", ctx, playerID, month).Return(int64(0), nil)
	commissionRepo.
		On("ListRecords", ctx, mock.MatchedBy(func(opts commissionrepo.CommissionRecordListOptions) bool {
			return opts.PlayerID != nil && *opts.PlayerID == playerID && opts.Page == 1 && opts.PageSize == 1
		})).
		Return([]model.CommissionRecord{}, int64(0), nil).Once()

	resp, err := svc.GetPlayerCommissionSummary(ctx, playerID, month)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), resp.TotalCommission)
	assert.Equal(t, int64(0), resp.TotalIncome)
	assert.Equal(t, int64(0), resp.TotalOrders)
	commissionRepo.AssertNumberOfCalls(t, "ListRecords", 1)
}

func TestCommissionService_GetCommissionRecords(t *testing.T) {
	ctx := context.Background()
	commissionRepo := new(MockCommissionRepo)
	orderRepo := new(MockOrderRepo)
	playerRepo := new(MockPlayerRepo)
	svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

	playerID := uint64(99)
	now := time.Now()
	records := []model.CommissionRecord{
		{ID: 1, OrderID: 10, TotalAmountCents: 10000, CommissionRate: 20, CommissionCents: 2000, PlayerIncomeCents: 8000, SettlementStatus: "pending", SettlementMonth: "2025-01", CreatedAt: now},
		{ID: 2, OrderID: 11, TotalAmountCents: 20000, CommissionRate: 15, CommissionCents: 3000, PlayerIncomeCents: 17000, SettlementStatus: "done", SettlementMonth: "2025-01", CreatedAt: now.Add(time.Hour)},
	}

	commissionRepo.
		On("ListRecords", ctx, mock.MatchedBy(func(opts commissionrepo.CommissionRecordListOptions) bool {
			return opts.PlayerID != nil && *opts.PlayerID == playerID && opts.Page == 2 && opts.PageSize == 5
		})).
		Return(records, int64(len(records)), nil).Once()

	resp, err := svc.GetCommissionRecords(ctx, playerID, 2, 5)
	assert.NoError(t, err)
	assert.Equal(t, int64(len(records)), resp.Total)
	assert.Len(t, resp.Records, 2)
	assert.Equal(t, records[0].OrderID, resp.Records[0].OrderID)
	assert.Equal(t, records[1].SettlementStatus, resp.Records[1].SettlementStatus)
}

func TestCommissionService_GetMonthlySettlements(t *testing.T) {
	ctx := context.Background()
	commissionRepo := new(MockCommissionRepo)
	orderRepo := new(MockOrderRepo)
	playerRepo := new(MockPlayerRepo)
	svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

	playerID := uint64(5)
	now := time.Now()
	settlements := []model.MonthlySettlement{
		{
			ID:                   1,
			PlayerID:             playerID,
			SettlementMonth:      "2025-01",
			TotalOrderCount:      3,
			TotalAmountCents:     30000,
			TotalCommissionCents: 6000,
			TotalIncomeCents:     24000,
			BonusCents:           1000,
			FinalIncomeCents:     25000,
			Status:               "completed",
			CreatedAt:            now,
			SettledAt:            &now,
		},
	}

	commissionRepo.
		On("ListSettlements", ctx, mock.MatchedBy(func(opts commissionrepo.SettlementListOptions) bool {
			return opts.PlayerID != nil && *opts.PlayerID == playerID && opts.Page == 1 && opts.PageSize == 10
		})).
		Return(settlements, int64(len(settlements)), nil).Once()

	resp, err := svc.GetMonthlySettlements(ctx, playerID, 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(len(settlements)), resp.Total)
	assert.Len(t, resp.Settlements, 1)
	assert.Equal(t, settlements[0].SettlementMonth, resp.Settlements[0].SettlementMonth)
	assert.Equal(t, settlements[0].FinalIncomeCents, resp.Settlements[0].FinalIncomeCents)
}

func TestCommissionService_UpdateCommissionRule(t *testing.T) {
	ctx := context.Background()
	commissionRepo := new(MockCommissionRepo)
	orderRepo := new(MockOrderRepo)
	playerRepo := new(MockPlayerRepo)
	svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

	rule := &model.CommissionRule{
		ID:       1,
		Name:     "Base",
		Rate:     20,
		IsActive: true,
	}
	newName := "Premium"
	newRate := 35
	newActive := false

	commissionRepo.On("GetRule", ctx, uint64(1)).Return(rule, nil)
	commissionRepo.On("UpdateRule", ctx, rule).Return(nil)

	err := svc.UpdateCommissionRule(ctx, 1, UpdateCommissionRuleRequest{
		Name:     &newName,
		Rate:     &newRate,
		IsActive: &newActive,
	})

	assert.NoError(t, err)
	assert.Equal(t, newName, rule.Name)
	assert.Equal(t, newRate, rule.Rate)
	assert.False(t, rule.IsActive)
}

func TestCommissionService_UpdateCommissionRule_InvalidRate(t *testing.T) {
	ctx := context.Background()
	commissionRepo := new(MockCommissionRepo)
	orderRepo := new(MockOrderRepo)
	playerRepo := new(MockPlayerRepo)
	svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

	rule := &model.CommissionRule{ID: 2, Rate: 10}
	commissionRepo.On("GetRule", ctx, uint64(2)).Return(rule, nil)

	badRate := 150
	err := svc.UpdateCommissionRule(ctx, 2, UpdateCommissionRuleRequest{Rate: &badRate})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "commission rate")
	commissionRepo.AssertNotCalled(t, "UpdateRule", ctx, mock.Anything)
}

func TestCommissionService_GetPlatformStats(t *testing.T) {
	ctx := context.Background()
	commissionRepo := new(MockCommissionRepo)
	orderRepo := new(MockOrderRepo)
	playerRepo := new(MockPlayerRepo)
	svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

	stats := &commissionrepo.MonthlyStats{
		TotalOrders:       50,
		TotalIncome:       100000,
		TotalCommission:   20000,
		TotalPlayerIncome: 80000,
	}

	commissionRepo.On("GetMonthlyStats", ctx, "2025-01").Return(stats, nil)

	resp, err := svc.GetPlatformStats(ctx, "2025-01")
	assert.NoError(t, err)
	assert.Equal(t, stats.TotalOrders, resp.TotalOrders)
	assert.Equal(t, stats.TotalCommission, resp.TotalCommission)
	assert.Equal(t, "2025-01", resp.Month)
}

func TestSelectLowestRate(t *testing.T) {
	defaultCandidate := selectLowestRate(nil)
	assert.Equal(t, 20, defaultCandidate.Rate)
	assert.Contains(t, defaultCandidate.Detail, "20")

	candidates := []CommissionCandidate{
		{Source: "A", Rate: 30},
		{Source: "B", Rate: 15},
		{Source: "C", Rate: 20},
	}
	chosen := selectLowestRate(candidates)
	assert.Equal(t, "B", chosen.Source)
	assert.Equal(t, 15, chosen.Rate)
}

func TestParseRankingCommissionRules_InvalidJSON(t *testing.T) {
	jsonStr := `[{"rankStart":1,"rankEnd":10,"commissionRate":15}]`
	rules, err := ParseRankingCommissionRules(jsonStr)
	assert.NoError(t, err)
	assert.Len(t, rules, 1)
	assert.Equal(t, 15, rules[0].CommissionRate)

	_, err = ParseRankingCommissionRules("invalid json")
	assert.Error(t, err)
}

func TestFindCommissionRateForRank_Bounds(t *testing.T) {
	rules := []model.RankingCommissionRule{
		{RankStart: 1, RankEnd: 5, CommissionRate: 20},
		{RankStart: 6, RankEnd: 10, CommissionRate: 15},
	}

	assert.Equal(t, 20, FindCommissionRateForRank(rules, 4))
	assert.Equal(t, 15, FindCommissionRateForRank(rules, 9))
	assert.Equal(t, 0, FindCommissionRateForRank(rules, 11))
}

func TestValidateRankingRules_InvalidInputs(t *testing.T) {
	valid := []model.RankingCommissionRule{
		{RankStart: 1, RankEnd: 5, CommissionRate: 20},
		{RankStart: 6, RankEnd: 10, CommissionRate: 15},
	}
	assert.NoError(t, ValidateRankingRules(valid))

	invalidRate := []model.RankingCommissionRule{
		{RankStart: 1, RankEnd: 5, CommissionRate: 120},
	}
	assert.Error(t, ValidateRankingRules(invalidRate))

	overlap := []model.RankingCommissionRule{
		{RankStart: 1, RankEnd: 5, CommissionRate: 20},
		{RankStart: 5, RankEnd: 10, CommissionRate: 15},
	}
	assert.Error(t, ValidateRankingRules(overlap))
}

func TestRangesOverlap(t *testing.T) {
	assert.True(t, rangesOverlap(1, 5, 3, 7))
	assert.True(t, rangesOverlap(1, 5, 5, 10))
	assert.False(t, rangesOverlap(1, 4, 5, 8))
}
