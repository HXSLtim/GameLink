package admin

import (
	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/admin/dto"
	"gamelink/internal/model"
	adminservice "gamelink/internal/service/admin"
	"gamelink/pkg/apierr"
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
// @Success      200  {array}   model.Player
// @Router       /admin/players [get]
//
// ListPlayers returns a paginated list of players.
func (h *PlayerHandler) ListPlayers(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	// Parse optional filters
	keyword := c.Query("keyword")
	statusStr := c.Query("status")

	var status *model.VerificationStatus
	if statusStr != "" {
		s := model.VerificationStatus(statusStr)
		status = &s
	}

	// Use filter method if any filter is provided
	if keyword != "" || status != nil {
		players, pagination, err := h.svc.ListPlayersPagedWithFilter(c.Request.Context(), page, pageSize, keyword, status)
		if err != nil {
			respondError(c, apierr.InternalError("list players failed").WithDetails(err.Error()))
			return
		}
		respondList(c, dto.ToPlayerResponseList(players), pagination)
		return
	}

	players, pagination, err := h.svc.ListPlayersPaged(c.Request.Context(), page, pageSize)
	if err != nil {
		respondError(c, apierr.InternalError("list players failed").WithDetails(err.Error()))
		return
	}
	respondList(c, dto.ToPlayerResponseList(players), pagination)
}

// GetPlayer
// @Summary      获取玩家资料
// @Tags         Admin/Players
// @Security     BearerAuth
// @Param        id   path  int  true  "玩家ID"
// @Produce      json
// @Success      200  {object}  Player
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/players/{id} [get]
//
// GetPlayer returns a single player by id.
func (h *PlayerHandler) GetPlayer(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}
	player, err := h.svc.GetPlayer(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, dto.ToPlayerResponse(player))
}

// CreatePlayer
// @Summary      新建玩家资料
// @Tags         Admin/Players
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  CreatePlayerPayload  true  "玩家信息"
// @Success      201  {object}  Player
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/players [post]
//
// CreatePlayer creates a new player profile.
func (h *PlayerHandler) CreatePlayer(c *gin.Context) {
	var payload CreatePlayerPayload
	if !ValidateAndRespond(c, &payload) {
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
		respondError(c, err)
		return
	}
	respondCreated(c, dto.ToPlayerResponse(player))
}

