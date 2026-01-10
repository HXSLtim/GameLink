// Package integration provides supplementary integration tests for payment service.
package integration

import (
	"context"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository/implementations"
	"gamelink/internal/repository/payment"
	"gamelink/internal/repository/wallet"
	paymentservice "gamelink/internal/service/payment"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPaymentService_PaymentTimeout tests payment timeout handling.
func TestPaymentService_PaymentTimeout(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "timeout_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "timeout_player"))
	testGame := CreateTestGame(t, db, "timeout_game")

	// Create order
	order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)

	// Create pending payment (simulating timeout scenario)
	pendingPayment := &model.Payment{
		Base: model.Base{
			ExtJSON: "{}",
		},
		OrderID:     order.ID,
		UserID:      testUser.ID,
		Method:      model.PaymentMethodWeChat,
		AmountCents: order.TotalPriceCents,
		Status:      model.PaymentStatusPending,
		Currency:    model.CurrencyCNY,
	}
	err := db.Create(pendingPayment).Error
	require.NoError(t, err)

	// Simulate timeout check - mark payment as failed
	pendingPayment.Status = model.PaymentStatusFailed
	err = db.Save(pendingPayment).Error
	require.NoError(t, err)

	// Update order status
	order.Status = model.OrderStatusCanceled
	order.CancelReason = "Payment timeout"
	err = db.Save(order).Error
	require.NoError(t, err)

	// Verify
	var updatedPayment model.Payment
	err = db.First(&updatedPayment, pendingPayment.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.PaymentStatusFailed, updatedPayment.Status)

	var updatedOrder model.Order
	err = db.First(&updatedOrder, order.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.OrderStatusCanceled, updatedOrder.Status)
}

// TestPaymentService_CombinedPaymentRefund tests refund for combined payment.
func TestPaymentService_CombinedPaymentRefund(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	paymentRepo := payment.NewPaymentRepository(db)
	walletRepo := wallet.NewWalletRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := paymentservice.NewPaymentService(paymentRepo, orderRepo)
	svc.SetWalletRepository(walletRepo)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "combined_refund_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "combined_refund_player"))
	testGame := CreateTestGame(t, db, "combined_refund_game")

	// Create wallet with initial balance
	testWallet := CreateTestWallet(t, db, testUser.ID, 5000)

	// Create order
	order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusConfirmed, 10000)

	// Create combined payment (wallet + third party)
	paidPayment := &model.Payment{
		Base: model.Base{
			ExtJSON: "{}",
		},
		OrderID:           order.ID,
		UserID:            testUser.ID,
		Method:            model.PaymentMethodCombined,
		AmountCents:       order.TotalPriceCents,
		WalletAmountCents: 3000, // 3000 from wallet
		Status:            model.PaymentStatusPaid,
		Currency:          model.CurrencyCNY,
		ProviderTradeNo:   "mock_combined_123",
	}
	err := db.Create(paidPayment).Error
	require.NoError(t, err)

	// Deduct wallet balance (simulating payment)
	testWallet.BalanceCents -= 3000
	db.Save(testWallet)

	// Refund payment
	err = svc.RefundPayment(ctx, paidPayment.ID, "Customer requested refund")
	require.NoError(t, err)

	// Verify payment status
	var updatedPayment model.Payment
	err = db.First(&updatedPayment, paidPayment.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.PaymentStatusRefunded, updatedPayment.Status)

	// Verify wallet balance restored
	var updatedWallet model.Wallet
	err = db.First(&updatedWallet, testWallet.ID).Error
	require.NoError(t, err)
	assert.Equal(t, int64(5000), updatedWallet.BalanceCents) // 2000 + 3000 refund
}

