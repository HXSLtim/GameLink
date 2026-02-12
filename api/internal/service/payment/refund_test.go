package payment

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// MockRefundRecordRepository is a mock implementation of RefundRecordRepository
type MockRefundRecordRepository struct {
	mock.Mock
}

func (m *MockRefundRecordRepository) Create(ctx context.Context, record *model.RefundRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockRefundRecordRepository) Get(ctx context.Context, id uint64) (*model.RefundRecord, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RefundRecord), args.Error(1)
}

func (m *MockRefundRecordRepository) Update(ctx context.Context, record *model.RefundRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockRefundRecordRepository) ListByPaymentID(ctx context.Context, paymentID uint64) ([]model.RefundRecord, error) {
	args := m.Called(ctx, paymentID)
	return args.Get(0).([]model.RefundRecord), args.Error(1)
}

func (m *MockRefundRecordRepository) ListByOrderID(ctx context.Context, orderID uint64) ([]model.RefundRecord, error) {
	args := m.Called(ctx, orderID)
	return args.Get(0).([]model.RefundRecord), args.Error(1)
}

// MockOperationLogRepository is a mock implementation of OperationLogRepository
type MockOperationLogRepository struct {
	mock.Mock
}

func (m *MockOperationLogRepository) Append(ctx context.Context, log *model.OperationLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

func (m *MockOperationLogRepository) ListByEntity(ctx context.Context, entityType string, entityID uint64, opts repository.OperationLogListOptions) ([]model.OperationLog, int64, error) {
	args := m.Called(ctx, entityType, entityID, opts)
	return args.Get(0).([]model.OperationLog), args.Get(1).(int64), args.Error(2)
}

func (m *MockOperationLogRepository) List(ctx context.Context, opts repository.OperationLogSearchOptions) ([]model.OperationLog, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.OperationLog), args.Get(1).(int64), args.Error(2)
}

// TestRefundService_ProcessRefund_Success tests successful refund processing
func TestRefundService_ProcessRefund_Success(t *testing.T) {
	ctx := context.Background()
	paymentID := uint64(1)
	orderID := uint64(1000)
	userID := uint64(100)
	amountCents := int64(5000)

	mockPayments := new(MockPaymentRepository)
	mockRefunds := new(MockRefundRecordRepository)
	mockOrders := new(MockOrderRepository)
	mockWallets := new(MockWalletRepository)
	mockOpLogs := new(MockOperationLogRepository)

	payment := createTestPayment(paymentID, orderID, userID, model.PaymentStatusPaid, 10000, model.PaymentMethodWeChat)
	payment.RefundedAmountCents = 0

	operatorID := uint64(200)

	mockPayments.On("Get", ctx, paymentID).Return(payment, nil)
	mockRefunds.On("Create", ctx, mock.AnythingOfType("*model.RefundRecord")).Return(nil).Run(func(args mock.Arguments) {
		record := args.Get(1).(*model.RefundRecord)
		record.ID = 1
	})
	mockRefunds.On("Update", ctx, mock.MatchedBy(func(r *model.RefundRecord) bool {
		return r.Status == model.RefundStatusProcessed
	})).Return(nil)
	mockPayments.On("Update", ctx, mock.MatchedBy(func(p *model.Payment) bool {
		return p.RefundedAmountCents == amountCents
	})).Return(nil)
	mockWallets.On("GetByUserID", ctx, userID).Return(&model.Wallet{UserID: userID, BalanceCents: 0}, nil)
	mockWallets.On("Save", ctx, mock.AnythingOfType("*model.Wallet")).Return(nil)
	mockOpLogs.On("Append", ctx, mock.AnythingOfType("*model.OperationLog")).Return(nil)

	service := NewRefundService(mockPayments, mockRefunds, mockOrders)
	service.wallets = mockWallets
	service.SetOperationLogRepository(mockOpLogs)

	req := model.RefundRequest{
		PaymentID:   paymentID,
		AmountCents: amountCents,
		Reason:      "customer request",
		OperatorID:  &operatorID,
	}

	resp, err := service.ProcessRefund(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, amountCents, resp.RefundRecord.AmountCents)
	assert.Equal(t, "customer request", resp.RefundRecord.Reason)
	assert.Equal(t, model.RefundStatusProcessed, resp.RefundRecord.Status)
	assert.Equal(t, int64(5000), resp.RemainingAmount) // 10000 - 5000

	mockPayments.AssertExpectations(t)
	mockRefunds.AssertExpectations(t)
	mockWallets.AssertExpectations(t)
	mockOpLogs.AssertExpectations(t)
}

