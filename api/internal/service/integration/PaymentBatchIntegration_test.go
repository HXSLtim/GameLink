// Package integration provides integration tests for admin payment batch operations.
package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/service/admin"
)

// ============================================================================
// BatchCapture Tests - 批量收款
// ============================================================================

// TestPaymentBatch_BatchCapture_Success tests successful batch capture of pending payments.
func TestPaymentBatch_BatchCapture_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "capture_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "capture_player"))
	testGame := CreateTestGame(t, db, "capture_game")

	// Create pending payments
	var paymentIDs []uint64
	for i := 0; i < 3; i++ {
		order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)
		payment := CreateTestPayment(t, db, order, model.PaymentStatusPending)
		paymentIDs = append(paymentIDs, payment.ID)
	}

	// Batch capture
	req := admin.BatchCaptureRequest{
		PaymentIDs: paymentIDs,
	}
	result, err := svc.BatchCapture(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Empty(t, result.FailedIDs)
	assert.Empty(t, result.Errors)

	// Verify database state
	for _, id := range paymentIDs {
		var payment model.Payment
		err := db.First(&payment, id).Error
		require.NoError(t, err)
		assert.Equal(t, model.PaymentStatusPaid, payment.Status)
		assert.NotNil(t, payment.PaidAt)
		assert.NotEmpty(t, payment.ProviderTradeNo)
	}
}

// TestPaymentBatch_BatchCapture_InvalidStatus tests capture with invalid payment status.
func TestPaymentBatch_BatchCapture_InvalidStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	testUser := CreateUniqueTestUser(t, db, "capture_invalid_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "capture_invalid_player"))
	testGame := CreateTestGame(t, db, "capture_invalid_game")

	// Create payments with different statuses
	var paymentIDs []uint64

	// pending - should succeed
	order1 := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)
	payment1 := CreateTestPayment(t, db, order1, model.PaymentStatusPending)
	paymentIDs = append(paymentIDs, payment1.ID)

	// already paid - should fail
	order2 := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusConfirmed, 10000)
	payment2 := CreateTestPayment(t, db, order2, model.PaymentStatusPaid)
	paymentIDs = append(paymentIDs, payment2.ID)

	// failed - should fail
	order3 := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCanceled, 10000)
	payment3 := CreateTestPayment(t, db, order3, model.PaymentStatusFailed)
	paymentIDs = append(paymentIDs, payment3.ID)

	// refunded - should fail
	order4 := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusRefunded, 10000)
	payment4 := CreateTestPayment(t, db, order4, model.PaymentStatusRefunded)
	paymentIDs = append(paymentIDs, payment4.ID)

	req := admin.BatchCaptureRequest{
		PaymentIDs: paymentIDs,
	}
	result, err := svc.BatchCapture(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 3, result.FailedCount)
	assert.Len(t, result.FailedIDs, 3)
	assert.Len(t, result.Errors, 3)

	// Verify error messages
	for _, e := range result.Errors {
		assert.Contains(t, e.Message, "invalid status for capture")
	}

	// Verify only pending payment was updated
	var payment model.Payment
	db.First(&payment, payment1.ID)
	assert.Equal(t, model.PaymentStatusPaid, payment.Status)

	db.First(&payment, payment2.ID)
	assert.Equal(t, model.PaymentStatusPaid, payment.Status) // unchanged

	db.First(&payment, payment3.ID)
	assert.Equal(t, model.PaymentStatusFailed, payment.Status) // unchanged

	db.First(&payment, payment4.ID)
	assert.Equal(t, model.PaymentStatusRefunded, payment.Status) // unchanged
}

