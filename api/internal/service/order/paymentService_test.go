package order

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	repoiface "gamelink/internal/repository/interfaces"
)

// MockPaymentRepository for PaymentService tests
type MockPaymentRepositoryForPayment struct {
	get              func(ctx context.Context, id uint64) (*model.Payment, error)
	create           func(ctx context.Context, payment *model.Payment) error
	update           func(ctx context.Context, payment *model.Payment) error
	list             func(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error)
	getWithRelations func(ctx context.Context, id uint64) (*model.Payment, error)
	getByOrderID     func(ctx context.Context, orderID uint64) ([]model.Payment, error)
	delete           func(ctx context.Context, id uint64) error
}

func (m *MockPaymentRepositoryForPayment) Create(ctx context.Context, payment *model.Payment) error {
	if m.create != nil {
		return m.create(ctx, payment)
	}
	return nil
}

func (m *MockPaymentRepositoryForPayment) Get(ctx context.Context, id uint64) (*model.Payment, error) {
	if m.get != nil {
		return m.get(ctx, id)
	}
	return nil, nil
}

func (m *MockPaymentRepositoryForPayment) GetWithRelations(ctx context.Context, id uint64) (*model.Payment, error) {
	if m.getWithRelations != nil {
		return m.getWithRelations(ctx, id)
	}
	return m.Get(ctx, id)
}

func (m *MockPaymentRepositoryForPayment) Update(ctx context.Context, payment *model.Payment) error {
	if m.update != nil {
		return m.update(ctx, payment)
	}
	return nil
}

func (m *MockPaymentRepositoryForPayment) Delete(ctx context.Context, id uint64) error {
	if m.delete != nil {
		return m.delete(ctx, id)
	}
	return nil
}

func (m *MockPaymentRepositoryForPayment) List(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
	if m.list != nil {
		return m.list(ctx, opts)
	}
	return nil, 0, nil
}

func (m *MockPaymentRepositoryForPayment) GetByOrderID(ctx context.Context, orderID uint64) ([]model.Payment, error) {
	if m.getByOrderID != nil {
		return m.getByOrderID(ctx, orderID)
	}
	return nil, nil
}

func (m *MockPaymentRepositoryForPayment) GetByRequestID(ctx context.Context, requestID string) (*model.Payment, error) {
	return nil, nil
}

// MockOrderReadWriter for PaymentService tests
type MockOrderReadWriter struct {
	repoiface.OrderRepository
	getOrder    func(ctx context.Context, id uint64) (*model.Order, error)
	updateOrder func(ctx context.Context, order *model.Order) error
}

func (m *MockOrderReadWriter) Get(ctx context.Context, id uint64) (*model.Order, error) {
	if m.getOrder != nil {
		return m.getOrder(ctx, id)
	}
	return nil, nil
}

func (m *MockOrderReadWriter) Update(ctx context.Context, order *model.Order) error {
	if m.updateOrder != nil {
		return m.updateOrder(ctx, order)
	}
	return nil
}

// MockWalletRepository for refund tests
type MockWalletRepository struct {
	getByUserID func(ctx context.Context, userID uint64) (*model.Wallet, error)
	save        func(ctx context.Context, wallet *model.Wallet) error
}

func (m *MockWalletRepository) GetByUserID(ctx context.Context, userID uint64) (*model.Wallet, error) {
	if m.getByUserID != nil {
		return m.getByUserID(ctx, userID)
	}
	return &model.Wallet{UserID: userID, BalanceCents: 10000}, nil
}

func (m *MockWalletRepository) Save(ctx context.Context, wallet *model.Wallet) error {
	if m.save != nil {
		return m.save(ctx, wallet)
	}
	return nil
}

func (m *MockWalletRepository) SaveWithOptimisticLock(ctx context.Context, wallet *model.Wallet) error {
	if m.save != nil {
		return m.save(ctx, wallet)
	}
	return nil
}

