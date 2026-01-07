package wallet

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
	repoiface "gamelink/internal/repository/interfaces"
)

// ============================================================================
// Mock Implementations
// ============================================================================

// MockWalletRepository is a mock implementation of wallet repository
type MockWalletRepository struct {
	mock.Mock
}

func (m *MockWalletRepository) GetByUserID(ctx context.Context, userID uint64) (*model.Wallet, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Wallet), args.Error(1)
}

func (m *MockWalletRepository) Save(ctx context.Context, wallet *model.Wallet) error {
	args := m.Called(ctx, wallet)
	return args.Error(0)
}

func (m *MockWalletRepository) SaveWithOptimisticLock(ctx context.Context, wallet *model.Wallet) error {
	args := m.Called(ctx, wallet)
	return args.Error(0)
}

func (m *MockWalletRepository) UpdateBalanceWithLock(ctx context.Context, userID uint64, delta int64, maxRetries int) (*model.Wallet, error) {
	args := m.Called(ctx, userID, delta, maxRetries)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Wallet), args.Error(1)
}

// MockPaymentRepository is a mock implementation of payment repository
type MockPaymentRepository struct {
	mock.Mock
}

func (m *MockPaymentRepository) Create(ctx context.Context, payment *model.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func (m *MockPaymentRepository) List(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]model.Payment), args.Get(1).(int64), args.Error(2)
}

func (m *MockPaymentRepository) Get(ctx context.Context, id uint64) (*model.Payment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Payment), args.Error(1)
}

func (m *MockPaymentRepository) GetWithRelations(ctx context.Context, id uint64) (*model.Payment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Payment), args.Error(1)
}

func (m *MockPaymentRepository) Update(ctx context.Context, payment *model.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func (m *MockPaymentRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPaymentRepository) GetByOrderID(ctx context.Context, orderID uint64) ([]model.Payment, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Payment), args.Error(1)
}

func (m *MockPaymentRepository) GetByRequestID(ctx context.Context, requestID string) (*model.Payment, error) {
	args := m.Called(ctx, requestID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Payment), args.Error(1)
}

// MockOrderRepository is a mock implementation of order repository
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(ctx context.Context, order *model.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) Get(ctx context.Context, id uint64) (*model.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderRepository) Update(ctx context.Context, order *model.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) UpdateWithCondition(ctx context.Context, orderID uint64, expectedStatus model.OrderStatus, updates map[string]any) (bool, error) {
	args := m.Called(ctx, orderID, expectedStatus, updates)
	return args.Bool(0), args.Error(1)
}

func (m *MockOrderRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOrderRepository) List(ctx context.Context, opts repoiface.OrderListOptions) ([]model.Order, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]model.Order), args.Get(1).(int64), args.Error(2)
}

func (m *MockOrderRepository) GetByUser(ctx context.Context, userID uint64, opts repoiface.OrderListOptions) ([]model.Order, int64, error) {
	args := m.Called(ctx, userID, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]model.Order), args.Get(1).(int64), args.Error(2)
}

func (m *MockOrderRepository) GetByPlayer(ctx context.Context, playerID uint64, opts repoiface.OrderListOptions) ([]model.Order, int64, error) {
	args := m.Called(ctx, playerID, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]model.Order), args.Get(1).(int64), args.Error(2)
}

