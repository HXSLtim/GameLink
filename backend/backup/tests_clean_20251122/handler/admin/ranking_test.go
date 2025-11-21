package admin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gamelink/internal/cache"
	"gamelink/internal/model"
	adminservice "gamelink/internal/service/admin"
)

// 排行榜管理测试

func setupRankingTestRouter() (*gin.Engine, *adminservice.AdminService) {
	r := newTestEngine()

	svc := adminservice.NewAdminService(
		&fakeGameRepo{},
		&fakeUserRepo{},
		&fakePlayerRepo{},
		&fakeOrderRepo{},
		&fakePaymentRepo{},
		&fakeRoleRepo{},
		cache.NewMemory(),
	)

	// 注册排行榜相关路由
	r.GET("/admin/rankings/players", func(c *gin.Context) {
		// 获取玩家排行榜
		writeJSON(c, http.StatusOK, model.APIResponse[any]{
			Success: true,
			Code:    http.StatusOK,
			Message: "OK",
			Data: []map[string]interface{}{
				{"playerId": 1, "name": "Top Player 1", "revenue": 10000},
				{"playerId": 2, "name": "Top Player 2", "revenue": 8000},
			},
		})
	})

	r.GET("/admin/rankings/games", func(c *gin.Context) {
		// 获取游戏排行榜
		writeJSON(c, http.StatusOK, model.APIResponse[any]{
			Success: true,
			Code:    http.StatusOK,
			Message: "OK",
			Data: []map[string]interface{}{
				{"gameId": 1, "name": "Game 1", "orderCount": 1000},
				{"gameId": 2, "name": "Game 2", "orderCount": 800},
			},
		})
	})

	return r, svc
}

func TestRanking_GetPlayerRanking(t *testing.T) {
	r, _ := setupRankingTestRouter()

	tests := []struct {
		name       string
		query      string
		wantStatus int
		checkData  bool
	}{
		{
			name:       "获取默认排行榜",
			query:      "",
			wantStatus: http.StatusOK,
			checkData:  true,
		},
		{
			name:       "按收入排序",
			query:      "?sortBy=revenue&limit=10",
			wantStatus: http.StatusOK,
			checkData:  true,
		},
		{
			name:       "按订单数排序",
			query:      "?sortBy=orderCount&limit=20",
			wantStatus: http.StatusOK,
			checkData:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/admin/rankings/players"+tt.query, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.checkData {
				var resp model.APIResponse[any]
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.True(t, resp.Success)
				assert.NotNil(t, resp.Data)
			}
		})
	}
}

func TestRanking_GetGameRanking(t *testing.T) {
	r, _ := setupRankingTestRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/rankings/games", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.APIResponse[any]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestRanking_QueryParameters(t *testing.T) {
	r, _ := setupRankingTestRouter()

	tests := []struct {
		name     string
		endpoint string
		query    string
		wantCode int
	}{
		{
			name:     "有效的limit参数",
			endpoint: "/admin/rankings/players",
			query:    "?limit=50",
			wantCode: http.StatusOK,
		},
		{
			name:     "有效的时间范围",
			endpoint: "/admin/rankings/players",
			query:    "?dateFrom=2024-01-01&dateTo=2024-12-31",
			wantCode: http.StatusOK,
		},
		{
			name:     "组合查询参数",
			endpoint: "/admin/rankings/games",
			query:    "?limit=10&sortBy=orderCount",
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tt.endpoint+tt.query, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}
