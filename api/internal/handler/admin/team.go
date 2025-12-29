package admin

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/middleware"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	teamrepo "gamelink/internal/repository/team"
	teamservice "gamelink/internal/service/team"
	"gamelink/pkg/apierr"
)

// TeamHandler 团队管理处理器
type TeamHandler struct {
	svc *teamservice.TeamService
}

// NewTeamHandler 创建团队处理器
func NewTeamHandler(svc *teamservice.TeamService) *TeamHandler {
	return &TeamHandler{svc: svc}
}

// ============================================================================
// 团队管理
// ============================================================================

// TeamCreateRequest 创建团队请求
type TeamCreateRequest struct {
	Name            string  `json:"name" binding:"required,max=64"`
	Description     string  `json:"description" binding:"max=255"`
	AvatarURL       string  `json:"avatarUrl"`
	LeaderPlayerID  uint64  `json:"leaderPlayerId" binding:"required"`
	MaxMembers      int     `json:"maxMembers"`
	IncomeShareType string  `json:"incomeShareType"`
	LeaderBonusRate float64 `json:"leaderBonusRate"`
}

// TeamUpdateRequest 更新团队请求
type TeamUpdateRequest struct {
	Name            string  `json:"name" binding:"required,max=64"`
	Description     string  `json:"description" binding:"max=255"`
	AvatarURL       string  `json:"avatarUrl"`
	MaxMembers      int     `json:"maxMembers"`
	IncomeShareType string  `json:"incomeShareType"`
	LeaderBonusRate float64 `json:"leaderBonusRate"`
}

// ListTeams 获取团队列表
func (h *TeamHandler) ListTeams(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	var status *model.TeamStatus
	if v := strings.TrimSpace(c.Query("status")); v != "" {
		s := model.TeamStatus(v)
		status = &s
	}

	var leaderID *uint64
	if v := strings.TrimSpace(c.Query("leaderId")); v != "" {
		id, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			respondBadRequest(c, "无效的队长ID")
			return
		}
		leaderID = &id
	}

	var minMember, maxMember *int
	if v := strings.TrimSpace(c.Query("minMember")); v != "" {
		n, err := strconv.Atoi(v)
		if err == nil {
			minMember = &n
		}
	}
	if v := strings.TrimSpace(c.Query("maxMember")); v != "" {
		n, err := strconv.Atoi(v)
		if err == nil {
			maxMember = &n
		}
	}

	opts := teamrepo.TeamListOptions{
		Page:      page,
		PageSize:  pageSize,
		Keyword:   c.Query("keyword"),
		Status:    status,
		LeaderID:  leaderID,
		MinMember: minMember,
		MaxMember: maxMember,
	}

	teams, total, err := h.svc.ListTeams(c.Request.Context(), opts)
	if err != nil {
		respondError(c, apierr.InternalError("获取团队列表失败").WithDetails(err.Error()))
		return
	}

	respondList(c, teams, &model.Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
	})
}

// GetTeam 获取团队详情
func (h *TeamHandler) GetTeam(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	team, err := h.svc.GetTeam(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("团队不存在"))
			return
		}
		respondError(c, apierr.InternalError("获取团队失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, team)
}

// CreateTeam 创建团队
func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var req TeamCreateRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	team := &model.Team{
		Name:            req.Name,
		Description:     req.Description,
		AvatarURL:       req.AvatarURL,
		MaxMembers:      req.MaxMembers,
		IncomeShareType: req.IncomeShareType,
		LeaderBonusRate: req.LeaderBonusRate,
	}

	if err := h.svc.CreateTeam(c.Request.Context(), team, req.LeaderPlayerID); err != nil {
		respondError(c, apierr.BadRequest(err.Error()))
		return
	}

	respondCreated(c, team)
}

