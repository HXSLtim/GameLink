package payment

import (
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestWeChatAPIClient_NewWeChatAPIClient 测试创建微信支付 API 客户端
func TestWeChatAPIClient_NewWeChatAPIClient(t *testing.T) {
	client := NewWeChatAPIClient(
		"test_app_id",
		"test_mch_id",
		"test_api_key",
		"/path/to/cert.pem",
		"/path/to/key.pem",
		"http://example.com/notify",
	)

	if client == nil {
		t.Fatal("NewWeChatAPIClient returned nil")
	}

	if client.appID != "test_app_id" {
		t.Errorf("expected appID 'test_app_id', got '%s'", client.appID)
	}

	if client.mchID != "test_mch_id" {
		t.Errorf("expected mchID 'test_mch_id', got '%s'", client.mchID)
	}

	if client.apiKey != "test_api_key" {
		t.Errorf("expected apiKey 'test_api_key', got '%s'", client.apiKey)
	}

	if client.certPath != "/path/to/cert.pem" {
		t.Errorf("expected certPath '/path/to/cert.pem', got '%s'", client.certPath)
	}

	if client.certKey != "/path/to/key.pem" {
		t.Errorf("expected certKey '/path/to/key.pem', got '%s'", client.certKey)
	}

	if client.notifyURL != "http://example.com/notify" {
		t.Errorf("expected notifyURL 'http://example.com/notify', got '%s'", client.notifyURL)
	}

	if client.httpClient == nil {
		t.Error("httpClient should not be nil")
	}

	if client.httpClient.Timeout != 30*time.Second {
		t.Errorf("expected timeout 30s, got %v", client.httpClient.Timeout)
	}
}

// TestWeChatAPIClient_UnifiedOrder 测试统一下单接口
func TestWeChatAPIClient_UnifiedOrder(t *testing.T) {
	// 创建 mock 服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求方法
		if r.Method != "POST" {
			t.Errorf("expected POST method, got %s", r.Method)
		}

		// 验证 Content-Type
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/xml" {
			t.Errorf("expected Content-Type 'application/xml', got '%s'", contentType)
		}

		// 解析请求体
		body := &struct {
			XMLName  xml.Name `xml:"xml"`
			AppID    string   `xml:"appid"`
			MchID    string   `xml:"mch_id"`
			NonceStr string   `xml:"nonce_str"`
			Sign     string   `xml:"sign"`
		}{}
		if err := xml.NewDecoder(r.Body).Decode(body); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		// 验证必需参数
		if body.AppID != "test_app_id" {
			t.Errorf("expected appID 'test_app_id', got '%s'", body.AppID)
		}
		if body.MchID != "test_mch_id" {
			t.Errorf("expected mchID 'test_mch_id', got '%s'", body.MchID)
		}
		if body.NonceStr == "" {
			t.Error("nonce_str should not be empty")
		}
		if body.Sign == "" {
			t.Error("sign should not be empty")
		}
		if len(body.Sign) != 32 {
			t.Errorf("expected sign length 32, got %d", len(body.Sign))
		}

		// 返回成功响应
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		response := `<?xml version="1.0" encoding="UTF-8"?>
<xml>
  <return_code><![CDATA[SUCCESS]]></return_code>
  <return_msg><![CDATA[OK]]></return_msg>
  <appid><![CDATA[test_app_id]]></appid>
  <mch_id><![CDATA[test_mch_id]]></mch_id>
  <nonce_str><![CDATA[test_nonce]]></nonce_str>
  <sign><![CDATA[TESTSIGN123456789012345678901234]]></sign>
  <result_code><![CDATA[SUCCESS]]></result_code>
  <prepay_id><![CDATA[wxprepays12345678901234567890]]></prepay_id>
  <trade_type><![CDATA[JSAPI]]></trade_type>
</xml>`
		w.Write([]byte(response))
	}))
	defer server.Close()

	// 创建客户端（使用 mock 服务器 URL）
	client := &WeChatAPIClient{
		appID:      "test_app_id",
		mchID:      "test_mch_id",
		apiKey:     "test_api_key",
		notifyURL:  "http://example.com/notify",
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	// 准备请求参数
	req := UnifiedOrderRequest{
		Body:           "测试商品",
		OutTradeNo:     "ORD20260209001",
		TotalFee:       100, // 1元
		SpbillCreateIP: "127.0.0.1",
		TradeType:      "JSAPI",
		OpenID:         "test_openid",
		NonceStr:       generateNonceStr(),
	}

	// 调用统一下单（使用 mock URL）
	// 注意：这里需要修改 doRequest 方法支持自定义 URL，或者在测试中直接调用
	// 为了测试，我们创建一个辅助方法
	resp, err := client.unifiedOrderWithServer(server.URL, req)
	if err != nil {
		t.Fatalf("UnifiedOrder failed: %v", err)
	}

	// 验证响应
	if resp.ReturnCode != "SUCCESS" {
		t.Errorf("expected return_code 'SUCCESS', got '%s'", resp.ReturnCode)
	}
	if resp.ReturnMsg != "OK" {
		t.Errorf("expected return_msg 'OK', got '%s'", resp.ReturnMsg)
	}
	if resp.ResultCode != "SUCCESS" {
		t.Errorf("expected result_code 'SUCCESS', got '%s'", resp.ResultCode)
	}
	if resp.PrepayID == "" {
		t.Error("prepay_id should not be empty")
	}
	if resp.TradeType != "JSAPI" {
		t.Errorf("expected trade_type 'JSAPI', got '%s'", resp.TradeType)
	}
}

