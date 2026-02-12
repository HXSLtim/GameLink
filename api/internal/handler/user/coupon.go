package user

import (
	"errors"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	"gamelink/internal/repository"
	couponservice "gamelink/internal/service/coupon"
	"gamelink/pkg/apierr"
)

// CouponHandler 用户优惠券处理器
type CouponHandler struct {
	svc *couponservice.Service
}

// NewCouponHandler 创建用户优惠券处理器
func NewCouponHandler(svc *couponservice.Service) *CouponHandler {
	return &CouponHandler{svc: svc}
}

// ListMyCoupons 获取我的优惠券列表
func (h *CouponHandler) ListMyCoupons(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	coupons, err := h.svc.GetUserAvailableCoupons(c.Request.Context(), userID)
	if err != nil {
		resp.Error(c, apierr.InternalError("获取优惠券列表失败").WithDetails(err.Error()))
		return
	}

	resp.OK(c, coupons)
}

// GetMyCouponStats 获取当前用户优惠券统计（兼容前端 /coupons/stats）
func (h *CouponHandler) GetMyCouponStats(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	stats, err := h.svc.GetUserCouponStats(c.Request.Context(), userID)
	if err != nil {
		resp.Error(c, apierr.InternalError("获取优惠券统计失败").WithDetails(err.Error()))
		return
	}

	resp.OK(c, stats)
}

// GetCoupon 获取优惠券详情
func (h *CouponHandler) GetCoupon(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	id, ok := resp.ParseIDOrFail(c, "id")
	if !ok {
		return
	}

	coupon, err := h.svc.GetCouponWithTemplate(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			resp.Error(c, apierr.NotFound("优惠券不存在"))
			return
		}
		resp.Error(c, apierr.InternalError("获取优惠券失败").WithDetails(err.Error()))
		return
	}

	// 验证所有权
	if coupon.UserID != userID {
		resp.Error(c, apierr.NotFound("优惠券不存在"))
		return
	}

	resp.OK(c, coupon)
}

// ClaimCouponRequest 领取优惠券请求
type ClaimCouponRequest struct {
	TemplateID uint64 `json:"templateId"`
	Link       string `json:"link"`
}

// ClaimCoupon 领取优惠券
func (h *CouponHandler) ClaimCoupon(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	var req ClaimCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, apierr.BadRequest("无效的请求参数").WithDetails(err.Error()))
		return
	}

	var coupon any
	var err error

	if req.Link != "" {
		// 通过链接领取
		coupon, err = h.svc.ClaimCouponByLink(c.Request.Context(), userID, req.Link)
	} else if req.TemplateID > 0 {
		// 通过模板ID领取
		coupon, err = h.svc.ClaimCoupon(c.Request.Context(), userID, req.TemplateID)
	} else {
		resp.Error(c, apierr.BadRequest("请提供模板ID或领取链接"))
		return
	}

	if err != nil {
		resp.Error(c, apierr.BadRequest(err.Error()))
		return
	}

	resp.Created(c, coupon)
}

// RegisterCouponRoutes 注册用户优惠券路由
func RegisterCouponRoutes(rg *gin.RouterGroup, svc *couponservice.Service, _ gin.HandlerFunc) {
	h := NewCouponHandler(svc)

	couponGroup := rg.Group("/coupons")
	{
		couponGroup.GET("", h.ListMyCoupons)
		couponGroup.GET("/stats", h.GetMyCouponStats)
		couponGroup.GET("/:id", h.GetCoupon)
		couponGroup.POST("/claim", h.ClaimCoupon)
	}
}
