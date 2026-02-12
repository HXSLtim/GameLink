package payment

import (
	"encoding/xml"
	"testing"
)

// TestWeChatCallbackHandler_ParseCallback 测试解析回调通知
func TestWeChatCallbackHandler_ParseCallback(t *testing.T) {
	handler := NewWeChatCallbackHandler("test_app_id", "test_mch_id", "test_api_key")

	// 构造成功的回调 XML
	xmlData := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<xml>
  <return_code><![CDATA[SUCCESS]]></return_code>
  <return_msg><![CDATA[OK]]></return_msg>
  <appid><![CDATA[test_app_id]]></appid>
  <mch_id><![CDATA[test_mch_id]]></mch_id>
  <device_info><![CDATA[WEB]]></device_info>
  <nonce_str><![CDATA[test_nonce_str]]></nonce_str>
  <sign><![CDATA[TESTSIGN123456789012345678901234]]></sign>
  <result_code><![CDATA[SUCCESS]]></result_code>
  <openid><![CDATA[test_openid]]></openid>
  <is_subscribe><![CDATA[Y]]></is_subscribe>
  <trade_type><![CDATA[JSAPI]]></trade_type>
  <bank_type><![CDATA[CFT]]></bank_type>
  <total_fee><![CDATA[100]]></total_fee>
  <fee_type><![CDATA[CNY]]></fee_type>
  <cash_fee><![CDATA[100]]></cash_fee>
  <cash_fee_type><![CDATA[CNY]]></cash_fee_type>
  <transaction_id><![CDATA[WX2026020900123456789012345678]]></transaction_id>
  <out_trade_no><![CDATA[ORD20260209001]]></out_trade_no>
  <attach><![CDATA[custom_data]]></attach>
  <time_end><![CDATA[20260209120000]]></time_end>
</xml>`)

	notification, err := handler.ParseCallback(xmlData)
	if err != nil {
		t.Fatalf("ParseCallback failed: %v", err)
	}

	// 验证解析结果
	if notification.ReturnCode != "SUCCESS" {
		t.Errorf("expected return_code 'SUCCESS', got '%s'", notification.ReturnCode)
	}
	if notification.ResultCode != "SUCCESS" {
		t.Errorf("expected result_code 'SUCCESS', got '%s'", notification.ResultCode)
	}
	if notification.OutTradeNo != "ORD20260209001" {
		t.Errorf("expected out_trade_no 'ORD20260209001', got '%s'", notification.OutTradeNo)
	}
	if notification.TransactionID != "WX2026020900123456789012345678" {
		t.Errorf("expected transaction_id 'WX2026020900123456789012345678', got '%s'", notification.TransactionID)
	}
	if notification.TotalFee != "100" {
		t.Errorf("expected total_fee '100', got '%s'", notification.TotalFee)
	}
	if notification.Attach != "custom_data" {
		t.Errorf("expected attach 'custom_data', got '%s'", notification.Attach)
	}
}

// TestWeChatCallbackHandler_ParseCallback_FailReturnCode 测试解析失败的回调
func TestWeChatCallbackHandler_ParseCallback_FailReturnCode(t *testing.T) {
	handler := NewWeChatCallbackHandler("test_app_id", "test_mch_id", "test_api_key")

	// 构造失败的回调 XML（return_code != SUCCESS）
	xmlData := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<xml>
  <return_code><![CDATA[FAIL]]></return_code>
  <return_msg><![CDATA[签名失败]]></return_msg>
</xml>`)

	_, err := handler.ParseCallback(xmlData)
	if err == nil {
		t.Fatal("expected error for FAIL return_code, got nil")
	}

	expectedErrMsg := "return_code is not SUCCESS"
	if !contains(err.Error(), expectedErrMsg) {
		t.Errorf("expected error message to contain '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

// TestWeChatCallbackHandler_ParseCallback_FailResultCode 测试解析失败的回调
func TestWeChatCallbackHandler_ParseCallback_FailResultCode(t *testing.T) {
	handler := NewWeChatCallbackHandler("test_app_id", "test_mch_id", "test_api_key")

	// 构造失败的回调 XML（result_code != SUCCESS）
	xmlData := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<xml>
  <return_code><![CDATA[SUCCESS]]></return_code>
  <return_msg><![CDATA[OK]]></return_msg>
  <appid><![CDATA[test_app_id]]></appid>
  <mch_id><![CDATA[test_mch_id]]></mch_id>
  <nonce_str><![CDATA[test_nonce_str]]></nonce_str>
  <sign><![CDATA[TESTSIGN123456789012345678901234]]></sign>
  <result_code><![CDATA[FAIL]]></result_code>
  <err_code><![CDATA[NOTENOUGH]]></err_code>
  <err_code_des><![CDATA[余额不足]]></err_code_des>
  <out_trade_no><![CDATA[ORD20260209001]]></out_trade_no>
</xml>`)

	_, err := handler.ParseCallback(xmlData)
	if err == nil {
		t.Fatal("expected error for FAIL result_code, got nil")
	}

	expectedErrMsg := "result_code is not SUCCESS"
	if !contains(err.Error(), expectedErrMsg) {
		t.Errorf("expected error message to contain '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

// TestWeChatCallbackHandler_VerifyCallback 测试签名验证
func TestWeChatCallbackHandler_VerifyCallback(t *testing.T) {
	handler := NewWeChatCallbackHandler("test_app_id", "test_mch_id", "test_api_key")

	// 构造回调通知
	notification := &WeChatCallbackNotification{
		ReturnCode:    "SUCCESS",
		ReturnMsg:     "OK",
		AppID:         "test_app_id",
		MchID:         "test_mch_id",
		NonceStr:      "test_nonce_str",
		OutTradeNo:    "ORD20260209001",
		TransactionID: "WX2026020900123456789012345678",
		TotalFee:      "100",
	}

	// 生成正确的签名
	signClient := NewWeChatClient("test_app_id", "test_mch_id", "test_api_key", "")
	params := map[string]string{
		"return_code":    "SUCCESS",
		"return_msg":     "OK",
		"appid":          "test_app_id",
		"mch_id":         "test_mch_id",
		"nonce_str":      "test_nonce_str",
		"out_trade_no":   "ORD20260209001",
		"transaction_id": "WX2026020900123456789012345678",
		"total_fee":      "100",
	}
	notification.Sign = signClient.GenerateSign(params)

	// 验证签名应该成功
	if !handler.VerifyCallback(notification) {
		t.Error("signature verification should succeed")
	}

	// 修改签名，验证应该失败
	notification.Sign = "INVALID_SIGNATURE"
	if handler.VerifyCallback(notification) {
		t.Error("signature verification should fail for invalid signature")
	}
}

// TestWeChatCallbackHandler_CheckReplay 测试重放攻击检测
func TestWeChatCallbackHandler_CheckReplay(t *testing.T) {
	handler := NewWeChatCallbackHandler("test_app_id", "test_mch_id", "test_api_key")

	notification := &WeChatCallbackNotification{
		TransactionID: "WX2026020900123456789012345678",
	}

	// 第一次检查，不应该检测到重放
	if handler.CheckReplay(notification) {
		t.Error("first callback should not be detected as replay")
	}

	// 第二次检查，应该检测到重放
	if !handler.CheckReplay(notification) {
		t.Error("second callback with same transaction_id should be detected as replay")
	}
}

// TestWeChatCallbackHandler_ValidateCallback 测试完整的回调验证流程
func TestWeChatCallbackHandler_ValidateCallback(t *testing.T) {
	handler := NewWeChatCallbackHandler("test_app_id", "test_mch_id", "test_api_key")

	// 构造完整的回调 XML
	signClient := NewWeChatClient("test_app_id", "test_mch_id", "test_api_key", "")
	params := map[string]string{
		"return_code":    "SUCCESS",
		"return_msg":     "OK",
		"appid":          "test_app_id",
		"mch_id":         "test_mch_id",
		"nonce_str":      "test_nonce_str",
		"result_code":    "SUCCESS",
		"out_trade_no":   "ORD20260209001",
		"transaction_id": "WX2026020900123456789012345678",
		"total_fee":      "100",
		"time_end":       "20260209120000",
	}
	sign := signClient.GenerateSign(params)

	type xmlCallback struct {
		XMLName       xml.Name `xml:"xml"`
		ReturnCode    string   `xml:"return_code"`
		ReturnMsg     string   `xml:"return_msg"`
		AppID         string   `xml:"appid"`
		MchID         string   `xml:"mch_id"`
		NonceStr      string   `xml:"nonce_str"`
		Sign          string   `xml:"sign"`
		ResultCode    string   `xml:"result_code"`
		OutTradeNo    string   `xml:"out_trade_no"`
		TransactionID string   `xml:"transaction_id"`
		TotalFee      string   `xml:"total_fee"`
		TimeEnd       string   `xml:"time_end"`
	}

	xmlData, err := xml.Marshal(xmlCallback{
		ReturnCode:    "SUCCESS",
		ReturnMsg:     "OK",
		AppID:         "test_app_id",
		MchID:         "test_mch_id",
		NonceStr:      "test_nonce_str",
		Sign:          sign,
		ResultCode:    "SUCCESS",
		OutTradeNo:    "ORD20260209001",
		TransactionID: "WX2026020900123456789012345678",
		TotalFee:      "100",
		TimeEnd:       "20260209120000",
	})
	if err != nil {
		t.Fatalf("failed to marshal XML: %v", err)
	}

	// 验证回调
	notification, err := handler.ValidateCallback(xmlData)
	if err != nil {
		t.Fatalf("ValidateCallback failed: %v", err)
	}

	// 验证解析结果
	if notification.OutTradeNo != "ORD20260209001" {
		t.Errorf("expected out_trade_no 'ORD20260209001', got '%s'", notification.OutTradeNo)
	}
	if notification.TransactionID != "WX2026020900123456789012345678" {
		t.Errorf("expected transaction_id 'WX2026020900123456789012345678', got '%s'", notification.TransactionID)
	}
}

// TestWeChatCallbackHandler_ValidateCallback_InvalidSignature 测试无效签名
func TestWeChatCallbackHandler_ValidateCallback_InvalidSignature(t *testing.T) {
	handler := NewWeChatCallbackHandler("test_app_id", "test_mch_id", "test_api_key")

	xmlData := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<xml>
  <return_code><![CDATA[SUCCESS]]></return_code>
  <return_msg><![CDATA[OK]]></return_msg>
  <appid><![CDATA[test_app_id]]></appid>
  <mch_id><![CDATA[test_mch_id]]></mch_id>
  <nonce_str><![CDATA[test_nonce_str]]></nonce_str>
  <sign><![CDATA[INVALID_SIGNATURE]]></sign>
  <result_code><![CDATA[SUCCESS]]></result_code>
  <out_trade_no><![CDATA[ORD20260209001]]></out_trade_no>
  <transaction_id><![CDATA[WX2026020900123456789012345678]]></transaction_id>
</xml>`)

	_, err := handler.ValidateCallback(xmlData)
	if err == nil {
		t.Fatal("expected error for invalid signature, got nil")
	}

	expectedErrMsg := "signature verification failed"
	if !contains(err.Error(), expectedErrMsg) {
		t.Errorf("expected error message to contain '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

// TestWeChatCallbackHandler_Getters 测试 getter 方法
func TestWeChatCallbackHandler_Getters(t *testing.T) {
	handler := NewWeChatCallbackHandler("test_app_id", "test_mch_id", "test_api_key")

	notification := &WeChatCallbackNotification{
		OutTradeNo:    "ORD20260209001",
		TransactionID: "WX2026020900123456789012345678",
		TotalFee:      "100",
		TimeEnd:       "20260209120000",
		Attach:        "custom_data",
	}

	// 测试 GetCallbackUniqueKey
	if handler.GetCallbackUniqueKey(notification) != "WX2026020900123456789012345678" {
		t.Error("GetCallbackUniqueKey failed")
	}

	// 测试 GetOrderID
	if handler.GetOrderID(notification) != "ORD20260209001" {
		t.Error("GetOrderID failed")
	}

	// 测试 GetTransactionID
	if handler.GetTransactionID(notification) != "WX2026020900123456789012345678" {
		t.Error("GetTransactionID failed")
	}

	// 测试 GetTotalFee
	if handler.GetTotalFee(notification) != 100 {
		t.Errorf("GetTotalFee failed: expected 100, got %d", handler.GetTotalFee(notification))
	}

	// 测试 GetPaymentTime
	if handler.GetPaymentTime(notification) != "20260209120000" {
		t.Error("GetPaymentTime failed")
	}

	// 测试 GetAttach
	if handler.GetAttach(notification) != "custom_data" {
		t.Error("GetAttach failed")
	}
}

// TestWeChatCallbackResponse_ToXML 测试响应 XML 序列化
func TestWeChatCallbackResponse_ToXML(t *testing.T) {
	response := NewWeChatCallbackResponse()
	xmlData, err := response.ToXML()
	if err != nil {
		t.Fatalf("ToXML failed: %v", err)
	}

	// 验证 XML 包含正确的标签
	xmlStr := string(xmlData)
	if !contains(xmlStr, "<return_code>") || !contains(xmlStr, "SUCCESS") {
		t.Error("XML should contain return_code SUCCESS")
	}
	if !contains(xmlStr, "<return_msg>") || !contains(xmlStr, "OK") {
		t.Error("XML should contain return_msg OK")
	}
}

// TestWeChatCallbackResponse_FailResponse 测试失败响应
func TestWeChatCallbackResponse_FailResponse(t *testing.T) {
	response := NewWeChatCallbackFailResponse("签名验证失败")
	if response.ReturnCode != "FAIL" {
		t.Errorf("expected return_code 'FAIL', got '%s'", response.ReturnCode)
	}
	if response.ReturnMsg != "签名验证失败" {
		t.Errorf("expected return_msg '签名验证失败', got '%s'", response.ReturnMsg)
	}
}

// 辅助函数：检查字符串是否包含子字符串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || indexOf(s, substr) >= 0))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
