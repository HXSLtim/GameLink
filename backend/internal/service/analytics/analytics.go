// Package analytics provides business analytics services.
package analytics

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// AnalyticsService provides business analytics functionality.
type AnalyticsService struct {
	db *gorm.DB
}

// NewAnalyticsService creates a new analytics service.
func NewAnalyticsService(db *gorm.DB) *AnalyticsService {
	return &AnalyticsService{db: db}
}

// DateRange represents a date range for queries.
type DateRange struct {
	StartDate time.Time
	EndDate   time.Time
}

// Granularity represents time granularity for aggregation.
type Granularity string

const (
	GranularityDay   Granularity = "day"
	GranularityWeek  Granularity = "week"
	GranularityMonth Granularity = "month"
)

// TrendPoint represents a data point in a trend.
type TrendPoint struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
	Type  string  `json:"type,omitempty"`
}

// ActiveUsersData represents active users statistics.
type ActiveUsersData struct {
	DAU         int          `json:"dau"`
	WAU         int          `json:"wau"`
	MAU         int          `json:"mau"`
	DAUMAURatio float64      `json:"dauMauRatio"`
	Trend       []TrendPoint `json:"trend"`
}

// RetentionMatrixRow represents a row in retention matrix.
type RetentionMatrixRow struct {
	CohortDate string    `json:"cohortDate"`
	CohortSize int       `json:"cohortSize"`
	Retention  []float64 `json:"retention"`
}

// RetentionData represents user retention statistics.
type RetentionData struct {
	Day1   float64              `json:"day1"`
	Day7   float64              `json:"day7"`
	Day30  float64              `json:"day30"`
	Matrix []RetentionMatrixRow `json:"matrix"`
}

// PaymentData represents payment analytics.
type PaymentData struct {
	PayingRate    float64      `json:"payingRate"`
	ARPU          float64      `json:"arpu"`
	ARPPU         float64      `json:"arppu"`
	AvgOrderValue float64      `json:"avgOrderValue"`
	Trend         []TrendPoint `json:"trend"`
}

// FunnelStep represents a step in conversion funnel.
type FunnelStep struct {
	Name  string  `json:"name"`
	Value int     `json:"value"`
	Rate  float64 `json:"rate"`
}

// ConversionFunnel represents conversion funnel data.
type ConversionFunnel struct {
	Steps []FunnelStep `json:"steps"`
}

// GetActiveUsers returns active users statistics.
func (s *AnalyticsService) GetActiveUsers(ctx context.Context, dateRange DateRange, granularity Granularity) (*ActiveUsersData, error) {
	result := &ActiveUsersData{}

	// Get DAU (today's active users)
	today := time.Now().Truncate(24 * time.Hour)
	var dauCount int64
	s.db.WithContext(ctx).Table("orders").
		Where("created_at >= ? AND created_at < ?", today, today.Add(24*time.Hour)).
		Distinct("user_id").
		Count(&dauCount)
	result.DAU = int(dauCount)

	// Get WAU (last 7 days active users)
	weekAgo := today.AddDate(0, 0, -7)
	var wauCount int64
	s.db.WithContext(ctx).Table("orders").
		Where("created_at >= ?", weekAgo).
		Distinct("user_id").
		Count(&wauCount)
	result.WAU = int(wauCount)

	// Get MAU (last 30 days active users)
	monthAgo := today.AddDate(0, 0, -30)
	var mauCount int64
	s.db.WithContext(ctx).Table("orders").
		Where("created_at >= ?", monthAgo).
		Distinct("user_id").
		Count(&mauCount)
	result.MAU = int(mauCount)

	// Calculate DAU/MAU ratio
	if result.MAU > 0 {
		result.DAUMAURatio = float64(result.DAU) / float64(result.MAU) * 100
	}

	// Get trend data
	result.Trend = s.getActiveUsersTrend(ctx, dateRange, granularity)

	return result, nil
}

