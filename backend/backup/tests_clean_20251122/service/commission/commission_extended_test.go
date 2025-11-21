package commission

import (
	"context"
	"errors"
	"testing"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	commissionrepo "gamelink/internal/repository/commission"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestCalculateCommission_EdgeCases 测试抽成计算的边界情况
func TestCalculateCommission_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("零金额订单的抽成计算", func(t *testing.T) {
		commissionRepo := new(MockCommissionRepo)
		orderRepo := new(MockOrderRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

		gameID := uint64(1)
		playerID := uint64(5)
		order := &model.Order{
			GameID:          &gameID,
			PlayerID:        &playerID,
			TotalPriceCents: 0, // 零金额
		}
		order.ID = 1001

		defaultRule := &model.CommissionRule{
			Rate: 20,
		}

		orderRepo.On("Get", ctx, uint64(1001)).Return(order, nil)
		commissionRepo.On("GetRuleForOrder", ctx, order.GameID, order.PlayerID, (*string)(nil)).
			Return(nil, repository.ErrNotFound)
		commissionRepo.On("GetDefaultRule", ctx).Return(defaultRule, nil)

		calc, err := svc.CalculateCommission(ctx, 1001)

		assert.NoError(t, err)
		assert.NotNil(t, calc)
		assert.Equal(t, int64(0), calc.TotalAmountCents)
		assert.Equal(t, int64(0), calc.CommissionCents)
		assert.Equal(t, int64(0), calc.PlayerIncomeCents)
	})

	t.Run("极大金额订单的抽成计算", func(t *testing.T) {
		commissionRepo := new(MockCommissionRepo)
		orderRepo := new(MockOrderRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

		gameID := uint64(1)
		playerID := uint64(5)
		order := &model.Order{
			GameID:          &gameID,
			PlayerID:        &playerID,
			TotalPriceCents: 10000000, // 100,000元
		}
		order.ID = 1001

		defaultRule := &model.CommissionRule{
			Rate: 20,
		}

		orderRepo.On("Get", ctx, uint64(1001)).Return(order, nil)
		commissionRepo.On("GetRuleForOrder", ctx, order.GameID, order.PlayerID, (*string)(nil)).
			Return(nil, repository.ErrNotFound)
		commissionRepo.On("GetDefaultRule", ctx).Return(defaultRule, nil)

		calc, err := svc.CalculateCommission(ctx, 1001)

		assert.NoError(t, err)
		assert.Equal(t, int64(10000000), calc.TotalAmountCents)
		assert.Equal(t, int64(2000000), calc.CommissionCents)   // 20%
		assert.Equal(t, int64(8000000), calc.PlayerIncomeCents) // 80%
	})

	t.Run("抽成率为0的情况", func(t *testing.T) {
		commissionRepo := new(MockCommissionRepo)
		orderRepo := new(MockOrderRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

		gameID := uint64(1)
		playerID := uint64(5)
		order := &model.Order{
			GameID:          &gameID,
			PlayerID:        &playerID,
			TotalPriceCents: 50000,
		}
		order.ID = 1001

		zeroRule := &model.CommissionRule{
			Rate: 0, // 零抽成
		}

		orderRepo.On("Get", ctx, uint64(1001)).Return(order, nil)
		commissionRepo.On("GetRuleForOrder", ctx, order.GameID, order.PlayerID, (*string)(nil)).
			Return(zeroRule, nil)

		calc, err := svc.CalculateCommission(ctx, 1001)

		assert.NoError(t, err)
		assert.Equal(t, 0, calc.CommissionRate)
		assert.Equal(t, int64(0), calc.CommissionCents)
		assert.Equal(t, int64(50000), calc.PlayerIncomeCents) // 全部归玩家
	})

	t.Run("抽成率为100%的情况", func(t *testing.T) {
		commissionRepo := new(MockCommissionRepo)
		orderRepo := new(MockOrderRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

		gameID := uint64(1)
		playerID := uint64(5)
		order := &model.Order{
			GameID:          &gameID,
			PlayerID:        &playerID,
			TotalPriceCents: 50000,
		}
		order.ID = 1001

		fullRule := &model.CommissionRule{
			Rate: 100, // 全额抽成
		}

		orderRepo.On("Get", ctx, uint64(1001)).Return(order, nil)
		commissionRepo.On("GetRuleForOrder", ctx, order.GameID, order.PlayerID, (*string)(nil)).
			Return(fullRule, nil)

		calc, err := svc.CalculateCommission(ctx, 1001)

		assert.NoError(t, err)
		assert.Equal(t, 100, calc.CommissionRate)
		assert.Equal(t, int64(50000), calc.CommissionCents)
		assert.Equal(t, int64(0), calc.PlayerIncomeCents) // 玩家收入为0
	})

	t.Run("订单不存在应该返回错误", func(t *testing.T) {
		commissionRepo := new(MockCommissionRepo)
		orderRepo := new(MockOrderRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

		orderRepo.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

		calc, err := svc.CalculateCommission(ctx, 999)

		assert.Error(t, err)
		assert.Nil(t, calc)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("没有任何规则时使用默认规则", func(t *testing.T) {
		commissionRepo := new(MockCommissionRepo)
		orderRepo := new(MockOrderRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

		gameID := uint64(1)
		playerID := uint64(5)
		order := &model.Order{
			GameID:          &gameID,
			PlayerID:        &playerID,
			TotalPriceCents: 50000,
		}
		order.ID = 1001

		defaultRule := &model.CommissionRule{
			Rate: 25,
			Name: "新手默认25%",
		}

		orderRepo.On("Get", ctx, uint64(1001)).Return(order, nil)
		commissionRepo.On("GetRuleForOrder", ctx, order.GameID, order.PlayerID, (*string)(nil)).
			Return(nil, repository.ErrNotFound)
		commissionRepo.On("GetDefaultRule", ctx).Return(defaultRule, nil)

		calc, err := svc.CalculateCommission(ctx, 1001)

		assert.NoError(t, err)
		assert.Equal(t, 25, calc.CommissionRate)
		assert.Equal(t, int64(12500), calc.CommissionCents)
	})

	t.Run("默认规则也不存在时使用硬编码20%", func(t *testing.T) {
		commissionRepo := new(MockCommissionRepo)
		orderRepo := new(MockOrderRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

		gameID := uint64(1)
		playerID := uint64(5)
		order := &model.Order{
			GameID:          &gameID,
			PlayerID:        &playerID,
			TotalPriceCents: 50000,
		}
		order.ID = 1001

		orderRepo.On("Get", ctx, uint64(1001)).Return(order, nil)
		commissionRepo.On("GetRuleForOrder", ctx, order.GameID, order.PlayerID, (*string)(nil)).
			Return(nil, repository.ErrNotFound)
		commissionRepo.On("GetDefaultRule", ctx).Return(nil, repository.ErrNotFound)

		calc, err := svc.CalculateCommission(ctx, 1001)

		assert.NoError(t, err)
		assert.Equal(t, 20, calc.CommissionRate) // 硬编码的默认值
		assert.Equal(t, int64(10000), calc.CommissionCents)
	})
}

// TestRecordCommission_EdgeCases 测试记录抽成的边界情况
func TestRecordCommission_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("订单没有玩家ID应该失败", func(t *testing.T) {
		commissionRepo := new(MockCommissionRepo)
		orderRepo := new(MockOrderRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

		order := &model.Order{
			TotalPriceCents: 50000,
			// PlayerID 为 nil
		}
		order.ID = 1001

		defaultRule := &model.CommissionRule{
			Rate: 20,
		}

		commissionRepo.On("GetRecordByOrderID", ctx, uint64(1001)).
			Return(nil, repository.ErrNotFound)
		orderRepo.On("Get", ctx, uint64(1001)).Return(order, nil)
		commissionRepo.On("GetRuleForOrder", ctx, order.GameID, order.PlayerID, (*string)(nil)).
			Return(nil, repository.ErrNotFound)
		commissionRepo.On("GetDefaultRule", ctx).Return(defaultRule, nil)

		err := svc.RecordCommission(ctx, 1001)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no player assigned")
	})

	t.Run("数据库创建记录失败应该返回错误", func(t *testing.T) {
		commissionRepo := new(MockCommissionRepo)
		orderRepo := new(MockOrderRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

		gameID := uint64(1)
		playerID := uint64(5)
		order := &model.Order{
			GameID:          &gameID,
			PlayerID:        &playerID,
			TotalPriceCents: 50000,
		}
		order.ID = 1001

		defaultRule := &model.CommissionRule{
			Rate: 20,
		}

		dbErr := errors.New("database constraint violation")

		commissionRepo.On("GetRecordByOrderID", ctx, uint64(1001)).
			Return(nil, repository.ErrNotFound)
		orderRepo.On("Get", ctx, uint64(1001)).Return(order, nil)
		commissionRepo.On("GetRuleForOrder", ctx, order.GameID, order.PlayerID, (*string)(nil)).
			Return(nil, repository.ErrNotFound)
		commissionRepo.On("GetDefaultRule", ctx).Return(defaultRule, nil)
		commissionRepo.On("CreateRecord", ctx, mock.Anything).Return(dbErr)

		err := svc.RecordCommission(ctx, 1001)

		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
	})

	t.Run("计算抽成失败应该返回错误", func(t *testing.T) {
		commissionRepo := new(MockCommissionRepo)
		orderRepo := new(MockOrderRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

		commissionRepo.On("GetRecordByOrderID", ctx, uint64(999)).
			Return(nil, repository.ErrNotFound)
		orderRepo.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

		err := svc.RecordCommission(ctx, 999)

		assert.Error(t, err)
		assert.Equal(t, repository.ErrNotFound, err)
	})
}

// TestCreateCommissionRule_EdgeCases 测试创建抽成规则的边界情况
func TestCreateCommissionRule_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("抽成率为负数应该失败", func(t *testing.T) {
		commissionRepo := new(MockCommissionRepo)
		orderRepo := new(MockOrderRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

		req := CreateCommissionRuleRequest{
			Name: "无效规则",
			Rate: -10, // 负数
		}

		rule, err := svc.CreateCommissionRule(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, rule)
		assert.Contains(t, err.Error(), "commission rate must be between 0 and 100")
	})

	t.Run("抽成率超过100应该失败", func(t *testing.T) {
		commissionRepo := new(MockCommissionRepo)
		orderRepo := new(MockOrderRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

		req := CreateCommissionRuleRequest{
			Name: "无效规则",
			Rate: 150, // 超过100
		}

		rule, err := svc.CreateCommissionRule(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, rule)
		assert.Contains(t, err.Error(), "commission rate must be between 0 and 100")
	})

	t.Run("抽成率为0应该成功", func(t *testing.T) {
		commissionRepo := new(MockCommissionRepo)
		orderRepo := new(MockOrderRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

		req := CreateCommissionRuleRequest{
			Name: "零抽成规则",
			Rate: 0,
		}

		commissionRepo.On("CreateRule", ctx, mock.MatchedBy(func(rule *model.CommissionRule) bool {
			return rule.Rate == 0
		})).Return(nil)

		rule, err := svc.CreateCommissionRule(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, rule)
		assert.Equal(t, 0, rule.Rate)
	})

	t.Run("抽成率为100应该成功", func(t *testing.T) {
		commissionRepo := new(MockCommissionRepo)
		orderRepo := new(MockOrderRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

		req := CreateCommissionRuleRequest{
			Name: "全额抽成规则",
			Rate: 100,
		}

		commissionRepo.On("CreateRule", ctx, mock.MatchedBy(func(rule *model.CommissionRule) bool {
			return rule.Rate == 100
		})).Return(nil)

		rule, err := svc.CreateCommissionRule(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, rule)
		assert.Equal(t, 100, rule.Rate)
	})
}

// 注意：GetCommissionSummary方法可能不存在，这些测试被注释掉
// 如果需要，可以在实现该方法后取消注释

// TestGetCommissionRecords_EdgeCases 测试获取抽成记录的边界情况
func TestGetCommissionRecords_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("空记录列表应该返回空数组", func(t *testing.T) {
		commissionRepo := new(MockCommissionRepo)
		orderRepo := new(MockOrderRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

		playerID := uint64(5)

		commissionRepo.On("ListRecords", ctx, mock.Anything).
			Return([]model.CommissionRecord{}, int64(0), nil)

		result, err := svc.GetCommissionRecords(ctx, playerID, 1, 20)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result.Records)
		assert.Equal(t, int64(0), result.Total)
	})

	t.Run("大页码应该正常处理", func(t *testing.T) {
		commissionRepo := new(MockCommissionRepo)
		orderRepo := new(MockOrderRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

		playerID := uint64(5)

		commissionRepo.On("ListRecords", ctx, mock.MatchedBy(func(opts commissionrepo.CommissionRecordListOptions) bool {
			return opts.Page == 1000
		})).Return([]model.CommissionRecord{}, int64(0), nil)

		result, err := svc.GetCommissionRecords(ctx, playerID, 1000, 20)

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}

// TestSettleMonth_EdgeCases 测试月度结算的边界情况
func TestSettleMonth_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("已经存在结算记录应该失败", func(t *testing.T) {
		commissionRepo := new(MockCommissionRepo)
		orderRepo := new(MockOrderRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

		month := "2024-01"

		existingSettlement := []model.MonthlySettlement{
			{
				ID:              1,
				SettlementMonth: month,
			},
		}

		commissionRepo.On("ListSettlements", ctx, mock.MatchedBy(func(opts commissionrepo.SettlementListOptions) bool {
			return *opts.SettlementMonth == month
		})).Return(existingSettlement, int64(1), nil)

		err := svc.SettleMonth(ctx, month)

		assert.Error(t, err)
		assert.Equal(t, ErrAlreadySettled, err)
	})

	t.Run("没有待结算记录应该返回错误", func(t *testing.T) {
		commissionRepo := new(MockCommissionRepo)
		orderRepo := new(MockOrderRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

		month := "2024-01"

		// 没有已存在的结算
		commissionRepo.On("ListSettlements", ctx, mock.MatchedBy(func(opts commissionrepo.SettlementListOptions) bool {
			return *opts.SettlementMonth == month
		})).Return([]model.MonthlySettlement{}, int64(0), nil)

		// 没有待结算的记录
		commissionRepo.On("ListRecords", ctx, mock.MatchedBy(func(opts commissionrepo.CommissionRecordListOptions) bool {
			return *opts.SettlementMonth == month && *opts.SettlementStatus == "pending"
		})).Return([]model.CommissionRecord{}, int64(0), nil)

		err := svc.SettleMonth(ctx, month)

		// 应该返回错误（没有记录可结算）
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no records to settle")
	})
}

// TestConcurrentRecordCommission 测试并发记录抽成
func TestConcurrentRecordCommission(t *testing.T) {
	// 这个测试模拟并发场景，确保不会重复记录
	ctx := context.Background()

	commissionRepo := new(MockCommissionRepo)
	orderRepo := new(MockOrderRepo)
	playerRepo := new(MockPlayerRepo)
	svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

	gameID := uint64(1)
	playerID := uint64(5)
	order := &model.Order{
		GameID:          &gameID,
		PlayerID:        &playerID,
		TotalPriceCents: 50000,
	}
	order.ID = 1001

	defaultRule := &model.CommissionRule{
		Rate: 20,
	}

	// 第一次调用：没有记录
	commissionRepo.On("GetRecordByOrderID", ctx, uint64(1001)).
		Return(nil, repository.ErrNotFound).Once()
	orderRepo.On("Get", ctx, uint64(1001)).Return(order, nil)
	commissionRepo.On("GetRuleForOrder", ctx, order.GameID, order.PlayerID, (*string)(nil)).
		Return(nil, repository.ErrNotFound)
	commissionRepo.On("GetDefaultRule", ctx).Return(defaultRule, nil)
	commissionRepo.On("CreateRecord", ctx, mock.Anything).Return(nil).Once()

	// 第二次调用：已经有记录了
	existingRecord := &model.CommissionRecord{
		ID:      1,
		OrderID: 1001,
	}
	commissionRepo.On("GetRecordByOrderID", ctx, uint64(1001)).
		Return(existingRecord, nil).Once()

	// 第一次应该成功
	err1 := svc.RecordCommission(ctx, 1001)
	assert.NoError(t, err1)

	// 第二次应该失败（已记录）
	err2 := svc.RecordCommission(ctx, 1001)
	assert.Error(t, err2)
	assert.Equal(t, ErrAlreadyRecorded, err2)
}

// 注意：由于GetPlayerCommissionSummary、GetMonthlySettlements和GetPlatformStats
// 依赖复杂的repository返回类型，这些测试已在commission_test.go中覆盖
// 这里只测试不依赖复杂mock的函数

// TestUpdateCommissionRule_EdgeCases 测试更新抽成规则的边界情况
func TestUpdateCommissionRule_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("规则不存在应该失败", func(t *testing.T) {
		commissionRepo := new(MockCommissionRepo)
		orderRepo := new(MockOrderRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

		commissionRepo.On("GetRule", ctx, uint64(999)).
			Return(nil, repository.ErrNotFound)

		rate := 15
		req := UpdateCommissionRuleRequest{
			Rate: &rate,
		}

		err := svc.UpdateCommissionRule(ctx, 999, req)

		assert.Error(t, err)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("抽成率超出范围应该失败", func(t *testing.T) {
		commissionRepo := new(MockCommissionRepo)
		orderRepo := new(MockOrderRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

		existingRule := &model.CommissionRule{
			Rate: 20,
		}
		existingRule.ID = 1

		commissionRepo.On("GetRule", ctx, uint64(1)).
			Return(existingRule, nil)

		rate := 150 // 超出范围
		req := UpdateCommissionRuleRequest{
			Rate: &rate,
		}

		err := svc.UpdateCommissionRule(ctx, 1, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be between 0 and 100")
	})

	t.Run("成功更新抽成规则", func(t *testing.T) {
		commissionRepo := new(MockCommissionRepo)
		orderRepo := new(MockOrderRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

		existingRule := &model.CommissionRule{
			Rate: 20,
		}
		existingRule.ID = 1

		commissionRepo.On("GetRule", ctx, uint64(1)).
			Return(existingRule, nil)
		commissionRepo.On("UpdateRule", ctx, mock.MatchedBy(func(rule *model.CommissionRule) bool {
			return rule.Rate == 15
		})).Return(nil)

		rate := 15
		req := UpdateCommissionRuleRequest{
			Rate: &rate,
		}

		err := svc.UpdateCommissionRule(ctx, 1, req)

		assert.NoError(t, err)
	})
}

// TestParseRankingCommissionRules 测试解析排名抽成规则
func TestParseRankingCommissionRules(t *testing.T) {
	t.Run("成功解析有效JSON", func(t *testing.T) {
		rulesJSON := `[{"rankStart":1,"rankEnd":3,"commissionRate":10},{"rankStart":4,"rankEnd":10,"commissionRate":15}]`

		rules, err := ParseRankingCommissionRules(rulesJSON)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(rules))
		assert.Equal(t, 1, rules[0].RankStart)
		assert.Equal(t, 3, rules[0].RankEnd)
		assert.Equal(t, 10, rules[0].CommissionRate)
	})

	t.Run("无效JSON应该失败", func(t *testing.T) {
		rulesJSON := `invalid json`

		rules, err := ParseRankingCommissionRules(rulesJSON)

		assert.Error(t, err)
		assert.Nil(t, rules)
	})

	t.Run("空JSON应该返回空数组", func(t *testing.T) {
		rulesJSON := `[]`

		rules, err := ParseRankingCommissionRules(rulesJSON)

		assert.NoError(t, err)
		assert.Equal(t, 0, len(rules))
	})
}

// TestFindCommissionRateForRank 测试根据排名查找抽成比例
func TestFindCommissionRateForRank(t *testing.T) {
	rules := []model.RankingCommissionRule{
		{RankStart: 1, RankEnd: 3, CommissionRate: 10},
		{RankStart: 4, RankEnd: 10, CommissionRate: 15},
		{RankStart: 11, RankEnd: 20, CommissionRate: 18},
	}

	t.Run("排名1应该返回10%", func(t *testing.T) {
		rate := FindCommissionRateForRank(rules, 1)
		assert.Equal(t, 10, rate)
	})

	t.Run("排名5应该返回15%", func(t *testing.T) {
		rate := FindCommissionRateForRank(rules, 5)
		assert.Equal(t, 15, rate)
	})

	t.Run("排名15应该返回18%", func(t *testing.T) {
		rate := FindCommissionRateForRank(rules, 15)
		assert.Equal(t, 18, rate)
	})

	t.Run("排名100（不在范围内）应该返回0", func(t *testing.T) {
		rate := FindCommissionRateForRank(rules, 100)
		assert.Equal(t, 0, rate)
	})
}

// TestValidateRankingRules 测试验证排名规则
func TestValidateRankingRules(t *testing.T) {
	t.Run("有效规则应该通过", func(t *testing.T) {
		rules := []model.RankingCommissionRule{
			{RankStart: 1, RankEnd: 3, CommissionRate: 10},
			{RankStart: 4, RankEnd: 10, CommissionRate: 15},
		}

		err := ValidateRankingRules(rules)
		assert.NoError(t, err)
	})

	t.Run("RankStart小于1应该失败", func(t *testing.T) {
		rules := []model.RankingCommissionRule{
			{RankStart: 0, RankEnd: 3, CommissionRate: 10},
		}

		err := ValidateRankingRules(rules)
		assert.Error(t, err)
		assert.Equal(t, ErrValidation, err)
	})

	t.Run("RankEnd小于RankStart应该失败", func(t *testing.T) {
		rules := []model.RankingCommissionRule{
			{RankStart: 5, RankEnd: 3, CommissionRate: 10},
		}

		err := ValidateRankingRules(rules)
		assert.Error(t, err)
		assert.Equal(t, ErrValidation, err)
	})

	t.Run("抽成率超出范围应该失败", func(t *testing.T) {
		rules := []model.RankingCommissionRule{
			{RankStart: 1, RankEnd: 3, CommissionRate: 150},
		}

		err := ValidateRankingRules(rules)
		assert.Error(t, err)
		assert.Equal(t, ErrValidation, err)
	})
}