func (m *MockOrderRepository) GetByNo(ctx context.Context, orderNo string) (*model.Order, error) {
	args := m.Called(ctx, orderNo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

// ============================================================================
// Helper Functions
// ============================================================================

// createTestWalletService creates a test wallet service with mocked dependencies
func createTestWalletService() (*WalletService, *MockWalletRepository, *MockPaymentRepository, *MockOrderRepository) {
	mockWallets := new(MockWalletRepository)
	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)
	svc := NewWalletService(mockWallets, mockPayments, mockOrders)
	return svc, mockWallets, mockPayments, mockOrders
}

// createTestWallet creates a test wallet with given balances
func createTestWallet(userID uint64, balanceCents, frozenCents int64) *model.Wallet {
	return &model.Wallet{
		UserID:       userID,
		BalanceCents: balanceCents,
		FrozenCents:  frozenCents,
	}
}

// ============================================================================
// Recharge Tests
// ============================================================================

func TestWalletService_Recharge_Success(t *testing.T) {
	svc, mockWallets, mockPayments, mockOrders := createTestWalletService()
	ctx := context.Background()
	userID := uint64(1001)
	req := RechargeRequest{
		AmountCents: 10000, // 100.00 CNY
		Method:      model.PaymentMethodAlipay,
	}

	// Mock order creation
	mockOrders.On("Create", ctx, mock.AnythingOfType("*model.Order")).
		Run(func(args mock.Arguments) {
			order := args.Get(1).(*model.Order)
			order.ID = 1001
			assert.Equal(t, userID, order.UserID)
			assert.Equal(t, int64(10000), order.TotalPriceCents)
			assert.Equal(t, model.OrderStatusCompleted, order.Status)
			assert.NotNil(t, order.CompletedAt)
		}).
		Return(nil).Once()

	// Mock payment creation
	mockPayments.On("Create", ctx, mock.AnythingOfType("*model.Payment")).
		Run(func(args mock.Arguments) {
			payment := args.Get(1).(*model.Payment)
			payment.ID = 2001
			assert.Equal(t, uint64(1001), payment.OrderID)
			assert.Equal(t, userID, payment.UserID)
			assert.Equal(t, int64(10000), payment.AmountCents)
			assert.Equal(t, model.PaymentStatusPaid, payment.Status)
			assert.NotNil(t, payment.PaidAt)
		}).
		Return(nil).Once()

	// Mock wallet not found (auto-create)
	mockWallets.On("GetByUserID", ctx, userID).
		Return(nil, repository.ErrNotFound).Once()

	// Mock wallet save
	mockWallets.On("Save", ctx, mock.AnythingOfType("*model.Wallet")).
		Run(func(args mock.Arguments) {
			wallet := args.Get(1).(*model.Wallet)
			assert.Equal(t, userID, wallet.UserID)
			assert.Equal(t, int64(10000), wallet.BalanceCents)
			assert.Equal(t, int64(0), wallet.FrozenCents)
		}).
		Return(nil).Once()

	resp, err := svc.Recharge(ctx, userID, req)
	require.NoError(t, err)
	assert.Equal(t, uint64(1001), resp.OrderID)
	assert.Equal(t, uint64(2001), resp.PaymentID)
	assert.Equal(t, int64(10000), resp.Balance)

	mockOrders.AssertExpectations(t)
	mockPayments.AssertExpectations(t)
	mockWallets.AssertExpectations(t)
}

func TestWalletService_Recharge_ExistingWallet(t *testing.T) {
	svc, mockWallets, mockPayments, mockOrders := createTestWalletService()
	ctx := context.Background()
	userID := uint64(1001)
	req := RechargeRequest{
		AmountCents: 5000, // 50.00 CNY
		Method:      model.PaymentMethodWeChat,
	}

	existingWallet := createTestWallet(userID, 8000, 2000)

	// Mock order creation
	mockOrders.On("Create", ctx, mock.AnythingOfType("*model.Order")).
		Return(nil).Once()

	// Mock payment creation
	mockPayments.On("Create", ctx, mock.AnythingOfType("*model.Payment")).
		Return(nil).Once()

	// Mock existing wallet found
	mockWallets.On("GetByUserID", ctx, userID).
		Return(existingWallet, nil).Once()

	// Mock wallet save with updated balance
	mockWallets.On("Save", ctx, mock.AnythingOfType("*model.Wallet")).
		Run(func(args mock.Arguments) {
			wallet := args.Get(1).(*model.Wallet)
			assert.Equal(t, int64(13000), wallet.BalanceCents) // 8000 + 5000
			assert.Equal(t, int64(2000), wallet.FrozenCents)   // unchanged
		}).
		Return(nil).Once()

	resp, err := svc.Recharge(ctx, userID, req)
	require.NoError(t, err)
	assert.Equal(t, int64(13000), resp.Balance)

	mockOrders.AssertExpectations(t)
	mockPayments.AssertExpectations(t)
	mockWallets.AssertExpectations(t)
}

func TestWalletService_Recharge_InvalidAmount(t *testing.T) {
	svc, _, _, _ := createTestWalletService()
	ctx := context.Background()
	userID := uint64(1001)

	testCases := []struct {
		name        string
		amountCents int64
	}{
		{"zero amount", 0},
		{"negative amount", -100},
		{"large negative", -10000},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := RechargeRequest{
				AmountCents: tc.amountCents,
				Method:      model.PaymentMethodAlipay,
			}

			resp, err := svc.Recharge(ctx, userID, req)
			require.Error(t, err)
			assert.Equal(t, ErrInvalidAmount, err)
			assert.Nil(t, resp)
		})
	}
}

