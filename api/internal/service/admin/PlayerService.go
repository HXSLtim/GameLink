package admin

import (
	"context"
	"strings"
	"time"
	
	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/pkg/apierr"
)

// --- Player management ---

// CreatePlayerInput 创建陪玩资料。
type CreatePlayerInput struct {
	UserID             uint64
	Nickname           string
	Bio                string
	Rank               string
	HourlyRateCents    int64
	MainGameID         uint64
	VerificationStatus model.VerificationStatus
}

// UpdatePlayerInput 更新陪玩资料。
type UpdatePlayerInput struct {
	Nickname           string
	Bio                string
	Rank               string
	HourlyRateCents    int64
	MainGameID         uint64
	VerificationStatus model.VerificationStatus
}

// UpdateVerificationInput 审核陪玩师请求参数
type UpdateVerificationInput struct {
	Nickname           string
	Bio                string
	HourlyRateCents    int64
	MainGameID         uint64
	VerificationStatus model.VerificationStatus
	VerifiedBy         uint64 // 审核人ID
	Remark             string // 审核备注
}

// ListPlayers 返回陪玩列表。
func (s *AdminService) ListPlayers(ctx context.Context) ([]model.Player, error) {
	return getCachedList(ctx, s.cache, cacheKeyPlayers, listCacheTTL, func() ([]model.Player, error) {
		return s.players.List(ctx)
	})
}

// ListPlayersPaged 返回分页陪玩列表。
func (s *AdminService) ListPlayersPaged(ctx context.Context, page, pageSize int) ([]model.Player, *model.Pagination, error) {
	page = repository.NormalizePage(page)
	pageSize = repository.NormalizePageSize(pageSize)
	items, total, err := s.players.ListPaged(ctx, page, pageSize)
	if err != nil {
		return nil, nil, err
	}
	p := buildPagination(page, pageSize, total)
	return items, &p, nil
}

// ListPlayersPagedWithFilter 返回带筛选的分页陪玩列表。
func (s *AdminService) ListPlayersPagedWithFilter(ctx context.Context, page, pageSize int, keyword string, status *model.VerificationStatus) ([]model.Player, *model.Pagination, error) {
	page = repository.NormalizePage(page)
	pageSize = repository.NormalizePageSize(pageSize)
	items, total, err := s.players.ListPagedWithFilter(ctx, page, pageSize, keyword, status)
	if err != nil {
		return nil, nil, err
	}
	p := buildPagination(page, pageSize, total)
	return items, &p, nil
}

// BatchUpdatePlayerStatus 批量更新陪玩师状态。
func (s *AdminService) BatchUpdatePlayerStatus(ctx context.Context, ids []uint64, status model.VerificationStatus) (int64, error) {
	if len(ids) == 0 {
		return 0, apierr.BadRequest("no player ids provided")
	}
	updated, err := s.players.BatchUpdateStatus(ctx, ids, status)
	if err != nil {
		return 0, WrapError(err, "batch update player status")
	}
	s.invalidateCache(ctx, cacheKeyPlayers)
	return updated, nil
}

// BatchDeletePlayers 批量删除陪玩师。
func (s *AdminService) BatchDeletePlayers(ctx context.Context, ids []uint64) (int64, error) {
	if len(ids) == 0 {
		return 0, apierr.BadRequest("no player ids provided")
	}
	deleted, err := s.players.BatchDelete(ctx, ids)
	if err != nil {
		return 0, WrapError(err, "batch delete players")
	}
	s.invalidateCache(ctx, cacheKeyPlayers)
	return deleted, nil
}

// GetPlayer 返回陪玩详情。
func (s *AdminService) GetPlayer(ctx context.Context, id uint64) (*model.Player, error) {
	player, err := s.players.Get(ctx, id)
	if err != nil {
		return nil, WrapError(err, "get player")
	}
	return player, nil
}

