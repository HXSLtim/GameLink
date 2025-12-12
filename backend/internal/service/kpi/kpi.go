// Package kpi provides KPI management and calculation services.
package kpi

import (
	"context"
	"time"

	"gorm.io/gorm"

	"gamelink/internal/model"
)

// KPIService provides KPI management functionality.
type KPIService struct {
	db *gorm.DB
}

// NewKPIService creates a new KPI service.
func NewKPIService(db *gorm.DB) *KPIService {
	return &KPIService{db: db}
}

// PeriodType represents KPI period types.
type PeriodType string

const (
	PeriodToday  PeriodType = "today"
	PeriodWeek   PeriodType = "week"
	PeriodMonth  PeriodType = "month"
	PeriodCustom PeriodType = "custom"
)

// CompareType represents comparison types.
type CompareType string

const (
	CompareMoM CompareType = "mom" // Month over Month
	CompareYoY CompareType = "yoy" // Year over Year
)

// KPIMetric represents a single KPI metric.
type KPIMetric struct {
	Name           string  `json:"name"`
	CurrentValue   float64 `json:"currentValue"`
	TargetValue    float64 `json:"targetValue"`
	PreviousValue  float64 `json:"previousValue"`
	CompletionRate float64 `json:"completionRate"`
	ChangeRate     float64 `json:"changeRate"`
	Trend          string  `json:"trend"` // up, down, flat
}

// KPIOverview represents the KPI overview data.
type KPIOverview struct {
	GMV            KPIMetric `json:"gmv"`
	Orders         KPIMetric `json:"orders"`
	NewUsers       KPIMetric `json:"newUsers"`
	NewPlayers     KPIMetric `json:"newPlayers"`
	DAU            KPIMetric `json:"dau"`
	Retention      KPIMetric `json:"retention"`
	RepurchaseRate KPIMetric `json:"repurchaseRate"`
}

// TrendPoint represents a data point in a trend.
type TrendPoint struct {
	Date   string   `json:"date"`
	Value  float64  `json:"value"`
	Target *float64 `json:"target,omitempty"`
}

// QueryParams represents KPI query parameters.
type QueryParams struct {
	Period    PeriodType
	StartDate *time.Time
	EndDate   *time.Time
	Compare   CompareType
}

// GetOverview returns KPI overview data.
func (s *KPIService) GetOverview(ctx context.Context, params QueryParams) (*KPIOverview, error) {
	startDate, endDate := s.getDateRange(params)
	prevStartDate, prevEndDate := s.getPreviousDateRange(startDate, endDate, params.Compare)

	overview := &KPIOverview{
		GMV:            s.calculateGMV(ctx, startDate, endDate, prevStartDate, prevEndDate),
		Orders:         s.calculateOrders(ctx, startDate, endDate, prevStartDate, prevEndDate),
		NewUsers:       s.calculateNewUsers(ctx, startDate, endDate, prevStartDate, prevEndDate),
		NewPlayers:     s.calculateNewPlayers(ctx, startDate, endDate, prevStartDate, prevEndDate),
		DAU:            s.calculateDAU(ctx, startDate, endDate, prevStartDate, prevEndDate),
		Retention:      s.calculateRetention(ctx, startDate, endDate, prevStartDate, prevEndDate),
		RepurchaseRate: s.calculateRepurchaseRate(ctx, startDate, endDate, prevStartDate, prevEndDate),
	}

	return overview, nil
}

// getDateRange calculates start and end dates based on period type.
func (s *KPIService) getDateRange(params QueryParams) (time.Time, time.Time) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	switch params.Period {
	case PeriodToday:
		return today, today.Add(24 * time.Hour)
	case PeriodWeek:
		weekStart := today.AddDate(0, 0, -int(today.Weekday()))
		return weekStart, today.Add(24 * time.Hour)
	case PeriodMonth:
		monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		return monthStart, today.Add(24 * time.Hour)
	case PeriodCustom:
		if params.StartDate != nil && params.EndDate != nil {
			return *params.StartDate, params.EndDate.Add(24 * time.Hour)
		}
	}

	// Default to last 30 days
	return today.AddDate(0, 0, -30), today.Add(24 * time.Hour)
}

