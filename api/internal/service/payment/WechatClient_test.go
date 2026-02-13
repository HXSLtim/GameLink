package payment

import (
	"testing"
)

// TestWeChatClient_NewWeChatClient 测试创建微信支付客户端
func TestWeChatClient_NewWeChatClient(t *testing.T) {
	client := NewWeChatClient(
		"test_app_id",
		"test_mch_id",
		"test_api_key",
		"https://example.com/notify",
	)

	if client == nil {
		t.Fatal("NewWeChatClient returned nil")
	}

	if client.appID != "test_app_id" {
		t.Errorf("expected app_id 'test_app_id', got '%s'", client.appID)
	}

	if client.mchID != "test_mch_id" {
		t.Errorf("expected mch_id 'test_mch_id', got '%s'", client.mchID)
	}

	if client.apiKey != "test_api_key" {
		t.Errorf("expected apiKey 'test_api_key', got '%s'", client.apiKey)
	}

	if client.notifyURL != "https://example.com/notify" {
		t.Errorf("expected notifyURL 'https://example.com/notify', got '%s'", client.notifyURL)
	}
}

// TestWeChatClient_GenerateSign 测试签名生成
func TestWeChatClient_GenerateSign(t *testing.T) {
	client := NewWeChatClient(
		"test_app_id",
		"test_mch_id",
		"test_api_key",
		"https://example.com/notify",
	)

	params := map[string]string{
		"appid":        "wx1234567890",
		"mch_id":       "1234567890",
		"nonce_str":    "test_nonce",
		"body":         "测试商品",
		"out_trade_no": "ORDER123",
		"total_fee":    "100",
	}

	sign := client.GenerateSign(params)

	if sign == "" {
		t.Error("GenerateSign returned empty string")
	}

	// 签名应该是 32 位大写 MD5 字符串
	if len(sign) != 32 {
		t.Errorf("expected signature length 32, got %d", len(sign))
	}

	// 签名应该是大写
	for _, ch := range sign {
		if ch >= 'a' && ch <= 'z' {
			t.Error("signature should be uppercase")
			break
		}
	}
}

// TestWeChatClient_GenerateSign_Sorted 测试签名生成（验证字典序）
func TestWeChatClient_GenerateSign_Sorted(t *testing.T) {
	client := NewWeChatClient(
		"test_app_id",
		"test_mch_id",
		"test_api_key",
		"https://example.com/notify",
	)

	// 相同参数，不同顺序
	params1 := map[string]string{
		"appid":        "wx1234567890",
		"mch_id":       "1234567890",
		"nonce_str":    "test_nonce",
		"body":         "测试商品",
		"out_trade_no": "ORDER123",
		"total_fee":    "100",
	}

	params2 := map[string]string{
		"total_fee":    "100",
		"appid":        "wx1234567890",
		"out_trade_no": "ORDER123",
		"nonce_str":    "test_nonce",
		"mch_id":       "1234567890",
		"body":         "测试商品",
	}

	sign1 := client.GenerateSign(params1)
	sign2 := client.GenerateSign(params2)

	if sign1 != sign2 {
		t.Errorf("signature should be same regardless of param order: %s != %s", sign1, sign2)
	}
}

// TestWeChatClient_VerifySign 测试签名验证
func TestWeChatClient_VerifySign(t *testing.T) {
	client := NewWeChatClient(
		"test_app_id",
		"test_mch_id",
		"test_api_key",
		"https://example.com/notify",
	)

	params := map[string]string{
		"appid":        "wx1234567890",
		"mch_id":       "1234567890",
		"nonce_str":    "test_nonce",
		"body":         "测试商品",
		"out_trade_no": "ORDER123",
		"total_fee":    "100",
	}

	// 生成签名
	sign := client.GenerateSign(params)
	params["sign"] = sign

	// 验证签名
	if !client.VerifySign(params) {
		t.Error("VerifySign failed for valid signature")
	}

	// 修改签名
	params["sign"] = "INVALID_SIGN"
	if client.VerifySign(params) {
		t.Error("VerifySign should fail for invalid signature")
	}

	// 删除签名
	delete(params, "sign")
	if client.VerifySign(params) {
		t.Error("VerifySign should fail when sign is missing")
	}
}