// TestRefundService_ProcessRefund_FullRefund tests full refund processing
func TestRefundService_ProcessRefund_FullRefund(t *testing.T) {
	ctx := context.Background()
	paymentID := uint64(1)
	orderID := uint64(1000)
	userID := uint64(100)
	totalCents := int64(10000)

	mockPayments := new(MockPaymentRepository)
	mockRefunds := new(MockRefundRecordRepository)
	mockOrders := new(MockOrderRepository)
	mockWallets := new(MockWalletRepository)
	mockOpLogs := new(MockOperationLogRepository)

	payment := createTestPayment(paymentID, orderID, userID, model.PaymentStatusPaid, totalCents, model.PaymentMethodWeChat)
	payment.RefundedAmountCents = 0

	order := createTestOrder(orderID, userID, model.OrderStatusConfirmed, totalCents)

	mockPayments.On("Get", ctx, paymentID).Return(payment, nil)
	mockRefunds.On("Create", ctx, mock.AnythingOfType("*model.RefundRecord")).Return(nil).Run(func(args mock.Arguments) {
		record := args.Get(1).(*model.RefundRecord)
		record.ID = 1
	})
	mockRefunds.On("Update", ctx, mock.MatchedBy(func(r *model.RefundRecord) bool {
		return r.Status == model.RefundStatusProcessed
	})).Return(nil)
	mockPayments.On("Update", ctx, mock.MatchedBy(func(p *model.Payment) bool {
		return p.RefundedAmountCents == totalCents && p.Status == model.PaymentStatusRefunded
	})).Return(nil)
	mockOrders.On("Get", ctx, orderID).Return(order, nil)
	mockOrders.On("Update", ctx, mock.MatchedBy(func(o *model.Order) bool {
		return o.Status == model.OrderStatusRefunded && o.RefundAmountCents == totalCents
	})).Return(nil)
	mockWallets.On("GetByUserID", ctx, userID).Return(&model.Wallet{UserID: userID, BalanceCents: 0}, nil)
	mockWallets.On("Save", ctx, mock.AnythingOfType("*model.Wallet")).Return(nil)
	mockOpLogs.On("Append", ctx, mock.AnythingOfType("*model.OperationLog")).Return(nil)

	service := NewRefundService(mockPayments, mockRefunds, mockOrders)
	service.wallets = mockWallets
	service.SetOperationLogRepository(mockOpLogs)

	req := model.RefundRequest{
		PaymentID:   paymentID,
		AmountCents: totalCents,
		Reason:      "full refund",
	}

	resp, err := service.ProcessRefund(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, totalCents, resp.RefundRecord.AmountCents)
	assert.Equal(t, model.PaymentStatusRefunded, resp.Payment.Status)
	assert.True(t, resp.Payment.IsFullyRefunded())
	assert.Equal(t, int64(0), resp.RemainingAmount)

	mockPayments.AssertExpectations(t)
	mockRefunds.AssertExpectations(t)
	mockOrders.AssertExpectations(t)
	mockWallets.AssertExpectations(t)
	mockOpLogs.AssertExpectations(t)
}

// TestRefundService_ProcessRefund_InvalidAmount tests refund with invalid amount
func TestRefundService_ProcessRefund_InvalidAmount(t *testing.T) {
	ctx := context.Background()
	paymentID := uint64(1)

	mockPayments := new(MockPaymentRepository)
	mockRefunds := new(MockRefundRecordRepository)
	mockOrders := new(MockOrderRepository)

	payment := createTestPayment(paymentID, 1000, 100, model.PaymentStatusPaid, 10000, model.PaymentMethodWeChat)

	mockPayments.On("Get", ctx, paymentID).Return(payment, nil)

	service := NewRefundService(mockPayments, mockRefunds, mockOrders)

	tests := []struct {
		name        string
		amountCents int64
		expectError string
	}{
		{
			name:        "Zero amount",
			amountCents: 0,
			expectError: "refund amount must be positive",
		},
		{
			name:        "Negative amount",
			amountCents: -100,
			expectError: "refund amount must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := model.RefundRequest{
				PaymentID:   paymentID,
				AmountCents: tt.amountCents,
				Reason:      "test",
			}

			resp, err := service.ProcessRefund(ctx, req)

			assert.Error(t, err)
			assert.Nil(t, resp)
			assert.Contains(t, err.Error(), tt.expectError)
		})
	}

	mockPayments.AssertExpectations(t)
}

