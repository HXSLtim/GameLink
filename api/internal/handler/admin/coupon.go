package admin

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/middleware"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	couponrepo "gamelink/internal/repository/coupon"
	couponservice "gamelink/internal/service/coupon"
	"gamelink/pkg/apierr"
)

// CouponHandler 优惠券管理处理器
type CouponHandler struct {
	svc *couponservice.Service
}

// NewCouponHandler 创建优惠券处理器
func NewCouponHandler(svc *couponservice.Service) *CouponHandler {
	return &CouponHandler{svc: svc}
}

// ============================================================================
// 优惠券模板管理
// ============================================================================

// TemplateCreateRequest 创建模板请求
type TemplateCreateRequest struct {
	Name              string             `json:"name" binding:"required,max=128"`
	Type              model.CouponType   `json:"type" binding:"required"`
	Source            model.CouponSource `json:"source" binding:"required"`
	Description       string             `json:"description"`
	MinAmountCents    int64              `json:"minAmountCents"`
	DeductAmountCents int64              `json:"deductAmountCents"`
	DiscountRate      float64            `json:"discountRate"`
	MaxDiscountCents  int64              `json:"maxDiscountCents"`
	Scope             model.CouponScope  `json:"scope"`
	GameIDs           string             `json:"gameIds"`
	ItemIDs           string             `json:"itemIds"`
	ValidityType      string             `json:"validityType"`
	ValidityDays      int                `json:"validityDays"`
	FixedExpireAt     string             `json:"fixedExpireAt"`
	TotalCount        int                `json:"totalCount"`
	PerUserLimit      int                `json:"perUserLimit"`
	ClaimLink         string             `json:"claimLink"`
	IsActive          bool               `json:"isActive"`
}

// TemplateUpdateRequest 更新模板请求
type TemplateUpdateRequest struct {
	Name              string             `json:"name" binding:"required,max=128"`
	Type              model.CouponType   `json:"type" binding:"required"`
	Source            model.CouponSource `json:"source" binding:"required"`
	Description       string             `json:"description"`
	MinAmountCents    int64              `json:"minAmountCents"`
	DeductAmountCents int64              `json:"deductAmountCents"`
	DiscountRate      float64            `json:"discountRate"`
	MaxDiscountCents  int64              `json:"maxDiscountCents"`
	Scope             model.CouponScope  `json:"scope"`
	GameIDs           string             `json:"gameIds"`
	ItemIDs           string             `json:"itemIds"`
	ValidityType      string             `json:"validityType"`
	ValidityDays      int                `json:"validityDays"`
	FixedExpireAt     string             `json:"fixedExpireAt"`
	TotalCount        int                `json:"totalCount"`
	PerUserLimit      int                `json:"perUserLimit"`
	ClaimLink         string             `json:"claimLink"`
	IsActive          bool               `json:"isActive"`
}

// ListTemplates 获取模板列表
func (h *CouponHandler) ListTemplates(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	var couponType *model.CouponType
	if v := strings.TrimSpace(c.Query("type")); v != "" {
		t := model.CouponType(v)
		couponType = &t
	}

	var source *model.CouponSource
	if v := strings.TrimSpace(c.Query("source")); v != "" {
		s := model.CouponSource(v)
		source = &s
	}

	var isActive *bool
	if v := strings.TrimSpace(c.Query("isActive")); v != "" {
		b := v == "true" || v == "1"
		isActive = &b
	}

	opts := couponrepo.TemplateListOptions{
		Page:     page,
		PageSize: pageSize,
		Keyword:  c.Query("keyword"),
		Type:     couponType,
		Source:   source,
		IsActive: isActive,
	}

	templates, total, err := h.svc.ListTemplates(c.Request.Context(), opts)
	if err != nil {
		respondError(c, apierr.InternalError("获取模板列表失败").WithDetails(err.Error()))
		return
	}

	respondList(c, templates, &model.Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
	})
}

