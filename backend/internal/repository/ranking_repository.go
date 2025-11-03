package repository

import (
	"context"

	"gamelink/internal/model"

	"gorm.io/gorm"
)

// RankingRepository 排名仓储接口
type RankingRepository interface {
	// 排名管理
	CreateRanking(ctx context.Context, ranking *model.PlayerRanking) error
	GetPlayerRanking(ctx context.Context, playerID uint64, rankingType model.RankingType, period, periodValue string) (*model.PlayerRanking, error)
	ListRankings(ctx context.Context, opts RankingListOptions) ([]model.PlayerRanking, int64, error)
	UpdateRanking(ctx context.Context, ranking *model.PlayerRanking) error

	// 排名奖励规则
	CreateReward(ctx context.Context, reward *model.RankingReward) error
	GetReward(ctx context.Context, id uint64) (*model.RankingReward, error)
	ListRewards(ctx context.Context, opts RewardListOptions) ([]model.RankingReward, int64, error)
	UpdateReward(ctx context.Context, reward *model.RankingReward) error
	DeleteReward(ctx context.Context, id uint64) error

	// 查询排名奖励
	GetRewardForRank(ctx context.Context, rankingType model.RankingType, period string, rank int) (*model.RankingReward, error)
}

// RankingListOptions 排名查询选项
type RankingListOptions struct {
	PlayerID    *uint64
	RankingType *model.RankingType
	Period      *string
	PeriodValue *string
	Page        int
	PageSize    int
}

// RewardListOptions 奖励规则查询选项
type RewardListOptions struct {
	RankingType *model.RankingType
	Period      *string
	IsActive    *bool
	Page        int
	PageSize    int
}

type rankingRepository struct {
	db *gorm.DB
}

// NewRankingRepository 创建排名仓储
func NewRankingRepository(db *gorm.DB) RankingRepository {
	return &rankingRepository{db: db}
}

// CreateRanking 创建排名记录
func (r *rankingRepository) CreateRanking(ctx context.Context, ranking *model.PlayerRanking) error {
	return r.db.WithContext(ctx).Create(ranking).Error
}

// GetPlayerRanking 获取陪玩师排名
func (r *rankingRepository) GetPlayerRanking(ctx context.Context, playerID uint64, rankingType model.RankingType, period, periodValue string) (*model.PlayerRanking, error) {
	var ranking model.PlayerRanking
	err := r.db.WithContext(ctx).
		Where("player_id = ? AND ranking_type = ? AND period = ? AND period_value = ?",
			playerID, rankingType, period, periodValue).
		First(&ranking).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &ranking, nil
}

// ListRankings 查询排名列表
func (r *rankingRepository) ListRankings(ctx context.Context, opts RankingListOptions) ([]model.PlayerRanking, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.PlayerRanking{})

	if opts.PlayerID != nil {
		query = query.Where("player_id = ?", *opts.PlayerID)
	}
	if opts.RankingType != nil {
		query = query.Where("ranking_type = ?", *opts.RankingType)
	}
	if opts.Period != nil {
		query = query.Where("period = ?", *opts.Period)
	}
	if opts.PeriodValue != nil {
		query = query.Where("period_value = ?", *opts.PeriodValue)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.PageSize < 1 {
		opts.PageSize = 20
	}
	offset := (opts.Page - 1) * opts.PageSize

	var rankings []model.PlayerRanking
	err := query.Order("rank ASC").Offset(offset).Limit(opts.PageSize).Find(&rankings).Error
	if err != nil {
		return nil, 0, err
	}

	return rankings, total, nil
}

// UpdateRanking 更新排名
func (r *rankingRepository) UpdateRanking(ctx context.Context, ranking *model.PlayerRanking) error {
	return r.db.WithContext(ctx).Save(ranking).Error
}

// CreateReward 创建奖励规则
func (r *rankingRepository) CreateReward(ctx context.Context, reward *model.RankingReward) error {
	return r.db.WithContext(ctx).Create(reward).Error
}

// GetReward 获取奖励规则
func (r *rankingRepository) GetReward(ctx context.Context, id uint64) (*model.RankingReward, error) {
	var reward model.RankingReward
	err := r.db.WithContext(ctx).First(&reward, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &reward, nil
}

// ListRewards 查询奖励规则列表
func (r *rankingRepository) ListRewards(ctx context.Context, opts RewardListOptions) ([]model.RankingReward, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.RankingReward{})

	if opts.RankingType != nil {
		query = query.Where("ranking_type = ?", *opts.RankingType)
	}
	if opts.Period != nil {
		query = query.Where("period = ?", *opts.Period)
	}
	if opts.IsActive != nil {
		query = query.Where("is_active = ?", *opts.IsActive)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.PageSize < 1 {
		opts.PageSize = 20
	}
	offset := (opts.Page - 1) * opts.PageSize

	var rewards []model.RankingReward
	err := query.Order("rank_start ASC").Offset(offset).Limit(opts.PageSize).Find(&rewards).Error
	if err != nil {
		return nil, 0, err
	}

	return rewards, total, nil
}

// UpdateReward 更新奖励规则
func (r *rankingRepository) UpdateReward(ctx context.Context, reward *model.RankingReward) error {
	return r.db.WithContext(ctx).Save(reward).Error
}

// DeleteReward 删除奖励规则
func (r *rankingRepository) DeleteReward(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.RankingReward{}, id).Error
}

// GetRewardForRank 获取指定排名对应的奖励
func (r *rankingRepository) GetRewardForRank(ctx context.Context, rankingType model.RankingType, period string, rank int) (*model.RankingReward, error) {
	var reward model.RankingReward
	err := r.db.WithContext(ctx).
		Where("ranking_type = ? AND period = ? AND is_active = ?", rankingType, period, true).
		Where("rank_start <= ? AND rank_end >= ?", rank, rank).
		First(&reward).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &reward, nil
}


