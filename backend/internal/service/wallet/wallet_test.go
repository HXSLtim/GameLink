package wallet

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/model"
	orderimpl "gamelink/internal/repository/implementations"
	orderrepo "gamelink/internal/repository/order"
	walletrepo "gamelink/internal/repository/wallet"
	"gamelink/pkg/testutil"
)

func setupWalletTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.User{},
		&model.Order{},
		&model.Payment{},
		&model.Wallet{},
	)
	return db
}

func createWalletTestUser(t *testing.T, db *gorm.DB, phone string) *model.User {
	t.Helper()
	user := &model.User{
		Phone:        phone,
		Email:        phone + "@test.com",
		Name:         "Test User",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(user).Error)
	return user
}

func createWalletService(db *gorm.DB) *WalletService {
	walletRepo := walletrepo.NewWalletRepository(db)
	paymentRepo := orderrepo.NewPaymentRepository(db)
	orderRepo := orderimpl.NewOrderRepository(db)
	return NewWalletService(walletRepo, paymentRepo, orderRepo)
}

func TestWalletService_Recharge(t *testing.T) {
	db := setupWalletTestDB(t)
	defer testutil.CleanDB(t, db)

	user := createWalletTestUser(t, db, "13800000001")
	svc := createWalletService(db)
	ctx := context.Background()

	t.Run("充值成功_微信支付", func(t *testing.T) {
		resp, err := svc.Recharge(ctx, user.ID, RechargeRequest{
			AmountCents: 10000,
			Method:      model.PaymentMethodWeChat,
		})
		require.NoError(t, err)
		assert.NotZero(t, resp.OrderID)
		assert.NotZero(t, resp.PaymentID)
		assert.Equal(t, int64(10000), resp.Balance)
	})

	t.Run("充值成功_支付宝", func(t *testing.T) {
		user2 := createWalletTestUser(t, db, "13800000002")
		resp, err := svc.Recharge(ctx, user2.ID, RechargeRequest{
			AmountCents: 5000,
			Method:      model.PaymentMethodAlipay,
		})
		require.NoError(t, err)
		assert.NotZero(t, resp.OrderID)
		assert.Equal(t, int64(5000), resp.Balance)
	})

	t.Run("多次充值累加余额", func(t *testing.T) {
		user3 := createWalletTestUser(t, db, "13800000003")

		// 第一次充值
		resp1, err := svc.Recharge(ctx, user3.ID, RechargeRequest{
			AmountCents: 10000,
			Method:      model.PaymentMethodWeChat,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(10000), resp1.Balance)

		// 第二次充值
		resp2, err := svc.Recharge(ctx, user3.ID, RechargeRequest{
			AmountCents: 5000,
			Method:      model.PaymentMethodAlipay,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(15000), resp2.Balance)
	})

	t.Run("充值金额为零应失败", func(t *testing.T) {
		_, err := svc.Recharge(ctx, user.ID, RechargeRequest{
			AmountCents: 0,
			Method:      model.PaymentMethodWeChat,
		})
		assert.ErrorIs(t, err, ErrInvalidAmount)
	})

	t.Run("充值金额为负数应失败", func(t *testing.T) {
		_, err := svc.Recharge(ctx, user.ID, RechargeRequest{
			AmountCents: -100,
			Method:      model.PaymentMethodWeChat,
		})
		assert.ErrorIs(t, err, ErrInvalidAmount)
	})

	t.Run("大额充值成功", func(t *testing.T) {
		user4 := createWalletTestUser(t, db, "13800000004")
		resp, err := svc.Recharge(ctx, user4.ID, RechargeRequest{
			AmountCents: 10000000, // 10万元
			Method:      model.PaymentMethodWeChat,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(10000000), resp.Balance)
	})
}

func TestWalletService_GetBalance(t *testing.T) {
	db := setupWalletTestDB(t)
	defer testutil.CleanDB(t, db)

	user := createWalletTestUser(t, db, "13800000010")
	svc := createWalletService(db)
	ctx := context.Background()

	t.Run("新用户余额为零", func(t *testing.T) {
		wallet, err := svc.GetBalance(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(0), wallet.BalanceCents)
		assert.Equal(t, int64(0), wallet.FrozenCents)
	})

	t.Run("充值后查询余额", func(t *testing.T) {
		_, err := svc.Recharge(ctx, user.ID, RechargeRequest{
			AmountCents: 10000,
			Method:      model.PaymentMethodWeChat,
		})
		require.NoError(t, err)

		wallet, err := svc.GetBalance(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(10000), wallet.BalanceCents)
	})

	t.Run("不存在的用户返回零余额", func(t *testing.T) {
		wallet, err := svc.GetBalance(ctx, 99999)
		require.NoError(t, err)
		assert.Equal(t, int64(0), wallet.BalanceCents)
		assert.Equal(t, uint64(99999), wallet.UserID)
	})
}

func TestWalletService_RechargeCreatesOrderAndPayment(t *testing.T) {
	db := setupWalletTestDB(t)
	defer testutil.CleanDB(t, db)

	user := createWalletTestUser(t, db, "13800000020")
	svc := createWalletService(db)
	ctx := context.Background()

	resp, err := svc.Recharge(ctx, user.ID, RechargeRequest{
		AmountCents: 10000,
		Method:      model.PaymentMethodWeChat,
	})
	require.NoError(t, err)

	// 验证订单创建
	var order model.Order
	require.NoError(t, db.First(&order, resp.OrderID).Error)
	assert.Equal(t, user.ID, order.UserID)
	assert.Equal(t, model.OrderStatusCompleted, order.Status)
	assert.Equal(t, "Wallet Recharge", order.Title)
	assert.Equal(t, int64(10000), order.TotalPriceCents)
	assert.NotNil(t, order.CompletedAt)

	// 验证支付记录创建
	var payment model.Payment
	require.NoError(t, db.First(&payment, resp.PaymentID).Error)
	assert.Equal(t, order.ID, payment.OrderID)
	assert.Equal(t, user.ID, payment.UserID)
	assert.Equal(t, model.PaymentMethodWeChat, payment.Method)
	assert.Equal(t, int64(10000), payment.AmountCents)
	assert.Equal(t, model.PaymentStatusPaid, payment.Status)
	assert.NotNil(t, payment.PaidAt)
}

// 原有的单元测试保留
func TestRechargeRequest_Fields(t *testing.T) {
	req := RechargeRequest{
		AmountCents: 10000,
		Method:      model.PaymentMethodAlipay,
	}

	assert.Equal(t, int64(10000), req.AmountCents)
	assert.Equal(t, model.PaymentMethodAlipay, req.Method)
}

func TestRechargeResponse_Fields(t *testing.T) {
	resp := RechargeResponse{
		OrderID:   1,
		PaymentID: 2,
		Balance:   15000,
	}

	assert.Equal(t, uint64(1), resp.OrderID)
	assert.Equal(t, uint64(2), resp.PaymentID)
	assert.Equal(t, int64(15000), resp.Balance)
}

func TestErrInvalidAmount(t *testing.T) {
	assert.Equal(t, "invalid amount", ErrInvalidAmount.Error())
}

func TestRechargeRequest_Validation(t *testing.T) {
	tests := []struct {
		name        string
		amountCents int64
		isValid     bool
	}{
		{"positive amount", 1000, true},
		{"zero amount", 0, false},
		{"negative amount", -100, false},
		{"large amount", 1000000, true},
		{"minimum valid", 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := RechargeRequest{
				AmountCents: tt.amountCents,
				Method:      model.PaymentMethodAlipay,
			}
			isValid := req.AmountCents > 0
			assert.Equal(t, tt.isValid, isValid)
		})
	}
}

func TestPaymentMethods(t *testing.T) {
	// Test that payment methods are valid
	methods := []model.PaymentMethod{
		model.PaymentMethodAlipay,
		model.PaymentMethodWeChat,
	}

	for _, method := range methods {
		assert.NotEmpty(t, string(method))
	}
}

func TestRechargeResponse_ZeroValues(t *testing.T) {
	resp := RechargeResponse{}

	assert.Equal(t, uint64(0), resp.OrderID)
	assert.Equal(t, uint64(0), resp.PaymentID)
	assert.Equal(t, int64(0), resp.Balance)
}
