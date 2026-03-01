package config

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"unicode"
)

// Default constants for safe configuration values
const (
	// Server defaults
	DefaultPort          = "8081"
	DefaultEnableSwagger = true

	// Database defaults
	DefaultDBType = "sqlite"

	// Cache defaults
	DefaultCacheType = "memory"
	DefaultRedisAddr = "127.0.0.1:6379"
	DefaultRedisDB   = 0

	// Auth defaults
	DefaultTokenTTL = 24 // hours

	// Crypto defaults
	DefaultCryptoEnabled      = false
	DefaultCryptoUseSignature = true
	DefaultCryptoMethods      = "POST,PUT,PATCH"
	DefaultCryptoExcludePaths = "/api/v1/health,/api/v1/ping,/api/v1/auth/refresh"

	// Admin defaults
	DefaultAdminName     = "Super Admin"
	DefaultAdminAuthMode = "admin"

	// Security constants
	MinJWTSecretLength               = 16
	MinAdminPasswordLength           = 6
	ProductionMinAdminPasswordLength = 8
	MinCryptoIVLength                = 16

	// Allowed crypto key lengths
	CryptoKeyLength16 = 16
	CryptoKeyLength24 = 24
	CryptoKeyLength32 = 32
)

// SampleDSNs returns sample DSNs for different database types
func SampleDSNs() map[string]string {
	return map[string]string{
		"sqlite":   "gamelink.db",
		"postgres": "host=localhost user=gamelink password=gamelink dbname=gamelink sslmode=disable",
		"mysql":    "gamelink:gamelink@tcp(127.0.0.1:3306)/gamelink?charset=utf8mb4&parseTime=True&loc=Local",
	}
}

// GetSampleDSN returns sample DSN for the given database type
func GetSampleDSN(dbType string) string {
	samples := SampleDSNs()
	if dsn, ok := samples[dbType]; ok {
		return dsn
	}
	return ""
}

// AllowedHTTPMethods returns allowed HTTP methods
func AllowedHTTPMethods() []string {
	return []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
}

// AllowedAdminAuthModes returns allowed admin authentication modes
func AllowedAdminAuthModes() []string {
	return []string{"admin", "jwt"}
}

// AllowedCacheTypes returns allowed cache types
func AllowedCacheTypes() []string {
	return []string{"memory", "redis"}
}

// AllowedDatabaseTypes returns allowed database types
func AllowedDatabaseTypes() []string {
	return []string{"sqlite", "postgres", "mysql"}
}

// IsValidHTTPMethod checks if the method is a valid HTTP method
func IsValidHTTPMethod(method string) bool {
	methods := AllowedHTTPMethods()
	for _, m := range methods {
		if m == method {
			return true
		}
	}
	return false
}

// IsValidAdminAuthMode checks if the auth mode is valid
func IsValidAdminAuthMode(mode string) bool {
	modes := AllowedAdminAuthModes()
	for _, m := range modes {
		if m == mode {
			return true
		}
	}
	return false
}

// IsValidCacheType checks if the cache type is valid
func IsValidCacheType(cacheType string) bool {
	types := AllowedCacheTypes()
	for _, t := range types {
		if t == cacheType {
			return true
		}
	}
	return false
}

// IsValidDatabaseType checks if the database type is valid
func IsValidDatabaseType(dbType string) bool {
	types := AllowedDatabaseTypes()
	for _, t := range types {
		if t == dbType {
			return true
		}
	}
	return false
}

// IsValidCryptoKeyLength checks if the crypto key length is valid
func IsValidCryptoKeyLength(length int) bool {
	return length == CryptoKeyLength16 || length == CryptoKeyLength24 || length == CryptoKeyLength32
}

