package config

import (
	"strings"
	"testing"
)

func TestValidateCryptoConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     AppConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "加密禁用时不验证",
			cfg: AppConfig{
				Crypto: CryptoConfig{
					Enabled:   false,
					SecretKey: "",
					IV:        "",
				},
			},
			wantErr: false,
		},
		{
			name: "加密启用时密钥为空应报错",
			cfg: AppConfig{
				Crypto: CryptoConfig{
					Enabled:   true,
					SecretKey: "",
					IV:        "",
				},
			},
			wantErr: true,
			errMsg:  "CRYPTO_SECRET_KEY is required",
		},
		{
			name: "加密启用时IV为空应报错",
			cfg: AppConfig{
				Crypto: CryptoConfig{
					Enabled:   true,
					SecretKey: "12345678901234567890123456789012", // 32 bytes
					IV:        "",
				},
			},
			wantErr: true,
			errMsg:  "CRYPTO_IV is required",
		},
		{
			name: "密钥长度不正确应报错",
			cfg: AppConfig{
				Crypto: CryptoConfig{
					Enabled:   true,
					SecretKey: "short",            // 5 bytes
					IV:        "1234567890123456", // 16 bytes
				},
			},
			wantErr: true,
			errMsg:  "must be 16, 24 or 32 bytes",
		},
		{
			name: "使用硬编码默认密钥应报错",
			cfg: AppConfig{
				Crypto: CryptoConfig{
					Enabled:   true,
					SecretKey: "GameLink2025SecretKey!@#", // 24 bytes but hardcoded
					IV:        "1234567890123456",
				},
			},
			wantErr: true,
			errMsg:  "hardcoded default value",
		},
		{
			name: "IV使用硬编码默认值应报错",
			cfg: AppConfig{
				Crypto: CryptoConfig{
					Enabled:   true,
					SecretKey: "12345678901234567890123456789012",
					IV:        "GameLink2025IV!!!",
				},
			},
			wantErr: true,
			errMsg:  "hardcoded default value",
		},
		{
			name: "IV长度不足应报错",
			cfg: AppConfig{
				Crypto: CryptoConfig{
					Enabled:   true,
					SecretKey: "12345678901234567890123456789012",
					IV:        "short",
				},
			},
			wantErr: true,
			errMsg:  "must be at least 16 bytes",
		},
		{
			name: "有效配置应通过验证（16字节密钥）",
			cfg: AppConfig{
				Crypto: CryptoConfig{
					Enabled:   true,
					SecretKey: "1234567890123456", // 16 bytes
					IV:        "1234567890123456", // 16 bytes
					Methods:   []string{"POST", "PUT"},
				},
			},
			wantErr: false,
		},
		{
			name: "有效配置应通过验证（24字节密钥）",
			cfg: AppConfig{
				Crypto: CryptoConfig{
					Enabled:   true,
					SecretKey: "123456789012345678901234", // 24 bytes
					IV:        "1234567890123456",         // 16 bytes
					Methods:   []string{"POST", "PUT"},
				},
			},
			wantErr: false,
		},
		{
			name: "有效配置应通过验证（32字节密钥）",
			cfg: AppConfig{
				Crypto: CryptoConfig{
					Enabled:   true,
					SecretKey: "12345678901234567890123456789012", // 32 bytes
					IV:        "1234567890123456",                 // 16 bytes
					Methods:   []string{"POST", "PUT"},
				},
			},
			wantErr: false,
		},
		{
			name: "密钥包含空白字符应视为空",
			cfg: AppConfig{
				Crypto: CryptoConfig{
					Enabled:   true,
					SecretKey: "   ", // only spaces
					IV:        "1234567890123456",
				},
			},
			wantErr: true,
			errMsg:  "CRYPTO_SECRET_KEY is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCryptoConfig(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateCryptoConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("validateCryptoConfig() error = %v, want error containing %q", err, tt.errMsg)
				}
			}
		})
	}
}

