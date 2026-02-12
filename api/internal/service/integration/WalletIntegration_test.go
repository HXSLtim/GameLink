// Package integration provides integration tests for Wallet service with PostgreSQL.
package integration

import (
	"context"
	"sync"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository/implementations"
	"gamelink/internal/repository/payment"
	walletrepo "gamelink/internal/repository/wallet"
	"gamelink/internal/service/wallet"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Recharge Tests
// ============================================================================

func TestWalletService_Recharge_Success_NewWallet(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	walletRepo := walletrepo.NewWalletRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := wallet.NewWalletService(walletRepo, paymentRepo, orderRepo)

	// Create test user
	testUser := CreateUniqueTestUser(t, db, "recharge_user")
	amountCents := int64(10000) // 100.00 CNY

	// Recharge with Alipay
	req := wallet.RechargeRequest{
		AmountCents: amountCents,
		Method:      model.PaymentMethodAlipay,
	}

	resp, err := svc.Recharge(ctx, testUser.ID, req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Verify response
	assert.NotZero(t, resp.OrderID)
	assert.NotZero(t, resp.PaymentID)
	assert.Equal(t, amountCents, resp.Balance)

	// Verify order was created
	var order model.Order
	err = db.First(&order, resp.OrderID).Error
	require.NoError(t, err)
	assert.Equal(t, testUser.ID, order.UserID)
	assert.Equal(t, amountCents, order.TotalPriceCents)
	assert.Equal(t, "Wallet Recharge", order.Title)
	assert.Equal(t, model.OrderStatusCompleted, order.Status)
	assert.NotNil(t, order.CompletedAt)

	// Verify payment was created
	var payment model.Payment
	err = db.First(&payment, resp.PaymentID).Error
	require.NoError(t, err)
	assert.Equal(t, resp.OrderID, payment.OrderID)
	assert.Equal(t, testUser.ID, payment.UserID)
	assert.Equal(t, amountCents, payment.AmountCents)
	assert.Equal(t, model.PaymentMethodAlipay, payment.Method)
	assert.Equal(t, model.PaymentStatusPaid, payment.Status)
	assert.NotNil(t, payment.PaidAt)

	// Verify wallet was created
	var wallet model.Wallet
	err = db.Where("user_id = ?", testUser.ID).First(&wallet).Error
	require.NoError(t, err)
	assert.Equal(t, testUser.ID, wallet.UserID)
	assert.Equal(t, amountCents, wallet.BalanceCents)
	assert.Equal(t, int64(0), wallet.FrozenCents)
}

func TestWalletService_Recharge_Success_ExistingWallet(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	walletRepo := walletrepo.NewWalletRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := wallet.NewWalletService(walletRepo, paymentRepo, orderRepo)

	// Create test user with existing wallet
	testUser := CreateUniqueTestUser(t, db, "recharge_existing_user")
	existingWallet := CreateTestWallet(t, db, testUser.ID, 5000) // 50.00 CNY

	rechargeAmount := int64(3000) // 30.00 CNY
	req := wallet.RechargeRequest{
		AmountCents: rechargeAmount,
		Method:      model.PaymentMethodWeChat,
	}

	resp, err := svc.Recharge(ctx, testUser.ID, req)
	require.NoError(t, err)

	// Verify balance increased
	expectedBalance := existingWallet.BalanceCents + rechargeAmount
	assert.Equal(t, expectedBalance, resp.Balance)

	// Verify wallet in database
	var updatedWallet model.Wallet
	err = db.Where("user_id = ?", testUser.ID).First(&updatedWallet).Error
	require.NoError(t, err)
	assert.Equal(t, expectedBalance, updatedWallet.BalanceCents)
	assert.Equal(t, existingWallet.FrozenCents, updatedWallet.FrozenCents)
}

func TestWalletService_Recharge_InvalidAmount(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	walletRepo := walletrepo.NewWalletRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := wallet.NewWalletService(walletRepo, paymentRepo, orderRepo)

	testUser := CreateUniqueTestUser(t, db, "invalid_amount_user")

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
			req := wallet.RechargeRequest{
				AmountCents: tc.amountCents,
				Method:      model.PaymentMethodAlipay,
			}

			resp, err := svc.Recharge(ctx, testUser.ID, req)
			require.Error(t, err)
			assert.Equal(t, wallet.ErrInvalidAmount, err)
			assert.Nil(t, resp)
		})
	}
}

