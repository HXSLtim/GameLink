package order

import (
	"context"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

type gormReviewRepository struct{ db *gorm.DB }

func NewReviewRepository(db *gorm.DB) repository.ReviewRepository {
	return &gormReviewRepository{db: db}
}

func (r *gormReviewRepository) List(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error) {
	page := repository.NormalizePage(opts.Page)
	size := repository.NormalizePageSize(opts.PageSize)
	offset := (page - 1) * size
	q := r.db.WithContext(ctx).Model(&model.Review{})
	if opts.OrderID != nil {
		q = q.Where("order_id = ?", *opts.OrderID)
	}
	if opts.UserID != nil {
		q = q.Where("user_id = ?", *opts.UserID)
	}
	if opts.PlayerID != nil {
		q = q.Where("player_id = ?", *opts.PlayerID)
	}
	if opts.Status != nil {
		q = q.Where("status = ?", *opts.Status)
	}
	if opts.IsPublic != nil {
		q = q.Where("is_public = ?", *opts.IsPublic)
	}
	if opts.Rating != nil {
		q = q.Where("score = ?", *opts.Rating)
	}
	if opts.DateFrom != nil {
		q = q.Where("created_at >= ?", *opts.DateFrom)
	}
	if opts.DateTo != nil {
		q = q.Where("created_at <= ?", *opts.DateTo)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []model.Review
	if err := q.Order("created_at DESC").Offset(offset).Limit(size).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *gormReviewRepository) Get(ctx context.Context, id uint64) (*model.Review, error) {
	var obj model.Review
	if err := r.db.WithContext(ctx).First(&obj, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &obj, nil
}

func (r *gormReviewRepository) Create(ctx context.Context, obj *model.Review) error {
	return r.db.WithContext(ctx).Create(obj).Error
}

func (r *gormReviewRepository) Update(ctx context.Context, obj *model.Review) error {
	tx := r.db.WithContext(ctx).Model(obj).Where("id = ?", obj.ID).Updates(map[string]any{
		"score":       obj.Score,
		"content":     obj.Content,
		"is_reported": obj.IsReported,
		"status":      obj.Status,
	})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *gormReviewRepository) Delete(ctx context.Context, id uint64) error {
	tx := r.db.WithContext(ctx).Delete(&model.Review{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// ListPending 获取待审核评价列表
func (r *gormReviewRepository) ListPending(ctx context.Context, page, pageSize int) ([]model.Review, int64, error) {
	page = repository.NormalizePage(page)
	size := repository.NormalizePageSize(pageSize)
	offset := (page - 1) * size

	q := r.db.WithContext(ctx).Model(&model.Review{}).Where("status = ?", model.ReviewStatusPending)

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []model.Review
	if err := q.Order("created_at ASC").Offset(offset).Limit(size).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// UpdateStatus 更新评价审核状态
func (r *gormReviewRepository) UpdateStatus(ctx context.Context, id uint64, status model.ReviewStatus, rejectionReason string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if status == model.ReviewStatusRejected && rejectionReason != "" {
		updates["rejection_reason"] = rejectionReason
	}

	tx := r.db.WithContext(ctx).Model(&model.Review{}).Where("id = ?", id).Updates(updates)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// BatchUpdateStatus 批量更新评价审核状态
func (r *gormReviewRepository) BatchUpdateStatus(ctx context.Context, ids []uint64, status model.ReviewStatus, rejectionReason string) error {
	if len(ids) == 0 {
		return nil
	}

	updates := map[string]interface{}{
		"status": status,
	}
	if status == model.ReviewStatusRejected && rejectionReason != "" {
		updates["rejection_reason"] = rejectionReason
	}

	tx := r.db.WithContext(ctx).Model(&model.Review{}).Where("id IN ?", ids).Updates(updates)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// GetStats 获取评价统计数据
func (r *gormReviewRepository) GetStats(ctx context.Context) (repository.ReviewStats, error) {
	var stats repository.ReviewStats

	// 获取总评价数
	if err := r.db.WithContext(ctx).Model(&model.Review{}).
		Where("status = ?", model.ReviewStatusApproved).
		Count(&stats.TotalReviews).Error; err != nil {
		return stats, err
	}

	// 获取平均评分
	var avgScore float64
	if err := r.db.WithContext(ctx).Model(&model.Review{}).
		Where("status = ?", model.ReviewStatusApproved).
		Select("COALESCE(AVG(score), 0)").
		Scan(&avgScore).Error; err != nil {
		return stats, err
	}
	stats.AverageRating = avgScore

	// 获取评分分布
	type scoreCount struct {
		Score int   `json:"score"`
		Count int64 `json:"count"`
	}
	var distribution []scoreCount
	if err := r.db.WithContext(ctx).Model(&model.Review{}).
		Where("status = ?", model.ReviewStatusApproved).
		Select("score, COUNT(*) as count").
		Group("score").
		Scan(&distribution).Error; err != nil {
		return stats, err
	}

	stats.RatingDistribution = make(map[int]int64)
	for _, d := range distribution {
		stats.RatingDistribution[d.Score] = d.Count
	}

	return stats, nil
}

// GetTrend 获取评价趋势（最近N天）
func (r *gormReviewRepository) GetTrend(ctx context.Context, days int) ([]repository.DateValue, error) {
	if days <= 0 {
		days = 30
	}

	since := repository.NowFunc().AddDate(0, 0, -days+1)
	var rows []repository.DateValue

	// 按日期分组统计评价数量
	// PostgreSQL: 使用 created_at::date 进行日期类型转换
	if err := r.db.WithContext(ctx).Model(&model.Review{}).
		Select("created_at::date as date, COUNT(*) as value").
		Where("created_at >= ? AND status = ?", since, model.ReviewStatusApproved).
		Group("created_at::date").
		Order("created_at::date").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}

// GetTopPlayersByReviewCount 获取评价最多的陪玩师
func (r *gormReviewRepository) GetTopPlayersByReviewCount(ctx context.Context, limit int) ([]repository.PlayerReviewStats, error) {
	if limit <= 0 {
		limit = 10
	}

	var stats []repository.PlayerReviewStats
	if err := r.db.WithContext(ctx).Table("reviews").
		Select("reviews.player_id, players.nickname as player_name, COUNT(*) as review_count, COALESCE(AVG(reviews.score), 0) as average_rating").
		Joins("LEFT JOIN players ON reviews.player_id = players.id").
		Where("reviews.status = ?", model.ReviewStatusApproved).
		Group("reviews.player_id, players.nickname").
		Order("review_count DESC, average_rating DESC").
		Limit(limit).
		Scan(&stats).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

// GetTopPlayersByRating 获取评分最高的陪玩师
func (r *gormReviewRepository) GetTopPlayersByRating(ctx context.Context, limit int) ([]repository.PlayerReviewStats, error) {
	if limit <= 0 {
		limit = 10
	}

	var stats []repository.PlayerReviewStats
	if err := r.db.WithContext(ctx).Table("reviews").
		Select("reviews.player_id, players.nickname as player_name, COUNT(*) as review_count, COALESCE(AVG(reviews.score), 0) as average_rating").
		Joins("LEFT JOIN players ON reviews.player_id = players.id").
		Where("reviews.status = ?", model.ReviewStatusApproved).
		Group("reviews.player_id, players.nickname").
		Having("COUNT(*) >= ?", 5). // 至少5条评价才参与排名
		Order("average_rating DESC, review_count DESC").
		Limit(limit).
		Scan(&stats).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

// GetGameStats 获取按游戏统计的评价数据
func (r *gormReviewRepository) GetGameStats(ctx context.Context) ([]repository.GameReviewStats, error) {
	var stats []repository.GameReviewStats
	if err := r.db.WithContext(ctx).Table("reviews").
		Select("orders.game_id, games.name as game_name, COUNT(*) as review_count, COALESCE(AVG(reviews.score), 0) as average_rating").
		Joins("LEFT JOIN orders ON reviews.order_id = orders.id").
		Joins("LEFT JOIN games ON orders.game_id = games.id").
		Where("reviews.status = ?", model.ReviewStatusApproved).
		Group("orders.game_id, games.name").
		Order("review_count DESC").
		Scan(&stats).Error; err != nil {
		return nil, err
	}

	return stats, nil
}
