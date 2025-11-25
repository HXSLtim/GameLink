package admin

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/apierr"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	adminservice "gamelink/internal/service/admin"
)

// Player 陪玩师模型（类型别名）
type Player = model.Player

// PlayerHandler 处理陪玩资料管理接口
type PlayerHandler struct {
	svc *adminservice.AdminService
}

// NewPlayerHandler 创建Handler
func NewPlayerHandler(svc *adminservice.AdminService) *PlayerHandler {
	return &PlayerHandler{svc: svc}
}

// ListPlayers
// @Summary      列出玩家资料
// @Tags         Admin/Players
// @Security     BearerAuth
// @Param        page       query  int  false  "页码"
// @Param        pageSize   query     int       false  "每页数量"
// @Produce      json
// @Success      200  {object}  model.APIResponse[[]model.Player]
// @Router       /admin/players [get]
//
// ListPlayers returns a paginated list of players.
func (h *PlayerHandler) ListPlayers(c *gin.Context) {
	page, err := queryIntDefault(c, "page", 1)
	if err != nil {
		respondAPIError(c, apierr.BadRequest("invalid page parameter"))
		return
	}
	pageSize, err := queryIntDefault(c, "page_size", 20)
	if err != nil {
		respondAPIError(c, apierr.BadRequest("invalid page size parameter"))
		return
	}

	players, pagination, err := h.svc.ListPlayersPaged(c.Request.Context(), page, pageSize)
	if err != nil {
		respondAPIError(c, apierr.InternalError("list players failed").WithDetails(err.Error()))
		return
	}
	players = ensureSlice(players)
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
// @Success      200  {object}  model.APIResponse[Player]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/players/{id} [get]
//
// GetPlayer returns a single player by id.
func (h *PlayerHandler) GetPlayer(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		respondAPIError(c, apierr.BadRequest("invalid player ID"))
		return
	}
	player, err := h.svc.GetPlayer(c.Request.Context(), id)
	if err != nil {
		if apierr.IsNotFound(err) {
			respondAPIError(c, err)
			return
		}
		respondAPIError(c, apierr.InternalError("get player failed").WithDetails(err.Error()))
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
// @Success      201  {object}  model.APIResponse[Player]
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/players [post]
//
// CreatePlayer creates a new player profile.
func (h *PlayerHandler) CreatePlayer(c *gin.Context) {
	var payload CreatePlayerPayload
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		respondAPIError(c, apierr.BadRequest("invalid JSON payload"))
		return
	}

	player, err := h.svc.CreatePlayer(c.Request.Context(), adminservice.CreatePlayerInput{
		UserID:             payload.UserID,
		Nickname:           payload.Nickname,
		Bio:                payload.Bio,
		Rank:               payload.Rank,
		HourlyRateCents:    payload.HourlyRateCents,
		MainGameID:         payload.MainGameID,
		VerificationStatus: model.VerificationStatus(payload.VerificationStatus),
	})
	if err != nil {
		if apierr.IsValidationError(err) {
			respondAPIError(c, err)
			return
		}
		respondAPIError(c, apierr.InternalError("create player failed").WithDetails(err.Error()))
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
// @Success      200  {object}  model.APIResponse[Player]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/players/{id} [put]
//
// UpdatePlayer updates player profile.
func (h *PlayerHandler) UpdatePlayer(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		respondAPIError(c, apierr.BadRequest("invalid player ID"))
		return
	}

	var payload UpdatePlayerPayload
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		respondAPIError(c, apierr.BadRequest("invalid JSON payload"))
		return
	}

	player, err := h.svc.UpdatePlayer(c.Request.Context(), id, adminservice.UpdatePlayerInput{
		Nickname:           payload.Nickname,
		Bio:                payload.Bio,
		Rank:               payload.Rank,
		HourlyRateCents:    payload.HourlyRateCents,
		MainGameID:         payload.MainGameID,
		VerificationStatus: model.VerificationStatus(payload.VerificationStatus),
	})
	if err != nil {
		if apierr.IsValidationError(err) {
			respondAPIError(c, err)
			return
		}
		if apierr.IsNotFound(err) {
			respondAPIError(c, err)
			return
		}
		respondAPIError(c, apierr.InternalError("update player failed").WithDetails(err.Error()))
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
// @Success      200  {object}  model.SuccessResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/players/{id} [delete]
//
// DeletePlayer deletes a player profile by id.
func (h *PlayerHandler) DeletePlayer(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		respondAPIError(c, apierr.BadRequest("invalid player ID"))
		return
	}
	if err := h.svc.DeletePlayer(c.Request.Context(), id); err != nil {
		if apierr.IsNotFound(err) {
			respondAPIError(c, err)
			return
		}
		respondAPIError(c, apierr.InternalError("delete player failed").WithDetails(err.Error()))
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "deleted",
	})
}

// ListPlayerLogs
// @Summary      获取玩家操作日志
// @Tags         Admin/Players
// @Security     BearerAuth
// @Produce      json
// @Param        id           path   int  true  "玩家ID"
// @Param        page         query  int  false "页码"
// @Param        pageSize   query     int       false  "每页数量"
// @Param        action       query  string false "动作过滤" Enums(create,update,delete)
// @Param        actor_user_id query int   false "操作者用户ID"
// @Param        dateFrom       query    string       false  "Start date (YYYY-MM-DD)"// @Param        dateTo     query     string    false  "End date (YYYY-MM-DD)"
// @Param        export       query  string false "导出格式" Enums(csv)
// @Param        fields         query    string       false  "Export fields (comma separated)"// @Param        header_lang  query  string false "列头语言" Enums(en,zh)
// @Success      200  {object}  model.APIResponse[[]model.OperationLog]
// @Router       /admin/players/{id}/logs [get]
func (h *PlayerHandler) ListPlayerLogs(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		respondAPIError(c, apierr.BadRequest("invalid player ID"))
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
	items, p, err := h.svc.ListOperationLogs(c.Request.Context(), "player", id, opts)
	if err != nil {
		respondAPIError(c, apierr.InternalError("list player logs failed").WithDetails(err.Error()))
		return
	}
	if strings.EqualFold(strings.TrimSpace(c.Query("export")), "csv") {
		exportOperationLogsCSV(c, "player", id, items)
		return
	}
	items = ensureSlice(items)
	writeJSON(c, 200, model.APIResponse[[]model.OperationLog]{Success: true, Code: 200, Message: "OK", Data: items, Pagination: p})
}

// UpdatePlayerVerification
// @Summary      更新玩家认证状
// @Tags         Admin/Players
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int  true  "玩家ID"
// @Param        request  body  map[string]string  true  "{verification_status}"
// @Success      200  {object}  model.APIResponse[Player]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/players/{id}/verification [put]
func (h *PlayerHandler) UpdatePlayerVerification(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		respondAPIError(c, apierr.BadRequest("invalid player ID"))
		return
	}
	var payload struct {
		VerificationStatus string `json:"verification_status" binding:"required"`
	}
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		respondAPIError(c, apierr.BadRequest("invalid JSON payload"))
		return
	}

	player, err := h.svc.GetPlayer(c.Request.Context(), id)
	if err != nil {
		if apierr.IsNotFound(err) {
			respondAPIError(c, err)
			return
		}
		respondAPIError(c, apierr.InternalError("get player failed").WithDetails(err.Error()))
		return
	}

	out, err := h.svc.UpdatePlayer(c.Request.Context(), id, adminservice.UpdatePlayerInput{
		Nickname:           player.Nickname,
		Bio:                player.Bio,
		HourlyRateCents:    player.HourlyRateCents,
		MainGameID:         player.MainGameID,
		VerificationStatus: model.VerificationStatus(payload.VerificationStatus),
	})
	if err != nil {
		if apierr.IsValidationError(err) {
			respondAPIError(c, err)
			return
		}
		respondAPIError(c, apierr.InternalError("update player verification failed").WithDetails(err.Error()))
		return
	}
	writeJSON(c, http.StatusOK, model.APIResponse[*model.Player]{Success: true, Code: http.StatusOK, Message: "updated", Data: out})
}

