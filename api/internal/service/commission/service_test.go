package commission

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	commissionrepo "gamelink/internal/repository/commission"
	repoiface "gamelink/internal/repository/interfaces"
)

// MockCommissionRepository is a mock implementation of CommissionRepository
type MockCommissionRepository struct {
	mock.Mock
}

func (m *MockCommissionRepository) CreateRule(ctx context.Context, rule *model.CommissionRule) error {
	args := m.Called(ctx, rule)
	return args.Error(0)
}

func (m *MockCommissionRepository) GetRule(ctx context.Context, id uint64) (*model.CommissionRule, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CommissionRule), args.Error(1)
}

func (m *MockCommissionRepository) GetDefaultRule(ctx context.Context) (*model.CommissionRule, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CommissionRule), args.Error(1)
}

func (m *MockCommissionRepository) GetRuleForOrder(ctx context.Context, gameID *uint64, playerID *uint64, serviceType *string) (*model.CommissionRule, error) {
	args := m.Called(ctx, gameID, playerID, serviceType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CommissionRule), args.Error(1)
}

func (m *MockCommissionRepository) ListRules(ctx context.Context, opts commissionrepo.CommissionRuleListOptions) ([]model.CommissionRule, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]model.CommissionRule), args.Get(1).(int64), args.Error(2)
}

func (m *MockCommissionRepository) UpdateRule(ctx context.Context, rule *model.CommissionRule) error {
	args := m.Called(ctx, rule)
	return args.Error(0)
}

func (m *MockCommissionRepository) DeleteRule(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCommissionRepository) CreateRecord(ctx context.Context, record *model.CommissionRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockCommissionRepository) GetRecord(ctx context.Context, id uint64) (*model.CommissionRecord, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CommissionRecord), args.Error(1)
}

func (m *MockCommissionRepository) GetRecordByOrderID(ctx context.Context, orderID uint64) (*model.CommissionRecord, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CommissionRecord), args.Error(1)
}

func (m *MockCommissionRepository) ListRecords(ctx context.Context, opts commissionrepo.CommissionRecordListOptions) ([]model.CommissionRecord, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]model.CommissionRecord), args.Get(1).(int64), args.Error(2)
}

func (m *MockCommissionRepository) UpdateRecord(ctx context.Context, record *model.CommissionRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockCommissionRepository) CreateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error {
	args := m.Called(ctx, settlement)
	return args.Error(0)
}

func (m *MockCommissionRepository) GetSettlement(ctx context.Context, id uint64) (*model.MonthlySettlement, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.MonthlySettlement), args.Error(1)
}

func (m *MockCommissionRepository) GetSettlementByPlayerMonth(ctx context.Context, playerID uint64, month string) (*model.MonthlySettlement, error) {
	args := m.Called(ctx, playerID, month)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.MonthlySettlement), args.Error(1)
}

func (m *MockCommissionRepository) ListSettlements(ctx context.Context, opts commissionrepo.SettlementListOptions) ([]model.MonthlySettlement, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]model.MonthlySettlement), args.Get(1).(int64), args.Error(2)
}

func (m *MockCommissionRepository) UpdateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error {
	args := m.Called(ctx, settlement)
	return args.Error(0)
}

func (m *MockCommissionRepository) GetMonthlyStats(ctx context.Context, month string) (*commissionrepo.MonthlyStats, error) {
	args := m.Called(ctx, month)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*commissionrepo.MonthlyStats), args.Error(1)
}

func (m *MockCommissionRepository) GetPlayerMonthlyIncome(ctx context.Context, playerID uint64, month string) (int64, error) {
	args := m.Called(ctx, playerID, month)
	return args.Get(0).(int64), args.Error(1)
}

// MockOrderReader is a mock implementation of OrderReader
type MockOrderReader struct {
	mock.Mock
}

func (m *MockOrderReader) Get(ctx context.Context, id uint64) (*model.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderReader) List(ctx context.Context, opts repoiface.OrderListOptions) ([]model.Order, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]model.Order), args.Get(1).(int64), args.Error(2)
}

func (m *MockOrderReader) GetByPlayerID(ctx context.Context, playerID uint64, page, pageSize int) ([]model.Order, int64, error) {
	args := m.Called(ctx, playerID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]model.Order), args.Get(1).(int64), args.Error(2)
}

func (m *MockOrderReader) CountByPlayerID(ctx context.Context, playerID uint64) (int64, error) {
	args := m.Called(ctx, playerID)
	return args.Get(0).(int64), args.Error(1)
}

// MockPlayerRepository is a mock implementation of PlayerRepository
type MockPlayerRepository struct {
	mock.Mock
}

func (m *MockPlayerRepository) List(ctx context.Context) ([]model.Player, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Player), args.Error(1)
}

func (m *MockPlayerRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.Player, int64, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]model.Player), args.Get(1).(int64), args.Error(2)
}

func (m *MockPlayerRepository) ListPagedWithFilter(ctx context.Context, page, pageSize int, keyword string, status *model.VerificationStatus) ([]model.Player, int64, error) {
	args := m.Called(ctx, page, pageSize, keyword, status)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]model.Player), args.Get(1).(int64), args.Error(2)
}

func (m *MockPlayerRepository) Get(ctx context.Context, id uint64) (*model.Player, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Player), args.Error(1)
}

func (m *MockPlayerRepository) GetByUserID(ctx context.Context, userID uint64) (*model.Player, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Player), args.Error(1)
}

func (m *MockPlayerRepository) Create(ctx context.Context, player *model.Player) error {
	args := m.Called(ctx, player)
	return args.Error(0)
}

func (m *MockPlayerRepository) Update(ctx context.Context, player *model.Player) error {
	args := m.Called(ctx, player)
	return args.Error(0)
}

