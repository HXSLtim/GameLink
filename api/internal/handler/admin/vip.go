package admin

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/middleware"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	vipservice "gamelink/internal/service/vip"
	"gamelink/pkg/apierr"
)

// VipHandler VIP管理处理器
type VipHandler struct {
	svc *vipservice.Service
}

// NewVipHandler 创建VIP处理器
func NewVipHandler(svc *vipservice.Service) *VipHandler {
	return &VipHandler{svc: svc}
}

// ============================================================================
// VIP等级管理
// ============================================================================

// VipCreateLevelRequest 创建VIP等级请求
type VipCreateLevelRequest struct {
	Slug                    string  `json:"slug" binding:"required,max=32"`
	Title                   string  `json:"title" binding:"required,max=64"`
	ExpRequired             int64   `json:"expRequired"`
	OrderDiscount           float64 `json:"orderDiscount"`
	MonthlyCouponTemplateID *uint64 `json:"monthlyCouponTemplateId"`
	MonthlyCouponCount      int     `json:"monthlyCouponCount"`
	IconURL                 string  `json:"iconUrl"`
	Color                   string  `json:"color"`
	Benefits                string  `json:"benefits"`
	SortOrder               int     `json:"sortOrder"`
	IsDefault               bool    `json:"isDefault"`
	IsActive                bool    `json:"isActive"`
}

// VipUpdateLevelRequest 更新VIP等级请求
type VipUpdateLevelRequest struct {
	Slug                    string  `json:"slug" binding:"required,max=32"`
	Title                   string  `json:"title" binding:"required,max=64"`
	ExpRequired             int64   `json:"expRequired"`
	OrderDiscount           float64 `json:"orderDiscount"`
	MonthlyCouponTemplateID *uint64 `json:"monthlyCouponTemplateId"`
	MonthlyCouponCount      int     `json:"monthlyCouponCount"`
	IconURL                 string  `json:"iconUrl"`
	Color                   string  `json:"color"`
	Benefits                string  `json:"benefits"`
	SortOrder               int     `json:"sortOrder"`
	IsDefault               bool    `json:"isDefault"`
	IsActive                bool    `json:"isActive"`
}

// ListLevels 获取VIP等级列表
func (h *VipHandler) ListLevels(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}
	keyword := c.Query("keyword")

	// 解析 isActive 参数
	var isActive *bool
	if v := strings.TrimSpace(c.Query("isActive")); v != "" {
		b := v == "true" || v == "1"
		isActive = &b
	}

	opts := repository.VipLevelListOptions{
		Page:     page,
		PageSize: pageSize,
		Keyword:  keyword,
		IsActive: isActive,
	}

	levels, total, err := h.svc.ListLevelsPaged(c.Request.Context(), opts)
	if err != nil {
		respondError(c, apierr.InternalError("获取VIP等级列表失败").WithDetails(err.Error()))
		return
	}

	respondList(c, levels, &model.Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
	})
}

// GetLevel 获取VIP等级详情
func (h *VipHandler) GetLevel(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	level, err := h.svc.GetLevel(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("VIP等级不存在"))
			return
		}
		respondError(c, apierr.InternalError("获取VIP等级失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, level)
}

// CreateLevel 创建VIP等级
func (h *VipHandler) CreateLevel(c *gin.Context) {
	var req VipCreateLevelRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	level := &model.VipLevel{
		Slug:                    req.Slug,
		Title:                   req.Title,
		ExpRequired:             req.ExpRequired,
		OrderDiscount:           req.OrderDiscount,
		MonthlyCouponTemplateID: req.MonthlyCouponTemplateID,
		MonthlyCouponCount:      req.MonthlyCouponCount,
		IconURL:                 req.IconURL,
		Color:                   req.Color,
		Benefits:                req.Benefits,
		SortOrder:               req.SortOrder,
		IsDefault:               req.IsDefault,
		IsActive:                req.IsActive,
	}

	if err := h.svc.CreateLevel(c.Request.Context(), level); err != nil {
		respondError(c, apierr.InternalError("创建VIP等级失败").WithDetails(err.Error()))
		return
	}

	respondCreated(c, level)
}

// UpdateLevel 更新VIP等级
func (h *VipHandler) UpdateLevel(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var req VipUpdateLevelRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	level := &model.VipLevel{
		Slug:                    req.Slug,
		Title:                   req.Title,
		ExpRequired:             req.ExpRequired,
		OrderDiscount:           req.OrderDiscount,
		MonthlyCouponTemplateID: req.MonthlyCouponTemplateID,
		MonthlyCouponCount:      req.MonthlyCouponCount,
		IconURL:                 req.IconURL,
		Color:                   req.Color,
		Benefits:                req.Benefits,
		SortOrder:               req.SortOrder,
		IsDefault:               req.IsDefault,
		IsActive:                req.IsActive,
	}
	level.ID = id

	if err := h.svc.UpdateLevel(c.Request.Context(), level); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("VIP等级不存在"))
			return
		}
		respondError(c, apierr.InternalError("更新VIP等级失败").WithDetails(err.Error()))
		return
	}

	respondUpdated(c, level)
}

