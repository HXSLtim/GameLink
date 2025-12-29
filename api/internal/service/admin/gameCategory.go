package admin

import (
	"context"
	"fmt"
	"strings"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/common"
	"gamelink/pkg/apierr"
)

// Input types for GameCategory operations

// CreateGameCategoryInput 创建游戏分类输入
type CreateGameCategoryInput struct {
	Name        string
	Description string
	IconURL     string
	SortOrder   int
}

// UpdateGameCategoryInput 更新游戏分类输入
type UpdateGameCategoryInput struct {
	Name        *string
	Description *string
	IconURL     *string
	SortOrder   *int
	IsActive    *bool
}

// CreateGameCategory 创建游戏分类
func (s *AdminService) CreateGameCategory(ctx context.Context, input CreateGameCategoryInput) (*model.GameCategory, error) {
	// Validate input
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrValidation.WithDetails("name is required")
	}
	if len(name) > 50 {
		return nil, ErrValidation.WithDetails("name must be at most 50 characters")
	}

	if len(input.Description) > 500 {
		return nil, ErrValidation.WithDetails("description must be at most 500 characters")
	}

	if len(input.IconURL) > 255 {
		return nil, ErrValidation.WithDetails("icon_url must be at most 255 characters")
	}

	if s.tx == nil {
		return nil, apierr.InternalError("transaction manager not configured")
	}

	var createdCategory *model.GameCategory
	err := s.tx.WithTx(ctx, func(r *common.Repos) error {
		category := &model.GameCategory{
			Name:        name,
			Description: strings.TrimSpace(input.Description),
			IconURL:     strings.TrimSpace(input.IconURL),
			SortOrder:   input.SortOrder,
			IsActive:    true,
		}

		if err := r.GameCategories.Create(ctx, category); err != nil {
			return wrapRepositoryError("create game category", err)
		}

		createdCategory = category
		return nil
	})

	if err != nil {
		return nil, err
	}

	s.appendLogAsync(ctx, string(model.OpEntityGameCategory), createdCategory.ID, string(model.OpActionCreate), map[string]any{
		"name": createdCategory.Name,
	})

	return createdCategory, nil
}

// GetGameCategory 获取游戏分类
func (s *AdminService) GetGameCategory(ctx context.Context, id uint64) (*model.GameCategory, error) {
	if s.tx == nil {
		return nil, apierr.InternalError("transaction manager not configured")
	}

	var category *model.GameCategory
	err := s.tx.WithTx(ctx, func(r *common.Repos) error {
		var err error
		category, err = r.GameCategories.Get(ctx, id)
		return err
	})

	if err != nil {
		return nil, wrapRepositoryError("get game category", err)
	}

	return category, nil
}

// ListGameCategoriesPaged 获取游戏分类列表（分页）
func (s *AdminService) ListGameCategoriesPaged(ctx context.Context, page, pageSize int, keyword string, isActive *bool) ([]*model.GameCategory, *model.Pagination, error) {
	if s.tx == nil {
		return nil, nil, apierr.InternalError("transaction manager not configured")
	}

	var categories []*model.GameCategory
	var total int64

	err := s.tx.WithTx(ctx, func(r *common.Repos) error {
		opts := repository.GameCategoryListOptions{
			Page:     page,
			PageSize: pageSize,
			Keyword:  strings.TrimSpace(keyword),
			IsActive: isActive,
		}

		var err error
		categories, total, err = r.GameCategories.List(ctx, opts)
		return err
	})

	if err != nil {
		return nil, nil, wrapRepositoryError("list game categories", err)
	}

	pagination := &model.Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
	}

	return categories, pagination, nil
}

// UpdateGameCategory 更新游戏分类
func (s *AdminService) UpdateGameCategory(ctx context.Context, id uint64, input UpdateGameCategoryInput) (*model.GameCategory, error) {
	if s.tx == nil {
		return nil, apierr.InternalError("transaction manager not configured")
	}

	var updatedCategory *model.GameCategory
	err := s.tx.WithTx(ctx, func(r *common.Repos) error {
		// Get existing category
		category, err := r.GameCategories.Get(ctx, id)
		if err != nil {
			return wrapRepositoryError("get game category", err)
		}

		// Apply updates
		if input.Name != nil {
			name := strings.TrimSpace(*input.Name)
			if name == "" {
				return ErrValidation.WithDetails("name cannot be empty")
			}
			if len(name) > 50 {
				return ErrValidation.WithDetails("name must be at most 50 characters")
			}
			category.Name = name
		}

		if input.Description != nil {
			if len(*input.Description) > 500 {
				return ErrValidation.WithDetails("description must be at most 500 characters")
			}
			category.Description = strings.TrimSpace(*input.Description)
		}

		if input.IconURL != nil {
			if len(*input.IconURL) > 255 {
				return ErrValidation.WithDetails("icon_url must be at most 255 characters")
			}
			category.IconURL = strings.TrimSpace(*input.IconURL)
		}

		if input.SortOrder != nil {
			category.SortOrder = *input.SortOrder
		}

		if input.IsActive != nil {
			category.IsActive = *input.IsActive
		}

		if err := r.GameCategories.Update(ctx, category); err != nil {
			return wrapRepositoryError("update game category", err)
		}

		updatedCategory = category
		return nil
	})

	if err != nil {
		return nil, err
	}

	s.appendLogAsync(ctx, string(model.OpEntityGameCategory), updatedCategory.ID, string(model.OpActionUpdate), map[string]any{
		"name":      updatedCategory.Name,
		"is_active": updatedCategory.IsActive,
	})

	return updatedCategory, nil
}

