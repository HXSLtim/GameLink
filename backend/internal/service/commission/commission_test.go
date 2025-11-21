package commission

import (
	"time"
	"testing"
	commissionrepo "gamelink/internal/repository/commission"
	"context"
	"errors"
	"gamelink/internal/model"
	"gamelink/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCommissionRepo struct {
	mock.Mock
}

func (m *MockCommissionRepo) CreateRule(ctx context.Context, rule *model.CommissionRule) error {
	args := m.Called(ctx, rule)
	if args.Get(0) != nil {
		rule.ID = 1
	}
	return args.Error(0)
}

func (m *MockCommissionRepo) GetRule(ctx context.Context, id uint64) (*model.CommissionRule, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CommissionRule), args.Error(1)
}

func (m *MockCommissionRepo) GetDefaultRule(ctx context.Context) (*model.CommissionRule, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CommissionRule), args.Error(1)
}

func (m *MockCommissionRepo) GetRuleForOrder(ctx context.Context, gameID, playerID *uint64, serviceType *string) (*model.CommissionRule, error) {
	args := m.Called(ctx, gameID, playerID, serviceType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CommissionRule), args.Error(1)
}

func (m *MockCommissionRepo) ListRules(ctx context.Context, opts commissionrepo.CommissionRuleListOptions) ([]model.CommissionRule, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.CommissionRule), args.Get(1).(int64), args.Error(2)
}

func (m *MockCommissionRepo) UpdateRule(ctx context.Context, rule *model.CommissionRule) error {
	args := m.Called(ctx, rule)
	return args.Error(0)
}

func (m *MockCommissionRepo) DeleteRule(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCommissionRepo) CreateRecord(ctx context.Context, record *model.CommissionRecord) error {
	args := m.Called(ctx, record)
	if args.Get(0) != nil {
		record.ID = 1
	}
	return args.Error(0)
}

func (m *MockCommissionRepo) GetRecord(ctx context.Context, id uint64) (*model.CommissionRecord, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CommissionRecord), args.Error(1)
}

func (m *MockCommissionRepo) GetRecordByOrderID(ctx context.Context, orderID uint64) (*model.CommissionRecord, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CommissionRecord), args.Error(1)
}

func (m *MockCommissionRepo) ListRecords(ctx context.Context, opts commissionrepo.CommissionRecordListOptions) ([]model.CommissionRecord, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.CommissionRecord), args.Get(1).(int64), args.Error(2)
}

func (m *MockCommissionRepo) UpdateRecord(ctx context.Context, record *model.CommissionRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockCommissionRepo) CreateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error {
	args := m.Called(ctx, settlement)
	return args.Error(0)
}

func (m *MockCommissionRepo) GetSettlement(ctx context.Context, id uint64) (*model.MonthlySettlement, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.MonthlySettlement), args.Error(1)
}

func (m *MockCommissionRepo) GetSettlementByPlayerMonth(ctx context.Context, playerID uint64, month string) (*model.MonthlySettlement, error) {
	args := m.Called(ctx, playerID, month)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.MonthlySettlement), args.Error(1)
}

func (m *MockCommissionRepo) ListSettlements(ctx context.Context, opts commissionrepo.SettlementListOptions) ([]model.MonthlySettlement, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.MonthlySettlement), args.Get(1).(int64), args.Error(2)
}

func (m *MockCommissionRepo) UpdateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error {
	args := m.Called(ctx, settlement)
	return args.Error(0)
}

func (m *MockCommissionRepo) GetMonthlyStats(ctx context.Context, month string) (*commissionrepo.MonthlyStats, error) {
	args := m.Called(ctx, month)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*commissionrepo.MonthlyStats), args.Error(1)
}

func (m *MockCommissionRepo) GetPlayerMonthlyIncome(ctx context.Context, playerID uint64, month string) (int64, error) {
	args := m.Called(ctx, playerID, month)
	return args.Get(0).(int64), args.Error(1)
}

type MockOrderRepo struct {
	mock.Mock
}

