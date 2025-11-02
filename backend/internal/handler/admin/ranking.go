package admin

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	commissionservice "gamelink/internal/service/commission"
)

// RegisterRankingCommissionRoutes 注册管理端排名抽成配置路�?
func RegisterRankingCommissionRoutes(router gin.IRouter, repo repository.RankingCommissionRepository) {
	group := router.Group("/admin/ranking-commission")
	{
		group.POST("/configs", func(c *gin.Context) { createRankingCommissionConfigHandler(c, repo) })
		group.GET("/configs", func(c *gin.Context) { listRankingCommissionConfigsHandler(c, repo) })
		group.GET("/configs/:id", func(c *gin.Context) { getRankingCommissionConfigHandler(c, repo) })
		group.PUT("/configs/:id", func(c *gin.Context) { updateRankingCommissionConfigHandler(c, repo) })
		group.DELETE("/configs/:id", func(c *gin.Context) { deleteRankingCommissionConfigHandler(c, repo) })
	}
}

// CreateRankingCommissionConfigRequest 创建排名抽成配置请求
type CreateRankingCommissionConfigRequest struct {
	Name        string                              `json:"name" binding:"required"`
	RankingType model.RankingType                   `json:"rankingType" binding:"required,oneof=income order_count"`
	Month       string                              `json:"month" binding:"required"` // YYYY-MM
	Rules       []model.RankingCommissionRule       `json:"rules" binding:"required,min=1"`
	Description string                              `json:"description"`
}

// createRankingCommissionConfigHandler 创建排名抽成配置
// @Summary      创建排名抽成配置
// @Description  管理员配置排名抽成规则（支持阶梯抽成�?
// @Tags         Admin - RankingCommission
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                                   true  "Bearer {token}"
// @Param        request        body      CreateRankingCommissionConfigRequest  true  "配置信息"
// @Success      200            {object}  model.APIResponse[model.RankingCommissionConfig]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /admin/ranking-commission/configs [post]
func createRankingCommissionConfigHandler(c *gin.Context, repo repository.RankingCommissionRepository) {
	var req CreateRankingCommissionConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	// 验证规则
	if err := commissionservice.ValidateRankingRules(req.Rules); err != nil {
		respondError(c, http.StatusBadRequest, "Invalid rules: "+err.Error())
		return
	}

	// 序列化规�?
	rulesJSON, err := json.Marshal(req.Rules)
	if err != nil {
		respondError(c, http.StatusBadRequest, "Failed to serialize rules")
		return
	}

	// 创建配置
	config := &model.RankingCommissionConfig{
		Name:        req.Name,
		RankingType: req.RankingType,
		Period:      "monthly",
		Month:       req.Month,
		RulesJSON:   string(rulesJSON),
		Description: req.Description,
		IsActive:    true,
	}

	if err := repo.CreateConfig(c.Request.Context(), config); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[model.RankingCommissionConfig]{
		Success: true,
		Code:    http.StatusOK,
		Message: "Ranking commission config created successfully",
		Data:    *config,
	})
}

// listRankingCommissionConfigsHandler 获取排名抽成配置列表
// @Summary      获取排名抽成配置列表
// @Description  管理员查看所有排名抽成配�?
// @Tags         Admin - RankingCommission
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer {token}"
// @Param        month          query     string  false  "月份筛�?
// @Param        rankingType    query     string  false  "排名类型"
// @Param        page           query     int     false  "页码"
// @Param        pageSize       query     int     false  "每页数量"
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /admin/ranking-commission/configs [get]
func listRankingCommissionConfigsHandler(c *gin.Context, repo repository.RankingCommissionRepository) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	opts := repository.RankingCommissionConfigListOptions{
		Page:     page,
		PageSize: pageSize,
	}

	if month := c.Query("month"); month != "" {
		opts.Month = &month
	}

	if rankingTypeStr := c.Query("rankingType"); rankingTypeStr != "" {
		rankingType := model.RankingType(rankingTypeStr)
		opts.RankingType = &rankingType
	}

	configs, total, err := repo.ListConfigs(c.Request.Context(), opts)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 解析每个配置的规�?
	type ConfigDTO struct {
		model.RankingCommissionConfig
		Rules []model.RankingCommissionRule `json:"rules"`
	}

	configDTOs := make([]ConfigDTO, 0, len(configs))
	for _, config := range configs {
		var rules []model.RankingCommissionRule
		json.Unmarshal([]byte(config.RulesJSON), &rules)

		configDTOs = append(configDTOs, ConfigDTO{
			RankingCommissionConfig: config,
			Rules:                   rules,
		})
	}

	respondJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data: map[string]interface{}{
			"configs": configDTOs,
			"total":   total,
		},
	})
}