// ValidateCryptoKey validates crypto key length
func ValidateCryptoKey(key string) error {
	keyLen := len(key)
	if !IsValidCryptoKeyLength(keyLen) {
		return fmt.Errorf("crypto key must be %d, %d, or %d bytes (current: %d)",
			CryptoKeyLength16, CryptoKeyLength24, CryptoKeyLength32, keyLen)
	}
	return nil
}

// ValidateIV validates initialization vector length
func ValidateIV(iv string) error {
	if len(iv) < MinCryptoIVLength {
		return fmt.Errorf("IV must be at least %d bytes (current: %d)", MinCryptoIVLength, len(iv))
	}
	return nil
}

// ValidateJWTSecret validates JWT secret length
func ValidateJWTSecret(secret string) error {
	secret = strings.TrimSpace(secret)
	if len(secret) < MinJWTSecretLength {
		return fmt.Errorf("JWT secret is too short (length: %d, minimum: %d)", len(secret), MinJWTSecretLength)
	}
	return nil
}

// ValidateAdminPassword validates admin password strength
func ValidateAdminPassword(password string, env string) error {
	if password == "" {
		return errors.New("password cannot be empty")
	}

	minLength := MinAdminPasswordLength
	if env == "production" {
		minLength = ProductionMinAdminPasswordLength
	}

	if len(password) < minLength {
		return fmt.Errorf("password must be at least %d characters", minLength)
	}

	// 生产环境要求更强的密码
	if env == "production" {
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

// ValidateEmail validates email format
func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}
	if !strings.Contains(email, "@") {
		return errors.New("email must contain @ symbol")
	}
	if !strings.Contains(email, ".") {
		return errors.New("email must contain domain part")
	}
	return nil
}

// ValidatePort validates port number
func ValidatePort(port string) error {
	if port == "" {
		return errors.New("port cannot be empty")
	}
	// 简单的端口验证，可以扩展为更严格的验证
	if len(port) > 5 {
		return errors.New("port number is too long")
	}
	return nil
}

// ValidateTokenTTL validates token TTL
func ValidateTokenTTL(ttl int) error {
	if ttl <= 0 {
		return errors.New("token TTL must be positive")
	}
	if ttl > 8760 { // 一年小时数
		return errors.New("token TTL cannot exceed one year")
	}
	return nil
}

// ValidateCryptoMethods validates crypto methods
func ValidateCryptoMethods(methods []string) error {
	if len(methods) == 0 {
		return errors.New("crypto methods cannot be empty when encryption is enabled")
	}

	for _, method := range methods {
		if !IsValidHTTPMethod(method) {
			return fmt.Errorf("invalid HTTP method: %s", method)
		}
	}
	return nil
}

// ValidateExcludePaths validates exclude paths
func ValidateExcludePaths(paths []string) error {
	for _, path := range paths {
		if path == "" {
			return errors.New("exclude path cannot be empty")
		}
		if !strings.HasPrefix(path, "/") {
			return fmt.Errorf("exclude path must start with /: %s", path)
		}
	}
	return nil
}

// IsProduction checks if the environment is production
func IsProduction(env string) bool {
	return env == "production"
}

// IsDevelopment checks if the environment is development
func IsDevelopment(env string) bool {
	return env == "development" || env == "dev" || env == ""
}

// IsTest checks if the environment is test
func IsTest(env string) bool {
	return env == "test" || env == "testing"
}

// GetEnvironment returns the current environment
func GetEnvironment() string {
	env := os.Getenv("APP_ENV")
	if env == "" {
		return "development"
	}
	return env
}

// IsSecureEnvironment checks if the environment requires security validation
func IsSecureEnvironment(env string) bool {
	return IsProduction(env) || env == "staging"
}

// ValidateEnvironment validates environment string
func ValidateEnvironment(env string) error {
	validEnvs := []string{"development", "dev", "test", "testing", "staging", "production", "prod"}
	for _, valid := range validEnvs {
		if env == valid {
			return nil
		}
	}
	return fmt.Errorf("invalid environment: %s", env)
}

