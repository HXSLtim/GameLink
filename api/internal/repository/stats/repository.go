package stats

import (
	"context"
	"time"

	"gorm.io/gorm"

	"gamelink/internal/repository"
)

type gormStatsRepository struct {
	db       *gorm.DB
	isSQLite bool
}

func NewStatsRepository(db *gorm.DB) repository.StatsRepository {
	// 检测数据库类型
	isSQLite := db.Dialector.Name() == "sqlite"
	return &gormStatsRepository{db: db, isSQLite: isSQLite}
}

// dateExpr 返回适合当前数据库的日期转换表达式
func (r *gormStatsRepository) dateExpr(column string) string {
	if r.isSQLite {
		return "date(" + column + ")"
	}
	return column + "::date"
}

func (r *gormStatsRepository) Dashboard(ctx context.Context) (repository.Dashboard, error) {
	var d repository.Dashboard
	// counts
	if err := r.db.WithContext(ctx).Table("users").Count(&d.TotalUsers).Error; err != nil {
		return d, err
	}
	if err := r.db.WithContext(ctx).Table("players").Count(&d.TotalPlayers).Error; err != nil {
		return d, err
	}
	if err := r.db.WithContext(ctx).Table("games").Count(&d.TotalGames).Error; err != nil {
		return d, err
	}
	if err := r.db.WithContext(ctx).Table("orders").Count(&d.TotalOrders).Error; err != nil {
		return d, err
	}

	// orders by status
	d.OrdersByStatus = map[string]int64{}
	type pair struct {
		K string
		V int64
	}
	var rows []pair
	if err := r.db.WithContext(ctx).Table("orders").Select("status as k, COUNT(1) as v").Group("status").Scan(&rows).Error; err != nil {
		return d, err
	}
	for _, r2 := range rows {
		d.OrdersByStatus[r2.K] = r2.V
	}

	// payments by status + total paid amount
	d.PaymentsByStatus = map[string]int64{}
	rows = nil
	if err := r.db.WithContext(ctx).Table("payments").Select("status as k, COUNT(1) as v").Group("status").Scan(&rows).Error; err != nil {
		return d, err
	}
	for _, r2 := range rows {
		d.PaymentsByStatus[r2.K] = r2.V
	}
	if err := r.db.WithContext(ctx).Table("payments").Where("status = ?", "paid").Select("COALESCE(SUM(amount_cents),0)").Scan(&d.TotalPaidAmountCents).Error; err != nil {
		return d, err
	}
	return d, nil
}