// getRankingCommissionConfigHandler 获取排名抽成配置详情
// @Summary      获取排名抽成配置详情
// @Description  管理员查看排名抽成配置详�?
// @Tags         Admin - RankingCommission
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        id             path      int     true  "配置ID"
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /admin/ranking-commission/configs/{id} [get]
func getRankingCommissionConfigHandler(c *gin.Context, repo repository.RankingCommissionRepository) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, "Invalid config ID")
		return
	}

	config, err := repo.GetConfig(c.Request.Context(), id)
	if err != nil {
		respondError(c, http.StatusNotFound, "Config not found")
		return
	}

	// 解析规则
	var rules []model.RankingCommissionRule
	json.Unmarshal([]byte(config.RulesJSON), &rules)

	respondJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data: map[string]interface{}{
			"config": config,
			"rules":  rules,
		},
	})
}

// UpdateRankingCommissionConfigRequest 更新配置请求
type UpdateRankingCommissionConfigRequest struct {
	Name        *string                            `json:"name"`
	Rules       *[]model.RankingCommissionRule     `json:"rules"`
	Description *string                            `json:"description"`
	IsActive    *bool                              `json:"isActive"`
}

// updateRankingCommissionConfigHandler 更新排名抽成配置
// @Summary      更新排名抽成配置
// @Description  管理员更新排名抽成配�?
// @Tags         Admin - RankingCommission
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                                   true  "Bearer {token}"
// @Param        id             path      int                                      true  "配置ID"
// @Param        request        body      UpdateRankingCommissionConfigRequest  true  "更新信息"
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /admin/ranking-commission/configs/{id} [put]
func updateRankingCommissionConfigHandler(c *gin.Context, repo repository.RankingCommissionRepository) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, "Invalid config ID")
		return
	}

	var req UpdateRankingCommissionConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	config, err := repo.GetConfig(c.Request.Context(), id)
	if err != nil {
		respondError(c, http.StatusNotFound, "Config not found")
		return
	}

	// 更新字段
	if req.Name != nil {
		config.Name = *req.Name
	}
	if req.Description != nil {
		config.Description = *req.Description
	}
	if req.IsActive != nil {
		config.IsActive = *req.IsActive
	}
	if req.Rules != nil {
		// 验证规则
		if err := commissionservice.ValidateRankingRules(*req.Rules); err != nil {
			respondError(c, http.StatusBadRequest, "Invalid rules: "+err.Error())
			return
		}

		// 序列化规�?
		rulesJSON, err := json.Marshal(*req.Rules)
		if err != nil {
			respondError(c, http.StatusBadRequest, "Failed to serialize rules")
			return
		}
		config.RulesJSON = string(rulesJSON)
	}

	if err := repo.UpdateConfig(c.Request.Context(), config); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "Config updated successfully",
	})
}

// deleteRankingCommissionConfigHandler 删除排名抽成配置
// @Summary      删除排名抽成配置
// @Description  管理员删除排名抽成配�?
// @Tags         Admin - RankingCommission
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        id             path      int     true  "配置ID"
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /admin/ranking-commission/configs/{id} [delete]
func deleteRankingCommissionConfigHandler(c *gin.Context, repo repository.RankingCommissionRepository) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, "Invalid config ID")
		return
	}

	if err := repo.DeleteConfig(c.Request.Context(), id); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "Config deleted successfully",
	})
}

