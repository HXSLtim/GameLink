package playerrank

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

var (
	// ErrNotFound 记录不存在
	ErrNotFound = repository.ErrNotFound
	// ErrValidation 输入校验失败
	ErrValidation = errors.New("validation failed")
	// ErrPlayerNotFound 陪玩师不存在
	ErrPlayerNotFound = errors.New("player not found")
	// ErrGameNotFound 游戏不存在
	ErrGameNotFound = errors.New("game not found")
	// ErrRankNotFound 段位不存在
	ErrRankNotFound = errors.New("rank not found")
	// ErrAlreadyApplied 已申请过该游戏段位认证
	ErrAlreadyApplied = errors.New("already applied for this game rank")
	// ErrInvalidStatus 无效的状态变更
	ErrInvalidStatus = errors.New("invalid status transition")
)

// PlayerRankService 陪玩师段位认证服务
type PlayerRankService struct {
	records repository.PlayerRankRepository
	ranks   repository.GameRankRepository
	players repository.PlayerRepository
	games   repository.GameRepository
}

// NewPlayerRankService 创建陪玩师段位认证服务
func NewPlayerRankService(
	records repository.PlayerRankRepository,
	ranks repository.GameRankRepository,
	players repository.PlayerRepository,
	games repository.GameRepository,
) *PlayerRankService {
	return &PlayerRankService{
		records: records,
		ranks:   ranks,
		players: players,
		games:   games,
	}
}

// ApplyInput 申请段位认证输入
type ApplyInput struct {
	PlayerID       uint64
	GameID         uint64
	RankID         uint64
	ScreenshotURLs string // JSON数组
	Remark         string
}

// VerifyInput 审核段位认证输入
type VerifyInput struct {
	RecordID     uint64
	Status       model.PlayerRankStatus // verified/rejected/revoked
	VerifiedBy   uint64
	RejectReason string
}

// Apply 申请段位认证
func (s *PlayerRankService) Apply(ctx context.Context, input ApplyInput) (*model.PlayerRankRecord, error) {
	// 验证陪玩师是否存在
	if _, err := s.players.Get(ctx, input.PlayerID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrPlayerNotFound
		}
		return nil, fmt.Errorf("get player: %w", err)
	}

	// 验证游戏是否存在
	if _, err := s.games.Get(ctx, input.GameID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrGameNotFound
		}
		return nil, fmt.Errorf("get game: %w", err)
	}

	// 验证段位是否存在
	rank, err := s.ranks.Get(ctx, input.RankID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrRankNotFound
		}
		return nil, fmt.Errorf("get rank: %w", err)
	}

	// 验证段位属于该游戏
	if rank.GameID != input.GameID {
		return nil, fmt.Errorf("%w: rank does not belong to this game", ErrValidation)
	}

	// 检查是否已申请过该游戏的段位认证（pending 状态）
	existing, err := s.records.GetByPlayerAndGame(ctx, input.PlayerID, input.GameID)
	if err == nil && existing.Status == model.PlayerRankStatusPending {
		return nil, ErrAlreadyApplied
	}

	record := &model.PlayerRankRecord{
		PlayerID:       input.PlayerID,
		GameID:         input.GameID,
		RankID:         input.RankID,
		Status:         model.PlayerRankStatusPending,
		ScreenshotURLs: input.ScreenshotURLs,
		Remark:         input.Remark,
	}

	if err := s.records.Create(ctx, record); err != nil {
		return nil, fmt.Errorf("create record: %w", err)
	}

	return record, nil
}

