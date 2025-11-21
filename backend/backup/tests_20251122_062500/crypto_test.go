package middleware

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/config"
	"gamelink/internal/model"
)

const (
	testSecret = "GameLink2025SecretKey!@#"
	testIV     = "GameLink2025IV!!!"
)

func TestCryptoMiddleware_DecryptsPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(Crypto(config.CryptoConfig{
		Enabled:      true,
		SecretKey:    testSecret,
		IV:           testIV,
		Methods:      []string{"POST"},
		ExcludePaths: []string{"/api/v1/health"},
		UseSignature: true,
	}))

	expected := `{"username":"admin","password":"123456"}`
	router.POST("/api/v1/auth/login", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			t.Fatalf("failed to read decrypted body: %v", err)
		}
		if string(body) != expected {
			t.Fatalf("unexpected body, want %s got %s", expected, string(body))
		}
		c.JSON(http.StatusOK, model.APIResponse[any]{Success: true, Code: http.StatusOK, Message: "OK"})
	})

	timestamp := time.Now().UnixMilli()
	payload, err := encryptPayload([]byte(expected), []byte(testSecret), []byte(testIV))
	if err != nil {
		t.Fatalf("encrypt payload: %v", err)
	}
	signature := generateSignature([]byte(expected), timestamp, testSecret)

	reqBody := map[string]any{
		"encrypted": true,
		"payload":   payload,
		"timestamp": timestamp,
		"signature": signature,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status code: %d, body: %s", w.Code, w.Body.String())
	}
}

func TestCryptoMiddleware_SignatureMismatch(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(Crypto(config.CryptoConfig{
		Enabled:      true,
		SecretKey:    testSecret,
		IV:           testIV,
		Methods:      []string{"POST"},
		ExcludePaths: nil,
		UseSignature: true,
	}))
	router.POST("/api/v1/orders", func(c *gin.Context) {
		t.Fatal("handler should not be called when signature mismatch occurs")
	})

	expected := `{"order_id":1001,"amount":199.99}`
	timestamp := time.Now().UnixMilli()
	payload, err := encryptPayload([]byte(expected), []byte(testSecret), []byte(testIV))
	if err != nil {
		t.Fatalf("encrypt payload: %v", err)
	}

	reqBody := map[string]any{
		"encrypted": true,
		"payload":   payload,
		"timestamp": timestamp,
		"signature": "invalid-signature",
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "请求签名验证失败") {
		t.Fatalf("expected signature error message, body: %s", w.Body.String())
	}
}

func TestCryptoMiddleware_SkipsExcludedPath(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(Crypto(config.CryptoConfig{
		Enabled:      true,
		SecretKey:    testSecret,
		IV:           testIV,
		Methods:      []string{"POST"},
		ExcludePaths: []string{"/api/v1/ping"},
		UseSignature: true,
	}))

	router.POST("/api/v1/ping", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		// 当路径在白名单时，原始加密内容应原样透传
		if !bytes.Contains(body, []byte(`"payload"`)) {
			t.Fatalf("expected encrypted payload to be forwarded, got %s", string(body))
		}
		c.JSON(http.StatusOK, model.APIResponse[any]{Success: true, Code: http.StatusOK, Message: "OK"})
	})

	expected := `{"message":"pong"}`
	timestamp := time.Now().UnixMilli()
	payload, err := encryptPayload([]byte(expected), []byte(testSecret), []byte(testIV))
	if err != nil {
		t.Fatalf("encrypt payload: %v", err)
	}
	signature := generateSignature([]byte(expected), timestamp, testSecret)

	reqBody := map[string]any{
		"encrypted": true,
		"payload":   payload,
		"timestamp": timestamp,
		"signature": signature,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/ping", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status code: %d", w.Code)
	}
}

func encryptPayload(plain, key, iv []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	if len(iv) < blockSize {
		return "", errors.New("iv length must be at least the block size")
	}
	padded := pkcs7Pad(plain, blockSize)

	ciphertext := make([]byte, len(padded))
	copyIV := make([]byte, blockSize)
	copy(copyIV, iv[:blockSize])
	mode := cipher.NewCBCEncrypter(block, copyIV)
	mode.CryptBlocks(ciphertext, padded)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	if padding == 0 {
		padding = blockSize
	}
	pad := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, pad...)
}
