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
	activityrepo "gamelink/internal/repository/activity"
	activityservice "gamelink/internal/service/activity"
	"gamelink/pkg/apierr"
)

// ActivityHandler 活动管理处理器
type ActivityHandler struct {
	svc *activityservice.Service
}

// NewActivityHandler 创建活动处理器
func NewActivityHandler(svc *activityservice.Service) *ActivityHandler {
	return &ActivityHandler{svc: svc}
}

// ============================================================================
// 活动管理
// ============================================================================

// ActivityCreateRequest 创建活动请求
type ActivityCreateRequest struct {
	Name          string               `json:"name" binding:"required,max=128"`
	Description   string               `json:"description"`
	Type          model.ActivityType   `json:"type" binding:"required"`
	CoverURL      string               `json:"coverUrl"`
	BannerURL     string               `json:"bannerUrl"`
	PreheatAt     *time.Time           `json:"preheatAt"`
	StartAt       time.Time            `json:"startAt" binding:"required"`
	EndAt         time.Time            `json:"endAt" binding:"required"`
	TotalLimit    int                  `json:"totalLimit"`
	DailyLimit    int                  `json:"dailyLimit"`
	PerUserLimit  int                  `json:"perUserLimit"`
	AllowVipStack bool                 `json:"allowVipStack"`
	Rules         string               `json:"rules"`
	SortOrder     int                  `json:"sortOrder"`
	IsVisible     bool                 `json:"isVisible"`
	Status        model.ActivityStatus `json:"status"`
}

// ActivityUpdateRequest 更新活动请求
type ActivityUpdateRequest struct {
	Name          string             `json:"name" binding:"required,max=128"`
	Description   string             `json:"description"`
	Type          model.ActivityType `json:"type" binding:"required"`
	CoverURL      string             `json:"coverUrl"`
	BannerURL     string             `json:"bannerUrl"`
	PreheatAt     *time.Time         `json:"preheatAt"`
	StartAt       time.Time          `json:"startAt" binding:"required"`
	EndAt         time.Time          `json:"endAt" binding:"required"`
	TotalLimit    int                `json:"totalLimit"`
	DailyLimit    int                `json:"dailyLimit"`
	PerUserLimit  int                `json:"perUserLimit"`
	AllowVipStack bool               `json:"allowVipStack"`
	Rules         string             `json:"rules"`
	SortOrder     int                `json:"sortOrder"`
	IsVisible     bool               `json:"isVisible"`
}

// ListActivities 获取活动列表
func (h *ActivityHandler) ListActivities(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	var activityType *model.ActivityType
	if v := strings.TrimSpace(c.Query("type")); v != "" {
		t := model.ActivityType(v)
		activityType = &t
	}

	var status *model.ActivityStatus
	if v := strings.TrimSpace(c.Query("status")); v != "" {
		s := model.ActivityStatus(v)
		status = &s
	}

	var isVisible *bool
	if v := strings.TrimSpace(c.Query("isVisible")); v != "" {
		b := v == "true" || v == "1"
		isVisible = &b
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

	opts := activityrepo.ActivityListOptions{
		Page:      page,
		PageSize:  pageSize,
		Keyword:   c.Query("keyword"),
		Type:      activityType,
		Status:    status,
		IsVisible: isVisible,
		StartTime: startTime,
		EndTime:   endTime,
	}

	activities, total, err := h.svc.ListActivities(c.Request.Context(), opts)
	if err != nil {
		respondError(c, apierr.InternalError("获取活动列表失败").WithDetails(err.Error()))
		return
	}

	respondList(c, activities, &model.Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
	})
}

