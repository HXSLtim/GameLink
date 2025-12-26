package referral

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// Repository 推荐仓库
type Repository struct {
	db *gorm.DB
}

// NewReferralRepository 创建推荐仓库
func NewReferralRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ============================================================================
// 推荐配置
// ============================================================================

// GetConfig 获取配置
func (r *Repository) GetConfig(ctx context.Context, key string) (*model.ReferralConfig, error) {
	var config model.ReferralConfig
	if err := r.db.WithContext(ctx).Where("config_key = ?", key).First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get config: %w", err)
	}
	return &config, nil
}

// GetAllConfigs 获取所有配置
func (r *Repository) GetAllConfigs(ctx context.Context) ([]model.ReferralConfig, error) {
	var configs []model.ReferralConfig
	if err := r.db.WithContext(ctx).Order("config_key ASC").Find(&configs).Error; err != nil {
		return nil, fmt.Errorf("get all configs: %w", err)
	}
	return configs, nil
}

// SetConfig 设置配置
func (r *Repository) SetConfig(ctx context.Context, key, value, description string) error {
	config := model.ReferralConfig{
		ConfigKey:   key,
		ConfigValue: value,
		Description: description,
	}
	result := r.db.WithContext(ctx).
		Where("config_key = ?", key).
		Assign(map[string]any{
			"config_value": value,
			"description":  description,
		}).
		FirstOrCreate(&config)
	if result.Error != nil {
		return fmt.Errorf("set config: %w", result.Error)
	}
	return nil
}

// ============================================================================
// 邀请码 CRUD
// ============================================================================

// CodeListOptions 邀请码列表查询选项
type CodeListOptions struct {
	Page     int
	PageSize int
	UserID   *uint64
	Type     *model.ReferralType
	IsActive *bool
	Keyword  string
}

// ListCodes 获取邀请码列表
func (r *Repository) ListCodes(ctx context.Context, opts CodeListOptions) ([]model.ReferralCode, int64, error) {
	var items []model.ReferralCode
	var total int64

	query := r.db.WithContext(ctx).Model(&model.ReferralCode{})

	if opts.UserID != nil {
		query = query.Where("user_id = ?", *opts.UserID)
	}
	if opts.Type != nil {
		query = query.Where("type = ?", *opts.Type)
	}
	if opts.IsActive != nil {
		query = query.Where("is_active = ?", *opts.IsActive)
	}
	if opts.Keyword != "" {
		query = query.Where("code ILIKE ?", "%"+opts.Keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count codes: %w", err)
	}

	offset := (opts.Page - 1) * opts.PageSize
	if err := query.Preload("User").Order("created_at DESC").Offset(offset).Limit(opts.PageSize).Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("list codes: %w", err)
	}

	return items, total, nil
}

// GetCodeByID 根据ID获取邀请码
func (r *Repository) GetCodeByID(ctx context.Context, id uint64) (*model.ReferralCode, error) {
	var item model.ReferralCode
	if err := r.db.WithContext(ctx).Preload("User").First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get code: %w", err)
	}
	return &item, nil
}

// GetCodeByCode 根据邀请码获取
func (r *Repository) GetCodeByCode(ctx context.Context, code string) (*model.ReferralCode, error) {
	var item model.ReferralCode
	if err := r.db.WithContext(ctx).Preload("User").Where("code = ?", code).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get code by code: %w", err)
	}
	return &item, nil
}

// GetUserCode 获取用户的邀请码
func (r *Repository) GetUserCode(ctx context.Context, userID uint64, refType model.ReferralType) (*model.ReferralCode, error) {
	var item model.ReferralCode
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND type = ? AND is_active = ?", userID, refType, true).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get user code: %w", err)
	}
	return &item, nil
}

// CreateCode 创建邀请码
func (r *Repository) CreateCode(ctx context.Context, item *model.ReferralCode) error {
	if item.Code == "" {
		item.Code = generateCode()
	}
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return fmt.Errorf("create code: %w", err)
	}
	return nil
}