func (m *MockWalletRepository) UpdateBalanceWithLock(ctx context.Context, userID uint64, delta int64, maxRetries int) (*model.Wallet, error) {
	wallet, err := m.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	wallet.BalanceCents += delta
	return wallet, m.Save(ctx, wallet)
}

// TestPaymentService_CreatePayment_Success tests successful payment creation
func TestPaymentService_CreatePayment_Success(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(1)

	testOrder := createTestOrder(orderID, userID, model.OrderStatusPending)

	var createdPayment *model.Payment
	var updatedOrder *model.Order

	orders := &MockOrderReadWriter{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			if id == orderID {
				return testOrder, nil
			}
			return nil, repository.ErrNotFound
		},
		updateOrder: func(ctx context.Context, order *model.Order) error {
			updatedOrder = order
			return nil
		},
	}

	payments := &MockPaymentRepositoryForPayment{
		list: func(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
			return []model.Payment{}, 0, nil // No existing payments
		},
		create: func(ctx context.Context, payment *model.Payment) error {
			payment.ID = 123
			createdPayment = payment
			return nil
		},
		get: func(ctx context.Context, id uint64) (*model.Payment, error) {
			if createdPayment != nil && id == createdPayment.ID {
				return createdPayment, nil
			}
			return &model.Payment{Base: model.Base{ID: id}, OrderID: orderID, Status: model.PaymentStatusPending}, nil
		},
		update: func(ctx context.Context, payment *model.Payment) error {
			if createdPayment != nil && payment.ID == createdPayment.ID {
				createdPayment = payment
			}
			return nil
		},
	}

	service := NewPaymentService(payments, orders)

	req := CreatePaymentRequest{
		OrderID: orderID,
		Method:  model.PaymentMethodAlipay,
	}

	resp, err := service.CreatePayment(ctx, userID, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, uint64(123), resp.PaymentID)
	assert.NotEmpty(t, resp.PayInfo)
	assert.NotNil(t, updatedOrder)
	assert.Equal(t, model.OrderStatusConfirmed, updatedOrder.Status)
}

// TestPaymentService_CreatePayment_OrderNotFound tests payment creation with non-existent order
func TestPaymentService_CreatePayment_OrderNotFound(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)

	orders := &MockOrderReadWriter{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return nil, repository.ErrNotFound
		},
	}

	payments := &MockPaymentRepositoryForPayment{}

	service := NewPaymentService(payments, orders)

	req := CreatePaymentRequest{
		OrderID: 999,
		Method:  model.PaymentMethodAlipay,
	}

	resp, err := service.CreatePayment(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, repository.ErrNotFound, err)
}

// TestPaymentService_CreatePayment_Unauthorized tests payment creation by unauthorized user
func TestPaymentService_CreatePayment_Unauthorized(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(1)
	orderUserID := uint64(999)

	testOrder := createTestOrder(orderID, orderUserID, model.OrderStatusPending)

	orders := &MockOrderReadWriter{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return testOrder, nil
		},
	}

	payments := &MockPaymentRepositoryForPayment{
		list: func(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
			return []model.Payment{}, 0, nil
		},
	}

	service := NewPaymentService(payments, orders)

	req := CreatePaymentRequest{
		OrderID: orderID,
		Method:  model.PaymentMethodAlipay,
	}

	resp, err := service.CreatePayment(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

// TestPaymentService_CreatePayment_InvalidOrderStatus tests payment creation with invalid order status
func TestPaymentService_CreatePayment_InvalidOrderStatus(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(1)

	testOrder := createTestOrder(orderID, userID, model.OrderStatusConfirmed)

	orders := &MockOrderReadWriter{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return testOrder, nil
		},
	}

	payments := &MockPaymentRepositoryForPayment{
		list: func(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
			return []model.Payment{}, 0, nil
		},
	}

	service := NewPaymentService(payments, orders)

	req := CreatePaymentRequest{
		OrderID: orderID,
		Method:  model.PaymentMethodAlipay,
	}

	resp, err := service.CreatePayment(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, ErrInvalidOrderStatus, err)
}