// unifiedOrderWithServer 使用指定服务器 URL 的测试辅助方法
func (c *WeChatAPIClient) unifiedOrderWithServer(serverURL string, req UnifiedOrderRequest) (*UnifiedOrderResponse, error) {
	// 构建请求参数
	type XMLParams struct {
		XMLName         xml.Name `xml:"xml"`
		AppID           string   `xml:"appid"`
		MchID           string   `xml:"mch_id"`
		NonceStr        string   `xml:"nonce_str"`
		Body            string   `xml:"body"`
		OutTradeNo      string   `xml:"out_trade_no"`
		TotalFee        string   `xml:"total_fee"`
		SpbillCreateIP  string   `xml:"spbill_create_ip"`
		NotifyURL       string   `xml:"notify_url"`
		TradeType       string   `xml:"trade_type"`
		OpenID          string   `xml:"openid,omitempty"`
		Sign            string   `xml:"sign"`
	}

	// 生成签名
	signClient := NewWeChatClient(c.appID, c.mchID, c.apiKey, c.notifyURL)
	params := map[string]string{
		"appid":            c.appID,
		"mch_id":           c.mchID,
		"nonce_str":        req.NonceStr,
		"body":             req.Body,
		"out_trade_no":     req.OutTradeNo,
		"total_fee":        "100",
		"spbill_create_ip": req.SpbillCreateIP,
		"notify_url":       c.notifyURL,
		"trade_type":       req.TradeType,
	}
	if req.OpenID != "" {
		params["openid"] = req.OpenID
	}
	sign := signClient.GenerateSign(params)

	xmlParams := XMLParams{
		AppID:           c.appID,
		MchID:           c.mchID,
		NonceStr:        req.NonceStr,
		Body:            req.Body,
		OutTradeNo:      req.OutTradeNo,
		TotalFee:        "100",
		SpbillCreateIP:  req.SpbillCreateIP,
		NotifyURL:       c.notifyURL,
		TradeType:       req.TradeType,
		OpenID:          req.OpenID,
		Sign:            sign,
	}

	xmlData, err := xml.Marshal(xmlParams)
	if err != nil {
		return nil, err
	}

	// 发送请求到 mock 服务器
	httpReq, err := http.NewRequest("POST", serverURL, strings.NewReader(string(xmlData)))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/xml")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var response UnifiedOrderResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// TestWeChatAPIClient_OrderQuery 测试订单查询接口
func TestWeChatAPIClient_OrderQuery(t *testing.T) {
	// 创建 mock 服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 返回成功响应
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		response := `<?xml version="1.0" encoding="UTF-8"?>
<xml>
  <return_code><![CDATA[SUCCESS]]></return_code>
  <return_msg><![CDATA[OK]]></return_msg>
  <appid><![CDATA[test_app_id]]></appid>
  <mch_id><![CDATA[test_mch_id]]></mch_id>
  <nonce_str><![CDATA[test_nonce]]></nonce_str>
  <sign><![CDATA[TESTSIGN123456789012345678901234]]></sign>
  <result_code><![CDATA[SUCCESS]]></result_code>
  <out_trade_no><![CDATA[ORD20260209001]]></out_trade_no>
  <transaction_id><![CDATA[WX2026020900123456789012345678]]></transaction_id>
  <trade_state><![CDATA[SUCCESS]]></trade_state>
  <total_fee><![CDATA[100]]></total_fee>
  <time_end><![CDATA[20260209120000]]></time_end>
</xml>`
		w.Write([]byte(response))
	}))
	defer server.Close()

	client := &WeChatAPIClient{
		appID:      "test_app_id",
		mchID:      "test_mch_id",
		apiKey:     "test_api_key",
		notifyURL:  "http://example.com/notify",
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	resp, err := client.orderQueryWithServer(server.URL, "ORD20260209001", "")
	if err != nil {
		t.Fatalf("OrderQuery failed: %v", err)
	}

	// 验证响应
	if resp.ReturnCode != "SUCCESS" {
		t.Errorf("expected return_code 'SUCCESS', got '%s'", resp.ReturnCode)
	}
	if resp.ResultCode != "SUCCESS" {
		t.Errorf("expected result_code 'SUCCESS', got '%s'", resp.ResultCode)
	}
	if resp.OutTradeNo != "ORD20260209001" {
		t.Errorf("expected out_trade_no 'ORD20260209001', got '%s'", resp.OutTradeNo)
	}
	if resp.TransactionID == "" {
		t.Error("transaction_id should not be empty")
	}
	if resp.TradeState != "SUCCESS" {
		t.Errorf("expected trade_state 'SUCCESS', got '%s'", resp.TradeState)
	}
	if resp.TotalFee != "100" {
		t.Errorf("expected total_fee '100', got '%s'", resp.TotalFee)
	}
}