func TestWalletService_Recharge_OrderCreationFails(t *testing.T) {
	svc, mockWallets, _, mockOrders := createTestWalletService()
	ctx := context.Background()
	userID := uint64(1001)
	req := RechargeRequest{
		AmountCents: 10000,
		Method:      model.PaymentMethodAlipay,
	}

	mockOrders.On("Create", ctx, mock.AnythingOfType("*model.Order")).
		Return(errors.New("database error")).Once()

	resp, err := svc.Recharge(ctx, userID, req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.Nil(t, resp)
	mockOrders.AssertExpectations(t)
	// Wallet should not be called if order creation fails
	mockWallets.AssertNotCalled(t, "GetByUserID", ctx, userID)
}

func TestWalletService_Recharge_PaymentCreationFails(t *testing.T) {
	svc, mockWallets, mockPayments, mockOrders := createTestWalletService()
	ctx := context.Background()
	userID := uint64(1001)
	req := RechargeRequest{
		AmountCents: 10000,
		Method:      model.PaymentMethodAlipay,
	}

	mockOrders.On("Create", ctx, mock.AnythingOfType("*model.Order")).
		Return(nil).Once()

	mockPayments.On("Create", ctx, mock.AnythingOfType("*model.Payment")).
		Return(errors.New("payment gateway error")).Once()

	resp, err := svc.Recharge(ctx, userID, req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "payment gateway error")
	assert.Nil(t, resp)

	mockOrders.AssertExpectations(t)
	mockPayments.AssertExpectations(t)
	mockWallets.AssertNotCalled(t, "GetByUserID", ctx, userID)
}

func TestWalletService_Recharge_WalletGetFails(t *testing.T) {
	svc, mockWallets, mockPayments, mockOrders := createTestWalletService()
	ctx := context.Background()
	userID := uint64(1001)
	req := RechargeRequest{
		AmountCents: 10000,
		Method:      model.PaymentMethodAlipay,
	}

	mockOrders.On("Create", ctx, mock.AnythingOfType("*model.Order")).
		Return(nil).Once()

	mockPayments.On("Create", ctx, mock.AnythingOfType("*model.Payment")).
		Return(nil).Once()

	mockWallets.On("GetByUserID", ctx, userID).
		Return(nil, errors.New("database connection error")).Once()

	resp, err := svc.Recharge(ctx, userID, req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "database connection error")
	assert.Nil(t, resp)

	mockOrders.AssertExpectations(t)
	mockPayments.AssertExpectations(t)
	mockWallets.AssertExpectations(t)
}

func TestWalletService_Recharge_WalletSaveFails(t *testing.T) {
	svc, mockWallets, mockPayments, mockOrders := createTestWalletService()
	ctx := context.Background()
	userID := uint64(1001)
	req := RechargeRequest{
		AmountCents: 10000,
		Method:      model.PaymentMethodAlipay,
	}

	mockOrders.On("Create", ctx, mock.AnythingOfType("*model.Order")).
		Return(nil).Once()

	mockPayments.On("Create", ctx, mock.AnythingOfType("*model.Payment")).
		Return(nil).Once()

	mockWallets.On("GetByUserID", ctx, userID).
		Return(nil, repository.ErrNotFound).Once()

	mockWallets.On("Save", ctx, mock.AnythingOfType("*model.Wallet")).
		Return(errors.New("database write error")).Once()

	resp, err := svc.Recharge(ctx, userID, req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "database write error")
	assert.Nil(t, resp)

	mockOrders.AssertExpectations(t)
	mockPayments.AssertExpectations(t)
	mockWallets.AssertExpectations(t)
}

func TestWalletService_Recharge_LargeAmount(t *testing.T) {
	svc, mockWallets, mockPayments, mockOrders := createTestWalletService()
	ctx := context.Background()
	userID := uint64(1001)
	req := RechargeRequest{
		AmountCents: 1000000, // 10000.00 CNY
		Method:      model.PaymentMethodWallet,
	}

	mockOrders.On("Create", ctx, mock.AnythingOfType("*model.Order")).
		Return(nil).Once()

	mockPayments.On("Create", ctx, mock.AnythingOfType("*model.Payment")).
		Return(nil).Once()

	mockWallets.On("GetByUserID", ctx, userID).
		Return(nil, repository.ErrNotFound).Once()

	mockWallets.On("Save", ctx, mock.AnythingOfType("*model.Wallet")).
		Run(func(args mock.Arguments) {
			wallet := args.Get(1).(*model.Wallet)
			assert.Equal(t, int64(1000000), wallet.BalanceCents)
		}).
		Return(nil).Once()

	resp, err := svc.Recharge(ctx, userID, req)
	require.NoError(t, err)
	assert.Equal(t, int64(1000000), resp.Balance)
}

func TestWalletService_Recharge_DifferentPaymentMethods(t *testing.T) {
	paymentMethods := []model.PaymentMethod{
		model.PaymentMethodWeChat,
		model.PaymentMethodAlipay,
		model.PaymentMethodWallet,
	}

	for _, method := range paymentMethods {
		t.Run(string(method), func(t *testing.T) {
			svc, mockWallets, mockPayments, mockOrders := createTestWalletService()
			ctx := context.Background()
			userID := uint64(1001)
			req := RechargeRequest{
				AmountCents: 5000,
				Method:      method,
			}

			mockOrders.On("Create", ctx, mock.AnythingOfType("*model.Order")).
				Return(nil).Once()

			mockPayments.On("Create", ctx, mock.AnythingOfType("*model.Payment")).
				Run(func(args mock.Arguments) {
					payment := args.Get(1).(*model.Payment)
					assert.Equal(t, method, payment.Method)
				}).
				Return(nil).Once()

			mockWallets.On("GetByUserID", ctx, userID).
				Return(nil, repository.ErrNotFound).Once()

			mockWallets.On("Save", ctx, mock.AnythingOfType("*model.Wallet")).
				Return(nil).Once()

			resp, err := svc.Recharge(ctx, userID, req)
			require.NoError(t, err)
			assert.NotNil(t, resp)

			mockOrders.AssertExpectations(t)
			mockPayments.AssertExpectations(t)
			mockWallets.AssertExpectations(t)
		})
	}
}

// ============================================================================
// GetBalance Tests
// ============================================================================

func TestWalletService_GetBalance_ExistingWallet(t *testing.T) {
	svc, mockWallets, _, _ := createTestWalletService()
	ctx := context.Background()
	userID := uint64(1001)

	expectedWallet := createTestWallet(userID, 50000, 10000)

	mockWallets.On("GetByUserID", ctx, userID).
		Return(expectedWallet, nil).Once()

	wallet, err := svc.GetBalance(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, userID, wallet.UserID)
	assert.Equal(t, int64(50000), wallet.BalanceCents)
	assert.Equal(t, int64(10000), wallet.FrozenCents)

	mockWallets.AssertExpectations(t)
}

func TestWalletService_GetBalance_NoWallet_ReturnsEmpty(t *testing.T) {
	svc, mockWallets, _, _ := createTestWalletService()
	ctx := context.Background()
	userID := uint64(1001)

	mockWallets.On("GetByUserID", ctx, userID).
		Return(nil, repository.ErrNotFound).Once()

	wallet, err := svc.GetBalance(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, userID, wallet.UserID)
	assert.Equal(t, int64(0), wallet.BalanceCents)
	assert.Equal(t, int64(0), wallet.FrozenCents)

	mockWallets.AssertExpectations(t)
}

func TestWalletService_GetBalance_DatabaseError(t *testing.T) {
	svc, mockWallets, _, _ := createTestWalletService()
	ctx := context.Background()
	userID := uint64(1001)

	mockWallets.On("GetByUserID", ctx, userID).
		Return(nil, errors.New("database connection lost")).Once()

	wallet, err := svc.GetBalance(ctx, userID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "database connection lost")
	assert.Nil(t, wallet)

	mockWallets.AssertExpectations(t)
}

func TestWalletService_GetBalance_WithFrozenFunds(t *testing.T) {
	svc, mockWallets, _, _ := createTestWalletService()
	ctx := context.Background()
	userID := uint64(1001)

	// Simulate T+7 settlement scenario
	expectedWallet := createTestWallet(userID, 30000, 20000)

	mockWallets.On("GetByUserID", ctx, userID).
		Return(expectedWallet, nil).Once()

	wallet, err := svc.GetBalance(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(30000), wallet.BalanceCents) // Available for withdrawal
	assert.Equal(t, int64(20000), wallet.FrozenCents)  // Frozen pending T+7

	// Calculate total balance
	totalBalance := wallet.BalanceCents + wallet.FrozenCents
	assert.Equal(t, int64(50000), totalBalance)

	mockWallets.AssertExpectations(t)
}

func TestWalletService_GetBalance_ZeroBalances(t *testing.T) {
	svc, mockWallets, _, _ := createTestWalletService()
	ctx := context.Background()
	userID := uint64(1001)

	expectedWallet := createTestWallet(userID, 0, 0)

	mockWallets.On("GetByUserID", ctx, userID).
		Return(expectedWallet, nil).Once()

	wallet, err := svc.GetBalance(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), wallet.BalanceCents)
	assert.Equal(t, int64(0), wallet.FrozenCents)

	mockWallets.AssertExpectations(t)
}

