package payment

import (
	"context"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"

	"github.com/stretchr/testify/assert"
)

// TestCreatePayment_EdgeCases 测试创建支付的边界情况
func TestCreatePayment_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("创建支付成功", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		// 创建待支付订单
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 10000,
			Currency:        model.CurrencyCNY,
		}
		orderRepo.Create(ctx, order)

		req := CreatePaymentRequest{
			OrderID: order.ID,
			Method:  model.PaymentMethodWeChat,
		}

		resp, err := svc.CreatePayment(ctx, 1, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotZero(t, resp.PaymentID)
		assert.NotEmpty(t, resp.PayInfo)
	})

	t.Run("订单不存在应该失败", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		req := CreatePaymentRequest{
			OrderID: 999,
			Method:  model.PaymentMethodWeChat,
		}

		resp, err := svc.CreatePayment(ctx, 1, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("无权限支付他人订单", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		// 创建用户1的订单
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 10000,
			Currency:        model.CurrencyCNY,
		}
		orderRepo.Create(ctx, order)

		req := CreatePaymentRequest{
			OrderID: order.ID,
			Method:  model.PaymentMethodWeChat,
		}

		// 用户2尝试支付
		resp, err := svc.CreatePayment(ctx, 2, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "unauthorized")
	})

	t.Run("非pending状态订单不能支付", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		// 创建已确认的订单
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusConfirmed,
			TotalPriceCents: 10000,
			Currency:        model.CurrencyCNY,
		}
		orderRepo.Create(ctx, order)

		req := CreatePaymentRequest{
			OrderID: order.ID,
			Method:  model.PaymentMethodWeChat,
		}

		resp, err := svc.CreatePayment(ctx, 1, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, ErrInvalidOrderStatus, err)
	})

	t.Run("订单已支付不能重复支付", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		// 创建订单
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 10000,
			Currency:        model.CurrencyCNY,
		}
		orderRepo.Create(ctx, order)

		// 创建已支付的支付记录
		paidAt := time.Now()
		payment := &model.Payment{
			OrderID:     order.ID,
			UserID:      1,
			Method:      model.PaymentMethodWeChat,
			AmountCents: 10000,
			Status:      model.PaymentStatusPaid,
			PaidAt:      &paidAt,
		}
		paymentRepo.Create(ctx, payment)

		req := CreatePaymentRequest{
			OrderID: order.ID,
			Method:  model.PaymentMethodWeChat,
		}

		resp, err := svc.CreatePayment(ctx, 1, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, ErrOrderAlreadyPaid, err)
	})

	t.Run("支付金额为0", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		// 创建0元订单
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 0,
			Currency:        model.CurrencyCNY,
		}
		orderRepo.Create(ctx, order)

		req := CreatePaymentRequest{
			OrderID: order.ID,
			Method:  model.PaymentMethodWeChat,
		}

		resp, err := svc.CreatePayment(ctx, 1, req)

		// 当前实现允许0元支付，但实际业务可能需要验证
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("支付金额为极大值", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		// 创建大额订单
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 10000000, // 100,000元
			Currency:        model.CurrencyCNY,
		}
		orderRepo.Create(ctx, order)

		req := CreatePaymentRequest{
			OrderID: order.ID,
			Method:  model.PaymentMethodWeChat,
		}

		resp, err := svc.CreatePayment(ctx, 1, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})
}

// TestGetPaymentStatus_EdgeCases 测试查询支付状态的边界情况
func TestGetPaymentStatus_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("查询支付状态成功", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		paidAt := time.Now()
		payment := &model.Payment{
			OrderID:     1,
			UserID:      1,
			Method:      model.PaymentMethodWeChat,
			AmountCents: 10000,
			Status:      model.PaymentStatusPaid,
			PaidAt:      &paidAt,
		}
		paymentRepo.Create(ctx, payment)

		resp, err := svc.GetPaymentStatus(ctx, payment.ID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, payment.ID, resp.PaymentID)
		assert.Equal(t, model.PaymentStatusPaid, resp.Status)
		assert.NotNil(t, resp.PaidAt)
	})

	t.Run("查询不存在的支付", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		resp, err := svc.GetPaymentStatus(ctx, 999)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("查询pending状态的支付", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		payment := &model.Payment{
			OrderID:     1,
			UserID:      1,
			Method:      model.PaymentMethodAlipay,
			AmountCents: 10000,
			Status:      model.PaymentStatusPending,
		}
		paymentRepo.Create(ctx, payment)

		resp, err := svc.GetPaymentStatus(ctx, payment.ID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, model.PaymentStatusPending, resp.Status)
		assert.Nil(t, resp.PaidAt)
	})
}