// UpdateTeam 更新团队
func (h *TeamHandler) UpdateTeam(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var req TeamUpdateRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	team := &model.Team{
		Name:            req.Name,
		Description:     req.Description,
		AvatarURL:       req.AvatarURL,
		MaxMembers:      req.MaxMembers,
		IncomeShareType: req.IncomeShareType,
		LeaderBonusRate: req.LeaderBonusRate,
	}
	team.ID = id

	if err := h.svc.UpdateTeam(c.Request.Context(), team); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("团队不存在"))
			return
		}
		respondError(c, apierr.BadRequest(err.Error()))
		return
	}

	respondUpdated(c, team)
}

// DeleteTeam 删除团队
func (h *TeamHandler) DeleteTeam(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	if err := h.svc.DeleteTeam(c.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("团队不存在"))
			return
		}
		respondError(c, apierr.BadRequest(err.Error()))
		return
	}

	respondDeleted(c)
}

// UpdateTeamStatusRequest 更新团队状态请求
type UpdateTeamStatusRequest struct {
	Status model.TeamStatus `json:"status" binding:"required"`
}

// UpdateTeamStatus 更新团队状态
func (h *TeamHandler) UpdateTeamStatus(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var req UpdateTeamStatusRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	if err := h.svc.UpdateTeamStatus(c.Request.Context(), id, req.Status); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("团队不存在"))
			return
		}
		respondError(c, apierr.BadRequest(err.Error()))
		return
	}

	respondMsg(c, "状态更新成功")
}

// GetTeamStats 获取团队统计
func (h *TeamHandler) GetTeamStats(c *gin.Context) {
	stats, err := h.svc.GetTeamStats(c.Request.Context())
	if err != nil {
		respondError(c, apierr.InternalError("获取统计失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, stats)
}

// ============================================================================
// 成员管理
// ============================================================================

// AddMemberRequest 添加成员请求
type AddMemberRequest struct {
	PlayerID uint64 `json:"playerId" binding:"required"`
}

// GetTeamMembers 获取团队成员
func (h *TeamHandler) GetTeamMembers(c *gin.Context) {
	teamID, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	members, err := h.svc.GetTeamMembers(c.Request.Context(), teamID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("团队不存在"))
			return
		}
		respondError(c, apierr.InternalError("获取成员列表失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, members)
}

// AddMember 添加成员
func (h *TeamHandler) AddMember(c *gin.Context) {
	teamID, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var req AddMemberRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	if err := h.svc.AddMember(c.Request.Context(), teamID, req.PlayerID); err != nil {
		respondError(c, apierr.BadRequest(err.Error()))
		return
	}

	respondMsg(c, "成员添加成功")
}

// RemoveMember 移除成员
func (h *TeamHandler) RemoveMember(c *gin.Context) {
	teamID, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	playerID, ok := ParseIDAndRespond(c, "playerId")
	if !ok {
		return
	}

	if err := h.svc.RemoveMember(c.Request.Context(), teamID, playerID, true); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("成员不存在"))
			return
		}
		respondError(c, apierr.BadRequest(err.Error()))
		return
	}

	respondMsg(c, "成员移除成功")
}

// TransferLeaderRequest 转让队长请求
type TransferLeaderRequest struct {
	NewLeaderPlayerID uint64 `json:"newLeaderPlayerId" binding:"required"`
}

// TransferLeader 转让队长
func (h *TeamHandler) TransferLeader(c *gin.Context) {
	teamID, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var req TransferLeaderRequest
	if !ValidateAndRespond(c, &req) {
		return
	}

	// 获取团队当前队长
	team, err := h.svc.GetTeam(c.Request.Context(), teamID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("团队不存在"))
			return
		}
		respondError(c, apierr.InternalError("获取团队失败").WithDetails(err.Error()))
		return
	}

	if err := h.svc.TransferLeader(c.Request.Context(), teamID, team.LeaderID, req.NewLeaderPlayerID); err != nil {
		respondError(c, apierr.BadRequest(err.Error()))
		return
	}

	respondMsg(c, "队长转让成功")
}

// ============================================================================
// 批量操作
// ============================================================================