// ============================================================================
// Balance Precision Tests (Cents)
// ============================================================================

func TestWalletService_BalancePrecision_Cents(t *testing.T) {
	svc, mockWallets, mockPayments, mockOrders := createTestWalletService()
	ctx := context.Background()
	userID := uint64(1001)

	// Test with amounts that require cents precision
	testCases := []struct {
		name           string
		amountCents    int64
		expectedString string
	}{
		{"one cent", 1, "0.01"},
		{"ten cents", 10, "0.10"},
		{"fifty cents", 50, "0.50"},
		{"ninety-nine cents", 99, "0.99"},
		{"exact yuan", 100, "1.00"},
		{"large with cents", 10099, "100.99"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := RechargeRequest{
				AmountCents: tc.amountCents,
				Method:      model.PaymentMethodAlipay,
			}

			mockOrders.On("Create", ctx, mock.AnythingOfType("*model.Order")).
				Return(nil).Once()
			mockPayments.On("Create", ctx, mock.AnythingOfType("*model.Payment")).
				Return(nil).Once()
			mockWallets.On("GetByUserID", ctx, userID).
				Return(nil, repository.ErrNotFound).Once()
			mockWallets.On("Save", ctx, mock.AnythingOfType("*model.Wallet")).
				Run(func(args mock.Arguments) {
					wallet := args.Get(1).(*model.Wallet)
					assert.Equal(t, tc.amountCents, wallet.BalanceCents)
				}).
				Return(nil).Once()

			resp, err := svc.Recharge(ctx, userID, req)
			require.NoError(t, err)
			assert.Equal(t, tc.amountCents, resp.Balance)

			mockOrders.AssertExpectations(t)
			mockPayments.AssertExpectations(t)
			mockWallets.AssertExpectations(t)
		})
	}
}