// UpdateCode 更新邀请码
func (r *Repository) UpdateCode(ctx context.Context, item *model.ReferralCode) error {
	result := r.db.WithContext(ctx).Model(item).Updates(item)
	if result.Error != nil {
		return fmt.Errorf("update code: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// DeleteCode 删除邀请码
func (r *Repository) DeleteCode(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).Delete(&model.ReferralCode{}, id)
	if result.Error != nil {
		return fmt.Errorf("delete code: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// IncrementCodeUseCount 增加邀请码使用次数
func (r *Repository) IncrementCodeUseCount(ctx context.Context, codeID uint64) error {
	result := r.db.WithContext(ctx).Model(&model.ReferralCode{}).
		Where("id = ?", codeID).
		UpdateColumn("use_count", gorm.Expr("use_count + 1"))
	if result.Error != nil {
		return fmt.Errorf("increment use count: %w", result.Error)
	}
	return nil
}

// generateCode 生成随机邀请码
func generateCode() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// ============================================================================
// 推荐记录 CRUD
// ============================================================================

// ReferralListOptions 推荐记录列表查询选项
type ReferralListOptions struct {
	Page       int
	PageSize   int
	ReferrerID *uint64
	RefereeID  *uint64
	Type       *model.ReferralType
	Status     *model.ReferralStatus
	StartTime  *time.Time
	EndTime    *time.Time
}

// ListReferrals 获取推荐记录列表
func (r *Repository) ListReferrals(ctx context.Context, opts ReferralListOptions) ([]model.Referral, int64, error) {
	var items []model.Referral
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Referral{})

	if opts.ReferrerID != nil {
		query = query.Where("referrer_id = ?", *opts.ReferrerID)
	}
	if opts.RefereeID != nil {
		query = query.Where("referee_id = ?", *opts.RefereeID)
	}
	if opts.Type != nil {
		query = query.Where("type = ?", *opts.Type)
	}
	if opts.Status != nil {
		query = query.Where("status = ?", *opts.Status)
	}
	if opts.StartTime != nil {
		query = query.Where("created_at >= ?", *opts.StartTime)
	}
	if opts.EndTime != nil {
		query = query.Where("created_at <= ?", *opts.EndTime)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count referrals: %w", err)
	}

	offset := (opts.Page - 1) * opts.PageSize
	if err := query.Preload("Referrer").Preload("Referee").Preload("Code").
		Order("created_at DESC").Offset(offset).Limit(opts.PageSize).Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("list referrals: %w", err)
	}

	return items, total, nil
}

// GetReferralByID 根据ID获取推荐记录
func (r *Repository) GetReferralByID(ctx context.Context, id uint64) (*model.Referral, error) {
	var item model.Referral
	if err := r.db.WithContext(ctx).Preload("Referrer").Preload("Referee").Preload("Code").First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get referral: %w", err)
	}
	return &item, nil
}

// GetReferralByReferee 根据被推荐人获取推荐记录
func (r *Repository) GetReferralByReferee(ctx context.Context, refereeID uint64) (*model.Referral, error) {
	var item model.Referral
	if err := r.db.WithContext(ctx).
		Where("referee_id = ?", refereeID).
		Preload("Referrer").
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get referral by referee: %w", err)
	}
	return &item, nil
}

// CreateReferral 创建推荐记录
func (r *Repository) CreateReferral(ctx context.Context, item *model.Referral) error {
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return fmt.Errorf("create referral: %w", err)
	}
	return nil
}

// UpdateReferral 更新推荐记录
func (r *Repository) UpdateReferral(ctx context.Context, item *model.Referral) error {
	result := r.db.WithContext(ctx).Model(item).Updates(item)
	if result.Error != nil {
		return fmt.Errorf("update referral: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// UpdateReferralStatus 更新推荐状态
func (r *Repository) UpdateReferralStatus(ctx context.Context, id uint64, status model.ReferralStatus) error {
	updates := map[string]any{"status": status}
	if status == model.ReferralStatusCompleted {
		now := time.Now()
		updates["completed_at"] = &now
	}
	result := r.db.WithContext(ctx).Model(&model.Referral{}).
		Where("id = ?", id).
		Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("update status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// GetUserReferrals 获取用户的推荐记录
func (r *Repository) GetUserReferrals(ctx context.Context, userID uint64, limit int) ([]model.Referral, error) {
	var items []model.Referral
	if err := r.db.WithContext(ctx).
		Where("referrer_id = ?", userID).
		Preload("Referee").
		Order("created_at DESC").
		Limit(limit).
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("get user referrals: %w", err)
	}
	return items, nil
}

// CountUserReferrals 统计用户推荐数量
func (r *Repository) CountUserReferrals(ctx context.Context, userID uint64, status *model.ReferralStatus) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&model.Referral{}).Where("referrer_id = ?", userID)
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count user referrals: %w", err)
	}
	return count, nil
}

// ============================================================================
// 推荐奖励 CRUD
// ============================================================================

// RewardListOptions 奖励记录列表查询选项
type RewardListOptions struct {
	Page       int
	PageSize   int
	UserID     *uint64
	ReferralID *uint64
	Type       *model.RewardType
	Status     *model.ReferralRewardStatus
}

// ListRewards 获取奖励记录列表
func (r *Repository) ListRewards(ctx context.Context, opts RewardListOptions) ([]model.ReferralReward, int64, error) {
	var items []model.ReferralReward
	var total int64

	query := r.db.WithContext(ctx).Model(&model.ReferralReward{})

	if opts.UserID != nil {
		query = query.Where("user_id = ?", *opts.UserID)
	}
	if opts.ReferralID != nil {
		query = query.Where("referral_id = ?", *opts.ReferralID)
	}
	if opts.Type != nil {
		query = query.Where("type = ?", *opts.Type)
	}
	if opts.Status != nil {
		query = query.Where("status = ?", *opts.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count rewards: %w", err)
	}

	offset := (opts.Page - 1) * opts.PageSize
	if err := query.Preload("Referral").Preload("User").Preload("Coupon").
		Order("created_at DESC").Offset(offset).Limit(opts.PageSize).Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("list rewards: %w", err)
	}

	return items, total, nil
}

// GetRewardByID 根据ID获取奖励记录
func (r *Repository) GetRewardByID(ctx context.Context, id uint64) (*model.ReferralReward, error) {
	var item model.ReferralReward
	if err := r.db.WithContext(ctx).Preload("Referral").Preload("User").Preload("Coupon").First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get reward: %w", err)
	}
	return &item, nil
}

