package gamerank

import (
	"context"
	"errors"
	"fmt"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

var (
	// ErrNotFound 段位不存在
	ErrNotFound = repository.ErrNotFound
	// ErrValidation 输入校验失败
	ErrValidation = errors.New("validation failed")
	// ErrGameNotFound 游戏不存在
	ErrGameNotFound = errors.New("game not found")
	// ErrDuplicateRank 段位重复
	ErrDuplicateRank = errors.New("duplicate rank for this game and level")
)

// GameRankService 游戏段位配置服务
type GameRankService struct {
	ranks repository.GameRankRepository
	games repository.GameRepository
}

// NewGameRankService 创建游戏段位配置服务
func NewGameRankService(
	ranks repository.GameRankRepository,
	games repository.GameRepository,
) *GameRankService {
	return &GameRankService{
		ranks: ranks,
		games: games,
	}
}

// CreateInput 创建段位输入
type CreateInput struct {
	GameID      uint64
	Name        string
	Level       int
	PriceCents  int64
	IconURL     string
	Color       string
	Description string
	SortOrder   int
	IsActive    bool
}

// UpdateInput 更新段位输入
type UpdateInput struct {
	Name        string
	Level       int
	PriceCents  int64
	IconURL     string
	Color       string
	Description string
	SortOrder   int
	IsActive    bool
}

// Create 创建游戏段位
func (s *GameRankService) Create(ctx context.Context, input CreateInput) (*model.GameRank, error) {
	// 验证游戏是否存在
	if _, err := s.games.Get(ctx, input.GameID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrGameNotFound
		}
		return nil, fmt.Errorf("get game: %w", err)
	}

	// 验证输入
	if input.Name == "" {
		return nil, fmt.Errorf("%w: name is required", ErrValidation)
	}
	if input.PriceCents < 0 {
		return nil, fmt.Errorf("%w: price must be non-negative", ErrValidation)
	}

	rank := &model.GameRank{
		GameID:      input.GameID,
		Name:        input.Name,
		Level:       input.Level,
		PriceCents:  input.PriceCents,
		IconURL:     input.IconURL,
		Color:       input.Color,
		Description: input.Description,
		SortOrder:   input.SortOrder,
		IsActive:    input.IsActive,
	}

	if err := s.ranks.Create(ctx, rank); err != nil {
		return nil, fmt.Errorf("create rank: %w", err)
	}

	return rank, nil
}

// Get 获取游戏段位
func (s *GameRankService) Get(ctx context.Context, id uint64) (*model.GameRank, error) {
	rank, err := s.ranks.GetWithGame(ctx, id)
	if err != nil {
		return nil, err
	}
	return rank, nil
}

// List 获取所有游戏段位
func (s *GameRankService) List(ctx context.Context) ([]model.GameRank, error) {
	return s.ranks.List(ctx)
}

// ListByGameID 根据游戏ID获取段位列表
func (s *GameRankService) ListByGameID(ctx context.Context, gameID uint64) ([]model.GameRank, error) {
	return s.ranks.ListByGameID(ctx, gameID)
}

// ListPaged 分页获取游戏段位
func (s *GameRankService) ListPaged(ctx context.Context, opts repository.GameRankListOptions) ([]model.GameRank, *model.Pagination, error) {
	ranks, total, err := s.ranks.ListPaged(ctx, opts)
	if err != nil {
		return nil, nil, err
	}

	pagination := &model.Pagination{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Total:    int(total),
	}

	return ranks, pagination, nil
}

// Update 更新游戏段位
func (s *GameRankService) Update(ctx context.Context, id uint64, input UpdateInput) (*model.GameRank, error) {
	rank, err := s.ranks.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// 验证输入
	if input.Name == "" {
		return nil, fmt.Errorf("%w: name is required", ErrValidation)
	}
	if input.PriceCents < 0 {
		return nil, fmt.Errorf("%w: price must be non-negative", ErrValidation)
	}

	rank.Name = input.Name
	rank.Level = input.Level
	rank.PriceCents = input.PriceCents
	rank.IconURL = input.IconURL
	rank.Color = input.Color
	rank.Description = input.Description
	rank.SortOrder = input.SortOrder
	rank.IsActive = input.IsActive

	if err := s.ranks.Update(ctx, rank); err != nil {
		return nil, fmt.Errorf("update rank: %w", err)
	}

	return rank, nil
}

// Delete 删除游戏段位
func (s *GameRankService) Delete(ctx context.Context, id uint64) error {
	return s.ranks.Delete(ctx, id)
}

// BatchDelete 批量删除游戏段位
func (s *GameRankService) BatchDelete(ctx context.Context, ids []uint64) (int64, error) {
	return s.ranks.BatchDelete(ctx, ids)
}

// BatchUpdateStatus 批量更新启用状态
func (s *GameRankService) BatchUpdateStatus(ctx context.Context, ids []uint64, isActive bool) (int64, error) {
	return s.ranks.BatchUpdateStatus(ctx, ids, isActive)
}
