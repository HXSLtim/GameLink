package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// JWTManagerTestSuite tests JWT token generation and validation.
type JWTManagerTestSuite struct {
	suite.Suite
	manager *JWTManager
}

// SetupSuite initializes the test suite.
func (s *JWTManagerTestSuite) SetupSuite() {
	s.manager = NewJWTManager("test-secret-key-for-jwt-token-generation-and-validation", 24*time.Hour)
}

// TestGenerateTokenSuccess tests successful token generation.
func (s *JWTManagerTestSuite) TestGenerateTokenSuccess() {
	token, err := s.manager.GenerateToken(1, "user")
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), token)
	assert.IsType(s.T(), "", token)
}

// TestGenerateTokenWithDifferentUserIDs tests token generation with different user IDs.
func (s *JWTManagerTestSuite) TestGenerateTokenWithDifferentUserIDs() {
	tests := []struct {
		name   string
		userID uint64
		role   string
	}{
		{"user 1", 1, "user"},
		{"user 999", 999, "admin"},
		{"user 12345", 12345, "player"},
		{"user 0", 0, "guest"},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			token, err := s.manager.GenerateToken(tt.userID, tt.role)
			assert.NoError(s.T(), err)
			assert.NotEmpty(s.T(), token)
		})
	}
}

// TestVerifyTokenSuccess tests successful token verification.
func (s *JWTManagerTestSuite) TestVerifyTokenSuccess() {
	token, err := s.manager.GenerateToken(1, "user")
	assert.NoError(s.T(), err)

	claims, err := s.manager.VerifyToken(token)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), claims)
	assert.Equal(s.T(), uint64(1), claims.UserID)
	assert.Equal(s.T(), "user", claims.Role)
}

// TestVerifyTokenInvalid tests verification of an invalid token.
func (s *JWTManagerTestSuite) TestVerifyTokenInvalid() {
	tests := []struct {
		name  string
		token string
	}{
		{"empty string", ""},
		{"invalid format", "invalid.token.format"},
		{"not a token", "this is not a jwt token"},
		{"truncated token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"},
		{"wrong signature", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJyb2xlIjoidXNlciIsImlzcyI6ImdhbWVsaW5rIn0.wrongsignature"},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			_, err := s.manager.VerifyToken(tt.token)
			assert.Error(s.T(), err)
		})
	}
}

// TestVerifyTokenExpired tests verification of an expired token.
func (s *JWTManagerTestSuite) TestVerifyTokenExpired() {
	// Create a manager with very short duration
	shortLivedManager := NewJWTManager("test-secret", 1*time.Nanosecond)

	token, err := shortLivedManager.GenerateToken(1, "user")
	assert.NoError(s.T(), err)

	// Wait for token to expire
	time.Sleep(10 * time.Millisecond)

	_, err = shortLivedManager.VerifyToken(token)
	assert.Error(s.T(), err)
}

// TestVerifyTokenWithDifferentKeys tests verification with different keys.
func (s *JWTManagerTestSuite) TestVerifyTokenWithDifferentKeys() {
	manager1 := NewJWTManager("key-one-for-testing-purpose", 24*time.Hour)
	manager2 := NewJWTManager("key-two-for-testing-purpose", 24*time.Hour)

	token1, err := manager1.GenerateToken(1, "user")
	assert.NoError(s.T(), err)

	// Verify with correct key
	_, err = manager1.VerifyToken(token1)
	assert.NoError(s.T(), err)

	// Verify with wrong key
	_, err = manager2.VerifyToken(token1)
	assert.Error(s.T(), err)
}

// TestTokenClaimsIntegrity tests that claims are not tampered with.
func (s *JWTManagerTestSuite) TestTokenClaimsIntegrity() {
	token, err := s.manager.GenerateToken(12345, "admin")
	assert.NoError(s.T(), err)

	claims, err := s.manager.VerifyToken(token)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), uint64(12345), claims.UserID)
	assert.Equal(s.T(), "admin", claims.Role)
}

// TestExtractTokenFromHeader tests extracting token from Authorization header.
func (s *JWTManagerTestSuite) TestExtractTokenFromHeader() {
	tests := []struct {
		name        string
		authHeader  string
		expectError bool
		expected    string
	}{
		{
			name:        "valid bearer token",
			authHeader:  "Bearer valid-token-string",
			expectError: false,
			expected:    "valid-token-string",
		},
		{
			name:        "missing bearer prefix",
			authHeader:  "valid-token-string",
			expectError: true,
		},
		{
			name:        "empty header",
			authHeader:  "",
			expectError: true,
		},
		{
			name:        "bearer with empty token",
			authHeader:  "Bearer ",
			expectError: true,
		},
		{
			name:        "bearer with extra spaces",
			authHeader:  "Bearer  token-with-spaces  ",
			expectError: false,
			expected:    "token-with-spaces",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			token, err := ExtractTokenFromHeader(tt.authHeader)
			if tt.expectError {
				assert.Error(s.T(), err)
			} else {
				assert.NoError(s.T(), err)
				assert.Equal(s.T(), tt.expected, token)
			}
		})
	}
}

// TestIsTokenExpired tests token expiration check.
func (s *JWTManagerTestSuite) TestIsTokenExpired() {
	// Create a manager with very short duration
	shortLivedManager := NewJWTManager("test-secret", 1*time.Nanosecond)

	token, err := shortLivedManager.GenerateToken(1, "user")
	assert.NoError(s.T(), err)

	// Wait for token to expire
	time.Sleep(10 * time.Millisecond)

	_, err = shortLivedManager.VerifyToken(token)
	// Should fail because it's expired
	assert.Error(s.T(), err)

	// Test with a valid token
	validManager := NewJWTManager("test-secret", 24*time.Hour)
	validToken, err := validManager.GenerateToken(1, "user")
	assert.NoError(s.T(), err)

	validClaims, err := validManager.VerifyToken(validToken)
	assert.NoError(s.T(), err)

	// Should not be expired
	assert.False(s.T(), IsTokenExpired(validClaims))
}

// TestGetTokenRemainingTime tests getting token remaining time.
func (s *JWTManagerTestSuite) TestGetTokenRemainingTime() {
	// Create a manager with 1 hour duration
	manager := NewJWTManager("test-secret", 1*time.Hour)

	token, err := manager.GenerateToken(1, "user")
	assert.NoError(s.T(), err)

	claims, err := manager.VerifyToken(token)
	assert.NoError(s.T(), err)

	remaining := GetTokenRemainingTime(claims)
	assert.Greater(s.T(), remaining, 0*time.Second)
	assert.LessOrEqual(s.T(), remaining, 1*time.Hour)
}

// TestJWTManagerWithRealisticDurations tests with realistic token durations.
func (s *JWTManagerTestSuite) TestJWTManagerWithRealisticDurations() {
	tests := []struct {
		name     string
		duration time.Duration
	}{
		{"15 minutes", 15 * time.Minute},
		{"1 hour", 1 * time.Hour},
		{"24 hours", 24 * time.Hour},
		{"7 days", 7 * 24 * time.Hour},
		{"30 days", 30 * 24 * time.Hour},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			manager := NewJWTManager("test-secret", tt.duration)
			token, err := manager.GenerateToken(1, "user")
			assert.NoError(s.T(), err)

			claims, err := manager.VerifyToken(token)
			assert.NoError(s.T(), err)
			assert.Equal(s.T(), uint64(1), claims.UserID)
		})
	}
}

// TestJWTManager runs the JWT manager test suite.
func TestJWTManager(t *testing.T) {
	suite.Run(t, new(JWTManagerTestSuite))
}