// TestRefundService_ProcessRefund_ExceedsRemaining tests refund exceeding remaining amount
func TestRefundService_ProcessRefund_ExceedsRemaining(t *testing.T) {
	ctx := context.Background()
	paymentID := uint64(1)

	mockPayments := new(MockPaymentRepository)
	mockRefunds := new(MockRefundRecordRepository)
	mockOrders := new(MockOrderRepository)

	payment := createTestPayment(paymentID, 1000, 100, model.PaymentStatusPaid, 10000, model.PaymentMethodWeChat)
	payment.RefundedAmountCents = 8000 // Only 2000 remaining

	mockPayments.On("Get", ctx, paymentID).Return(payment, nil)

	service := NewRefundService(mockPayments, mockRefunds, mockOrders)

	req := model.RefundRequest{
		PaymentID:   paymentID,
		AmountCents: 3000, // More than remaining 2000
		Reason:      "test",
	}

	resp, err := service.ProcessRefund(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "exceeds remaining refundable amount")

	mockPayments.AssertExpectations(t)
}

// TestRefundService_ProcessRefund_InvalidStatus tests refund with invalid payment status
func TestRefundService_ProcessRefund_InvalidStatus(t *testing.T) {
	ctx := context.Background()
	paymentID := uint64(1)

	tests := []struct {
		name          string
		paymentStatus model.PaymentStatus
		expectError   string
	}{
		{
			name:          "Pending payment",
			paymentStatus: model.PaymentStatusPending,
			expectError:   "payment status must be paid",
		},
		{
			name:          "Failed payment",
			paymentStatus: model.PaymentStatusFailed,
			expectError:   "payment status must be paid",
		},
		{
			name:          "Refunded payment",
			paymentStatus: model.PaymentStatusRefunded,
			expectError:   "payment status must be paid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payment := createTestPayment(paymentID, 1000, 100, tt.paymentStatus, 10000, model.PaymentMethodWeChat)

			mockPayments := new(MockPaymentRepository)
			mockRefunds := new(MockRefundRecordRepository)
			mockOrders := new(MockOrderRepository)

			mockPayments.On("Get", ctx, paymentID).Return(payment, nil)

			service := NewRefundService(mockPayments, mockRefunds, mockOrders)

			req := model.RefundRequest{
				PaymentID:   paymentID,
				AmountCents: 5000,
				Reason:      "test",
			}

			resp, err := service.ProcessRefund(ctx, req)

			assert.Error(t, err)
			assert.Nil(t, resp)
			assert.Contains(t, err.Error(), tt.expectError)

			mockPayments.AssertExpectations(t)
		})
	}
}

// TestRefundService_ProcessRefund_AlreadyFullyRefunded tests refund of already fully refunded payment
func TestRefundService_ProcessRefund_AlreadyFullyRefunded(t *testing.T) {
	ctx := context.Background()
	paymentID := uint64(1)

	mockPayments := new(MockPaymentRepository)
	mockRefunds := new(MockRefundRecordRepository)
	mockOrders := new(MockOrderRepository)

	payment := createTestPayment(paymentID, 1000, 100, model.PaymentStatusPaid, 10000, model.PaymentMethodWeChat)
	payment.RefundedAmountCents = 10000 // Already fully refunded

	mockPayments.On("Get", ctx, paymentID).Return(payment, nil)

	service := NewRefundService(mockPayments, mockRefunds, mockOrders)

	req := model.RefundRequest{
		PaymentID:   paymentID,
		AmountCents: 1000,
		Reason:      "test",
	}

	resp, err := service.ProcessRefund(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "already been fully refunded")

	mockPayments.AssertExpectations(t)
}

