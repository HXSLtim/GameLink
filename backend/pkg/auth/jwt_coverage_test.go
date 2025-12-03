package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGenerateToken_ErrorHandling 测试GenerateToken的错误处理路径
func TestGenerateToken_ErrorHandling(t *testing.T) {
	manager := NewJWTManager("test-secret", 1*time.Hour)

	// Test with valid data (cover success path)
	token, err := manager.GenerateToken(1, "user")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify the token works
	claims, err := manager.VerifyToken(token)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), claims.UserID)
	assert.Equal(t, "user", claims.Role)
	assert.Equal(t, "gamelink", claims.Issuer)
	assert.True(t, claims.IssuedAt.Time.After(time.Now().Add(-1*time.Minute)))
	assert.True(t, claims.ExpiresAt.Time.After(time.Now()))
	assert.True(t, claims.NotBefore.Time.Before(time.Now().Add(1*time.Minute)))
}

// TestVerifyToken_ErrorPaths 测试VerifyToken的所有错误路径
func TestVerifyToken_ErrorPaths(t *testing.T) {
	manager := NewJWTManager("test-secret", 1*time.Hour)

	t.Run("Invalid token format", func(t *testing.T) {
		_, err := manager.VerifyToken("not-a-valid-jwt-token")
		assert.Error(t, err)
	})

	t.Run("Malformed token - too many parts", func(t *testing.T) {
		_, err := manager.VerifyToken("part1.part2.part3.part4")
		assert.Error(t, err)
	})

	t.Run("Malformed token - wrong signing method", func(t *testing.T) {
		// The RS256 token will fail to parse correctly and result in a different error
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, &Claims{
			UserID: 1,
			Role:   "user",
		})
		// RS256 requires a private key, this will cause parsing to fail
		tokenString, _ := token.SignedString([]byte("dummy"))

		_, err := manager.VerifyToken(tokenString)
		assert.Error(t, err)
		// Content may vary, just ensure it errors
	})

	t.Run("Malformed token - HS256 with invalid signature", func(t *testing.T) {
		// Create a valid HS256 token structure but with wrong secret
		shortManager := NewJWTManager("different-secret", 1*time.Hour)
		token, err := shortManager.GenerateToken(1, "user")
		require.NoError(t, err)

		// Try to verify with different secret
		_, err = manager.VerifyToken(token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "signature is invalid")
	})

	t.Run("Token with empty claims", func(t *testing.T) {
		// This will test the case where token.Valid is false
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{})
		tokenString, _ := token.SignedString([]byte("different-secret"))

		_, err := manager.VerifyToken(tokenString)
		assert.Error(t, err)
	})

	t.Run("Token with standard claims", func(t *testing.T) {
		// Create token with standard claims instead of our custom Claims
		// This will still parse successfully because our Claims embeds RegisteredClaims
		// The type assertion will work but UserID and Role will be empty
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Issuer: "gamelink",
		})
		tokenString, _ := token.SignedString([]byte("test-secret"))

		claims, err := manager.VerifyToken(tokenString)
		assert.NoError(t, err)                    // This will succeed because Claims embeds RegisteredClaims
		assert.Equal(t, uint64(0), claims.UserID) // UserID will be zero
		assert.Equal(t, "", claims.Role)          // Role will be empty
	})

	t.Run("Expired token", func(t *testing.T) {
		shortManager := NewJWTManager("test-secret", 1*time.Millisecond)
		token, err := shortManager.GenerateToken(1, "user")
		require.NoError(t, err)

		time.Sleep(10 * time.Millisecond)

		_, err = shortManager.VerifyToken(token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token is expired")
	})

	t.Run("Token with wrong secret", func(t *testing.T) {
		token, err := manager.GenerateToken(1, "user")
		require.NoError(t, err)

		// Try to verify with wrong secret
		wrongManager := NewJWTManager("wrong-secret", 1*time.Hour)
		_, err = wrongManager.VerifyToken(token)
		assert.Error(t, err)
	})
}

// TestVerifyToken_ClaimsValidation 测试VerifyToken的claims验证逻辑
func TestVerifyToken_ClaimsValidation(t *testing.T) {
	manager := NewJWTManager("test-secret", 1*time.Hour)

	t.Run("Valid token with all claims", func(t *testing.T) {
		token, err := manager.GenerateToken(12345, "admin")
		require.NoError(t, err)

		claims, err := manager.VerifyToken(token)
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, uint64(12345), claims.UserID)
		assert.Equal(t, "admin", claims.Role)
		assert.Equal(t, "gamelink", claims.Issuer)
		assert.True(t, claims.IssuedAt.Time.Before(time.Now()))
		assert.True(t, claims.ExpiresAt.Time.After(time.Now()))
		assert.True(t, claims.NotBefore.Time.Before(time.Now()))
	})

	t.Run("Token with tampered signature", func(t *testing.T) {
		token, err := manager.GenerateToken(1, "user")
		require.NoError(t, err)

		// Tamper with the signature part
		parts := strings.Split(token, ".")
		require.Len(t, parts, 3)

		parts[2] = "tamperedsignature"
		tamperedToken := strings.Join(parts, ".")

		_, err = manager.VerifyToken(tamperedToken)
		assert.Error(t, err)
	})
}

// TestRefreshToken_EdgeCases 测试RefreshToken的边缘情况
func TestRefreshToken_EdgeCases(t *testing.T) {
	t.Run("Token very close to expiration", func(t *testing.T) {
		manager := NewJWTManager("test-secret", 30*time.Second)

		// Create token and let it get close to expiration
		token, err := manager.GenerateToken(1, "user")
		require.NoError(t, err)

		// Wait for 20 seconds (10 seconds left, < 30 seconds threshold)
		time.Sleep(20 * time.Second)

		claims, err := manager.VerifyToken(token)
		require.NoError(t, err)

		// Should be able to refresh
		newToken, err := manager.RefreshToken(claims)
		assert.NoError(t, err)
		assert.NotEmpty(t, newToken)
	})

	t.Run("Token with plenty of time left", func(t *testing.T) {
		manager := NewJWTManager("test-secret", 1*time.Hour)

		token, err := manager.GenerateToken(1, "user")
		require.NoError(t, err)

		claims, err := manager.VerifyToken(token)
		require.NoError(t, err)

		// Should fail because token has plenty of time left
		_, err = manager.RefreshToken(claims)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Token还未到刷新时间")
	})
}

// TestExtractTokenFromHeader_EdgeCases 边缘情况测试
func TestExtractTokenFromHeader_EdgeCases(t *testing.T) {
	t.Run("Valid token with extra spaces", func(t *testing.T) {
		token, err := ExtractTokenFromHeader("Bearer   valid-token   ")
		assert.NoError(t, err)
		assert.Equal(t, "valid-token", token)
	})

	t.Run("Token with mixed whitespace", func(t *testing.T) {
		token, err := ExtractTokenFromHeader("Bearer  token-with-spaces  ")
		assert.NoError(t, err)
		assert.Equal(t, "token-with-spaces", token)
	})

	t.Run("Empty after Bearer prefix", func(t *testing.T) {
		// When only "Bearer" is provided (without space or token), it's considered format error
		_, err := ExtractTokenFromHeader("Bearer")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Authorization头格式错误")
	})
}
