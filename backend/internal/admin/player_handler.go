package admin

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/service"
)

// PlayerHandler 处理陪玩资料管理接口。
type PlayerHandler struct {
	svc *service.AdminService
}

// NewPlayerHandler 创建 Handler。
func NewPlayerHandler(svc *service.AdminService) *PlayerHandler {
	return &PlayerHandler{svc: svc}
}

// ListPlayers returns a paginated list of players.
func (h *PlayerHandler) ListPlayers(c *gin.Context) {
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

	players, pagination, err := h.svc.ListPlayersPaged(c.Request.Context(), page, pageSize)
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(c, http.StatusOK, model.APIResponse[[]model.Player]{
		Success:    true,
		Code:       http.StatusOK,
		Message:    "OK",
		Data:       players,
		Pagination: pagination,
	})
}

// GetPlayer returns a single player by id.
func (h *PlayerHandler) GetPlayer(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	player, err := h.svc.GetPlayer(c.Request.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		writeJSONError(c, http.StatusNotFound, "player not found")
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(c, http.StatusOK, model.APIResponse[*model.Player]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    player,
	})
}

// CreatePlayer creates a new player profile.
func (h *PlayerHandler) CreatePlayer(c *gin.Context) {
	var payload CreatePlayerPayload
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid json payload")
		return
	}

	player, err := h.svc.CreatePlayer(c.Request.Context(), service.CreatePlayerInput{
		UserID:             payload.UserID,
		Nickname:           payload.Nickname,
		Bio:                payload.Bio,
		HourlyRateCents:    payload.HourlyRateCents,
		MainGameID:         payload.MainGameID,
		VerificationStatus: model.VerificationStatus(payload.VerificationStatus),
	})
	if errors.Is(err, service.ErrValidation) {
		writeJSONError(c, http.StatusBadRequest, "缺少必填字段")
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusCreated, model.APIResponse[*model.Player]{
		Success: true,
		Code:    http.StatusCreated,
		Message: "created",
		Data:    player,
	})
}

// UpdatePlayer updates player profile.
func (h *PlayerHandler) UpdatePlayer(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid id")
		return
	}

	var payload UpdatePlayerPayload
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid json payload")
		return
	}

	player, err := h.svc.UpdatePlayer(c.Request.Context(), id, service.UpdatePlayerInput{
		Nickname:           payload.Nickname,
		Bio:                payload.Bio,
		HourlyRateCents:    payload.HourlyRateCents,
		MainGameID:         payload.MainGameID,
		VerificationStatus: model.VerificationStatus(payload.VerificationStatus),
	})
	if errors.Is(err, service.ErrValidation) {
		writeJSONError(c, http.StatusBadRequest, "缺少必填字段")
		return
	}
	if errors.Is(err, service.ErrNotFound) {
		writeJSONError(c, http.StatusNotFound, "player not found")
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[*model.Player]{
		Success: true,
		Code:    http.StatusOK,
		Message: "updated",
		Data:    player,
	})
}

// DeletePlayer deletes a player profile by id.
func (h *PlayerHandler) DeletePlayer(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.DeletePlayer(c.Request.Context(), id); errors.Is(err, service.ErrNotFound) {
		writeJSONError(c, http.StatusNotFound, "player not found")
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

// CreatePlayerPayload defines the request body for creating a player.
type CreatePlayerPayload struct {
	UserID             uint64 `json:"user_id" binding:"required"`
	Nickname           string `json:"nickname"`
	Bio                string `json:"bio"`
	HourlyRateCents    int64  `json:"hourly_rate_cents"`
	MainGameID         uint64 `json:"main_game_id"`
	VerificationStatus string `json:"verification_status" binding:"required"`
}

// UpdatePlayerPayload defines the request body for updating a player.
type UpdatePlayerPayload struct {
	Nickname           string `json:"nickname"`
	Bio                string `json:"bio"`
	HourlyRateCents    int64  `json:"hourly_rate_cents"`
	MainGameID         uint64 `json:"main_game_id"`
	VerificationStatus string `json:"verification_status" binding:"required"`
}
