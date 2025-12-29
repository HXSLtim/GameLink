package admin

import (
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/middleware"
	"gamelink/internal/model"
	referralrepo "gamelink/internal/repository/referral"
	referralservice "gamelink/internal/service/referral"
)

// ReferralHandler 推荐管理处理器
type ReferralHandler struct {
	svc *referralservice.Service
}

// NewReferralHandler 创建推荐管理处理器
func NewReferralHandler(svc *referralservice.Service) *ReferralHandler {
	return &ReferralHandler{svc: svc}
}

// RegisterRoutes 注册路由
func (h *ReferralHandler) RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/referrals")
	{
		// 配置管理
		g.GET("/configs", h.ListConfigs)
		g.PUT("/configs/:key", h.UpdateConfig)

		// 邀请码管理
		g.GET("/codes", h.ListCodes)
		g.GET("/codes/:id", h.GetCode)
		g.POST("/codes", h.CreateCode)
		g.PUT("/codes/:id", h.UpdateCode)
		g.DELETE("/codes/:id", h.DeleteCode)

		// 推荐记录管理
		g.GET("", h.ListReferrals)
		g.GET("/:id", h.GetReferral)
		g.PUT("/:id/status", h.UpdateReferralStatus)

		// 奖励管理
		g.GET("/rewards", h.ListRewards)
		g.GET("/rewards/:id", h.GetReward)
		g.POST("/rewards/:id/issue", h.IssueReward)
		g.POST("/rewards/:id/fail", h.FailReward)

		// 统计
		g.GET("/stats", h.GetStats)
	}
}

// ============================================================================
// 配置管理
// ============================================================================

// ListConfigs 获取所有配置
func (h *ReferralHandler) ListConfigs(c *gin.Context) {
	configs, err := h.svc.GetAllConfigs(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, configs)
}

// updateConfigRequest 更新配置请求
type updateConfigRequest struct {
	Value       string `json:"value" binding:"required"`
	Description string `json:"description"`
}

// UpdateConfig 更新配置
func (h *ReferralHandler) UpdateConfig(c *gin.Context) {
	key := c.Param("key")

	var req updateConfigRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	if err := h.svc.SetConfig(c.Request.Context(), key, req.Value, req.Description); err != nil {
		respondError(c, err)
		return
	}
	respondMsg(c, "config updated")
}

// ============================================================================
// 邀请码管理
// ============================================================================

// ListCodes 获取邀请码列表
func (h *ReferralHandler) ListCodes(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	opts := referralrepo.CodeListOptions{
		Page:     page,
		PageSize: pageSize,
		Keyword:  c.Query("keyword"),
	}

	if userID := c.Query("userId"); userID != "" {
		id, err := parseUint64(userID)
		if err == nil {
			opts.UserID = &id
		}
	}
	if refType := c.Query("type"); refType != "" {
		t := model.ReferralType(refType)
		opts.Type = &t
	}
	if isActive := c.Query("isActive"); isActive != "" {
		active := isActive == "true"
		opts.IsActive = &active
	}

	items, total, err := h.svc.ListCodes(c.Request.Context(), opts)
	if err != nil {
		respondError(c, err)
		return
	}

	respondList(c, items, &model.Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
	})
}

// GetCode 获取邀请码详情
func (h *ReferralHandler) GetCode(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	code, err := h.svc.GetCodeByID(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, code)
}

// createCodeRequest 创建邀请码请求
type createCodeRequest struct {
	UserID   uint64             `json:"userId" binding:"required"`
	Type     model.ReferralType `json:"type" binding:"required"`
	MaxUse   int                `json:"maxUse"`
	ExpireAt *time.Time         `json:"expireAt"`
}

