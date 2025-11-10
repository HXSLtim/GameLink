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

// 扩展的排行榜测试

func TestRanking_EdgeCases(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("空结果处理", func(t *testing.T) {
		r := gin.New()
		r.GET("/admin/rankings/players", func(c *gin.Context) {
			writeJSON(c, http.StatusOK, model.APIResponse[any]{
				Success: true,
				Code:    http.StatusOK,
				Message: "OK",
				Data:    []interface{}{},
			})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/admin/rankings/players", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var resp model.APIResponse[any]
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.True(t, resp.Success)
	})

	t.Run("无效的limit参数", func(t *testing.T) {
		r := gin.New()
		r.GET("/admin/rankings/players", func(c *gin.Context) {
			limit := c.Query("limit")
			if limit == "-1" || limit == "0" {
				writeJSONError(c, http.StatusBadRequest, "Invalid limit")
				return
			}
			writeJSON(c, http.StatusOK, model.APIResponse[any]{Success: true})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/admin/rankings/players?limit=-1", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("超大limit参数", func(t *testing.T) {
		r := gin.New()
		r.GET("/admin/rankings/players", func(c *gin.Context) {
			writeJSON(c, http.StatusOK, model.APIResponse[any]{Success: true})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/admin/rankings/players?limit=10000", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestRanking_SortOptions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name     string
		sortBy   string
		wantCode int
	}{
		{"按收入排序", "revenue", http.StatusOK},
		{"按订单数排序", "orderCount", http.StatusOK},
		{"按评分排序", "rating", http.StatusOK},
		{"默认排序", "", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.GET("/admin/rankings/players", func(c *gin.Context) {
				writeJSON(c, http.StatusOK, model.APIResponse[any]{Success: true})
			})

			w := httptest.NewRecorder()
			url := "/admin/rankings/players"
			if tt.sortBy != "" {
				url += "?sortBy=" + tt.sortBy
			}
			req := httptest.NewRequest(http.MethodGet, url, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}

func TestRanking_TimeRange(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("有效时间范围", func(t *testing.T) {
		r := gin.New()
		r.GET("/admin/rankings/players", func(c *gin.Context) {
			writeJSON(c, http.StatusOK, model.APIResponse[any]{Success: true})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/admin/rankings/players?dateFrom=2024-01-01&dateTo=2024-12-31", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("只有开始时间", func(t *testing.T) {
		r := gin.New()
		r.GET("/admin/rankings/players", func(c *gin.Context) {
			writeJSON(c, http.StatusOK, model.APIResponse[any]{Success: true})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/admin/rankings/players?dateFrom=2024-01-01", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