// getActiveUsersTrend returns active users trend data.
func (s *AnalyticsService) getActiveUsersTrend(ctx context.Context, dateRange DateRange, granularity Granularity) []TrendPoint {
	var trend []TrendPoint

	// Generate dates based on granularity
	current := dateRange.StartDate
	for !current.After(dateRange.EndDate) {
		var nextDate time.Time
		var dateFormat string

		switch granularity {
		case GranularityWeek:
			nextDate = current.AddDate(0, 0, 7)
			dateFormat = "2006-01-02"
		case GranularityMonth:
			nextDate = current.AddDate(0, 1, 0)
			dateFormat = "2006-01"
		default:
			nextDate = current.AddDate(0, 0, 1)
			dateFormat = "2006-01-02"
		}

		var count int64
		s.db.WithContext(ctx).Table("orders").
			Where("created_at >= ? AND created_at < ?", current, nextDate).
			Distinct("user_id").
			Count(&count)

		trend = append(trend, TrendPoint{
			Date:  current.Format(dateFormat),
			Value: float64(count),
		})

		current = nextDate
	}

	return trend
}

// GetRetention returns user retention statistics.
func (s *AnalyticsService) GetRetention(ctx context.Context, dateRange DateRange, _ Granularity) (*RetentionData, error) {
	result := &RetentionData{}

	// Calculate retention rates
	result.Day1 = s.calculateRetentionRate(ctx, 1)
	result.Day7 = s.calculateRetentionRate(ctx, 7)
	result.Day30 = s.calculateRetentionRate(ctx, 30)

	// Generate retention matrix (last 7 cohorts)
	result.Matrix = s.generateRetentionMatrix(ctx, 7)

	return result, nil
}

// calculateRetentionRate calculates retention rate for given days.
func (s *AnalyticsService) calculateRetentionRate(ctx context.Context, days int) float64 {
	targetDate := time.Now().AddDate(0, 0, -days).Truncate(24 * time.Hour)

	// Users who registered on target date
	var registeredCount int64
	s.db.WithContext(ctx).Table("users").
		Where("DATE(created_at) = DATE(?)", targetDate).
		Count(&registeredCount)

	if registeredCount == 0 {
		return 0
	}

	// Users who were active after registration
	var activeCount int64
	s.db.WithContext(ctx).Raw(`
		SELECT COUNT(DISTINCT o.user_id)
		FROM orders o
		INNER JOIN users u ON o.user_id = u.id
		WHERE DATE(u.created_at) = DATE(?)
		AND o.created_at > u.created_at
	`, targetDate).Scan(&activeCount)

	return float64(activeCount) / float64(registeredCount) * 100
}

// generateRetentionMatrix generates retention cohort matrix.
func (s *AnalyticsService) generateRetentionMatrix(ctx context.Context, cohortCount int) []RetentionMatrixRow {
	var matrix []RetentionMatrixRow

	for i := cohortCount; i >= 1; i-- {
		cohortDate := time.Now().AddDate(0, 0, -i*7).Truncate(24 * time.Hour)
		row := RetentionMatrixRow{
			CohortDate: cohortDate.Format("2006-01-02"),
		}

		// Get cohort size (users registered in that week)
		var cohortSize int64
		s.db.WithContext(ctx).Table("users").
			Where("created_at >= ? AND created_at < ?", cohortDate, cohortDate.AddDate(0, 0, 7)).
			Count(&cohortSize)
		row.CohortSize = int(cohortSize)

		if cohortSize > 0 {
			// Calculate retention for weeks 1-4
			row.Retention = make([]float64, 4)
			for week := 1; week <= 4 && week <= i; week++ {
				weekStart := cohortDate.AddDate(0, 0, week*7)
				weekEnd := weekStart.AddDate(0, 0, 7)

				var retainedCount int64
				s.db.WithContext(ctx).Raw(`
					SELECT COUNT(DISTINCT o.user_id)
					FROM orders o
					INNER JOIN users u ON o.user_id = u.id
					WHERE u.created_at >= ? AND u.created_at < ?
					AND o.created_at >= ? AND o.created_at < ?
				`, cohortDate, cohortDate.AddDate(0, 0, 7), weekStart, weekEnd).Scan(&retainedCount)

				row.Retention[week-1] = float64(retainedCount) / float64(cohortSize) * 100
			}
		}

		matrix = append(matrix, row)
	}

	return matrix
}