func (r *gormStatsRepository) RevenueTrend(ctx context.Context, days int) ([]repository.DateValue, error) {
	if days <= 0 {
		days = 7
	}
	since := time.Now().AddDate(0, 0, -days+1)
	var rows []repository.DateValue
	dateCol := r.dateExpr("paid_at")
	q := r.db.WithContext(ctx).Table("payments").Select(dateCol+" as date, COALESCE(SUM(amount_cents),0) as value").
		Where("status = ? AND paid_at IS NOT NULL AND paid_at >= ?", "paid", since).
		Group(dateCol).Order(dateCol)
	if err := q.Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *gormStatsRepository) UserGrowth(ctx context.Context, days int) ([]repository.DateValue, error) {
	if days <= 0 {
		days = 7
	}
	since := time.Now().AddDate(0, 0, -days+1)
	var rows []repository.DateValue
	dateCol := r.dateExpr("created_at")
	q := r.db.WithContext(ctx).Table("users").Select(dateCol+" as date, COUNT(1) as value").
		Where("created_at >= ?", since).Group(dateCol).Order(dateCol)
	if err := q.Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *gormStatsRepository) UserBehaviorStats(ctx context.Context) (repository.UserBehaviorMetrics, error) {
	var metrics repository.UserBehaviorMetrics

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if err := r.db.WithContext(ctx).
		Table("users").
		Where("last_login_at IS NOT NULL AND last_login_at >= ?", startOfDay).
		Count(&metrics.DAU).Error; err != nil {
		return metrics, err
	}

	avgOnlineQuery := r.db.WithContext(ctx).Table("user_login_histories").
		Where("created_at >= ? AND logout_at IS NOT NULL", startOfDay)
	if r.isSQLite {
		if err := avgOnlineQuery.Select("COALESCE(AVG((julianday(logout_at) - julianday(created_at)) * 86400.0), 0)").Scan(&metrics.AvgOnlineDurationSecond).Error; err != nil {
			return metrics, err
		}
	} else {
		if err := avgOnlineQuery.Select("COALESCE(AVG(EXTRACT(EPOCH FROM (logout_at - created_at))), 0)").Scan(&metrics.AvgOnlineDurationSecond).Error; err != nil {
			return metrics, err
		}
	}

	avgConsumptionQuery := r.db.WithContext(ctx).
		Table("orders").
		Where("status = ?", "completed")
	if r.isSQLite {
		if err := avgConsumptionQuery.Select("COALESCE(SUM(total_price_cents) * 1.0 / NULLIF(COUNT(DISTINCT user_id), 0), 0)").Scan(&metrics.AvgConsumptionCents).Error; err != nil {
			return metrics, err
		}
	} else {
		if err := avgConsumptionQuery.Select("COALESCE(SUM(total_price_cents)::numeric / NULLIF(COUNT(DISTINCT user_id), 0), 0)").Scan(&metrics.AvgConsumptionCents).Error; err != nil {
			return metrics, err
		}
	}

	return metrics, nil
}

func (r *gormStatsRepository) UserActivityTrend(ctx context.Context, days int) ([]repository.DateValue, error) {
	if days <= 0 {
		days = 7
	}

	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -days+1)

	var rows []repository.DateValue
	dateCol := r.dateExpr("last_login_at")
	if err := r.db.WithContext(ctx).
		Table("users").
		Select(dateCol+" as date, COUNT(1) as value").
		Where("last_login_at IS NOT NULL AND last_login_at >= ?", startDate).
		Group(dateCol).
		Order(dateCol).
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	byDate := make(map[string]int64, len(rows))
	for _, row := range rows {
		normalizedDate := row.Date
		if len(normalizedDate) >= 10 {
			normalizedDate = normalizedDate[:10]
		}
		byDate[normalizedDate] = row.Value
	}

	result := make([]repository.DateValue, 0, days)
	for i := 0; i < days; i++ {
		date := startDate.AddDate(0, 0, i).Format("2006-01-02")
		result = append(result, repository.DateValue{
			Date:  date,
			Value: byDate[date],
		})
	}

	return result, nil
}

func (r *gormStatsRepository) UserDistribution(ctx context.Context) (repository.UserDistributionMetrics, error) {
	var metrics repository.UserDistributionMetrics

	var regions []repository.DistributionValue
	if err := r.db.WithContext(ctx).
		Table("user_login_histories").
		Select("COALESCE(NULLIF(TRIM(location), ''), '未知') as name, COUNT(DISTINCT user_id) as value").
		Group("COALESCE(NULLIF(TRIM(location), ''), '未知')").
		Order("value DESC").
		Limit(10).
		Scan(&regions).Error; err != nil {
		return metrics, err
	}

	if len(regions) == 0 {
		var totalUsers int64
		if err := r.db.WithContext(ctx).Table("users").Count(&totalUsers).Error; err != nil {
			return metrics, err
		}
		if totalUsers > 0 {
			regions = append(regions, repository.DistributionValue{
				Name:  "未知",
				Value: totalUsers,
			})
		}
	}
	metrics.ByRegion = regions

	now := time.Now()
	day30 := now.AddDate(0, 0, -30)
	day90 := now.AddDate(0, 0, -90)
	day180 := now.AddDate(0, 0, -180)

	var ageRows []repository.DistributionValue
	if err := r.db.WithContext(ctx).Raw(`
		SELECT age_bucket AS name, COUNT(1) AS value
		FROM (
			SELECT CASE
				WHEN created_at >= ? THEN '0-30天'
				WHEN created_at >= ? THEN '31-90天'
				WHEN created_at >= ? THEN '91-180天'
				ELSE '180天以上'
			END AS age_bucket
			FROM users
		) age_data
		GROUP BY age_bucket
	`, day30, day90, day180).Scan(&ageRows).Error; err != nil {
		return metrics, err
	}

	ageByBucket := make(map[string]int64, len(ageRows))
	for _, row := range ageRows {
		ageByBucket[row.Name] = row.Value
	}

	ageOrder := []string{"0-30天", "31-90天", "91-180天", "180天以上"}
	ages := make([]repository.DistributionValue, 0, len(ageOrder))
	for _, bucket := range ageOrder {
		ages = append(ages, repository.DistributionValue{
			Name:  bucket,
			Value: ageByBucket[bucket],
		})
	}
	metrics.ByAge = ages

	return metrics, nil
}

