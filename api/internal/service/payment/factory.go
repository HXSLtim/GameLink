package payment

import (
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
			return NewAlipayProvider(f.config)
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
		provider, _ := NewAlipayProvider(f.config)
		providers[model.PaymentMethodAlipay] = provider
	} else {
		providers[model.PaymentMethodAlipay] = alipayProvider{}
	}

	return providers
}