// TestCancelPayment_EdgeCases 测试取消支付的边界情况
func TestCancelPayment_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("取消pending状态的支付", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		payment := &model.Payment{
			OrderID:     1,
			UserID:      1,
			Method:      model.PaymentMethodWeChat,
			AmountCents: 10000,
			Status:      model.PaymentStatusPending,
		}
		paymentRepo.Create(ctx, payment)

		err := svc.CancelPayment(ctx, 1, payment.ID)

		assert.NoError(t, err)

		// 验证状态已更新
		updated, _ := paymentRepo.Get(ctx, payment.ID)
		assert.Equal(t, model.PaymentStatusFailed, updated.Status)
	})

	t.Run("无权限取消他人支付", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		payment := &model.Payment{
			OrderID:     1,
			UserID:      1,
			Method:      model.PaymentMethodWeChat,
			AmountCents: 10000,
			Status:      model.PaymentStatusPending,
		}
		paymentRepo.Create(ctx, payment)

		// 用户2尝试取消用户1的支付
		err := svc.CancelPayment(ctx, 2, payment.ID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	})

	t.Run("不能取消已支付的支付", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		paidAt := time.Now()
		payment := &model.Payment{
			OrderID:     1,
			UserID:      1,
			Method:      model.PaymentMethodWeChat,
			AmountCents: 10000,
			Status:      model.PaymentStatusPaid,
			PaidAt:      &paidAt,
		}
		paymentRepo.Create(ctx, payment)

		err := svc.CancelPayment(ctx, 1, payment.ID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot cancel payment")
	})

	t.Run("取消不存在的支付", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		err := svc.CancelPayment(ctx, 1, 999)

		assert.Error(t, err)
		assert.Equal(t, repository.ErrNotFound, err)
	})
}

// TestHandlePaymentCallback_EdgeCases 测试支付回调的边界情况
func TestHandlePaymentCallback_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("成功处理支付回调", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		// 创建订单和支付
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 10000,
		}
		orderRepo.Create(ctx, order)

		payment := &model.Payment{
			OrderID:     order.ID,
			UserID:      1,
			Method:      model.PaymentMethodWeChat,
			AmountCents: 10000,
			Status:      model.PaymentStatusPending,
		}
		paymentRepo.Create(ctx, payment)

		// 模拟支付回调
		callbackData := map[string]interface{}{
			"payment_id": float64(payment.ID),
			"status":     "success",
			"trade_no":   "wx_trade_123",
		}

		err := svc.HandlePaymentCallback(ctx, "wechat", callbackData)

		assert.NoError(t, err)

		// 验证支付状态已更新
		updatedPayment, _ := paymentRepo.Get(ctx, payment.ID)
		assert.Equal(t, model.PaymentStatusPaid, updatedPayment.Status)
		assert.NotNil(t, updatedPayment.PaidAt)
	})

	t.Run("重复回调应该幂等", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		// 创建已支付的支付
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusConfirmed,
			TotalPriceCents: 10000,
		}
		orderRepo.Create(ctx, order)

		paidAt := time.Now()
		payment := &model.Payment{
			OrderID:     order.ID,
			UserID:      1,
			Method:      model.PaymentMethodWeChat,
			AmountCents: 10000,
			Status:      model.PaymentStatusPaid,
			PaidAt:      &paidAt,
		}
		paymentRepo.Create(ctx, payment)

		// 重复回调
		callbackData := map[string]interface{}{
			"payment_id": float64(payment.ID),
			"status":     "success",
		}

		err := svc.HandlePaymentCallback(ctx, "wechat", callbackData)

		// 应该成功返回，不报错
		assert.NoError(t, err)
	})

	t.Run("缺少payment_id应该失败", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		callbackData := map[string]interface{}{
			"status": "success",
		}

		err := svc.HandlePaymentCallback(ctx, "wechat", callbackData)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing payment_id")
	})

	t.Run("支付不存在应该失败", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		callbackData := map[string]interface{}{
			"payment_id": float64(999),
			"status":     "success",
		}

		err := svc.HandlePaymentCallback(ctx, "wechat", callbackData)

		assert.Error(t, err)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("支付方式不匹配应该失败", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		payment := &model.Payment{
			OrderID:     1,
			UserID:      1,
			Method:      model.PaymentMethodWeChat,
			AmountCents: 10000,
			Status:      model.PaymentStatusPending,
		}
		paymentRepo.Create(ctx, payment)

		callbackData := map[string]interface{}{
			"payment_id": float64(payment.ID),
			"status":     "success",
		}

		// 使用alipay回调wechat支付
		err := svc.HandlePaymentCallback(ctx, "alipay", callbackData)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "provider mismatch")
	})
}

// TestPaymentMethods 测试不同支付方式
func TestPaymentMethods(t *testing.T) {
	ctx := context.Background()

	t.Run("微信支付", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 10000,
			Currency:        model.CurrencyCNY,
		}
		orderRepo.Create(ctx, order)

		req := CreatePaymentRequest{
			OrderID: order.ID,
			Method:  model.PaymentMethodWeChat,
		}

		resp, err := svc.CreatePayment(ctx, 1, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Contains(t, resp.PayInfo, "prepay_id")
		assert.Contains(t, resp.PayInfo, "code_url")
	})

	t.Run("支付宝支付", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 10000,
			Currency:        model.CurrencyCNY,
		}
		orderRepo.Create(ctx, order)

		req := CreatePaymentRequest{
			OrderID: order.ID,
			Method:  model.PaymentMethodAlipay,
		}

		resp, err := svc.CreatePayment(ctx, 1, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Contains(t, resp.PayInfo, "trade_no")
		assert.Contains(t, resp.PayInfo, "qr_code")
	})
}
