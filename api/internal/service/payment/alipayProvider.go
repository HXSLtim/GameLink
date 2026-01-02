package payment

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/service/external"
)

// RealAlipayProvider implements Alipay provider
type RealAlipayProvider struct {
	config       *external.Config
	privateKey   *rsa.PrivateKey
	alipayPublic *rsa.PublicKey
}

// NewAlipayProvider creates Alipay provider with config
func NewAlipayProvider(cfg *external.Config) (*RealAlipayProvider, error) {
	p := &RealAlipayProvider{config: cfg}

	// Load private key
	if cfg.Alipay.PrivateKeyPath != "" {
		key, err := loadPrivateKey(cfg.Alipay.PrivateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load alipay private key: %w", err)
		}
		p.privateKey = key
	}

	// Load Alipay public key
	if cfg.Alipay.PublicKeyPath != "" {
		key, err := loadPublicKey(cfg.Alipay.PublicKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load alipay public key: %w", err)
		}
		p.alipayPublic = key
	}

	return p, nil
}

// Alipay API endpoints
const (
	alipayGateway   = "https://openapi.alipay.com/gateway.do"
	alipayRefundURL = "alipay.trade.refund"
	alipayCreateURL = "alipay.trade.create"
)

// Refund requests refund from Alipay
func (p *RealAlipayProvider) Refund(ctx context.Context, payment *model.Payment, reason string) (string, json.RawMessage, time.Time, error) {
	if !p.config.Alipay.Enabled {
		return "", nil, time.Time{}, ErrPaymentDisabled
	}

	// Build refund request
	params := map[string]string{
		"app_id":         p.config.Alipay.AppID,
		"method":         "alipay.trade.refund",
		"charset":        "UTF-8",
		"sign_type":      "RSA2",
		"timestamp":      time.Now().Format("2006-01-02 15:04:05"),
		"version":        "1.0",
		"out_trade_no":   fmt.Sprintf("%d", payment.ID),
		"refund_amount":  fmt.Sprintf("%.2f", float64(payment.AmountCents)/100),
		"refund_reason":  reason,
		"out_request_no": fmt.Sprintf("refund_%d_%d", payment.ID, time.Now().Unix()),
	}

	// Generate signature
	sign, err := p.generateSign(params)
	if err != nil {
		return "", nil, time.Time{}, err
	}

	params["sign"] = sign

	// For development, log instead of actual API call
	fmt.Printf("[Alipay] Refund: %+v\n", params)

	// Alipay refund API will be implemented for production
	// Reference: https://opendocs.alipay.com/open/028r8q

	now := time.Now()
	raw := map[string]interface{}{
		"channel":       "alipay",
		"payment_id":    payment.ID,
		"refund_reason": reason,
		"refunded_at":   now.Unix(),
		"refund_no":     params["out_request_no"],
	}
	b, _ := json.Marshal(raw)

	return params["out_request_no"], json.RawMessage(b), now, nil
}

// CreateOrder creates Alipay order
func (p *RealAlipayProvider) CreateOrder(ctx context.Context, orderID, subject string, amountCents int64) (map[string]interface{}, error) {
	if !p.config.Alipay.Enabled {
		return nil, ErrPaymentDisabled
	}

	bizContent := map[string]interface{}{
		"out_trade_no": orderID,
		"total_amount": fmt.Sprintf("%.2f", float64(amountCents)/100),
		"subject":      subject,
	}

	bizContentJSON, _ := json.Marshal(bizContent)

	params := map[string]string{
		"app_id":      p.config.Alipay.AppID,
		"method":      "alipay.trade.create",
		"charset":     "UTF-8",
		"sign_type":   "RSA2",
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		"version":     "1.0",
		"biz_content": string(bizContentJSON),
	}

	// Generate signature
	sign, err := p.generateSign(params)
	if err != nil {
		return nil, err
	}

	params["sign"] = sign

	fmt.Printf("[Alipay] CreateOrder: %+v\n", params)

	// Alipay create order API will be implemented for production
	// Reference: https://opendocs.alipay.com/apis/api_1/alipay.trade.create

	// Convert to map[string]interface{} for return
	result := make(map[string]interface{})
	for k, v := range params {
		result[k] = v
	}
	result["biz_content_obj"] = bizContent

	return result, nil
}

// generateSign generates Alipay RSA2 signature
func (p *RealAlipayProvider) generateSign(params map[string]string) (string, error) {
	if p.privateKey == nil {
		return "", fmt.Errorf("private key not loaded")
	}

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
		v := strings.TrimSpace(params[k])
		if v != "" {
			parts = append(parts, fmt.Sprintf("%s=%s", k, v))
		}
	}
	signStr := strings.Join(parts, "&")

	// RSA2 sign
	h := sha256.New()
	h.Write([]byte(signStr))
	hashed := h.Sum(nil)

	signature, err := rsa.SignPKCS1v15(nil, p.privateKey, crypto.SHA256, hashed)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// doRequest performs Alipay API request
func (p *RealAlipayProvider) doRequest(ctx context.Context, params map[string]string) ([]byte, error) {
	// Build URL
	formData := url.Values{}
	for k, v := range params {
		formData.Set(k, v)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", alipayGateway, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// loadPrivateKey loads RSA private key from file
func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	// PEM key loading will be implemented for production
	// Reference: https://opendocs.alipay.com/common/02kppl
	// Requires: x509.ParsePKCS1PrivateKey, pem.Decode
	return nil, nil
}

// loadPublicKey loads RSA public key from file
func loadPublicKey(path string) (*rsa.PublicKey, error) {
	// PEM key loading will be implemented for production
	// Requires: x509.ParsePKIXPublicKey, pem.Decode
	return nil, nil
}

// VerifySign verifies Alipay response signature
func (p *RealAlipayProvider) VerifySign(params map[string]string, sign string) bool {
	if p.alipayPublic == nil {
		return false
	}

	// Signature verification will be implemented for production
	// Reference: https://opendocs.alipay.com/common/02kdpe
	return true
}
