package payment

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	orderimpl "gamelink/internal/repository/implementations"
	orderrepo "gamelink/internal/repository/order"
	paymentrepo "gamelink/internal/repository/payment"
	"gamelink/internal/repository/user"
	"gamelink/pkg/testutil"
)

func setupPaymentTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.User{},
		&model.Player{},
		&model.Game{},
		&model.Order{},
		&model.Payment{},
		&model.Wallet{},
	)
	return db
}

func createPaymentTestData(t *testing.T, db *gorm.DB) (customer *model.User, player *model.Player, gameModel *model.Game, order *model.Order) {
	t.Helper()

	// 创建用户
	customer = &model.User{
		Phone:        "13800000001",
		Email:        "customer@test.com",
		Name:         "Customer",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(customer).Error)

	// 创建陪玩师用户
	playerUser := &model.User{
		Phone:        "13800000002",
		Email:        "player@test.com",
		Name:         "Player",
		Role:         model.RolePlayer,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(playerUser).Error)

	// 创建游戏
	gameModel = &model.Game{Key: "lol", Name: "英雄联盟", Category: "moba"}
	require.NoError(t, db.Create(gameModel).Error)

	// 创建陪玩师
	player = &model.Player{
		UserID:             playerUser.ID,
		Nickname:           "Pro Player",
		MainGameID:         gameModel.ID,
		HourlyRateCents:    5000,
		VerificationStatus: model.VerificationVerified,
	}
	require.NoError(t, db.Create(player).Error)

	// 创建订单
	scheduledStart := time.Now().Add(24 * time.Hour)
	order = &model.Order{
		UserID:          customer.ID,
		ItemID:          gameModel.ID,
		Title:           "测试订单",
		Status:          model.OrderStatusPending,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
	}
	order.SetPlayerID(player.ID)
	order.SetGameID(gameModel.ID)
	require.NoError(t, db.Create(order).Error)

	return
}

func createPaymentService(db *gorm.DB) *PaymentService {
	paymentRepo := orderrepo.NewPaymentRepository(db)
	orderRepo := orderimpl.NewOrderRepository(db)
	return NewPaymentService(paymentRepo, orderRepo)
}

func TestPaymentService_CreatePayment(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, _, order := createPaymentTestData(t, db)
	svc := createPaymentService(db)

	t.Run("创建微信支付成功", func(t *testing.T) {
		resp, err := svc.CreatePayment(context.Background(), customer.ID, CreatePaymentRequest{
			OrderID: order.ID,
			Method:  model.PaymentMethodWeChat,
		})
		require.NoError(t, err)
		assert.NotZero(t, resp.PaymentID)
		assert.NotNil(t, resp.PayInfo)
		assert.Contains(t, resp.PayInfo, "prepay_id")
	})

	t.Run("重复支付应失败", func(t *testing.T) {
		// 创建新订单
		scheduledStart := time.Now().Add(24 * time.Hour)
		newOrder := &model.Order{
			UserID:          customer.ID,
			ItemID:          1,
			Title:           "新订单",
			Status:          model.OrderStatusPending,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
		}
		require.NoError(t, db.Create(newOrder).Error)

		// 第一次支付
		_, err := svc.CreatePayment(context.Background(), customer.ID, CreatePaymentRequest{
			OrderID: newOrder.ID,
			Method:  model.PaymentMethodAlipay,
		})
		require.NoError(t, err)

		// 第二次支付应失败
		_, err = svc.CreatePayment(context.Background(), customer.ID, CreatePaymentRequest{
			OrderID: newOrder.ID,
			Method:  model.PaymentMethodAlipay,
		})
		assert.Error(t, err)
	})

	t.Run("非订单所有者无法支付", func(t *testing.T) {
		// 创建另一个用户
		otherUser := &model.User{
			Phone:        "13800000003",
			Email:        "other@test.com",
			Name:         "Other",
			Role:         model.RoleUser,
			Status:       model.UserStatusActive,
			PasswordHash: "hashed",
		}
		require.NoError(t, db.Create(otherUser).Error)

		// 创建新订单
		scheduledStart := time.Now().Add(24 * time.Hour)
		newOrder := &model.Order{
			UserID:          customer.ID,
			ItemID:          1,
			Title:           "另一个订单",
			Status:          model.OrderStatusPending,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
		}
		require.NoError(t, db.Create(newOrder).Error)

		_, err := svc.CreatePayment(context.Background(), otherUser.ID, CreatePaymentRequest{
			OrderID: newOrder.ID,
			Method:  model.PaymentMethodWeChat,
		})
		assert.Error(t, err)
	})
}

func TestPaymentService_GetPaymentStatus(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, _, order := createPaymentTestData(t, db)
	svc := createPaymentService(db)

	// 创建支付
	resp, err := svc.CreatePayment(context.Background(), customer.ID, CreatePaymentRequest{
		OrderID: order.ID,
		Method:  model.PaymentMethodWeChat,
	})
	require.NoError(t, err)

	t.Run("获取支付状态成功", func(t *testing.T) {
		status, err := svc.GetPaymentStatus(context.Background(), resp.PaymentID)
		require.NoError(t, err)
		assert.Equal(t, resp.PaymentID, status.PaymentID)
		assert.Equal(t, order.ID, status.OrderID)
		// Mock支付会自动标记为已支付
		assert.Equal(t, model.PaymentStatusPaid, status.Status)
	})

	t.Run("支付不存在", func(t *testing.T) {
		_, err := svc.GetPaymentStatus(context.Background(), 99999)
		assert.ErrorIs(t, err, ErrNotFound)
	})
}