func (m *MockPlayerRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPlayerRepository) BatchUpdateRank(ctx context.Context, ids []uint64, rank string) (int64, error) {
	args := m.Called(ctx, ids, rank)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPlayerRepository) BatchUpdateHourlyRate(ctx context.Context, ids []uint64, rateCents int64) (int64, error) {
	args := m.Called(ctx, ids, rateCents)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPlayerRepository) BatchUpdateStatus(ctx context.Context, ids []uint64, status model.VerificationStatus) (int64, error) {
	args := m.Called(ctx, ids, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPlayerRepository) BatchDelete(ctx context.Context, ids []uint64) (int64, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPlayerRepository) GetByIDs(ctx context.Context, ids []uint64) ([]model.Player, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Player), args.Error(1)
}

func (m *MockPlayerRepository) ListFeatured(ctx context.Context, limit int, status *model.VerificationStatus) ([]model.Player, int64, error) {
	args := m.Called(ctx, limit, status)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]model.Player), args.Get(1).(int64), args.Error(2)
}

// Helper function to create test order
func createTestOrder(playerID uint64, totalPriceCents int64) *model.Order {
	pID := playerID
	gID := uint64(1)
	return &model.Order{
		ItemID:          1,
		GameID:          &gID,
		TotalPriceCents: totalPriceCents,
		PlayerID:        &pID,
		Status:          model.OrderStatusCompleted,
	}
}

// Helper function to create test commission rule
func createTestCommissionRule(id uint64, ruleType model.CommissionRuleType, rate int, isActive bool) *model.CommissionRule {
	return &model.CommissionRule{
		ID:       id,
		Name:     "Test Rule",
		Type:     ruleType,
		Rate:     rate,
		IsActive: isActive,
	}
}

// TestCommissionService_CalculateCommission_WithDefaultRate tests commission calculation with default rate
func TestCommissionService_CalculateCommission_WithDefaultRate(t *testing.T) {
	ctx := context.Background()
	mockCommissions := new(MockCommissionRepository)
	mockOrders := new(MockOrderReader)
	mockPlayers := new(MockPlayerRepository)

	service := NewCommissionService(mockCommissions, mockOrders, mockPlayers)

	order := createTestOrder(100, 10000) // ¥100 order

	mockOrders.On("Get", ctx, uint64(1)).Return(order, nil)
	mockCommissions.On("GetDefaultRule", ctx).Return(createTestCommissionRule(1, model.CommissionRuleTypeDefault, 20, true), nil)
	mockCommissions.On("GetRuleForOrder", ctx, mock.Anything, mock.Anything, mock.Anything).Return(nil, repository.ErrNotFound)

	result, err := service.CalculateCommission(ctx, 1)

	require.NoError(t, err)
	assert.Equal(t, int64(10000), result.TotalAmountCents)
	assert.Equal(t, 20, result.CommissionRate)
	assert.Equal(t, int64(2000), result.CommissionCents)
	assert.Equal(t, int64(8000), result.PlayerIncomeCents)
	assert.Equal(t, "默认规则", result.AppliedRule)

	mockOrders.AssertExpectations(t)
}

// TestCommissionService_CalculateCommission_WithPlayerIndividualRate tests commission calculation with player individual rate
func TestCommissionService_CalculateCommission_WithPlayerIndividualRate(t *testing.T) {
	ctx := context.Background()
	mockCommissions := new(MockCommissionRepository)
	mockOrders := new(MockOrderReader)
	mockPlayers := new(MockPlayerRepository)

	service := NewCommissionService(mockCommissions, mockOrders, mockPlayers)

	order := createTestOrder(100, 10000) // ¥100 order
	playerID := uint64(100)

	mockOrders.On("Get", ctx, uint64(1)).Return(order, nil)

	// Player individual rate of 15%
	playerRule := createTestCommissionRule(1, model.CommissionRuleTypeSpecial, 15, true)
	playerRule.PlayerID = &playerID
	mockCommissions.On("GetRuleForOrder", ctx, mock.Anything, &playerID, mock.Anything).Return(playerRule, nil)

	result, err := service.CalculateCommission(ctx, 1)

	require.NoError(t, err)
	assert.Equal(t, 15, result.CommissionRate)
	assert.Equal(t, int64(1500), result.CommissionCents)
	assert.Equal(t, int64(8500), result.PlayerIncomeCents)
	assert.Equal(t, "陪玩师专属", result.AppliedRule)

	mockOrders.AssertExpectations(t)
}

// TestCommissionService_CalculateCommission_ThreeTierCalculation tests three-tier commission calculation
func TestCommissionService_CalculateCommission_ThreeTierCalculation(t *testing.T) {
	ctx := context.Background()
	mockCommissions := new(MockCommissionRepository)
	mockOrders := new(MockOrderReader)
	mockPlayers := new(MockPlayerRepository)

	service := NewCommissionService(mockCommissions, mockOrders, mockPlayers)

	t.Run("Player individual rate has lowest commission (15%)", func(t *testing.T) {
		order := createTestOrder(100, 10000)
		playerID := uint64(100)

		mockOrders.On("Get", ctx, uint64(1)).Return(order, nil)

		// Player individual rate: 15% (lowest)
		playerRule := createTestCommissionRule(1, model.CommissionRuleTypeSpecial, 15, true)
		playerRule.PlayerID = &playerID
		mockCommissions.On("GetRuleForOrder", ctx, mock.Anything, &playerID, mock.Anything).Return(playerRule, nil)

		result, err := service.CalculateCommission(ctx, 1)

		require.NoError(t, err)
		assert.Equal(t, 15, result.CommissionRate)
		assert.Equal(t, int64(1500), result.CommissionCents)
		assert.Equal(t, int64(8500), result.PlayerIncomeCents)

		mockOrders.AssertExpectations(t)
		mockCommissions.AssertExpectations(t)
	})

	t.Run("Player individual rate higher than default", func(t *testing.T) {
		mockCommissions := new(MockCommissionRepository)
		mockOrders := new(MockOrderReader)
		mockPlayers := new(MockPlayerRepository)
		service := NewCommissionService(mockCommissions, mockOrders, mockPlayers)

		order := createTestOrder(101, 10000)
		playerID := uint64(101)

		mockOrders.On("Get", ctx, uint64(2)).Return(order, nil)

		// Player individual rate: 25% (higher than default 20%)
		playerRule := createTestCommissionRule(2, model.CommissionRuleTypeSpecial, 25, true)
		playerRule.PlayerID = &playerID
		mockCommissions.On("GetRuleForOrder", ctx, mock.Anything, &playerID, mock.Anything).Return(playerRule, nil)

		result, err := service.CalculateCommission(ctx, 2)

		require.NoError(t, err)
		// Uses player rule even though it's higher
		assert.Equal(t, 25, result.CommissionRate)

		mockOrders.AssertExpectations(t)
		mockCommissions.AssertExpectations(t)
	})
}

