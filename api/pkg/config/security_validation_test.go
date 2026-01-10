package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateJWTSecret_Length(t *testing.T) {
	tests := []struct {
		name        string
		secret      string
		expectError bool
	}{
		{
			name:        "valid secret - 32 characters",
			secret:      "this-is-a-very-secure-jwt-secret-32",
			expectError: false,
		},
		{
			name:        "valid secret - more than 32 characters",
			secret:      "this-is-an-extremely-secure-jwt-secret-key-for-production-use-123",
			expectError: false,
		},
		{
			name:        "too short - 31 characters",
			secret:      "short-jwt-secret-only-31-chars-",
			expectError: true,
		},
		{
			name:        "too short - 16 characters (old minimum)",
			secret:      "short-secret-16",
			expectError: true,
		},
		{
			name:        "empty secret (allowed for development)",
			secret:      "",
			expectError: false,
		},
		{
			name:        "whitespace only (allowed for development)",
			secret:      "   ",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := AppConfig{
				Auth: AuthConfig{
					JWTSecret: tt.secret,
				},
			}

			err := Validate("development", cfg)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "too short")
				assert.Contains(t, err.Error(), "32")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateJWTSecret_DeprecatedValue(t *testing.T) {
	cfg := AppConfig{
		Auth: AuthConfig{
			JWTSecret: deprecatedDefaultJWTSecret,
		},
	}

	err := Validate("development", cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "deprecated default value")
}

func TestValidateCryptoConfig_Security(t *testing.T) {
	tests := []struct {
		name        string
		enabled     bool
		secretKey   string
		iv          string
		expectError bool
	}{
		{
			name:        "disabled crypto",
			enabled:     false,
			secretKey:   "",
			iv:          "",
			expectError: false,
		},
		{
			name:        "valid crypto config - 32 bytes",
			enabled:     true,
			secretKey:   "H/oguKMv23lWlivgq8snNZmTzSUp6KSH",
			iv:          "16-byte-iv-12345",
			expectError: false,
		},
		{
			name:        "valid crypto config - 16 bytes",
			enabled:     true,
			secretKey:   "16-byte-key-1234",
			iv:          "16-byte-iv-12345",
			expectError: false,
		},
		{
			name:        "invalid length - 15 bytes",
			enabled:     true,
			secretKey:   "only-15-bytes-123",
			iv:          "16-byte-iv-12345",
			expectError: true,
		},
		{
			name:        "invalid length - 17 bytes",
			enabled:     true,
			secretKey:   "17-byte-key-123456",
			iv:          "16-byte-iv-12345",
			expectError: true,
		},
		{
			name:        "IV too short",
			enabled:     true,
			secretKey:   "16-byte-key-1234",
			iv:          "short",
			expectError: true,
		},
		{
			name:        "hardcoded default secret key",
			enabled:     true,
			secretKey:   "H/oguKMv23lWlivgq8snNZmTzSUp6KSHZnEEo1c0Ook=",
			iv:          "hTeObHJQ3nGDNs4H4O778A==",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := AppConfig{
				Crypto: CryptoConfig{
					Enabled:   tt.enabled,
					SecretKey: tt.secretKey,
					IV:        tt.iv,
					Methods:   []string{"POST", "PUT"},
				},
			}

			err := Validate("development", cfg)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