// getPreviousDateRange calculates the previous period date range for comparison.
func (s *KPIService) getPreviousDateRange(start, end time.Time, compare CompareType) (time.Time, time.Time) {
	duration := end.Sub(start)

	switch compare {
	case CompareYoY:
		return start.AddDate(-1, 0, 0), end.AddDate(-1, 0, 0)
	default: // MoM
		return start.Add(-duration), start
	}
}

// calculateGMV calculates GMV (Gross Merchandise Volume).
func (s *KPIService) calculateGMV(ctx context.Context, start, end, prevStart, prevEnd time.Time) KPIMetric {
	metric := KPIMetric{Name: "GMV"}

	// Current value (amount_cents / 100 to convert to yuan)
	var currentGMVCents int64
	s.db.WithContext(ctx).Table("payments").
		Where("status = ? AND created_at >= ? AND created_at < ?", "completed", start, end).
		Select("COALESCE(SUM(amount_cents), 0)").
		Scan(&currentGMVCents)
	metric.CurrentValue = float64(currentGMVCents) / 100

	// Previous value
	var previousGMVCents int64
	s.db.WithContext(ctx).Table("payments").
		Where("status = ? AND created_at >= ? AND created_at < ?", "completed", prevStart, prevEnd).
		Select("COALESCE(SUM(amount_cents), 0)").
		Scan(&previousGMVCents)
	metric.PreviousValue = float64(previousGMVCents) / 100

	// Get target from kpi_targets table
	metric.TargetValue = s.getTargetValue(ctx, "gmv", start)

	// Calculate completion rate
	if metric.TargetValue > 0 {
		metric.CompletionRate = metric.CurrentValue / metric.TargetValue * 100
	}

	// Calculate change rate and trend
	s.calculateChangeRateAndTrend(&metric)

	return metric
}

// calculateOrders calculates order count.
func (s *KPIService) calculateOrders(ctx context.Context, start, end, prevStart, prevEnd time.Time) KPIMetric {
	metric := KPIMetric{Name: "订单数"}

	// Current value
	var currentCount int64
	s.db.WithContext(ctx).Table("orders").
		Where("created_at >= ? AND created_at < ?", start, end).
		Count(&currentCount)
	metric.CurrentValue = float64(currentCount)

	// Previous value
	var previousCount int64
	s.db.WithContext(ctx).Table("orders").
		Where("created_at >= ? AND created_at < ?", prevStart, prevEnd).
		Count(&previousCount)
	metric.PreviousValue = float64(previousCount)

	// Get target
	metric.TargetValue = s.getTargetValue(ctx, "orders", start)

	// Calculate completion rate
	if metric.TargetValue > 0 {
		metric.CompletionRate = metric.CurrentValue / metric.TargetValue * 100
	}

	s.calculateChangeRateAndTrend(&metric)
	return metric
}

// calculateNewUsers calculates new user count.
func (s *KPIService) calculateNewUsers(ctx context.Context, start, end, prevStart, prevEnd time.Time) KPIMetric {
	metric := KPIMetric{Name: "新用户"}

	var currentCount int64
	s.db.WithContext(ctx).Table("users").
		Where("created_at >= ? AND created_at < ?", start, end).
		Count(&currentCount)
	metric.CurrentValue = float64(currentCount)

	var previousCount int64
	s.db.WithContext(ctx).Table("users").
		Where("created_at >= ? AND created_at < ?", prevStart, prevEnd).
		Count(&previousCount)
	metric.PreviousValue = float64(previousCount)

	metric.TargetValue = s.getTargetValue(ctx, "new_users", start)
	if metric.TargetValue > 0 {
		metric.CompletionRate = metric.CurrentValue / metric.TargetValue * 100
	}

	s.calculateChangeRateAndTrend(&metric)
	return metric
}

// calculateNewPlayers calculates new player (companion) count.
func (s *KPIService) calculateNewPlayers(ctx context.Context, start, end, prevStart, prevEnd time.Time) KPIMetric {
	metric := KPIMetric{Name: "新陪玩师"}

	var currentCount int64
	s.db.WithContext(ctx).Table("players").
		Where("created_at >= ? AND created_at < ?", start, end).
		Count(&currentCount)
	metric.CurrentValue = float64(currentCount)

	var previousCount int64
	s.db.WithContext(ctx).Table("players").
		Where("created_at >= ? AND created_at < ?", prevStart, prevEnd).
		Count(&previousCount)
	metric.PreviousValue = float64(previousCount)

	metric.TargetValue = s.getTargetValue(ctx, "new_players", start)
	if metric.TargetValue > 0 {
		metric.CompletionRate = metric.CurrentValue / metric.TargetValue * 100
	}

	s.calculateChangeRateAndTrend(&metric)
	return metric
}

