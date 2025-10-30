package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("生成新的请求ID", func(t *testing.T) {
		router := gin.New()
		router.Use(RequestID())
		router.GET("/api/test", func(c *gin.Context) {
			requestID, exists := c.Get("request_id")
			if !exists {
				t.Error("Expected request_id to be set in context")
			}
			if requestID == "" {
				t.Error("Expected non-empty request ID")
			}
			c.JSON(http.StatusOK, gin.H{"request_id": requestID})
		})

		req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// 检查响应头是否包含 X-Request-ID
		requestID := w.Header().Get("X-Request-ID")
		if requestID == "" {
			t.Error("Expected X-Request-ID header to be set")
		}

		// 检查请求ID长度（32个字符的十六进制）
		if len(requestID) != 32 {
			t.Errorf("Expected request ID length 32, got %d", len(requestID))
		}
	})

	t.Run("使用客户端提供的请求ID", func(t *testing.T) {
		clientRequestID := "client-provided-request-id-123"

		router := gin.New()
		router.Use(RequestID())
		router.GET("/api/test", func(c *gin.Context) {
			requestID, _ := c.Get("request_id")
			c.JSON(http.StatusOK, gin.H{"request_id": requestID})
		})

		req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
		req.Header.Set("X-Request-ID", clientRequestID)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// 检查响应头是否使用了客户端提供的请求ID
		requestID := w.Header().Get("X-Request-ID")
		if requestID != clientRequestID {
			t.Errorf("Expected request ID %s, got %s", clientRequestID, requestID)
		}
	})

	t.Run("每次请求生成不同的ID", func(t *testing.T) {
		router := gin.New()
		router.Use(RequestID())
		router.GET("/api/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// 第一次请求
		req1 := httptest.NewRequest(http.MethodGet, "/api/test", nil)
		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req1)
		requestID1 := w1.Header().Get("X-Request-ID")

		// 第二次请求
		req2 := httptest.NewRequest(http.MethodGet, "/api/test", nil)
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		requestID2 := w2.Header().Get("X-Request-ID")

		// 两次请求的ID应该不同
		if requestID1 == requestID2 {
			t.Error("Expected different request IDs for different requests")
		}
	})
}

func TestRandomID(t *testing.T) {
	t.Run("生成32字符的十六进制字符串", func(t *testing.T) {
		id := randomID()
		if len(id) != 32 {
			t.Errorf("Expected ID length 32, got %d", len(id))
		}

		// 检查是否为有效的十六进制字符串
		for _, c := range id {
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
				t.Errorf("Expected hex character, got %c", c)
			}
		}
	})

	t.Run("每次调用生成不同的ID", func(t *testing.T) {
		id1 := randomID()
		id2 := randomID()

		if id1 == id2 {
			t.Error("Expected different IDs")
		}
	})
}

