package config

import (
	"testing"
)

// TestValidateAllConfig_SuccessAndFailure 覆盖 ValidateAllConfig 的成功与失败分支。
func TestValidateAllConfig_SuccessAndFailure(t *testing.T) {
	validCfg := AppConfig{
		Port:          "8080",
		EnableSwagger: true,
		Database: DatabaseConfig{
			Type: "sqlite",
			DSN:  "file:test.db",
		},
		Cache: CacheConfig{
			Type: "memory",
			Redis: RedisConfig{
				Addr: "127.0.0.1:6379",
				DB:   0,
			},
		},
		Auth: AuthConfig{
			JWTSecret:     "1234567890ABCDEF",
			TokenTTLHours: 24,
		},
		Crypto: CryptoConfig{
			Enabled: false,
		},
		Seed: SeedConfig{Enabled: false},
		SuperAdmin: SuperAdminConfig{
			Email:    "admin@example.com",
			Password: "Admin#123", // development 环境下长度足够
			Name:     "Admin",
		},
		AdminAuth: AdminAuthConfig{Mode: "admin"},
	}

	if err := ValidateAllConfig(validCfg, "development"); err != nil {
		t.Fatalf("期望 valid 配置通过验证，实际错误: %v", err)
	}

	// 构造无效配置：数据库类型非法或 DSN 为空
	invalidCfg := validCfg
	invalidCfg.Database.Type = "invalid-db"
	invalidCfg.Database.DSN = ""
	if err := ValidateAllConfig(invalidCfg, "development"); err == nil {
		t.Fatalf("期望 invalid 配置验证失败，但返回 nil")
	}
}

// TestSanitizeConfig 确认敏感字段被脱敏。
func TestSanitizeConfig(t *testing.T) {
	cfg := AppConfig{
		Port: "8080",
		Auth: AuthConfig{
			JWTSecret: "secret-jwt",
		},
		Crypto: CryptoConfig{
			SecretKey: "crypto-key",
			IV:        "crypto-iv",
		},
		SuperAdmin: SuperAdminConfig{
			Password: "admin-password",
		},
		Cache: CacheConfig{
			Type:  "redis",
			Redis: RedisConfig{Password: "redis-pass"},
		},
	}

	sanitized := SanitizeConfig(cfg)

	if sanitized.Auth.JWTSecret != "[REDACTED]" {
		t.Errorf("JWTSecret 未被脱敏: %q", sanitized.Auth.JWTSecret)
	}
	if sanitized.Crypto.SecretKey != "[REDACTED]" || sanitized.Crypto.IV != "[REDACTED]" {
		t.Errorf("Crypto Secret/IV 未被脱敏: %q / %q", sanitized.Crypto.SecretKey, sanitized.Crypto.IV)
	}
	if sanitized.SuperAdmin.Password != "[REDACTED]" {
		t.Errorf("SuperAdmin.Password 未被脱敏: %q", sanitized.SuperAdmin.Password)
	}
	if sanitized.Cache.Redis.Password != "[REDACTED]" {
		t.Errorf("Redis 密码未被脱敏: %q", sanitized.Cache.Redis.Password)
	}

	// 非敏感字段应保持不变
	if sanitized.Port != cfg.Port {
		t.Errorf("非敏感字段 Port 被意外修改: %q -> %q", cfg.Port, sanitized.Port)
	}
}

// TestValidateWithWarnings_NonProduction 覆盖非生产环境下的 warning 逻辑。
func TestValidateWithWarnings_NonProduction(t *testing.T) {
	cfg := AppConfig{
		Port: "8080",
		Database: DatabaseConfig{
			Type: "sqlite",
			DSN:  "file:test.db",
		},
		Cache: CacheConfig{
			Type:  "memory",
			Redis: RedisConfig{Addr: "127.0.0.1:6379"},
		},
		// 故意留空 JWTSecret 以触发 warning
		Auth: AuthConfig{
			JWTSecret:     "",
			TokenTTLHours: 24,
		},
		Crypto: CryptoConfig{
			Enabled:   true,
			SecretKey: "GameLink2025SecretKey!@#", // 使用默认值触发 warning
			IV:        "GameLink2025IV!!!",
			Methods:   []string{"POST"},
			ExcludePaths: []string{
				"/api/v1/health",
			},
		},
		SuperAdmin: SuperAdminConfig{
			Email:    "admin@example.com",
			Password: "weak", // 在非生产环境下触发弱密码 warning
		},
		AdminAuth: AdminAuthConfig{Mode: "admin"},
	}

	errs, warns := ValidateWithWarnings(cfg, "development")
	if len(warns) == 0 {
		t.Fatalf("预期至少有一个 warning，但为 0，errors=%v", errs)
	}
}

// TestValidateAndFix 确认 ValidateAndFix 会填充常见缺省并返回可通过验证的配置。
func TestValidateAndFix(t *testing.T) {
	cfg := AppConfig{
		// 故意留空，待自动修复
		Port: "",
		Database: DatabaseConfig{
			Type: "sqlite",
			DSN:  "file:test.db",
		},
		Cache: CacheConfig{
			Type:  "memory",
			Redis: RedisConfig{Addr: "127.0.0.1:6379"},
		},
		Auth: AuthConfig{
			JWTSecret:     "1234567890ABCDEF", // 长度有效
			TokenTTLHours: 0,                  // 待修复
		},
		Crypto: CryptoConfig{Enabled: false},
		SuperAdmin: SuperAdminConfig{
			Name: "", // 待修复
		},
		AdminAuth: AdminAuthConfig{Mode: ""}, // 待修复
	}

	fixed, err := ValidateAndFix(cfg, "development")
	if err != nil {
		t.Fatalf("期望配置可被修复，但返回错误: %v", err)
	}

	if fixed.Port == "" || fixed.Port != DefaultPort {
		t.Errorf("端口未被修复为默认值，got=%q, want=%q", fixed.Port, DefaultPort)
	}
	if fixed.Auth.TokenTTLHours != DefaultTokenTTL {
		t.Errorf("TokenTTL 未被修复，got=%d, want=%d", fixed.Auth.TokenTTLHours, DefaultTokenTTL)
	}
	if fixed.AdminAuth.Mode != DefaultAdminAuthMode {
		t.Errorf("AdminAuth.Mode 未被修复，got=%q, want=%q", fixed.AdminAuth.Mode, DefaultAdminAuthMode)
	}
	if fixed.SuperAdmin.Name != DefaultAdminName {
		t.Errorf("SuperAdmin.Name 未被修复，got=%q, want=%q", fixed.SuperAdmin.Name, DefaultAdminName)
	}
}
