package payment

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

// MockPaymentRepository is a mock implementation of PaymentRepository
type MockPaymentRepository struct {
	mock.Mock
}

func (m *MockPaymentRepository) Create(ctx context.Context, payment *model.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func (m *MockPaymentRepository) List(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
	args := m.Called(ctx, opts)
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
	return args.Get(0).([]model.Payment), args.Error(1)
}

func (m *MockPaymentRepository) GetByRequestID(ctx context.Context, requestID string) (*model.Payment, error) {
	args := m.Called(ctx, requestID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Payment), args.Error(1)
}

// MockOrderRepository is a mock implementation of OrderRepository
type MockOrderRepository struct {
	mock.Mock
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

// Implement other required methods from OrderReadWriter interface
func (m *MockOrderRepository) Create(ctx context.Context, order *model.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOrderRepository) List(ctx context.Context, opts repoiface.OrderListOptions) ([]model.Order, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.Order), args.Get(1).(int64), args.Error(2)
}

func (m *MockOrderRepository) GetByOrderNo(ctx context.Context, orderNo string) (*model.Order, error) {
	args := m.Called(ctx, orderNo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderRepository) GetByUserID(ctx context.Context, userID uint64, opts repoiface.OrderListOptions) ([]model.Order, int64, error) {
	args := m.Called(ctx, userID, opts)
	return args.Get(0).([]model.Order), args.Get(1).(int64), args.Error(2)
}

func (m *MockOrderRepository) UpdateStatus(ctx context.Context, id uint64, status model.OrderStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockOrderRepository) CountByStatus(ctx context.Context, statuses []model.OrderStatus) (map[model.OrderStatus]int64, error) {
	args := m.Called(ctx, statuses)
	return args.Get(0).(map[model.OrderStatus]int64), args.Error(1)
}

func (m *MockOrderRepository) GetPlayerOrderStats(ctx context.Context, playerID uint64, startDate, endDate time.Time) (map[string]int64, error) {
	args := m.Called(ctx, playerID, startDate, endDate)
	return args.Get(0).(map[string]int64), args.Error(1)
}

func (m *MockOrderRepository) GetUserOrderStats(ctx context.Context, userID uint64, startDate, endDate time.Time) (map[string]int64, error) {
	args := m.Called(ctx, userID, startDate, endDate)
	return args.Get(0).(map[string]int64), args.Error(1)
}

func (m *MockOrderRepository) UpdateWithCondition(ctx context.Context, orderID uint64, expectedStatus model.OrderStatus, updates map[string]any) (bool, error) {
	args := m.Called(ctx, orderID, expectedStatus, updates)
	return args.Bool(0), args.Error(1)
}

// MockWalletRepository is a mock implementation of WalletRepository
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

// Helper function to create test order
func createTestOrder(id, userID uint64, status model.OrderStatus, priceCents int64) *model.Order {
	return &model.Order{
		Base:            model.Base{ID: id},
		OrderNo:         "ORD123456",
		UserID:          userID,
		ItemID:          1,
		Status:          status,
		TotalPriceCents: priceCents,
		Currency:        model.CurrencyCNY,
	}
}

// Helper function to create test payment
func createTestPayment(id, orderID, userID uint64, status model.PaymentStatus, amountCents int64, method model.PaymentMethod) *model.Payment {
	return &model.Payment{
		Base:        model.Base{ID: id},
		OrderID:     orderID,
		UserID:      userID,
		Status:      status,
		AmountCents: amountCents,
		Method:      method,
		Currency:    model.CurrencyCNY,
	}
}

// Helper function to create test wallet
func createTestWallet(userID uint64, balanceCents, frozenCents int64) *model.Wallet {
	return &model.Wallet{
		Base:         model.Base{ID: 1},
		UserID:       userID,
		BalanceCents: balanceCents,
		FrozenCents:  frozenCents,
	}
}

// TestPaymentService_CreatePayment_ThirdParty_Success tests successful third-party payment creation
func TestPaymentService_CreatePayment_ThirdParty_Success(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	orderID := uint64(1000)

	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)

	order := createTestOrder(orderID, userID, model.OrderStatusPending, 10000)

	mockOrders.On("Get", ctx, orderID).Return(order, nil)
	mockPayments.On("List", ctx, mock.AnythingOfType("repository.PaymentListOptions")).Return([]model.Payment{}, int64(0), nil)
	mockPayments.On("Create", ctx, mock.AnythingOfType("*model.Payment")).Return(nil).Run(func(args mock.Arguments) {
		payment := args.Get(1).(*model.Payment)
		payment.ID = 1
	})
	mockPayments.On("Get", ctx, uint64(1)).Return(createTestPayment(1, orderID, userID, model.PaymentStatusPending, 10000, model.PaymentMethodWeChat), nil) // For mockPaymentSuccess
	mockPayments.On("Update", ctx, mock.AnythingOfType("*model.Payment")).Return(nil)
	mockOrders.On("Update", ctx, mock.AnythingOfType("*model.Order")).Return(nil)

	service := NewPaymentService(mockPayments, mockOrders)

	req := CreatePaymentRequest{
		OrderID: orderID,
		Method:  model.PaymentMethodWeChat,
	}

	resp, err := service.CreatePayment(ctx, userID, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Greater(t, resp.PaymentID, uint64(0))
	assert.NotNil(t, resp.PayInfo)
	assert.Equal(t, int64(10000), resp.ThirdPartyAmount)

	mockOrders.AssertExpectations(t)
	mockPayments.AssertExpectations(t)
}

// TestPaymentService_CreatePayment_InvalidOrderStatus tests payment creation with invalid order status
func TestPaymentService_CreatePayment_InvalidOrderStatus(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	orderID := uint64(1000)

	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)

	// Order already confirmed
	order := createTestOrder(orderID, userID, model.OrderStatusConfirmed, 10000)

	mockOrders.On("Get", ctx, orderID).Return(order, nil)

	service := NewPaymentService(mockPayments, mockOrders)

	req := CreatePaymentRequest{
		OrderID: orderID,
		Method:  model.PaymentMethodWeChat,
	}

	resp, err := service.CreatePayment(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, ErrInvalidOrderStatus, err)

	mockOrders.AssertExpectations(t)
}

// TestPaymentService_CreatePayment_UnauthorizedUser tests payment creation by unauthorized user
func TestPaymentService_CreatePayment_UnauthorizedUser(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	orderID := uint64(1000)
	differentUserID := uint64(999)

	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)

	// Order belongs to different user
	order := createTestOrder(orderID, differentUserID, model.OrderStatusPending, 10000)

	mockOrders.On("Get", ctx, orderID).Return(order, nil)

	service := NewPaymentService(mockPayments, mockOrders)

	req := CreatePaymentRequest{
		OrderID: orderID,
		Method:  model.PaymentMethodWeChat,
	}

	resp, err := service.CreatePayment(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "unauthorized")

	mockOrders.AssertExpectations(t)
}

