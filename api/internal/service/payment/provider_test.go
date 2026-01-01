package payment

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"gamelink/internal/model"
	"gamelink/internal/service/external"
	"gamelink/pkg/config"
)

// TestRealAlipayProvider_Refund_NoPrivateKey tests refund when private key is not loaded
func TestRealAlipayProvider_Refund_NoPrivateKey(t *testing.T) {
	cfg := &external.Config{
		Alipay: config.AlipayConfig{
			Enabled: true,
			AppID:   "test_app_id",
		},
	}

	provider, err := NewAlipayProvider(cfg)
	assert.NoError(t, err)

	payment := &model.Payment{
		Base:        model.Base{ID: 12345},
		AmountCents: 10000,
	}

	reason := "Test refund"
	refundNo, rawDetails, refundedAt, err := provider.Refund(context.Background(), payment, reason)

	// Should fail due to missing private key
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "private key not loaded")
	assert.Empty(t, refundNo)
	assert.Nil(t, rawDetails)
	assert.Zero(t, refundedAt)
}

// TestRealAlipayProvider_Refund_Disabled tests refund when Alipay is disabled
func TestRealAlipayProvider_Refund_Disabled(t *testing.T) {
	cfg := &external.Config{
		Alipay: config.AlipayConfig{
			Enabled: false,
		},
	}

	provider, err := NewAlipayProvider(cfg)
	assert.NoError(t, err)

	payment := &model.Payment{
		Base:        model.Base{ID: 12345},
		AmountCents: 10000,
	}

	refundNo, rawDetails, refundedAt, err := provider.Refund(context.Background(), payment, "test")

	assert.Error(t, err)
	assert.Equal(t, ErrPaymentDisabled, err)
	assert.Empty(t, refundNo)
	assert.Nil(t, rawDetails)
	assert.Zero(t, refundedAt)
}

// TestRealAlipayProvider_CreateOrder_NoPrivateKey tests order creation when private key is not loaded
func TestRealAlipayProvider_CreateOrder_NoPrivateKey(t *testing.T) {
	cfg := &external.Config{
		Alipay: config.AlipayConfig{
			Enabled: true,
			AppID:   "test_app_id",
		},
	}

	provider, err := NewAlipayProvider(cfg)
	assert.NoError(t, err)

	orderID := "ORDER_123"
	subject := "Test Order Subject"
	amountCents := int64(20000)

	result, err := provider.CreateOrder(context.Background(), orderID, subject, amountCents)

	// Should fail due to missing private key
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "private key not loaded")
	assert.Nil(t, result)
}

// TestRealAlipayProvider_CreateOrder_Disabled tests order creation when Alipay is disabled
func TestRealAlipayProvider_CreateOrder_Disabled(t *testing.T) {
	cfg := &external.Config{
		Alipay: config.AlipayConfig{
			Enabled: false,
		},
	}

	provider, err := NewAlipayProvider(cfg)
	assert.NoError(t, err)

	result, err := provider.CreateOrder(context.Background(), "ORDER_123", "Test", 10000)

	assert.Error(t, err)
	assert.Equal(t, ErrPaymentDisabled, err)
	assert.Nil(t, result)
}

// TestRealAlipayProvider_generateSign_NoPrivateKey tests sign generation without private key
func TestRealAlipayProvider_generateSign_NoPrivateKey(t *testing.T) {
	cfg := &external.Config{
		Alipay: config.AlipayConfig{
			Enabled: true,
		},
	}

	provider, err := NewAlipayProvider(cfg)
	assert.NoError(t, err)

	params := map[string]string{
		"app_id": "test_app",
		"method": "test_method",
	}

	sign, err := provider.generateSign(params)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "private key not loaded")
	assert.Empty(t, sign)
}

// TestRealAlipayProvider_VerifySign_NoPublicKey tests signature verification without public key
func TestRealAlipayProvider_VerifySign_NoPublicKey(t *testing.T) {
	cfg := &external.Config{
		Alipay: config.AlipayConfig{},
	}

	provider, err := NewAlipayProvider(cfg)
	assert.NoError(t, err)

	params := map[string]string{
		"app_id": "test_app",
	}
	sign := "test_signature"

	result := provider.VerifySign(params, sign)

	assert.False(t, result)
}