// TestRefundService_ProcessRefund_PaymentNotFound tests refund for non-existent payment
func TestRefundService_ProcessRefund_PaymentNotFound(t *testing.T) {
	ctx := context.Background()
	paymentID := uint64(999)

	mockPayments := new(MockPaymentRepository)
	mockRefunds := new(MockRefundRecordRepository)
	mockOrders := new(MockOrderRepository)

	mockPayments.On("Get", ctx, paymentID).Return(nil, repository.ErrNotFound)

	service := NewRefundService(mockPayments, mockRefunds, mockOrders)

	req := model.RefundRequest{
		PaymentID:   paymentID,
		AmountCents: 5000,
		Reason:      "test",
	}

	resp, err := service.ProcessRefund(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "payment not found")

	mockPayments.AssertExpectations(t)
}

// TestRefundService_ProcessRefund_ProviderFailure tests refund when provider fails
func TestRefundService_ProcessRefund_ProviderFailure(t *testing.T) {
	ctx := context.Background()
	paymentID := uint64(1)

	mockPayments := new(MockPaymentRepository)
	mockRefunds := new(MockRefundRecordRepository)
	mockOrders := new(MockOrderRepository)
	mockOpLogs := new(MockOperationLogRepository)

	payment := createTestPayment(paymentID, 1000, 100, model.PaymentStatusPaid, 10000, model.PaymentMethodWeChat)

	mockPayments.On("Get", ctx, paymentID).Return(payment, nil)
	mockRefunds.On("Create", ctx, mock.AnythingOfType("*model.RefundRecord")).Return(nil).Run(func(args mock.Arguments) {
		record := args.Get(1).(*model.RefundRecord)
		record.ID = 1
	})
	mockRefunds.On("Update", ctx, mock.MatchedBy(func(r *model.RefundRecord) bool {
		return r.Status == model.RefundStatusFailed
	})).Return(nil)
	mockOpLogs.On("Append", ctx, mock.AnythingOfType("*model.OperationLog")).Return(nil)

	service := NewRefundService(mockPayments, mockRefunds, mockOrders)
	service.SetOperationLogRepository(mockOpLogs)

	// Create a provider map with a failing provider
	providers := make(map[model.PaymentMethod]ProviderClient)
	providers[model.PaymentMethodWeChat] = &mockFailingProvider{}
	service.providers = providers

	req := model.RefundRequest{
		PaymentID:   paymentID,
		AmountCents: 5000,
		Reason:      "test",
	}

	resp, err := service.ProcessRefund(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "refund processing failed")

	mockPayments.AssertExpectations(t)
	mockRefunds.AssertExpectations(t)
	mockOpLogs.AssertExpectations(t)
}

// TestRefundService_GetRefundHistory_Success tests successful refund history retrieval
func TestRefundService_GetRefundHistory_Success(t *testing.T) {
	ctx := context.Background()
	paymentID := uint64(1)

	mockPayments := new(MockPaymentRepository)
	mockRefunds := new(MockRefundRecordRepository)
	mockOrders := new(MockOrderRepository)

	payment := createTestPayment(paymentID, 1000, 100, model.PaymentStatusPaid, 10000, model.PaymentMethodWeChat)

	refundRecords := []model.RefundRecord{
		{
			Base:        model.Base{ID: 1},
			PaymentID:   paymentID,
			OrderID:     1000,
			UserID:      100,
			AmountCents: 5000,
			Reason:      "partial refund",
			Status:      model.RefundStatusProcessed,
		},
		{
			Base:        model.Base{ID: 2},
			PaymentID:   paymentID,
			OrderID:     1000,
			UserID:      100,
			AmountCents: 3000,
			Reason:      "another partial refund",
			Status:      model.RefundStatusProcessed,
		},
	}

	mockPayments.On("Get", ctx, paymentID).Return(payment, nil)
	mockRefunds.On("ListByPaymentID", ctx, paymentID).Return(refundRecords, nil)

	service := NewRefundService(mockPayments, mockRefunds, mockOrders)

	records, err := service.GetRefundHistory(ctx, paymentID)

	require.NoError(t, err)
	assert.Len(t, records, 2)
	assert.Equal(t, uint64(1), records[0].ID)
	assert.Equal(t, uint64(2), records[1].ID)

	mockPayments.AssertExpectations(t)
	mockRefunds.AssertExpectations(t)
}

