package admin

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

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

// ListGames 返回全部游戏。
func (h *GameHandler) ListGames(c *gin.Context) {
	page, err := queryIntDefault(c, "page", 1)
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid page")
		return
	}
	pageSize, err := queryIntDefault(c, "page_size", 20)
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid page_size")
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

// GetGame 获取单个游戏。
func (h *GameHandler) GetGame(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	game, err := h.svc.GetGame(c.Request.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		writeJSONError(c, http.StatusNotFound, "game not found")
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

// CreateGame 创建新游戏。
func (h *GameHandler) CreateGame(c *gin.Context) {
	var payload GamePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid json payload")
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
		writeJSONError(c, http.StatusBadRequest, "key 与 name 为必填字段")
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

// UpdateGame 更新游戏信息。
func (h *GameHandler) UpdateGame(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid id")
		return
	}

	var payload GamePayload
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid json payload")
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
		writeJSONError(c, http.StatusBadRequest, "key 与 name 为必填字段")
		return
	}
	if errors.Is(err, service.ErrNotFound) {
		writeJSONError(c, http.StatusNotFound, "game not found")
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

// DeleteGame 删除游戏。
func (h *GameHandler) DeleteGame(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.svc.DeleteGame(c.Request.Context(), id); errors.Is(err, service.ErrNotFound) {
		writeJSONError(c, http.StatusNotFound, "game not found")
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