// TestPaymentBatch_BatchCapture_NotFound tests capture with non-existent payments.
func TestPaymentBatch_BatchCapture_NotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	testUser := CreateUniqueTestUser(t, db, "capture_notfound_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "capture_notfound_player"))
	testGame := CreateTestGame(t, db, "capture_notfound_game")

	// Create one valid payment
	order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)
	payment := CreateTestPayment(t, db, order, model.PaymentStatusPending)

	// Include non-existent IDs
	paymentIDs := []uint64{payment.ID, 99999, 88888}

	req := admin.BatchCaptureRequest{
		PaymentIDs: paymentIDs,
	}
	result, err := svc.BatchCapture(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 2, result.FailedCount)
	assert.Len(t, result.FailedIDs, 2)
	assert.Len(t, result.Errors, 2)

	// Verify error messages
	for _, e := range result.Errors {
		assert.Contains(t, e.Message, "not found")
	}
}

// TestPaymentBatch_BatchCapture_WithCustomTradeNo tests capture with custom provider trade no.
func TestPaymentBatch_BatchCapture_WithCustomTradeNo(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	testUser := CreateUniqueTestUser(t, db, "capture_trade_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "capture_trade_player"))
	testGame := CreateTestGame(t, db, "capture_trade_game")

	order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)
	payment := CreateTestPayment(t, db, order, model.PaymentStatusPending)

	customTradeNo := "WX_20250129_1234567890"
	req := admin.BatchCaptureRequest{
		PaymentIDs:      []uint64{payment.ID},
		ProviderTradeNo: customTradeNo,
	}
	result, err := svc.BatchCapture(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 1, result.SuccessCount)

	var updatedPayment model.Payment
	db.First(&updatedPayment, payment.ID)
	assert.Equal(t, customTradeNo, updatedPayment.ProviderTradeNo)
}

// TestPaymentBatch_BatchCapture_WithCustomPaidAt tests capture with custom paid at time.
func TestPaymentBatch_BatchCapture_WithCustomPaidAt(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	testUser := CreateUniqueTestUser(t, db, "capture_time_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "capture_time_player"))
	testGame := CreateTestGame(t, db, "capture_time_game")

	order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)
	payment := CreateTestPayment(t, db, order, model.PaymentStatusPending)

	customTime := time.Now().Add(-1 * time.Hour)
	req := admin.BatchCaptureRequest{
		PaymentIDs: []uint64{payment.ID},
		PaidAt:     &customTime,
	}
	result, err := svc.BatchCapture(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 1, result.SuccessCount)

	var updatedPayment model.Payment
	db.First(&updatedPayment, payment.ID)
	assert.NotNil(t, updatedPayment.PaidAt)
	// Allow small time difference
	assert.WithinDuration(t, customTime, *updatedPayment.PaidAt, time.Second)
}

// TestPaymentBatch_BatchCapture_EmptyIDs tests capture with empty payment IDs.
func TestPaymentBatch_BatchCapture_EmptyIDs(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	req := admin.BatchCaptureRequest{
		PaymentIDs: []uint64{},
	}
	result, err := svc.BatchCapture(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "payment ids cannot be empty")
}

// TestPaymentBatch_BatchCapture_TooManyIDs tests capture with too many payment IDs.
func TestPaymentBatch_BatchCapture_TooManyIDs(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	// Create 501 payment IDs
	paymentIDs := make([]uint64, 501)
	for i := 0; i < 501; i++ {
		paymentIDs[i] = uint64(i + 1)
	}

	req := admin.BatchCaptureRequest{
		PaymentIDs: paymentIDs,
	}
	result, err := svc.BatchCapture(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "maximum 500 payments allowed")
}

// ============================================================================
// BatchRefund Tests - 批量退款
// ============================================================================