func TestPaymentService_CancelPayment(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, _, _ := createPaymentTestData(t, db)
	svc := createPaymentService(db)

	t.Run("取消待支付的支付", func(t *testing.T) {
		// 直接创建一个pending状态的支付记录
		payment := &model.Payment{
			OrderID:     1,
			UserID:      customer.ID,
			Method:      model.PaymentMethodWeChat,
			AmountCents: 10000,
			Status:      model.PaymentStatusPending,
		}
		require.NoError(t, db.Create(payment).Error)

		err := svc.CancelPayment(context.Background(), customer.ID, payment.ID)
		require.NoError(t, err)

		// 验证状态
		var updated model.Payment
		require.NoError(t, db.First(&updated, payment.ID).Error)
		assert.Equal(t, model.PaymentStatusFailed, updated.Status)
	})

	t.Run("无法取消已支付的支付", func(t *testing.T) {
		now := time.Now()
		payment := &model.Payment{
			OrderID:     2,
			UserID:      customer.ID,
			Method:      model.PaymentMethodAlipay,
			AmountCents: 10000,
			Status:      model.PaymentStatusPaid,
			PaidAt:      &now,
		}
		require.NoError(t, db.Create(payment).Error)

		err := svc.CancelPayment(context.Background(), customer.ID, payment.ID)
		assert.Error(t, err)
	})

	t.Run("非支付所有者无法取消", func(t *testing.T) {
		payment := &model.Payment{
			OrderID:     3,
			UserID:      customer.ID,
			Method:      model.PaymentMethodWeChat,
			AmountCents: 10000,
			Status:      model.PaymentStatusPending,
		}
		require.NoError(t, db.Create(payment).Error)

		// 创建另一个用户
		otherUser := &model.User{
			Phone:        "13800000004",
			Email:        "other2@test.com",
			Name:         "Other2",
			Role:         model.RoleUser,
			Status:       model.UserStatusActive,
			PasswordHash: "hashed",
		}
		require.NoError(t, db.Create(otherUser).Error)

		err := svc.CancelPayment(context.Background(), otherUser.ID, payment.ID)
		assert.Error(t, err)
	})
}

func TestPaymentService_HandlePaymentCallback(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, gameModel, _ := createPaymentTestData(t, db)
	svc := createPaymentService(db)

	t.Run("处理支付回调成功", func(t *testing.T) {
		// 创建订单
		scheduledStart := time.Now().Add(24 * time.Hour)
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          gameModel.ID,
			Title:           "回调测试订单",
			Status:          model.OrderStatusPending,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
		}
		require.NoError(t, db.Create(order).Error)

		// 创建pending状态的支付
		payment := &model.Payment{
			OrderID:     order.ID,
			UserID:      customer.ID,
			Method:      model.PaymentMethodWeChat,
			AmountCents: 10000,
			Status:      model.PaymentStatusPending,
		}
		require.NoError(t, db.Create(payment).Error)

		// 处理回调
		err := svc.HandlePaymentCallback(context.Background(), "wechat", map[string]interface{}{
			"payment_id":   float64(payment.ID),
			"amount_cents": int64(10000),
			"trade_no":     "wx_trade_123",
		})
		require.NoError(t, err)

		// 验证支付状态
		var updatedPayment model.Payment
		require.NoError(t, db.First(&updatedPayment, payment.ID).Error)
		assert.Equal(t, model.PaymentStatusPaid, updatedPayment.Status)
		assert.NotNil(t, updatedPayment.PaidAt)

		// 验证订单状态
		var updatedOrder model.Order
		require.NoError(t, db.First(&updatedOrder, order.ID).Error)
		assert.Equal(t, model.OrderStatusConfirmed, updatedOrder.Status)
	})

	t.Run("重复回调应幂等", func(t *testing.T) {
		// 创建已支付的支付记录
		now := time.Now()
		payment := &model.Payment{
			OrderID:     100,
			UserID:      customer.ID,
			Method:      model.PaymentMethodAlipay,
			AmountCents: 10000,
			Status:      model.PaymentStatusPaid,
			PaidAt:      &now,
		}
		require.NoError(t, db.Create(payment).Error)

		// 重复回调应该成功（幂等）
		err := svc.HandlePaymentCallback(context.Background(), "alipay", map[string]interface{}{
			"payment_id": float64(payment.ID),
		})
		assert.NoError(t, err)
	})

	t.Run("提供商不匹配应失败", func(t *testing.T) {
		payment := &model.Payment{
			OrderID:     101,
			UserID:      customer.ID,
			Method:      model.PaymentMethodWeChat,
			AmountCents: 10000,
			Status:      model.PaymentStatusPending,
		}
		require.NoError(t, db.Create(payment).Error)

		err := svc.HandlePaymentCallback(context.Background(), "alipay", map[string]interface{}{
			"payment_id": float64(payment.ID),
		})
		assert.Error(t, err)
	})
}

func TestPaymentService_RefundPayment(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, gameModel, _ := createPaymentTestData(t, db)
	svc := createPaymentService(db)

	// 设置钱包仓储
	walletRepo := user.NewWalletRepository(db)
	svc.SetWalletRepository(walletRepo)

	t.Run("退款成功", func(t *testing.T) {
		// 创建订单
		scheduledStart := time.Now().Add(24 * time.Hour)
		order := &model.Order{
			UserID:          customer.ID,
			ItemID:          gameModel.ID,
			Title:           "退款测试订单",
			Status:          model.OrderStatusConfirmed,
			UnitPriceCents:  5000,
			TotalPriceCents: 10000,
			ScheduledStart:  &scheduledStart,
		}
		require.NoError(t, db.Create(order).Error)

		// 创建已支付的支付记录
		now := time.Now()
		payment := &model.Payment{
			OrderID:     order.ID,
			UserID:      customer.ID,
			Method:      model.PaymentMethodWeChat,
			AmountCents: 10000,
			Status:      model.PaymentStatusPaid,
			PaidAt:      &now,
		}
		require.NoError(t, db.Create(payment).Error)

		// 执行退款
		err := svc.RefundPayment(context.Background(), payment.ID, "用户申请退款")
		require.NoError(t, err)

		// 验证支付状态
		var updatedPayment model.Payment
		require.NoError(t, db.First(&updatedPayment, payment.ID).Error)
		assert.Equal(t, model.PaymentStatusRefunded, updatedPayment.Status)
		assert.NotNil(t, updatedPayment.RefundedAt)

		// 验证订单状态
		var updatedOrder model.Order
		require.NoError(t, db.First(&updatedOrder, order.ID).Error)
		assert.Equal(t, model.OrderStatusRefunded, updatedOrder.Status)
		assert.Equal(t, int64(10000), updatedOrder.RefundAmountCents)

		// 验证钱包余额
		wallet, err := walletRepo.GetByUserID(context.Background(), customer.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(10000), wallet.BalanceCents)
	})

	t.Run("未支付的订单无法退款", func(t *testing.T) {
		payment := &model.Payment{
			OrderID:     200,
			UserID:      customer.ID,
			Method:      model.PaymentMethodAlipay,
			AmountCents: 10000,
			Status:      model.PaymentStatusPending,
		}
		require.NoError(t, db.Create(payment).Error)

		err := svc.RefundPayment(context.Background(), payment.ID, "测试退款")
		assert.Error(t, err)
	})
}

