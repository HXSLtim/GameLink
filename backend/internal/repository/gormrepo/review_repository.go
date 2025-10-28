package gormrepo

import (
    "context"

    "gorm.io/gorm"

    "gamelink/internal/model"
    "gamelink/internal/repository"
)

type ReviewRepository struct{ db *gorm.DB }

func NewReviewRepository(db *gorm.DB) *ReviewRepository { return &ReviewRepository{db: db} }

func (r *ReviewRepository) List(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error) {
    page := repository.NormalizePage(opts.Page)
    size := repository.NormalizePageSize(opts.PageSize)
    offset := (page - 1) * size
    q := r.db.WithContext(ctx).Model(&model.Review{})
    if opts.OrderID != nil { q = q.Where("order_id = ?", *opts.OrderID) }
    if opts.UserID != nil { q = q.Where("user_id = ?", *opts.UserID) }
    if opts.PlayerID != nil { q = q.Where("player_id = ?", *opts.PlayerID) }
    if opts.DateFrom != nil { q = q.Where("created_at >= ?", *opts.DateFrom) }
    if opts.DateTo != nil { q = q.Where("created_at <= ?", *opts.DateTo) }
    var total int64
    if err := q.Count(&total).Error; err != nil { return nil, 0, err }
    var items []model.Review
    if err := q.Order("created_at DESC").Offset(offset).Limit(size).Find(&items).Error; err != nil { return nil, 0, err }
    return items, total, nil
}

func (r *ReviewRepository) Get(ctx context.Context, id uint64) (*model.Review, error) {
    var obj model.Review
    if err := r.db.WithContext(ctx).First(&obj, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound { return nil, repository.ErrNotFound }
        return nil, err
    }
    return &obj, nil
}

func (r *ReviewRepository) Create(ctx context.Context, obj *model.Review) error {
    return r.db.WithContext(ctx).Create(obj).Error
}

func (r *ReviewRepository) Update(ctx context.Context, obj *model.Review) error {
    tx := r.db.WithContext(ctx).Model(obj).Where("id = ?", obj.ID).Updates(map[string]any{
        "score":   obj.Score,
        "content": obj.Content,
    })
    if tx.Error != nil { return tx.Error }
    if tx.RowsAffected == 0 { return repository.ErrNotFound }
    return nil
}

func (r *ReviewRepository) Delete(ctx context.Context, id uint64) error {
    tx := r.db.WithContext(ctx).Delete(&model.Review{}, id)
    if tx.Error != nil { return tx.Error }
    if tx.RowsAffected == 0 { return repository.ErrNotFound }
    return nil
}

var _ repository.ReviewRepository = (*ReviewRepository)(nil)