// TestPaymentService_CreatePayment_AlreadyPaid tests payment creation when order is already paid
func TestPaymentService_CreatePayment_AlreadyPaid(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(1)

	testOrder := createTestOrder(orderID, userID, model.OrderStatusPending)

	orders := &MockOrderReadWriter{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return testOrder, nil
		},
	}

	paidTime := time.Now()
	payments := &MockPaymentRepositoryForPayment{
		list: func(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
			return []model.Payment{
				{Base: model.Base{ID: 1}, OrderID: orderID, Status: model.PaymentStatusPaid, PaidAt: &paidTime},
			}, 1, nil
		},
	}

	service := NewPaymentService(payments, orders)

	req := CreatePaymentRequest{
		OrderID: orderID,
		Method:  model.PaymentMethodAlipay,
	}

	resp, err := service.CreatePayment(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, ErrOrderAlreadyPaid, err)
}

// TestPaymentService_CreatePayment_PendingExists tests payment creation when pending payment exists
func TestPaymentService_CreatePayment_PendingExists(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(1)

	testOrder := createTestOrder(orderID, userID, model.OrderStatusPending)

	orders := &MockOrderReadWriter{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return testOrder, nil
		},
	}

	payments := &MockPaymentRepositoryForPayment{
		list: func(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
			return []model.Payment{
				{Base: model.Base{ID: 1}, OrderID: orderID, Status: model.PaymentStatusPending},
			}, 1, nil
		},
	}

	service := NewPaymentService(payments, orders)

	req := CreatePaymentRequest{
		OrderID: orderID,
		Method:  model.PaymentMethodAlipay,
	}

	resp, err := service.CreatePayment(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

// TestPaymentService_GetPaymentStatus_Success tests successful payment status retrieval
func TestPaymentService_GetPaymentStatus_Success(t *testing.T) {
	ctx := context.Background()
	paymentID := uint64(1)

	paidTime := time.Now()

	payments := &MockPaymentRepositoryForPayment{
		get: func(ctx context.Context, id uint64) (*model.Payment, error) {
			return &model.Payment{
				Base:        model.Base{ID: paymentID},
				OrderID:     100,
				Status:      model.PaymentStatusPaid,
				PaidAt:      &paidTime,
				AmountCents: 5000,
			}, nil
		},
	}

	service := NewPaymentService(payments, nil)

	resp, err := service.GetPaymentStatus(ctx, paymentID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, paymentID, resp.PaymentID)
	assert.Equal(t, uint64(100), resp.OrderID)
	assert.Equal(t, model.PaymentStatusPaid, resp.Status)
	assert.NotNil(t, resp.PaidAt)
}

// TestPaymentService_GetPaymentStatus_NotFound tests payment status when payment doesn't exist
func TestPaymentService_GetPaymentStatus_NotFound(t *testing.T) {
	ctx := context.Background()

	payments := &MockPaymentRepositoryForPayment{
		get: func(ctx context.Context, id uint64) (*model.Payment, error) {
			return nil, repository.ErrNotFound
		},
	}

	service := NewPaymentService(payments, nil)

	resp, err := service.GetPaymentStatus(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, ErrPaymentNotFound, err)
}

// TestPaymentService_CancelPayment_Success tests successful payment cancellation
func TestPaymentService_CancelPayment_Success(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	paymentID := uint64(1)

	var updatedPayment *model.Payment

	payments := &MockPaymentRepositoryForPayment{
		get: func(ctx context.Context, id uint64) (*model.Payment, error) {
			return &model.Payment{
				Base:        model.Base{ID: paymentID},
				UserID:      userID,
				OrderID:     100,
				Status:      model.PaymentStatusPending,
				AmountCents: 5000,
			}, nil
		},
		update: func(ctx context.Context, payment *model.Payment) error {
			updatedPayment = payment
			return nil
		},
	}

	service := NewPaymentService(payments, nil)

	err := service.CancelPayment(ctx, userID, paymentID)

	require.NoError(t, err)
	assert.NotNil(t, updatedPayment)
	assert.Equal(t, model.PaymentStatusFailed, updatedPayment.Status)
}

// TestPaymentService_CancelPayment_NotFound tests cancellation when payment doesn't exist
func TestPaymentService_CancelPayment_NotFound(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)

	payments := &MockPaymentRepositoryForPayment{
		get: func(ctx context.Context, id uint64) (*model.Payment, error) {
			return nil, repository.ErrNotFound
		},
	}

	service := NewPaymentService(payments, nil)

	err := service.CancelPayment(ctx, userID, 999)

	assert.Error(t, err)
	assert.Equal(t, repository.ErrNotFound, err)
}