// GetTemplate 获取模板详情
func (h *CouponHandler) GetTemplate(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	template, err := h.svc.GetTemplate(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("模板不存在"))
			return
		}
		respondError(c, apierr.InternalError("获取模板失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, template)
}

// CreateTemplate 创建模板
func (h *CouponHandler) CreateTemplate(c *gin.Context) {
	var req TemplateCreateRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	template := &model.CouponTemplate{
		Name:              req.Name,
		Type:              req.Type,
		Source:            req.Source,
		Description:       req.Description,
		MinAmountCents:    req.MinAmountCents,
		DeductAmountCents: req.DeductAmountCents,
		DiscountRate:      req.DiscountRate,
		MaxDiscountCents:  req.MaxDiscountCents,
		Scope:             req.Scope,
		GameIDs:           req.GameIDs,
		ItemIDs:           req.ItemIDs,
		ValidityType:      req.ValidityType,
		ValidityDays:      req.ValidityDays,
		TotalCount:        req.TotalCount,
		PerUserLimit:      req.PerUserLimit,
		ClaimLink:         req.ClaimLink,
		IsActive:          req.IsActive,
	}

	if err := h.svc.CreateTemplate(c.Request.Context(), template); err != nil {
		respondError(c, apierr.InternalError("创建模板失败").WithDetails(err.Error()))
		return
	}

	respondCreated(c, template)
}

// UpdateTemplate 更新模板
func (h *CouponHandler) UpdateTemplate(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var req TemplateUpdateRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	template := &model.CouponTemplate{
		Name:              req.Name,
		Type:              req.Type,
		Source:            req.Source,
		Description:       req.Description,
		MinAmountCents:    req.MinAmountCents,
		DeductAmountCents: req.DeductAmountCents,
		DiscountRate:      req.DiscountRate,
		MaxDiscountCents:  req.MaxDiscountCents,
		Scope:             req.Scope,
		GameIDs:           req.GameIDs,
		ItemIDs:           req.ItemIDs,
		ValidityType:      req.ValidityType,
		ValidityDays:      req.ValidityDays,
		TotalCount:        req.TotalCount,
		PerUserLimit:      req.PerUserLimit,
		ClaimLink:         req.ClaimLink,
		IsActive:          req.IsActive,
	}
	template.ID = id

	if err := h.svc.UpdateTemplate(c.Request.Context(), template); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("模板不存在"))
			return
		}
		respondError(c, apierr.InternalError("更新模板失败").WithDetails(err.Error()))
		return
	}

	respondUpdated(c, template)
}

// DeleteTemplate 删除模板
func (h *CouponHandler) DeleteTemplate(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	if err := h.svc.DeleteTemplate(c.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("模板不存在"))
			return
		}
		respondError(c, apierr.InternalError("删除模板失败").WithDetails(err.Error()))
		return
	}

	respondDeleted(c)
}

// TemplateBatchStatusRequest 批量更新模板状态请求
type TemplateBatchStatusRequest struct {
	IDs      []uint64 `json:"ids" binding:"required,min=1"`
	IsActive bool     `json:"isActive"`
}

// TemplateBatchDeleteRequest 批量删除模板请求
type TemplateBatchDeleteRequest struct {
	IDs []uint64 `json:"ids" binding:"required,min=1"`
}

