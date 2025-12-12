package model_test

import (
	"encoding/json"
	"testing"
	"time"

	"gamelink/internal/model"

	"github.com/stretchr/testify/assert"
)

func TestPaymentModel(t *testing.T) {
	now := time.Now()
	rawData := json.RawMessage(`{"key": "value"}`)

	payment := &model.Payment{
		Base: model.Base{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		OrderID:         100,
		UserID:          200,
		Method:          model.PaymentMethodWeChat,
		AmountCents:     10000,
		Currency:        model.CurrencyCNY,
		Status:          model.PaymentStatusPaid,
		ProviderTradeNo: "TRADE123456",
		ProviderRaw:     rawData,
		PaidAt:          &now,
		RefundedAt:      nil,
	}

	assert.Equal(t, uint64(1), payment.ID)
	assert.Equal(t, uint64(100), payment.OrderID)
	assert.Equal(t, uint64(200), payment.UserID)
	assert.Equal(t, model.PaymentMethodWeChat, payment.Method)
	assert.Equal(t, int64(10000), payment.AmountCents)
	assert.Equal(t, model.CurrencyCNY, payment.Currency)
	assert.Equal(t, model.PaymentStatusPaid, payment.Status)
	assert.Equal(t, "TRADE123456", payment.ProviderTradeNo)
	assert.Equal(t, rawData, payment.ProviderRaw)
	assert.Equal(t, &now, payment.PaidAt)
	assert.Nil(t, payment.RefundedAt)
}

func TestPaymentJSONSerialization(t *testing.T) {
	now := time.Now()
	rawData := json.RawMessage(`{"provider": "wechat", "amount": 10000}`)

	payment := &model.Payment{
		Base: model.Base{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		OrderID:         100,
		UserID:          200,
		Method:          model.PaymentMethodAlipay,
		AmountCents:     5000,
		Currency:        model.CurrencyCNY,
		Status:          model.PaymentStatusPending,
		ProviderTradeNo: "TRADE789",
		ProviderRaw:     rawData,
		PaidAt:          nil,
		RefundedAt:      nil,
	}

	// 序列化
	data, err := json.Marshal(payment)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "TRADE789")
	assert.Contains(t, string(data), "alipay")

	// 反序列化
	var decoded model.Payment
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, payment.OrderID, decoded.OrderID)
	assert.Equal(t, payment.UserID, decoded.UserID)
	assert.Equal(t, payment.Method, decoded.Method)
	assert.Equal(t, payment.AmountCents, decoded.AmountCents)
	assert.Equal(t, payment.Status, decoded.Status)
}

func TestPaymentConstants(t *testing.T) {
	// 测试支付方式常量
	assert.Equal(t, model.PaymentMethod("wechat"), model.PaymentMethodWeChat)
	assert.Equal(t, model.PaymentMethod("alipay"), model.PaymentMethodAlipay)

	// 测试支付状态常量
	assert.Equal(t, model.PaymentStatus("pending"), model.PaymentStatusPending)
	assert.Equal(t, model.PaymentStatus("paid"), model.PaymentStatusPaid)
	assert.Equal(t, model.PaymentStatus("failed"), model.PaymentStatusFailed)
	assert.Equal(t, model.PaymentStatus("refunded"), model.PaymentStatusRefunded)
}

func TestPaymentWithPaidAt(t *testing.T) {
	now := time.Now()
	payment := &model.Payment{
		PaidAt: &now,
	}

	assert.NotNil(t, payment.PaidAt)
	assert.Equal(t, now.Unix(), payment.PaidAt.Unix())
}

func TestPaymentWithRefundedAt(t *testing.T) {
	now := time.Now()
	payment := &model.Payment{
		RefundedAt: &now,
	}

	assert.NotNil(t, payment.RefundedAt)
	assert.Equal(t, now.Unix(), payment.RefundedAt.Unix())
}

func TestPaymentNilTimes(t *testing.T) {
	payment := &model.Payment{
		PaidAt:     nil,
		RefundedAt: nil,
	}

	assert.Nil(t, payment.PaidAt)
	assert.Nil(t, payment.RefundedAt)
}