// TestPaymentService_CancelPayment_Unauthorized tests payment cancellation by unauthorized user
func TestPaymentService_CancelPayment_Unauthorized(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	paymentID := uint64(1)

	payments := &MockPaymentRepositoryForPayment{
		get: func(ctx context.Context, id uint64) (*model.Payment, error) {
			return &model.Payment{
				Base:        model.Base{ID: paymentID},
				UserID:      999, // Different user
				Status:      model.PaymentStatusPending,
				AmountCents: 5000,
			}, nil
		},
	}

	service := NewPaymentService(payments, nil)

	err := service.CancelPayment(ctx, userID, paymentID)

	assert.Error(t, err)
}

// TestPaymentService_CancelPayment_InvalidStatus tests cancellation with invalid status
func TestPaymentService_CancelPayment_InvalidStatus(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	paymentID := uint64(1)

	paidTime := time.Now()

	payments := &MockPaymentRepositoryForPayment{
		get: func(ctx context.Context, id uint64) (*model.Payment, error) {
			return &model.Payment{
				Base:        model.Base{ID: paymentID},
				UserID:      userID,
				Status:      model.PaymentStatusPaid,
				PaidAt:      &paidTime,
				AmountCents: 5000,
			}, nil
		},
	}

	service := NewPaymentService(payments, nil)

	err := service.CancelPayment(ctx, userID, paymentID)

	assert.Error(t, err)
}

// TestPaymentService_HandlePaymentCallback_Success tests successful payment callback handling
func TestPaymentService_HandlePaymentCallback_Success(t *testing.T) {
	ctx := context.Background()
	paymentID := uint64(1)
	orderID := uint64(100)

	var updatedPayment *model.Payment
	var updatedOrder *model.Order

	payments := &MockPaymentRepositoryForPayment{
		get: func(ctx context.Context, id uint64) (*model.Payment, error) {
			return &model.Payment{
				Base:        model.Base{ID: paymentID},
				OrderID:     orderID,
				Method:      model.PaymentMethodAlipay,
				Status:      model.PaymentStatusPending,
				AmountCents: 5000,
			}, nil
		},
		update: func(ctx context.Context, payment *model.Payment) error {
			updatedPayment = payment
			return nil
		},
	}

	orders := &MockOrderReadWriter{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return createTestOrder(orderID, 1, model.OrderStatusPending), nil
		},
		updateOrder: func(ctx context.Context, order *model.Order) error {
			updatedOrder = order
			return nil
		},
	}

	service := NewPaymentService(payments, orders)

	data := map[string]interface{}{
		"payment_id":   paymentID,
		"amount_cents": int64(5000),
		"trade_no":     "ali_test_123",
	}

	err := service.HandlePaymentCallback(ctx, "alipay", data)

	require.NoError(t, err)
	assert.NotNil(t, updatedPayment)
	assert.Equal(t, model.PaymentStatusPaid, updatedPayment.Status)
	assert.NotNil(t, updatedPayment.PaidAt)
	assert.NotNil(t, updatedOrder)
	assert.Equal(t, model.OrderStatusConfirmed, updatedOrder.Status)
}