// BatchDeleteTeamsRequest 批量删除团队请求
type BatchDeleteTeamsRequest struct {
	TeamIDs []uint64 `json:"team_ids" binding:"required,min=1,max=100"`
}

// BatchDeleteTeams 批量删除团队
// @Summary      批量删除团队
// @Description  批量删除多个团队（会检查成员数量和进行中的订单）
// @Tags         Admin/Teams
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BatchDeleteTeamsRequest  true  "团队ID列表"
// @Success      200  {object}  model.APIResponse[BatchOperationResponse]
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/teams/batch/delete [post]
func (h *TeamHandler) BatchDeleteTeams(c *gin.Context) {
	var req BatchDeleteTeamsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}

	if len(req.TeamIDs) == 0 {
		respondAPIError(c, apierr.BadRequest("team_ids is required"))
		return
	}
	if len(req.TeamIDs) > 100 {
		respondAPIError(c, apierr.BadRequest("maximum 100 teams per batch"))
		return
	}

	result, err := h.svc.BatchDeleteTeams(c.Request.Context(), req.TeamIDs)
	if err != nil {
		respondAPIError(c, apierr.InternalError("batch delete teams failed").WithDetails(err.Error()))
		return
	}

	// Convert service result to BatchOperationResponse
	response := convertBatchDeleteResultToResponse(result)
	respondSuccess(c, response)
}

// BatchUpdateTeamsStatusRequest 批量更新团队状态请求
type BatchUpdateTeamsStatusRequest struct {
	TeamIDs []uint64         `json:"team_ids" binding:"required,min=1,max=100"`
	Status  model.TeamStatus `json:"status" binding:"required,oneof=active inactive busy"`
}

// BatchUpdateTeamsStatus 批量更新团队状态
// @Summary      批量更新团队状态
// @Description  批量更新多个团队的状态（active/inactive/busy）
// @Tags         Admin/Teams
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BatchUpdateTeamsStatusRequest  true  "团队ID列表和新状态"
// @Success      200  {object}  model.APIResponse[BatchOperationResponse]
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/teams/batch/status [put]
func (h *TeamHandler) BatchUpdateTeamsStatus(c *gin.Context) {
	var req BatchUpdateTeamsStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}

	if len(req.TeamIDs) == 0 {
		respondAPIError(c, apierr.BadRequest("team_ids is required"))
		return
	}
	if len(req.TeamIDs) > 100 {
		respondAPIError(c, apierr.BadRequest("maximum 100 teams per batch"))
		return
	}
	if req.Status == "" {
		respondAPIError(c, apierr.BadRequest("status is required"))
		return
	}

	// Validate status
	validStatuses := map[model.TeamStatus]bool{
		model.TeamStatusActive:   true,
		model.TeamStatusInactive: true,
		model.TeamStatusBusy:     true,
	}
	if !validStatuses[req.Status] {
		respondAPIError(c, apierr.BadRequest("invalid status, must be one of: active, inactive, busy"))
		return
	}

	result, err := h.svc.BatchUpdateTeamStatus(c.Request.Context(), req.TeamIDs, req.Status)
	if err != nil {
		respondAPIError(c, apierr.InternalError("batch update team status failed").WithDetails(err.Error()))
		return
	}

	// Convert service result to BatchOperationResponse
	response := convertBatchStatusResultToResponse(result)
	respondSuccess(c, response)
}

// BatchAddTeamMembersRequest 批量添加团队成员请求
type BatchAddTeamMembersRequest struct {
	TeamID    uint64   `json:"team_id" binding:"required"`
	PlayerIDs []uint64 `json:"player_ids" binding:"required,min=1,max=100"`
}

