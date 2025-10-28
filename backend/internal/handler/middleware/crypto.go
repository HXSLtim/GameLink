package middleware

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"gamelink/internal/config"
	"gamelink/internal/model"
)

type encryptedRequest struct {
	Encrypted bool   `json:"encrypted"`
	Payload   string `json:"payload"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
}

// Crypto 提供与前端 AES-256-CBC 加密协议对接的请求解密能力。
// 满足 config.Crypto 中的配置时会尝试解密请求体，并在成功后将明文请求体重新写入 gin.Context。
func Crypto(cfg config.CryptoConfig) gin.HandlerFunc {
	if !cfg.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	block, err := aes.NewCipher([]byte(cfg.SecretKey))
	if err != nil {
		slog.Error("failed to init crypto middleware", "error", err)
		return func(c *gin.Context) {
			abortWithCryptoError(c, http.StatusInternalServerError, "后端加密配置错误")
		}
	}

	methods := make(map[string]struct{}, len(cfg.Methods))
	for _, m := range cfg.Methods {
		methods[strings.ToUpper(strings.TrimSpace(m))] = struct{}{}
	}
	excludePaths := cfg.ExcludePaths
	iv := []byte(cfg.IV)

	return func(c *gin.Context) {
		if !shouldProcessRequest(c, methods, excludePaths) {
			c.Next()
			return
		}

		if c.Request.Body == nil {
			c.Next()
			return
		}

		originalBody := c.Request.Body
		bodyBytes, err := io.ReadAll(originalBody)
		if err != nil {
			slog.Warn("crypto middleware: read body failed", "error", err)
			abortWithCryptoError(c, http.StatusBadRequest, "无法读取请求体")
			return
		}
		_ = originalBody.Close()

		if len(bytes.TrimSpace(bodyBytes)) == 0 {
			restoreRequestBody(c, bodyBytes)
			c.Next()
			return
		}

		var req encryptedRequest
		if err := json.Unmarshal(bodyBytes, &req); err != nil || !req.Encrypted {
			restoreRequestBody(c, bodyBytes)
			c.Next()
			return
		}

		plain, err := decryptPayload(block, iv, req.Payload)
		if err != nil {
			slog.Warn("crypto middleware: decrypt failed", "error", err)
			abortWithCryptoError(c, http.StatusBadRequest, "请求数据解密失败")
			return
		}

		if cfg.UseSignature {
			if req.Signature == "" || req.Timestamp == 0 {
				abortWithCryptoError(c, http.StatusBadRequest, "缺少签名或时间戳")
				return
			}
			expected := generateSignature(plain, req.Timestamp, cfg.SecretKey)
			if !strings.EqualFold(expected, req.Signature) {
				slog.Warn("crypto middleware: signature mismatch", "expected", expected, "got", req.Signature)
				abortWithCryptoError(c, http.StatusBadRequest, "请求签名验证失败")
				return
			}
		}

		restoreRequestBody(c, plain)
		c.Set("crypto.encrypted", true)
		c.Set("crypto.timestamp", req.Timestamp)
		if req.Signature != "" {
			c.Set("crypto.signature", req.Signature)
		}
		c.Set("crypto.raw_payload", string(plain))
		c.Next()
	}
}

func shouldProcessRequest(c *gin.Context, methods map[string]struct{}, excludePaths []string) bool {
	method := strings.ToUpper(c.Request.Method)
	if len(methods) > 0 {
		if _, ok := methods[method]; !ok {
			return false
		}
	}
	path := c.Request.URL.Path
	for _, exclude := range excludePaths {
		exclude = strings.TrimSpace(exclude)
		if exclude == "" {
			continue
		}
		if strings.Contains(path, exclude) {
			return false
		}
	}
	return true
}

func decryptPayload(block cipher.Block, iv []byte, payload string) ([]byte, error) {
	raw, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return nil, err
	}
	if len(iv) < block.BlockSize() {
		return nil, errors.New("invalid iv length")
	}
	if len(raw) == 0 || len(raw)%block.BlockSize() != 0 {
		return nil, errors.New("ciphertext length is not aligned with block size")
	}

	plain := make([]byte, len(raw))
	copyIV := make([]byte, block.BlockSize())
	copy(copyIV, iv[:block.BlockSize()])
	mode := cipher.NewCBCDecrypter(block, copyIV)
	mode.CryptBlocks(plain, raw)

	return pkcs7Unpad(plain, block.BlockSize())
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, errors.New("invalid padded data size")
	}
	padding := int(data[len(data)-1])
	if padding == 0 || padding > blockSize {
		return nil, errors.New("invalid padding size")
	}
	for i := len(data) - padding; i < len(data); i++ {
		if data[i] != byte(padding) {
			return nil, errors.New("invalid padding bytes")
		}
	}
	return data[:len(data)-padding], nil
}

func generateSignature(plain []byte, timestamp int64, secret string) string {
	message := string(plain) + strconv.FormatInt(timestamp, 10) + secret
	hash := sha256.Sum256([]byte(message))
	return hex.EncodeToString(hash[:])
}

func restoreRequestBody(c *gin.Context, data []byte) {
	c.Request.Body = io.NopCloser(bytes.NewReader(data))
	c.Request.ContentLength = int64(len(data))
	c.Request.Header.Set("Content-Length", strconv.Itoa(len(data)))
}

func abortWithCryptoError(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, model.APIResponse[any]{
		Success: false,
		Code:    status,
		Message: message,
	})
}
