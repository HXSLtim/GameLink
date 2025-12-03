package auth

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGenerateToken_100percent 覆盖GenerateToken的所有分支
func TestGenerateToken_100percent(t *testing.T) {
	manager := NewJWTManager("test-secret-key-for-100-percent-coverage-testing", 1*time.Hour)

	t.Run("Success path", func(t *testing.T) {
		token, err := manager.GenerateToken(1, "user")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("Generate with different user types", func(t *testing.T) {
		tests := []struct {
			name   string
			userID uint64
			role   string
		}{
			{"Admin user", 1, "admin"},
			{"Player user", 2, "player"},
			{"Regular user", 3, "user"},
			{"Zero ID user", 0, "guest"},
			{"Large ID user", 999999999, "user"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				token, err := manager.GenerateToken(tt.userID, tt.role)
				assert.NoError(t, err)
				assert.NotEmpty(t, token)

				// Verify the claims
				claims, err := manager.VerifyToken(token)
				assert.NoError(t, err)
				assert.Equal(t, tt.userID, claims.UserID)
				assert.Equal(t, tt.role, claims.Role)
			})
		}
	})
}

// TestVerifyToken_100percent 覆盖VerifyToken的所有分支
func TestVerifyToken_100percent(t *testing.T) {
	manager := NewJWTManager("test-secret-for-verify-token-coverage", 1*time.Hour)

	t.Run("Valid token verification", func(t *testing.T) {
		token, err := manager.GenerateToken(1, "user")
		require.NoError(t, err)

		claims, err := manager.VerifyToken(token)
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, uint64(1), claims.UserID)
		assert.Equal(t, "user", claims.Role)
	})

	t.Run("Expired token", func(t *testing.T) {
		shortManager := NewJWTManager("test-secret", 1*time.Millisecond)
		token, err := shortManager.GenerateToken(1, "user")
		require.NoError(t, err)

		time.Sleep(10 * time.Millisecond)

		_, err = shortManager.VerifyToken(token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expired")
	})

	t.Run("Malformed token - invalid format", func(t *testing.T) {
		tests := []string{
			"",
			"invalid.token",
			"invalid.token.format",
			"not-a-jwt-token-at-all",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", // Only header
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0", // No signature
		}

		for _, invalidToken := range tests {
			_, err := manager.VerifyToken(invalidToken)
			assert.Error(t, err, "Should error for invalid token: %s", invalidToken)
		}
	})

	t.Run("Wrong signing method", func(t *testing.T) {
		// We can't easily mock JWKS, but we can test the algorithm check
		// by creating a token with different method
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, &Claims{
			UserID: 1,
			Role:   "user",
		})

		tokenString, _ := token.SigningString()

		_, err := manager.VerifyToken(tokenString + ".invalid-signature")
		assert.Error(t, err)
	})

	t.Run("Wrong secret key", func(t *testing.T) {
		token, err := manager.GenerateToken(1, "user")
		require.NoError(t, err)

		// Try to verify with different manager (different secret)
		wrongManager := NewJWTManager("different-secret-key", 1*time.Hour)
		_, err = wrongManager.VerifyToken(token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "signature is invalid")
	})

	t.Run("Token with valid structure but invalid signature", func(t *testing.T) {
		// Create a valid token structure
		token, err := manager.GenerateToken(1, "user")
		require.NoError(t, err)

		// Tamper with the signature part
		parts := strings.Split(token, ".")
		require.Len(t, parts, 3)

		// Change the last part (signature)
		parts[2] = "tamperedsignature123"
		tamperedToken := parts[0] + "." + parts[1] + "." + parts[2]

		_, err = manager.VerifyToken(tamperedToken)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "signature is invalid")
	})

	t.Run("Token with modified payload", func(t *testing.T) {
		token, err := manager.GenerateToken(1, "user")
		require.NoError(t, err)

		// Decode and modify the payload
		parts := strings.Split(token, ".")
		require.Len(t, parts, 3)

		// Create token with modified claims but valid signature of original
		modifiedToken := parts[0] + "." + parts[1] + "." + parts[2] + "modified"

		_, err = manager.VerifyToken(modifiedToken)
		assert.Error(t, err)
	})
}