func TestValidatePasswordStrength(t *testing.T) {
	tests := []struct {
		name    string
		pass    string
		env     string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "开发环境弱密码可通过",
			pass:    "123456",
			env:     "development",
			wantErr: false,
		},
		{
			name:    "密码少于6字符应报错",
			pass:    "12345",
			env:     "development",
			wantErr: true,
			errMsg:  "at least 6 characters",
		},
		{
			name:    "生产环境少于8字符应报错",
			pass:    "pass123",
			env:     "production",
			wantErr: true,
			errMsg:  "at least 8 characters",
		},
		{
			name:    "生产环境弱密码应报错",
			pass:    "12345678",
			env:     "production",
			wantErr: true,
			errMsg:  "must contain uppercase",
		},
		{
			name:    "生产环境强密码应通过",
			pass:    "NNLeRYZN1IF3A/T80C7+Q6mU3xBZtdnu",
			env:     "production",
			wantErr: false,
		},
		{
			name:    "生产环境复杂密码应通过",
			pass:    "SecurePass123!@#",
			env:     "production",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePasswordStrength(tt.pass, tt.env)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePasswordStrength() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("validatePasswordStrength() error = %v, want error containing %q", err, tt.errMsg)
				}
			}
		})
	}
}

func TestValidateProductionConfig(t *testing.T) {
	tests := []struct {
		name    string
		env     string
		cfg     AppConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "非生产环境不验证",
			env:  "development",
			cfg: AppConfig{
				Database: DatabaseConfig{DSN: ""},
			},
			wantErr: false,
		},
		{
			name: "生产环境DSN为空应报错",
			env:  "production",
			cfg: AppConfig{
				Database: DatabaseConfig{DSN: ""},
			},
			wantErr: true,
			errMsg:  "DB_DSN is required",
		},
		{
			name: "生产环境JWT密钥为空应报错",
			env:  "production",
			cfg: AppConfig{
				Database: DatabaseConfig{DSN: "postgres://localhost"},
				Auth:     AuthConfig{JWTSecret: ""},
			},
			wantErr: true,
			errMsg:  "JWT_SECRET_KEY is required",
		},
		{
			name: "生产环境超级管理员密码为空应报错",
			env:  "production",
			cfg: AppConfig{
				Database:   DatabaseConfig{DSN: "postgres://localhost"},
				Auth:       AuthConfig{JWTSecret: "secret123456789012345678"},
				SuperAdmin: SuperAdminConfig{Email: "admin@test.com", Password: ""},
			},
			wantErr: true,
			errMsg:  "SUPER_ADMIN_PASSWORD is required",
		},
		{
			name: "生产环境超级管理员密码弱应报错",
			env:  "production",
			cfg: AppConfig{
				Database:   DatabaseConfig{DSN: "postgres://localhost"},
				Auth:       AuthConfig{JWTSecret: "secret123456789012345678"},
				SuperAdmin: SuperAdminConfig{Email: "admin@test.com", Password: "weak"},
			},
			wantErr: true,
			errMsg:  "at least 8 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateProductionConfig(tt.env, tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateProductionConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("validateProductionConfig() error = %v, want error containing %q", err, tt.errMsg)
				}
			}
		})
	}
}

func TestValidateJWTSecret(t *testing.T) {
	tests := []struct {
		name    string
		secret  string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "空密钥不报错",
			secret:  "",
			wantErr: false,
		},
		{
			name:    "短密钥应报错",
			secret:  "short",
			wantErr: true,
			errMsg:  "too short",
		},
		{
			name:    "使用废弃默认密钥应报错",
			secret:  deprecatedDefaultJWTSecret,
			wantErr: true,
			errMsg:  "deprecated default value",
		},
		{
			name:    "有效密钥应通过",
			secret:  "MiRSQJJKEW2euVXKpvxRzjS1C5TCFlXx4RXGUXSdWpJ",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := AppConfig{Auth: AuthConfig{JWTSecret: tt.secret}}
			err := validateJWTSecret(cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateJWTSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("validateJWTSecret() error = %v, want error containing %q", err, tt.errMsg)
				}
			}
		})
	}
}
