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

// OnlineStatusResponse 在线状态响应
type OnlineStatusResponse struct {
	Online bool `json:"online"`
}

// UpdateOnlineStatusRequest 更新在线状态请求
type UpdateOnlineStatusRequest struct {
	Online bool `json:"online"`
}

// getOnlineStatusHandler 获取在线状态
// @Summary      获取在线状态
// @Description  获取当前陪玩师的在线状态
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

	online, err := svc.GetPlayerOnlineStatusByUserID(c.Request.Context(), userID)
	if err != nil {
		if err == player.ErrNotFound {
			respondAPIError(c, apierr.NotFound("陪玩师资料不存在"))
			return
		}
		respondAPIError(c, apierr.InternalError("获取在线状态失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "OK", OnlineStatusResponse{Online: online})
}

// updateOnlineStatusHandler 更新在线状态
// @Summary      更新在线状态
// @Description  更新当前陪玩师的在线状态
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

	if err := svc.SetPlayerOnlineStatus(c.Request.Context(), userID, req.Online); err != nil {
		if err == player.ErrNotFound {
			respondAPIError(c, apierr.NotFound("陪玩师资料不存在"))
			return
		}
		respondAPIError(c, apierr.InternalError("更新在线状态失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "在线状态已更新", OnlineStatusResponse{Online: req.Online})
}
