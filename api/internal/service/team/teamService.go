package team

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	teamrepo "gamelink/internal/repository/team"
)

// TeamRepository defines the interface for team repository operations
type TeamRepository interface {
	Create(ctx context.Context, team *model.Team) error
	GetByID(ctx context.Context, id uint64) (*model.Team, error)
	Update(ctx context.Context, team *model.Team) error
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, opts teamrepo.TeamListOptions) ([]model.Team, int64, error)
	UpdateStatus(ctx context.Context, id uint64, status model.TeamStatus) error
	UpdateMemberCount(ctx context.Context, id uint64, delta int) error
	UpdateLeader(ctx context.Context, id uint64, leaderID uint64) error
	GetTeamStats(ctx context.Context) (*teamrepo.TeamStats, error)
	CreateMember(ctx context.Context, member *model.TeamMember) error
	GetMemberByTeamAndPlayer(ctx context.Context, teamID, playerID uint64) (*model.TeamMember, error)
	GetActiveMemberByPlayer(ctx context.Context, playerID uint64) (*model.TeamMember, error)
	UpdateMember(ctx context.Context, member *model.TeamMember) error
	GetActiveMembers(ctx context.Context, teamID uint64) ([]model.TeamMember, error)
	ListMembers(ctx context.Context, opts teamrepo.MemberListOptions) ([]model.TeamMember, int64, error)
	GetNextLeader(ctx context.Context, teamID uint64, excludePlayerID uint64) (*model.TeamMember, error)
	CreateInvite(ctx context.Context, invite *model.TeamInvite) error
	GetInviteByID(ctx context.Context, id uint64) (*model.TeamInvite, error)
	GetPendingInvite(ctx context.Context, teamID, playerID uint64) (*model.TeamInvite, error)
	UpdateInvite(ctx context.Context, invite *model.TeamInvite) error
	ListInvites(ctx context.Context, opts teamrepo.InviteListOptions) ([]model.TeamInvite, int64, error)
}

// TeamService 团队服务
type TeamService struct {
	repo TeamRepository
}

// NewTeamService 创建团队服务
func NewTeamService(repo TeamRepository) *TeamService {
	return &TeamService{repo: repo}
}

// ============================================================================
// 团队管理
// ============================================================================

// CreateTeam 创建团队
func (s *TeamService) CreateTeam(ctx context.Context, team *model.Team, leaderPlayerID uint64) error {
	// 检查陪玩师是否已在其他团队
	_, err := s.repo.GetActiveMemberByPlayer(ctx, leaderPlayerID)
	if err == nil {
		return errors.New("陪玩师已在其他团队中")
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return fmt.Errorf("检查陪玩师团队状态失败: %w", err)
	}

	// 设置默认值
	team.LeaderID = leaderPlayerID
	team.Status = model.TeamStatusActive
	team.MemberCount = 1
	if team.MaxMembers == 0 {
		team.MaxMembers = 5
	}
	if team.IncomeShareType == "" {
		team.IncomeShareType = "equal"
	}

	// 创建团队
	if err := s.repo.Create(ctx, team); err != nil {
		return fmt.Errorf("创建团队失败: %w", err)
	}

	// 创建队长成员记录
	member := &model.TeamMember{
		TeamID:    team.ID,
		PlayerID:  leaderPlayerID,
		Role:      model.TeamMemberRoleLeader,
		Status:    model.TeamMemberStatusActive,
		JoinedAt:  time.Now(),
		SortOrder: 0,
	}
	if err := s.repo.CreateMember(ctx, member); err != nil {
		return fmt.Errorf("创建队长成员记录失败: %w", err)
	}

	return nil
}

// GetTeam 获取团队详情
func (s *TeamService) GetTeam(ctx context.Context, id uint64) (*model.Team, error) {
	return s.repo.GetByID(ctx, id)
}

