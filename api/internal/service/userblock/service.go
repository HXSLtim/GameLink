package userblock

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
	// ErrAlreadyBlocked 已拉黑
	ErrAlreadyBlocked = errors.New("already blocked")
	// ErrNotBlocked 未拉黑
	ErrNotBlocked = errors.New("not blocked")
	// ErrCannotBlockSelf 不能拉黑自己
	ErrCannotBlockSelf = errors.New("cannot block yourself")
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = errors.New("user not found")
)

// UserBlockService 用户拉黑服务
type UserBlockService struct {
	repo  repository.UserBlockRepository
	users repository.UserRepository
}

// NewUserBlockService 创建用户拉黑服务
func NewUserBlockService(
	repo repository.UserBlockRepository,
	users repository.UserRepository,
) *UserBlockService {
	return &UserBlockService{
		repo:  repo,
		users: users,
	}
}

// ============================================================================
// 拉黑操作
// ============================================================================

// BlockInput 拉黑输入
type BlockInput struct {
	BlockerID   uint64
	BlockerType model.BlockUserType
	BlockedID   uint64
	BlockedType model.BlockUserType
	Reason      string
}

// Block 拉黑用户
func (s *UserBlockService) Block(ctx context.Context, input BlockInput) (*model.UserBlock, error) {
	// 不能拉黑自己
	if input.BlockerID == input.BlockedID {
		return nil, ErrCannotBlockSelf
	}

	// 验证被拉黑用户是否存在
	if _, err := s.users.Get(ctx, input.BlockedID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get blocked user: %w", err)
	}

	// 检查是否已拉黑
	existing, err := s.repo.GetActiveByBlockerAndBlocked(ctx, input.BlockerID, input.BlockedID)
	if err == nil && existing != nil {
		return nil, ErrAlreadyBlocked
	}

	block := &model.UserBlock{
		BlockerID:   input.BlockerID,
		BlockerType: input.BlockerType,
		BlockedID:   input.BlockedID,
		BlockedType: input.BlockedType,
		Reason:      input.Reason,
		Status:      model.BlockStatusActive,
		BlockedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, block); err != nil {
		return nil, fmt.Errorf("create block: %w", err)
	}

	return block, nil
}

// Unblock 取消拉黑（用户主动）
func (s *UserBlockService) Unblock(ctx context.Context, blockerID, blockedID uint64) error {
	block, err := s.repo.GetActiveByBlockerAndBlocked(ctx, blockerID, blockedID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotBlocked
		}
		return fmt.Errorf("get block: %w", err)
	}

	block.Cancel()
	if err := s.repo.Update(ctx, block); err != nil {
		return fmt.Errorf("update block: %w", err)
	}

	return nil
}

// AdminUnblock 管理员强制取消拉黑
func (s *UserBlockService) AdminUnblock(ctx context.Context, id uint64, adminID uint64, remark string) error {
	block, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}

	if block.Status != model.BlockStatusActive {
		return ErrNotBlocked
	}

	block.AdminCancel(adminID, remark)
	if err := s.repo.Update(ctx, block); err != nil {
		return fmt.Errorf("update block: %w", err)
	}

	return nil
}

// ============================================================================
// 查询操作
// ============================================================================

// Get 获取拉黑记录
func (s *UserBlockService) Get(ctx context.Context, id uint64) (*model.UserBlock, error) {
	return s.repo.GetWithRelations(ctx, id)
}

// IsBlocked 检查两个用户之间是否存在拉黑关系（任一方向）
func (s *UserBlockService) IsBlocked(ctx context.Context, userID1, userID2 uint64) (bool, error) {
	return s.repo.IsBlocked(ctx, userID1, userID2)
}

// IsBlockedBy 检查 blockedID 是否被 blockerID 拉黑
func (s *UserBlockService) IsBlockedBy(ctx context.Context, blockerID, blockedID uint64) (bool, error) {
	return s.repo.IsBlockedBy(ctx, blockerID, blockedID)
}

// ListByBlocker 获取用户拉黑的列表
func (s *UserBlockService) ListByBlocker(ctx context.Context, blockerID uint64, activeOnly bool) ([]model.UserBlock, error) {
	var status *model.BlockStatus
	if activeOnly {
		active := model.BlockStatusActive
		status = &active
	}
	return s.repo.ListByBlockerID(ctx, blockerID, status)
}

// ListByBlocked 获取被拉黑的列表
func (s *UserBlockService) ListByBlocked(ctx context.Context, blockedID uint64, activeOnly bool) ([]model.UserBlock, error) {
	var status *model.BlockStatus
	if activeOnly {
		active := model.BlockStatusActive
		status = &active
	}
	return s.repo.ListByBlockedID(ctx, blockedID, status)
}

// ListPaged 分页获取拉黑记录
func (s *UserBlockService) ListPaged(ctx context.Context, opts repository.UserBlockListOptions) ([]model.UserBlock, *model.Pagination, error) {
	blocks, total, err := s.repo.ListPaged(ctx, opts)
	if err != nil {
		return nil, nil, err
	}

	pagination := &model.Pagination{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Total:    int(total),
	}

	return blocks, pagination, nil
}

// GetBlockedUserIDs 获取用户拉黑的所有用户ID列表
func (s *UserBlockService) GetBlockedUserIDs(ctx context.Context, blockerID uint64) ([]uint64, error) {
	return s.repo.GetBlockedUserIDs(ctx, blockerID)
}

// GetBlockerUserIDs 获取拉黑该用户的所有用户ID列表
func (s *UserBlockService) GetBlockerUserIDs(ctx context.Context, blockedID uint64) ([]uint64, error) {
	return s.repo.GetBlockerUserIDs(ctx, blockedID)
}

// GetAllBlockRelatedUserIDs 获取与用户有拉黑关系的所有用户ID（双向）
func (s *UserBlockService) GetAllBlockRelatedUserIDs(ctx context.Context, userID uint64) ([]uint64, error) {
	return s.repo.GetAllBlockRelatedUserIDs(ctx, userID)
}

// ============================================================================
// 统计操作
// ============================================================================

// GetStats 获取拉黑统计
func (s *UserBlockService) GetStats(ctx context.Context) (map[model.BlockStatus]int64, error) {
	return s.repo.CountByStatus(ctx)
}

// GetActiveCount 获取生效中的拉黑记录数量
func (s *UserBlockService) GetActiveCount(ctx context.Context) (int64, error) {
	return s.repo.GetActiveCount(ctx)
}

// ============================================================================
// 管理操作
// ============================================================================

// Delete 删除拉黑记录
func (s *UserBlockService) Delete(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}

// BatchUnblock 批量取消拉黑（管理员）
func (s *UserBlockService) BatchUnblock(ctx context.Context, ids []uint64, adminID uint64, remark string) (int, error) {
	var successCount int
	for _, id := range ids {
		if err := s.AdminUnblock(ctx, id, adminID, remark); err != nil {
			continue
		}
		successCount++
	}
	return successCount, nil
}
