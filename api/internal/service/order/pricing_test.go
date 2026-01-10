package order

import (
	"context"
	"testing"

	"gamelink/internal/model"
	"gamelink/internal/repository/commission"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockCommissionRepo implements commission.CommissionRepository for pricing tests
type mockCommissionRepo struct {
	mock.Mock
}

func (m *mockCommissionRepo) CreateRule(ctx context.Context, rule *model.CommissionRule) error {
	return m.Called(ctx, rule).Error(0)
}

func (m *mockCommissionRepo) GetRule(ctx context.Context, id uint64) (*model.CommissionRule, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CommissionRule), args.Error(1)
}

func (m *mockCommissionRepo) GetDefaultRule(ctx context.Context) (*model.CommissionRule, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CommissionRule), args.Error(1)
}

func (m *mockCommissionRepo) GetRuleForOrder(ctx context.Context, gameID, playerID *uint64, serviceType *string) (*model.CommissionRule, error) {
	args := m.Called(ctx, gameID, playerID, serviceType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CommissionRule), args.Error(1)
}

func (m *mockCommissionRepo) ListRules(ctx context.Context, opts commission.CommissionRuleListOptions) ([]model.CommissionRule, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.CommissionRule), args.Get(1).(int64), args.Error(2)
}

func (m *mockCommissionRepo) UpdateRule(ctx context.Context, rule *model.CommissionRule) error {
	return m.Called(ctx, rule).Error(0)
}

func (m *mockCommissionRepo) DeleteRule(ctx context.Context, id uint64) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockCommissionRepo) CreateRecord(ctx context.Context, record *model.CommissionRecord) error {
	return m.Called(ctx, record).Error(0)
}

func (m *mockCommissionRepo) GetRecord(ctx context.Context, id uint64) (*model.CommissionRecord, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CommissionRecord), args.Error(1)
}

func (m *mockCommissionRepo) GetRecordByOrderID(ctx context.Context, orderID uint64) (*model.CommissionRecord, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CommissionRecord), args.Error(1)
}

func (m *mockCommissionRepo) ListRecords(ctx context.Context, opts commission.CommissionRecordListOptions) ([]model.CommissionRecord, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.CommissionRecord), args.Get(1).(int64), args.Error(2)
}

func (m *mockCommissionRepo) UpdateRecord(ctx context.Context, record *model.CommissionRecord) error {
	return m.Called(ctx, record).Error(0)
}

func (m *mockCommissionRepo) CreateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error {
	return m.Called(ctx, settlement).Error(0)
}

func (m *mockCommissionRepo) GetSettlement(ctx context.Context, id uint64) (*model.MonthlySettlement, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.MonthlySettlement), args.Error(1)
}

func (m *mockCommissionRepo) GetSettlementByPlayerMonth(ctx context.Context, playerID uint64, month string) (*model.MonthlySettlement, error) {
	args := m.Called(ctx, playerID, month)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.MonthlySettlement), args.Error(1)
}

func (m *mockCommissionRepo) ListSettlements(ctx context.Context, opts commission.SettlementListOptions) ([]model.MonthlySettlement, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.MonthlySettlement), args.Get(1).(int64), args.Error(2)
}

func (m *mockCommissionRepo) UpdateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error {
	return m.Called(ctx, settlement).Error(0)
}

func (m *mockCommissionRepo) GetMonthlyStats(ctx context.Context, month string) (*commission.MonthlyStats, error) {
	args := m.Called(ctx, month)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*commission.MonthlyStats), args.Error(1)
}

func (m *mockCommissionRepo) GetPlayerMonthlyIncome(ctx context.Context, playerID uint64, month string) (int64, error) {
	args := m.Called(ctx, playerID, month)
	return args.Get(0).(int64), args.Error(1)
}

func TestCalculateOrderPricing_WithPlayerSpecificRule(t *testing.T) {
	mockCommissions := new(mockCommissionRepo)
	svc := &OrderService{commissions: mockCommissions}

	player := &model.Player{
		Base:            model.Base{ID: 1},
		HourlyRateCents: 5000, // ¥50/hour
	}
	req := CreateOrderRequest{
		GameID:        1,
		PlayerID:      1,
		DurationHours: 2,
	}

	// 陪玩师个人抽成 15%
	mockCommissions.On("GetRuleForOrder", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&model.CommissionRule{Rate: 15}, nil)

	totalPrice, commissionCents, playerIncomeCents := svc.calculateOrderPricing(player, req)

	assert.Equal(t, int64(10000), totalPrice)       // ¥100
	assert.Equal(t, int64(1500), commissionCents)   // 15% = ¥15
	assert.Equal(t, int64(8500), playerIncomeCents) // ¥85
	mockCommissions.AssertExpectations(t)
}