// BatchUpdateTemplateStatus 批量更新模板状态
func (h *CouponHandler) BatchUpdateTemplateStatus(c *gin.Context) {
	var req TemplateBatchStatusRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	affected, err := h.svc.BatchUpdateTemplateStatus(c.Request.Context(), req.IDs, req.IsActive)
	if err != nil {
		respondError(c, apierr.InternalError("批量更新状态失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, gin.H{"affected": affected})
}

// BatchDeleteTemplates 批量删除模板
func (h *CouponHandler) BatchDeleteTemplates(c *gin.Context) {
	var req TemplateBatchDeleteRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	affected, err := h.svc.BatchDeleteTemplates(c.Request.Context(), req.IDs)
	if err != nil {
		respondError(c, apierr.InternalError("批量删除失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, gin.H{"affected": affected})
}

// ============================================================================
// 用户优惠券管理
// ============================================================================

// ListCoupons 获取优惠券列表
func (h *CouponHandler) ListCoupons(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	var userID *uint64
	if v := strings.TrimSpace(c.Query("userId")); v != "" {
		id, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			respondBadRequest(c, "无效的用户ID")
			return
		}
		userID = &id
	}

	var templateID *uint64
	if v := strings.TrimSpace(c.Query("templateId")); v != "" {
		id, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			respondBadRequest(c, "无效的模板ID")
			return
		}
		templateID = &id
	}

	var state *model.CouponState
	if v := strings.TrimSpace(c.Query("state")); v != "" {
		s := model.CouponState(v)
		state = &s
	}

	var couponType *model.CouponType
	if v := strings.TrimSpace(c.Query("type")); v != "" {
		t := model.CouponType(v)
		couponType = &t
	}

	opts := couponrepo.CouponListOptions{
		Page:       page,
		PageSize:   pageSize,
		UserID:     userID,
		TemplateID: templateID,
		State:      state,
		Type:       couponType,
	}

	coupons, total, err := h.svc.ListCoupons(c.Request.Context(), opts)
	if err != nil {
		respondError(c, apierr.InternalError("获取优惠券列表失败").WithDetails(err.Error()))
		return
	}

	respondList(c, coupons, &model.Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
	})
}

// GetCoupon 获取优惠券详情
func (h *CouponHandler) GetCoupon(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	coupon, err := h.svc.GetCouponWithTemplate(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("优惠券不存在"))
			return
		}
		respondError(c, apierr.InternalError("获取优惠券失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, coupon)
}

// DeleteCoupon 删除优惠券
func (h *CouponHandler) DeleteCoupon(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	if err := h.svc.DeleteCoupon(c.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("优惠券不存在"))
			return
		}
		respondError(c, apierr.InternalError("删除优惠券失败").WithDetails(err.Error()))
		return
	}

	respondDeleted(c)
}

// IssueCouponRequest 发放优惠券请求
type IssueCouponRequest struct {
	UserID     uint64             `json:"userId" binding:"required"`
	TemplateID uint64             `json:"templateId" binding:"required"`
	Source     model.CouponSource `json:"source"`
}

// IssueCoupon 发放优惠券
func (h *CouponHandler) IssueCoupon(c *gin.Context) {
	var req IssueCouponRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	source := req.Source
	if source == "" {
		source = model.CouponSourceManual
	}

	coupon, err := h.svc.IssueCoupon(c.Request.Context(), req.UserID, req.TemplateID, source)
	if err != nil {
		respondError(c, apierr.BadRequest(err.Error()))
		return
	}

	respondCreated(c, coupon)
}

// GetCouponStats 获取优惠券统计
func (h *CouponHandler) GetCouponStats(c *gin.Context) {
	stats, err := h.svc.GetCouponStats(c.Request.Context())
	if err != nil {
		respondError(c, apierr.InternalError("获取统计失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, stats)
}

// ============================================================================
// 路由注册
// ============================================================================

// RegisterCouponRoutes 注册优惠券管理路由
func RegisterCouponRoutes(rg *gin.RouterGroup, svc *couponservice.Service, pm *middleware.PermissionMiddleware) {
	h := NewCouponHandler(svc)

	couponGroup := rg.Group("/coupons")
	couponGroup.Use(pm.RequireAuth())
	{
		// 优惠券模板管理
		couponGroup.GET("/templates", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/coupons/templates"), h.ListTemplates)
		couponGroup.GET("/templates/:id", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/coupons/templates/:id"), h.GetTemplate)
		couponGroup.POST("/templates", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/coupons/templates"), h.CreateTemplate)
		couponGroup.PUT("/templates/:id", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/coupons/templates/:id"), h.UpdateTemplate)
		couponGroup.DELETE("/templates/:id", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/coupons/templates/:id"), h.DeleteTemplate)
		couponGroup.POST("/templates/batch-status", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/coupons/templates/batch-status"), h.BatchUpdateTemplateStatus)
		couponGroup.POST("/templates/batch-delete", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/coupons/templates/batch-delete"), h.BatchDeleteTemplates)

		// 用户优惠券管理
		couponGroup.GET("", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/coupons"), h.ListCoupons)
		couponGroup.GET("/stats", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/coupons/stats"), h.GetCouponStats)
		couponGroup.GET("/:id", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/coupons/:id"), h.GetCoupon)
		couponGroup.POST("/issue", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/coupons/issue"), h.IssueCoupon)
		couponGroup.DELETE("/:id", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/coupons/:id"), h.DeleteCoupon)
	}
}
