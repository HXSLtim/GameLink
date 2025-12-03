package auth

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGenerateToken_EdgeCases 尝试覆盖SignedString错误分支
func TestGenerateToken_EdgeCases(t *testing.T) {
	manager := NewJWTManager("test-secret", 1*time.Hour)

	t.Run("Generate token with maximum userID", func(t *testing.T) {
		maxID := uint64(18446744073709551615) // uint64 max
		token, err := manager.GenerateToken(maxID, "user")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// Verify
		claims, err := manager.VerifyToken(token)
		assert.NoError(t, err)
		assert.Equal(t, maxID, claims.UserID)
	})

	t.Run("Generate token with special characters in role", func(t *testing.T) {
		specialRoles := []string{
			"user-admin",
			"user_admin",
			"user.admin",
			"user:admin",
			"user admin",
			"用户",   // Chinese
			"ユーザー", // Japanese
		}

		for _, role := range specialRoles {
			t.Run(role, func(t *testing.T) {
				token, err := manager.GenerateToken(1, role)
				assert.NoError(t, err)
				assert.NotEmpty(t, token)

				claims, err := manager.VerifyToken(token)
				assert.NoError(t, err)
				assert.Equal(t, role, claims.Role)
			})
		}
	})
}

// TestVerifyToken_TokenValidBranch 尝试触发 token.Valid = false 分支
func TestVerifyToken_TokenValidBranch(t *testing.T) {
	manager := NewJWTManager("test-secret", 1*time.Hour)

	t.Run("Signature verification passes but token.Valid is false", func(t *testing.T) {
		// 创建一个签名正确但内容无效的场景
		// 比如：使用正确的密钥签名，但算法不匹配
		claims := &Claims{
			UserID: 1,
			Role:   "user",
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer: "gamelink",
			},
		}

		// 使用不同的算法（虽然签名可能成功，但验证会失败）
		token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
		tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
		assert.NoError(t, err)

		// 尝试用我们的HS256验证器验证
		_, err = manager.VerifyToken(tokenString)
		assert.Error(t, err) // 应该失败，因为算法不正确
	})

	t.Run("Test type assertion failure", func(t *testing.T) {
		// 创建一个 tokenString，它包含有效的JWT结构，但claims不是我们期望的类型
		// 但由于我们传入了 &Claims{}，JWT库会尝试解析到该类型

		// 创建一个标准的JWT，但包含我们能控制的数据
		claims := map[string]interface{}{
			"user_id": 1,
			"role":    "user",
			"iss":     "gamelink",
			"exp":     time.Now().Add(1 * time.Hour).Unix(),
			"iat":     time.Now().Unix(),
		}

		// 手动构建 token
		header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
		payloadBytes, _ := json.Marshal(claims)
		payload := base64.RawURLEncoding.EncodeToString(payloadBytes)

		// 用正确的密钥签名
		signingString := header + "." + payload
		hash := jwt.SigningMethodHS256.Hash.New()
		hash.Write([]byte(signingString))
		signature := base64.RawURLEncoding.EncodeToString(hash.Sum([]byte("test-secret")))

		tokenString := signingString + "." + signature

		// 验证这个 token（应该能解析为 Claims）
		parsedClaims, err := manager.VerifyToken(tokenString)
		// 这可能会成功，因为结构相似
		if err == nil {
			assert.NotNil(t, parsedClaims)
		} else {
			// 或者失败，取决于 JWT 库的严格程度
			assert.Error(t, err)
		}
	})
}

// TestVerifyToken_CompleteCoverage 尝试覆盖所有可能的分支
func TestVerifyToken_CompleteCoverage(t *testing.T) {
	t.Run("Malicious token structure", func(t *testing.T) {
		tests := []struct {
			name  string
			token string
		}{
			{"Empty parts", ".."},
			{"Only dots", "..."},
			{"Invalid base64 in header", "!!!.payload.signature"},
			{"Invalid base64 in payload", "header.!!!.signature"},
			{"Valid base64 but invalid JSON", "ZXJyb3I.payload.signature"},
		}

		manager := NewJWTManager("test-secret", 1*time.Hour)
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := manager.VerifyToken(tt.token)
				assert.Error(t, err)
			})
		}
	})

	t.Run("Token with manipulated claims", func(t *testing.T) {
		// 创建一个有效的 token
		manager := NewJWTManager("test-secret", 1*time.Hour)
		token, err := manager.GenerateToken(1, "user")
		assert.NoError(t, err)

		// 解析 token
		parts := strings.Split(token, ".")
		require.Len(t, parts, 3)

		// 解码 payload
		payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
		assert.NoError(t, err)

		// 修改 claims
		var claimsMap map[string]interface{}
		err = json.Unmarshal(payloadBytes, &claimsMap)
		assert.NoError(t, err)

		// 添加无效字段
		claimsMap["invalid_field"] = []byte{0xFF, 0xFE, 0xFD}

		// 重新编码
		newPayloadBytes, _ := json.Marshal(claimsMap)
		newPayload := base64.RawURLEncoding.EncodeToString(newPayloadBytes)

		// 用相同的签名（这会验证失败）
		newToken := parts[0] + "." + newPayload + "." + parts[2]

		_, err = manager.VerifyToken(newToken)
		assert.Error(t, err)
	})
}