// TestPaymentBatch_BatchRefund_Success tests successful batch refund of paid payments.
func TestPaymentBatch_BatchRefund_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	testUser := CreateUniqueTestUser(t, db, "refund_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "refund_player"))
	testGame := CreateTestGame(t, db, "refund_game")

	// Create paid payments
	var paymentIDs []uint64
	for i := 0; i < 3; i++ {
		order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusConfirmed, 10000)
		payment := CreateTestPayment(t, db, order, model.PaymentStatusPaid)
		paymentIDs = append(paymentIDs, payment.ID)
	}

	req := admin.BatchRefundRequest{
		PaymentIDs: paymentIDs,
		Reason:     "Customer requested refund",
	}
	result, err := svc.BatchRefund(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Empty(t, result.FailedIDs)
	assert.Empty(t, result.Errors)

	// Verify database state
	for _, id := range paymentIDs {
		var payment model.Payment
		err := db.First(&payment, id).Error
		require.NoError(t, err)
		assert.Equal(t, model.PaymentStatusRefunded, payment.Status)
		assert.NotNil(t, payment.RefundedAt)
		assert.Equal(t, int64(10000), payment.RefundedAmountCents)
		assert.NotEmpty(t, payment.ProviderTradeNo)

		// Verify associated order status
		var order model.Order
		err = db.First(&order, payment.OrderID).Error
		require.NoError(t, err)
		assert.Equal(t, model.OrderStatusRefunded, order.Status)
		assert.Equal(t, int64(10000), order.RefundAmountCents)
		assert.Equal(t, "Customer requested refund", order.RefundReason)
	}
}

// TestPaymentBatch_BatchRefund_InvalidStatus tests refund with invalid payment status.
func TestPaymentBatch_BatchRefund_InvalidStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	testUser := CreateUniqueTestUser(t, db, "refund_invalid_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "refund_invalid_player"))
	testGame := CreateTestGame(t, db, "refund_invalid_game")

	var paymentIDs []uint64

	// paid - should succeed
	order1 := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusConfirmed, 10000)
	payment1 := CreateTestPayment(t, db, order1, model.PaymentStatusPaid)
	paymentIDs = append(paymentIDs, payment1.ID)

	// pending - should fail
	order2 := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)
	payment2 := CreateTestPayment(t, db, order2, model.PaymentStatusPending)
	paymentIDs = append(paymentIDs, payment2.ID)

	// failed - should fail
	order3 := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCanceled, 10000)
	payment3 := CreateTestPayment(t, db, order3, model.PaymentStatusFailed)
	paymentIDs = append(paymentIDs, payment3.ID)

	req := admin.BatchRefundRequest{
		PaymentIDs: paymentIDs,
		Reason:     "Test refund",
	}
	result, err := svc.BatchRefund(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 2, result.FailedCount)

	// Verify error messages
	for _, e := range result.Errors {
		assert.Contains(t, e.Message, "invalid status for refund")
	}
}

// TestPaymentBatch_BatchRefund_AlreadyRefunded tests refund with already fully refunded payments.
func TestPaymentBatch_BatchRefund_AlreadyRefunded(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	testUser := CreateUniqueTestUser(t, db, "refund_already_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "refund_already_player"))
	testGame := CreateTestGame(t, db, "refund_already_game")

	var paymentIDs []uint64

	// paid - should succeed
	order1 := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusConfirmed, 10000)
	payment1 := CreateTestPayment(t, db, order1, model.PaymentStatusPaid)
	paymentIDs = append(paymentIDs, payment1.ID)

	// already refunded - should fail
	order2 := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusRefunded, 10000)
	payment2 := CreateTestPayment(t, db, order2, model.PaymentStatusRefunded)
	payment2.RefundedAmountCents = 10000 // fully refunded
	db.Save(payment2)
	paymentIDs = append(paymentIDs, payment2.ID)

	req := admin.BatchRefundRequest{
		PaymentIDs: paymentIDs,
		Reason:     "Test refund",
	}
	result, err := svc.BatchRefund(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)

	// Verify error message for already refunded
	var foundAlreadyRefunded bool
	for _, e := range result.Errors {
		if e.Message == "payment is already fully refunded" {
			foundAlreadyRefunded = true
			break
		}
	}
	assert.True(t, foundAlreadyRefunded, "Expected 'already fully refunded' error message")
}