// GetDeprecatedDefaultJWTSecret returns the deprecated default JWT secret for validation
func GetDeprecatedDefaultJWTSecret() string {
	return deprecatedDefaultJWTSecret
}

// IsDeprecatedJWTSecret checks if the JWT secret is the deprecated default
func IsDeprecatedJWTSecret(secret string) bool {
	return strings.TrimSpace(secret) == deprecatedDefaultJWTSecret
}

// ValidateRedisConfig validates Redis configuration
func ValidateRedisConfig(cfg RedisConfig) error {
	if cfg.Addr == "" {
		return errors.New("redis address cannot be empty")
	}
	if cfg.DB < 0 {
		return errors.New("redis database number cannot be negative")
	}
	return nil
}

// ValidateCacheConfig validates cache configuration
func ValidateCacheConfig(cfg CacheConfig) error {
	if !IsValidCacheType(cfg.Type) {
		return fmt.Errorf("invalid cache type: %s", cfg.Type)
	}

	if cfg.Type == "redis" {
		if err := ValidateRedisConfig(cfg.Redis); err != nil {
			return fmt.Errorf("redis configuration error: %v", err)
		}
	}

	return nil
}

// ValidateDatabaseConfig validates database configuration
func ValidateDatabaseConfig(cfg DatabaseConfig) error {
	if !IsValidDatabaseType(cfg.Type) {
		return fmt.Errorf("invalid database type: %s", cfg.Type)
	}

	if cfg.DSN == "" {
		return errors.New("database DSN cannot be empty")
	}

	return nil
}

// ValidateAuthConfig validates authentication configuration
func ValidateAuthConfig(cfg AuthConfig) error {
	if err := ValidateJWTSecret(cfg.JWTSecret); err != nil {
		return fmt.Errorf("JWT secret validation failed: %v", err)
	}

	if err := ValidateTokenTTL(cfg.TokenTTLHours); err != nil {
		return fmt.Errorf("token TTL validation failed: %v", err)
	}

	if IsDeprecatedJWTSecret(cfg.JWTSecret) {
		return fmt.Errorf("JWT secret cannot be the deprecated default value")
	}

	return nil
}

// ValidateCryptoConfig validates crypto configuration
func ValidateCryptoConfig(cfg CryptoConfig) error {
	if !cfg.Enabled {
		return nil
	}

	if err := ValidateCryptoKey(cfg.SecretKey); err != nil {
		return fmt.Errorf("crypto key validation failed: %v", err)
	}

	if err := ValidateIV(cfg.IV); err != nil {
		return fmt.Errorf("IV validation failed: %v", err)
	}

	if err := ValidateCryptoMethods(cfg.Methods); err != nil {
		return fmt.Errorf("crypto methods validation failed: %v", err)
	}

	if err := ValidateExcludePaths(cfg.ExcludePaths); err != nil {
		return fmt.Errorf("exclude paths validation failed: %v", err)
	}

	// 检查是否为硬编码的默认值
	if cfg.SecretKey == "GameLink2025SecretKey!@#" {
		return errors.New("crypto secret key is using hardcoded default value, please set a secure key")
	}
	if cfg.IV == "GameLink2025IV!!!" {
		return errors.New("crypto IV is using hardcoded default value, please set a secure IV")
	}

	return nil
}

// ValidateSuperAdminConfig validates super admin configuration
func ValidateSuperAdminConfig(cfg SuperAdminConfig, env string) error {
	if cfg.Email == "" {
		return errors.New("super admin email cannot be empty")
	}

	if err := ValidateEmail(cfg.Email); err != nil {
		return fmt.Errorf("super admin email validation failed: %v", err)
	}

	if cfg.Password == "" {
		return errors.New("super admin password cannot be empty")
	}

	if err := ValidateAdminPassword(cfg.Password, env); err != nil {
		return fmt.Errorf("super admin password validation failed: %v", err)
	}

	return nil
}

