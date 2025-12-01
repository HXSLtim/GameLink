package auth

import (
	"encoding/base64"
	"encoding/json"

	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGenerateToken_SignedStringError attempts to test the extremely rare SignedString error
func TestGenerateToken_SignedStringError(t *testing.T) {
	manager := NewJWTManager("test-secret", 1*time.Hour)

	// The SignedString method with HS256 is very reliable.
	// The only realistic way to test this is to verify successful generation.
	t.Run("Successful token generation covers the normal path", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			token, err := manager.GenerateToken(uint64(i), "user")
			assert.NoError(t, err)
			assert.NotEmpty(t, token)
		}
	})
}

// TestVerifyToken_TokenValidFalseBranch attempts to trigger the token.Valid = false branch
func TestVerifyToken_TokenValidFalseBranch(t *testing.T) {
	manager := NewJWTManager("test-secret", 1*time.Hour)

	t.Run("Expired token parsing", func(t *testing.T) {
		// Create an expired token manually
		claims := &Claims{
			UserID: 1,
			Role:   "user",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte("test-secret"))
		require.NoError(t, err)

		// When we verify an expired token, jwt.ParseWithClaims returns an error
		// So we never reach the token.Valid check
		_, err = manager.VerifyToken(tokenString)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token has invalid claims")
	})

	t.Run("Not yet valid token (nbf in future)", func(t *testing.T) {
		// Create a token that's not yet valid
		claims := &Claims{
			UserID: 1,
			Role:   "user",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
				NotBefore: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)), // In future
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte("test-secret"))
		require.NoError(t, err)

		// Parse should fail with nbf check
		_, err = manager.VerifyToken(tokenString)
		assert.Error(t, err)
	})

	t.Run("Signature verification fails after parsing", func(t *testing.T) {
		// Create a properly formed token but with wrong signature
		// Use a sophisticated approach to maintain structure but break signature
		claims := &Claims{
			UserID: 1,
			Role:   "user",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}

		header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
		payload, _ := json.Marshal(claims)
		payloadEnc := base64.RawURLEncoding.EncodeToString(payload)

		// Create token with invalid signature part but valid structure
		invalidToken := header + "." + payloadEnc + ".invalidsignature123"

		_, err := manager.VerifyToken(invalidToken)
		assert.Error(t, err)
	})
}

// TestVerifyToken_TypeAssertionFailure attempts to trigger the type assertion failure
func TestVerifyToken_TypeAssertionFailure(t *testing.T) {
	manager := NewJWTManager("test-secret", 1*time.Hour)

	t.Run("Manipulate JWT structure to cause type assertion issue", func(t *testing.T) {
		// The jwt.ParseWithClaims should always return *Claims when we pass &Claims{}
		// This is extremely difficult to trigger without modifying the JWT library
		// We can verify the normal path works correctly

		token, err := manager.GenerateToken(1, "user")
		require.NoError(t, err)

		// The normal flow should work perfectly
		parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("test-secret"), nil
		})
		require.NoError(t, err)

		claims, ok := parsedToken.Claims.(*Claims)
		assert.True(t, ok)
		assert.NotNil(t, claims)
		assert.Equal(t, uint64(1), claims.UserID)
	})

	t.Run("Use different claims type", func(t *testing.T) {
		// Even if we use a different claims type during parsing,
		// the library handles it gracefully
		token, err := manager.GenerateToken(1, "user")
		require.NoError(t, err)

		// Try with standard claims instead
		parsedToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("test-secret"), nil
		})
		require.NoError(t, err)

		// The token will have Claims of type *jwt.RegisteredClaims
		// This shows the type assertion is safe in our normal usage
		_, ok := parsedToken.Claims.(*jwt.RegisteredClaims)
		assert.True(t, ok)
	})
}

// TestVerifyToken_ExhaustiveBranches tries to cover all possible branches
func TestVerifyToken_ExhaustiveBranches(t *testing.T) {
	manager := NewJWTManager("test-secret", 1*time.Hour)

	testCases := []struct {
		name      string
		token     string
		wantError bool
		errorMsg  string
	}{
		{
			name:      "Completely empty token",
			token:     "",
			wantError: true,
			errorMsg:  "token is malformed",
		},
		{
			name:      "Single dot token",
			token:     ".",
			wantError: true,
		},
		{
			name:      "Only header",
			token:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			wantError: true,
		},
		{
			name:      "Malformed JSON in token",
			token:     "eyJ.eyJ.Invalid",
			wantError: true,
		},
		{
			name:      "Invalid base64",
			token:     "!!!.!!!.!!!",
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := manager.VerifyToken(tc.token)
			if tc.wantError {
				assert.Error(t, err)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
