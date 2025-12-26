package team

import (
	"context"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"

	"gorm.io/gorm"
)

// Repository 团队仓库
type Repository struct {
	db *gorm.DB
}

// NewTeamRepository 创建团队仓库
func NewTeamRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ============================================================================
// 团队 CRUD
// ============================================================================

// TeamListOptions 团队列表查询选项
type TeamListOptions struct {
	Page      int
	PageSize  int
	Keyword   string
	Status    *model.TeamStatus
	LeaderID  *uint64
	MinMember *int
	MaxMember *int
}

// Create 创建团队
func (r *Repository) Create(ctx context.Context, team *model.Team) error {
	return r.db.WithContext(ctx).Create(team).Error
}

// GetByID 根据ID获取团队
func (r *Repository) GetByID(ctx context.Context, id uint64) (*model.Team, error) {
	var team model.Team
	err := r.db.WithContext(ctx).
		Preload("Leader").
		Preload("Members", "status = ?", model.TeamMemberStatusActive).
		Preload("Members.Player").
		First(&team, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &team, nil
}

// Update 更新团队
func (r *Repository) Update(ctx context.Context, team *model.Team) error {
	result := r.db.WithContext(ctx).Model(team).Updates(map[string]interface{}{
		"name":              team.Name,
		"description":       team.Description,
		"avatar_url":        team.AvatarURL,
		"max_members":       team.MaxMembers,
		"income_share_type": team.IncomeShareType,
		"leader_bonus_rate": team.LeaderBonusRate,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// Delete 删除团队（软删除）
func (r *Repository) Delete(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).Delete(&model.Team{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// List 获取团队列表
func (r *Repository) List(ctx context.Context, opts TeamListOptions) ([]model.Team, int64, error) {
	var teams []model.Team
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Team{})

	if opts.Keyword != "" {
		query = query.Where("name ILIKE ?", "%"+opts.Keyword+"%")
	}
	if opts.Status != nil {
		query = query.Where("status = ?", *opts.Status)
	}
	if opts.LeaderID != nil {
		query = query.Where("leader_id = ?", *opts.LeaderID)
	}
	if opts.MinMember != nil {
		query = query.Where("member_count >= ?", *opts.MinMember)
	}
	if opts.MaxMember != nil {
		query = query.Where("member_count <= ?", *opts.MaxMember)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (opts.Page - 1) * opts.PageSize
	err := query.
		Preload("Leader").
		Order("created_at DESC").
		Offset(offset).
		Limit(opts.PageSize).
		Find(&teams).Error

	return teams, total, err
}

// UpdateStatus 更新团队状态
func (r *Repository) UpdateStatus(ctx context.Context, id uint64, status model.TeamStatus) error {
	result := r.db.WithContext(ctx).Model(&model.Team{}).
		Where("id = ?", id).
		Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// UpdateCurrentOrder 更新当前订单
func (r *Repository) UpdateCurrentOrder(ctx context.Context, id uint64, orderID *uint64, status model.TeamStatus) error {
	return r.db.WithContext(ctx).Model(&model.Team{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"current_order_id": orderID,
			"status":           status,
		}).Error
}

// IncrementStats 增加统计数据
func (r *Repository) IncrementStats(ctx context.Context, id uint64, orderCount int, incomeCents int64) error {
	return r.db.WithContext(ctx).Model(&model.Team{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"total_order_count":  gorm.Expr("total_order_count + ?", orderCount),
			"total_income_cents": gorm.Expr("total_income_cents + ?", incomeCents),
		}).Error
}

// UpdateMemberCount 更新成员数量
func (r *Repository) UpdateMemberCount(ctx context.Context, id uint64, delta int) error {
	return r.db.WithContext(ctx).Model(&model.Team{}).
		Where("id = ?", id).
		Update("member_count", gorm.Expr("member_count + ?", delta)).Error
}

// UpdateLeader 更新队长
func (r *Repository) UpdateLeader(ctx context.Context, id uint64, leaderID uint64) error {
	return r.db.WithContext(ctx).Model(&model.Team{}).
		Where("id = ?", id).
		Update("leader_id", leaderID).Error
}

// ============================================================================
// 团队成员 CRUD
// ============================================================================

// MemberListOptions 成员列表查询选项
type MemberListOptions struct {
	Page     int
	PageSize int
	TeamID   *uint64
	PlayerID *uint64
	Role     *model.TeamMemberRole
	Status   *model.TeamMemberStatus
}

// CreateMember 创建团队成员
func (r *Repository) CreateMember(ctx context.Context, member *model.TeamMember) error {
	return r.db.WithContext(ctx).Create(member).Error
}

// GetMemberByID 根据ID获取成员
func (r *Repository) GetMemberByID(ctx context.Context, id uint64) (*model.TeamMember, error) {
	var member model.TeamMember
	err := r.db.WithContext(ctx).
		Preload("Team").
		Preload("Player").
		First(&member, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &member, nil
}

// GetMemberByTeamAndPlayer 根据团队ID和陪玩师ID获取成员
func (r *Repository) GetMemberByTeamAndPlayer(ctx context.Context, teamID, playerID uint64) (*model.TeamMember, error) {
	var member model.TeamMember
	err := r.db.WithContext(ctx).
		Where("team_id = ? AND player_id = ?", teamID, playerID).
		First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &member, nil
}

// GetActiveMemberByPlayer 获取陪玩师当前所在的活跃团队成员记录
func (r *Repository) GetActiveMemberByPlayer(ctx context.Context, playerID uint64) (*model.TeamMember, error) {
	var member model.TeamMember
	err := r.db.WithContext(ctx).
		Where("player_id = ? AND status = ?", playerID, model.TeamMemberStatusActive).
		Preload("Team").
		First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &member, nil
}

// UpdateMember 更新成员
func (r *Repository) UpdateMember(ctx context.Context, member *model.TeamMember) error {
	return r.db.WithContext(ctx).Model(member).Updates(map[string]interface{}{
		"role":       member.Role,
		"status":     member.Status,
		"sort_order": member.SortOrder,
		"left_at":    member.LeftAt,
	}).Error
}

// ListMembers 获取成员列表
func (r *Repository) ListMembers(ctx context.Context, opts MemberListOptions) ([]model.TeamMember, int64, error) {
	var members []model.TeamMember
	var total int64

	query := r.db.WithContext(ctx).Model(&model.TeamMember{})

	if opts.TeamID != nil {
		query = query.Where("team_id = ?", *opts.TeamID)
	}
	if opts.PlayerID != nil {
		query = query.Where("player_id = ?", *opts.PlayerID)
	}
	if opts.Role != nil {
		query = query.Where("role = ?", *opts.Role)
	}
	if opts.Status != nil {
		query = query.Where("status = ?", *opts.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (opts.Page - 1) * opts.PageSize
	err := query.
		Preload("Player").
		Order("sort_order ASC, joined_at ASC").
		Offset(offset).
		Limit(opts.PageSize).
		Find(&members).Error

	return members, total, err
}

// GetActiveMembers 获取团队活跃成员
func (r *Repository) GetActiveMembers(ctx context.Context, teamID uint64) ([]model.TeamMember, error) {
	var members []model.TeamMember
	err := r.db.WithContext(ctx).
		Where("team_id = ? AND status = ?", teamID, model.TeamMemberStatusActive).
		Preload("Player").
		Order("sort_order ASC").
		Find(&members).Error
	return members, err
}

// GetNextLeader 获取下一个队长候选人
func (r *Repository) GetNextLeader(ctx context.Context, teamID uint64, excludePlayerID uint64) (*model.TeamMember, error) {
	var member model.TeamMember
	err := r.db.WithContext(ctx).
		Where("team_id = ? AND player_id != ? AND status = ?", teamID, excludePlayerID, model.TeamMemberStatusActive).
		Order("sort_order ASC, joined_at ASC").
		First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &member, nil
}

// IncrementMemberStats 增加成员统计数据
func (r *Repository) IncrementMemberStats(ctx context.Context, id uint64, orderCount int, incomeCents int64) error {
	return r.db.WithContext(ctx).Model(&model.TeamMember{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"order_count":  gorm.Expr("order_count + ?", orderCount),
			"income_cents": gorm.Expr("income_cents + ?", incomeCents),
		}).Error
}

// ============================================================================
// 团队邀请 CRUD
// ============================================================================

// InviteListOptions 邀请列表查询选项
type InviteListOptions struct {
	Page     int
	PageSize int
	TeamID   *uint64
	PlayerID *uint64
	Status   *model.TeamInviteStatus
}

// CreateInvite 创建邀请
func (r *Repository) CreateInvite(ctx context.Context, invite *model.TeamInvite) error {
	return r.db.WithContext(ctx).Create(invite).Error
}

// GetInviteByID 根据ID获取邀请
func (r *Repository) GetInviteByID(ctx context.Context, id uint64) (*model.TeamInvite, error) {
	var invite model.TeamInvite
	err := r.db.WithContext(ctx).
		Preload("Team").
		Preload("Player").
		Preload("Inviter").
		First(&invite, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &invite, nil
}

// GetPendingInvite 获取待处理的邀请
func (r *Repository) GetPendingInvite(ctx context.Context, teamID, playerID uint64) (*model.TeamInvite, error) {
	var invite model.TeamInvite
	err := r.db.WithContext(ctx).
		Where("team_id = ? AND player_id = ? AND status = ? AND expire_at > ?",
			teamID, playerID, model.TeamInviteStatusPending, time.Now()).
		First(&invite).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &invite, nil
}

// UpdateInvite 更新邀请
func (r *Repository) UpdateInvite(ctx context.Context, invite *model.TeamInvite) error {
	return r.db.WithContext(ctx).Model(invite).Updates(map[string]interface{}{
		"status": invite.Status,
	}).Error
}

// ListInvites 获取邀请列表
func (r *Repository) ListInvites(ctx context.Context, opts InviteListOptions) ([]model.TeamInvite, int64, error) {
	var invites []model.TeamInvite
	var total int64

	query := r.db.WithContext(ctx).Model(&model.TeamInvite{})

	if opts.TeamID != nil {
		query = query.Where("team_id = ?", *opts.TeamID)
	}
	if opts.PlayerID != nil {
		query = query.Where("player_id = ?", *opts.PlayerID)
	}
	if opts.Status != nil {
		query = query.Where("status = ?", *opts.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (opts.Page - 1) * opts.PageSize
	err := query.
		Preload("Team").
		Preload("Player").
		Preload("Inviter").
		Order("created_at DESC").
		Offset(offset).
		Limit(opts.PageSize).
		Find(&invites).Error

	return invites, total, err
}

// ExpireInvites 过期邀请
func (r *Repository) ExpireInvites(ctx context.Context) (int64, error) {
	result := r.db.WithContext(ctx).Model(&model.TeamInvite{}).
		Where("status = ? AND expire_at <= ?", model.TeamInviteStatusPending, time.Now()).
		Update("status", model.TeamInviteStatusExpired)
	return result.RowsAffected, result.Error
}

// ============================================================================
// 统计查询
// ============================================================================

// GetTeamStats 获取团队统计
func (r *Repository) GetTeamStats(ctx context.Context) (*TeamStats, error) {
	var stats TeamStats

	// 总团队数
	r.db.WithContext(ctx).Model(&model.Team{}).Count(&stats.TotalTeams)

	// 活跃团队数
	r.db.WithContext(ctx).Model(&model.Team{}).
		Where("status = ?", model.TeamStatusActive).
		Count(&stats.ActiveTeams)

	// 忙碌团队数
	r.db.WithContext(ctx).Model(&model.Team{}).
		Where("status = ?", model.TeamStatusBusy).
		Count(&stats.BusyTeams)

	// 总成员数
	r.db.WithContext(ctx).Model(&model.TeamMember{}).
		Where("status = ?", model.TeamMemberStatusActive).
		Count(&stats.TotalMembers)

	// 总订单数和收入
	r.db.WithContext(ctx).Model(&model.Team{}).
		Select("COALESCE(SUM(total_order_count), 0) as total_orders, COALESCE(SUM(total_income_cents), 0) as total_income_cents").
		Row().Scan(&stats.TotalOrders, &stats.TotalIncomeCents)

	return &stats, nil
}

// TeamStats 团队统计
type TeamStats struct {
	TotalTeams       int64 `json:"totalTeams"`
	ActiveTeams      int64 `json:"activeTeams"`
	BusyTeams        int64 `json:"busyTeams"`
	TotalMembers     int64 `json:"totalMembers"`
	TotalOrders      int64 `json:"totalOrders"`
	TotalIncomeCents int64 `json:"totalIncomeCents"`
}