// ============================================================================
// Constructor Tests
// ============================================================================

func TestNewWalletService(t *testing.T) {
	mockWallets := new(MockWalletRepository)
	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)

	svc := NewWalletService(mockWallets, mockPayments, mockOrders)

	require.NotNil(t, svc)
	assert.Equal(t, mockWallets, svc.wallets)
	assert.Equal(t, mockPayments, svc.payments)
	assert.Equal(t, mockOrders, svc.orders)
}

// ============================================================================
// T+7 Settlement Related Tests (Frozen Balance)
// ============================================================================

func TestWalletService_Recharge_WithFrozenFunds(t *testing.T) {
	svc, mockWallets, mockPayments, mockOrders := createTestWalletService()
	ctx := context.Background()
	userID := uint64(1001)
	req := RechargeRequest{
		AmountCents: 10000,
		Method:      model.PaymentMethodAlipay,
	}

	// Existing wallet with frozen funds from T+7 settlement
	existingWallet := createTestWallet(userID, 5000, 15000)

	mockOrders.On("Create", ctx, mock.AnythingOfType("*model.Order")).
		Return(nil).Once()

	mockPayments.On("Create", ctx, mock.AnythingOfType("*model.Payment")).
		Return(nil).Once()

	mockWallets.On("GetByUserID", ctx, userID).
		Return(existingWallet, nil).Once()

	mockWallets.On("Save", ctx, mock.AnythingOfType("*model.Wallet")).
		Run(func(args mock.Arguments) {
			wallet := args.Get(1).(*model.Wallet)
			// Only balance_cents should increase, frozen_cents unchanged
			assert.Equal(t, int64(15000), wallet.BalanceCents) // 5000 + 10000
			assert.Equal(t, int64(15000), wallet.FrozenCents)  // unchanged
		}).
		Return(nil).Once()

	resp, err := svc.Recharge(ctx, userID, req)
	require.NoError(t, err)
	assert.Equal(t, int64(15000), resp.Balance)

	mockOrders.AssertExpectations(t)
	mockPayments.AssertExpectations(t)
	mockWallets.AssertExpectations(t)
}

