package admin

import (
	"context"
	"fmt"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/pkg/apierr"
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

// ============================================================================
// Batch Operation Input/Response Types
// ============================================================================

// BatchDeleteMenusInput 批量删除菜单输入
type BatchDeleteMenusInput struct {
	MenuIDs []uint64
}

// BatchUpdateMenuStatusInput 批量更新菜单状态输入
type BatchUpdateMenuStatusInput struct {
	MenuIDs []uint64
	Status  string // "enabled" or "disabled"
}

// BatchUpdateMenuOrderInput 批量更新菜单排序输入
type BatchUpdateMenuOrderInput struct {
	MenuOrders []MenuOrderUpdate
}

// MenuOrderUpdate 菜单排序更新
type MenuOrderUpdate struct {
	MenuID uint64
	Order  int
}

// ============================================================================
// Batch Operations
// ============================================================================

// BatchDeleteMenus 批量删除菜单
func (s *MenuService) BatchDeleteMenus(ctx context.Context, input BatchDeleteMenusInput) (*BatchOperationResponse, error) {
	if len(input.MenuIDs) == 0 {
		return nil, apierr.BadRequest("menu_ids is required")
	}
	if len(input.MenuIDs) > 100 {
		return nil, apierr.BadRequest("maximum 100 menus per batch")
	}

	result := &BatchOperationResponse{
		TotalCount:   len(input.MenuIDs),
		FailedItems:  make([]BatchErrorItem, 0),
		SuccessItems: make([]uint64, 0),
	}

	for _, menuID := range input.MenuIDs {
		// 检查是否有子菜单
		hasChildren, err := s.menus.HasChildren(ctx, menuID)
		if err != nil {
			result.FailedCount++
			result.FailedItems = append(result.FailedItems, BatchErrorItem{
				ID:      menuID,
				Message: "failed to check children: " + err.Error(),
			})
			continue
		}
		if hasChildren {
			result.FailedCount++
			result.FailedItems = append(result.FailedItems, BatchErrorItem{
				ID:      menuID,
				Message: "cannot delete menu with child menus",
			})
			continue
		}

		// 删除菜单
		if err := s.menus.Delete(ctx, menuID); err != nil {
			result.FailedCount++
			result.FailedItems = append(result.FailedItems, BatchErrorItem{
				ID:      menuID,
				Message: err.Error(),
			})
		} else {
			result.SuccessCount++
			result.SuccessItems = append(result.SuccessItems, menuID)
		}
	}

	return result, nil
}

// BatchUpdateMenuStatus 批量更新菜单状态
func (s *MenuService) BatchUpdateMenuStatus(ctx context.Context, input BatchUpdateMenuStatusInput) (*BatchOperationResponse, error) {
	if len(input.MenuIDs) == 0 {
		return nil, apierr.BadRequest("menu_ids is required")
	}
	if len(input.MenuIDs) > 100 {
		return nil, apierr.BadRequest("maximum 100 menus per batch")
	}
	if input.Status != "enabled" && input.Status != "disabled" {
		return nil, apierr.BadRequest("status must be 'enabled' or 'disabled'")
	}

	// 将 status 转换为 hidden 值
	hidden := input.Status == "disabled"

	result := &BatchOperationResponse{
		TotalCount:   len(input.MenuIDs),
		FailedItems:  make([]BatchErrorItem, 0),
		SuccessItems: make([]uint64, 0),
	}

	for _, menuID := range input.MenuIDs {
		menu, err := s.menus.Get(ctx, menuID)
		if err != nil {
			result.FailedCount++
			result.FailedItems = append(result.FailedItems, BatchErrorItem{
				ID:      menuID,
				Message: "menu not found",
			})
			continue
		}

		menu.Hidden = hidden
		if err := s.menus.Update(ctx, menu); err != nil {
			result.FailedCount++
			result.FailedItems = append(result.FailedItems, BatchErrorItem{
				ID:      menuID,
				Message: err.Error(),
			})
		} else {
			result.SuccessCount++
			result.SuccessItems = append(result.SuccessItems, menuID)
		}
	}

	return result, nil
}

// BatchUpdateMenuOrder 批量更新菜单排序
func (s *MenuService) BatchUpdateMenuOrder(ctx context.Context, input BatchUpdateMenuOrderInput) (*BatchOperationResponse, error) {
	if len(input.MenuOrders) == 0 {
		return nil, apierr.BadRequest("menu_orders is required")
	}
	if len(input.MenuOrders) > 100 {
		return nil, apierr.BadRequest("maximum 100 menus per batch")
	}

	// 验证 order 值的唯一性（在同一父菜单下）
	orderMap := make(map[int]uint64) // order -> menuID
	for _, mo := range input.MenuOrders {
		if existingID, ok := orderMap[mo.Order]; ok {
			return nil, apierr.BadRequest(fmt.Sprintf("duplicate order value %d for menus %d and %d", mo.Order, existingID, mo.MenuID))
		}
		orderMap[mo.Order] = mo.MenuID
	}

	result := &BatchOperationResponse{
		TotalCount:   len(input.MenuOrders),
		FailedItems:  make([]BatchErrorItem, 0),
		SuccessItems: make([]uint64, 0),
	}

	for _, menuOrder := range input.MenuOrders {
		menu, err := s.menus.Get(ctx, menuOrder.MenuID)
		if err != nil {
			result.FailedCount++
			result.FailedItems = append(result.FailedItems, BatchErrorItem{
				ID:      menuOrder.MenuID,
				Message: "menu not found",
			})
			continue
		}

		menu.Order = menuOrder.Order
		if err := s.menus.Update(ctx, menu); err != nil {
			result.FailedCount++
			result.FailedItems = append(result.FailedItems, BatchErrorItem{
				ID:      menuOrder.MenuID,
				Message: err.Error(),
			})
		} else {
			result.SuccessCount++
			result.SuccessItems = append(result.SuccessItems, menuOrder.MenuID)
		}
	}

	return result, nil
}

// ============================================================================
// Legacy Batch Methods (保持向后兼容)
// ============================================================================

// BatchDeleteResult 批量删除结果 (旧格式，保持兼容)
type BatchDeleteResult struct {
	SuccessCount int      `json:"successCount"`
	FailedCount  int      `json:"failedCount"`
	FailedIDs    []uint64 `json:"failedIds,omitempty"`
}

// BatchDelete 批量删除菜单（软删除）- 旧方法，保持兼容
func (s *MenuService) BatchDelete(ctx context.Context, ids []uint64) (*BatchDeleteResult, error) {
	if len(ids) == 0 {
		return nil, apierr.BadRequest("菜单ID列表不能为空")
	}

	result := &BatchDeleteResult{
		FailedIDs: make([]uint64, 0),
	}

	for _, id := range ids {
		if err := s.menus.Delete(ctx, id); err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, id)
		} else {
			result.SuccessCount++
		}
	}

	return result, nil
}