// BatchAddTeamMembers 批量添加团队成员
// @Summary      批量添加团队成员
// @Description  批量添加多个陪玩师到指定团队
// @Tags         Admin/Teams
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BatchAddTeamMembersRequest  true  "团队ID和陪玩师ID列表"
// @Success      200  {object}  model.APIResponse[BatchOperationResponse]
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/teams/batch/members [post]
func (h *TeamHandler) BatchAddTeamMembers(c *gin.Context) {
	var req BatchAddTeamMembersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}

	if req.TeamID == 0 {
		respondAPIError(c, apierr.BadRequest("team_id is required"))
		return
	}
	if len(req.PlayerIDs) == 0 {
		respondAPIError(c, apierr.BadRequest("player_ids is required"))
		return
	}
	if len(req.PlayerIDs) > 100 {
		respondAPIError(c, apierr.BadRequest("maximum 100 players per batch"))
		return
	}

	result, err := h.svc.BatchAddMembers(c.Request.Context(), req.TeamID, req.PlayerIDs)
	if err != nil {
		respondAPIError(c, apierr.InternalError("batch add team members failed").WithDetails(err.Error()))
		return
	}

	// Convert service result to BatchOperationResponse
	response := convertBatchAddMembersResultToResponse(result)
	respondSuccess(c, response)
}

// Helper functions to convert service results to unified response format

func convertBatchDeleteResultToResponse(result *teamservice.BatchDeleteTeamsResult) BatchOperationResponse {
	response := BatchOperationResponse{
		SuccessCount: result.SuccessCount,
		FailedCount:  result.FailedCount,
		TotalCount:   result.SuccessCount + result.FailedCount,
		FailedItems:  make([]BatchOperationError, 0, len(result.FailedIDs)),
		SuccessItems: make([]uint64, 0, result.SuccessCount),
	}

	for i, id := range result.FailedIDs {
		response.FailedItems = append(response.FailedItems, BatchOperationError{
			ID:      id,
			Message: result.Errors[i],
		})
	}

	return response
}

func convertBatchStatusResultToResponse(result *teamservice.BatchUpdateTeamStatusResult) BatchOperationResponse {
	response := BatchOperationResponse{
		SuccessCount: result.SuccessCount,
		FailedCount:  result.FailedCount,
		TotalCount:   result.SuccessCount + result.FailedCount,
		FailedItems:  make([]BatchOperationError, 0, len(result.FailedIDs)),
		SuccessItems: make([]uint64, 0, result.SuccessCount),
	}

	for i, id := range result.FailedIDs {
		response.FailedItems = append(response.FailedItems, BatchOperationError{
			ID:      id,
			Message: result.Errors[i],
		})
	}

	return response
}

func convertBatchAddMembersResultToResponse(result *teamservice.BatchAddMembersResult) BatchOperationResponse {
	response := BatchOperationResponse{
		SuccessCount: result.SuccessCount,
		FailedCount:  result.FailedCount,
		TotalCount:   result.SuccessCount + result.FailedCount,
		FailedItems:  make([]BatchOperationError, 0, len(result.FailedPlayerIDs)),
		SuccessItems: make([]uint64, 0, result.SuccessCount),
	}

	for i, id := range result.FailedPlayerIDs {
		response.FailedItems = append(response.FailedItems, BatchOperationError{
			ID:      id,
			Message: result.Errors[i],
		})
	}

	return response
}

// ListMembers 获取成员列表（分页）
func (h *TeamHandler) ListMembers(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	var teamID *uint64
	if v := strings.TrimSpace(c.Query("teamId")); v != "" {
		id, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			respondBadRequest(c, "无效的团队ID")
			return
		}
		teamID = &id
	}

	var playerID *uint64
	if v := strings.TrimSpace(c.Query("playerId")); v != "" {
		id, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			respondBadRequest(c, "无效的陪玩师ID")
			return
		}
		playerID = &id
	}

	var role *model.TeamMemberRole
	if v := strings.TrimSpace(c.Query("role")); v != "" {
		r := model.TeamMemberRole(v)
		role = &r
	}

	var status *model.TeamMemberStatus
	if v := strings.TrimSpace(c.Query("status")); v != "" {
		s := model.TeamMemberStatus(v)
		status = &s
	}

	opts := teamrepo.MemberListOptions{
		Page:     page,
		PageSize: pageSize,
		TeamID:   teamID,
		PlayerID: playerID,
		Role:     role,
		Status:   status,
	}

	members, total, err := h.svc.ListMembers(c.Request.Context(), opts)
	if err != nil {
		respondError(c, apierr.InternalError("获取成员列表失败").WithDetails(err.Error()))
		return
	}

	respondList(c, members, &model.Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
	})
}