func TestWalletService_GetBalance_WithdrawableCalculation(t *testing.T) {
	svc, mockWallets, _, _ := createTestWalletService()
	ctx := context.Background()
	userID := uint64(1001)

	// Scenario: Player has completed orders but some are in T+7 freeze
	expectedWallet := createTestWallet(userID, 8000, 12000)

	mockWallets.On("GetByUserID", ctx, userID).
		Return(expectedWallet, nil).Once()

	wallet, err := svc.GetBalance(ctx, userID)
	require.NoError(t, err)

	// Only balance_cents is withdrawable
	withdrawable := wallet.BalanceCents
	pendingSettlement := wallet.FrozenCents

	assert.Equal(t, int64(8000), withdrawable)
	assert.Equal(t, int64(12000), pendingSettlement)

	// Verify player can only withdraw the non-frozen amount
	totalWithdrawable := withdrawable
	assert.Equal(t, int64(8000), totalWithdrawable)

	mockWallets.AssertExpectations(t)
}

// ============================================================================
// Edge Cases and Error Handling
// ============================================================================

func TestWalletService_Recharge_MultipleRecharges(t *testing.T) {
	svc, mockWallets, mockPayments, mockOrders := createTestWalletService()
	ctx := context.Background()
	userID := uint64(1001)

	// First recharge
	req1 := RechargeRequest{
		AmountCents: 5000,
		Method:      model.PaymentMethodAlipay,
	}

	mockOrders.On("Create", ctx, mock.AnythingOfType("*model.Order")).
		Return(nil).Once()
	mockPayments.On("Create", ctx, mock.AnythingOfType("*model.Payment")).
		Return(nil).Once()
	mockWallets.On("GetByUserID", ctx, userID).
		Return(nil, repository.ErrNotFound).Once()
	mockWallets.On("Save", ctx, mock.AnythingOfType("*model.Wallet")).
		Return(nil).Once()

	resp1, err := svc.Recharge(ctx, userID, req1)
	require.NoError(t, err)
	assert.Equal(t, int64(5000), resp1.Balance)

	// Second recharge (wallet now exists)
	req2 := RechargeRequest{
		AmountCents: 3000,
		Method:      model.PaymentMethodWeChat,
	}

	updatedWallet := createTestWallet(userID, 5000, 0)
	mockOrders.On("Create", ctx, mock.AnythingOfType("*model.Order")).
		Return(nil).Once()
	mockPayments.On("Create", ctx, mock.AnythingOfType("*model.Payment")).
		Return(nil).Once()
	mockWallets.On("GetByUserID", ctx, userID).
		Return(updatedWallet, nil).Once()
	mockWallets.On("Save", ctx, mock.AnythingOfType("*model.Wallet")).
		Return(nil).Once()

	resp2, err := svc.Recharge(ctx, userID, req2)
	require.NoError(t, err)
	assert.Equal(t, int64(8000), resp2.Balance)

	mockOrders.AssertExpectations(t)
	mockPayments.AssertExpectations(t)
	mockWallets.AssertExpectations(t)
}

func TestWalletService_GetBalance_ContextCancellation(t *testing.T) {
	svc, mockWallets, _, _ := createTestWalletService()

	// Create a canceled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	userID := uint64(1001)

	mockWallets.On("GetByUserID", ctx, userID).
		Return(nil, context.Canceled).Once()

	_, err := svc.GetBalance(ctx, userID)
	require.Error(t, err)
	assert.Equal(t, context.Canceled, err)

	mockWallets.AssertExpectations(t)
}

