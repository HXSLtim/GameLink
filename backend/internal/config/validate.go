package config

import "errors"

// Validate checks configuration for required values in production.
func Validate(env string, cfg AppConfig) error {
	if env == "production" {
		if cfg.Database.DSN == "" {
			return errors.New("DB_DSN is required in production")
		}
	}
	return nil
}
