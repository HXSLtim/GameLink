package activity

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// Repository 活动仓库
type Repository struct {
	db *gorm.DB
}

// NewActivityRepository 创建活动仓库
func NewActivityRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ============================================================================
// 活动 CRUD
// ============================================================================

// ActivityListOptions 活动列表查询选项
type ActivityListOptions struct {
	Page      int
	PageSize  int
	Keyword   string
	Type      *model.ActivityType
	Status    *model.ActivityStatus
	IsVisible *bool
	StartTime *time.Time
	EndTime   *time.Time
}

// ListActivities 获取活动列表
func (r *Repository) ListActivities(ctx context.Context, opts ActivityListOptions) ([]model.Activity, int64, error) {
	var items []model.Activity
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Activity{})

	if opts.Keyword != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+opts.Keyword+"%", "%"+opts.Keyword+"%")
	}
	if opts.Type != nil {
		query = query.Where("type = ?", *opts.Type)
	}
	if opts.Status != nil {
		query = query.Where("status = ?", *opts.Status)
	}
	if opts.IsVisible != nil {
		query = query.Where("is_visible = ?", *opts.IsVisible)
	}
	if opts.StartTime != nil {
		query = query.Where("start_at >= ?", *opts.StartTime)
	}
	if opts.EndTime != nil {
		query = query.Where("end_at <= ?", *opts.EndTime)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count activities: %w", err)
	}

	offset := (opts.Page - 1) * opts.PageSize
	if err := query.Preload("Rewards").Order("sort_order ASC, created_at DESC").Offset(offset).Limit(opts.PageSize).Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("list activities: %w", err)
	}

	return items, total, nil
}

// GetActiveActivities 获取进行中的活动列表（用户端）
func (r *Repository) GetActiveActivities(ctx context.Context) ([]model.Activity, error) {
	var items []model.Activity
	now := time.Now()

	if err := r.db.WithContext(ctx).
		Where("status = ? AND is_visible = ? AND start_at <= ? AND end_at >= ?",
			model.ActivityStatusActive, true, now, now).
		Preload("Rewards").
		Order("sort_order ASC").
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("get active activities: %w", err)
	}
	return items, nil
}

// GetVisibleActivities 获取可见的活动列表（包含预热期）
func (r *Repository) GetVisibleActivities(ctx context.Context) ([]model.Activity, error) {
	var items []model.Activity
	now := time.Now()

	if err := r.db.WithContext(ctx).
		Where("is_visible = ? AND (status = ? OR status = ?) AND end_at >= ?",
			true, model.ActivityStatusActive, model.ActivityStatusPreheat, now).
		Preload("Rewards").
		Order("sort_order ASC").
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("get visible activities: %w", err)
	}
	return items, nil
}

// GetActivityByID 根据ID获取活动
func (r *Repository) GetActivityByID(ctx context.Context, id uint64) (*model.Activity, error) {
	var item model.Activity
	if err := r.db.WithContext(ctx).Preload("Rewards").First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get activity: %w", err)
	}
	return &item, nil
}

// CreateActivity 创建活动
func (r *Repository) CreateActivity(ctx context.Context, item *model.Activity) error {
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return fmt.Errorf("create activity: %w", err)
	}
	return nil
}

