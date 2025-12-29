package sms

import (
	"context"
	"fmt"
)

// AliyunProvider implements SMS provider for Aliyun
type AliyunProvider struct {
	AccessKey string
	SecretKey string
	SignName  string
}

// SendCode sends verification code via Aliyun SMS
func (p *AliyunProvider) SendCode(ctx context.Context, phone, code string) error {
	// TODO: Implement Aliyun SMS API
	// Reference: https://help.aliyun.com/document_detail/101414.html
	fmt.Printf("[Aliyun SMS] SendCode: phone=%s, code=%s, sign=%s\n", phone, code, p.SignName)
	return nil
}

// SendNotification sends templated notification via Aliyun SMS
func (p *AliyunProvider) SendNotification(ctx context.Context, phone, templateID string, params map[string]string) error {
	// TODO: Implement Aliyun SMS API
	fmt.Printf("[Aliyun SMS] SendNotification: phone=%s, template=%s, params=%v, sign=%s\n",
		phone, templateID, params, p.SignName)
	return nil
}
