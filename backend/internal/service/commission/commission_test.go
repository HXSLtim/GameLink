package commission

import (
	"context"
	"testing"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	commissionrepo "gamelink/internal/repository/commission"

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
