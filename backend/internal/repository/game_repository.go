package repository

import (
	"context"

	"gamelink/internal/model"
)

// GameRepository 定义后台管理对游戏实体的数据库操作。
type GameRepository interface {
	List(ctx context.Context) ([]model.Game, error)
	// ListPaged 支持分页列表与总数统计
	ListPaged(ctx context.Context, page, pageSize int) ([]model.Game, int64, error)
	Get(ctx context.Context, id uint64) (*model.Game, error)
	Create(ctx context.Context, game *model.Game) error
	Update(ctx context.Context, game *model.Game) error
	Delete(ctx context.Context, id uint64) error
}
