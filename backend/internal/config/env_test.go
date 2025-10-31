package config

import (
	"os"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("APP_ENV", "unit-test")
	cfg := Load()
	if cfg.Port == "" {
		t.Fatal("expected default port to be set")
	}
	if cfg.Database.Type == "" {
		t.Fatal("expected database type to be set")
	}
}

func TestValidateProductionRequiresDSN(t *testing.T) {
	cfg := AppConfig{}
	if err := Validate("production", cfg); err == nil {
		t.Fatal("expected validation error when DSN missing")
	}
	cfg.Database.DSN = "postgres://example"
	if err := Validate("production", cfg); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestValidateCrypto(t *testing.T) {
	tests := []struct {
		name      string
		cfg       AppConfig
		expectErr bool
	}{
		{
			name: "Crypto disabled - no validation",
			cfg: AppConfig{
				Crypto: CryptoConfig{Enabled: false},
			},
			expectErr: false,
		},
		{
			name: "Crypto enabled with valid 16-byte key",
			cfg: AppConfig{
				Crypto: CryptoConfig{
					Enabled:   true,
					SecretKey: "1234567890123456",
					IV:        "1234567890123456",
					Methods:   []string{"POST"},
				},
			},
			expectErr: false,
		},
		{
			name: "Crypto enabled with valid 24-byte key",
			cfg: AppConfig{
				Crypto: CryptoConfig{
					Enabled:   true,
					SecretKey: "123456789012345678901234",
					IV:        "1234567890123456",
					Methods:   []string{"POST"},
				},
			},
			expectErr: false,
		},
		{
			name: "Crypto enabled with valid 32-byte key",
			cfg: AppConfig{
				Crypto: CryptoConfig{
					Enabled:   true,
					SecretKey: "12345678901234567890123456789012",
					IV:        "1234567890123456",
					Methods:   []string{"POST"},
				},
			},
			expectErr: false,
		},
		{
			name: "Crypto enabled with invalid key length",
			cfg: AppConfig{
				Crypto: CryptoConfig{
					Enabled:   true,
					SecretKey: "123",
					IV:        "1234567890123456",
					Methods:   []string{"POST"},
				},
			},
			expectErr: true,
		},
		{
			name: "Crypto enabled with short IV",
			cfg: AppConfig{
				Crypto: CryptoConfig{
					Enabled:   true,
					SecretKey: "1234567890123456",
					IV:        "123",
					Methods:   []string{"POST"},
				},
			},
			expectErr: true,
		},
		{
			name: "Crypto enabled with no methods",
			cfg: AppConfig{
				Crypto: CryptoConfig{
					Enabled:   true,
					SecretKey: "1234567890123456",
					IV:        "1234567890123456",
					Methods:   []string{},
				},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate("development", tt.cfg)
			if tt.expectErr && err == nil {
				t.Error("Expected validation error")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected validation error: %v", err)
			}
		})
	}
}

func TestNormalizeDBType(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"sqlite", "sqlite"},
		{"SQLITE", "sqlite"},
		{"  SQLite  ", "sqlite"},
		{"postgres", "postgres"},
		{"mysql", "mysql"},
		{"sqlserver", "sqlserver"},
		{"unknown", "sqlite"}, // defaults to sqlite
		{"", "sqlite"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := normalizeDBType(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeDBType(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNormalizeHTTPMethods(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "Normal methods",
			input:    []string{"post", "put", "patch"},
			expected: []string{"POST", "PUT", "PATCH"},
		},
		{
			name:     "Mixed case",
			input:    []string{"Post", "PUT", "pAtCh"},
			expected: []string{"POST", "PUT", "PATCH"},
		},
		{
			name:     "With whitespace",
			input:    []string{" post ", "  put", "patch  "},
			expected: []string{"POST", "PUT", "PATCH"},
		},
		{
			name:     "Empty strings filtered",
			input:    []string{"post", "", "put", "  ", "patch"},
			expected: []string{"POST", "PUT", "PATCH"},
		},
		{
			name:     "All empty returns default",
			input:    []string{"", "  ", ""},
			expected: []string{"POST", "PUT", "PATCH"},
		},
		{
			name:     "Empty input returns default",
			input:    []string{},
			expected: []string{"POST", "PUT", "PATCH"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeHTTPMethods(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("Length mismatch: got %d, want %d", len(result), len(tt.expected))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("Index %d: got %q, want %q", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestNormalizePaths(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "Normal paths",
			input:    []string{"/api/v1/health", "/api/v1/ping"},
			expected: []string{"/api/v1/health", "/api/v1/ping"},
		},
		{
			name:     "With whitespace",
			input:    []string{" /api/v1/health ", "  /api/v1/ping"},
			expected: []string{"/api/v1/health", "/api/v1/ping"},
		},
		{
			name:     "Empty strings filtered",
			input:    []string{"/api/v1/health", "", "/api/v1/ping", "  "},
			expected: []string{"/api/v1/health", "/api/v1/ping"},
		},
		{
			name:     "All empty returns empty",
			input:    []string{"", "  ", ""},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizePaths(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("Length mismatch: got %d, want %d", len(result), len(tt.expected))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("Index %d: got %q, want %q", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestOverrideFromEnv(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		validate func(*testing.T, *AppConfig)
	}{
		{
			name: "Override port",
			envVars: map[string]string{
				"SERVICE_PORT": "9090",
			},
			validate: func(t *testing.T, cfg *AppConfig) {
				if cfg.Port != "9090" {
					t.Errorf("Port = %q, want 9090", cfg.Port)
				}
			},
		},
		{
			name: "Override swagger enabled",
			envVars: map[string]string{
				"ENABLE_SWAGGER": "false",
			},
			validate: func(t *testing.T, cfg *AppConfig) {
				if cfg.EnableSwagger != false {
					t.Error("EnableSwagger should be false")
				}
			},
		},
		{
			name: "Override database config",
			envVars: map[string]string{
				"DB_TYPE": "postgres",
				"DB_DSN":  "postgres://test",
			},
			validate: func(t *testing.T, cfg *AppConfig) {
				if cfg.Database.Type != "postgres" {
					t.Errorf("Database.Type = %q, want postgres", cfg.Database.Type)
				}
				if cfg.Database.DSN != "postgres://test" {
					t.Errorf("Database.DSN = %q, want postgres://test", cfg.Database.DSN)
				}
			},
		},
		{
			name: "Override cache config",
			envVars: map[string]string{
				"CACHE_TYPE":     "redis",
				"REDIS_ADDR":     "localhost:6379",
				"REDIS_PASSWORD": "secret",
				"REDIS_DB":       "5",
			},
			validate: func(t *testing.T, cfg *AppConfig) {
				if cfg.Cache.Type != "redis" {
					t.Errorf("Cache.Type = %q, want redis", cfg.Cache.Type)
				}
				if cfg.Cache.Redis.Addr != "localhost:6379" {
					t.Errorf("Cache.Redis.Addr = %q, want localhost:6379", cfg.Cache.Redis.Addr)
				}
				if cfg.Cache.Redis.Password != "secret" {
					t.Errorf("Cache.Redis.Password = %q, want secret", cfg.Cache.Redis.Password)
				}
				if cfg.Cache.Redis.DB != 5 {
					t.Errorf("Cache.Redis.DB = %d, want 5", cfg.Cache.Redis.DB)
				}
			},
		},
		{
			name: "Override crypto config",
			envVars: map[string]string{
				"CRYPTO_ENABLED":       "true",
				"CRYPTO_SECRET_KEY":    "1234567890123456",
				"CRYPTO_IV":            "1234567890123456",
				"CRYPTO_METHODS":       "POST,GET",
				"CRYPTO_EXCLUDE_PATHS": "/health,/ping",
				"CRYPTO_USE_SIGNATURE": "false",
			},
			validate: func(t *testing.T, cfg *AppConfig) {
				if !cfg.Crypto.Enabled {
					t.Error("Crypto.Enabled should be true")
				}
				if cfg.Crypto.SecretKey != "1234567890123456" {
					t.Errorf("Crypto.SecretKey = %q", cfg.Crypto.SecretKey)
				}
				if cfg.Crypto.IV != "1234567890123456" {
					t.Errorf("Crypto.IV = %q", cfg.Crypto.IV)
				}
				if len(cfg.Crypto.Methods) != 2 {
					t.Errorf("Crypto.Methods length = %d, want 2", len(cfg.Crypto.Methods))
				}
				if len(cfg.Crypto.ExcludePaths) != 2 {
					t.Errorf("Crypto.ExcludePaths length = %d, want 2", len(cfg.Crypto.ExcludePaths))
				}
				if cfg.Crypto.UseSignature {
					t.Error("Crypto.UseSignature should be false")
				}
			},
		},
		{
			name: "Override auth config",
			envVars: map[string]string{
				"JWT_SECRET_KEY":      "test-secret",
				"JWT_TOKEN_TTL_HOURS": "48",
			},
			validate: func(t *testing.T, cfg *AppConfig) {
				if cfg.Auth.JWTSecret != "test-secret" {
					t.Errorf("Auth.JWTSecret = %q, want test-secret", cfg.Auth.JWTSecret)
				}
				if cfg.Auth.TokenTTLHours != 48 {
					t.Errorf("Auth.TokenTTLHours = %d, want 48", cfg.Auth.TokenTTLHours)
				}
			},
		},
		{
			name: "Override seed config",
			envVars: map[string]string{
				"SEED_ENABLED": "true",
			},
			validate: func(t *testing.T, cfg *AppConfig) {
				if !cfg.Seed.Enabled {
					t.Error("Seed.Enabled should be true")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set env vars
			for k, v := range tt.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			// Create default config
			cfg := &AppConfig{
				Port:          "8080",
				EnableSwagger: true,
				Database: DatabaseConfig{
					Type: "sqlite",
					DSN:  "",
				},
				Cache: CacheConfig{
					Type: "memory",
					Redis: RedisConfig{
						Addr:     "127.0.0.1:6379",
						Password: "",
						DB:       0,
					},
				},
				Crypto: CryptoConfig{
					Enabled:      false,
					SecretKey:    "",
					IV:           "",
					Methods:      []string{"POST"},
					ExcludePaths: []string{},
					UseSignature: true,
				},
				Auth: AuthConfig{
					JWTSecret:     "",
					TokenTTLHours: 24,
				},
				Seed: SeedConfig{Enabled: false},
			}

			// Apply env overrides
			overrideFromEnv(cfg)

			// Validate
			tt.validate(t, cfg)
		})
	}
}