// GetPaymentAnalytics returns payment analytics data.
func (s *AnalyticsService) GetPaymentAnalytics(ctx context.Context, dateRange DateRange, granularity Granularity) (*PaymentData, error) {
	result := &PaymentData{}

	// Get total users in date range
	var totalUsers int64
	s.db.WithContext(ctx).Table("users").
		Where("created_at <= ?", dateRange.EndDate).
		Count(&totalUsers)

	// Get paying users
	var payingUsers int64
	s.db.WithContext(ctx).Table("payments").
		Where("status = ? AND created_at >= ? AND created_at <= ?", "completed", dateRange.StartDate, dateRange.EndDate).
		Distinct("user_id").
		Count(&payingUsers)

	// Calculate paying rate
	if totalUsers > 0 {
		result.PayingRate = float64(payingUsers) / float64(totalUsers) * 100
	}

	// Get total revenue (amount_cents / 100 to convert to yuan)
	var totalRevenueCents int64
	s.db.WithContext(ctx).Table("payments").
		Where("status = ? AND created_at >= ? AND created_at <= ?", "completed", dateRange.StartDate, dateRange.EndDate).
		Select("COALESCE(SUM(amount_cents), 0)").
		Scan(&totalRevenueCents)
	totalRevenue := float64(totalRevenueCents) / 100

	// Calculate ARPU (Average Revenue Per User)
	if totalUsers > 0 {
		result.ARPU = totalRevenue / float64(totalUsers)
	}

	// Calculate ARPPU (Average Revenue Per Paying User)
	if payingUsers > 0 {
		result.ARPPU = totalRevenue / float64(payingUsers)
	}

	// Get total orders
	var totalOrders int64
	s.db.WithContext(ctx).Table("payments").
		Where("status = ? AND created_at >= ? AND created_at <= ?", "completed", dateRange.StartDate, dateRange.EndDate).
		Count(&totalOrders)

	// Calculate average order value
	if totalOrders > 0 {
		result.AvgOrderValue = totalRevenue / float64(totalOrders)
	}

	// Get trend data
	result.Trend = s.getPaymentTrend(ctx, dateRange, granularity)

	return result, nil
}

// getPaymentTrend returns payment trend data.
func (s *AnalyticsService) getPaymentTrend(ctx context.Context, dateRange DateRange, granularity Granularity) []TrendPoint {
	var trend []TrendPoint

	current := dateRange.StartDate
	for !current.After(dateRange.EndDate) {
		var nextDate time.Time
		var dateFormat string

		switch granularity {
		case GranularityWeek:
			nextDate = current.AddDate(0, 0, 7)
			dateFormat = "2006-01-02"
		case GranularityMonth:
			nextDate = current.AddDate(0, 1, 0)
			dateFormat = "2006-01"
		default:
			nextDate = current.AddDate(0, 0, 1)
			dateFormat = "2006-01-02"
		}

		var revenueCents int64
		s.db.WithContext(ctx).Table("payments").
			Where("status = ? AND created_at >= ? AND created_at < ?", "completed", current, nextDate).
			Select("COALESCE(SUM(amount_cents), 0)").
			Scan(&revenueCents)

		trend = append(trend, TrendPoint{
			Date:  current.Format(dateFormat),
			Value: float64(revenueCents) / 100, // Convert cents to yuan
		})

		current = nextDate
	}

	return trend
}

