package admin

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/middleware"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	rechargerepo "gamelink/internal/repository/recharge"
	rechargeservice "gamelink/internal/service/recharge"
	"gamelink/pkg/apierr"
)

// RechargeHandler 充值管理处理器
type RechargeHandler struct {
	svc *rechargeservice.Service
}

// NewRechargeHandler 创建充值处理器
func NewRechargeHandler(svc *rechargeservice.Service) *RechargeHandler {
	return &RechargeHandler{svc: svc}
}

// ============================================================================
// 充值档位管理
// ============================================================================

// OptionCreateRequest 创建档位请求
type OptionCreateRequest struct {
	Name             string  `json:"name" binding:"required,max=64"`
	AmountCents      int64   `json:"amountCents" binding:"required,gt=0"`
	BonusCents       int64   `json:"bonusCents"`
	OriginalCents    *int64  `json:"originalCents"`
	DiscountPercent  *int    `json:"discountPercent"`
	Description      string  `json:"description"`
	Tag              string  `json:"tag"`
	IconURL          string  `json:"iconUrl"`
	SortOrder        int     `json:"sortOrder"`
	IsActive         bool    `json:"isActive"`
	IsRecommended    bool    `json:"isRecommended"`
	CouponTemplateID *uint64 `json:"couponTemplateId"`
	CouponCount      int     `json:"couponCount"`
	MinVipLevel      *uint64 `json:"minVipLevel"`
	PerUserLimit     int     `json:"perUserLimit"`
	TotalLimit       int     `json:"totalLimit"`
}

// OptionUpdateRequest 更新档位请求
type OptionUpdateRequest struct {
	Name             string  `json:"name" binding:"required,max=64"`
	AmountCents      int64   `json:"amountCents" binding:"required,gt=0"`
	BonusCents       int64   `json:"bonusCents"`
	OriginalCents    *int64  `json:"originalCents"`
	DiscountPercent  *int    `json:"discountPercent"`
	Description      string  `json:"description"`
	Tag              string  `json:"tag"`
	IconURL          string  `json:"iconUrl"`
	SortOrder        int     `json:"sortOrder"`
	IsActive         bool    `json:"isActive"`
	IsRecommended    bool    `json:"isRecommended"`
	CouponTemplateID *uint64 `json:"couponTemplateId"`
	CouponCount      int     `json:"couponCount"`
	MinVipLevel      *uint64 `json:"minVipLevel"`
	PerUserLimit     int     `json:"perUserLimit"`
	TotalLimit       int     `json:"totalLimit"`
}

// ListOptions 获取档位列表
func (h *RechargeHandler) ListOptions(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	var isActive *bool
	if v := strings.TrimSpace(c.Query("isActive")); v != "" {
		b := v == "true" || v == "1"
		isActive = &b
	}

	var isRecommended *bool
	if v := strings.TrimSpace(c.Query("isRecommended")); v != "" {
		b := v == "true" || v == "1"
		isRecommended = &b
	}

	opts := rechargerepo.OptionListOptions{
		Page:          page,
		PageSize:      pageSize,
		Keyword:       c.Query("keyword"),
		IsActive:      isActive,
		IsRecommended: isRecommended,
	}

	options, total, err := h.svc.ListOptions(c.Request.Context(), opts)
	if err != nil {
		respondError(c, apierr.InternalError("获取档位列表失败").WithDetails(err.Error()))
		return
	}

	respondList(c, options, &model.Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
	})
}