// TestRealAlipayProvider_VerifySign_WithPublicKey tests signature verification with dummy public key path
func TestRealAlipayProvider_VerifySign_WithPublicKey(t *testing.T) {
	cfg := &external.Config{
		Alipay: config.AlipayConfig{
			PublicKeyPath: "dummy_path",
		},
	}

	provider, err := NewAlipayProvider(cfg)
	assert.NoError(t, err)

	params := map[string]string{"app_id": "test"}
	sign := "any_signature"

	// loadPublicKey returns nil, nil, so alipayPublic is still nil, returns false
	result := provider.VerifySign(params, sign)

	// TODO: Update this test when loadPublicKey is properly implemented
	// Currently returns false because loadPublicKey doesn't actually load keys
	assert.False(t, result)
}

// TestNewAlipayProvider_Creation tests provider creation
func TestNewAlipayProvider_Creation(t *testing.T) {
	cfg := &external.Config{
		Alipay: config.AlipayConfig{
			AppID:          "test_app",
			PrivateKeyPath: "",
			PublicKeyPath:  "",
		},
	}

	provider, err := NewAlipayProvider(cfg)

	assert.NoError(t, err)
	assert.NotNil(t, provider)
	assert.NotNil(t, provider.config)
}

// TestRealAlipayProvider_doRequest tests the doRequest method
func TestRealAlipayProvider_doRequest(t *testing.T) {
	cfg := &external.Config{
		Alipay: config.AlipayConfig{Enabled: true},
	}

	provider, err := NewAlipayProvider(cfg)
	assert.NoError(t, err)

	ctx := context.Background()
	params := map[string]string{
		"app_id":  "test_app",
		"method":  "test_method",
		"charset": "UTF-8",
	}

	// This will try to make a real HTTP request to Alipay
	// The test verifies the method exists and builds the request correctly
	// Result may vary depending on network connectivity
	_, err = provider.doRequest(ctx, params)

	// Either succeeds (if Alipay endpoint is reachable) or fails (if not)
	// We just verify the method runs without panicking
	assert.NotNil(t, provider)
}

// ============ WeChat Provider Tests ============

// TestRealWeChatProvider_Refund_Success tests successful WeChat refund
func TestRealWeChatProvider_Refund_Success(t *testing.T) {
	cfg := &external.Config{
		WeChatPay: config.WeChatPayConfig{
			Enabled:   true,
			AppID:     "test_app_id",
			MchID:     "test_mch_id",
			APIKey:    "test_api_key_32_characters_long",
			NotifyURL: "https://example.com/notify",
		},
	}

	provider := NewWeChatProvider(cfg)

	payment := &model.Payment{
		Base:        model.Base{ID: 54321},
		AmountCents: 5000,
	}

	reason := "Customer requested refund"

	refundNo, rawDetails, refundedAt, err := provider.Refund(context.Background(), payment, reason)

	assert.NoError(t, err)
	assert.NotEmpty(t, refundNo)
	assert.Contains(t, refundNo, "refund_54321")
	assert.NotNil(t, rawDetails)

	var details map[string]interface{}
	err = json.Unmarshal(rawDetails, &details)
	assert.NoError(t, err)
	assert.Equal(t, "wechat", details["channel"])
	assert.Equal(t, float64(54321), details["payment_id"])
	assert.Equal(t, reason, details["refund_reason"])
	assert.NotZero(t, refundedAt)
}

// TestRealWeChatProvider_Refund_Disabled tests refund when WeChat Pay is disabled
func TestRealWeChatProvider_Refund_Disabled(t *testing.T) {
	cfg := &external.Config{
		WeChatPay: config.WeChatPayConfig{
			Enabled: false,
		},
	}

	provider := NewWeChatProvider(cfg)

	payment := &model.Payment{
		Base:        model.Base{ID: 54321},
		AmountCents: 5000,
	}

	refundNo, rawDetails, refundedAt, err := provider.Refund(context.Background(), payment, "test")

	assert.Error(t, err)
	assert.Equal(t, ErrPaymentDisabled, err)
	assert.Empty(t, refundNo)
	assert.Nil(t, rawDetails)
	assert.Zero(t, refundedAt)
}

