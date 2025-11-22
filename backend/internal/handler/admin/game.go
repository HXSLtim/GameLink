package admin

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler"
	apierr "gamelink/internal/apierr"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	adminservice "gamelink/internal/service/admin"
)

// GameHandler 处理后台游戏管理接口�?
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
// @Param        page       query  int  false  "页码"
// @Param        pageSize   query     int       false  "每页数量"
// @Produce      json
// @Success      200  {object}  model.SuccessResponse
// @Router       /admin/games [get]
//
// ListGames 返回全部游戏�?
func (h *GameHandler) ListGames(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	games, pagination, err := h.svc.ListGamesPaged(c.Request.Context(), page, pageSize)
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	games = ensureSlice(games)
	writeJSON(c, http.StatusOK, model.APIResponse[[]model.Game]{
		Success:    true,
		Code:       http.StatusOK,
		Message:    "OK",
		Data:       games,
		Pagination: pagination,
	})
}

// GetGame
// @Summary      获取游戏
// @Tags         Admin/Games
// @Security     BearerAuth
// @Param        id   path  int  true  "游戏ID"
// @Produce      json
// @Success      200  {object}  model.SuccessResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/games/{id} [get]
//
// GetGame 获取单个游戏�?
func (h *GameHandler) GetGame(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}
	
	game, err := h.svc.GetGame(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			handler.RespondWithServiceError(c, apierr.NotFound(apierr.ErrGameNotFound))
			return
		}
		handler.RespondWithServiceError(c, err)
		return
	}
	
	handler.RespondSuccess(c, "OK", game)
}

// CreateGame
// @Summary      创建游戏
// @Tags         Admin/Games
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  GamePayload  true  "游戏信息"
// @Success      201  {object}  model.SuccessResponse
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/games [post]
//
// CreateGame 创建新游戏�?
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
		if errors.Is(err, adminservice.ErrValidation) {
			handler.RespondWithServiceError(c, apierr.BadRequest("validation failed"))
			return
		}
		handler.RespondWithServiceError(c, err)
		return
	}

	handler.RespondCreated(c, game)
}

// UpdateGame
// @Summary      更新游戏
// @Tags         Admin/Games
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int         true  "游戏ID"
// @Param        request  body  GamePayload true  "游戏信息"
// @Success      200  {object}  model.SuccessResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/games/{id} [put]
//
// UpdateGame 更新游戏信息�?
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
		if errors.Is(err, adminservice.ErrValidation) {
			handler.RespondWithServiceError(c, apierr.BadRequest("validation failed"))
			return
		}
		if errors.Is(err, repository.ErrNotFound) {
			handler.RespondWithServiceError(c, apierr.NotFound("game not found"))
			return
		}
		handler.RespondWithServiceError(c, err)
		return
	}

	handler.RespondSuccess(c, "updated", game)
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
//
// DeleteGame 删除游戏�?
func (h *GameHandler) DeleteGame(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}

	err = h.svc.DeleteGame(c.Request.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		_ = c.Error(adminservice.ErrNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "deleted",
	})
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
// @Success      200  {object}  model.SuccessResponse
// @Router       /admin/games/{id}/logs [get]
func (h *GameHandler) ListGameLogs(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, 400, apierr.ErrInvalidID)
		return
	}
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}
	var actorID *uint64
	if v, err := queryUint64Ptr(c, "actor_user_id"); err == nil {
		actorID = v
	}
	var dateFrom, dateTo *time.Time
	if v, err := queryTimePtr(c, "date_from"); err == nil {
		dateFrom = v
	} else {
		writeJSONError(c, 400, apierr.ErrInvalidDateFrom)
		return
	}
	if v, err := queryTimePtr(c, "date_to"); err == nil {
		dateTo = v
	} else {
		writeJSONError(c, 400, apierr.ErrInvalidDateTo)
		return
	}
	opts := repository.OperationLogListOptions{Page: page, PageSize: pageSize, Action: strings.TrimSpace(c.Query("action")), ActorUserID: actorID, DateFrom: dateFrom, DateTo: dateTo}
	items, p, err := h.svc.ListOperationLogs(c.Request.Context(), "game", id, opts)
	if err != nil {
		writeJSONError(c, 500, err.Error())
		return
	}
	if strings.EqualFold(strings.TrimSpace(c.Query("export")), "csv") {
		exportOperationLogsCSV(c, "game", id, items)
		return
	}
	items = ensureSlice(items)
	writeJSON(c, 200, model.APIResponse[[]model.OperationLog]{Success: true, Code: 200, Message: "OK", Data: items, Pagination: p})
}

// GamePayload defines request body for creating/updating a game.
type GamePayload struct {
	Key         string `json:"key" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Category    string `json:"category"`
	IconURL     string `json:"icon_url"`
	Description string `json:"description"`
}
