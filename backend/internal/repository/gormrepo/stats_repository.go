package gormrepo

import (
    "context"
    "time"

    "gorm.io/gorm"

    "gamelink/internal/repository"
)

type StatsRepository struct{ db *gorm.DB }

func NewStatsRepository(db *gorm.DB) *StatsRepository { return &StatsRepository{db: db} }

func (r *StatsRepository) Dashboard(ctx context.Context) (repository.Dashboard, error) {
    var d repository.Dashboard
    // counts
    if err := r.db.WithContext(ctx).Table("users").Count(&d.TotalUsers).Error; err != nil { return d, err }
    if err := r.db.WithContext(ctx).Table("players").Count(&d.TotalPlayers).Error; err != nil { return d, err }
    if err := r.db.WithContext(ctx).Table("games").Count(&d.TotalGames).Error; err != nil { return d, err }
    if err := r.db.WithContext(ctx).Table("orders").Count(&d.TotalOrders).Error; err != nil { return d, err }

    // orders by status
    d.OrdersByStatus = map[string]int64{}
    type pair struct{ K string; V int64 }
    var rows []pair
    if err := r.db.WithContext(ctx).Table("orders").Select("status as k, COUNT(1) as v").Group("status").Scan(&rows).Error; err != nil { return d, err }
    for _, r2 := range rows { d.OrdersByStatus[r2.K] = r2.V }

    // payments by status + total paid amount
    d.PaymentsByStatus = map[string]int64{}
    rows = nil
    if err := r.db.WithContext(ctx).Table("payments").Select("status as k, COUNT(1) as v").Group("status").Scan(&rows).Error; err != nil { return d, err }
    for _, r2 := range rows { d.PaymentsByStatus[r2.K] = r2.V }
    if err := r.db.WithContext(ctx).Table("payments").Where("status = ?", "paid").Select("COALESCE(SUM(amount_cents),0)").Scan(&d.TotalPaidAmountCents).Error; err != nil { return d, err }
    return d, nil
}

func (r *StatsRepository) RevenueTrend(ctx context.Context, days int) ([]repository.DateValue, error) {
    if days <= 0 { days = 7 }
    since := time.Now().AddDate(0, 0, -days+1)
    var rows []repository.DateValue
    // GROUP BY date(paid_at)
    q := r.db.WithContext(ctx).Table("payments").Select("DATE(paid_at) as date, COALESCE(SUM(amount_cents),0) as value").
        Where("status = ? AND paid_at IS NOT NULL AND paid_at >= ?", "paid", since).
        Group("DATE(paid_at)").Order("DATE(paid_at)")
    if err := q.Scan(&rows).Error; err != nil { return nil, err }
    return rows, nil
}

func (r *StatsRepository) UserGrowth(ctx context.Context, days int) ([]repository.DateValue, error) {
    if days <= 0 { days = 7 }
    since := time.Now().AddDate(0, 0, -days+1)
    var rows []repository.DateValue
    q := r.db.WithContext(ctx).Table("users").Select("DATE(created_at) as date, COUNT(1) as value").
        Where("created_at >= ?", since).Group("DATE(created_at)").Order("DATE(created_at)")
    if err := q.Scan(&rows).Error; err != nil { return nil, err }
    return rows, nil
}

func (r *StatsRepository) OrdersByStatus(ctx context.Context) (map[string]int64, error) {
    type pair struct{ K string; V int64 }
    var rows []pair
    if err := r.db.WithContext(ctx).Table("orders").Select("status as k, COUNT(1) as v").Group("status").Scan(&rows).Error; err != nil {
        return nil, err
    }
    m := map[string]int64{}
    for _, p := range rows { m[p.K] = p.V }
    return m, nil
}

func (r *StatsRepository) TopPlayers(ctx context.Context, limit int) ([]repository.PlayerTop, error) {
    if limit <= 0 { limit = 10 }
    var rows []repository.PlayerTop
    // 简化：按 rating_count 排序
    if err := r.db.WithContext(ctx).Table("players").Select("id as player_id, nickname, rating_average, rating_count").
        Order("rating_count DESC, rating_average DESC").Limit(limit).Scan(&rows).Error; err != nil { return nil, err }
    return rows, nil
}

// AuditOverview returns counts grouped by entity_type and action.
func (r *StatsRepository) AuditOverview(ctx context.Context, from, to *time.Time) (map[string]int64, map[string]int64, error) {
    byEntity := map[string]int64{}
    byAction := map[string]int64{}
    q := r.db.WithContext(ctx).Table("operation_logs")
    if from != nil { q = q.Where("created_at >= ?", *from) }
    if to != nil { q = q.Where("created_at <= ?", *to) }
    type pair struct{ K string; V int64 }
    var rows []pair
    if err := q.Select("entity_type as k, COUNT(1) as v").Group("entity_type").Scan(&rows).Error; err != nil { return nil, nil, err }
    for _, p := range rows { byEntity[p.K] = p.V }
    rows = nil
    if err := q.Select("action as k, COUNT(1) as v").Group("action").Scan(&rows).Error; err != nil { return nil, nil, err }
    for _, p := range rows { byAction[p.K] = p.V }
    return byEntity, byAction, nil
}

// AuditTrend returns per-day counts within range, with optional entity/action filters.
func (r *StatsRepository) AuditTrend(ctx context.Context, from, to *time.Time, entity, action string) ([]repository.DateValue, error) {
    q := r.db.WithContext(ctx).Table("operation_logs")
    if from != nil { q = q.Where("created_at >= ?", *from) }
    if to != nil { q = q.Where("created_at <= ?", *to) }
    if entity != "" { q = q.Where("entity_type = ?", entity) }
    if action != "" { q = q.Where("action = ?", action) }
    var rows []repository.DateValue
    if err := q.Select("DATE(created_at) as date, COUNT(1) as value").Group("DATE(created_at)").Order("DATE(created_at)").Scan(&rows).Error; err != nil { return nil, err }
    return rows, nil
}

// compile-time assertion
var _ repository.StatsRepository = (*StatsRepository)(nil)