// calculateDAU calculates daily active users.
func (s *KPIService) calculateDAU(ctx context.Context, start, end, prevStart, prevEnd time.Time) KPIMetric {
	metric := KPIMetric{Name: "日活用户"}

	var currentDAU int64
	s.db.WithContext(ctx).Table("orders").
		Where("created_at >= ? AND created_at < ?", start, end).
		Distinct("user_id").
		Count(&currentDAU)
	metric.CurrentValue = float64(currentDAU)

	var previousDAU int64
	s.db.WithContext(ctx).Table("orders").
		Where("created_at >= ? AND created_at < ?", prevStart, prevEnd).
		Distinct("user_id").
		Count(&previousDAU)
	metric.PreviousValue = float64(previousDAU)

	metric.TargetValue = s.getTargetValue(ctx, "dau", start)
	if metric.TargetValue > 0 {
		metric.CompletionRate = metric.CurrentValue / metric.TargetValue * 100
	}

	s.calculateChangeRateAndTrend(&metric)
	return metric
}

// calculateRetention calculates retention rate.
func (s *KPIService) calculateRetention(ctx context.Context, start, end, prevStart, prevEnd time.Time) KPIMetric {
	metric := KPIMetric{Name: "留存率"}

	// Calculate 7-day retention for current period
	metric.CurrentValue = s.calculateRetentionRate(ctx, start.AddDate(0, 0, -7), 7)
	metric.PreviousValue = s.calculateRetentionRate(ctx, prevStart.AddDate(0, 0, -7), 7)

	metric.TargetValue = s.getTargetValue(ctx, "retention", start)
	if metric.TargetValue > 0 {
		metric.CompletionRate = metric.CurrentValue / metric.TargetValue * 100
	}

	s.calculateChangeRateAndTrend(&metric)
	return metric
}

// calculateRetentionRate calculates retention rate for given days.
func (s *KPIService) calculateRetentionRate(ctx context.Context, targetDate time.Time, days int) float64 {
	var registeredCount int64
	// PostgreSQL: 使用 ::date 进行日期类型转换
	targetDateOnly := targetDate.Truncate(24 * time.Hour)
	nextDay := targetDateOnly.Add(24 * time.Hour)
	s.db.WithContext(ctx).Table("users").
		Where("created_at >= ? AND created_at < ?", targetDateOnly, nextDay).
		Count(&registeredCount)

	if registeredCount == 0 {
		return 0
	}

	var activeCount int64
	// PostgreSQL: 使用时间范围比较代替 DATE() 函数
	s.db.WithContext(ctx).Raw(`
		SELECT COUNT(DISTINCT o.user_id)
		FROM orders o
		INNER JOIN users u ON o.user_id = u.id
		WHERE u.created_at >= ? AND u.created_at < ?
		AND o.created_at > u.created_at
	`, targetDateOnly, nextDay).Scan(&activeCount)

	return float64(activeCount) / float64(registeredCount) * 100
}

// calculateRepurchaseRate calculates repurchase rate.
func (s *KPIService) calculateRepurchaseRate(ctx context.Context, start, end, prevStart, prevEnd time.Time) KPIMetric {
	metric := KPIMetric{Name: "复购率"}

	// Users with more than 1 completed order
	metric.CurrentValue = s.getRepurchaseRate(ctx, start, end)
	metric.PreviousValue = s.getRepurchaseRate(ctx, prevStart, prevEnd)

	metric.TargetValue = s.getTargetValue(ctx, "repurchase", start)
	if metric.TargetValue > 0 {
		metric.CompletionRate = metric.CurrentValue / metric.TargetValue * 100
	}

	s.calculateChangeRateAndTrend(&metric)
	return metric
}

