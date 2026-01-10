package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/middleware"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/service/gamerank"
	"gamelink/pkg/apierr"
)

// GameRankHandler 游戏段位配置管理接口
type GameRankHandler struct {
	svc *gamerank.GameRankService
}

// NewGameRankHandler 创建Handler
func NewGameRankHandler(svc *gamerank.GameRankService) *GameRankHandler {
	return &GameRankHandler{svc: svc}
}

// GameRankPayload 创建/更新段位请求体
type GameRankPayload struct {
	GameID      uint64 `json:"gameId" binding:"required"`
	Name        string `json:"name" binding:"required,max=64"`
	Level       int    `json:"level"`
	PriceCents  int64  `json:"priceCents"`
	IconURL     string `json:"iconUrl"`
	Color       string `json:"color"`
	Description string `json:"description"`
	SortOrder   int    `json:"sortOrder"`
	IsActive    bool   `json:"isActive"`
}

// GameRankUpdatePayload 更新段位请求体（不含GameID）
type GameRankUpdatePayload struct {
	Name        string `json:"name" binding:"required,max=64"`
	Level       int    `json:"level"`
	PriceCents  int64  `json:"priceCents"`
	IconURL     string `json:"iconUrl"`
	Color       string `json:"color"`
	Description string `json:"description"`
	SortOrder   int    `json:"sortOrder"`
	IsActive    bool   `json:"isActive"`
}

// ListGameRanks
// @Summary      列出游戏段位
// @Tags         Admin/GameRanks
// @Security     BearerAuth
// @Param        page       query  int     false  "页码"
// @Param        pageSize   query  int     false  "每页数量"
// @Param        gameId     query  int     false  "游戏ID筛选"
// @Param        keyword    query  string  false  "关键词搜索"
// @Param        isActive   query  bool    false  "是否启用"
// @Produce      json
// @Success      200  {object}  model.APIResponse[[]model.GameRank]
// @Router       /admin/game-ranks [get]
func (h *GameRankHandler) ListGameRanks(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	gameID, ok := QueryUint64PtrAndRespond(c, "gameId", apierr.ErrInvalidGameID)
	if !ok {
		return
	}

	var isActive *bool
	if v := c.Query("isActive"); v != "" {
		b := v == "true"
		isActive = &b
	}

	opts := repository.GameRankListOptions{
		Page:     page,
		PageSize: pageSize,
		GameID:   gameID,
		Keyword:  c.Query("keyword"),
		IsActive: isActive,
	}

	ranks, pagination, err := h.svc.ListPaged(c.Request.Context(), opts)
	if err != nil {
		respondError(c, err)
		return
	}
	respondList(c, ranks, pagination)
}

// ListGameRanksByGame
// @Summary      根据游戏ID获取段位列表
// @Tags         Admin/GameRanks
// @Security     BearerAuth
// @Param        gameId   path  int  true  "游戏ID"
// @Produce      json
// @Success      200  {object}  model.APIResponse[[]model.GameRank]
// @Router       /admin/games/{gameId}/ranks [get]
func (h *GameRankHandler) ListGameRanksByGame(c *gin.Context) {
	gameID, ok := ParseIDAndRespond(c, "gameId")
	if !ok {
		return
	}

	ranks, err := h.svc.ListByGameID(c.Request.Context(), gameID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, ranks)
}

// GetGameRank
// @Summary      获取游戏段位详情
// @Tags         Admin/GameRanks
// @Security     BearerAuth
// @Param        id   path  int  true  "段位ID"
// @Produce      json
// @Success      200  {object}  model.APIResponse[model.GameRank]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/game-ranks/{id} [get]
func (h *GameRankHandler) GetGameRank(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	rank, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, rank)
}

