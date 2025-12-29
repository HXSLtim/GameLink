/**
 * @file game category service
 * @description 游戏分类业务逻辑层
 */

package gamecategory

import (
	"context"
	"errors"
	"fmt"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// GameCategoryService 游戏分类服务
type GameCategoryService struct {
	categoryRepo repository.GameCategoryRepository
	gameRepo     repository.GameRepository
	itemRepo     repository.ServiceItemRepository
}

// NewGameCategoryService 创建游戏分类服务
func NewGameCategoryService(
	categoryRepo repository.GameCategoryRepository,
	gameRepo repository.GameRepository,
	itemRepo repository.ServiceItemRepository,
) *GameCategoryService {
	return &GameCategoryService{
		categoryRepo: categoryRepo,
		gameRepo:     gameRepo,
		itemRepo:     itemRepo,
	}
}

// ============================================================================
// 请求/响应类型定义
// ============================================================================

// CreateCategoryRequest 创建分类请求
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=50"`
	Description string `json:"description" binding:"max=1000"`
	IconURL     string `json:"iconUrl" binding:"omitempty,url,max=255"`
	SortOrder   int    `json:"sortOrder"`
}

// UpdateCategoryRequest 更新分类请求
type UpdateCategoryRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=1,max=50"`
	Description *string `json:"description" binding:"omitempty,max=1000"`
	IconURL     *string `json:"iconUrl" binding:"omitempty,url,max=255"`
	SortOrder   *int    `json:"sortOrder"`
	IsActive    *bool   `json:"isActive"`
}

// CategoryResponse 分类响应
type CategoryResponse struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IconURL     string `json:"iconUrl,omitempty"`
	SortOrder   int    `json:"sortOrder"`
	IsActive    bool   `json:"isActive"`
	GameCount   int64  `json:"gameCount,omitempty"`
}

// BatchOperationResult 批量操作结果
type BatchOperationResult struct {
	SuccessCount int               `json:"successCount"`
	FailedCount  int               `json:"failedCount"`
	FailedIDs    []uint64          `json:"failedIds,omitempty"`
	Errors       []string          `json:"errors,omitempty"`
}

// ============================================================================
// CRUD 操作
// ============================================================================

// CreateCategory 创建游戏分类
// 业务规则：
//   - 分类名称必须唯一
//   - 名称长度 1-50 字符
//   - 描述最多 1000 字符
//   - 图标 URL 必须是有效 URL
func (s *GameCategoryService) CreateCategory(ctx context.Context, name, description, iconURL string, sortOrder int) (*model.GameCategory, error) {
	// 参数验证
	if name == "" {
		return nil, errors.New("分类名称不能为空")
	}
	if len(name) > 50 {
		return nil, errors.New("分类名称长度不能超过50个字符")
	}
	if len(description) > 1000 {
		return nil, errors.New("分类描述长度不能超过1000个字符")
	}

	// 检查名称是否已存在
	if _, err := s.categoryRepo.GetByName(ctx, name); err == nil {
		return nil, errors.New("分类名称已存在")
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("检查分类名称失败: %w", err)
	}

	// 创建分类
	category := &model.GameCategory{
		Name:        name,
		Description: description,
		IconURL:     iconURL,
		SortOrder:   sortOrder,
		IsActive:    true,
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, fmt.Errorf("创建分类失败: %w", err)
	}

	return category, nil
}

// GetCategory 获取分类详情
func (s *GameCategoryService) GetCategory(ctx context.Context, id uint64) (*model.GameCategory, error) {
	category, err := s.categoryRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, errors.New("分类不存在")
		}
		return nil, fmt.Errorf("获取分类失败: %w", err)
	}
	return category, nil
}

// ListCategories 获取分类列表
// 参数:
//   - page: 页码（从1开始）
//   - pageSize: 每页数量
//   - isActive: 启用状态过滤（nil表示不过滤）
//   - keyword: 关键词搜索（匹配名称和描述）
func (s *GameCategoryService) ListCategories(ctx context.Context, page, pageSize int, isActive *bool, keyword string) ([]*model.GameCategory, int64, error) {
	opts := repository.GameCategoryListOptions{
		Page:     page,
		PageSize: pageSize,
		IsActive: isActive,
		Keyword:  keyword,
	}

	categories, total, err := s.categoryRepo.List(ctx, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("获取分类列表失败: %w", err)
	}

	return categories, total, nil
}

