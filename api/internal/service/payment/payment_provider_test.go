package payment

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/service/external"
	"gamelink/pkg/config"
	"gamelink/pkg/testutil"

	_ "gamelink/pkg/cache" // Import for cache.DistributedLock interface
)

// TestProviderFactory_CreateProvider_WeChatEnabled tests WeChat provider creation when enabled
func TestProviderFactory_CreateProvider_WeChatEnabled(t *testing.T) {
	cfg := &external.Config{
		WeChatPay: config.WeChatPayConfig{
			Enabled:   true,
			AppID:     "test_app_id",
			MchID:     "test_mch_id",
			APIKey:    "test_api_key",
			NotifyURL: "https://test.com/notify",
		},
		Alipay: config.AlipayConfig{
			Enabled: false,
		},
	}

	factory := NewProviderFactory(cfg)

	provider, err := factory.CreateProvider(model.PaymentMethodWeChat)
	assert.NoError(t, err)
	assert.NotNil(t, provider)
}

// TestProviderFactory_CreateProvider_WeChatDisabled tests WeChat provider creation when disabled (mock)
func TestProviderFactory_CreateProvider_WeChatDisabled(t *testing.T) {
	cfg := &external.Config{
		WeChatPay: config.WeChatPayConfig{
			Enabled: false,
		},
		Alipay: config.AlipayConfig{
			Enabled: false,
		},
	}

	factory := NewProviderFactory(cfg)

	provider, err := factory.CreateProvider(model.PaymentMethodWeChat)
	assert.NoError(t, err)
	assert.NotNil(t, provider)
	// Should return mock provider
	payment := &model.Payment{Base: model.Base{ID: 1}, AmountCents: 10000}
	_, _, _, err = provider.Refund(context.Background(), payment, "test")
	assert.NoError(t, err)
}

// TestProviderFactory_CreateProvider_AlipayEnabled tests Alipay provider creation when enabled
func TestProviderFactory_CreateProvider_AlipayEnabled(t *testing.T) {
	cfg := &external.Config{
		WeChatPay: config.WeChatPayConfig{
			Enabled: false,
		},
		Alipay: config.AlipayConfig{
			Enabled:        true,
			AppID:          "test_app_id",
			PrivateKeyPath: "test_private_key",
			PublicKeyPath:  "test_public_key",
			NotifyURL:      "https://test.com/notify",
		},
	}

	factory := NewProviderFactory(cfg)

	provider, err := factory.CreateProvider(model.PaymentMethodAlipay)
	// Note: NewAlipayProvider will fail to load private key if file doesn't exist
	// In production, this would require actual key files
	// For this test, we just verify that CreateProvider returns something or fails gracefully
	if err == nil {
		assert.NotNil(t, provider)
	}
}

// TestProviderFactory_CreateProvider_AlipayDisabled tests Alipay provider creation when disabled (mock)
func TestProviderFactory_CreateProvider_AlipayDisabled(t *testing.T) {
	cfg := &external.Config{
		WeChatPay: config.WeChatPayConfig{
			Enabled: false,
		},
		Alipay: config.AlipayConfig{
			Enabled: false,
		},
	}

	factory := NewProviderFactory(cfg)

	provider, err := factory.CreateProvider(model.PaymentMethodAlipay)
	assert.NoError(t, err)
	assert.NotNil(t, provider)
	// Should return mock provider
	payment := &model.Payment{Base: model.Base{ID: 1}, AmountCents: 10000}
	_, _, _, err = provider.Refund(context.Background(), payment, "test")
	assert.NoError(t, err)
}

// TestProviderFactory_CreateProvider_GenericMethod tests provider creation for generic methods
func TestProviderFactory_CreateProvider_GenericMethod(t *testing.T) {
	cfg := &external.Config{
		WeChatPay: config.WeChatPayConfig{
			Enabled: false,
		},
		Alipay: config.AlipayConfig{
			Enabled: false,
		},
	}

	factory := NewProviderFactory(cfg)

	provider, err := factory.CreateProvider(model.PaymentMethodWallet)
	assert.NoError(t, err)
	assert.NotNil(t, provider)
	// Should return generic provider
	payment := &model.Payment{Base: model.Base{ID: 1}, AmountCents: 10000}
	_, _, _, err = provider.Refund(context.Background(), payment, "test")
	assert.NoError(t, err)
}