// UpdateTeam 更新团队
func (s *TeamService) UpdateTeam(ctx context.Context, team *model.Team) error {
	// 检查团队是否存在
	existing, err := s.repo.GetByID(ctx, team.ID)
	if err != nil {
		return err
	}

	// 不能减少最大成员数到当前成员数以下
	if team.MaxMembers < existing.MemberCount {
		return errors.New("最大成员数不能小于当前成员数")
	}

	return s.repo.Update(ctx, team)
}

// DeleteTeam 删除团队
func (s *TeamService) DeleteTeam(ctx context.Context, id uint64) error {
	team, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 检查是否有进行中的订单
	if team.CurrentOrderID != nil {
		return errors.New("团队有进行中的订单，无法删除")
	}

	return s.repo.Delete(ctx, id)
}

// ListTeams 获取团队列表
func (s *TeamService) ListTeams(ctx context.Context, opts teamrepo.TeamListOptions) ([]model.Team, int64, error) {
	return s.repo.List(ctx, opts)
}

// UpdateTeamStatus 更新团队状态
func (s *TeamService) UpdateTeamStatus(ctx context.Context, id uint64, status model.TeamStatus) error {
	team, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 如果有进行中的订单，不能设置为 inactive
	if status == model.TeamStatusInactive && team.CurrentOrderID != nil {
		return errors.New("团队有进行中的订单，无法停用")
	}

	return s.repo.UpdateStatus(ctx, id, status)
}

// GetTeamStats 获取团队统计
func (s *TeamService) GetTeamStats(ctx context.Context) (*teamrepo.TeamStats, error) {
	return s.repo.GetTeamStats(ctx)
}

// ============================================================================
// 成员管理
// ============================================================================

// AddMember 添加成员
func (s *TeamService) AddMember(ctx context.Context, teamID, playerID uint64) error {
	// 获取团队
	team, err := s.repo.GetByID(ctx, teamID)
	if err != nil {
		return err
	}

	// 检查团队是否已满
	if team.MemberCount >= team.MaxMembers {
		return errors.New("团队已满")
	}

	// 检查陪玩师是否已在其他团队
	_, err = s.repo.GetActiveMemberByPlayer(ctx, playerID)
	if err == nil {
		return errors.New("陪玩师已在其他团队中")
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return fmt.Errorf("检查陪玩师团队状态失败: %w", err)
	}

	// 检查是否曾经是该团队成员
	existingMember, err := s.repo.GetMemberByTeamAndPlayer(ctx, teamID, playerID)
	if err == nil {
		// 重新激活
		existingMember.Status = model.TeamMemberStatusActive
		existingMember.LeftAt = nil
		existingMember.JoinedAt = time.Now()
		if err := s.repo.UpdateMember(ctx, existingMember); err != nil {
			return fmt.Errorf("重新激活成员失败: %w", err)
		}
	} else if errors.Is(err, repository.ErrNotFound) {
		// 创建新成员
		member := &model.TeamMember{
			TeamID:    teamID,
			PlayerID:  playerID,
			Role:      model.TeamMemberRoleMember,
			Status:    model.TeamMemberStatusActive,
			JoinedAt:  time.Now(),
			SortOrder: team.MemberCount, // 排在最后
		}
		if err := s.repo.CreateMember(ctx, member); err != nil {
			return fmt.Errorf("创建成员失败: %w", err)
		}
	} else {
		return fmt.Errorf("检查成员记录失败: %w", err)
	}

	// 更新团队成员数
	return s.repo.UpdateMemberCount(ctx, teamID, 1)
}