func TestWalletService_Recharge_OrderTimestamps(t *testing.T) {
	svc, mockWallets, mockPayments, mockOrders := createTestWalletService()
	ctx := context.Background()
	userID := uint64(1001)
	req := RechargeRequest{
		AmountCents: 10000,
		Method:      model.PaymentMethodAlipay,
	}

	beforeTest := time.Now().UTC()

	mockOrders.On("Create", ctx, mock.AnythingOfType("*model.Order")).
		Run(func(args mock.Arguments) {
			order := args.Get(1).(*model.Order)
			assert.NotNil(t, order.CompletedAt)
			assert.True(t, order.CompletedAt.After(beforeTest) || order.CompletedAt.Equal(beforeTest))
		}).
		Return(nil).Once()

	mockPayments.On("Create", ctx, mock.AnythingOfType("*model.Payment")).
		Run(func(args mock.Arguments) {
			payment := args.Get(1).(*model.Payment)
			assert.NotNil(t, payment.PaidAt)
			assert.True(t, payment.PaidAt.After(beforeTest) || payment.PaidAt.Equal(beforeTest))
		}).
		Return(nil).Once()

	mockWallets.On("GetByUserID", ctx, userID).
		Return(nil, repository.ErrNotFound).Once()
	mockWallets.On("Save", ctx, mock.AnythingOfType("*model.Wallet")).
		Return(nil).Once()

	resp, err := svc.Recharge(ctx, userID, req)
	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestWalletService_Recharge_OrderFields(t *testing.T) {
	svc, mockWallets, mockPayments, mockOrders := createTestWalletService()
	ctx := context.Background()
	userID := uint64(1001)
	req := RechargeRequest{
		AmountCents: 10000,
		Method:      model.PaymentMethodAlipay,
	}

	mockOrders.On("Create", ctx, mock.AnythingOfType("*model.Order")).
		Run(func(args mock.Arguments) {
			order := args.Get(1).(*model.Order)
			assert.Equal(t, userID, order.UserID)
			assert.Equal(t, uint64(0), order.ItemID) // No service item for recharge
			assert.Equal(t, "Wallet Recharge", order.Title)
			assert.Equal(t, int64(10000), order.UnitPriceCents)
			assert.Equal(t, int64(10000), order.TotalPriceCents)
			assert.Equal(t, model.OrderStatusCompleted, order.Status)
		}).
		Return(nil).Once()

	mockPayments.On("Create", ctx, mock.AnythingOfType("*model.Payment")).
		Run(func(args mock.Arguments) {
			payment := args.Get(1).(*model.Payment)
			assert.Equal(t, model.PaymentMethodAlipay, payment.Method)
			assert.Equal(t, model.CurrencyCNY, payment.Currency)
			assert.Equal(t, model.PaymentStatusPaid, payment.Status)
		}).
		Return(nil).Once()

	mockWallets.On("GetByUserID", ctx, userID).
		Return(nil, repository.ErrNotFound).Once()
	mockWallets.On("Save", ctx, mock.AnythingOfType("*model.Wallet")).
		Return(nil).Once()

	resp, err := svc.Recharge(ctx, userID, req)
	require.NoError(t, err)
	assert.NotNil(t, resp)
}

// ============================================================================
// Currency Tests
// ============================================================================

func TestWalletService_Recharge_CurrencyDefault(t *testing.T) {
	svc, mockWallets, mockPayments, mockOrders := createTestWalletService()
	ctx := context.Background()
	userID := uint64(1001)
	req := RechargeRequest{
		AmountCents: 10000,
		Method:      model.PaymentMethodAlipay,
	}

	mockOrders.On("Create", ctx, mock.AnythingOfType("*model.Order")).
		Return(nil).Once()

	mockPayments.On("Create", ctx, mock.AnythingOfType("*model.Payment")).
		Run(func(args mock.Arguments) {
			payment := args.Get(1).(*model.Payment)
			// Verify default currency is CNY
			assert.Equal(t, model.CurrencyCNY, payment.Currency)
		}).
		Return(nil).Once()

	mockWallets.On("GetByUserID", ctx, userID).
		Return(nil, repository.ErrNotFound).Once()
	mockWallets.On("Save", ctx, mock.AnythingOfType("*model.Wallet")).
		Return(nil).Once()

	resp, err := svc.Recharge(ctx, userID, req)
	require.NoError(t, err)
	assert.NotNil(t, resp)
}