// TestPaymentService_CreatePayment_OrderAlreadyPaid tests payment creation for already paid order
func TestPaymentService_CreatePayment_OrderAlreadyPaid(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	orderID := uint64(1000)

	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)

	order := createTestOrder(orderID, userID, model.OrderStatusPending, 10000)
	existingPayment := createTestPayment(1, orderID, userID, model.PaymentStatusPaid, 10000, model.PaymentMethodWeChat)

	mockOrders.On("Get", ctx, orderID).Return(order, nil)
	mockPayments.On("List", ctx, mock.AnythingOfType("repository.PaymentListOptions")).Return(
		[]model.Payment{*existingPayment}, int64(1), nil)

	service := NewPaymentService(mockPayments, mockOrders)

	req := CreatePaymentRequest{
		OrderID: orderID,
		Method:  model.PaymentMethodWeChat,
	}

	resp, err := service.CreatePayment(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, ErrOrderAlreadyPaid, err)

	mockOrders.AssertExpectations(t)
	mockPayments.AssertExpectations(t)
}

// TestPaymentService_CreatePayment_Wallet_Success tests successful wallet payment
func TestPaymentService_CreatePayment_Wallet_Success(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	orderID := uint64(1000)
	priceCents := int64(10000)

	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)
	mockWallets := new(MockWalletRepository)

	order := createTestOrder(orderID, userID, model.OrderStatusPending, priceCents)
	wallet := createTestWallet(userID, 20000, 0)
	updatedWallet := createTestWallet(userID, 10000, 0) // After deduction

	mockOrders.On("Get", ctx, orderID).Return(order, nil)
	mockPayments.On("List", ctx, mock.AnythingOfType("repository.PaymentListOptions")).Return([]model.Payment{}, int64(0), nil)
	mockWallets.On("GetByUserID", ctx, userID).Return(wallet, nil)
	mockWallets.On("UpdateBalanceWithLock", ctx, userID, int64(-10000), 3).Return(updatedWallet, nil)
	mockPayments.On("Create", ctx, mock.AnythingOfType("*model.Payment")).Return(nil).Run(func(args mock.Arguments) {
		payment := args.Get(1).(*model.Payment)
		payment.ID = 1
	})
	mockOrders.On("Update", ctx, mock.MatchedBy(func(o *model.Order) bool {
		return o.Status == model.OrderStatusConfirmed
	})).Return(nil)

	service := NewPaymentService(mockPayments, mockOrders)
	service.SetWalletRepository(mockWallets)

	req := CreatePaymentRequest{
		OrderID: orderID,
		Method:  model.PaymentMethodWallet,
	}

	resp, err := service.CreatePayment(ctx, userID, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, uint64(1), resp.PaymentID)
	assert.True(t, resp.WalletPaidDirect)
	assert.Equal(t, priceCents, resp.WalletDeducted)

	mockOrders.AssertExpectations(t)
	mockPayments.AssertExpectations(t)
	mockWallets.AssertExpectations(t)
}

// TestPaymentService_CreatePayment_Wallet_InsufficientBalance tests wallet payment with insufficient balance
func TestPaymentService_CreatePayment_Wallet_InsufficientBalance(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	orderID := uint64(1000)
	priceCents := int64(10000)

	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)
	mockWallets := new(MockWalletRepository)

	order := createTestOrder(orderID, userID, model.OrderStatusPending, priceCents)
	wallet := createTestWallet(userID, 5000, 0) // Only 5000 available

	mockOrders.On("Get", ctx, orderID).Return(order, nil)
	mockPayments.On("List", ctx, mock.AnythingOfType("repository.PaymentListOptions")).Return([]model.Payment{}, int64(0), nil)
	mockWallets.On("GetByUserID", ctx, userID).Return(wallet, nil)

	service := NewPaymentService(mockPayments, mockOrders)
	service.SetWalletRepository(mockWallets)

	req := CreatePaymentRequest{
		OrderID: orderID,
		Method:  model.PaymentMethodWallet,
	}

	resp, err := service.CreatePayment(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "insufficient wallet balance")

	mockOrders.AssertExpectations(t)
	mockPayments.AssertExpectations(t)
	mockWallets.AssertExpectations(t)
}

// TestPaymentService_CreatePayment_Wallet_NoWallet tests wallet payment when user has no wallet
func TestPaymentService_CreatePayment_Wallet_NoWallet(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	orderID := uint64(1000)
	priceCents := int64(10000)

	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)
	mockWallets := new(MockWalletRepository)

	order := createTestOrder(orderID, userID, model.OrderStatusPending, priceCents)

	mockOrders.On("Get", ctx, orderID).Return(order, nil)
	mockPayments.On("List", ctx, mock.AnythingOfType("repository.PaymentListOptions")).Return([]model.Payment{}, int64(0), nil)
	mockWallets.On("GetByUserID", ctx, userID).Return(nil, repository.ErrNotFound)

	service := NewPaymentService(mockPayments, mockOrders)
	service.SetWalletRepository(mockWallets)

	req := CreatePaymentRequest{
		OrderID: orderID,
		Method:  model.PaymentMethodWallet,
	}

	resp, err := service.CreatePayment(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "insufficient wallet balance")

	mockOrders.AssertExpectations(t)
	mockPayments.AssertExpectations(t)
	mockWallets.AssertExpectations(t)
}

