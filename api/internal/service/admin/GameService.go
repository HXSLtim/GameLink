package admin

import (
	"context"
	"strings"
	
	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/pkg/apierr"
)

// --- Game management ---

// CreateGameInput 创建游戏时使用的参数。
type CreateGameInput struct {
	Key         string
	Name        string
	Category    string
	IconURL     string
	Description string
}

// UpdateGameInput 修改游戏资料。
type UpdateGameInput struct {
	Key         string
	Name        string
	Category    string
	IconURL     string
	Description string
}

// ListGames 返回全部游戏。
func (s *AdminService) ListGames(ctx context.Context) ([]model.Game, error) {
	return getCachedList(ctx, s.cache, cacheKeyGames, listCacheTTL, func() ([]model.Game, error) {
		return s.games.List(ctx)
	})
}

// ListGamesPaged 返回分页游戏列表。
func (s *AdminService) ListGamesPaged(ctx context.Context, page, pageSize int) ([]model.Game, *model.Pagination, error) {
	page = repository.NormalizePage(page)
	pageSize = repository.NormalizePageSize(pageSize)
	items, total, err := s.games.ListPaged(ctx, page, pageSize)
	if err != nil {
		return nil, nil, err
	}
	p := buildPagination(page, pageSize, total)
	return items, &p, nil
}

// ListGamesPagedWithFilter 返回带筛选的分页游戏列表。
func (s *AdminService) ListGamesPagedWithFilter(ctx context.Context, page, pageSize int, keyword string) ([]model.Game, *model.Pagination, error) {
	page = repository.NormalizePage(page)
	pageSize = repository.NormalizePageSize(pageSize)
	items, total, err := s.games.ListPagedWithFilter(ctx, page, pageSize, keyword)
	if err != nil {
		return nil, nil, err
	}
	p := buildPagination(page, pageSize, total)
	return items, &p, nil
}

// BatchDeleteGames 批量删除游戏。
func (s *AdminService) BatchDeleteGames(ctx context.Context, ids []uint64) (int64, error) {
	if len(ids) == 0 {
		return 0, apierr.BadRequest("no game ids provided")
	}
	deleted, err := s.games.BatchDelete(ctx, ids)
	if err != nil {
		return 0, WrapError(err, "batch delete games")
	}
	s.invalidateCache(ctx, cacheKeyGames)
	return deleted, nil
}

// GetGame 获取单个游戏详情。
func (s *AdminService) GetGame(ctx context.Context, id uint64) (*model.Game, error) {
	game, err := s.games.Get(ctx, id)
	if err != nil {
		return nil, WrapError(err, "get game")
	}
	return game, nil
}

// CreateGame 创建游戏。
func (s *AdminService) CreateGame(ctx context.Context, input CreateGameInput) (*model.Game, error) {
	if err := validateGameInput(input.Key, input.Name); err != nil {
		return nil, err
	}

	game := &model.Game{
		Key:         strings.TrimSpace(input.Key),
		Name:        strings.TrimSpace(input.Name),
		Category:    strings.TrimSpace(input.Category),
		IconURL:     strings.TrimSpace(input.IconURL),
		Description: strings.TrimSpace(input.Description),
	}

	if err := s.games.Create(ctx, game); err != nil {
		return nil, WrapError(err, "create game")
	}

	s.invalidateCache(ctx, cacheKeyGames)
	// audit
	s.appendLogAsync(ctx, string(model.OpEntityGame), game.ID, string(model.OpActionCreate), map[string]any{"key": game.Key})

	return game, nil
}

// UpdateGame 更新游戏。
func (s *AdminService) UpdateGame(ctx context.Context, id uint64, input UpdateGameInput) (*model.Game, error) {
	if err := validateGameInput(input.Key, input.Name); err != nil {
		return nil, err
	}

	game, err := s.games.Get(ctx, id)
	if err != nil {
		return nil, WrapError(err, "get game")
	}

	game.Key = strings.TrimSpace(input.Key)
	game.Name = strings.TrimSpace(input.Name)
	game.Category = strings.TrimSpace(input.Category)
	game.IconURL = strings.TrimSpace(input.IconURL)
	game.Description = strings.TrimSpace(input.Description)

	if err := s.games.Update(ctx, game); err != nil {
		return nil, WrapError(err, "update game")
	}

	s.invalidateCache(ctx, cacheKeyGames)
	// audit
	s.appendLogAsync(ctx, string(model.OpEntityGame), game.ID, string(model.OpActionUpdate), map[string]any{"key": game.Key})

	return game, nil
}

// DeleteGame 删除游戏。
func (s *AdminService) DeleteGame(ctx context.Context, id uint64) error {
	if err := s.games.Delete(ctx, id); err != nil {
		return WrapError(err, "delete game")
	}
	s.invalidateCache(ctx, cacheKeyGames)
	// audit
	s.appendLogAsync(ctx, string(model.OpEntityGame), id, string(model.OpActionDelete), nil)
	return nil
}

func validateGameInput(key, name string) error {
	if strings.TrimSpace(key) == "" || strings.TrimSpace(name) == "" {
		return ErrValidation
	}
	return nil
}

// BatchUpdateGamesStatus 批量更新游戏状态（启用/禁用）。
func (s *AdminService) BatchUpdateGamesStatus(ctx context.Context, ids []uint64, isActive bool) (int64, error) {
	if len(ids) == 0 {
		return 0, apierr.BadRequest("no game ids provided")
	}
	updated, err := s.games.BatchUpdateStatus(ctx, ids, isActive)
	if err != nil {
		return 0, WrapError(err, "batch update games status")
	}
	s.invalidateCache(ctx, cacheKeyGames)
	for _, id := range ids {
		s.appendLogAsync(ctx, string(model.OpEntityGame), id, string(model.OpActionUpdate), map[string]any{"is_active": isActive})
	}
	return updated, nil
}

// BatchUpdateGamesSortOrder 批量更新游戏排序。
func (s *AdminService) BatchUpdateGamesSortOrder(ctx context.Context, updates map[uint64]int) (int64, error) {
	if len(updates) == 0 {
		return 0, apierr.BadRequest("no updates provided")
	}
	updated, err := s.games.BatchUpdateSortOrder(ctx, updates)
	if err != nil {
		return 0, WrapError(err, "batch update games sort order")
	}
	s.invalidateCache(ctx, cacheKeyGames)
	for id := range updates {
		s.appendLogAsync(ctx, string(model.OpEntityGame), id, string(model.OpActionUpdate), map[string]any{"sort_order": updates[id]})
	}
	return updated, nil
}