// TestPaymentService_HandlePaymentCallback_AlreadyPaid tests callback when payment is already paid
func TestPaymentService_HandlePaymentCallback_AlreadyPaid(t *testing.T) {
	ctx := context.Background()
	paymentID := uint64(1)

	paidTime := time.Now()

	payments := &MockPaymentRepositoryForPayment{
		get: func(ctx context.Context, id uint64) (*model.Payment, error) {
			return &model.Payment{
				Base:        model.Base{ID: paymentID},
				OrderID:     100,
				Method:      model.PaymentMethodAlipay,
				Status:      model.PaymentStatusPaid,
				PaidAt:      &paidTime,
				AmountCents: 5000,
			}, nil
		},
	}

	orders := &MockOrderReadWriter{}

	service := NewPaymentService(payments, orders)

	data := map[string]interface{}{
		"payment_id": paymentID,
	}

	err := service.HandlePaymentCallback(ctx, "alipay", data)

	// Should not error - already paid payments are idempotent
	require.NoError(t, err)
}

// TestPaymentService_HandlePaymentCallback_ProviderMismatch tests callback with provider mismatch
func TestPaymentService_HandlePaymentCallback_ProviderMismatch(t *testing.T) {
	ctx := context.Background()
	paymentID := uint64(1)

	payments := &MockPaymentRepositoryForPayment{
		get: func(ctx context.Context, id uint64) (*model.Payment, error) {
			return &model.Payment{
				Base:        model.Base{ID: paymentID},
				OrderID:     100,
				Method:      model.PaymentMethodAlipay,
				Status:      model.PaymentStatusPending,
				AmountCents: 5000,
			}, nil
		},
	}

	orders := &MockOrderReadWriter{}

	service := NewPaymentService(payments, orders)

	data := map[string]interface{}{
		"payment_id": paymentID,
	}

	err := service.HandlePaymentCallback(ctx, "wechat", data)

	assert.Error(t, err)
}

// TestPaymentService_HandlePaymentCallback_AmountMismatch tests callback with amount mismatch
func TestPaymentService_HandlePaymentCallback_AmountMismatch(t *testing.T) {
	ctx := context.Background()
	paymentID := uint64(1)

	payments := &MockPaymentRepositoryForPayment{
		get: func(ctx context.Context, id uint64) (*model.Payment, error) {
			return &model.Payment{
				Base:        model.Base{ID: paymentID},
				OrderID:     100,
				Method:      model.PaymentMethodAlipay,
				Status:      model.PaymentStatusPending,
				AmountCents: 5000,
			}, nil
		},
	}

	orders := &MockOrderReadWriter{}

	service := NewPaymentService(payments, orders)

	data := map[string]interface{}{
		"payment_id":   paymentID,
		"amount_cents": int64(3000), // Wrong amount
	}

	err := service.HandlePaymentCallback(ctx, "alipay", data)

	assert.Error(t, err)
}

// TestPaymentService_RefundPayment_Success tests successful refund
func TestPaymentService_RefundPayment_Success(t *testing.T) {
	ctx := context.Background()
	paymentID := uint64(1)
	orderID := uint64(100)

	paidTime := time.Now()

	var updatedPayment *model.Payment
	var updatedOrder *model.Order
	var walletCredited bool

	payments := &MockPaymentRepositoryForPayment{
		get: func(ctx context.Context, id uint64) (*model.Payment, error) {
			return &model.Payment{
				Base:            model.Base{ID: paymentID},
				OrderID:         orderID,
				UserID:          1,
				Method:          model.PaymentMethodAlipay,
				Status:          model.PaymentStatusPaid,
				PaidAt:          &paidTime,
				AmountCents:     5000,
				ProviderTradeNo: "ali_test_123",
			}, nil
		},
		update: func(ctx context.Context, payment *model.Payment) error {
			updatedPayment = payment
			return nil
		},
	}

	orders := &MockOrderReadWriter{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return createTestOrder(orderID, 1, model.OrderStatusConfirmed), nil
		},
		updateOrder: func(ctx context.Context, order *model.Order) error {
			updatedOrder = order
			return nil
		},
	}

	wallets := &MockWalletRepository{
		save: func(ctx context.Context, wallet *model.Wallet) error {
			walletCredited = true
			return nil
		},
	}

	service := NewPaymentService(payments, orders)
	service.SetWalletRepository(wallets)

	err := service.RefundPayment(ctx, paymentID, "User requested refund")

	require.NoError(t, err)
	assert.NotNil(t, updatedPayment)
	assert.Equal(t, model.PaymentStatusRefunded, updatedPayment.Status)
	assert.NotNil(t, updatedPayment.RefundedAt)
	assert.NotNil(t, updatedOrder)
	assert.Equal(t, model.OrderStatusRefunded, updatedOrder.Status)
	assert.Equal(t, int64(5000), updatedOrder.RefundAmountCents)
	assert.True(t, walletCredited)
}

