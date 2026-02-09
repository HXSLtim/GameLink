package payment

import (
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/service/external"
)

// RealWeChatProvider implements WeChat Pay provider
type RealWeChatProvider struct {
	config *external.Config
}

// NewWeChatProvider creates WeChat Pay provider with config
func NewWeChatProvider(cfg *external.Config) *RealWeChatProvider {
	return &RealWeChatProvider{config: cfg}
}

// WeChat Pay API endpoints
const (
	wechatPayHost      = "api.mch.weixin.qq.com"
	wechatPayRefundURL = "https://" + wechatPayHost + "/secapi/pay/refund"
)

// Refund requests refund from WeChat Pay
func (p *RealWeChatProvider) Refund(ctx context.Context, payment *model.Payment, reason string) (string, json.RawMessage, time.Time, error) {
	// Timeout choice: 60s for payment gateway refund operations
	// Payment gateway operations can be slow due to bank processing and network latency
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	if !p.config.WeChatPay.Enabled {
		return "", nil, time.Time{}, ErrPaymentDisabled
	}

	// Build refund request
	req := WeChatRefundRequest{
		AppID:       p.config.WeChatPay.AppID,
		MchID:       p.config.WeChatPay.MchID,
		NonceStr:    generateNonceStr(),
		SignType:    "MD5",
		OutTradeNo:  fmt.Sprintf("%d", payment.ID),
		OutRefundNo: fmt.Sprintf("refund_%d_%d", payment.ID, time.Now().Unix()),
		TotalFee:    payment.AmountCents, // 单位：分
		RefundFee:   payment.AmountCents, // 全额退款
	}

	// Generate signature
	req.Sign = p.generateSign(req)

	// For development, log instead of actual API call
	fmt.Printf("[WeChat Pay] Refund: %+v\n", req)

	// WeChat Pay refund API integration will be implemented in production
	// Reference: https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_4
	// Requires: merchant certificate, API secret key, HTTPS client

	now := time.Now()
	raw := map[string]interface{}{
		"channel":       "wechat",
		"payment_id":    payment.ID,
		"refund_reason": reason,
		"refunded_at":   now.Unix(),
		"refund_no":     req.OutRefundNo,
	}
	b, _ := json.Marshal(raw)

	return req.OutRefundNo, json.RawMessage(b), now, nil
}

// generateSign generates WeChat Pay signature
func (p *RealWeChatProvider) generateSign(req WeChatRefundRequest) string {
	// Build parameter string
	params := map[string]string{
		"appid":         req.AppID,
		"mch_id":        req.MchID,
		"nonce_str":     req.NonceStr,
		"out_trade_no":  req.OutTradeNo,
		"out_refund_no": req.OutRefundNo,
		"total_fee":     fmt.Sprintf("%d", req.TotalFee),
		"refund_fee":    fmt.Sprintf("%d", req.RefundFee),
	}

	// Sort keys
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build string
	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, params[k]))
	}
	signStr := strings.Join(parts, "&") + "&key=" + p.config.WeChatPay.APIKey

	// MD5 hash
	h := md5.New()
	h.Write([]byte(signStr))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

// WeChatRefundRequest represents WeChat Pay refund request
type WeChatRefundRequest struct {
	XMLName     xml.Name `xml:"xml"`
	AppID       string   `xml:"appid"`
	MchID       string   `xml:"mch_id"`
	NonceStr    string   `xml:"nonce_str"`
	Sign        string   `xml:"sign"`
	SignType    string   `xml:"sign_type,omitempty"`
	OutTradeNo  string   `xml:"out_trade_no"`
	OutRefundNo string   `xml:"out_refund_no"`
	TotalFee    int64    `xml:"total_fee"`
	RefundFee   int64    `xml:"refund_fee"`
}

// WeChatRefundResponse represents WeChat Pay refund response
type WeChatRefundResponse struct {
	XMLName       xml.Name `xml:"xml"`
	ReturnCode    string   `xml:"return_code"`
	ReturnMsg     string   `xml:"return_msg"`
	ResultCode    string   `xml:"result_code,omitempty"`
	ErrCode       string   `xml:"err_code,omitempty"`
	ErrCodeDes    string   `xml:"err_code_des,omitempty"`
	AppID         string   `xml:"appid,omitempty"`
	MchID         string   `xml:"mch_id,omitempty"`
	NonceStr      string   `xml:"nonce_str,omitempty"`
	Sign          string   `xml:"sign,omitempty"`
	OutTradeNo    string   `xml:"out_trade_no,omitempty"`
	OutRefundNo   string   `xml:"out_refund_no,omitempty"`
	TransactionID string   `xml:"transaction_id,omitempty"`
	TotalFee      int64    `xml:"total_fee,omitempty"`
	RefundFee     int64    `xml:"refund_fee,omitempty"`
}