// ValidateAdminAuthConfig validates admin authentication configuration
func ValidateAdminAuthConfig(cfg AdminAuthConfig) error {
	if !IsValidAdminAuthMode(cfg.Mode) {
		return fmt.Errorf("invalid admin auth mode: %s", cfg.Mode)
	}
	return nil
}

// ValidatePortConfig validates port configuration
func ValidatePortConfig(port string) error {
	if err := ValidatePort(port); err != nil {
		return fmt.Errorf("port validation failed: %v", err)
	}
	return nil
}

// ValidateSwaggerConfig validates swagger configuration
func ValidateSwaggerConfig(enabled bool) error {
	// Swagger配置目前没有太多验证规则
	return nil
}

// ValidateSeedConfig validates seed configuration
func ValidateSeedConfig(enabled bool) error {
	// Seed配置目前没有太多验证规则
	return nil
}

// NormalizeEnvironment normalizes environment string
func NormalizeEnvironment(env string) string {
	switch env {
	case "dev":
		return "development"
	case "prod":
		return "production"
	case "testing":
		return "test"
	default:
		return env
	}
}

// NormalizeDatabaseType normalizes database type
func NormalizeDatabaseType(dbType string) string {
	switch strings.ToLower(dbType) {
	case "postgresql", "pgsql":
		return "postgres"
	default:
		return strings.ToLower(dbType)
	}
}

// NormalizeCacheType normalizes cache type
func NormalizeCacheType(cacheType string) string {
	return strings.ToLower(cacheType)
}

// NormalizeHTTPMethods normalizes HTTP methods
func NormalizeHTTPMethods(methods []string) []string {
	var normalized []string
	for _, method := range methods {
		trimmed := strings.TrimSpace(strings.ToUpper(method))
		if trimmed == "" || !IsValidHTTPMethod(trimmed) {
			continue
		}
		normalized = append(normalized, trimmed)
	}

	if len(normalized) == 0 {
		return []string{"POST", "PUT", "PATCH"}
	}
	return normalized
}

// NormalizePaths normalizes paths
func NormalizePaths(paths []string) []string {
	var normalized []string
	for _, path := range paths {
		trimmed := strings.TrimSpace(path)
		if trimmed == "" {
			continue
		}
		normalized = append(normalized, trimmed)
	}
	return normalized
}

// NormalizeAdminAuthMode normalizes admin auth mode
func NormalizeAdminAuthMode(mode string) string {
	return strings.ToLower(mode)
}

// GetConfigValidationRules returns all validation rules for documentation
func GetConfigValidationRules() map[string]interface{} {
	return map[string]interface{}{
		"environment": map[string]interface{}{
			"allowed": []string{"development", "test", "staging", "production"},
			"default": "development",
		},
		"port": map[string]interface{}{
			"required": true,
			"format":   "string",
			"default":  DefaultPort,
		},
		"database": map[string]interface{}{
			"type": map[string]interface{}{
				"allowed": AllowedDatabaseTypes(),
				"default": DefaultDBType,
			},
			"dsn": map[string]interface{}{
				"required": true,
				"format":   "string",
			},
		},
		"auth": map[string]interface{}{
			"jwt_secret": map[string]interface{}{
				"required":    true,
				"min_length":  MinJWTSecretLength,
				"not_allowed": []string{GetDeprecatedDefaultJWTSecret()},
			},
			"token_ttl_hours": map[string]interface{}{
				"default": DefaultTokenTTL,
				"min":     1,
				"max":     8760,
			},
		},
		"crypto": map[string]interface{}{
			"enabled": map[string]interface{}{
				"default": DefaultCryptoEnabled,
			},
			"secret_key": map[string]interface{}{
				"when_enabled":    true,
				"allowed_lengths": []int{CryptoKeyLength16, CryptoKeyLength24, CryptoKeyLength32},
				"not_allowed":     []string{"GameLink2025SecretKey!@#"},
			},
			"iv": map[string]interface{}{
				"when_enabled": true,
				"min_length":   MinCryptoIVLength,
				"not_allowed":  []string{"GameLink2025IV!!!"},
			},
			"methods": map[string]interface{}{
				"when_enabled": true,
				"allowed":      AllowedHTTPMethods(),
				"default":      DefaultCryptoMethods,
			},
		},
		"super_admin": map[string]interface{}{
			"email": map[string]interface{}{
				"required": true,
				"format":   "email",
			},
			"password": map[string]interface{}{
				"required": true,
				"min_length": map[string]interface{}{
					"development": MinAdminPasswordLength,
					"production":  ProductionMinAdminPasswordLength,
				},
			},
		},
		"admin_auth": map[string]interface{}{
			"mode": map[string]interface{}{
				"allowed": AllowedAdminAuthModes(),
				"default": DefaultAdminAuthMode,
			},
		},
	}
}