// TestWeChatClient_BuildUnifiedOrderParams 测试构建统一下单参数
func TestWeChatClient_BuildUnifiedOrderParams(t *testing.T) {
	client := NewWeChatClient(
		"wx_test_app_id",
		"test_mch_id",
		"test_api_key",
		"https://example.com/notify",
	)

	params := client.BuildUnifiedOrderParams(
		"ORDER123",
		"测试商品",
		"127.0.0.1",
		"JSAPI",
		100,
		"test_openid",
	)

	// 验证必需参数
	requiredFields := map[string]string{
		"appid":            "wx_test_app_id",
		"mch_id":           "test_mch_id",
		"body":             "测试商品",
		"out_trade_no":     "ORDER123",
		"total_fee":        "100",
		"spbill_create_ip": "127.0.0.1",
		"notify_url":       "https://example.com/notify",
		"trade_type":       "JSAPI",
		"openid":           "test_openid",
	}

	for key, expectedValue := range requiredFields {
		if params[key] != expectedValue {
			t.Errorf("expected %s '%s', got '%s'", key, expectedValue, params[key])
		}
	}

	// 验证签名已生成
	if params["sign"] == "" {
		t.Error("sign should be generated")
	}

	// 验证签名正确性
	if !client.VerifySign(params) {
		t.Error("generated signature is invalid")
	}
}

// TestWeChatClient_BuildUnifiedOrderParams_Native 测试 Native 支付（不需要 openID）
func TestWeChatClient_BuildUnifiedOrderParams_Native(t *testing.T) {
	client := NewWeChatClient(
		"wx_test_app_id",
		"test_mch_id",
		"test_api_key",
		"https://example.com/notify",
	)

	params := client.BuildUnifiedOrderParams(
		"ORDER123",
		"测试商品",
		"127.0.0.1",
		"NATIVE",
		100,
		"", // Native 支付不需要 openID
	)

	if params["trade_type"] != "NATIVE" {
		t.Errorf("expected trade_type 'NATIVE', got '%s'", params["trade_type"])
	}

	if params["openid"] != "" {
		t.Error("openid should be empty for NATIVE payment")
	}
}

// TestWeChatClient_BuildOrderQueryParams 测试构建订单查询参数
func TestWeChatClient_BuildOrderQueryParams(t *testing.T) {
	client := NewWeChatClient(
		"wx_test_app_id",
		"test_mch_id",
		"test_api_key",
		"https://example.com/notify",
	)

	params := client.BuildOrderQueryParams("ORDER123", "WX_TRANSACTION")

	// 验证必需参数
	requiredFields := map[string]string{
		"appid":          "wx_test_app_id",
		"mch_id":         "test_mch_id",
		"out_trade_no":   "ORDER123",
		"transaction_id": "WX_TRANSACTION",
	}

	for key, expectedValue := range requiredFields {
		if params[key] != expectedValue {
			t.Errorf("expected %s '%s', got '%s'", key, expectedValue, params[key])
		}
	}

	// 验证签名已生成
	if params["sign"] == "" {
		t.Error("sign should be generated")
	}

	// 验证签名正确性
	if !client.VerifySign(params) {
		t.Error("generated signature is invalid")
	}
}

// TestWeChatClient_BuildRefundParams 测试构建退款参数
func TestWeChatClient_BuildRefundParams(t *testing.T) {
	client := NewWeChatClient(
		"wx_test_app_id",
		"test_mch_id",
		"test_api_key",
		"https://example.com/notify",
	)

	params := client.BuildRefundParams("ORDER123", "REFUND123", 100, 100)

	// 验证必需参数
	requiredFields := map[string]string{
		"appid":         "wx_test_app_id",
		"mch_id":        "test_mch_id",
		"out_trade_no":  "ORDER123",
		"out_refund_no": "REFUND123",
		"total_fee":     "100",
		"refund_fee":    "100",
	}

	for key, expectedValue := range requiredFields {
		if params[key] != expectedValue {
			t.Errorf("expected %s '%s', got '%s'", key, expectedValue, params[key])
		}
	}

	// 验证签名已生成
	if params["sign"] == "" {
		t.Error("sign should be generated")
	}

	// 验证签名正确性
	if !client.VerifySign(params) {
		t.Error("generated signature is invalid")
	}
}