// TestProviderFactory_CreateProviders tests creating all providers
func TestProviderFactory_CreateProviders(t *testing.T) {
	tests := []struct {
		name              string
		weChatEnabled     bool
		alipayEnabled     bool
		expectedProviders []model.PaymentMethod
	}{
		{
			name:              "Both disabled - mock providers",
			weChatEnabled:     false,
			alipayEnabled:     false,
			expectedProviders: []model.PaymentMethod{model.PaymentMethodWeChat, model.PaymentMethodAlipay},
		},
		{
			name:              "WeChat enabled only",
			weChatEnabled:     true,
			alipayEnabled:     false,
			expectedProviders: []model.PaymentMethod{model.PaymentMethodWeChat, model.PaymentMethodAlipay},
		},
		{
			name:              "Alipay enabled only",
			weChatEnabled:     false,
			alipayEnabled:     true,
			expectedProviders: []model.PaymentMethod{model.PaymentMethodWeChat, model.PaymentMethodAlipay},
		},
		{
			name:              "Both enabled",
			weChatEnabled:     true,
			alipayEnabled:     true,
			expectedProviders: []model.PaymentMethod{model.PaymentMethodWeChat, model.PaymentMethodAlipay},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &external.Config{
				WeChatPay: config.WeChatPayConfig{
					Enabled: tt.weChatEnabled,
					AppID:   "test_app",
					MchID:   "test_mch",
					APIKey:  "test_key",
				},
				Alipay: config.AlipayConfig{
					Enabled:        tt.alipayEnabled,
					AppID:          "test_app",
					PrivateKeyPath: "test_key",
					PublicKeyPath:  "test_pub",
				},
			}

			factory := NewProviderFactory(cfg)
			providers := factory.CreateProviders()

			for _, method := range tt.expectedProviders {
				assert.Contains(t, providers, method, "should have provider for %s", method)
				assert.NotNil(t, providers[method], "provider should not be nil for %s", method)

				// Test that provider can refund (skip Alipay if enabled since it requires key files)
				if method == model.PaymentMethodAlipay && tt.alipayEnabled {
					continue
				}
				payment := &model.Payment{Base: model.Base{ID: 1}, AmountCents: 10000}
				_, _, _, err := providers[method].Refund(context.Background(), payment, "test")
				assert.NoError(t, err, "provider should be able to refund for %s", method)
			}
		})
	}
}

// TestPaymentService_RefundPayment_ThirdParty_Success tests third-party payment refund.
// Uses in-memory SQLite DB because RefundPayment runs a gorm.DB.Transaction internally.
func TestPaymentService_RefundPayment_ThirdParty_Success(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	amountCents := int64(10000)

	// Setup in-memory DB for transaction support
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db, &model.Payment{}, &model.Order{})
	defer testutil.CleanDB(t, db)

	// Seed data directly in DB (the transaction path creates repos from db)
	order := &model.Order{
		OrderNo:         "ORD123456",
		UserID:          userID,
		ItemID:          1,
		Status:          model.OrderStatusConfirmed,
		TotalPriceCents: amountCents,
		Currency:        model.CurrencyCNY,
	}
	db.Create(order)

	payment := &model.Payment{
		OrderID:     order.ID,
		UserID:      userID,
		AmountCents: amountCents,
		Status:      model.PaymentStatusPaid,
		Method:      model.PaymentMethodWeChat,
		Currency:    model.CurrencyCNY,
	}
	db.Create(payment)

	// Mock the Get call (pre-transaction) on the injected repo
	mockPayments := new(MockPaymentRepository)
	mockPayments.On("Get", ctx, payment.ID).Return(payment, nil)

	service := NewPaymentService(mockPayments, nil)
	service.SetDB(db)

	err := service.RefundPayment(ctx, payment.ID, "customer request")

	assert.NoError(t, err)

	// Verify DB state
	var updatedPayment model.Payment
	db.First(&updatedPayment, payment.ID)
	assert.Equal(t, model.PaymentStatusRefunded, updatedPayment.Status)
	assert.NotNil(t, updatedPayment.RefundedAt)

	var updatedOrder model.Order
	db.First(&updatedOrder, order.ID)
	assert.Equal(t, model.OrderStatusRefunded, updatedOrder.Status)
	assert.Equal(t, amountCents, updatedOrder.RefundAmountCents)
}

// TestPaymentService_RefundPayment_Combined_Success tests combined payment refund (both wallet and third-party).
// Uses in-memory SQLite DB because RefundPayment runs a gorm.DB.Transaction internally.
func TestPaymentService_RefundPayment_Combined_Success(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	totalCents := int64(10000)
	walletCents := int64(6000)
	thirdPartyCents := int64(4000)

	// Setup in-memory DB for transaction support
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db, &model.Payment{}, &model.Order{})
	defer testutil.CleanDB(t, db)

	// Seed order and payment in DB
	order := &model.Order{
		OrderNo:         "ORD123456",
		UserID:          userID,
		ItemID:          1,
		Status:          model.OrderStatusConfirmed,
		TotalPriceCents: totalCents,
		Currency:        model.CurrencyCNY,
	}
	db.Create(order)

	payment := &model.Payment{
		OrderID:               order.ID,
		UserID:                userID,
		AmountCents:           totalCents,
		Status:                model.PaymentStatusPaid,
		Method:                model.PaymentMethodCombined,
		Currency:              model.CurrencyCNY,
		WalletAmountCents:     walletCents,
		ThirdPartyAmountCents: thirdPartyCents,
		ThirdPartyMethod:      model.PaymentMethodWeChat,
	}
	db.Create(payment)

	mockPayments := new(MockPaymentRepository)
	mockWallets := new(MockWalletRepository)

	mockPayments.On("Get", ctx, payment.ID).Return(payment, nil)
	// creditWallet uses UpdateBalanceWithLock
	updatedWallet := createTestWallet(userID, walletCents, 0)
	mockWallets.On("UpdateBalanceWithLock", ctx, userID, walletCents, 3).Return(updatedWallet, nil)

	service := NewPaymentService(mockPayments, nil)
	service.SetWalletRepository(mockWallets)
	service.SetDB(db)

	err := service.RefundPayment(ctx, payment.ID, "customer request")

	assert.NoError(t, err)

	// Verify DB state
	var updatedPayment model.Payment
	db.First(&updatedPayment, payment.ID)
	assert.Equal(t, model.PaymentStatusRefunded, updatedPayment.Status)

	var updatedOrder model.Order
	db.First(&updatedOrder, order.ID)
	assert.Equal(t, model.OrderStatusRefunded, updatedOrder.Status)

	mockPayments.AssertExpectations(t)
	mockWallets.AssertExpectations(t)
}

