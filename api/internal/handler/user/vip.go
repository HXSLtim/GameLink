package user

import (
	"errors"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	"gamelink/internal/repository"
	vipservice "gamelink/internal/service/vip"
	"gamelink/pkg/apierr"
)

// VipHandler 用户VIP处理器
type VipHandler struct {
	svc *vipservice.Service
}

// NewVipHandler 创建用户VIP处理器
func NewVipHandler(svc *vipservice.Service) *VipHandler {
	return &VipHandler{svc: svc}
}

// ListLevels 获取VIP等级列表（用户端）
// @Summary 获取VIP等级列表
// @Tags 用户-VIP
// @Produce json
// @Success 200 {object} resp.Response{data=[]model.VipLevel}
// @Router /user/vip/levels [get]
func (h *VipHandler) ListLevels(c *gin.Context) {
	levels, err := h.svc.ListActiveLevels(c.Request.Context())
	if err != nil {
		resp.Error(c, apierr.InternalError("获取VIP等级列表失败").WithDetails(err.Error()))
		return
	}

	resp.OK(c, levels)
}

// GetLevel 获取VIP等级详情
// @Summary 获取VIP等级详情
// @Tags 用户-VIP
// @Produce json
// @Param id path int true "等级ID"
// @Success 200 {object} resp.Response{data=model.VipLevel}
// @Router /user/vip/levels/{id} [get]
func (h *VipHandler) GetLevel(c *gin.Context) {
	id, ok := resp.ParseIDOrFail(c, "id")
	if !ok {
		return
	}

	level, err := h.svc.GetLevel(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			resp.Error(c, apierr.NotFound("VIP等级不存在"))
			return
		}
		resp.Error(c, apierr.InternalError("获取VIP等级失败").WithDetails(err.Error()))
		return
	}

	// 只返回启用的等级
	if !level.IsActive {
		resp.Error(c, apierr.NotFound("VIP等级不存在"))
		return
	}

	resp.OK(c, level)
}

// GetUnlockThreshold 获取VIP解锁门槛
// @Summary 获取VIP解锁门槛
// @Tags 用户-VIP
// @Produce json
// @Success 200 {object} resp.Response
// @Router /user/vip/threshold [get]
func (h *VipHandler) GetUnlockThreshold(c *gin.Context) {
	consumeThreshold, rechargeThreshold, err := h.svc.GetUnlockThreshold(c.Request.Context())
	if err != nil {
		resp.Error(c, apierr.InternalError("获取VIP解锁门槛失败").WithDetails(err.Error()))
		return
	}

	resp.OK(c, gin.H{
		"consumeThreshold":  consumeThreshold,
		"rechargeThreshold": rechargeThreshold,
	})
}

// RegisterVipRoutes 注册用户VIP路由
func RegisterVipRoutes(rg *gin.RouterGroup, svc *vipservice.Service, _ gin.HandlerFunc) {
	h := NewVipHandler(svc)

	vipGroup := rg.Group("/vip")
	{
		vipGroup.GET("/levels", h.ListLevels)
		vipGroup.GET("/levels/:id", h.GetLevel)
		vipGroup.GET("/threshold", h.GetUnlockThreshold)
	}
}
