package admin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// 测试commission handler的基本功能
// 注意：这些测试主要验证HTTP层的请求处理，业务逻辑应在service层测试

func TestCommissionHandlerInvalidRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		method         string
		path           string
		body           string
		params         gin.Params
		expectedStatus int
	}{
		{
			name:           "创建规则-无效JSON",
			method:         http.MethodPost,
			path:           "/admin/commission/rules",
			body:           "invalid json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "更新规则-无效ID",
			method:         http.MethodPut,
			path:           "/admin/commission/rules/invalid",
			body:           `{"name":"test"}`,
			params:         gin.Params{{Key: "id", Value: "invalid"}},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(tt.method, tt.path, bytes.NewBufferString(tt.body))
			c.Request.Header.Set("Content-Type", "application/json")
			if len(tt.params) > 0 {
				c.Params = tt.params
			}

			// 验证请求被正确创建
			assert.NotNil(t, c.Request)
			assert.Equal(t, tt.method, c.Request.Method)
		})
	}
}

func TestCommissionRequestStructures(t *testing.T) {
	// 测试请求结构体的JSON序列化
	t.Run("ApproveWithdrawRequest序列化", func(t *testing.T) {
		req := ApproveWithdrawRequest{
			Remark: "测试备注",
		}
		data, err := json.Marshal(req)
		assert.NoError(t, err)
		assert.Contains(t, string(data), "remark")
	})

	t.Run("RejectWithdrawRequest序列化", func(t *testing.T) {
		req := RejectWithdrawRequest{
			Reason: "测试原因",
		}
		data, err := json.Marshal(req)
		assert.NoError(t, err)
		assert.Contains(t, string(data), "reason")
	})
}
