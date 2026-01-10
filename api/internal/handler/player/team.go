package player

import (
	"errors"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	teamrepo "gamelink/internal/repository/team"
	teamservice "gamelink/internal/service/team"
	"gamelink/pkg/apierr"
)

// TeamHandler 陪玩师团队处理器
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
	Name            string `json:"name" binding:"required,max=64"`
	Description     string `json:"description" binding:"max=255"`
	AvatarURL       string `json:"avatarUrl"`
	MaxMembers      int    `json:"maxMembers"`
	IncomeShareType string `json:"incomeShareType"`
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

// GetMyTeam 获取我的团队
func (h *TeamHandler) GetMyTeam(c *gin.Context) {
	playerID := resp.GetUserID(c)
	if playerID == 0 {
		resp.Error(c, apierr.Unauthorized("未登录"))
		return
	}

	team, err := h.svc.GetPlayerTeam(c.Request.Context(), playerID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			resp.OK(c, (*model.Team)(nil)) // 未加入团队
			return
		}
		resp.Error(c, apierr.InternalError("获取团队失败").WithDetails(err.Error()))
		return
	}

	resp.OK(c, team)
}

// CreateTeam 创建团队
func (h *TeamHandler) CreateTeam(c *gin.Context) {
	playerID := resp.GetUserID(c)
	if playerID == 0 {
		resp.Error(c, apierr.Unauthorized("未登录"))
		return
	}

	var req TeamCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, apierr.BadRequest("参数错误").WithDetails(err.Error()))
		return
	}

	team := &model.Team{
		Name:            req.Name,
		Description:     req.Description,
		AvatarURL:       req.AvatarURL,
		MaxMembers:      req.MaxMembers,
		IncomeShareType: req.IncomeShareType,
	}

	if err := h.svc.CreateTeam(c.Request.Context(), team, playerID); err != nil {
		resp.Error(c, apierr.BadRequest(err.Error()))
		return
	}

	resp.Created(c, team)
}

// UpdateTeam 更新团队（仅队长）
func (h *TeamHandler) UpdateTeam(c *gin.Context) {
	playerID := resp.GetUserID(c)
	if playerID == 0 {
		resp.Error(c, apierr.Unauthorized("未登录"))
		return
	}

	// 获取我的团队
	team, err := h.svc.GetPlayerTeam(c.Request.Context(), playerID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			resp.Error(c, apierr.NotFound("您未加入任何团队"))
			return
		}
		resp.Error(c, apierr.InternalError("获取团队失败").WithDetails(err.Error()))
		return
	}

	// 验证是队长
	if team.LeaderID != playerID {
		resp.Error(c, apierr.Forbidden("只有队长可以修改团队信息"))
		return
	}

	var req TeamUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, apierr.BadRequest("参数错误").WithDetails(err.Error()))
		return
	}

	team.Name = req.Name
	team.Description = req.Description
	team.AvatarURL = req.AvatarURL
	team.MaxMembers = req.MaxMembers
	team.IncomeShareType = req.IncomeShareType
	team.LeaderBonusRate = req.LeaderBonusRate

	if err := h.svc.UpdateTeam(c.Request.Context(), team); err != nil {
		resp.Error(c, apierr.BadRequest(err.Error()))
		return
	}

	resp.Updated(c, team)
}

// LeaveTeam 离开团队
func (h *TeamHandler) LeaveTeam(c *gin.Context) {
	playerID := resp.GetUserID(c)
	if playerID == 0 {
		resp.Error(c, apierr.Unauthorized("未登录"))
		return
	}

	// 获取我的团队
	team, err := h.svc.GetPlayerTeam(c.Request.Context(), playerID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			resp.Error(c, apierr.NotFound("您未加入任何团队"))
			return
		}
		resp.Error(c, apierr.InternalError("获取团队失败").WithDetails(err.Error()))
		return
	}

	if err := h.svc.RemoveMember(c.Request.Context(), team.ID, playerID, false); err != nil {
		resp.Error(c, apierr.BadRequest(err.Error()))
		return
	}

	resp.OK(c, gin.H{"message": "已离开团队"})
}

