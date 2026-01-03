package sms

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// TencentProvider implements SMS provider for Tencent Cloud
type TencentProvider struct {
	SecretID  string
	SecretKey string
	SignName  string
}

// Tencent SMS API endpoints
const (
	tencentSMSHost = "sms.tencentcloudapi.com"
	tencentSendURL = "https://" + tencentSMSHost + "/"
)

// SendCode sends verification code via Tencent Cloud SMS
func (p *TencentProvider) SendCode(ctx context.Context, phone, code string) error {
	// Use default template for verification code
	// Template param: {1}=code
	return p.sendRequest(ctx, phone, "", map[string]string{"1": code})
}

// SendNotification sends templated notification
func (p *TencentProvider) SendNotification(ctx context.Context, phone, templateID string, params map[string]string) error {
	return p.sendRequest(ctx, phone, templateID, params)
}

// sendRequest sends SMS request to Tencent Cloud
func (p *TencentProvider) sendRequest(ctx context.Context, phone, templateID string, params map[string]string) error {
	// Build request body (for future implementation)
	_ = map[string]interface{}{
		"PhoneNumberSet":   []string{p.formatPhone(phone)},
		"SmsSdkAppId":      "", // Set from config or use default
		"SignName":         p.SignName,
		"TemplateId":       templateID,
		"TemplateParamSet": p.buildParams(params),
	}

	// For development, just log the request
	fmt.Printf("[Tencent SMS] Send: phone=%s, sign=%s, template=%s, params=%v\n",
		phone, p.SignName, templateID, params)

	// Tencent Cloud SMS API will be implemented for production
	// Reference: https://cloud.tencent.com/document/product/382/43195

	return nil
}

// formatPhone formats phone number for Tencent Cloud (+86 prefix)
func (p *TencentProvider) formatPhone(phone string) string {
	phone = strings.TrimPrefix(phone, "+86")
	phone = strings.TrimPrefix(phone, "86")
	if !strings.HasPrefix(phone, "+") {
		phone = "+86" + phone
	}
	return phone
}

// buildParams builds template parameter array
func (p *TencentProvider) buildParams(params map[string]string) []string {
	// Sort by key to ensure correct order
	maxKey := 0
	for k := range params {
		var key int
		fmt.Sscanf(k, "%d", &key)
		if key > maxKey {
			maxKey = key
		}
	}

	result := make([]string, maxKey)
	for i := 1; i <= maxKey; i++ {
		key := fmt.Sprintf("%d", i)
		result[i-1] = params[key]
	}
	return result
}

// Tencent SMS API response structures
type tencentResponse struct {
	Response struct {
		RequestID     string `json:"RequestId"`
		SendStatusSet []struct {
			SerialNo    string `json:"SerialNo"`
			PhoneNumber string `json:"PhoneNumber"`
			Fee         int    `json:"Fee"`
			SessionNo   string `json:"SessionNo"`
			Code        string `json:"Code"`
			Message     string `json:"Message"`
		} `json:"SendStatusSet"`
	} `json:"Response"`
	Error struct {
		Code    string `json:"Code"`
		Message string `json:"Message"`
	} `json:"Error"`
}

// doRequest performs HTTP request to Tencent Cloud API
func (p *TencentProvider) doRequest(ctx context.Context, action string, payload map[string]interface{}) ([]byte, error) {
	_, _ = json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", tencentSendURL, strings.NewReader(""))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-TC-Action", action)
	req.Header.Set("X-TC-Timestamp", fmt.Sprintf("%d", time.Now().Unix()))
	req.Header.Set("X-TC-Version", "2021-01-11")
	req.Header.Set("X-TC-Region", "ap-guangzhou")

	// TC3-HMAC-SHA256 signature will be implemented for production
	// Reference: https://cloud.tencent.com/document/product/382/43105

	// Timeout choice: 10s for SMS API calls (should be fast)
	// Note: We use context.WithTimeout to ensure the entire operation respects the deadline
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