// TestPaymentService_CreatePayment_Combined_Success tests successful combined payment
func TestPaymentService_CreatePayment_Combined_Success(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	orderID := uint64(1000)
	totalCents := int64(10000)
	walletCents := int64(6000)
	thirdPartyCents := int64(4000)

	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)
	mockWallets := new(MockWalletRepository)

	order := createTestOrder(orderID, userID, model.OrderStatusPending, totalCents)
	wallet := createTestWallet(userID, 10000, 0)
	updatedWallet := createTestWallet(userID, 4000, 0) // After deduction

	mockOrders.On("Get", ctx, orderID).Return(order, nil)
	mockPayments.On("List", ctx, mock.AnythingOfType("repository.PaymentListOptions")).Return([]model.Payment{}, int64(0), nil)
	mockWallets.On("GetByUserID", ctx, userID).Return(wallet, nil)
	mockWallets.On("UpdateBalanceWithLock", ctx, userID, int64(-6000), 3).Return(updatedWallet, nil)
	mockPayments.On("Create", ctx, mock.AnythingOfType("*model.Payment")).Return(nil).Run(func(args mock.Arguments) {
		payment := args.Get(1).(*model.Payment)
		payment.ID = 1
	})
	mockPayments.On("Get", ctx, uint64(1)).Return(createTestPayment(1, orderID, userID, model.PaymentStatusPending, totalCents, model.PaymentMethodCombined), nil) // For mockPaymentSuccess
	mockPayments.On("Update", ctx, mock.AnythingOfType("*model.Payment")).Return(nil)
	mockOrders.On("Update", ctx, mock.AnythingOfType("*model.Order")).Return(nil)

	service := NewPaymentService(mockPayments, mockOrders)
	service.SetWalletRepository(mockWallets)

	req := CreatePaymentRequest{
		OrderID:           orderID,
		Method:            model.PaymentMethodCombined,
		WalletAmountCents: walletCents,
		ThirdPartyMethod:  model.PaymentMethodWeChat,
	}

	resp, err := service.CreatePayment(ctx, userID, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, walletCents, resp.WalletDeducted)
	assert.Equal(t, thirdPartyCents, resp.ThirdPartyAmount)

	mockOrders.AssertExpectations(t)
	mockPayments.AssertExpectations(t)
	mockWallets.AssertExpectations(t)
}

// TestPaymentService_CreatePayment_Combined_InvalidWalletAmount tests combined payment with invalid wallet amount
func TestPaymentService_CreatePayment_Combined_InvalidWalletAmount(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	orderID := uint64(1000)

	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)
	mockWallets := new(MockWalletRepository)

	order := createTestOrder(orderID, userID, model.OrderStatusPending, 10000)

	mockOrders.On("Get", ctx, orderID).Return(order, nil)
	mockPayments.On("List", ctx, mock.AnythingOfType("repository.PaymentListOptions")).Return([]model.Payment{}, int64(0), nil)

	service := NewPaymentService(mockPayments, mockOrders)
	service.SetWalletRepository(mockWallets)

	tests := []struct {
		name              string
		walletAmountCents int64
		expectError       string
	}{
		{
			name:              "Zero wallet amount",
			walletAmountCents: 0,
			expectError:       "wallet amount must be positive",
		},
		{
			name:              "Wallet amount equals or exceeds total",
			walletAmountCents: 10000,
			expectError:       "wallet amount must be less than total price",
		},
		{
			name:              "Wallet amount exceeds total",
			walletAmountCents: 15000,
			expectError:       "wallet amount must be less than total price",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := CreatePaymentRequest{
				OrderID:           orderID,
				Method:            model.PaymentMethodCombined,
				WalletAmountCents: tt.walletAmountCents,
				ThirdPartyMethod:  model.PaymentMethodWeChat,
			}

			resp, err := service.CreatePayment(ctx, userID, req)

			assert.Error(t, err)
			assert.Nil(t, resp)
			assert.Contains(t, err.Error(), tt.expectError)
		})
	}

	mockOrders.AssertExpectations(t)
	mockPayments.AssertExpectations(t)
}

// TestPaymentService_GetPaymentStatus_Success tests successful payment status retrieval
func TestPaymentService_GetPaymentStatus_Success(t *testing.T) {
	ctx := context.Background()
	paymentID := uint64(1)

	mockPayments := new(MockPaymentRepository)

	now := time.Now()
	payment := createTestPayment(paymentID, 1000, 100, model.PaymentStatusPaid, 10000, model.PaymentMethodWeChat)
	payment.PaidAt = &now

	mockPayments.On("Get", ctx, paymentID).Return(payment, nil)

	service := NewPaymentService(mockPayments, nil)

	resp, err := service.GetPaymentStatus(ctx, paymentID)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, paymentID, resp.PaymentID)
	assert.Equal(t, uint64(1000), resp.OrderID)
	assert.Equal(t, model.PaymentStatusPaid, resp.Status)
	assert.NotNil(t, resp.PaidAt)

	mockPayments.AssertExpectations(t)
}

// TestPaymentService_GetPaymentStatus_NotFound tests payment status retrieval for non-existent payment
func TestPaymentService_GetPaymentStatus_NotFound(t *testing.T) {
	ctx := context.Background()
	paymentID := uint64(999)

	mockPayments := new(MockPaymentRepository)

	mockPayments.On("Get", ctx, paymentID).Return(nil, repository.ErrNotFound)

	service := NewPaymentService(mockPayments, nil)

	resp, err := service.GetPaymentStatus(ctx, paymentID)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, ErrNotFound, err)

	mockPayments.AssertExpectations(t)
}

// TestPaymentService_CancelPayment_Success tests successful payment cancellation
func TestPaymentService_CancelPayment_Success(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	paymentID := uint64(1)

	mockPayments := new(MockPaymentRepository)

	payment := createTestPayment(paymentID, 1000, userID, model.PaymentStatusPending, 10000, model.PaymentMethodWeChat)

	mockPayments.On("Get", ctx, paymentID).Return(payment, nil)
	mockPayments.On("Update", ctx, mock.MatchedBy(func(p *model.Payment) bool {
		return p.Status == model.PaymentStatusFailed
	})).Return(nil)

	service := NewPaymentService(mockPayments, nil)

	err := service.CancelPayment(ctx, userID, paymentID)

	assert.NoError(t, err)

	mockPayments.AssertExpectations(t)
}