// GetTeamMembers 获取团队成员
func (h *TeamHandler) GetTeamMembers(c *gin.Context) {
	playerID := resp.GetUserID(c)
	if playerID == 0 {
		resp.Error(c, apierr.Unauthorized("未登录"))
		return
	}

	// 获取我的团队
	team, err := h.svc.GetPlayerTeam(c.Request.Context(), playerID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			resp.Error(c, apierr.NotFound("您未加入任何团队"))
			return
		}
		resp.Error(c, apierr.InternalError("获取团队失败").WithDetails(err.Error()))
		return
	}

	members, err := h.svc.GetTeamMembers(c.Request.Context(), team.ID)
	if err != nil {
		resp.Error(c, apierr.InternalError("获取成员列表失败").WithDetails(err.Error()))
		return
	}

	resp.OK(c, members)
}

// ============================================================================
// 成员管理（队长操作）
// ============================================================================

// KickMemberRequest 踢出成员请求
type KickMemberRequest struct {
	PlayerID uint64 `json:"playerId" binding:"required"`
}

// KickMember 踢出成员（仅队长）
func (h *TeamHandler) KickMember(c *gin.Context) {
	playerID := resp.GetUserID(c)
	if playerID == 0 {
		resp.Error(c, apierr.Unauthorized("未登录"))
		return
	}

	// 获取我的团队
	team, err := h.svc.GetPlayerTeam(c.Request.Context(), playerID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			resp.Error(c, apierr.NotFound("您未加入任何团队"))
			return
		}
		resp.Error(c, apierr.InternalError("获取团队失败").WithDetails(err.Error()))
		return
	}

	var req KickMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, apierr.BadRequest("参数错误").WithDetails(err.Error()))
		return
	}

	if err := h.svc.KickMember(c.Request.Context(), team.ID, playerID, req.PlayerID); err != nil {
		resp.Error(c, apierr.BadRequest(err.Error()))
		return
	}

	resp.OK(c, gin.H{"message": "成员已踢出"})
}

// TransferLeaderRequest 转让队长请求
type TransferLeaderRequest struct {
	NewLeaderPlayerID uint64 `json:"newLeaderPlayerId" binding:"required"`
}

// TransferLeader 转让队长（仅队长）
func (h *TeamHandler) TransferLeader(c *gin.Context) {
	playerID := resp.GetUserID(c)
	if playerID == 0 {
		resp.Error(c, apierr.Unauthorized("未登录"))
		return
	}

	// 获取我的团队
	team, err := h.svc.GetPlayerTeam(c.Request.Context(), playerID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			resp.Error(c, apierr.NotFound("您未加入任何团队"))
			return
		}
		resp.Error(c, apierr.InternalError("获取团队失败").WithDetails(err.Error()))
		return
	}

	var req TransferLeaderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, apierr.BadRequest("参数错误").WithDetails(err.Error()))
		return
	}

	if err := h.svc.TransferLeader(c.Request.Context(), team.ID, playerID, req.NewLeaderPlayerID); err != nil {
		resp.Error(c, apierr.BadRequest(err.Error()))
		return
	}

	resp.OK(c, gin.H{"message": "队长已转让"})
}

// ============================================================================
// 邀请管理
// ============================================================================

// InviteRequest 邀请请求
type InviteRequest struct {
	PlayerID uint64 `json:"playerId" binding:"required"`
	Message  string `json:"message"`
}

// InviteMember 邀请成员
func (h *TeamHandler) InviteMember(c *gin.Context) {
	playerID := resp.GetUserID(c)
	if playerID == 0 {
		resp.Error(c, apierr.Unauthorized("未登录"))
		return
	}

	// 获取我的团队
	team, err := h.svc.GetPlayerTeam(c.Request.Context(), playerID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			resp.Error(c, apierr.NotFound("您未加入任何团队"))
			return
		}
		resp.Error(c, apierr.InternalError("获取团队失败").WithDetails(err.Error()))
		return
	}

	var req InviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, apierr.BadRequest("参数错误").WithDetails(err.Error()))
		return
	}

	if err := h.svc.CreateInvite(c.Request.Context(), team.ID, playerID, req.PlayerID, req.Message); err != nil {
		resp.Error(c, apierr.BadRequest(err.Error()))
		return
	}

	resp.OK(c, gin.H{"message": "邀请已发送"})
}