// TestRefundService_GetRefundHistory_PaymentNotFound tests refund history for non-existent payment
func TestRefundService_GetRefundHistory_PaymentNotFound(t *testing.T) {
	ctx := context.Background()
	paymentID := uint64(999)

	mockPayments := new(MockPaymentRepository)
	mockRefunds := new(MockRefundRecordRepository)
	mockOrders := new(MockOrderRepository)

	mockPayments.On("Get", ctx, paymentID).Return(nil, repository.ErrNotFound)

	service := NewRefundService(mockPayments, mockRefunds, mockOrders)

	records, err := service.GetRefundHistory(ctx, paymentID)

	assert.Error(t, err)
	assert.Nil(t, records)
	assert.Contains(t, err.Error(), "payment not found")

	mockPayments.AssertExpectations(t)
}

// TestRefundService_GetRefundsByOrder_Success tests successful refund retrieval by order
func TestRefundService_GetRefundsByOrder_Success(t *testing.T) {
	ctx := context.Background()
	orderID := uint64(1000)

	mockPayments := new(MockPaymentRepository)
	mockRefunds := new(MockRefundRecordRepository)
	mockOrders := new(MockOrderRepository)

	refundRecords := []model.RefundRecord{
		{
			Base:        model.Base{ID: 1},
			PaymentID:   1,
			OrderID:     orderID,
			UserID:      100,
			AmountCents: 5000,
			Reason:      "refund 1",
			Status:      model.RefundStatusProcessed,
		},
		{
			Base:        model.Base{ID: 2},
			PaymentID:   2,
			OrderID:     orderID,
			UserID:      100,
			AmountCents: 5000,
			Reason:      "refund 2",
			Status:      model.RefundStatusProcessed,
		},
	}

	mockRefunds.On("ListByOrderID", ctx, orderID).Return(refundRecords, nil)

	service := NewRefundService(mockPayments, mockRefunds, mockOrders)

	records, err := service.GetRefundsByOrder(ctx, orderID)

	require.NoError(t, err)
	assert.Len(t, records, 2)

	mockRefunds.AssertExpectations(t)
}