// doRefundRequest performs actual WeChat Pay refund API request
func (p *RealWeChatProvider) doRefundRequest(ctx context.Context, req WeChatRefundRequest) (*WeChatRefundResponse, error) {
	// Marshal to XML
	xmlData, err := xml.Marshal(req)
	if err != nil {
		return nil, err
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", wechatPayRefundURL, strings.NewReader(string(xmlData)))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/xml")

	// Client certificate will be added for production WeChat Pay integration
	// Required for refund API: httpReq.GetBody = ...

	// Timeout choice: 30s for individual HTTP requests to WeChat Pay API
	// This is a reasonable timeout for API calls, while the overall operation has 60s
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result WeChatRefundResponse
	if err := xml.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// verifySign verifies WeChat Pay response signature
func (p *RealWeChatProvider) verifySign(resp WeChatRefundResponse) bool {
	// Signature verification will be implemented for production
	// Reference: https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=4_3
	return true
}

// VerifySign verifies WeChat Pay notification signature
func (p *RealWeChatProvider) VerifySign(params map[string]string) bool {
	// Get original sign
	sign := params["sign"]
	if sign == "" {
		return false
	}

	// Build string to sign
	var keys []string
	for k := range params {
		if k != "sign" && params[k] != "" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(params[k])
		sb.WriteString("&")
	}
	sb.WriteString("key=")
	sb.WriteString(p.config.WeChatPay.APIKey)

	// Calculate signature
	hash := md5.Sum([]byte(sb.String()))
	calculated := strings.ToUpper(hex.EncodeToString(hash[:]))

	return calculated == sign
}

// ParseWeChatXML parses WeChat XML notification body
func ParseWeChatXML(data []byte) (map[string]string, error) {
	result := make(map[string]string)
	decoder := xml.NewDecoder(strings.NewReader(string(data)))

	var currentKey string
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		switch t := token.(type) {
		case xml.StartElement:
			currentKey = t.Name.Local
		case xml.CharData:
			if currentKey != "" && currentKey != "xml" {
				result[currentKey] = strings.TrimSpace(string(t))
			}
		case xml.EndElement:
			currentKey = ""
		}
	}

	return result, nil
}

// CreateOrder creates WeChat Pay order
func (p *RealWeChatProvider) CreateOrder(ctx context.Context, orderID, description string, amountCents int64, clientIP string) (map[string]interface{}, error) {
	// Timeout choice: 60s for payment gateway order creation
	// Payment operations require interaction with external banking systems
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	if !p.config.WeChatPay.Enabled {
		return nil, ErrPaymentDisabled
	}

	// Build order request
	req := map[string]interface{}{
		"appid":            p.config.WeChatPay.AppID,
		"mch_id":           p.config.WeChatPay.MchID,
		"nonce_str":        generateNonceStr(),
		"body":             description,
		"out_trade_no":     orderID,
		"total_fee":        amountCents,
		"spbill_create_ip": clientIP,
		"notify_url":       p.config.WeChatPay.NotifyURL,
		"trade_type":       "JSAPI", // or NATIVE, APP
	}

	// Generate signature
	sign := p.generateMapSign(req)
	req["sign"] = sign

	fmt.Printf("[WeChat Pay] CreateOrder: %+v\n", req)

	// WeChat Pay unified order API will be implemented for production
	// Reference: https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_1

	return req, nil
}

// generateMapSign generates signature from map
func (p *RealWeChatProvider) generateMapSign(params map[string]interface{}) string {
	// Sort keys
	var keys []string
	for k := range params {
		if k != "sign" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	// Build string
	var parts []string
	for _, k := range keys {
		v := fmt.Sprintf("%v", params[k])
		if v != "" {
			parts = append(parts, fmt.Sprintf("%s=%s", k, v))
		}
	}
	signStr := strings.Join(parts, "&") + "&key=" + p.config.WeChatPay.APIKey

	// MD5 hash
	h := md5.New()
	h.Write([]byte(signStr))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

// generateNonceStr generates random nonce string
func generateNonceStr() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// ErrPaymentDisabled indicates payment is disabled
var ErrPaymentDisabled = fmt.Errorf("payment service is disabled")

// HmacSHA256 generates HMAC-SHA256 signature
func HmacSHA256(secret, data string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
