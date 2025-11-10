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

func TestGetDashboardOverview_EdgeCases(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("空数据库", func(t *testing.T) {
		router := gin.New()
		
		userRepo := &dashboardUserRepo{users: []model.User{}}
		playerRepo := &fakePlayerRepoForHandler{
			listPaged: func(page, size int) ([]model.Player, int64, error) {
				return []model.Player{}, 0, nil
			},
		}
		orderRepo := &fakeOrderRepoForHandler{items: []model.Order{}}
		withdrawRepo := &dashboardWithdrawRepo{}
		serviceItemRepo := &dashboardServiceItemRepo{}
		
		router.GET("/dashboard", func(c *gin.Context) {
			getDashboardOverviewHandler(c, userRepo, playerRepo, orderRepo, withdrawRepo, serviceItemRepo)
		})

		req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var resp model.APIResponse[DashboardOverviewStats]
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, int64(0), resp.Data.TotalUsers)
		assert.Equal(t, int64(0), resp.Data.TotalPlayers)
	})

	t.Run("大量数据", func(t *testing.T) {
		router := gin.New()
		
		// 创建1000个用户
		users := make([]model.User, 1000)
		for i := 0; i < 1000; i++ {
			users[i] = model.User{Base: model.Base{ID: uint64(i + 1)}}
		}
		
		userRepo := &dashboardUserRepo{users: users}
		playerRepo := &fakePlayerRepoForHandler{
			listPaged: func(page, size int) ([]model.Player, int64, error) {
				return []model.Player{}, 500, nil
			},
		}
		orderRepo := &fakeOrderRepoForHandler{items: []model.Order{}}
		withdrawRepo := &dashboardWithdrawRepo{}
		serviceItemRepo := &dashboardServiceItemRepo{}
		
		router.GET("/dashboard", func(c *gin.Context) {
			getDashboardOverviewHandler(c, userRepo, playerRepo, orderRepo, withdrawRepo, serviceItemRepo)
		})

		req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var resp model.APIResponse[DashboardOverviewStats]
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, int64(1000), resp.Data.TotalUsers)
		assert.Equal(t, int64(500), resp.Data.TotalPlayers)
	})
}

func TestGetRecentOrders_EdgeCases(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("limit=0", func(t *testing.T) {
		router := gin.New()
		orderRepo := &fakeOrderRepoForHandler{items: []model.Order{}}
		
		router.GET("/orders", func(c *gin.Context) {
			getRecentOrdersHandler(c, orderRepo)
		})

		req := httptest.NewRequest(http.MethodGet, "/orders?limit=0", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("limit=负数", func(t *testing.T) {
		router := gin.New()
		orderRepo := &fakeOrderRepoForHandler{items: []model.Order{}}
		
		router.GET("/orders", func(c *gin.Context) {
			getRecentOrdersHandler(c, orderRepo)
		})

		req := httptest.NewRequest(http.MethodGet, "/orders?limit=-1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("limit=101超过最大值", func(t *testing.T) {
		router := gin.New()
		orderRepo := &fakeOrderRepoForHandler{items: []model.Order{}}
		
		router.GET("/orders", func(c *gin.Context) {
			getRecentOrdersHandler(c, orderRepo)
		})

		req := httptest.NewRequest(http.MethodGet, "/orders?limit=101", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("limit=非数字", func(t *testing.T) {
		router := gin.New()
		orderRepo := &fakeOrderRepoForHandler{items: []model.Order{}}
		
		router.GET("/orders", func(c *gin.Context) {
			getRecentOrdersHandler(c, orderRepo)
		})

		req := httptest.NewRequest(http.MethodGet, "/orders?limit=abc", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestGetRecentWithdraws_EdgeCases(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("limit参数测试", func(t *testing.T) {
		tests := []struct {
			name  string
			limit string
		}{
			{"默认值", ""},
			{"有效值", "20"},
			{"最大值", "100"},
			{"超过最大值", "200"},
			{"零", "0"},
			{"负数", "-5"},
			{"非数字", "xyz"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				router := gin.New()
				withdrawRepo := &dashboardWithdrawRepo{}
				
				router.GET("/withdraws", func(c *gin.Context) {
					getRecentWithdrawsHandler(c, withdrawRepo)
				})

				url := "/withdraws"
				if tt.limit != "" {
					url += "?limit=" + tt.limit
				}
				
				req := httptest.NewRequest(http.MethodGet, url, nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				assert.Equal(t, http.StatusOK, w.Code)
			})
		}
	})
}

func TestGetMonthlyRevenue_EdgeCases(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("months参数测试", func(t *testing.T) {
		tests := []struct {
			name   string
			months string
		}{
			{"默认12个月", ""},
			{"1个月", "1"},
			{"6个月", "6"},
			{"最大24个月", "24"},
			{"超过最大值", "30"},
			{"零", "0"},
			{"负数", "-1"},
			{"非数字", "abc"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				router := gin.New()
				commissionRepo := &dashboardCommissionRepo{}
				
				router.GET("/revenue", func(c *gin.Context) {
					getMonthlyRevenueHandler(c, commissionRepo)
				})

				url := "/revenue"
				if tt.months != "" {
					url += "?months=" + tt.months
				}
				
				req := httptest.NewRequest(http.MethodGet, url, nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				assert.Equal(t, http.StatusOK, w.Code)
				
				var resp model.APIResponse[map[string]interface{}]
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.True(t, resp.Success)
			})
		}
	})
}

// Mock repositories are defined in dashboard_test.go