// TestPaymentBatch_BatchRefund_NotFound tests refund with non-existent payments.
func TestPaymentBatch_BatchRefund_NotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	testUser := CreateUniqueTestUser(t, db, "refund_notfound_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "refund_notfound_player"))
	testGame := CreateTestGame(t, db, "refund_notfound_game")

	order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusConfirmed, 10000)
	payment := CreateTestPayment(t, db, order, model.PaymentStatusPaid)

	paymentIDs := []uint64{payment.ID, 99999, 88888}

	req := admin.BatchRefundRequest{
		PaymentIDs: paymentIDs,
		Reason:     "Test refund",
	}
	result, err := svc.BatchRefund(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 2, result.FailedCount)

	// Verify error messages
	for _, e := range result.Errors {
		assert.Contains(t, e.Message, "not found")
	}
}

// TestPaymentBatch_BatchRefund_WithCustomRefundedAt tests refund with custom refunded at time.
func TestPaymentBatch_BatchRefund_WithCustomRefundedAt(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	testUser := CreateUniqueTestUser(t, db, "refund_time_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "refund_time_player"))
	testGame := CreateTestGame(t, db, "refund_time_game")

	order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusConfirmed, 10000)
	payment := CreateTestPayment(t, db, order, model.PaymentStatusPaid)

	customTime := time.Now().Add(-2 * time.Hour)
	req := admin.BatchRefundRequest{
		PaymentIDs: []uint64{payment.ID},
		Reason:     "Test refund",
		RefundedAt: &customTime,
	}
	result, err := svc.BatchRefund(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 1, result.SuccessCount)

	var updatedPayment model.Payment
	db.First(&updatedPayment, payment.ID)
	assert.NotNil(t, updatedPayment.RefundedAt)
	assert.WithinDuration(t, customTime, *updatedPayment.RefundedAt, time.Second)

	// Verify order also has same RefundedAt
	var updatedOrder model.Order
	db.First(&updatedOrder, order.ID)
	assert.NotNil(t, updatedOrder.RefundedAt)
	assert.WithinDuration(t, customTime, *updatedOrder.RefundedAt, time.Second)
}

// TestPaymentBatch_BatchRefund_EmptyReason tests refund with empty reason.
func TestPaymentBatch_BatchRefund_EmptyReason(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	testUser := CreateUniqueTestUser(t, db, "refund_empty_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "refund_empty_player"))
	testGame := CreateTestGame(t, db, "refund_empty_game")

	order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusConfirmed, 10000)
	payment := CreateTestPayment(t, db, order, model.PaymentStatusPaid)

	req := admin.BatchRefundRequest{
		PaymentIDs: []uint64{payment.ID},
		Reason:     "", // empty reason
	}
	result, err := svc.BatchRefund(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "refund reason is required")
}

// TestPaymentBatch_BatchRefund_WhitespaceReason tests refund with whitespace-only reason.
func TestPaymentBatch_BatchRefund_WhitespaceReason(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	testUser := CreateUniqueTestUser(t, db, "refund_whitespace_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "refund_whitespace_player"))
	testGame := CreateTestGame(t, db, "refund_whitespace_game")

	order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusConfirmed, 10000)
	payment := CreateTestPayment(t, db, order, model.PaymentStatusPaid)

	req := admin.BatchRefundRequest{
		PaymentIDs: []uint64{payment.ID},
		Reason:     "   ", // whitespace only
	}
	result, err := svc.BatchRefund(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "refund reason is required")
}