// TestCommissionService_CalculateCommission_WithRankingDiscount tests commission calculation with ranking discount
func TestCommissionService_CalculateCommission_WithRankingDiscount(t *testing.T) {
	ctx := context.Background()
	mockCommissions := new(MockCommissionRepository)
	mockOrders := new(MockOrderReader)
	mockPlayers := new(MockPlayerRepository)

	service := NewCommissionService(mockCommissions, mockOrders, mockPlayers)

	order := createTestOrder(100, 10000) // ¥100 order

	mockOrders.On("Get", ctx, uint64(1)).Return(order, nil)
	mockCommissions.On("GetDefaultRule", ctx).Return(createTestCommissionRule(1, model.CommissionRuleTypeDefault, 20, true), nil)
	mockCommissions.On("GetRuleForOrder", ctx, mock.Anything, mock.Anything, mock.Anything).Return(nil, repository.ErrNotFound)

	// Note: Ranking discount is not fully implemented yet (returns 0)
	// This test will use default rule
	result, err := service.CalculateCommission(ctx, 1)

	require.NoError(t, err)
	assert.Equal(t, 20, result.CommissionRate)

	mockOrders.AssertExpectations(t)
}

// TestCommissionService_CalculateCommission_OrderNotFound tests commission calculation when order is not found
func TestCommissionService_CalculateCommission_OrderNotFound(t *testing.T) {
	ctx := context.Background()
	mockCommissions := new(MockCommissionRepository)
	mockOrders := new(MockOrderReader)
	mockPlayers := new(MockPlayerRepository)

	service := NewCommissionService(mockCommissions, mockOrders, mockPlayers)

	mockOrders.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

	result, err := service.CalculateCommission(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, repository.ErrNotFound, err)

	mockOrders.AssertExpectations(t)
}

// TestCommissionService_CalculateCommission_ZeroAmount tests commission calculation with zero amount
func TestCommissionService_CalculateCommission_ZeroAmount(t *testing.T) {
	ctx := context.Background()
	mockCommissions := new(MockCommissionRepository)
	mockOrders := new(MockOrderReader)
	mockPlayers := new(MockPlayerRepository)

	service := NewCommissionService(mockCommissions, mockOrders, mockPlayers)

	order := createTestOrder(100, 0) // ¥0 order

	mockOrders.On("Get", ctx, uint64(1)).Return(order, nil)
	mockCommissions.On("GetDefaultRule", ctx).Return(createTestCommissionRule(1, model.CommissionRuleTypeDefault, 20, true), nil)
	mockCommissions.On("GetRuleForOrder", ctx, mock.Anything, mock.Anything, mock.Anything).Return(nil, repository.ErrNotFound)

	result, err := service.CalculateCommission(ctx, 1)

	require.NoError(t, err)
	assert.Equal(t, int64(0), result.CommissionCents)
	assert.Equal(t, int64(0), result.PlayerIncomeCents)

	mockOrders.AssertExpectations(t)
}

// TestCommissionService_RecordCommission_Success tests successful commission recording
func TestCommissionService_RecordCommission_Success(t *testing.T) {
	ctx := context.Background()
	mockCommissions := new(MockCommissionRepository)
	mockOrders := new(MockOrderReader)
	mockPlayers := new(MockPlayerRepository)

	service := NewCommissionService(mockCommissions, mockOrders, mockPlayers)

	order := createTestOrder(100, 10000)

	mockOrders.On("Get", ctx, uint64(1)).Return(order, nil).Twice()
	mockCommissions.On("GetRecordByOrderID", ctx, uint64(1)).Return(nil, repository.ErrNotFound)
	mockCommissions.On("GetDefaultRule", ctx).Return(createTestCommissionRule(1, model.CommissionRuleTypeDefault, 20, true), nil)
	mockCommissions.On("GetRuleForOrder", ctx, mock.Anything, mock.Anything, mock.Anything).Return(nil, repository.ErrNotFound)
	mockCommissions.On("CreateRecord", ctx, mock.AnythingOfType("*model.CommissionRecord")).Return(nil)

	err := service.RecordCommission(ctx, 1)

	require.NoError(t, err)

	mockOrders.AssertExpectations(t)
	mockCommissions.AssertExpectations(t)
}

// TestCommissionService_RecordCommission_AlreadyRecorded tests commission recording when already recorded
func TestCommissionService_RecordCommission_AlreadyRecorded(t *testing.T) {
	ctx := context.Background()
	mockCommissions := new(MockCommissionRepository)
	mockOrders := new(MockOrderReader)
	mockPlayers := new(MockPlayerRepository)

	service := NewCommissionService(mockCommissions, mockOrders, mockPlayers)

	existingRecord := &model.CommissionRecord{ID: 1, OrderID: 1}
	mockCommissions.On("GetRecordByOrderID", ctx, uint64(1)).Return(existingRecord, nil)

	err := service.RecordCommission(ctx, 1)

	assert.Error(t, err)
	assert.Equal(t, ErrAlreadyRecorded, err)

	mockCommissions.AssertExpectations(t)
}