// TestPaymentService_CancelPayment_Unauthorized tests payment cancellation by unauthorized user
func TestPaymentService_CancelPayment_Unauthorized(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	differentUserID := uint64(999)
	paymentID := uint64(1)

	mockPayments := new(MockPaymentRepository)

	payment := createTestPayment(paymentID, 1000, differentUserID, model.PaymentStatusPending, 10000, model.PaymentMethodWeChat)

	mockPayments.On("Get", ctx, paymentID).Return(payment, nil)

	service := NewPaymentService(mockPayments, nil)

	err := service.CancelPayment(ctx, userID, paymentID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")

	mockPayments.AssertExpectations(t)
}

// TestPaymentService_CancelPayment_InvalidStatus tests payment cancellation with invalid status
func TestPaymentService_CancelPayment_InvalidStatus(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	paymentID := uint64(1)

	mockPayments := new(MockPaymentRepository)

	payment := createTestPayment(paymentID, 1000, userID, model.PaymentStatusPaid, 10000, model.PaymentMethodWeChat)

	mockPayments.On("Get", ctx, paymentID).Return(payment, nil)

	service := NewPaymentService(mockPayments, nil)

	err := service.CancelPayment(ctx, userID, paymentID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot cancel payment")

	mockPayments.AssertExpectations(t)
}

// TestPaymentService_RefundPayment_Wallet_Success tests successful wallet payment refund
func TestPaymentService_RefundPayment_Wallet_Success(t *testing.T) {
	ctx := context.Background()
	paymentID := uint64(1)
	orderID := uint64(1000)
	userID := uint64(100)
	amountCents := int64(10000)

	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)
	mockWallets := new(MockWalletRepository)

	payment := createTestPayment(paymentID, orderID, userID, model.PaymentStatusPaid, amountCents, model.PaymentMethodWallet)
	payment.WalletAmountCents = amountCents

	order := createTestOrder(orderID, userID, model.OrderStatusConfirmed, amountCents)
	wallet := createTestWallet(userID, 0, 0)

	mockPayments.On("Get", ctx, paymentID).Return(payment, nil)
	mockWallets.On("GetByUserID", ctx, userID).Return(wallet, nil)
	mockWallets.On("Save", ctx, mock.MatchedBy(func(w *model.Wallet) bool {
		return w.BalanceCents == amountCents
	})).Return(nil)
	mockPayments.On("Update", ctx, mock.MatchedBy(func(p *model.Payment) bool {
		return p.Status == model.PaymentStatusRefunded
	})).Return(nil)
	mockOrders.On("Get", ctx, orderID).Return(order, nil)
	mockOrders.On("Update", ctx, mock.MatchedBy(func(o *model.Order) bool {
		return o.Status == model.OrderStatusRefunded
	})).Return(nil)

	service := NewPaymentService(mockPayments, mockOrders)
	service.SetWalletRepository(mockWallets)

	err := service.RefundPayment(ctx, paymentID, "customer request")

	assert.NoError(t, err)

	mockPayments.AssertExpectations(t)
	mockOrders.AssertExpectations(t)
	mockWallets.AssertExpectations(t)
}

// TestPaymentService_RefundPayment_InvalidStatus tests refund with invalid payment status
func TestPaymentService_RefundPayment_InvalidStatus(t *testing.T) {
	ctx := context.Background()
	paymentID := uint64(1)

	mockPayments := new(MockPaymentRepository)

	payment := createTestPayment(paymentID, 1000, 100, model.PaymentStatusPending, 10000, model.PaymentMethodWeChat)

	mockPayments.On("Get", ctx, paymentID).Return(payment, nil)

	service := NewPaymentService(mockPayments, nil)

	err := service.RefundPayment(ctx, paymentID, "test refund")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "payment status must be paid")

	mockPayments.AssertExpectations(t)
}

// TestPaymentService_List_Success tests successful payment list retrieval
func TestPaymentService_List_Success(t *testing.T) {
	ctx := context.Background()

	mockPayments := new(MockPaymentRepository)

	payments := []model.Payment{
		*createTestPayment(1, 1000, 100, model.PaymentStatusPaid, 10000, model.PaymentMethodWeChat),
		*createTestPayment(2, 1001, 100, model.PaymentStatusPending, 5000, model.PaymentMethodAlipay),
	}

	mockPayments.On("List", ctx, mock.MatchedBy(func(opts repository.PaymentListOptions) bool {
		return opts.Page == 1 && opts.PageSize == 20
	})).Return(payments, int64(2), nil)

	service := NewPaymentService(mockPayments, nil)

	result, total, err := service.List(ctx, repository.PaymentListOptions{
		Page:     1,
		PageSize: 20,
	})

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)

	mockPayments.AssertExpectations(t)
}

// TestPaymentService_List_DefaultPagination tests default pagination values
func TestPaymentService_List_DefaultPagination(t *testing.T) {
	ctx := context.Background()

	mockPayments := new(MockPaymentRepository)

	mockPayments.On("List", ctx, mock.MatchedBy(func(opts repository.PaymentListOptions) bool {
		return opts.Page == 1 && opts.PageSize == 20
	})).Return([]model.Payment{}, int64(0), nil)

	service := NewPaymentService(mockPayments, nil)

	// Test with page 0 (should default to 1)
	_, _, err := service.List(ctx, repository.PaymentListOptions{
		Page: 0,
	})

	require.NoError(t, err)

	// Test with page size > 100 (should default to 20)
	mockPayments.On("List", ctx, mock.MatchedBy(func(opts repository.PaymentListOptions) bool {
		return opts.Page == 1 && opts.PageSize == 20
	})).Return([]model.Payment{}, int64(0), nil)

	_, _, err = service.List(ctx, repository.PaymentListOptions{
		Page:     1,
		PageSize: 150,
	})

	require.NoError(t, err)

	mockPayments.AssertExpectations(t)
}