// UpdatePlayer
// @Summary      更新玩家资料
// @Tags         Admin/Players
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                   true  "玩家ID"
// @Param        request  body  UpdatePlayerPayload   true  "玩家信息"
// @Success      200  {object}  Player
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/players/{id} [put]
//
// UpdatePlayer updates player profile.
func (h *PlayerHandler) UpdatePlayer(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var payload UpdatePlayerPayload
	if !ValidateAndRespond(c, &payload) {
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
		respondError(c, err)
		return
	}
	respondUpdated(c, dto.ToPlayerResponse(player))
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
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}
	if err := h.svc.DeletePlayer(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	respondDeleted(c)
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
// @Success      200  {array}   model.OperationLog
// @Router       /admin/players/{id}/logs [get]
func (h *PlayerHandler) ListPlayerLogs(c *gin.Context) {
	handleOperationLogList(c, "player", h.svc.ListOperationLogs)
}

// UpdatePlayerVerification
// @Summary      更新玩家认证状态（审核）
// @Tags         Admin/Players
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int  true  "玩家ID"
// @Param        request  body  UpdateVerificationPayload  true  "审核信息"
// @Success      200  {object}  Player
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/players/{id}/verification [put]
func (h *PlayerHandler) UpdatePlayerVerification(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	// 获取审核人ID
	adminID, ok := getAdminUserID(c)
	if !ok {
		return
	}

	var payload UpdateVerificationPayload
	if !ValidateAndRespond(c, &payload) {
		return
	}

	player, err := h.svc.GetPlayer(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}

	out, err := h.svc.UpdatePlayerVerification(c.Request.Context(), id, adminservice.UpdateVerificationInput{
		Nickname:           player.Nickname,
		Bio:                player.Bio,
		HourlyRateCents:    player.HourlyRateCents,
		MainGameID:         player.MainGameID,
		VerificationStatus: model.VerificationStatus(payload.VerificationStatus),
		VerifiedBy:         adminID,
		Remark:             payload.Remark,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	respondUpdated(c, dto.ToPlayerResponse(out))
}

// UpdateVerificationPayload
type UpdateVerificationPayload struct {
	VerificationStatus string `json:"verification_status" binding:"required,oneof=pending verified rejected"`
	Remark             string `json:"remark"`
}

// UpdatePlayerGames
// @Summary      更新玩家主游
// @Tags         Admin/Players
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int  true  "玩家ID"
// @Param        request  body  map[string]uint64  true  "{main_game_id}"
// @Success      200  {object}  Player
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/players/{id}/games [put]
func (h *PlayerHandler) UpdatePlayerGames(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}
	var payload struct {
		MainGameID uint64 `json:"main_game_id" binding:"required"`
	}
	if !ValidateAndRespond(c, &payload) {
		return
	}

	player, err := h.svc.GetPlayer(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
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
		respondError(c, err)
		return
	}
	respondUpdated(c, dto.ToPlayerResponse(out))
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
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}
	var body SkillTagsBody
	if !ValidateAndRespond(c, &body) {
		return
	}
	// Ensure player exists first
	if _, err := h.svc.GetPlayer(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	if err := h.svc.UpdatePlayerSkillTags(c.Request.Context(), id, body.Tags); err != nil {
		respondError(c, apierr.InternalError("update player skill tags failed").WithDetails(err.Error()))
		return
	}
	respondSuccessWithMsg[any](c, "updated", nil)
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

// BatchUpdateStatusPayload 批量更新状态请求
type BatchUpdateStatusPayload struct {
	PlayerIDs []string `json:"playerIds" binding:"required"`
	Status    string   `json:"status" binding:"required,oneof=pending verified rejected"`
}

// BatchUpdatePlayerStatus
// @Summary      批量更新陪玩师状态
// @Tags         Admin/Players
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BatchUpdateStatusPayload  true  "批量更新请求"
// @Success      200  {object}  map[string]int64
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/players/batch/status [put]
func (h *PlayerHandler) BatchUpdatePlayerStatus(c *gin.Context) {
	var payload BatchUpdateStatusPayload
	if !ValidateAndRespond(c, &payload) {
		return
	}

	// Convert string IDs to uint64
	ids := make([]uint64, 0, len(payload.PlayerIDs))
	for _, idStr := range payload.PlayerIDs {
		id, err := parseUint64(idStr)
		if err != nil {
			respondError(c, apierr.BadRequest("invalid player id: "+idStr))
			return
		}
		ids = append(ids, id)
	}

	status := model.VerificationStatus(payload.Status)
	updated, err := h.svc.BatchUpdatePlayerStatus(c.Request.Context(), ids, status)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, map[string]int64{"updated": updated})
}

// BatchDeletePlayersPayload 批量删除请求
type BatchDeletePlayersPayload struct {
	PlayerIDs []string `json:"playerIds" binding:"required"`
}

// BatchDeletePlayers
// @Summary      批量删除陪玩师
// @Tags         Admin/Players
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BatchDeletePlayersPayload  true  "批量删除请求"
// @Success      200  {object}  map[string]int64
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/players/batch/delete [post]
func (h *PlayerHandler) BatchDeletePlayers(c *gin.Context) {
	var payload BatchDeletePlayersPayload
	if !ValidateAndRespond(c, &payload) {
		return
	}

	// Convert string IDs to uint64
	ids := make([]uint64, 0, len(payload.PlayerIDs))
	for _, idStr := range payload.PlayerIDs {
		id, err := parseUint64(idStr)
		if err != nil {
			respondError(c, apierr.BadRequest("invalid player id: "+idStr))
			return
		}
		ids = append(ids, id)
	}

	deleted, err := h.svc.BatchDeletePlayers(c.Request.Context(), ids)
	if err != nil {
		respondError(c, err)
		return
	}
	respondSuccess(c, map[string]int64{"deleted": deleted})
}