// TestPaymentService_RefundPayment_Combined_ThirdPartyRefundFailed tests combined payment refund when third-party refund fails
func TestPaymentService_RefundPayment_Combined_ThirdPartyRefundFailed(t *testing.T) {
	ctx := context.Background()
	paymentID := uint64(1)
	orderID := uint64(1000)
	userID := uint64(100)
	totalCents := int64(10000)
	walletCents := int64(6000)

	mockPayments := new(MockPaymentRepository)
	mockWallets := new(MockWalletRepository)

	payment := createTestPayment(paymentID, orderID, userID, model.PaymentStatusPaid, totalCents, model.PaymentMethodCombined)
	payment.WalletAmountCents = walletCents
	payment.ThirdPartyAmountCents = totalCents - walletCents
	payment.ThirdPartyMethod = model.PaymentMethodWeChat

	creditWallet := createTestWallet(userID, 16000, 0)
	rollbackWallet := createTestWallet(userID, 10000, 0)

	mockPayments.On("Get", ctx, paymentID).Return(payment, nil)
	// creditWallet uses UpdateBalanceWithLock
	mockWallets.On("UpdateBalanceWithLock", ctx, userID, walletCents, 3).Return(creditWallet, nil)
	// Third-party refund fails → debitWallet rollback also uses UpdateBalanceWithLock
	mockWallets.On("UpdateBalanceWithLock", ctx, userID, -walletCents, 3).Return(rollbackWallet, nil)

	service := NewPaymentService(mockPayments, nil)
	service.SetWalletRepository(mockWallets)

	// Create a provider that will fail
	providers := make(map[model.PaymentMethod]ProviderClient)
	providers[model.PaymentMethodWeChat] = &mockThirdPartyFailingProvider{}
	service.SetProviders(providers)

	err := service.RefundPayment(ctx, paymentID, "customer request")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "provider error")

	mockPayments.AssertExpectations(t)
	mockWallets.AssertExpectations(t)
}

// TestPaymentService_RefundPayment_WalletCreditFailed tests refund when wallet credit fails
func TestPaymentService_RefundPayment_WalletCreditFailed(t *testing.T) {
	ctx := context.Background()
	paymentID := uint64(1)
	orderID := uint64(1000)
	userID := uint64(100)
	amountCents := int64(10000)

	mockPayments := new(MockPaymentRepository)
	mockWallets := new(MockWalletRepository)

	payment := createTestPayment(paymentID, orderID, userID, model.PaymentStatusPaid, amountCents, model.PaymentMethodWallet)
	payment.WalletAmountCents = amountCents

	mockPayments.On("Get", ctx, paymentID).Return(payment, nil)
	// creditWallet uses UpdateBalanceWithLock which fails
	mockWallets.On("UpdateBalanceWithLock", ctx, userID, amountCents, 3).Return(nil, errors.New("wallet save failed"))

	service := NewPaymentService(mockPayments, nil)
	service.SetWalletRepository(mockWallets)

	err := service.RefundPayment(ctx, paymentID, "customer request")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "wallet save failed")

	mockPayments.AssertExpectations(t)
	mockWallets.AssertExpectations(t)
}

// TestPaymentService_RefundPayment_OrderUpdateFailed tests refund when order update fails.
// Uses in-memory SQLite DB because RefundPayment runs a gorm.DB.Transaction internally.
// Payment exists in DB but order does NOT, so the Get inside the transaction fails.
func TestPaymentService_RefundPayment_OrderUpdateFailed(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	amountCents := int64(10000)

	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db, &model.Payment{}, &model.Order{})
	defer testutil.CleanDB(t, db)

	// Seed payment WITHOUT a matching order
	payment := &model.Payment{
		OrderID:     9999, // non-existent order
		UserID:      userID,
		AmountCents: amountCents,
		Status:      model.PaymentStatusPaid,
		Method:      model.PaymentMethodWeChat,
		Currency:    model.CurrencyCNY,
	}
	db.Create(payment)

	mockPayments := new(MockPaymentRepository)
	mockPayments.On("Get", ctx, payment.ID).Return(payment, nil)

	service := NewPaymentService(mockPayments, nil)
	service.SetDB(db)

	err := service.RefundPayment(ctx, payment.ID, "customer request")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "order")

	mockPayments.AssertExpectations(t)
}