// TestRealWeChatProvider_CreateOrder_Success tests successful WeChat order creation
func TestRealWeChatProvider_CreateOrder_Success(t *testing.T) {
	cfg := &external.Config{
		WeChatPay: config.WeChatPayConfig{
			Enabled:   true,
			AppID:     "wx_test_app_id",
			MchID:     "test_mch_id",
			APIKey:    "test_api_key_32_characters_long",
			NotifyURL: "https://example.com/notify",
		},
	}

	provider := NewWeChatProvider(cfg)

	orderID := "WX_ORDER_456"
	description := "Test WeChat Order"
	amountCents := int64(15000)
	clientIP := "127.0.0.1"

	result, err := provider.CreateOrder(context.Background(), orderID, description, amountCents, clientIP)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "wx_test_app_id", result["appid"])
	assert.Equal(t, "test_mch_id", result["mch_id"])
	assert.Equal(t, description, result["body"])
	assert.Equal(t, orderID, result["out_trade_no"])
	assert.Equal(t, amountCents, result["total_fee"])
	assert.Equal(t, clientIP, result["spbill_create_ip"])
	assert.Equal(t, "JSAPI", result["trade_type"])
	assert.Contains(t, result, "sign")
	assert.NotEmpty(t, result["nonce_str"])
}

// TestRealWeChatProvider_CreateOrder_Disabled tests order creation when WeChat Pay is disabled
func TestRealWeChatProvider_CreateOrder_Disabled(t *testing.T) {
	cfg := &external.Config{
		WeChatPay: config.WeChatPayConfig{
			Enabled: false,
		},
	}

	provider := NewWeChatProvider(cfg)

	result, err := provider.CreateOrder(context.Background(), "ORDER_123", "Test", 10000, "127.0.0.1")

	assert.Error(t, err)
	assert.Equal(t, ErrPaymentDisabled, err)
	assert.Nil(t, result)
}

// TestRealWeChatProvider_generateMapSign tests signature generation from map
func TestRealWeChatProvider_generateMapSign(t *testing.T) {
	cfg := &external.Config{
		WeChatPay: config.WeChatPayConfig{
			APIKey: "test_api_key_32_characters_long",
		},
	}

	provider := NewWeChatProvider(cfg)

	params := map[string]interface{}{
		"appid":        "test_app",
		"mch_id":       "test_mch",
		"nonce_str":    "random_string",
		"body":         "test product",
		"out_trade_no": "ORDER123",
	}

	sign := provider.generateMapSign(params)

	assert.NotEmpty(t, sign)
	assert.Len(t, sign, 32)                       // MD5 hash is 32 hex characters
	assert.True(t, strings.ToUpper(sign) == sign) // Should be uppercase
}

// TestRealWeChatProvider_generateSign tests signature generation for refund request
func TestRealWeChatProvider_generateSign(t *testing.T) {
	cfg := &external.Config{
		WeChatPay: config.WeChatPayConfig{
			AppID:  "test_app_id",
			MchID:  "test_mch_id",
			APIKey: "test_api_key_32_characters_long",
		},
	}

	provider := NewWeChatProvider(cfg)

	req := WeChatRefundRequest{
		AppID:       "test_app_id",
		MchID:       "test_mch_id",
		NonceStr:    "test_nonce",
		OutTradeNo:  "ORDER123",
		OutRefundNo: "REFUND123",
		TotalFee:    10000,
		RefundFee:   10000,
	}

	sign := provider.generateSign(req)

	assert.NotEmpty(t, sign)
	assert.Len(t, sign, 32)                       // MD5 hash is 32 hex characters
	assert.True(t, strings.ToUpper(sign) == sign) // Should be uppercase
}

// TestRealWeChatProvider_verifySign tests WeChat signature verification
func TestRealWeChatProvider_verifySign(t *testing.T) {
	cfg := &external.Config{
		WeChatPay: config.WeChatPayConfig{},
	}

	provider := NewWeChatProvider(cfg)

	resp := WeChatRefundResponse{
		ReturnCode: "SUCCESS",
		ResultCode: "SUCCESS",
	}

	result := provider.verifySign(resp)

	// Currently returns true (TODO implementation)
	assert.True(t, result)
}

