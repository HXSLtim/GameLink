package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/pkg/config"
)

const (
	// SignatureHeader is the default header name for request signature
	SignatureHeader = "X-Signature"
)

// Signature provides HMAC-SHA256 request signature validation middleware.
// The signature is calculated as: HMAC-SHA256(METHOD:PATH:BODY, secret_key)
// This ensures request integrity and authenticity.
//
// Signature Format:
//   - Header: X-Signature (configurable via signature.header_name)
//   - Algorithm: HMAC-SHA256
//   - Message: METHOD:PATH:BODY (e.g., "POST:/api/v1/users:{"name":"test"}")
//   - Key: signature_secret_key from config
//
// Example (using curl with openssl):
//
//	# Calculate signature
//	BODY='{"name":"test"}'
//	MESSAGE="POST:/api/v1/users:${BODY}"
//	SIGNATURE=$(echo -n "$MESSAGE" | openssl dgst -sha256 -hmac "your-secret-key" | awk '{print $2}')
//	curl -X POST http://localhost:8080/api/v1/users \
//	  -H "Content-Type: application/json" \
//	  -H "X-Signature: $SIGNATURE" \
//	  -d "$BODY"
func Signature(cfg config.SignatureConfig) gin.HandlerFunc {
	// If signature validation is disabled, return a no-op middleware
	if !cfg.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	// Validate configuration
	if cfg.SecretKey == "" {
		slog.Error("signature middleware enabled but secret key is empty, disabling middleware")
		return func(c *gin.Context) {
			c.Next()
		}
	}

	// Normalize header name (ensure it starts with "X-" if custom)
	headerName := normalizeHeaderName(cfg.HeaderName)
	if headerName == "" {
		headerName = SignatureHeader
	}

	// Build methods set for quick lookup
	methods := make(map[string]struct{}, len(cfg.Methods))
	for _, m := range cfg.Methods {
		methods[strings.ToUpper(strings.TrimSpace(m))] = struct{}{}
	}

	// Pre-normalize exclude paths for pattern matching
	excludePaths := cfg.ExcludePaths

	secretKey := cfg.SecretKey

	slog.Info("signature middleware enabled",
		"header", headerName,
		"methods", cfg.Methods,
		"exclude_paths", len(excludePaths))

	return func(c *gin.Context) {
		// Skip if method not in configured list
		if len(methods) > 0 {
			method := strings.ToUpper(c.Request.Method)
			if _, ok := methods[method]; !ok {
				c.Next()
				return
			}
		}

		// Skip if path is in exclude list
		path := c.Request.URL.Path
		for _, exclude := range excludePaths {
			if shouldExcludePath(path, exclude) {
				c.Next()
				return
			}
		}

		// Get signature from header
		clientSignature := c.GetHeader(headerName)
		if clientSignature == "" {
			slog.Warn("signature validation failed: missing signature header",
				"path", path,
				"method", c.Request.Method,
				"header", headerName)
			respondWithError(c, http.StatusUnauthorized, "missing signature header")
			return
		}

		// Read request body
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			slog.Warn("signature validation failed: unable to read request body",
				"path", path,
				"error", err)
			respondWithError(c, http.StatusBadRequest, "unable to read request body")
			return
		}
		_ = c.Request.Body.Close()

		// Restore request body for subsequent handlers
		c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		// Calculate expected signature
		expectedSignature := calculateSignature(c.Request.Method, path, bodyBytes, secretKey)

		// Compare signatures using constant-time comparison to prevent timing attacks
		if !hmac.Equal([]byte(expectedSignature), []byte(clientSignature)) {
			slog.Warn("signature validation failed: signature mismatch",
				"path", path,
				"method", c.Request.Method,
				"client_ip", c.ClientIP())
			respondWithError(c, http.StatusUnauthorized, "invalid signature")
			return
		}

		// Signature is valid, continue to next handler
		slog.Debug("signature validation successful",
			"path", path,
			"method", c.Request.Method)

		// Store signature validation result in context for downstream handlers
		c.Set("signature.valid", true)
		c.Set("signature.value", clientSignature)

		c.Next()
	}
}

// calculateSignature computes HMAC-SHA256 signature for the request
// Format: HMAC-SHA256(METHOD:PATH:BODY, secret_key)
func calculateSignature(method, path string, body []byte, secretKey string) string {
	// Build message: METHOD:PATH:BODY
	message := strings.ToUpper(method) + ":" + path + ":" + string(body)

	// Create HMAC-SHA256 hash
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(message))

	// Return hex-encoded signature
	return hex.EncodeToString(h.Sum(nil))
}

// shouldExcludePath checks if a path should be excluded based on pattern matching.
// Supports three matching modes:
// 1. Exact match: "/api/v1/health" only matches exactly "/api/v1/health"
// 2. Prefix match with wildcard: "/api/v1/public/*" matches "/api/v1/public/" and all subpaths
// 3. Catch-all: "*" matches any path
func shouldExcludePath(path string, exclude string) bool {
	exclude = strings.TrimSpace(exclude)
	if exclude == "" {
		return false
	}

	// Catch-all: * matches any path
	if exclude == "*" {
		return true
	}

	// Exact match
	if path == exclude {
		return true
	}

	// Prefix match (ends with /*)
	if strings.HasSuffix(exclude, "/*") {
		prefix := strings.TrimSuffix(exclude, "/*")
		// Match prefix/ or exact prefix (for "/api/v1/public/*" matches "/api/v1/public" and "/api/v1/public/...")
		return strings.HasPrefix(path, prefix+"/") || path == prefix
	}

	// Suffix match (starts with *)
	if strings.HasPrefix(exclude, "*") {
		suffix := strings.TrimPrefix(exclude, "*")
		return strings.HasSuffix(path, suffix)
	}

	// Default: no match
	return false
}

// normalizeHeaderName ensures the header name is properly formatted
func normalizeHeaderName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}
	// Ensure header name starts with "X-"
	if !strings.HasPrefix(name, "X-") && !strings.HasPrefix(name, "x-") {
		return "X-" + name
	}
	// Capitalize first letter after X-
	if len(name) > 2 {
		return "X-" + strings.ToUpper(name[2:3])
	}
	return name
}

// respondWithError sends a standardized error response
func respondWithError(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, model.APIResponse[any]{
		Success: false,
		Code:    status,
		Message: message,
	})
}

// GetSignatureValidation retrieves the signature validation result from context
func GetSignatureValidation(c *gin.Context) (valid bool, signature string) {
	if val, exists := c.Get("signature.valid"); exists {
		valid, _ = val.(bool)
	}
	if val, exists := c.Get("signature.value"); exists {
		signature, _ = val.(string)
	}
	return
}

// IsSignatureValid checks if the request signature was validated successfully
func IsSignatureValid(c *gin.Context) bool {
	if val, exists := c.Get("signature.valid"); exists {
		valid, _ := val.(bool)
		return valid
	}
	return false
}