// PrintValidationRules prints validation rules for debugging
func PrintValidationRules() {
	rules := GetConfigValidationRules()
	fmt.Println("Configuration Validation Rules:")
	printRules(rules, 0)
}

func printRules(rules map[string]interface{}, indent int) {
	for key, value := range rules {
		for i := 0; i < indent; i++ {
			fmt.Print("  ")
		}
		fmt.Printf("%s: ", key)

		switch v := value.(type) {
		case map[string]interface{}:
			fmt.Println()
			printRules(v, indent+1)
		case []interface{}:
			fmt.Printf("%v\n", v)
		default:
			fmt.Printf("%v\n", v)
		}
	}
}

// ValidateAllConfig validates all configuration sections
func ValidateAllConfig(cfg AppConfig, env string) error {
	// 基本配置验证
	if err := ValidatePortConfig(cfg.Port); err != nil {
		return fmt.Errorf("port configuration error: %v", err)
	}

	if err := ValidateSwaggerConfig(cfg.EnableSwagger); err != nil {
		return fmt.Errorf("swagger configuration error: %v", err)
	}

	if err := ValidateSeedConfig(cfg.Seed.Enabled); err != nil {
		return fmt.Errorf("seed configuration error: %v", err)
	}

	// 数据库配置验证
	if err := ValidateDatabaseConfig(cfg.Database); err != nil {
		return fmt.Errorf("database configuration error: %v", err)
	}

	// 缓存配置验证
	if err := ValidateCacheConfig(cfg.Cache); err != nil {
		return fmt.Errorf("cache configuration error: %v", err)
	}

	// 认证配置验证
	if err := ValidateAuthConfig(cfg.Auth); err != nil {
		return fmt.Errorf("auth configuration error: %v", err)
	}

	// 加密配置验证
	if err := ValidateCryptoConfig(cfg.Crypto); err != nil {
		return fmt.Errorf("crypto configuration error: %v", err)
	}

	// 超级管理员配置验证
	if err := ValidateSuperAdminConfig(cfg.SuperAdmin, env); err != nil {
		return fmt.Errorf("super admin configuration error: %v", err)
	}

	// 管理员认证配置验证
	if err := ValidateAdminAuthConfig(cfg.AdminAuth); err != nil {
		return fmt.Errorf("admin auth configuration error: %v", err)
	}

	// 生产环境特殊验证
	if IsProduction(env) {
		if err := ValidateProductionConfig(cfg); err != nil {
			return fmt.Errorf("production configuration error: %v", err)
		}
	}

	return nil
}

// ValidateProductionConfig validates production-specific configuration
func ValidateProductionConfig(cfg AppConfig) error {
	// 数据库DSN必须配置
	if cfg.Database.DSN == "" {
		return errors.New("database DSN is required in production")
	}

	// JWT Secret必须配置且不能是默认值
	if strings.TrimSpace(cfg.Auth.JWTSecret) == "" {
		return errors.New("JWT secret is required in production")
	}

	// 超级管理员凭据必须配置
	if strings.TrimSpace(cfg.SuperAdmin.Email) == "" {
		return errors.New("super admin email is required in production")
	}
	if strings.TrimSpace(cfg.SuperAdmin.Password) == "" {
		return errors.New("super admin password is required in production")
	}

	return nil
}