// TestPaymentBatch_BatchRefund_TooLongReason tests refund with reason exceeding max length.
func TestPaymentBatch_BatchRefund_TooLongReason(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	testUser := CreateUniqueTestUser(t, db, "refund_long_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "refund_long_player"))
	testGame := CreateTestGame(t, db, "refund_long_game")

	order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusConfirmed, 10000)
	payment := CreateTestPayment(t, db, order, model.PaymentStatusPaid)

	// Create reason longer than 500 characters
	longReason := string(make([]byte, 501))
	for i := range longReason {
		longReason = longReason[:i] + "A"
	}

	req := admin.BatchRefundRequest{
		PaymentIDs: []uint64{payment.ID},
		Reason:     longReason,
	}
	result, err := svc.BatchRefund(ctx, req)

	// The binding in handler should catch this, but service validates too
	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================================================
// BatchCancel Tests - 批量取消支付
// ============================================================================

// TestPaymentBatch_BatchCancel_Success tests successful batch cancel of pending payments.
func TestPaymentBatch_BatchCancel_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	testUser := CreateUniqueTestUser(t, db, "cancel_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "cancel_player"))
	testGame := CreateTestGame(t, db, "cancel_game")

	// Create pending payments
	var paymentIDs []uint64
	for i := 0; i < 3; i++ {
		order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)
		payment := CreateTestPayment(t, db, order, model.PaymentStatusPending)
		paymentIDs = append(paymentIDs, payment.ID)
	}

	req := admin.BatchCancelRequest{
		PaymentIDs: paymentIDs,
	}
	result, err := svc.BatchCancel(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Empty(t, result.FailedIDs)
	assert.Empty(t, result.Errors)

	// Verify database state
	for _, id := range paymentIDs {
		var payment model.Payment
		err := db.First(&payment, id).Error
		require.NoError(t, err)
		assert.Equal(t, model.PaymentStatusFailed, payment.Status)
	}
}

// TestPaymentBatch_BatchCancel_InvalidStatus tests cancel with invalid payment status.
func TestPaymentBatch_BatchCancel_InvalidStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	testUser := CreateUniqueTestUser(t, db, "cancel_invalid_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "cancel_invalid_player"))
	testGame := CreateTestGame(t, db, "cancel_invalid_game")

	var paymentIDs []uint64

	// pending - should succeed
	order1 := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)
	payment1 := CreateTestPayment(t, db, order1, model.PaymentStatusPending)
	paymentIDs = append(paymentIDs, payment1.ID)

	// paid - should fail
	order2 := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusConfirmed, 10000)
	payment2 := CreateTestPayment(t, db, order2, model.PaymentStatusPaid)
	paymentIDs = append(paymentIDs, payment2.ID)

	// already failed - should fail
	order3 := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCanceled, 10000)
	payment3 := CreateTestPayment(t, db, order3, model.PaymentStatusFailed)
	paymentIDs = append(paymentIDs, payment3.ID)

	// refunded - should fail
	order4 := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusRefunded, 10000)
	payment4 := CreateTestPayment(t, db, order4, model.PaymentStatusRefunded)
	paymentIDs = append(paymentIDs, payment4.ID)

	req := admin.BatchCancelRequest{
		PaymentIDs: paymentIDs,
	}
	result, err := svc.BatchCancel(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 3, result.FailedCount)

	// Verify error messages
	for _, e := range result.Errors {
		assert.Contains(t, e.Message, "invalid status for cancel")
	}
}

// TestPaymentBatch_BatchCancel_NotFound tests cancel with non-existent payments.
func TestPaymentBatch_BatchCancel_NotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	testUser := CreateUniqueTestUser(t, db, "cancel_notfound_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "cancel_notfound_player"))
	testGame := CreateTestGame(t, db, "cancel_notfound_game")

	order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)
	payment := CreateTestPayment(t, db, order, model.PaymentStatusPending)

	paymentIDs := []uint64{payment.ID, 99999, 88888}

	req := admin.BatchCancelRequest{
		PaymentIDs: paymentIDs,
	}
	result, err := svc.BatchCancel(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 2, result.FailedCount)

	// Verify error messages
	for _, e := range result.Errors {
		assert.Contains(t, e.Message, "not found")
	}
}