// TestCommissionService_RecordCommission_NoPlayerAssigned tests commission recording when no player is assigned
func TestCommissionService_RecordCommission_NoPlayerAssigned(t *testing.T) {
	ctx := context.Background()
	mockCommissions := new(MockCommissionRepository)
	mockOrders := new(MockOrderReader)
	mockPlayers := new(MockPlayerRepository)

	service := NewCommissionService(mockCommissions, mockOrders, mockPlayers)

	// Order with no player assigned
	order := createTestOrder(100, 10000)
	order.PlayerID = nil

	mockCommissions.On("GetRecordByOrderID", ctx, uint64(1)).Return(nil, repository.ErrNotFound)
	// Called twice: once in CalculateCommission, once in RecordCommission after player check
	mockOrders.On("Get", ctx, uint64(1)).Return(order, nil).Twice()
	mockCommissions.On("GetDefaultRule", ctx).Return(createTestCommissionRule(1, model.CommissionRuleTypeDefault, 20, true), nil)
	// GetRuleForOrder won't be called because playerID is 0 (line 469: if playerID > 0)

	err := service.RecordCommission(ctx, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "订单未分配打手")

	mockOrders.AssertExpectations(t)
	mockCommissions.AssertExpectations(t)
}

// TestCommissionService_SettleMonth_Success tests successful monthly settlement
func TestCommissionService_SettleMonth_Success(t *testing.T) {
	ctx := context.Background()
	mockCommissions := new(MockCommissionRepository)
	mockOrders := new(MockOrderReader)
	mockPlayers := new(MockPlayerRepository)

	service := NewCommissionService(mockCommissions, mockOrders, mockPlayers)

	month := "2025-01"

	playerID1 := uint64(100)
	playerID2 := uint64(101)

	records := []model.CommissionRecord{
		{
			ID:                1,
			OrderID:           1,
			PlayerID:          playerID1,
			TotalAmountCents:  10000,
			CommissionRate:    20,
			CommissionCents:   2000,
			PlayerIncomeCents: 8000,
			SettlementStatus:  model.SettlementStatusPending,
			SettlementMonth:   month,
		},
		{
			ID:                2,
			OrderID:           2,
			PlayerID:          playerID2,
			TotalAmountCents:  15000,
			CommissionRate:    20,
			CommissionCents:   3000,
			PlayerIncomeCents: 12000,
			SettlementStatus:  model.SettlementStatusPending,
			SettlementMonth:   month,
		},
	}

	statusPending := "pending"
	mockCommissions.On("ListSettlements", ctx, mock.MatchedBy(func(opts commissionrepo.SettlementListOptions) bool {
		return opts.SettlementMonth != nil && *opts.SettlementMonth == month && opts.Page == 1 && opts.PageSize == 1
	})).Return([]model.MonthlySettlement{}, int64(0), nil)

	mockCommissions.On("ListRecords", ctx, mock.MatchedBy(func(opts commissionrepo.CommissionRecordListOptions) bool {
		return opts.SettlementMonth != nil && *opts.SettlementMonth == month &&
			opts.SettlementStatus != nil && *opts.SettlementStatus == statusPending &&
			opts.Page == 1 && opts.PageSize == 10000
	})).Return(records, int64(2), nil)

	mockCommissions.On("CreateSettlement", ctx, mock.MatchedBy(func(s *model.MonthlySettlement) bool {
		return s.SettlementMonth == month
	})).Return(nil).Times(2)

	mockCommissions.On("UpdateRecord", ctx, mock.MatchedBy(func(r *model.CommissionRecord) bool {
		return r.SettlementStatus == "settled" && r.SettledAt != nil
	})).Return(nil).Times(2)

	err := service.SettleMonth(ctx, month)

	require.NoError(t, err)

	mockCommissions.AssertExpectations(t)
}

// TestCommissionService_SettleMonth_AlreadySettled tests monthly settlement when already settled
func TestCommissionService_SettleMonth_AlreadySettled(t *testing.T) {
	ctx := context.Background()
	mockCommissions := new(MockCommissionRepository)
	mockOrders := new(MockOrderReader)
	mockPlayers := new(MockPlayerRepository)

	service := NewCommissionService(mockCommissions, mockOrders, mockPlayers)

	month := "2025-01"

	existingSettlement := []model.MonthlySettlement{
		{ID: 1, SettlementMonth: month},
	}

	mockCommissions.On("ListSettlements", ctx, mock.MatchedBy(func(opts commissionrepo.SettlementListOptions) bool {
		return opts.SettlementMonth != nil && *opts.SettlementMonth == month && opts.Page == 1 && opts.PageSize == 1
	})).Return(existingSettlement, int64(1), nil)

	err := service.SettleMonth(ctx, month)

	assert.Error(t, err)
	assert.Equal(t, ErrAlreadySettled, err)

	mockCommissions.AssertExpectations(t)
}

// TestCommissionService_SettleMonth_NoRecords tests monthly settlement with no records
func TestCommissionService_SettleMonth_NoRecords(t *testing.T) {
	ctx := context.Background()
	mockCommissions := new(MockCommissionRepository)
	mockOrders := new(MockOrderReader)
	mockPlayers := new(MockPlayerRepository)

	service := NewCommissionService(mockCommissions, mockOrders, mockPlayers)

	month := "2025-01"

	statusPending := "pending"
	mockCommissions.On("ListSettlements", ctx, mock.MatchedBy(func(opts commissionrepo.SettlementListOptions) bool {
		return opts.SettlementMonth != nil && *opts.SettlementMonth == month && opts.Page == 1 && opts.PageSize == 1
	})).Return([]model.MonthlySettlement{}, int64(0), nil)

	mockCommissions.On("ListRecords", ctx, mock.MatchedBy(func(opts commissionrepo.CommissionRecordListOptions) bool {
		return opts.SettlementMonth != nil && *opts.SettlementMonth == month &&
			opts.SettlementStatus != nil && *opts.SettlementStatus == statusPending &&
			opts.Page == 1 && opts.PageSize == 10000
	})).Return([]model.CommissionRecord{}, int64(0), nil)

	err := service.SettleMonth(ctx, month)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "该月无待结算记录")

	mockCommissions.AssertExpectations(t)
}