// TestPaymentService_GetWalletBalance_Success tests successful wallet balance retrieval
func TestPaymentService_GetWalletBalance_Success(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)

	mockWallets := new(MockWalletRepository)

	wallet := createTestWallet(userID, 50000, 10000)

	mockWallets.On("GetByUserID", ctx, userID).Return(wallet, nil)

	service := NewPaymentService(nil, nil)
	service.SetWalletRepository(mockWallets)

	resp, err := service.GetWalletBalance(ctx, userID)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, int64(50000), resp.BalanceCents)
	assert.Equal(t, int64(10000), resp.FrozenCents)

	mockWallets.AssertExpectations(t)
}

// TestPaymentService_GetWalletBalance_NoWallet tests wallet balance when user has no wallet
func TestPaymentService_GetWalletBalance_NoWallet(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)

	mockWallets := new(MockWalletRepository)

	mockWallets.On("GetByUserID", ctx, userID).Return(nil, repository.ErrNotFound)

	service := NewPaymentService(nil, nil)
	service.SetWalletRepository(mockWallets)

	resp, err := service.GetWalletBalance(ctx, userID)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, int64(0), resp.BalanceCents)
	assert.Equal(t, int64(0), resp.FrozenCents)

	mockWallets.AssertExpectations(t)
}

// TestPaymentService_GetWalletBalance_NoRepository tests wallet balance without wallet repository
func TestPaymentService_GetWalletBalance_NoRepository(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)

	service := NewPaymentService(nil, nil)

	resp, err := service.GetWalletBalance(ctx, userID)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, int64(0), resp.BalanceCents)
	assert.Equal(t, int64(0), resp.FrozenCents)
}

// TestPaymentService_CalculateCombinedPayment_Success tests successful combined payment calculation
func TestPaymentService_CalculateCombinedPayment_Success(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	orderID := uint64(1000)
	totalCents := int64(10000)
	walletBalance := int64(15000)

	mockOrders := new(MockOrderRepository)
	mockWallets := new(MockWalletRepository)

	order := createTestOrder(orderID, userID, model.OrderStatusPending, totalCents)
	wallet := createTestWallet(userID, walletBalance, 0)

	mockOrders.On("Get", ctx, orderID).Return(order, nil)
	mockWallets.On("GetByUserID", ctx, userID).Return(wallet, nil)

	service := NewPaymentService(nil, mockOrders)
	service.SetWalletRepository(mockWallets)

	req := CalculateCombinedPaymentRequest{
		OrderID:           orderID,
		WalletAmountCents: 8000, // User wants to use 8000 from wallet
	}

	resp, err := service.CalculateCombinedPayment(ctx, userID, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, totalCents, resp.OrderTotalCents)
	assert.Equal(t, walletBalance, resp.WalletBalanceCents)
	assert.Equal(t, int64(8000), resp.WalletUsableCents)
	assert.Equal(t, int64(2000), resp.ThirdPartyAmountCents)
	assert.True(t, resp.CanPayWithWalletOnly)

	mockOrders.AssertExpectations(t)
	mockWallets.AssertExpectations(t)
}

// TestPaymentService_CalculateCombinedPayment_Unauthorized tests combined payment calculation by unauthorized user
func TestPaymentService_CalculateCombinedPayment_Unauthorized(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	orderID := uint64(1000)
	differentUserID := uint64(999)

	mockOrders := new(MockOrderRepository)

	order := createTestOrder(orderID, differentUserID, model.OrderStatusPending, 10000)

	mockOrders.On("Get", ctx, orderID).Return(order, nil)

	service := NewPaymentService(nil, mockOrders)

	req := CalculateCombinedPaymentRequest{
		OrderID: orderID,
	}

	resp, err := service.CalculateCombinedPayment(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "unauthorized")

	mockOrders.AssertExpectations(t)
}

// TestPaymentService_HandlePaymentCallback_Success tests successful payment callback handling
func TestPaymentService_HandlePaymentCallback_Success(t *testing.T) {
	ctx := context.Background()
	paymentID := uint64(1)
	orderID := uint64(1000)
	userID := uint64(100)
	amountCents := int64(10000)

	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)

	payment := createTestPayment(paymentID, orderID, userID, model.PaymentStatusPending, amountCents, model.PaymentMethodWeChat)
	order := createTestOrder(orderID, userID, model.OrderStatusPending, amountCents)

	mockPayments.On("Get", ctx, paymentID).Return(payment, nil)
	mockPayments.On("Update", ctx, mock.MatchedBy(func(p *model.Payment) bool {
		return p.Status == model.PaymentStatusPaid && p.PaidAt != nil
	})).Return(nil)
	mockOrders.On("Get", ctx, orderID).Return(order, nil)
	mockOrders.On("Update", ctx, mock.MatchedBy(func(o *model.Order) bool {
		return o.Status == model.OrderStatusConfirmed
	})).Return(nil)

	service := NewPaymentService(mockPayments, mockOrders)

	callbackData := map[string]interface{}{
		"payment_id":   paymentID,
		"amount_cents": amountCents,
		"trade_no":     "wx_trade_123",
	}

	err := service.HandlePaymentCallback(ctx, "wechat", callbackData)

	assert.NoError(t, err)

	mockPayments.AssertExpectations(t)
	mockOrders.AssertExpectations(t)
}

// TestPaymentService_HandlePaymentCallback_AmountMismatch tests payment callback with amount mismatch
func TestPaymentService_HandlePaymentCallback_AmountMismatch(t *testing.T) {
	ctx := context.Background()
	paymentID := uint64(1)
	orderID := uint64(1000)
	userID := uint64(100)
	amountCents := int64(10000)

	mockPayments := new(MockPaymentRepository)

	payment := createTestPayment(paymentID, orderID, userID, model.PaymentStatusPending, amountCents, model.PaymentMethodWeChat)

	mockPayments.On("Get", ctx, paymentID).Return(payment, nil)

	service := NewPaymentService(mockPayments, nil)

	callbackData := map[string]interface{}{
		"payment_id":   paymentID,
		"amount_cents": int64(5000), // Different amount
		"trade_no":     "wx_trade_123",
	}

	err := service.HandlePaymentCallback(ctx, "wechat", callbackData)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "amount mismatch")

	mockPayments.AssertExpectations(t)
}

