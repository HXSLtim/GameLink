package config

import "errors"

// Validate checks configuration for required values in production.
func Validate(env string, cfg AppConfig) error {
	if env == "production" {
		if cfg.Database.DSN == "" {
			return errors.New("DB_DSN is required in production")
		}
	}
	if cfg.Crypto.Enabled {
		keyLen := len(cfg.Crypto.SecretKey)
		if keyLen != 16 && keyLen != 24 && keyLen != 32 {
			return errors.New("CRYPTO_SECRET_KEY must be 16, 24 or 32 bytes when encryption is enabled")
		}
		if len(cfg.Crypto.IV) < 16 {
			return errors.New("CRYPTO_IV must be at least 16 bytes when encryption is enabled")
		}
		if len(cfg.Crypto.Methods) == 0 {
			return errors.New("crypto methods must not be empty when encryption is enabled")
		}
	}
	return nil
}
