package user

import (
	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	trtcservice "gamelink/internal/service/trtc"
	"gamelink/pkg/apierr"
)

// VoiceHandler 语音处理器
type VoiceHandler struct {
	svc *trtcservice.Service
}

// NewVoiceHandler 创建语音处理器
func NewVoiceHandler(svc *trtcservice.Service) *VoiceHandler {
	return &VoiceHandler{svc: svc}
}

// GetUserSig 获取TRTC UserSig
// @Summary 获取TRTC UserSig
// @Tags 用户-语音
// @Produce json
// @Security BearerAuth
// @Param id path int true "房间ID"
// @Success 200 {object} trtcservice.UserSigResponse
// @Router /user/rooms/{id}/voice/token [get]
func (h *VoiceHandler) GetUserSig(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	roomID, ok := resp.ParseIDOrFail(c, "id")
	if !ok {
		return
	}

	result, err := h.svc.GetUserSig(c.Request.Context(), roomID, userID)
	if err != nil {
		resp.Error(c, err)
		return
	}

	resp.OK(c, result)
}

// StartVoice 开启语音
// @Summary 开启房间语音
// @Tags 用户-语音
// @Produce json
// @Security BearerAuth
// @Param id path int true "房间ID"
// @Success 200 {object} model.SuccessResponse
// @Router /user/rooms/{id}/voice/start [post]
func (h *VoiceHandler) StartVoice(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	roomID, ok := resp.ParseIDOrFail(c, "id")
	if !ok {
		return
	}

	if err := h.svc.StartVoice(c.Request.Context(), roomID, userID); err != nil {
		resp.Error(c, err)
		return
	}

	resp.OK(c, gin.H{"message": "语音已开启"})
}

// StopVoice 关闭语音
// @Summary 关闭房间语音
// @Tags 用户-语音
// @Produce json
// @Security BearerAuth
// @Param id path int true "房间ID"
// @Success 200 {object} model.SuccessResponse
// @Router /user/rooms/{id}/voice/stop [post]
func (h *VoiceHandler) StopVoice(c *gin.Context) {
	userID, ok := resp.GetUserIDOrFail(c)
	if !ok {
		return
	}

	roomID, ok := resp.ParseIDOrFail(c, "id")
	if !ok {
		return
	}

	if err := h.svc.StopVoice(c.Request.Context(), roomID, userID); err != nil {
		resp.Error(c, err)
		return
	}

	resp.OK(c, gin.H{"message": "语音已关闭"})
}

// GetVoiceStatus 获取语音状态
// @Summary 获取房间语音状态
// @Tags 用户-语音
// @Produce json
// @Param id path int true "房间ID"
// @Success 200 {object} trtcservice.VoiceStatusResponse
// @Router /user/rooms/{id}/voice/status [get]
func (h *VoiceHandler) GetVoiceStatus(c *gin.Context) {
	roomID, ok := resp.ParseIDOrFail(c, "id")
	if !ok {
		return
	}

	status, err := h.svc.GetVoiceStatus(c.Request.Context(), roomID)
	if err != nil {
		resp.Error(c, err)
		return
	}

	resp.OK(c, status)
}

// ============================================================================
// Route Registration
// ============================================================================

// RegisterVoiceRoutes 注册语音路由
func RegisterVoiceRoutes(rg *gin.RouterGroup, svc *trtcservice.Service, authMiddleware gin.HandlerFunc) {
	if svc == nil {
		// TRTC服务未配置，跳过注册
		return
	}

	h := NewVoiceHandler(svc)

	// 语音路由挂载在房间下
	roomsGroup := rg.Group("/rooms")
	{
		// 公开路由
		roomsGroup.GET("/:id/voice/status", h.GetVoiceStatus)

		// 需要认证的路由
		authGroup := roomsGroup.Group("")
		authGroup.Use(authMiddleware)
		{
			authGroup.GET("/:id/voice/token", h.GetUserSig)
			authGroup.POST("/:id/voice/start", h.StartVoice)
			authGroup.POST("/:id/voice/stop", h.StopVoice)
		}
	}
}

// VoiceNotConfiguredHandler 语音未配置时的占位处理器
func VoiceNotConfiguredHandler(c *gin.Context) {
	resp.Error(c, apierr.ServiceUnavailable("语音服务未配置"))
}