func TestWalletService_Recharge_DifferentPaymentMethods(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	walletRepo := walletrepo.NewWalletRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := wallet.NewWalletService(walletRepo, paymentRepo, orderRepo)

	paymentMethods := []model.PaymentMethod{
		model.PaymentMethodWeChat,
		model.PaymentMethodAlipay,
		model.PaymentMethodWallet,
	}

	for _, method := range paymentMethods {
		t.Run(string(method), func(t *testing.T) {
			testUser := CreateUniqueTestUser(t, db, "method_"+string(method))

			req := wallet.RechargeRequest{
				AmountCents: 5000,
				Method:      method,
			}

			resp, err := svc.Recharge(ctx, testUser.ID, req)
			require.NoError(t, err)
			assert.NotNil(t, resp)

			// Verify payment method
			var payment model.Payment
			err = db.First(&payment, resp.PaymentID).Error
			require.NoError(t, err)
			assert.Equal(t, method, payment.Method)
		})
	}
}

func TestWalletService_Recharge_LargeAmount(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	walletRepo := walletrepo.NewWalletRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := wallet.NewWalletService(walletRepo, paymentRepo, orderRepo)

	testUser := CreateUniqueTestUser(t, db, "large_amount_user")
	largeAmount := int64(1000000) // 10000.00 CNY

	req := wallet.RechargeRequest{
		AmountCents: largeAmount,
		Method:      model.PaymentMethodWallet,
	}

	resp, err := svc.Recharge(ctx, testUser.ID, req)
	require.NoError(t, err)
	assert.Equal(t, largeAmount, resp.Balance)

	// Verify wallet
	var wallet model.Wallet
	err = db.Where("user_id = ?", testUser.ID).First(&wallet).Error
	require.NoError(t, err)
	assert.Equal(t, largeAmount, wallet.BalanceCents)
}