// TestPaymentService_HandlePaymentCallback_ProviderMismatch tests payment callback with provider mismatch
func TestPaymentService_HandlePaymentCallback_ProviderMismatch(t *testing.T) {
	ctx := context.Background()
	paymentID := uint64(1)

	mockPayments := new(MockPaymentRepository)

	payment := createTestPayment(paymentID, 1000, 100, model.PaymentStatusPending, 10000, model.PaymentMethodWeChat)

	mockPayments.On("Get", ctx, paymentID).Return(payment, nil)

	service := NewPaymentService(mockPayments, nil)

	callbackData := map[string]interface{}{
		"payment_id":   paymentID,
		"amount_cents": int64(10000),
		"trade_no":     "ali_trade_123",
	}

	err := service.HandlePaymentCallback(ctx, "alipay", callbackData) // Wrong provider

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "provider mismatch")

	mockPayments.AssertExpectations(t)
}

// TestPaymentService_HandlePaymentCallback_AlreadyPaid tests payment callback for already paid payment
func TestPaymentService_HandlePaymentCallback_AlreadyPaid(t *testing.T) {
	ctx := context.Background()
	paymentID := uint64(1)

	mockPayments := new(MockPaymentRepository)

	payment := createTestPayment(paymentID, 1000, 100, model.PaymentStatusPaid, 10000, model.PaymentMethodWeChat)

	mockPayments.On("Get", ctx, paymentID).Return(payment, nil)

	service := NewPaymentService(mockPayments, nil)

	callbackData := map[string]interface{}{
		"payment_id":   paymentID,
		"amount_cents": int64(10000),
		"trade_no":     "wx_trade_123",
	}

	err := service.HandlePaymentCallback(ctx, "wechat", callbackData)

	// Should not error, should be idempotent
	assert.NoError(t, err)

	mockPayments.AssertExpectations(t)
}

// TestWrapError tests error wrapping functionality
func TestWrapError(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		operation string
		wantType  string
	}{
		{
			name:      "Nil error",
			err:       nil,
			operation: "test",
			wantType:  "",
		},
		{
			name:      "Repository not found error",
			err:       repository.ErrNotFound,
			operation: "get payment",
			wantType:  "*apierr.APIError",
		},
		{
			name:      "Validation error",
			err:       ErrValidation,
			operation: "validate",
			wantType:  "*apierr.APIError",
		},
		{
			name:      "Generic error",
			err:       errors.New("generic error"),
			operation: "process",
			wantType:  "*apierr.APIError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapError(tt.err, tt.operation)

			if tt.err == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
			}
		})
	}
}

// TestPaymentService_SetDistributedLock tests setting distributed lock
func TestPaymentService_SetDistributedLock(t *testing.T) {
	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)

	service := NewPaymentService(mockPayments, mockOrders)

	// Create a mock distributed lock
	mockLock := &MockDistributedLock{}

	// Set the lock
	service.SetDistributedLock(mockLock)

	// Service should still be valid
	assert.NotNil(t, service)
}

// TestPaymentService_SetWalletRepository tests setting wallet repository
func TestPaymentService_SetWalletRepository(t *testing.T) {
	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)

	service := NewPaymentService(mockPayments, mockOrders)

	// Create a mock wallet repository
	mockWallets := new(MockWalletRepository)

	// Set the wallet repository
	service.SetWalletRepository(mockWallets)

	// Service should still be valid
	assert.NotNil(t, service)
}

// MockDistributedLock is a mock implementation of DistributedLock
type MockDistributedLock struct {
	mock.Mock
}

func (m *MockDistributedLock) TryLock(ctx context.Context, key string, ttl time.Duration, retries int, retryDelay time.Duration) (bool, error) {
	args := m.Called(ctx, key, ttl, retries, retryDelay)
	return args.Bool(0), args.Error(1)
}

func (m *MockDistributedLock) Lock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	args := m.Called(ctx, key, ttl)
	return args.Bool(0), args.Error(1)
}

func (m *MockDistributedLock) Unlock(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

// TestPaymentService_RoutingEngineNotInitialized tests behavior when routing engine is not initialized
func TestPaymentService_RoutingEngineNotInitialized(t *testing.T) {
	ctx := context.Background()
	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)

	service := NewPaymentService(mockPayments, mockOrders)

	// Don't initialize routing engine
	// Call a method that uses routing engine internally
	// The service should handle this gracefully

	order := createTestOrder(1000, 100, model.OrderStatusPending, 10000)

	mockOrders.On("Get", ctx, uint64(1000)).Return(order, nil)
	mockPayments.On("List", ctx, mock.AnythingOfType("repository.PaymentListOptions")).Return([]model.Payment{}, int64(0), nil)
	mockPayments.On("Create", ctx, mock.AnythingOfType("*model.Payment")).Return(nil).Run(func(args mock.Arguments) {
		payment := args.Get(1).(*model.Payment)
		payment.ID = 1
	})
	mockPayments.On("Get", ctx, uint64(1)).Return(createTestPayment(1, 1000, 100, model.PaymentStatusPending, 10000, model.PaymentMethodWeChat), nil)
	mockPayments.On("Update", ctx, mock.AnythingOfType("*model.Payment")).Return(nil)
	mockOrders.On("Update", ctx, mock.AnythingOfType("*model.Order")).Return(nil)

	req := CreatePaymentRequest{
		OrderID: 1000,
		Method:  model.PaymentMethodWeChat,
	}

	// Should work without routing engine
	resp, err := service.CreatePayment(ctx, 100, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Greater(t, resp.PaymentID, uint64(0))

	mockOrders.AssertExpectations(t)
	mockPayments.AssertExpectations(t)
}

// TestPaymentService_GetPaymentRoutingLog_WithoutEngine tests getting routing log without engine
func TestPaymentService_GetPaymentRoutingLog_WithoutEngine(t *testing.T) {
	ctx := context.Background()
	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)

	service := NewPaymentService(mockPayments, mockOrders)

	// Don't initialize routing engine
	_, err := service.GetPaymentRoutingLog(ctx, 1)

	// Should return an error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "routing engine not initialized")
}