// Verify 审核段位认证
func (s *PlayerRankService) Verify(ctx context.Context, input VerifyInput) (*model.PlayerRankRecord, error) {
	record, err := s.records.Get(ctx, input.RecordID)
	if err != nil {
		return nil, err
	}

	// 验证状态变更是否合法
	if !isValidStatusTransition(record.Status, input.Status) {
		return nil, fmt.Errorf("%w: cannot change from %s to %s", ErrInvalidStatus, record.Status, input.Status)
	}

	// 更新状态
	now := time.Now()
	record.Status = input.Status
	record.VerifiedBy = &input.VerifiedBy
	record.RejectReason = input.RejectReason

	if input.Status == model.PlayerRankStatusVerified {
		record.VerifiedAt = &now
	}

	if err := s.records.Update(ctx, record); err != nil {
		return nil, fmt.Errorf("update record: %w", err)
	}

	// 如果认证通过，更新陪玩师的主段位信息
	if input.Status == model.PlayerRankStatusVerified {
		if err := s.updatePlayerRank(ctx, record); err != nil {
			// 记录错误但不影响审核结果
			fmt.Printf("update player rank failed: %v\n", err)
		}
	}

	return record, nil
}

// updatePlayerRank 更新陪玩师的段位信息
func (s *PlayerRankService) updatePlayerRank(ctx context.Context, record *model.PlayerRankRecord) error {
	player, err := s.players.Get(ctx, record.PlayerID)
	if err != nil {
		return err
	}

	rank, err := s.ranks.Get(ctx, record.RankID)
	if err != nil {
		return err
	}

	// 更新陪玩师的段位名称（冗余字段）
	player.Rank = rank.Name
	player.HourlyRateCents = rank.PriceCents

	return s.players.Update(ctx, player)
}

// isValidStatusTransition 检查状态变更是否合法
func isValidStatusTransition(from, to model.PlayerRankStatus) bool {
	switch from {
	case model.PlayerRankStatusPending:
		// pending 可以变为 verified 或 rejected
		return to == model.PlayerRankStatusVerified || to == model.PlayerRankStatusRejected
	case model.PlayerRankStatusVerified:
		// verified 可以变为 revoked 或 expired
		return to == model.PlayerRankStatusRevoked || to == model.PlayerRankStatusExpired
	case model.PlayerRankStatusRejected:
		// rejected 可以重新申请（创建新记录），不能直接变更
		return false
	case model.PlayerRankStatusRevoked:
		// revoked 可以重新申请（创建新记录），不能直接变更
		return false
	case model.PlayerRankStatusExpired:
		// expired 可以重新申请（创建新记录），不能直接变更
		return false
	default:
		return false
	}
}

// Get 获取段位认证记录
func (s *PlayerRankService) Get(ctx context.Context, id uint64) (*model.PlayerRankRecord, error) {
	return s.records.GetWithRelations(ctx, id)
}

// ListByPlayerID 获取陪玩师的所有段位认证记录
func (s *PlayerRankService) ListByPlayerID(ctx context.Context, playerID uint64) ([]model.PlayerRankRecord, error) {
	return s.records.ListByPlayerID(ctx, playerID)
}

// ListPaged 分页获取段位认证记录
func (s *PlayerRankService) ListPaged(ctx context.Context, opts repository.PlayerRankListOptions) ([]model.PlayerRankRecord, *model.Pagination, error) {
	records, total, err := s.records.ListPaged(ctx, opts)
	if err != nil {
		return nil, nil, err
	}

	pagination := &model.Pagination{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Total:    int(total),
	}

	return records, pagination, nil
}

// ListPending 获取待审核的段位认证记录
func (s *PlayerRankService) ListPending(ctx context.Context, page, pageSize int) ([]model.PlayerRankRecord, *model.Pagination, error) {
	records, total, err := s.records.ListPending(ctx, page, pageSize)
	if err != nil {
		return nil, nil, err
	}

	pagination := &model.Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
	}

	return records, pagination, nil
}

// Delete 删除段位认证记录
func (s *PlayerRankService) Delete(ctx context.Context, id uint64) error {
	return s.records.Delete(ctx, id)
}

// GetStats 获取段位认证统计
func (s *PlayerRankService) GetStats(ctx context.Context) (map[model.PlayerRankStatus]int64, error) {
	return s.records.CountByStatus(ctx)
}

// GetPendingCount 获取待审核数量
func (s *PlayerRankService) GetPendingCount(ctx context.Context) (int64, error) {
	return s.records.GetPendingCount(ctx)
}
