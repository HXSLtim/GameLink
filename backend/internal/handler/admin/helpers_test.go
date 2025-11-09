package admin

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gamelink/internal/model"
)

func setupHelpersTestRouter() *gin.Engine {
	return newTestEngine()
}

func TestParseUintParam(t *testing.T) {
	r := setupHelpersTestRouter()
	r.GET("/test/:id", func(c *gin.Context) {
		id, err := parseUintParam(c, "id")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id})
	})

	// 测试有效ID
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test/123", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// 测试无效ID
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/test/invalid", nil)
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusBadRequest, w2.Code)
}

func TestQueryIntDefault(t *testing.T) {
	r := setupHelpersTestRouter()
	r.GET("/test", func(c *gin.Context) {
		value, err := queryIntDefault(c, "page", 1)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"value": value})
	})

	// 测试有值
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test?page=5", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// 测试默认值
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)

	// 测试无效值
	w3 := httptest.NewRecorder()
	req3 := httptest.NewRequest(http.MethodGet, "/test?page=invalid", nil)
	r.ServeHTTP(w3, req3)

	assert.Equal(t, http.StatusBadRequest, w3.Code)
}

func TestQueryUint64Ptr(t *testing.T) {
	r := setupHelpersTestRouter()
	r.GET("/test", func(c *gin.Context) {
		value, err := queryUint64Ptr(c, "user_id")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if value == nil {
			c.JSON(http.StatusOK, gin.H{"value": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"value": *value})
	})

	// 测试有值
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test?user_id=123", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// 测试空值
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)

	// 测试无效值
	w3 := httptest.NewRecorder()
	req3 := httptest.NewRequest(http.MethodGet, "/test?user_id=invalid", nil)
	r.ServeHTTP(w3, req3)

	assert.Equal(t, http.StatusBadRequest, w3.Code)
}

func TestQueryTimePtr(t *testing.T) {
	r := setupHelpersTestRouter()
	r.GET("/test", func(c *gin.Context) {
		value, err := queryTimePtr(c, "date_from")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if value == nil {
			c.JSON(http.StatusOK, gin.H{"value": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"value": value.Format(time.RFC3339)})
	})

	// 测试RFC3339格式
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test?date_from=2024-01-01T00:00:00Z", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// 测试日期格式
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/test?date_from=2024-01-01", nil)
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)

	// 测试空值
	w3 := httptest.NewRecorder()
	req3 := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w3, req3)

	assert.Equal(t, http.StatusOK, w3.Code)

	// 测试无效格式
	w4 := httptest.NewRecorder()
	req4 := httptest.NewRequest(http.MethodGet, "/test?date_from=invalid", nil)
	r.ServeHTTP(w4, req4)

	assert.Equal(t, http.StatusBadRequest, w4.Code)
}

func TestParseCSVParams(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "单个值",
			input:    []string{"pending"},
			expected: []string{"pending"},
		},
		{
			name:     "多个值",
			input:    []string{"pending", "confirmed"},
			expected: []string{"pending", "confirmed"},
		},
		{
			name:     "CSV格式",
			input:    []string{"pending,confirmed"},
			expected: []string{"pending", "confirmed"},
		},
		{
			name:     "混合格式",
			input:    []string{"pending", "confirmed,completed"},
			expected: []string{"pending", "confirmed", "completed"},
		},
		{
			name:     "空值",
			input:    []string{""},
			expected: []string{},
		},
		{
			name:     "带空格",
			input:    []string{" pending , confirmed "},
			expected: []string{"pending", "confirmed"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseCSVParams(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParsePagination(t *testing.T) {
	r := setupHelpersTestRouter()
	r.GET("/test", func(c *gin.Context) {
		page, pageSize, ok := parsePagination(c)
		if !ok {
			return
		}
		c.JSON(http.StatusOK, gin.H{"page": page, "page_size": pageSize})
	})

	// 测试默认值
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// 测试指定值
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/test?page=2&page_size=10", nil)
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)

	// 测试无效page
	w3 := httptest.NewRecorder()
	req3 := httptest.NewRequest(http.MethodGet, "/test?page=invalid", nil)
	r.ServeHTTP(w3, req3)

	assert.Equal(t, http.StatusBadRequest, w3.Code)
}

func TestNormalizeOrderStatus(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected model.OrderStatus
	}{
		{
			name:     "正常状态",
			input:    "pending",
			expected: model.OrderStatusPending,
		},
		{
			name:     "legacy cancelled",
			input:    "cancelled",
			expected: model.OrderStatusCanceled,
		},
		{
			name:     "带空格",
			input:    "  pending  ",
			expected: model.OrderStatusPending,
		},
		{
			name:     "大写",
			input:    "PENDING",
			expected: model.OrderStatusPending,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeOrderStatus(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildOrderListOptions(t *testing.T) {
	r := setupHelpersTestRouter()
	r.GET("/test", func(c *gin.Context) {
		opts, ok := buildOrderListOptions(c)
		if !ok {
			return
		}
		c.JSON(http.StatusOK, gin.H{"opts": opts})
	})

	// 测试基本参数
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test?page=1&page_size=20&status=pending&user_id=1&player_id=2&game_id=3&date_from=2024-01-01&date_to=2024-12-31&keyword=test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBuildPaymentListOptions(t *testing.T) {
	r := setupHelpersTestRouter()
	r.GET("/test", func(c *gin.Context) {
		opts, ok := buildPaymentListOptions(c)
		if !ok {
			return
		}
		c.JSON(http.StatusOK, gin.H{"opts": opts})
	})

	// 测试基本参数
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test?page=1&page_size=20&status=pending&method=alipay&user_id=1&order_id=2&date_from=2024-01-01&date_to=2024-12-31", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBuildUserListOptions(t *testing.T) {
	r := setupHelpersTestRouter()
	r.GET("/test", func(c *gin.Context) {
		opts, ok := buildUserListOptions(c)
		if !ok {
			return
		}
		c.JSON(http.StatusOK, gin.H{"opts": opts})
	})

	// 测试基本参数
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test?page=1&page_size=20&role=user&status=active&date_from=2024-01-01&date_to=2024-12-31&keyword=test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestEnsureSlice(t *testing.T) {
	// 测试nil切片
	var nilSlice []string
	result := ensureSlice(nilSlice)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)

	// 测试非nil切片
	nonNilSlice := []string{"a", "b"}
	result2 := ensureSlice(nonNilSlice)
	assert.Equal(t, nonNilSlice, result2)
}

func TestWriteJSON(t *testing.T) {
	r := setupHelpersTestRouter()
	r.GET("/test", func(c *gin.Context) {
		writeJSON(c, http.StatusOK, model.APIResponse[string]{
			Success: true,
			Code:    http.StatusOK,
			Message: "OK",
			Data:    "test",
		})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "test")
}

func TestWriteJSONError(t *testing.T) {
	r := setupHelpersTestRouter()
	r.GET("/test", func(c *gin.Context) {
		writeJSONError(c, http.StatusBadRequest, "test error")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "test error")
}