// UpdatePlayerGames
// @Summary      更新玩家主游
// @Tags         Admin/Players
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int  true  "玩家ID"
// @Param        request  body  map[string]uint64  true  "{main_game_id}"
// @Success      200  {object}  model.APIResponse[Player]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/players/{id}/games [put]
func (h *PlayerHandler) UpdatePlayerGames(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		respondAPIError(c, apierr.BadRequest("invalid player ID"))
		return
	}
	var payload struct {
		MainGameID uint64 `json:"main_game_id" binding:"required"`
	}
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		respondAPIError(c, apierr.BadRequest("invalid JSON payload"))
		return
	}

	player, err := h.svc.GetPlayer(c.Request.Context(), id)
	if err != nil {
		if apierr.IsNotFound(err) {
			respondAPIError(c, err)
			return
		}
		respondAPIError(c, apierr.InternalError("get player failed").WithDetails(err.Error()))
		return
	}

	out, err := h.svc.UpdatePlayer(c.Request.Context(), id, adminservice.UpdatePlayerInput{
		Nickname:           player.Nickname,
		Bio:                player.Bio,
		HourlyRateCents:    player.HourlyRateCents,
		MainGameID:         payload.MainGameID,
		VerificationStatus: player.VerificationStatus,
	})
	if err != nil {
		if apierr.IsValidationError(err) {
			respondAPIError(c, err)
			return
		}
		respondAPIError(c, apierr.InternalError("update player games failed").WithDetails(err.Error()))
		return
	}
	writeJSON(c, http.StatusOK, model.APIResponse[*model.Player]{Success: true, Code: http.StatusOK, Message: "updated", Data: out})
}