func TestWalletService_Recharge_MultipleSequential(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	walletRepo := walletrepo.NewWalletRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := wallet.NewWalletService(walletRepo, paymentRepo, orderRepo)

	testUser := CreateUniqueTestUser(t, db, "multi_recharge_user")

	// First recharge
	resp1, err := svc.Recharge(ctx, testUser.ID, wallet.RechargeRequest{
		AmountCents: 5000,
		Method:      model.PaymentMethodAlipay,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(5000), resp1.Balance)

	// Second recharge
	resp2, err := svc.Recharge(ctx, testUser.ID, wallet.RechargeRequest{
		AmountCents: 3000,
		Method:      model.PaymentMethodWeChat,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(8000), resp2.Balance)

	// Third recharge
	resp3, err := svc.Recharge(ctx, testUser.ID, wallet.RechargeRequest{
		AmountCents: 2000,
		Method:      model.PaymentMethodWallet,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(10000), resp3.Balance)

	// Verify final wallet state
	var wallet model.Wallet
	err = db.Where("user_id = ?", testUser.ID).First(&wallet).Error
	require.NoError(t, err)
	assert.Equal(t, int64(10000), wallet.BalanceCents)
}

// ============================================================================
// GetBalance Tests
// ============================================================================

func TestWalletService_GetBalance_ExistingWallet(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	walletRepo := walletrepo.NewWalletRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := wallet.NewWalletService(walletRepo, paymentRepo, orderRepo)

	testUser := CreateUniqueTestUser(t, db, "balance_user")
	balanceCents := int64(50000)
	frozenCents := int64(10000)

	// Create wallet with specific balances
	testWallet := &model.Wallet{
		Base:         model.Base{ExtJSON: "{}"},
		UserID:       testUser.ID,
		BalanceCents: balanceCents,
		FrozenCents:  frozenCents,
	}
	err := db.Create(testWallet).Error
	require.NoError(t, err)

	// Get balance
	walletData, err := svc.GetBalance(ctx, testUser.ID)
	require.NoError(t, err)
	assert.Equal(t, testUser.ID, walletData.UserID)
	assert.Equal(t, balanceCents, walletData.BalanceCents)
	assert.Equal(t, frozenCents, walletData.FrozenCents)

	// Verify total balance calculation
	totalBalance := walletData.BalanceCents + walletData.FrozenCents
	assert.Equal(t, int64(60000), totalBalance)
}

func TestWalletService_GetBalance_NoWallet_ReturnsEmpty(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	walletRepo := walletrepo.NewWalletRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := wallet.NewWalletService(walletRepo, paymentRepo, orderRepo)

	testUser := CreateUniqueTestUser(t, db, "no_wallet_user")

	// Get balance for user with no wallet
	walletData, err := svc.GetBalance(ctx, testUser.ID)
	require.NoError(t, err)
	assert.Equal(t, testUser.ID, walletData.UserID)
	assert.Equal(t, int64(0), walletData.BalanceCents)
	assert.Equal(t, int64(0), walletData.FrozenCents)
}

func TestWalletService_GetBalance_OnlyBalance_NoFrozen(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	walletRepo := walletrepo.NewWalletRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := wallet.NewWalletService(walletRepo, paymentRepo, orderRepo)

	testUser := CreateUniqueTestUser(t, db, "only_balance_user")
	CreateTestWallet(t, db, testUser.ID, 30000)

	walletData, err := svc.GetBalance(ctx, testUser.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(30000), walletData.BalanceCents)
	assert.Equal(t, int64(0), walletData.FrozenCents)
}

func TestWalletService_GetBalance_OnlyFrozen_NoBalance(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	walletRepo := walletrepo.NewWalletRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := wallet.NewWalletService(walletRepo, paymentRepo, orderRepo)

	testUser := CreateUniqueTestUser(t, db, "only_frozen_user")

	// Create wallet with only frozen funds
	testWallet := &model.Wallet{
		Base:         model.Base{ExtJSON: "{}"},
		UserID:       testUser.ID,
		BalanceCents: 0,
		FrozenCents:  20000,
	}
	err := db.Create(testWallet).Error
	require.NoError(t, err)

	walletData, err := svc.GetBalance(ctx, testUser.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), walletData.BalanceCents)
	assert.Equal(t, int64(20000), walletData.FrozenCents)
}

func TestWalletService_GetBalance_BothZero(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	walletRepo := walletrepo.NewWalletRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := wallet.NewWalletService(walletRepo, paymentRepo, orderRepo)

	testUser := CreateUniqueTestUser(t, db, "zero_balance_user")
	CreateTestWallet(t, db, testUser.ID, 0)

	walletData, err := svc.GetBalance(ctx, testUser.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), walletData.BalanceCents)
	assert.Equal(t, int64(0), walletData.FrozenCents)
}

// ============================================================================
// T+7 Settlement Tests (Frozen Balance)
// ============================================================================

func TestWalletService_TPlus7_FrozenBalance_Simulation(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	walletRepo := walletrepo.NewWalletRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := wallet.NewWalletService(walletRepo, paymentRepo, orderRepo)

	playerUser := CreateUniqueTestUser(t, db, "player_user")
	player := CreateTestPlayer(t, db, playerUser)
	testUser := CreateUniqueTestUser(t, db, "customer_user")
	testGame := CreateTestGame(t, db, "test_game")
	serviceItem := CreateTestServiceItem(t, db, testGame, "Test Service", 10000)

	// Simulate completed order - player should get frozen funds
	now := time.Now().UTC()
	scheduledStart := now.Add(time.Hour)
	scheduledEnd := now.Add(2 * time.Hour)

	playerIncomeCents := int64(8000) // 80% of 10000 after commission

	order := &model.Order{
		Base:              model.Base{ExtJSON: "{}"},
		OrderNo:           "TPLUS7" + now.Format("20060102150405"),
		UserID:            testUser.ID,
		PlayerID:          &player.ID,
		ItemID:            serviceItem.ID,
		Quantity:          1,
		UnitPriceCents:    10000,
		TotalPriceCents:   10000,
		CommissionCents:   2000,
		PlayerIncomeCents: playerIncomeCents,
		Currency:          model.CurrencyCNY,
		Status:            model.OrderStatusCompleted,
		Title:             "Test Order",
		ScheduledStart:    &scheduledStart,
		ScheduledEnd:      &scheduledEnd,
		CompletedAt:       &now,
		GameID:            &testGame.ID,
		OrderConfig:       "{}",
	}
	err := db.Create(order).Error
	require.NoError(t, err)

	// Simulate T+7: Create wallet with frozen funds
	playerWallet := &model.Wallet{
		Base:         model.Base{ExtJSON: "{}"},
		UserID:       player.ID,
		BalanceCents: 5000, // Existing withdrawable balance
		FrozenCents:  8000, // Frozen from completed order (T+7 pending)
	}
	err = db.Create(playerWallet).Error
	require.NoError(t, err)

	// Get balance - should show frozen funds
	walletData, err := svc.GetBalance(ctx, player.ID)
	require.NoError(t, err)

	// Verify T+7 scenario
	withdrawable := walletData.BalanceCents
	pendingSettlement := walletData.FrozenCents
	totalBalance := withdrawable + pendingSettlement

	assert.Equal(t, int64(5000), withdrawable, "withdrawable amount")
	assert.Equal(t, int64(8000), pendingSettlement, "frozen pending T+7")
	assert.Equal(t, int64(13000), totalBalance, "total balance")
}

func TestWalletService_TPlus7_Recharge_WithFrozenFunds(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	walletRepo := walletrepo.NewWalletRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := wallet.NewWalletService(walletRepo, paymentRepo, orderRepo)

	testUser := CreateUniqueTestUser(t, db, "frozen_recharge_user")

	// Create wallet with frozen funds (T+7 pending)
	testWallet := &model.Wallet{
		Base:         model.Base{ExtJSON: "{}"},
		UserID:       testUser.ID,
		BalanceCents: 5000,
		FrozenCents:  15000,
	}
	err := db.Create(testWallet).Error
	require.NoError(t, err)

	// Recharge - should only increase BalanceCents
	rechargeAmount := int64(10000)
	req := wallet.RechargeRequest{
		AmountCents: rechargeAmount,
		Method:      model.PaymentMethodAlipay,
	}

	resp, err := svc.Recharge(ctx, testUser.ID, req)
	require.NoError(t, err)

	// Verify only balance increased, frozen unchanged
	expectedBalance := int64(15000) // 5000 + 10000
	assert.Equal(t, expectedBalance, resp.Balance)

	var updatedWallet model.Wallet
	err = db.Where("user_id = ?", testUser.ID).First(&updatedWallet).Error
	require.NoError(t, err)
	assert.Equal(t, expectedBalance, updatedWallet.BalanceCents)
	assert.Equal(t, int64(15000), updatedWallet.FrozenCents, "frozen should be unchanged")
}

func TestWalletService_TPlus7_WithdrawableCalculation(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	walletRepo := walletrepo.NewWalletRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := wallet.NewWalletService(walletRepo, paymentRepo, orderRepo)

	playerUser := CreateUniqueTestUser(t, db, "withdrawable_user")

	// Scenario: Player with mixed balance states
	testWallet := &model.Wallet{
		Base:         model.Base{ExtJSON: "{}"},
		UserID:       playerUser.ID,
		BalanceCents: 8000,  // Withdrawable
		FrozenCents:  12000, // Pending T+7 settlement
	}
	err := db.Create(testWallet).Error
	require.NoError(t, err)

	walletData, err := svc.GetBalance(ctx, playerUser.ID)
	require.NoError(t, err)

	// Calculate withdrawable amount
	withdrawable := walletData.BalanceCents
	pendingSettlement := walletData.FrozenCents

	assert.Equal(t, int64(8000), withdrawable, "only balance_cents is withdrawable")
	assert.Equal(t, int64(12000), pendingSettlement, "frozen_cents pending T+7")

	// Total balance includes frozen funds
	totalBalance := withdrawable + pendingSettlement
	assert.Equal(t, int64(20000), totalBalance, "total balance")
}

// ============================================================================
// Precision Tests (Cents)
// ============================================================================

func TestWalletService_Precision_CentsAccuracy(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	walletRepo := walletrepo.NewWalletRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := wallet.NewWalletService(walletRepo, paymentRepo, orderRepo)

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
		{"one yuan one cent", 101, "1.01"},
		{"large with cents", 10099, "100.99"},
		{"exact large", 100000, "1000.00"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testUser := CreateUniqueTestUser(t, db, "precision_"+tc.name)

			req := wallet.RechargeRequest{
				AmountCents: tc.amountCents,
				Method:      model.PaymentMethodAlipay,
			}

			resp, err := svc.Recharge(ctx, testUser.ID, req)
			require.NoError(t, err)
			assert.Equal(t, tc.amountCents, resp.Balance)

			// Verify database stores exact cents value
			var wallet model.Wallet
			err = db.Where("user_id = ?", testUser.ID).First(&wallet).Error
			require.NoError(t, err)
			assert.Equal(t, tc.amountCents, wallet.BalanceCents)
		})
	}
}

