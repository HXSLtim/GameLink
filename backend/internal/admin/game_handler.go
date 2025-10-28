package admin

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	apierr "gamelink/internal/handler"
	"gamelink/internal/model"
	"gamelink/internal/service"
)

// GameHandler 处理后台游戏管理接口。
type GameHandler struct {
	svc *service.AdminService
}

// NewGameHandler 创建 Handler。
func NewGameHandler(svc *service.AdminService) *GameHandler {
	return &GameHandler{svc: svc}
}

// ListGames
// @Summary      列出游戏
// @Tags         Admin/Games
// @Security     BearerAuth
// @Param        page       query  int  false  "页码"
// @Param        page_size  query  int  false  "每页数量"
// @Produce      json
// @Success      200  {object}  map[string]any
// @Router       /admin/games [get]
//
// ListGames 返回全部游戏。
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
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/games/{id} [get]
//
// GetGame 获取单个游戏。
func (h *GameHandler) GetGame(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	game, err := h.svc.GetGame(c.Request.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		writeJSONError(c, http.StatusNotFound, apierr.ErrGameNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(c, http.StatusOK, model.APIResponse[*model.Game]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    game,
	})
}

// CreateGame
// @Summary      创建游戏
// @Tags         Admin/Games
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  GamePayload  true  "游戏信息"
// @Success      201  {object}  map[string]any
// @Failure      400  {object}  map[string]any
// @Router       /admin/games [post]
//
// CreateGame 创建新游戏。
func (h *GameHandler) CreateGame(c *gin.Context) {
	var payload GamePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload)
		return
	}

	game, err := h.svc.CreateGame(c.Request.Context(), service.CreateGameInput{
		Key:         payload.Key,
		Name:        payload.Name,
		Category:    payload.Category,
		IconURL:     payload.IconURL,
		Description: payload.Description,
	})
    if errors.Is(err, service.ErrValidation) {
        _ = c.Error(service.ErrValidation)
        return
    }
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusCreated, model.APIResponse[*model.Game]{
		Success: true,
		Code:    http.StatusCreated,
		Message: "created",
		Data:    game,
	})
}

// UpdateGame
// @Summary      更新游戏
// @Tags         Admin/Games
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int         true  "游戏ID"
// @Param        request  body  GamePayload true  "游戏信息"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/games/{id} [put]
//
// UpdateGame 更新游戏信息。
func (h *GameHandler) UpdateGame(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}

	var payload GamePayload
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload)
		return
	}

	game, err := h.svc.UpdateGame(c.Request.Context(), id, service.UpdateGameInput{
		Key:         payload.Key,
		Name:        payload.Name,
		Category:    payload.Category,
		IconURL:     payload.IconURL,
		Description: payload.Description,
	})
    if errors.Is(err, service.ErrValidation) {
        _ = c.Error(service.ErrValidation)
        return
    }
    if errors.Is(err, service.ErrNotFound) {
        _ = c.Error(service.ErrNotFound)
        return
    }
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[*model.Game]{
		Success: true,
		Code:    http.StatusOK,
		Message: "updated",
		Data:    game,
	})
}

// DeleteGame
// @Summary      删除游戏
// @Tags         Admin/Games
// @Security     BearerAuth
// @Param        id   path  int  true  "游戏ID"
// @Produce      json
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/games/{id} [delete]
//
// DeleteGame 删除游戏。
func (h *GameHandler) DeleteGame(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}

    if err := h.svc.DeleteGame(c.Request.Context(), id); errors.Is(err, service.ErrNotFound) {
        _ = c.Error(service.ErrNotFound)
        return
    } else if err != nil {
        writeJSONError(c, http.StatusInternalServerError, err.Error())
        return
    }

	writeJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "deleted",
	})
}

// GamePayload defines request body for creating/updating a game.
type GamePayload struct {
	Key         string `json:"key" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Category    string `json:"category"`
	IconURL     string `json:"icon_url"`
	Description string `json:"description"`
}
