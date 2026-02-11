//go:build integration

package payment

import (
	"context"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository/common"
	"gamelink/internal/repository/implementations"
)

// TestPaymentService_HandleWeChatPaymentCallback 测试微信支付回调处理
func TestPaymentService_HandleWeChatPaymentCallback(t *testing.T) {
	// 创建测试数据库
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	// 创建仓储
	paymentRepo := implementations.NewPaymentRepository(db)
	orderRepo := implementations.NewOrderRepository(db)

	// 创建支付服务
	paymentService := NewPaymentService(paymentRepo, orderRepo)
	paymentService.SetTxManager(common.NewUnitOfWork(db))

	// 创建微信回调处理器
	wechatHandler := NewWeChatCallbackHandler(
		"test_app_id",
		"test_mch_id",
		"test_api_key",
	)
	paymentService.SetWeChatCallbackHandler(wechatHandler)

	// 创建测试订单
	ctx := context.Background()
	order := &model.Order{
		UserID:          1,
		PlayerID:        2,
		ServiceID:       1,
		Status:          model.OrderStatusPending,
		TotalPriceCents: 10000, // 100元
		Currency:        "CNY",
		ScheduledAt:     time.Now(),
	}
	if err := orderRepo.Create(ctx, order); err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	// 创建待支付的支付记录
	payment := &model.Payment{
		RequestID:             "test_request_123",
		OrderID:               order.ID,
		UserID:                1,
		Method:                model.PaymentMethodWeChat,
		AmountCents:           10000,
		ThirdPartyAmountCents: 10000,
		Currency:              "CNY",
		Status:                model.PaymentStatusPending,
	}
	if err := paymentRepo.Create(ctx, payment); err != nil {
		t.Fatalf("failed to create payment: %v", err)
	}

	// 构造微信支付回调 XML
	callbackXML := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<xml>
  <return_code><![CDATA[SUCCESS]]></return_code>
  <return_msg><![CDATA[OK]]></return_msg>
  <appid><![CDATA[test_app_id]]></appid>
  <mch_id><![CDATA[test_mch_id]]></mch_id>
  <nonce_str><![CDATA[test_nonce_str]]></nonce_str>
  <sign><![CDATA[TESTSIGN123456789012345678901234]]></sign>
  <result_code><![CDATA[SUCCESS]]></result_code>
  <openid><![CDATA[test_openid]]></openid>
  <trade_type><![CDATA[JSAPI]]></trade_type>
  <bank_type><![CDATA[CFT]]></bank_type>
  <total_fee><![CDATA[10000]]></total_fee>
  <fee_type><![CDATA[CNY]]></fee_type>
  <transaction_id><![CDATA[WX2026020900123456789012345678]]></transaction_id>
  <out_trade_no><![CDATA[` + string(rune(order.ID)) + `]]></out_trade_no>
  <time_end><![CDATA[20260209120000]]></time_end>
</xml>`)

	// 注意：上面的 out_trade_no 应该是订单 ID，但在实际场景中应该是商户订单号
	// 为了测试，我们需要修改这部分，使用正确的订单号

	t.Logf("Created order with ID: %d", order.ID)
	t.Logf("Created payment with ID: %d", payment.ID)

	// TODO: 完善这个测试
	// 由于需要生成正确的签名和订单号映射，这个测试需要进一步实现

	_ = callbackXML // 暂时避免未使用变量警告
}

// TestPaymentService_HandleWeChatPaymentCallback_Duplicate 测试重复回调处理
func TestPaymentService_HandleWeChatPaymentCallback_Duplicate(t *testing.T) {
	// TODO: 实现重复回调测试
	// 验证分布式锁机制
}

// TestPaymentService_HandleWeChatPaymentCallback_InvalidSignature 测试无效签名
func TestPaymentService_HandleWeChatPaymentCallback_InvalidSignature(t *testing.T) {
	// TODO: 实现无效签名测试
	// 验证签名验证失败时返回错误响应
}

// TestPaymentService_HandleWeChatPaymentCallback_AmountMismatch 测试金额不匹配
func TestPaymentService_HandleWeChatPaymentCallback_AmountMismatch(t *testing.T) {
	// TODO: 实现金额不匹配测试
	// 验证回调金额与订单金额不一致时返回错误
}
