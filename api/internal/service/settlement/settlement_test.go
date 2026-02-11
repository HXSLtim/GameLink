package settlement

import (
	"context"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository/commission"
	"gamelink/internal/repository/common"
	"gamelink/pkg/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock CommissionRepository
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

// Mock WalletRepository
type mockWalletRepo struct {
	mock.Mock
}

func (m *mockWalletRepo) GetByUserID(ctx context.Context, userID uint64) (*model.Wallet, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Wallet), args.Error(1)
}

func (m *mockWalletRepo) Save(ctx context.Context, wallet *model.Wallet) error {
	return m.Called(ctx, wallet).Error(0)
}

func (m *mockWalletRepo) SaveWithOptimisticLock(ctx context.Context, wallet *model.Wallet) error {
	return m.Called(ctx, wallet).Error(0)
}

func (m *mockWalletRepo) UpdateBalanceWithLock(ctx context.Context, userID uint64, delta int64, maxRetries int) (*model.Wallet, error) {
	args := m.Called(ctx, userID, delta, maxRetries)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Wallet), args.Error(1)
}

// TestProcessT7Unfreeze_Success uses in-memory SQLite DB because ProcessT7Unfreeze runs
// a gorm.DB.Transaction internally which creates real repositories from the DB connection.
func TestProcessT7Unfreeze_Success(t *testing.T) {
	ctx := context.Background()

	// Setup in-memory DB for transaction support
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db, &model.Wallet{}, &model.CommissionRecord{})
	defer testutil.CleanDB(t, db)

	// Seed wallet in DB (the transaction creates real wallet repo from db)
	wallet := &model.Wallet{
		UserID:       1,
		BalanceCents: 10000,
		FrozenCents:  15000, // 8000 + 5000 + 2000 extra
		Version:      1,
	}
	db.Create(wallet)

	mockCommissions := new(mockCommissionRepo)

	// No player repo - will use playerID as userID
	svc := NewSettlementService(mockCommissions, nil, nil)
	svc.SetTxManager(common.NewUnitOfWork(db))

	// Create test records (8 days old)
	oldDate := time.Now().AddDate(0, 0, -8)
	records := []model.CommissionRecord{
		{
			ID:                1,
			OrderID:           100,
			PlayerID:          1,
			PlayerIncomeCents: 8000,
			SettlementStatus:  model.SettlementStatusPending,
			CreatedAt:         oldDate,
		},
		{
			ID:                2,
			OrderID:           101,
			PlayerID:          1,
			PlayerIncomeCents: 5000,
			SettlementStatus:  model.SettlementStatusPending,
			CreatedAt:         oldDate,
		},
	}

	// Seed commission records in DB (the transaction creates real commission repo from db)
	for _, r := range records {
		db.Create(&r)
	}

	// Mock ListRecords (called on the injected repo, before the transaction)
	mockCommissions.On("ListRecords", ctx, mock.Anything).Return(records, int64(2), nil)

	result, err := svc.ProcessT7Unfreeze(ctx)

	require.NoError(t, err)
	assert.Equal(t, 2, result.ProcessedCount)
	assert.Equal(t, 2, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Equal(t, int64(13000), result.TotalUnfrozen)

	// Verify wallet state in DB
	var updatedWallet model.Wallet
	db.Where("user_id = ?", uint64(1)).First(&updatedWallet)
	assert.Equal(t, int64(23000), updatedWallet.BalanceCents)
	assert.Equal(t, int64(2000), updatedWallet.FrozenCents)

	mockCommissions.AssertExpectations(t)
}