// TestPaymentBatch_BatchCancel_EmptyIDs tests cancel with empty payment IDs.
func TestPaymentBatch_BatchCancel_EmptyIDs(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	req := admin.BatchCancelRequest{
		PaymentIDs: []uint64{},
	}
	result, err := svc.BatchCancel(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "payment ids cannot be empty")
}

// TestPaymentBatch_BatchCancel_TooManyIDs tests cancel with too many payment IDs.
func TestPaymentBatch_BatchCancel_TooManyIDs(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	paymentIDs := make([]uint64, 501)
	for i := 0; i < 501; i++ {
		paymentIDs[i] = uint64(i + 1)
	}

	req := admin.BatchCancelRequest{
		PaymentIDs: paymentIDs,
	}
	result, err := svc.BatchCancel(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "maximum 500 payments allowed")
}

// ============================================================================
// BatchUpdateStatus Tests - 批量更新支付状态
// ============================================================================

// TestPaymentBatch_BatchUpdateStatus_Success tests successful batch status update.
func TestPaymentBatch_BatchUpdateStatus_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	testUser := CreateUniqueTestUser(t, db, "update_status_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "update_status_player"))
	testGame := CreateTestGame(t, db, "update_status_game")

	// Create pending payments
	var paymentIDs []uint64
	for i := 0; i < 3; i++ {
		order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)
		payment := CreateTestPayment(t, db, order, model.PaymentStatusPending)
		paymentIDs = append(paymentIDs, payment.ID)
	}

	req := admin.BatchUpdateStatusRequest{
		PaymentIDs: paymentIDs,
		Status:     model.PaymentStatusPaid,
	}
	result, err := svc.BatchUpdateStatus(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)

	// Verify database state
	for _, id := range paymentIDs {
		var payment model.Payment
		err := db.First(&payment, id).Error
		require.NoError(t, err)
		assert.Equal(t, model.PaymentStatusPaid, payment.Status)
		assert.NotNil(t, payment.PaidAt) // Should be set automatically
	}
}

// TestPaymentBatch_BatchUpdateStatus_InvalidTransition tests status update with invalid transitions.
func TestPaymentBatch_BatchUpdateStatus_InvalidTransition(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	testUser := CreateUniqueTestUser(t, db, "update_invalid_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "update_invalid_player"))
	testGame := CreateTestGame(t, db, "update_invalid_game")

	var paymentIDs []uint64

	// pending -> paid (valid)
	order1 := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)
	payment1 := CreateTestPayment(t, db, order1, model.PaymentStatusPending)
	paymentIDs = append(paymentIDs, payment1.ID)

	// paid -> refunded (valid)
	order2 := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusConfirmed, 10000)
	payment2 := CreateTestPayment(t, db, order2, model.PaymentStatusPaid)
	paymentIDs = append(paymentIDs, payment2.ID)

	// failed -> paid (invalid, failed is terminal)
	order3 := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCanceled, 10000)
	payment3 := CreateTestPayment(t, db, order3, model.PaymentStatusFailed)
	paymentIDs = append(paymentIDs, payment3.ID)

	// refunded -> paid (invalid, refunded is terminal)
	order4 := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusRefunded, 10000)
	payment4 := CreateTestPayment(t, db, order4, model.PaymentStatusRefunded)
	paymentIDs = append(paymentIDs, payment4.ID)

	req := admin.BatchUpdateStatusRequest{
		PaymentIDs: paymentIDs,
		Status:     model.PaymentStatusPaid,
	}
	result, err := svc.BatchUpdateStatus(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 2, result.SuccessCount)
	assert.Equal(t, 2, result.FailedCount)

	// Verify error messages for invalid transitions
	for _, e := range result.Errors {
		assert.Contains(t, e.Message, "invalid status transition")
	}
}