func TestPaymentProviderRawData(t *testing.T) {
	// 测试空的RawMessage
	payment1 := &model.Payment{
		ProviderRaw: json.RawMessage(nil),
	}
	assert.Nil(t, payment1.ProviderRaw)

	// 测试有效的JSON数据
	jsonData := `{"status": "success", "code": 200}`
	payment2 := &model.Payment{
		ProviderRaw: json.RawMessage(jsonData),
	}
	assert.NotNil(t, payment2.ProviderRaw)
	assert.Equal(t, jsonData, string(payment2.ProviderRaw))

	// 测试复杂JSON数据
	complexData := `{"provider": "wechat", "transaction": {"id": "123", "amount": 10000, "currency": "CNY"}}`
	payment3 := &model.Payment{
		ProviderRaw: json.RawMessage(complexData),
	}
	assert.NotNil(t, payment3.ProviderRaw)
	assert.Equal(t, complexData, string(payment3.ProviderRaw))
}

func TestPaymentEdgeCases(t *testing.T) {
	// 测试零值
	payment := &model.Payment{
		OrderID:         0,
		UserID:          0,
		AmountCents:     0,
		ProviderTradeNo: "",
	}

	assert.Equal(t, uint64(0), payment.OrderID)
	assert.Equal(t, uint64(0), payment.UserID)
	assert.Equal(t, int64(0), payment.AmountCents)
	assert.Equal(t, "", payment.ProviderTradeNo)

	// 测试最大整数值
	payment2 := &model.Payment{
		OrderID:     ^uint64(0),
		UserID:      ^uint64(0),
		AmountCents: ^int64(0),
	}

	assert.Equal(t, ^uint64(0), payment2.OrderID)
	assert.Equal(t, ^uint64(0), payment2.UserID)
	assert.Equal(t, ^int64(0), payment2.AmountCents)
}

func TestPaymentRelations(t *testing.T) {
	payment := &model.Payment{
		Order: model.Order{
			OrderNo: "ORDER123",
		},
		User: model.User{
			Name: "test user",
		},
	}

	assert.Equal(t, "ORDER123", payment.Order.OrderNo)
	assert.Equal(t, "test user", payment.User.Name)
}

func TestPaymentMethodValues(t *testing.T) {
	// 测试所有支付方式
	payment1 := &model.Payment{
		Method: model.PaymentMethodWeChat,
	}
	assert.Equal(t, model.PaymentMethodWeChat, payment1.Method)

	payment2 := &model.Payment{
		Method: model.PaymentMethodAlipay,
	}
	assert.Equal(t, model.PaymentMethodAlipay, payment2.Method)
}

func TestPaymentStatusValues(t *testing.T) {
	// 测试所有支付状态
	statuses := []model.PaymentStatus{
		model.PaymentStatusPending,
		model.PaymentStatusPaid,
		model.PaymentStatusFailed,
		model.PaymentStatusRefunded,
	}

	for _, status := range statuses {
		payment := &model.Payment{Status: status}
		assert.Equal(t, status, payment.Status)
	}
}

func TestPaymentCurrency(t *testing.T) {
	// 测试货币字段
	payment := &model.Payment{
		Currency: model.CurrencyUSD,
	}
	assert.Equal(t, model.CurrencyUSD, payment.Currency)

	payment2 := &model.Payment{
		Currency: model.CurrencyCNY,
	}
	assert.Equal(t, model.CurrencyCNY, payment2.Currency)
}

func TestPaymentTimeSequence(t *testing.T) {
	// 测试时间顺序
	createdAt := time.Now()
	paidAt := createdAt.Add(1 * time.Hour)
	refundedAt := paidAt.Add(30 * time.Minute)

	payment := &model.Payment{
		Base: model.Base{
			CreatedAt: createdAt,
		},
		PaidAt:     &paidAt,
		RefundedAt: &refundedAt,
	}

	assert.True(t, payment.PaidAt.After(payment.CreatedAt))
	assert.True(t, payment.RefundedAt.After(*payment.PaidAt))
}

func TestPaymentJSONWithRawMessage(t *testing.T) {
	// 测试包含RawMessage的JSON序列化
	payment := &model.Payment{
		ProviderRaw: json.RawMessage(`{"test":"data"}`), // 移除空格以匹配JSON标准格式
	}

	data, err := json.Marshal(payment)
	assert.NoError(t, err)

	// 反序列化回来
	var decoded model.Payment
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	// 由于JSON序列化会标准化格式，我们验证内容而不是精确的字节匹配
	assert.Equal(t, string(payment.ProviderRaw), string(decoded.ProviderRaw))
}