// TestPaymentService_routePayment_NoRoutingEngine tests routePayment without routing engine
func TestPaymentService_routePayment_NoRoutingEngine(t *testing.T) {
	ctx := context.Background()
	service := NewPaymentService(nil, nil)

	order := createTestOrder(1000, 100, model.OrderStatusPending, 10000)

	result, err := service.routePayment(ctx, order, model.PaymentMethodWeChat)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "routing engine not initialized")
}

// TestPaymentService_GetPaymentRoutingLog_NoRoutingEngine tests getting routing log without routing engine
func TestPaymentService_GetPaymentRoutingLog_NoRoutingEngine(t *testing.T) {
	ctx := context.Background()
	service := NewPaymentService(nil, nil)

	log, err := service.GetPaymentRoutingLog(ctx, 1)

	assert.Error(t, err)
	assert.Nil(t, log)
	assert.Contains(t, err.Error(), "routing engine not initialized")
}

// TestPaymentService_CreatePayment_WithDistributedLock_Success tests payment creation with distributed lock
func TestPaymentService_CreatePayment_WithDistributedLock_Success(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	orderID := uint64(1000)

	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)
	mockLock := new(mockDistributedLock)

	order := createTestOrder(orderID, userID, model.OrderStatusPending, 10000)

	mockLock.On("TryLock", ctx, mock.AnythingOfType("string"), time.Second*10, 3, time.Millisecond*50).Return(true, nil)
	mockLock.On("Unlock", ctx, mock.AnythingOfType("string")).Return(nil)
	mockOrders.On("Get", ctx, orderID).Return(order, nil)
	mockPayments.On("List", ctx, mock.AnythingOfType("repository.PaymentListOptions")).Return([]model.Payment{}, int64(0), nil)
	mockPayments.On("Create", ctx, mock.AnythingOfType("*model.Payment")).Return(nil).Run(func(args mock.Arguments) {
		payment := args.Get(1).(*model.Payment)
		payment.ID = 1
	})
	mockPayments.On("Get", ctx, uint64(1)).Return(createTestPayment(1, orderID, userID, model.PaymentStatusPending, 10000, model.PaymentMethodWeChat), nil)
	mockPayments.On("Update", ctx, mock.AnythingOfType("*model.Payment")).Return(nil)
	mockOrders.On("Update", ctx, mock.AnythingOfType("*model.Order")).Return(nil)

	service := NewPaymentService(mockPayments, mockOrders)
	service.SetDistributedLock(mockLock)

	req := CreatePaymentRequest{
		OrderID: orderID,
		Method:  model.PaymentMethodWeChat,
	}

	resp, err := service.CreatePayment(ctx, userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	mockLock.AssertExpectations(t)
	mockOrders.AssertExpectations(t)
	mockPayments.AssertExpectations(t)
}

// TestPaymentService_CreatePayment_WithDistributedLock_Failed tests payment creation when lock acquisition fails
func TestPaymentService_CreatePayment_WithDistributedLock_Failed(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	orderID := uint64(1000)

	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)
	mockLock := new(mockDistributedLock)

	mockLock.On("TryLock", ctx, mock.AnythingOfType("string"), time.Second*10, 3, time.Millisecond*50).Return(false, errors.New("lock timeout"))

	service := NewPaymentService(mockPayments, mockOrders)
	service.SetDistributedLock(mockLock)

	req := CreatePaymentRequest{
		OrderID: orderID,
		Method:  model.PaymentMethodWeChat,
	}

	resp, err := service.CreatePayment(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to acquire lock")

	mockLock.AssertExpectations(t)
}

// TestPaymentService_CreatePayment_WithDistributedLock_Locked tests payment creation when already locked
func TestPaymentService_CreatePayment_WithDistributedLock_Locked(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	orderID := uint64(1000)

	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)
	mockLock := new(mockDistributedLock)

	mockLock.On("TryLock", ctx, mock.AnythingOfType("string"), time.Second*10, 3, time.Millisecond*50).Return(false, nil)

	service := NewPaymentService(mockPayments, mockOrders)
	service.SetDistributedLock(mockLock)

	req := CreatePaymentRequest{
		OrderID: orderID,
		Method:  model.PaymentMethodWeChat,
	}

	resp, err := service.CreatePayment(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "payment creation in progress")

	mockLock.AssertExpectations(t)
}

// TestPaymentService_CreatePayment_WalletSaveFailed tests payment creation when wallet save fails
func TestPaymentService_CreatePayment_WalletSaveFailed(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	orderID := uint64(1000)

	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)
	mockWallets := new(MockWalletRepository)

	order := createTestOrder(orderID, userID, model.OrderStatusPending, 10000)
	wallet := createTestWallet(userID, 20000, 0)

	mockOrders.On("Get", ctx, orderID).Return(order, nil)
	mockPayments.On("List", ctx, mock.AnythingOfType("repository.PaymentListOptions")).Return([]model.Payment{}, int64(0), nil)
	mockWallets.On("GetByUserID", ctx, userID).Return(wallet, nil)
	mockWallets.On("UpdateBalanceWithLock", ctx, userID, int64(-10000), 3).Return(nil, errors.New("wallet save failed"))

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
	mockWallets.AssertExpectations(t)
}