// UpdatePlayerSkillTags
// @Summary      更新玩家技能标
// @Tags         Admin/Players
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int            true  "玩家ID"
// @Param        request  body  SkillTagsBody  true  "标签集合"
// @Success      200  {object}  model.SuccessResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/players/{id}/skill-tags [put]
func (h *PlayerHandler) UpdatePlayerSkillTags(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		respondAPIError(c, apierr.BadRequest("invalid player ID"))
		return
	}
	var body SkillTagsBody
	if bindErr := c.ShouldBindJSON(&body); bindErr != nil {
		respondAPIError(c, apierr.BadRequest("invalid JSON payload"))
		return
	}
	// Ensure player exists first
	if _, err := h.svc.GetPlayer(c.Request.Context(), id); err != nil {
		if apierr.IsNotFound(err) {
			respondAPIError(c, err)
			return
		}
		respondAPIError(c, apierr.InternalError("get player failed").WithDetails(err.Error()))
		return
	}
	if err := h.svc.UpdatePlayerSkillTags(c.Request.Context(), id, body.Tags); err != nil {
		respondAPIError(c, apierr.InternalError("update player skill tags failed").WithDetails(err.Error()))
		return
	}
	writeJSON(c, http.StatusOK, model.APIResponse[any]{Success: true, Code: http.StatusOK, Message: "updated"})
}

type SkillTagsBody struct {
	Tags []string `json:"tags" binding:"required"`
}

// CreatePlayerPayload defines the request body for creating a player.
type CreatePlayerPayload struct {
	UserID             uint64 `json:"user_id" binding:"required"`
	Nickname           string `json:"nickname"`
	Bio                string `json:"bio"`
	Rank               string `json:"rank"`
	HourlyRateCents    int64  `json:"hourly_rate_cents"`
	MainGameID         uint64 `json:"main_game_id"`
	VerificationStatus string `json:"verification_status" binding:"required"`
}

// UpdatePlayerPayload defines the request body for updating a player.
type UpdatePlayerPayload struct {
	Nickname           string `json:"nickname"`
	Bio                string `json:"bio"`
	Rank               string `json:"rank"`
	HourlyRateCents    int64  `json:"hourly_rate_cents"`
	MainGameID         uint64 `json:"main_game_id"`
	VerificationStatus string `json:"verification_status" binding:"required"`
}