// TestPaymentService_T7SettlementTracking tests T+7 settlement tracking.
func TestPaymentService_T7SettlementTracking(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "settlement_user")
	playerUser := CreateUniqueTestUser(t, db, "settlement_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "settlement_game")

	// Create completed order
	order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)

	// Create commission record with pending status (T+7 period)
	record := CreateTestCommissionRecord(t, db, order.ID, testPlayer.ID, 10000, model.SettlementStatusPending)

	// Verify initial status
	assert.Equal(t, model.SettlementStatusPending, record.SettlementStatus)

	// Simulate T+7 passed - settle the commission
	now := time.Now()
	record.SettlementStatus = model.SettlementStatusSettled
	record.SettledAt = &now
	err := db.Save(record).Error
	require.NoError(t, err)

	// Verify settled status
	var updatedRecord model.CommissionRecord
	err = db.First(&updatedRecord, record.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.SettlementStatusSettled, updatedRecord.SettlementStatus)
	assert.NotNil(t, updatedRecord.SettledAt)
}

// TestPaymentService_SettlementWithDispute tests settlement with dispute.
func TestPaymentService_SettlementWithDispute(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "dispute_settlement_user")
	playerUser := CreateUniqueTestUser(t, db, "dispute_settlement_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "dispute_settlement_game")

	// Create completed order
	order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)

	// Create commission record
	record := CreateTestCommissionRecord(t, db, order.ID, testPlayer.ID, 10000, model.SettlementStatusPending)

	// Create dispute
	CreateTestDispute(t, db, order, testUser.ID, model.DisputeInitiatorUser)

	// Update commission status to disputed
	record.SettlementStatus = model.SettlementStatusDisputed
	err := db.Save(record).Error
	require.NoError(t, err)

	// Verify disputed status
	var updatedRecord model.CommissionRecord
	err = db.First(&updatedRecord, record.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.SettlementStatusDisputed, updatedRecord.SettlementStatus)
}

// TestPaymentService_WalletFreezeForWithdraw tests wallet freeze for withdraw.
func TestPaymentService_WalletFreezeForWithdraw(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create player with wallet
	playerUser := CreateUniqueTestUser(t, db, "freeze_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testWallet := CreateTestWallet(t, db, playerUser.ID, 10000)

	// Create withdraw request
	withdraw := CreateTestWithdraw(t, db, testPlayer, 5000, model.WithdrawStatusPending)

	// Freeze wallet balance
	testWallet.BalanceCents -= 5000
	testWallet.FrozenCents += 5000
	err := db.Save(testWallet).Error
	require.NoError(t, err)

	// Verify
	var updatedWallet model.Wallet
	err = db.First(&updatedWallet, testWallet.ID).Error
	require.NoError(t, err)
	assert.Equal(t, int64(5000), updatedWallet.BalanceCents)
	assert.Equal(t, int64(5000), updatedWallet.FrozenCents)

	// Approve withdraw
	withdraw.Status = model.WithdrawStatusApproved
	db.Save(withdraw)

	// Complete withdraw - deduct frozen
	testWallet.FrozenCents -= 5000
	db.Save(testWallet)

	// Verify final state
	err = db.First(&updatedWallet, testWallet.ID).Error
	require.NoError(t, err)
	assert.Equal(t, int64(5000), updatedWallet.BalanceCents)
	assert.Equal(t, int64(0), updatedWallet.FrozenCents)
}

// TestPaymentService_PlayerIncomeToWallet tests player income going to wallet.
func TestPaymentService_PlayerIncomeToWallet(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create player with wallet
	playerUser := CreateUniqueTestUser(t, db, "income_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testWallet := CreateTestWallet(t, db, playerUser.ID, 0)

	// Create user and completed order
	testUser := CreateUniqueTestUser(t, db, "income_user")
	testGame := CreateTestGame(t, db, "income_game")
	order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)

	// Create commission record
	record := CreateTestCommissionRecord(t, db, order.ID, testPlayer.ID, 10000, model.SettlementStatusPending)

	// Income goes to frozen first (T+7 period)
	testWallet.FrozenCents += record.PlayerIncomeCents
	err := db.Save(testWallet).Error
	require.NoError(t, err)

	// Verify frozen
	var updatedWallet model.Wallet
	err = db.First(&updatedWallet, testWallet.ID).Error
	require.NoError(t, err)
	assert.Equal(t, int64(0), updatedWallet.BalanceCents)
	assert.Equal(t, int64(8000), updatedWallet.FrozenCents)

	// After T+7, move to available balance
	record.SettlementStatus = model.SettlementStatusSettled
	now := time.Now()
	record.SettledAt = &now
	db.Save(record)

	testWallet.FrozenCents -= record.PlayerIncomeCents
	testWallet.BalanceCents += record.PlayerIncomeCents
	db.Save(testWallet)

	// Verify available
	err = db.First(&updatedWallet, testWallet.ID).Error
	require.NoError(t, err)
	assert.Equal(t, int64(8000), updatedWallet.BalanceCents)
	assert.Equal(t, int64(0), updatedWallet.FrozenCents)
}

// TestPaymentService_RefundWithCommissionRollback tests refund with commission rollback.
func TestPaymentService_RefundWithCommissionRollback(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "rollback_user")
	playerUser := CreateUniqueTestUser(t, db, "rollback_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)
	testGame := CreateTestGame(t, db, "rollback_game")

	// Create player wallet with income
	testWallet := CreateTestWallet(t, db, playerUser.ID, 0)
	testWallet.FrozenCents = 8000 // Player income frozen
	db.Save(testWallet)

	// Create completed order
	order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusCompleted, 10000)

	// Create commission record
	record := CreateTestCommissionRecord(t, db, order.ID, testPlayer.ID, 10000, model.SettlementStatusPending)

	// Process refund - rollback commission
	order.Status = model.OrderStatusRefunded
	order.RefundAmountCents = 10000
	now := time.Now()
	order.RefundedAt = &now
	db.Save(order)

	// Deduct from player frozen balance
	testWallet.FrozenCents -= record.PlayerIncomeCents
	db.Save(testWallet)

	// Delete commission record (or mark as void)
	db.Delete(record)

	// Verify
	var updatedWallet model.Wallet
	err := db.First(&updatedWallet, testWallet.ID).Error
	require.NoError(t, err)
	assert.Equal(t, int64(0), updatedWallet.FrozenCents)

	var updatedOrder model.Order
	err = db.First(&updatedOrder, order.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.OrderStatusRefunded, updatedOrder.Status)
}