// DeleteLevel 删除VIP等级
func (h *VipHandler) DeleteLevel(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	if err := h.svc.DeleteLevel(c.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("VIP等级不存在"))
			return
		}
		respondError(c, apierr.InternalError("删除VIP等级失败").WithDetails(err.Error()))
		return
	}

	respondDeleted(c)
}

// SetDefaultLevel 设置默认VIP等级
func (h *VipHandler) SetDefaultLevel(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	if err := h.svc.SetDefaultLevel(c.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("VIP等级不存在"))
			return
		}
		respondError(c, apierr.InternalError("设置默认等级失败").WithDetails(err.Error()))
		return
	}

	respondMsg(c, "设置成功")
}

// VipBatchUpdateStatusRequest 批量更新状态请求
type VipBatchUpdateStatusRequest struct {
	IDs      []uint64 `json:"ids" binding:"required,min=1"`
	IsActive bool     `json:"isActive"`
}

// VipBatchDeleteRequest 批量删除请求
type VipBatchDeleteRequest struct {
	IDs []uint64 `json:"ids" binding:"required,min=1"`
}

// BatchUpdateLevelStatus 批量更新VIP等级状态
func (h *VipHandler) BatchUpdateLevelStatus(c *gin.Context) {
	var req VipBatchUpdateStatusRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	affected, err := h.svc.BatchUpdateLevelStatus(c.Request.Context(), req.IDs, req.IsActive)
	if err != nil {
		respondError(c, apierr.InternalError("批量更新状态失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, gin.H{"affected": affected})
}

// BatchDeleteLevels 批量删除VIP等级
func (h *VipHandler) BatchDeleteLevels(c *gin.Context) {
	var req VipBatchDeleteRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	affected, err := h.svc.BatchDeleteLevels(c.Request.Context(), req.IDs)
	if err != nil {
		respondError(c, apierr.InternalError("批量删除失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, gin.H{"affected": affected})
}

// ============================================================================
// VIP配置管理
// ============================================================================

// VipSaveConfigRequest 保存VIP配置请求
type VipSaveConfigRequest struct {
	ConfigKey   string `json:"configKey" binding:"required,max=64"`
	ConfigValue string `json:"configValue" binding:"required"`
	Description string `json:"description"`
}

// ListConfigs 获取VIP配置列表
func (h *VipHandler) ListConfigs(c *gin.Context) {
	configs, err := h.svc.ListConfigs(c.Request.Context())
	if err != nil {
		respondError(c, apierr.InternalError("获取VIP配置列表失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, configs)
}

// GetConfig 获取VIP配置详情
func (h *VipHandler) GetConfig(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		respondBadRequest(c, "配置键不能为空")
		return
	}

	config, err := h.svc.GetConfig(c.Request.Context(), key)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("VIP配置不存在"))
			return
		}
		respondError(c, apierr.InternalError("获取VIP配置失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, config)
}

// SaveConfig 保存VIP配置
func (h *VipHandler) SaveConfig(c *gin.Context) {
	var req VipSaveConfigRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	config := &model.VipConfig{
		ConfigKey:   req.ConfigKey,
		ConfigValue: req.ConfigValue,
		Description: req.Description,
	}

	if err := h.svc.SaveConfig(c.Request.Context(), config); err != nil {
		respondError(c, apierr.InternalError("保存VIP配置失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, config)
}

// DeleteConfig 删除VIP配置
func (h *VipHandler) DeleteConfig(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		respondBadRequest(c, "配置键不能为空")
		return
	}

	if err := h.svc.DeleteConfig(c.Request.Context(), key); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("VIP配置不存在"))
			return
		}
		respondError(c, apierr.InternalError("删除VIP配置失败").WithDetails(err.Error()))
		return
	}

	respondDeleted(c)
}

// ============================================================================
// 路由注册
// ============================================================================

// RegisterVipRoutes 注册VIP管理路由
func RegisterVipRoutes(rg gin.IRouter, svc *vipservice.Service, pm *middleware.PermissionMiddleware) {
	h := NewVipHandler(svc)

	vipGroup := rg.Group("/vip")
	vipGroup.Use(pm.RequireAuth())
	{
		// VIP等级管理
		vipGroup.GET("/levels", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/vip/levels"), h.ListLevels)
		vipGroup.GET("/levels/:id", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/vip/levels/:id"), h.GetLevel)
		vipGroup.POST("/levels", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/vip/levels"), h.CreateLevel)
		vipGroup.PUT("/levels/:id", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/vip/levels/:id"), h.UpdateLevel)
		vipGroup.DELETE("/levels/:id", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/vip/levels/:id"), h.DeleteLevel)
		vipGroup.POST("/levels/:id/default", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/vip/levels/:id/default"), h.SetDefaultLevel)
		vipGroup.POST("/levels/batch-status", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/vip/levels/batch-status"), h.BatchUpdateLevelStatus)
		vipGroup.POST("/levels/batch-delete", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/vip/levels/batch-delete"), h.BatchDeleteLevels)

		// VIP配置管理
		vipGroup.GET("/configs", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/vip/configs"), h.ListConfigs)
		vipGroup.GET("/configs/:key", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/vip/configs/:key"), h.GetConfig)
		vipGroup.POST("/configs", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/vip/configs"), h.SaveConfig)
		vipGroup.DELETE("/configs/:key", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/vip/configs/:key"), h.DeleteConfig)
	}
}