// GetOption 获取档位详情
func (h *RechargeHandler) GetOption(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	option, err := h.svc.GetOption(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("档位不存在"))
			return
		}
		respondError(c, apierr.InternalError("获取档位失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, option)
}

// CreateOption 创建档位
func (h *RechargeHandler) CreateOption(c *gin.Context) {
	var req OptionCreateRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	option := &model.RechargeOption{
		Name:             req.Name,
		AmountCents:      req.AmountCents,
		BonusCents:       req.BonusCents,
		OriginalCents:    req.OriginalCents,
		DiscountPercent:  req.DiscountPercent,
		Description:      req.Description,
		Tag:              req.Tag,
		IconURL:          req.IconURL,
		SortOrder:        req.SortOrder,
		IsActive:         req.IsActive,
		IsRecommended:    req.IsRecommended,
		CouponTemplateID: req.CouponTemplateID,
		CouponCount:      req.CouponCount,
		MinVipLevel:      req.MinVipLevel,
		PerUserLimit:     req.PerUserLimit,
		TotalLimit:       req.TotalLimit,
	}

	if err := h.svc.CreateOption(c.Request.Context(), option); err != nil {
		respondError(c, apierr.InternalError("创建档位失败").WithDetails(err.Error()))
		return
	}

	respondCreated(c, option)
}

// UpdateOption 更新档位
func (h *RechargeHandler) UpdateOption(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var req OptionUpdateRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	option := &model.RechargeOption{
		Name:             req.Name,
		AmountCents:      req.AmountCents,
		BonusCents:       req.BonusCents,
		OriginalCents:    req.OriginalCents,
		DiscountPercent:  req.DiscountPercent,
		Description:      req.Description,
		Tag:              req.Tag,
		IconURL:          req.IconURL,
		SortOrder:        req.SortOrder,
		IsActive:         req.IsActive,
		IsRecommended:    req.IsRecommended,
		CouponTemplateID: req.CouponTemplateID,
		CouponCount:      req.CouponCount,
		MinVipLevel:      req.MinVipLevel,
		PerUserLimit:     req.PerUserLimit,
		TotalLimit:       req.TotalLimit,
	}
	option.ID = id

	if err := h.svc.UpdateOption(c.Request.Context(), option); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("档位不存在"))
			return
		}
		respondError(c, apierr.InternalError("更新档位失败").WithDetails(err.Error()))
		return
	}

	respondUpdated(c, option)
}

// DeleteOption 删除档位
func (h *RechargeHandler) DeleteOption(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	if err := h.svc.DeleteOption(c.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("档位不存在"))
			return
		}
		respondError(c, apierr.InternalError("删除档位失败").WithDetails(err.Error()))
		return
	}

	respondDeleted(c)
}

// OptionBatchStatusRequest 批量更新档位状态请求
type OptionBatchStatusRequest struct {
	IDs      []uint64 `json:"ids" binding:"required,min=1"`
	IsActive bool     `json:"isActive"`
}

// OptionBatchDeleteRequest 批量删除档位请求
type OptionBatchDeleteRequest struct {
	IDs []uint64 `json:"ids" binding:"required,min=1"`
}