func TestWalletService_Precision_Accumulation(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	walletRepo := walletrepo.NewWalletRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := wallet.NewWalletService(walletRepo, paymentRepo, orderRepo)

	testUser := CreateUniqueTestUser(t, db, "accumulate_user")

	// Recharge multiple times with cent values
	amounts := []int64{99, 1, 50, 50} // Total: 200 cents = 2.00
	expectedTotal := int64(0)

	for _, amount := range amounts {
		_, err := svc.Recharge(ctx, testUser.ID, wallet.RechargeRequest{
			AmountCents: amount,
			Method:      model.PaymentMethodAlipay,
		})
		require.NoError(t, err)
		expectedTotal += amount
	}

	// Verify final balance
	walletData, err := svc.GetBalance(ctx, testUser.ID)
	require.NoError(t, err)
	assert.Equal(t, expectedTotal, walletData.BalanceCents)
	assert.Equal(t, int64(200), expectedTotal, "total should be exactly 200 cents")
}

func TestWalletService_Precision_NoFloatingPointErrors(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	walletRepo := walletrepo.NewWalletRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := wallet.NewWalletService(walletRepo, paymentRepo, orderRepo)

	testUser := CreateUniqueTestUser(t, db, "precision_user")

	// Values that commonly cause floating point errors
	amounts := []int64{33, 33, 34} // Total: 100 cents
	expectedTotal := int64(100)

	for _, amount := range amounts {
		_, err := svc.Recharge(ctx, testUser.ID, wallet.RechargeRequest{
			AmountCents: amount,
			Method:      model.PaymentMethodAlipay,
		})
		require.NoError(t, err)
	}

	walletData, err := svc.GetBalance(ctx, testUser.ID)
	require.NoError(t, err)
	assert.Equal(t, expectedTotal, walletData.BalanceCents, "should sum exactly without floating point errors")
}