// RemoveMember 移除成员
func (s *TeamService) RemoveMember(ctx context.Context, teamID, playerID uint64, kicked bool) error {
	// 获取团队
	team, err := s.repo.GetByID(ctx, teamID)
	if err != nil {
		return err
	}

	// 检查是否有进行中的订单
	if team.CurrentOrderID != nil {
		return errors.New("团队有进行中的订单，无法移除成员")
	}

	// 获取成员
	member, err := s.repo.GetMemberByTeamAndPlayer(ctx, teamID, playerID)
	if err != nil {
		return err
	}

	if member.Status != model.TeamMemberStatusActive {
		return errors.New("成员已不在团队中")
	}

	// 如果是队长离队，需要转让队长
	if member.Role == model.TeamMemberRoleLeader {
		// 查找下一个队长
		nextLeader, err := s.repo.GetNextLeader(ctx, teamID, playerID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				// 没有其他成员，删除团队
				return s.repo.Delete(ctx, teamID)
			}
			return fmt.Errorf("查找下一个队长失败: %w", err)
		}

		// 转让队长
		nextLeader.Role = model.TeamMemberRoleLeader
		if err := s.repo.UpdateMember(ctx, nextLeader); err != nil {
			return fmt.Errorf("转让队长失败: %w", err)
		}
		if err := s.repo.UpdateLeader(ctx, teamID, nextLeader.PlayerID); err != nil {
			return fmt.Errorf("更新团队队长失败: %w", err)
		}
	}

	// 更新成员状态
	now := time.Now()
	member.LeftAt = &now
	if kicked {
		member.Status = model.TeamMemberStatusKicked
	} else {
		member.Status = model.TeamMemberStatusLeft
	}
	if err := s.repo.UpdateMember(ctx, member); err != nil {
		return fmt.Errorf("更新成员状态失败: %w", err)
	}

	// 更新团队成员数
	return s.repo.UpdateMemberCount(ctx, teamID, -1)
}

// KickMember 踢出成员（队长操作）
func (s *TeamService) KickMember(ctx context.Context, teamID, leaderPlayerID, targetPlayerID uint64) error {
	// 验证操作者是队长
	leader, err := s.repo.GetMemberByTeamAndPlayer(ctx, teamID, leaderPlayerID)
	if err != nil {
		return fmt.Errorf("获取队长信息失败: %w", err)
	}
	if leader.Role != model.TeamMemberRoleLeader {
		return errors.New("只有队长可以踢出成员")
	}

	// 不能踢出自己
	if leaderPlayerID == targetPlayerID {
		return errors.New("不能踢出自己")
	}

	return s.RemoveMember(ctx, teamID, targetPlayerID, true)
}

// TransferLeader 转让队长
func (s *TeamService) TransferLeader(ctx context.Context, teamID, currentLeaderID, newLeaderID uint64) error {
	// 获取团队
	team, err := s.repo.GetByID(ctx, teamID)
	if err != nil {
		return err
	}

	// 验证当前队长
	if team.LeaderID != currentLeaderID {
		return errors.New("只有队长可以转让队长")
	}

	// 获取当前队长成员记录
	currentLeader, err := s.repo.GetMemberByTeamAndPlayer(ctx, teamID, currentLeaderID)
	if err != nil {
		return fmt.Errorf("获取当前队长信息失败: %w", err)
	}

	// 获取新队长成员记录
	newLeader, err := s.repo.GetMemberByTeamAndPlayer(ctx, teamID, newLeaderID)
	if err != nil {
		return fmt.Errorf("获取新队长信息失败: %w", err)
	}

	if newLeader.Status != model.TeamMemberStatusActive {
		return errors.New("目标成员不在团队中")
	}

	// 更新角色
	currentLeader.Role = model.TeamMemberRoleMember
	newLeader.Role = model.TeamMemberRoleLeader

	if err := s.repo.UpdateMember(ctx, currentLeader); err != nil {
		return fmt.Errorf("更新原队长角色失败: %w", err)
	}
	if err := s.repo.UpdateMember(ctx, newLeader); err != nil {
		return fmt.Errorf("更新新队长角色失败: %w", err)
	}

	// 更新团队队长
	return s.repo.UpdateLeader(ctx, teamID, newLeaderID)
}

// GetTeamMembers 获取团队成员列表
func (s *TeamService) GetTeamMembers(ctx context.Context, teamID uint64) ([]model.TeamMember, error) {
	return s.repo.GetActiveMembers(ctx, teamID)
}

// ListMembers 获取成员列表（分页）
func (s *TeamService) ListMembers(ctx context.Context, opts teamrepo.MemberListOptions) ([]model.TeamMember, int64, error) {
	return s.repo.ListMembers(ctx, opts)
}