func TestPaymentNewFields(t *testing.T) {
	// 测试新增的收款主体ID和商户号字段
	entityID := uint64(123)
	payment := &model.Payment{
		CollectionEntityID:  &entityID,
		MerchantNo:          "MCH123456",
		RefundedAmountCents: 5000,
	}

	assert.NotNil(t, payment.CollectionEntityID)
	assert.Equal(t, uint64(123), *payment.CollectionEntityID)
	assert.Equal(t, "MCH123456", payment.MerchantNo)
	assert.Equal(t, int64(5000), payment.RefundedAmountCents)
}

func TestPaymentHasRequiredFields(t *testing.T) {
	now := time.Now()

	// 测试完整的支付记录
	completePayment := &model.Payment{
		Base: model.Base{
			ID:        1,
			CreatedAt: now,
		},
		OrderID: 100,
		UserID:  200,
		Method:  model.PaymentMethodWeChat,
		Status:  model.PaymentStatusPaid,
	}
	assert.True(t, completePayment.HasRequiredFields())

	// 测试缺少OrderID
	missingOrderID := &model.Payment{
		Base: model.Base{
			ID:        1,
			CreatedAt: now,
		},
		UserID: 200,
		Method: model.PaymentMethodWeChat,
		Status: model.PaymentStatusPaid,
	}
	assert.False(t, missingOrderID.HasRequiredFields())

	// 测试缺少UserID
	missingUserID := &model.Payment{
		Base: model.Base{
			ID:        1,
			CreatedAt: now,
		},
		OrderID: 100,
		Method:  model.PaymentMethodWeChat,
		Status:  model.PaymentStatusPaid,
	}
	assert.False(t, missingUserID.HasRequiredFields())

	// 测试缺少Method
	missingMethod := &model.Payment{
		Base: model.Base{
			ID:        1,
			CreatedAt: now,
		},
		OrderID: 100,
		UserID:  200,
		Status:  model.PaymentStatusPaid,
	}
	assert.False(t, missingMethod.HasRequiredFields())

	// 测试缺少Status
	missingStatus := &model.Payment{
		Base: model.Base{
			ID:        1,
			CreatedAt: now,
		},
		OrderID: 100,
		UserID:  200,
		Method:  model.PaymentMethodWeChat,
	}
	assert.False(t, missingStatus.HasRequiredFields())

	// 测试缺少CreatedAt
	missingCreatedAt := &model.Payment{
		Base: model.Base{
			ID: 1,
		},
		OrderID: 100,
		UserID:  200,
		Method:  model.PaymentMethodWeChat,
		Status:  model.PaymentStatusPaid,
	}
	assert.False(t, missingCreatedAt.HasRequiredFields())
}

func TestPaymentRemainingRefundableAmount(t *testing.T) {
	// 测试已支付状态的可退款金额
	paidPayment := &model.Payment{
		AmountCents:         10000,
		RefundedAmountCents: 3000,
		Status:              model.PaymentStatusPaid,
	}
	assert.Equal(t, int64(7000), paidPayment.RemainingRefundableAmount())

	// 测试未支付状态的可退款金额（应为0）
	pendingPayment := &model.Payment{
		AmountCents:         10000,
		RefundedAmountCents: 0,
		Status:              model.PaymentStatusPending,
	}
	assert.Equal(t, int64(0), pendingPayment.RemainingRefundableAmount())

	// 测试已退款状态的可退款金额（应为0）
	refundedPayment := &model.Payment{
		AmountCents:         10000,
		RefundedAmountCents: 10000,
		Status:              model.PaymentStatusRefunded,
	}
	assert.Equal(t, int64(0), refundedPayment.RemainingRefundableAmount())

	// 测试全额可退款
	fullRefundable := &model.Payment{
		AmountCents:         10000,
		RefundedAmountCents: 0,
		Status:              model.PaymentStatusPaid,
	}
	assert.Equal(t, int64(10000), fullRefundable.RemainingRefundableAmount())
}

func TestPaymentCollectionEntityIDNil(t *testing.T) {
	// 测试收款主体ID为nil的情况
	payment := &model.Payment{
		CollectionEntityID: nil,
		MerchantNo:         "",
	}

	assert.Nil(t, payment.CollectionEntityID)
	assert.Equal(t, "", payment.MerchantNo)
}
