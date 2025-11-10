package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	rankingrepo "gamelink/internal/repository/ranking"
)

// fakeRankingCommissionRepo 实现RankingCommissionRepository接口
type fakeRankingCommissionRepo struct {
	configs     []model.RankingCommissionConfig
	createError error
	getError    error
}

func (f *fakeRankingCommissionRepo) CreateConfig(ctx context.Context, config *model.RankingCommissionConfig) error {
	if f.createError != nil {
		return f.createError
	}
	if config.ID == 0 {
		config.ID = uint64(len(f.configs) + 1)
	}
	f.configs = append(f.configs, *config)
	return nil
}

func (f *fakeRankingCommissionRepo) GetConfig(ctx context.Context, id uint64) (*model.RankingCommissionConfig, error) {
	if f.getError != nil {
		return nil, f.getError
	}
	for i := range f.configs {
		if f.configs[i].ID == id {
			return &f.configs[i], nil
		}
	}
	return nil, repository.ErrNotFound
}

func (f *fakeRankingCommissionRepo) ListConfigs(ctx context.Context, opts rankingrepo.RankingCommissionConfigListOptions) ([]model.RankingCommissionConfig, int64, error) {
	return f.configs, int64(len(f.configs)), nil
}

func (f *fakeRankingCommissionRepo) UpdateConfig(ctx context.Context, config *model.RankingCommissionConfig) error {
	for i := range f.configs {
		if f.configs[i].ID == config.ID {
			f.configs[i] = *config
			return nil
		}
	}
	return repository.ErrNotFound
}

func (f *fakeRankingCommissionRepo) DeleteConfig(ctx context.Context, id uint64) error {
	for i := range f.configs {
		if f.configs[i].ID == id {
			f.configs = append(f.configs[:i], f.configs[i+1:]...)
			return nil
		}
	}
	return repository.ErrNotFound
}

func (f *fakeRankingCommissionRepo) GetActiveConfigForMonth(ctx context.Context, rankingType model.RankingType, month string) (*model.RankingCommissionConfig, error) {
	for i := range f.configs {
		if f.configs[i].RankingType == rankingType && f.configs[i].Month == month && f.configs[i].IsActive {
			return &f.configs[i], nil
		}
	}
	return nil, repository.ErrNotFound
}