// GetPlayerTeam 获取陪玩师所在团队
func (s *TeamService) GetPlayerTeam(ctx context.Context, playerID uint64) (*model.Team, error) {
	member, err := s.repo.GetActiveMemberByPlayer(ctx, playerID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, member.TeamID)
}

// ============================================================================
// 邀请管理
// ============================================================================

// CreateInvite 创建邀请
func (s *TeamService) CreateInvite(ctx context.Context, teamID, inviterPlayerID, targetPlayerID uint64, message string) error {
	// 获取团队
	team, err := s.repo.GetByID(ctx, teamID)
	if err != nil {
		return err
	}

	// 验证邀请者是团队成员
	inviter, err := s.repo.GetMemberByTeamAndPlayer(ctx, teamID, inviterPlayerID)
	if err != nil {
		return fmt.Errorf("邀请者不是团队成员: %w", err)
	}
	if inviter.Status != model.TeamMemberStatusActive {
		return errors.New("邀请者不在团队中")
	}

	// 检查团队是否已满
	if team.MemberCount >= team.MaxMembers {
		return errors.New("团队已满")
	}

	// 检查目标陪玩师是否已在其他团队
	_, err = s.repo.GetActiveMemberByPlayer(ctx, targetPlayerID)
	if err == nil {
		return errors.New("目标陪玩师已在其他团队中")
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return fmt.Errorf("检查目标陪玩师团队状态失败: %w", err)
	}

	// 检查是否已有待处理的邀请
	_, err = s.repo.GetPendingInvite(ctx, teamID, targetPlayerID)
	if err == nil {
		return errors.New("已有待处理的邀请")
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return fmt.Errorf("检查邀请状态失败: %w", err)
	}

	// 创建邀请
	invite := &model.TeamInvite{
		TeamID:    teamID,
		PlayerID:  targetPlayerID,
		InviterID: inviterPlayerID,
		Status:    model.TeamInviteStatusPending,
		ExpireAt:  time.Now().Add(7 * 24 * time.Hour), // 7天有效期
		Message:   message,
	}

	return s.repo.CreateInvite(ctx, invite)
}

// AcceptInvite 接受邀请
func (s *TeamService) AcceptInvite(ctx context.Context, inviteID, playerID uint64) error {
	invite, err := s.repo.GetInviteByID(ctx, inviteID)
	if err != nil {
		return err
	}

	// 验证是被邀请人
	if invite.PlayerID != playerID {
		return errors.New("无权操作此邀请")
	}

	// 检查邀请状态
	if invite.Status != model.TeamInviteStatusPending {
		return errors.New("邀请已处理或已过期")
	}

	// 检查是否过期
	if time.Now().After(invite.ExpireAt) {
		invite.Status = model.TeamInviteStatusExpired
		_ = s.repo.UpdateInvite(ctx, invite)
		return errors.New("邀请已过期")
	}

	// 添加成员
	if err := s.AddMember(ctx, invite.TeamID, playerID); err != nil {
		return err
	}

	// 更新邀请状态
	invite.Status = model.TeamInviteStatusAccepted
	return s.repo.UpdateInvite(ctx, invite)
}

// RejectInvite 拒绝邀请
func (s *TeamService) RejectInvite(ctx context.Context, inviteID, playerID uint64) error {
	invite, err := s.repo.GetInviteByID(ctx, inviteID)
	if err != nil {
		return err
	}

	// 验证是被邀请人
	if invite.PlayerID != playerID {
		return errors.New("无权操作此邀请")
	}

	// 检查邀请状态
	if invite.Status != model.TeamInviteStatusPending {
		return errors.New("邀请已处理或已过期")
	}

	invite.Status = model.TeamInviteStatusRejected
	return s.repo.UpdateInvite(ctx, invite)
}

// ListInvites 获取邀请列表
func (s *TeamService) ListInvites(ctx context.Context, opts teamrepo.InviteListOptions) ([]model.TeamInvite, int64, error) {
	return s.repo.ListInvites(ctx, opts)
}