// BatchUpdateOptionStatus 批量更新档位状态
func (h *RechargeHandler) BatchUpdateOptionStatus(c *gin.Context) {
	var req OptionBatchStatusRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	affected, err := h.svc.BatchUpdateOptionStatus(c.Request.Context(), req.IDs, req.IsActive)
	if err != nil {
		respondError(c, apierr.InternalError("批量更新状态失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, gin.H{"affected": affected})
}

// BatchDeleteOptions 批量删除档位
func (h *RechargeHandler) BatchDeleteOptions(c *gin.Context) {
	var req OptionBatchDeleteRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	affected, err := h.svc.BatchDeleteOptions(c.Request.Context(), req.IDs)
	if err != nil {
		respondError(c, apierr.InternalError("批量删除失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, gin.H{"affected": affected})
}

// ============================================================================
// 充值记录管理
// ============================================================================

// ListRecords 获取充值记录列表
func (h *RechargeHandler) ListRecords(c *gin.Context) {
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

	var optionID *uint64
	if v := strings.TrimSpace(c.Query("optionId")); v != "" {
		id, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			respondBadRequest(c, "无效的档位ID")
			return
		}
		optionID = &id
	}

	var status *model.RechargeStatus
	if v := strings.TrimSpace(c.Query("status")); v != "" {
		s := model.RechargeStatus(v)
		status = &s
	}

	var paymentChannel *string
	if v := strings.TrimSpace(c.Query("paymentChannel")); v != "" {
		paymentChannel = &v
	}

	var startTime, endTime *time.Time
	if v := strings.TrimSpace(c.Query("startTime")); v != "" {
		t, err := time.Parse("2006-01-02", v)
		if err == nil {
			startTime = &t
		}
	}
	if v := strings.TrimSpace(c.Query("endTime")); v != "" {
		t, err := time.Parse("2006-01-02", v)
		if err == nil {
			t = t.Add(24*time.Hour - time.Second)
			endTime = &t
		}
	}

	opts := rechargerepo.RecordListOptions{
		Page:           page,
		PageSize:       pageSize,
		UserID:         userID,
		OptionID:       optionID,
		Status:         status,
		PaymentChannel: paymentChannel,
		OrderNo:        c.Query("orderNo"),
		StartTime:      startTime,
		EndTime:        endTime,
	}

	records, total, err := h.svc.ListRecords(c.Request.Context(), opts)
	if err != nil {
		respondError(c, apierr.InternalError("获取充值记录失败").WithDetails(err.Error()))
		return
	}

	respondList(c, records, &model.Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
	})
}

// GetRecord 获取充值记录详情
func (h *RechargeHandler) GetRecord(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	record, err := h.svc.GetRecord(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("充值记录不存在"))
			return
		}
		respondError(c, apierr.InternalError("获取充值记录失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, record)
}

// RefundRequest 退款请求
type RefundRequest struct {
	Reason string `json:"reason" binding:"required,max=255"`
}

// RefundRecord 退款
func (h *RechargeHandler) RefundRecord(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var req RefundRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	if err := h.svc.RefundRecord(c.Request.Context(), id, req.Reason); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("充值记录不存在"))
			return
		}
		respondError(c, apierr.BadRequest(err.Error()))
		return
	}

	respondSuccess(c, gin.H{"message": "退款成功"})
}

// GetRechargeStats 获取充值统计
func (h *RechargeHandler) GetRechargeStats(c *gin.Context) {
	stats, err := h.svc.GetRechargeStats(c.Request.Context())
	if err != nil {
		respondError(c, apierr.InternalError("获取统计失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, stats)
}

// ============================================================================
// 路由注册
// ============================================================================

// RegisterRechargeRoutes 注册充值管理路由
func RegisterRechargeRoutes(rg *gin.RouterGroup, svc *rechargeservice.Service, pm *middleware.PermissionMiddleware) {
	h := NewRechargeHandler(svc)

	rechargeGroup := rg.Group("/recharge")
	rechargeGroup.Use(pm.RequireAuth())
	{
		// 充值档位管理
		rechargeGroup.GET("/options", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/recharge/options"), h.ListOptions)
		rechargeGroup.GET("/options/:id", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/recharge/options/:id"), h.GetOption)
		rechargeGroup.POST("/options", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/recharge/options"), h.CreateOption)
		rechargeGroup.PUT("/options/:id", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/recharge/options/:id"), h.UpdateOption)
		rechargeGroup.DELETE("/options/:id", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/recharge/options/:id"), h.DeleteOption)
		rechargeGroup.POST("/options/batch-status", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/recharge/options/batch-status"), h.BatchUpdateOptionStatus)
		rechargeGroup.POST("/options/batch-delete", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/recharge/options/batch-delete"), h.BatchDeleteOptions)

		// 充值记录管理
		rechargeGroup.GET("/records", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/recharge/records"), h.ListRecords)
		rechargeGroup.GET("/records/:id", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/recharge/records/:id"), h.GetRecord)
		rechargeGroup.POST("/records/:id/refund", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/recharge/records/:id/refund"), h.RefundRecord)
		rechargeGroup.GET("/stats", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/recharge/stats"), h.GetRechargeStats)
	}
}