// TestPaymentService_routePayment_NoEngine tests routePayment without routing engine
func TestPaymentService_routePayment_NoEngine(t *testing.T) {
	ctx := context.Background()
	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)

	service := NewPaymentService(mockPayments, mockOrders)

	// Don't initialize routing engine
	order := &model.Order{
		Base:            model.Base{ID: 1000},
		TotalPriceCents: 10000,
		Game: &model.Game{
			Name: "TestGame",
		},
		ServiceItem: &model.ServiceItem{
			Name: "TestService",
		},
	}

	result, err := service.routePayment(ctx, order, model.PaymentMethodWeChat)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "routing engine not initialized")
	assert.Nil(t, result)
}

// TestPaymentService_CreatePayment_Combined_ThirdPartyMethodRequired tests combined payment without third party method
func TestPaymentService_CreatePayment_Combined_ThirdPartyMethodRequired(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	orderID := uint64(1000)

	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)
	mockWallets := new(MockWalletRepository)

	order := createTestOrder(orderID, userID, model.OrderStatusPending, 10000)

	mockOrders.On("Get", ctx, orderID).Return(order, nil)
	mockPayments.On("List", ctx, mock.AnythingOfType("repository.PaymentListOptions")).Return([]model.Payment{}, int64(0), nil)

	service := NewPaymentService(mockPayments, mockOrders)
	service.SetWalletRepository(mockWallets)

	req := CreatePaymentRequest{
		OrderID:           orderID,
		Method:            model.PaymentMethodCombined,
		WalletAmountCents: 5000,
		// Missing ThirdPartyMethod
	}

	resp, err := service.CreatePayment(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "third party method must be wechat or alipay")

	mockOrders.AssertExpectations(t)
	mockPayments.AssertExpectations(t)
}

// TestPaymentService_CreatePayment_Wallet_SaveFailure tests wallet payment when save fails
func TestPaymentService_CreatePayment_Wallet_SaveFailure(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	orderID := uint64(1000)
	priceCents := int64(10000)

	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)
	mockWallets := new(MockWalletRepository)

	order := createTestOrder(orderID, userID, model.OrderStatusPending, priceCents)
	wallet := createTestWallet(userID, 20000, 0)

	mockOrders.On("Get", ctx, orderID).Return(order, nil)
	mockPayments.On("List", ctx, mock.AnythingOfType("repository.PaymentListOptions")).Return([]model.Payment{}, int64(0), nil)
	mockWallets.On("GetByUserID", ctx, userID).Return(wallet, nil)
	mockWallets.On("UpdateBalanceWithLock", ctx, userID, int64(-10000), 3).Return(nil, errors.New("database error"))

	service := NewPaymentService(mockPayments, mockOrders)
	service.SetWalletRepository(mockWallets)

	req := CreatePaymentRequest{
		OrderID: orderID,
		Method:  model.PaymentMethodWallet,
	}

	resp, err := service.CreatePayment(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to deduct wallet balance")

	mockOrders.AssertExpectations(t)
	mockPayments.AssertExpectations(t)
	mockWallets.AssertExpectations(t)
}

// TestPaymentService_creditWallet_CreateNew tests creditWallet when user has no wallet
func TestPaymentService_creditWallet_CreateNew(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	amountCents := int64(5000)

	mockPayments := new(MockPaymentRepository)
	mockWallets := new(MockWalletRepository)

	mockWallets.On("GetByUserID", ctx, userID).Return(nil, repository.ErrNotFound)
	mockWallets.On("Save", ctx, mock.MatchedBy(func(w *model.Wallet) bool {
		return w.UserID == userID && w.BalanceCents == amountCents
	})).Return(nil)

	service := NewPaymentService(mockPayments, nil)
	service.SetWalletRepository(mockWallets)

	err := service.creditWallet(ctx, userID, amountCents)

	assert.NoError(t, err)
	mockWallets.AssertExpectations(t)
}

// TestPaymentService_creditWallet_Error tests creditWallet when save fails
func TestPaymentService_creditWallet_Error(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	amountCents := int64(5000)

	mockPayments := new(MockPaymentRepository)
	mockWallets := new(MockWalletRepository)

	wallet := createTestWallet(userID, 1000, 0)

	mockWallets.On("GetByUserID", ctx, userID).Return(wallet, nil)
	mockWallets.On("Save", ctx, mock.AnythingOfType("*model.Wallet")).Return(errors.New("database error"))

	service := NewPaymentService(mockPayments, nil)
	service.SetWalletRepository(mockWallets)

	err := service.creditWallet(ctx, userID, amountCents)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	mockWallets.AssertExpectations(t)
}

// TestPaymentService_CreatePayment_Combined_InsufficientWallet tests combined payment with insufficient wallet
func TestPaymentService_CreatePayment_Combined_InsufficientWallet(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	orderID := uint64(1000)

	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)
	mockWallets := new(MockWalletRepository)

	order := createTestOrder(orderID, userID, model.OrderStatusPending, 10000)
	wallet := createTestWallet(userID, 3000, 0) // Only 3000 available

	mockOrders.On("Get", ctx, orderID).Return(order, nil)
	mockPayments.On("List", ctx, mock.AnythingOfType("repository.PaymentListOptions")).Return([]model.Payment{}, int64(0), nil)
	mockWallets.On("GetByUserID", ctx, userID).Return(wallet, nil)

	service := NewPaymentService(mockPayments, mockOrders)
	service.SetWalletRepository(mockWallets)

	req := CreatePaymentRequest{
		OrderID:           orderID,
		Method:            model.PaymentMethodCombined,
		WalletAmountCents: 6000,
		ThirdPartyMethod:  model.PaymentMethodWeChat,
	}

	resp, err := service.CreatePayment(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "insufficient wallet balance")

	mockOrders.AssertExpectations(t)
	mockPayments.AssertExpectations(t)
	mockWallets.AssertExpectations(t)
}

