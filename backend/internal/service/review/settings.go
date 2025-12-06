package review

import (
	"context"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// SettingsService 评价展示设置服务
type SettingsService struct {
	repo repository.ReviewDisplaySettingsRepository
}

// NewSettingsService 创建评价展示设置服务实例
func NewSettingsService(repo repository.ReviewDisplaySettingsRepository) *SettingsService {
	return &SettingsService{repo: repo}
}

// GetSettings 获取当前评价展示设置
func (s *SettingsService) GetSettings(ctx context.Context) (*model.ReviewDisplaySettings, error) {
	return s.repo.Get(ctx)
}

// UpdateSettings 更新评价展示设置
func (s *SettingsService) UpdateSettings(ctx context.Context, settings *model.ReviewDisplaySettings) (*model.ReviewDisplaySettings, error) {
	// 验证设置
	if err := settings.Validate(); err != nil {
		return nil, err
	}

	// 保存设置
	if err := s.repo.Save(ctx, settings); err != nil {
		return nil, err
	}

	// 返回更新后的设置
	return s.repo.Get(ctx)
}

// UpdateSettingsInput 更新设置的输入参数
type UpdateSettingsInput struct {
	SortBy               *model.ReviewSortBy `json:"sortBy"`
	MinScore             *int                `json:"minScore"`
	ShowAnonymous        *bool               `json:"showAnonymous"`
	PageSize             *int                `json:"pageSize"`
	AutoApprove          *bool               `json:"autoApprove"`
	AutoApproveMinRating *int                `json:"autoApproveMinRating"`
}

// UpdateSettingsPartial 部分更新评价展示设置
func (s *SettingsService) UpdateSettingsPartial(ctx context.Context, input UpdateSettingsInput) (*model.ReviewDisplaySettings, error) {
	// 获取当前设置
	current, err := s.repo.Get(ctx)
	if err != nil {
		return nil, err
	}

	// 应用部分更新
	if input.SortBy != nil {
		current.SortBy = *input.SortBy
	}
	if input.MinScore != nil {
		current.MinScore = *input.MinScore
	}
	if input.ShowAnonymous != nil {
		current.ShowAnonymous = *input.ShowAnonymous
	}
	if input.PageSize != nil {
		current.PageSize = *input.PageSize
	}
	if input.AutoApprove != nil {
		current.AutoApprove = *input.AutoApprove
	}
	if input.AutoApproveMinRating != nil {
		current.AutoApproveMinRating = *input.AutoApproveMinRating
	}

	// 验证并保存
	if err := current.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, current); err != nil {
		return nil, err
	}

	return current, nil
}