// GetInvite 获取邀请详情
func (s *TeamService) GetInvite(ctx context.Context, id uint64) (*model.TeamInvite, error) {
	return s.repo.GetInviteByID(ctx, id)
}

// ============================================================================
// 批量操作
// ============================================================================

// BatchDeleteTeamsResult 批量删除团队结果
type BatchDeleteTeamsResult struct {
	SuccessCount int      `json:"successCount"`
	FailedCount  int      `json:"failedCount"`
	FailedIDs    []uint64 `json:"failedIds"`
	Errors       []string `json:"errors"`
}

// BatchDeleteTeams 批量删除团队
func (s *TeamService) BatchDeleteTeams(ctx context.Context, ids []uint64) (*BatchDeleteTeamsResult, error) {
	if len(ids) == 0 {
		return nil, errors.New("团队ID列表不能为空")
	}

	result := &BatchDeleteTeamsResult{
		FailedIDs: make([]uint64, 0),
		Errors:    make([]string, 0),
	}

	for _, id := range ids {
		err := s.DeleteTeam(ctx, id)
		if err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, id)
			result.Errors = append(result.Errors, fmt.Sprintf("团队%d: %s", id, err.Error()))
		} else {
			result.SuccessCount++
		}
	}

	return result, nil
}

// BatchUpdateTeamStatusResult 批量更新团队状态结果
type BatchUpdateTeamStatusResult struct {
	SuccessCount int      `json:"successCount"`
	FailedCount  int      `json:"failedCount"`
	FailedIDs    []uint64 `json:"failedIds"`
	Errors       []string `json:"errors"`
}

// BatchUpdateTeamStatus 批量更新团队状态
func (s *TeamService) BatchUpdateTeamStatus(ctx context.Context, ids []uint64, status model.TeamStatus) (*BatchUpdateTeamStatusResult, error) {
	if len(ids) == 0 {
		return nil, errors.New("团队ID列表不能为空")
	}

	// 验证状态值
	validStatuses := map[model.TeamStatus]bool{
		model.TeamStatusActive:   true,
		model.TeamStatusInactive: true,
		model.TeamStatusBusy:     true,
	}
	if !validStatuses[status] {
		return nil, errors.New("无效的团队状态")
	}

	result := &BatchUpdateTeamStatusResult{
		FailedIDs: make([]uint64, 0),
		Errors:    make([]string, 0),
	}

	for _, id := range ids {
		err := s.UpdateTeamStatus(ctx, id, status)
		if err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, id)
			result.Errors = append(result.Errors, fmt.Sprintf("团队%d: %s", id, err.Error()))
		} else {
			result.SuccessCount++
		}
	}

	return result, nil
}

// BatchAddMembersResult 批量添加成员结果
type BatchAddMembersResult struct {
	SuccessCount  int      `json:"successCount"`
	FailedCount   int      `json:"failedCount"`
	FailedPlayerIDs []uint64 `json:"failedPlayerIds"`
	Errors        []string `json:"errors"`
}

// BatchAddMembers 批量添加成员到团队
func (s *TeamService) BatchAddMembers(ctx context.Context, teamID uint64, playerIDs []uint64) (*BatchAddMembersResult, error) {
	if len(playerIDs) == 0 {
		return nil, errors.New("陪玩师ID列表不能为空")
	}

	// 检查团队是否存在
	_, err := s.repo.GetByID(ctx, teamID)
	if err != nil {
		return nil, err
	}

	result := &BatchAddMembersResult{
		FailedPlayerIDs: make([]uint64, 0),
		Errors:          make([]string, 0),
	}

	for _, playerID := range playerIDs {
		err := s.AddMember(ctx, teamID, playerID)
		if err != nil {
			result.FailedCount++
			result.FailedPlayerIDs = append(result.FailedPlayerIDs, playerID)
			result.Errors = append(result.Errors, fmt.Sprintf("陪玩师%d: %s", playerID, err.Error()))
		} else {
			result.SuccessCount++
		}
	}

	return result, nil
}
