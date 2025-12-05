package contentcategory

import (
	"context"
	"errors"
	"strings"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

var (
	// ErrNotFound 分类不存在
	ErrNotFound = repository.ErrNotFound
	// ErrValidation 输入校验失败
	ErrValidation = errors.New("validation failed")
	// ErrDuplicate 分类名称已存在
	ErrDuplicate = errors.New("category name already exists")
	// ErrHasFeeds 分类下有动态，无法删除
	ErrHasFeeds = errors.New("category has feeds, cannot delete without migration")
)

// ContentCategoryService 内容分类服务
type ContentCategoryService struct {
	repo repository.ContentCategoryRepository
}

// NewContentCategoryService 创建内容分类服务
func NewContentCategoryService(repo repository.ContentCategoryRepository) *ContentCategoryService {
	return &ContentCategoryService{repo: repo}
}

// CreateRequest 创建分类请求
type CreateRequest struct {
	Name        string `json:"name" binding:"required,max=64"`
	Description string `json:"description"`
	SortOrder   int    `json:"sortOrder"`
	IconURL     string `json:"iconUrl"`
}

// UpdateRequest 更新分类请求
type UpdateRequest struct {
	Name        string                      `json:"name" binding:"required,max=64"`
	Description string                      `json:"description"`
	SortOrder   int                         `json:"sortOrder"`
	Status      model.ContentCategoryStatus `json:"status"`
	IconURL     string                      `json:"iconUrl"`
}

// CategoryDTO 分类DTO
type CategoryDTO struct {
	ID          uint64                      `json:"id"`
	Name        string                      `json:"name"`
	Description string                      `json:"description"`
	SortOrder   int                         `json:"sortOrder"`
	Status      model.ContentCategoryStatus `json:"status"`
	IconURL     string                      `json:"iconUrl"`
	FeedCount   int64                       `json:"feedCount"`
	CreatedAt   string                      `json:"createdAt"`
	UpdatedAt   string                      `json:"updatedAt"`
}

// ListRequest 列出分类请求
type ListRequest struct {
	Page     int                          `form:"page"`
	PageSize int                          `form:"pageSize"`
	Keyword  string                       `form:"keyword"`
	Status   *model.ContentCategoryStatus `form:"status"`
}

// ListResponse 列出分类响应
type ListResponse struct {
	Categories []CategoryDTO `json:"categories"`
	Total      int64         `json:"total"`
}

// Create 创建分类
func (s *ContentCategoryService) Create(ctx context.Context, req CreateRequest) (*CategoryDTO, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, ErrValidation
	}

	category := &model.ContentCategory{
		Name:        strings.TrimSpace(req.Name),
		Description: req.Description,
		SortOrder:   req.SortOrder,
		Status:      model.ContentCategoryStatusActive,
		IconURL:     req.IconURL,
	}

	if err := s.repo.Create(ctx, category); err != nil {
		if strings.Contains(err.Error(), "已存在") {
			return nil, ErrDuplicate
		}
		return nil, err
	}

	return s.toDTO(category, 0), nil
}

// Get 获取分类详情
func (s *ContentCategoryService) Get(ctx context.Context, id uint64) (*CategoryDTO, error) {
	category, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	feedCount, _ := s.repo.GetFeedCount(ctx, id)
	return s.toDTO(category, feedCount), nil
}

// List 列出分类
func (s *ContentCategoryService) List(ctx context.Context, req ListRequest) (*ListResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	categories, total, err := s.repo.List(ctx, repository.ContentCategoryListOptions{
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
		Status:   req.Status,
	})
	if err != nil {
		return nil, err
	}

	dtos := make([]CategoryDTO, 0, len(categories))
	for _, c := range categories {
		feedCount, _ := s.repo.GetFeedCount(ctx, c.ID)
		dtos = append(dtos, *s.toDTO(&c, feedCount))
	}

	return &ListResponse{
		Categories: dtos,
		Total:      total,
	}, nil
}

// Update 更新分类
func (s *ContentCategoryService) Update(ctx context.Context, id uint64, req UpdateRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return ErrValidation
	}

	category, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}

	category.Name = strings.TrimSpace(req.Name)
	category.Description = req.Description
	category.SortOrder = req.SortOrder
	category.IconURL = req.IconURL
	if req.Status.Valid() {
		category.Status = req.Status
	}

	if err := s.repo.Update(ctx, category); err != nil {
		if strings.Contains(err.Error(), "已存在") {
			return ErrDuplicate
		}
		return err
	}

	return nil
}

// Delete 删除分类（可选迁移）
func (s *ContentCategoryService) Delete(ctx context.Context, id uint64, migrateToCategoryID *uint64) error {
	// 检查分类是否存在
	_, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}

	// 检查分类下是否有动态
	feedCount, err := s.repo.GetFeedCount(ctx, id)
	if err != nil {
		return err
	}

	if feedCount > 0 {
		if migrateToCategoryID == nil {
			return ErrHasFeeds
		}
		// 迁移动态到目标分类
		if err := s.repo.MigrateFeeds(ctx, id, *migrateToCategoryID); err != nil {
			return err
		}
	}

	return s.repo.Delete(ctx, id)
}

func (s *ContentCategoryService) toDTO(category *model.ContentCategory, feedCount int64) *CategoryDTO {
	return &CategoryDTO{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		SortOrder:   category.SortOrder,
		Status:      category.Status,
		IconURL:     category.IconURL,
		FeedCount:   feedCount,
		CreatedAt:   category.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   category.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