func TestProcessT7Unfreeze_NoRecords(t *testing.T) {
	ctx := context.Background()
	mockCommissions := new(mockCommissionRepo)
	mockWallets := new(mockWalletRepo)

	svc := NewSettlementService(mockCommissions, mockWallets, nil)

	// No pending records
	mockCommissions.On("ListRecords", ctx, mock.Anything).Return([]model.CommissionRecord{}, int64(0), nil)

	result, err := svc.ProcessT7Unfreeze(ctx)

	require.NoError(t, err)
	assert.Equal(t, 0, result.ProcessedCount)
	assert.Equal(t, 0, result.SuccessCount)

	mockCommissions.AssertExpectations(t)
}

func TestUnfreezeByOrderID_Success(t *testing.T) {
	ctx := context.Background()

	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db, &model.Wallet{}, &model.CommissionRecord{})
	defer testutil.CleanDB(t, db)

	// Seed wallet in DB (unfreezePlayerIncomeWithTx uses real wallet repo from DB)
	wallet := &model.Wallet{
		UserID:       1,
		BalanceCents: 10000,
		FrozenCents:  10000,
		Version:      1,
	}
	db.Create(wallet)

	mockCommissions := new(mockCommissionRepo)

	record := &model.CommissionRecord{
		ID:                1,
		OrderID:           100,
		PlayerID:          1,
		PlayerIncomeCents: 8000,
		SettlementStatus:  model.SettlementStatusPending,
	}

	mockCommissions.On("GetRecordByOrderID", ctx, uint64(100)).Return(record, nil)
	mockCommissions.On("UpdateRecord", ctx, mock.MatchedBy(func(r *model.CommissionRecord) bool {
		return r.SettlementStatus == model.SettlementStatusSettled
	})).Return(nil)

	svc := NewSettlementService(mockCommissions, nil, nil)
	svc.SetTxManager(common.NewUnitOfWork(db))

	err := svc.UnfreezeByOrderID(ctx, 100)

	require.NoError(t, err)
	mockCommissions.AssertExpectations(t)

	// Verify wallet was updated in DB
	var updatedWallet model.Wallet
	db.Where("user_id = ?", uint64(1)).First(&updatedWallet)
	assert.Equal(t, int64(18000), updatedWallet.BalanceCents)
	assert.Equal(t, int64(2000), updatedWallet.FrozenCents)
}

func TestUnfreezeByOrderID_AlreadySettled(t *testing.T) {
	ctx := context.Background()
	mockCommissions := new(mockCommissionRepo)
	mockWallets := new(mockWalletRepo)

	svc := NewSettlementService(mockCommissions, mockWallets, nil)

	record := &model.CommissionRecord{
		ID:               1,
		OrderID:          100,
		SettlementStatus: model.SettlementStatusSettled,
	}

	mockCommissions.On("GetRecordByOrderID", ctx, uint64(100)).Return(record, nil)

	err := svc.UnfreezeByOrderID(ctx, 100)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "already settled")
}

func TestGetPendingSettlementStats(t *testing.T) {
	ctx := context.Background()
	mockCommissions := new(mockCommissionRepo)
	mockWallets := new(mockWalletRepo)

	svc := NewSettlementService(mockCommissions, mockWallets, nil)

	oldDate := time.Now().AddDate(0, 0, -10)
	newDate := time.Now().AddDate(0, 0, -3)

	records := []model.CommissionRecord{
		{ID: 1, PlayerIncomeCents: 8000, CreatedAt: oldDate}, // Ready
		{ID: 2, PlayerIncomeCents: 5000, CreatedAt: oldDate}, // Ready
		{ID: 3, PlayerIncomeCents: 3000, CreatedAt: newDate}, // Not ready
	}

	mockCommissions.On("ListRecords", ctx, mock.Anything).Return(records, int64(3), nil)

	stats, err := svc.GetPendingSettlementStats(ctx)

	require.NoError(t, err)
	assert.Equal(t, int64(3), stats.TotalPendingCount)
	assert.Equal(t, int64(16000), stats.TotalFrozenCents)
	assert.Equal(t, int64(2), stats.ReadyToUnfreezeCount)
	assert.Equal(t, int64(13000), stats.ReadyToUnfreezeCents)
}