// TestCreateRankingCommissionConfigHandler_WithActualCall 测试创建排名抽成配置
func TestCreateRankingCommissionConfigHandler_WithActualCall(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("成功创建配置", func(t *testing.T) {
		repo := &fakeRankingCommissionRepo{}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := CreateRankingCommissionConfigRequest{
			Name:        "2024年11月收入排名",
			RankingType: "income",
			Month:       "2024-11",
			Rules: []model.RankingCommissionRule{
				{RankStart: 1, RankEnd: 3, CommissionRate: 10},
				{RankStart: 4, RankEnd: 10, CommissionRate: 5},
			},
			Description: "月度收入排名奖励",
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest(http.MethodPost, "/admin/ranking-commission/configs", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		createRankingCommissionConfigHandler(c, repo)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.APIResponse[model.RankingCommissionConfig]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "2024年11月收入排名", response.Data.Name)
		assert.Equal(t, uint64(1), response.Data.ID)
	})

	t.Run("无效JSON返回400", func(t *testing.T) {
		repo := &fakeRankingCommissionRepo{}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/admin/ranking-commission/configs", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")

		createRankingCommissionConfigHandler(c, repo)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("无效规则返回400", func(t *testing.T) {
		repo := &fakeRankingCommissionRepo{}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := CreateRankingCommissionConfigRequest{
			Name:        "无效配置",
			RankingType: "income",
			Month:       "2024-11",
			Rules: []model.RankingCommissionRule{
				{RankStart: 5, RankEnd: 3, CommissionRate: 10}, // rankStart > rankEnd
			},
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest(http.MethodPost, "/admin/ranking-commission/configs", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		createRankingCommissionConfigHandler(c, repo)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestListRankingCommissionConfigsHandler_WithActualCall 测试获取配置列表
func TestListRankingCommissionConfigsHandler_WithActualCall(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("成功获取列表", func(t *testing.T) {
		rulesJSON1, _ := json.Marshal([]model.RankingCommissionRule{
			{RankStart: 1, RankEnd: 3, CommissionRate: 10},
		})
		rulesJSON2, _ := json.Marshal([]model.RankingCommissionRule{
			{RankStart: 1, RankEnd: 5, CommissionRate: 8},
		})

		repo := &fakeRankingCommissionRepo{
			configs: []model.RankingCommissionConfig{
				{
					ID:          1,
					Name:        "配置1",
					RankingType: "income",
					Month:       "2024-11",
					RulesJSON:   string(rulesJSON1),
					IsActive:    true,
				},
				{
					ID:          2,
					Name:        "配置2",
					RankingType: "order_count",
					Month:       "2024-11",
					RulesJSON:   string(rulesJSON2),
					IsActive:    true,
				},
			},
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodGet, "/admin/ranking-commission/configs?page=1&pageSize=10", nil)

		listRankingCommissionConfigsHandler(c, repo)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.APIResponse[any]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
	})

	t.Run("带筛选条件获取列表", func(t *testing.T) {
		repo := &fakeRankingCommissionRepo{configs: []model.RankingCommissionConfig{}}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodGet, "/admin/ranking-commission/configs?month=2024-11&rankingType=income", nil)

		listRankingCommissionConfigsHandler(c, repo)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

// TestGetRankingCommissionConfigHandler_WithActualCall 测试获取配置详情
func TestGetRankingCommissionConfigHandler_WithActualCall(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("成功获取配置", func(t *testing.T) {
		rulesJSON, _ := json.Marshal([]model.RankingCommissionRule{
			{RankStart: 1, RankEnd: 3, CommissionRate: 100000},
		})

		repo := &fakeRankingCommissionRepo{
			configs: []model.RankingCommissionConfig{
				{
					ID:          1,
					Name:        "测试配置",
					RankingType: "income",
					Month:       "2024-11",
					RulesJSON:   string(rulesJSON),
					IsActive:    true,
				},
			},
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodGet, "/admin/ranking-commission/configs/1", nil)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		getRankingCommissionConfigHandler(c, repo)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.APIResponse[any]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
	})

	t.Run("无效ID返回400", func(t *testing.T) {
		repo := &fakeRankingCommissionRepo{}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodGet, "/admin/ranking-commission/configs/invalid", nil)
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}

		getRankingCommissionConfigHandler(c, repo)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("配置不存在返回404", func(t *testing.T) {
		repo := &fakeRankingCommissionRepo{configs: []model.RankingCommissionConfig{}}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodGet, "/admin/ranking-commission/configs/999", nil)
		c.Params = gin.Params{{Key: "id", Value: "999"}}

		getRankingCommissionConfigHandler(c, repo)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// TestUpdateRankingCommissionConfigHandler_WithActualCall 测试更新配置
func TestUpdateRankingCommissionConfigHandler_WithActualCall(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("成功更新配置", func(t *testing.T) {
		rulesJSON, _ := json.Marshal([]model.RankingCommissionRule{
			{RankStart: 1, RankEnd: 3, CommissionRate: 100000},
		})

		repo := &fakeRankingCommissionRepo{
			configs: []model.RankingCommissionConfig{
				{
					ID:          1,
					Name:        "旧名称",
					RankingType: "income",
					Month:       "2024-11",
					RulesJSON:   string(rulesJSON),
					IsActive:    true,
				},
			},
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		newName := "新名称"
		isActive := false
		reqBody := UpdateRankingCommissionConfigRequest{
			Name:     &newName,
			IsActive: &isActive,
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest(http.MethodPut, "/admin/ranking-commission/configs/1", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		updateRankingCommissionConfigHandler(c, repo)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("无效ID返回400", func(t *testing.T) {
		repo := &fakeRankingCommissionRepo{}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := UpdateRankingCommissionConfigRequest{}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest(http.MethodPut, "/admin/ranking-commission/configs/invalid", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}

		updateRankingCommissionConfigHandler(c, repo)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("配置不存在返回404", func(t *testing.T) {
		repo := &fakeRankingCommissionRepo{configs: []model.RankingCommissionConfig{}}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := UpdateRankingCommissionConfigRequest{}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest(http.MethodPut, "/admin/ranking-commission/configs/999", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: "999"}}

		updateRankingCommissionConfigHandler(c, repo)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("更新无效规则返回400", func(t *testing.T) {
		rulesJSON, _ := json.Marshal([]model.RankingCommissionRule{
			{RankStart: 1, RankEnd: 3, CommissionRate: 100000},
		})

		repo := &fakeRankingCommissionRepo{
			configs: []model.RankingCommissionConfig{
				{
					ID:          1,
					Name:        "测试",
					RankingType: "income",
					Month:       "2024-11",
					RulesJSON:   string(rulesJSON),
					IsActive:    true,
				},
			},
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		invalidRules := []model.RankingCommissionRule{
			{RankStart: 5, RankEnd: 3, CommissionRate: 10}, // invalid
		}
		reqBody := UpdateRankingCommissionConfigRequest{
			Rules: &invalidRules,
		}
		body, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest(http.MethodPut, "/admin/ranking-commission/configs/1", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		updateRankingCommissionConfigHandler(c, repo)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestDeleteRankingCommissionConfigHandler_WithActualCall 测试删除配置
func TestDeleteRankingCommissionConfigHandler_WithActualCall(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("成功删除配置", func(t *testing.T) {
		repo := &fakeRankingCommissionRepo{
			configs: []model.RankingCommissionConfig{
				{ID: 1, Name: "测试配置"},
			},
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodDelete, "/admin/ranking-commission/configs/1", nil)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		deleteRankingCommissionConfigHandler(c, repo)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, 0, len(repo.configs)) // 验证已删除
	})

	t.Run("无效ID返回400", func(t *testing.T) {
		repo := &fakeRankingCommissionRepo{}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodDelete, "/admin/ranking-commission/configs/invalid", nil)
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}

		deleteRankingCommissionConfigHandler(c, repo)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestRegisterRankingCommissionRoutes_Coverage 测试路由注册
func TestRegisterRankingCommissionRoutes_Coverage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	repo := &fakeRankingCommissionRepo{}

	RegisterRankingCommissionRoutes(router, repo)

	routes := router.Routes()
	assert.NotEmpty(t, routes)

	routeMap := make(map[string]bool)
	for _, route := range routes {
		routeMap[route.Method+":"+route.Path] = true
	}

	assert.True(t, routeMap["POST:/admin/ranking-commission/configs"])
	assert.True(t, routeMap["GET:/admin/ranking-commission/configs"])
	assert.True(t, routeMap["GET:/admin/ranking-commission/configs/:id"])
	assert.True(t, routeMap["PUT:/admin/ranking-commission/configs/:id"])
	assert.True(t, routeMap["DELETE:/admin/ranking-commission/configs/:id"])
}