// GetMyInvites 获取我收到的邀请
func (h *TeamHandler) GetMyInvites(c *gin.Context) {
	playerID := resp.GetUserID(c)
	if playerID == 0 {
		resp.Error(c, apierr.Unauthorized("未登录"))
		return
	}

	page, _ := resp.ParseUintParam(c, "page")
	if page == 0 {
		page = 1
	}
	pageSize, _ := resp.ParseUintParam(c, "pageSize")
	if pageSize == 0 {
		pageSize = 20
	}

	status := model.TeamInviteStatusPending
	invites, total, err := h.svc.ListInvites(c.Request.Context(), teamrepo.InviteListOptions{
		Page:     int(page),
		PageSize: int(pageSize),
		PlayerID: &playerID,
		Status:   &status,
	})
	if err != nil {
		resp.Error(c, apierr.InternalError("获取邀请列表失败").WithDetails(err.Error()))
		return
	}

	resp.List(c, invites, &model.Pagination{
		Page:     int(page),
		PageSize: int(pageSize),
		Total:    int(total),
	})
}

// AcceptInvite 接受邀请
func (h *TeamHandler) AcceptInvite(c *gin.Context) {
	playerID := resp.GetUserID(c)
	if playerID == 0 {
		resp.Error(c, apierr.Unauthorized("未登录"))
		return
	}

	inviteID, err := resp.ParseUintParam(c, "id")
	if err != nil {
		resp.Error(c, apierr.BadRequest("无效的邀请ID"))
		return
	}

	if err := h.svc.AcceptInvite(c.Request.Context(), inviteID, playerID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			resp.Error(c, apierr.NotFound("邀请不存在"))
			return
		}
		resp.Error(c, apierr.BadRequest(err.Error()))
		return
	}

	resp.OK(c, gin.H{"message": "已加入团队"})
}

// RejectInvite 拒绝邀请
func (h *TeamHandler) RejectInvite(c *gin.Context) {
	playerID := resp.GetUserID(c)
	if playerID == 0 {
		resp.Error(c, apierr.Unauthorized("未登录"))
		return
	}

	inviteID, err := resp.ParseUintParam(c, "id")
	if err != nil {
		resp.Error(c, apierr.BadRequest("无效的邀请ID"))
		return
	}

	if err := h.svc.RejectInvite(c.Request.Context(), inviteID, playerID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			resp.Error(c, apierr.NotFound("邀请不存在"))
			return
		}
		resp.Error(c, apierr.BadRequest(err.Error()))
		return
	}

	resp.OK(c, gin.H{"message": "已拒绝邀请"})
}

// ============================================================================
// 路由注册
// ============================================================================

// RegisterTeamRoutes 注册陪玩师团队路由
func RegisterTeamRoutes(rg *gin.RouterGroup, svc *teamservice.TeamService, authMiddleware gin.HandlerFunc) {
	h := NewTeamHandler(svc)

	teamGroup := rg.Group("/team")
	teamGroup.Use(authMiddleware)
	{
		// 我的团队
		teamGroup.GET("", h.GetMyTeam)
		teamGroup.POST("", h.CreateTeam)
		teamGroup.PUT("", h.UpdateTeam)
		teamGroup.DELETE("", h.LeaveTeam)

		// 成员管理
		teamGroup.GET("/members", h.GetTeamMembers)
		teamGroup.POST("/kick", h.KickMember)
		teamGroup.POST("/transfer-leader", h.TransferLeader)

		// 邀请管理
		teamGroup.POST("/invite", h.InviteMember)
		teamGroup.GET("/invites", h.GetMyInvites)
		teamGroup.POST("/invites/:id/accept", h.AcceptInvite)
		teamGroup.POST("/invites/:id/reject", h.RejectInvite)
	}
}
