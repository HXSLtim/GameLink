// Package integration provides integration tests for services.
package integration

import (
	"context"
	"testing"

	"gamelink/internal/model"
	"gamelink/internal/repository/implementations"
	"gamelink/internal/repository/payment"
	"gamelink/internal/repository/wallet"
	paymentservice "gamelink/internal/service/payment"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaymentService_CreatePayment_ThirdParty(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	paymentRepo := payment.NewPaymentRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := paymentservice.NewPaymentService(paymentRepo, orderRepo)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "payment_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "payment_player"))
	testGame := CreateTestGame(t, db, "payment_game")

	// Create real order in database
	order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)

	// Create payment
	req := paymentservice.CreatePaymentRequest{
		OrderID: order.ID,
		Method:  model.PaymentMethodWeChat,
	}

	resp, err := svc.CreatePayment(ctx, testUser.ID, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotZero(t, resp.PaymentID)
	assert.NotNil(t, resp.PayInfo)
	assert.Equal(t, order.TotalPriceCents, resp.ThirdPartyAmount)
}

func TestPaymentService_CreatePayment_WalletOnly(t *testing.T) {
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
	testUser := CreateUniqueTestUser(t, db, "wallet_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "wallet_player"))
	testGame := CreateTestGame(t, db, "wallet_game")

	// Create wallet with sufficient balance
	CreateTestWallet(t, db, testUser.ID, 20000)

	// Create real order in database
	order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)

	// Create wallet payment
	req := paymentservice.CreatePaymentRequest{
		OrderID: order.ID,
		Method:  model.PaymentMethodWallet,
	}

	resp, err := svc.CreatePayment(ctx, testUser.ID, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotZero(t, resp.PaymentID)
	assert.True(t, resp.WalletPaidDirect)
	assert.Equal(t, order.TotalPriceCents, resp.WalletDeducted)

	// Verify wallet balance deducted
	updatedWallet, err := walletRepo.GetByUserID(ctx, testUser.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(10000), updatedWallet.BalanceCents) // 20000 - 10000
}

func TestPaymentService_CreatePayment_WalletInsufficientBalance(t *testing.T) {
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
	testUser := CreateUniqueTestUser(t, db, "insufficient_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "insufficient_player"))
	testGame := CreateTestGame(t, db, "insufficient_game")

	// Create wallet with insufficient balance
	CreateTestWallet(t, db, testUser.ID, 5000)

	// Create real order in database
	order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)

	// Try wallet payment
	req := paymentservice.CreatePaymentRequest{
		OrderID: order.ID,
		Method:  model.PaymentMethodWallet,
	}

	_, err := svc.CreatePayment(ctx, testUser.ID, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient wallet balance")
}

func TestPaymentService_CreatePayment_Combined(t *testing.T) {
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
	testUser := CreateUniqueTestUser(t, db, "combined_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "combined_player"))
	testGame := CreateTestGame(t, db, "combined_game")

	// Create wallet with partial balance
	CreateTestWallet(t, db, testUser.ID, 5000)

	// Create real order in database
	order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)

	// Create combined payment
	req := paymentservice.CreatePaymentRequest{
		OrderID:           order.ID,
		Method:            model.PaymentMethodCombined,
		WalletAmountCents: 3000,
		ThirdPartyMethod:  model.PaymentMethodWeChat,
	}

	resp, err := svc.CreatePayment(ctx, testUser.ID, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotZero(t, resp.PaymentID)
	assert.Equal(t, int64(3000), resp.WalletDeducted)
	assert.Equal(t, int64(7000), resp.ThirdPartyAmount)

	// Verify wallet balance deducted
	updatedWallet, err := walletRepo.GetByUserID(ctx, testUser.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(2000), updatedWallet.BalanceCents) // 5000 - 3000
}

func TestPaymentService_GetPaymentStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	paymentRepo := payment.NewPaymentRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := paymentservice.NewPaymentService(paymentRepo, orderRepo)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "status_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "status_player"))
	testGame := CreateTestGame(t, db, "status_game")

	// Create real order in database
	order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)

	// Create payment first
	createReq := paymentservice.CreatePaymentRequest{
		OrderID: order.ID,
		Method:  model.PaymentMethodWeChat,
	}

	createResp, err := svc.CreatePayment(ctx, testUser.ID, createReq)
	require.NoError(t, err)

	// Get payment status
	statusResp, err := svc.GetPaymentStatus(ctx, createResp.PaymentID)
	require.NoError(t, err)
	assert.Equal(t, createResp.PaymentID, statusResp.PaymentID)
	assert.Equal(t, order.ID, statusResp.OrderID)
	// Mock payment auto-completes, so status should be paid
	assert.Equal(t, model.PaymentStatusPaid, statusResp.Status)
}

func TestPaymentService_GetPaymentStatus_NotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	paymentRepo := payment.NewPaymentRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := paymentservice.NewPaymentService(paymentRepo, orderRepo)

	// Try to get non-existent payment
	_, err := svc.GetPaymentStatus(ctx, 99999)
	assert.Error(t, err)
	assert.Equal(t, paymentservice.ErrNotFound, err)
}

func TestPaymentService_CancelPayment(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	paymentRepo := payment.NewPaymentRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := paymentservice.NewPaymentService(paymentRepo, orderRepo)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "cancel_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "cancel_player"))
	testGame := CreateTestGame(t, db, "cancel_game")

	// Create real order in database
	order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)

	// Create a pending payment directly in DB (bypass mock auto-complete)
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
	err := paymentRepo.Create(ctx, pendingPayment)
	require.NoError(t, err)

	// Cancel payment
	err = svc.CancelPayment(ctx, testUser.ID, pendingPayment.ID)
	require.NoError(t, err)

	// Verify payment status
	statusResp, err := svc.GetPaymentStatus(ctx, pendingPayment.ID)
	require.NoError(t, err)
	assert.Equal(t, model.PaymentStatusFailed, statusResp.Status)
}