func (m *MockOrderRepo) Get(ctx context.Context, id uint64) (*model.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderRepo) Create(ctx context.Context, order *model.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepo) Update(ctx context.Context, order *model.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepo) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOrderRepo) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.Order), args.Get(1).(int64), args.Error(2)
}

type MockPlayerRepo struct {
	mock.Mock
}

func (m *MockPlayerRepo) Get(ctx context.Context, id uint64) (*model.Player, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Player), args.Error(1)
}

func (m *MockPlayerRepo) Create(ctx context.Context, player *model.Player) error { return nil }
func (m *MockPlayerRepo) Update(ctx context.Context, player *model.Player) error { return nil }

func (m *MockPlayerRepo) Delete(ctx context.Context, id uint64) error            { return nil }
func (m *MockPlayerRepo) List(ctx context.Context) ([]model.Player, error)       { return nil, nil }

func (m *MockPlayerRepo) ListPaged(ctx context.Context, page, pageSize int) ([]model.Player, int64, error) {
	return nil, 0, nil
}

func (m *MockPlayerRepo) GetByUserID(ctx context.Context, userID uint64) (*model.Player, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Player), args.Error(1)
}

func TestCommissionService_CalculateCommission(t *testing.T) {
	ctx := context.Background()

	commissionRepo := new(MockCommissionRepo)
	orderRepo := new(MockOrderRepo)
	playerRepo := new(MockPlayerRepo)

	svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

	t.Run("使用默认抽成规则", func(t *testing.T) {
		gameID := uint64(1)
		playerID := uint64(5)
		order := &model.Order{
			GameID:          &gameID,
			PlayerID:        &playerID,
			TotalPriceCents: 50000,
		}
		order.ID = 1001

		defaultRule := &model.CommissionRule{
			Rate: 20, // 20%
		}

		orderRepo.On("Get", ctx, uint64(1001)).Return(order, nil)
		commissionRepo.On("GetRuleForOrder", ctx, order.GameID, order.PlayerID, (*string)(nil)).
			Return(nil, repository.ErrNotFound)
		commissionRepo.On("GetDefaultRule", ctx).Return(defaultRule, nil)

		// 计算抽成
		calc, err := svc.CalculateCommission(ctx, 1001)

		// 验证
		assert.NoError(t, err)
		assert.NotNil(t, calc)
		assert.Equal(t, uint64(1001), calc.OrderID)
		assert.Equal(t, int64(50000), calc.TotalAmountCents)
		assert.Equal(t, 20, calc.CommissionRate)
		assert.Equal(t, int64(10000), calc.CommissionCents)   // 20%
		assert.Equal(t, int64(40000), calc.PlayerIncomeCents) // 80%
	})

	t.Run("使用特殊抽成规则", func(t *testing.T) {
		commissionRepo := new(MockCommissionRepo)
		orderRepo := new(MockOrderRepo)
		playerRepo := new(MockPlayerRepo)

		svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

		gameID := uint64(1)
		playerID := uint64(5)
		order := &model.Order{
			GameID:          &gameID,
			PlayerID:        &playerID,
			TotalPriceCents: 100000,
		}
		order.ID = 1002

		// 特殊规则：15%抽成
		specialRule := &model.CommissionRule{
			Rate:   15,
			GameID: &gameID,
		}

		orderRepo.On("Get", ctx, uint64(1002)).Return(order, nil)
		commissionRepo.On("GetRuleForOrder", ctx, order.GameID, order.PlayerID, (*string)(nil)).
			Return(specialRule, nil)

		// 计算抽成
		calc, err := svc.CalculateCommission(ctx, 1002)

		// 验证使用了15%的特殊规则
		assert.NoError(t, err)
		assert.Equal(t, 15, calc.CommissionRate)
		assert.Equal(t, int64(15000), calc.CommissionCents)   // 15%
		assert.Equal(t, int64(85000), calc.PlayerIncomeCents) // 85%
	})
}