// orderQueryWithServer 使用指定服务器 URL 的测试辅助方法
func (c *WeChatAPIClient) orderQueryWithServer(serverURL, outTradeNo, transactionID string) (*OrderQueryResponse, error) {
	type XMLParams struct {
		XMLName       xml.Name `xml:"xml"`
		AppID         string   `xml:"appid"`
		MchID         string   `xml:"mch_id"`
		NonceStr      string   `xml:"nonce_str"`
		OutTradeNo    string   `xml:"out_trade_no"`
		TransactionID string   `xml:"transaction_id,omitempty"`
		Sign          string   `xml:"sign"`
	}

	nonceStr := generateNonceStr()
	signClient := NewWeChatClient(c.appID, c.mchID, c.apiKey, c.notifyURL)
	params := map[string]string{
		"appid":       c.appID,
		"mch_id":      c.mchID,
		"nonce_str":   nonceStr,
		"out_trade_no": outTradeNo,
	}
	if transactionID != "" {
		params["transaction_id"] = transactionID
	}
	sign := signClient.GenerateSign(params)

	xmlParams := XMLParams{
		AppID:         c.appID,
		MchID:         c.mchID,
		NonceStr:      nonceStr,
		OutTradeNo:    outTradeNo,
		TransactionID: transactionID,
		Sign:          sign,
	}

	xmlData, err := xml.Marshal(xmlParams)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", serverURL, strings.NewReader(string(xmlData)))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/xml")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response OrderQueryResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// TestWeChatAPIClient_Refund 测试退款接口
func TestWeChatAPIClient_Refund(t *testing.T) {
	// 创建 mock 服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求方法
		if r.Method != "POST" {
			t.Errorf("expected POST method, got %s", r.Method)
		}

		// 返回成功响应
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		response := `<?xml version="1.0" encoding="UTF-8"?>
<xml>
  <return_code><![CDATA[SUCCESS]]></return_code>
  <return_msg><![CDATA[OK]]></return_msg>
  <appid><![CDATA[test_app_id]]></appid>
  <mch_id><![CDATA[test_mch_id]]></mch_id>
  <nonce_str><![CDATA[test_nonce]]></nonce_str>
  <sign><![CDATA[TESTSIGN123456789012345678901234]]></sign>
  <result_code><![CDATA[SUCCESS]]></result_code>
  <out_trade_no><![CDATA[ORD20260209001]]></out_trade_no>
  <out_refund_no><![CDATA[REF20260209001]]></out_refund_no>
  <refund_id><![CDATA[WXREFUND20260209001234567890]]></refund_id>
  <refund_channel><![CDATA[ORIGINAL]]></refund_channel>
  <refund_status><![CDATA[SUCCESS]]></refund_status>
</xml>`
		w.Write([]byte(response))
	}))
	defer server.Close()

	client := &WeChatAPIClient{
		appID:      "test_app_id",
		mchID:      "test_mch_id",
		apiKey:     "test_api_key",
		notifyURL:  "http://example.com/notify",
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}

	req := RefundRequest{
		OutTradeNo:  "ORD20260209001",
		OutRefundNo: "REF20260209001",
		TotalFee:    100,
		RefundFee:   100,
		NonceStr:    generateNonceStr(),
	}

	resp, err := client.refundWithServer(server.URL, req)
	if err != nil {
		t.Fatalf("Refund failed: %v", err)
	}

	// 验证响应
	if resp.ReturnCode != "SUCCESS" {
		t.Errorf("expected return_code 'SUCCESS', got '%s'", resp.ReturnCode)
	}
	if resp.ResultCode != "SUCCESS" {
		t.Errorf("expected result_code 'SUCCESS', got '%s'", resp.ResultCode)
	}
	if resp.OutTradeNo != "ORD20260209001" {
		t.Errorf("expected out_trade_no 'ORD20260209001', got '%s'", resp.OutTradeNo)
	}
	if resp.OutRefundNo != "REF20260209001" {
		t.Errorf("expected out_refund_no 'REF20260209001', got '%s'", resp.OutRefundNo)
	}
	if resp.RefundID == "" {
		t.Error("refund_id should not be empty")
	}
	if resp.RefundStatus != "SUCCESS" {
		t.Errorf("expected refund_status 'SUCCESS', got '%s'", resp.RefundStatus)
	}
}