func TestPaymentService_CancelPayment_Unauthorized(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	paymentRepo := payment.NewPaymentRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := paymentservice.NewPaymentService(paymentRepo, orderRepo)

	// Create test data
	testUser := CreateUniqueTestUser(t, db, "unauth_user")
	otherUser := CreateUniqueTestUser(t, db, "other_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "unauth_player"))
	testGame := CreateTestGame(t, db, "unauth_game")

	// Create real order in database
	order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)

	// Create a pending payment
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
	err := paymentRepo.Create(ctx, pendingPayment)
	require.NoError(t, err)

	// Try to cancel with different user
	err = svc.CancelPayment(ctx, otherUser.ID, pendingPayment.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
}

func TestPaymentService_RefundPayment(t *testing.T) {
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
	testUser := CreateUniqueTestUser(t, db, "refund_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "refund_player"))
	testGame := CreateTestGame(t, db, "refund_game")

	// Create wallet
	CreateTestWallet(t, db, testUser.ID, 0)

	// Create real order in database with confirmed status
	order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusConfirmed, 10000)

	// Create a paid payment
	paidPayment := &model.Payment{
		Base: model.Base{
			ExtJSON: "{}",
		},
		OrderID:         order.ID,
		UserID:          testUser.ID,
		Method:          model.PaymentMethodWeChat,
		AmountCents:     order.TotalPriceCents,
		Status:          model.PaymentStatusPaid,
		Currency:        model.CurrencyCNY,
		ProviderTradeNo: "mock_trade_123",
	}
	err := paymentRepo.Create(ctx, paidPayment)
	require.NoError(t, err)

	// Refund payment
	err = svc.RefundPayment(ctx, paidPayment.ID, "Test refund reason")
	require.NoError(t, err)

	// Verify payment status
	statusResp, err := svc.GetPaymentStatus(ctx, paidPayment.ID)
	require.NoError(t, err)
	assert.Equal(t, model.PaymentStatusRefunded, statusResp.Status)

	// Verify order status
	updatedOrder, err := orderRepo.Get(ctx, order.ID)
	require.NoError(t, err)
	assert.Equal(t, model.OrderStatusRefunded, updatedOrder.Status)
}

func TestPaymentService_GetWalletBalance(t *testing.T) {
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
	testUser := CreateUniqueTestUser(t, db, "balance_user")

	// Create wallet with balance
	testWallet := CreateTestWallet(t, db, testUser.ID, 15000)
	testWallet.FrozenCents = 5000
	db.Save(testWallet)

	// Get wallet balance
	resp, err := svc.GetWalletBalance(ctx, testUser.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(15000), resp.BalanceCents)
	assert.Equal(t, int64(5000), resp.FrozenCents)
}

func TestPaymentService_GetWalletBalance_NoWallet(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	paymentRepo := payment.NewPaymentRepository(db)
	walletRepo := wallet.NewWalletRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := paymentservice.NewPaymentService(paymentRepo, orderRepo)
	svc.SetWalletRepository(walletRepo)

	// Create test user without wallet
	testUser := CreateUniqueTestUser(t, db, "no_wallet_user")

	// Get wallet balance (should return 0)
	resp, err := svc.GetWalletBalance(ctx, testUser.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), resp.BalanceCents)
	assert.Equal(t, int64(0), resp.FrozenCents)
}

func TestPaymentService_CalculateCombinedPayment(t *testing.T) {
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
	testUser := CreateUniqueTestUser(t, db, "calc_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "calc_player"))
	testGame := CreateTestGame(t, db, "calc_game")

	// Create wallet with balance
	CreateTestWallet(t, db, testUser.ID, 8000)

	// Create real order in database
	order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)

	// Calculate combined payment
	req := paymentservice.CalculateCombinedPaymentRequest{
		OrderID:           order.ID,
		WalletAmountCents: 5000,
	}

	resp, err := svc.CalculateCombinedPayment(ctx, testUser.ID, req)
	require.NoError(t, err)
	assert.Equal(t, int64(10000), resp.OrderTotalCents)
	assert.Equal(t, int64(8000), resp.WalletBalanceCents)
	assert.Equal(t, int64(5000), resp.WalletUsableCents)
	assert.Equal(t, int64(5000), resp.ThirdPartyAmountCents)
	assert.False(t, resp.CanPayWithWalletOnly)
}

func TestPaymentService_CalculateCombinedPayment_CanPayWithWalletOnly(t *testing.T) {
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
	testUser := CreateUniqueTestUser(t, db, "full_wallet_user")
	testPlayer := CreateTestPlayer(t, db, CreateUniqueTestUser(t, db, "full_wallet_player"))
	testGame := CreateTestGame(t, db, "full_wallet_game")

	// Create wallet with sufficient balance
	CreateTestWallet(t, db, testUser.ID, 15000)

	// Create real order in database
	order := CreateTestOrderWithDetails(t, db, testUser, testPlayer, testGame, model.OrderStatusPending, 10000)

	// Calculate combined payment
	req := paymentservice.CalculateCombinedPaymentRequest{
		OrderID: order.ID,
	}

	resp, err := svc.CalculateCombinedPayment(ctx, testUser.ID, req)
	require.NoError(t, err)
	assert.True(t, resp.CanPayWithWalletOnly)
	assert.Equal(t, int64(10000), resp.WalletUsableCents) // Capped at order total
	assert.Equal(t, int64(0), resp.ThirdPartyAmountCents)
}
