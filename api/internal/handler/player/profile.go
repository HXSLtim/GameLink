package player

import (
	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	serviceplayer "gamelink/internal/service/player"
	"gamelink/pkg/apierr"
)

// ApplyPlayerRequest 申请成为陪玩师请求（类型别名）
type ApplyPlayerRequest = serviceplayer.ApplyPlayerRequest

// ApplyPlayerResponse 申请陪玩师响应（类型别名）
type ApplyPlayerResponse = serviceplayer.ApplyPlayerResponse

// PlayerDetailResponse 陪玩师详情响应（类型别名）
type PlayerDetailResponse = serviceplayer.PlayerDetailResponse

// UpdatePlayerProfileRequest 更新陪玩师资料请求（类型别名）
type UpdatePlayerProfileRequest = serviceplayer.UpdatePlayerProfileRequest

// SetPlayerStatusRequest 设置在线状态请求（类型别名）
type SetPlayerStatusRequest = serviceplayer.SetPlayerStatusRequest

// ApplyPlayerResponseSwagger Swagger文档用的申请陪玩师响应类型
type ApplyPlayerResponseSwagger struct {
	PlayerID           uint64                   `json:"playerId"`
	VerificationStatus model.VerificationStatus `json:"verificationStatus"`
}

// RegisterProfileRoutes 注册陪玩师端资料管理路由
func RegisterProfileRoutes(router gin.IRouter, svc *serviceplayer.PlayerService, authMiddleware gin.HandlerFunc) {
	group := router.Group("")
	group.Use(authMiddleware) // 需要认证
	group.POST("/apply", func(c *gin.Context) { applyAsPlayerHandler(c, svc) })
	group.GET("/profile", func(c *gin.Context) { getPlayerProfileHandler(c, svc) })
	group.PUT("/profile", func(c *gin.Context) { updatePlayerProfileHandler(c, svc) })
	group.PUT("/status", func(c *gin.Context) { setPlayerStatusHandler(c, svc) })
}

// applyAsPlayerHandler 申请成为陪玩师
// @Summary      申请成为陪玩师
// @Description  用户申请成为陪玩师
// @Tags         Player - Profile
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  ApplyPlayerRequest  true  "申请信息"
// @Success      200      {object}  ApplyPlayerResponseSwagger
// @Failure      400      {object}  apierr.APIError
// @Failure      401      {object}  apierr.APIError
// @Failure      409      {object}  apierr.APIError
// @Failure      500      {object}  apierr.APIError
// @Router       /player/apply [post]
func applyAsPlayerHandler(c *gin.Context, svc *serviceplayer.PlayerService) {
	userID := getUserIDFromContext(c)

	var req serviceplayer.ApplyPlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest(apierr.ErrInvalidJSONPayload).WithDetails(err.Error()))
		return
	}

	resp, err := svc.ApplyAsPlayer(c.Request.Context(), userID, req)
	if err != nil {
		if err == serviceplayer.ErrAlreadyPlayer {
			respondAPIError(c, apierr.BadRequest(err.Error()))
			return
		}
		respondAPIError(c, apierr.InternalError("申请提交失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "申请提交成功", *resp)
}

// getPlayerProfileHandler 获取陪玩师资料
// @Summary      获取陪玩师资料
// @Description  获取陪玩师个人资料
// @Tags         Player - Profile
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  PlayerDetailResponse
// @Failure      400  {object}  apierr.APIError
// @Failure      401  {object}  apierr.APIError
// @Failure      404  {object}  apierr.APIError
// @Failure      500  {object}  apierr.APIError
// @Router       /player/profile [get]
func getPlayerProfileHandler(c *gin.Context, svc *serviceplayer.PlayerService) {
	userID := getUserIDFromContext(c)

	resp, err := svc.GetPlayerProfile(c.Request.Context(), userID)
	if err != nil {
		if err == serviceplayer.ErrNotFound {
			respondAPIError(c, apierr.NotFound(err.Error()))
			return
		}
		respondAPIError(c, apierr.InternalError("获取陪玩师资料失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "OK", *resp)
}

// updatePlayerProfileHandler 更新陪玩师资料
// @Summary      更新陪玩师资料
// @Description  更新陪玩师个人资料
// @Tags         Player - Profile
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  UpdatePlayerProfileRequest  true  "更新信息"
// @Success      200      {object}  model.SuccessResponse
// @Failure      400      {object}  apierr.APIError
// @Failure      401      {object}  apierr.APIError
// @Failure      404      {object}  apierr.APIError
// @Failure      422      {object}  apierr.APIError
// @Failure      500      {object}  apierr.APIError
// @Router       /player/profile [put]
func updatePlayerProfileHandler(c *gin.Context, svc *serviceplayer.PlayerService) {
	userID := getUserIDFromContext(c)

	var req serviceplayer.UpdatePlayerProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest(apierr.ErrInvalidJSONPayload).WithDetails(err.Error()))
		return
	}

	if err := svc.UpdatePlayerProfile(c.Request.Context(), userID, req); err != nil {
		if err == serviceplayer.ErrNotFound {
			respondAPIError(c, apierr.NotFound(err.Error()))
			return
		}
		if err == serviceplayer.ErrValidation {
			respondAPIError(c, apierr.BadRequest(err.Error()))
			return
		}
		respondAPIError(c, apierr.InternalError("更新资料失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "资料更新成功", struct{}{})
}

// setPlayerStatusHandler 设置陪玩师在线状态
// @Summary      设置陪玩师在线状态
// @Description  设置陪玩师在线/离线状态
// @Tags         Player - Profile
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  serviceplayer.SetPlayerStatusRequest  true  "状态请求"
// @Success      200      {object}  model.SuccessResponse
// @Failure      400      {object}  apierr.APIError
// @Failure      401      {object}  apierr.APIError
// @Failure      500      {object}  apierr.APIError
// @Router       /player/status [put]
func setPlayerStatusHandler(c *gin.Context, svc *serviceplayer.PlayerService) {
	userID := getUserIDFromContext(c)

	var req serviceplayer.SetPlayerStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest(apierr.ErrInvalidJSONPayload).WithDetails(err.Error()))
		return
	}

	// 转换布尔值为状态字符串
	status := "offline"
	if req.Online {
		status = "online"
	}

	if err := svc.SetPlayerOnlineStatus(c.Request.Context(), userID, status); err != nil {
		respondAPIError(c, apierr.InternalError("状态更新失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "状态更新成功", struct{}{})
}