// GetValidationSummary returns a summary of configuration validation
func GetValidationSummary(cfg AppConfig, env string) map[string]interface{} {
	summary := map[string]interface{}{
		"environment": env,
		"validation": map[string]interface{}{
			"database":    ValidateDatabaseConfig(cfg.Database) == nil,
			"cache":       ValidateCacheConfig(cfg.Cache) == nil,
			"auth":        ValidateAuthConfig(cfg.Auth) == nil,
			"crypto":      ValidateCryptoConfig(cfg.Crypto) == nil,
			"super_admin": ValidateSuperAdminConfig(cfg.SuperAdmin, env) == nil,
			"admin_auth":  ValidateAdminAuthConfig(cfg.AdminAuth) == nil,
			"port":        ValidatePortConfig(cfg.Port) == nil,
			"swagger":     ValidateSwaggerConfig(cfg.EnableSwagger) == nil,
			"seed":        ValidateSeedConfig(cfg.Seed.Enabled) == nil,
		},
	}

	if IsProduction(env) {
		summary["production_validation"] = ValidateProductionConfig(cfg) == nil
	}

	return summary
}

// PrintValidationSummary prints configuration validation summary
func PrintValidationSummary(cfg AppConfig, env string) {
	summary := GetValidationSummary(cfg, env)
	fmt.Printf("Configuration Validation Summary (Environment: %s):\n", env)

	if validation, ok := summary["validation"].(map[string]interface{}); ok {
		for key, valid := range validation {
			status := "✓"
			if !valid.(bool) {
				status = "✗"
			}
			fmt.Printf("  %s %s\n", status, key)
		}
	}

	if prodValid, ok := summary["production_validation"]; ok {
		status := "✓"
		if !prodValid.(bool) {
			status = "✗"
		}
		fmt.Printf("  %s production validation\n", status)
	}
}

// SanitizeConfig sanitizes configuration for logging (removes sensitive data)
func SanitizeConfig(cfg AppConfig) AppConfig {
	sanitized := cfg

	// 移除敏感信息
	if cfg.Auth.JWTSecret != "" {
		sanitized.Auth.JWTSecret = "[REDACTED]"
	}

	if cfg.Crypto.SecretKey != "" {
		sanitized.Crypto.SecretKey = "[REDACTED]"
	}

	if cfg.Crypto.IV != "" {
		sanitized.Crypto.IV = "[REDACTED]"
	}

	if cfg.SuperAdmin.Password != "" {
		sanitized.SuperAdmin.Password = "[REDACTED]"
	}

	if cfg.Cache.Redis.Password != "" {
		sanitized.Cache.Redis.Password = "[REDACTED]"
	}

	return sanitized
}

// GetConfigDescription returns a description of the configuration
func GetConfigDescription(cfg AppConfig, env string) string {
	return fmt.Sprintf("Configuration for %s environment:\n"+
		"  Server Port: %s\n"+
		"  Database Type: %s\n"+
		"  Cache Type: %s\n"+
		"  Crypto Enabled: %v\n"+
		"  Swagger Enabled: %v\n"+
		"  Seed Enabled: %v\n"+
		"  Admin Auth Mode: %s",
		env,
		cfg.Port,
		cfg.Database.Type,
		cfg.Cache.Type,
		cfg.Crypto.Enabled,
		cfg.EnableSwagger,
		cfg.Seed.Enabled,
		cfg.AdminAuth.Mode,
	)
}