// MenuStatusUpdate 菜单状态更新请求 (旧格式，保持兼容)
type MenuStatusUpdate struct {
	ID     uint64 `json:"id"`
	Hidden *bool  `json:"hidden"`
}

// BatchUpdateStatusResult 批量状态更新结果 (旧格式，保持兼容)
type BatchUpdateStatusResult struct {
	SuccessCount int      `json:"successCount"`
	FailedCount  int      `json:"failedCount"`
	FailedIDs    []uint64 `json:"failedIds,omitempty"`
}

// BatchUpdateStatus 批量更新菜单显示/隐藏状态 - 旧方法，保持兼容
func (s *MenuService) BatchUpdateStatus(ctx context.Context, updates []MenuStatusUpdate) (*BatchUpdateStatusResult, error) {
	if len(updates) == 0 {
		return nil, apierr.BadRequest("更新列表不能为空")
	}

	result := &BatchUpdateStatusResult{
		FailedIDs: make([]uint64, 0),
	}

	for _, update := range updates {
		menu, err := s.menus.Get(ctx, update.ID)
		if err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, update.ID)
			continue
		}

		if update.Hidden != nil {
			menu.Hidden = *update.Hidden
		}

		if err := s.menus.Update(ctx, menu); err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, update.ID)
		} else {
			result.SuccessCount++
		}
	}

	return result, nil
}

// MenuSortUpdate 菜单排序更新请求 (旧格式，保持兼容)
type MenuSortUpdate struct {
	ID        uint64 `json:"id"`
	SortOrder int    `json:"sortOrder"`
}

// BatchUpdateSortResult 批量排序更新结果 (旧格式，保持兼容)
type BatchUpdateSortResult struct {
	SuccessCount int      `json:"successCount"`
	FailedCount  int      `json:"failedCount"`
	FailedIDs    []uint64 `json:"failedIds,omitempty"`
}

// BatchUpdateSort 批量更新菜单排序 - 旧方法，保持兼容
func (s *MenuService) BatchUpdateSort(ctx context.Context, updates []MenuSortUpdate) (*BatchUpdateSortResult, error) {
	if len(updates) == 0 {
		return nil, apierr.BadRequest("更新列表不能为空")
	}

	result := &BatchUpdateSortResult{
		FailedIDs: make([]uint64, 0),
	}

	for _, update := range updates {
		menu, err := s.menus.Get(ctx, update.ID)
		if err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, update.ID)
			continue
		}

		menu.Order = update.SortOrder

		if err := s.menus.Update(ctx, menu); err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, update.ID)
		} else {
			result.SuccessCount++
		}
	}

	return result, nil
}