// TestCommissionService_GetPlayerCommissionSummary_Success tests getting player commission summary
func TestCommissionService_GetPlayerCommissionSummary_Success(t *testing.T) {
	ctx := context.Background()
	mockCommissions := new(MockCommissionRepository)
	mockOrders := new(MockOrderReader)
	mockPlayers := new(MockPlayerRepository)

	service := NewCommissionService(mockCommissions, mockOrders, mockPlayers)

	playerID := uint64(100)
	month := "2025-01"

	records := []model.CommissionRecord{
		{
			ID:                1,
			OrderID:           1,
			PlayerID:          playerID,
			CommissionCents:   2000,
			PlayerIncomeCents: 8000,
		},
		{
			ID:                2,
			OrderID:           2,
			PlayerID:          playerID,
			CommissionCents:   3000,
			PlayerIncomeCents: 12000,
		},
	}

	mockCommissions.On("GetPlayerMonthlyIncome", ctx, playerID, month).Return(int64(8000), nil)
	mockCommissions.On("ListRecords", ctx, mock.MatchedBy(func(opts commissionrepo.CommissionRecordListOptions) bool {
		return opts.PlayerID != nil && *opts.PlayerID == playerID && opts.Page == 1 && opts.PageSize == 1
	})).Return(records[:1], int64(2), nil)

	mockCommissions.On("ListRecords", ctx, mock.MatchedBy(func(opts commissionrepo.CommissionRecordListOptions) bool {
		return opts.PlayerID != nil && *opts.PlayerID == playerID && opts.Page == 1 && opts.PageSize == 10000
	})).Return(records, int64(2), nil)

	summary, err := service.GetPlayerCommissionSummary(ctx, playerID, month)

	require.NoError(t, err)
	assert.Equal(t, int64(8000), summary.MonthlyIncome)
	assert.Equal(t, int64(5000), summary.TotalCommission)
	assert.Equal(t, int64(20000), summary.TotalIncome)
	assert.Equal(t, int64(2), summary.TotalOrders)

	mockCommissions.AssertExpectations(t)
}

// TestCommissionService_GetCommissionRecords_Success tests getting commission records
func TestCommissionService_GetCommissionRecords_Success(t *testing.T) {
	ctx := context.Background()
	mockCommissions := new(MockCommissionRepository)
	mockOrders := new(MockOrderReader)
	mockPlayers := new(MockPlayerRepository)

	service := NewCommissionService(mockCommissions, mockOrders, mockPlayers)

	playerID := uint64(100)
	page := 1
	pageSize := 20

	records := []model.CommissionRecord{
		{
			ID:                1,
			OrderID:           1,
			PlayerID:          playerID,
			TotalAmountCents:  10000,
			CommissionRate:    20,
			CommissionCents:   2000,
			PlayerIncomeCents: 8000,
			SettlementStatus:  model.SettlementStatusPending,
			SettlementMonth:   "2025-01",
			CreatedAt:         time.Now(),
		},
	}

	mockCommissions.On("ListRecords", ctx, mock.MatchedBy(func(opts commissionrepo.CommissionRecordListOptions) bool {
		return opts.PlayerID != nil && *opts.PlayerID == playerID &&
			opts.Page == page && opts.PageSize == pageSize
	})).Return(records, int64(1), nil)

	result, err := service.GetCommissionRecords(ctx, playerID, page, pageSize)

	require.NoError(t, err)
	assert.Equal(t, int64(1), result.Total)
	assert.Len(t, result.Records, 1)
	assert.Equal(t, uint64(1), result.Records[0].ID)
	assert.Equal(t, 20, result.Records[0].CommissionRate)
	assert.Equal(t, int64(2000), result.Records[0].CommissionCents)

	mockCommissions.AssertExpectations(t)
}

// TestCommissionService_GetMonthlySettlements_Success tests getting monthly settlements
func TestCommissionService_GetMonthlySettlements_Success(t *testing.T) {
	ctx := context.Background()
	mockCommissions := new(MockCommissionRepository)
	mockOrders := new(MockOrderReader)
	mockPlayers := new(MockPlayerRepository)

	service := NewCommissionService(mockCommissions, mockOrders, mockPlayers)

	playerID := uint64(100)
	page := 1
	pageSize := 20

	settlements := []model.MonthlySettlement{
		{
			ID:                   1,
			PlayerID:             playerID,
			SettlementMonth:      "2025-01",
			TotalOrderCount:      10,
			TotalAmountCents:     100000,
			TotalCommissionCents: 20000,
			TotalIncomeCents:     80000,
			BonusCents:           0,
			FinalIncomeCents:     80000,
			Status:               model.MonthlySettlementStatusPending,
			CreatedAt:            time.Now(),
		},
	}

	mockCommissions.On("ListSettlements", ctx, mock.MatchedBy(func(opts commissionrepo.SettlementListOptions) bool {
		return opts.PlayerID != nil && *opts.PlayerID == playerID &&
			opts.Page == page && opts.PageSize == pageSize
	})).Return(settlements, int64(1), nil)

	result, err := service.GetMonthlySettlements(ctx, playerID, page, pageSize)

	require.NoError(t, err)
	assert.Equal(t, int64(1), result.Total)
	assert.Len(t, result.Settlements, 1)
	assert.Equal(t, uint64(1), result.Settlements[0].ID)
	assert.Equal(t, "2025-01", result.Settlements[0].SettlementMonth)
	assert.Equal(t, int64(10), result.Settlements[0].TotalOrderCount)
	assert.Equal(t, int64(20000), result.Settlements[0].TotalCommissionCents)

	mockCommissions.AssertExpectations(t)
}