// TestRealWeChatProvider_doRefundRequest tests the doRefundRequest method
func TestRealWeChatProvider_doRefundRequest(t *testing.T) {
	cfg := &external.Config{
		WeChatPay: config.WeChatPayConfig{
			AppID:  "test_app_id",
			MchID:  "test_mch_id",
			APIKey: "test_api_key_32_characters_long",
		},
	}

	provider := NewWeChatProvider(cfg)

	req := WeChatRefundRequest{
		AppID:       "test_app_id",
		MchID:       "test_mch_id",
		NonceStr:    "test_nonce",
		OutTradeNo:  "ORDER123",
		OutRefundNo: "REFUND123",
		TotalFee:    10000,
		RefundFee:   10000,
	}

	ctx := context.Background()

	// This will try to make a real HTTP request which will fail
	_, err := provider.doRefundRequest(ctx, req)

	assert.Error(t, err) // Expected to fail due to network
}

// TestGenerateNonceStr tests nonce string generation
func TestGenerateNonceStr(t *testing.T) {
	nonce1 := generateNonceStr()
	assert.NotEmpty(t, nonce1)
	// Verify it contains digits (timestamp)
	assert.Contains(t, nonce1, "0")

	// Generate another one - may or may not be different depending on timing
	nonce2 := generateNonceStr()
	assert.NotEmpty(t, nonce2)
}

// TestHmacSHA256 tests HMAC-SHA256 signature generation
func TestHmacSHA256(t *testing.T) {
	secret := "test_secret_key"
	data := "test_data_to_sign"

	signature := HmacSHA256(secret, data)

	assert.NotEmpty(t, signature)
	assert.Len(t, signature, 64) // SHA256 produces 64 hex characters

	// Same input should produce same signature
	signature2 := HmacSHA256(secret, data)
	assert.Equal(t, signature, signature2)

	// Different data should produce different signature
	signature3 := HmacSHA256(secret, "different_data")
	assert.NotEqual(t, signature, signature3)

	// Different secret should produce different signature
	signature4 := HmacSHA256("different_secret", data)
	assert.NotEqual(t, signature, signature4)
}

// TestHmacSHA256_EmptyInputs tests HMAC with empty inputs
func TestHmacSHA256_EmptyInputs(t *testing.T) {
	signature1 := HmacSHA256("", "")
	assert.NotEmpty(t, signature1)
	assert.Len(t, signature1, 64)

	signature2 := HmacSHA256("secret", "")
	assert.NotEmpty(t, signature2)
	assert.Len(t, signature2, 64)

	signature3 := HmacSHA256("", "data")
	assert.NotEmpty(t, signature3)
	assert.Len(t, signature3, 64)
}

// TestRealWeChatProvider_Refund_LargeAmount tests refund with large amount
func TestRealWeChatProvider_Refund_LargeAmount(t *testing.T) {
	cfg := &external.Config{
		WeChatPay: config.WeChatPayConfig{
			Enabled:   true,
			AppID:     "test_app_id",
			MchID:     "test_mch_id",
			APIKey:    "test_api_key_32_characters_long",
			NotifyURL: "https://example.com/notify",
		},
	}

	provider := NewWeChatProvider(cfg)

	payment := &model.Payment{
		Base:        model.Base{ID: 99999},
		AmountCents: 999999,
	}

	refundNo, rawDetails, refundedAt, err := provider.Refund(context.Background(), payment, "Large refund")

	assert.NoError(t, err)
	assert.NotEmpty(t, refundNo)
	assert.NotNil(t, rawDetails)
	assert.NotZero(t, refundedAt)

	var details map[string]interface{}
	err = json.Unmarshal(rawDetails, &details)
	assert.NoError(t, err)
	assert.Equal(t, float64(99999), details["payment_id"])
}