func TestCommissionService_RecordCommission(t *testing.T) {
	ctx := context.Background()

	commissionRepo := new(MockCommissionRepo)
	orderRepo := new(MockOrderRepo)
	playerRepo := new(MockPlayerRepo)

	svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

	t.Run("成功记录抽成", func(t *testing.T) {
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

		commissionRepo.On("GetRecordByOrderID", ctx, uint64(1001)).
			Return(nil, repository.ErrNotFound)
		orderRepo.On("Get", ctx, uint64(1001)).Return(order, nil)
		commissionRepo.On("GetRuleForOrder", ctx, order.GameID, order.PlayerID, (*string)(nil)).
			Return(nil, repository.ErrNotFound)
		commissionRepo.On("GetDefaultRule", ctx).Return(defaultRule, nil)
		commissionRepo.On("CreateRecord", ctx, mock.MatchedBy(func(record *model.CommissionRecord) bool {
			assert.Equal(t, uint64(1001), record.OrderID)
			assert.Equal(t, uint64(5), record.PlayerID)
			assert.Equal(t, int64(50000), record.TotalAmountCents)
			assert.Equal(t, 20, record.CommissionRate)
			assert.Equal(t, int64(10000), record.CommissionCents)
			assert.Equal(t, int64(40000), record.PlayerIncomeCents)
			assert.Equal(t, "pending", record.SettlementStatus)
			assert.NotEmpty(t, record.SettlementMonth)
			return true
		})).Return(nil)

		// 记录抽成
		err := svc.RecordCommission(ctx, 1001)

		// 验证
		assert.NoError(t, err)
		commissionRepo.AssertExpectations(t)
	})

	t.Run("已经记录过抽成", func(t *testing.T) {
		commissionRepo := new(MockCommissionRepo)
		orderRepo := new(MockOrderRepo)
		playerRepo := new(MockPlayerRepo)

		svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

		existingRecord := &model.CommissionRecord{
			ID:      1,
			OrderID: 1001,
		}

		commissionRepo.On("GetRecordByOrderID", ctx, uint64(1001)).
			Return(existingRecord, nil)

		// 尝试重复记录
		err := svc.RecordCommission(ctx, 1001)

		// 应该返回错误
		assert.Error(t, err)
		assert.Equal(t, ErrAlreadyRecorded, err)
	})
}

func TestCommissionService_SettleMonth(t *testing.T) {
	ctx := context.Background()

	commissionRepo := new(MockCommissionRepo)
	orderRepo := new(MockOrderRepo)
	playerRepo := new(MockPlayerRepo)

	svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

	t.Run("成功执行月度结算", func(t *testing.T) {
		month := "2024-11"
		status := "pending"

		// Mock: 该月没有结算记录
		commissionRepo.On("ListSettlements", ctx, mock.MatchedBy(func(opts commissionrepo.SettlementListOptions) bool {
			return *opts.SettlementMonth == month
		})).Return([]model.MonthlySettlement{}, int64(0), nil)

		// Mock: 有3条待结算记录（2个陪玩师）
		records := []model.CommissionRecord{
			{ID: 1, OrderID: 101, PlayerID: 5, TotalAmountCents: 50000, CommissionCents: 10000, PlayerIncomeCents: 40000},
			{ID: 2, OrderID: 102, PlayerID: 5, TotalAmountCents: 30000, CommissionCents: 6000, PlayerIncomeCents: 24000},
			{ID: 3, OrderID: 103, PlayerID: 6, TotalAmountCents: 40000, CommissionCents: 8000, PlayerIncomeCents: 32000},
		}

		commissionRepo.On("ListRecords", ctx, mock.MatchedBy(func(opts commissionrepo.CommissionRecordListOptions) bool {
			return *opts.SettlementMonth == month && *opts.SettlementStatus == status
		})).Return(records, int64(3), nil)

		// Mock: 创建结算记录
		commissionRepo.On("CreateSettlement", ctx, mock.MatchedBy(func(settlement *model.MonthlySettlement) bool {
			if settlement.PlayerID == 5 {
				// 陪玩师5: 2笔订单
				assert.Equal(t, int64(2), settlement.TotalOrderCount)
				assert.Equal(t, int64(80000), settlement.TotalAmountCents)     // 50000+30000
				assert.Equal(t, int64(16000), settlement.TotalCommissionCents) // 10000+6000
				assert.Equal(t, int64(64000), settlement.TotalIncomeCents)     // 40000+24000
			} else if settlement.PlayerID == 6 {
				// 陪玩师6: 1笔订单
				assert.Equal(t, int64(1), settlement.TotalOrderCount)
				assert.Equal(t, int64(40000), settlement.TotalAmountCents)
			}
			return true
		})).Return(nil)

		// Mock: 更新记录状态
		commissionRepo.On("UpdateRecord", ctx, mock.AnythingOfType("*model.CommissionRecord")).
			Return(nil)

		// 执行结算
		err := svc.SettleMonth(ctx, month)

		// 验证
		assert.NoError(t, err)

		// 验证创建了2个结算记录（2个陪玩师）
		commissionRepo.AssertNumberOfCalls(t, "CreateSettlement", 2)

		// 验证更新了3条记录
		commissionRepo.AssertNumberOfCalls(t, "UpdateRecord", 3)
	})

	t.Run("月份已经结算过", func(t *testing.T) {
		commissionRepo := new(MockCommissionRepo)
		orderRepo := new(MockOrderRepo)
		playerRepo := new(MockPlayerRepo)

		svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

		month := "2024-11"

		// Mock: 已有结算记录
		existingSettlement := []model.MonthlySettlement{
			{ID: 1, PlayerID: 5, SettlementMonth: month},
		}

		commissionRepo.On("ListSettlements", ctx, mock.MatchedBy(func(opts commissionrepo.SettlementListOptions) bool {
			return *opts.SettlementMonth == month
		})).Return(existingSettlement, int64(1), nil)

		// 尝试重复结算
		err := svc.SettleMonth(ctx, month)

		// 应该返回错误
		assert.Error(t, err)
		assert.Equal(t, ErrAlreadySettled, err)
	})
}