// TestPaymentModel_ValidateRefundAmount tests payment model refund validation
func TestPaymentModel_ValidateRefundAmount(t *testing.T) {
	tests := []struct {
		name        string
		payment     *model.Payment
		amountCents int64
		expectError bool
		errorCode   string
	}{
		{
			name: "Valid partial refund",
			payment: &model.Payment{
				Status:              model.PaymentStatusPaid,
				AmountCents:         10000,
				RefundedAmountCents: 0,
			},
			amountCents: 5000,
			expectError: false,
		},
		{
			name: "Valid full refund",
			payment: &model.Payment{
				Status:              model.PaymentStatusPaid,
				AmountCents:         10000,
				RefundedAmountCents: 0,
			},
			amountCents: 10000,
			expectError: false,
		},
		{
			name: "Zero amount",
			payment: &model.Payment{
				Status:              model.PaymentStatusPaid,
				AmountCents:         10000,
				RefundedAmountCents: 0,
			},
			amountCents: 0,
			expectError: true,
			errorCode:   model.RefundErrCodeInvalidAmount,
		},
		{
			name: "Negative amount",
			payment: &model.Payment{
				Status:              model.PaymentStatusPaid,
				AmountCents:         10000,
				RefundedAmountCents: 0,
			},
			amountCents: -100,
			expectError: true,
			errorCode:   model.RefundErrCodeInvalidAmount,
		},
		{
			name: "Exceeds remaining amount",
			payment: &model.Payment{
				Status:              model.PaymentStatusPaid,
				AmountCents:         10000,
				RefundedAmountCents: 8000,
			},
			amountCents: 3000, // Only 2000 remaining
			expectError: true,
			errorCode:   model.RefundErrCodeExceedsRemaining,
		},
		{
			name: "Invalid payment status - pending",
			payment: &model.Payment{
				Status:              model.PaymentStatusPending,
				AmountCents:         10000,
				RefundedAmountCents: 0,
			},
			amountCents: 5000,
			expectError: true,
			errorCode:   model.RefundErrCodeInvalidStatus,
		},
		{
			name: "Invalid payment status - failed",
			payment: &model.Payment{
				Status:              model.PaymentStatusFailed,
				AmountCents:         10000,
				RefundedAmountCents: 0,
			},
			amountCents: 5000,
			expectError: true,
			errorCode:   model.RefundErrCodeInvalidStatus,
		},
		{
			name: "Already fully refunded",
			payment: &model.Payment{
				Status:              model.PaymentStatusPaid,
				AmountCents:         10000,
				RefundedAmountCents: 10000,
			},
			amountCents: 1000,
			expectError: true,
			errorCode:   model.RefundErrCodeAlreadyRefunded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.payment.ValidateRefundAmount(tt.amountCents)

			if tt.expectError {
				assert.Error(t, err)
				refundErr, ok := err.(*model.RefundValidationError)
				require.True(t, ok, "error should be RefundValidationError type")
				assert.Equal(t, tt.errorCode, refundErr.Code)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestPaymentModel_RemainingRefundableAmount tests remaining refundable amount calculation
func TestPaymentModel_RemainingRefundableAmount(t *testing.T) {
	tests := []struct {
		name                    string
		payment                 *model.Payment
		expectedRemainingAmount int64
	}{
		{
			name: "No refunds yet",
			payment: &model.Payment{
				Status:              model.PaymentStatusPaid,
				AmountCents:         10000,
				RefundedAmountCents: 0,
			},
			expectedRemainingAmount: 10000,
		},
		{
			name: "Partial refund",
			payment: &model.Payment{
				Status:              model.PaymentStatusPaid,
				AmountCents:         10000,
				RefundedAmountCents: 3000,
			},
			expectedRemainingAmount: 7000,
		},
		{
			name: "Fully refunded",
			payment: &model.Payment{
				Status:              model.PaymentStatusPaid,
				AmountCents:         10000,
				RefundedAmountCents: 10000,
			},
			expectedRemainingAmount: 0,
		},
		{
			name: "Payment not paid",
			payment: &model.Payment{
				Status:              model.PaymentStatusPending,
				AmountCents:         10000,
				RefundedAmountCents: 0,
			},
			expectedRemainingAmount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.payment.RemainingRefundableAmount()
			assert.Equal(t, tt.expectedRemainingAmount, result)
		})
	}
}

// TestPaymentModel_CanRefund tests CanRefund method
func TestPaymentModel_CanRefund(t *testing.T) {
	tests := []struct {
		name     string
		payment  *model.Payment
		expected bool
	}{
		{
			name: "Paid with no refunds",
			payment: &model.Payment{
				Status:              model.PaymentStatusPaid,
				AmountCents:         10000,
				RefundedAmountCents: 0,
			},
			expected: true,
		},
		{
			name: "Paid with partial refund",
			payment: &model.Payment{
				Status:              model.PaymentStatusPaid,
				AmountCents:         10000,
				RefundedAmountCents: 5000,
			},
			expected: true,
		},
		{
			name: "Fully refunded",
			payment: &model.Payment{
				Status:              model.PaymentStatusPaid,
				AmountCents:         10000,
				RefundedAmountCents: 10000,
			},
			expected: false,
		},
		{
			name: "Pending payment",
			payment: &model.Payment{
				Status:              model.PaymentStatusPending,
				AmountCents:         10000,
				RefundedAmountCents: 0,
			},
			expected: false,
		},
		{
			name: "Failed payment",
			payment: &model.Payment{
				Status:              model.PaymentStatusFailed,
				AmountCents:         10000,
				RefundedAmountCents: 0,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.payment.CanRefund()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestPaymentModel_IsFullyRefunded tests IsFullyRefunded method
func TestPaymentModel_IsFullyRefunded(t *testing.T) {
	tests := []struct {
		name     string
		payment  *model.Payment
		expected bool
	}{
		{
			name: "No refunds",
			payment: &model.Payment{
				AmountCents:         10000,
				RefundedAmountCents: 0,
			},
			expected: false,
		},
		{
			name: "Partial refund",
			payment: &model.Payment{
				AmountCents:         10000,
				RefundedAmountCents: 5000,
			},
			expected: false,
		},
		{
			name: "Full refund",
			payment: &model.Payment{
				AmountCents:         10000,
				RefundedAmountCents: 10000,
			},
			expected: true,
		},
		{
			name: "Over-refunded (edge case)",
			payment: &model.Payment{
				AmountCents:         10000,
				RefundedAmountCents: 11000,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.payment.IsFullyRefunded()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestPaymentModel_IsPartiallyRefunded tests IsPartiallyRefunded method
func TestPaymentModel_IsPartiallyRefunded(t *testing.T) {
	tests := []struct {
		name     string
		payment  *model.Payment
		expected bool
	}{
		{
			name: "No refunds",
			payment: &model.Payment{
				AmountCents:         10000,
				RefundedAmountCents: 0,
			},
			expected: false,
		},
		{
			name: "Partial refund",
			payment: &model.Payment{
				AmountCents:         10000,
				RefundedAmountCents: 5000,
			},
			expected: true,
		},
		{
			name: "Full refund",
			payment: &model.Payment{
				AmountCents:         10000,
				RefundedAmountCents: 10000,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.payment.IsPartiallyRefunded()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestPaymentModel_IsValidStatusTransition tests valid payment status transitions
func TestPaymentModel_IsValidStatusTransition(t *testing.T) {
	tests := []struct {
		name     string
		from     model.PaymentStatus
		to       model.PaymentStatus
		expected bool
	}{
		{
			name:     "Pending to Paid - valid",
			from:     model.PaymentStatusPending,
			to:       model.PaymentStatusPaid,
			expected: true,
		},
		{
			name:     "Pending to Failed - valid",
			from:     model.PaymentStatusPending,
			to:       model.PaymentStatusFailed,
			expected: true,
		},
		{
			name:     "Paid to Refunded - valid",
			from:     model.PaymentStatusPaid,
			to:       model.PaymentStatusRefunded,
			expected: true,
		},
		{
			name:     "Pending to Refunded - invalid",
			from:     model.PaymentStatusPending,
			to:       model.PaymentStatusRefunded,
			expected: false,
		},
		{
			name:     "Failed to Paid - invalid",
			from:     model.PaymentStatusFailed,
			to:       model.PaymentStatusPaid,
			expected: false,
		},
		{
			name:     "Refunded to Paid - invalid",
			from:     model.PaymentStatusRefunded,
			to:       model.PaymentStatusPaid,
			expected: false,
		},
		{
			name:     "Same status - valid",
			from:     model.PaymentStatusPaid,
			to:       model.PaymentStatusPaid,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := model.IsValidStatusTransition(tt.from, tt.to)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestPaymentModel_ValidateStatusTransition tests payment status transition validation
func TestPaymentModel_ValidateStatusTransition(t *testing.T) {
	tests := []struct {
		name        string
		payment     *model.Payment
		newStatus   model.PaymentStatus
		expectError bool
	}{
		{
			name: "Valid transition - pending to paid",
			payment: &model.Payment{
				Status: model.PaymentStatusPending,
			},
			newStatus:   model.PaymentStatusPaid,
			expectError: false,
		},
		{
			name: "Valid transition - paid to refunded",
			payment: &model.Payment{
				Status: model.PaymentStatusPaid,
			},
			newStatus:   model.PaymentStatusRefunded,
			expectError: false,
		},
		{
			name: "Invalid transition - pending to refunded",
			payment: &model.Payment{
				Status: model.PaymentStatusPending,
			},
			newStatus:   model.PaymentStatusRefunded,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.payment.ValidateStatusTransition(tt.newStatus)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestPaymentModel_GetAllowedTransitions tests getting allowed status transitions
func TestPaymentModel_GetAllowedTransitions(t *testing.T) {
	tests := []struct {
		name                string
		payment             *model.Payment
		expectedTransitions []model.PaymentStatus
	}{
		{
			name: "Pending payment transitions",
			payment: &model.Payment{
				Status: model.PaymentStatusPending,
			},
			expectedTransitions: []model.PaymentStatus{
				model.PaymentStatusPaid,
				model.PaymentStatusFailed,
			},
		},
		{
			name: "Paid payment transitions",
			payment: &model.Payment{
				Status: model.PaymentStatusPaid,
			},
			expectedTransitions: []model.PaymentStatus{
				model.PaymentStatusRefunded,
			},
		},
		{
			name: "Failed payment transitions",
			payment: &model.Payment{
				Status: model.PaymentStatusFailed,
			},
			expectedTransitions: []model.PaymentStatus{},
		},
		{
			name: "Refunded payment transitions",
			payment: &model.Payment{
				Status: model.PaymentStatusRefunded,
			},
			expectedTransitions: []model.PaymentStatus{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transitions := tt.payment.GetAllowedTransitions()
			assert.ElementsMatch(t, tt.expectedTransitions, transitions)
		})
	}
}

// Mock helpers

type mockFailingProvider struct{}

func (m *mockFailingProvider) Refund(ctx context.Context, p *model.Payment, reason string) (string, json.RawMessage, time.Time, error) {
	return "", nil, time.Time{}, errors.New("provider error")
}
