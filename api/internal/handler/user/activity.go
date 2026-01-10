package user

import (
	"errors"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	"gamelink/internal/repository"
	activityservice "gamelink/internal/service/activity"
	"gamelink/pkg/apierr"
)

// ActivityHandler 用户端活动处理器
type ActivityHandler struct {
	svc *activityservice.Service
}

// NewActivityHandler 创建活动处理器
func NewActivityHandler(svc *activityservice.Service) *ActivityHandler {
	return &ActivityHandler{svc: svc}
}

// ListActivities 获取可见的活动列表
func (h *ActivityHandler) ListActivities(c *gin.Context) {
	activities, err := h.svc.GetVisibleActivities(c.Request.Context())
	if err != nil {
		resp.Error(c, apierr.InternalError("获取活动列表失败").WithDetails(err.Error()))
		return
	}

	resp.OK(c, activities)
}

// GetActivity 获取活动详情
func (h *ActivityHandler) GetActivity(c *gin.Context) {
	id, ok := resp.ParseIDOrFail(c, "id")
	if !ok {
		return
	}

	activity, err := h.svc.GetActivity(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			resp.Error(c, apierr.NotFound("活动不存在"))
			return
		}
		resp.Error(c, apierr.InternalError("获取活动失败").WithDetails(err.Error()))
		return
	}

	// 检查活动是否可见
	if !activity.IsVisible {
		resp.Error(c, apierr.NotFound("活动不存在"))
		return
	}

	resp.OK(c, activity)
}

// ParticipateRequest 参与活动请求
type ParticipateRequest struct {
	RewardID uint64 `json:"rewardId" binding:"required"`
}

// Participate 参与活动
func (h *ActivityHandler) Participate(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	activityID, ok := resp.ParseIDOrFail(c, "id")
	if !ok {
		return
	}

	var req ParticipateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, apierr.BadRequest("参数错误").WithDetails(err.Error()))
		return
	}

	// 获取客户端IP
	clientIP := c.ClientIP()

	participation, err := h.svc.ParticipateActivity(c.Request.Context(), userID, activityID, req.RewardID, clientIP)
	if err != nil {
		resp.Error(c, apierr.BadRequest(err.Error()))
		return
	}

	resp.OK(c, participation)
}

// GetMyParticipations 获取我的参与记录
func (h *ActivityHandler) GetMyParticipations(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	participations, err := h.svc.GetUserParticipations(c.Request.Context(), userID, 50)
	if err != nil {
		resp.Error(c, apierr.InternalError("获取参与记录失败").WithDetails(err.Error()))
		return
	}

	resp.OK(c, participations)
}

// RegisterActivityRoutes 注册用户端活动路由
func RegisterActivityRoutes(rg *gin.RouterGroup, svc *activityservice.Service, authMiddleware gin.HandlerFunc) {
	h := NewActivityHandler(svc)

	activityGroup := rg.Group("/activities")
	{
		// 公开接口（无需登录）
		activityGroup.GET("", h.ListActivities)
		activityGroup.GET("/:id", h.GetActivity)

		// 需要登录的接口
		activityGroup.POST("/:id/participate", authMiddleware, h.Participate)
		activityGroup.GET("/my/participations", authMiddleware, h.GetMyParticipations)
	}
}