// UpdateActivity 更新活动
func (r *Repository) UpdateActivity(ctx context.Context, item *model.Activity) error {
	result := r.db.WithContext(ctx).Model(item).Updates(item)
	if result.Error != nil {
		return fmt.Errorf("update activity: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// DeleteActivity 删除活动
func (r *Repository) DeleteActivity(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).Delete(&model.Activity{}, id)
	if result.Error != nil {
		return fmt.Errorf("delete activity: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// UpdateActivityStatus 更新活动状态
func (r *Repository) UpdateActivityStatus(ctx context.Context, id uint64, status model.ActivityStatus) error {
	result := r.db.WithContext(ctx).Model(&model.Activity{}).
		Where("id = ?", id).
		Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("update status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// IncrementParticipants 增加参与人数
func (r *Repository) IncrementParticipants(ctx context.Context, activityID uint64) error {
	result := r.db.WithContext(ctx).Model(&model.Activity{}).
		Where("id = ?", activityID).
		UpdateColumns(map[string]any{
			"total_participants": gorm.Expr("total_participants + 1"),
			"total_claimed":      gorm.Expr("total_claimed + 1"),
		})
	if result.Error != nil {
		return fmt.Errorf("increment participants: %w", result.Error)
	}
	return nil
}

// ResetTodayParticipants 重置今日参与人数（定时任务）
func (r *Repository) ResetTodayParticipants(ctx context.Context) error {
	result := r.db.WithContext(ctx).Model(&model.Activity{}).
		Where("status = ?", model.ActivityStatusActive).
		Update("today_participants", 0)
	if result.Error != nil {
		return fmt.Errorf("reset today participants: %w", result.Error)
	}
	return nil
}

// ============================================================================
// 活动奖励 CRUD
// ============================================================================

// GetRewardByID 根据ID获取奖励
func (r *Repository) GetRewardByID(ctx context.Context, id uint64) (*model.ActivityReward, error) {
	var item model.ActivityReward
	if err := r.db.WithContext(ctx).Preload("CouponTemplate").First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get reward: %w", err)
	}
	return &item, nil
}

// GetRewardsByActivityID 获取活动的奖励列表
func (r *Repository) GetRewardsByActivityID(ctx context.Context, activityID uint64) ([]model.ActivityReward, error) {
	var items []model.ActivityReward
	if err := r.db.WithContext(ctx).
		Where("activity_id = ?", activityID).
		Preload("CouponTemplate").
		Order("sort_order ASC").
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("get rewards: %w", err)
	}
	return items, nil
}

// CreateReward 创建奖励
func (r *Repository) CreateReward(ctx context.Context, item *model.ActivityReward) error {
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return fmt.Errorf("create reward: %w", err)
	}
	return nil
}

// UpdateReward 更新奖励
func (r *Repository) UpdateReward(ctx context.Context, item *model.ActivityReward) error {
	result := r.db.WithContext(ctx).Model(item).Updates(item)
	if result.Error != nil {
		return fmt.Errorf("update reward: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// DeleteReward 删除奖励
func (r *Repository) DeleteReward(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).Delete(&model.ActivityReward{}, id)
	if result.Error != nil {
		return fmt.Errorf("delete reward: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// DecrementRewardStock 减少奖励库存
func (r *Repository) DecrementRewardStock(ctx context.Context, rewardID uint64) error {
	result := r.db.WithContext(ctx).Model(&model.ActivityReward{}).
		Where("id = ? AND (total_stock = 0 OR remaining_stock > 0)", rewardID).
		UpdateColumn("remaining_stock", gorm.Expr("CASE WHEN total_stock = 0 THEN remaining_stock ELSE remaining_stock - 1 END"))
	if result.Error != nil {
		return fmt.Errorf("decrement stock: %w", result.Error)
	}
	return nil
}

// ============================================================================
// 活动参与记录
// ============================================================================

// ParticipationListOptions 参与记录列表查询选项
type ParticipationListOptions struct {
	Page       int
	PageSize   int
	ActivityID *uint64
	UserID     *uint64
	StartTime  *time.Time
	EndTime    *time.Time
}

// ListParticipations 获取参与记录列表
func (r *Repository) ListParticipations(ctx context.Context, opts ParticipationListOptions) ([]model.ActivityParticipation, int64, error) {
	var items []model.ActivityParticipation
	var total int64

	query := r.db.WithContext(ctx).Model(&model.ActivityParticipation{})

	if opts.ActivityID != nil {
		query = query.Where("activity_id = ?", *opts.ActivityID)
	}
	if opts.UserID != nil {
		query = query.Where("user_id = ?", *opts.UserID)
	}
	if opts.StartTime != nil {
		query = query.Where("claimed_at >= ?", *opts.StartTime)
	}
	if opts.EndTime != nil {
		query = query.Where("claimed_at <= ?", *opts.EndTime)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count participations: %w", err)
	}

	offset := (opts.Page - 1) * opts.PageSize
	if err := query.Preload("Activity").Preload("User").Preload("Reward").
		Order("claimed_at DESC").Offset(offset).Limit(opts.PageSize).Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("list participations: %w", err)
	}

	return items, total, nil
}

// CreateParticipation 创建参与记录
func (r *Repository) CreateParticipation(ctx context.Context, item *model.ActivityParticipation) error {
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return fmt.Errorf("create participation: %w", err)
	}
	return nil
}

// CountUserParticipations 统计用户参与某活动次数
func (r *Repository) CountUserParticipations(ctx context.Context, userID, activityID uint64) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.ActivityParticipation{}).
		Where("user_id = ? AND activity_id = ?", userID, activityID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count user participations: %w", err)
	}
	return count, nil
}

// CountTodayParticipations 统计今日参与次数
func (r *Repository) CountTodayParticipations(ctx context.Context, activityID uint64) (int64, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	if err := r.db.WithContext(ctx).Model(&model.ActivityParticipation{}).
		Where("activity_id = ? AND claimed_at >= ?", activityID, today).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count today participations: %w", err)
	}
	return count, nil
}

// GetUserParticipations 获取用户参与记录
func (r *Repository) GetUserParticipations(ctx context.Context, userID uint64, limit int) ([]model.ActivityParticipation, error) {
	var items []model.ActivityParticipation
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("Activity").
		Preload("Reward").
		Order("claimed_at DESC").
		Limit(limit).
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("get user participations: %w", err)
	}
	return items, nil
}

// ============================================================================
// 活动每日统计
// ============================================================================

// GetOrCreateDailyStats 获取或创建每日统计
func (r *Repository) GetOrCreateDailyStats(ctx context.Context, activityID uint64, date time.Time) (*model.ActivityDailyStats, error) {
	statsDate := date.Truncate(24 * time.Hour)
	var stats model.ActivityDailyStats

	err := r.db.WithContext(ctx).
		Where("activity_id = ? AND stats_date = ?", activityID, statsDate).
		First(&stats).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		stats = model.ActivityDailyStats{
			ActivityID: activityID,
			StatsDate:  statsDate,
		}
		if err := r.db.WithContext(ctx).Create(&stats).Error; err != nil {
			return nil, fmt.Errorf("create daily stats: %w", err)
		}
		return &stats, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get daily stats: %w", err)
	}
	return &stats, nil
}

// IncrementDailyStats 增加每日统计
func (r *Repository) IncrementDailyStats(ctx context.Context, activityID uint64) error {
	today := time.Now().Truncate(24 * time.Hour)

	// 先尝试更新
	result := r.db.WithContext(ctx).Model(&model.ActivityDailyStats{}).
		Where("activity_id = ? AND stats_date = ?", activityID, today).
		UpdateColumns(map[string]any{
			"participants": gorm.Expr("participants + 1"),
			"claim_count":  gorm.Expr("claim_count + 1"),
		})

	if result.Error != nil {
		return fmt.Errorf("increment daily stats: %w", result.Error)
	}

	// 如果没有更新到记录，创建新记录
	if result.RowsAffected == 0 {
		stats := model.ActivityDailyStats{
			ActivityID:   activityID,
			StatsDate:    today,
			Participants: 1,
			ClaimCount:   1,
		}
		if err := r.db.WithContext(ctx).Create(&stats).Error; err != nil {
			// 可能是并发创建，忽略唯一索引冲突
			return nil
		}
	}
	return nil
}

// GetActivityStats 获取活动统计
func (r *Repository) GetActivityStats(ctx context.Context, activityID uint64) (map[string]any, error) {
	stats := make(map[string]any)

	// 获取活动基本信息
	activity, err := r.GetActivityByID(ctx, activityID)
	if err != nil {
		return nil, err
	}

	stats["totalParticipants"] = activity.TotalParticipants
	stats["totalClaimed"] = activity.TotalClaimed
	stats["todayParticipants"] = activity.TodayParticipants

	// 获取最近7天统计
	var dailyStats []model.ActivityDailyStats
	sevenDaysAgo := time.Now().AddDate(0, 0, -7).Truncate(24 * time.Hour)
	if err := r.db.WithContext(ctx).
		Where("activity_id = ? AND stats_date >= ?", activityID, sevenDaysAgo).
		Order("stats_date ASC").
		Find(&dailyStats).Error; err != nil {
		return nil, fmt.Errorf("get daily stats: %w", err)
	}
	stats["dailyStats"] = dailyStats

	return stats, nil
}

// GetAllActivityStats 获取所有活动统计概览
func (r *Repository) GetAllActivityStats(ctx context.Context) (map[string]any, error) {
	stats := make(map[string]any)

	// 活动总数
	var totalCount int64
	if err := r.db.WithContext(ctx).Model(&model.Activity{}).Count(&totalCount).Error; err != nil {
		return nil, fmt.Errorf("count activities: %w", err)
	}
	stats["totalCount"] = totalCount

	// 进行中活动数
	var activeCount int64
	if err := r.db.WithContext(ctx).Model(&model.Activity{}).
		Where("status = ?", model.ActivityStatusActive).
		Count(&activeCount).Error; err != nil {
		return nil, fmt.Errorf("count active: %w", err)
	}
	stats["activeCount"] = activeCount

	// 总参与人次
	var totalParticipants int64
	if err := r.db.WithContext(ctx).Model(&model.Activity{}).
		Select("COALESCE(SUM(total_participants), 0)").
		Scan(&totalParticipants).Error; err != nil {
		return nil, fmt.Errorf("sum participants: %w", err)
	}
	stats["totalParticipants"] = totalParticipants

	// 今日参与人次
	var todayParticipants int64
	if err := r.db.WithContext(ctx).Model(&model.Activity{}).
		Select("COALESCE(SUM(today_participants), 0)").
		Scan(&todayParticipants).Error; err != nil {
		return nil, fmt.Errorf("sum today participants: %w", err)
	}
	stats["todayParticipants"] = todayParticipants

	return stats, nil
}
