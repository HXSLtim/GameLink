package sms

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
)

var ErrSMSIntegrationNotImplemented = errors.New("sms integration is not implemented in production")

func isProductionEnv() bool {
	env := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	return env == "production" || env == "prod"
}

// AliyunProvider implements SMS provider for Aliyun
type AliyunProvider struct {
	AccessKey string
	SecretKey string
	SignName  string
}

// SendCode sends verification code via Aliyun SMS
func (p *AliyunProvider) SendCode(ctx context.Context, phone, code string) error {
	if isProductionEnv() {
		return fmt.Errorf("%w: aliyun SendCode", ErrSMSIntegrationNotImplemented)
	}

	// Aliyun SMS API will be implemented for production
	// Reference: https://help.aliyun.com/document_detail/101414.html
	fmt.Printf("[Aliyun SMS] SendCode: phone=%s, code=%s, sign=%s\n", phone, code, p.SignName)
	return nil
}

// SendNotification sends templated notification via Aliyun SMS
func (p *AliyunProvider) SendNotification(ctx context.Context, phone, templateID string, params map[string]string) error {
	if isProductionEnv() {
		return fmt.Errorf("%w: aliyun SendNotification", ErrSMSIntegrationNotImplemented)
	}

	// Aliyun SMS API will be implemented for production
	fmt.Printf("[Aliyun SMS] SendNotification: phone=%s, template=%s, params=%v, sign=%s\n",
		phone, templateID, params, p.SignName)
	return nil
}