// refundWithServer 使用指定服务器 URL 的测试辅助方法
func (c *WeChatAPIClient) refundWithServer(serverURL string, req RefundRequest) (*RefundResponse, error) {
	type XMLParams struct {
		XMLName      xml.Name `xml:"xml"`
		AppID        string   `xml:"appid"`
		MchID        string   `xml:"mch_id"`
		NonceStr     string   `xml:"nonce_str"`
		OutTradeNo   string   `xml:"out_trade_no"`
		OutRefundNo  string   `xml:"out_refund_no"`
		TotalFee     string   `xml:"total_fee"`
		RefundFee    string   `xml:"refund_fee"`
		Sign         string   `xml:"sign"`
	}

	signClient := NewWeChatClient(c.appID, c.mchID, c.apiKey, c.notifyURL)
	params := map[string]string{
		"appid":         c.appID,
		"mch_id":        c.mchID,
		"nonce_str":     req.NonceStr,
		"out_trade_no":  req.OutTradeNo,
		"out_refund_no": req.OutRefundNo,
		"total_fee":     "100",
		"refund_fee":    "100",
	}
	sign := signClient.GenerateSign(params)

	xmlParams := XMLParams{
		AppID:       c.appID,
		MchID:       c.mchID,
		NonceStr:    req.NonceStr,
		OutTradeNo:  req.OutTradeNo,
		OutRefundNo: req.OutRefundNo,
		TotalFee:    "100",
		RefundFee:   "100",
		Sign:        sign,
	}

	xmlData, err := xml.Marshal(xmlParams)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", serverURL, strings.NewReader(string(xmlData)))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/xml")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response RefundResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// TestWeChatAPIClient_ResponseTypes 测试响应类型定义
func TestWeChatAPIClient_ResponseTypes(t *testing.T) {
	// 测试 UnifiedOrderResponse
	unifiedResp := UnifiedOrderResponse{
		ReturnCode: "SUCCESS",
		ReturnMsg:  "OK",
		AppID:      "test_app_id",
		MchID:      "test_mch_id",
		NonceStr:   "test_nonce",
		Sign:       "TESTSIGN123456789012345678901234",
		ResultCode: "SUCCESS",
		PrepayID:   "wxprepays12345678901234567890",
		CodeURL:    "weixin://wxpay/bizpayurl?pr=xxxx",
		TradeType:  "NATIVE",
	}

	if unifiedResp.ReturnCode != "SUCCESS" {
		t.Error("UnifiedOrderResponse struct definition error")
	}

	// 测试 OrderQueryResponse
	queryResp := OrderQueryResponse{
		ReturnCode:    "SUCCESS",
		ReturnMsg:     "OK",
		ResultCode:    "SUCCESS",
		TradeState:    "SUCCESS",
		TransactionID: "WX2026020900123456789012345678",
		OutTradeNo:    "ORD20260209001",
		TotalFee:      "100",
		TimeEnd:       "20260209120000",
	}

	if queryResp.TradeState != "SUCCESS" {
		t.Error("OrderQueryResponse struct definition error")
	}

	// 测试 RefundResponse
	refundResp := RefundResponse{
		ReturnCode:    "SUCCESS",
		ReturnMsg:     "OK",
		ResultCode:    "SUCCESS",
		OutTradeNo:    "ORD20260209001",
		OutRefundNo:   "REF20260209001",
		RefundID:      "WXREFUND20260209001234567890",
		RefundStatus:  "SUCCESS",
		RefundChannel: "ORIGINAL",
	}

	if refundResp.RefundStatus != "SUCCESS" {
		t.Error("RefundResponse struct definition error")
	}
}

