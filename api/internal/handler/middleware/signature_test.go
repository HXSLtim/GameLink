package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gamelink/pkg/config"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestSignature_Disabled(t *testing.T) {
	cfg := config.SignatureConfig{
		Enabled: false,
	}

	middleware := Signature(cfg)

	router := gin.New()
	router.Use(middleware)
	router.POST("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("POST", "/test", bytes.NewReader([]byte(`{"test":"data"}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSignature_MissingSecretKey(t *testing.T) {
	cfg := config.SignatureConfig{
		Enabled:   true,
		SecretKey: "",
	}

	middleware := Signature(cfg)

	router := gin.New()
	router.Use(middleware)
	router.POST("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("POST", "/test", bytes.NewReader([]byte(`{"test":"data"}`)))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should pass through when secret key is missing (middleware disables itself)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSignature_ValidSignature(t *testing.T) {
	secretKey := "test-secret-key-for-signature-validation"
	cfg := config.SignatureConfig{
		Enabled:    true,
		SecretKey:  secretKey,
		HeaderName: "X-Signature",
		Methods:    []string{"POST"},
	}

	middleware := Signature(cfg)

	router := gin.New()
	router.Use(middleware)
	router.POST("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	body := []byte(`{"test":"data"}`)
	signature := calculateSignature("POST", "/test", body, secretKey)

	req := httptest.NewRequest("POST", "/test", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Signature", signature)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSignature_InvalidSignature(t *testing.T) {
	secretKey := "test-secret-key-for-signature-validation"
	cfg := config.SignatureConfig{
		Enabled:    true,
		SecretKey:  secretKey,
		HeaderName: "X-Signature",
		Methods:    []string{"POST"},
	}

	middleware := Signature(cfg)

	router := gin.New()
	router.Use(middleware)
	router.POST("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("POST", "/test", bytes.NewReader([]byte(`{"test":"data"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Signature", "invalid-signature")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid signature")
}

func TestSignature_MissingSignature(t *testing.T) {
	cfg := config.SignatureConfig{
		Enabled:    true,
		SecretKey:  "test-secret-key",
		HeaderName: "X-Signature",
		Methods:    []string{"POST"},
	}

	middleware := Signature(cfg)

	router := gin.New()
	router.Use(middleware)
	router.POST("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("POST", "/test", bytes.NewReader([]byte(`{"test":"data"}`)))
	req.Header.Set("Content-Type", "application/json")
	// Don't set X-Signature header
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "missing signature header")
}

func TestSignature_ExcludePath(t *testing.T) {
	cfg := config.SignatureConfig{
		Enabled:      true,
		SecretKey:    "test-secret-key",
		HeaderName:   "X-Signature",
		Methods:      []string{"POST"},
		ExcludePaths: []string{"/api/v1/health", "/api/v1/ping"},
	}

	middleware := Signature(cfg)

	router := gin.New()
	router.Use(middleware)
	router.POST("/api/v1/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("POST", "/api/v1/health", bytes.NewReader([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	// Don't set X-Signature header - should still pass
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSignature_MethodNotInList(t *testing.T) {
	cfg := config.SignatureConfig{
		Enabled:    true,
		SecretKey:  "test-secret-key",
		HeaderName: "X-Signature",
		Methods:    []string{"POST", "PUT"}, // GET not included
	}

	middleware := Signature(cfg)

	router := gin.New()
	router.Use(middleware)
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSignature_SignatureAlgorithm(t *testing.T) {
	secretKey := "test-secret-key"
	body := []byte(`{"name":"test","value":123}`)

	tests := []struct {
		name     string
		method   string
		path     string
		body     []byte
		expected string
	}{
		{
			name:     "POST request",
			method:   "POST",
			path:     "/api/v1/users",
			body:     body,
			expected: calculateSignature("POST", "/api/v1/users", body, secretKey),
		},
		{
			name:     "PUT request",
			method:   "PUT",
			path:     "/api/v1/users/123",
			body:     body,
			expected: calculateSignature("PUT", "/api/v1/users/123", body, secretKey),
		},
		{
			name:     "PATCH request",
			method:   "PATCH",
			path:     "/api/v1/users/123",
			body:     body,
			expected: calculateSignature("PATCH", "/api/v1/users/123", body, secretKey),
		},
		{
			name:     "DELETE request with empty body",
			method:   "DELETE",
			path:     "/api/v1/users/123",
			body:     []byte{},
			expected: calculateSignature("DELETE", "/api/v1/users/123", []byte{}, secretKey),
		},
		{
			name:     "Different body produces different signature",
			method:   "POST",
			path:     "/api/v1/users",
			body:     []byte(`{"different":"body"}`),
			expected: calculateSignature("POST", "/api/v1/users", []byte(`{"different":"body"}`), secretKey),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify signature is deterministic
			sig1 := calculateSignature(tt.method, tt.path, tt.body, secretKey)
			sig2 := calculateSignature(tt.method, tt.path, tt.body, secretKey)
			assert.Equal(t, sig1, sig2, "signature should be deterministic")
			assert.Equal(t, tt.expected, sig1)

			// Verify different inputs produce different signatures
			if tt.name == "POST request" {
				differentSig := calculateSignature("GET", tt.path, tt.body, secretKey)
				assert.NotEqual(t, sig1, differentSig, "different method should produce different signature")

				differentSig = calculateSignature(tt.method, "/different/path", tt.body, secretKey)
				assert.NotEqual(t, sig1, differentSig, "different path should produce different signature")

				differentSig = calculateSignature(tt.method, tt.path, []byte(`{}`), secretKey)
				assert.NotEqual(t, sig1, differentSig, "different body should produce different signature")
			}
		})
	}
}

func TestSignature_ContextValues(t *testing.T) {
	secretKey := "test-secret-key"
	cfg := config.SignatureConfig{
		Enabled:    true,
		SecretKey:  secretKey,
		HeaderName: "X-Signature",
		Methods:    []string{"POST"},
	}

	middleware := Signature(cfg)

	router := gin.New()
	router.Use(middleware)
	router.POST("/test", func(c *gin.Context) {
		valid := IsSignatureValid(c)
		_, signature := GetSignatureValidation(c)

		assert.True(t, valid, "signature should be valid")
		assert.NotEmpty(t, signature, "signature value should be set")

		c.Status(http.StatusOK)
	})

	body := []byte(`{"test":"data"}`)
	signature := calculateSignature("POST", "/test", body, secretKey)

	req := httptest.NewRequest("POST", "/test", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Signature", signature)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSignature_WildcardExcludePath(t *testing.T) {
	cfg := config.SignatureConfig{
		Enabled:      true,
		SecretKey:    "test-secret-key",
		HeaderName:   "X-Signature",
		Methods:      []string{"POST"},
		ExcludePaths: []string{"/api/v1/public/*", "/api/v1/health"},
	}

	middleware := Signature(cfg)

	router := gin.New()
	router.Use(middleware)
	router.POST("/api/v1/public/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	router.POST("/api/v1/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	router.POST("/api/v1/private/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Test wildcard exclude - should pass without signature
	req1 := httptest.NewRequest("POST", "/api/v1/public/test", bytes.NewReader([]byte(`{}`)))
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Test exact path exclude - should pass without signature
	req2 := httptest.NewRequest("POST", "/api/v1/health", bytes.NewReader([]byte(`{}`)))
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	// Test non-excluded path - should fail without signature
	req3 := httptest.NewRequest("POST", "/api/v1/private/test", bytes.NewReader([]byte(`{}`)))
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusUnauthorized, w3.Code)
}

func TestSignature_CustomHeaderName(t *testing.T) {
	secretKey := "test-secret-key"
	customHeader := "X-Custom-Signature"
	cfg := config.SignatureConfig{
		Enabled:    true,
		SecretKey:  secretKey,
		HeaderName: customHeader,
		Methods:    []string{"POST"},
	}

	middleware := Signature(cfg)

	router := gin.New()
	router.Use(middleware)
	router.POST("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	body := []byte(`{"test":"data"}`)
	signature := calculateSignature("POST", "/test", body, secretKey)

	req := httptest.NewRequest("POST", "/test", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(customHeader, signature)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSignature_EmptyBody(t *testing.T) {
	secretKey := "test-secret-key"
	cfg := config.SignatureConfig{
		Enabled:    true,
		SecretKey:  secretKey,
		HeaderName: "X-Signature",
		Methods:    []string{"POST", "DELETE"},
	}

	middleware := Signature(cfg)

	router := gin.New()
	router.Use(middleware)
	router.POST("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	router.DELETE("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Test POST with empty body
	body := []byte{}
	signature := calculateSignature("POST", "/test", body, secretKey)

	req := httptest.NewRequest("POST", "/test", bytes.NewReader(body))
	req.Header.Set("X-Signature", signature)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test DELETE with empty body
	signature = calculateSignature("DELETE", "/test", body, secretKey)

	req = httptest.NewRequest("DELETE", "/test", bytes.NewReader(body))
	req.Header.Set("X-Signature", signature)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
