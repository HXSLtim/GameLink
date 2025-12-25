package admin

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	rankingrepo "gamelink/internal/repository/ranking"
	commissionservice "gamelink/internal/service/commission"
	apierr "gamelink/pkg/apierr"
)

// RankingCommissionConfig 排名抽成配置模型（类型别名）
type RankingCommissionConfig = model.RankingCommissionConfig

// RegisterRankingCommissionRoutes 注册管理端排名抽成配置路由
func RegisterRankingCommissionRoutes(router gin.IRouter, repo rankingrepo.RankingCommissionRepository) {
	group := router.Group("/ranking-commission")
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
	Name        string                        `json:"name" binding:"required"`
	RankingType model.RankingType             `json:"rankingType" binding:"required,oneof=income order_count"`
	Month       string                        `json:"month" binding:"required"` // YYYY-MM
	Rules       []model.RankingCommissionRule `json:"rules" binding:"required,min=1"`
	Description string                        `json:"description"`
}

// createRankingCommissionConfigHandler 创建排名抽成配置
// @Summary      创建排名抽成配置
// @Description  API endpoint// @Tags         Admin - RankingCommission
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                                   true  "Bearer {token}"
// @Param        request        body      CreateRankingCommissionConfigRequest  true  "配置信息"
// @Success      200            {object}  model.APIResponse[RankingCommissionConfig]
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Router       /admin/ranking-commission/configs [post]
func createRankingCommissionConfigHandler(c *gin.Context, repo rankingrepo.RankingCommissionRepository) {
	var req CreateRankingCommissionConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	// 验证规则
	if err := commissionservice.ValidateRankingRules(req.Rules); err != nil {
		respondBadRequest(c, "Invalid rules: "+err.Error())
		return
	}

	// 序列化规
	rulesJSON, err := json.Marshal(req.Rules)
	if err != nil {
		respondBadRequest(c, "Failed to serialize rules")
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
		respondError(c, err)
		return
	}

	respondSuccessWithMsg(c, "Ranking commission config created successfully", *config)
}

// listRankingCommissionConfigsHandler 获取排名抽成配置列表
// @Summary      获取排名抽成配置列表
// @Description  API endpoint// @Tags         Admin - RankingCommission
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer {token}"
// @Param        month          query    string       false  "Month filter (YYYY-MM)"// @Param        rankingType    query     string  false  "排名类型"
// @Param        page           query     int     false  "页码"
// @Param        pageSize       query     int     false  "每页数量"
// @Success      200            {object}  model.SuccessResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Router       /admin/ranking-commission/configs [get]
func listRankingCommissionConfigsHandler(c *gin.Context, repo rankingrepo.RankingCommissionRepository) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	opts := rankingrepo.RankingCommissionConfigListOptions{
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
		respondError(c, err)
		return
	}

	// 解析每个配置的规
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

	respondSuccess(c, map[string]interface{}{
		"configs": configDTOs,
		"total":   total,
	})
}

// getRankingCommissionConfigHandler 获取排名抽成配置详情
// @Summary      获取排名抽成配置详情
// @Description  API endpoint// @Tags         Admin - RankingCommission
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        id             path      int     true  "配置ID"
// @Success      200            {object}  model.SuccessResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Router       /admin/ranking-commission/configs/{id} [get]
func getRankingCommissionConfigHandler(c *gin.Context, repo rankingrepo.RankingCommissionRepository) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	config, err := repo.GetConfig(c.Request.Context(), id)
	if err != nil {
		respondError(c, apierr.NotFound(apierr.ErrRankingConfigNotFound))
		return
	}

	// 解析规则
	var rules []model.RankingCommissionRule
	json.Unmarshal([]byte(config.RulesJSON), &rules)

	respondSuccess(c, map[string]interface{}{
		"config": config,
		"rules":  rules,
	})
}

// UpdateRankingCommissionConfigRequest 更新配置请求
type UpdateRankingCommissionConfigRequest struct {
	Name        *string                        `json:"name"`
	Rules       *[]model.RankingCommissionRule `json:"rules"`
	Description *string                        `json:"description"`
	IsActive    *bool                          `json:"isActive"`
}

// updateRankingCommissionConfigHandler 更新排名抽成配置
// @Summary      更新排名抽成配置
// @Description  API endpoint// @Tags         Admin - RankingCommission
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                                   true  "Bearer {token}"
// @Param        id             path      int                                      true  "配置ID"
// @Param        request        body      UpdateRankingCommissionConfigRequest  true  "更新信息"
// @Success      200            {object}  model.SuccessResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Router       /admin/ranking-commission/configs/{id} [put]
func updateRankingCommissionConfigHandler(c *gin.Context, repo rankingrepo.RankingCommissionRepository) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var req UpdateRankingCommissionConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	config, err := repo.GetConfig(c.Request.Context(), id)
	if err != nil {
		respondError(c, apierr.NotFound("Config not found"))
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
			respondBadRequest(c, "Invalid rules: "+err.Error())
			return
		}

		// 序列化规
		rulesJSON, err := json.Marshal(*req.Rules)
		if err != nil {
			respondBadRequest(c, "Failed to serialize rules")
			return
		}
		config.RulesJSON = string(rulesJSON)
	}

	if err := repo.UpdateConfig(c.Request.Context(), config); err != nil {
		respondError(c, err)
		return
	}

	respondMsg(c, "Config updated successfully")
}

// deleteRankingCommissionConfigHandler 删除排名抽成配置
// @Summary      删除排名抽成配置
// @Description  API endpoint// @Tags         Admin - RankingCommission
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        id             path      int     true  "配置ID"
// @Success      200            {object}  model.SuccessResponse
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Router       /admin/ranking-commission/configs/{id} [delete]
func deleteRankingCommissionConfigHandler(c *gin.Context, repo rankingrepo.RankingCommissionRepository) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	if err := repo.DeleteConfig(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}

	respondMsg(c, "Config deleted successfully")
}
