package player

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/service/player"
)

// RegisterProfileRoutes 注册陪玩师端资料管理路由
func RegisterProfileRoutes(router gin.IRouter, svc *player.PlayerService, authMiddleware gin.HandlerFunc) {
    group := router.Group("/player")
    group.Use(authMiddleware) // 需要认证
    group.POST("/apply", func(c *gin.Context) { applyAsPlayerHandler(c, svc) })
    group.GET("/profile", func(c *gin.Context) { getPlayerProfileHandler(c, svc) })
    group.PUT("/profile", func(c *gin.Context) { updatePlayerProfileHandler(c, svc) })
    group.PUT("/status", func(c *gin.Context) { setPlayerStatusHandler(c, svc) })
}

// applyAsPlayerHandler 申请成为陪玩�?// @Summary      申请成为陪玩�?// @Description  用户申请成为陪玩�?// @Tags         Player - Profile
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                       true  "Bearer {token}"
// @Param        request        body      player.ApplyPlayerRequest    true  "申请信息"
// @Success      200            {object}  model.APIResponse[player.ApplyPlayerResponse]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /player/apply [post]
func applyAsPlayerHandler(c *gin.Context, svc *player.PlayerService) {
	userID := getUserIDFromContext(c)

	var req player.ApplyPlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := svc.ApplyAsPlayer(c.Request.Context(), userID, req)
	if err != nil {
		if err == player.ErrAlreadyPlayer {
			respondError(c, http.StatusBadRequest, err.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[player.ApplyPlayerResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "申请提交成功",
		Data:    *resp,
	})
}

// getPlayerProfileHandler 获取陪玩师资�?// @Summary      获取陪玩师资�?// @Description  获取当前用户的陪玩师资料
// @Tags         Player - Profile
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Success      200            {object}  model.APIResponse[player.PlayerDetailResponse]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Failure      404            {object}  model.APIResponse[any]
// @Router       /player/profile [get]
func getPlayerProfileHandler(c *gin.Context, svc *player.PlayerService) {
	userID := getUserIDFromContext(c)

	resp, err := svc.GetPlayerProfile(c.Request.Context(), userID)
	if err != nil {
		if err == player.ErrNotFound {
			respondError(c, http.StatusNotFound, err.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[player.PlayerDetailResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    *resp,
	})
}

// updatePlayerProfileHandler 更新陪玩师资�?// @Summary      更新陪玩师资�?// @Description  更新陪玩师资�?// @Tags         Player - Profile
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                              true  "Bearer {token}"
// @Param        request        body      player.UpdatePlayerProfileRequest   true  "更新信息"
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /player/profile [put]
func updatePlayerProfileHandler(c *gin.Context, svc *player.PlayerService) {
	userID := getUserIDFromContext(c)

	var req player.UpdatePlayerProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := svc.UpdatePlayerProfile(c.Request.Context(), userID, req); err != nil {
		if err == player.ErrNotFound {
			respondError(c, http.StatusNotFound, err.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "资料更新成功",
	})
}

// setPlayerStatusHandler 设置在线状�?// @Summary      设置在线状�?// @Description  设置陪玩师在�?离线状�?// @Tags         Player - Profile
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                          true  "Bearer {token}"
// @Param        request        body      player.SetPlayerStatusRequest   true  "在线状�?
// @Success      200            {object}  model.APIResponse[any]
// @Failure      400            {object}  model.APIResponse[any]
// @Failure      401            {object}  model.APIResponse[any]
// @Router       /player/status [put]
func setPlayerStatusHandler(c *gin.Context, svc *player.PlayerService) {
	userID := getUserIDFromContext(c)

	var req player.SetPlayerStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := svc.SetPlayerOnlineStatus(c.Request.Context(), userID, req.Online); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "状态更新成功",
	})
}