// TestPaymentBatch_BatchUpdateStatus_NotFound tests status update with non-existent payments.
func TestPaymentBatch_BatchUpdateStatus_NotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	testUser := CreateUniqueTestUser(t, db, "update_notfound_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "update_notfound_player"))
	testGame := CreateTestGame(t, db, "update_notfound_game")

	order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)
	payment := CreateTestPayment(t, db, order, model.PaymentStatusPending)

	paymentIDs := []uint64{payment.ID, 99999, 88888}

	req := admin.BatchUpdateStatusRequest{
		PaymentIDs: paymentIDs,
		Status:     model.PaymentStatusPaid,
	}
	result, err := svc.BatchUpdateStatus(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 2, result.FailedCount)

	// Verify error messages
	for _, e := range result.Errors {
		assert.Contains(t, e.Message, "not found")
	}
}

// TestPaymentBatch_BatchUpdateStatus_ToRefunded tests status update to refunded.
func TestPaymentBatch_BatchUpdateStatus_ToRefunded(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	testUser := CreateUniqueTestUser(t, db, "update_refunded_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "update_refunded_player"))
	testGame := CreateTestGame(t, db, "update_refunded_game")

	// Create paid payment
	order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusConfirmed, 10000)
	payment := CreateTestPayment(t, db, order, model.PaymentStatusPaid)

	req := admin.BatchUpdateStatusRequest{
		PaymentIDs: []uint64{payment.ID},
		Status:     model.PaymentStatusRefunded,
	}
	result, err := svc.BatchUpdateStatus(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 1, result.SuccessCount)

	var updatedPayment model.Payment
	db.First(&updatedPayment, payment.ID)
	assert.Equal(t, model.PaymentStatusRefunded, updatedPayment.Status)
	assert.NotNil(t, updatedPayment.RefundedAt)
	assert.Equal(t, int64(10000), updatedPayment.RefundedAmountCents)
}

// TestPaymentBatch_BatchUpdateStatus_ToFailed tests status update to failed.
func TestPaymentBatch_BatchUpdateStatus_ToFailed(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	testUser := CreateUniqueTestUser(t, db, "update_failed_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "update_failed_player"))
	testGame := CreateTestGame(t, db, "update_failed_game")

	// Create pending payment
	order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)
	payment := CreateTestPayment(t, db, order, model.PaymentStatusPending)

	req := admin.BatchUpdateStatusRequest{
		PaymentIDs: []uint64{payment.ID},
		Status:     model.PaymentStatusFailed,
	}
	result, err := svc.BatchUpdateStatus(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 1, result.SuccessCount)

	var updatedPayment model.Payment
	db.First(&updatedPayment, payment.ID)
	assert.Equal(t, model.PaymentStatusFailed, updatedPayment.Status)
}

// TestPaymentBatch_BatchUpdateStatus_EmptyIDs tests status update with empty payment IDs.
func TestPaymentBatch_BatchUpdateStatus_EmptyIDs(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	req := admin.BatchUpdateStatusRequest{
		PaymentIDs: []uint64{},
		Status:     model.PaymentStatusPaid,
	}
	result, err := svc.BatchUpdateStatus(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "payment ids cannot be empty")
}

// TestPaymentBatch_BatchUpdateStatus_InvalidStatus tests status update with invalid status.
func TestPaymentBatch_BatchUpdateStatus_InvalidStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	testUser := CreateUniqueTestUser(t, db, "update_status_invalid_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "update_status_invalid_player"))
	testGame := CreateTestGame(t, db, "update_status_invalid_game")

	order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)
	payment := CreateTestPayment(t, db, order, model.PaymentStatusPending)

	// Use invalid status
	req := admin.BatchUpdateStatusRequest{
		PaymentIDs: []uint64{payment.ID},
		Status:     model.PaymentStatus("invalid_status"),
	}
	result, err := svc.BatchUpdateStatus(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid payment status")
}

