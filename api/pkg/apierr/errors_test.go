package apierr

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIError_WithInternalError(t *testing.T) {
	// Test in development mode
	os.Setenv("APP_ENV", "development")
	defer os.Unsetenv("APP_ENV")

	err := errors.New("database connection failed")
	apiErr := BadRequest("Invalid input").WithInternalError(err)

	assert.NotNil(t, apiErr)
	assert.Equal(t, "Invalid input", apiErr.Message)
	assert.Equal(t, err, apiErr.internalError)
}

func TestAPIError_WithDetails_Development(t *testing.T) {
	// Test in development mode
	os.Setenv("APP_ENV", "development")
	defer os.Unsetenv("APP_ENV")

	apiErr := BadRequest("Invalid input").WithDetails("field 'email' is required")

	assert.Equal(t, "field 'email' is required", apiErr.Details)
}

func TestAPIError_WithDetails_Production(t *testing.T) {
	// Test in production mode
	os.Setenv("APP_ENV", "production")
	defer os.Unsetenv("APP_ENV")

	apiErr := BadRequest("Invalid input").WithDetails("field 'email' is required")

	// In production, details should not be exposed
	assert.Empty(t, apiErr.Details)
}

func TestAPIError_Sanitize_Development(t *testing.T) {
	os.Setenv("APP_ENV", "development")
	defer os.Unsetenv("APP_ENV")

	apiErr := InternalError("Database error").WithDetails("connection timeout")
	sanitized := apiErr.Sanitize()

	// In development, should return original
	assert.Equal(t, apiErr, sanitized)
	assert.Equal(t, "Database error", sanitized.Message)
}

func TestAPIError_Sanitize_Production_4xx(t *testing.T) {
	os.Setenv("APP_ENV", "production")
	defer os.Unsetenv("APP_ENV")

	apiErr := BadRequest("Invalid email format").WithDetails("email validation failed")
	sanitized := apiErr.Sanitize()

	// 4xx errors keep original message
	assert.Equal(t, "Invalid email format", sanitized.Message)
	// Details should be removed
	assert.Empty(t, sanitized.Details)
}

func TestAPIError_Sanitize_Production_5xx(t *testing.T) {
	os.Setenv("APP_ENV", "production")
	defer os.Unsetenv("APP_ENV")

	apiErr := InternalError("Database connection failed").WithDetails("timeout after 30s")
	sanitized := apiErr.Sanitize()

	// 5xx errors use generic message in production
	assert.Equal(t, "服务器内部错误，请稍后重试", sanitized.Message)
	// Details should be removed
	assert.Empty(t, sanitized.Details)
}

func TestAPIError_Sanitize_Staging(t *testing.T) {
	os.Setenv("APP_ENV", "staging")
	defer os.Unsetenv("APP_ENV")

	apiErr := InternalError("Database error").WithDetails("connection timeout")
	sanitized := apiErr.Sanitize()

	// Staging should behave like production
	assert.Equal(t, "服务器内部错误，请稍后重试", sanitized.Message)
	assert.Empty(t, sanitized.Details)
}

func TestIsProduction(t *testing.T) {
	tests := []struct {
		name     string
		env      string
		expected bool
	}{
		{"production", "production", true},
		{"staging", "staging", true},
		{"development", "development", false},
		{"empty", "", false},
		{"test", "test", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("APP_ENV", tt.env)
			defer os.Unsetenv("APP_ENV")

			result := isProduction()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAPIError_WithField(t *testing.T) {
	apiErr := BadRequest("Validation failed").WithField("email")

	assert.Equal(t, "email", apiErr.Field)
}

func TestAPIError_WithRequestID(t *testing.T) {
	apiErr := BadRequest("Invalid input").WithRequestID("req-123")

	assert.Equal(t, "req-123", apiErr.RequestID)
}

func TestAPIError_WithExtension(t *testing.T) {
	apiErr := BadRequest("Invalid input").
		WithExtension("retry_after", 60).
		WithExtension("error_code", "RATE_LIMIT")

	assert.Equal(t, 60, apiErr.Extensions["retry_after"])
	assert.Equal(t, "RATE_LIMIT", apiErr.Extensions["error_code"])
}
