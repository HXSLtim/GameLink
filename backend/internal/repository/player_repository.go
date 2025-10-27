package repository

import (
	"context"

	"gamelink/internal/model"
)

// PlayerRepository 管理陪玩资料。
type PlayerRepository interface {
	List(ctx context.Context) ([]model.Player, error)
	// ListPaged 支持分页列表与总数统计
	ListPaged(ctx context.Context, page, pageSize int) ([]model.Player, int64, error)
	Get(ctx context.Context, id uint64) (*model.Player, error)
	Create(ctx context.Context, player *model.Player) error
	Update(ctx context.Context, player *model.Player) error
	Delete(ctx context.Context, id uint64) error
}