// ============================================================================
// Concurrency Tests
// ============================================================================

func TestWalletService_Concurrent_Recharge(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	walletRepo := walletrepo.NewWalletRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := wallet.NewWalletService(walletRepo, paymentRepo, orderRepo)

	testUser := CreateUniqueTestUser(t, db, "concurrent_user")

	numGoroutines := 10
	amountPerRecharge := int64(1000)
	expectedTotal := int64(numGoroutines) * amountPerRecharge

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines)

	// Concurrent recharges
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			_, err := svc.Recharge(ctx, testUser.ID, wallet.RechargeRequest{
				AmountCents: amountPerRecharge,
				Method:      model.PaymentMethodAlipay,
			})
			if err != nil {
				errors <- err
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	for err := range errors {
		t.Errorf("Concurrent recharge error: %v", err)
	}

	// Verify final balance is correct
	walletData, err := svc.GetBalance(ctx, testUser.ID)
	require.NoError(t, err)
	assert.Equal(t, expectedTotal, walletData.BalanceCents)
}

func TestWalletService_Concurrent_GetBalance(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	walletRepo := walletrepo.NewWalletRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := wallet.NewWalletService(walletRepo, paymentRepo, orderRepo)

	testUser := CreateUniqueTestUser(t, db, "concurrent_read_user")
	CreateTestWallet(t, db, testUser.ID, 5000)

	numGoroutines := 50
	var wg sync.WaitGroup

	// Concurrent balance reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := svc.GetBalance(ctx, testUser.ID)
			assert.NoError(t, err)
		}()
	}

	wg.Wait()
}

