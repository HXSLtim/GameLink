package admin

import (
	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/middleware"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/service/playercertification"
	"gamelink/pkg/apierr"
)

// PlayerCertificationHandler 陪玩师实名认证管理接口
type PlayerCertificationHandler struct {
	svc *playercertification.PlayerCertificationService
}

// NewPlayerCertificationHandler 创建Handler
func NewPlayerCertificationHandler(svc *playercertification.PlayerCertificationService) *PlayerCertificationHandler {
	return &PlayerCertificationHandler{svc: svc}
}

// VerifyPlayerCertificationRequest 审核实名认证请求
type VerifyPlayerCertificationRequest struct {
	Status       string `json:"status" binding:"required,oneof=verified rejected"`
	RejectReason string `json:"rejectReason"`
}

// ListPlayerCertifications
// @Summary      列出实名认证记录
// @Tags         Admin/PlayerCertifications
// @Security     BearerAuth
// @Param        page       query  int     false  "页码"
// @Param        pageSize   query  int     false  "每页数量"
// @Param        playerId   query  int     false  "陪玩师ID筛选"
// @Param        status     query  string  false  "状态筛选" Enums(pending,verified,rejected)
// @Param        keyword    query  string  false  "关键词搜索（真实姓名）"
// @Produce      json
// @Success      200  {array}   model.PlayerCertification
// @Router       /admin/player-certifications [get]
func (h *PlayerCertificationHandler) ListPlayerCertifications(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	playerID, ok := QueryUint64PtrAndRespond(c, "playerId", apierr.ErrInvalidPlayerID)
	if !ok {
		return
	}

	var status *model.CertificationStatus
	if v := c.Query("status"); v != "" {
		s := model.CertificationStatus(v)
		status = &s
	}

	opts := repository.PlayerCertificationListOptions{
		Page:     page,
		PageSize: pageSize,
		PlayerID: playerID,
		Status:   status,
		Keyword:  c.Query("keyword"),
	}

	certs, pagination, err := h.svc.ListPaged(c.Request.Context(), opts)
	if err != nil {
		respondError(c, err)
		return
	}
	respondList(c, certs, pagination)
}

// ListPendingPlayerCertifications
// @Summary      列出待审核的实名认证
// @Tags         Admin/PlayerCertifications
// @Security     BearerAuth
// @Param        page       query  int  false  "页码"
// @Param        pageSize   query  int  false  "每页数量"
// @Produce      json
// @Success      200  {array}   model.PlayerCertification
// @Router       /admin/player-certifications/pending [get]
func (h *PlayerCertificationHandler) ListPendingPlayerCertifications(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	certs, pagination, err := h.svc.ListPending(c.Request.Context(), page, pageSize)
	if err != nil {
		respondError(c, err)
		return
	}
	respondList(c, certs, pagination)
}

// GetPlayerCertification
// @Summary      获取实名认证详情
// @Tags         Admin/PlayerCertifications
// @Security     BearerAuth
// @Param        id   path  int  true  "认证记录ID"
// @Produce      json
// @Success      200  {object}  model.PlayerCertification
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/player-certifications/{id} [get]
func (h *PlayerCertificationHandler) GetPlayerCertification(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	cert, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, cert)
}

// VerifyPlayerCertification
// @Summary      审核实名认证
// @Tags         Admin/PlayerCertifications
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                               true  "认证记录ID"
// @Param        request  body  VerifyPlayerCertificationRequest  true  "审核信息"
// @Success      200  {object}  model.PlayerCertification
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/player-certifications/{id}/verify [post]
func (h *PlayerCertificationHandler) VerifyPlayerCertification(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var req VerifyPlayerCertificationRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	// 获取当前管理员ID
	adminID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	cert, err := h.svc.Verify(c.Request.Context(), playercertification.VerifyInput{
		CertID:       id,
		Status:       model.CertificationStatus(req.Status),
		VerifiedBy:   adminID,
		RejectReason: req.RejectReason,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	respondUpdated(c, cert)
}

// DeletePlayerCertification
// @Summary      删除实名认证记录
// @Tags         Admin/PlayerCertifications
// @Security     BearerAuth
// @Param        id   path  int  true  "认证记录ID"
// @Produce      json
// @Success      200  {object}  model.SuccessResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/player-certifications/{id} [delete]
func (h *PlayerCertificationHandler) DeletePlayerCertification(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	respondDeleted(c)
}

// GetPlayerCertificationStats
// @Summary      获取实名认证统计
// @Tags         Admin/PlayerCertifications
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]int64
// @Router       /admin/player-certifications/stats [get]
func (h *PlayerCertificationHandler) GetPlayerCertificationStats(c *gin.Context) {
	stats, err := h.svc.GetStats(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}

	// 转换为 string key 的 map
	result := make(map[string]int64)
	for k, v := range stats {
		result[string(k)] = v
	}
	respondSuccess(c, result)
}

// GetPlayerCertificationPendingCount
// @Summary      获取待审核实名认证数量
// @Tags         Admin/PlayerCertifications
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]int64
// @Router       /admin/player-certifications/pending/count [get]
func (h *PlayerCertificationHandler) GetPlayerCertificationPendingCount(c *gin.Context) {
	count, err := h.svc.GetPendingCount(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, map[string]int64{"count": count})
}

// GetPlayerCertificationByPlayer
// @Summary      获取陪玩师的实名认证
// @Tags         Admin/PlayerCertifications
// @Security     BearerAuth
// @Param        playerId   path  int  true  "陪玩师ID"
// @Produce      json
// @Success      200  {object}  model.PlayerCertification
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/players/{playerId}/certification [get]
func (h *PlayerCertificationHandler) GetPlayerCertificationByPlayer(c *gin.Context) {
	playerID, ok := ParseIDAndRespond(c, "playerId")
	if !ok {
		return
	}

	cert, err := h.svc.GetByPlayerID(c.Request.Context(), playerID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, cert)
}

// RegisterPlayerCertificationRoutes 注册陪玩师实名认证管理路由
func RegisterPlayerCertificationRoutes(router gin.IRouter, svc *playercertification.PlayerCertificationService, pm *middleware.PermissionMiddleware) {
	h := NewPlayerCertificationHandler(svc)

	group := router.Group("/player-certifications")
	group.Use(pm.RequireAuth())
	{
		group.GET("", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/player-certifications"), h.ListPlayerCertifications)
		group.GET("/pending", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/player-certifications/pending"), h.ListPendingPlayerCertifications)
		group.GET("/pending/count", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/player-certifications/pending/count"), h.GetPlayerCertificationPendingCount)
		group.GET("/stats", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/player-certifications/stats"), h.GetPlayerCertificationStats)
		group.GET("/:id", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/player-certifications/:id"), h.GetPlayerCertification)
		group.POST("/:id/verify", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/player-certifications/:id/verify"), h.VerifyPlayerCertification)
		group.DELETE("/:id", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/player-certifications/:id"), h.DeletePlayerCertification)
		group.GET("/player/:playerId", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/player-certifications/player/:playerId"), h.GetPlayerCertificationByPlayer)
	}
}