// UpdateCategory 更新分类
// 业务规则：
//   - 更新名称时需检查重复（排除自身）
//   - 只更新传入的非空字段
func (s *GameCategoryService) UpdateCategory(ctx context.Context, id uint64, input *UpdateCategoryRequest) error {
	// 获取现有分类
	category, err := s.categoryRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return errors.New("分类不存在")
		}
		return fmt.Errorf("获取分类失败: %w", err)
	}

	// 如果更新名称，检查名称是否重复
	if input.Name != nil && *input.Name != category.Name {
		if len(*input.Name) > 50 {
			return errors.New("分类名称长度不能超过50个字符")
		}

		existing, err := s.categoryRepo.GetByName(ctx, *input.Name)
		if err != nil && !errors.Is(err, repository.ErrNotFound) {
			return fmt.Errorf("检查分类名称失败: %w", err)
		}
		if existing != nil && existing.ID != id {
			return errors.New("分类名称已存在")
		}
		category.Name = *input.Name
	}

	// 更新其他字段
	if input.Description != nil {
		if len(*input.Description) > 1000 {
			return errors.New("分类描述长度不能超过1000个字符")
		}
		category.Description = *input.Description
	}
	if input.IconURL != nil {
		category.IconURL = *input.IconURL
	}
	if input.SortOrder != nil {
		category.SortOrder = *input.SortOrder
	}
	if input.IsActive != nil {
		category.IsActive = *input.IsActive
	}

	// 执行更新
	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return fmt.Errorf("更新分类失败: %w", err)
	}

	return nil
}

// DeleteCategory 删除分类
// 业务规则：
//   - 检查是否有游戏关联到此分类（通过 category_id 字段）
//   - 检查是否有服务项目使用此分类名称
func (s *GameCategoryService) DeleteCategory(ctx context.Context, id uint64) error {
	// 获取分类信息（验证分类存在）
	_, err := s.categoryRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return errors.New("分类不存在")
		}
		return fmt.Errorf("获取分类失败: %w", err)
	}

	// 检查是否有关联的游戏
	gameCount, err := s.categoryRepo.CountGames(ctx, id)
	if err != nil {
		return fmt.Errorf("检查游戏关联失败: %w", err)
	}
	if gameCount > 0 {
		return fmt.Errorf("该分类下还有 %d 个游戏，无法删除", gameCount)
	}

	// 检查是否有关联的服务项目（通过分类名称匹配）
	itemCount, err := s.categoryRepo.CountServiceItems(ctx, id)
	if err != nil {
		return fmt.Errorf("检查服务项目关联失败: %w", err)
	}
	if itemCount > 0 {
		return fmt.Errorf("该分类下还有 %d 个服务项目，无法删除", itemCount)
	}

	// 删除分类
	if err := s.categoryRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("删除分类失败: %w", err)
	}

	return nil
}

// ============================================================================
// 批量操作
// ============================================================================

// BatchUpdateStatus 批量更新分类启用状态
// 业务规则：
//   - 返回统一格式的批量操作结果
//   - 部分失败不影响其他操作
func (s *GameCategoryService) BatchUpdateStatus(ctx context.Context, categoryIDs []uint64, isActive bool) (*BatchOperationResult, error) {
	if len(categoryIDs) == 0 {
		return nil, errors.New("分类ID列表不能为空")
	}

	result := &BatchOperationResult{
		FailedIDs: make([]uint64, 0),
		Errors:    make([]string, 0),
	}

	// 执行批量更新
	err := s.categoryRepo.BatchUpdateStatus(ctx, categoryIDs, isActive)
	if err != nil {
		return nil, fmt.Errorf("批量更新状态失败: %w", err)
	}

	// 成功：所有ID都更新了
	result.SuccessCount = len(categoryIDs)
	result.FailedCount = 0

	return result, nil
}

// BatchDeleteCategories 批量删除分类
// 业务规则：
//   - 检查每个分类是否有游戏或服务项目关联
//   - 返回统一格式的批量操作结果
//   - 部分失败不影响其他操作
func (s *GameCategoryService) BatchDeleteCategories(ctx context.Context, categoryIDs []uint64) (*BatchOperationResult, error) {
	if len(categoryIDs) == 0 {
		return nil, errors.New("分类ID列表不能为空")
	}

	result := &BatchOperationResult{
		FailedIDs: make([]uint64, 0),
		Errors:    make([]string, 0),
	}

	// 逐个检查并删除
	for _, id := range categoryIDs {
		err := s.DeleteCategory(ctx, id)
		if err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, id)
			result.Errors = append(result.Errors, fmt.Sprintf("分类%d: %s", id, err.Error()))
		} else {
			result.SuccessCount++
		}
	}

	return result, nil
}

// ============================================================================
// 辅助方法
// ============================================================================

// ToCategoryResponse 转换为响应格式
func (s *GameCategoryService) ToCategoryResponse(ctx context.Context, category *model.GameCategory) (*CategoryResponse, error) {
	// 获取游戏数量
	gameCount, _ := s.categoryRepo.CountGames(ctx, category.ID)

	return &CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		IconURL:     category.IconURL,
		SortOrder:   category.SortOrder,
		IsActive:    category.IsActive,
		GameCount:   gameCount,
	}, nil
}

// ToCategoryResponseList 批量转换为响应格式
func (s *GameCategoryService) ToCategoryResponseList(ctx context.Context, categories []*model.GameCategory) ([]CategoryResponse, error) {
	responses := make([]CategoryResponse, 0, len(categories))

	for _, category := range categories {
		resp, err := s.ToCategoryResponse(ctx, category)
		if err != nil {
			continue
		}
		responses = append(responses, *resp)
	}

	return responses, nil
}
