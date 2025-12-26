package admin

import (
	"context"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// MenuService 提供菜单管理
type MenuService struct {
	menus repository.MenuRepository
}

func NewMenuService(menus repository.MenuRepository) *MenuService {
	return &MenuService{menus: menus}
}

// Create 创建菜单
func (s *MenuService) Create(ctx context.Context, menu *model.Menu) error {
	return s.menus.Create(ctx, menu)
}

// Update 更新菜单
func (s *MenuService) Update(ctx context.Context, menu *model.Menu) error {
	if menu.ID == 0 {
		return repository.ErrNotFound
	}
	return s.menus.Update(ctx, menu)
}

// Delete 删除菜单
func (s *MenuService) Delete(ctx context.Context, id uint64) error {
	return s.menus.Delete(ctx, id)
}

// Get 获取菜单
func (s *MenuService) Get(ctx context.Context, id uint64) (*model.Menu, error) {
	return s.menus.Get(ctx, id)
}

// List 按 parentID 列表
func (s *MenuService) List(ctx context.Context, parentID *uint64) ([]model.Menu, error) {
	return s.menus.List(ctx, parentID)
}

// ListPaged 分页列表
func (s *MenuService) ListPaged(ctx context.Context, page, pageSize int, parentID *uint64) ([]model.Menu, int64, error) {
	return s.menus.ListPaged(ctx, page, pageSize, parentID)
}

// ListAccessible 根据权限码筛选可见菜单
func (s *MenuService) ListAccessible(ctx context.Context, codes []string) ([]model.Menu, error) {
	if len(codes) == 0 {
		return []model.Menu{}, nil
	}
	return s.menus.ListByPermission(ctx, codes)
}
