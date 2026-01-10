package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	adminservice "gamelink/internal/service/admin"
)

// Game 游戏模型（类型别名）
type Game = model.Game

// GameHandler 处理后台游戏管理接口
type GameHandler struct {
	svc *adminservice.AdminService
}

// NewGameHandler 创建Handler
func NewGameHandler(svc *adminservice.AdminService) *GameHandler {
	return &GameHandler{svc: svc}
}

// ListGames
// @Summary      列出游戏
// @Tags         Admin/Games
// @Security     BearerAuth
// @Param        page       query  int     false  "页码"
// @Param        pageSize   query  int     false  "每页数量"
// @Param        keyword    query  string  false  "关键词搜索"
// @Produce      json
// @Success      200  {object}  model.APIResponse[[]Game]
// @Router       /admin/games [get]
func (h *GameHandler) ListGames(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	keyword := c.Query("keyword")
	games, pagination, err := h.svc.ListGamesPagedWithFilter(c.Request.Context(), page, pageSize, keyword)
	if err != nil {
		respondError(c, err)
		return
	}
	respondList(c, games, pagination)
}

// GetGame
// @Summary      获取游戏
// @Tags         Admin/Games
// @Security     BearerAuth
// @Param        id   path  int  true  "游戏ID"
// @Produce      json
// @Success      200  {object}  model.APIResponse[Game]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/games/{id} [get]
func (h *GameHandler) GetGame(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	game, err := h.svc.GetGame(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, game)
}

// CreateGame
// @Summary      创建游戏
// @Tags         Admin/Games
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  GamePayload  true  "游戏信息"
// @Success      201  {object}  model.APIResponse[Game]
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/games [post]
func (h *GameHandler) CreateGame(c *gin.Context) {
	var payload GamePayload
	if !ValidateAndRespond(c, &payload) {
		return
	}

	game, err := h.svc.CreateGame(c.Request.Context(), adminservice.CreateGameInput{
		Key:         payload.Key,
		Name:        payload.Name,
		Category:    payload.Category,
		IconURL:     payload.IconURL,
		Description: payload.Description,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	respondCreated(c, game)
}

// UpdateGame
// @Summary      更新游戏
// @Tags         Admin/Games
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int         true  "游戏ID"
// @Param        request  body  GamePayload true  "游戏信息"
// @Success      200  {object}  model.APIResponse[Game]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/games/{id} [put]
func (h *GameHandler) UpdateGame(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var payload GamePayload
	if !ValidateAndRespond(c, &payload) {
		return
	}

	game, err := h.svc.UpdateGame(c.Request.Context(), id, adminservice.UpdateGameInput{
		Key:         payload.Key,
		Name:        payload.Name,
		Category:    payload.Category,
		IconURL:     payload.IconURL,
		Description: payload.Description,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	respondUpdated(c, game)
}

// DeleteGame
// @Summary      删除游戏
// @Tags         Admin/Games
// @Security     BearerAuth
// @Param        id   path  int  true  "游戏ID"
// @Produce      json
// @Success      200  {object}  model.SuccessResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/games/{id} [delete]
func (h *GameHandler) DeleteGame(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	err := h.svc.DeleteGame(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondDeleted(c)
}

// BatchDeleteGamesRequest 批量删除游戏请求
type BatchDeleteGamesRequest struct {
	GameIDs []string `json:"gameIds" binding:"required,min=1"`
}

// BatchDeleteGames
// @Summary      批量删除游戏
// @Tags         Admin/Games
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BatchDeleteGamesRequest  true  "游戏ID列表"
// @Success      200  {object}  model.APIResponse[map[string]int64]
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/games/batch/delete [post]
func (h *GameHandler) BatchDeleteGames(c *gin.Context) {
	var req BatchDeleteGamesRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	// 转换 string[] 到 uint64[]
	ids := make([]uint64, 0, len(req.GameIDs))
	for _, idStr := range req.GameIDs {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			respondBadRequest(c, "invalid game id: "+idStr)
			return
		}
		ids = append(ids, id)
	}

	deleted, err := h.svc.BatchDeleteGames(c.Request.Context(), ids)
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, map[string]int64{"deleted": deleted})
}

// ListGameLogs
// @Summary      获取游戏操作日志
// @Tags         Admin/Games
// @Security     BearerAuth
// @Produce      json
// @Param        id           path   int  true  "游戏ID"
// @Param        page         query  int  false "页码"
// @Param        pageSize   query     int       false  "每页数量"
// @Param        action       query  string false "动作过滤" Enums(create,update,delete)
// @Param        actor_user_id query int   false "操作者用户ID"
// @Param        dateFrom   query     string    false  "Start date (YYYY-MM-DD)"
// @Param        dateTo     query     string    false  "End date (YYYY-MM-DD)"
// @Param        export       query  string false "导出格式" Enums(csv)
// @Param        fields         query    string       false  "Export fields (comma separated)"// @Param        header_lang  query  string false "列头语言" Enums(en,zh)
// @Success      200  {object}  model.APIResponse[[]model.OperationLog]
// @Router       /admin/games/{id}/logs [get]
func (h *GameHandler) ListGameLogs(c *gin.Context) {
	handleOperationLogList(c, "game", h.svc.ListOperationLogs)
}

// GamePayload defines request body for creating/updating a game.
type GamePayload struct {
	Key         string `json:"key" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Category    string `json:"category"`
	IconURL     string `json:"icon_url"`
	Description string `json:"description"`
}
