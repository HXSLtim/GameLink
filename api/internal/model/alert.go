package model

import (
	"context"
	"time"
)

// AlertLevel defines alert severity levels.
type AlertLevel string

const (
	AlertLevelHigh   AlertLevel = "high"
	AlertLevelMedium AlertLevel = "medium"
	AlertLevelLow    AlertLevel = "low"
)

// AlertType defines alert categories.
type AlertType string

const (
	AlertTypeSystem   AlertType = "system"
	AlertTypeBusiness AlertType = "business"
	AlertTypeSecurity AlertType = "security"
)

// Alert represents a system or business alert.
type Alert struct {
	Base
	Level   AlertLevel `json:"level" gorm:"type:varchar(10);not null;index:idx_level_created"`
	Type    AlertType  `json:"type" gorm:"type:varchar(20);not null"`
	Title   string     `json:"title" gorm:"type:varchar(200);not null"`
	Message string     `json:"message" gorm:"type:text"`
	Source  string     `json:"source" gorm:"type:varchar(100)"`
	IsRead  bool       `json:"isRead" gorm:"default:false;index:idx_is_read"`
	ReadBy  *uint64    `json:"readBy" gorm:"index"`
	ReadAt  *time.Time `json:"readAt"`
}

// TableName returns the table name for Alert model.
func (Alert) TableName() string {
	return "alerts"
}

// KPITarget represents a KPI target configuration.
type KPITarget struct {
	Base
	PeriodType  string    `json:"periodType" gorm:"type:varchar(10);not null"` // daily, weekly, monthly
	MetricName  string    `json:"metricName" gorm:"type:varchar(50);not null"` // gmv, orders, users, etc.
	TargetValue float64   `json:"targetValue" gorm:"type:decimal(15,2);not null"`
	StartDate   time.Time `json:"startDate" gorm:"type:date;not null"`
	EndDate     time.Time `json:"endDate" gorm:"type:date;not null"`
	CreatedBy   uint64    `json:"createdBy" gorm:"not null"`
}

// TableName returns the table name for KPITarget model.
func (KPITarget) TableName() string {
	return "kpi_targets"
}

// UserActivityDaily represents daily user activity statistics.
type UserActivityDaily struct {
	Base
	StatDate       time.Time `json:"statDate" gorm:"type:date;not null;uniqueIndex:uk_stat_date"`
	DAU            int       `json:"dau" gorm:"default:0"`            // Daily Active Users
	NewUsers       int       `json:"newUsers" gorm:"default:0"`       // New users registered
	ReturningUsers int       `json:"returningUsers" gorm:"default:0"` // Returning users
	PayingUsers    int       `json:"payingUsers" gorm:"default:0"`    // Users who made payment
}

// TableName returns the table name for UserActivityDaily model.
func (UserActivityDaily) TableName() string {
	return "user_activity_daily"
}

// KPIPeriodType defines KPI period types.
type KPIPeriodType string

const (
	KPIPeriodDaily   KPIPeriodType = "daily"
	KPIPeriodWeekly  KPIPeriodType = "weekly"
	KPIPeriodMonthly KPIPeriodType = "monthly"
)

// KPIMetricName defines standard KPI metric names.
type KPIMetricName string

const (
	KPIMetricGMV        KPIMetricName = "gmv"
	KPIMetricOrders     KPIMetricName = "orders"
	KPIMetricNewUsers   KPIMetricName = "new_users"
	KPIMetricNewPlayers KPIMetricName = "new_players"
	KPIMetricDAU        KPIMetricName = "dau"
	KPIMetricMAU        KPIMetricName = "mau"
	KPIMetricRetention  KPIMetricName = "retention"
	KPIMetricConversion KPIMetricName = "conversion"
	KPIMetricRepurchase KPIMetricName = "repurchase"
)

// AlertQueryOptions defines options for querying alerts.
type AlertQueryOptions struct {
	Page     int
	PageSize int
	Level    string
	Type     string
	IsRead   *bool
	DateFrom *time.Time
	DateTo   *time.Time
}

// AlertRepository defines the interface for alert data access.
type AlertRepository interface {
	// Create creates a new alert.
	Create(ctx context.Context, alert *Alert) error
	// GetByID retrieves an alert by ID.
	GetByID(ctx context.Context, id uint) (*Alert, error)
	// List retrieves alerts with options.
	List(ctx context.Context, opts AlertQueryOptions) ([]Alert, int64, error)
	// MarkAsRead marks an alert as read.
	MarkAsRead(ctx context.Context, id uint) error
	// BatchMarkAsRead marks multiple alerts as read.
	BatchMarkAsRead(ctx context.Context, ids []uint) error
	// GetUnreadCount returns the count of unread alerts.
	GetUnreadCount(ctx context.Context) (int64, error)
}
