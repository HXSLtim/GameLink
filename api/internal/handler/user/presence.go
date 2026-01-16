package user

import (
	"errors"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	presenceservice "gamelink/internal/service/presence"
	"gamelink/pkg/apierr"
)

// PresenceHandler 陪玩师在线状态处理器
type PresenceHandler struct {
	svc *presenceservice.Service
}

// NewPresenceHandler 创建在线状态处理器
func NewPresenceHandler(svc *presenceservice.Service) *PresenceHandler {
	return &PresenceHandler{svc: svc}
}

// GetMyPresence 获取自己的在线状态
// @Summary 获取自己的在线状态
// @Tags 用户-在线状态
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.PlayerPresence
// @Router /user/presence [get]
func (h *PresenceHandler) GetMyPresence(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	// 需要获取玩家ID（假设用户ID和玩家ID关联）
	presence, err := h.svc.GetOrCreatePresence(c.Request.Context(), userID)
	if err != nil {
		resp.Error(c, apierr.InternalError("获取在线状态失败").WithDetails(err.Error()))
		return
	}

	resp.OK(c, presence)
}

// UpdateMyPresence 更新自己的在线状态
// @Summary 更新自己的在线状态
// @Tags 用户-在线状态
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body UpdatePresenceRequest true "状态更新请求"
// @Success 200 {object} model.PlayerPresence
// @Router /user/presence [put]
func (h *PresenceHandler) UpdateMyPresence(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	var req UpdatePresenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, apierr.BadRequest("参数错误").WithDetails(err.Error()))
		return
	}

	// 验证状态值
	if req.Status != "" && !isValidPresenceStatus(req.Status) {
		resp.Error(c, apierr.BadRequest("无效的状态值"))
		return
	}

	presence, err := h.svc.UpdatePresence(c.Request.Context(), userID, &presenceservice.UpdatePresenceRequest{
		Status:          req.Status,
		CurrentGameID:   req.CurrentGameID,
		CurrentGameName: req.CurrentGameName,
		CustomStatus:    req.CustomStatus,
		DeviceType:      req.DeviceType,
	})
	if err != nil {
		resp.Error(c, apierr.InternalError("更新在线状态失败").WithDetails(err.Error()))
		return
	}

	resp.OK(c, presence)
}

// SetStatus 设置在线状态
// @Summary 设置在线状态
// @Tags 用户-在线状态
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body SetStatusRequest true "状态设置请求"
// @Success 200 {object} model.SuccessResponse
// @Router /user/presence/status [put]
func (h *PresenceHandler) SetStatus(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	var req SetStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, apierr.BadRequest("参数错误").WithDetails(err.Error()))
		return
	}

	if !isValidPresenceStatus(req.Status) {
		resp.Error(c, apierr.BadRequest("无效的状态值"))
		return
	}

	if err := h.svc.SetStatus(c.Request.Context(), userID, req.Status); err != nil {
		resp.Error(c, apierr.InternalError("设置状态失败").WithDetails(err.Error()))
		return
	}

	resp.OK(c, gin.H{"message": "状态已更新"})
}

// Heartbeat 心跳保活
// @Summary 心跳保活
// @Tags 用户-在线状态
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.SuccessResponse
// @Router /user/presence/heartbeat [post]
func (h *PresenceHandler) Heartbeat(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	if err := h.svc.Heartbeat(c.Request.Context(), userID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			// 如果没有状态记录，创建一个
			_, err = h.svc.GoOnline(c.Request.Context(), userID, "")
			if err != nil {
				resp.Error(c, apierr.InternalError("心跳失败").WithDetails(err.Error()))
				return
			}
		} else {
			resp.Error(c, apierr.InternalError("心跳失败").WithDetails(err.Error()))
			return
		}
	}

	resp.OK(c, gin.H{"message": "心跳成功"})
}

// GetPlayerPresence 获取指定陪玩师的在线状态
// @Summary 获取指定陪玩师的在线状态
// @Tags 用户-在线状态
// @Produce json
// @Param id path int true "陪玩师ID"
// @Success 200 {object} model.PlayerPresence
// @Router /user/players/{id}/presence [get]
func (h *PresenceHandler) GetPlayerPresence(c *gin.Context) {
	playerID, ok := resp.ParseIDOrFail(c, "id")
	if !ok {
		return
	}

	presence, err := h.svc.GetPresence(c.Request.Context(), playerID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			// 返回默认离线状态
			resp.OK(c, gin.H{
				"playerId": playerID,
				"status":   model.PresenceOffline,
			})
			return
		}
		resp.Error(c, apierr.InternalError("获取在线状态失败").WithDetails(err.Error()))
		return
	}

	// 如果是隐身状态，对外显示为离线
	if presence.Status == model.PresenceInvisible {
		presence.Status = model.PresenceOffline
	}

	resp.OK(c, presence)
}