// CreateCode 创建邀请码
func (h *ReferralHandler) CreateCode(c *gin.Context) {
	var req createCodeRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	code, err := h.svc.CreateCode(c.Request.Context(), referralservice.CreateCodeRequest{
		UserID:   req.UserID,
		Type:     req.Type,
		MaxUse:   req.MaxUse,
		ExpireAt: req.ExpireAt,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	respondCreated(c, code)
}

// updateCodeRequest 更新邀请码请求
type updateCodeRequest struct {
	IsActive *bool      `json:"isActive"`
	MaxUse   *int       `json:"maxUse"`
	ExpireAt *time.Time `json:"expireAt"`
}

// UpdateCode 更新邀请码
func (h *ReferralHandler) UpdateCode(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var req updateCodeRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	code, err := h.svc.UpdateCode(c.Request.Context(), referralservice.UpdateCodeRequest{
		ID:       id,
		IsActive: req.IsActive,
		MaxUse:   req.MaxUse,
		ExpireAt: req.ExpireAt,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	respondUpdated(c, code)
}

// DeleteCode 删除邀请码
func (h *ReferralHandler) DeleteCode(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	if err := h.svc.DeleteCode(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	respondDeleted(c)
}

// ============================================================================
// 推荐记录管理
// ============================================================================

// ListReferrals 获取推荐记录列表
func (h *ReferralHandler) ListReferrals(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	opts := referralrepo.ReferralListOptions{
		Page:     page,
		PageSize: pageSize,
	}

	if referrerID := c.Query("referrerId"); referrerID != "" {
		id, err := parseUint64(referrerID)
		if err == nil {
			opts.ReferrerID = &id
		}
	}
	if refereeID := c.Query("refereeId"); refereeID != "" {
		id, err := parseUint64(refereeID)
		if err == nil {
			opts.RefereeID = &id
		}
	}
	if refType := c.Query("type"); refType != "" {
		t := model.ReferralType(refType)
		opts.Type = &t
	}
	if status := c.Query("status"); status != "" {
		s := model.ReferralStatus(status)
		opts.Status = &s
	}

	items, total, err := h.svc.ListReferrals(c.Request.Context(), opts)
	if err != nil {
		respondError(c, err)
		return
	}

	respondList(c, items, &model.Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
	})
}

// GetReferral 获取推荐记录详情
func (h *ReferralHandler) GetReferral(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	referral, err := h.svc.GetReferralByID(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, referral)
}

// updateReferralStatusRequest 更新推荐状态请求
type updateReferralStatusRequest struct {
	Status model.ReferralStatus `json:"status" binding:"required"`
}

// UpdateReferralStatus 更新推荐状态
func (h *ReferralHandler) UpdateReferralStatus(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var req updateReferralStatusRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	if err := h.svc.UpdateReferralStatus(c.Request.Context(), id, req.Status); err != nil {
		respondError(c, err)
		return
	}
	respondMsg(c, "status updated")
}

// ============================================================================
// 奖励管理
// ============================================================================

// ListRewards 获取奖励记录列表
func (h *ReferralHandler) ListRewards(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	opts := referralrepo.RewardListOptions{
		Page:     page,
		PageSize: pageSize,
	}

	if userID := c.Query("userId"); userID != "" {
		id, err := parseUint64(userID)
		if err == nil {
			opts.UserID = &id
		}
	}
	if referralID := c.Query("referralId"); referralID != "" {
		id, err := parseUint64(referralID)
		if err == nil {
			opts.ReferralID = &id
		}
	}
	if rewardType := c.Query("type"); rewardType != "" {
		t := model.RewardType(rewardType)
		opts.Type = &t
	}
	if status := c.Query("status"); status != "" {
		s := model.ReferralRewardStatus(status)
		opts.Status = &s
	}

	items, total, err := h.svc.ListRewards(c.Request.Context(), opts)
	if err != nil {
		respondError(c, err)
		return
	}

	respondList(c, items, &model.Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
	})
}

// GetReward 获取奖励记录详情
func (h *ReferralHandler) GetReward(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	reward, err := h.svc.GetRewardByID(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, reward)
}

// IssueReward 发放奖励
func (h *ReferralHandler) IssueReward(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	if err := h.svc.IssueReward(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	respondMsg(c, "reward issued")
}

// failRewardRequest 标记奖励失败请求
type failRewardRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// FailReward 标记奖励发放失败
func (h *ReferralHandler) FailReward(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var req failRewardRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	if err := h.svc.FailReward(c.Request.Context(), id, req.Reason); err != nil {
		respondError(c, err)
		return
	}
	respondMsg(c, "reward marked as failed")
}

// ============================================================================
// 统计
// ============================================================================

// GetStats 获取推荐统计
func (h *ReferralHandler) GetStats(c *gin.Context) {
	stats, err := h.svc.GetReferralStats(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, stats)
}

// ============================================================================
// 批量操作 - 邀请码
// ============================================================================

// batchUpdateCodesStatusRequest 批量更新邀请码状态请求
type batchUpdateCodesStatusRequest struct {
	IDs      []uint64 `json:"ids" binding:"required,min=1"`
	IsActive bool     `json:"isActive" binding:"required"`
}

// BatchUpdateCodesStatus 批量更新邀请码状态
func (h *ReferralHandler) BatchUpdateCodesStatus(c *gin.Context) {
	var req batchUpdateCodesStatusRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	result, err := h.svc.BatchUpdateCodesStatus(c.Request.Context(), req.IDs, req.IsActive)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, result)
}

// batchDeleteRequest 批量删除请求
type batchDeleteRequest struct {
	IDs []uint64 `json:"ids" binding:"required,min=1"`
}

// BatchDeleteCodes 批量删除邀请码
func (h *ReferralHandler) BatchDeleteCodes(c *gin.Context) {
	var req batchDeleteRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	result, err := h.svc.BatchDeleteCodes(c.Request.Context(), req.IDs)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, result)
}

// ============================================================================
// 批量操作 - 推荐记录
// ============================================================================

// batchUpdateReferralsStatusRequest 批量更新推荐状态请求
type batchUpdateReferralsStatusRequest struct {
	IDs    []uint64                `json:"ids" binding:"required,min=1"`
	Status model.ReferralStatus    `json:"status" binding:"required"`
}

// BatchUpdateReferralsStatus 批量更新推荐状态
func (h *ReferralHandler) BatchUpdateReferralsStatus(c *gin.Context) {
	var req batchUpdateReferralsStatusRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	result, err := h.svc.BatchUpdateReferralsStatus(c.Request.Context(), req.IDs, req.Status)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, result)
}

// BatchDeleteReferrals 批量删除推荐记录
func (h *ReferralHandler) BatchDeleteReferrals(c *gin.Context) {
	var req batchDeleteRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	result, err := h.svc.BatchDeleteReferrals(c.Request.Context(), req.IDs)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, result)
}

