package player

import (
	"github.com/gin-gonic/gin"

	"gamelink/internal/service/player"
	"gamelink/pkg/apierr"
)

// RegisterStatusRoutes 注册陪玩师在线状态路由
func RegisterStatusRoutes(router gin.IRouter, svc *player.PlayerService, authMiddleware gin.HandlerFunc) {
	group := router.Group("/online-status")
	group.Use(authMiddleware)
	group.GET("", func(c *gin.Context) { getOnlineStatusHandler(c, svc) })
	group.PUT("", func(c *gin.Context) { updateOnlineStatusHandler(c, svc) })
}

// OnlineStatus 在线状态枚举
type OnlineStatus string

const (
	StatusOffline OnlineStatus = "offline" // 离线
	StatusOnline  OnlineStatus = "online"  // 在线
	StatusBusy    OnlineStatus = "busy"    // 忙碌（服务中）
)

// OnlineStatusResponse 在线状态响应
type OnlineStatusResponse struct {
	Status OnlineStatus `json:"status"` // offline/online/busy
	Online bool         `json:"online"` // 兼容旧版本
}

// UpdateOnlineStatusRequest 更新在线状态请求
type UpdateOnlineStatusRequest struct {
	Status OnlineStatus `json:"status"` // offline/online/busy
	Online *bool        `json:"online"` // 兼容旧版本（deprecated）
}

// getOnlineStatusHandler 获取在线状态
// @Summary      获取在线状态
// @Description  获取当前陪玩师的在线状态（offline/online/busy）
// @Tags         Player - Status
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  SuccessResponseSwagger
// @Failure      401  {object}  apierr.APIError
// @Failure      404  {object}  apierr.APIError
// @Failure      500  {object}  apierr.APIError
// @Router       /player/online-status [get]
func getOnlineStatusHandler(c *gin.Context, svc *player.PlayerService) {
	userID := getUserIDFromContext(c)

	status, err := svc.GetPlayerOnlineStatusByUserID(c.Request.Context(), userID)
	if err != nil {
		if err == player.ErrNotFound {
			respondAPIError(c, apierr.NotFound("陪玩师资料不存在"))
			return
		}
		respondAPIError(c, apierr.InternalError("获取在线状态失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "OK", OnlineStatusResponse{
		Status: OnlineStatus(status),
		Online: status != "offline",
	})
}

// updateOnlineStatusHandler 更新在线状态
// @Summary      更新在线状态
// @Description  更新当前陪玩师的在线状态（offline/online/busy）
// @Tags         Player - Status
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body  UpdateOnlineStatusRequest  true  "在线状态"
// @Success      200  {object}  SuccessResponseSwagger
// @Failure      400  {object}  apierr.APIError
// @Failure      401  {object}  apierr.APIError
// @Failure      404  {object}  apierr.APIError
// @Failure      500  {object}  apierr.APIError
// @Router       /player/online-status [put]
func updateOnlineStatusHandler(c *gin.Context, svc *player.PlayerService) {
	userID := getUserIDFromContext(c)

	var req UpdateOnlineStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("无效的请求参数").WithDetails(err.Error()))
		return
	}

	// 兼容旧版本：如果使用 online 字段
	status := req.Status
	if status == "" && req.Online != nil {
		if *req.Online {
			status = StatusOnline
		} else {
			status = StatusOffline
		}
	}

	// 验证状态值
	if status != StatusOffline && status != StatusOnline && status != StatusBusy {
		respondAPIError(c, apierr.BadRequest("无效的状态值，支持: offline/online/busy"))
		return
	}

	if err := svc.SetPlayerOnlineStatus(c.Request.Context(), userID, string(status)); err != nil {
		if err == player.ErrNotFound {
			respondAPIError(c, apierr.NotFound("陪玩师资料不存在"))
			return
		}
		respondAPIError(c, apierr.InternalError("更新在线状态失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "在线状态已更新", OnlineStatusResponse{
		Status: status,
		Online: status != StatusOffline,
	})
}
