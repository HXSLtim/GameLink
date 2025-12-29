package sms

import (
	"context"
	"errors"
	"fmt"

	"gamelink/internal/service/external"
)

// Provider defines SMS provider interface
type Provider interface {
	SendCode(ctx context.Context, phone, code string) error
	SendNotification(ctx context.Context, phone, templateID string, params map[string]string) error
}

// Service handles SMS operations
type Service struct {
	config *external.Config
}

// NewService creates SMS service
func NewService(cfg *external.Config) *Service {
	return &Service{config: cfg}
}

var (
	// ErrSMSDisabled SMS service is disabled
	ErrSMSDisabled = errors.New("sms service is disabled")
	// ErrUnsupportedProvider unsupported SMS provider
	ErrUnsupportedProvider = errors.New("unsupported sms provider")
)

// SendCode sends verification code
func (s *Service) SendCode(ctx context.Context, phone, code string) error {
	if !s.config.SMS.Enabled {
		return ErrSMSDisabled
	}

	provider, err := s.getProvider()
	if err != nil {
		return err
	}

	return provider.SendCode(ctx, phone, code)
}

// SendNotification sends templated notification
func (s *Service) SendNotification(ctx context.Context, phone, templateID string, params map[string]string) error {
	if !s.config.SMS.Enabled {
		return ErrSMSDisabled
	}

	provider, err := s.getProvider()
	if err != nil {
		return err
	}

	return provider.SendNotification(ctx, phone, templateID, params)
}

func (s *Service) getProvider() (Provider, error) {
	switch s.config.SMS.Provider {
	case "tencent":
		return &TencentProvider{
			SecretID:  s.config.SMS.AccessKey,
			SecretKey: s.config.SMS.SecretKey,
			SignName:  s.config.SMS.SignName,
		}, nil
	case "aliyun":
		return &AliyunProvider{
			AccessKey: s.config.SMS.AccessKey,
			SecretKey: s.config.SMS.SecretKey,
			SignName:  s.config.SMS.SignName,
		}, nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedProvider, s.config.SMS.Provider)
	}
}
