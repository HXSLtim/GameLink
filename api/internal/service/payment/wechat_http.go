package payment

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// WeChatAPIClient 微信支付 API HTTP 客户端
type WeChatAPIClient struct {
	appID      string
	mchID      string
	apiKey     string
	certPath   string
	certKey    string
	notifyURL  string
	httpClient *http.Client
}

// NewWeChatAPIClient 创建微信支付 API 客户端
func NewWeChatAPIClient(appID, mchID, apiKey, certPath, certKey, notifyURL string) *WeChatAPIClient {
	return &WeChatAPIClient{
		appID:     appID,
		mchID:     mchID,
		apiKey:    apiKey,
		certPath:  certPath,
		certKey:   certKey,
		notifyURL: notifyURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// UnifiedOrderResponse 统一下单响应
type UnifiedOrderResponse struct {
	ReturnCode  string `xml:"return_code"`
	ReturnMsg   string `xml:"return_msg"`
	AppID       string `xml:"appid,omitempty"`
	MchID       string `xml:"mch_id,omitempty"`
	NonceStr    string `xml:"nonce_str,omitempty"`
	Sign        string `xml:"sign,omitempty"`
	ResultCode  string `xml:"result_code,omitempty"`
	ErrCode     string `xml:"err_code,omitempty"`
	ErrCodeDes  string `xml:"err_code_des,omitempty"`
	TradeType   string `xml:"trade_type,omitempty"`
	PrepayID    string `xml:"prepay_id,omitempty"`
	CodeURL     string `xml:"code_url,omitempty"`
	MwebURL     string `xml:"mweb_url,omitempty"`
}

// OrderQueryResponse 订单查询响应
type OrderQueryResponse struct {
	ReturnCode    string `xml:"return_code"`
	ReturnMsg     string `xml:"return_msg"`
	AppID        string `xml:"appid,omitempty"`
	MchID        string `xml:"mch_id,omitempty`
	NonceStr     string `xml:"nonce_str,omitempty"`
	Sign         string `xml:"sign,omitempty"`
	ResultCode    string `xml:"result_code,omitempty"`
	ErrCode       string `xml:"err_code,omitempty"`
	ErrCodeDes    string `xml:"err_code_des,omitempty"`
	TradeState    string `xml:"trade_state,omitempty"`
	TransactionID string `xml:"transaction_id,omitempty"`
	OutTradeNo    string `xml:"out_trade_no,omitempty"`
	TotalFee      string `xml:"total_fee,omitempty"`
	TimeEnd       string `xml:"time_end,omitempty"`
}

// RefundResponse 退款响应
type RefundResponse struct {
	ReturnCode    string `xml:"return_code"`
	ReturnMsg     string `xml:"return_msg"`
	AppID        string `xml:"appid,omitempty"`
	MchID        string `xml:"mch_id,omitempty`
	NonceStr     string `xml:"nonce_str,omitempty"`
	Sign         string `xml:"sign,omitempty"`
	ResultCode    string `xml:"result_code,omitempty"`
	ErrCode       string `xml:"err_code,omitempty"`
	ErrCodeDes    string `xml:"err_code_des,omitempty"`
	TransactionID string `xml:"transaction_id,omitempty"`
	OutTradeNo    string `xml:"out_trade_no,omitempty"`
	OutRefundNo   string `xml:"out_refund_no,omitempty"`
	RefundID      string `xml:"refund_id,omitempty"`
	RefundChannel string `xml:"refund_channel,omitempty"`
	RefundStatus  string `xml:"refund_status,omitempty"`
}

// UnifiedOrder 统一下单
// 参考：https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_1
func (c *WeChatAPIClient) UnifiedOrder(req UnifiedOrderRequest) (*UnifiedOrderResponse, error) {
	// 构建请求参数
	type xmlRequest struct {
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
		"total_fee":        fmt.Sprintf("%d", req.TotalFee),
		"spbill_create_ip": req.SpbillCreateIP,
		"notify_url":       c.notifyURL,
		"trade_type":       req.TradeType,
	}
	if req.OpenID != "" {
		params["openid"] = req.OpenID
	}
	sign := signClient.GenerateSign(params)

	xmlReq := xmlRequest{
		AppID:           c.appID,
		MchID:           c.mchID,
		NonceStr:        req.NonceStr,
		Body:            req.Body,
		OutTradeNo:      req.OutTradeNo,
		TotalFee:        fmt.Sprintf("%d", req.TotalFee),
		SpbillCreateIP:  req.SpbillCreateIP,
		NotifyURL:       c.notifyURL,
		TradeType:       req.TradeType,
		OpenID:          req.OpenID,
		Sign:            sign,
	}

	// 编码 XML
	xmlData, err := xml.Marshal(xmlReq)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	// 发送请求
	url := "https://api.mch.weixin.qq.com/pay/unifiedorder"
	respBody, err := c.doRequest(url, xmlData, false)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	// 解析响应
	var resp UnifiedOrderResponse
	if err := xml.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &resp, nil
}

// OrderQuery 查询订单
// 参考：https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_2
func (c *WeChatAPIClient) OrderQuery(outTradeNo, transactionID string) (*OrderQueryResponse, error) {
	// 构建请求参数
	type xmlRequest struct {
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

	xmlReq := xmlRequest{
		AppID:         c.appID,
		MchID:         c.mchID,
		NonceStr:      nonceStr,
		OutTradeNo:    outTradeNo,
		TransactionID: transactionID,
		Sign:          sign,
	}

	// 编码 XML
	xmlData, err := xml.Marshal(xmlReq)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	// 发送请求
	url := "https://api.mch.weixin.qq.com/pay/orderquery"
	respBody, err := c.doRequest(url, xmlData, false)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	// 解析响应
	var resp OrderQueryResponse
	if err := xml.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &resp, nil
}

// Refund 申请退款
// 参考：https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_4
func (c *WeChatAPIClient) Refund(req RefundRequest) (*RefundResponse, error) {
	// 构建请求参数
	type xmlRequest struct {
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

	// 生成签名
	signClient := NewWeChatClient(c.appID, c.mchID, c.apiKey, c.notifyURL)
	params := map[string]string{
		"appid":         c.appID,
		"mch_id":        c.mchID,
		"nonce_str":     req.NonceStr,
		"out_trade_no":  req.OutTradeNo,
		"out_refund_no": req.OutRefundNo,
		"total_fee":     fmt.Sprintf("%d", req.TotalFee),
		"refund_fee":    fmt.Sprintf("%d", req.RefundFee),
	}
	sign := signClient.GenerateSign(params)

	xmlReq := xmlRequest{
		AppID:       c.appID,
		MchID:       c.mchID,
		NonceStr:    req.NonceStr,
		OutTradeNo:  req.OutTradeNo,
		OutRefundNo: req.OutRefundNo,
		TotalFee:    fmt.Sprintf("%d", req.TotalFee),
		RefundFee:   fmt.Sprintf("%d", req.RefundFee),
		Sign:        sign,
	}

	// 编码 XML
	xmlData, err := xml.Marshal(xmlReq)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	// 发送请求（退款需要双向证书）
	url := "https://api.mch.weixin.qq.com/secapi/pay/refund"
	respBody, err := c.doRequest(url, xmlData, true)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	// 解析响应
	var resp RefundResponse
	if err := xml.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &resp, nil
}

// doRequest 发送 HTTP 请求
func (c *WeChatAPIClient) doRequest(url string, xmlData []byte, useCert bool) ([]byte, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(string(xmlData)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/xml")

	client := c.httpClient
	if useCert && c.certPath != "" && c.certKey != "" {
		// 加载客户端证书（退款接口需要）
		cert, err := tls.LoadX509KeyPair(c.certPath, c.certKey)
		if err != nil {
			return nil, fmt.Errorf("load certificate: %w", err)
		}

		client = &http.Client{
			Timeout: 60 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					Certificates: []tls.Certificate{cert},
					// 生产环境应该验证证书
					InsecureSkipVerify: false,
				},
			},
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// UnifiedOrderRequest 统一下单请求
type UnifiedOrderRequest struct {
	AppID          string `xml:"appid"`                       // 应用 ID
	MchID          string `xml:"mch_id"`                      // 商户号
	NonceStr       string `xml:"nonce_str"`                   // 随机字符串
	Body           string `xml:"body"`                        // 商品描述
	OutTradeNo     string `xml:"out_trade_no"`                 // 商户订单号
	TotalFee       int64  `xml:"total_fee"`                    // 订单金额（分）
	SpbillCreateIP string `xml:"spbill_create_ip"`             // 终端 IP
	NotifyURL      string `xml:"notify_url"`                   // 通知 URL
	TradeType      string `xml:"trade_type"`                   // 交易类型
	OpenID         string `xml:"openid,omitempty"`             // 用户标识（JSAPI）
}

// RefundRequest 退款请求
type RefundRequest struct {
	AppID        string `xml:"appid"`        // 应用 ID
	MchID        string `xml:"mch_id"`       // 商户号
	NonceStr     string `xml:"nonce_str"`    // 随机字符串
	OutTradeNo   string `xml:"out_trade_no"`  // 商户订单号
	OutRefundNo  string `xml:"out_refund_no"` // 商户退款单号
	TotalFee     int64  `xml:"total_fee"`    // 订单金额（分）
	RefundFee    int64  `xml:"refund_fee"`   // 退款金额（分）
}