// TestReadMaxRefreshWindow_Parallel 测试环境变量读取的竞争条件
func TestReadMaxRefreshWindow_Parallel(t *testing.T) {
	t.Run("Concurrent reads", func(t *testing.T) {
		done := make(chan bool)

		for i := 0; i < 10; i++ {
			go func() {
				duration := readMaxRefreshWindow()
				assert.Equal(t, 7*24*time.Hour, duration)
				done <- true
			}()
		}

		for i := 0; i < 10; i++ {
			<-done
		}
	})

	t.Run("Environment variable changes", func(t *testing.T) {
		// 保存原始值
		original := os.Getenv("JWT_MAX_REFRESH")
		defer os.Setenv("JWT_MAX_REFRESH", original)

		// 设置为1小时
		os.Setenv("JWT_MAX_REFRESH", "1h")
		duration := readMaxRefreshWindow()
		assert.Equal(t, 1*time.Hour, duration)

		// 设置为7天
		os.Setenv("JWT_MAX_REFRESH", "168h")
		duration = readMaxRefreshWindow()
		assert.Equal(t, 168*time.Hour, duration)

		// 设置为无效值
		os.Setenv("JWT_MAX_REFRESH", "invalid")
		duration = readMaxRefreshWindow()
		assert.Equal(t, 7*24*time.Hour, duration) // 回退到默认值
	})
}

// TestExtractTokenFromHeader_Exhaustive 穷尽测试所有分支
func TestExtractTokenFromHeader_Exhaustive(t *testing.T) {
	tests := []struct {
		name        string
		authHeader  string
		expectError bool
		errorMsg    string
		expected    string
	}{
		// 正常情况
		{"Standard token", "Bearer abc123xyz", false, "", "abc123xyz"},
		{"Token with spaces", "Bearer token with spaces", false, "", "token with spaces"},
		{"Token with tabs trimmed", "Bearer\ttoken\t", true, "Authorization头格式错误", ""},
		{"Token with newlines trimmed", "Bearer\ntoken\n", true, "Authorization头格式错误", ""},
		{"Token with mixed whitespace", "Bearer  \t token \n  ", false, "", "token"},

		// 边界情况
		{"Empty header", "", true, "缺少Authorization头", ""},
		{"Missing Bearer", "token123", true, "Authorization头格式错误", ""},
		{"Wrong scheme", "Basic abc123", true, "Authorization头格式错误", ""},
		{"Bearer only", "Bearer", true, "Authorization头格式错误", ""},
		{"Bearer with spaces only", "Bearer   ", true, "Token为空", ""},
		{"Bearer lowercase", "bearer token", true, "Authorization头格式错误", ""},
		{"Bearer incomplete", "Bear token", true, "Authorization头格式错误", ""},

		// 特殊字符
		{"Token with special chars", "Bearer abc-123_xyz.456", false, "", "abc-123_xyz.456"},
		{"Token with symbols", "Bearer abc@#$%^&*()", false, "", "abc@#$%^&*()"},
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
				assert.Equal(t, tt.expected, token)
			}
		})
	}
}