// GetActivity 获取活动详情
func (h *ActivityHandler) GetActivity(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	activity, err := h.svc.GetActivity(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("活动不存在"))
			return
		}
		respondError(c, apierr.InternalError("获取活动失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, activity)
}

// CreateActivity 创建活动
func (h *ActivityHandler) CreateActivity(c *gin.Context) {
	var req ActivityCreateRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	activity := &model.Activity{
		Name:          req.Name,
		Description:   req.Description,
		Type:          req.Type,
		Status:        req.Status,
		CoverURL:      req.CoverURL,
		BannerURL:     req.BannerURL,
		PreheatAt:     req.PreheatAt,
		StartAt:       req.StartAt,
		EndAt:         req.EndAt,
		TotalLimit:    req.TotalLimit,
		DailyLimit:    req.DailyLimit,
		PerUserLimit:  req.PerUserLimit,
		AllowVipStack: req.AllowVipStack,
		Rules:         req.Rules,
		SortOrder:     req.SortOrder,
		IsVisible:     req.IsVisible,
	}

	if err := h.svc.CreateActivity(c.Request.Context(), activity); err != nil {
		respondError(c, apierr.BadRequest(err.Error()))
		return
	}

	respondCreated(c, activity)
}

// UpdateActivity 更新活动
func (h *ActivityHandler) UpdateActivity(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var req ActivityUpdateRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	activity := &model.Activity{
		Name:          req.Name,
		Description:   req.Description,
		Type:          req.Type,
		CoverURL:      req.CoverURL,
		BannerURL:     req.BannerURL,
		PreheatAt:     req.PreheatAt,
		StartAt:       req.StartAt,
		EndAt:         req.EndAt,
		TotalLimit:    req.TotalLimit,
		DailyLimit:    req.DailyLimit,
		PerUserLimit:  req.PerUserLimit,
		AllowVipStack: req.AllowVipStack,
		Rules:         req.Rules,
		SortOrder:     req.SortOrder,
		IsVisible:     req.IsVisible,
	}
	activity.ID = id

	if err := h.svc.UpdateActivity(c.Request.Context(), activity); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("活动不存在"))
			return
		}
		respondError(c, apierr.BadRequest(err.Error()))
		return
	}

	respondUpdated(c, activity)
}

// DeleteActivity 删除活动
func (h *ActivityHandler) DeleteActivity(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	if err := h.svc.DeleteActivity(c.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("活动不存在"))
			return
		}
		respondError(c, apierr.BadRequest(err.Error()))
		return
	}

	respondDeleted(c)
}

// UpdateStatusRequest 更新状态请求
type UpdateStatusRequest struct {
	Status model.ActivityStatus `json:"status" binding:"required"`
}

// UpdateActivityStatus 更新活动状态
func (h *ActivityHandler) UpdateActivityStatus(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var req UpdateStatusRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	if err := h.svc.UpdateActivityStatus(c.Request.Context(), id, req.Status); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("活动不存在"))
			return
		}
		respondError(c, apierr.BadRequest(err.Error()))
		return
	}

	respondSuccess(c, gin.H{"message": "状态更新成功"})
}

// ============================================================================
// 活动奖励管理
// ============================================================================

// RewardCreateRequest 创建奖励请求
type RewardCreateRequest struct {
	ActivityID       uint64 `json:"activityId" binding:"required"`
	CouponTemplateID uint64 `json:"couponTemplateId" binding:"required"`
	CouponCount      int    `json:"couponCount" binding:"required,gt=0"`
	Probability      int    `json:"probability" binding:"required,min=1,max=100"`
	TotalStock       int    `json:"totalStock"`
	SortOrder        int    `json:"sortOrder"`
}

// RewardUpdateRequest 更新奖励请求
type RewardUpdateRequest struct {
	CouponTemplateID uint64 `json:"couponTemplateId" binding:"required"`
	CouponCount      int    `json:"couponCount" binding:"required,gt=0"`
	Probability      int    `json:"probability" binding:"required,min=1,max=100"`
	TotalStock       int    `json:"totalStock"`
	SortOrder        int    `json:"sortOrder"`
}

// GetRewardsByActivity 获取活动的奖励列表
func (h *ActivityHandler) GetRewardsByActivity(c *gin.Context) {
	activityID, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	rewards, err := h.svc.GetRewardsByActivityID(c.Request.Context(), activityID)
	if err != nil {
		respondError(c, apierr.InternalError("获取奖励列表失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, rewards)
}

// GetReward 获取奖励详情
func (h *ActivityHandler) GetReward(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "rewardId")
	if !ok {
		return
	}

	reward, err := h.svc.GetReward(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("奖励不存在"))
			return
		}
		respondError(c, apierr.InternalError("获取奖励失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, reward)
}

// CreateReward 创建奖励
func (h *ActivityHandler) CreateReward(c *gin.Context) {
	var req RewardCreateRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	reward := &model.ActivityReward{
		ActivityID:       req.ActivityID,
		CouponTemplateID: req.CouponTemplateID,
		CouponCount:      req.CouponCount,
		Probability:      req.Probability,
		TotalStock:       req.TotalStock,
		SortOrder:        req.SortOrder,
	}

	if err := h.svc.CreateReward(c.Request.Context(), reward); err != nil {
		respondError(c, apierr.BadRequest(err.Error()))
		return
	}

	respondCreated(c, reward)
}

// UpdateReward 更新奖励
func (h *ActivityHandler) UpdateReward(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "rewardId")
	if !ok {
		return
	}

	var req RewardUpdateRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	reward := &model.ActivityReward{
		CouponTemplateID: req.CouponTemplateID,
		CouponCount:      req.CouponCount,
		Probability:      req.Probability,
		TotalStock:       req.TotalStock,
		SortOrder:        req.SortOrder,
	}
	reward.ID = id

	if err := h.svc.UpdateReward(c.Request.Context(), reward); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("奖励不存在"))
			return
		}
		respondError(c, apierr.BadRequest(err.Error()))
		return
	}

	respondUpdated(c, reward)
}

