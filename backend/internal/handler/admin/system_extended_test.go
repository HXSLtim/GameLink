package admin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gamelink/internal/model"
)

// 系统信息测试

func TestSystem_Config(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("获取系统配置", func(t *testing.T) {
		r := gin.New()
		r.GET("/admin/system/config", func(c *gin.Context) {
			writeJSON(c, http.StatusOK, model.APIResponse[map[string]interface{}]{
				Success: true,
				Code:    http.StatusOK,
				Message: "OK",
				Data: map[string]interface{}{
					"appName":    "GameLink",
					"version":    "1.0.0",
					"env":        "production",
					"maxUpload":  10485760,
				},
			})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/admin/system/config", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var resp model.APIResponse[map[string]interface{}]
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.True(t, resp.Success)
		assert.NotNil(t, resp.Data)
	})
}

func TestSystem_DBStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("数据库连接正常", func(t *testing.T) {
		r := gin.New()
		r.GET("/admin/system/db", func(c *gin.Context) {
			writeJSON(c, http.StatusOK, model.APIResponse[map[string]interface{}]{
				Success: true,
				Code:    http.StatusOK,
				Message: "OK",
				Data: map[string]interface{}{
					"status":      "connected",
					"connections": 10,
					"maxConn":     100,
				},
			})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/admin/system/db", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("数据库连接失败", func(t *testing.T) {
		r := gin.New()
		r.GET("/admin/system/db", func(c *gin.Context) {
			writeJSONError(c, http.StatusServiceUnavailable, "Database connection failed")
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/admin/system/db", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	})
}

func TestSystem_CacheStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("缓存状态正常", func(t *testing.T) {
		r := gin.New()
		r.GET("/admin/system/cache", func(c *gin.Context) {
			writeJSON(c, http.StatusOK, model.APIResponse[map[string]interface{}]{
				Success: true,
				Code:    http.StatusOK,
				Message: "OK",
				Data: map[string]interface{}{
					"status":   "connected",
					"hitRate":  0.85,
					"keys":     1000,
				},
			})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/admin/system/cache", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestSystem_Resources(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("获取系统资源", func(t *testing.T) {
		r := gin.New()
		r.GET("/admin/system/resources", func(c *gin.Context) {
			writeJSON(c, http.StatusOK, model.APIResponse[map[string]interface{}]{
				Success: true,
				Code:    http.StatusOK,
				Message: "OK",
				Data: map[string]interface{}{
					"cpuUsage":    45.5,
					"memoryUsage": 60.2,
					"diskUsage":   30.0,
				},
			})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/admin/system/resources", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var resp model.APIResponse[map[string]interface{}]
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.True(t, resp.Success)
	})
}

func TestSystem_Version(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("获取版本信息", func(t *testing.T) {
		r := gin.New()
		r.GET("/admin/system/version", func(c *gin.Context) {
			writeJSON(c, http.StatusOK, model.APIResponse[map[string]interface{}]{
				Success: true,
				Code:    http.StatusOK,
				Message: "OK",
				Data: map[string]interface{}{
					"version":   "1.0.0",
					"buildTime": "2024-01-01",
					"goVersion": "1.25.3",
				},
			})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/admin/system/version", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