// TestCommissionService_CreateCommissionRule_Success tests successful commission rule creation
func TestCommissionService_CreateCommissionRule_Success(t *testing.T) {
	ctx := context.Background()
	mockCommissions := new(MockCommissionRepository)
	mockOrders := new(MockOrderReader)
	mockPlayers := new(MockPlayerRepository)

	service := NewCommissionService(mockCommissions, mockOrders, mockPlayers)

	req := CreateCommissionRuleRequest{
		Name:        "Test Rule",
		Description: "Test Description",
		Type:        model.CommissionRuleTypeSpecial,
		Rate:        15,
	}

	mockCommissions.On("CreateRule", ctx, mock.MatchedBy(func(rule *model.CommissionRule) bool {
		return rule.Name == "Test Rule" && rule.Rate == 15
	})).Return(nil)

	rule, err := service.CreateCommissionRule(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, rule)
	assert.Equal(t, "Test Rule", rule.Name)
	assert.Equal(t, 15, rule.Rate)
	assert.True(t, rule.IsActive)

	mockCommissions.AssertExpectations(t)
}

// TestCommissionService_CreateCommissionRule_InvalidRate tests commission rule creation with invalid rate
func TestCommissionService_CreateCommissionRule_InvalidRate(t *testing.T) {
	ctx := context.Background()
	mockCommissions := new(MockCommissionRepository)
	mockOrders := new(MockOrderReader)
	mockPlayers := new(MockPlayerRepository)

	service := NewCommissionService(mockCommissions, mockOrders, mockPlayers)

	t.Run("Rate below 0", func(t *testing.T) {
		req := CreateCommissionRuleRequest{
			Name: "Test Rule",
			Type: model.CommissionRuleTypeSpecial,
			Rate: -1,
		}

		rule, err := service.CreateCommissionRule(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, rule)
	})

	t.Run("Rate above 100", func(t *testing.T) {
		req := CreateCommissionRuleRequest{
			Name: "Test Rule",
			Type: model.CommissionRuleTypeSpecial,
			Rate: 101,
		}

		rule, err := service.CreateCommissionRule(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, rule)
	})
}

// TestCommissionService_UpdateCommissionRule_Success tests successful commission rule update
func TestCommissionService_UpdateCommissionRule_Success(t *testing.T) {
	ctx := context.Background()
	mockCommissions := new(MockCommissionRepository)
	mockOrders := new(MockOrderReader)
	mockPlayers := new(MockPlayerRepository)

	service := NewCommissionService(mockCommissions, mockOrders, mockPlayers)

	existingRule := createTestCommissionRule(1, model.CommissionRuleTypeSpecial, 20, true)
	newName := "Updated Rule"
	newRate := 15

	mockCommissions.On("GetRule", ctx, uint64(1)).Return(existingRule, nil)
	mockCommissions.On("UpdateRule", ctx, mock.MatchedBy(func(rule *model.CommissionRule) bool {
		return rule.Name == newName && rule.Rate == newRate
	})).Return(nil)

	req := UpdateCommissionRuleRequest{
		Name: &newName,
		Rate: &newRate,
	}

	err := service.UpdateCommissionRule(ctx, 1, req)

	require.NoError(t, err)

	mockCommissions.AssertExpectations(t)
}