// TestPaymentService_RefundPayment_InvalidStatus tests refund with invalid status
func TestPaymentService_RefundPayment_InvalidStatus(t *testing.T) {
	ctx := context.Background()
	paymentID := uint64(1)

	payments := &MockPaymentRepositoryForPayment{
		get: func(ctx context.Context, id uint64) (*model.Payment, error) {
			return &model.Payment{
				Base:        model.Base{ID: paymentID},
				Method:      model.PaymentMethodAlipay,
				Status:      model.PaymentStatusPending,
				AmountCents: 5000,
			}, nil
		},
	}

	orders := &MockOrderReadWriter{}

	service := NewPaymentService(payments, orders)

	err := service.RefundPayment(ctx, paymentID, "Test refund")

	assert.Error(t, err)
}

// TestPaymentService_ListPayments_Success tests successful payment list retrieval
func TestPaymentService_ListPayments_Success(t *testing.T) {
	ctx := context.Background()

	payments := &MockPaymentRepositoryForPayment{
		list: func(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
			return []model.Payment{
				{Base: model.Base{ID: 1}, Status: model.PaymentStatusPaid},
				{Base: model.Base{ID: 2}, Status: model.PaymentStatusPending},
			}, 2, nil
		},
	}

	service := NewPaymentService(payments, nil)

	opts := repository.PaymentListOptions{
		Page:     1,
		PageSize: 10,
	}

	paymentsList, total, err := service.ListPayments(ctx, opts)

	require.NoError(t, err)
	assert.Equal(t, 2, len(paymentsList))
	assert.Equal(t, int64(2), total)
}

// TestPaymentService_ListPayments_DefaultPagination tests default pagination
func TestPaymentService_ListPayments_DefaultPagination(t *testing.T) {
	ctx := context.Background()

	var actualOpts repository.PaymentListOptions

	payments := &MockPaymentRepositoryForPayment{
		list: func(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
			actualOpts = opts
			return []model.Payment{}, 0, nil
		},
	}

	service := NewPaymentService(payments, nil)

	// Invalid page and page size
	opts := repository.PaymentListOptions{
		Page:     0,
		PageSize: 0,
	}

	_, _, err := service.ListPayments(ctx, opts)

	require.NoError(t, err)
	assert.Equal(t, 1, actualOpts.Page)
	assert.Equal(t, 20, actualOpts.PageSize)
}

// TestPaymentService_creditWallet_CreatesNewWallet tests wallet credit when wallet doesn't exist
func TestPaymentService_creditWallet_CreatesNewWallet(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	amount := int64(5000)

	var savedWallet *model.Wallet

	wallets := &MockWalletRepository{
		getByUserID: func(ctx context.Context, id uint64) (*model.Wallet, error) {
			return nil, repository.ErrNotFound
		},
		save: func(ctx context.Context, wallet *model.Wallet) error {
			savedWallet = wallet
			return nil
		},
	}

	service := NewPaymentService(nil, nil)
	service.SetWalletRepository(wallets)

	err := service.creditWallet(ctx, userID, amount)

	require.NoError(t, err)
	assert.NotNil(t, savedWallet)
	assert.Equal(t, userID, savedWallet.UserID)
	assert.Equal(t, amount, savedWallet.BalanceCents)
}