func TestCalculateOrderPricing_WithDefaultRule(t *testing.T) {
	mockCommissions := new(mockCommissionRepo)
	svc := &OrderService{commissions: mockCommissions}

	player := &model.Player{
		Base:            model.Base{ID: 2},
		HourlyRateCents: 6000, // ¥60/hour
	}
	req := CreateOrderRequest{
		GameID:        2,
		PlayerID:      2,
		DurationHours: 1,
	}

	// 无特定规则，使用默认规则 20%
	mockCommissions.On("GetRuleForOrder", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil)
	mockCommissions.On("GetDefaultRule", mock.Anything).
		Return(&model.CommissionRule{Rate: 20}, nil)

	totalPrice, commissionCents, playerIncomeCents := svc.calculateOrderPricing(player, req)

	assert.Equal(t, int64(6000), totalPrice)        // ¥60
	assert.Equal(t, int64(1200), commissionCents)   // 20% = ¥12
	assert.Equal(t, int64(4800), playerIncomeCents) // ¥48
	mockCommissions.AssertExpectations(t)
}

func TestCalculateOrderPricing_FallbackToHardcodedDefault(t *testing.T) {
	mockCommissions := new(mockCommissionRepo)
	svc := &OrderService{commissions: mockCommissions}

	player := &model.Player{
		Base:            model.Base{ID: 3},
		HourlyRateCents: 4000, // ¥40/hour
	}
	req := CreateOrderRequest{
		GameID:        3,
		PlayerID:      3,
		DurationHours: 1.5,
	}

	// 所有规则查询都失败，使用硬编码默认值 20%
	mockCommissions.On("GetRuleForOrder", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil)
	mockCommissions.On("GetDefaultRule", mock.Anything).
		Return(nil, nil)

	totalPrice, commissionCents, playerIncomeCents := svc.calculateOrderPricing(player, req)

	assert.Equal(t, int64(6000), totalPrice)        // ¥60
	assert.Equal(t, int64(1200), commissionCents)   // 20% = ¥12
	assert.Equal(t, int64(4800), playerIncomeCents) // ¥48
	mockCommissions.AssertExpectations(t)
}

func TestCalculateOrderPricing_NoCommissionRepo(t *testing.T) {
	// 没有注入 commissions repository
	svc := &OrderService{commissions: nil}

	player := &model.Player{
		Base:            model.Base{ID: 4},
		HourlyRateCents: 3000, // ¥30/hour
	}
	req := CreateOrderRequest{
		GameID:        4,
		PlayerID:      4,
		DurationHours: 2,
	}

	totalPrice, commissionCents, playerIncomeCents := svc.calculateOrderPricing(player, req)

	// 使用硬编码默认值 20%
	assert.Equal(t, int64(6000), totalPrice)        // ¥60
	assert.Equal(t, int64(1200), commissionCents)   // 20% = ¥12
	assert.Equal(t, int64(4800), playerIncomeCents) // ¥48
}

func TestGetCommissionRate_Priority(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func(*mockCommissionRepo)
		expectedRate int
	}{
		{
			name: "player specific rule takes priority",
			setupMock: func(m *mockCommissionRepo) {
				m.On("GetRuleForOrder", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(&model.CommissionRule{Rate: 15}, nil)
			},
			expectedRate: 15,
		},
		{
			name: "default rule when no specific rule",
			setupMock: func(m *mockCommissionRepo) {
				m.On("GetRuleForOrder", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil, nil)
				m.On("GetDefaultRule", mock.Anything).
					Return(&model.CommissionRule{Rate: 18}, nil)
			},
			expectedRate: 18,
		},
		{
			name: "hardcoded default when all rules fail",
			setupMock: func(m *mockCommissionRepo) {
				m.On("GetRuleForOrder", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil, nil)
				m.On("GetDefaultRule", mock.Anything).
					Return(nil, nil)
			},
			expectedRate: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCommissions := new(mockCommissionRepo)
			tt.setupMock(mockCommissions)

			svc := &OrderService{commissions: mockCommissions}
			gameID := uint64(1)
			playerID := uint64(1)

			rate := svc.getCommissionRate(context.Background(), &gameID, &playerID)

			assert.Equal(t, tt.expectedRate, rate)
			mockCommissions.AssertExpectations(t)
		})
	}
}