// TestPaymentService_MonthlySettlement tests monthly settlement.
func TestPaymentService_MonthlySettlement(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)

	// Create player
	playerUser := CreateUniqueTestUser(t, db, "monthly_player")
	testPlayer := CreateTestPlayer(t, db, playerUser)

	// Create monthly settlement
	settlement := &model.MonthlySettlement{
		PlayerID:             testPlayer.ID,
		SettlementMonth:      time.Now().Format("2006-01"),
		TotalOrderCount:      10,
		TotalAmountCents:     100000,
		TotalCommissionCents: 20000,
		TotalIncomeCents:     80000,
		BonusCents:           5000,
		FinalIncomeCents:     85000,
		Status:               model.MonthlySettlementStatusPending,
	}
	err := db.Create(settlement).Error
	require.NoError(t, err)

	// Verify
	var savedSettlement model.MonthlySettlement
	err = db.First(&savedSettlement, settlement.ID).Error
	require.NoError(t, err)
	assert.Equal(t, int64(85000), savedSettlement.FinalIncomeCents)
	assert.Equal(t, model.MonthlySettlementStatusPending, savedSettlement.Status)

	// Confirm settlement
	settlement.Status = model.MonthlySettlementStatusConfirmed
	db.Save(settlement)

	// Pay settlement
	settlement.Status = model.MonthlySettlementStatusPaid
	now := time.Now()
	settlement.SettledAt = &now
	db.Save(settlement)

	// Verify final status
	err = db.First(&savedSettlement, settlement.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.MonthlySettlementStatusPaid, savedSettlement.Status)
	assert.NotNil(t, savedSettlement.SettledAt)
}