// ============================================================================
// 邀请管理
// ============================================================================

// ListInvites 获取邀请列表
func (h *TeamHandler) ListInvites(c *gin.Context) {
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}

	var teamID *uint64
	if v := strings.TrimSpace(c.Query("teamId")); v != "" {
		id, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			respondBadRequest(c, "无效的团队ID")
			return
		}
		teamID = &id
	}

	var playerID *uint64
	if v := strings.TrimSpace(c.Query("playerId")); v != "" {
		id, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			respondBadRequest(c, "无效的陪玩师ID")
			return
		}
		playerID = &id
	}

	var status *model.TeamInviteStatus
	if v := strings.TrimSpace(c.Query("status")); v != "" {
		s := model.TeamInviteStatus(v)
		status = &s
	}

	opts := teamrepo.InviteListOptions{
		Page:     page,
		PageSize: pageSize,
		TeamID:   teamID,
		PlayerID: playerID,
		Status:   status,
	}

	invites, total, err := h.svc.ListInvites(c.Request.Context(), opts)
	if err != nil {
		respondError(c, apierr.InternalError("获取邀请列表失败").WithDetails(err.Error()))
		return
	}

	respondList(c, invites, &model.Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
	})
}

// GetInvite 获取邀请详情
func (h *TeamHandler) GetInvite(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "inviteId")
	if !ok {
		return
	}

	invite, err := h.svc.GetInvite(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(c, apierr.NotFound("邀请不存在"))
			return
		}
		respondError(c, apierr.InternalError("获取邀请失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, invite)
}

// ============================================================================
// 路由注册
// ============================================================================

// RegisterTeamRoutes 注册团队管理路由
func RegisterTeamRoutes(rg *gin.RouterGroup, svc *teamservice.TeamService, pm *middleware.PermissionMiddleware) {
	h := NewTeamHandler(svc)

	teamGroup := rg.Group("/teams")
	teamGroup.Use(pm.RequireAuth())
	{
		// 团队管理
		teamGroup.GET("", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/teams"), h.ListTeams)
		teamGroup.GET("/stats", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/teams/stats"), h.GetTeamStats)
		teamGroup.GET("/:id", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/teams/:id"), h.GetTeam)
		teamGroup.POST("", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/teams"), h.CreateTeam)
		teamGroup.PUT("/:id", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/teams/:id"), h.UpdateTeam)
		teamGroup.DELETE("/:id", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/teams/:id"), h.DeleteTeam)
		teamGroup.PUT("/:id/status", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/teams/:id/status"), h.UpdateTeamStatus)

		// 批量操作
		teamGroup.DELETE("/batch", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/teams/batch"), h.BatchDeleteTeams)
		teamGroup.PUT("/batch/status", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/teams/batch/status"), h.BatchUpdateTeamsStatus)
		teamGroup.POST("/batch/members", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/teams/batch/members"), h.BatchAddTeamMembers)

		// 成员管理
		teamGroup.GET("/:id/members", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/teams/:id/members"), h.GetTeamMembers)
		teamGroup.POST("/:id/members", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/teams/:id/members"), h.AddMember)
		teamGroup.DELETE("/:id/members/:playerId", pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/teams/:id/members/:playerId"), h.RemoveMember)
		teamGroup.POST("/:id/transfer-leader", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/teams/:id/transfer-leader"), h.TransferLeader)

		// 成员列表（全局）
		teamGroup.GET("/members", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/teams/members"), h.ListMembers)

		// 邀请管理
		teamGroup.GET("/invites", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/teams/invites"), h.ListInvites)
		teamGroup.GET("/invites/:inviteId", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/teams/invites/:inviteId"), h.GetInvite)
	}
}