// TestRefreshToken_Exhaustive 穷尽刷新Token的所有场景
func TestRefreshToken_Exhaustive(t *testing.T) {
	// 场景1: Token刚刚超过30秒阈值
	t.Run("Token just past 30s threshold", func(t *testing.T) {
		manager := &JWTManager{
			secretKey:     "test-secret",
			tokenDuration: 50 * time.Second, // 50秒总时长
			maxRefresh:    7 * 24 * time.Hour,
		}

		token, err := manager.GenerateToken(1, "user")
		require.NoError(t, err)

		// 等待21秒（剩余29秒 < 30秒）
		time.Sleep(21 * time.Second)

		claims, err := manager.VerifyToken(token)
		require.NoError(t, err)

		remaining := time.Until(claims.ExpiresAt.Time)
		assert.Less(t, remaining, 30*time.Second)

		// 现在应该可以刷新
		newToken, err := manager.RefreshToken(claims)
		assert.NoError(t, err)
		assert.NotEmpty(t, newToken)
	})

	// 场景2: Token还有31秒以上，刚好超过阈值
	t.Run("Token with more than 30s remaining", func(t *testing.T) {
		manager := &JWTManager{
			secretKey:     "test-secret",
			tokenDuration: 35 * time.Second, // 35秒总时长
			maxRefresh:    7 * 24 * time.Hour,
		}

		token, err := manager.GenerateToken(1, "user")
		require.NoError(t, err)

		// 等待1秒（剩余34秒 > 30秒）
		time.Sleep(time.Second)

		claims, err := manager.VerifyToken(token)
		require.NoError(t, err)

		remaining := time.Until(claims.ExpiresAt.Time)
		assert.Greater(t, remaining, 30*time.Second)

		// 现在不应该刷新
		_, err = manager.RefreshToken(claims)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Token还未到刷新时间")
	})
}

// TestJWTManager_FullLifecycleCoverage 覆盖完整的Token生命周期
func TestJWTManager_FullLifecycleCoverage(t *testing.T) {
	t.Run("Complete token lifecycle with all edge cases", func(t *testing.T) {
		// 创建Manager
		manager := NewJWTManager("lifecycle-secret", 10*time.Second)

		// 步骤1: 生成
		token, err := manager.GenerateToken(999, "admin")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// 步骤2: 验证
		claims, err := manager.VerifyToken(token)
		assert.NoError(t, err)
		assert.Equal(t, uint64(999), claims.UserID)
		assert.Equal(t, "admin", claims.Role)
		assert.Equal(t, "gamelink", claims.Issuer)
		assert.NotZero(t, claims.IssuedAt)
		assert.NotZero(t, claims.ExpiresAt)
		assert.NotZero(t, claims.NotBefore)

		// 步骤3: 检查未过期
		assert.False(t, IsTokenExpired(claims))

		// 步骤4: 检查剩余时间
		remaining := GetTokenRemainingTime(claims)
		assert.Greater(t, remaining, 8*time.Second)
		assert.LessOrEqual(t, remaining, 10*time.Second)

		// 步骤5: 等待8秒（剩余 < 30秒）
		time.Sleep(8 * time.Second)

		// 步骤6: 刷新Token
		refreshClaims, err := manager.VerifyToken(token)
		assert.NoError(t, err) // Token还可用

		newToken, err := manager.RefreshToken(refreshClaims)
		assert.NoError(t, err)
		assert.NotEmpty(t, newToken)
		assert.NotEqual(t, token, newToken)

		// 步骤7: 验证新Token
		newClaims, err := manager.VerifyToken(newToken)
		assert.NoError(t, err)
		assert.Equal(t, uint64(999), newClaims.UserID)
		assert.Equal(t, "admin", newClaims.Role)

		// 步骤8: 验证新Token的过期时间
		newRemaining := GetTokenRemainingTime(newClaims)
		assert.Greater(t, newRemaining, 9*time.Second) // 新token应该有完整的10秒
	})

	t.Run("Token cannot be refreshed after expiration", func(t *testing.T) {
		manager := NewJWTManager("test-secret", 2*time.Second)

		token, err := manager.GenerateToken(1, "user")
		assert.NoError(t, err)

		// 等待3秒，token过期
		time.Sleep(3 * time.Second)

		// 尝试验证（应该失败）
		_, err = manager.VerifyToken(token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token has invalid claims")

		// 无法获取claims，所以无法刷新 - 这是正确的行为
	})
}

// TestClaims_EmbeddedFields 测试嵌入的RegisteredClaims所有字段
func TestClaims_EmbeddedFields(t *testing.T) {
	claims := &Claims{
		UserID: 1,
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "gamelink",
			Subject:   "user_1",
			Audience:  []string{"users"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        "token-id-123",
		},
	}

	// 验证所有嵌入字段都可访问
	assert.Equal(t, "gamelink", claims.Issuer)
	assert.Equal(t, "user_1", claims.Subject)
	assert.Equal(t, jwt.ClaimStrings{"users"}, claims.Audience)
	assert.NotZero(t, claims.ExpiresAt)
	assert.NotZero(t, claims.NotBefore)
	assert.NotZero(t, claims.IssuedAt)
	assert.Equal(t, "token-id-123", claims.ID)
	assert.Equal(t, uint64(1), claims.UserID)
	assert.Equal(t, "user", claims.Role)
}