// TestWeChatAPIClient_UnifiedOrder_Native 测试 NATIVE 支付类型
func TestWeChatAPIClient_UnifiedOrder_Native(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 返回 NATIVE 支付响应（包含 code_url）
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		response := `<?xml version="1.0" encoding="UTF-8"?>
<xml>
  <return_code><![CDATA[SUCCESS]]></return_code>
  <return_msg><![CDATA[OK]]></return_msg>
  <result_code><![CDATA[SUCCESS]]></result_code>
  <prepay_id><![CDATA[wxprepays12345678901234567890]]></prepay_id>
  <code_url><![CDATA[weixin://wxpay/bizpayurl?pr=xxxx]]></code_url>
  <trade_type><![CDATA[NATIVE]]></trade_type>
</xml>`
		w.Write([]byte(response))
	}))
	defer server.Close()

	client := &WeChatAPIClient{
		appID:      "test_app_id",
		mchID:      "test_mch_id",
		apiKey:     "test_api_key",
		notifyURL:  "http://example.com/notify",
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	req := UnifiedOrderRequest{
		Body:           "测试商品",
		OutTradeNo:     "ORD20260209002",
		TotalFee:       100,
		SpbillCreateIP: "127.0.0.1",
		TradeType:      "NATIVE",
		NonceStr:       generateNonceStr(),
	}

	resp, err := client.unifiedOrderWithServer(server.URL, req)
	if err != nil {
		t.Fatalf("UnifiedOrder NATIVE failed: %v", err)
	}

	if resp.CodeURL == "" {
		t.Error("code_url should not be empty for NATIVE payment")
	}
	if !strings.HasPrefix(resp.CodeURL, "weixin://wxpay/bizpayurl") {
		t.Errorf("expected code_url to start with 'weixin://wxpay/bizpayurl', got '%s'", resp.CodeURL)
	}
}

// TestWeChatAPIClient_doRequest_ErrorHandling 测试错误处理
func TestWeChatAPIClient_doRequest_ErrorHandling(t *testing.T) {
	// 创建一个会返回错误的测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 返回错误响应
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		response := `<?xml version="1.0" encoding="UTF-8"?>
<xml>
  <return_code><![CDATA[FAIL]]></return_code>
  <return_msg><![CDATA[签名错误]]></return_msg>
</xml>`
		w.Write([]byte(response))
	}))
	defer server.Close()

	client := &WeChatAPIClient{
		appID:      "test_app_id",
		mchID:      "test_mch_id",
		apiKey:     "test_api_key",
		notifyURL:  "http://example.com/notify",
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	req := UnifiedOrderRequest{
		Body:           "测试商品",
		OutTradeNo:     "ORD20260209003",
		TotalFee:       100,
		SpbillCreateIP: "127.0.0.1",
		TradeType:      "JSAPI",
		NonceStr:       generateNonceStr(),
	}

	resp, err := client.unifiedOrderWithServer(server.URL, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.ReturnCode != "FAIL" {
		t.Errorf("expected return_code 'FAIL', got '%s'", resp.ReturnCode)
	}
	if resp.ReturnMsg != "签名错误" {
		t.Errorf("expected return_msg '签名错误', got '%s'", resp.ReturnMsg)
	}
}

// TestCreateTempCerts 创建临时证书文件用于测试
func TestCreateTempCerts(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()

	// 创建临时证书文件（用于测试，非真实证书）
	certPath := filepath.Join(tempDir, "cert.pem")
	keyPath := filepath.Join(tempDir, "key.pem")

	// 创建假的证书文件（仅用于测试文件存在性）
	if err := os.WriteFile(certPath, []byte("fake cert"), 0644); err != nil {
		t.Fatalf("failed to create temp cert file: %v", err)
	}
	if err := os.WriteFile(keyPath, []byte("fake key"), 0644); err != nil {
		t.Fatalf("failed to create temp key file: %v", err)
	}

	// 验证文件存在
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		t.Error("cert file was not created")
	}
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		t.Error("key file was not created")
	}
}