// TestCommissionService_UpdateCommissionRule_RuleNotFound tests commission rule update when rule not found
func TestCommissionService_UpdateCommissionRule_RuleNotFound(t *testing.T) {
	ctx := context.Background()
	mockCommissions := new(MockCommissionRepository)
	mockOrders := new(MockOrderReader)
	mockPlayers := new(MockPlayerRepository)

	service := NewCommissionService(mockCommissions, mockOrders, mockPlayers)

	mockCommissions.On("GetRule", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

	req := UpdateCommissionRuleRequest{
		Name: strPtr("Updated Name"),
	}

	err := service.UpdateCommissionRule(ctx, 999, req)

	assert.Error(t, err)
	assert.Equal(t, repository.ErrNotFound, err)

	mockCommissions.AssertExpectations(t)
}

// TestCommissionService_UpdateCommissionRule_InvalidRate tests commission rule update with invalid rate
func TestCommissionService_UpdateCommissionRule_InvalidRate(t *testing.T) {
	ctx := context.Background()
	mockCommissions := new(MockCommissionRepository)
	mockOrders := new(MockOrderReader)
	mockPlayers := new(MockPlayerRepository)

	service := NewCommissionService(mockCommissions, mockOrders, mockPlayers)

	existingRule := createTestCommissionRule(1, model.CommissionRuleTypeSpecial, 20, true)

	mockCommissions.On("GetRule", ctx, uint64(1)).Return(existingRule, nil)

	t.Run("Rate below 0", func(t *testing.T) {
		invalidRate := -1
		req := UpdateCommissionRuleRequest{
			Rate: &invalidRate,
		}

		err := service.UpdateCommissionRule(ctx, 1, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "0-100")
	})

	t.Run("Rate above 100", func(t *testing.T) {
		invalidRate := 101
		req := UpdateCommissionRuleRequest{
			Rate: &invalidRate,
		}

		err := service.UpdateCommissionRule(ctx, 1, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "0-100")
	})
}

// TestCommissionService_GetPlatformStats_Success tests getting platform statistics
func TestCommissionService_GetPlatformStats_Success(t *testing.T) {
	ctx := context.Background()
	mockCommissions := new(MockCommissionRepository)
	mockOrders := new(MockOrderReader)
	mockPlayers := new(MockPlayerRepository)

	service := NewCommissionService(mockCommissions, mockOrders, mockPlayers)

	month := "2025-01"

	stats := &commissionrepo.MonthlyStats{
		TotalOrders:       100,
		TotalIncome:       100000,
		TotalCommission:   20000,
		TotalPlayerIncome: 80000,
	}

	mockCommissions.On("GetMonthlyStats", ctx, month).Return(stats, nil)

	result, err := service.GetPlatformStats(ctx, month)

	require.NoError(t, err)
	assert.Equal(t, month, result.Month)
	assert.Equal(t, int64(100), result.TotalOrders)
	assert.Equal(t, int64(100000), result.TotalIncome)
	assert.Equal(t, int64(20000), result.TotalCommission)
	assert.Equal(t, int64(80000), result.TotalPlayerIncome)

	mockCommissions.AssertExpectations(t)
}

// TestCommissionService_BatchDeleteCommissionRules_Success tests successful batch deletion
func TestCommissionService_BatchDeleteCommissionRules_Success(t *testing.T) {
	ctx := context.Background()
	mockCommissions := new(MockCommissionRepository)
	mockOrders := new(MockOrderReader)
	mockPlayers := new(MockPlayerRepository)

	service := NewCommissionService(mockCommissions, mockOrders, mockPlayers)

	ids := []uint64{1, 2, 3}

	mockCommissions.On("DeleteRule", ctx, uint64(1)).Return(nil)
	mockCommissions.On("DeleteRule", ctx, uint64(2)).Return(nil)
	mockCommissions.On("DeleteRule", ctx, uint64(3)).Return(nil)

	result, err := service.BatchDeleteCommissionRules(ctx, ids)

	require.NoError(t, err)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Len(t, result.FailedIDs, 0)
	assert.Len(t, result.Errors, 0)

	mockCommissions.AssertExpectations(t)
}

// TestCommissionService_BatchDeleteCommissionRules_PartialFailure tests batch deletion with partial failures
func TestCommissionService_BatchDeleteCommissionRules_PartialFailure(t *testing.T) {
	ctx := context.Background()
	mockCommissions := new(MockCommissionRepository)
	mockOrders := new(MockOrderReader)
	mockPlayers := new(MockPlayerRepository)

	service := NewCommissionService(mockCommissions, mockOrders, mockPlayers)

	ids := []uint64{1, 2, 3}

	mockCommissions.On("DeleteRule", ctx, uint64(1)).Return(nil)
	mockCommissions.On("DeleteRule", ctx, uint64(2)).Return(errors.New("not found"))
	mockCommissions.On("DeleteRule", ctx, uint64(3)).Return(nil)

	result, err := service.BatchDeleteCommissionRules(ctx, ids)

	require.NoError(t, err)
	assert.Equal(t, 2, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Len(t, result.FailedIDs, 1)
	assert.Len(t, result.Errors, 1)

	mockCommissions.AssertExpectations(t)
}

// TestCommissionService_BatchDeleteCommissionRules_EmptyList tests batch deletion with empty list
func TestCommissionService_BatchDeleteCommissionRules_EmptyList(t *testing.T) {
	ctx := context.Background()
	mockCommissions := new(MockCommissionRepository)
	mockOrders := new(MockOrderReader)
	mockPlayers := new(MockPlayerRepository)

	service := NewCommissionService(mockCommissions, mockOrders, mockPlayers)

	result, err := service.BatchDeleteCommissionRules(ctx, []uint64{})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "不能为空")
}

// TestCommissionService_BatchUpdateCommissionRuleStatus_Success tests successful batch status update
func TestCommissionService_BatchUpdateCommissionRuleStatus_Success(t *testing.T) {
	ctx := context.Background()
	mockCommissions := new(MockCommissionRepository)
	mockOrders := new(MockOrderReader)
	mockPlayers := new(MockPlayerRepository)

	service := NewCommissionService(mockCommissions, mockOrders, mockPlayers)

	ids := []uint64{1, 2, 3}
	isActive := false

	rule1 := createTestCommissionRule(1, model.CommissionRuleTypeSpecial, 20, true)
	rule2 := createTestCommissionRule(2, model.CommissionRuleTypeSpecial, 20, true)
	rule3 := createTestCommissionRule(3, model.CommissionRuleTypeSpecial, 20, true)

	mockCommissions.On("GetRule", ctx, uint64(1)).Return(rule1, nil)
	mockCommissions.On("GetRule", ctx, uint64(2)).Return(rule2, nil)
	mockCommissions.On("GetRule", ctx, uint64(3)).Return(rule3, nil)

	mockCommissions.On("UpdateRule", ctx, mock.MatchedBy(func(rule *model.CommissionRule) bool {
		return rule.IsActive == false
	})).Return(nil).Times(3)

	result, err := service.BatchUpdateCommissionRuleStatus(ctx, ids, isActive)

	require.NoError(t, err)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)

	mockCommissions.AssertExpectations(t)
}

// TestCommissionService_BatchUpdateCommissionRuleStatus_EmptyList tests batch status update with empty list
func TestCommissionService_BatchUpdateCommissionRuleStatus_EmptyList(t *testing.T) {
	ctx := context.Background()
	mockCommissions := new(MockCommissionRepository)
	mockOrders := new(MockOrderReader)
	mockPlayers := new(MockPlayerRepository)

	service := NewCommissionService(mockCommissions, mockOrders, mockPlayers)

	result, err := service.BatchUpdateCommissionRuleStatus(ctx, []uint64{}, true)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "不能为空")
}

// TestParseRankingCommissionRules_Success tests successful parsing of ranking commission rules
func TestParseRankingCommissionRules_Success(t *testing.T) {
	rulesJSON := `[
		{
			"rankStart": 1,
			"rankEnd": 3,
			"commissionRate": 15
		},
		{
			"rankStart": 4,
			"rankEnd": 10,
			"commissionRate": 18
		}
	]`

	rules, err := ParseRankingCommissionRules(rulesJSON)

	require.NoError(t, err)
	assert.Len(t, rules, 2)
	assert.Equal(t, 1, rules[0].RankStart)
	assert.Equal(t, 3, rules[0].RankEnd)
	assert.Equal(t, 15, rules[0].CommissionRate)
	assert.Equal(t, 4, rules[1].RankStart)
	assert.Equal(t, 10, rules[1].RankEnd)
	assert.Equal(t, 18, rules[1].CommissionRate)
}