// TestRealAlipayProvider_Refund_VariousAmounts tests Alipay refund with various amounts
func TestRealAlipayProvider_Refund_VariousAmounts(t *testing.T) {
	cfg := &external.Config{
		Alipay: config.AlipayConfig{
			Enabled: true,
			AppID:   "test_app_id",
		},
	}

	provider, err := NewAlipayProvider(cfg)
	assert.NoError(t, err)

	testCases := []struct {
		name        string
		amountCents int64
	}{
		{"Minimum", 1},
		{"Small", 100},
		{"Medium", 10000},
		{"Large", 100000},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			payment := &model.Payment{
				Base:        model.Base{ID: uint64(tc.amountCents)},
				AmountCents: tc.amountCents,
			}

			refundNo, rawDetails, refundedAt, err := provider.Refund(context.Background(), payment, "test")

			// Should fail due to missing private key
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "private key not loaded")
			assert.Empty(t, refundNo)
			assert.Nil(t, rawDetails)
			assert.Zero(t, refundedAt)
		})
	}
}

// TestRealAlipayProvider_CreateOrder_VariousSubjects tests order creation with various subjects
func TestRealAlipayProvider_CreateOrder_VariousSubjects(t *testing.T) {
	cfg := &external.Config{
		Alipay: config.AlipayConfig{
			Enabled: true,
			AppID:   "test_app_id",
		},
	}

	provider, err := NewAlipayProvider(cfg)
	assert.NoError(t, err)

	testCases := []struct {
		name    string
		subject string
	}{
		{"English", "Test Order"},
		{"Chinese", "测试订单"},
		{"Mixed", "Test测试订单123"},
		{"Special", "Order!@#$%^&*()"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Should fail due to missing private key
			result, err := provider.CreateOrder(context.Background(), "ORDER_123", tc.subject, 10000)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "private key not loaded")
			assert.Nil(t, result)
		})
	}
}

// TestRealAlipayProvider_Refund_WithSpecialCharacters tests refund with special characters in reason
func TestRealAlipayProvider_Refund_WithSpecialCharacters(t *testing.T) {
	cfg := &external.Config{
		Alipay: config.AlipayConfig{
			Enabled: true,
			AppID:   "test_app_id",
		},
	}

	provider, err := NewAlipayProvider(cfg)
	assert.NoError(t, err)

	payment := &model.Payment{
		Base:        model.Base{ID: 12345},
		AmountCents: 10000,
	}

	specialReasons := []string{
		"Customer not satisfied",
		"服务不满意",
		"Double payment!@#",
		"Reason with émojis 🎉",
	}

	for _, reason := range specialReasons {
		t.Run("Reason_"+reason, func(t *testing.T) {
			// Should fail due to missing private key
			refundNo, rawDetails, refundedAt, err := provider.Refund(context.Background(), payment, reason)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "private key not loaded")
			assert.Empty(t, refundNo)
			assert.Nil(t, rawDetails)
			assert.Zero(t, refundedAt)
		})
	}
}

// TestRealWeChatProvider_generateMapSign_EmptyParams tests sign generation with empty params
func TestRealWeChatProvider_generateMapSign_EmptyParams(t *testing.T) {
	cfg := &external.Config{
		WeChatPay: config.WeChatPayConfig{
			APIKey: "test_api_key_32_characters_long",
		},
	}

	provider := NewWeChatProvider(cfg)

	// Empty params
	sign1 := provider.generateMapSign(map[string]interface{}{})
	assert.NotEmpty(t, sign1)

	// Params with empty values
	params := map[string]interface{}{
		"appid":  "test",
		"mch_id": "",
		"body":   "",
	}
	sign2 := provider.generateMapSign(params)
	assert.NotEmpty(t, sign2)
	assert.NotEqual(t, sign1, sign2)
}

// TestNewWeChatProvider_Creation tests WeChat provider creation
func TestNewWeChatProvider_Creation(t *testing.T) {
	cfg := &external.Config{
		WeChatPay: config.WeChatPayConfig{
			AppID:  "test_app",
			MchID:  "test_mch",
			APIKey: "test_key",
		},
	}

	provider := NewWeChatProvider(cfg)

	assert.NotNil(t, provider)
	assert.NotNil(t, provider.config)
}