func (r *gormStatsRepository) OrdersByStatus(ctx context.Context) (map[string]int64, error) {
	type pair struct {
		K string
		V int64
	}
	var rows []pair
	if err := r.db.WithContext(ctx).Table("orders").Select("status as k, COUNT(1) as v").Group("status").Scan(&rows).Error; err != nil {
		return nil, err
	}
	m := map[string]int64{}
	for _, p := range rows {
		m[p.K] = p.V
	}
	return m, nil
}

func (r *gormStatsRepository) TopPlayers(ctx context.Context, limit int) ([]repository.PlayerTop, error) {
	if limit <= 0 {
		limit = 10
	}
	var rows []repository.PlayerTop
	// ç®åï¼æ?rating_count æåº
	if err := r.db.WithContext(ctx).Table("players").Select("id as player_id, nickname, rating_average, rating_count").
		Order("rating_count DESC, rating_average DESC").Limit(limit).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// AuditOverview returns counts grouped by entity_type and action.
func (r *gormStatsRepository) AuditOverview(ctx context.Context, from, to *time.Time) (map[string]int64, map[string]int64, error) {
	byEntity := map[string]int64{}
	byAction := map[string]int64{}
	q := r.db.WithContext(ctx).Table("operation_logs")
	if from != nil {
		q = q.Where("created_at >= ?", *from)
	}
	if to != nil {
		q = q.Where("created_at <= ?", *to)
	}
	type pair struct {
		K string
		V int64
	}
	var rows []pair
	if err := q.Select("entity_type as k, COUNT(1) as v").Group("entity_type").Scan(&rows).Error; err != nil {
		return nil, nil, err
	}
	for _, p := range rows {
		byEntity[p.K] = p.V
	}
	rows = nil
	if err := q.Select("action as k, COUNT(1) as v").Group("action").Scan(&rows).Error; err != nil {
		return nil, nil, err
	}
	for _, p := range rows {
		byAction[p.K] = p.V
	}
	return byEntity, byAction, nil
}

// AuditTrend returns per-day counts within range, with optional entity/action filters.
func (r *gormStatsRepository) AuditTrend(ctx context.Context, from, to *time.Time, entity, action string) ([]repository.DateValue, error) {
	q := r.db.WithContext(ctx).Table("operation_logs")
	if from != nil {
		q = q.Where("created_at >= ?", *from)
	}
	if to != nil {
		q = q.Where("created_at <= ?", *to)
	}
	if entity != "" {
		q = q.Where("entity_type = ?", entity)
	}
	if action != "" {
		q = q.Where("action = ?", action)
	}
	var rows []repository.DateValue
	dateCol := r.dateExpr("created_at")
	if err := q.Select(dateCol + " as date, COUNT(1) as value").Group(dateCol).Order(dateCol).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// compile-time assertion
