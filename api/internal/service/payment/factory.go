package payment

import (
	"fmt"
	"log/slog"

	"gamelink/internal/model"
	"gamelink/internal/service/external"
)

// ProviderFactory creates payment providers based on config
type ProviderFactory struct {
	config *external.Config
}

// NewProviderFactory creates provider factory
func NewProviderFactory(cfg *external.Config) *ProviderFactory {
	return &ProviderFactory{config: cfg}
}

// CreateProvider creates provider for payment method
func (f *ProviderFactory) CreateProvider(method model.PaymentMethod) (ProviderClient, error) {
	switch method {
	case model.PaymentMethodWeChat:
		if f.config.WeChatPay.Enabled {
			return NewWeChatProvider(f.config), nil
		}
		return wechatProvider{}, nil // Use mock
	case model.PaymentMethodAlipay:
		if f.config.Alipay.Enabled {
			provider, err := NewAlipayProvider(f.config)
			if err != nil {
				if isProductionLikeEnv() {
					return nil, fmt.Errorf("alipay provider initialization failed in production: %w", err)
				}
				slog.Error("failed to create alipay provider, falling back to mock", "error", err)
				return alipayProvider{}, nil
			}
			return provider, nil
		}
		return alipayProvider{}, nil // Use mock
	default:
		return genericProvider{}, nil
	}
}

// CreateProviders creates all providers based on config
func (f *ProviderFactory) CreateProviders() map[model.PaymentMethod]ProviderClient {
	providers := make(map[model.PaymentMethod]ProviderClient)

	// WeChat Pay
	if f.config.WeChatPay.Enabled {
		providers[model.PaymentMethodWeChat] = NewWeChatProvider(f.config)
	} else {
		providers[model.PaymentMethodWeChat] = wechatProvider{}
	}

	// Alipay
	if f.config.Alipay.Enabled {
		provider, err := NewAlipayProvider(f.config)
		if err != nil || provider == nil {
			initErr := err
			if initErr == nil {
				initErr = fmt.Errorf("provider is nil")
			}
			if isProductionLikeEnv() {
				slog.Error("alipay enabled but provider creation failed, using fail-closed provider",
					"error", initErr,
					"publicKeyPath", f.config.Alipay.PublicKeyPath,
					"privateKeyPath", f.config.Alipay.PrivateKeyPath,
				)
				providers[model.PaymentMethodAlipay] = failClosedProvider{
					method: model.PaymentMethodAlipay,
					cause:  initErr,
				}
				return providers
			}

			slog.Error("alipay enabled but provider creation failed, falling back to mock provider",
				"error", initErr,
				"publicKeyPath", f.config.Alipay.PublicKeyPath,
				"privateKeyPath", f.config.Alipay.PrivateKeyPath,
			)
			providers[model.PaymentMethodAlipay] = alipayProvider{}
		} else {
			providers[model.PaymentMethodAlipay] = provider
		}
	} else {
		providers[model.PaymentMethodAlipay] = alipayProvider{}
	}

	return providers
}
