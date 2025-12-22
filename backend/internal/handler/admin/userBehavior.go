/**
 * @file user behavior handler
 * @description 用户行为分析API接口
 */

package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"

	adminservice "gamelink/internal/service/admin"
	"gamelink/pkg/apierr"
)

// UserBehaviorHandler 处理用户行为分析接口
type UserBehaviorHandler struct {
	statsService *adminservice.StatsService
}

// NewUserBehaviorHandler 创建Handler
func NewUserBehaviorHandler(statsService *adminservice.StatsService) *UserBehaviorHandler {
	return &UserBehaviorHandler{statsService: statsService}
}

// GetBehaviorStats
// @Summary      获取用户行为统计
// @Description  获取DAU、平均在线时长、人均消费等用户行为统计数据
// @Tags         Admin/UserBehavior
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  model.APIResponse[adminservice.UserBehaviorStatsResponse]
// @Router       /admin/users/behavior/stats [get]
func (h *UserBehaviorHandler) GetBehaviorStats(c *gin.Context) {
	stats, err := h.statsService.UserBehaviorStats(c.Request.Context())
	if err != nil {
		respondError(c, apierr.InternalError("获取用户行为统计失败").WithDetails(err.Error()))
		return
	}
	respondSuccess(c, stats)
}

// GetActivityTrend
// @Summary      获取用户活动趋势
// @Description  获取最近N天的用户活动趋势数据
// @Tags         Admin/UserBehavior
// @Security     BearerAuth
// @Produce      json
// @Param        days  query  int  false  "统计天数（默认7天）"
// @Success      200  {object}  model.SuccessResponse
// @Router       /admin/users/behavior/trend [get]
func (h *UserBehaviorHandler) GetActivityTrend(c *gin.Context) {
	days := 7
	if daysParam := c.Query("days"); daysParam != "" {
		if parsed, err := strconv.Atoi(daysParam); err == nil && parsed > 0 {
			days = parsed
		}
	}

	trend, err := h.statsService.UserActivityTrend(c.Request.Context(), days)
	if err != nil {
		respondError(c, apierr.InternalError("获取用户活动趋势失败").WithDetails(err.Error()))
		return
	}
	respondSuccess(c, trend)
}

// GetUserDistribution
// @Summary      获取用户分布
// @Description  获取用户地域分布、年龄分布等统计数据
// @Tags         Admin/UserBehavior
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  model.APIResponse[adminservice.UserDistributionResponse]
// @Router       /admin/users/behavior/distribution [get]
func (h *UserBehaviorHandler) GetUserDistribution(c *gin.Context) {
	distribution, err := h.statsService.UserDistribution(c.Request.Context())
	if err != nil {
		respondError(c, apierr.InternalError("获取用户分布失败").WithDetails(err.Error()))
		return
	}
	respondSuccess(c, distribution)
}