// CreateGameRank
// @Summary      创建游戏段位
// @Tags         Admin/GameRanks
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  GameRankPayload  true  "段位信息"
// @Success      201  {object}  model.APIResponse[model.GameRank]
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/game-ranks [post]
func (h *GameRankHandler) CreateGameRank(c *gin.Context) {
	var payload GameRankPayload
	if !ValidateAndRespond(c, &payload) {
		return
	}

	rank, err := h.svc.Create(c.Request.Context(), gamerank.CreateInput{
		GameID:      payload.GameID,
		Name:        payload.Name,
		Level:       payload.Level,
		PriceCents:  payload.PriceCents,
		IconURL:     payload.IconURL,
		Color:       payload.Color,
		Description: payload.Description,
		SortOrder:   payload.SortOrder,
		IsActive:    payload.IsActive,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	respondCreated(c, rank)
}

// UpdateGameRank
// @Summary      更新游戏段位
// @Tags         Admin/GameRanks
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                    true  "段位ID"
// @Param        request  body  GameRankUpdatePayload  true  "段位信息"
// @Success      200  {object}  model.APIResponse[model.GameRank]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/game-ranks/{id} [put]
func (h *GameRankHandler) UpdateGameRank(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var payload GameRankUpdatePayload
	if !ValidateAndRespond(c, &payload) {
		return
	}

	rank, err := h.svc.Update(c.Request.Context(), id, gamerank.UpdateInput{
		Name:        payload.Name,
		Level:       payload.Level,
		PriceCents:  payload.PriceCents,
		IconURL:     payload.IconURL,
		Color:       payload.Color,
		Description: payload.Description,
		SortOrder:   payload.SortOrder,
		IsActive:    payload.IsActive,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	respondUpdated(c, rank)
}

// DeleteGameRank
// @Summary      删除游戏段位
// @Tags         Admin/GameRanks
// @Security     BearerAuth
// @Param        id   path  int  true  "段位ID"
// @Produce      json
// @Success      200  {object}  model.SuccessResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/game-ranks/{id} [delete]
func (h *GameRankHandler) DeleteGameRank(c *gin.Context) {
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

// BatchDeleteGameRanksRequest 批量删除请求
type BatchDeleteGameRanksRequest struct {
	IDs []string `json:"ids" binding:"required,min=1"`
}

// BatchDeleteGameRanks
// @Summary      批量删除游戏段位
// @Tags         Admin/GameRanks
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BatchDeleteGameRanksRequest  true  "段位ID列表"
// @Success      200  {object}  model.APIResponse[map[string]int64]
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/game-ranks/batch/delete [post]
func (h *GameRankHandler) BatchDeleteGameRanks(c *gin.Context) {
	var req BatchDeleteGameRanksRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	ids := make([]uint64, 0, len(req.IDs))
	for _, idStr := range req.IDs {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			respondBadRequest(c, "invalid id: "+idStr)
			return
		}
		ids = append(ids, id)
	}

	deleted, err := h.svc.BatchDelete(c.Request.Context(), ids)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, map[string]int64{"deleted": deleted})
}

// BatchUpdateStatusRequest 批量更新状态请求
type BatchUpdateStatusRequest struct {
	IDs      []string `json:"ids" binding:"required,min=1"`
	IsActive bool     `json:"isActive"`
}

// BatchUpdateGameRankStatus
// @Summary      批量更新游戏段位状态
// @Tags         Admin/GameRanks
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BatchUpdateStatusRequest  true  "请求体"
// @Success      200  {object}  model.APIResponse[map[string]int64]
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/game-ranks/batch/status [post]
func (h *GameRankHandler) BatchUpdateGameRankStatus(c *gin.Context) {
	var req BatchUpdateStatusRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	ids := make([]uint64, 0, len(req.IDs))
	for _, idStr := range req.IDs {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			respondBadRequest(c, "invalid id: "+idStr)
			return
		}
		ids = append(ids, id)
	}

	updated, err := h.svc.BatchUpdateStatus(c.Request.Context(), ids, req.IsActive)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, map[string]int64{"updated": updated})
}

// RegisterGameRankRoutes 注册游戏段位管理路由
func RegisterGameRankRoutes(router gin.IRouter, svc *gamerank.GameRankService, pm *middleware.PermissionMiddleware) {
	h := NewGameRankHandler(svc)

	group := router.Group("/game-ranks")
	group.Use(pm.RequireAuth())
	{
		group.GET("", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/game-ranks"), h.ListGameRanks)
		group.POST("", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/game-ranks"), h.CreateGameRank)
		group.GET("/:id", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/game-ranks/:id"), h.GetGameRank)
		group.PUT("/:id", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/game-ranks/:id"), h.UpdateGameRank)
		group.DELETE("/:id", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/game-ranks/:id"), h.DeleteGameRank)
		group.GET("/game/:gameId", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/game-ranks/game/:gameId"), h.ListGameRanksByGame)
		group.POST("/batch/delete", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/game-ranks/batch/delete"), h.BatchDeleteGameRanks)
		group.POST("/batch/status", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/game-ranks/batch/status"), h.BatchUpdateGameRankStatus)
	}
}