// TestPaymentService_creditWallet_UpdatesExisting tests wallet credit when wallet exists
func TestPaymentService_creditWallet_UpdatesExisting(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	amount := int64(5000)

	var savedWallet *model.Wallet

	wallets := &MockWalletRepository{
		getByUserID: func(ctx context.Context, id uint64) (*model.Wallet, error) {
			return &model.Wallet{UserID: userID, BalanceCents: 3000}, nil
		},
		save: func(ctx context.Context, wallet *model.Wallet) error {
			savedWallet = wallet
			return nil
		},
	}

	service := NewPaymentService(nil, nil)
	service.SetWalletRepository(wallets)

	err := service.creditWallet(ctx, userID, amount)

	require.NoError(t, err)
	assert.NotNil(t, savedWallet)
	assert.Equal(t, int64(8000), savedWallet.BalanceCents) // 3000 + 5000
}

// TestPaymentService_SetDistributedLock tests distributed lock injection
func TestPaymentService_SetDistributedLock(t *testing.T) {
	payments := &MockPaymentRepositoryForPayment{}
	orders := &MockOrderReadWriter{}

	service := NewPaymentService(payments, orders)

	// Should be nil initially
	assert.Nil(t, service.distributedLock)

	// Create a mock lock
	mockLock := &MockDistributedLock{}
	service.distributedLock = mockLock

	// Should now be set
	assert.NotNil(t, service.distributedLock)
}

// TestPaymentService_SetWalletRepository tests wallet repository injection
func TestPaymentService_SetWalletRepository(t *testing.T) {
	payments := &MockPaymentRepositoryForPayment{}
	orders := &MockOrderReadWriter{}

	service := NewPaymentService(payments, orders)

	// Should be nil initially
	assert.Nil(t, service.wallets)

	// Create a mock wallet repo
	wallets := &MockWalletRepository{}
	service.SetWalletRepository(wallets)

	// Should now be set
	assert.NotNil(t, service.wallets)
}

// TestWechatProvider_Refund tests WeChat provider refund
func TestWechatProvider_Refund(t *testing.T) {
	ctx := context.Background()
	provider := wechatProvider{}

	payment := &model.Payment{
		Base:   model.Base{ID: 123},
		Method: model.PaymentMethodWeChat,
	}

	tradeNo, raw, refundedAt, err := provider.Refund(ctx, payment, "Test refund")

	require.NoError(t, err)
	assert.NotEmpty(t, tradeNo)
	assert.Contains(t, tradeNo, "wx_refund_123")
	assert.NotNil(t, raw)
	assert.NotNil(t, refundedAt)
}

// TestAlipayProvider_Refund tests Alipay provider refund
func TestAlipayProvider_Refund(t *testing.T) {
	ctx := context.Background()
	provider := alipayProvider{}

	payment := &model.Payment{
		Base:   model.Base{ID: 456},
		Method: model.PaymentMethodAlipay,
	}

	tradeNo, raw, refundedAt, err := provider.Refund(ctx, payment, "Test refund")

	require.NoError(t, err)
	assert.NotEmpty(t, tradeNo)
	assert.Contains(t, tradeNo, "ali_refund_456")
	assert.NotNil(t, raw)
	assert.NotNil(t, refundedAt)
}

// TestGenericProvider_Refund tests generic provider refund
func TestGenericProvider_Refund(t *testing.T) {
	ctx := context.Background()
	provider := genericProvider{}

	payment := &model.Payment{
		Base:   model.Base{ID: 789},
		Method: "unknown",
	}

	tradeNo, raw, refundedAt, err := provider.Refund(ctx, payment, "Test refund")

	require.NoError(t, err)
	assert.NotEmpty(t, tradeNo)
	assert.Contains(t, tradeNo, "refund_789")
	assert.NotNil(t, raw)
	assert.NotNil(t, refundedAt)
}