// GetConversionFunnel returns conversion funnel data.
func (s *AnalyticsService) GetConversionFunnel(ctx context.Context, dateRange DateRange, _ Granularity) (*ConversionFunnel, error) {
	result := &ConversionFunnel{
		Steps: make([]FunnelStep, 0),
	}

	// Step 1: Visitors (registered users)
	var registeredUsers int64
	s.db.WithContext(ctx).Table("users").
		Where("created_at >= ? AND created_at <= ?", dateRange.StartDate, dateRange.EndDate).
		Count(&registeredUsers)

	result.Steps = append(result.Steps, FunnelStep{
		Name:  "注册用户",
		Value: int(registeredUsers),
		Rate:  100,
	})

	if registeredUsers == 0 {
		return result, nil
	}

	// Step 2: Browse players (users who viewed player profiles - simulate with order creation)
	var browsingUsers int64
	s.db.WithContext(ctx).Raw(`
		SELECT COUNT(DISTINCT user_id) FROM orders
		WHERE created_at >= ? AND created_at <= ?
		AND user_id IN (SELECT id FROM users WHERE created_at >= ? AND created_at <= ?)
	`, dateRange.StartDate, dateRange.EndDate, dateRange.StartDate, dateRange.EndDate).Scan(&browsingUsers)

	result.Steps = append(result.Steps, FunnelStep{
		Name:  "浏览陪玩师",
		Value: int(browsingUsers),
		Rate:  float64(browsingUsers) / float64(registeredUsers) * 100,
	})

	// Step 3: Create order
	var orderCreators int64
	s.db.WithContext(ctx).Raw(`
		SELECT COUNT(DISTINCT user_id) FROM orders
		WHERE created_at >= ? AND created_at <= ?
		AND user_id IN (SELECT id FROM users WHERE created_at >= ? AND created_at <= ?)
	`, dateRange.StartDate, dateRange.EndDate, dateRange.StartDate, dateRange.EndDate).Scan(&orderCreators)

	result.Steps = append(result.Steps, FunnelStep{
		Name:  "创建订单",
		Value: int(orderCreators),
		Rate:  float64(orderCreators) / float64(registeredUsers) * 100,
	})

	// Step 4: Complete payment
	var paidUsers int64
	s.db.WithContext(ctx).Raw(`
		SELECT COUNT(DISTINCT p.user_id) FROM payments p
		INNER JOIN orders o ON p.order_id = o.id
		WHERE p.status = 'completed'
		AND p.created_at >= ? AND p.created_at <= ?
		AND p.user_id IN (SELECT id FROM users WHERE created_at >= ? AND created_at <= ?)
	`, dateRange.StartDate, dateRange.EndDate, dateRange.StartDate, dateRange.EndDate).Scan(&paidUsers)

	result.Steps = append(result.Steps, FunnelStep{
		Name:  "完成支付",
		Value: int(paidUsers),
		Rate:  float64(paidUsers) / float64(registeredUsers) * 100,
	})

	// Step 5: Complete order
	var completedUsers int64
	s.db.WithContext(ctx).Raw(`
		SELECT COUNT(DISTINCT user_id) FROM orders
		WHERE status = 'completed'
		AND created_at >= ? AND created_at <= ?
		AND user_id IN (SELECT id FROM users WHERE created_at >= ? AND created_at <= ?)
	`, dateRange.StartDate, dateRange.EndDate, dateRange.StartDate, dateRange.EndDate).Scan(&completedUsers)

	result.Steps = append(result.Steps, FunnelStep{
		Name:  "完成订单",
		Value: int(completedUsers),
		Rate:  float64(completedUsers) / float64(registeredUsers) * 100,
	})

	// Step 6: Repurchase (users with more than 1 completed order)
	var repurchaseUsers int64
	s.db.WithContext(ctx).Raw(`
		SELECT COUNT(*) FROM (
			SELECT user_id FROM orders
			WHERE status = 'completed'
			AND created_at >= ? AND created_at <= ?
			AND user_id IN (SELECT id FROM users WHERE created_at >= ? AND created_at <= ?)
			GROUP BY user_id
			HAVING COUNT(*) > 1
		) AS repurchase
	`, dateRange.StartDate, dateRange.EndDate, dateRange.StartDate, dateRange.EndDate).Scan(&repurchaseUsers)

	result.Steps = append(result.Steps, FunnelStep{
		Name:  "复购用户",
		Value: int(repurchaseUsers),
		Rate:  float64(repurchaseUsers) / float64(registeredUsers) * 100,
	})

	return result, nil
}
