package user

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	"gamelink/internal/model"
	referralservice "gamelink/internal/service/referral"
	"gamelink/pkg/apierr"
)

// ReferralHandler 用户端推荐处理器
type ReferralHandler struct {
	svc *referralservice.Service
}

// NewReferralHandler 创建用户端推荐处理器
func NewReferralHandler(svc *referralservice.Service) *ReferralHandler {
	return &ReferralHandler{svc: svc}
}

// RegisterRoutes 注册路由
func (h *ReferralHandler) RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/referrals")
	{
		// 邀请码
		g.GET("/code", h.GetMyCode)
		g.POST("/code", h.CreateMyCode)

		// 使用邀请码
		g.POST("/use", h.UseCode)
		g.GET("/validate/:code", h.ValidateCode)

		// 我的推荐
		g.GET("/my", h.GetMyReferrals)
		g.GET("/my/stats", h.GetMyStats)

		// 我的奖励
		g.GET("/my/rewards", h.GetMyRewards)
	}
}

// GetMyCode 获取我的邀请码
func (h *ReferralHandler) GetMyCode(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	// 默认获取用户邀请用户类型的邀请码
	refType := model.ReferralType(c.DefaultQuery("type", string(model.ReferralTypeUserToUser)))

	code, err := h.svc.GetOrCreateUserCode(c.Request.Context(), userID, refType)
	if err != nil {
		resp.Error(c, apierr.InternalError("get code failed").WithDetails(err.Error()))
		return
	}
	resp.OK(c, code)
}

// createMyCodeRequest 创建邀请码请求
type createMyCodeRequest struct {
	Type model.ReferralType `json:"type"`
}

// CreateMyCode 创建我的邀请码
func (h *ReferralHandler) CreateMyCode(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	var req createMyCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, apierr.BadRequest("invalid request").WithDetails(err.Error()))
		return
	}

	if req.Type == "" {
		req.Type = model.ReferralTypeUserToUser
	}

	code, err := h.svc.GetOrCreateUserCode(c.Request.Context(), userID, req.Type)
	if err != nil {
		resp.Error(c, apierr.InternalError("create code failed").WithDetails(err.Error()))
		return
	}
	resp.Created(c, code)
}

// useCodeRequest 使用邀请码请求
type useCodeRequest struct {
	Code string `json:"code" binding:"required"`
}

// UseCode 使用邀请码（注册时调用）
func (h *ReferralHandler) UseCode(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	var req useCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, apierr.BadRequest("invalid request").WithDetails(err.Error()))
		return
	}

	referral, err := h.svc.UseCode(c.Request.Context(), referralservice.UseCodeRequest{
		Code:      req.Code,
		RefereeID: userID,
	})
	if err != nil {
		resp.Error(c, apierr.BadRequest("use code failed").WithDetails(err.Error()))
		return
	}
	resp.OK(c, referral)
}

// ValidateCode 验证邀请码
func (h *ReferralHandler) ValidateCode(c *gin.Context) {
	codeStr := c.Param("code")
	if codeStr == "" {
		resp.Error(c, apierr.BadRequest("code is required"))
		return
	}

	code, err := h.svc.ValidateCode(c.Request.Context(), codeStr)
	if err != nil {
		resp.Error(c, apierr.BadRequest("invalid code").WithDetails(err.Error()))
		return
	}

	// 返回简化信息，不暴露敏感数据
	resp.OK(c, gin.H{
		"valid":    true,
		"type":     code.Type,
		"referrer": code.User,
	})
}

// GetMyReferrals 获取我的推荐记录
func (h *ReferralHandler) GetMyReferrals(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	limit := 20
	if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 {
		limit = l
	}

	referrals, err := h.svc.GetUserReferrals(c.Request.Context(), userID, limit)
	if err != nil {
		resp.Error(c, apierr.InternalError("get referrals failed").WithDetails(err.Error()))
		return
	}
	resp.OK(c, referrals)
}

// GetMyStats 获取我的推荐统计
func (h *ReferralHandler) GetMyStats(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	stats, err := h.svc.GetUserReferralStats(c.Request.Context(), userID)
	if err != nil {
		resp.Error(c, apierr.InternalError("get stats failed").WithDetails(err.Error()))
		return
	}
	resp.OK(c, stats)
}

// GetMyRewards 获取我的奖励记录
func (h *ReferralHandler) GetMyRewards(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	limit := 20
	if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 {
		limit = l
	}

	rewards, err := h.svc.GetUserRewards(c.Request.Context(), userID, limit)
	if err != nil {
		resp.Error(c, apierr.InternalError("get rewards failed").WithDetails(err.Error()))
		return
	}
	resp.OK(c, rewards)
}

// RegisterReferralRoutes 注册用户端推荐路由
func RegisterReferralRoutes(r *gin.RouterGroup, svc *referralservice.Service, _ gin.HandlerFunc) {
	h := NewReferralHandler(svc)
	h.RegisterRoutes(r)
}
