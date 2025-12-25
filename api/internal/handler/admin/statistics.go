package admin

import (
	"net/http"
	"strconv"

	"gamelink/internal/service/statistics"

	"github.com/gin-gonic/gin"
)

// StatisticsHandler 统计接口处理器
type StatisticsHandler struct {
	svc       *statistics.Service
	evaluator *statistics.TagEvaluator
}

// NewStatisticsHandler 创建统计处理器
func NewStatisticsHandler(svc *statistics.Service, evaluator *statistics.TagEvaluator) *StatisticsHandler {
	return &StatisticsHandler{svc: svc, evaluator: evaluator}
}

// RefreshUserStatistics 刷新用户统计
// @Summary 刷新用户统计
// @Tags 统计管理
// @Param id path int true "用户ID"
// @Success 200 {object} map[string]interface{}
// @Router /admin/statistics/user/{id}/refresh [post]
func (h *StatisticsHandler) RefreshUserStatistics(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的用户ID"})
		return
	}

	if err := h.svc.UpdateUserStatistics(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "用户统计已刷新"})
}

// RefreshPlayerStatistics 刷新陪玩师统计
// @Summary 刷新陪玩师统计
// @Tags 统计管理
// @Param id path int true "陪玩师ID"
// @Success 200 {object} map[string]interface{}
// @Router /admin/statistics/player/{id}/refresh [post]
func (h *StatisticsHandler) RefreshPlayerStatistics(c *gin.Context) {
	playerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的陪玩师ID"})
		return
	}

	if err := h.svc.UpdatePlayerStatistics(c.Request.Context(), playerID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "陪玩师统计已刷新"})
}

// RefreshAllStatistics 刷新所有统计
// @Summary 刷新所有统计数据
// @Tags 统计管理
// @Success 200 {object} map[string]interface{}
// @Router /admin/statistics/refresh-all [post]
func (h *StatisticsHandler) RefreshAllStatistics(c *gin.Context) {
	ctx := c.Request.Context()

	// 更新所有用户统计
	if err := h.svc.UpdateAllUserStatistics(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新用户统计失败: " + err.Error()})
		return
	}

	// 更新所有陪玩师统计
	if err := h.svc.UpdateAllPlayerStatistics(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新陪玩师统计失败: " + err.Error()})
		return
	}

	// 更新所有服务项目统计
	if err := h.svc.UpdateAllServiceItemStatistics(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新服务项目统计失败: " + err.Error()})
		return
	}

	// 更新所有游戏统计
	if err := h.svc.UpdateAllGameStatistics(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新游戏统计失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "所有统计数据已刷新"})
}

// SyncUserTags 同步用户标签
// @Summary 根据统计数据同步用户标签
// @Tags 统计管理
// @Param id path int true "用户ID"
// @Success 200 {object} map[string]interface{}
// @Router /admin/statistics/user/{id}/sync-tags [post]
func (h *StatisticsHandler) SyncUserTags(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的用户ID"})
		return
	}

	if err := h.evaluator.SyncUserTags(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "用户标签已同步"})
}

// EvaluateUserTags 评估用户标签
// @Summary 评估用户应该拥有的标签（不实际修改）
// @Tags 统计管理
// @Param id path int true "用户ID"
// @Success 200 {object} map[string]interface{}
// @Router /admin/statistics/user/{id}/evaluate-tags [get]
func (h *StatisticsHandler) EvaluateUserTags(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的用户ID"})
		return
	}

	tagIDs, err := h.evaluator.EvaluateUserTags(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"tagIds": tagIDs}})
}
