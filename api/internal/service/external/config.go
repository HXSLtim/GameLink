package external

import (
	"gamelink/pkg/config"
)

// Config holds external API credentials
type Config struct {
	WeChatPay config.WeChatPayConfig
	Alipay    config.AlipayConfig
	SMS       config.SMSConfig
	OSS       config.OSSConfig
}

// NewConfig creates external API config from app config
func NewConfig(appCfg config.ExternalAPIConfig) *Config {
	return &Config{
		WeChatPay: appCfg.WeChatPay,
		Alipay:    appCfg.Alipay,
		SMS:       appCfg.SMS,
		OSS:       appCfg.OSS,
	}
}

// IsEnabled checks if any external API is enabled
func (c *Config) IsEnabled() bool {
	return c.WeChatPay.Enabled || c.Alipay.Enabled || c.SMS.Enabled || c.OSS.Enabled
}