// ValidateAndGetErrors validates configuration and returns all errors
func ValidateAndGetErrors(cfg AppConfig, env string) []error {
	var errors []error

	// 基本配置验证
	if err := ValidatePortConfig(cfg.Port); err != nil {
		errors = append(errors, fmt.Errorf("port: %v", err))
	}

	// 数据库配置验证
	if err := ValidateDatabaseConfig(cfg.Database); err != nil {
		errors = append(errors, fmt.Errorf("database: %v", err))
	}

	// 缓存配置验证
	if err := ValidateCacheConfig(cfg.Cache); err != nil {
		errors = append(errors, fmt.Errorf("cache: %v", err))
	}

	// 认证配置验证
	if err := ValidateAuthConfig(cfg.Auth); err != nil {
		errors = append(errors, fmt.Errorf("auth: %v", err))
	}

	// 加密配置验证
	if err := ValidateCryptoConfig(cfg.Crypto); err != nil {
		errors = append(errors, fmt.Errorf("crypto: %v", err))
	}

	// 超级管理员配置验证
	if err := ValidateSuperAdminConfig(cfg.SuperAdmin, env); err != nil {
		errors = append(errors, fmt.Errorf("super_admin: %v", err))
	}

	// 管理员认证配置验证
	if err := ValidateAdminAuthConfig(cfg.AdminAuth); err != nil {
		errors = append(errors, fmt.Errorf("admin_auth: %v", err))
	}

	// 生产环境特殊验证
	if IsProduction(env) {
		if err := ValidateProductionConfig(cfg); err != nil {
			errors = append(errors, fmt.Errorf("production: %v", err))
		}
	}

	return errors
}

// ValidateWithWarnings validates configuration and returns errors and warnings
func ValidateWithWarnings(cfg AppConfig, env string) (errors []error, warnings []string) {
	errors = ValidateAndGetErrors(cfg, env)

	// 添加警告
	if !IsProduction(env) {
		if cfg.Auth.JWTSecret == "" {
			warnings = append(warnings, "JWT secret is not set (acceptable in non-production environments)")
		}
		if cfg.SuperAdmin.Password != "" && len(cfg.SuperAdmin.Password) < ProductionMinAdminPasswordLength {
			warnings = append(warnings, "Super admin password is weak (acceptable in non-production environments)")
		}
		if cfg.Crypto.Enabled && (cfg.Crypto.SecretKey == "GameLink2025SecretKey!@#" || cfg.Crypto.IV == "GameLink2025IV!!!") {
			warnings = append(warnings, "Crypto configuration is using default values (acceptable in non-production environments)")
		}
	}

	if cfg.EnableSwagger && IsProduction(env) {
		warnings = append(warnings, "Swagger is enabled in production environment")
	}

	if cfg.Seed.Enabled && IsProduction(env) {
		warnings = append(warnings, "Seed data is enabled in production environment")
	}

	return errors, warnings
}

// ValidateAndPrint validates configuration and prints results
func ValidateAndPrint(cfg AppConfig, env string) bool {
	errors, warnings := ValidateWithWarnings(cfg, env)

	fmt.Printf("Configuration Validation Results (Environment: %s):\n", env)

	if len(errors) == 0 {
		fmt.Println("✓ All validations passed")
	} else {
		fmt.Printf("✗ Found %d validation errors:\n", len(errors))
		for i, err := range errors {
			fmt.Printf("  %d. %v\n", i+1, err)
		}
	}

	if len(warnings) > 0 {
		fmt.Printf("⚠ Found %d warnings:\n", len(warnings))
		for i, warning := range warnings {
			fmt.Printf("  %d. %s\n", i+1, warning)
		}
	}

	return len(errors) == 0
}

// ValidateWithContext validates configuration with context
func ValidateWithContext(ctx context.Context, cfg AppConfig, env string) error {
	// 添加超时控制
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	// 创建验证完成通道
	done := make(chan error, 1)

	go func() {
		done <- Validate(env, cfg)
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return fmt.Errorf("configuration validation timeout: %v", ctx.Err())
	}
}