// TestPaymentBatch_BatchUpdateStatus_SameStatus tests status update to same status (no-op).
func TestPaymentBatch_BatchUpdateStatus_SameStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	testUser := CreateUniqueTestUser(t, db, "update_same_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "update_same_player"))
	testGame := CreateTestGame(t, db, "update_same_game")

	order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)
	payment := CreateTestPayment(t, db, order, model.PaymentStatusPending)

	// Update to same status
	req := admin.BatchUpdateStatusRequest{
		PaymentIDs: []uint64{payment.ID},
		Status:     model.PaymentStatusPending,
	}
	result, err := svc.BatchUpdateStatus(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)

	// Verify no changes
	var updatedPayment model.Payment
	db.First(&updatedPayment, payment.ID)
	assert.Equal(t, model.PaymentStatusPending, updatedPayment.Status)
}

// TestPaymentBatch_BatchUpdateStatus_TooManyIDs tests status update with too many payment IDs.
func TestPaymentBatch_BatchUpdateStatus_TooManyIDs(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	paymentIDs := make([]uint64, 501)
	for i := 0; i < 501; i++ {
		paymentIDs[i] = uint64(i + 1)
	}

	req := admin.BatchUpdateStatusRequest{
		PaymentIDs: paymentIDs,
		Status:     model.PaymentStatusPaid,
	}
	result, err := svc.BatchUpdateStatus(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "maximum 500 payments allowed")
}

// ============================================================================
// Edge Cases Tests
// ============================================================================

// TestPaymentBatch_MixedScenario tests batch operations with mixed scenarios.
func TestPaymentBatch_MixedScenario(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	testUser := CreateUniqueTestUser(t, db, "mixed_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "mixed_player"))
	testGame := CreateTestGame(t, db, "mixed_game")

	// Create various payment states
	var paymentIDs []uint64

	// pending (valid for capture)
	order1 := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)
	payment1 := CreateTestPayment(t, db, order1, model.PaymentStatusPending)
	paymentIDs = append(paymentIDs, payment1.ID)

	// pending (valid for capture)
	order2 := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)
	payment2 := CreateTestPayment(t, db, order2, model.PaymentStatusPending)
	paymentIDs = append(paymentIDs, payment2.ID)

	// non-existent
	paymentIDs = append(paymentIDs, 99999)

	// paid (invalid for capture)
	order3 := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusConfirmed, 10000)
	payment3 := CreateTestPayment(t, db, order3, model.PaymentStatusPaid)
	paymentIDs = append(paymentIDs, payment3.ID)

	req := admin.BatchCaptureRequest{
		PaymentIDs: paymentIDs,
	}
	result, err := svc.BatchCapture(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 2, result.SuccessCount)
	assert.Equal(t, 2, result.FailedCount)
	assert.Len(t, result.FailedIDs, 2)
	assert.Len(t, result.Errors, 2)
}

// TestPaymentBatch_LargeBatch tests handling of large batches (within limits).
func TestPaymentBatch_LargeBatch(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()
	svc := setupOrderBatchAdminService(t, db)

	testUser := CreateUniqueTestUser(t, db, "large_batch_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "large_batch_player"))
	testGame := CreateTestGame(t, db, "large_batch_game")

	// Create 100 pending payments
	var paymentIDs []uint64
	for i := 0; i < 100; i++ {
		order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)
		payment := CreateTestPayment(t, db, order, model.PaymentStatusPending)
		paymentIDs = append(paymentIDs, payment.ID)
	}

	req := admin.BatchCaptureRequest{
		PaymentIDs: paymentIDs,
	}
	result, err := svc.BatchCapture(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 100, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)

	// Verify all were updated
	for _, id := range paymentIDs {
		var payment model.Payment
		err := db.First(&payment, id).Error
		require.NoError(t, err)
		assert.Equal(t, model.PaymentStatusPaid, payment.Status)
	}
}