func TestCommissionService_CreateCommissionRule(t *testing.T) {
	ctx := context.Background()

	commissionRepo := new(MockCommissionRepo)
	orderRepo := new(MockOrderRepo)
	playerRepo := new(MockPlayerRepo)

	svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

	t.Run("创建默认规则", func(t *testing.T) {
		req := CreateCommissionRuleRequest{
			Name:        "默认抽成",
			Description: "平台默认20%抽成",
			Type:        "default",
			Rate:        20,
		}

		commissionRepo.On("CreateRule", ctx, mock.MatchedBy(func(rule *model.CommissionRule) bool {
			assert.Equal(t, "默认抽成", rule.Name)
			assert.Equal(t, 20, rule.Rate)
			assert.Equal(t, "default", rule.Type)
			assert.True(t, rule.IsActive)
			return true
		})).Return(nil)

		rule, err := svc.CreateCommissionRule(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, rule)
	})

	t.Run("抽成比例超出范围", func(t *testing.T) {
		req := CreateCommissionRuleRequest{
			Name: "无效规则",
			Type: "default",
			Rate: 150, // 超过100%
		}

		rule, err := svc.CreateCommissionRule(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, rule)
		assert.Contains(t, err.Error(), "between 0 and 100")
	})

	t.Run("创建游戏专属规则", func(t *testing.T) {
		commissionRepo := new(MockCommissionRepo)
		orderRepo := new(MockOrderRepo)
		playerRepo := new(MockPlayerRepo)

		svc := NewCommissionService(commissionRepo, orderRepo, playerRepo)

		gameID := uint64(1)
		req := CreateCommissionRuleRequest{
			Name:        "王者荣耀特殊抽成",
			Description: "王者荣耀15%抽成",
			Type:        "special",
			Rate:        15,
			GameID:      &gameID,
		}

		commissionRepo.On("CreateRule", ctx, mock.MatchedBy(func(rule *model.CommissionRule) bool {
			assert.Equal(t, 15, rule.Rate)
			assert.Equal(t, &gameID, rule.GameID)
			return true
		})).Return(nil)

		rule, err := svc.CreateCommissionRule(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, rule)
	})
}

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