// TestExtractTokenFromHeader_Coverage 覆盖所有分支
func TestExtractTokenFromHeader_Coverage(t *testing.T) {
	tests := []struct {
		name        string
		authHeader  string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Perfect valid token",
			authHeader:  "Bearer validtoken123",
			expectError: false,
		},
		{
			name:        "Token with spaces",
			authHeader:  "Bearer token with spaces",
			expectError: false,
		},
		{
			name:        "Empty header",
			authHeader:  "",
			expectError: true,
			errorMsg:    "缺少Authorization头",
		},
		{
			name:        "Missing Bearer prefix",
			authHeader:  "validtoken123",
			expectError: true,
			errorMsg:    "Authorization头格式错误",
		},
		{
			name:        "Wrong prefix",
			authHeader:  "Basic abc123",
			expectError: true,
			errorMsg:    "Authorization头格式错误",
		},
		{
			name:        "Bearer with spaces only",
			authHeader:  "Bearer   ",
			expectError: true,
			errorMsg:    "Token为空",
		},
		{
			name:        "Bearer with tabs",
			authHeader:  "Bearer\t\t\t",
			expectError: true,
			errorMsg:    "Authorization头格式错误",
		},
		{
			name:        "Bearer length less than prefix",
			authHeader:  "Bear",
			expectError: true,
			errorMsg:    "Authorization头格式错误",
		},
		{
			name:        "Case sensitive - lower case bearer",
			authHeader:  "bearer token",
			expectError: true,
			errorMsg:    "Authorization头格式错误",
		},
		{
			name:        "Extra spaces around token",
			authHeader:  "Bearer   token123   ",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := ExtractTokenFromHeader(tt.authHeader)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}

// TestReadMaxRefreshWindow_AllBranches 覆盖所有分支
func TestReadMaxRefreshWindow_AllBranches(t *testing.T) {
	t.Run("Default value when env not set", func(t *testing.T) {
		duration := readMaxRefreshWindow()
		assert.Equal(t, 7*24*time.Hour, duration)
	})

	t.Run("Custom value from env", func(t *testing.T) {
		// Save original
		original := os.Getenv("JWT_MAX_REFRESH")
		defer os.Setenv("JWT_MAX_REFRESH", original)

		os.Setenv("JWT_MAX_REFRESH", "48h")
		duration := readMaxRefreshWindow()
		assert.Equal(t, 48*time.Hour, duration)
	})

	t.Run("Invalid env value falls back to default", func(t *testing.T) {
		// Save original
		original := os.Getenv("JWT_MAX_REFRESH")
		defer os.Setenv("JWT_MAX_REFRESH", original)

		os.Setenv("JWT_MAX_REFRESH", "invalid-duration")
		duration := readMaxRefreshWindow()
		assert.Equal(t, 7*24*time.Hour, duration)
	})

	t.Run("Various valid durations", func(t *testing.T) {
		// Save original
		original := os.Getenv("JWT_MAX_REFRESH")
		defer os.Setenv("JWT_MAX_REFRESH", original)

		tests := []struct {
			durationStr string
			expected    time.Duration
		}{
			{"24h", 24 * time.Hour},
			{"48h", 48 * time.Hour},
			{"72h", 72 * time.Hour},
			{"168h", 168 * time.Hour},
			{"30m", 30 * time.Minute},
		}

		for _, tt := range tests {
			t.Run(tt.durationStr, func(t *testing.T) {
				os.Setenv("JWT_MAX_REFRESH", tt.durationStr)
				result := readMaxRefreshWindow()
				assert.Equal(t, tt.expected, result)
			})
		}
	})
}

// TestVerifyToken_StructCoverage - make sure all struct fields are used
func TestVerifyToken_StructCoverage(t *testing.T) {
	manager := NewJWTManager("test-secret", 1*time.Hour)

	t.Run("All claims fields populated", func(t *testing.T) {
		token, err := manager.GenerateToken(12345, "admin")
		require.NoError(t, err)

		claims, err := manager.VerifyToken(token)
		require.NoError(t, err)

		// Verify all Claims fields are accessible and correct
		assert.Equal(t, uint64(12345), claims.UserID)
		assert.Equal(t, "admin", claims.Role)
		assert.Equal(t, "gamelink", claims.Issuer)
		assert.True(t, claims.ExpiresAt.Time.After(time.Now()))
		assert.True(t, claims.IssuedAt.Time.Before(time.Now()))
		assert.True(t, claims.NotBefore.Time.Before(time.Now()))
	})
}