// TestParseRankingCommissionRules_InvalidJSON tests parsing with invalid JSON
func TestParseRankingCommissionRules_InvalidJSON(t *testing.T) {
	rulesJSON := `{invalid json}`

	rules, err := ParseRankingCommissionRules(rulesJSON)

	assert.Error(t, err)
	assert.Nil(t, rules)
}

// TestFindCommissionRateForRank_Success tests finding commission rate for a rank
func TestFindCommissionRateForRank_Success(t *testing.T) {
	rules := []model.RankingCommissionRule{
		{RankStart: 1, RankEnd: 3, CommissionRate: 15},
		{RankStart: 4, RankEnd: 10, CommissionRate: 18},
		{RankStart: 11, RankEnd: 50, CommissionRate: 20},
	}

	t.Run("Rank in first tier", func(t *testing.T) {
		rate := FindCommissionRateForRank(rules, 2)
		assert.Equal(t, 15, rate)
	})

	t.Run("Rank in second tier", func(t *testing.T) {
		rate := FindCommissionRateForRank(rules, 5)
		assert.Equal(t, 18, rate)
	})

	t.Run("Rank in third tier", func(t *testing.T) {
		rate := FindCommissionRateForRank(rules, 25)
		assert.Equal(t, 20, rate)
	})

	t.Run("Rank at boundary", func(t *testing.T) {
		rate := FindCommissionRateForRank(rules, 3)
		assert.Equal(t, 15, rate)
	})

	t.Run("Rank not in any range", func(t *testing.T) {
		rate := FindCommissionRateForRank(rules, 51)
		assert.Equal(t, 0, rate)
	})
}

// TestValidateRankingRules_ValidRules tests validation of valid ranking rules
func TestValidateRankingRules_ValidRules(t *testing.T) {
	rules := []model.RankingCommissionRule{
		{RankStart: 1, RankEnd: 3, CommissionRate: 15},
		{RankStart: 4, RankEnd: 10, CommissionRate: 18},
	}

	err := ValidateRankingRules(rules)

	assert.NoError(t, err)
}

// TestValidateRankingRules_InvalidRankRange tests validation with invalid rank ranges
func TestValidateRankingRules_InvalidRankRange(t *testing.T) {
	t.Run("Rank start < 1", func(t *testing.T) {
		rules := []model.RankingCommissionRule{
			{RankStart: 0, RankEnd: 3, CommissionRate: 15},
		}

		err := ValidateRankingRules(rules)

		assert.Error(t, err)
		assert.Equal(t, ErrValidation, err)
	})

	t.Run("Rank end < start", func(t *testing.T) {
		rules := []model.RankingCommissionRule{
			{RankStart: 5, RankEnd: 3, CommissionRate: 15},
		}

		err := ValidateRankingRules(rules)

		assert.Error(t, err)
		assert.Equal(t, ErrValidation, err)
	})
}

// TestValidateRankingRules_InvalidCommissionRate tests validation with invalid commission rate
func TestValidateRankingRules_InvalidCommissionRate(t *testing.T) {
	t.Run("Rate < 0", func(t *testing.T) {
		rules := []model.RankingCommissionRule{
			{RankStart: 1, RankEnd: 3, CommissionRate: -1},
		}

		err := ValidateRankingRules(rules)

		assert.Error(t, err)
		assert.Equal(t, ErrValidation, err)
	})

	t.Run("Rate > 100", func(t *testing.T) {
		rules := []model.RankingCommissionRule{
			{RankStart: 1, RankEnd: 3, CommissionRate: 101},
		}

		err := ValidateRankingRules(rules)

		assert.Error(t, err)
		assert.Equal(t, ErrValidation, err)
	})
}

// TestValidateRankingRules_OverlappingRanges tests validation with overlapping ranges
func TestValidateRankingRules_OverlappingRanges(t *testing.T) {
	rules := []model.RankingCommissionRule{
		{RankStart: 1, RankEnd: 5, CommissionRate: 15},
		{RankStart: 3, RankEnd: 10, CommissionRate: 18}, // Overlaps with first rule
	}

	err := ValidateRankingRules(rules)

	assert.Error(t, err)
	assert.Equal(t, ErrValidation, err)
}

// TestSelectLowestRate tests selecting the lowest commission rate
func TestSelectLowestRate(t *testing.T) {
	t.Run("Single candidate", func(t *testing.T) {
		candidates := []CommissionCandidate{
			{Source: "Test", Rate: 20, Detail: "Test rule"},
		}

		lowest := selectLowestRate(candidates)

		assert.Equal(t, 20, lowest.Rate)
		assert.Equal(t, "Test", lowest.Source)
	})

	t.Run("Multiple candidates", func(t *testing.T) {
		candidates := []CommissionCandidate{
			{Source: "Service Item", Rate: 20, Detail: "Standard"},
			{Source: "Player", Rate: 15, Detail: "VIP Player"},
			{Source: "Ranking", Rate: 18, Detail: "Top 10"},
		}

		lowest := selectLowestRate(candidates)

		assert.Equal(t, 15, lowest.Rate)
		assert.Equal(t, "Player", lowest.Source)
	})

	t.Run("Empty candidates", func(t *testing.T) {
		candidates := []CommissionCandidate{}

		lowest := selectLowestRate(candidates)

		assert.Equal(t, 20, lowest.Rate)
		assert.Equal(t, "默认规则", lowest.Source)
	})
}

// Helper function
func strPtr(s string) *string {
	return &s
}