// GetPlayersPresence 批量获取陪玩师在线状态
// @Summary 批量获取陪玩师在线状态
// @Tags 用户-在线状态
// @Accept json
// @Produce json
// @Param body body GetPlayersPresenceRequest true "陪玩师ID列表"
// @Success 200 {array} model.PlayerPresence
// @Router /user/players/presence [post]
func (h *PresenceHandler) GetPlayersPresence(c *gin.Context) {
	var req GetPlayersPresenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, apierr.BadRequest("参数错误").WithDetails(err.Error()))
		return
	}

	if len(req.PlayerIDs) == 0 {
		resp.OK(c, []model.PlayerPresence{})
		return
	}

	if len(req.PlayerIDs) > 100 {
		resp.Error(c, apierr.BadRequest("最多查询100个陪玩师"))
		return
	}

	presences, err := h.svc.GetPresencesByPlayerIDs(c.Request.Context(), req.PlayerIDs)
	if err != nil {
		resp.Error(c, apierr.InternalError("获取在线状态失败").WithDetails(err.Error()))
		return
	}

	// 隐身状态对外显示为离线
	for i := range presences {
		if presences[i].Status == model.PresenceInvisible {
			presences[i].Status = model.PresenceOffline
		}
	}

	resp.OK(c, presences)
}

// ListOnlinePlayers 获取在线陪玩师列表
// @Summary 获取在线陪玩师列表
// @Tags 用户-在线状态
// @Produce json
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param gameId query int false "游戏ID"
// @Success 200 {object} resp.PagedResponse
// @Router /user/players/online [get]
func (h *PresenceHandler) ListOnlinePlayers(c *gin.Context) {
	page := resp.ParseQueryInt(c, "page", 1)
	pageSize := resp.ParseQueryInt(c, "pageSize", 20)
	gameID := resp.ParseQueryUint64Ptr(c, "gameId")

	presences, total, err := h.svc.ListOnline(c.Request.Context(), repository.PlayerPresenceListOptions{
		Page:     page,
		PageSize: pageSize,
		GameID:   gameID,
	})
	if err != nil {
		resp.Error(c, apierr.InternalError("获取在线列表失败").WithDetails(err.Error()))
		return
	}

	resp.List(c, presences, &model.Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
	})
}

// GetOnlineCount 获取在线陪玩师数量
// @Summary 获取在线陪玩师数量
// @Tags 用户-在线状态
// @Produce json
// @Success 200 {object} model.SuccessResponse
// @Router /user/players/online/count [get]
func (h *PresenceHandler) GetOnlineCount(c *gin.Context) {
	count, err := h.svc.CountOnline(c.Request.Context())
	if err != nil {
		resp.Error(c, apierr.InternalError("获取在线数量失败").WithDetails(err.Error()))
		return
	}

	resp.OK(c, gin.H{"count": count})
}

// ============================================================================
// Request/Response DTOs
// ============================================================================

// UpdatePresenceRequest 更新状态请求
type UpdatePresenceRequest struct {
	Status          model.PlayerPresenceStatus `json:"status,omitempty"`
	CurrentGameID   *uint64                    `json:"currentGameId,omitempty"`
	CurrentGameName string                     `json:"currentGameName,omitempty"`
	CustomStatus    string                     `json:"customStatus,omitempty"`
	DeviceType      string                     `json:"deviceType,omitempty"`
}

// SetStatusRequest 设置状态请求
type SetStatusRequest struct {
	Status model.PlayerPresenceStatus `json:"status" binding:"required"`
}

// GetPlayersPresenceRequest 批量获取状态请求
type GetPlayersPresenceRequest struct {
	PlayerIDs []uint64 `json:"playerIds" binding:"required"`
}

// ============================================================================
// Helper functions
// ============================================================================

// isValidPresenceStatus 验证状态值是否有效
func isValidPresenceStatus(status model.PlayerPresenceStatus) bool {
	switch status {
	case model.PresenceOnline,
		model.PresenceAccepting,
		model.PresenceInGame,
		model.PresenceMatching,
		model.PresenceResting,
		model.PresenceOffline,
		model.PresenceInvisible:
		return true
	default:
		return false
	}
}

// RegisterPresenceRoutes 注册在线状态路由
func RegisterPresenceRoutes(rg *gin.RouterGroup, svc *presenceservice.Service, authMiddleware gin.HandlerFunc) {
	h := NewPresenceHandler(svc)

	// 需要认证的路由
	presenceGroup := rg.Group("/presence")
	presenceGroup.Use(authMiddleware)
	{
		presenceGroup.GET("", h.GetMyPresence)
		presenceGroup.PUT("", h.UpdateMyPresence)
		presenceGroup.PUT("/status", h.SetStatus)
		presenceGroup.POST("/heartbeat", h.Heartbeat)
	}

	// 玩家状态查询（部分需要认证）
	playersGroup := rg.Group("/players")
	{
		playersGroup.GET("/:id/presence", h.GetPlayerPresence)
		playersGroup.POST("/presence", h.GetPlayersPresence)
		playersGroup.GET("/online", h.ListOnlinePlayers)
		playersGroup.GET("/online/count", h.GetOnlineCount)
	}
}