// CreatePlayer 新建陪玩档案。
func (s *AdminService) CreatePlayer(ctx context.Context, input CreatePlayerInput) (*model.Player, error) {
	if err := validatePlayerInput(input.UserID, input.VerificationStatus); err != nil {
		return nil, err
	}

	player := &model.Player{
		UserID:             input.UserID,
		Nickname:           strings.TrimSpace(input.Nickname),
		Bio:                strings.TrimSpace(input.Bio),
		HourlyRateCents:    input.HourlyRateCents,
		MainGameID:         input.MainGameID,
		VerificationStatus: input.VerificationStatus,
	}

	if err := s.players.Create(ctx, player); err != nil {
		return nil, err
	}
	s.invalidateCache(ctx, cacheKeyPlayers)
	// audit
	s.appendLogAsync(ctx, string(model.OpEntityPlayer), player.ID, string(model.OpActionCreate), map[string]any{"user_id": player.UserID})
	return player, nil
}

// UpdatePlayer 调整陪玩信息。
func (s *AdminService) UpdatePlayer(ctx context.Context, id uint64, input UpdatePlayerInput) (*model.Player, error) {
	player, err := s.players.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := validatePlayerInput(player.UserID, input.VerificationStatus); err != nil {
		return nil, err
	}

	player.Nickname = strings.TrimSpace(input.Nickname)
	player.Bio = strings.TrimSpace(input.Bio)
	player.HourlyRateCents = input.HourlyRateCents
	player.MainGameID = input.MainGameID
	player.VerificationStatus = input.VerificationStatus

	if err := s.players.Update(ctx, player); err != nil {
		return nil, err
	}
	s.invalidateCache(ctx, cacheKeyPlayers)
	// audit
	s.appendLogAsync(ctx, string(model.OpEntityPlayer), player.ID, string(model.OpActionUpdate), map[string]any{"main_game_id": player.MainGameID})
	return player, nil
}

// UpdatePlayerVerification 审核陪玩师（保存审核记录）
func (s *AdminService) UpdatePlayerVerification(ctx context.Context, id uint64, input UpdateVerificationInput) (*model.Player, error) {
	player, err := s.players.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := validatePlayerInput(player.UserID, input.VerificationStatus); err != nil {
		return nil, err
	}

	// 保留原有信息
	player.Nickname = strings.TrimSpace(input.Nickname)
	player.Bio = strings.TrimSpace(input.Bio)
	player.HourlyRateCents = input.HourlyRateCents
	player.MainGameID = input.MainGameID

	// 更新审核状态和记录
	oldStatus := player.VerificationStatus
	player.VerificationStatus = input.VerificationStatus

	// 记录审核信息
	now := time.Now()
	player.VerifiedAt = &now
	player.VerifiedBy = &input.VerifiedBy
	player.VerifyRemark = strings.TrimSpace(input.Remark)

	// 如果是拒绝，保存拒绝原因
	if input.VerificationStatus == model.VerificationRejected {
		player.RejectReason = strings.TrimSpace(input.Remark)
	} else {
		player.RejectReason = ""
	}

	if err := s.players.Update(ctx, player); err != nil {
		return nil, err
	}
	s.invalidateCache(ctx, cacheKeyPlayers)

	// 审计日志
	s.appendLogAsync(ctx, string(model.OpEntityPlayer), player.ID, "verify", map[string]any{
		"old_status":  oldStatus,
		"new_status":  input.VerificationStatus,
		"verified_by": input.VerifiedBy,
		"remark":      input.Remark,
	})

	// Notification to player will be sent when notification service is integrated

	return player, nil
}

// DeletePlayer 删除陪玩档案。
func (s *AdminService) DeletePlayer(ctx context.Context, id uint64) error {
	if err := s.players.Delete(ctx, id); err != nil {
		return err
	}
	s.invalidateCache(ctx, cacheKeyPlayers)
	// audit
	s.appendLogAsync(ctx, string(model.OpEntityPlayer), id, string(model.OpActionDelete), nil)
	return nil
}

func validatePlayerInput(userID uint64, verification model.VerificationStatus) error {
	if userID == 0 {
		return ErrValidation
	}
	if verification == "" {
		return ErrValidation
	}
	return nil
}