func TestPaymentService_List(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, _, _ := createPaymentTestData(t, db)
	svc := createPaymentService(db)

	// 创建多个支付记录
	for i := 0; i < 5; i++ {
		payment := &model.Payment{
			OrderID:     uint64(i + 1),
			UserID:      customer.ID,
			Method:      model.PaymentMethodWeChat,
			AmountCents: int64((i + 1) * 1000),
			Status:      model.PaymentStatusPending,
		}
		require.NoError(t, db.Create(payment).Error)
	}

	t.Run("获取支付列表", func(t *testing.T) {
		payments, total, err := svc.List(context.Background(), repository.PaymentListOptions{
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, payments, 5)
	})

	t.Run("分页获取", func(t *testing.T) {
		payments, total, err := svc.List(context.Background(), repository.PaymentListOptions{
			Page:     1,
			PageSize: 2,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, payments, 2)
	})
}

func TestGenerateMockPayInfo(t *testing.T) {
	svc := &PaymentService{}

	t.Run("微信支付参数", func(t *testing.T) {
		payInfo := svc.generateMockPayInfo(1, model.PaymentMethodWeChat, 10000)
		assert.Contains(t, payInfo, "prepay_id")
		assert.Contains(t, payInfo, "code_url")
		assert.Equal(t, int64(10000), payInfo["amountCents"])
	})

	t.Run("支付宝支付参数", func(t *testing.T) {
		payInfo := svc.generateMockPayInfo(2, model.PaymentMethodAlipay, 20000)
		assert.Contains(t, payInfo, "trade_no")
		assert.Contains(t, payInfo, "qr_code")
		assert.Equal(t, int64(20000), payInfo["amountCents"])
	})
}

func TestPaymentService_GetPaymentByOrderID(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, _, order := createPaymentTestData(t, db)
	svc := createPaymentService(db)

	// 创建支付记录
	payment := &model.Payment{
		OrderID:     order.ID,
		UserID:      customer.ID,
		Method:      model.PaymentMethodWeChat,
		AmountCents: 10000,
		Status:      model.PaymentStatusPaid,
	}
	require.NoError(t, db.Create(payment).Error)

	t.Run("按订单ID获取支付记录", func(t *testing.T) {
		orderIDPtr := &order.ID
		payments, total, err := svc.List(context.Background(), repository.PaymentListOptions{
			OrderID:  orderIDPtr,
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, payments, 1)
		assert.Equal(t, order.ID, payments[0].OrderID)
	})
}

func TestPaymentService_ListByStatus(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, _, _ := createPaymentTestData(t, db)
	svc := createPaymentService(db)

	// 创建不同状态的支付记录
	statuses := []model.PaymentStatus{
		model.PaymentStatusPending,
		model.PaymentStatusPaid,
		model.PaymentStatusPaid,
		model.PaymentStatusFailed,
	}

	for i, status := range statuses {
		payment := &model.Payment{
			OrderID:     uint64(i + 100),
			UserID:      customer.ID,
			Method:      model.PaymentMethodWeChat,
			AmountCents: int64((i + 1) * 1000),
			Status:      status,
		}
		require.NoError(t, db.Create(payment).Error)
	}

	t.Run("获取所有支付记录", func(t *testing.T) {
		payments, total, err := svc.List(context.Background(), repository.PaymentListOptions{
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(4), total)
		assert.Len(t, payments, 4)
	})
}

func TestPaymentService_ListByUserID(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, _, _ := createPaymentTestData(t, db)
	svc := createPaymentService(db)

	// 创建另一个用户
	otherUser := &model.User{
		Phone:        "13800000099",
		Email:        "other99@test.com",
		Name:         "Other User",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(otherUser).Error)

	// 为两个用户创建支付记录
	for i := 0; i < 3; i++ {
		payment := &model.Payment{
			OrderID:     uint64(i + 200),
			UserID:      customer.ID,
			Method:      model.PaymentMethodWeChat,
			AmountCents: int64((i + 1) * 1000),
			Status:      model.PaymentStatusPending,
		}
		require.NoError(t, db.Create(payment).Error)
	}

	for i := 0; i < 2; i++ {
		payment := &model.Payment{
			OrderID:     uint64(i + 300),
			UserID:      otherUser.ID,
			Method:      model.PaymentMethodAlipay,
			AmountCents: int64((i + 1) * 2000),
			Status:      model.PaymentStatusPending,
		}
		require.NoError(t, db.Create(payment).Error)
	}

	t.Run("按用户ID筛选支付记录", func(t *testing.T) {
		payments, total, err := svc.List(context.Background(), repository.PaymentListOptions{
			UserID:   &customer.ID,
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, payments, 3)
	})
}

func TestPaymentService_ListPagination(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, _, _ := createPaymentTestData(t, db)
	svc := createPaymentService(db)

	// 创建多个支付记录
	for i := 0; i < 25; i++ {
		payment := &model.Payment{
			OrderID:     uint64(i + 400),
			UserID:      customer.ID,
			Method:      model.PaymentMethodWeChat,
			AmountCents: int64((i + 1) * 100),
			Status:      model.PaymentStatusPending,
		}
		require.NoError(t, db.Create(payment).Error)
	}

	t.Run("分页获取第一页", func(t *testing.T) {
		payments, total, err := svc.List(context.Background(), repository.PaymentListOptions{
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(25), total)
		assert.Len(t, payments, 10)
	})

	t.Run("分页获取第二页", func(t *testing.T) {
		payments, total, err := svc.List(context.Background(), repository.PaymentListOptions{
			Page:     2,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(25), total)
		assert.Len(t, payments, 10)
	})

	t.Run("分页获取最后一页", func(t *testing.T) {
		payments, total, err := svc.List(context.Background(), repository.PaymentListOptions{
			Page:     3,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(25), total)
		assert.Len(t, payments, 5)
	})

	t.Run("默认分页参数", func(t *testing.T) {
		payments, _, err := svc.List(context.Background(), repository.PaymentListOptions{
			Page:     0,
			PageSize: 0,
		})
		require.NoError(t, err)
		assert.Len(t, payments, 20) // 默认 pageSize 是 20
	})

	t.Run("超大分页参数被限制", func(t *testing.T) {
		payments, _, err := svc.List(context.Background(), repository.PaymentListOptions{
			Page:     1,
			PageSize: 200,
		})
		require.NoError(t, err)
		assert.Len(t, payments, 20) // 最大 pageSize 是 100，但这里只有25条
	})
}

func TestPaymentService_Errors(t *testing.T) {
	t.Run("错误变量定义正确", func(t *testing.T) {
		assert.NotNil(t, ErrNotFound)
		assert.NotNil(t, ErrValidation)
		assert.NotNil(t, ErrOrderAlreadyPaid)
		assert.NotNil(t, ErrInvalidOrderStatus)
	})
}

func TestPaymentService_CancelPaymentNotFound(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, _, _ := createPaymentTestData(t, db)
	svc := createPaymentService(db)

	t.Run("取消不存在的支付", func(t *testing.T) {
		err := svc.CancelPayment(context.Background(), customer.ID, 99999)
		assert.Error(t, err)
	})
}

func TestPaymentService_HandlePaymentCallbackMissingPaymentID(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	svc := createPaymentService(db)

	t.Run("回调缺少payment_id", func(t *testing.T) {
		err := svc.HandlePaymentCallback(context.Background(), "wechat", map[string]interface{}{
			"amount_cents": int64(10000),
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing payment_id")
	})
}

func TestPaymentService_HandlePaymentCallbackAmountMismatch(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, gameModel, _ := createPaymentTestData(t, db)
	svc := createPaymentService(db)

	// 创建订单
	scheduledStart := time.Now().Add(24 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          gameModel.ID,
		Title:           "金额不匹配测试",
		Status:          model.OrderStatusPending,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
	}
	require.NoError(t, db.Create(order).Error)

	// 创建支付
	payment := &model.Payment{
		OrderID:     order.ID,
		UserID:      customer.ID,
		Method:      model.PaymentMethodWeChat,
		AmountCents: 10000,
		Status:      model.PaymentStatusPending,
	}
	require.NoError(t, db.Create(payment).Error)

	t.Run("回调金额不匹配", func(t *testing.T) {
		err := svc.HandlePaymentCallback(context.Background(), "wechat", map[string]interface{}{
			"payment_id":   float64(payment.ID),
			"amount_cents": int64(5000), // 金额不匹配
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "amount mismatch")
	})
}

func TestPaymentService_RefundPaymentNotFound(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	svc := createPaymentService(db)

	t.Run("退款不存在的支付", func(t *testing.T) {
		err := svc.RefundPayment(context.Background(), 99999, "测试退款")
		assert.Error(t, err)
	})
}

func TestPaymentService_CreatePaymentOrderNotFound(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, _, _ := createPaymentTestData(t, db)
	svc := createPaymentService(db)

	t.Run("订单不存在", func(t *testing.T) {
		_, err := svc.CreatePayment(context.Background(), customer.ID, CreatePaymentRequest{
			OrderID: 99999,
			Method:  model.PaymentMethodWeChat,
		})
		assert.Error(t, err)
	})
}

func TestPaymentService_CreatePaymentInvalidOrderStatus(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, gameModel, _ := createPaymentTestData(t, db)
	svc := createPaymentService(db)

	// 创建已完成的订单
	scheduledStart := time.Now().Add(24 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          gameModel.ID,
		Title:           "已完成订单",
		Status:          model.OrderStatusCompleted, // 非pending状态
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
	}
	require.NoError(t, db.Create(order).Error)

	t.Run("非pending状态订单无法支付", func(t *testing.T) {
		_, err := svc.CreatePayment(context.Background(), customer.ID, CreatePaymentRequest{
			OrderID: order.ID,
			Method:  model.PaymentMethodWeChat,
		})
		assert.Error(t, err)
	})
}

func TestPaymentService_SetDistributedLock(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	svc := createPaymentService(db)

	t.Run("设置分布式锁", func(t *testing.T) {
		// 这里只测试方法不会panic
		svc.SetDistributedLock(nil)
	})
}

func TestPaymentService_SetWalletRepository(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	svc := createPaymentService(db)
	walletRepo := user.NewWalletRepository(db)

	t.Run("设置钱包仓储", func(t *testing.T) {
		svc.SetWalletRepository(walletRepo)
		// 验证设置成功
		assert.NotNil(t, svc.wallets)
	})
}

func TestPaymentService_HandlePaymentCallbackWithTradeNo(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, gameModel, _ := createPaymentTestData(t, db)
	svc := createPaymentService(db)

	// 创建订单
	scheduledStart := time.Now().Add(24 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          gameModel.ID,
		Title:           "交易号测试",
		Status:          model.OrderStatusPending,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
	}
	require.NoError(t, db.Create(order).Error)

	// 创建支付
	payment := &model.Payment{
		OrderID:     order.ID,
		UserID:      customer.ID,
		Method:      model.PaymentMethodWeChat,
		AmountCents: 10000,
		Status:      model.PaymentStatusPending,
	}
	require.NoError(t, db.Create(payment).Error)

	t.Run("回调包含交易号", func(t *testing.T) {
		err := svc.HandlePaymentCallback(context.Background(), "wechat", map[string]interface{}{
			"payment_id": float64(payment.ID),
			"trade_no":   "wx_trade_12345",
		})
		require.NoError(t, err)

		// 验证交易号
		var updated model.Payment
		require.NoError(t, db.First(&updated, payment.ID).Error)
		assert.Equal(t, "wx_trade_12345", updated.ProviderTradeNo)
	})
}

func TestPaymentService_HandlePaymentCallbackWithoutTradeNo(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, gameModel, _ := createPaymentTestData(t, db)
	svc := createPaymentService(db)

	// 创建订单
	scheduledStart := time.Now().Add(24 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          gameModel.ID,
		Title:           "无交易号测试",
		Status:          model.OrderStatusPending,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
	}
	require.NoError(t, db.Create(order).Error)

	// 创建支付
	payment := &model.Payment{
		OrderID:     order.ID,
		UserID:      customer.ID,
		Method:      model.PaymentMethodAlipay,
		AmountCents: 10000,
		Status:      model.PaymentStatusPending,
	}
	require.NoError(t, db.Create(payment).Error)

	t.Run("回调不包含交易号时自动生成", func(t *testing.T) {
		err := svc.HandlePaymentCallback(context.Background(), "alipay", map[string]interface{}{
			"payment_id": float64(payment.ID),
		})
		require.NoError(t, err)

		// 验证自动生成的交易号
		var updated model.Payment
		require.NoError(t, db.First(&updated, payment.ID).Error)
		assert.Contains(t, updated.ProviderTradeNo, "alipay_")
	})
}

func TestPaymentService_RefundWithWalletCredit(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, gameModel, _ := createPaymentTestData(t, db)
	svc := createPaymentService(db)

	// 设置钱包仓储
	walletRepo := user.NewWalletRepository(db)
	svc.SetWalletRepository(walletRepo)

	// 创建用户钱包（初始余额为0）
	wallet := &model.Wallet{
		UserID:       customer.ID,
		BalanceCents: 0,
	}
	require.NoError(t, db.Create(wallet).Error)

	// 创建订单
	scheduledStart := time.Now().Add(24 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          gameModel.ID,
		Title:           "钱包退款测试",
		Status:          model.OrderStatusConfirmed,
		UnitPriceCents:  5000,
		TotalPriceCents: 15000,
		ScheduledStart:  &scheduledStart,
	}
	require.NoError(t, db.Create(order).Error)

	// 创建已支付的支付记录
	now := time.Now()
	payment := &model.Payment{
		OrderID:     order.ID,
		UserID:      customer.ID,
		Method:      model.PaymentMethodWeChat,
		AmountCents: 15000,
		Status:      model.PaymentStatusPaid,
		PaidAt:      &now,
	}
	require.NoError(t, db.Create(payment).Error)

	t.Run("退款后钱包余额增加", func(t *testing.T) {
		err := svc.RefundPayment(context.Background(), payment.ID, "测试退款")
		require.NoError(t, err)

		// 验证钱包余额
		updatedWallet, err := walletRepo.GetByUserID(context.Background(), customer.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(15000), updatedWallet.BalanceCents)
	})
}

func TestPaymentService_GetPaymentStatusNotFound(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	svc := createPaymentService(db)

	t.Run("获取不存在的支付状态", func(t *testing.T) {
		_, err := svc.GetPaymentStatus(context.Background(), 99999)
		assert.ErrorIs(t, err, ErrNotFound)
	})
}

// ============================================
// Additional tests to improve coverage to 80%+
// ============================================

func TestPaymentService_RoutePaymentNotInitialized(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	svc := createPaymentService(db)

	t.Run("route payment without routing engine", func(t *testing.T) {
		// routingEngine is nil by default
		_, err := svc.GetPaymentRoutingLog(context.Background(), 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "routing engine not initialized")
	})
}

func TestPaymentService_SetRoutingEngine(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	svc := createPaymentService(db)

	t.Run("set routing engine to nil", func(t *testing.T) {
		svc.SetRoutingEngine(nil)
		assert.Nil(t, svc.routingEngine)
	})
}

func TestPaymentService_CreditWallet(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, _, _ := createPaymentTestData(t, db)
	svc := createPaymentService(db)

	// Set wallet repository
	walletRepo := user.NewWalletRepository(db)
	svc.SetWalletRepository(walletRepo)

	t.Run("credit wallet for existing user", func(t *testing.T) {
		// Create wallet
		wallet := &model.Wallet{
			UserID:       customer.ID,
			BalanceCents: 5000,
		}
		require.NoError(t, db.Create(wallet).Error)

		err := svc.creditWallet(context.Background(), customer.ID, 3000)
		require.NoError(t, err)

		// Verify balance
		updatedWallet, err := walletRepo.GetByUserID(context.Background(), customer.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(8000), updatedWallet.BalanceCents)
	})

	t.Run("credit wallet for new user", func(t *testing.T) {
		// Create new user without wallet
		newUser := &model.User{
			Phone:        "13800000088",
			Email:        "newuser@test.com",
			Name:         "New User",
			Role:         model.RoleUser,
			Status:       model.UserStatusActive,
			PasswordHash: "hashed",
		}
		require.NoError(t, db.Create(newUser).Error)

		err := svc.creditWallet(context.Background(), newUser.ID, 5000)
		require.NoError(t, err)

		// Verify wallet was created
		wallet, err := walletRepo.GetByUserID(context.Background(), newUser.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(5000), wallet.BalanceCents)
	})
}

func TestPaymentService_HandlePaymentCallbackWithUint64PaymentID(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, gameModel, _ := createPaymentTestData(t, db)
	svc := createPaymentService(db)

	// Create order
	scheduledStart := time.Now().Add(24 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          gameModel.ID,
		Title:           "uint64测试",
		Status:          model.OrderStatusPending,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
	}
	require.NoError(t, db.Create(order).Error)

	// Create payment
	payment := &model.Payment{
		OrderID:     order.ID,
		UserID:      customer.ID,
		Method:      model.PaymentMethodWeChat,
		AmountCents: 10000,
		Status:      model.PaymentStatusPending,
	}
	require.NoError(t, db.Create(payment).Error)

	t.Run("callback with uint64 payment_id", func(t *testing.T) {
		err := svc.HandlePaymentCallback(context.Background(), "wechat", map[string]interface{}{
			"payment_id": uint64(payment.ID),
		})
		require.NoError(t, err)
	})
}

func TestPaymentService_RefundPaymentAlreadyRefunded(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, _, _ := createPaymentTestData(t, db)
	svc := createPaymentService(db)

	// Create already refunded payment
	refundedAt := time.Now()
	payment := &model.Payment{
		OrderID:     500,
		UserID:      customer.ID,
		Method:      model.PaymentMethodWeChat,
		AmountCents: 10000,
		Status:      model.PaymentStatusRefunded,
		RefundedAt:  &refundedAt,
	}
	require.NoError(t, db.Create(payment).Error)

	t.Run("cannot refund already refunded payment", func(t *testing.T) {
		err := svc.RefundPayment(context.Background(), payment.ID, "重复退款")
		assert.Error(t, err)
	})
}

func TestPaymentService_RefundPaymentFailed(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, _, _ := createPaymentTestData(t, db)
	svc := createPaymentService(db)

	// Create failed payment
	payment := &model.Payment{
		OrderID:     501,
		UserID:      customer.ID,
		Method:      model.PaymentMethodWeChat,
		AmountCents: 10000,
		Status:      model.PaymentStatusFailed,
	}
	require.NoError(t, db.Create(payment).Error)

	t.Run("cannot refund failed payment", func(t *testing.T) {
		err := svc.RefundPayment(context.Background(), payment.ID, "退款失败支付")
		assert.Error(t, err)
	})
}

func TestPaymentService_CreatePaymentWithDistributedLock(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, gameModel, _ := createPaymentTestData(t, db)
	svc := createPaymentService(db)

	// Create a mock distributed lock
	mockLock := &mockPaymentDistributedLock{
		tryLockResult: true,
	}
	svc.SetDistributedLock(mockLock)

	// Create order
	scheduledStart := time.Now().Add(24 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          gameModel.ID,
		Title:           "带锁支付订单",
		Status:          model.OrderStatusPending,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
	}
	require.NoError(t, db.Create(order).Error)

	t.Run("create payment with distributed lock", func(t *testing.T) {
		resp, err := svc.CreatePayment(context.Background(), customer.ID, CreatePaymentRequest{
			OrderID: order.ID,
			Method:  model.PaymentMethodWeChat,
		})
		require.NoError(t, err)
		assert.NotZero(t, resp.PaymentID)
	})
}

// mockPaymentDistributedLock is a mock implementation of cache.DistributedLock
type mockPaymentDistributedLock struct {
	tryLockResult bool
	tryLockErr    error
}

func (m *mockPaymentDistributedLock) Lock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	return m.tryLockResult, m.tryLockErr
}

func (m *mockPaymentDistributedLock) TryLock(ctx context.Context, key string, ttl time.Duration, retries int, retryDelay time.Duration) (bool, error) {
	return m.tryLockResult, m.tryLockErr
}

func (m *mockPaymentDistributedLock) Unlock(ctx context.Context, key string) error {
	return nil
}

func TestPaymentService_CreatePaymentLockFailed(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, gameModel, _ := createPaymentTestData(t, db)
	svc := createPaymentService(db)

	// Create a mock distributed lock that fails
	mockLock := &mockPaymentDistributedLock{
		tryLockResult: false,
	}
	svc.SetDistributedLock(mockLock)

	// Create order
	scheduledStart := time.Now().Add(24 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          gameModel.ID,
		Title:           "锁失败订单",
		Status:          model.OrderStatusPending,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
	}
	require.NoError(t, db.Create(order).Error)

	t.Run("create payment when lock acquisition fails", func(t *testing.T) {
		_, err := svc.CreatePayment(context.Background(), customer.ID, CreatePaymentRequest{
			OrderID: order.ID,
			Method:  model.PaymentMethodWeChat,
		})
		assert.Error(t, err)
	})
}

func TestPaymentService_CreatePaymentLockError(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, gameModel, _ := createPaymentTestData(t, db)
	svc := createPaymentService(db)

	// Create a mock distributed lock that returns error
	mockLock := &mockPaymentDistributedLock{
		tryLockResult: false,
		tryLockErr:    errors.New("lock error"),
	}
	svc.SetDistributedLock(mockLock)

	// Create order
	scheduledStart := time.Now().Add(24 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          gameModel.ID,
		Title:           "锁错误订单",
		Status:          model.OrderStatusPending,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
	}
	require.NoError(t, db.Create(order).Error)

	t.Run("create payment when lock returns error", func(t *testing.T) {
		_, err := svc.CreatePayment(context.Background(), customer.ID, CreatePaymentRequest{
			OrderID: order.ID,
			Method:  model.PaymentMethodWeChat,
		})
		assert.Error(t, err)
	})
}

func TestPaymentService_HandlePaymentCallbackPaymentNotFound(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	svc := createPaymentService(db)

	t.Run("callback with non-existent payment", func(t *testing.T) {
		err := svc.HandlePaymentCallback(context.Background(), "wechat", map[string]interface{}{
			"payment_id": float64(99999),
		})
		assert.Error(t, err)
	})
}

func TestPaymentService_ListByMethod(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, _, _ := createPaymentTestData(t, db)
	svc := createPaymentService(db)

	// Create payments with different methods
	for i := 0; i < 3; i++ {
		payment := &model.Payment{
			OrderID:     uint64(i + 600),
			UserID:      customer.ID,
			Method:      model.PaymentMethodWeChat,
			AmountCents: int64((i + 1) * 1000),
			Status:      model.PaymentStatusPending,
		}
		require.NoError(t, db.Create(payment).Error)
	}

	for i := 0; i < 2; i++ {
		payment := &model.Payment{
			OrderID:     uint64(i + 700),
			UserID:      customer.ID,
			Method:      model.PaymentMethodAlipay,
			AmountCents: int64((i + 1) * 2000),
			Status:      model.PaymentStatusPending,
		}
		require.NoError(t, db.Create(payment).Error)
	}

	t.Run("list all payments", func(t *testing.T) {
		payments, total, err := svc.List(context.Background(), repository.PaymentListOptions{
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, payments, 5)
	})
}

func TestPaymentService_RefundWithoutWalletRepo(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, gameModel, _ := createPaymentTestData(t, db)
	svc := createPaymentService(db)
	// Don't set wallet repository

	// Create order
	scheduledStart := time.Now().Add(24 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          gameModel.ID,
		Title:           "无钱包退款测试",
		Status:          model.OrderStatusConfirmed,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
	}
	require.NoError(t, db.Create(order).Error)

	// Create paid payment
	now := time.Now()
	payment := &model.Payment{
		OrderID:     order.ID,
		UserID:      customer.ID,
		Method:      model.PaymentMethodWeChat,
		AmountCents: 10000,
		Status:      model.PaymentStatusPaid,
		PaidAt:      &now,
	}
	require.NoError(t, db.Create(payment).Error)

	t.Run("refund without wallet repo", func(t *testing.T) {
		err := svc.RefundPayment(context.Background(), payment.ID, "测试退款")
		require.NoError(t, err)

		// Verify payment status
		var updated model.Payment
		require.NoError(t, db.First(&updated, payment.ID).Error)
		assert.Equal(t, model.PaymentStatusRefunded, updated.Status)
	})
}

func TestPaymentService_CreatePaymentAlipay(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, gameModel, _ := createPaymentTestData(t, db)
	svc := createPaymentService(db)

	// Create order
	scheduledStart := time.Now().Add(24 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          gameModel.ID,
		Title:           "支付宝支付订单",
		Status:          model.OrderStatusPending,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
	}
	require.NoError(t, db.Create(order).Error)

	t.Run("create alipay payment", func(t *testing.T) {
		resp, err := svc.CreatePayment(context.Background(), customer.ID, CreatePaymentRequest{
			OrderID: order.ID,
			Method:  model.PaymentMethodAlipay,
		})
		require.NoError(t, err)
		assert.NotZero(t, resp.PaymentID)
		assert.Contains(t, resp.PayInfo, "trade_no")
		assert.Contains(t, resp.PayInfo, "qr_code")
	})
}

func TestPaymentService_RefundWithAlipay(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, gameModel, _ := createPaymentTestData(t, db)
	svc := createPaymentService(db)

	// Create order
	scheduledStart := time.Now().Add(24 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          gameModel.ID,
		Title:           "支付宝退款测试",
		Status:          model.OrderStatusConfirmed,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
	}
	require.NoError(t, db.Create(order).Error)

	// Create paid alipay payment
	now := time.Now()
	payment := &model.Payment{
		OrderID:     order.ID,
		UserID:      customer.ID,
		Method:      model.PaymentMethodAlipay,
		AmountCents: 10000,
		Status:      model.PaymentStatusPaid,
		PaidAt:      &now,
	}
	require.NoError(t, db.Create(payment).Error)

	t.Run("refund alipay payment", func(t *testing.T) {
		err := svc.RefundPayment(context.Background(), payment.ID, "支付宝退款")
		require.NoError(t, err)

		// Verify payment status
		var updated model.Payment
		require.NoError(t, db.First(&updated, payment.ID).Error)
		assert.Equal(t, model.PaymentStatusRefunded, updated.Status)
	})
}

func TestPaymentService_HandlePaymentCallbackWithIntPaymentID(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	customer, _, gameModel, _ := createPaymentTestData(t, db)
	svc := createPaymentService(db)

	// Create order
	scheduledStart := time.Now().Add(24 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          gameModel.ID,
		Title:           "int测试",
		Status:          model.OrderStatusPending,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
	}
	require.NoError(t, db.Create(order).Error)

	// Create payment
	payment := &model.Payment{
		OrderID:     order.ID,
		UserID:      customer.ID,
		Method:      model.PaymentMethodWeChat,
		AmountCents: 10000,
		Status:      model.PaymentStatusPending,
	}
	require.NoError(t, db.Create(payment).Error)

	t.Run("callback with int payment_id (invalid type)", func(t *testing.T) {
		// int type is not handled, should fail
		err := svc.HandlePaymentCallback(context.Background(), "wechat", map[string]interface{}{
			"payment_id": int(payment.ID),
		})
		assert.Error(t, err)
	})
}

func TestPaymentService_WrapError(t *testing.T) {
	t.Run("wrap error function", func(t *testing.T) {
		originalErr := errors.New("original error")
		wrappedErr := WrapError(originalErr, "context")
		assert.Error(t, wrappedErr)
		assert.Contains(t, wrappedErr.Error(), "context")
	})
}

// ============================================
// Tests for provider.go
// ============================================

func TestWechatProvider_Refund(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	provider := wechatProvider{}

	t.Run("refund success", func(t *testing.T) {
		// Create payment in DB to get valid ID
		payment := &model.Payment{
			OrderID:     100,
			UserID:      1,
			Method:      model.PaymentMethodWeChat,
			AmountCents: 10000,
			Status:      model.PaymentStatusPaid,
		}
		require.NoError(t, db.Create(payment).Error)

		tradeNo, raw, refundedAt, err := provider.Refund(context.Background(), payment, "test refund")
		require.NoError(t, err)
		assert.Contains(t, tradeNo, "wx_refund_")
		assert.NotNil(t, raw)
		assert.False(t, refundedAt.IsZero())
	})
}

func TestAlipayProvider_Refund(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	provider := alipayProvider{}

	t.Run("refund success", func(t *testing.T) {
		payment := &model.Payment{
			OrderID:     101,
			UserID:      1,
			Method:      model.PaymentMethodAlipay,
			AmountCents: 20000,
			Status:      model.PaymentStatusPaid,
		}
		require.NoError(t, db.Create(payment).Error)

		tradeNo, raw, refundedAt, err := provider.Refund(context.Background(), payment, "test refund")
		require.NoError(t, err)
		assert.Contains(t, tradeNo, "ali_refund_")
		assert.NotNil(t, raw)
		assert.False(t, refundedAt.IsZero())
	})
}

func TestGenericProvider_Refund(t *testing.T) {
	db := setupPaymentTestDB(t)
	defer testutil.CleanDB(t, db)

	provider := genericProvider{}

	t.Run("refund success", func(t *testing.T) {
		payment := &model.Payment{
			OrderID:     102,
			UserID:      1,
			AmountCents: 30000,
			Status:      model.PaymentStatusPaid,
		}
		require.NoError(t, db.Create(payment).Error)

		tradeNo, raw, refundedAt, err := provider.Refund(context.Background(), payment, "test refund")
		require.NoError(t, err)
		assert.Contains(t, tradeNo, "refund_")
		assert.NotNil(t, raw)
		assert.False(t, refundedAt.IsZero())
	})
}

// ============================================
// Tests for wrap.go
// ============================================

func TestWrapError_Nil(t *testing.T) {
	t.Run("nil error returns nil", func(t *testing.T) {
		result := WrapError(nil, "test")
		assert.Nil(t, result)
	})
}

func TestWrapError_NotFound(t *testing.T) {
	t.Run("not found error", func(t *testing.T) {
		result := WrapError(repository.ErrNotFound, "get payment")
		assert.Error(t, result)
	})
}

func TestWrapError_Validation(t *testing.T) {
	t.Run("validation error", func(t *testing.T) {
		result := WrapError(ErrValidation, "create payment")
		assert.Error(t, result)
	})
}

func TestWrapError_Generic(t *testing.T) {
	t.Run("generic error", func(t *testing.T) {
		result := WrapError(errors.New("some error"), "operation")
		assert.Error(t, result)
	})
}

// ============================================
// Tests for RefundService
// ============================================

func setupRefundTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.User{},
		&model.Player{},
		&model.Game{},
		&model.Order{},
		&model.Payment{},
		&model.Wallet{},
		&model.RefundRecord{},
		&model.OperationLog{},
	)
	return db
}

func createRefundService(db *gorm.DB) *RefundService {
	paymentRepo := orderrepo.NewPaymentRepository(db)
	refundRepo := paymentrepo.NewRefundRecordRepository(db)
	orderRepo := orderimpl.NewOrderRepository(db)
	return NewRefundService(paymentRepo, refundRepo, orderRepo)
}

func TestRefundService_ProcessRefund(t *testing.T) {
	db := setupRefundTestDB(t)
	defer testutil.CleanDB(t, db)

	// Create test user
	customer := &model.User{
		Phone:        "13800000001",
		Email:        "customer@test.com",
		Name:         "Customer",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(customer).Error)

	// Create test order
	scheduledStart := time.Now().Add(24 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          1,
		Title:           "测试订单",
		Status:          model.OrderStatusConfirmed,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
	}
	require.NoError(t, db.Create(order).Error)

	// Create paid payment
	now := time.Now()
	payment := &model.Payment{
		OrderID:     order.ID,
		UserID:      customer.ID,
		Method:      model.PaymentMethodWeChat,
		AmountCents: 10000,
		Status:      model.PaymentStatusPaid,
		PaidAt:      &now,
	}
	require.NoError(t, db.Create(payment).Error)

	svc := createRefundService(db)

	t.Run("process full refund", func(t *testing.T) {
		resp, err := svc.ProcessRefund(context.Background(), model.RefundRequest{
			PaymentID:   payment.ID,
			AmountCents: 10000,
			Reason:      "用户申请退款",
		})
		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, model.RefundStatusProcessed, resp.RefundRecord.Status)
		assert.Equal(t, int64(0), resp.RemainingAmount)
	})
}

func TestRefundService_ProcessRefund_PaymentNotFound(t *testing.T) {
	db := setupRefundTestDB(t)
	defer testutil.CleanDB(t, db)

	svc := createRefundService(db)

	t.Run("payment not found", func(t *testing.T) {
		_, err := svc.ProcessRefund(context.Background(), model.RefundRequest{
			PaymentID:   99999,
			AmountCents: 10000,
			Reason:      "测试",
		})
		assert.Error(t, err)
	})
}

func TestRefundService_GetRefundHistory(t *testing.T) {
	db := setupRefundTestDB(t)
	defer testutil.CleanDB(t, db)

	// Create test user
	customer := &model.User{
		Phone:        "13800000002",
		Email:        "customer2@test.com",
		Name:         "Customer2",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(customer).Error)

	// Create paid payment
	now := time.Now()
	payment := &model.Payment{
		OrderID:     1,
		UserID:      customer.ID,
		Method:      model.PaymentMethodWeChat,
		AmountCents: 10000,
		Status:      model.PaymentStatusPaid,
		PaidAt:      &now,
	}
	require.NoError(t, db.Create(payment).Error)

	// Create refund record
	refund := &model.RefundRecord{
		PaymentID:   payment.ID,
		OrderID:     1,
		UserID:      customer.ID,
		AmountCents: 5000,
		Reason:      "部分退款",
		Status:      model.RefundStatusProcessed,
	}
	require.NoError(t, db.Create(refund).Error)

	svc := createRefundService(db)

	t.Run("get refund history", func(t *testing.T) {
		records, err := svc.GetRefundHistory(context.Background(), payment.ID)
		require.NoError(t, err)
		assert.Len(t, records, 1)
		assert.Equal(t, int64(5000), records[0].AmountCents)
	})

	t.Run("payment not found", func(t *testing.T) {
		_, err := svc.GetRefundHistory(context.Background(), 99999)
		assert.Error(t, err)
	})
}

func TestRefundService_GetRefundsByOrder(t *testing.T) {
	db := setupRefundTestDB(t)
	defer testutil.CleanDB(t, db)

	// Create test user
	customer := &model.User{
		Phone:        "13800000003",
		Email:        "customer3@test.com",
		Name:         "Customer3",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(customer).Error)

	// Create refund records for an order
	for i := 0; i < 3; i++ {
		refund := &model.RefundRecord{
			PaymentID:   uint64(i + 1),
			OrderID:     100,
			UserID:      customer.ID,
			AmountCents: int64((i + 1) * 1000),
			Reason:      "退款",
			Status:      model.RefundStatusProcessed,
		}
		require.NoError(t, db.Create(refund).Error)
	}

	svc := createRefundService(db)

	t.Run("get refunds by order", func(t *testing.T) {
		records, err := svc.GetRefundsByOrder(context.Background(), 100)
		require.NoError(t, err)
		assert.Len(t, records, 3)
	})
}

func TestRefundService_SetWalletRepository(t *testing.T) {
	db := setupRefundTestDB(t)
	defer testutil.CleanDB(t, db)

	svc := createRefundService(db)
	walletRepo := user.NewWalletRepository(db)

	t.Run("set wallet repository", func(t *testing.T) {
		svc.SetWalletRepository(walletRepo)
		assert.NotNil(t, svc.wallets)
	})
}

func TestRefundService_SetOperationLogRepository(t *testing.T) {
	db := setupRefundTestDB(t)
	defer testutil.CleanDB(t, db)

	svc := createRefundService(db)

	t.Run("set operation log repository to nil", func(t *testing.T) {
		svc.SetOperationLogRepository(nil)
		assert.Nil(t, svc.opLogs)
	})
}

func TestRefundService_ProcessRefundWithWallet(t *testing.T) {
	db := setupRefundTestDB(t)
	defer testutil.CleanDB(t, db)

	// Create test user
	customer := &model.User{
		Phone:        "13800000004",
		Email:        "customer4@test.com",
		Name:         "Customer4",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(customer).Error)

	// Create wallet
	wallet := &model.Wallet{
		UserID:       customer.ID,
		BalanceCents: 0,
	}
	require.NoError(t, db.Create(wallet).Error)

	// Create test order
	scheduledStart := time.Now().Add(24 * time.Hour)
	order := &model.Order{
		UserID:          customer.ID,
		ItemID:          1,
		Title:           "测试订单",
		Status:          model.OrderStatusConfirmed,
		UnitPriceCents:  5000,
		TotalPriceCents: 10000,
		ScheduledStart:  &scheduledStart,
	}
	require.NoError(t, db.Create(order).Error)

	// Create paid payment
	now := time.Now()
	payment := &model.Payment{
		OrderID:     order.ID,
		UserID:      customer.ID,
		Method:      model.PaymentMethodAlipay,
		AmountCents: 10000,
		Status:      model.PaymentStatusPaid,
		PaidAt:      &now,
	}
	require.NoError(t, db.Create(payment).Error)

	svc := createRefundService(db)
	walletRepo := user.NewWalletRepository(db)
	svc.SetWalletRepository(walletRepo)

	t.Run("refund credits wallet", func(t *testing.T) {
		resp, err := svc.ProcessRefund(context.Background(), model.RefundRequest{
			PaymentID:   payment.ID,
			AmountCents: 10000,
			Reason:      "退款测试",
		})
		require.NoError(t, err)
		assert.NotNil(t, resp)

		// Verify wallet balance
		updatedWallet, err := walletRepo.GetByUserID(context.Background(), customer.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(10000), updatedWallet.BalanceCents)
	})
}