// RegisterReferralRoutes 注册推荐管理路由
func RegisterReferralRoutes(r *gin.RouterGroup, svc *referralservice.Service, perm *middleware.PermissionMiddleware) {
	h := NewReferralHandler(svc)

	g := r.Group("/referrals")
	g.Use(perm.RequireAuth())
	{
		// 配置管理
		g.GET("/configs", perm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/referrals/configs"), h.ListConfigs)
		g.PUT("/configs/:key", perm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/referrals/configs/:key"), h.UpdateConfig)

		// 邀请码管理
		g.GET("/codes", perm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/referrals/codes"), h.ListCodes)
		g.GET("/codes/:id", perm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/referrals/codes/:id"), h.GetCode)
		g.POST("/codes", perm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/referrals/codes"), h.CreateCode)
		g.PUT("/codes/:id", perm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/referrals/codes/:id"), h.UpdateCode)
		g.DELETE("/codes/:id", perm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/referrals/codes/:id"), h.DeleteCode)
		g.PUT("/codes/batch/status", perm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/referrals/codes/batch/status"), h.BatchUpdateCodesStatus)
		g.DELETE("/codes/batch", perm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/referrals/codes/batch"), h.BatchDeleteCodes)

		// 推荐记录管理
		g.GET("", perm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/referrals"), h.ListReferrals)
		g.GET("/:id", perm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/referrals/:id"), h.GetReferral)
		g.PUT("/:id/status", perm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/referrals/:id/status"), h.UpdateReferralStatus)
		g.DELETE("/batch", perm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/referrals/batch"), h.BatchDeleteReferrals)
		g.PUT("/batch/status", perm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/referrals/batch/status"), h.BatchUpdateReferralsStatus)

		// 奖励管理
		g.GET("/rewards", perm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/referrals/rewards"), h.ListRewards)
		g.GET("/rewards/:id", perm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/referrals/rewards/:id"), h.GetReward)
		g.POST("/rewards/:id/issue", perm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/referrals/rewards/:id/issue"), h.IssueReward)
		g.POST("/rewards/:id/fail", perm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/referrals/rewards/:id/fail"), h.FailReward)

		// 统计
		g.GET("/stats", perm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/referrals/stats"), h.GetStats)
	}
}