func TestWalletService_Concurrent_RechargeAndRead(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	walletRepo := walletrepo.NewWalletRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := wallet.NewWalletService(walletRepo, paymentRepo, orderRepo)

	testUser := CreateUniqueTestUser(t, db, "mixed_concurrent_user")

	numRecharges := 5
	numReads := 20
	var wg sync.WaitGroup

	// Concurrent recharges
	for i := 0; i < numRecharges; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			svc.Recharge(ctx, testUser.ID, wallet.RechargeRequest{
				AmountCents: 1000,
				Method:      model.PaymentMethodAlipay,
			})
		}()
	}

	// Concurrent reads
	for i := 0; i < numReads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			svc.GetBalance(ctx, testUser.ID)
		}()
	}

	wg.Wait()

	// Verify final state
	walletData, err := svc.GetBalance(ctx, testUser.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(5000), walletData.BalanceCents)
}

// ============================================================================
// Edge Cases
// ============================================================================

func TestWalletService_OrderTimestamps_AreSet(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	walletRepo := walletrepo.NewWalletRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := wallet.NewWalletService(walletRepo, paymentRepo, orderRepo)

	testUser := CreateUniqueTestUser(t, db, "timestamp_user")
	beforeTest := time.Now().UTC()

	req := wallet.RechargeRequest{
		AmountCents: 5000,
		Method:      model.PaymentMethodAlipay,
	}

	resp, err := svc.Recharge(ctx, testUser.ID, req)
	require.NoError(t, err)

	// Verify order timestamp
	var order model.Order
	err = db.First(&order, resp.OrderID).Error
	require.NoError(t, err)
	assert.NotNil(t, order.CompletedAt)
	assert.True(t, order.CompletedAt.After(beforeTest) || order.CompletedAt.Equal(beforeTest))

	// Verify payment timestamp
	var payment model.Payment
	err = db.First(&payment, resp.PaymentID).Error
	require.NoError(t, err)
	assert.NotNil(t, payment.PaidAt)
	assert.True(t, payment.PaidAt.After(beforeTest) || payment.PaidAt.Equal(beforeTest))
}

func TestWalletService_Recharge_OrderFields_Correct(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	walletRepo := walletrepo.NewWalletRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := wallet.NewWalletService(walletRepo, paymentRepo, orderRepo)

	testUser := CreateUniqueTestUser(t, db, "fields_user")
	amountCents := int64(10000)

	req := wallet.RechargeRequest{
		AmountCents: amountCents,
		Method:      model.PaymentMethodWeChat,
	}

	resp, err := svc.Recharge(ctx, testUser.ID, req)
	require.NoError(t, err)

	// Verify order fields
	var order model.Order
	err = db.First(&order, resp.OrderID).Error
	require.NoError(t, err)
	assert.Equal(t, testUser.ID, order.UserID)
	assert.Equal(t, uint64(0), order.ItemID, "recharge orders have no service item")
	assert.Equal(t, "Wallet Recharge", order.Title)
	assert.Equal(t, amountCents, order.UnitPriceCents)
	assert.Equal(t, amountCents, order.TotalPriceCents)
	assert.Equal(t, model.OrderStatusCompleted, order.Status)

	// Verify payment fields
	var payment model.Payment
	err = db.First(&payment, resp.PaymentID).Error
	require.NoError(t, err)
	assert.Equal(t, resp.OrderID, payment.OrderID)
	assert.Equal(t, testUser.ID, payment.UserID)
	assert.Equal(t, model.PaymentMethodWeChat, payment.Method)
	assert.Equal(t, model.CurrencyCNY, payment.Currency)
	assert.Equal(t, model.PaymentStatusPaid, payment.Status)
}

