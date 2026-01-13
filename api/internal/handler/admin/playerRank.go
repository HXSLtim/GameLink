package admin

import (
	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/middleware"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/service/playerrank"
	"gamelink/pkg/apierr"
)

// PlayerRankHandler 陪玩师段位认证管理接口
type PlayerRankHandler struct {
	svc *playerrank.PlayerRankService
}

// NewPlayerRankHandler 创建Handler
func NewPlayerRankHandler(svc *playerrank.PlayerRankService) *PlayerRankHandler {
	return &PlayerRankHandler{svc: svc}
}

// VerifyPlayerRankRequest 审核段位认证请求
type VerifyPlayerRankRequest struct {
	Status       string `json:"status" binding:"required,oneof=verified rejected revoked"`
	RejectReason string `json:"rejectReason"`
}

// ListPlayerRanks
// @Summary      列出段位认证记录
// @Tags         Admin/PlayerRanks
// @Security     BearerAuth
// @Param        page       query  int     false  "页码"
// @Param        pageSize   query  int     false  "每页数量"
// @Param        playerId   query  int     false  "陪玩师ID筛选"
// @Param        gameId     query  int     false  "游戏ID筛选"
// @Param        status     query  string  false  "状态筛选" Enums(pending,verified,rejected,revoked,expired)
// @Produce      json
// @Success      200  {array}   model.PlayerRankRecord
// @Router       /admin/player-ranks [get]
func (h *PlayerRankHandler) ListPlayerRanks(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	playerID, ok := QueryUint64PtrAndRespond(c, "playerId", apierr.ErrInvalidPlayerID)
	if !ok {
		return
	}

	gameID, ok := QueryUint64PtrAndRespond(c, "gameId", apierr.ErrInvalidGameID)
	if !ok {
		return
	}

	var status *model.PlayerRankStatus
	if v := c.Query("status"); v != "" {
		s := model.PlayerRankStatus(v)
		status = &s
	}

	opts := repository.PlayerRankListOptions{
		Page:     page,
		PageSize: pageSize,
		PlayerID: playerID,
		GameID:   gameID,
		Status:   status,
	}

	records, pagination, err := h.svc.ListPaged(c.Request.Context(), opts)
	if err != nil {
		respondError(c, err)
		return
	}
	respondList(c, records, pagination)
}

// ListPendingPlayerRanks
// @Summary      列出待审核的段位认证
// @Tags         Admin/PlayerRanks
// @Security     BearerAuth
// @Param        page       query  int  false  "页码"
// @Param        pageSize   query  int  false  "每页数量"
// @Produce      json
// @Success      200  {array}   model.PlayerRankRecord
// @Router       /admin/player-ranks/pending [get]
func (h *PlayerRankHandler) ListPendingPlayerRanks(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	records, pagination, err := h.svc.ListPending(c.Request.Context(), page, pageSize)
	if err != nil {
		respondError(c, err)
		return
	}
	respondList(c, records, pagination)
}

// GetPlayerRank
// @Summary      获取段位认证详情
// @Tags         Admin/PlayerRanks
// @Security     BearerAuth
// @Param        id   path  int  true  "认证记录ID"
// @Produce      json
// @Success      200  {object}  model.PlayerRankRecord
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/player-ranks/{id} [get]
func (h *PlayerRankHandler) GetPlayerRank(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	record, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, record)
}

// VerifyPlayerRank
// @Summary      审核段位认证
// @Tags         Admin/PlayerRanks
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                      true  "认证记录ID"
// @Param        request  body  VerifyPlayerRankRequest  true  "审核信息"
// @Success      200  {object}  model.PlayerRankRecord
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/player-ranks/{id}/verify [post]
func (h *PlayerRankHandler) VerifyPlayerRank(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var req VerifyPlayerRankRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	// 获取当前管理员ID
	adminID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	record, err := h.svc.Verify(c.Request.Context(), playerrank.VerifyInput{
		RecordID:     id,
		Status:       model.PlayerRankStatus(req.Status),
		VerifiedBy:   adminID,
		RejectReason: req.RejectReason,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	respondUpdated(c, record)
}

// DeletePlayerRank
// @Summary      删除段位认证记录
// @Tags         Admin/PlayerRanks
// @Security     BearerAuth
// @Param        id   path  int  true  "认证记录ID"
// @Produce      json
// @Success      200  {object}  model.SuccessResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/player-ranks/{id} [delete]
func (h *PlayerRankHandler) DeletePlayerRank(c *gin.Context) {
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

// GetPlayerRankStats
// @Summary      获取段位认证统计
// @Tags         Admin/PlayerRanks
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]int64
// @Router       /admin/player-ranks/stats [get]
func (h *PlayerRankHandler) GetPlayerRankStats(c *gin.Context) {
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

// GetPlayerRankPendingCount
// @Summary      获取待审核段位认证数量
// @Tags         Admin/PlayerRanks
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]int64
// @Router       /admin/player-ranks/pending/count [get]
func (h *PlayerRankHandler) GetPlayerRankPendingCount(c *gin.Context) {
	count, err := h.svc.GetPendingCount(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, map[string]int64{"count": count})
}

// ListPlayerRanksByPlayer
// @Summary      获取陪玩师的所有段位认证
// @Tags         Admin/PlayerRanks
// @Security     BearerAuth
// @Param        playerId   path  int  true  "陪玩师ID"
// @Produce      json
// @Success      200  {array}   model.PlayerRankRecord
// @Router       /admin/players/{playerId}/ranks [get]
func (h *PlayerRankHandler) ListPlayerRanksByPlayer(c *gin.Context) {
	playerID, ok := ParseIDAndRespond(c, "playerId")
	if !ok {
		return
	}

	records, err := h.svc.ListByPlayerID(c.Request.Context(), playerID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, records)
}

// RegisterPlayerRankRoutes 注册陪玩师段位认证管理路由
func RegisterPlayerRankRoutes(router gin.IRouter, svc *playerrank.PlayerRankService, pm *middleware.PermissionMiddleware) {
	h := NewPlayerRankHandler(svc)

	group := router.Group("/player-ranks")
	group.Use(pm.RequireAuth())
	{
		group.GET("", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/player-ranks"), h.ListPlayerRanks)
		group.GET("/pending", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/player-ranks/pending"), h.ListPendingPlayerRanks)
		group.GET("/pending/count", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/player-ranks/pending/count"), h.GetPlayerRankPendingCount)
		group.GET("/stats", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/player-ranks/stats"), h.GetPlayerRankStats)
		group.GET("/:id", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/player-ranks/:id"), h.GetPlayerRank)
		group.POST("/:id/verify", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/player-ranks/:id/verify"), h.VerifyPlayerRank)
		group.DELETE("/:id", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/player-ranks/:id"), h.DeletePlayerRank)
		group.GET("/player/:playerId", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/player-ranks/player/:playerId"), h.ListPlayerRanksByPlayer)
	}
}
