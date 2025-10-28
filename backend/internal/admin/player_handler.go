package admin

import (
    "errors"
    "net/http"

    "github.com/gin-gonic/gin"

    apierr "gamelink/internal/handler"
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

// ListPlayers
// @Summary      列出玩家资料
// @Tags         Admin/Players
// @Security     BearerAuth
// @Param        page       query  int  false  "页码"
// @Param        page_size  query  int  false  "每页数量"
// @Produce      json
// @Success      200  {object}  map[string]any
// @Router       /admin/players [get]
//
// ListPlayers returns a paginated list of players.
func (h *PlayerHandler) ListPlayers(c *gin.Context) {
    page, err := queryIntDefault(c, "page", 1)
    if err != nil {
        writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidPage)
        return
    }
    pageSize, err := queryIntDefault(c, "page_size", 20)
    if err != nil {
        writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidPageSize)
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

// GetPlayer
// @Summary      获取玩家资料
// @Tags         Admin/Players
// @Security     BearerAuth
// @Param        id   path  int  true  "玩家ID"
// @Produce      json
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/players/{id} [get]
//
// GetPlayer returns a single player by id.
func (h *PlayerHandler) GetPlayer(c *gin.Context) {
    id, err := parseUintParam(c, "id")
    if err != nil {
        writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
        return
    }
	player, err := h.svc.GetPlayer(c.Request.Context(), id)
    if errors.Is(err, service.ErrNotFound) {
        _ = c.Error(service.ErrNotFound)
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

// CreatePlayer
// @Summary      新建玩家资料
// @Tags         Admin/Players
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  CreatePlayerPayload  true  "玩家信息"
// @Success      201  {object}  map[string]any
// @Failure      400  {object}  map[string]any
// @Router       /admin/players [post]
//
// CreatePlayer creates a new player profile.
func (h *PlayerHandler) CreatePlayer(c *gin.Context) {
	var payload CreatePlayerPayload
    if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
        writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload)
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
        _ = c.Error(service.ErrValidation)
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

// UpdatePlayer
// @Summary      更新玩家资料
// @Tags         Admin/Players
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                   true  "玩家ID"
// @Param        request  body  UpdatePlayerPayload   true  "玩家信息"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/players/{id} [put]
//
// UpdatePlayer updates player profile.
func (h *PlayerHandler) UpdatePlayer(c *gin.Context) {
    id, err := parseUintParam(c, "id")
    if err != nil {
        writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
        return
    }

	var payload UpdatePlayerPayload
    if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
        writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload)
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

	writeJSON(c, http.StatusOK, model.APIResponse[*model.Player]{
		Success: true,
		Code:    http.StatusOK,
		Message: "updated",
		Data:    player,
	})
}

// DeletePlayer
// @Summary      删除玩家资料
// @Tags         Admin/Players
// @Security     BearerAuth
// @Param        id   path  int  true  "玩家ID"
// @Produce      json
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/players/{id} [delete]
//
// DeletePlayer deletes a player profile by id.
func (h *PlayerHandler) DeletePlayer(c *gin.Context) {
    id, err := parseUintParam(c, "id")
    if err != nil {
        writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
        return
    }
    if err := h.svc.DeletePlayer(c.Request.Context(), id); errors.Is(err, service.ErrNotFound) {
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

// UpdatePlayerVerification
// @Summary      更新玩家认证状态
// @Tags         Admin/Players
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int  true  "玩家ID"
// @Param        request  body  map[string]string  true  "{verification_status}"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/players/{id}/verification [put]
func (h *PlayerHandler) UpdatePlayerVerification(c *gin.Context) {
    id, err := parseUintParam(c, "id")
    if err != nil { writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID); return }
    var payload struct{ VerificationStatus string `json:"verification_status" binding:"required"` }
    if bindErr := c.ShouldBindJSON(&payload); bindErr != nil { writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload); return }

    player, err := h.svc.GetPlayer(c.Request.Context(), id)
    if errors.Is(err, service.ErrNotFound) { _ = c.Error(service.ErrNotFound); return }
    if err != nil { writeJSONError(c, http.StatusInternalServerError, err.Error()); return }

    out, err := h.svc.UpdatePlayer(c.Request.Context(), id, service.UpdatePlayerInput{
        Nickname:           player.Nickname,
        Bio:                player.Bio,
        HourlyRateCents:    player.HourlyRateCents,
        MainGameID:         player.MainGameID,
        VerificationStatus: model.VerificationStatus(payload.VerificationStatus),
    })
    if errors.Is(err, service.ErrValidation) { _ = c.Error(service.ErrValidation); return }
    if err != nil { writeJSONError(c, http.StatusInternalServerError, err.Error()); return }
    writeJSON(c, http.StatusOK, model.APIResponse[*model.Player]{ Success: true, Code: http.StatusOK, Message: "updated", Data: out })
}

// UpdatePlayerGames
// @Summary      更新玩家主游戏
// @Tags         Admin/Players
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int  true  "玩家ID"
// @Param        request  body  map[string]uint64  true  "{main_game_id}"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/players/{id}/games [put]
func (h *PlayerHandler) UpdatePlayerGames(c *gin.Context) {
    id, err := parseUintParam(c, "id")
    if err != nil { writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID); return }
    var payload struct{ MainGameID uint64 `json:"main_game_id" binding:"required"` }
    if bindErr := c.ShouldBindJSON(&payload); bindErr != nil { writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload); return }

    player, err := h.svc.GetPlayer(c.Request.Context(), id)
    if errors.Is(err, service.ErrNotFound) { _ = c.Error(service.ErrNotFound); return }
    if err != nil { writeJSONError(c, http.StatusInternalServerError, err.Error()); return }

    out, err := h.svc.UpdatePlayer(c.Request.Context(), id, service.UpdatePlayerInput{
        Nickname:           player.Nickname,
        Bio:                player.Bio,
        HourlyRateCents:    player.HourlyRateCents,
        MainGameID:         payload.MainGameID,
        VerificationStatus: player.VerificationStatus,
    })
    if errors.Is(err, service.ErrValidation) { _ = c.Error(service.ErrValidation); return }
    if err != nil { writeJSONError(c, http.StatusInternalServerError, err.Error()); return }
    writeJSON(c, http.StatusOK, model.APIResponse[*model.Player]{ Success: true, Code: http.StatusOK, Message: "updated", Data: out })
}

// UpdatePlayerSkillTags
// @Summary      更新玩家技能标签
// @Tags         Admin/Players
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int            true  "玩家ID"
// @Param        request  body  SkillTagsBody  true  "标签集合"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/players/{id}/skill-tags [put]
func (h *PlayerHandler) UpdatePlayerSkillTags(c *gin.Context) {
    id, err := parseUintParam(c, "id")
    if err != nil { writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID); return }
    var body SkillTagsBody
    if bindErr := c.ShouldBindJSON(&body); bindErr != nil { writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload); return }
    // Ensure player exists first
    if _, err := h.svc.GetPlayer(c.Request.Context(), id); errors.Is(err, service.ErrNotFound) { _ = c.Error(service.ErrNotFound); return } else if err != nil { writeJSONError(c, http.StatusInternalServerError, err.Error()); return }
    if err := h.svc.UpdatePlayerSkillTags(c.Request.Context(), id, body.Tags); err != nil { writeJSONError(c, http.StatusInternalServerError, err.Error()); return }
    writeJSON(c, http.StatusOK, model.APIResponse[any]{ Success: true, Code: http.StatusOK, Message: "updated" })
}

type SkillTagsBody struct { Tags []string `json:"tags" binding:"required"` }

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