// getRepurchaseRate calculates repurchase rate for a period.
func (s *KPIService) getRepurchaseRate(ctx context.Context, start, end time.Time) float64 {
	// Total paying users
	var totalPayingUsers int64
	s.db.WithContext(ctx).Table("orders").
		Where("status = ? AND created_at >= ? AND created_at < ?", "completed", start, end).
		Distinct("user_id").
		Count(&totalPayingUsers)

	if totalPayingUsers == 0 {
		return 0
	}

	// Users with more than 1 order
	var repurchaseUsers int64
	s.db.WithContext(ctx).Raw(`
		SELECT COUNT(*) FROM (
			SELECT user_id FROM orders
			WHERE status = 'completed' AND created_at >= ? AND created_at < ?
			GROUP BY user_id
			HAVING COUNT(*) > 1
		) AS repurchase
	`, start, end).Scan(&repurchaseUsers)

	return float64(repurchaseUsers) / float64(totalPayingUsers) * 100
}

// getTargetValue gets target value from kpi_targets table.
func (s *KPIService) getTargetValue(ctx context.Context, metricName string, date time.Time) float64 {
	var target model.KPITarget
	err := s.db.WithContext(ctx).
		Where("metric_name = ? AND start_date <= ? AND end_date >= ?", metricName, date, date).
		Order("created_at DESC").
		First(&target).Error

	if err != nil {
		return 0
	}
	return target.TargetValue
}

// calculateChangeRateAndTrend calculates change rate and trend direction.
func (s *KPIService) calculateChangeRateAndTrend(metric *KPIMetric) {
	if metric.PreviousValue > 0 {
		metric.ChangeRate = (metric.CurrentValue - metric.PreviousValue) / metric.PreviousValue * 100
	} else if metric.CurrentValue > 0 {
		metric.ChangeRate = 100
	}

	if metric.ChangeRate > 1 {
		metric.Trend = "up"
	} else if metric.ChangeRate < -1 {
		metric.Trend = "down"
	} else {
		metric.Trend = "flat"
	}
}

// GetTrend returns trend data for a specific metric.
func (s *KPIService) GetTrend(ctx context.Context, metricName string, params QueryParams) ([]TrendPoint, error) {
	startDate, endDate := s.getDateRange(params)
	var trend []TrendPoint

	current := startDate
	for !current.After(endDate) {
		nextDate := current.AddDate(0, 0, 1)
		dateStr := current.Format("2006-01-02")

		var value float64
		switch metricName {
		case "gmv":
			var gmvCents int64
			s.db.WithContext(ctx).Table("payments").
				Where("status = ? AND created_at >= ? AND created_at < ?", "completed", current, nextDate).
				Select("COALESCE(SUM(amount_cents), 0)").
				Scan(&gmvCents)
			value = float64(gmvCents) / 100
		case "orders":
			var count int64
			s.db.WithContext(ctx).Table("orders").
				Where("created_at >= ? AND created_at < ?", current, nextDate).
				Count(&count)
			value = float64(count)
		case "new_users":
			var count int64
			s.db.WithContext(ctx).Table("users").
				Where("created_at >= ? AND created_at < ?", current, nextDate).
				Count(&count)
			value = float64(count)
		case "dau":
			var count int64
			s.db.WithContext(ctx).Table("orders").
				Where("created_at >= ? AND created_at < ?", current, nextDate).
				Distinct("user_id").
				Count(&count)
			value = float64(count)
		}

		// Get target for this date
		targetValue := s.getTargetValue(ctx, metricName, current)
		point := TrendPoint{
			Date:  dateStr,
			Value: value,
		}
		if targetValue > 0 {
			point.Target = &targetValue
		}

		trend = append(trend, point)
		current = nextDate
	}

	return trend, nil
}

// GetTargets returns KPI target configurations.
func (s *KPIService) GetTargets(ctx context.Context, periodType, metricName string) ([]model.KPITarget, error) {
	var targets []model.KPITarget
	q := s.db.WithContext(ctx).Model(&model.KPITarget{})

	if periodType != "" {
		q = q.Where("period_type = ?", periodType)
	}
	if metricName != "" {
		q = q.Where("metric_name = ?", metricName)
	}

	err := q.Order("created_at DESC").Find(&targets).Error
	return targets, err
}

// CreateTarget creates a new KPI target.
func (s *KPIService) CreateTarget(ctx context.Context, target *model.KPITarget) error {
	return s.db.WithContext(ctx).Create(target).Error
}

// UpdateTarget updates an existing KPI target.
func (s *KPIService) UpdateTarget(ctx context.Context, id uint, updates map[string]interface{}) error {
	return s.db.WithContext(ctx).Model(&model.KPITarget{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteTarget deletes a KPI target.
func (s *KPIService) DeleteTarget(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&model.KPITarget{}, id).Error
}
