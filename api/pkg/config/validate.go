package config

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

// Validate checks configuration for required values in production.
// This function is deprecated, use ValidateAllConfig instead for comprehensive validation.
func Validate(env string, cfg AppConfig) error {
	if err := validateProductionConfig(env, cfg); err != nil {
		return err
	}
	if err := validateJWTSecret(cfg); err != nil {
		return err
	}
	if err := validateCryptoConfig(cfg); err != nil {
		return err
	}
	if err := validateSuperAdminConfig(env, cfg); err != nil {
		return err
	}
	return nil
}

// validateProductionConfig 验证生产环境必需配置
func validateProductionConfig(env string, cfg AppConfig) error {
	if env != "production" {
		return nil
	}
	if cfg.Database.DSN == "" {
		return errors.New("DB_DSN is required in production")
	}
	if strings.TrimSpace(cfg.Auth.JWTSecret) == "" {
		return errors.New("JWT_SECRET_KEY is required in production")
	}
	if strings.TrimSpace(cfg.SuperAdmin.Email) == "" {
		return errors.New("SUPER_ADMIN_EMAIL is required in production")
	}
	if strings.TrimSpace(cfg.SuperAdmin.Password) == "" {
		return errors.New("SUPER_ADMIN_PASSWORD is required in production")
	}
	if len(cfg.SuperAdmin.Password) < 8 {
		return errors.New("SUPER_ADMIN_PASSWORD must be at least 8 characters in production")
	}
	// 生产环境必须启用加密
	if !cfg.Crypto.Enabled {
		return errors.New("CRYPTO_ENABLED must be true in production for security")
	}
	// 生产环境必须配置缓存为 Redis
	if cfg.Cache.Type != "redis" {
		return errors.New("CACHE_TYPE must be 'redis' in production (not 'memory')")
	}
	return nil
}

// validateJWTSecret 验证 JWT Secret 安全性
func validateJWTSecret(cfg AppConfig) error {
	jwtSecret := strings.TrimSpace(cfg.Auth.JWTSecret)
	if jwtSecret == "" {
		return nil
	}
	// Enforce minimum 32 characters for JWT secret (increased from 16)
	if len(jwtSecret) < 32 {
		return fmt.Errorf("JWT secret is too short (length: %d, minimum: 32 characters for production security)", len(jwtSecret))
	}
	if jwtSecret == deprecatedDefaultJWTSecret {
		return fmt.Errorf("JWT secret cannot be the deprecated default value '%s'", deprecatedDefaultJWTSecret)
	}
	return nil
}

// validateCryptoConfig 验证加密配置
func validateCryptoConfig(cfg AppConfig) error {
	if !cfg.Crypto.Enabled {
		return nil
	}
	// 检查密钥是否为空（P0安全漏洞）
	if strings.TrimSpace(cfg.Crypto.SecretKey) == "" {
		return errors.New("CRYPTO_SECRET_KEY is required when encryption is enabled")
	}
	if strings.TrimSpace(cfg.Crypto.IV) == "" {
		return errors.New("CRYPTO_IV is required when encryption is enabled")
	}
	// 先检查是否使用了已知的硬编码默认值（必须在长度检查之前）
	deprecatedKeys := []string{
		"H/oguKMv23lWlivgq8snNZmTzSUp6KSHZnEEo1c0Ook=",
		"hTeObHJQ3nGDNs4H4O778A==",
		"GameLink2025SecretKey!@#",
		"GameLink2025IV!!!",
	}
	for _, deprecatedKey := range deprecatedKeys {
		if cfg.Crypto.SecretKey == deprecatedKey {
			return errors.New("CRYPTO_SECRET_KEY is using a hardcoded default value, please set a secure key")
		}
		if cfg.Crypto.IV == deprecatedKey {
			return errors.New("CRYPTO_IV is using a hardcoded default value, please set a secure IV")
		}
	}
	// 验证密钥长度（AES支持16、24、32字节，分别对应AES-128、AES-192、AES-256）
	// 推荐使用32字节以获得最佳安全性（AES-256-CBC）
	keyLen := len(cfg.Crypto.SecretKey)
	if keyLen != 16 && keyLen != 24 && keyLen != 32 {
		return fmt.Errorf("CRYPTO_SECRET_KEY must be 16, 24 or 32 bytes for AES (current: %d). Generate with: openssl rand -base64 32", keyLen)
	}
	if len(cfg.Crypto.IV) < 16 {
		return fmt.Errorf("CRYPTO_IV must be at least 16 bytes for AES-CBC (current: %d). Generate with: openssl rand -base64 16", len(cfg.Crypto.IV))
	}
	if len(cfg.Crypto.Methods) == 0 {
		return errors.New("CRYPTO_METHODS must not be empty when encryption is enabled")
	}
	return nil
}

// validateSuperAdminConfig 验证超级管理员配置
func validateSuperAdminConfig(env string, cfg AppConfig) error {
	if cfg.SuperAdmin.Password != "" {
		if err := validatePasswordStrength(cfg.SuperAdmin.Password, env); err != nil {
			return fmt.Errorf("SUPER_ADMIN_PASSWORD does not meet security requirements: %v", err)
		}
	}
	if cfg.SuperAdmin.Email != "" && !strings.Contains(cfg.SuperAdmin.Email, "@") {
		return errors.New("SUPER_ADMIN_EMAIL must be a valid email address")
	}
	return nil
}

// validatePasswordStrength 验证密码强度
func validatePasswordStrength(password string, env string) error {
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	// 生产环境要求更强的密码
	if env == "production" {
		if len(password) < 8 {
			return errors.New("password must be at least 8 characters in production")
		}

		hasUpper := false
		hasLower := false
		hasDigit := false
		hasSpecial := false

		for _, char := range password {
			switch {
			case unicode.IsUpper(char):
				hasUpper = true
			case unicode.IsLower(char):
				hasLower = true
			case unicode.IsDigit(char):
				hasDigit = true
			case unicode.IsPunct(char) || unicode.IsSymbol(char):
				hasSpecial = true
			}
		}

		if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
			return errors.New("password must contain uppercase, lowercase, digit, and special character in production")
		}
	}

	return nil
}