// DeleteGameCategory 删除游戏分类
func (s *AdminService) DeleteGameCategory(ctx context.Context, id uint64) error {
	if s.tx == nil {
		return apierr.InternalError("transaction manager not configured")
	}

	err := s.tx.WithTx(ctx, func(r *common.Repos) error {
		if err := r.GameCategories.Delete(ctx, id); err != nil {
			return wrapRepositoryError("delete game category", err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	s.appendLogAsync(ctx, string(model.OpEntityGameCategory), id, string(model.OpActionDelete), nil)
	return nil
}

// BatchUpdateGameCategoriesStatus 批量更新游戏分类状态
func (s *AdminService) BatchUpdateGameCategoriesStatus(ctx context.Context, ids []uint64, isActive bool) (*BatchOperationResponse, error) {
	if s.tx == nil {
		return nil, apierr.InternalError("transaction manager not configured")
	}

	response := &BatchOperationResponse{
		TotalCount:   len(ids),
		SuccessItems: make([]uint64, 0),
		FailedItems:  make([]BatchErrorItem, 0),
	}

	if len(ids) == 0 {
		response.FailedCount = len(ids)
		return response, apierr.BadRequest("no category ids provided")
	}

	err := s.tx.WithTx(ctx, func(r *common.Repos) error {
		for _, categoryID := range ids {
			// Verify category exists
			_, err := r.GameCategories.Get(ctx, categoryID)
			if err != nil {
				if isNotFound(err) {
					response.FailedItems = append(response.FailedItems, BatchErrorItem{
						ID:      categoryID,
						Message: "category not found",
					})
					response.FailedCount++
					continue
				}
				response.FailedItems = append(response.FailedItems, BatchErrorItem{
					ID:      categoryID,
					Message: fmt.Sprintf("get category failed: %v", err),
				})
				response.FailedCount++
				continue
			}

			response.SuccessItems = append(response.SuccessItems, categoryID)
			response.SuccessCount++
		}

		if response.SuccessCount > 0 {
			if err := r.GameCategories.BatchUpdateStatus(ctx, response.SuccessItems, isActive); err != nil {
				// Batch update failed, mark all as failed
				for _, id := range response.SuccessItems {
					response.FailedItems = append(response.FailedItems, BatchErrorItem{
						ID:      id,
						Message: fmt.Sprintf("batch update failed: %v", err),
					})
				}
				response.FailedCount += response.SuccessCount
				response.SuccessCount = 0
				response.SuccessItems = make([]uint64, 0)
				return nil
			}
		}

		return nil
	})

	if err != nil {
		return response, err
	}

	for _, id := range response.SuccessItems {
		s.appendLogAsync(ctx, string(model.OpEntityGameCategory), id, string(model.OpActionUpdate), map[string]any{"is_active": isActive})
	}

	return response, nil
}

// BatchDeleteGameCategories 批量删除游戏分类
func (s *AdminService) BatchDeleteGameCategories(ctx context.Context, ids []uint64) (*BatchOperationResponse, error) {
	if s.tx == nil {
		return nil, apierr.InternalError("transaction manager not configured")
	}

	response := &BatchOperationResponse{
		TotalCount:   len(ids),
		SuccessItems: make([]uint64, 0),
		FailedItems:  make([]BatchErrorItem, 0),
	}

	if len(ids) == 0 {
		response.FailedCount = len(ids)
		return response, apierr.BadRequest("no category ids provided")
	}

	err := s.tx.WithTx(ctx, func(r *common.Repos) error {
		for _, categoryID := range ids {
			// Verify category exists
			_, err := r.GameCategories.Get(ctx, categoryID)
			if err != nil {
				if isNotFound(err) {
					response.FailedItems = append(response.FailedItems, BatchErrorItem{
						ID:      categoryID,
						Message: "category not found",
					})
					response.FailedCount++
					continue
				}
				response.FailedItems = append(response.FailedItems, BatchErrorItem{
					ID:      categoryID,
					Message: fmt.Sprintf("get category failed: %v", err),
				})
				response.FailedCount++
				continue
			}

			response.SuccessItems = append(response.SuccessItems, categoryID)
			response.SuccessCount++
		}

		if response.SuccessCount > 0 {
			if err := r.GameCategories.BatchDelete(ctx, response.SuccessItems); err != nil {
				// Batch delete failed, mark all as failed
				for _, id := range response.SuccessItems {
					response.FailedItems = append(response.FailedItems, BatchErrorItem{
						ID:      id,
						Message: fmt.Sprintf("batch delete failed: %v", err),
					})
				}
				response.FailedCount += response.SuccessCount
				response.SuccessCount = 0
				response.SuccessItems = make([]uint64, 0)
				return nil
			}
		}

		return nil
	})

	if err != nil {
		return response, err
	}

	for _, id := range response.SuccessItems {
		s.appendLogAsync(ctx, string(model.OpEntityGameCategory), id, string(model.OpActionDelete), nil)
	}

	return response, nil
}