// CreateReward 创建奖励记录
func (r *Repository) CreateReward(ctx context.Context, item *model.ReferralReward) error {
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return fmt.Errorf("create reward: %w", err)
	}
	return nil
}

// UpdateRewardStatus 更新奖励状态
func (r *Repository) UpdateRewardStatus(ctx context.Context, id uint64, status model.ReferralRewardStatus, failureReason string) error {
	updates := map[string]any{"status": status}
	if status == model.ReferralRewardStatusIssued {
		now := time.Now()
		updates["issued_at"] = &now
	}
	if failureReason != "" {
		updates["failure_reason"] = failureReason
	}
	result := r.db.WithContext(ctx).Model(&model.ReferralReward{}).
		Where("id = ?", id).
		Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("update reward status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// GetUserRewards 获取用户的奖励记录
func (r *Repository) GetUserRewards(ctx context.Context, userID uint64, limit int) ([]model.ReferralReward, error) {
	var items []model.ReferralReward
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("Referral").
		Order("created_at DESC").
		Limit(limit).
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("get user rewards: %w", err)
	}
	return items, nil
}

// SumUserRewards 统计用户奖励总额
func (r *Repository) SumUserRewards(ctx context.Context, userID uint64) (int64, error) {
	var sum int64
	if err := r.db.WithContext(ctx).Model(&model.ReferralReward{}).
		Where("user_id = ? AND status = ?", userID, model.ReferralRewardStatusIssued).
		Select("COALESCE(SUM(amount_cents), 0)").
		Scan(&sum).Error; err != nil {
		return 0, fmt.Errorf("sum user rewards: %w", err)
	}
	return sum, nil
}

// ============================================================================
// 统计
// ============================================================================

// GetReferralStats 获取推荐统计
func (r *Repository) GetReferralStats(ctx context.Context) (map[string]any, error) {
	stats := make(map[string]any)

	// 总推荐数
	var totalCount int64
	if err := r.db.WithContext(ctx).Model(&model.Referral{}).Count(&totalCount).Error; err != nil {
		return nil, fmt.Errorf("count referrals: %w", err)
	}
	stats["totalCount"] = totalCount

	// 已完成推荐数
	var completedCount int64
	if err := r.db.WithContext(ctx).Model(&model.Referral{}).
		Where("status IN ?", []model.ReferralStatus{model.ReferralStatusCompleted, model.ReferralStatusRewarded}).
		Count(&completedCount).Error; err != nil {
		return nil, fmt.Errorf("count completed: %w", err)
	}
	stats["completedCount"] = completedCount

	// 已发放奖励数
	var rewardedCount int64
	if err := r.db.WithContext(ctx).Model(&model.ReferralReward{}).
		Where("status = ?", model.ReferralRewardStatusIssued).
		Count(&rewardedCount).Error; err != nil {
		return nil, fmt.Errorf("count rewarded: %w", err)
	}
	stats["rewardedCount"] = rewardedCount

	// 总奖励金额
	var totalRewardCents int64
	if err := r.db.WithContext(ctx).Model(&model.ReferralReward{}).
		Where("status = ?", model.ReferralRewardStatusIssued).
		Select("COALESCE(SUM(amount_cents), 0)").
		Scan(&totalRewardCents).Error; err != nil {
		return nil, fmt.Errorf("sum rewards: %w", err)
	}
	stats["totalRewardCents"] = totalRewardCents

	// 活跃邀请码数
	var activeCodeCount int64
	if err := r.db.WithContext(ctx).Model(&model.ReferralCode{}).
		Where("is_active = ?", true).
		Count(&activeCodeCount).Error; err != nil {
		return nil, fmt.Errorf("count active codes: %w", err)
	}
	stats["activeCodeCount"] = activeCodeCount

	return stats, nil
}

// GetUserReferralStats 获取用户推荐统计
func (r *Repository) GetUserReferralStats(ctx context.Context, userID uint64) (map[string]any, error) {
	stats := make(map[string]any)

	// 总推荐数
	totalCount, err := r.CountUserReferrals(ctx, userID, nil)
	if err != nil {
		return nil, err
	}
	stats["totalCount"] = totalCount

	// 已完成推荐数
	completedStatus := model.ReferralStatusCompleted
	completedCount, err := r.CountUserReferrals(ctx, userID, &completedStatus)
	if err != nil {
		return nil, err
	}
	rewardedStatus := model.ReferralStatusRewarded
	rewardedCount, err := r.CountUserReferrals(ctx, userID, &rewardedStatus)
	if err != nil {
		return nil, err
	}
	stats["completedCount"] = completedCount + rewardedCount

	// 总奖励金额
	totalRewardCents, err := r.SumUserRewards(ctx, userID)
	if err != nil {
		return nil, err
	}
	stats["totalRewardCents"] = totalRewardCents

	return stats, nil
}
