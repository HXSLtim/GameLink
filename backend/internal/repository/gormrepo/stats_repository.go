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

// compile-time assertion
var _ repository.StatsRepository = (*StatsRepository)(nil)