func TestWalletService_MultipleUsers_IsolatedWallets(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	walletRepo := walletrepo.NewWalletRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := wallet.NewWalletService(walletRepo, paymentRepo, orderRepo)

	// Create multiple users
	user1 := CreateUniqueTestUser(t, db, "user1_isolated")
	user2 := CreateUniqueTestUser(t, db, "user2_isolated")
	user3 := CreateUniqueTestUser(t, db, "user3_isolated")

	// Recharge different amounts
	_, err := svc.Recharge(ctx, user1.ID, wallet.RechargeRequest{
		AmountCents: 1000,
		Method:      model.PaymentMethodAlipay,
	})
	require.NoError(t, err)

	_, err = svc.Recharge(ctx, user2.ID, wallet.RechargeRequest{
		AmountCents: 5000,
		Method:      model.PaymentMethodAlipay,
	})
	require.NoError(t, err)

	_, err = svc.Recharge(ctx, user3.ID, wallet.RechargeRequest{
		AmountCents: 10000,
		Method:      model.PaymentMethodAlipay,
	})
	require.NoError(t, err)

	// Verify each user's wallet is independent
	wallet1, _ := svc.GetBalance(ctx, user1.ID)
	wallet2, _ := svc.GetBalance(ctx, user2.ID)
	wallet3, _ := svc.GetBalance(ctx, user3.ID)

	assert.Equal(t, int64(1000), wallet1.BalanceCents)
	assert.Equal(t, int64(5000), wallet2.BalanceCents)
	assert.Equal(t, int64(10000), wallet3.BalanceCents)
}

func TestWalletService_WalletUniqueness_PerUser(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	walletRepo := walletrepo.NewWalletRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := wallet.NewWalletService(walletRepo, paymentRepo, orderRepo)

	testUser := CreateUniqueTestUser(t, db, "unique_wallet_user")

	// Multiple recharges should update same wallet, not create new ones
	_, err := svc.Recharge(ctx, testUser.ID, wallet.RechargeRequest{
		AmountCents: 1000,
		Method:      model.PaymentMethodAlipay,
	})
	require.NoError(t, err)

	_, err = svc.Recharge(ctx, testUser.ID, wallet.RechargeRequest{
		AmountCents: 2000,
		Method:      model.PaymentMethodWeChat,
	})
	require.NoError(t, err)

	// Verify only one wallet exists for the user
	var wallets []model.Wallet
	err = db.Where("user_id = ?", testUser.ID).Find(&wallets).Error
	require.NoError(t, err)
	assert.Len(t, wallets, 1, "should have exactly one wallet per user")
	assert.Equal(t, int64(3000), wallets[0].BalanceCents)
}

func TestWalletService_Currency_DefaultCNY(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	walletRepo := walletrepo.NewWalletRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	svc := wallet.NewWalletService(walletRepo, paymentRepo, orderRepo)

	testUser := CreateUniqueTestUser(t, db, "currency_user")

	req := wallet.RechargeRequest{
		AmountCents: 5000,
		Method:      model.PaymentMethodAlipay,
	}

	resp, err := svc.Recharge(ctx, testUser.ID, req)
	require.NoError(t, err)

	// Verify payment currency is CNY
	var payment model.Payment
	err = db.First(&payment, resp.PaymentID).Error
	require.NoError(t, err)
	assert.Equal(t, model.CurrencyCNY, payment.Currency, "default currency should be CNY")
}
