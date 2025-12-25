package chat

import (
	"context"
	"time"

	"gamelink/internal/repository"
)

// CleanupService 负责聊天组清理任务（例如定期删除已停用的订单群）
type CleanupService struct {
	groups repository.ChatGroupRepository
	now    func() time.Time
}

// NewCleanupService 创建清理服务
func NewCleanupService(groups repository.ChatGroupRepository) *CleanupService {
	return &CleanupService{
		groups: groups,
		now:    time.Now,
	}
}

// CleanupInactiveOrderGroups 删除在 cutoff 之前已停用的订单群，返回删除数量
func (s *CleanupService) CleanupInactiveOrderGroups(ctx context.Context, cutoff time.Time, limit int) (int, error) {
	groups, err := s.groups.ListDeactivatedBefore(ctx, cutoff, limit)
	if err != nil {
		return 0, err
	}
	if len(groups) == 0 {
		return 0, nil
	}
	ids := make([]uint64, 0, len(groups))
	for _, g := range groups {
		ids = append(ids, g.ID)
	}
	if err := s.groups.DeleteByIDs(ctx, ids); err != nil {
		return 0, err
	}
	return len(ids), nil
}