// ValidateAsync validates configuration asynchronously
func ValidateAsync(cfg AppConfig, env string, callback func(error)) {
	go func() {
		err := Validate(env, cfg)
		if callback != nil {
			callback(err)
		}
	}()
}

// ValidateAndPanic validates configuration and panics if validation fails
func ValidateAndPanic(env string, cfg AppConfig) {
	if err := Validate(env, cfg); err != nil {
		panic(fmt.Sprintf("Configuration validation failed: %v", err))
	}
}

// ValidateOrDefault validates configuration and returns default values if invalid
func ValidateOrDefault(cfg AppConfig, env string) AppConfig {
	if err := Validate(env, cfg); err != nil {
		log.Printf("Configuration validation failed: %v, using default values", err)
		return Load()
	}
	return cfg
}

// ValidateWithCustomRules validates configuration with custom validation rules
func ValidateWithCustomRules(cfg AppConfig, env string, customRules func(AppConfig, string) error) error {
	// 先执行标准验证
	if err := Validate(env, cfg); err != nil {
		return err
	}

	// 执行自定义验证
	if customRules != nil {
		if err := customRules(cfg, env); err != nil {
			return fmt.Errorf("custom validation failed: %v", err)
		}
	}

	return nil
}

// ValidateAndFix validates configuration and attempts to fix common issues
func ValidateAndFix(cfg AppConfig, env string) (AppConfig, error) {
	fixed := cfg

	// 修复常见问题
	if cfg.Port == "" {
		fixed.Port = DefaultPort
		log.Printf("Fixed empty port, using default: %s", DefaultPort)
	}

	if cfg.Auth.TokenTTLHours <= 0 {
		fixed.Auth.TokenTTLHours = DefaultTokenTTL
		log.Printf("Fixed invalid token TTL, using default: %d hours", DefaultTokenTTL)
	}

	if cfg.AdminAuth.Mode == "" {
		fixed.AdminAuth.Mode = DefaultAdminAuthMode
		log.Printf("Fixed empty admin auth mode, using default: %s", DefaultAdminAuthMode)
	}

	if cfg.SuperAdmin.Name == "" {
		fixed.SuperAdmin.Name = DefaultAdminName
		log.Printf("Fixed empty admin name, using default: %s", DefaultAdminName)
	}

	// 验证修复后的配置
	if err := Validate(env, fixed); err != nil {
		return cfg, fmt.Errorf("configuration cannot be fixed: %v", err)
	}

	return fixed, nil
}

// ValidateAndSuggest validates configuration and provides suggestions
func ValidateAndSuggest(cfg AppConfig, env string) (error, []string) {
	errors, warnings := ValidateWithWarnings(cfg, env)

	var suggestions []string

	// 基于警告和错误提供建议
	if len(errors) > 0 {
		suggestions = append(suggestions, "Fix the validation errors above before proceeding")
	}

	if len(warnings) > 0 {
		suggestions = append(suggestions, "Consider addressing the warnings for better security and performance")
	}

	// 额外的建议
	if !cfg.Crypto.Enabled && IsProduction(env) {
		suggestions = append(suggestions, "Consider enabling crypto for better security in production")
	}

	if cfg.EnableSwagger && IsProduction(env) {
		suggestions = append(suggestions, "Consider disabling Swagger in production for security")
	}

	if cfg.Seed.Enabled && IsProduction(env) {
		suggestions = append(suggestions, "Consider disabling seed data in production")
	}

	if len(cfg.Auth.JWTSecret) < 32 {
		suggestions = append(suggestions, "Consider using a longer JWT secret (32+ characters) for better security")
	}

	var mainError error
	if len(errors) > 0 {
		mainError = fmt.Errorf("configuration validation failed with %d errors", len(errors))
	}

	return mainError, suggestions
}