// TestPaymentService_CreatePayment_Combined_CreatePaymentError tests combined payment when payment creation fails
func TestPaymentService_CreatePayment_Combined_CreatePaymentError(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	orderID := uint64(1000)

	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)
	mockWallets := new(MockWalletRepository)

	order := createTestOrder(orderID, userID, model.OrderStatusPending, 10000)
	wallet := createTestWallet(userID, 10000, 0)
	updatedWallet := createTestWallet(userID, 4000, 0)   // After deduction
	rolledBackWallet := createTestWallet(userID, 10000, 0) // After rollback

	mockOrders.On("Get", ctx, orderID).Return(order, nil)
	mockPayments.On("List", ctx, mock.AnythingOfType("repository.PaymentListOptions")).Return([]model.Payment{}, int64(0), nil)
	mockWallets.On("GetByUserID", ctx, userID).Return(wallet, nil)
	mockWallets.On("UpdateBalanceWithLock", ctx, userID, int64(-6000), 3).Return(updatedWallet, nil)
	mockPayments.On("Create", ctx, mock.AnythingOfType("*model.Payment")).Return(errors.New("database error"))
	// Rollback wallet
	mockWallets.On("UpdateBalanceWithLock", ctx, userID, int64(6000), 3).Return(rolledBackWallet, nil)

	service := NewPaymentService(mockPayments, mockOrders)
	service.SetWalletRepository(mockWallets)

	req := CreatePaymentRequest{
		OrderID:           orderID,
		Method:            model.PaymentMethodCombined,
		WalletAmountCents: 6000,
		ThirdPartyMethod:  model.PaymentMethodWeChat,
	}

	resp, err := service.CreatePayment(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "database error")

	mockOrders.AssertExpectations(t)
	mockPayments.AssertExpectations(t)
	mockWallets.AssertExpectations(t)
}

// TestPaymentService_HandlePaymentCallback_FloatPaymentID tests callback with float payment ID
func TestPaymentService_HandlePaymentCallback_FloatPaymentID(t *testing.T) {
	ctx := context.Background()
	paymentID := uint64(1)
	orderID := uint64(1000)
	userID := uint64(100)
	amountCents := int64(10000)

	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)

	payment := createTestPayment(paymentID, orderID, userID, model.PaymentStatusPending, amountCents, model.PaymentMethodWeChat)
	order := createTestOrder(orderID, userID, model.OrderStatusPending, amountCents)

	mockPayments.On("Get", ctx, paymentID).Return(payment, nil)
	mockPayments.On("Update", ctx, mock.MatchedBy(func(p *model.Payment) bool {
		return p.Status == model.PaymentStatusPaid && p.PaidAt != nil
	})).Return(nil)
	mockOrders.On("Get", ctx, orderID).Return(order, nil)
	mockOrders.On("Update", ctx, mock.MatchedBy(func(o *model.Order) bool {
		return o.Status == model.OrderStatusConfirmed
	})).Return(nil)

	service := NewPaymentService(mockPayments, mockOrders)

	// Pass paymentID as float (common in JSON unmarshaling)
	callbackData := map[string]interface{}{
		"payment_id":   float64(paymentID),
		"amount_cents": amountCents,
		"trade_no":     "wx_trade_123",
	}

	err := service.HandlePaymentCallback(ctx, "wechat", callbackData)

	assert.NoError(t, err)

	mockPayments.AssertExpectations(t)
	mockOrders.AssertExpectations(t)
}

// TestPaymentService_HandlePaymentCallback_NoTradeNo tests callback without trade number
func TestPaymentService_HandlePaymentCallback_NoTradeNo(t *testing.T) {
	ctx := context.Background()
	paymentID := uint64(1)
	orderID := uint64(1000)
	userID := uint64(100)
	amountCents := int64(10000)

	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)

	payment := createTestPayment(paymentID, orderID, userID, model.PaymentStatusPending, amountCents, model.PaymentMethodWeChat)
	order := createTestOrder(orderID, userID, model.OrderStatusPending, amountCents)

	mockPayments.On("Get", ctx, paymentID).Return(payment, nil)
	mockPayments.On("Update", ctx, mock.MatchedBy(func(p *model.Payment) bool {
		return p.Status == model.PaymentStatusPaid && p.ProviderTradeNo != ""
	})).Return(nil)
	mockOrders.On("Get", ctx, orderID).Return(order, nil)
	mockOrders.On("Update", ctx, mock.MatchedBy(func(o *model.Order) bool {
		return o.Status == model.OrderStatusConfirmed
	})).Return(nil)

	service := NewPaymentService(mockPayments, mockOrders)

	callbackData := map[string]interface{}{
		"payment_id":   paymentID,
		"amount_cents": amountCents,
		// No trade_no - should generate one
	}

	err := service.HandlePaymentCallback(ctx, "wechat", callbackData)

	assert.NoError(t, err)

	mockPayments.AssertExpectations(t)
	mockOrders.AssertExpectations(t)
}

// TestPaymentService_HandlePaymentCallback_NoAmount tests callback without amount validation
func TestPaymentService_HandlePaymentCallback_NoAmount(t *testing.T) {
	ctx := context.Background()
	paymentID := uint64(1)
	orderID := uint64(1000)
	userID := uint64(100)
	amountCents := int64(10000)

	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)

	payment := createTestPayment(paymentID, orderID, userID, model.PaymentStatusPending, amountCents, model.PaymentMethodWeChat)
	order := createTestOrder(orderID, userID, model.OrderStatusPending, amountCents)

	mockPayments.On("Get", ctx, paymentID).Return(payment, nil)
	mockPayments.On("Update", ctx, mock.MatchedBy(func(p *model.Payment) bool {
		return p.Status == model.PaymentStatusPaid
	})).Return(nil)
	mockOrders.On("Get", ctx, orderID).Return(order, nil)
	mockOrders.On("Update", ctx, mock.MatchedBy(func(o *model.Order) bool {
		return o.Status == model.OrderStatusConfirmed
	})).Return(nil)

	service := NewPaymentService(mockPayments, mockOrders)

	callbackData := map[string]interface{}{
		"payment_id": paymentID,
		// No amount_cents field - should skip validation
		"trade_no":   "wx_trade_123",
	}

	err := service.HandlePaymentCallback(ctx, "wechat", callbackData)

	assert.NoError(t, err)

	mockPayments.AssertExpectations(t)
	mockOrders.AssertExpectations(t)
}