// DeleteReward 删除奖励
func (h *ActivityHandler) DeleteReward(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "rewardId")
	if !ok {
		return
	}

	if err := h.svc.DeleteReward(c.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("奖励不存在"))
			return
		}
		respondError(c, apierr.InternalError("删除奖励失败").WithDetails(err.Error()))
		return
	}

	respondDeleted(c)
}

// ============================================================================
// 参与记录管理
// ============================================================================

// ListParticipations 获取参与记录列表
func (h *ActivityHandler) ListParticipations(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	var activityID *uint64
	if v := strings.TrimSpace(c.Query("activityId")); v != "" {
		id, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			respondBadRequest(c, "无效的活动ID")
			return
		}
		activityID = &id
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

	opts := activityrepo.ParticipationListOptions{
		Page:       page,
		PageSize:   pageSize,
		ActivityID: activityID,
		UserID:     userID,
		StartTime:  startTime,
		EndTime:    endTime,
	}

	participations, total, err := h.svc.ListParticipations(c.Request.Context(), opts)
	if err != nil {
		respondError(c, apierr.InternalError("获取参与记录失败").WithDetails(err.Error()))
		return
	}

	respondList(c, participations, &model.Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
	})
}

// ============================================================================
// 统计
// ============================================================================

// GetActivityStats 获取活动统计
func (h *ActivityHandler) GetActivityStats(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	stats, err := h.svc.GetActivityStats(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("活动不存在"))
			return
		}
		respondError(c, apierr.InternalError("获取统计失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, stats)
}

// GetAllActivityStats 获取所有活动统计概览
func (h *ActivityHandler) GetAllActivityStats(c *gin.Context) {
	stats, err := h.svc.GetAllActivityStats(c.Request.Context())
	if err != nil {
		respondError(c, apierr.InternalError("获取统计失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, stats)
}

// ============================================================================
// 路由注册
// ============================================================================

// RegisterActivityRoutes 注册活动管理路由
func RegisterActivityRoutes(rg *gin.RouterGroup, svc *activityservice.Service, pm *middleware.PermissionMiddleware) {
	h := NewActivityHandler(svc)

	activityGroup := rg.Group("/activities")
	activityGroup.Use(pm.RequireAuth())
	{
		// 活动管理
		activityGroup.GET("", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/activities"), h.ListActivities)
		activityGroup.GET("/stats", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/activities/stats"), h.GetAllActivityStats)
		activityGroup.GET("/:id", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/activities/:id"), h.GetActivity)
		activityGroup.GET("/:id/stats", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/activities/:id/stats"), h.GetActivityStats)
		activityGroup.POST("", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/activities"), h.CreateActivity)
		activityGroup.PUT("/:id", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/activities/:id"), h.UpdateActivity)
		activityGroup.DELETE("/:id", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/activities/:id"), h.DeleteActivity)
		activityGroup.PUT("/:id/status", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/activities/:id/status"), h.UpdateActivityStatus)

		// 活动奖励管理
		activityGroup.GET("/:id/rewards", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/activities/:id/rewards"), h.GetRewardsByActivity)
		activityGroup.GET("/:id/rewards/:rewardId", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/activities/:id/rewards/:rewardId"), h.GetReward)
		activityGroup.POST("/rewards", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/activities/rewards"), h.CreateReward)
		activityGroup.PUT("/rewards/:rewardId", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/activities/rewards/:rewardId"), h.UpdateReward)
		activityGroup.DELETE("/rewards/:rewardId", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/activities/rewards/:rewardId"), h.DeleteReward)

		// 参与记录
		activityGroup.GET("/participations", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/activities/participations"), h.ListParticipations)
	}
}