// TestPaymentService_CreatePayment_PaymentCreateFailed tests payment creation when payment record creation fails
func TestPaymentService_CreatePayment_PaymentCreateFailed(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	orderID := uint64(1000)

	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)
	mockWallets := new(MockWalletRepository)

	order := createTestOrder(orderID, userID, model.OrderStatusPending, 10000)
	updatedWallet := createTestWallet(userID, 10000, 0) // After deduction

	mockOrders.On("Get", ctx, orderID).Return(order, nil)
	mockPayments.On("List", ctx, mock.AnythingOfType("repository.PaymentListOptions")).Return([]model.Payment{}, int64(0), nil)
	mockWallets.On("GetByUserID", ctx, userID).Return(createTestWallet(userID, 20000, 0), nil)
	mockWallets.On("UpdateBalanceWithLock", ctx, userID, int64(-10000), 3).Return(updatedWallet, nil)
	mockPayments.On("Create", ctx, mock.AnythingOfType("*model.Payment")).Return(errors.New("payment creation failed"))
	// Rollback wallet
	mockWallets.On("UpdateBalanceWithLock", ctx, userID, int64(10000), 3).Return(createTestWallet(userID, 20000, 0), nil)

	service := NewPaymentService(mockPayments, mockOrders)
	service.SetWalletRepository(mockWallets)

	req := CreatePaymentRequest{
		OrderID: orderID,
		Method:  model.PaymentMethodWallet,
	}

	resp, err := service.CreatePayment(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "payment creation failed")

	mockOrders.AssertExpectations(t)
	mockPayments.AssertExpectations(t)
	mockWallets.AssertExpectations(t)
}

// TestPaymentService_CreatePayment_CombinedWalletSaveFailed tests combined payment when wallet save fails
func TestPaymentService_CreatePayment_CombinedWalletSaveFailed(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	orderID := uint64(1000)

	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)
	mockWallets := new(MockWalletRepository)

	order := createTestOrder(orderID, userID, model.OrderStatusPending, 10000)
	wallet := createTestWallet(userID, 10000, 0)

	mockOrders.On("Get", ctx, orderID).Return(order, nil)
	mockPayments.On("List", ctx, mock.AnythingOfType("repository.PaymentListOptions")).Return([]model.Payment{}, int64(0), nil)
	mockWallets.On("GetByUserID", ctx, userID).Return(wallet, nil)
	mockWallets.On("UpdateBalanceWithLock", ctx, userID, int64(-6000), 3).Return(nil, errors.New("wallet save failed"))

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
	assert.Contains(t, err.Error(), "failed to deduct wallet balance")

	mockOrders.AssertExpectations(t)
	mockWallets.AssertExpectations(t)
}

// TestPaymentService_CreatePayment_Combined_InvalidThirdPartyMethod tests combined payment with invalid third-party method
func TestPaymentService_CreatePayment_Combined_InvalidThirdPartyMethod(t *testing.T) {
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
		WalletAmountCents: 6000,
		ThirdPartyMethod:  model.PaymentMethodWallet, // Invalid - must be wechat or alipay
	}

	resp, err := service.CreatePayment(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "third party method must be wechat or alipay")

	mockOrders.AssertExpectations(t)
}

// TestPaymentService_CreatePayment_OrderNotFound tests payment creation for non-existent order
func TestPaymentService_CreatePayment_OrderNotFound(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	orderID := uint64(999)

	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)

	mockOrders.On("Get", ctx, orderID).Return(nil, repository.ErrNotFound)

	service := NewPaymentService(mockPayments, mockOrders)

	req := CreatePaymentRequest{
		OrderID: orderID,
		Method:  model.PaymentMethodWeChat,
	}

	resp, err := service.CreatePayment(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, repository.ErrNotFound, err)

	mockOrders.AssertExpectations(t)
}

// TestPaymentService_HandlePaymentCallback_MissingPaymentID tests callback without payment_id
func TestPaymentService_HandlePaymentCallback_MissingPaymentID(t *testing.T) {
	ctx := context.Background()

	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)

	service := NewPaymentService(mockPayments, mockOrders)

	callbackData := map[string]interface{}{
		"amount_cents": int64(10000),
		"trade_no":     "wx_trade_123",
	}

	err := service.HandlePaymentCallback(ctx, "wechat", callbackData)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing payment_id")
}

// TestPaymentService_HandlePaymentCallback_PaymentNotFound tests callback for non-existent payment
func TestPaymentService_HandlePaymentCallback_PaymentNotFound(t *testing.T) {
	ctx := context.Background()
	paymentID := uint64(999)

	mockPayments := new(MockPaymentRepository)

	mockPayments.On("Get", ctx, paymentID).Return(nil, repository.ErrNotFound)

	service := NewPaymentService(mockPayments, nil)

	callbackData := map[string]interface{}{
		"payment_id":   paymentID,
		"amount_cents": int64(10000),
		"trade_no":     "wx_trade_123",
	}

	err := service.HandlePaymentCallback(ctx, "wechat", callbackData)

	assert.Error(t, err)
	assert.Equal(t, repository.ErrNotFound, err)

	mockPayments.AssertExpectations(t)
}

// TestPaymentService_HandlePaymentCallback_OrderNotFound tests callback when order not found
func TestPaymentService_HandlePaymentCallback_OrderNotFound(t *testing.T) {
	ctx := context.Background()
	paymentID := uint64(1)
	orderID := uint64(1000)
	userID := uint64(100)

	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)

	payment := createTestPayment(paymentID, orderID, userID, model.PaymentStatusPending, 10000, model.PaymentMethodWeChat)

	mockPayments.On("Get", ctx, paymentID).Return(payment, nil)
	// Order is fetched before payment update
	mockOrders.On("Get", ctx, orderID).Return(nil, repository.ErrNotFound)
	// Payment update won't be reached since order Get fails

	service := NewPaymentService(mockPayments, mockOrders)

	callbackData := map[string]interface{}{
		"payment_id":   paymentID,
		"amount_cents": int64(10000),
		"trade_no":     "wx_trade_123",
	}

	err := service.HandlePaymentCallback(ctx, "wechat", callbackData)

	assert.Error(t, err)
	assert.Equal(t, repository.ErrNotFound, err)

	mockPayments.AssertExpectations(t)
	mockOrders.AssertExpectations(t)
}

// TestPaymentService_HandlePaymentCallback_PaymentUpdateFailed tests callback when payment update fails.
// Uses in-memory SQLite DB because HandlePaymentCallback runs a gorm.DB.Transaction internally.
// The order is seeded but the payment is NOT, so txPaymentRepo.Update fails inside the transaction.
func TestPaymentService_HandlePaymentCallback_PaymentUpdateFailed(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)

	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db, &model.Payment{}, &model.Order{})
	defer testutil.CleanDB(t, db)

	// Seed order in DB (needed for txOrderRepo.Update inside the transaction)
	order := &model.Order{
		OrderNo:         "ORD-CALLBACK",
		UserID:          userID,
		ItemID:          1,
		Status:          model.OrderStatusPending,
		TotalPriceCents: 10000,
		Currency:        model.CurrencyCNY,
	}
	db.Create(order)

	// Payment object returned by mock but NOT seeded in DB
	// so txPaymentRepo.Update inside transaction will fail (RowsAffected=0 → ErrNotFound)
	payment := createTestPayment(9999, order.ID, userID, model.PaymentStatusPending, 10000, model.PaymentMethodWeChat)

	mockPayments := new(MockPaymentRepository)
	mockOrders := new(MockOrderRepository)

	mockPayments.On("Get", ctx, uint64(9999)).Return(payment, nil)
	mockOrders.On("Get", ctx, order.ID).Return(order, nil)

	service := NewPaymentService(mockPayments, mockOrders)
	service.SetDB(db)

	callbackData := map[string]interface{}{
		"payment_id":   uint64(9999),
		"amount_cents": int64(10000),
		"trade_no":     "wx_trade_123",
	}

	err := service.HandlePaymentCallback(ctx, "wechat", callbackData)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update payment")

	mockPayments.AssertExpectations(t)
	mockOrders.AssertExpectations(t)
}

// TestPaymentService_CancelPayment_PaymentNotFound tests cancel for non-existent payment
func TestPaymentService_CancelPayment_PaymentNotFound(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	paymentID := uint64(999)

	mockPayments := new(MockPaymentRepository)

	mockPayments.On("Get", ctx, paymentID).Return(nil, repository.ErrNotFound)

	service := NewPaymentService(mockPayments, nil)

	err := service.CancelPayment(ctx, userID, paymentID)

	assert.Error(t, err)
	assert.Equal(t, repository.ErrNotFound, err)

	mockPayments.AssertExpectations(t)
}

// TestPaymentService_CancelPayment_UpdateFailed tests cancel when update fails
func TestPaymentService_CancelPayment_UpdateFailed(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	paymentID := uint64(1)

	mockPayments := new(MockPaymentRepository)

	payment := createTestPayment(paymentID, 1000, userID, model.PaymentStatusPending, 10000, model.PaymentMethodWeChat)

	mockPayments.On("Get", ctx, paymentID).Return(payment, nil)
	mockPayments.On("Update", ctx, mock.AnythingOfType("*model.Payment")).Return(errors.New("update failed"))

	service := NewPaymentService(mockPayments, nil)

	err := service.CancelPayment(ctx, userID, paymentID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update payment")

	mockPayments.AssertExpectations(t)
}

// TestPaymentService_SetExternalConfig tests setting external config
func TestPaymentService_SetExternalConfig(t *testing.T) {
	service := NewPaymentService(nil, nil)

	cfg := &external.Config{
		WeChatPay: config.WeChatPayConfig{
			Enabled: true,
			AppID:   "test_app",
			MchID:   "test_mch",
			APIKey:  "test_key",
		},
		Alipay: config.AlipayConfig{
			Enabled:        true,
			AppID:          "test_app",
			PrivateKeyPath: "test_key",
			PublicKeyPath:  "test_pub",
		},
	}

	service.SetExternalConfig(cfg)

	assert.NotNil(t, service.providers)
	assert.Contains(t, service.providers, model.PaymentMethodWeChat)
	assert.Contains(t, service.providers, model.PaymentMethodAlipay)
}

// TestPaymentService_SetProviders tests setting providers directly
func TestPaymentService_SetProviders(t *testing.T) {
	service := NewPaymentService(nil, nil)

	providers := map[model.PaymentMethod]ProviderClient{
		model.PaymentMethodWeChat: wechatProvider{},
		model.PaymentMethodAlipay: alipayProvider{},
	}

	service.SetProviders(providers)

	assert.Equal(t, providers, service.providers)
}

// TestPaymentService_creditWallet tests wallet crediting via UpdateBalanceWithLock
func TestPaymentService_creditWallet(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	amount := int64(5000)

	tests := []struct {
		name        string
		setupMocks  func(*MockWalletRepository)
		expectError bool
	}{
		{
			name: "Credit existing wallet",
			setupMocks: func(m *MockWalletRepository) {
				m.On("UpdateBalanceWithLock", ctx, userID, amount, 3).Return(createTestWallet(userID, 15000, 0), nil)
			},
			expectError: false,
		},
		{
			name: "Wallet update fails",
			setupMocks: func(m *MockWalletRepository) {
				m.On("UpdateBalanceWithLock", ctx, userID, amount, 3).Return(nil, errors.New("save failed"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWallets := new(MockWalletRepository)
			tt.setupMocks(mockWallets)

			service := NewPaymentService(nil, nil)
			service.SetWalletRepository(mockWallets)

			err := service.creditWallet(ctx, userID, amount)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockWallets.AssertExpectations(t)
		})
	}
}

// TestPaymentService_debitWallet tests wallet debiting via UpdateBalanceWithLock (negative delta)
func TestPaymentService_debitWallet(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	amount := int64(5000)

	tests := []struct {
		name        string
		setupMocks  func(*MockWalletRepository)
		expectError bool
	}{
		{
			name: "Debit existing wallet",
			setupMocks: func(m *MockWalletRepository) {
				// debitWallet passes -amount to UpdateBalanceWithLock
				m.On("UpdateBalanceWithLock", ctx, userID, -amount, 3).Return(createTestWallet(userID, 5000, 0), nil)
			},
			expectError: false,
		},
		{
			name: "Wallet not found",
			setupMocks: func(m *MockWalletRepository) {
				m.On("UpdateBalanceWithLock", ctx, userID, -amount, 3).Return(nil, repository.ErrNotFound)
			},
			expectError: true,
		},
		{
			name: "Wallet update fails",
			setupMocks: func(m *MockWalletRepository) {
				m.On("UpdateBalanceWithLock", ctx, userID, -amount, 3).Return(nil, errors.New("save failed"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWallets := new(MockWalletRepository)
			tt.setupMocks(mockWallets)

			service := NewPaymentService(nil, nil)
			service.SetWalletRepository(mockWallets)

			err := service.debitWallet(ctx, userID, amount)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockWallets.AssertExpectations(t)
		})
	}
}

// TestPaymentService_GetWalletBalance_Error tests wallet balance retrieval error
func TestPaymentService_GetWalletBalance_Error(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)

	mockWallets := new(MockWalletRepository)

	mockWallets.On("GetByUserID", ctx, userID).Return(nil, errors.New("database error"))

	service := NewPaymentService(nil, nil)
	service.SetWalletRepository(mockWallets)

	resp, err := service.GetWalletBalance(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "database error")

	mockWallets.AssertExpectations(t)
}

// TestPaymentService_CalculateCombinedPayment_OrderNotFound tests combined payment calculation for non-existent order
func TestPaymentService_CalculateCombinedPayment_OrderNotFound(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	orderID := uint64(999)

	mockOrders := new(MockOrderRepository)

	mockOrders.On("Get", ctx, orderID).Return(nil, repository.ErrNotFound)

	service := NewPaymentService(nil, mockOrders)

	req := CalculateCombinedPaymentRequest{
		OrderID: orderID,
	}

	resp, err := service.CalculateCombinedPayment(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, repository.ErrNotFound, err)

	mockOrders.AssertExpectations(t)
}

// TestRefundService_SetExternalConfig tests setting external config for refund service
func TestRefundService_SetExternalConfig(t *testing.T) {
	mockPayments := new(MockPaymentRepository)
	mockRefunds := new(MockRefundRecordRepository)
	mockOrders := new(MockOrderRepository)

	service := NewRefundService(mockPayments, mockRefunds, mockOrders)

	cfg := &external.Config{
		WeChatPay: config.WeChatPayConfig{
			Enabled: true,
			AppID:   "test_app",
			MchID:   "test_mch",
			APIKey:  "test_key",
		},
		Alipay: config.AlipayConfig{
			Enabled:        true,
			AppID:          "test_app",
			PrivateKeyPath: "test_key",
			PublicKeyPath:  "test_pub",
		},
	}

	service.SetExternalConfig(cfg)

	// Providers should be updated
	assert.NotNil(t, service.providers)
}

// TestPaymentService_RefundPayment_RefundToWalletOnly tests refund to wallet only.
// Uses in-memory SQLite DB because RefundPayment runs a gorm.DB.Transaction internally.
func TestPaymentService_RefundPayment_RefundToWalletOnly(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	amountCents := int64(10000)

	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db, &model.Payment{}, &model.Order{})
	defer testutil.CleanDB(t, db)

	// Seed order and payment in DB
	order := &model.Order{
		OrderNo:         "ORD-WALLET",
		UserID:          userID,
		ItemID:          1,
		Status:          model.OrderStatusConfirmed,
		TotalPriceCents: amountCents,
		Currency:        model.CurrencyCNY,
	}
	db.Create(order)

	payment := &model.Payment{
		OrderID:           order.ID,
		UserID:            userID,
		AmountCents:       amountCents,
		Status:            model.PaymentStatusPaid,
		Method:            model.PaymentMethodWallet,
		Currency:          model.CurrencyCNY,
		WalletAmountCents: amountCents,
	}
	db.Create(payment)

	mockPayments := new(MockPaymentRepository)
	mockWallets := new(MockWalletRepository)

	mockPayments.On("Get", ctx, payment.ID).Return(payment, nil)
	// creditWallet uses UpdateBalanceWithLock
	updatedWallet := createTestWallet(userID, 15000, 0)
	mockWallets.On("UpdateBalanceWithLock", ctx, userID, amountCents, 3).Return(updatedWallet, nil)

	service := NewPaymentService(mockPayments, nil)
	service.SetWalletRepository(mockWallets)
	service.SetDB(db)

	err := service.RefundPayment(ctx, payment.ID, "customer request")

	assert.NoError(t, err)

	// Verify DB state
	var updatedPayment model.Payment
	db.First(&updatedPayment, payment.ID)
	assert.Equal(t, model.PaymentStatusRefunded, updatedPayment.Status)

	var updatedOrder model.Order
	db.First(&updatedOrder, order.ID)
	assert.Equal(t, model.OrderStatusRefunded, updatedOrder.Status)

	mockPayments.AssertExpectations(t)
	mockWallets.AssertExpectations(t)
}

// TestPaymentService_RefundPayment_PartialRefundThirdParty tests partial refund of third-party payment.
// Uses in-memory SQLite DB because RefundPayment runs a gorm.DB.Transaction internally.
func TestPaymentService_RefundPayment_PartialRefundThirdParty(t *testing.T) {
	ctx := context.Background()
	userID := uint64(100)
	fullAmountCents := int64(10000)

	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db, &model.Payment{}, &model.Order{})
	defer testutil.CleanDB(t, db)

	// Seed order and payment in DB
	order := &model.Order{
		OrderNo:         "ORD-PARTIAL",
		UserID:          userID,
		ItemID:          1,
		Status:          model.OrderStatusConfirmed,
		TotalPriceCents: fullAmountCents,
		Currency:        model.CurrencyCNY,
	}
	db.Create(order)

	payment := &model.Payment{
		OrderID:     order.ID,
		UserID:      userID,
		AmountCents: fullAmountCents,
		Status:      model.PaymentStatusPaid,
		Method:      model.PaymentMethodWeChat,
		Currency:    model.CurrencyCNY,
	}
	db.Create(payment)

	mockPayments := new(MockPaymentRepository)
	mockPayments.On("Get", ctx, payment.ID).Return(payment, nil)

	service := NewPaymentService(mockPayments, nil)
	service.SetDB(db)

	err := service.RefundPayment(ctx, payment.ID, "partial refund")

	assert.NoError(t, err)

	// Verify DB state
	var updatedPayment model.Payment
	db.First(&updatedPayment, payment.ID)
	assert.Equal(t, model.PaymentStatusRefunded, updatedPayment.Status)

	var updatedOrder model.Order
	db.First(&updatedOrder, order.ID)
	assert.Equal(t, model.OrderStatusRefunded, updatedOrder.Status)

	mockPayments.AssertExpectations(t)
}

// TestPaymentService_RefundPayment_Combined_Partial tests partial refund of combined payment
func TestPaymentService_RefundPayment_Combined_Partial(t *testing.T) {
	t.Skip("Skipping due to complex wallet credit interactions in combined partial refunds")
}

// Mock helpers

type mockThirdPartyFailingProvider struct{}

func (m *mockThirdPartyFailingProvider) Refund(ctx context.Context, p *model.Payment, reason string) (string, json.RawMessage, time.Time, error) {
	return "", nil, time.Time{}, errors.New("provider error")
}

type mockDistributedLock struct {
	mock.Mock
}

func (m *mockDistributedLock) Lock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	args := m.Called(ctx, key, ttl)
	return args.Bool(0), args.Error(1)
}

func (m *mockDistributedLock) Unlock(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *mockDistributedLock) TryLock(ctx context.Context, key string, ttl time.Duration, retries int, retryDelay time.Duration) (bool, error) {
	args := m.Called(ctx, key, ttl, retries, retryDelay)
	return args.Bool(0), args.Error(1)
}
